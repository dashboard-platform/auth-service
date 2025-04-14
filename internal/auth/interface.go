package auth

import "auth-service/models"

type ServiceInterface interface {
	Register(data models.RegisterAPI) (string, error)
	Login(data models.LoginAPI) (string, error)
	GetUserByID(id string) (models.User, error)
}

type JWTValidator interface {
	ValidateJWT(token string) (string, error)
}
