# Project Roadmap

## Current State (Phase 1 — Complete)

Core infrastructure fully implemented and stable.

| Area | Status | Notes |
|------|--------|-------|
| Go workspace + go.work | ✅ Done | 5 modules (root + 4 feature modules) |
| system/ infrastructure layer | ✅ Done | config, database, eventbus, extension, http, logger, security, fx |
| modules/core | ✅ Done | lifecycle domain, service, repository, handler, migrations |
| modules/user | ✅ Done | CRUD, bcrypt, events, extension point |
| modules/authentication | ✅ Done | stateless JWT login/refresh, event publishing |
| modules/authorization | ✅ Done | RBAC roles/permissions, permission cache |
| Module manifest (module.yaml) | ✅ Done | parsed, validated, topological sort |
| Topological sort + cycle detection | ✅ Done | Kahn's algorithm, 5 unit tests |
| Per-module SQL migrations | ✅ Done | golang-migrate, dependency-ordered |
| Integration test harness | ✅ Done | pkg/testapp, skip without DB |
| CLI smoke tests | ✅ Done | tests/cli via go run |
| Docker Compose (Postgres) | ✅ Done | deployment/docker-compose.yml |
| Multi-stage Dockerfile | ✅ Done | static binary, alpine final image |

---

## Phase 2 — Complete Module CLI (Near-term)

All 6 `module` subcommands are currently stubs. Wire them to the actual service layer.

| Command | Priority | Blocker |
|---------|----------|---------|
| `module list` | High | needs fx app initialization for disk discovery |
| `module install <name>` | High | needs full lifecycle service wiring |
| `module enable <name>` | High | same |
| `module disable <name>` | High | same |
| `module uninstall <name>` | High | same |
| `module make <name>` | Medium | scaffold generator (template) |

**Key work:** Initialize a minimal fx app within CLI context (no HTTP server) to access DB and service layer. Factor `pkg/cliapp` as a variant of `pkg/testapp` without HTTP.

---

## Phase 3 — Async Event Bus

Current `localEventBus` is synchronous (handlers run on the caller's goroutine).

| Item | Description |
|------|-------------|
| Goroutine-per-publish | Move handler invocation off the caller's call stack |
| Ordered delivery | Configurable per-topic FIFO channel |
| Dead-letter handling | Log + store failed events for retry |
| Backpressure | Bounded channel with configurable capacity |
| Context cancellation | Respect `ctx.Done()` in handler loops |

Migration path: keep `EventBus` interface unchanged; swap `localEventBus` implementation.

---

## Phase 4 — Generic Extension Point Registry

`system/extension.ExtensionRegistry` currently uses `any` for handler types (no compile-time safety).

| Item | Description |
|------|-------------|
| `ExtensionPoint[T]` generic interface | Type-safe register/resolve |
| Priority-ordered execution | Currently appended, priority field unused |
| Per-point registration | Named point factory, not flat string map |
| Error aggregation | Collect all handler errors, not just log |

Blocked by: Go generics support in fx (requires fx v1.22+ with `fx.Module` + generic annotations).

---

## Phase 5 — Worker Infrastructure

`worker start` command is an empty stub. Background job processing needs design.

| Item | Description |
|------|-------------|
| Worker entrypoint | Full fx app composition without HTTP server |
| Job queue abstraction | Interface for queue backends (Redis, NATS, SQS) |
| Event-driven jobs | Subscribe to domain events → enqueue background jobs |
| Graceful shutdown | Drain in-flight jobs on SIGTERM |
| Retry + backoff | Configurable per job type |

Decision needed: queue backend selection (Redis Streams vs NATS JetStream vs PostgreSQL LISTEN/NOTIFY).

---

## Phase 6 — Additional Business Modules

| Module | Priority | Dependencies | Description |
|--------|----------|--------------|-------------|
| `notification` | High | user, authentication | Email/SMS/push via event subscriptions |
| `audit-log` | High | all modules | Subscribe to all domain events, persist audit trail |
| `payment` | Medium | user, order | Stripe/PayPal adapter pattern |
| `order` | Medium | user, inventory | Order creation and lifecycle |
| `inventory` | Medium | — | Stock management |
| `report` | Low | order, inventory | Aggregation queries |

Each follows the same module structure pattern as existing modules.

---

## Known Issues to Fix

| # | Issue | Severity | File |
|---|-------|----------|------|
| 1 | Dockerfile references `modules/module_manager/` (old name) | High | `deployment/Dockerfile` |
| 2 | `module` CLI: 6 commands exist, all stubs (no service wiring yet) | High | `cmd/module/*.go` |
| 3 | Default role hook in authz is a no-op TODO | Medium | `modules/authorization/fx/module.go` |
| 4 | `POST /users/:id/roles` handler exists but not wired in router | Medium | `modules/authorization/internal/handler/role/router.go` |
| 5 | ~~`system/http` uses `slog.Default()` in goroutine instead of injected logger~~ Fixed: `system/server/fx.go` now injects `*logger.Logger` | Resolved | `system/server/fx.go` |
| 6 | `core` Upsert does not write lifecycle timestamps | Low | `modules/core/internal/repository/module/module_repository.go` |
| 7 | `worker start` completely unimplemented | Low | `cmd/worker/start.go` |
| 8 | `pkg/testapp.WithConfig()` is a no-op placeholder | Low | `pkg/testapp/options.go` |
| 9 | Integration tests only validate fx wiring; no behavioral assertions yet | Low | `tests/integration/` |

---

## Guiding Principles (Non-Negotiable)

- No ORM — raw `database/sql` + explicit SQL migrations
- No runtime plugin loading — compiled-only modules
- No cross-module `internal/` imports — contracts boundary enforced
- 200 LOC per file limit
- Interface at every layer boundary
- YAGNI / KISS / DRY
