// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"testing"
	"time"

	"github.com/dashboard-platform/auth-service/models"

	"github.com/stretchr/testify/require"
)

// TestNewMockService tests the creation of a new mock authentication service.
// It verifies that the mock service and its internal data structures are properly initialized.
func TestNewMockService(t *testing.T) {
	mockService := NewMockService()
	require.NotNil(t, mockService)
	require.NotNil(t, mockService.users)
	require.NotNil(t, mockService.logins)
}

// TestMockRegister tests the Register method of the mock authentication service.
// It verifies that valid registration data creates a new user and invalid data
// returns appropriate errors.
func TestMockRegister(t *testing.T) {
	mockService := NewMockService()

	tests := []struct {
		name          string             // Name of the test case.
		input         models.RegisterAPI // Input data for the Register method.
		expectedError bool               // Expected error: true if an error is expected.
	}{
		{
			name: "Valid registration",
			input: models.RegisterAPI{
				Email:    "test@test.com",
				Password: "test",
			},
			expectedError: false,
		},
		{
			name: "Invalid email",
			input: models.RegisterAPI{
				Email:    "invalid-email",
				Password: "test",
			},
			expectedError: true,
		},
		{
			name: "Empty password",
			input: models.RegisterAPI{
				Email:    "test@test.com",
				Password: "",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := mockService.Register(tt.input)
			if tt.expectedError {
				require.Error(t, err)
				require.Empty(t, id)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, id)
		})
	}
}

// TestMockLogin tests the Login method of the mock authentication service.
// It verifies that valid login credentials return the correct user ID and
// invalid credentials return appropriate errors.
func TestMockLogin(t *testing.T) {
	mockService := NewMockService()
	user := models.RegisterAPI{
		Email:    "test@test.com",
		Password: "test",
	}

	id, err := mockService.Register(user)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string          // Name of the test case.
		input         models.LoginAPI // Input data for the Login method.
		expectedError bool            // Expected error: true if an error is expected.
	}{
		{
			name: "Valid login",
			input: models.LoginAPI{
				Email:    "test@test.com",
				Password: "test",
			},
			expectedError: false,
		},
		{
			name: "Invalid login",
			input: models.LoginAPI{
				Email:    "invalid@test.com",
				Password: "test",
			},
			expectedError: true,
		},
		{
			name: "Invalid email",
			input: models.LoginAPI{
				Email:    "",
				Password: "test",
			},
			expectedError: true,
		},
		{
			name: "Invalid password",
			input: models.LoginAPI{
				Email:    "test@test.com",
				Password: "",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginId, err := mockService.Login(tt.input)
			if tt.expectedError {
				require.Error(t, err)
				require.Empty(t, loginId)
				return
			}
			require.NoError(t, err)
			require.Equal(t, id, loginId)
		})
	}
}

// TestMockGetUserByID tests the GetUserByID method of the mock authentication service.
// It verifies that valid user IDs return the correct user data and invalid IDs return errors.
func TestMockGetUserByID(t *testing.T) {
	mockService := NewMockService()
	user := models.RegisterAPI{
		Email:    "test@test.com",
		Password: "test",
	}

	id, err := mockService.Register(user)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string // Name of the test case.
		input         string // Input user ID for the GetUserByID method.
		expectedError bool   // Expected error: true if an error is expected.
	}{
		{
			name:          "Valid ID",
			input:         id,
			expectedError: false,
		},
		{
			name:          "Invalid ID",
			input:         "invalid-id",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := mockService.GetUserByID(tt.input)
			if tt.expectedError {
				require.Error(t, err)
				require.Empty(t, user)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, user)
			require.Equal(t, id, user.ID)
		})
	}
}

// TestMockAddUser tests the AddUser method of the mock authentication service.
// It verifies that a user can be added to the mock service and is stored correctly.
func TestMockAddUser(t *testing.T) {
	mockService := NewMockService()
	user := models.User{
		ID:           "test",
		Email:        "test@test.com",
		PasswordHash: "test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockService.AddUser(user)
	require.NotEmpty(t, mockService.users)
	require.Equal(t, user.ID, mockService.users[user.ID].ID)
}

// TestMockDeleteUser tests the DeleteUser method of the mock authentication service.
// It verifies that a user can be deleted from the mock service and is no longer stored.
func TestMockDeleteUser(t *testing.T) {
	mockService := NewMockService()
	user := models.User{
		ID:           "test",
		Email:        "test@test.com",
		PasswordHash: "test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockService.AddUser(user)

	mockService.DeleteUser(user.ID)
	require.Empty(t, mockService.users[user.ID])
}
