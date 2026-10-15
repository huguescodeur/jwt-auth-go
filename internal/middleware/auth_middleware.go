package middleware

import (
	"context"
	"jwt-auth/internal/auth"
	"jwt-auth/internal/pkg/ctxkeys"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtManager *auth.JWTManager, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return

		}

		tokenString := strings.TrimSpace(
			strings.TrimPrefix(header, "Bearer "),
		)

		if tokenString == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := jwtManager.Verify(tokenString)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
