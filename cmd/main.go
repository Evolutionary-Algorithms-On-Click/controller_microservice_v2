/*
Then main file of the binary executable that will be compiled
*/
package main

import (
	"context"
	"log"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
)

func main() {
	logger, err := pkg.NewLogger("DEVELOPMENT")
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Logger = &logger
	pkg.Logger.Info().Msg("[MSG]: Application starting...")
	// since under development, the url and token is hardcoded
	client, err := jupyterclient.NewClient("http://localhost:8888", "HelloThereHowAreyou!")
	if err != nil {
		pkg.Logger.Error().Msg(err.Error())
		// should i return here? or just run the server until service is discovered?
	}

	// starting a kernel to see if it works
	_, err = client.StartKernel(context.Background(), "python3")
	if err != nil {
		pkg.Logger.Error().Msg(err.Error())
	}
	_, err = client.StartKernel(context.Background(), "python3")
	if err != nil {
		pkg.Logger.Error().Msg(err.Error())
	}
	_, err = client.GetKernels(context.Background())
	if err != nil {
		pkg.Logger.Error().Msg(err.Error())
	}
}
