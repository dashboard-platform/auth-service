package handler

import (
	"github.com/dashboard-platform/auth-service/internal/auth"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type HTTPHandler struct {
	auth auth.ServiceInterface
	jwt  *auth.JWTObj
}

func New(authService auth.ServiceInterface, jwtObj *auth.JWTObj) HTTPHandler {
	return HTTPHandler{
		auth: authService,
		jwt:  jwtObj,
	}
}

func (h *HTTPHandler) Healthcheck(ctx *fiber.Ctx) error {
	log.Info().Msg("Healthcheck called")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "auth-service is alive",
	})
}

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

	if id == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"data":  "invalid login credentials",
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

/* TODO: move this to api-gateway.
func (h *HTTPHandler) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  "logged out successfully",
	})
}
*/

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
