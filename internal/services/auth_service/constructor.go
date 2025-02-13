package auth_service

import (
	"os"
)

type service struct {
	privateKey []byte
}

type Service interface {
	GenerateJWT(userID int) (string, error)
	GetUserID(jwtToken string) (int, error)
}

func New() Service {
	return &service{privateKey: []byte(os.Getenv("JWT_PRIVATE_KEY"))}
}
