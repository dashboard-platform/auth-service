package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoogleOAuthFlow(t *testing.T) {
	const baseURL = "http://localhost:8081"

	// Google OAuth code for testing. This code should be replaced
	// with a valid OAuth code obtained from the Google OAuth flow.
	code := "4/0Ab_5qlndT2vww1kwSiwNM-fuV72Sxi0gmeNY-QiwtDJu8U9TSUJ7dPZt8tHUJP8bIYuNSw"

	t.Run("Real Google OAuth", func(t *testing.T) {
		body := map[string]string{
			"code": code,
		}
		payload, _ := json.Marshal(body)

		resp, err := http.Post(baseURL+"/auth/google", "application/json", bytes.NewBuffer(payload))
		require.NoError(t, err)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("auth/google failed: status %d", resp.StatusCode)
		}

		var result struct {
			Error bool                   `json:"error"`
			Data  map[string]interface{} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		require.False(t, result.Error)
		require.Contains(t, result.Data, "token")
		require.Contains(t, result.Data, "email")

		token := result.Data["token"].(string)
		req, err := http.NewRequest("GET", baseURL+"/auth/me", nil)
		require.NoError(t, err)

		// Set Authorization header with the token.
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("auth/me failed: status %d", resp.StatusCode)
		}
	})
}
