# CLAUDE.md - System Instructions for Bookmark Service

## 🎯 Core Philosophy
- **Clean Architecture**: Strict separation of concerns (Handler -> Service -> Repository).
- **Dependency Injection**: Use constructor injection for all layers. No global state.
- **DevOps-First**: Maintainability, observability (logging), and resiliency are priority.
- **Fail Fast**: Wrap errors, handle them at the service layer, and map to appropriate HTTP codes.

## 🏗️ Architectural Rules
1. **Handler (`internal/handler/`)**: Handles HTTP I/O, validation, and status codes. Never access DB directly.
2. **Service (`internal/service/`)**: Contains pure business logic. Orchestrates domain operations.
3. **Repository (`internal/repository/`)**: Data access layer (PostgreSQL/Redis).
4. **DTO (`internal/dto/`)**: Input/Output contracts. Use `go-playground/validator` tags.
5. **Flow**: `Router` -> `Handler` -> `Service` -> `Repository`.

## 🛠️ Implementation Standards
- **Language**: Go 1.26.
- **Logging**: Use `rs/zerolog` for all structured logging.
- **Security**: JWT (RSA) for auth, `bcrypt` for passwords.
- **Testing**:
    - Unit tests (`_test.go`) use `testify/mock`.
    - Integration tests (`internal/integration/`) use real DB/Redis.
    - Coverage goal: > 96.5%.
- **Documentation**: Swagger annotations required for all new endpoints.

## 🧠 Claude Interaction Protocol
- **Context First**: Always reference files using `@filename` before asking for changes.
- **Analyze Before Code**: When refactoring, ask for an analysis of dependency impacts first.
- **Constraint Checklist**:
    - [ ] Is this logic in the correct layer?
    - [ ] Are all dependencies injected?
    - [ ] Is error handling wrapped correctly?
    - [ ] Did I add/update Swagger docs?
- **Proactive Debugging**: When an error occurs, paste the terminal log and the relevant service file for context.

## 🚀 Key Commands for Claude
- `/clear`: Reset context when switching features.
- `/compact`: Use when history grows too large.
- `@file`: Use to anchor specific files for Claude to analyze.

## 🚦 Pre-Flight Checklist
1. **Migration**: New DB field? Run `migrations/`.
2. **DI**: Is the new service registered in `internal/bootstrap/app.go`?
3. **Test**: Run `make test` before finalizing code.
4. **Fmt**: Run `make fmt` and `make vet`.

---
*Bookmark Service - Keep code clean, resilient, and production-ready.*