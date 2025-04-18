// Package handler provides HTTP handlers for the authentication service.
// These handlers process incoming requests, interact with the service layer, and return responses.
package handler

import (
	"net/http"
	"os"

	"github.com/dashboard-platform/auth-service/internal/auth"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// HTTPHandler represents the HTTP handlers for the authentication service.
// It includes methods for health checks, user registration, login, and retrieving user details.
type HTTPHandler struct {
	auth auth.ServiceInterface // The authentication service interface.
	jwt  *auth.JWTObj          // The JWT utility object.
}

// New creates a new instance of HTTPHandler.
//
// Parameters:
//   - authService: The authentication service implementation.
//   - jwtObj: The JWT utility object.
//
// Returns:
//   - HTTPHandler: A new instance of the HTTPHandler.
func New(authService auth.ServiceInterface, jwtObj *auth.JWTObj) HTTPHandler {
	return HTTPHandler{
		auth: authService,
		jwt:  jwtObj,
	}
}

// Healthcheck handles the health check endpoint.
//
// Returns:
//   - fiber.StatusOK: If the service is running.
func (h *HTTPHandler) Healthcheck(ctx *fiber.Ctx) error {
	log.Info().Msg("Healthcheck called")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "auth-service is alive",
	})
}

// Register handles user registration requests.
//
// Parameters:
//   - ctx: The Fiber context containing the request data.
//
// Returns:
//   - fiber.StatusCreated: If the user is successfully registered.
//   - fiber.StatusBadRequest: If the request data is invalid.
//   - fiber.StatusInternalServerError: If an internal error occurs.
func (h *HTTPHandler) Register(ctx *fiber.Ctx) error {
	var data models.RegisterAPI

	if err := ctx.BodyParser(&data); err != nil {
		log.Error().Err(err).Msg("error reading/parsing HTTP request body data")
		return ctx.Status(fiber.StatusBadRequest).JSON(Response{
			Error: true,
			Data:  "invalid data provided",
		})
	}

	userID, err := h.auth.Register(data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	ctx.Locals("user_id", userID)

	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Error: false,
		Data:  data.Email,
	})
}

// Login handles user login requests.
//
// Parameters:
//   - ctx: The Fiber context containing the request data.
//
// Returns:
//   - fiber.StatusOK: If the user is successfully authenticated.
//   - fiber.StatusBadRequest: If the request data is invalid.
//   - fiber.StatusInternalServerError: If an internal error occurs.
func (h *HTTPHandler) Login(ctx *fiber.Ctx) error {
	var data models.LoginAPI

	if err := ctx.BodyParser(&data); err != nil {
		log.Error().Err(err).Msg("error reading/parsing HTTP request body data")
		return ctx.Status(fiber.StatusBadRequest).JSON(Response{
			Error: true,
			Data:  "invalid data provided",
		})
	}

	id, err := h.auth.Login(data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "invalid login credentials",
		})
	}

	token, err := h.jwt.CreateJWT(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") != "dev", // use Secure only in production
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   60 * 60 * 1, // 1 hour
	})
	ctx.Locals("user_id", id)

	return ctx.Status(fiber.StatusOK).JSON(Response{
		Error: false,
		Data: fiber.Map{
			"token": token,
		},
	})
}

func (h *HTTPHandler) GoogleRedirect(ctx *fiber.Ctx) error {
	// Încarcă din env
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL") // ex: http://localhost:5173/auth/google/callback

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return ctx.Status(fiber.StatusInternalServerError).SendString("OAuth2 Google misconfigured")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	state := "some-random-state"
	authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Redirect către Google
	return ctx.Status(http.StatusTemporaryRedirect).Redirect(authURL)
}

// GoogleLogin handles user login via Google OAuth.
//
// This method processes the OAuth code provided by the client, exchanges it for
// an access token, retrieves the user's email, and either logs in or registers
// the user in the system. A JWT token is then generated and returned to the client.
//
// Parameters:
//   - ctx: The Fiber context containing the request data.
//
// Returns:
//   - fiber.StatusOK: If the user is successfully authenticated or registered.
//   - fiber.StatusBadRequest: If the request data is invalid.
//   - fiber.StatusInternalServerError: If an internal error occurs.
func (h *HTTPHandler) GoogleLogin(ctx *fiber.Ctx) error {
	var body struct {
		Code string `json:"code"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("error reading/parsing HTTP request body data")
		return ctx.Status(fiber.StatusBadRequest).JSON(Response{
			Error: true,
			Data:  "invalid data provided",
		})
	}

	client := auth.GoogleOAuthAPI{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}

	if client.ClientID == "" || client.ClientSecret == "" || client.RedirectURL == "" {
		log.Error().Msg("missing Google OAuth credentials in environment variables")
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	if body.Code == "" {
		log.Error().Msg("missing OAuth code in request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(Response{
			Error: true,
			Data:  "missing code in request body",
		})
	}

	email, err := client.ExchangeCode(body.Code)
	if err != nil {
		log.Error().Err(err).Msg("error exchanging code for token")
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	userID, err := h.auth.LoginOrRegisterOAuth(email)
	if err != nil {
		log.Error().Err(err).Msg("error logging in or registering user")
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	token, err := h.jwt.CreateJWT(userID)
	if err != nil {
		log.Error().Err(err).Msg("error creating JWT token")
		return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
			Error: true,
			Data:  "unexpected internal error",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") != "dev", // use Secure only in production
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   60 * 60 * 1, // 1 hour
	})
	ctx.Locals("user_id", userID)

	return ctx.Status(fiber.StatusOK).JSON(Response{
		Error: false,
		Data: fiber.Map{
			"token": token,
			"email": email,
		},
	})
}

// GetMe retrieves the details of the authenticated user.
//
// Parameters:
//   - ctx: The Fiber context containing the user ID.
//
// Returns:
//   - fiber.StatusOK: If the user details are successfully retrieved.
//   - fiber.StatusUnauthorized: If the user is not authenticated.
//   - fiber.StatusInternalServerError: If an internal error occurs.
func (h *HTTPHandler) GetMe(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"data":  "unauthorized",
		})
	}

	user, err := h.auth.GetUserByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"data":  "unexpected internal error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(Response{
		Error: false,
		Data: fiber.Map{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
