package auth

import (
	"MerchStore/src/cmd"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var cfg, _ = cmd.Load()

func GenerateRefreshToken(userID int, username string) (string, error) {
	secret := cfg.Service.SecretKey
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseRefreshToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
