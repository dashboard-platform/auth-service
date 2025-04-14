package auth

import (
	"auth-service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewMockService(t *testing.T) {
	mockService := NewMockService()
	require.NotNil(t, mockService)
	require.NotNil(t, mockService.users)
	require.NotNil(t, mockService.logins)
}

func TestMockRegister(t *testing.T) {
	mockService := NewMockService()

	tests := []struct {
		name          string
		input         models.RegisterAPI
		expectedError bool
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
		name          string
		input         models.LoginAPI
		expectedError bool
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
		name          string
		input         string
		expectedError bool
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
