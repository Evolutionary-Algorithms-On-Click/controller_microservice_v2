package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

// LoadSchema initializes the database schema in development environment
func LoadSchema(ctx context.Context) error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "DEVELOPMENT" {
		pkg.Logger.Info().Msg("[DB]: Skipping schema initialization, not in development mode.")
		return nil
	}

	// Step 1: Run DROP statements first (reverse FK order recommended)
	if err := executeSchemaFile(ctx, "db/schema_drop.sql"); err != nil {
		return fmt.Errorf("failed during schema_drop.sql: %w", err)
	}

	// Step 2: Wait for schema changes to finish (ignore GC jobs)
	if err := waitForSchemaChanges(ctx); err != nil {
		return fmt.Errorf("waiting for schema drop changes failed: %w", err)
	}

	// Step 3: Run CREATE statements
	if err := executeSchemaFile(ctx, "db/schema_create.sql"); err != nil {
		return fmt.Errorf("failed during schema_create.sql: %w", err)
	}

	// Step 4: Wait for create schema changes
	if err := waitForSchemaChanges(ctx); err != nil {
		return fmt.Errorf("waiting for schema create changes failed: %w", err)
	}

	pkg.Logger.Info().Msg("[DB]: Schema initialization completed successfully.")
	return nil
}

// executeSchemaFile reads and executes the given SQL file
func executeSchemaFile(ctx context.Context, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire DB connection: %w", err)
	}
	defer conn.Release()

	pkg.Logger.Info().Msgf("[DB]: Executing %s ...", filePath)
	_, err = conn.Exec(ctx, string(content))
	if err != nil {
		pkg.Logger.Warn().Err(err).Msgf("[DB]: Error executing %s (non-fatal, Cockroach may be applying changes)", filePath)
	} else {
		pkg.Logger.Info().Msgf("[DB]: %s executed successfully.", filePath)
	}

	return nil
}

// waitForSchemaChanges waits for all active CREATE/ALTER schema changes to finish.
// DROP jobs or jobs in GC TTL are ignored.
func waitForSchemaChanges(ctx context.Context) error {
	pkg.Logger.Info().Msg("[DB]: Waiting for active schema changes to complete...")

	timeout := time.After(5 * time.Minute)
	pollInterval := 2 * time.Second

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for schema changes to finish")
		case <-time.After(pollInterval):
			conn, err := Pool.Acquire(ctx)
			if err != nil {
				pkg.Logger.Error().Err(err).Msg("[DB]: Failed to acquire connection during polling")
				pollInterval *= 2 // exponential backoff
				if pollInterval > 10*time.Second {
					pollInterval = 10 * time.Second
				}
				continue
			}

			var runningJobs int

			// Only wait for CREATE/ALTER jobs; skip DROP and GC TTL jobs
			query := `
				SELECT count(*)
				FROM [SHOW JOBS]
				WHERE job_type = 'SCHEMA CHANGE'
				  AND status = 'running'
				  AND (running_status IS NULL OR running_status NOT LIKE '%waiting for GC TTL%')
				  AND description NOT LIKE '%DROP%'
			`

			err = conn.QueryRow(ctx, query).Scan(&runningJobs)
			conn.Release()

			if err != nil {
				pkg.Logger.Error().Err(err).Msg("[DB]: Error while checking schema jobs")
				continue
			}

			if runningJobs == 0 {
				pkg.Logger.Info().Msg("[DB]: All active schema changes completed.")
				return nil
			}

			pkg.Logger.Info().
				Int("running_jobs", runningJobs).
				Msg("[DB]: Still waiting for schema changes...")

			// Optional: log all active CREATE/ALTER jobs for debugging
			if runningJobs > 0 {
				conn, err := Pool.Acquire(ctx)
				if err == nil {
					rows, err := conn.Query(ctx, `
						SELECT job_id, description, status, running_status
						FROM [SHOW JOBS]
						WHERE status='running'
						  AND job_type='SCHEMA CHANGE'
						  AND description NOT LIKE '%DROP%'
					`)
					if err == nil {
						for rows.Next() {
							var jobID int64
							var desc, status, runStatus string
							if err := rows.Scan(&jobID, &desc, &status, &runStatus); err != nil {
								pkg.Logger.Err(err).Msg("[DB]: Falied to execute job check query")
							}
							pkg.Logger.Warn().Msgf("[DB]: Active job: %d | %s | %s | %s", jobID, desc, status, runStatus)
						}
						rows.Close()
					}
					conn.Release()
				}
			}
		}
	}
}
