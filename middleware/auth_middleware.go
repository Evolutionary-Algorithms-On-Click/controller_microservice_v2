package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/Thanus-Kumaar/controller_microservice_v2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type contextKey string

const (
	UserContextKey = contextKey("user")
)

type User struct {
	ID       string
	Role     string
	Email    string
	UserName string
	FullName string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("t")
		if err != nil {
			http.Error(w, "Authorization cookie required", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(w, "Authorization cookie is empty", http.StatusUnauthorized)
			return
		}

		conn, err := grpc.NewClient(os.Getenv("AUTH_GRPC_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			http.Error(w, "Failed to connect to auth service", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := proto.NewAuthenticateClient(conn)
		res, err := client.Auth(context.Background(), &proto.TokenValidateRequest{Token: tokenString})
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusInternalServerError)
			return
		}

		if !res.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user := &User{
			ID:       res.Id,
			Role:     res.Role,
			Email:    res.Email,
			UserName: res.UserName,
			FullName: res.FullName,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
