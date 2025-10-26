# MODULES.md

## List of Packages

- **cmd/server**  
  Entry point for the application; starts the HTTP server and wires dependencies.

- **internal/config**  
  Loads and manages application configuration from environment variables.

- **internal/db**  
  Handles database connection pooling and setup for PostgreSQL.

- **internal/domain**  
  Contains core business logic, domain models, and repository/service abstractions.

- **internal/http**  
  Sets up HTTP routing and middleware for the API.

- **internal/http/handlers**  
  Implements HTTP handlers for authentication, meters, and readings endpoints.

- **internal/pkg/userctx**  
  Provides utilities for managing user context within requests.

- **internal/util**  
  Utility functions, such as JWT signing and verification.

- **migrations**  
  SQL migration scripts for initializing and updating the database schema.