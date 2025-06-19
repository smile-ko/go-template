# Go Template Microservice

## Introduction

This is a microservice template for Golang, designed to help you quickly bootstrap projects with common components: HTTP (Fiber), gRPC, Kafka, Postgres, logging, migration, swagger, and docker support.

## Project Structure

```
├── cmd/server/           # Application entrypoint
├── config/               # Application configuration (env variables)
├── internal/
│   ├── app/              # Service initialization, graceful shutdown
│   ├── controller/       # HTTP/gRPC handlers, routers, middleware
│   ├── entity/           # Data models
│   ├── repo/             # Data repositories
│   └── service/          # Business logic
├── pkg/                  # Utility packages (grpc, httpserver, kafka, logger, postgres)
├── docs/
│   ├── proto/            # gRPC API definitions (proto, generated files)
│   └── swagger/          # HTTP API definitions (swagger yaml/json)
├── migrations/           # Database migration files
├── Dockerfile            # Docker build file
├── docker-compose.yaml   # Compose file for supporting services (Postgres, ...)
├── Makefile              # Build, run, and generate commands
```

## Main Features

- **HTTP API**:
  - `GET /api/v1/healthz` — Health check endpoint
  - `GET /api/v1/hello?name=abc` — Returns "Hello abc"
- **gRPC API**:
  - `UserService` with RPCs: `GetUserById`, `GetUsersByIds`, `CreateUser`
- **Kafka**: Consumer for demo-topic (prints received messages)
- **Postgres**: Connects via environment variables
- **Swagger**: API documentation for dev environment at `/v1/swagger/index.html`
- **Logging**: Configurable by environment
- **Migration**: SQL migration support

## Getting Started

### 1. Prerequisites

- Go 1.19+
- Docker, docker-compose (if using containers)
- Postgres, Kafka (for full feature set)

### 2. Configuration

Create a `.env` file in the project root, for example:

```
APP_NAME=go-template
APP_VERSION=1.0.0
ENV_NAME=dev
HTTP_PORT=8000
LOG_LEVEL=debug
PG_POOL_MAX=10
PG_URL=postgres://postgres:postgres@localhost:5432/go_template?sslmode=disable
GRPC_PORT=9000
METRICS_ENABLED=true
SWAGGER_ENABLED=true
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=go-template-group
```

### 3. Run the Application

- Build proto files:
  ```sh
  make proto-user
  ```
- Run locally:
  ```sh
  make run
  ```
- Start Postgres with docker-compose:
  ```sh
  docker-compose up -d
  ```

### 4. Database Migration

- Place your SQL files in the `migrations/` directory
- Run migration using your preferred tool (e.g. [golang-migrate](https://github.com/golang-migrate/migrate))

### 5. Swagger

- Access swagger docs at: `http://localhost:8000/v1/swagger/index.html` (enabled only if ENV_NAME=dev and SWAGGER_ENABLED=true)

## Contribution

- Fork, create a branch, and submit a PR as with any open-source project.

## License

MIT
