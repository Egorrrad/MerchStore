package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestPurchaseItem(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser1",
		"password": "testpassword",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&authResp)
	token := authResp.Token

	// Шаг 2: Покупка предмета
	req, _ := http.NewRequest("GET", baseURL+"/api/buy/pink-hoody", nil)
	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBuyItem_Unauthorized(t *testing.T) {
	// Шаг 1: Попытка покупки без авторизации
	req, _ := http.NewRequest("GET", baseURL+"/api/buy/pink-hoody", nil)
	req.Header.Set("Authorization", "Bearer invalid_token") // Неверный токен
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestBuyItem_InternalServerError(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser11",
		"password": "testpassword",
	}
	authBody, _ := json.Marshal(authReq)
	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&authResp)
	token := authResp.Token

	// Шаг 2: Попытка купить предмет, что вызовет ошибку на сервере
	req, _ := http.NewRequest("GET", baseURL+"/api/buy/broken-item", nil) // Предмет, который вызывает ошибку
	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
