# System Architecture

## Layer Diagram

```
┌─────────────────────────────────────────────────┐
│          main.go → cmd/ (lee-goo CLI)            │
│   cmd/api/serve.go      cmd/module/*.go           │
│   (HTTP server)         (module CLI)             │
└──────────────┬──────────────────────────────────┘
               │ composes via fx.Option
┌──────────────▼──────────────────────────────────┐
│               modules/ (feature modules)         │
│  ┌──────────┐ ┌──────┐ ┌───────────┐ ┌───────┐ │
│  │  module  │ │ user │ │   auth    │ │ authz │ │
│  └──────────┘ └──────┘ └───────────┘ └───────┘ │
│  Each module: contracts/ + internal/ + fx/       │
└──────────────┬──────────────────────────────────┘
               │ depends on
┌──────────────▼──────────────────────────────────┐
│              system/ (shared infrastructure)   │
│  config   database  eventbus  extension           │
│  http     logger   security  fx                  │
└─────────────────────────────────────────────────┘
```

## Hexagonal Architecture per Module

Each module enforces the hexagonal (ports & adapters) pattern:

```
contracts/              ← PUBLIC ports (importable by other modules)
  service.go            ← service interface
  event.go              ← event payload types
  error.go              ← sentinel errors

internal/
  domain/{entity}/      ← CORE: types, use-case interfaces, domain errors
                          NO framework dependencies allowed here
  service/{entity}/     ← ADAPTERS-IN: use-case implementations
                          Depends on domain interfaces only
  repository/{entity}/  ← ADAPTERS-OUT: pgx persistence
                          Implements domain port interfaces
  handler/{entity}/     ← ADAPTERS-IN: Echo HTTP handlers
                          Calls use-case interfaces only
  router/               ← mounts handlers on Echo instance
```

Dependency rule: `handler → usecase interface ← service → port interface ← repository`
Nothing in `internal/` imports from another module's `internal/`.

## Dependency Injection (Uber Fx)

All wiring happens in `fx/module.go`. The pattern per module:

```go
func Module() fx.Option {
    return fx.Module("name",
        fx.Provide(
            repository.NewRepository,   // DB → domain port
            service.NewService,         // port → use case
            service.NewAdapter,         // use case → contracts interface
            handler.NewHandler,         // use case → handler
            handler.NewRouter,          // handler → HandlerRouter
        ),
        fx.Invoke(router.Register),     // mounts routes at startup
    )
}
```

`cmd/api/serve.go` composes: `infrafx.Options()` + all module `Module()` calls.

## Event Bus Flow

```
Publisher (service)
  └─▶ eventbus.EventBus.Publish("user.created", payload)
        └─▶ LocalEventBus (in-process, goroutine per subscriber)
              ├─▶ Subscriber A handler(ctx, payload)
              └─▶ Subscriber B handler(ctx, payload)
```

- Events are string-keyed (e.g. `"user.created"`, `"user.updated"`)
- Payload types are defined in `contracts/event.go`
- Subscribers register in their module's `fx.Invoke` call

## Extension Point Mechanism

Extension points allow modules to hook into other modules' business logic without direct coupling:

```
Authorization module registers:
  registry.Register("user.after_created", priority=100, handler)

User service invokes:
  registry.Resolve("user.after_created")  →  calls all registered handlers
```

- Registry is provided by `system/extension`
- Priority controls handler execution order (lower = earlier)
- Handlers are `func(context.Context, T) error`

## Per-Module Database Migrations

Each module owns its own SQL migrations in `migrations/`:

```
modules/{name}/migrations/
  000001_create_{table}.up.sql
  000001_create_{table}.down.sql
```

The module CLI (`go run . module migrate`) runs golang-migrate against each module's migration directory in dependency order (topological sort).

## Module Manifest (module.yaml)

```yaml
name: user
version: 1.0.0
description: User management module
status: stable
dependencies:
  required: []
  optional: []
provides:
  services:
    - UserService
  events:
    - user.created
    - user.updated
extension_points:
  - user.after_created
migrations:
  path: migrations/
  transactional: true
config:
  prefix: user
```

The module management service reads manifests to:
1. Validate dependency declarations
2. Compute topological install order
3. Detect circular dependencies (Kahn's algorithm)
4. Track install/enable/disable state in DB

## Go Workspace Structure

All modules are separate Go modules with isolated `go.mod` files, linked via `go.work`:

```
go.work
  use .                          ← github.com/dinhtp/lee-goo (root)
  use ./modules/core
  use ./modules/user
  use ./modules/authentication
  use ./modules/authorization
```

Each module `go.mod` has a `replace` directive pointing back to root:
```
replace github.com/dinhtp/lee-goo => ../..
```

This enforces that module-to-platform imports go through the defined public API and never cross module boundaries at the `internal/` level.
