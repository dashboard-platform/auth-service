package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTObj struct {
	Secret []byte
}

func (j *JWTObj) CreateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

func (j *JWTObj) ValidateJWT(tokenStr string) (string, error) {
	errToken := errors.New("invalid token")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errToken
		}
		return j.Secret, nil
	})

	if err != nil || !token.Valid {
		return "", errToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", errToken
	}

	return sub, nil
}
