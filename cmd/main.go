/*
Then main file of the binary executable that will be compiled
*/
package main

import (
	"log"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

func main() {
	logger, err := pkg.NewLogger("DEVELOPMENT")
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	logger.Info().Msg("[MSG]: Application starting...")
}
