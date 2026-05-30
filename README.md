# SS-CatalogService

## Overview

SS-CatalogService is a core microservice responsible for managing the product catalog and taxonomy data for the SamStore e-commerce platform. Built with Go 1.26+, the service relies on high-performance frameworks and integrations to handle product searching, indexing, pricing, inventory updates, and multi-variant taxonomy.

The service connects with downstream components via RabbitMQ, ensuring that changes to the catalog trigger corresponding updates across other microservices, and utilizes PostgreSQL for data persistence.

## Features

- **Product Taxonomy & Catalog Management**: Extensive support for products, brands, categories, attributes, options, and multi-level product variants.
- **Pricing & Stock Controls**: Handles product pricing, inventory tracking, stock adjustments, and promotional pricing rules.
- **High-Performance Search & Indexing**: Direct integration with Meilisearch for ultra-fast full-text product search and updates.
- **In-Memory Caching**: Implements local caching with Bigcache to reduce PostgreSQL query load for high-traffic read paths.
- **Event-Driven Microservices Sync**: Transactional outbox event publishing via RabbitMQ to keep downstream services up to date.
- **Database Migrations & Data Seeding**: Versioned schema migrations using golang-migrate and an automated data seed utility.

## Tech Stack

| Category       | Technology                             |
| -------------- | -------------------------------------- |
| Language       | Go (version 1.26.2)                    |
| Web Framework  | Gin-Gonic                              |
| Database / ORM | PostgreSQL / GORM                      |
| Search Engine  | Meilisearch (meilisearch-go)           |
| Caching        | Bigcache (allegro/bigcache)            |
| Message Broker | RabbitMQ (amqp091-go)                  |
| Telemetry      | OpenTelemetry with Gin instrumentation |

## Project Structure

```text
SS-CatalogService/
├── cmd/
│   ├── api/                  # HTTP API server entry point (main.go)
│   └── seed/                 # Catalog initial data seeding utility (main.go)
├── config/                   # Configuration parsing and environment loading (config.go)
├── db/                       # SQL migrations and schemas
│   └── migrations/           # Versioned SQL migration scripts
├── internal/                 # Core private application codebase
│   ├── delivery/             # Presentation layer (HTTP / Gin handlers)
│   ├── domain/               # Business entities and repository/usecase interfaces
│   ├── infrastructure/       # External drivers (cache, DB initialization, messaging, telemetry)
│   ├── repository/           # GORM repository implementation
│   └── usecase/              # Core business application logic
└── pkg/                      # Public helper utilities (custom logger, responses)
```

## Requirements

- Go 1.26+
- PostgreSQL
- Meilisearch
- RabbitMQ

## Installation

```bash
git clone <repository>
cd SamStore/SS-CatalogService
```

Download the required Go modules:

```bash
go mod download
```

## Configuration

The service maps environment variables via a local `.env` file. Key properties include:

```env
APP_PORT=                 # HTTP Port for Gin web server (default: 8081)
DB_DSN=                   # PostgreSQL DSN (e.g. host=localhost user=postgres password=123456 dbname=ss_catalog_db port=5432 sslmode=disable)
RABBITMQ_URL=             # RabbitMQ broker connection URL
MEILISEARCH_HOST=         # Meilisearch server host URL
MEILISEARCH_API_KEY=      # Meilisearch API key
JWT_SECRET=               # Secret key used for local token verification
ENVIRONMENT=              # Execution stage (e.g., development, production)
```

## Running Locally

To run the REST API server:

```bash
go run cmd/api/main.go
```

To run the initial database data seeding script:

```bash
go run cmd/seed/main.go
```

## Build

To compile the service binary:

```bash
go build -o bin/catalog-service cmd/api/main.go
```

Or build the containerized Docker image:

```bash
docker build -t ss-catalog-service .
```

## Testing

```bash
go test ./...
```

## API Documentation

Endpoint actions are exposed via the Gin HTTP server:

| Method | Endpoint | Description                                                      |
| ------ | -------- | ---------------------------------------------------------------- |
| GET    | /health  | Check if the service and its backing database/broker are online  |
| GET    | /metrics | Prometheus compatible metrics endpoint exposed via OpenTelemetry |

## Database

- **Database Type**: PostgreSQL.
- **ORM**: GORM.
- **Migrations**: Executed via versioned SQL migrations in `db/migrations/` using `golang-migrate`.
- **Seed Data**: Populated using the custom seeding application inside `cmd/seed/`.

## Deployment

- **Docker**: Deployed using a multi-stage [Dockerfile](Dockerfile) (builder stage compiles the Go binary, runner stage copies it to a lightweight Alpine container).
- **Docker Compose**: Handled alongside other microservices under the orchestrator configurations.

## Architecture Notes

- **Clean Hexagonal Architecture**: Strict boundary separation. Domain models are isolated from GORM model mapping and transport layers.
- **Event-Driven Outbox**: Outbox pattern guarantees that catalog updates (like stock additions or pricing changes) are reliably broadcasted to RabbitMQ without risking partial failure.

## Known Issues

Not identified from source code.

## Future Improvements

- Add automatic migration runs on application startup.
- Incorporate distributed caching (like Redis) alongside the existing local Bigcache in-memory system.

## License

```text
License information not specified.
```
