# Codebase Summary

## Directory Structure

```
/
├── cmd/
│   ├── cmd.go                   — Root Cobra command; registers api/module/worker
│   ├── api/
│   │   ├── cmd.go               — "api" subcommand
│   │   └── serve.go             — "api serve": composes all modules into HTTP server
│   ├── module/                  — Cobra CLI for module lifecycle management
│   │   ├── cmd.go
│   │   └── *.go                 — list, make, install, enable, disable, uninstall
│   └── worker/
│       ├── cmd.go               — "worker" subcommand
│       └── start.go             — "worker start" (not yet implemented)
├── system/                      — Shared infrastructure, no business logic
│   ├── config/config.go         — Viper config: ServerConfig, DatabaseConfig, AuthConfig, LogConfig
│   ├── database/                — sqlx Connection adapter + postgresql/ sub-package + fx provider
│   ├── eventbus/                — EventBus interface, NoopEventBus, LocalEventBus + fx provider
│   ├── extension/               — Priority-ordered ExtensionRegistry + fx provider
│   ├── http/                    — Echo v4 server lifecycle + fx provider
│   ├── logger/                  — slog.Logger (JSON handler) + fx provider
│   ├── security/                — Signer/Verifier JWT interfaces (HMAC-SHA256) + fx provider
│   └── fx/options.go            — Options() + TestOptions() composing all system providers
├── modules/
│   ├── core/                    — Module lifecycle management (go.mod: github.com/dinhtp/lee-goo/modules/core)
│   │   ├── contracts/           — ModuleService interface, event/error types
│   │   ├── config/              — ModuleConfig struct
│   │   ├── internal/
│   │   │   ├── domain/module/   — Module entity, UseCase interface, ErrCircularDependency
│   │   │   ├── service/module/  — TopologicalSort, install/enable/disable/sync/doctor logic
│   │   │   ├── repository/module/ — sqlx persistence
│   │   │   └── handler/module/  — Echo handler + router
│   │   ├── migrations/          — SQL migrations for modules table
│   │   ├── tests/               — TopologicalSort unit tests
│   │   └── fx/module.go         — Fx wiring
│   ├── user/                    — User CRUD (go.mod: github.com/dinhtp/lee-goo/modules/user)
│   │   ├── contracts/           — UserService interface, UserCreatedEvent, error types
│   │   ├── config/              — UserConfig struct
│   │   ├── internal/
│   │   │   ├── domain/user/     — User entity, UserPort, UseCase interfaces
│   │   │   ├── service/user/    — Service + UserServiceAdapter
│   │   │   ├── repository/user/ — sqlx persistence (UserPort implementation)
│   │   │   └── handler/user/    — Echo handler + router
│   │   ├── migrations/          — SQL migrations for users table
│   │   └── fx/module.go         — Fx wiring
│   ├── authentication/          — Stateless JWT auth (go.mod: github.com/dinhtp/lee-goo/modules/authentication)
│   │   ├── contracts/           — AuthService interface, TokenPair, event/error types
│   │   ├── config/              — AuthConfig (access/refresh TTLs)
│   │   ├── internal/
│   │   │   ├── domain/auth/     — Auth entity, UseCase interface
│   │   │   ├── service/auth/    — Service + AuthServiceAdapter; bcrypt + JWT signing
│   │   │   └── handler/auth/    — Echo handler (login/refresh/logout) + router
│   │   └── fx/module.go         — Fx wiring
│   └── authorization/           — RBAC roles + permissions (go.mod: github.com/dinhtp/lee-goo/modules/authorization)
│       ├── contracts/           — RoleService/PolicyService interfaces, event types
│       ├── config/              — AuthzConfig (DefaultRole)
│       ├── internal/
│       │   ├── domain/role/     — Role entity, RoleUseCase/PolicyUseCase interfaces
│       │   ├── service/role/    — Service with in-memory permission cache (sync.Map)
│       │   ├── repository/      — role/ and permission/ sqlx repositories
│       │   └── handler/role/    — Echo handler + router
│       ├── migrations/          — roles, permissions, role_permissions tables
│       └── fx/module.go         — Fx wiring (+ user.after_created extension hook)
├── pkg/
│   ├── validate/validate.go     — Echo bind + validator helper
│   └── testapp/                 — Integration test fx harness
│       ├── testapp.go           — App struct, New(), Start()
│       └── options.go           — WithModules(), WithConfig()
├── tests/
│   ├── integration/             — DB integration tests (skip without DATABASE_HOST)
│   │   ├── user_flow_test.go
│   │   ├── auth_flow_test.go
│   │   ├── authz_flow_test.go
│   │   └── core_flow_test.go
│   └── cli/
│       └── module_cli_test.go   — CLI smoke tests (module list exit 0)
└── docs/                        — Project documentation
```

## Module List

| Module | Go Module Path | Status | Key Dependencies |
|--------|---------------|--------|------------------|
| core | `github.com/dinhtp/lee-goo/modules/core` | stable | system/database, system/http, system/eventbus |
| user | `github.com/dinhtp/lee-goo/modules/user` | stable | system/database, system/eventbus, system/extension |
| authentication | `github.com/dinhtp/lee-goo/modules/authentication` | stable | modules/user (contracts), system/security, system/eventbus |
| authorization | `github.com/dinhtp/lee-goo/modules/authorization` | stable | system/database, system/extension, system/eventbus |

## Tech Stack Versions

| Library | Version | Purpose |
|---------|---------|---------|
| `go.uber.org/fx` | v1.24.0 | Dependency injection |
| `github.com/labstack/echo/v4` | v4.15.2 | HTTP server |
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/spf13/viper` | v1.21.0 | Config loading |
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | JWT signing/verification |
| `github.com/jackc/pgx/v5` | v5.9.2 | PostgreSQL driver (stdlib adapter) |
| `github.com/jmoiron/sqlx` | v1.4.0 | SQL extension (named queries, struct scanning) |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | DB migrations (in modules/core/go.mod) |
| `log/slog` | stdlib | Structured logging (JSON handler) |
| `github.com/stretchr/testify` | v1.10.0 | Test assertions |

## Key Interfaces

| Interface | Package | Purpose |
|-----------|---------|---------|
| `contracts.UserService` | `modules/user/contracts` | Public user API for inter-module use |
| `contracts.AuthService` | `modules/authentication/contracts` | Public auth API |
| `domainUser.UseCase` | `modules/user/internal/domain/user` | User business logic boundary |
| `domainUser.UserPort` | `modules/user/internal/domain/user` | Repository boundary |
| `domainRole.RoleUseCase` | `modules/authorization/internal/domain/role` | Role management |
| `domainRole.PolicyUseCase` | `modules/authorization/internal/domain/role` | Policy assignment |
| `eventbus.EventBus` | `system/eventbus` | Async event publishing/subscribing |
| `security.Signer` | `system/security` | JWT signing |
| `security.Verifier` | `system/security` | JWT verification |
| `extension.ExtensionRegistry` | `system/extension` | Extension point hooks |
| `router.HandlerRouter` | per-module `internal/router` | Route registration contract |
