package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ (в проде лучше хранить в .env)
var secretKey = []byte("supersecretkey")

// Claims структура для хранения данных в токене
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Генерация JWT-токена
func GenerateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // 3 дня
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Проверка и разбор токена
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
