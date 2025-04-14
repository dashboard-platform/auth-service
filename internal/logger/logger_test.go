package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
		env  string
	}{
		{name: "dev", env: "dev"},
		{name: "prod", env: "prod"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Running test for environment:", tt.env)
			logger := Init(tt.env)
			require.NotNil(t, logger)
		})
	}
}
