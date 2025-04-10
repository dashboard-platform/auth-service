package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func Init(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	if env == "dev" {
		return zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewComponentLogger(logger zerolog.Logger, name string) zerolog.Logger {
	return logger.With().Str("component", name).Logger()
}
