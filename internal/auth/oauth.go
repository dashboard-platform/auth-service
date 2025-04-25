// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GoogleOAuthAPI handles Google OAuth token exchange and ID token parsing.
type GoogleOAuthAPI struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type googleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func (g *GoogleOAuthAPI) ExchangeCode(code string) (googleClaims, error) {
	values := make(url.Values)
	values.Set("client_id", g.ClientID)
	values.Set("client_secret", g.ClientSecret)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", g.RedirectURL)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", values)
	if err != nil {
		return googleClaims{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return googleClaims{}, err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		Email       string `json:"email"`
	}
	err = json.Unmarshal(bodyBytes, &tokenResp)
	if err != nil {
		return googleClaims{}, err
	}

	parts := strings.Split(tokenResp.IDToken, ".")
	if len(parts) != 3 {
		return googleClaims{}, fmt.Errorf("invalid ID token format")
	}

	payload := parts[1]
	decoded, err := decodeBase64URL(payload)
	if err != nil {
		return googleClaims{}, err
	}

	var claims googleClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return googleClaims{}, err
	}

	if !claims.EmailVerified {
		return googleClaims{}, fmt.Errorf("email not verified")
	}

	return claims, nil
}

func decodeBase64URL(data string) ([]byte, error) {
	pad := len(data) % 4
	if pad > 0 {
		data += strings.Repeat("=", 4-pad)
	}
	ret, err := url.QueryUnescape(data)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(ret)
}
