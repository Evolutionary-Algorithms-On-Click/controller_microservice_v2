package middleware

import (
	"context"
	"net/http"
	"strings"

	pb "github.com/Thanus-Kumaar/controller_microservice_v2/proto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type ctxKey string

const (
    ctxUserIDKey   ctxKey = "userID"
    ctxUserRoleKey ctxKey = "userRole"
)

// AuthMiddleware holds the dependencies for the authentication middleware.
type AuthMiddleware struct {
	AuthClient pb.AuthenticateClient
	Logger     zerolog.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(authConn grpc.ClientConnInterface, logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		AuthClient: pb.NewAuthenticateClient(authConn),
		Logger:     logger,
	}
}

// Authenticate is the middleware handler.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Call the gRPC service
		authResponse, err := m.AuthClient.Auth(context.Background(), &pb.TokenValidateRequest{Token: token})
		if err != nil {
			m.Logger.Error().Err(err).Msg("gRPC call to auth service failed")
			http.Error(w, "Authentication service failed", http.StatusInternalServerError)
			return
		}

		if !authResponse.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context for downstream handlers
		ctx := context.WithValue(r.Context(), ctxUserIDKey, authResponse.Id)
		ctx = context.WithValue(ctx, ctxUserRoleKey, authResponse.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
