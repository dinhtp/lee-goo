# 2026-06-04 — golang-migrate Integration

## What shipped

Wired `github.com/golang-migrate/migrate/v4` into the modular monorepo. Each module now controls its own DB schema via embedded SQL files registered at compile time. CLI trigger: `go run . module install <name>`.

## Key decisions

**Per-module tracking table** — each module gets `schema_migrations_<name>` (not a shared table). Prevents any cross-module migration interference and makes rollback scoped.

**iofs driver, zero file I/O** — `//go:embed *.sql` in each module's `migrations/fs.go`. No runtime filesystem access; migrations ship in the binary.

**`MigrationSource()` separate from `Module()`** — CLI apps compose a lean `MigrateOptions()` fx.App (config + logger + db + migrator only) without pulling in HTTP server, eventbus, extension, or security. Keeps `module install` startup fast and side-effect free.

**`app.Start()` not `app.Run()`** — `fx.Invoke` fires during Start. Using Run() would block indefinitely. Explicit `app.Stop()` after ensures DB connection lifecycle is cleaned up.

**`runErr` closure pattern** — `fx.Invoke` is synchronous within `app.Start`, so capturing errors via a closure variable is safe and avoids needing to return errors from the invoke function.

**authentication module skipped** — owns no tables, so no `migrations/fs.go` or `MigrationSource()`. Intentional; authentication is stateless JWT.

## Non-obvious detail

`sqlx.DB` embeds `*sql.DB` as the `.DB` field. The migrator's fx params accept `*sqlx.DB` (what's in the container) and passes `p.DB.DB` to `migratepg.WithInstance`. This reuses the existing connection pool rather than opening a second one.

## Result

`go build ./...` passes clean. `go run . module install <name>` builds a minimal fx.App, runs UP migrations for the named module only, and exits.
