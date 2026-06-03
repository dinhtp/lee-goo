# 2026-06-03 — Echo HTTP Engine Migration: system/http → system/server

## Summary

Replaced the thin `system/http` wrapper (`Server{Echo *echo.Echo}`) with a proper `system/server`
package using the Engine interface pattern. Source adapted from `gotility/server/echo`; no gotility
import — code is self-contained in lee-goo.

## Changes

| Area | Before | After |
|------|--------|-------|
| Package | `system/http` (pkg `http`) | `system/server` (pkg `server`) |
| Abstraction | `*Server` struct, public `.Echo` field | `Engine` interface + `*echo.Echo` via fx |
| Module routers | `*systemHTTP.Server` param | `*echo.Echo` param — no system import needed |
| Tests | none | 5 unit tests in `engine_test.go` |

**Files created:** `system/server/{config,engine,engine_test,fx}.go`  
**Files deleted:** `system/http/{server,fx}.go`  
**Files updated:** `system/fx/options.go`, 4 module `router/register.go` files, README, codebase-summary

## Key Decision

`fx.Provide(NewEchoEngine)` returns `(Engine, *echo.Echo, error)`. fx injects both automatically:
- `Engine` → consumed only by `RegisterLifecycle` (OnStart/OnStop)
- `*echo.Echo` → consumed by module routers for `app.Group()`

Module routers no longer import `system/server` at all — clean separation.

## Result

- `go build ./...` clean
- `go test ./system/server/... -v` — 5/5 pass
- Full suite: all non-DB tests pass (DB tests require running PG, pre-existing)
