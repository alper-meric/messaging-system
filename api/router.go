package api

import (
	"log"

	"github.com/alper.meric/messaging-system/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, controller *handlers.MessageController) {
	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// API Endpoints - Using controller methods
	api := app.Group("/api")
	api.Post("/service", controller.ServiceControl)
	api.Get("/service/status", controller.ServiceStatus)
	api.Get("/messages", controller.GetSentMessages)

	// Swagger UI - serve static files
	app.Static("/swagger", "./docs/swagger-ui")

	// Serve swagger.yaml
	app.Get("/docs/swagger.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.yaml")
	})

	// Home page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./views/index.html", false)
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("404 Not Found: %s %s", c.Method(), c.Path())
		return c.Status(fiber.StatusNotFound).SendString("Page not found")
	})
}

// ErrorHandler handles errors returned from routes
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error
	log.Printf("Error: %v", err)

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
