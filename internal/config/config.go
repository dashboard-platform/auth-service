package config

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env       string
	Port      string
	JWTSecret []byte
	DBUrl     string
}

const (
	envKey        = "ENV"
	portEnv       = "PORT"
	jwtSecretKey  = "JWT_SECRET"
	dbUrlKey      = "DB_URL"
	defaultEnvKey = "dev"
)

func Load() (Config, error) {
	var c Config

	c.Env = os.Getenv(envKey)
	if c.Env == "" {
		c.Env = defaultEnvKey
	}

	c.Port = getEnv(portEnv)
	if c.Port == "" {
		return Config{}, errors.New("empty key")
	}

	c.JWTSecret = []byte(getEnv(jwtSecretKey))
	if len(c.JWTSecret) == 0 {
		return Config{}, errors.New("empty key")
	}

	c.DBUrl = getEnv(dbUrlKey)
	if c.DBUrl == "" {
		return Config{}, errors.New("empty key")
	}

	return c, nil
}

func getEnv(key string) string {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	val := os.Getenv(key)
	if val == "" {
		log.Error().Str("var", key).Msg("Failed to load environment")
	}
	return val
}
