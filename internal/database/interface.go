// Package database provides functionality for managing database connections and operations.
// It includes methods for initializing the database, performing migrations, and executing
// CRUD operations on user data. This package is essential for interacting with the application's
// persistent storage layer.
package database

import "github.com/dashboard-platform/auth-service/models"

// Repository defines the interface for database operations.
// It includes methods for schema migration, user creation, and user retrieval.
type Repository interface {
	// AutoMigrate performs database migrations.
	AutoMigrate() error

	// Create inserts a new user record into the database.
	//
	// Parameters:
	//   - user: The user data to insert.
	//
	// Returns:
	//   - error: An error if the insertion fails.
	Create(user models.User) error

	// Fetch retrieves a user record from the database by email.
	//
	// Parameters:
	//   - email: The email of the user to retrieve.
	//
	// Returns:
	//   - models.User: The retrieved user data.
	//   - error: An error if the retrieval fails.
	Fetch(email string) (models.User, error)

	// FetchByID retrieves a user record from the database by ID.
	//
	// Parameters:
	//   - id: The ID of the user to retrieve.
	//
	// Returns:
	//   - models.User: The retrieved user data.
	//   - error: An error if the retrieval fails.
	FetchByID(id string) (models.User, error)

	// Close closes the database connection.
	Close()
}
