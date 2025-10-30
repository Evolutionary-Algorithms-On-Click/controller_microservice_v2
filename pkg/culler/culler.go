package culler

import (
	"context"
	"time"
	"os"
	"strconv"

	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

// StartCuller starts a background process that periodically removes idle kernels.
func StartCuller(ctx context.Context, jupyterClient *jupyterclient.Client) {
	logger := pkg.Logger

	// Load intervals from env (with fallback defaults)
	cullIntervalMin, _ := strconv.Atoi(getEnvOrDefault("CULL_INTERVAL_MINUTES", "10"))
	idleThresholdMin, _ := strconv.Atoi(getEnvOrDefault("IDLE_THRESHOLD_MINUTES", "30"))

	cullInterval := time.Duration(cullIntervalMin) * time.Minute
	idleThreshold := time.Duration(idleThresholdMin) * time.Minute

	ticker := time.NewTicker(cullInterval)
	logger.Info().Msgf("[CULLER]: Started. Checking every %d minutes, idle threshold = %d minutes", cullIntervalMin, idleThresholdMin)

	go func() {
		for {
			select {
			case <-ticker.C:
				cullIdleKernels(ctx, jupyterClient, idleThreshold)
			case <-ctx.Done():
				logger.Warn().Msg("[CULLER]: Context cancelled, stopping culler.")
				ticker.Stop()
				return
			}
		}
	}()
}

func cullIdleKernels(ctx context.Context, client *jupyterclient.Client, threshold time.Duration) {
	logger := pkg.Logger
	logger.Info().Msg("[CULLER]: Running idle kernel check...")

	kernels, err := client.GetKernels(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("[CULLER]: Failed to get kernels, attempting reconnect...")

		// Attempt to reinitialize the Jupyter client connection
		newClient, connErr := jupyterclient.NewClient("http://localhost:8888", "YOUR_SECRET_TOKEN")
		if connErr != nil {
			logger.Error().Err(connErr).Msg("[CULLER]: Reconnection to Jupyter Gateway failed.")
			return
		}
		logger.Info().Msg("[CULLER]: Reconnected to Jupyter Gateway successfully.")
		*client = *newClient
		return
	}

	now := time.Now().UTC()
	for _, k := range *kernels {
		idleDuration := now.Sub(k.LastActivity)
		if idleDuration > threshold {
			logger.Warn().Str("kernel_id", k.ID).Str("language", k.Name).Dur("idle_for", idleDuration).
				Msg("[CULLER]: Kernel exceeded idle threshold. Deleting...")

			if err := client.DeleteKernel(ctx, k.ID); err != nil {
				logger.Error().Err(err).Str("kernel_id", k.ID).Msg("Failed to delete idle kernel")
				continue
			}

			logger.Info().Str("kernel_id", k.ID).Msg("[CULLER]: Kernel deleted successfully.")
		}
	}

	logger.Info().Msg("[CULLER]: Idle kernel check complete.")
}

// TODO: Useful function. Can be seperated and used elsewhere if needed
func getEnvOrDefault(key, def string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return def
}
