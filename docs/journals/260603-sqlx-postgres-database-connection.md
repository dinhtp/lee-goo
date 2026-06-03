# sqlx Postgres Database Connection Migration

**Date:** 2026-06-03
**Plan:** `plans/260603-0939-sqlx-postgres-database-connection/`

## What Changed

Replaced `database/sql`-backed `system/database/connection.go` with `jmoiron/sqlx` sourced from gotility (copied, not imported).

New layout:
```
system/database/
├── connection.go      ← adapter shim: wraps postgresql.Connection, exposes DB() *sql.DB
├── fx.go              ← builds DSN + pool config, provides database.Connection via Fx
└── postgresql/
    ├── config.go      ← Config struct (pool settings) + sentinel errors
    ├── connection.go  ← Connection interface returning *sqlx.DB, NewConnection(dsn, cfg)
    └── transaction.go ← Transact() helper (commit/rollback wrapper)
```

`system/config/config.go` `DatabaseConfig` extended with 4 pool fields (`MaxOpenConnections`, `MaxIdleConnections`, `ConnectionMaxTime`, `ConnectionIdleTime`); defaults `max_open=4`, `max_idle=2`.

## Key Decision: Adapter Shim

Module repositories (`modules/*/internal/repository/`) accepted `database.Connection` and called `.DB() *sql.DB`. The brainstorm incorrectly assumed no consumers existed.

Rather than updating all repositories, a thin `connectionAdapter` struct wraps `postgresql.Connection` and returns the embedded `sql.DB` field from `sqlx.DB`. No repository changes needed. New code that needs `*sqlx.DB` directly can inject `postgresql.Connection`.

## Tests

6 unit tests pass without a DB (sentinel-error paths, default-config, ping-failure). 1 integration test (`TestTransactRollsBackOnCallbackError`) skipped with `-short`; requires real Postgres.

## Unresolved

- Existing repositories still use `*sql.DB` (`ExecContext`, `QueryRowContext`) — sqlx named-query and struct-scan features are available but unused until repos are migrated.
