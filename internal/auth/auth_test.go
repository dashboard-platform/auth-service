package auth

import (
	"testing"

	"github.com/dashboard-platform/auth-service/internal/database"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Register(t *testing.T) {
	auth := NewService(database.NewMockupDatabase())

	tests := []struct {
		name         string
		data         models.RegisterAPI
		expectedData bool
		expectedErr  bool
	}{
		{
			name: "Register valid",
			data: models.RegisterAPI{
				Email:    "test@test.com",
				Password: "securepass",
			},
			expectedData: true,
			expectedErr:  false,
		},
		{
			name: "Register invalid email",
			data: models.RegisterAPI{
				Email:    "test",
				Password: "securepass",
			},
			expectedData: false,
			expectedErr:  true,
		},
		{
			name: "Register invalid password",
			data: models.RegisterAPI{
				Email:    "test@test.com",
				Password: "",
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
		name         string
		data         models.LoginAPI
		expectedData bool
		expectedErr  bool
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
