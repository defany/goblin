package pg

import (
	"context"
	"fmt"

	"github.com/defany/goblin/errfmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool        *pgxpool.Pool
	middlewares []Middleware
}

func New(opts ...Option) (*Postgres, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	dsn := cfg.buildDSN()
	if dsn == "" {
		return nil, errfmt.WithSource(fmt.Errorf("pg: dsn or connection params required"))
	}

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errfmt.WithSource(fmt.Errorf("pg: parse config: %w", err))
	}

	if cfg.maxConns != nil {
		poolCfg.MaxConns = *cfg.maxConns
	}
	if cfg.minConns != nil {
		poolCfg.MinConns = *cfg.minConns
	}
	if cfg.maxConnIdleTime != nil {
		poolCfg.MaxConnIdleTime = *cfg.maxConnIdleTime
	}
	if cfg.maxConnLifetime != nil {
		poolCfg.MaxConnLifetime = *cfg.maxConnLifetime
	}
	if cfg.healthCheckPeriod != nil {
		poolCfg.HealthCheckPeriod = *cfg.healthCheckPeriod
	}
	if cfg.tracer != nil {
		poolCfg.ConnConfig.Tracer = cfg.tracer
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, errfmt.WithSource(fmt.Errorf("pg: connect: %w", err))
	}

	return &Postgres{
		pool:        pool,
		middlewares: cfg.middlewares,
	}, nil
}

func (p *Postgres) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	ctx, query, args, err := applyMiddlewares(ctx, p.middlewares, query, args)
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	if tx := ExtractTx(ctx); tx != nil {
		rows, err := tx.Query(ctx, query, args...)
		return rows, errfmt.WithSource(err)
	}

	rows, err := p.pool.Query(ctx, query, args...)
	return rows, errfmt.WithSource(err)
}

func (p *Postgres) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	ctx, query, args, err := applyMiddlewares(ctx, p.middlewares, query, args)
	if err != nil {
		return errorRow{err: errfmt.WithSource(err)}
	}

	if tx := ExtractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.pool.QueryRow(ctx, query, args...)
}

func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	ctx, query, args, err := applyMiddlewares(ctx, p.middlewares, query, args)
	if err != nil {
		return pgconn.CommandTag{}, errfmt.WithSource(err)
	}

	if tx := ExtractTx(ctx); tx != nil {
		tag, err := tx.Exec(ctx, query, args...)
		return tag, errfmt.WithSource(err)
	}

	tag, err := p.pool.Exec(ctx, query, args...)
	return tag, errfmt.WithSource(err)
}

func (p *Postgres) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, opts)
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *Postgres) Close() {
	p.pool.Close()
}

type errorRow struct {
	err error
}

func (r errorRow) Scan(_ ...any) error {
	return r.err
}
