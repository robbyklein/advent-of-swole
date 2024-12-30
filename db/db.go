package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/robbyklein/swole/sqlc"
)

var Queries *sqlc.Queries
var CTX context.Context

func ConnectToDatabase() {
	// Get connection string
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL not set, cannot connect to database")
	}

	// Create context
	CTX = context.Background()

	// Establish a connection
	conn, err := pgx.Connect(CTX, dsn)

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Create an sqlc instance with it
	Queries = sqlc.New(conn)
}
