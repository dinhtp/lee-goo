# 2026-06-04 — Config: Env Field + IsDevelopment Helper

## What

Compared `golang-scaffold/config` (raw `os.Getenv` + `sync.Once` singleton) against lee-goo's Viper-based config. Adopted the single useful pattern: `IsDevelopment()` helper.

## Changes

- `system/config/config.go`: added `Env string \`mapstructure:"env"\`` field, `v.SetDefault("env", "local")` default, and `IsDevelopment()` method
- `.env.example`: added `ENV=production` so operators know the variable exists

## Key Decisions

**lee-goo's Viper approach wins on every axis** — proper `(*Config, error)` return, sensible defaults, type-safe `time.Duration` pool settings, DI-friendly (no global singleton). golang-scaffold's `sync.Once` global makes test isolation impossible.

**Only `IsDevelopment()` was worth backporting.** golang-scaffold has `IsDevelopment()`, `Database.DataSourceName()`, and `Server.URL()` — but lee-goo already has `DataSourceName()` on the `Connection` interface (better separation), so only the env helper was genuinely missing.

**Used `static.EnvLocal`/`static.EnvDev` constants** instead of inline strings to avoid DRY violation; the constants already existed in the codebase.

**Skipped case normalization** (`strings.ToLower`) — YAGNI; operators supply lowercase values by convention.
