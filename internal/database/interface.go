package database

import "auth-service/models"

type Repository interface {
	AutoMigrate() error
	Create(user models.User) error
	Fetch(email string) (models.User, error)
	FetchByID(id string) (models.User, error)
	Close()
}
