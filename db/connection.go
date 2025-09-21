package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(ctx context.Context) error {
	// This in future will be gathered from consul!
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	Pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Ping the database to verify the connection
	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	return nil
}