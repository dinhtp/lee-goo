# 2026-06-05 — Healthz Liveness API

## What

Added `GET /healthz` liveness probe to `modules/core`. Returns `{"status":"ok"}` with HTTP 200. No auth, no DB.

## Files

- **Created:** `internal/handler/health/contract.go`, `handler.go`, `router.go`
- **Modified:** `fx/module.go` — wired `healthHandler.NewHandler` + `healthHandler.Register`
- **Test:** `tests/health_handler_test.go` — `TestHealthzEndpoint` via Echo httptest

## Key Decisions

**Bypassed `HandlerRouter` interface.** The existing module handler uses `HandlerRouter.Register(*echo.Group)` + a parent `Register(*echo.Echo, HandlerRouter)` that creates a prefixed group. `/healthz` is root-level with no middleware, so that indirection adds nothing — `Register(app *echo.Echo, h *Handler)` calling `app.GET` directly is simpler.

**No domain/service/repository layers.** Pure handler — no state, no DB. Adding layers would be YAGNI.

**Plain httptest for the test.** The existing testapp harness requires DB env vars. The health handler has zero external dependencies, so a direct Echo httptest is cleaner and runs without any infrastructure.

## Result

All 6 tests pass. Build clean. External monitors, k8s liveness probes, and load balancers can now verify process liveness at `/healthz`.
