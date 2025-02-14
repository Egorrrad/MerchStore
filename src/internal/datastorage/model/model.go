package model

import "time"

type User struct {
	UserID       int
	Username     string
	PasswordHash string
	Role         string
	Coins        int
	CreatedAt    time.Time
}

type Product struct {
	ProductID int
	Name      string
	Price     int
	Quantity  int
}

type Purchase struct {
	PurchaseID    int
	UserID        int
	ProductID     int
	ProductName   string
	Quantity      int
	OperationDate time.Time
}

type Operation struct {
	OperationID      int
	SenderUserID     int
	SenderUsername   string
	ReceiverUserID   int
	ReceiverUsername string
	Amount           int
	OperationDate    time.Time
}

type RefreshToken struct {
	TokenID   int
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
