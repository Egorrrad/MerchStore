package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser4",
		"password": "testpassword",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Шаг 2: Проверка, что токен получен
	var authResp struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NotEmpty(t, authResp.Token)
}

func TestAuthRepeat(t *testing.T) {
	// повторная авторизация того же пользователя
	TestAuth(t)
}

func TestAuth_EmptyUsername(t *testing.T) {
	// Шаг 1: Аутентификация с пустым именем пользователя
	authReq := map[string]string{
		"username": "",
		"password": "testpassword",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // Ожидаем ошибку из-за пустого имени пользователя

	// Шаг 2: Проверка, что в ответе есть сообщение об ошибке
	var authResp map[string]string
	json.NewDecoder(resp.Body).Decode(&authResp)
	t.Log(authResp)
	assert.Equal(t, "username: is required", authResp["errors"])
}

func TestAuth_EmptyPassword(t *testing.T) {
	// Шаг 1: Аутентификация с пустым паролем
	authReq := map[string]string{
		"username": "testuser4",
		"password": "",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // Ожидаем ошибку из-за пустого пароля

	// Шаг 2: Проверка, что в ответе есть сообщение об ошибке
	var authResp map[string]string
	json.NewDecoder(resp.Body).Decode(&authResp)
	t.Log(authResp)
	assert.Equal(t, "password: is required", authResp["errors"])
}
