// Package models defines the core data structures used across the application.
// These include API request/response models and database entity models.
package models

import "time"

// RegisterAPI represents the data structure for user registration requests.
type RegisterAPI struct {
	Email    string `json:"email"`    // The email address of the user.
	Password string `json:"password"` // The password of the user.
}

// LoginAPI represents the data structure for user login requests.
type LoginAPI struct {
	Email    string `json:"email"`    // The email address of the user.
	Password string `json:"password"` // The password of the user.
}

// User represents the database entity for a user.
type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`        // The unique identifier of the user.
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`     // The email address of the user.
	PasswordHash string    `gorm:"not null" json:"-"`                     // The hashed password of the user.
	AuthProvider string    `gorm:"default:password" json:"auth_provider"` // The authentication provider used by the user (e.g., password, google).
	CreatedAt    time.Time `json:"created_at"`                            // The timestamp when the user was created.
	UpdatedAt    time.Time `json:"updated_at"`                            // The timestamp when the user was last updated.
}
