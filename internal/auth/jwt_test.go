package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWT_Validate_Table(t *testing.T) {
	j := &JWTObj{Secret: []byte("super-secret")}
	validToken, _ := j.CreateJWT("user123")

	wrongSigner := &JWTObj{Secret: []byte("wrong")}
	invalidToken, _ := wrongSigner.CreateJWT("user456")

	expiredClaims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredSigned, _ := expiredToken.SignedString(j.Secret)

	tests := []struct {
		name      string
		token     string
		wantID    string
		expectErr bool
	}{
		{
			name:   "valid token",
			token:  validToken,
			wantID: "user123",
		},
		{
			name:      "wrong signer",
			token:     invalidToken,
			expectErr: true,
		},
		{
			name:      "expired token",
			token:     expiredSigned,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := j.ValidateJWT(tt.token)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantID, id)
			}
		})
	}
}
