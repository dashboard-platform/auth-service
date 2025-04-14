package database

import (
	"testing"
	"time"

	"auth-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMockupDatabase verifies that a new instance of MockupDatabase is properly initialized.
func TestNewMockupDatabase(t *testing.T) {
	// Create a new mockup database.
	db := NewMockupDatabase()
	// Assert that the database instance is not nil.
	require.NotNil(t, db)
	// Assert that the internal map is initialized and empty.
	assert.Equal(t, 0, len(db.db))
}

// TestAutoMigrate verifies that AutoMigrate runs without error.
func TestAutoMigrate(t *testing.T) {
	db := NewMockupDatabase()
	err := db.AutoMigrate()
	// Since AutoMigrate does nothing in the mock, we expect no error.
	assert.NoError(t, err)
}

// TestCreateAndFetch tests the Create, Fetch (by email) and FetchByID functions.
func TestCreateAndFetch(t *testing.T) {
	db := NewMockupDatabase()

	// Create a dummy user.
	now := time.Now()
	user := models.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Test the Create method.
	err := db.Create(user)
	require.NoError(t, err, "creating a user should not produce an error")

	// Test the Fetch method by email.
	fetchedByEmail, err := db.Fetch("test@example.com")
	require.NoError(t, err, "should find user by email")
	assert.Equal(t, user.Email, fetchedByEmail.Email, "emails should match")

	// Test the FetchByID method.
	fetchedByID, err := db.FetchByID("user-1")
	require.NoError(t, err, "should find user by ID")
	assert.Equal(t, user.ID, fetchedByID.ID, "IDs should match")
}

// TestFetchNonExistingUser tests that fetching a non-existent user returns an error.
func TestFetchNonExistingUser(t *testing.T) {
	db := NewMockupDatabase()

	// Attempt to fetch a user by email that hasn't been created.
	_, err := db.Fetch("nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error(), "error message should indicate user not found")

	// Attempt to fetch a user by a non-existent ID.
	_, err = db.FetchByID("nonexistent-id")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error(), "error should indicate user not found")
}

// TestClose simply verifies that the Close method can be called without error.
// (In this mock implementation, Close does nothing.)
func TestClose(t *testing.T) {
	db := NewMockupDatabase()
	// Call Close; since it does nothing, we only check for absence of panic.
	db.Close()
}
