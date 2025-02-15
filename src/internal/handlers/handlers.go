package handlers

import (
	"MerchStore/src/internal/storage/model"
	"context"
	"encoding/json"
	"net/http"
)

type Repository interface {
	BuyItem(ctx context.Context, username, productName string) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserInfo(ctx context.Context, username string) ([]model.Purchase, []model.Operation, error)
	SendCoins(ctx context.Context, fromUser, toUser string, amount int) error
	PostAuthUser(ctx context.Context, username, password string) (*string, error)
}

type Server struct {
	repo Repository
}

func NewServer(repository Repository) Server {
	return Server{
		repo: repository,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
