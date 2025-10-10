# DEPENDENCIES.md

## External Libraries

- **github.com/go-chi/chi/v5**  
  HTTP router for building RESTful APIs with middleware support.

- **github.com/go-chi/cors**  
  Middleware for handling Cross-Origin Resource Sharing (CORS) in HTTP requests.

- **github.com/jmoiron/sqlx**  
  SQL database access library providing extensions on top of database/sql for easier database interactions.

- **github.com/lib/pq**  
  PostgreSQL database driver for Go.

- **github.com/golang-jwt/jwt/v5**  
  Library for creating and verifying JSON Web Tokens (JWT) for authentication.

- **go.uber.org/zap**  
  Structured, high-performance logging library.

- **golang.org/x/crypto/bcrypt**  
  Library for securely hashing and verifying passwords.

## Outbound Services/APIs

- **PostgreSQL Database**  
  - Used for persistent storage of users, meters, and readings.
  - Authenticated via connection string with user/password in environment variables.

- **No external cache, message bus, or third-party outbound APIs**  
  - The current codebase does not integrate with external caches (e.g., Redis), message buses (e.g., Kafka), or other APIs.

## Authentication for Outbound Services

- **Database**:  
  - Authenticated using credentials provided in the `DATABASE_URL` environment variable.
  - Example:  
    `postgres://postgres:postgres@localhost:5432/energy?sslmode=disable`

- **JWT**:  
  - Used for authenticating API requests from clients.
  - Tokens are signed with a secret (`JWT_SECRET` env var) and validated in middleware.

## Notes

- All authentication secrets and credentials are managed via environment variables and not hardcoded.
- If new outbound services are added (e.g., cache, message bus), they should be configured and authenticated via environment variables as well.