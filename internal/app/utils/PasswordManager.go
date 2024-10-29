package utils

import (
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseToken(tokenString, secret string) (claims *entity.Claims, err error) {

	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
