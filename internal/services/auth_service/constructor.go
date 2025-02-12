package auth_service

import (
	"os"
)

type service struct {
	privateKey []byte
}

type Service interface {
	GenerateJWT(userID string) (string, error)
	GetUserID(jwtToken string) (string, error)
}

func New() Service {
	return &service{privateKey: []byte(os.Getenv("JWT_PRIVATE_KEY"))}
}
