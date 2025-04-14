// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import "github.com/dashboard-platform/auth-service/models"

type ServiceInterface interface {
	Register(data models.RegisterAPI) (string, error)
	Login(data models.LoginAPI) (string, error)
	GetUserByID(id string) (models.User, error)
}

type JWTValidator interface {
	ValidateJWT(token string) (string, error)
}
