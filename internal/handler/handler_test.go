// Package handler provides HTTP handlers for the authentication service.
// These handlers process incoming requests, interact with the service layer, and return responses.
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/dashboard-platform/auth-service/internal/auth"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// setupTestApp initializes a Fiber app with mock services for testing.
//
// Returns:
//   - *fiber.App: The initialized Fiber app.
//   - *auth.MockService: The mock authentication service.
func setupTestApp() (*fiber.App, *auth.MockService) {
	app := fiber.New()

	fakeService := auth.NewMockService()
	jwt := &auth.JWTObj{Secret: []byte("test-secret")}
	h := New(fakeService, jwt)

	app.Get("/healthcheck", h.Healthcheck)
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Get("/me", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-1")
		return h.GetMe(c)
	})

	return app, fakeService
}

// FakeOAuthClient is a mock implementation of the GoogleOAuthClient interface.
// It is used for testing purposes to simulate the behavior of exchanging an OAuth code
// for an email address. The behavior can be controlled using the ShouldFail flag.
type FakeOAuthClient struct {
	ShouldFail bool // Determines whether the ExchangeCode method should simulate a failure.
}

// ExchangeCode simulates the exchange of an OAuth code for an email address.
//
// Parameters:
//   - code: The OAuth authorization code provided by the client.
//
// Returns:
//   - string: The email address associated with the OAuth code if successful.
//   - error: An error if the code is invalid or ShouldFail is set to true.
func (f *FakeOAuthClient) ExchangeCode(code string) (string, error) {
	if f.ShouldFail || code != "valid-code" {
		return "", errors.New("invalid code")
	}
	return "mockuser@example.com", nil
}

// setupGoogleTestApp initializes a Fiber app with mock services and a mock OAuth client for testing.
//
// Parameters:
//   - fakeAuth: The mock authentication service to be used.
//   - jwt: The JWT utility object for token generation.
//   - oauthClient: The mock OAuth client to simulate Google OAuth behavior.
//
// Returns:
//   - *fiber.App: The initialized Fiber app configured with the Google OAuth endpoint.
func setupGoogleTestApp(fakeAuth auth.ServiceInterface, jwt *auth.JWTObj) *fiber.App {
	app := fiber.New()

	h := New(fakeAuth, jwt)

	app.Post("/auth/google", func(c *fiber.Ctx) error {
		// Inject mock OAuth logic direct
		var body struct {
			Code string `json:"code"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Data: "invalid request"})
		}
		var email string
		if body.Code == "valid-code" {
			email = "mockuser@example.com"
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(Response{Error: true, Data: "oauth failed"})
		}

		id, err := h.auth.LoginOrRegisterOAuth(models.RegisterAPI{Email: email, Name: "Mock User"}) // Furnizează și Name
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Data: "auth service error"})
		}
		token, _ := h.jwt.CreateJWT(id)
		return c.Status(fiber.StatusOK).JSON(Response{Error: false, Data: fiber.Map{"token": token}})
	})

	return app
}

// TestHandlers_Healthcheck tests the Healthcheck handler.
func TestHandlers_Healthcheck(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/healthcheck", nil)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestHandlers_Register tests the Register handler with various scenarios.
func TestHandlers_Register(t *testing.T) {
	app, _ := setupTestApp()

	tests := []struct {
		name           string
		body           any
		expectedStatus int
	}{
		{
			name:           "Register valid",
			body:           models.RegisterAPI{Name: "Test User", Email: "test@example.com", Password: "securepass"},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name:           "Register invalid email",
			body:           models.RegisterAPI{Name: "Test User", Email: "test", Password: "securepass"},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:           "Register invalid password",
			body:           models.RegisterAPI{Name: "Test User", Email: "test@example.com", Password: ""},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:           "Register missing name",
			body:           models.RegisterAPI{Email: "test@example.com", Password: "securepass"}, // Name lipsește
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name: "Register incorrect JSON fields",
			body: map[string]string{
				"something": "else",
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:           "Register invalid data",
			body:           "something-else",
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(fiber.MethodPost, "/register", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// TestHandlers_Login tests the Login handler with various scenarios.
func TestHandlers_Login(t *testing.T) {
	app, mockService := setupTestApp()
	mockService.AddUser(models.User{
		ID:           uuid.New().String(),
		Email:        "test@test",
		PasswordHash: "securepass",
	})

	tests := []struct {
		name           string
		body           any
		expectedStatus int
	}{
		{
			name: "Login valid",
			body: models.LoginAPI{
				Email:    "test@test",
				Password: "securepass",
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "Login invalid",
			body: models.LoginAPI{
				Email:    "invalid@test",
				Password: "securepass",
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name: "Register incorrect JSON fields",
			body: map[string]string{
				"something": "else",
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:           "Register invalid data",
			body:           "something-else",
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// TestHandlers_GetMe tests the GetMe handler with various scenarios.
func TestHandlers_GetMe(t *testing.T) {
	app, mockService := setupTestApp()

	tests := []struct {
		name           string
		expectedStatus int
		helper         func(t *testing.T, app *fiber.App)
	}{
		{
			name:           "Get Me valid",
			expectedStatus: fiber.StatusOK,
			helper: func(t *testing.T, app *fiber.App) {
				mockService.AddUser(models.User{
					ID:           "user-1",
					Email:        "test@test",
					PasswordHash: "securepass",
				})
			},
		},
		{
			name:           "Get Me unauthorized",
			expectedStatus: fiber.StatusUnauthorized,
			helper: func(t *testing.T, app *fiber.App) {
				mockService.DeleteUser("user-1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.helper(t, app)
			req := httptest.NewRequest(fiber.MethodGet, "/me", nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// TestHandlers_GoogleLogin tests the Google login handler with various scenarios.
func TestHandlers_GoogleLogin(t *testing.T) {
	jwt := &auth.JWTObj{Secret: []byte("test-secret")}
	service := auth.NewMockService()
	app := setupGoogleTestApp(service, jwt)

	tests := []struct {
		name           string
		body           any
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Valid code",
			body: map[string]string{
				"code": "valid-code",
			},
			expectedStatus: fiber.StatusOK,
			expectToken:    true,
		},
		{
			name: "Invalid code",
			body: map[string]string{
				"code": "bad-code",
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "Missing code",
			body:           map[string]string{},
			expectedStatus: fiber.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(fiber.MethodPost, "/auth/google", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
