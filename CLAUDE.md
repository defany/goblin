# Goblin

Utility library for Go projects. Module: `github.com/defany/goblin`

## Build & Test

```bash
go build ./...
go test ./...
```

## Code Style

### General
- Go 1.26 — use all modern features (see below)
- Flat package structure, no `pkg/` prefix
- Public interfaces, private implementations where appropriate
- Functional options pattern for configuration (`WithX` functions)
- Context-based dependency injection (e.g. `lc.Defer(ctx, fn)`)

### Modern Go Features (mandatory)
- `any` instead of `interface{}`
- `cmp.Or(opts...)` for first non-zero value instead of manual if/else
- `new(val)` for pointer creation (Go 1.26): `new(30)` not `x := 30; &x`
- `errors.AsType[T](err)` instead of `errors.As(err, &target)`
- `for i := range n` instead of `for i := 0; i < n; i++`
- `slices`, `maps`, `cmp` packages over manual loops
- `arrutil.Map` from `github.com/gookit/goutil/arrutil` for slice transformations instead of manual make+for loops
- `math/rand/v2` instead of `math/rand`

### File Structure
- Unexported (private) functions must be placed at the bottom of the file, after all exported functions

### Naming
- Short package names: `pg`, `tx`, `lc`, not `postgres`, `transaction`, `lifecycle`
- `opt` for resolved options variable, not `o`
- Functional options: `WithX` prefix
- Errors: `ErrXxx` format with package prefix in message: `fmt.Errorf("pg: ...")`

### Dependencies
- `github.com/sourcegraph/conc/pool` for concurrent goroutine management
- `github.com/jackc/pgx/v5` for PostgreSQL
- `github.com/riverqueue/river` for job queues
- `github.com/pressly/goose/v3` for migrations
- `github.com/gookit/goutil/arrutil` for slice utilities
- `golang.org/x/sync/singleflight` for dedup

### Architecture
- Packages must not have circular dependencies
- Context is used for cross-cutting concerns (tx injection, lifecycle registration)
- `pg/` contains InjectTx/ExtractTx — shared between `pg`, `tx`, `river`
- `tx/` is database-agnostic; `pg/tx/` adapts it for PostgreSQL
- `lc/` lifecycle manager: `Go()` for heads, `Defer(ctx, fn)` for cleanup (LIFO)
