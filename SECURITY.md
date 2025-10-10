# SECURITY.md

## Authentication & Authorization

- **Authentication**:  
  - Uses JWT (JSON Web Tokens) for stateless authentication.
  - Users log in with email and password; on success, a signed JWT is issued.
- **Authorization**:  
  - JWT is required for all protected endpoints.
  - User identity is extracted from the token and enforced in handlers.

## Token Flow

- Clients obtain a JWT via `/v1/auth/login`.
- JWT must be sent in the `Authorization: Bearer <token>` header for API requests.
- Tokens are signed with a server-side secret and validated on each request.

## PII Boundaries

- Personally Identifiable Information (PII) includes user email addresses and meter data.
- PII is only accessible to authenticated users and is not exposed to other users.
- PII is stored in the PostgreSQL database and not logged.

## Data Retention

- No automated data retention or deletion policy is implemented.
- Data persists in the database until manually deleted or purged by an admin.

## Encryption

- **In Transit**:  
  - All API traffic should be served over HTTPS in production to ensure encryption in transit.
- **At Rest**:  
  - Data is stored in PostgreSQL; encryption at rest depends on DB/server configuration (not handled by app).

---