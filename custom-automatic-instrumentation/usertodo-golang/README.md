# Gin Hexagonal API

This project exposes a small Gin API using a hexagonal architecture, PostgreSQL, and GORM.

## Structure

- `cmd/api`: application entrypoint
- `internal/domain`: business entities
- `internal/application`: use cases and ports
- `internal/adapters`: HTTP and PostgreSQL adapters
- `internal/infrastructure`: config and database bootstrap

## Run

1. Start the full stack with Docker:

```bash
docker compose up --build
```

2. Or run only PostgreSQL in Docker and start the API locally:

```bash
docker compose up -d postgres
go mod tidy
go run ./cmd/api
```

The API uses `DATABASE_DSN` and `PORT`, with defaults defined in `.env.example`.
Database schema is managed by GORM auto-migration during startup.

## Docker

Build the image:

```bash
docker build -t usertodo-golang .
```

Run the container against a PostgreSQL instance:

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_DSN="host=host.docker.internal user=postgres password=postgres dbname=hexagonal_api port=5432 sslmode=disable" \
  usertodo-golang
```

## Endpoints

Create a task:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Prepare interview","description":"Review hexagonal architecture"}'
```

Create a user and associate the user to existing tasks:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Lukas","email":"lukas@example.com","task_ids":[1]}'
```
