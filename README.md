# Go Postgres Health Check API

A minimal, production-ready Go microservice that exposes **liveness** and **readiness** health endpoints backed by PostgreSQL.

This project follows a clean folder structure and is designed to be compatible with **Docker**, **Docker Compose**, and **Kubernetes-style health checks**.

---

## Features

- Liveness endpoint (`/live`) — checks if the process is running
- Readiness endpoint (`/ready`) — checks PostgreSQL connectivity
- Optional `/health` alias for readiness
- Graceful shutdown support
- Connection-pooled PostgreSQL access
- Docker & Docker Compose support
- Clean, scalable project structure

---

## Project Structure

go-postgres-micro/
├── cmd/
│ └── api/
│ └── main.go # Application entry point
├── internal/
│ ├── config/
│ │ └── config.go # Environment configuration
│ ├── db/
│ │ └── db.go # Database connection logic
│ └── http/
│ └── handlers/
│ └── health.go # Health check handlers
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum

---
## Health Endpoints

### `GET /live`
Checks whether the application process is running.

**Response**
```json
{
  "status": "ok",
  "time": "2026-02-01T11:52:26+06:00"
}
Used for:

Liveness probes

Load balancer basic checks

GET /ready
Checks whether the application is ready to serve traffic by verifying PostgreSQL connectivity.

Response (DB OK)

{
  "status": "ok",
  "time": "2026-02-01T11:52:34+06:00",
  "db": "ok"
}
Response (DB Down)

HTTP 503 Service Unavailable

{
  "status": "down",
  "time": "2026-02-01T11:52:34+06:00",
  "db": "down"
}
GET /health
Alias of /ready for simpler monitoring setups.

Configuration
The application is configured using environment variables.

Variable	Description	Required
PORT	HTTP server port (default: 8080)	No
DATABASE_URL	PostgreSQL connection string	Yes
Example:

DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
Running Locally (Without Docker)
Start PostgreSQL separately, then:

export PORT=8080
export DATABASE_URL="postgres://app:app123@localhost:5432/appdb?sslmode=disable"

go mod tidy
go run ./cmd/api
Running with Docker Compose
docker compose up --build
Services:

API: http://localhost:8080

PostgreSQL: localhost:5432

Graceful Shutdown
The service listens for SIGINT and SIGTERM signals and shuts down gracefully, allowing in-flight requests to complete before exiting.

Use Cases

Kubernetes liveness & readiness probes

Load balancer health checks

Infrastructure monitoring

Boilerplate for Go microservices