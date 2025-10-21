/*
Then main file of the binary executable that will be compiled
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/routes"
	"github.com/rs/cors"
)

func main() {
	// Initialization of Logger
	logger, err := pkg.NewLogger("DEVELOPMENT")
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Logger = &logger
	pkg.Logger.Info().Msg("[MSG]: Logger initialization successful!")

	// Initialization of Database
	ctx := context.Background()
	if err := db.InitDB(ctx); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize database connection")
	}
	defer db.Pool.Close()

	// Initialization of HTTP Server
	mux := http.NewServeMux()
	routes.RegisterAPIRoutes(mux)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // e.g. http://localhost:3000 or should get the frontend url from consul
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)
	addr := ":8080"
	pkg.Logger.Info().Msg(fmt.Sprintf("Server listening on %s", addr))
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		pkg.Logger.Fatal().Err(err).Msg("[CRASH]: Failed to initialize the http server")
	}

	// since under development, the url and token is hardcoded
	_, err = jupyterclient.NewClient("http://localhost:8888", "HelloThereHowAreyou!")
	if err != nil {
		pkg.Logger.Error().Msg(err.Error())
		// should i return here? or just run the server until service is discovered?
	}

	
}
