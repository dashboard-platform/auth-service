# Auth Service

![coverage](https://img.shields.io/badge/coverage-35.4%25-orange)

The `auth-service` is a microservice responsible for handling user registration, login, and identity verification within a distributed system. It generates and validates JWT tokens for authentication and communicates via HTTP.

## Features

- User registration with hashed password (Argon2ID)
- Secure login flow with JWT token generation
- JWT validation for authenticated requests
- Healthcheck endpoint
- Designed for integration behind an API Gateway

## Endpoints

| Method | Path         | Auth Required | Description                       |
|--------|--------------|----------------|-----------------------------------|
| GET    | `/healthcheck` | ❌             | Basic service status              |
| POST   | `/register`    | ❌             | Create new user account           |
| POST   | `/login`       | ❌             | Authenticate user and return JWT  |
| GET    | `/me`          | ✅             | Return current user's information |

> ⚠️ The `/logout` logic should be handled by the API Gateway, not this service.

## Environment Variables

| Variable     | Description                             |
|--------------|-----------------------------------------|
| `PORT`       | Port on which the service runs          |
| `DB_URL`     | `"host=localhost user=postgres password=password dbname=db port=5432 sslmode=disable"`|
| `JWT_SECRET` | Secret used for signing JWTs            |
| `ENV`        | Application environment (`dev`, `prod`) |

`.env` file example:
```env
   PORT=8080
   DB_URL=host=localhost user=postgres password=secret dbname=authdb port=5432
   JWT_SECRET=your-super-secret
   ENV=dev
```

## Run Locally

### Option 1: Run with Go (requires PostgreSQL running locally)

1. Create a `.env` file.
2. Start a local PostgreSQL instance with database `authdb`.
3. Run:

```bash
go run cmd/main.go
```

### Option 2: Run with Docker Compose (recommended for testing)
```bash
docker-compose up --build
```

This will:
- Start a PostgreSQL container with `authdb`
- Start the auth-service on port `8081`
- Use default credentials from the compose environment

Access healthcheck:
```bash
curl http://localhost:8081/healthcheck
```
