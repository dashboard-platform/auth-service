// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"testing"

	"github.com/dashboard-platform/auth-service/internal/database"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService_Register tests the Register method of the authentication service.
// It verifies that valid registration data creates a new user and invalid data
// returns appropriate errors.
func TestService_Register(t *testing.T) {
	auth := NewService(database.NewMockupDatabase())

	tests := []struct {
		name         string             // Name of the test case.
		data         models.RegisterAPI // Input data for the Register method.
		expectedData bool               // Expected result: true if registration succeeds.
		expectedErr  bool               // Expected error: true if an error is expected.
	}{
		{
			name: "Register valid",
			data: models.RegisterAPI{
				Name:     "Test User",
				Email:    "test@test.com",
				Password: "securepass",
			},
			expectedData: true,
			expectedErr:  false,
		},
		{
			name: "Register invalid email",
			data: models.RegisterAPI{
				Name:     "Test User",
				Email:    "test",
				Password: "securepass",
			},
			expectedData: false,
			expectedErr:  true,
		},
		{
			name: "Register invalid password",
			data: models.RegisterAPI{
				Name:     "Test User",
				Email:    "test@test.com",
				Password: "",
			},
			expectedData: false,
			expectedErr:  true,
		},
		{
			name: "Register empty name",
			data: models.RegisterAPI{
				Name:     "",
				Email:    "test@test.com",
				Password: "securepassword",
			},
			expectedData: false,
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := auth.Register(tt.data)

			assert.Equal(t, tt.expectedData, id != "")
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

// TestService_Login tests the Login method of the authentication service.
// It verifies that valid login credentials return the correct user ID and
// invalid credentials return appropriate errors.
func TestService_Login(t *testing.T) {
	auth := NewService(database.NewMockupDatabase())
	expectedID, err := auth.Register(models.RegisterAPI{
		Email:    "valid@test.com",
		Password: "securepass",
	})
	if err != nil {
		t.Fatalf("error registering test data: %v", err)
	}

	tests := []struct {
		name         string          // Name of the test case.
		data         models.LoginAPI // Input data for the Login method.
		expectedData bool            // Expected result: true if login succeeds.
		expectedErr  bool            // Expected error: true if an error is expected.
	}{
		{
			name: "Login valid",
			data: models.LoginAPI{
				Email:    "valid@test.com",
				Password: "securepass",
			},
			expectedData: true,
			expectedErr:  false,
		},
		{
			name: "Login invalid email",
			data: models.LoginAPI{
				Email:    "",
				Password: "securepass",
			},
			expectedData: false,
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := auth.Login(tt.data)

			require.Equal(t, tt.expectedData, id != "")
			require.Equal(t, tt.expectedErr, err != nil)

			if id != "" {
				require.Equal(t, expectedID, id)
			}
		})
	}
}
