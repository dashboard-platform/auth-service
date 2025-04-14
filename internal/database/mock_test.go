// Package database provides functionality for managing database connections and operations.
// It includes methods for initializing the database, performing migrations, and executing
// CRUD operations on user data. This package is essential for interacting with the application's
// persistent storage layer.
package database

import (
	"testing"
	"time"

	"github.com/dashboard-platform/auth-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMockupDatabase verifies that a new instance of MockupDatabase is properly initialized.
func TestNewMockupDatabase(t *testing.T) {
	db := NewMockupDatabase()
	require.NotNil(t, db)
	assert.Equal(t, 0, len(db.db))
}

// TestAutoMigrate verifies that AutoMigrate runs without error.
func TestAutoMigrate(t *testing.T) {
	db := NewMockupDatabase()
	err := db.AutoMigrate()
	assert.NoError(t, err)
}

// TestCreateAndFetch tests the Create, Fetch (by email), and FetchByID functions.
func TestCreateAndFetch(t *testing.T) {
	db := NewMockupDatabase()

	now := time.Now()
	user := models.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := db.Create(user)
	require.NoError(t, err)

	fetchedByEmail, err := db.Fetch("test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.Email, fetchedByEmail.Email)

	fetchedByID, err := db.FetchByID("user-1")
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetchedByID.ID)
}

// TestFetchNonExistingUser tests that fetching a non-existent user returns an error.
func TestFetchNonExistingUser(t *testing.T) {
	db := NewMockupDatabase()

	_, err := db.Fetch("nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	_, err = db.FetchByID("nonexistent-id")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

// TestClose verifies that the Close method can be called without error.
func TestClose(t *testing.T) {
	db := NewMockupDatabase()
	db.Close()
}
