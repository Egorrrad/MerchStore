package middleware

import (
	"MerchStore/src/internal/auth"
	"MerchStore/src/internal/repository"
	"context"
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

			// Парсинг токена
			claims, err := auth.ParseRefreshToken(token, secretKey)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Извлекаем user_id и username из токена
			userID := int(claims["user_id"].(float64))
			username := claims["username"].(string)

			err = repo.ValidateRefreshToken(r.Context(), userID, token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Сохраняем username в контекст
			ctx := r.Context()
			ctx = context.WithValue(ctx, usernameKey, username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
