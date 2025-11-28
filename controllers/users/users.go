package users

import (
	"aadhaar-user-service/internals/dto"
	"aadhaar-user-service/internals/validator"
	"aadhaar-user-service/services/users"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                      `json:"error"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

// Add creates a new user
func Add(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var input dto.UserCreate

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Validate input
	if validationErrors := validator.Payload(input); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
	}

	// Create user via service
	svc := users.New()
	if err := svc.Create(ctx, input); err != nil {
		switch err {
		case users.ErrEmailExists:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "Email already exists",
			})
		case users.ErrAadhaarIDExists:
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "Aadhaar application ID already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to create user",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(svc.User)
}

// Get retrieves a user by ID
func Get(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "User ID is required",
		})
	}

	svc := users.New()
	if err := svc.GetByID(ctx, id); err != nil {
		switch err {
		case users.ErrInvalidUUID:
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Invalid user ID format",
			})
		case users.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "User not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to retrieve user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(svc.User)
}

// GetAll retrieves all users with pagination and sorting
func GetAll(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Get default params
	params := dto.DefaultPaginationParams()

	// Parse query parameters
	if page := c.QueryInt("page", 1); page > 0 {
		params.Page = page
	}
	if limit := c.QueryInt("limit", 10); limit > 0 {
		params.Limit = limit
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if order := c.Query("order"); order != "" {
		params.Order = order
	}
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// Validate and normalize pagination
	params.Page, params.Limit = validator.ValidatePagination(params.Page, params.Limit)

	// Validate sort parameters
	validSortFields := map[string]bool{
		"name": true, "email": true, "created_at": true, "aadhaar_application_id": true,
	}
	if !validSortFields[params.SortBy] {
		params.SortBy = "created_at"
	}

	validOrders := map[string]bool{"asc": true, "desc": true}
	if !validOrders[params.Order] {
		params.Order = "desc"
	}

	svc := users.New()
	if err := svc.GetAllPaginated(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to retrieve users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(svc.Users)
}

// Delete removes a user by ID
func Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "User ID is required",
		})
	}

	svc := users.New()
	if err := svc.Delete(ctx, id); err != nil {
		switch err {
		case users.ErrInvalidUUID:
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Invalid user ID format",
			})
		case users.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "User not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to delete user",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
