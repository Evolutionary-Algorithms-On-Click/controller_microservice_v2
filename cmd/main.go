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
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/culler"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/routes"
	"github.com/joho/godotenv"
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
	for i := range MaxDBRetries {
		err := db.InitDB(ctx)
		if err == nil {
			logger.Info().Msg("[DB]: Database connection pool successfully initialized.")
			return nil
		}

		logger.Warn().Err(err).Msgf("Database connection failed (Attempt %d/%d). Retrying in %s...",
			i+1, MaxDBRetries, DBRetryDelay)
		time.Sleep(DBRetryDelay)
	}
	return fmt.Errorf("failed to initialize database after %d retries", MaxDBRetries)
}

func main() {
	// Initializing the environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("[CRASH]: Could not load .env file: %v", err)
	}

	// Initialization of Logger
	logger, err := pkg.NewLogger(os.Getenv("APP_ENV"))
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Logger = &logger
	pkg.Logger.Info().Msg("[LOGGER]: Logger initialization successful!")

	// Initialization of Database with Retry Logic
	ctx := context.Background()
	if err := initDBWithRetry(ctx, logger); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize database connection")
	}
	defer db.Pool.Close()

	if os.Getenv("APP_ENV") == "DEVELOPMENT" {
		if err := db.LoadSchema(ctx); err != nil {
			pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize database schema")
		} else {
			pkg.Logger.Info().Msg("[DB]: Database schema successfully initialized (development mode).")
		}
	}

	jupyterGatewayURL := os.Getenv("JUPYTER_GATEWAY_URL")
  jupyterAuthToken := os.Getenv("JUPYTER_AUTH_TOKEN")
	jupyterGateway, err := jupyterclient.NewClient(jupyterGatewayURL, jupyterAuthToken)
	if err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Could not create Jupyter client/connection check failed")
		return
	}
	pkg.Logger.Info().Msg("[MSG]: Connection with jupyter kernel gateway initialized successfully!")

	// Initializing culler
	go culler.StartCuller(context.Background(), jupyterGateway)

	// Initialization of HTTP Server
	mux := http.NewServeMux()
	routes.RegisterAPIRoutes(mux, jupyterGateway)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allows all origins for development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)
	addr := ":8080"
	pkg.Logger.Info().Msg(fmt.Sprintf("[MSG]: Server listening on %s", addr))
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize the http server")
	}
}
