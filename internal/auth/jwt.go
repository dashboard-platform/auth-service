// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTObj represents the JWT utility object.
// It provides methods for creating and validating JWT tokens.
type JWTObj struct {
	Secret []byte // Secret key used for signing and validating JWT tokens.
}

// CreateJWT generates a new JWT token for a given user ID.
//
// Parameters:
//   - userID: The ID of the user for whom the token is being created.
//
// Returns:
//   - string: The signed JWT token.
//   - error: An error if the token creation fails.
func (j *JWTObj) CreateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                                // Subject (user ID).
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Expiration time.
		"iat": time.Now().Unix(),                     // Issued at time.
		"iss": "auth-service",                        // Issuer.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

// ValidateJWT validates a given JWT token and extracts the user ID.
//
// Parameters:
//   - tokenStr: The JWT token string to validate.
//
// Returns:
//   - string: The user ID extracted from the token.
//   - error: An error if the token is invalid or verification fails.
func (j *JWTObj) ValidateJWT(tokenStr string) (string, error) {
	errToken := errors.New("invalid token")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
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
