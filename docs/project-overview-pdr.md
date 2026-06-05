# Project Overview — PDR (Product Design & Requirements)

## Project Summary

**Lee-Goo** is a reusable Go backend boilerplate implementing a compiled modular architecture inspired by Magento 2's module system. It enables teams to build enterprise-grade Go services where features are packaged as independent, composable modules that declare dependencies, expose typed service contracts, publish/subscribe to domain events, and extend other modules via registered hook points — all without tight coupling.

## Goals

1. **Modular isolation**: Each module has its own `go.mod`, migrations, config, and domain layer with zero cross-module `internal/` imports.
2. **Compiled safety**: All modules are compiled into a single binary; no runtime dynamic loading. Module "installation" is a DB state + migration step, not a binary operation.
3. **Lifecycle management**: A CLI (`cmd/module`) controls module install, enable, disable, upgrade, and doctor workflows with dependency-order enforcement via topological sort.
4. **Extension without coupling**: Modules extend each other via event bus subscriptions and named extension points, not direct function calls.
5. **Enterprise readiness**: Stateless JWT auth, RBAC authorization, per-module migrations, structured logging, and config from environment variables.
6. **Developer ergonomics**: Clear `make` targets, module scaffold generator, and integration test harness that skips gracefully without a DB.

## Constraints

- **No ORM**: Raw SQL via `pgx/v5` — maintains explicit control over queries and migrations.
- **No runtime plugin loading**: Modules are compile-time; the compiled binary is the final artifact.
- **No framework in domain layer**: `internal/domain/` has zero imports from Echo, pgx, fx, or Cobra.
- **Interface at every boundary**: Concrete types never cross layer boundaries; only interfaces do.
- **File size cap**: 200 lines per `.go` file to maintain LLM context manageability.

## Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|-------------|
| 1 | All packages compile without errors | `go build ./...` exits 0 |
| 2 | Topological sort correctly orders dependencies | `go test ./modules/core/...` — 5 tests pass |
| 3 | Circular dependency detection returns `ErrCircularDependency` | `TestTopologicalSort_DetectsCycle` passes |
| 4 | Module CLI commands exit 0 | `go test ./tests/cli/...` passes |
| 5 | Integration tests skip gracefully without DB | `go test ./tests/integration/...` skips (not fails) without `DATABASE_HOST` |
| 6 | User, auth, authz unit tests pass | `go test ./modules/user/... ./modules/authentication/... ./modules/authorization/...` |
| 7 | `go vet ./...` passes cleanly | No vet errors |
| 8 | Platform fx wires all providers | `cmd/api/main.go` compiles with all four modules composed |
| 9 | Liveness probe responds correctly | `GET /healthz` returns HTTP 200 `{"status":"ok"}` |

## Module Roadmap (Future)

Modules declared in the TRD not yet implemented in this boilerplate:

| Module | Priority | Notes |
|--------|----------|-------|
| `notification` | High | Email/SMS/push via event subscriptions |
| `audit-log` | High | Subscribes to all domain events, persists audit trail |
| `payment` | Medium | Stripe/PayPal adapter pattern |
| `order` | Medium | Depends on user + payment + inventory |
| `inventory` | Medium | Stock management |
| `report` | Low | Aggregation queries, depends on order + inventory |

## Key Design Decisions

### Why Uber Fx?
Fx provides constructor-based DI with lifecycle hooks (`OnStart`/`OnStop`). Its module grouping (`fx.Module`) maps cleanly to the feature module concept and produces a clear dependency graph for debugging.

### Why separate go.mod per module?
Isolates dependency graphs — a module can upgrade `golang-migrate` independently without affecting the root module. `go.work` stitches them together for local development without publishing.

### Why no ORM?
ORMs obscure query costs and migration complexity. Raw `pgx` with explicit SQL keeps migrations predictable and queries auditable.

### Why stateless JWT?
Eliminates session storage, simplifies horizontal scaling. The refresh-token pattern provides revocability without a session table.

### Why topological sort for module install order?
Module B's migration may depend on tables created by module A's migration (e.g., foreign keys). Kahn's algorithm guarantees correct migration order and detects cycles at install time.
