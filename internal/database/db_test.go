// Package database provides functionality for managing database connections and operations.
// It includes methods for initializing the database, performing migrations, and executing
// CRUD operations on user data. This package is essential for interacting with the application's
// persistent storage layer.
package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dashboard-platform/auth-service/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgresContainer starts a PostgreSQL container using testcontainers-go
// and returns the container along with a DSN (Data Source Name) string.
//
// Parameters:
//   - t: The testing context.
//
// Returns:
//   - tc.Container: The started PostgreSQL container.
//   - string: The DSN string for connecting to the database.
func setupPostgresContainer(t *testing.T) (tc.Container, string) {
	ctx := context.Background()

	// Define container request with the PostgreSQL image and environment variables.
	req := tc.ContainerRequest{
		Image:        "postgres:15",        // Use PostgreSQL version 15.
		ExposedPorts: []string{"5432/tcp"}, // Expose port 5432.
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_DB":       "authdb",
		},
		// Wait until PostgreSQL is ready to accept connections.
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Build the DSN based on the container's host and port.
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=secret dbname=authdb sslmode=disable", host, mappedPort.Port())
	return container, dsn
}

// TestDatabaseOperations tests the functionality in db.go.
// It covers initialization, automigration, creation, and query operations for a user.
//
// This test uses a PostgreSQL container to simulate a real database environment.
func TestDatabaseOperations(t *testing.T) {
	ctx := context.Background()
	// Start the PostgreSQL container.
	container, dsn := setupPostgresContainer(t)
	// Ensure that the container is terminated after the test finishes.
	defer container.Terminate(ctx)

	// Create a logger; here, os.Stdout is used for simplicity.
	logger := zerolog.New(os.Stdout)

	// Initialize the database connection using the Init function.
	db, err := Init(dsn, logger)
	require.NoError(t, err)
	defer db.Close()

	// Run AutoMigrate to create the required schema (e.g., the User table).
	err = db.AutoMigrate()
	require.NoError(t, err)

	// Create a dummy user to test database operations.
	now := time.Now()
	dummyUser := models.User{
		ID:           uuid.New().String(),
		Email:        "dummy@example.com",
		PasswordHash: "dummyhash",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = db.Create(dummyUser)
	require.NoError(t, err)

	// Test the Fetch method by email.
	fetchedUser, err := db.Fetch(dummyUser.Email)
	require.NoError(t, err)
	assert.Equal(t, dummyUser.Email, fetchedUser.Email, "Emails should match")

	// Test the FetchByID method by user ID.
	fetchedByID, err := db.FetchByID(dummyUser.ID)
	require.NoError(t, err)
	assert.Equal(t, dummyUser.ID, fetchedByID.ID, "User IDs should match")

	// Test the FetchByID method with an invalid ID.
	invalidID := uuid.New().String()
	fetchedUser, err = db.FetchByID(invalidID)
	require.Empty(t, fetchedUser, "Expected empty user when fetching with invalid ID")
	require.Error(t, err, "Expected error when fetching with invalid ID")

	// Test the Fetch method with an invalid email.
	dummyUser.Email = "wrong@example.com"
	fetchedUser, err = db.Fetch(dummyUser.Email)
	require.Empty(t, fetchedUser, "Expected empty user when fetching with invalid email")
	require.Error(t, err, "Expected error when fetching with invalid email")
}
