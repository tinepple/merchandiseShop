package handler

import (
	"MerchandiseShop/internal/storage"
	"context"
)

type Storage interface {
	GetUserByUsername(ctx context.Context, userName string) (storage.User, error)
	CreateUser(ctx context.Context, username string, password string) (storage.User, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
	GetItem(ctx context.Context, name string) (storage.Item, error)
	CreatePurchase(ctx context.Context, userID, itemID, newBalance int) error
	CreateTransaction(ctx context.Context, transaction storage.Transaction, NewBalanceUserFrom int, NewBalanceUserTo int) error
	GetPurchasesByUserID(ctx context.Context, userID int) ([]storage.Inventory, error)
	GetTransactionsByUserID(ctx context.Context, userID int) ([]storage.CoinsHistory, error)
}

type authService interface {
	GenerateJWT(userID int) (string, error)
	GetUserID(jwtToken string) (int, error)
}
