# Deployment Guide

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.25+ | Build and run |
| Docker | 24+ | PostgreSQL via Compose |
| Docker Compose | v2+ | Local DB stack |
| golangci-lint | latest | Linting (optional) |
| make | any | Convenience targets |

---

## Local Development

### 1. Start PostgreSQL

```bash
make dev
# runs: docker compose -f deployment/docker-compose.yml up -d
# postgres:15-alpine, port 5432, db=leegoo, user=leegoo, pass=leegoo
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` — required values:

```env
DATABASE_PASSWORD=leegoo
AUTH_JWT_SECRET=your-secret-here  # min 32 chars recommended
```

Full variable reference:

| Variable | Default | Required |
|----------|---------|----------|
| `ENV` | `local` | No |
| `SERVER_PORT` | `8080` | No |
| `DATABASE_HOST` | `localhost` | Yes |
| `DATABASE_PORT` | `5432` | No |
| `DATABASE_USER` | — | Yes |
| `DATABASE_PASSWORD` | — | Yes |
| `DATABASE_DBNAME` | — | Yes |
| `DATABASE_SSLMODE` | `disable` | No |
| `DATABASE_MAX_OPEN_CONNECTIONS` | `4` | No |
| `DATABASE_MAX_IDLE_CONNECTIONS` | `2` | No |
| `DATABASE_CONNECTION_MAX_TIME` | `0` | No |
| `DATABASE_CONNECTION_IDLE_TIME` | `0` | No |
| `AUTH_JWT_SECRET` | — | Yes (prod) |
| `AUTH_ACCESS_TOKEN_TTL` | `15m` | No |
| `AUTH_REFRESH_TOKEN_TTL` | `168h` | No |
| `SERVER_LOG_LEVEL` | `info` | No |

### 3. Run Migrations

```bash
# Run all module migrations in dependency order
make migrate-up
# equivalent: go run . module migrate --all
```

Migration order (topological): core → user → authorization (authentication has no migrations).

### 4. Start the API Server

```bash
make run-api
# equivalent: go run . api serve
# server listens on :8080
```

### 5. Use the Module CLI

```bash
make run-module
# equivalent: go run . module

go run . module list            # list discovered modules
go run . module doctor          # validate all manifests
go run . module graph           # print dependency graph
go run . module status core     # show module DB status
```

---

## Make Targets Reference

| Target | Command | Description |
|--------|---------|-------------|
| `make dev` | `docker compose up -d` | Start Postgres container |
| `make run-api` | `go run . api serve` | Start HTTP API server |
| `make run-module` | `go run . module` | Module CLI (no subcommand = help) |
| `make test` | `go test ./... -count=1 -timeout=120s` | Full test suite |
| `make lint` | `golangci-lint run ./...` | Static analysis |
| `make build` | `go build -o lee-goo .` | Build single binary |
| `make migrate-up` | `go run . module migrate --all` | Run all pending migrations |

---

## Running Tests

### Unit Tests (no DB required)

```bash
go test ./modules/core/...
go test ./modules/user/...
go test ./modules/authentication/...
go test ./modules/authorization/...
```

### CLI Smoke Tests (no DB required)

```bash
go test ./tests/cli/... -v -timeout=60s
# tests: module list, module doctor, module graph exit 0
```

### Integration Tests (require DB)

```bash
# Ensure DB is running and env vars are set
make dev
export DATABASE_HOST=localhost
export DATABASE_USER=leegoo
export DATABASE_PASSWORD=leegoo
export DATABASE_DBNAME=leegoo
export AUTH_JWT_SECRET=testsecret

go test ./tests/integration/... -v -timeout=120s
```

Tests skip automatically (not fail) when `DATABASE_HOST` or `TEST_DATABASE_URL` is absent.

### Full Suite

```bash
make test
# go test ./... -count=1 -timeout=120s
```

---

## Docker Build

### Build the Image

```bash
docker build -f deployment/Dockerfile -t lee-goo:latest .
```

Multi-stage build:
- **Builder**: `golang:1.25-alpine` — compiles static binary (`CGO_ENABLED=0 -trimpath -ldflags="-w -s"`)
- **Final**: `alpine:3.20` + `ca-certificates tzdata` — binary only (~15MB image)

Binary artifacts copied into final image:
- `/app/lee-goo` — compiled binary
- `/app/modules/*/migrations/` — SQL migration files
- `/app/modules/*/module.yaml` — module manifests

> **Known issue:** Dockerfile currently references `modules/module_manager/` (old name). Update to `modules/core/` before building. See `docs/project-roadmap.md` issue #1.

### Run the Container

```bash
docker run -p 8080:8080 \
  -e DATABASE_HOST=host.docker.internal \
  -e DATABASE_PORT=5432 \
  -e DATABASE_USER=leegoo \
  -e DATABASE_PASSWORD=leegoo \
  -e DATABASE_DBNAME=leegoo \
  -e AUTH_JWT_SECRET=your-secret \
  lee-goo:latest
# default CMD: api serve
```

---

## Database Migration Workflow

### Initial Setup

```bash
# 1. Start postgres
make dev

# 2. Run all migrations (dependency order: core → user → authorization)
make migrate-up
```

### Per-Module Migration

```bash
# Run migrations for a single module
go run . module migrate core
go run . module migrate user
go run . module migrate authorization
```

### Migration Files Location

```
modules/core/migrations/
  000001_create_modules.up.sql
  000001_create_modules.down.sql

modules/user/migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql

modules/authorization/migrations/
  000001_create_roles.up.sql
  000002_create_permissions.up.sql
  000003_create_role_permissions.up.sql
  (+ corresponding .down.sql files)
```

### Adding Migrations to a New Module

1. Create numbered files in `modules/{name}/migrations/`
2. Follow the pattern: `000001_create_{table}.up.sql` / `.down.sql`
3. Run `go run . module migrate <name>`

The migration runner builds its DSN from the standard `DATABASE_*` config fields (no `DATABASE_URL` required).

---

## Adding a New Module

1. Scaffold the skeleton:
   ```bash
   go run . module make <name>
   ```

2. Implement layers in order:
   - `contracts/` — service interface, events, errors
   - `internal/domain/{entity}/` — domain types + interfaces
   - `internal/repository/{entity}/` — pgx persistence
   - `internal/service/{entity}/` — business logic
   - `internal/handler/{entity}/` — Echo handlers
   - `internal/router/` — route registration
   - `fx/module.go` — fx wiring

3. Add migrations to `migrations/`

4. Register in `go.work`:
   ```
   use ./modules/<name>
   ```

5. Compose in `cmd/api/serve.go`:
   ```go
   import nameModule "github.com/dinhtp/lee-goo/modules/<name>/fx"
   // ...
   nameModule.Module(),
   ```

6. Add `module.yaml` manifest with dependencies declared

7. Run migrations: `go run . module migrate <name>`

---

## API Endpoints Reference

| Method | Path | Module | Description |
|--------|------|--------|-------------|
| `POST` | `/users` | user | Create user |
| `GET` | `/users` | user | List users |
| `GET` | `/users/:id` | user | Get user |
| `PUT` | `/users/:id` | user | Update user |
| `DELETE` | `/users/:id` | user | Delete user |
| `POST` | `/auth/login` | authentication | Login → JWT token pair |
| `POST` | `/auth/refresh` | authentication | Refresh access token |
| `POST` | `/auth/logout` | authentication | Logout (no-op, stateless) |
| `POST` | `/roles` | authorization | Create role |
| `GET` | `/roles` | authorization | List roles |
| `GET` | `/healthz` | core | Liveness probe — returns `{"status":"ok"}` |
| `GET` | `/admin/modules` | core | List discovered modules |
| `GET` | `/admin/modules/:name` | core | Get module status |
