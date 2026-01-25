package pkg

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"time"
)

var Logger *zerolog.Logger

func NewLogger(env string, logLevel string) (zerolog.Logger, error) {
	var output io.Writer

	switch env {
	case "DEVELOPMENT":
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	case "PRODUCTION":
		output = os.Stderr // For production, output to stderr by default. Consider file or other sink for actual deployments.
	default:
		return zerolog.Logger{}, fmt.Errorf("invalid environment for logger setup: %s (Allowed: DEVELOPMENT, PRODUCTION)", env)
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	// Set global log level
	switch strings.ToLower(logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default to info level
	}

	return logger, nil
}
