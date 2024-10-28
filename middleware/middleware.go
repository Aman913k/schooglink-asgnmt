package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Aman913k/utils"
)

// Defining custom type for context keys
type contextKey string

const (
	EmailContextKey  = contextKey("email")
	NameContextKey   = contextKey("name")
)

// JWTAuth middleware function for validating JWT tokens
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieving the token from the Authorization header
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Removing "Bearer " prefix to get the actual token
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validating JWT token
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Storing the email from claims in the context with the custom key
		ctx := context.WithValue(r.Context(), EmailContextKey, claims.Email)
		ctx = context.WithValue(ctx, NameContextKey, claims.Name)

		// Passing control to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
