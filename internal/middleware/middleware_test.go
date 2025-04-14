package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// TestRequestLogger tests the RequestLogger middleware
// It checks if the middleware correctly logs the request details and handles errors.
func TestRequestLogger_Success(t *testing.T) {
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	app := fiber.New()
	app.Use(RequestLogger(logger))
	app.Get("/", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	logOutput := logBuf.String()
	assert.True(t, strings.Contains(logOutput, `"method":"GET"`), "Expected log to contain method GET")
}

// TestRequestLogger_Error tests if the logger correctly logs an error response.
func TestRequestLogger_FiberError(t *testing.T) {
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	app := fiber.New()
	app.Use(RequestLogger(logger))
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "Bad Request")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", result["error"])
}

// TestRequestLogger_GenericError tests if the logger correctly logs a generic error response.
func TestRequestLogger_GenericError(t *testing.T) {
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	app := fiber.New()
	app.Use(RequestLogger(logger))
	app.Get("/generic", func(c *fiber.Ctx) error {
		return errors.New("custom error")
	})

	req := httptest.NewRequest("GET", "/generic", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Internal Server Error", result["error"])
}

// TestRequestLogger_WithUserID tests if the logger correctly logs the user ID from the context.
func TestRequestLogger_WithUserID(t *testing.T) {
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	app := fiber.New()
	app.Use(RequestLogger(logger))
	app.Get("/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", "12345")
		c.Status(fiber.StatusOK)
		return c.SendString("User OK")
	})

	req := httptest.NewRequest("GET", "/user", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	logOutput := logBuf.String()
	assert.True(t, strings.Contains(logOutput, `"user_id":"12345"`), "Expected log to contain user_id")
}

// FakeJWT is a fake implementation of the JWTValidator interface for testing.
// It simulates token validation: if the token is "valid-token", it returns "user123"; otherwise, it returns an error.
type FakeJWT struct{}

// ValidateJWT simulates token validation.
func (fj *FakeJWT) ValidateJWT(token string) (string, error) {
	if token == "valid-token" {
		return "user123", nil
	}
	return "", errors.New("invalid or expired token")
}

// parseJSONBody is a helper function to unmarshal a JSON response body into a map.
func parseJSONBody(body []byte) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal(body, &result)
	return result, err
}

// TestRequireAuth_NoToken tests the case when no token is provided (neither cookie nor header).
func TestRequireAuth_NoToken(t *testing.T) {
	app := fiber.New()

	fakeJWT := &FakeJWT{}
	app.Use(RequireAuth(fakeJWT))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	bodyBytes := make([]byte, resp.ContentLength)
	resp.Body.Read(bodyBytes)
	result, err := parseJSONBody(bodyBytes)
	assert.NoError(t, err)
	assert.Equal(t, "authentication required", result["error"])
}

// TestRequireAuth_InvalidToken tests the scenario where an invalid token is provided via cookie.
func TestRequireAuth_InvalidToken(t *testing.T) {
	app := fiber.New()

	fakeJWT := &FakeJWT{}
	app.Use(RequireAuth(fakeJWT))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Cookie", "access_token=invalid-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	bodyBytes := make([]byte, resp.ContentLength)
	resp.Body.Read(bodyBytes)
	result, err := parseJSONBody(bodyBytes)
	assert.NoError(t, err)
	assert.Equal(t, "invalid or expired token", result["error"])
}

// TestRequireAuth_ValidTokenFromCookie tests a valid token provided via cookie.
func TestRequireAuth_ValidTokenFromCookie(t *testing.T) {
	app := fiber.New()

	fakeJWT := &FakeJWT{}
	app.Use(RequireAuth(fakeJWT))
	app.Get("/", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		return c.SendString(uid.(string))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Cookie", "access_token=valid-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	buf := make([]byte, resp.ContentLength)
	resp.Body.Read(buf)
	responseText := string(buf)
	assert.Equal(t, "user123", responseText)
}

// TestRequireAuth_ValidTokenFromHeader tests a valid token provided via Authorization header.
func TestRequireAuth_ValidTokenFromHeader(t *testing.T) {
	app := fiber.New()

	fakeJWT := &FakeJWT{}
	app.Use(RequireAuth(fakeJWT))
	app.Get("/", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		return c.SendString(uid.(string))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	buf := make([]byte, resp.ContentLength)
	resp.Body.Read(buf)
	responseText := string(buf)
	assert.Equal(t, "user123", responseText)
}

// TestRequireAuth_CookiePrecedence tests that when both cookie and header are provided,
// the token from the cookie is used (even if invalid), and the header is ignored.
func TestRequireAuth_CookiePrecedence(t *testing.T) {
	app := fiber.New()

	fakeJWT := &FakeJWT{}
	app.Use(RequireAuth(fakeJWT))
	app.Get("/", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		return c.SendString(uid.(string))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Cookie", "access_token=invalid-token")
	req.Header.Set("Authorization", "Bearer valid-token")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	buf := make([]byte, resp.ContentLength)
	resp.Body.Read(buf)
	result, err := parseJSONBody(buf)
	assert.NoError(t, err)
	assert.Equal(t, "invalid or expired token", result["error"])
}
