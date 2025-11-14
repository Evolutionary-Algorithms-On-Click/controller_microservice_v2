package middleware

import (
	"net/http"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RequestLogger logs each incoming HTTP request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // default
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		pkg.Logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Int("status", rec.status).
			Dur("duration", duration).
			Msg("incoming HTTP request")
	})
}
