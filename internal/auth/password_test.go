package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyPassword_Table(t *testing.T) {
	raw := "MyStrongPass123!"
	hash, err := hashPassword(raw)
	require.NoError(t, err)

	tests := []struct {
		name      string
		hash      string
		input     string
		wantMatch bool
		expectErr bool
	}{
		{
			name:      "valid password",
			hash:      hash,
			input:     raw,
			wantMatch: true,
		},
		{
			name:      "wrong password",
			hash:      hash,
			input:     "WrongPass!",
			wantMatch: false,
		},
		{
			name:      "corrupt hash",
			hash:      "not_a_real_hash",
			input:     raw,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := verifyPassword(tt.hash, tt.input)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantMatch, match)
		})
	}
}
