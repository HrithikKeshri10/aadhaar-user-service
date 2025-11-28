# Aadhaar User Profile Backend Service

A high-performance RESTful API built with Go and Fiber framework for managing Aadhaar application user profiles with pagination, sorting, and validation capabilities.

## 📱 Overview

The Aadhaar User Profile Backend Service is a modern backend service built with Go that provides a complete solution for storing, retrieving, and managing user details for new Aadhaar applications. It supports secure data storage, efficient retrieval with pagination and sorting, and comprehensive input validation.

## ✨ Features

- **User Management**
  - Create new Aadhaar application user records
  - Retrieve user by UUID
  - List users with pagination and sorting
  - Delete user records
  - Unique constraints on email and Aadhaar Application ID

- **Pagination & Sorting**
  - Configurable page size (1-100 items)
  - Sort by multiple fields (name, email, created_at, aadhaar_application_id)
  - Ascending/descending order support
  - Search functionality across multiple fields

- **Security & Validation**
  - Input validation using struct tags
  - SQL injection protection via parameterized queries
  - UUID-based identifiers
  - Comprehensive error handling
  - Request body size limits (16MB)

- **Performance**
  - Efficient database queries with GORM
  - Indexed columns for fast lookups
  - Connection pooling

## 🛠️ Tech Stack

- **Backend Framework:** Fiber v2.52.9
- **Language:** Go 1.21+
- **Database:** PostgreSQL
- **ORM:** GORM v1.31.0
- **Validation:** go-playground/validator v10
- **UUID:** google/uuid

## 📋 Prerequisites

Before running this application, ensure you have the following installed:

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd aadhaar-user-service
```

### 2. Configure PostgreSQL Database

Create a PostgreSQL database and run the migration:

```sql
-- Create database
CREATE DATABASE aadhaar_db;

-- Connect to the database and run migration
\c aadhaar_db
\i migrations/001_create_users_table.sql
```

Or manually:

```sql
CREATE DATABASE aadhaar_db;
CREATE USER your_db_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE aadhaar_db TO your_db_user;
```

### 3. Configure Environment Variables

Create a `.env` file or set environment variables:

```bash
# Database Configuration
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=aadhaar_db
export DB_SSLMODE=disable
```

### 4. Install Dependencies

```bash
go mod download
go mod tidy
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The application will start on `http://localhost:3015`

## 📡 API Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Service health status |

### User Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/aadhaar/users` | Create a new user |
| GET | `/aadhaar/users` | List users with pagination |
| GET | `/aadhaar/users/:id` | Get user by ID |
| DELETE | `/aadhaar/users/:id` | Delete user by ID |

## 📝 API Request Examples

### Create a User

```bash
POST /aadhaar/users
Content-Type: application/json

{
    "aadhaar_application_id": "12345678901234",
    "name": "Rahul Kumar",
    "email": "rahul.kumar@example.com",
    "phone": "9876543210",
    "address": "123 MG Road, Bangalore, Karnataka 560001",
    "date_of_birth": "1990-05-15",
    "gender": "male"
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "aadhaar_application_id": "12345678901234",
    "name": "Rahul Kumar",
    "email": "rahul.kumar@example.com",
    "phone": "9876543210",
    "address": "123 MG Road, Bangalore, Karnataka 560001",
    "date_of_birth": "1990-05-15",
    "gender": "male",
    "created_at": "2024-12-01T10:30:00Z"
}
```

### Get User by ID

```bash
GET /aadhaar/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "aadhaar_application_id": "12345678901234",
    "name": "Rahul Kumar",
    "email": "rahul.kumar@example.com",
    "phone": "9876543210",
    "address": "123 MG Road, Bangalore, Karnataka 560001",
    "date_of_birth": "1990-05-15",
    "gender": "male",
    "created_at": "2024-12-01T10:30:00Z",
    "updated_at": "2024-12-01T10:30:00Z"
}
```

### List Users with Pagination and Sorting

```bash
GET /aadhaar/users?page=1&limit=10&sort_by=created_at&order=desc&search=rahul
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 10 | Items per page (max: 100) |
| sort_by | string | created_at | Sort field (name, email, created_at, aadhaar_application_id) |
| order | string | desc | Sort order (asc, desc) |
| search | string | - | Search term (searches name, email, aadhaar_application_id) |

**Response (200 OK):**
```json
{
    "users": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "aadhaar_application_id": "12345678901234",
            "name": "Rahul Kumar",
            "email": "rahul.kumar@example.com",
            "phone": "9876543210",
            "address": "123 MG Road, Bangalore, Karnataka 560001",
            "date_of_birth": "1990-05-15",
            "gender": "male",
            "created_at": "2024-12-01T10:30:00Z",
            "updated_at": "2024-12-01T10:30:00Z"
        }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
}
```

### Delete User

```bash
DELETE /aadhaar/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (204 No Content)**

## 🗄️ Database Schema

### Users Table

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, auto-generated | Unique identifier |
| aadhaar_application_id | VARCHAR(14) | UNIQUE, NOT NULL | 14-character Aadhaar application ID |
| name | VARCHAR(100) | NOT NULL | Full name |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email address |
| phone | VARCHAR(10) | NOT NULL | 10-digit phone number |
| address | VARCHAR(500) | NOT NULL | Residential address |
| date_of_birth | VARCHAR(10) | NOT NULL | Date of birth (YYYY-MM-DD) |
| gender | VARCHAR(10) | NOT NULL, CHECK | Gender (male/female/other) |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Record creation time |
| updated_at | TIMESTAMP | AUTO-UPDATED | Last update time |

### Indexes

- `idx_users_email` - Unique index on email
- `idx_users_aadhaar_application_id` - Unique index on Aadhaar Application ID
- `idx_users_name` - Index on name for search
- `idx_users_created_at` - Index on created_at for sorting

## 📂 Project Structure

```
aadhaar-user-service/
├── cmd/
│   ├── app/
│   │   └── app.go              # Application initialization
│   └── main.go                 # Entry point
├── controllers/
│   └── users/
│       └── users.go            # User HTTP handlers
├── internals/
│   ├── config/
│   │   └── db.go               # Database migrations config
│   ├── database/
│   │   └── db.go               # PostgreSQL connection
│   ├── dto/
│   │   └── users.go            # Data Transfer Objects
│   ├── server/
│   │   ├── handlers.go         # Route handlers
│   │   ├── middleware.go       # Middleware setup
│   │   └── server.go           # Server configuration
│   └── validator/
│       ├── users.go            # User validation
│       └── utils.go            # Validation utilities
├── migrations/
│   └── 001_create_users_table.sql  # Database migration
├── models/
│   └── users/
│       └── users.go            # User database model
├── routes/
│   └── users.go                # User routes
├── services/
│   └── users/
│       └── users.go            # User business logic
├── .gitignore
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## ⚙️ Configuration

### Server Configuration
```go
Port: ":3015"
BodyLimit: 16 * 1024 * 1024 // 16MB
```

### Validation Rules
```
Aadhaar Application ID: exactly 14 characters
Name: 2-100 characters
Email: valid email format
Phone: exactly 10 numeric characters
Address: max 500 characters
Date of Birth: required (YYYY-MM-DD format)
Gender: one of (male, female, other)
```

### Pagination Limits
```
Default Page: 1
Default Limit: 10
Max Limit: 100
```

## 🔒 Security Features

- **UUID Identifiers:** Prevents sequential ID attacks
- **SQL Injection Protection:** Parameterized queries via GORM
- **Input Validation:** Comprehensive validation for all inputs
- **Safe Column Mapping:** Whitelisted sort columns
- **Request Size Limits:** 16MB body limit to prevent DoS
- **Error Handling:** Centralized error handling with appropriate status codes

## 🚦 Error Responses

All errors follow this format:

```json
{
    "error": "Error message",
    "details": [
        {
            "field": "FieldName",
            "message": "Validation message"
        }
    ]
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (successful delete) |
| 400 | Bad Request (validation error) |
| 404 | Not Found |
| 409 | Conflict (duplicate email/aadhaar_id) |
| 500 | Internal Server Error |

## 👨‍💻 Author

**Hrithik Keshri**
- LinkedIn: [Hrithik Keshri](https://linkedin.com/in/hrithikkeshri10)

---

Made with ❤️ using Go and Fiber
