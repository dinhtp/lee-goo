# Lee-Goo — Go Modular Monorepo Boilerplate

A reusable Go backend boilerplate with compiled modular architecture, modeled after Magento 2's module system: each module declares its own dependencies, exposes service contracts, provides extension points, and extends the behavior of other modules without tight coupling.

## Tech Stack

| Component | Library | Version |
|-----------|---------|---------|
| Go workspace | `go.work` | Go 1.25 |
| Dependency injection | `go.uber.org/fx` | v1.24 |
| HTTP server | `github.com/labstack/echo/v4` | v4.15 |
| CLI framework | `github.com/spf13/cobra` | v1.10 |
| DB migrations | `github.com/golang-migrate/migrate/v4` | v4.19 |
| Database | PostgreSQL 15+ | — |
| Auth | `github.com/golang-jwt/jwt/v5` (stateless JWT) | v5.3 |
| Config | `github.com/spf13/viper` | v1.21 |

## Quick Start

```bash
# 1. Start PostgreSQL
make dev

# 2. Copy env config
cp .env.example .env
# Edit .env: set DATABASE_PASSWORD, AUTH_JWT_SECRET

# 3. Run the API
make run-api

# 4. Use the module CLI
go run . module list
go run . module make <name>
```

## Module CLI

```bash
go run . module list             # list all modules with status
go run . module make <name>      # scaffold a new module skeleton
go run . module install <name>   # install a module (validate, migrate, register routes)
go run . module enable <name>    # enable a module (cascade-enables deps)
go run . module disable <name>   # disable a module (blocks if dependents are enabled)
go run . module uninstall <name> # uninstall a module (rollback, unregister, remove source)
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make dev` | Start Docker Compose (PostgreSQL) |
| `make run-api` | Run HTTP API server |
| `make run-module` | Run module CLI (no subcommand = help) |
| `make test` | Run all tests |
| `make build` | Build single `lee-goo` binary at repo root |
| `make lint` | Run golangci-lint |
| `make migrate-up` | Run all module migrations |

## Modules

| Module | Go Module | Description | HTTP Routes |
|--------|-----------|-------------|-------------|
| `core` | `modules/core` | Module lifecycle manager | `GET /admin/modules` |
| `user` | `modules/user` | User CRUD + events | `POST /users`, `GET /users/:id`, `PUT /users/:id`, `DELETE /users/:id` |
| `authentication` | `modules/authentication` | Stateless JWT login/refresh | `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout` |
| `authorization` | `modules/authorization` | RBAC roles + permissions | `POST /roles`, `GET /roles`, `GET /roles/:id` |

## Repository Structure

```
/
├── main.go                  — sole entry point; calls cmd.Execute()
├── cmd/
│   ├── cmd.go               — lee-goo root Cobra command; registers api/module/worker
│   ├── api/
│   │   ├── cmd.go           — "api" subcommand
│   │   └── serve.go         — "api serve": composes all modules into HTTP server
│   ├── module/
│   │   ├── cmd.go           — "module" subcommand
│   │   └── *.go             — module lifecycle subcommands (list, install, …)
│   └── worker/
│       ├── cmd.go           — "worker" subcommand
│       └── start.go         — "worker start" (not yet implemented)
├── system/                  — shared infrastructure (no business logic)
│   ├── config/              — Viper-based config loading
│   ├── database/            — pgx stdlib connection (database/sql driver)
│   ├── eventbus/            — async event bus (local in-process)
│   ├── extension/           — extension point registry
│   ├── server/              — Echo server engine (Engine interface + fx lifecycle)
│   ├── logger/              — slog logger (JSON handler)
│   ├── security/            — JWT signer/verifier interfaces
│   └── fx/                  — platform fx.Option composition
├── modules/
│   ├── core/                — module lifecycle management
│   ├── user/                — user CRUD (reference implementation)
│   ├── authentication/      — stateless JWT auth
│   └── authorization/       — RBAC roles + permissions
├── pkg/
│   ├── validate/            — Echo bind + validate helper
│   └── testapp/             — fx test harness for integration tests
├── tests/
│   ├── integration/         — integration tests (require DB)
│   └── cli/                 — CLI smoke tests
└── docs/                    — project documentation
```

## Module Structure

Each module follows the hexagonal architecture pattern:

```
modules/{name}/
├── module.yaml              # manifest: version, dependencies, events
├── go.mod                   # isolated dependency boundary
├── contracts/               # public interfaces (importable by other modules)
│   ├── service.go           # service interface (the public API)
│   ├── event.go             # event payload types
│   └── error.go             # sentinel errors
├── config/                  # module config struct + defaults
├── internal/
│   ├── domain/{entity}/     # domain types + interfaces (no framework deps)
│   ├── service/{entity}/    # business logic + use-case implementations
│   ├── repository/{entity}/ # pgx DB persistence
│   ├── handler/{entity}/    # Echo HTTP handlers
│   └── router/              # route registration
├── migrations/              # per-module SQL migrations (golang-migrate)
└── fx/module.go             # Uber Fx wiring (single composition root)
```

## Running Tests

```bash
# Unit + compile tests (no DB required)
go test ./modules/core/...
go test ./modules/user/...
go test ./modules/authentication/...
go test ./modules/authorization/...

# CLI smoke tests (no DB required)
go test ./tests/cli/... -v -timeout=60s

# Integration tests (require docker compose DB)
make dev
export DATABASE_HOST=localhost DATABASE_USER=leegoo DATABASE_PASSWORD=leegoo DATABASE_DBNAME=leegoo
export AUTH_JWT_SECRET=testsecret
go test ./tests/integration/... -v -timeout=120s

# Full test suite
go test ./... -count=1 -timeout=120s
```

## Adding a New Module

1. Run the scaffold generator: `go run . module make <name>`
2. Implement the domain, service, repository, and handler layers
3. Wire with fx in `fx/module.go`
4. Add the module to `go.work` and compose it in `cmd/api/serve.go`
5. Register migrations in `migrations/`

## Architecture Overview

See [docs/system-architecture.md](docs/system-architecture.md) for the full architecture diagram and detailed design decisions.
