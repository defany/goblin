package migrate

import (
	"context"
	"os"

	"github.com/defany/goblin/errfmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	provider *goose.Provider
}

func New(pool *pgxpool.Pool, dir string, opts ...goose.ProviderOption) (*Migrator, error) {
	provider, err := goose.NewProvider(goose.DialectPostgres, stdlib.OpenDBFromPool(pool), os.DirFS(dir), opts...)
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	return &Migrator{provider: provider}, nil
}

func (m *Migrator) Up(ctx context.Context) ([]*goose.MigrationResult, error) {
	res, err := m.provider.Up(ctx)
	return res, errfmt.WithSource(err)
}

func (m *Migrator) Down(ctx context.Context) (*goose.MigrationResult, error) {
	res, err := m.provider.Down(ctx)
	return res, errfmt.WithSource(err)
}
