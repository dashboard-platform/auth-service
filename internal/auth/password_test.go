// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestVerifyPassword_Table tests the verifyPassword function using a table-driven approach.
// It verifies that valid passwords match their hashes, invalid passwords do not match,
// and corrupt hashes return an error.
func TestVerifyPassword_Table(t *testing.T) {
	raw := "MyStrongPass123!" // The raw password to be hashed and tested.
	hash, err := hashPassword(raw)
	require.NoError(t, err)

	tests := []struct {
		name      string // Name of the test case.
		hash      string // The hash to verify against.
		input     string // The input password to verify.
		wantMatch bool   // Expected result: true if the password matches the hash.
		expectErr bool   // Expected error: true if an error is expected.
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
