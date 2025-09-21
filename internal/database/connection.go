package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/deedubs/choochoo/internal/db"
)

// Connection manages database connection
type Connection struct {
	conn    *pgx.Conn
	queries *db.Queries
}

// NewConnection creates a new database connection
func NewConnection(ctx context.Context) (*Connection, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/choochoo?sslmode=disable"
		log.Printf("Warning: DATABASE_URL not set, using default: %s", dbURL)
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := db.New(conn)

	return &Connection{
		conn:    conn,
		queries: queries,
	}, nil
}

// Queries returns the sqlc queries instance
func (c *Connection) Queries() *db.Queries {
	return c.queries
}

// Close closes the database connection
func (c *Connection) Close(ctx context.Context) error {
	if c.conn != nil {
		return c.conn.Close(ctx)
	}
	return nil
}

// IsConnected checks if the database connection is active
func (c *Connection) IsConnected(ctx context.Context) bool {
	if c.conn == nil {
		return false
	}
	return c.conn.Ping(ctx) == nil
}