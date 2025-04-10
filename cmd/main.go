package main

import (
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	c, err := config.Load()
	if err != nil {
		return
	}

	baseLogger := logger.Init(c.Env)
	db, err := database.Init(c.DBUrl, baseLogger)
	if err != nil {
		return
	}
	defer db.Close()

	if err = db.AutoMigrate(); err != nil {
		return
	}

	log.Info().Msg("Starting Authentication Service...")

	httpLogger := logger.NewComponentLogger(baseLogger, "http")

	app := fiber.New()
	app.Use(
		// Set CORS correct configuration.
		cors.New(cors.Config{
			// AllowOrigins:     "https://API-GATEWAY:1234",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			AllowMethods: "GET, POST, PUT, DELETE",
		}),

		// Add security headers.
		helmet.New(),

		// Implement rate limiting.
		limiter.New(limiter.Config{
			Max:        20,
			Expiration: 1 * time.Minute,
		}),

		// Add custom logger
		middleware.RequestLogger(httpLogger),
	)

	jwtObj := auth.JWTObj{Secret: c.JWTSecret}
	h := handler.New(auth.NewService(db), &jwtObj)

	app.Get("/healthcheck", h.Healthcheck)
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Get("/me", middleware.RequireAuth(&jwtObj), h.GetMe)

	log.Info().Msgf("Authentication Service started on %s", c.Port)
	if err = app.Listen(c.Port); err != nil {
		log.Error().Msgf("Error starting authentication service: %v", err)
		return
	}
}
