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
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/middleware"
	"github.com/Thanus-Kumaar/controller_microservice_v2/routes"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

const (
	MaxDBRetries = 5
	DBRetryDelay = 3 * time.Second
)

// Initialize DB with retry logic for Docker race conditions
func initDBWithRetry(ctx context.Context, logger zerolog.Logger) error {
	for i := range MaxDBRetries {
		err := db.InitDB(ctx)
		if err == nil {
			logger.Info().Msg("[DB]: Database connection pool successfully initialized.")
			return nil
		}

		logger.Warn().Err(err).Msgf(
			"Database connection failed (Attempt %d/%d). Retrying in %s...",
			i+1, MaxDBRetries, DBRetryDelay,
		)
		time.Sleep(DBRetryDelay)
	}
	return fmt.Errorf("failed to initialize database after %d retries", MaxDBRetries)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("[CRASH]: Could not load .env file: %v", err)
	}

	// Initialize Logger
	logger, err := pkg.NewLogger(os.Getenv("APP_ENV"))
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Logger = &logger
	pkg.Logger.Info().Msg("[LOGGER]: Logger initialization successful!")

	// Initialize DB
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

	// === gRPC Connection to Auth Microservice ==========================

	// authGrpcAddr := os.Getenv("AUTH_GRPC_SERVER_ADDRESS")
	// if authGrpcAddr == "" {
	// 	authGrpcAddr = "localhost:5001"
	// 	pkg.Logger.Warn().Msgf(
	// 		"AUTH_GRPC_SERVER_ADDRESS not set, defaulting to %s",
	// 		authGrpcAddr,
	// 	)
	// }

	// var authConn *grpc.ClientConn

	// for i := range MaxDBRetries {
	// 	pkg.Logger.Info().Msgf(
	// 		"Attempting to connect to Auth gRPC service at %s (Attempt %d/%d)...",
	// 		authGrpcAddr, i+1, MaxDBRetries,
	// 	)

	// 	// Replaces deprecated grpc.WithTimeout()
	// 	ctx, cancel := context.WithTimeout(context.Background(), DBRetryDelay)
	// 	defer cancel()

	// 	//nolint:staticcheck
	// 	conn, err := grpc.DialContext(
	// 		ctx,
	// 		authGrpcAddr,
	// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 		grpc.WithBlock(),
	// 	)

	// 	if err == nil {
	// 		authConn = conn
	// 		pkg.Logger.Info().Msg("[AUTH]: Auth gRPC service connection successful.")
	// 		break
	// 	}

	// 	pkg.Logger.Warn().Err(err).Msgf(
	// 		"Auth gRPC connection failed. Retrying in %s...",
	// 		DBRetryDelay,
	// 	)

	// 	time.Sleep(DBRetryDelay)
	// }

	// if authConn == nil {
	// 	pkg.Logger.Fatal().Msg("[CRASH]: Failed to connect to Auth gRPC service after multiple retries")
	// }
	// defer authConn.Close()

	// authMiddleware := middleware.NewAuthMiddleware(authConn, logger)

	// === Jupyter Gateway Initialization =================================

	jupyterGatewayURL := os.Getenv("JUPYTER_GATEWAY_URL")
	jupyterAuthToken := os.Getenv("JUPYTER_AUTH_TOKEN")

	jupyterGateway, err := jupyterclient.NewClient(jupyterGatewayURL, jupyterAuthToken)
	if err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Could not create Jupyter client / connection check failed")
		return
	}
	pkg.Logger.Info().Msg("[MSG]: Connection with Jupyter kernel gateway initialized successfully!")

	// Start kernel culler
	go culler.StartCuller(context.Background(), jupyterGateway)

	// === HTTP Server =====================================================

	mux := http.NewServeMux()
	// routes.RegisterAPIRoutes(mux, jupyterGateway, authMiddleware)
	routes.RegisterAPIRoutes(mux, jupyterGateway)
	loggedMux := middleware.RequestLogger(mux)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(loggedMux)

	addr := ":8080"

	pkg.Logger.Info().Msg(fmt.Sprintf("[MSG]: Server listening on %s", addr))

	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize the HTTP server")
	}
}
