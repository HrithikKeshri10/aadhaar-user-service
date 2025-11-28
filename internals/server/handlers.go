package server

import (
	"aadhaar-user-service/routes"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// errHandler handles errors globally
func errHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(ErrorResponse{
		Error: err.Error(),
	})
}

// notFoundHandler handles 404 errors
var notFoundHandler = func(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Error: "Requested resource not found",
	})
}

// addRoutes registers all routes
func addRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "aadhaar-user-service",
		})
	})

	// API routes
	baseRouter := app.Group("/aadhaar")
	routes.Users(baseRouter)
}
