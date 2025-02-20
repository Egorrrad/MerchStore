package middleware

import (
	"MerchStore/src/internal/repository"
	"context"
	"errors"
	"net/http"
)

const usernameKey string = "username"

// AuthMiddleware возвращает middleware для проверки JWT-токена.
func AuthMiddleware(repo repository.Repository, secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notAuth := []string{"/api/auth"}
			requestPath := r.URL.Path

			for _, value := range notAuth {
				if value == requestPath {
					next.ServeHTTP(w, r)
					return
				}
			}

			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			username, err := repo.ValidateRefreshToken(r.Context(), token, secretKey)
			if err != nil {
				switch {
				case errors.Is(err, repository.ErrMsgTokenExpired):
					http.Error(w, "Token expired", http.StatusUnauthorized)
				default:
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
				return
			}

			// Сохраняем username в контекст
			ctx := r.Context()
			ctx = context.WithValue(ctx, usernameKey, *username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
