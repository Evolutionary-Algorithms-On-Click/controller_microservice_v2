/*
Then main file of the binary executable that will be compiled
*/
package main

import (
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"log"
)

func main() {
	logger, err := pkg.NewLogger("DEVELOPMENT")
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	logger.Info().Msg("[MSG]: Application starting...")
	// since under development, the url and token is hardcoded
	_, err = jupyterclient.NewClient("http://localhost:8888", "HelloThereHowAreyou!")
	if err != nil {
		logger.Error().Msg(err.Error())
		// should i return here? or just run the server until service is discovered?
	}
}
