package models

import "time"

type RegisterAPI struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAPI struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
