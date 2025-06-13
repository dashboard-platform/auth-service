package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegisterLoginMeFlow(t *testing.T) {
	const baseURL = "http://localhost:8081"

	email := fmt.Sprintf("testuser+%d@example.com", time.Now().UnixNano())
	password := "securepassword123"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	t.Run("Register User", func(t *testing.T) {
		body := map[string]string{
			"name":     "Test User E2E",
			"email":    email,
			"password": password,
		}
		payload, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, baseURL+"/auth/register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)
		require.NotEmpty(t, req)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result struct {
			Error bool   `json:"error"`
			Data  string `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.False(t, result.Error)
		require.Equal(t, email, result.Data)
	})

	var token string

	t.Run("Login User", func(t *testing.T) {
		body := map[string]string{
			"email":    email,
			"password": password,
		}
		payload, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, baseURL+"/auth/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)
		require.NotEmpty(t, req)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Error bool `json:"error"`
			Data  struct {
				Token string `json:"token"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.False(t, result.Error)
		require.NotEmpty(t, result.Data.Token)

		token = result.Data.Token
	})

	t.Run("Get Me", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Error bool                   `json:"error"`
			Data  map[string]interface{} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.False(t, result.Error)
		require.Equal(t, email, result.Data["email"])
	})
}
