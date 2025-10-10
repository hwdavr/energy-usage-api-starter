# Energy Usage API - Architecture

## Purpose & Scope

This service provides a RESTful API for managing users, meters, and energy readings. It allows users to sign up, log in, register energy meters, submit readings, and query usage statistics. The service is designed for extensibility, security, and ease of integration with frontend or other backend systems.

---

## High-Level Architecture Diagram

```
+-------------------+
|    HTTP Server    |
| (cmd/server/main) |
+-------------------+
          |
          v
+-------------------+
|   Router &        |
| Middleware        |
| (internal/http/)  |
+-------------------+
          |
          v
+-------------------+
|   Handlers        |
| (auth, meters,    |
|  readings)        |
+-------------------+
          |
          v
+-------------------+
|   Domain Service  |
| (business logic)  |
| (internal/domain) |
+-------------------+
          |
          v
+-------------------+
|   Repository      |
| (DB access)       |
| (internal/domain) |
+-------------------+
          |
          v
+-------------------+
|   PostgreSQL DB   |
+-------------------+
```

---

## Key Flows

### Example: Add Meter Reading

1. **Request**:  
   `POST /v1/meters/{meterID}/readings`  
   → HTTP Server (`main.go`)
2. **Router**:  
   → Auth Middleware (JWT validation)  
   → Route matched in `internal/http/router.go`
3. **Handler**:  
   → `ReadingsHandler.Add` parses request, extracts user/meter info
4. **Service**:  
   → `Service.AddReading` validates and processes business logic
5. **Repository**:  
   → `Repository.AddReading` persists to DB
6. **DB**:  
   → PostgreSQL stores reading

### Other Flows

- **Signup/Login**:  
  → `AuthHandler` → `Service.Signup/VerifyUser` → `Repository.CreateUser/FindUserByEmail`
- **List Meters**:  
  → `MetersHandler.List` → `Service.ListMeters` → `Repository.ListMeters`
- **Usage Query**:  
  → `ReadingsHandler.Usage` → `Service.Usage` → `Repository.UsageSum`

---

## Configuration, Feature Flags, and Secrets

- **Config Source**:  
  - Environment variables, loaded via [`internal/config/config.go`](internal/config/config.go:1)
  - Example: `ADDR`, `DATABASE_URL`, `JWT_SECRET`
  - `.env.local` for local development
- **Feature Flags**:  
  - Not implemented, but could be added via config/env
- **Secrets Boundary**:  
  - JWT secret and DB credentials are only loaded in the backend, not exposed to clients

---

## Error Strategy, Logging, and Metrics

- **Error Handling**:  
  - Errors are wrapped and returned up the stack, with HTTP handlers mapping them to appropriate status codes and messages
  - Sentinel errors (e.g., "invalid credentials") are used for business logic
- **Logging**:  
  - [`zap`](https://github.com/uber-go/zap) logger in [`cmd/server/main.go`](cmd/server/main.go:1)
  - Logs server start, errors, and authentication failures
- **Metrics**:  
  - Not implemented, but can be added via middleware (e.g., Prometheus)

---

## Concurrency Model

- **HTTP Server**:  
  - Each request handled in its own goroutine (Go standard)
- **Database**:  
  - Connection pool managed by `sqlx` with limits set in [`internal/db/db.go`](internal/db/db.go:1)
- **Context Usage**:  
  - All major flows pass `context.Context` for cancellation/timeouts
  - Middleware injects user context for request-scoped data
- **Workers/Queues**:  
  - Not present in current design, but can be added for background jobs

---

## Summary

This service is a modular, layered Go application following best practices for configuration, error handling, and concurrency. It is ready for extension with additional features such as metrics, feature flags, and background processing.