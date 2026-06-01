# Code Standards

Go coding conventions enforced in this project.

## File & Package Naming

- **Files**: `snake_case.go` — e.g. `user_repository.go`, `jwt_security.go`
- **Packages**: match the directory leaf name, lowercase, no underscores where avoidable
  - `package fx` (in `modules/user/fx/`)
  - `package tests` (in `modules/user/tests/`)
  - Exception: avoid collision with stdlib — prefix as needed (e.g. `infrafx`, `infrahttp`)
- **File size limit**: keep individual `.go` files under 200 lines; split by concern if exceeded

## Interface-Driven Design

Every layer boundary is expressed as an interface, never a concrete type:

```go
// domain layer defines the port
type UserPort interface {
    Create(ctx context.Context, u *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

// compile-time assertion in the implementing file
var _ UserPort = (*Repository)(nil)
```

Compile-time interface assertions (`var _ Interface = (*Impl)(nil)`) are mandatory in every file that claims to implement a public interface.

## No Framework in Domain Layer

`internal/domain/{entity}/` must NOT import:
- `echo`, `cobra`, `pgx`, `fx`, or any HTTP/DB/DI package
- Only standard library + domain types are allowed

```go
// Good — pure domain
type User struct {
    ID        string
    Email     string
    CreatedAt time.Time
}

// Bad — framework leaked into domain
type User struct {
    gorm.Model          // violates domain purity
    Email string `json:"email"`
}
```

## Error Handling

- Use sentinel errors (`var ErrNotFound = errors.New("not found")`) in `contracts/error.go`
- Wrap errors with `fmt.Errorf("context: %w", err)` to preserve the chain
- Never swallow errors silently; every `err != nil` check must log or return
- HTTP handlers map domain errors to HTTP status codes explicitly — no generic 500s

```go
// Good
if err != nil {
    return fmt.Errorf("user repository Create: %w", err)
}

// Bad
if err != nil {
    return nil  // silent failure
}
```

## Fx Module Pattern

Every module must expose a single `Module() fx.Option` function in `fx/module.go`:

```go
func Module() fx.Option {
    return fx.Module("name",
        fx.Provide(
            repository.New,   // concrete type, satisfies port interface
            service.New,      // port interface → use case
            handler.New,      // use case → handler
            handler.NewRouter,
        ),
        fx.Invoke(router.Register),
    )
}
```

Rules:
- `fx.Module("name", ...)` — always name the fx module for debug output
- `fx.Provide` only; never call constructors directly in `fx.Invoke`
- Use `fx.Annotate` + `fx.As` when providing a concrete type as an interface to avoid ambiguous binding
- `fx.NopLogger` in tests to suppress fx startup logs

## Context Propagation

- Every function that does I/O must accept `context.Context` as the first parameter
- Never store a context in a struct
- Pass the request context from Echo handlers through the full call chain

## Configuration

- Config structs live in `config/config.go` per module
- All config is loaded from environment variables via Viper (system/config)
- No hardcoded values in business logic; use the config struct

## Testing

- Unit tests: file suffix `_test.go`, package `tests` (external test package for `internal/`)
- Integration tests: live in `tests/integration/`, skip with `t.Skip` when no DB
- CLI tests: live in `tests/cli/`, use `os/exec` to invoke `go run ./cmd/...`
- Use `github.com/stretchr/testify/require` for fatal assertions, `assert` for non-fatal
- Compile-time smoke tests: `func TestXModuleCompiles(t *testing.T) { t.Log("...") }`

## Commit Style

Conventional commits:
```
feat: add user password reset endpoint
fix: correct JWT expiry calculation for refresh tokens
refactor: extract topological sort into standalone function
test: add diamond dependency test case for TopologicalSort
```

No `chore` or `docs` prefixes for `.claude/` directory changes.
