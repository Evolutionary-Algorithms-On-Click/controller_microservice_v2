package pkg

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

var Logger *zerolog.Logger

func NewLogger(env string) (zerolog.Logger, error) {
	var output io.Writer

	switch env {
	case "DEVELOPMENT":
		// Use a human-readable console writer for development
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	case "PRODUCTION":
		// Use JSON output for production, which is machine-readable
		// and can be easily consumed by log aggregation systems.
		file, err := os.OpenFile("prod.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return zerolog.Logger{}, err
		}
		output = file
	default:
		return zerolog.Logger{}, fmt.Errorf("invalid environment for logger setup: %s (Allowed: DEVELOPMENT, PRODUCTION)", env)
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger, nil
}
