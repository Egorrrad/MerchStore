package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetInfo(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser3",
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

	// Шаг 2: Запрос информации
	req, _ := http.NewRequest("GET", baseURL+"/api/info", nil)
	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Шаг 3: Проверка ответа
	var infoResp struct {
		Coins     int `json:"coins"`
		Inventory []struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		} `json:"inventory"`
		CoinHistory struct {
			Received []struct {
				FromUser string `json:"fromUser"`
				Amount   int    `json:"amount"`
			} `json:"received"`
			Sent []struct {
				ToUser string `json:"toUser"`
				Amount int    `json:"amount"`
			} `json:"sent"`
		} `json:"coinHistory"`
	}
	json.NewDecoder(resp.Body).Decode(&infoResp)

	// Проверка того, что данные в ответе корректные
	assert.Equal(t, infoResp.Coins, 1000)
	assert.NotNil(t, infoResp.Inventory)
	assert.NotNil(t, infoResp.CoinHistory)
}

func TestInfo_Unauthorized(t *testing.T) {
	// Шаг 1: Аутентификация с некорректным токеном
	req, _ := http.NewRequest("GET", baseURL+"/api/info", nil)
	req.Header.Set("Authorization", "Bearer invalid_token") // Неверный токен
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
