package auth_service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const tokenTTL = 99999

func (s *service) GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(s.privateKey)
}

func (s *service) GetUserID(jwtToken string) (int, error) {
	token, err := s.getToken(jwtToken)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token provided")
	}

	id, err := strconv.Atoi(claims["id"].(string))
	if err != nil {
		return 0, errors.New("invalid id value")
	}

	return id, nil
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
