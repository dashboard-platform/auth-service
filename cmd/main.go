// Package main is the entry point for the Authentication Service.
// It initializes the configuration, database, middleware, and HTTP routes.
package main

import (
	"time"

	"github.com/dashboard-platform/auth-service/internal/auth"
	"github.com/dashboard-platform/auth-service/internal/config"
	"github.com/dashboard-platform/auth-service/internal/database"
	"github.com/dashboard-platform/auth-service/internal/handler"
	"github.com/dashboard-platform/auth-service/internal/logger"
	"github.com/dashboard-platform/auth-service/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

// main is the entry point of the Authentication Service.
// It loads the configuration, initializes the database, sets up middleware,
// and starts the HTTP server.
func main() {
	// Load configuration from environment variables or .env file.
	c, err := config.Load()
	if err != nil {
		return
	}

	// Initialize the base logger for the application.
	baseLogger := logger.Init(c.Env)

	// Initialize the database connection.
	db, err := database.Init(c.DBUrl, baseLogger)
	if err != nil {
		return
	}
	defer db.Close()

	// Perform database migrations.
	if err = db.AutoMigrate(); err != nil {
		return
	}

	log.Info().Msg("Starting Authentication Service...")

	// Create a logger for HTTP requests.
	httpLogger := logger.NewComponentLogger(baseLogger, "http")

	// Initialize the Fiber app with middleware.
	app := fiber.New()
	app.Use(
		// Set CORS configuration.
		cors.New(cors.Config{
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

		// Add custom request logger middleware.
		middleware.RequestLogger(httpLogger),
	)

	// Initialize the JWT object and HTTP handlers.
	jwtObj := auth.JWTObj{Secret: c.JWTSecret}
	h := handler.New(auth.NewService(db), &jwtObj)

	// Define HTTP routes.
	app.Get("/healthcheck", h.Healthcheck)
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Get("/me", middleware.RequireAuth(&jwtObj), h.GetMe)

	// Start the HTTP server.
	log.Info().Msgf("Authentication Service started on %s", c.Port)
	if err = app.Listen(c.Port); err != nil {
		log.Error().Msgf("Error starting authentication service: %v", err)
		return
	}
}
