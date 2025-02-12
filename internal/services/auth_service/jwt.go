package auth_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const tokenTTL = 99999

func (s *service) GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(s.privateKey)
}

func (s *service) GetUserID(jwtToken string) (string, error) {
	token, err := s.getToken(jwtToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token provided")
	}

	return claims["id"].(string), nil
}

func (s *service) getToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.privateKey, nil
	})
	return jwtToken, err
}
