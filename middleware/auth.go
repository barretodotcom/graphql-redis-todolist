package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/barretodotcom/graphql-redis-todolist/pkg/jwt"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.Split(authorizationHeader, " ")[1]

			userId, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userId)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}
