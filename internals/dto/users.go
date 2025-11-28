package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserCreate represents the request body for creating a new user
type UserCreate struct {
	AadhaarApplicationID string `json:"aadhaar_application_id" validate:"required,len=14"`
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Phone                string `json:"phone" validate:"required,len=10,numeric"`
	Address              string `json:"address" validate:"required,max=500"`
	DateOfBirth          string `json:"date_of_birth" validate:"required"`
	Gender               string `json:"gender" validate:"required,oneof=male female other"`
}

// User represents a user response
type User struct {
	ID                   uuid.UUID  `json:"id"`
	AadhaarApplicationID string     `json:"aadhaar_application_id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	Phone                string     `json:"phone"`
	Address              string     `json:"address"`
	DateOfBirth          string     `json:"date_of_birth"`
	Gender               string     `json:"gender"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

// Users represents a collection of users with pagination metadata
type Users struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

// PaginationParams represents pagination and sorting parameters
type PaginationParams struct {
	Page    int    `query:"page" validate:"min=1"`
	Limit   int    `query:"limit" validate:"min=1,max=100"`
	SortBy  string `query:"sort_by" validate:"omitempty,oneof=name email created_at aadhaar_application_id"`
	Order   string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search  string `query:"search"`
}

// DefaultPaginationParams returns default pagination values
func DefaultPaginationParams() PaginationParams {
	return PaginationParams{
		Page:   1,
		Limit:  10,
		SortBy: "created_at",
		Order:  "desc",
	}
}
