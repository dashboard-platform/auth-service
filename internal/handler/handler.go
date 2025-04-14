// Package handler provides HTTP handlers for the authentication service.
// These handlers process incoming requests, interact with the service layer, and return responses.
package handler

import (
	"github.com/dashboard-platform/auth-service/internal/auth"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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

	ctx.Locals("user_id", id)

	return ctx.Status(fiber.StatusOK).JSON(Response{
		Error: false,
		Data: fiber.Map{
			"token": token,
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
