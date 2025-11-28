package users

import (
	"context"
	"fmt"
	"time"

	"aadhaar-user-service/internals/database"
	"aadhaar-user-service/internals/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the database model for users table
type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AadhaarApplicationID string    `gorm:"uniqueIndex;size:14;not null" json:"aadhaar_application_id"`
	Name                 string    `gorm:"size:100;not null" json:"name"`
	Email                string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Phone                string    `gorm:"size:10;not null" json:"phone"`
	Address              string    `gorm:"size:500;not null" json:"address"`
	DateOfBirth          string    `gorm:"size:10;not null" json:"date_of_birth"`
	Gender               string    `gorm:"size:10;not null" json:"gender"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Non-database fields for DTO mapping
	UserDTO  *dto.User  `gorm:"-"`
	UsersDTO *dto.Users `gorm:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// New creates a new User instance
func New() *User {
	return &User{}
}

// Create inserts a new user record into the database
func (u *User) Create(ctx context.Context) error {
	if err := database.Client().WithContext(ctx).Create(u).Error; err != nil {
		fmt.Printf("Unable to create user: %v\n", err)
		return err
	}
	return nil
}

// GetByID retrieves a user by their UUID
func (u *User) GetByID(ctx context.Context) error {
	if err := database.Client().WithContext(ctx).First(u, "id = ?", u.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("User not found: %v\n", err)
			return err
		}
		fmt.Printf("Error getting user: %v\n", err)
		return err
	}
	return nil
}

// GetByEmail retrieves a user by their email
func (u *User) GetByEmail(ctx context.Context, email string) error {
	if err := database.Client().WithContext(ctx).First(u, "email = ?", email).Error; err != nil {
		return err
	}
	return nil
}

// GetByAadhaarApplicationID retrieves a user by their Aadhaar application ID
func (u *User) GetByAadhaarApplicationID(ctx context.Context, aadhaarID string) error {
	if err := database.Client().WithContext(ctx).First(u, "aadhaar_application_id = ?", aadhaarID).Error; err != nil {
		return err
	}
	return nil
}

// GetAllPaginated retrieves users with pagination, sorting, and optional search
func (u *User) GetAllPaginated(ctx context.Context, params dto.PaginationParams) ([]User, int64, error) {
	var users []User
	var total int64

	db := database.Client().WithContext(ctx).Model(&User{})

	// Apply search filter if provided (searches name, email, or aadhaar_application_id)
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ? OR aadhaar_application_id ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count before pagination
	if err := db.Count(&total).Error; err != nil {
		fmt.Printf("Error counting users: %v\n", err)
		return nil, 0, err
	}

	// Apply sorting - using safe column mapping to prevent SQL injection
	sortColumn := getSafeColumnName(params.SortBy)
	sortOrder := getSafeSortOrder(params.Order)
	orderClause := fmt.Sprintf("%s %s", sortColumn, sortOrder)

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Apply pagination with LIMIT and OFFSET
	if err := db.Order(orderClause).
		Limit(params.Limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		fmt.Printf("Error getting users: %v\n", err)
		return nil, 0, err
	}

	return users, total, nil
}

// Delete removes a user from the database
func (u *User) Delete(ctx context.Context) error {
	if err := database.Client().WithContext(ctx).Delete(u).Error; err != nil {
		fmt.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

// getSafeColumnName maps user input to safe column names to prevent SQL injection
func getSafeColumnName(column string) string {
	safeColumns := map[string]string{
		"name":                    "name",
		"email":                   "email",
		"created_at":              "created_at",
		"aadhaar_application_id":  "aadhaar_application_id",
	}

	if safe, ok := safeColumns[column]; ok {
		return safe
	}
	return "created_at" // default
}

// getSafeSortOrder validates sort order to prevent SQL injection
func getSafeSortOrder(order string) string {
	if order == "asc" || order == "ASC" {
		return "ASC"
	}
	return "DESC" // default
}
