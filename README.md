# Auth Service

![coverage](https://img.shields.io/badge/coverage-80.5%25-green)

The `auth-service` is a microservice responsible for handling user registration, login, and identity verification within a distributed system. It generates and validates JWT tokens for authentication and communicates via HTTP.

```bash
git clone https://github.com/dashboard-platform/auth-service.git
cd auth-service
```

## Run Locally

### Option 1: Run with Go (requires PostgreSQL running locally)

1. Create a `.env` file.
2. Start a local PostgreSQL instance with database `authdb`.
3. Run:

```bash
go run cmd/main.go
```

`.env` file example: (for details, see below)
```env
   PORT=:8080
   DB_URL=host=localhost user=postgres password=secret dbname=authdb port=5432
   JWT_SECRET=your-super-secret
   ENV=dev
```

### Option 2: Run with Docker Compose (recommended for testing)

You can add `-d` flag to run the application in `detached` mode.
```bash
docker-compose up (-d)
```

This will:
- Start a PostgreSQL container with `authdb`
- Start the auth-service on port `8081`
- Use default credentials from the compose environment

Access healthcheck:
```bash
curl http://localhost:8081/healthcheck
```

## Run Tests

To run the whole test suite, use:
```bash
cd auth-service
go tests -v ./...
```

## Environment Variables

| Variable     | Description                             |
|--------------|-----------------------------------------|
| `PORT`       | Port on which the service runs (`:8080`)          |
| `DB_URL`     | PostgreSQL database connection URL (`"host=localhost user=postgres password=password dbname=db port=5432 sslmode=disable"`)|
| `JWT_SECRET` | Secret used for signing JWTs (`secret`)        |
| `ENV`        | Application environment (`dev`/``prod``) |

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
