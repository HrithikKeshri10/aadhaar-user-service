package users

import (
	"context"
	"errors"
	"math"

	"aadhaar-user-service/internals/dto"
	"aadhaar-user-service/models/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrAadhaarIDExists = errors.New("aadhaar application id already exists")
	ErrInvalidUUID     = errors.New("invalid uuid format")
)

// UserService handles user business logic
type UserService struct {
	User  *dto.User
	Users *dto.Users
}

// New creates a new UserService instance
func New() *UserService {
	return &UserService{}
}

// Create creates a new user after validation
func (s *UserService) Create(ctx context.Context, input dto.UserCreate) error {
	// Check if email already exists
	existingUser := users.New()
	if err := existingUser.GetByEmail(ctx, input.Email); err == nil {
		return ErrEmailExists
	}

	// Check if Aadhaar Application ID already exists
	if err := existingUser.GetByAadhaarApplicationID(ctx, input.AadhaarApplicationID); err == nil {
		return ErrAadhaarIDExists
	}

	// Create new user
	user := users.New()
	user.AadhaarApplicationID = input.AadhaarApplicationID
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone
	user.Address = input.Address
	user.DateOfBirth = input.DateOfBirth
	user.Gender = input.Gender

	if err := user.Create(ctx); err != nil {
		return err
	}

	// Map to DTO
	s.User = &dto.User{
		ID:                   user.ID,
		AadhaarApplicationID: user.AadhaarApplicationID,
		Name:                 user.Name,
		Email:                user.Email,
		Phone:                user.Phone,
		Address:              user.Address,
		DateOfBirth:          user.DateOfBirth,
		Gender:               user.Gender,
		CreatedAt:            &user.CreatedAt,
	}

	return nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) error {
	user := users.New()

	// Parse UUID
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUID
	}
	user.ID = parsedID

	if err := user.GetByID(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return err
	}

	// Map to DTO
	s.User = &dto.User{
		ID:                   user.ID,
		AadhaarApplicationID: user.AadhaarApplicationID,
		Name:                 user.Name,
		Email:                user.Email,
		Phone:                user.Phone,
		Address:              user.Address,
		DateOfBirth:          user.DateOfBirth,
		Gender:               user.Gender,
		CreatedAt:            &user.CreatedAt,
		UpdatedAt:            &user.UpdatedAt,
	}

	return nil
}

// GetAllPaginated retrieves users with pagination
func (s *UserService) GetAllPaginated(ctx context.Context, params dto.PaginationParams) error {
	user := users.New()

	userList, total, err := user.GetAllPaginated(ctx, params)
	if err != nil {
		return err
	}

	// Map to DTOs
	userDTOs := make([]dto.User, len(userList))
	for i, u := range userList {
		userDTOs[i] = dto.User{
			ID:                   u.ID,
			AadhaarApplicationID: u.AadhaarApplicationID,
			Name:                 u.Name,
			Email:                u.Email,
			Phone:                u.Phone,
			Address:              u.Address,
			DateOfBirth:          u.DateOfBirth,
			Gender:               u.Gender,
			CreatedAt:            &u.CreatedAt,
			UpdatedAt:            &u.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	s.Users = &dto.Users{
		Users:      userDTOs,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return nil
}

// Delete removes a user by ID
func (s *UserService) Delete(ctx context.Context, id string) error {
	user := users.New()

	// Parse UUID
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUID
	}
	user.ID = parsedID

	// Check if user exists
	if err := user.GetByID(ctx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return err
	}

	return user.Delete(ctx)
}
