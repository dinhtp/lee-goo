# 2026-06-03 — Replace slog with zap-based Logger

## What changed

Replaced the 37-line `slog` stub in `system/logger/` with a richer zap-based `Logger`
copied from `gotility/logger`. Five new files: `level.go`, `option.go`, `logger.go`,
`echo.go`, `fx.go`. Two call-sites updated: `system/eventbus` and `system/server`.

## Key decisions

- **Copied, not imported.** Source is in-tree; no dependency on `gotility` module.
- **filter/ and casbin.go excluded.** Lee-goo has no use for field-redaction filters or
  Casbin log adapter. Stripping them simplified the port significantly.
- **fx.go normalises level strings.** `strings.ToUpper(cfg.Log.Level)` maps config
  values like `"info"` to `Level` constants — avoids a switch statement in the provider.
- **go mod tidy skipped.** The go.work monorepo setup makes `go mod tidy` fail without
  network access to workspace-local modules. Deps promoted manually in `go.mod`; verified
  with `go build ./...`.

## Outcome

- `go build ./...` clean
- `system/server` and `tests/cli` pass
- Pre-existing DB test failure unrelated (requires Docker)
- `docs/codebase-summary.md` and `docs/project-roadmap.md` updated
- Roadmap item 5 (slog.Default() in server goroutine) resolved as a side effect
