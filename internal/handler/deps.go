package handler

import (
	"MerchandiseShop/internal/storage"
	"context"
)

type Storage interface {
	GetUserByUsername(ctx context.Context, userID string) (storage.User, error)
	CreateUser(ctx context.Context, username string, password string) (storage.User, error)
}

type authService interface {
	GenerateJWT(userID string) (string, error)
	GetUserID(jwtToken string) (string, error)
}
