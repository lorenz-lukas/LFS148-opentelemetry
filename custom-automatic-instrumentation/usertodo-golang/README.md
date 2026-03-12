# Gin Hexagonal API

This project exposes a small Gin API using a hexagonal architecture, PostgreSQL, and GORM. The Go API manages users and user-task relationships, while the Java service owns the `todo` table.

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

Swagger docs are available at `http://localhost:8080/swagger`, and the raw spec is served at `http://localhost:8080/swagger/doc.yaml`.

## Docker

Build the image:

```bash
docker build -t usertodo-golang .
```

Run the container against a PostgreSQL instance:

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_DSN="host=host.docker.internal user=postgres password=postgres dbname=todos port=5432 sslmode=disable" \
  usertodo-golang
```

## Swagger

The repository includes a checked-in Swagger 2.0 document used by the runtime docs endpoint:

```bash
curl http://localhost:8080/swagger/doc.yaml
```

To regenerate the spec with `go-swagger` after changing annotations:

```bash
go run github.com/go-swagger/go-swagger/cmd/swagger@v0.33.1 generate spec \
  --scan-models \
  -o ./internal/adapters/http/swagger.generated.json
```

## Endpoints

List todos:

```bash
curl http://localhost:8080/todos
```

Create a user:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Lukas","email":"lukas@example.com"}'
```

Associate an existing todo to an existing user:

```bash
curl -X POST http://localhost:8080/users/1/tasks/1
```
