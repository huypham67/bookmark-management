# Bookmark Service

A production-ready REST API service for user authentication, profile management, and bookmark management built with Go, Gin framework, PostgreSQL, and Redis. Follows clean architecture principles with strict separation of concerns.

## Overview

Bookmark Service is a modern, scalable microservice designed for user management and bookmark operations with a focus on reliability, performance, and maintainability. It provides JWT-based authentication, comprehensive error handling, structured logging with Zerolog, extensive testing (>96% coverage), and complete API documentation using Swagger/OpenAPI.

## 🎯 Features

- **User Authentication**: Registration and login with JWT tokens (RSA-based)
- **User Profiles**: Get and update user information
- **Bookmark Management**: Create, read, update, and delete bookmarks with pagination
- **URL Shortening**: Create shortened URLs and redirect functionality
- **Health Check Endpoint**: Monitor service status and database connectivity
- **PostgreSQL Backend**: Persistent data storage with GORM ORM
- **Redis Cache**: Optional caching layer for performance
- **Swagger/OpenAPI Documentation**: Interactive API docs at `/swagger/`
- **Environment Configuration**: Flexible setup via environment variables
- **Structured Logging**: Zerolog integration for comprehensive logging
- **Comprehensive Testing**: Unit and integration tests with 96.5%+ coverage
- **Docker Ready**: Optimized Dockerfile for containerization
- **Database Migrations**: Golang-migrate for schema management
- **Cross-Platform Build**: Support for Linux, macOS, and Windows

## 📋 Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.26 |
| Web Framework | Gin | v1.12.0 |
| Database | PostgreSQL | (via GORM) |
| ORM | GORM | v1.6.0 |
| Cache | Redis | v9.19.0 |
| Logger | Zerolog | v1.35.1 |
| Auth | JWT (RSA) | v5.3.1 |
| Password Hashing | bcrypt | (golang.org/x/crypto) |
| Migrations | golang-migrate | v4.19.1 |
| API Documentation | Swagger/OpenAPI | v1.16.6 |
| Testing | Testify | v1.11.1 |
| UUID Generation | google/uuid | v1.6.0 |
| Config Management | envconfig | v1.4.0 |

## 🚀 Quick Start

### Prerequisites

- **Go 1.26** or higher
- **Git**
- **PostgreSQL 12+** (required)
- **Redis** (optional, for caching)
- **Make** or **PowerShell** (for Windows)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/huypham67/bookmark-service.git
   cd bookmark-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Install development tools** (Optional)
   ```bash
   make install-tools    # Linux/macOS/WSL
   ```

### Environment Setup

Create a `.env` file in the project root:

```env
# Application Configuration
APP_PORT=8080
SERVICE_NAME=bookmark-service
INSTANCE_ID=instance-1

# JWT Configuration
JWT_PRIVATE_KEY_PATH=/keys/private.pem
JWT_PUBLIC_KEY_PATH=/keys/public.pem
JWT_ISSUER=bookmark-service
JWT_AUDIENCE=bookmark-app
JWT_EXPIRATION_SECONDS=3600

# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=bookmark_db
DB_SSLMODE=disable
DB_TIMEZONE=UTC

# Redis Configuration (optional - defaults to localhost:6379)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Configuration Reference:**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_PORT` | No | 8080 | Port on which the API server runs |
| `SERVICE_NAME` | Yes | - | Name of the service for health checks |
| `INSTANCE_ID` | No | Auto-generated UUID | Unique identifier for this service instance |
| `JWT_PRIVATE_KEY_PATH` | Yes | - | Path to JWT private key (RSA) |
| `JWT_PUBLIC_KEY_PATH` | Yes | - | Path to JWT public key (RSA) |
| `JWT_ISSUER` | No | bookmark-service | JWT issuer claim |
| `JWT_AUDIENCE` | No | bookmark-app | JWT audience claim |
| `JWT_EXPIRATION_SECONDS` | No | 3600 | JWT token expiry in seconds (default 1 hour) |
| `DB_HOST` | No | localhost | PostgreSQL host |
| `DB_PORT` | No | 5432 | PostgreSQL port |
| `DB_USER` | No | postgres | PostgreSQL username |
| `DB_PASSWORD` | No | postgres | PostgreSQL password |
| `DB_NAME` | No | bookmark_service | PostgreSQL database name |
| `DB_SSL_MODE` | No | disable | PostgreSQL SSL mode |
| `REDIS_ADDR` | No | localhost:6379 | Redis connection address |
| `REDIS_PASSWORD` | No | Empty | Redis password |
| `REDIS_DATABASE` | No | 0 | Redis database number |
| `DB_SSLMODE` | No | disable | PostgreSQL SSL mode |
| `DB_TIMEZONE` | No | UTC | PostgreSQL timezone |

### Database Setup

1. **Create PostgreSQL database**
   ```bash
   createdb bookmark_db
   ```

2. **Run migrations** (required to create tables)
   ```bash
   make migrate-up    # Apply all migrations
   make migrate-down  # Rollback migrations (use with caution)
   ```

### Running the Application

**Using Make (Linux/macOS/WSL):**
```bash
make run           # Run the application
make dev           # Full development workflow (fmt + vet + test + swagger + run)
make test          # Run all tests with coverage
```

**Direct Go command:**
```bash
go run ./cmd/api/main.go
```

The API will be available at `http://localhost:8080/api/bookmark_service/v1` and Swagger docs at `http://localhost:8080/swagger/`

## 📁 Project Structure

```
bookmark-service/
├── cmd/
│   └── api/
│       └── main.go                      # Application entry point
├── internal/
│   ├── api/
│   │   └── router.go                    # Route definitions and setup
│   ├── bootstrap/
│   │   └── app.go                       # Application initialization and DI
│   ├── config/
│   │   └── config.go                    # Configuration management
│   ├── dto/
│   │   ├── auth/
│   │   │   ├── request.go               # Login/register DTOs
│   │   │   └── response.go
│   │   ├── bookmark/
│   │   │   ├── request.go               # Bookmark request DTOs
│   │   │   └── response.go
│   │   ├── health/
│   │   │   └── response.go
│   │   ├── link/
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   └── profile/
│   │       ├── request.go
│   │       └── response.go
│   ├── handler/
│   │   ├── auth/
│   │   │   ├── handler.go               # Auth handler interface & DI
│   │   │   ├── register.go              # User registration handler
│   │   │   ├── login.go                 # User login handler
│   │   │   ├── register_test.go
│   │   │   └── login_test.go
│   │   ├── bookmark/
│   │   │   ├── handler.go               # Bookmark handler interface & DI
│   │   │   ├── create.go                # Create bookmark
│   │   │   ├── list.go                  # List bookmarks with pagination
│   │   │   ├── update.go                # Update bookmark
│   │   │   └── delete.go                # Delete bookmark
│   │   ├── health/
│   │   │   ├── handler.go
│   │   │   ├── check.go                 # Health check handler
│   │   │   └── check_test.go
│   │   ├── link/
│   │   │   ├── handler.go               # Link handler interface & DI
│   │   │   ├── shorten.go               # URL shortening handler
│   │   │   ├── redirect.go              # URL redirect handler
│   │   │   ├── shorten_test.go
│   │   │   └── redirect_test.go
│   │   └── profile/
│   │       ├── handler.go               # Profile handler interface & DI
│   │       ├── get.go                   # Get user profile
│   │       ├── update.go                # Update user profile
│   │       ├── get_test.go
│   │       └── update_test.go
│   ├── model/
│   │   ├── base.go                      # Base model with timestamps
│   │   ├── user.go                      # User domain model
│   │   └── bookmark.go                  # Bookmark domain model
│   ├── repository/
│   │   ├── bookmark/
│   │   │   ├── repo.go                  # Bookmark repository interface & DI
│   │   │   ├── read.go                  # Read operations (list, get)
│   │   │   ├── write.go                 # Create/update/delete operations
│   │   │   ├── write_test.go
│   │   │   └── mocks/
│   │   ├── link/
│   │   │   ├── repo.go
│   │   │   ├── read.go
│   │   │   ├── write.go
│   │   │   ├── write_test.go
│   │   │   └── mocks/
│   │   ├── ping/
│   │   │   └── repo.go                  # Redis ping for health checks
│   │   └── user/
│   │       ├── repo.go                  # User repository interface & DI
│   │       ├── read.go                  # User read operations
│   │       ├── write.go                 # User create/update operations
│   │       └── mocks/
│   ├── service/
│   │   ├── auth/
│   │   │   ├── service.go               # Auth service interface & DI
│   │   │   ├── register.go              # Registration business logic
│   │   │   ├── login.go                 # Login business logic
│   │   │   ├── register_test.go
│   │   │   └── login_test.go
│   │   ├── bookmark/
│   │   │   ├── service.go               # Bookmark service interface & DI
│   │   │   ├── create.go                # Create business logic
│   │   │   ├── list.go                  # List with pagination logic
│   │   │   ├── update.go                # Update business logic
│   │   │   ├── delete.go                # Delete business logic
│   │   │   └── mocks/
│   │   ├── health/
│   │   │   ├── service.go               # Health check service
│   │   │   ├── check.go
│   │   │   ├── check_test.go
│   │   │   └── mocks/
│   │   ├── link/
│   │   │   ├── service.go               # Link service interface & DI
│   │   │   ├── shorten.go               # URL shortening logic
│   │   │   └── redirect.go              # URL redirect logic
│   │   └── profile/
│   │       ├── service.go               # Profile service interface & DI
│   │       ├── get.go                   # Get user info logic
│   │       └── update.go                # Update user info logic
│   ├── integration/
│   │   ├── test_helper.go               # Integration test setup
│   │   ├── user_create_test.go
│   │   ├── login_test.go
│   │   ├── update_user_info_test.go
│   │   ├── link_shorten_test.go
│   │   └── ...
│   └── testutil/
│       └── ...                          # Test utilities
├── pkg/
│   ├── common/
│   │   └── ...                          # Common utilities
│   ├── jwtutils/
│   │   ├── config.go                    # JWT configuration
│   │   ├── claims.go                    # JWT claims handling
│   │   └── ...                          # JWT utilities and validation
│   ├── logger/
│   │   ├── logger.go                    # Zerolog configuration
│   │   └── config.go
│   ├── password/
│   │   └── ...                          # Password hashing utilities
│   ├── redis/
│   │   ├── config.go                    # Redis configuration
│   │   └── ...                          # Redis client wrapper
│   ├── requestutils/
│   │   └── ...                          # Request binding utilities
│   ├── response/
│   │   └── ...                          # Response formatting utilities
│   ├── security/
│   │   └── ...                          # Security utilities
│   ├── sqldb/
│   │   └── config.go                    # PostgreSQL configuration
│   └── utils/
│       └── ...                          # General utilities
├── middleware/
│   └── jwt.go                           # JWT authentication middleware
├── migrations/
│   └── *.sql                            # Database migration files
├── docs/
│   ├── docs.go                          # Generated Swagger documentation
│   ├── swagger.json
│   └── swagger.yaml
├── coverage/                            # Test coverage reports
├── Dockerfile                           # Docker configuration
├── docker-compose.yml                   # Docker Compose setup
├── Makefile                             # Build automation
├── CLAUDE.md                            # Development guidelines
├── go.mod                               # Go module definition
├── go.sum                               # Module checksums
├── .env                                 # Environment variables (local)
├── .gitignore
└── README.md                            # This file
```

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080/api/bookmark_service/v1
```

### Authentication Endpoints

#### Register User
Create a new user account.

```http
POST /users/register
```

**Request Body:**
```json
{
  "display_name": "John Doe",
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Parameters:**
- `display_name` (string, required): User's display name (2-100 chars)
- `username` (string, required): Username (3-50 chars, alphanumeric)
- `email` (string, required): Email address (must be valid email)
- `password` (string, required): Password (minimum 8 chars)

**Response (201 Created):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "display_name": "John Doe",
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2025-06-06T12:00:00Z"
  },
  "message": "Register an user successfully!"
}
```

**Status Codes:**
- `201 Created` - User registered successfully
- `400 Bad Request` - Invalid input or validation failed
- `409 Conflict` - Username or email already exists
- `500 Internal Server Error` - Server error

---

#### Login User
Authenticate and receive JWT token.

```http
POST /users/login
```

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "SecurePassword123"
}
```

**Parameters:**
- `username` (string, required): Username
- `password` (string, required): Password (minimum 8 chars)

**Response (200 OK):**
```json
{
  "data": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ...",
  "message": "Logged in successfully!"
}
```

**Status Codes:**
- `200 OK` - Login successful
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Invalid credentials
- `500 Internal Server Error` - Server error

---

### Profile Endpoints
*Requires JWT token in Authorization header: `Bearer <token>`*

#### Get User Profile
Retrieve current user information.

```http
GET /self/info
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "display_name": "John Doe",
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2025-06-06T12:00:00Z"
  },
  "message": "User information retrieved successfully!"
}
```

**Status Codes:**
- `200 OK` - Profile retrieved
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

#### Update User Profile
Update current user information.

```http
PUT /self/info
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "display_name": "John Doe Updated",
  "email": "john.new@example.com"
}
```

**Parameters:**
- `display_name` (string, optional): Display name (2-100 chars)
- `email` (string, optional): Email address (must be valid email)

**Response (200 OK):**
```json
{
  "message": "Edit current user successfully!"
}
```

**Status Codes:**
- `200 OK` - Profile updated
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing or invalid token
- `409 Conflict` - Email already exists
- `500 Internal Server Error` - Server error

---

### Bookmark Endpoints
*Requires JWT token in Authorization header: `Bearer <token>`*

#### Create Bookmark
Create a new bookmark for the authenticated user.

```http
POST /bookmarks
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "url": "https://example.com",
  "description": "A useful example website"
}
```

**Parameters:**
- `url` (string, required): The URL to bookmark (must be valid URL)
- `description` (string, required): Bookmark description (2-500 chars)

**Response (201 Created):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "code": "ab12cd34",
    "url": "https://example.com",
    "description": "A useful example website",
    "created_at": "2025-06-06T12:00:00Z",
    "updated_at": "2025-06-06T12:00:00Z"
  },
  "message": "Bookmark created successfully!"
}
```

**Status Codes:**
- `201 Created` - Bookmark created
- `400 Bad Request` - Invalid input or validation failed
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

#### List Bookmarks
Get paginated list of user's bookmarks.

```http
GET /bookmarks?page=1&limit=10&sort=created_at
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer, optional, default=1): Page number (minimum 1)
- `limit` (integer, optional, default=10): Items per page (max 100)
- `sort` (string, optional, default=created_at): Sort field (created_at, updated_at, code, url)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "code": "ab12cd34",
      "url": "https://example.com",
      "description": "A useful example website",
      "created_at": "2025-06-06T12:00:00Z",
      "updated_at": "2025-06-06T12:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "code": "ef56gh78",
      "url": "https://another-example.com",
      "description": "Another example",
      "created_at": "2025-06-05T10:00:00Z",
      "updated_at": "2025-06-05T10:00:00Z"
    }
  ],
  "message": "Bookmarks retrieved successfully!",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2
  }
}
```

**Status Codes:**
- `200 OK` - Bookmarks retrieved
- `400 Bad Request` - Invalid pagination parameters
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

#### Update Bookmark
Update an existing bookmark.

```http
PUT /bookmarks/:id
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "url": "https://example.com/updated",
  "description": "Updated description"
}
```

**Parameters:**
- `url` (string, optional): Updated URL (must be valid URL if provided)
- `description` (string, optional): Updated description (2-500 chars if provided)

**Response (200 OK):**
```json
{
  "message": "Success"
}
```

**Status Codes:**
- `200 OK` - Bookmark updated
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Bookmark not found
- `500 Internal Server Error` - Server error

---

#### Delete Bookmark
Delete a bookmark.

```http
DELETE /bookmarks/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Success"
}
```

**Status Codes:**
- `200 OK` - Bookmark deleted
- `400 Bad Request` - Invalid bookmark ID
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Bookmark not found
- `500 Internal Server Error` - Server error

---

### Link Endpoints

#### Shorten URL
Create a shortened URL.

```http
POST /links/shorten
```

**Request Body:**
```json
{
  "url": "https://example.com/very/long/url",
  "exp": 86400
}
```

**Parameters:**
- `url` (string, required): The original URL to shorten
- `exp` (integer, optional): Expiration time in seconds

**Response (200 OK):**
```json
{
  "code": "ab12cd34",
  "message": "Shorten URL generated successfully"
}
```

**Note:** The shortened code can be used with the redirect endpoint: `GET /links/redirect/{code}`

**Status Codes:**
- `200 OK` - URL shortened successfully
- `400 Bad Request` - Invalid input
- `500 Internal Server Error` - Server error

---

#### Redirect to Original URL
Redirect from shortened code to original URL.

```http
GET /links/redirect/:code
```

**Parameters:**
- `code` (path parameter, required): The shortened code

**Response:**
- `302 Found` - Redirects to the original URL
- `404 Not Found` - Code doesn't exist or expired
- `500 Internal Server Error` - Server error

---

### Health Check Endpoint

#### Check Service Health
Check application health status and database connectivity.

```http
GET /health-check
```

**Response (200 OK):**
```json
{
  "message": "OK",
  "service_name": "bookmark-service",
  "instance_id": "instance-1"
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `500 Internal Server Error` - Database or service error

---

## Interactive API Documentation

Access Swagger UI at:
```
http://localhost:8080/swagger/
```

This provides an interactive interface to test all API endpoints.

## 🏗️ Architecture

### Clean Architecture Pattern
The project follows clean architecture principles with clear separation of concerns across four layers:

```
HTTP Request
    ↓
Handler Layer (HTTP I/O)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Database (PostgreSQL/Redis)
```

### Layer Responsibilities

- **Handler Layer** (`internal/handler/`): HTTP request/response handling, input binding, validation, and status codes. Never accesses database directly.
- **Service Layer** (`internal/service/`): Business logic, orchestration, error wrapping, and domain operations.
- **Repository Layer** (`internal/repository/`): Data persistence abstraction. Handles queries, filtering, and pagination.
- **DTO Layer** (`internal/dto/`): Data transfer objects with validation rules via struct tags.
- **Model Layer** (`internal/model/`): Domain models representing database entities.

### Design Patterns Used

- **Repository Pattern**: Abstracts database access
- **Service Pattern**: Encapsulates business logic
- **Dependency Injection**: Loose coupling via constructor injection
- **Handler Pattern**: Clean HTTP request handling
- **Middleware Pattern**: JWT authentication at the middleware layer

### Dependency Injection

The application uses constructor-based dependency injection initialized in the bootstrap layer for loose coupling and testability.

## 🧪 Testing

### Run Tests

**Using Make:**
```bash
make test              # Run all tests with coverage
make test-coverage     # Generate HTML coverage report
```

**Direct Go:**
```bash
go test -v -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage

Current coverage: **96.5%+** on testable code

**Coverage breakdown by layer:**
- **Handlers**: HTTP request/response handling, validation, error cases
- **Services**: Business logic, error handling, domain operations
- **Repository**: Database operations, query building, error classification
- **Integration**: End-to-end API flows with real database
- **Utilities**: Helper functions

**Excluded from coverage (infrastructure):**
- `cmd/` - Application entry point
- `internal/bootstrap/` - Dependency injection setup
- `internal/config/` - Configuration initialization
- `pkg/logger/`, `pkg/redis/` - Infrastructure boilerplate
- `docs/` - Auto-generated Swagger documentation

### Test Types

1. **Unit Tests**: Individual layer testing with mocks
   - Location: `*_test.go` files alongside source code
   - Uses: `testify/assert`, `testify/mock`, `testify/require`
   - Example: `internal/handler/auth/login_test.go`

2. **Integration Tests**: Full API flow testing with real database
   - Location: `internal/integration/`
   - Database: PostgreSQL with test fixtures
   - Features: Real request/response handling and database operations
   - Example: `internal/integration/login_test.go`

3. **Mock Generation**: Interfaces use mockery for testing
   - Command: `make generate-mocks`
   - Generated mocks: `**/mocks/` directories

### Key Test Files

| Package | Test Coverage | Notes |
|---------|---|---|
| `internal/handler/auth` | ✅ Registration & login endpoints | Uses mocks for service layer |
| `internal/handler/bookmark` | ✅ CRUD operations | All bookmark operations |
| `internal/handler/profile` | ✅ User profile management | Get and update operations |
| `internal/service/auth` | ✅ Registration & login logic | Database and password validation |
| `internal/service/bookmark` | ✅ Bookmark business logic | Create, list, update, delete |
| `internal/repository/user` | ✅ Database operations | Real PostgreSQL with test data |
| `internal/integration` | ✅ End-to-end API flows | Full workflows with real database |

## 📦 Building

### Development Build

```bash
make build        # Build binary for current OS
```

### Cross-Platform Builds

```bash
make build-linux    # Build for Linux (amd64)
make build-macos    # Build for macOS (arm64)
make build-windows  # Build for Windows (amd64)
make release        # Build for all platforms
```

### Code Quality

**Format Code:**
```bash
make fmt          # Format all Go files
```

**Run Go Vet:**
```bash
make vet          # Run go vet analysis
```

## 🐳 Docker

### Build Docker Image

```bash
# Build production binary
make build-prod

# Build Docker image
docker build -t bookmark-service:latest .
```

### Run Docker Container

```bash
docker run -d \
  -e APP_PORT=8080 \
  -e SERVICE_NAME=bookmark-service \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=admin \
  -e DB_PASSWORD=admin \
  -e DB_NAME=bookmark_db \
  -e JWT_PRIVATE_KEY_PATH=/keys/private.pem \
  -e JWT_PUBLIC_KEY_PATH=/keys/public.pem \
  -v /path/to/keys:/keys \
  -p 8080:8080 \
  --name bookmark-service \
  bookmark-service:latest
```

### Docker Compose

```bash
docker-compose up -d
```

## 🔄 Development Workflow

### Quick Development

```bash
make fmt           # Format code
make vet           # Check code issues
make test          # Run tests
go run ./cmd/api/main.go  # Run application
```

### Full Development Workflow

```bash
make dev           # Format → Vet → Test → Swagger → Run
```

## 🧹 Cleanup

**Remove build artifacts:**
```bash
make clean         # Remove binaries and coverage
```

**Remove documentation:**
```bash
make clean-docs    # Remove Swagger docs
```

**Full cleanup:**
```bash
make clean-all     # Remove everything including vendor
```

## 🛠️ Available Make Targets

```bash
make help          # Display all available targets
```

### Common Targets

| Target | Description |
|--------|-------------|
| `make run` | Run the application |
| `make test` | Run tests with coverage |
| `make build` | Build binary for current OS |
| `make fmt` | Format code with gofmt |
| `make vet` | Run go vet analysis |
| `make swagger` | Generate Swagger documentation |
| `make clean` | Remove build artifacts |
| `make help` | Show all available targets |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Rollback database migrations |

## 🚧 Development Guidelines

### Adding New Features

1. Create DTOs in `internal/dto/{feature}/` (request.go, response.go)
2. Create handler in `internal/handler/{feature}/` with proper HTTP methods
3. Create service in `internal/service/{feature}/` with business logic
4. Create/update repository in `internal/repository/{feature}/` for data access
5. Register routes in `internal/api/router.go`
6. Add Swagger documentation comments
7. Write unit tests for handlers, services, and repositories
8. Write integration tests in `internal/integration/`
9. Regenerate Swagger docs: `make swagger`
10. Update this README.md if new endpoints added

### Code Organization Standards

**Follow the CLAUDE.md file for:**
- Clean architecture principles
- Layer separation and responsibilities
- Error handling patterns
- Logging standards
- Testing conventions
- Swagger documentation requirements

### Pre-Commit Checklist

- [ ] Code formatted: `make fmt`
- [ ] No linting issues: `make vet`
- [ ] All tests pass: `make test`
- [ ] Coverage > 96.5%
- [ ] Swagger docs updated
- [ ] No debug prints or commented code

## 📝 Configuration Reference

See `.env` file for all available configuration variables.

## 🔐 Security Considerations

- **JWT Authentication**: RSA-based JWT tokens for secure user sessions
- **Password Hashing**: bcrypt with appropriate salt rounds
- **Environment Variables**: Store secrets in `.env` (not in version control)
- **Input Validation**: All inputs validated at DTO layer
- **SQL Injection Prevention**: GORM ORM prevents SQL injection
- **CORS**: Configure as needed for production
- **HTTPS**: Use reverse proxy (nginx, etc.) in production
- **Rate Limiting**: Consider adding middleware for production

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

For issues, questions, or suggestions, please create an issue on GitHub.

## 🔗 Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Swagger/OpenAPI](https://swagger.io/)
- [JWT Introduction](https://jwt.io/introduction)
- [Redis Documentation](https://redis.io/docs/)
- [Zerolog Logger](https://github.com/rs/zerolog)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
