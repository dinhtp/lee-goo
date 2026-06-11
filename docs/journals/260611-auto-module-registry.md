# 2026-06-11 — Auto-generated module registry

## Context

Yesterday's refactor landed `moduleRegistry` in `cmd/module/runner.go` as a hand-maintained
map — already an improvement over per-command fx wiring, but still required a manual edit
every time a new module was added. One forgotten entry = silent failure at runtime.

## What changed

- `cmd/module/gen/main.go` (NEW) — generator that walks `modules/*/fx/migration.go`,
  derives import paths and aliases, then writes `registry.go`
- `modules/registry.go` (NEW, generated, committed to git) — registry for
  3 current modules: `authorization`, `core`, `user`
- `cmd/module/runner.go` — removed manual `moduleRegistry` var and module-specific imports;
  added `//go:generate go run ./gen` directive
- `Makefile` — added `generate` target: `go generate ./cmd/module/...`

## Key decisions

- **Discovery rule:** module included iff `modules/{name}/fx/migration.go` exists. Explicit
  sentinel file beats directory-walk heuristics — no ambiguity about what "is a module".
- **Import alias:** `{name}Module` camelCase (e.g. `authorizationModule`). Consistent with
  existing manual aliases; generator enforces it.
- **Alias collision detection:** generator calls `log.Fatal` on duplicate aliases. Loud failure
  beats a silently broken registry.
- **Generated file committed** — not gitignored. Keeps CI builds reproducible without running
  `make generate` in the pipeline; reviewers see registry diffs in PRs.

## Developer contract

Add a module: ensure `modules/{name}/fx/migration.go` exists → run `make generate` → done.
No manual registry edits. The old failure mode (forgot to add entry, wrong name at runtime)
is gone.

## Lessons

The manual map was exactly the kind of accidental complexity that accumulates between
refactors. It only took one day to become stale. Code-gen for a deterministic, file-based
discovery rule is the right call — the generator is ~50 lines and the contract is obvious.
