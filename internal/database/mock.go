// Package database provides functionality for managing database connections and operations.
// It includes methods for initializing the database, performing migrations, and executing
// CRUD operations on user data. This package is essential for interacting with the application's
// persistent storage layer.
package database

import (
	"errors"
	"sync"

	"github.com/dashboard-platform/auth-service/models"
)

// MockupDatabase represents a mock implementation of the database.
// It is used for testing purposes and simulates CRUD operations on user data.
type MockupDatabase struct {
	mu sync.Mutex             // Mutex to ensure thread-safe operations.
	db map[string]models.User // Map to store user data by user ID.
}

// NewMockupDatabase creates a new instance of the mock database.
//
// Returns:
//   - *MockupDatabase: A pointer to the newly created mock database.
func NewMockupDatabase() *MockupDatabase {
	return &MockupDatabase{
		mu: sync.Mutex{},
		db: make(map[string]models.User),
	}
}

// AutoMigrate simulates the migration of database schemas.
//
// Returns:
//   - error: Always returns nil as no real migration is performed.
func (m *MockupDatabase) AutoMigrate() error {
	return nil
}

// Create inserts a new user record into the mock database.
//
// Parameters:
//   - user: The user data to insert.
//
// Returns:
//   - error: Always returns nil as no real database operation is performed.
func (m *MockupDatabase) Create(user models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.db[user.ID] = user
	return nil
}

// Fetch retrieves a user record from the mock database by email.
//
// Parameters:
//   - email: The email of the user to retrieve.
//
// Returns:
//   - models.User: The retrieved user data.
//   - error: An error if the user is not found.
func (m *MockupDatabase) Fetch(email string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.db {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

// FetchByID retrieves a user record from the mock database by ID.
//
// Parameters:
//   - id: The ID of the user to retrieve.
//
// Returns:
//   - models.User: The retrieved user data.
//   - error: An error if the user is not found.
func (m *MockupDatabase) FetchByID(id string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, ok := m.db[id]; ok {
		return user, nil
	}
	return models.User{}, errors.New("user not found")
}

// Close simulates closing the mock database connection.
// This method does nothing in the mock implementation.
func (m *MockupDatabase) Close() {
	// No action needed for mockup database
}
