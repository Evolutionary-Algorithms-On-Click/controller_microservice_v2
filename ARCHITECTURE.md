# Architecture and Development Guide

This document outlines the architectural patterns, conventions, and development guidelines for the `controller_microservice_v2` project. Its purpose is to ensure that new code contributions, whether from human developers or AI agents, are consistent with the existing design.

## 1. Core Architecture: Layered Approach

The service follows a classic layered architecture to ensure a clean separation of concerns. The flow of a request is as follows:

**`main.go` -> `routes` -> `controllers` -> `modules` -> `repository` -> `db`**

- **`cmd/main.go`**: The application entry point. It is responsible for initializing all dependencies (Logger, Database Pool, Clients) and injecting them into the layers that require them.
- **`routes`**: Defines the API endpoints, mapping HTTP methods and URL patterns to specific handler functions in the `controllers`.
- **`controllers`**: The HTTP layer. Handlers in this layer are responsible for parsing and validating requests, calling the appropriate business logic in the `modules` layer, and formatting the response (e.g., writing a JSON payload and an HTTP status code). **Controllers should not contain business logic.**
- **`modules`**: The business logic layer. This layer orchestrates tasks by calling data access functions and other services. It is completely unaware of the HTTP layer.
- **`repository` (To Be Implemented):** The data access layer. This layer abstracts the database. It provides interfaces for data operations (e.g., `CreateNotebook`, `GetNotebookByID`). The `modules` layer should depend on these interfaces, not on a concrete database implementation.
- **`db`**: The database implementation layer. It manages the database connection (`pgxpool`) and contains the concrete implementation of the repository interfaces, including the raw SQL queries.

## 2. API Design and Routing

- **Versioning**: All API endpoints are versioned under `/api/v1/`.
- **RESTful Principles**: The API should adhere to RESTful principles, using appropriate HTTP methods (`GET`, `POST`, `PUT`, `DELETE`) for corresponding actions.
- **Routing**: We use the standard library's `http.ServeMux` (Go 1.22+). Routes must be registered explicitly with both the HTTP method and the path pattern in `routes/api.go`.

  ```go
  // Good:
  mux.HandleFunc("POST /api/v1/notebooks", notebookController.CreateNotebookHandler)
  mux.HandleFunc("GET /api/v1/notebooks/{id}", notebookController.GetNotebookByIDHandler)

  // Bad (Avoid):
  // mux.HandleFunc("/api/v1/notebooks/", monolithicHandler)
  ```

## 3. Database and Data Access

- **Schema**: The single source of truth for the database schema is `db/schema.sql`.
- **Repository Pattern**: All database access from the business logic layer (`modules`) **must** go through a repository interface. This decouples the business logic from the database, making the code easier to test and maintain. Do not use `db.Pool` directly within the `modules` package.

  **Example:**
  ```go
  // 1. Define the interface (e.g., in a new `storage` or `repository` package)
  type NotebookRepository interface {
      GetByID(ctx context.Context, id string) (*models.Notebook, error)
  }

  // 2. The module depends on the interface
  type NotebookModule struct {
      Repo NotebookRepository
  }

  // 3. The implementation with SQL lives in the db/repository layer
  type PostgresNotebookRepo struct {
      DB *pgxpool.Pool
  }

  func (p *PostgresNotebookRepo) GetByID(ctx context.Context, id string) (*models.Notebook, error) {
      // SQL query logic goes here...
  }
  ```

## 4. Logging

- **Library**: We use `zerolog` for structured, high-performance logging.
- **No `fmt.Print*`**: The use of `fmt.Println`, `fmt.Printf`, or `log.Print*` is strictly forbidden in the application code. This is enforced by a pre-commit hook.
- **Dependency Injection**: The `zerolog.Logger` instance is initialized in `main.go` and passed as a dependency to any struct that needs to log messages. Always use this injected logger instance.

  ```go
  // In a controller method:
  c.Logger.Info().Str("kernel_id", kernel.ID).Msg("Kernel started successfully")
  c.Logger.Error().Err(err).Msg("Failed to retrieve kernel list")
  ```

## 5. Configuration

- **Environment Variables**: All configuration (database URLs, auth tokens, ports) **must** be supplied via environment variables.
- **No Hardcoded Values**: Do not hardcode configuration values in the source code. In development, these are loaded from the `.env` file by `godotenv` in `main.go`.

## 6. Directory Structure Overview

- **`/cmd`**: Main application entry point.
- **`/db`**: Database schema, connection, and repository implementations.
- **`/routes`**: API route definitions.
- **`/controllers`**: HTTP request/response handlers.
- **`/modules`**: Business logic layer.
- **`/pkg`**: Shared libraries and utilities safe for external use.
  - **`/pkg/models`**: Core data structures (structs).
  - **`/pkg/jupyter_client`**: A dedicated client for the Jupyter Gateway service.
  - **`/pkg/culler`**: Background process for cleaning up idle kernels.
- **`/scripts`**: Helper scripts for development and CI.
- **`/docker`**: Dockerfiles for containerizing services.
