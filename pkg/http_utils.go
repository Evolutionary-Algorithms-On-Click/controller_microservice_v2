package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

// WriteJSONResponseWithLogger is a helper function to write JSON responses with logging.
func WriteJSONResponseWithLogger(w http.ResponseWriter, status int, data interface{}, logger *zerolog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(err).Msg("failed to write json response")
	}
}
