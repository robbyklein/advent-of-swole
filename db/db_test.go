package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/robbyklein/swole/sqlc"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabase(t *testing.T) {
	// Create a context
	CTX = context.Background()

	// Get the database connection string
	dsn := os.Getenv("TEST_DATABASE_URL")

	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping DB integration test")
	}

	// Open up a connection
	conn, err := pgx.Connect(CTX, dsn)
	require.NoError(t, err, "Should connect to test database without error")

	// Create an sqlc instance
	Queries = sqlc.New(conn)

	// Run a simple query
	err = conn.Ping(CTX)
	require.NoError(t, err, "Database should respond to ping")

	// Close the connection
	conn.Close(CTX)
}
