# 2026-06-01 — CMD Root Command Refactor

## Summary

Unified three separate `cmd/` entry points into a single `lee-goo` root Cobra command. The repo now has one binary built from `main.go` at the repo root; `api`, `module`, and `worker` are level-2 subcommands with a max depth of 3.

## Changes

- **go.mod**: Added `github.com/spf13/cobra v1.10.2` as an explicit direct dependency (was only transitive via workspace).
- **main.go** (new): Sole entry point at repo root — calls `cmd.Execute()`.
- **cmd/root.go** (new): `lee-goo` root command; registers `api`, `module`, `worker` via `AddCommand`.
- **cmd/api/**: Created `cmd.go` + `serve.go`; deleted `main.go`. Fx app composition moved to `ServeCmd()`.
- **cmd/module/**: Created `cmd.go`; moved 15 command files from `commands/` up one level, renamed `package commands` → `package module`; deleted `main.go` and `commands/` dir.
- **cmd/worker/**: Created `cmd.go` + `start.go` (TODO stub); deleted `main.go`.
- **Makefile**: Updated `run-api`, `run-module`, `migrate-up`, and `build` targets.
- **tests/cli/**: Updated smoke tests from `go run ./cmd/module` → `go run . module`. All 3 pass.
- **README**: Updated Quick Start, Module CLI, Make Targets, Repository Structure, and Adding a Module sections.
- **.gitignore**: Added `/lee-goo` binary.

## Key Decisions

- Parent commands (`api`, `module`, `worker`) use `RunE: cmd.Help()` — prints subcommand list without extra dependencies, satisfying the "no-arg shows help" requirement.
- `go work sync` used over `go mod tidy` — `tidy` can't resolve workspace-local module imports and errors out.
- `cmd/module/commands/` entirely removed; flat layout in `cmd/module/` is simpler and avoids an intermediate package.

## Outcome

`go run . --help`, `go run . module list`, `go run . worker start` all work. Build and CLI tests pass cleanly.
