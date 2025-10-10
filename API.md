# API.md

## OpenAPI Specification

This service exposes a RESTful API.  
**No OpenAPI/Swagger YAML is present in the repository.**  
Below is an OpenAPI 3.0 YAML specification inferred from the codebase.

---

```yaml
openapi: 3.0.0
info:
  title: Energy Usage API
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /v1/auth/signup:
    post:
      summary: User signup
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '409':
          description: Conflict (user exists)
  /v1/auth/login:
    post:
      summary: User login
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: JWT token
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
        '401':
          description: Invalid credentials
  /v1/meters:
    get:
      summary: List meters for user
      security:
        - bearerAuth: []
      responses:
        '200':
          description: List of meters
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Meter'
    post:
      summary: Create a new meter
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                label:
                  type: string
      responses:
        '200':
          description: Meter created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Meter'
  /v1/meters/{meterID}/readings:
    post:
      summary: Add a reading to a meter
      security:
        - bearerAuth: []
      parameters:
        - name: meterID
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                reading_kwh:
                  type: number
                reading_at:
                  type: string
                  format: date-time
      responses:
        '200':
          description: Reading added
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Reading'
  /v1/meters/{meterID}/usage:
    get:
      summary: Get usage for a meter
      security:
        - bearerAuth: []
      parameters:
        - name: meterID
          in: path
          required: true
          schema:
            type: string
        - name: from
          in: query
          required: true
          schema:
            type: string
            format: date-time
        - name: to
          in: query
          required: true
          schema:
            type: string
            format: date-time
      responses:
        '200':
          description: Usage summary
          content:
            application/json:
              schema:
                type: object
                properties:
                  meter_id:
                    type: string
                  kwh:
                    type: number
                  from:
                    type: string
                    format: date-time
                  to:
                    type: string
                    format: date-time
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
        email:
          type: string
        created_at:
          type: string
          format: date-time
    Meter:
      type: object
      properties:
        id:
          type: string
        user_id:
          type: string
        label:
          type: string
        created_at:
          type: string
          format: date-time
    Reading:
      type: object
      properties:
        id:
          type: string
        meter_id:
          type: string
        reading_kwh:
          type: number
        reading_at:
          type: string
          format: date-time
        created_at:
          type: string
          format: date-time
```

---

## Other API Types

- **GraphQL**: Not used.
- **gRPC**: Not used.
