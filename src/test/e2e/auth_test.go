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
