package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
	"github.com/rs/zerolog"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("hijack not supported")
	}
	return h.Hijack()
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		// Select log level based on status code
		var event *zerolog.Event
		switch {
		case rec.status >= 500:
			event = pkg.Logger.Error()
		case rec.status >= 400:
			event = pkg.Logger.Warn()
		case rec.status >= 300:
			event = pkg.Logger.Info() // or Warn() if redirects matter
		default:
			event = pkg.Logger.Info()
		}

		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Int("status", rec.status).
			Dur("duration", duration).
			Msg("incoming HTTP request")
	})
}
