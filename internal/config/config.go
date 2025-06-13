// Package config provides functionality for loading and managing application configuration.
// It retrieves configuration values from environment variables and ensures that all required
// settings are properly initialized. This package is essential for setting up the application's
// runtime environment, including database connections, JWT secrets, and server settings.
package config

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config represents the application configuration.
// It contains environment-specific settings such as the environment name,
// server port, JWT secret, and database URL.
type Config struct {
	Env       string // The current environment (e.g., "dev", "prod").
	Port      string // The port on which the server will run.
	JWTSecret []byte // The secret key used for signing JWT tokens.
	DBUrl     string // The URL for connecting to the database.
}

const (
	envKey        = "ENV"        // Environment variable key for the environment name.
	portEnv       = "PORT"       // Environment variable key for the server port.
	jwtSecretKey  = "JWT_SECRET" // Environment variable key for the JWT secret.
	dbUrlKey      = "DB_URL"     // Environment variable key for the database URL.
	defaultEnvKey = "dev"        // Default environment name if none is provided.
)

// Load retrieves the application configuration from environment variables.
// It ensures that all required configuration values are set and returns an error
// if any mandatory value is missing.
//
// Returns:
//   - Config: The loaded application configuration.
//   - error: An error if any required configuration value is missing.
func Load() (Config, error) {
	var c Config

	c.Env = os.Getenv(envKey)
	if c.Env == "" {
		c.Env = defaultEnvKey
	}

	c.Port = getEnv(portEnv, true)
	if c.Port == "" {
		return Config{}, errors.New("PORT environment variable is missing or empty")
	}

	jwtSecretStr := getEnv(jwtSecretKey, true)
	c.JWTSecret = []byte(jwtSecretStr)
	if len(c.JWTSecret) == 0 {
		return Config{}, errors.New("JWT_SECRET environment variable is missing or empty")
	}

	c.DBUrl = getEnv(dbUrlKey, true)
	if c.DBUrl == "" {
		return Config{}, errors.New("DB_URL environment variable is missing or empty")
	}
	return c, nil
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set, it logs an error and returns an empty string.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//
// Returns:
//   - string: The value of the environment variable, or an empty string if not set.
func getEnv(key string, required bool) string {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	val := os.Getenv(key)
	if val == "" && required {
		log.Error().Str("var", key).Msg("Failed to load environment")
	}
	return val
}
