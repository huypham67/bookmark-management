# CLAUDE.md - System Instructions for Bookmark Service

## 🎯 Core Philosophy

- **Clean Architecture**: Strict separation of concerns (Handler → Service → Repository).
- **Dependency Injection**: Use constructor injection for all layers. No global state.
- **DevOps-First**: Maintainability, observability (logging), and resiliency are priority.
- **Fail Fast**: Wrap errors, handle them at the service layer, and map to appropriate HTTP codes.

---

## 🏗️ Architecture & Module Organization

### Layer Structure

Each feature is organized in parallel hierarchies across three layers:

```
internal/
├── dto/                          # Data Transfer Objects (Input/Output contracts)
│   └── {feature}/
│       ├── request.go            # Input DTOs with validation tags
│       └── response.go           # Output DTOs
├── handler/                      # HTTP I/O layer
│   └── {feature}/
│       ├── handler.go            # Interface & DI constructor
│       ├── create.go             # POST operations
│       ├── list.go               # GET operations (paginated)
│       ├── get.go                # GET single resource
│       └── update.go             # PUT operations
├── service/                      # Business logic layer
│   └── {feature}/
│       ├── service.go            # Interface & DI constructor + constants/errors
│       ├── create.go             # Create business logic
│       ├── list.go               # List with pagination logic
│       └── {operation}.go        # Other operations
└── repository/                   # Data access layer (DB/Cache)
    └── {feature}/
        ├── repo.go               # Interface & DI constructor
        ├── read.go               # Read operations (queries, pagination)
        ├── write.go              # Create/Update operations
        └── mocks/                # Mockery-generated mocks
```

### Architectural Flow

```
Router → Handler → Service → Repository → Database
  ↓         ↓         ↓           ↓
 HTTP     Validation Business   Query
 I/O      & Binding  Logic      Builder
```

### Layer Responsibilities

| Layer | Responsibility | Example |
|-------|---|---|
| **Handler** | HTTP I/O, request binding, status codes. Never access DB. | Extract JWT, use `Bind()`, return 200/400/401 |
| **Service** | Pure business logic, orchestration, error wrapping. | Validate sort fields, calculate pagination offset, call repo |
| **Repository** | Data access (queries, aggregations). Return domain models. | Filter, count, order, limit, join tables |
| **DTO** | Input/Output contracts with validation rules via struct tags. | `binding:"required,min=1,max=100"`, `form:"page"` |

---

## 🛠️ Implementation Standards

### Language & Dependencies
- **Language**: Go 1.26
- **Framework**: Gin (HTTP routing)
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
- **Logging**: `rs/zerolog` for structured logging
- **Security**: JWT (RSA) for auth, `bcrypt` for passwords

### Code Organization by Operation

#### Create Flow (POST)
```
POST /v1/{feature}
  ↓
handler.go: Create(c *gin.Context)
  - Extract JWT user ID
  - Bind[CreateRequest]()  ← validation happens here
  - service.Create(ctx, userID, req)
  - Return 201 with response DTO
  ↓
service.go: Create(ctx, userID, req)
  - Validate business rules
  - Generate codes/IDs if needed
  - repo.Create(ctx, model)
  - Log & return model
  ↓
repo.go: Create(ctx, model)
  - db.Create(model)  ← GORM saves with auto-generated ID
```

#### List Flow (GET with Pagination)
```
GET /v1/{feature}?page=1&limit=10&sort=created_at
  ↓
handler.go: List(c *gin.Context)
  - Extract JWT user ID
  - Bind[ListRequest]()  ← validates page (min=1), limit (max=100), sort (oneof)
  - service.List(ctx, userID, page, limit, sort)
  - Build ListResponse with pagination metadata
  ↓
service.go: List(ctx, userID, page, limit, sort)
  - Calculate offset = (page - 1) * limit
  - repo.GetPaginatedByUserID(ctx, userID, offset, limit, sort)
  - repo.CountByUserID(ctx, userID)
  - Return both for pagination struct
  ↓
repo.go: GetPaginatedByUserID() + CountByUserID()
  - TWO separate queries: one for data, one for total count
  - Filter by user_id, order by sort field DESC, limit
```

### Validation Strategy

**Request validation happens in two places** (defense in depth):

1. **DTO Layer** (Struct Tags via `Bind()`):
   ```go
   type ListRequest struct {
     Page  int64  `form:"page" binding:"omitempty,min=1"`
     Limit int64  `form:"limit" binding:"omitempty,min=1,max=100"`
     Sort  string `form:"sort" binding:"omitempty,oneof=created_at updated_at"`
   }
   ```
   - Handles type conversion & basic constraints
   - Returns 400 immediately if invalid

2. **Handler Layer** (Semantic Validation):
   - Check authentication state
   - Verify user authorization
   - Return 401/403 for auth failures

3. **Service Layer** (Business Rules):
   - Domain-specific validation only (rarely needed after DTO validation)
   - Return domain errors (e.g., `ErrInternalServerError`)

### Error Handling Pattern

```go
// Service layer: Define domain errors
var (
    ErrInternalServerError = errors.New("internal server error")
)

// Handler layer: Map domain errors to HTTP codes
if err != nil {
    switch {
    case errors.Is(err, service.ErrInternalServerError):
        c.JSON(http.StatusInternalServerError, gin.H{"error": "..."})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "..."})
    }
    return
}
```

### Logging Pattern

Use `rs/zerolog` for structured logging at each layer:

```go
// Info: Success paths
log.Info().
    Str("user_id", userID).
    Str("bookmark_id", bm.ID).
    Msg("bookmark created successfully")

// Error: Failures
log.Error().
    Err(err).
    Str("user_id", userID).
    Msg("failed to create bookmark")

// Warn: Suspicious but recoverable
log.Warn().
    Str("sort_field", sort).
    Msg("invalid sort field requested")
```

### Testing Standards

- **Unit tests** (`_test.go`): Use `testify/mock` for repository/service mocks
- **Integration tests** (`internal/integration/`): Use real PostgreSQL & Redis
- **Coverage goal**: > 96.5%
- **No tests for handlers** (yet) — use integration tests instead

### Documentation

- **Swagger annotations** required for ALL new endpoints
- Use `/v1/` path prefix in router
- Tag by feature (e.g., `@Tags bookmarks`)
- Document query params with defaults and constraints

---

## 🧠 Claude Interaction Protocol

### Before Asking for Code Changes

1. **Reference files** using `@internal/path/file.go` notation
2. **Ask for analysis** before refactoring large features
3. **Validate constraints**:
   - [ ] Logic in correct layer?
   - [ ] All dependencies injected?
   - [ ] Errors wrapped correctly?
   - [ ] Swagger docs added/updated?

### When Debugging

Provide:
- Error log output from terminal
- Relevant service/handler file
- Current git status (`git diff`)

### Clarification Commands

- `/clear` — Reset context when switching features
- `/compact` — Use when history becomes large
- `@file` — Anchor specific files for analysis

---

## 🚀 Pre-Flight Checklist

Before finalizing a feature:

### 1. Database Changes
- [ ] Migration written in `migrations/`
- [ ] Model updated in `internal/model/`
- [ ] Repository updated with new methods

### 2. Dependency Injection
- [ ] New service registered in `internal/bootstrap/app.go`
- [ ] Handler initialized with service dependency
- [ ] Routes registered in `internal/api/router.go`

### 3. Code Quality
- [ ] Run `make fmt` (gofmt)
- [ ] Run `make vet` (staticcheck)
- [ ] Run `make test` (unit + integration tests)
- [ ] Coverage > 96.5%

### 4. API Documentation
- [ ] Swagger `@Summary` added
- [ ] `@Param` documented with types/constraints
- [ ] Request/response DTOs use struct tags for validation

### 5. Git Commit
- [ ] Changes are atomic (one feature = one commit)
- [ ] Commit message references layer (e.g., "feat(bookmark): add list endpoint")
- [ ] No leftover debug prints or commented code

---

## 📋 Common Patterns by Feature

### Adding a New Endpoint

**Template**: `POST /v1/widgets` (Create)

1. **DTO** (`internal/dto/widget/request.go`):
   ```go
   type CreateWidgetRequest struct {
       Name string `json:"name" binding:"required,min=2,max=100"`
       Type string `json:"type" binding:"required,oneof=typeA typeB"`
   }
   ```

2. **Handler** (`internal/handler/widget/create.go`):
   ```go
   func (h *handler) Create(c *gin.Context) {
       req, err := requestutils.Bind[widgetDTO.CreateWidgetRequest](c)
       if err != nil { /* 400 */ }
       widget, err := h.service.Create(c, *req)
       if err != nil { /* 500 */ }
       c.JSON(http.StatusCreated, response)
   }
   ```

3. **Service** (`internal/service/widget/create.go`):
   ```go
   func (s *service) Create(ctx context.Context, req widgetDTO.CreateWidgetRequest) (*model.Widget, error) {
       widget := &model.Widget{Name: req.Name, Type: req.Type}
       if err := s.repo.Create(ctx, widget); err != nil {
           log.Error().Err(err).Msg("failed to create widget")
           return nil, ErrInternalServerError
       }
       return widget, nil
   }
   ```

4. **Repository** (`internal/repository/widget/write.go`):
   ```go
   func (r *repository) Create(ctx context.Context, widget *model.Widget) error {
       return r.db.WithContext(ctx).Create(widget).Error
   }
   ```

---

## 🔐 Security Checklist

- JWT tokens validated at middleware layer
- Passwords hashed with `bcrypt`
- User ID extracted from JWT, not user input
- SQL injection prevented by GORM parameterization
- No sensitive data in logs

---

*Bookmark Service — Keep code clean, resilient, and production-ready.*
