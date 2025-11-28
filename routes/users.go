package routes

import (
	"aadhaar-user-service/controllers/users"

	"github.com/gofiber/fiber/v2"
)

// Users registers user routes
func Users(r fiber.Router) {
	u := r.Group("/users")

	u.Post("/", users.Add)      // Create a new user
	u.Get("/", users.GetAll)    // List users with pagination and sorting
	u.Get("/:id", users.Get)    // Get user by ID
	u.Delete("/:id", users.Delete) // Delete user by ID
}
