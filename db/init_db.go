package db

import (
	"context"
	"fmt"
	"os"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

// LoadSchema initializes the database schema in development environment
func LoadSchema(ctx context.Context) error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "DEVELOPMENT" {
		pkg.Logger.Info().Msg("[DB]: Skipping schema initialization, not in development mode.")
		return nil
	}

	schemaFile := "db/schema.sql"
	content, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire DB connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema.sql: %w", err)
	}
	return nil
}
