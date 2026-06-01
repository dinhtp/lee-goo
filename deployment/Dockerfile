# syntax=docker/dockerfile:1
# Multi-stage build for lee-goo Go workspace

# ─── Builder stage ────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

# Copy workspace descriptor first (cache layer)
COPY go.work go.work.sum ./

# Copy root module dependency manifests
COPY go.mod go.sum ./

# Copy each workspace module's dependency manifests
COPY modules/authentication/go.mod  modules/authentication/go.sum  ./modules/authentication/
COPY modules/authorization/go.mod   modules/authorization/go.sum   ./modules/authorization/
COPY modules/module_manager/go.mod  modules/module_manager/go.sum  ./modules/module_manager/
COPY modules/user/go.mod            modules/user/go.sum            ./modules/user/

# Pre-download all workspace dependencies (cached unless go files change)
RUN go work sync && go mod download

# Copy full source tree
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -trimpath \
    -ldflags="-w -s" \
    -o lee-goo .

# ─── Final stage ──────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /build/lee-goo ./lee-goo

# Copy SQL migration files needed at runtime by golang-migrate
COPY --from=builder /build/modules/authentication/migrations  ./modules/authentication/migrations
COPY --from=builder /build/modules/authorization/migrations   ./modules/authorization/migrations
COPY --from=builder /build/modules/module_manager/migrations  ./modules/module_manager/migrations
COPY --from=builder /build/modules/user/migrations            ./modules/user/migrations

# Copy module manifests (module.yaml files used by the module manager)
COPY --from=builder /build/modules/authentication/module.yaml  ./modules/authentication/module.yaml
COPY --from=builder /build/modules/authorization/module.yaml   ./modules/authorization/module.yaml
COPY --from=builder /build/modules/module_manager/module.yaml  ./modules/module_manager/module.yaml
COPY --from=builder /build/modules/user/module.yaml            ./modules/user/module.yaml

EXPOSE 8080

ENTRYPOINT ["./lee-goo"]
CMD ["api", "serve"]
