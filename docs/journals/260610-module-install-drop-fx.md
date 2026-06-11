# 2026-06-10 — Drop fx from module install/uninstall

## Context

`module install` and `module uninstall` were one-shot CLI commands that booted a full fx
app — wiring config, logger, DB pool, and all three module migration sources — just to
select one source by name and run a single migration. The fx lifecycle added ~30ms startup
cost and required every `install.go` import list to be manually updated for each new module.

## What changed

**New pattern:** `buildRunner(name)` in `cmd/module/runner.go` —
fails fast for unknown names before opening any DB connection, then constructs a
`migrator.Runner` scoped to exactly one source.

- Each module fx package now exposes `Source() migrator.Source` alongside `MigrationSource()`.
  `MigrationSource()` delegates to `Source()` — no duplicated struct literals, server app
  fx wiring untouched.
- `cmd/module/dsn.go` — shared `buildDSN` helper (avoids DSN format string living in two places).
- `cmd/module/runner.go` — `moduleRegistry` map + `buildRunner`. Adding a new module = one
  registry entry, not a new fx.Option in install.go.
- `install.go` / `uninstall.go` — zero fx imports, identical structure using `buildRunner`.
- `--force` flag removed from uninstall (was declared, never read — YAGNI).

## Fixes from code review

Two pre-existing bugs in `system/migrator/migrator.go` surfaced during review:

1. **Context timeout was silently ignored.** `runUp`/`runDown` discarded the context param;
   golang-migrate v4 has no context-aware `Up`/`Down`. Fixed by adding `stopOnCancel` —
   wires `ctx.Done()` to `m.GracefulStop` so the CLI's 30s deadline is actually enforced.

2. **Advisory-lock connection leak on `newMigrate` failure.** Error path called `db.Close()`
   but the pg driver holds a separate `*sql.Conn` for the advisory lock. Fixed to call
   `driver.Close()` which closes both handles. Practically unreachable today but formally
   correct.

## Key decisions

- Kept `MigrationSource()` intact — breaking the fx DI path for the server app was out of scope.
- Extracted `buildRunner` (DRY) because install and uninstall share identical
  config/logger/DSN/runner construction — ~12 lines that would have been duplicated.
- Did not consolidate DSN format string in `system/migrator/fx.go` — pre-existing duplication,
  out of surgical scope for this refactor.
