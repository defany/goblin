# goblin

Utility library for Go projects.

## Packages

### Core

| Package | Import | Description |
|---------|--------|-------------|
| [`lc`](lc/example_test.go) | `github.com/defany/goblin/lc` | Application lifecycle manager (graceful startup & shutdown) |
| [`inject`](inject/example_test.go) | `github.com/defany/goblin/inject` | Lazy per-call-site dependency initialization |
| [`retry`](retry/example_test.go) | `github.com/defany/goblin/retry` | Retry with exponential backoff and jitter |
| [`cond`](cond/example_test.go) | `github.com/defany/goblin/cond` | Ternary operator |

### Database

| Package | Import | Description |
|---------|--------|-------------|
| [`pg`](pg/example_test.go) | `github.com/defany/goblin/pg` | PostgreSQL client (pgx pool, middleware, tx-aware queries) |
| `pg/tx` | `github.com/defany/goblin/pg/tx` | PostgreSQL adapter for transaction manager |
| `pg/migrate` | `github.com/defany/goblin/pg/migrate` | Database migrations via goose |
| [`tx`](tx/example_test.go) | `github.com/defany/goblin/tx` | Database-agnostic transaction manager with retry |
| `river` | `github.com/defany/goblin/river` | River job queue repository |

### Logging & Errors

| Package | Import | Description |
|---------|--------|-------------|
| [`slogx`](slogx/example_test.go) | `github.com/defany/goblin/slogx` | slog handler factories (JSON, Text, Pretty), ErrAttr, noop logger |
| [`errfmt`](errfmt/example_test.go) | `github.com/defany/goblin/errfmt` | Error wrapping with caller source info |
| [`rt`](rt/example_test.go) | `github.com/defany/goblin/rt` | Runtime utilities (caller name, file, unique key) |

## Requirements

Go 1.26+
