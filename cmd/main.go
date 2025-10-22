package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/routes"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

// Global constants for retry logic
const (
	MaxDBRetries = 5
	DBRetryDelay = 3 * time.Second
)

// A wrapper to initialize the database with retry logic to solve Docker race conditions.
func initDBWithRetry(ctx context.Context, logger zerolog.Logger) error {
	for i := 0; i < MaxDBRetries; i++ {
		err := db.InitDB(ctx)
		if err == nil {
			logger.Info().Msg("Database connection pool successfully initialized.")
			return nil
		}

		logger.Warn().Err(err).Msgf("Database connection failed (Attempt %d/%d). Retrying in %s...",
			i+1, MaxDBRetries, DBRetryDelay)
		time.Sleep(DBRetryDelay)
	}
	return fmt.Errorf("failed to initialize database after %d retries", MaxDBRetries)
}

func main() {
	// Initialization of Logger
	logger, err := pkg.NewLogger(os.Getenv("APP_ENV"))
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Logger = &logger
	pkg.Logger.Info().Msg("[MSG]: Logger initialization successful!")

	// Initialization of Database with Retry Logic
	ctx := context.Background()
	if err := initDBWithRetry(ctx, logger); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize database connection")
	}
	defer db.Pool.Close()

	_, err = jupyterclient.NewClient("http://localhost:8888", "YOUR_SECRET_TOKEN")
	if err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Could not create Jupyter client/connection check failed")
		return
	}

	// Initialization of HTTP Server
	mux := http.NewServeMux()
	routes.RegisterAPIRoutes(mux)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allows all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)
	addr := ":8080"
	pkg.Logger.Info().Msg(fmt.Sprintf("Server listening on %s", addr))
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize the http server")
	}
}
