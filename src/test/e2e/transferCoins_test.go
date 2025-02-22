package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// Структура запроса для отправки монетки
type SendCoinRequest struct {
	ReceiverName string `json:"toUser"`
	Amount       int    `json:"amount"`
}

func TestTransferCoins(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser2",
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

	// Шаг 2: Отправка монеток
	sendCoinReq := SendCoinRequest{
		ReceiverName: "testuser1", // Имя получателя
		Amount:       10,          // Количество монет
	}
	sendCoinBody, _ := json.Marshal(sendCoinReq)
	req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(sendCoinBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendCoin_BadRequest(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser22",
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

	// Шаг 2: Отправка некорректного запроса (отсутствует поле amount)
	sendCoinReq := SendCoinRequest{
		ReceiverName: "testuser2", // Некорректно, не указано количество монет
	}
	sendCoinBody, _ := json.Marshal(sendCoinReq)
	req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(sendCoinBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSendCoin_Unauthorized(t *testing.T) {
	// Шаг 1: Отправка запроса без авторизации
	sendCoinReq := SendCoinRequest{
		ReceiverName: "testuser2",
		Amount:       10,
	}
	sendCoinBody, _ := json.Marshal(sendCoinReq)
	req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(sendCoinBody))
	req.Header.Set("Authorization", "Bearer invalid_token") // Неправильный токен
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTransferCoinsToSelf(t *testing.T) {
	// Шаг 1: Аутентификация
	authReq := map[string]string{
		"username": "testuser25",
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

	// Шаг 2: Отправка монеток
	sendCoinReq := SendCoinRequest{
		ReceiverName: "testuser25", // Имя получателя
		Amount:       10,           // Количество монет
	}
	sendCoinBody, _ := json.Marshal(sendCoinReq)
	req, _ := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewBuffer(sendCoinBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
