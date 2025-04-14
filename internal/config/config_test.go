// Package config provides functionality for loading and managing application configuration.
// It retrieves configuration values from environment variables and ensures that all required
// settings are properly initialized. This package is essential for setting up the application's
// runtime environment, including database connections, JWT secrets, and server settings.
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetEnv tests the getEnv function to ensure it retrieves the correct environment variable value.
// It verifies that the function returns the expected value for existing environment variables
// and an empty string for non-existent variables.
func TestGetEnv(t *testing.T) {
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	tests := []struct {
		name string // Name of the test case.
		key  string // The environment variable key to retrieve.
		want string // The expected value of the environment variable.
	}{
		{
			name: "Test ENV variable",
			key:  "test",
			want: "test",
		},
		{
			name: "Test non-existent ENV variable",
			key:  "non_existent",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEnv(tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestLoad tests the Load function to ensure it correctly loads the configuration from environment variables.
// It verifies that the function returns a valid configuration when all required environment variables are set
// and returns an error when any mandatory variable is missing.
func TestLoad(t *testing.T) {
	envs := []string{envKey, portEnv, jwtSecretKey, dbUrlKey, "test"}
	for _, env := range envs {
		os.Setenv(env, "test")
		defer os.Unsetenv(env)
	}

	tests := []struct {
		name    string   // Name of the test case.
		envs    []string // List of environment variables to set for the test.
		wantErr bool     // Expected error: true if an error is expected.
	}{
		{
			name:    "Test Load with all envs set",
			envs:    []string{envKey, portEnv, jwtSecretKey, dbUrlKey},
			wantErr: false,
		},
		{
			name:    "Test Load with missing envKey",
			envs:    []string{portEnv, jwtSecretKey, dbUrlKey},
			wantErr: false,
		},
		{
			name:    "Test Load with missing portEnv",
			envs:    []string{envKey, jwtSecretKey, dbUrlKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing jwtSecretKey",
			envs:    []string{envKey, portEnv, dbUrlKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing dbUrlKey",
			envs:    []string{envKey, portEnv, jwtSecretKey},
			wantErr: true,
		},
		{
			name:    "Test Load with empty envs",
			envs:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, env := range tt.envs {
				os.Setenv(env, "test")
				defer os.Unsetenv(env)
			}

			cfg, err := Load()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, cfg)
		})
	}
}
