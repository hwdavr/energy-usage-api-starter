# OPERATIONS.md

## Running Locally

- **Prerequisites**:  
  - Go (latest stable)
  - PostgreSQL (local or Docker)
  - [migrate](https://github.com/golang-migrate/migrate) CLI for DB migrations

- **Environment Setup**:  
  - Copy `.env.local` and set `DATABASE_URL`, `JWT_SECRET`, and `ADDR` as needed.

- **Start the Server**:  
  - Use the Makefile:
    ```sh
    make run
    ```
    This runs the server with environment variables from `.env.local`.

- **Database Migrations**:  
  - Apply migrations:
    ```sh
    make migrate-up
    ```
  - Roll back last migration:
    ```sh
    make migrate-down
    ```

- **Seed Data**:  
  - No dedicated seed script provided. Insert data manually via SQL or add a custom Makefile target as needed.

## Health Checks

- **No explicit health check endpoint implemented.**  
  - To add: create `/healthz` endpoint in the router for readiness/liveness checks.

## Observability

- **Logging**:  
  - Uses [`zap`](https://github.com/uber-go/zap) for structured JSON logs.
  - Logs include server start, errors, and authentication failures.

- **Metrics**:  
  - No metrics endpoint implemented.
  - To add: integrate Prometheus middleware and expose `/metrics`.

- **Tracing**:  
  - No distributed tracing or trace headers implemented.
  - To add: propagate trace headers (e.g., `X-Request-Id`, `traceparent`) via middleware.

## Migration Policy

- **Tool**:  
  - Uses [`migrate`](https://github.com/golang-migrate/migrate) for schema migrations.
  - Migration files are in the [`migrations/`](migrations/) directory.

- **Policy**:  
  - All schema changes must be made via migration files.
  - Apply migrations before running the server in any environment.

## Versioning & Release Process

- **Versioning**:  
  - No explicit versioning in codebase.  
  - To add: use Git tags or a `VERSION` file.

- **Release Process**:  
  - Manual: build and deploy from main branch.
  - To automate: add CI/CD pipeline for build, test, and deploy.

---