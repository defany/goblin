package pg_test

import (
	"time"

	"github.com/defany/goblin/pg"
)

// This example shows how to create a Postgres client with a DSN.
// A running PostgreSQL instance is required.
func Example_withDSN() {
	client, err := pg.New(
		pg.WithDSN("postgresql://user:pass@localhost:5432/mydb"),
		pg.WithMaxConns(10),
		pg.WithMinConns(2),
		pg.WithMaxConnIdleTime(5 * time.Minute),
	)
	if err != nil {
		// handle error
		return
	}
	defer client.Close()

	// Use client.Query, client.QueryRow, or client.Exec to interact with the database.
	_ = client
}

// This example shows how to create a Postgres client with individual connection parameters.
func Example_withParams() {
	client, err := pg.New(
		pg.WithHost("localhost"),
		pg.WithPort(5432),
		pg.WithUser("user"),
		pg.WithPassword("pass"),
		pg.WithDatabase("mydb"),
	)
	if err != nil {
		// handle error
		return
	}
	defer client.Close()

	_ = client
}
