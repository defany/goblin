package pg

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Option func(*config)

type config struct {
	dsn      string
	host     string
	port     int
	user     string
	password string
	database string

	maxConns          *int32
	minConns          *int32
	maxConnIdleTime   *time.Duration
	maxConnLifetime   *time.Duration
	healthCheckPeriod *time.Duration

	tracer      pgx.QueryTracer
	middlewares []Middleware
}

func (c *config) buildDSN() string {
	if c.dsn != "" {
		return c.dsn
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		c.user, c.password, c.host, c.port, c.database,
	)
}

func WithDSN(dsn string) Option {
	return func(c *config) {
		c.dsn = dsn
	}
}

func WithHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

func WithPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

func WithUser(user string) Option {
	return func(c *config) {
		c.user = user
	}
}

func WithPassword(password string) Option {
	return func(c *config) {
		c.password = password
	}
}

func WithDatabase(database string) Option {
	return func(c *config) {
		c.database = database
	}
}

func WithMaxConns(n int32) Option {
	return func(c *config) {
		c.maxConns = new(n)
	}
}

func WithMinConns(n int32) Option {
	return func(c *config) {
		c.minConns = new(n)
	}
}

func WithMaxConnIdleTime(d time.Duration) Option {
	return func(c *config) {
		c.maxConnIdleTime = new(d)
	}
}

func WithMaxConnLifetime(d time.Duration) Option {
	return func(c *config) {
		c.maxConnLifetime = new(d)
	}
}

func WithHealthCheckPeriod(d time.Duration) Option {
	return func(c *config) {
		c.healthCheckPeriod = new(d)
	}
}

func WithTracer(tracer pgx.QueryTracer) Option {
	return func(c *config) {
		c.tracer = tracer
	}
}

func WithMiddlewares(mws ...Middleware) Option {
	return func(c *config) {
		c.middlewares = append(c.middlewares, mws...)
	}
}
