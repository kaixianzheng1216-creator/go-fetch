# go-fetch

Go analytics service with a React frontend. The backend owns the API, storage, migrations, OpenAPI generation, and embedded production web assets.

## Project Layout

- `cmd/server`: HTTP server entrypoint.
- `cmd/openapi`: OpenAPI JSON generator.
- `internal/server`: routes, handlers, OpenAPI configuration, and frontend serving.
- `internal/middleware`: HTTP and API middleware.
- `internal/session`: session manager configuration.
- `internal/domain`: application types and business rules, split by domain.
- `internal/store`: PostgreSQL access, sqlc generated code, and migrations.
- `internal/collector`: request parsing and analytics event construction.
- `internal/static`: embedded tracking script and production frontend build output.
- `frontend`: Vite React application. It has its own empty `go.mod` so root-level Go commands do not scan `node_modules`.
- `api`: generated OpenAPI artifact.

## Development

Start PostgreSQL:

```sh
docker compose up -d postgres
```

Run the backend:

```sh
go run ./cmd/server
```

Run the frontend dev server:

```sh
npm --prefix frontend install
npm --prefix frontend run dev
```

Useful checks:

```sh
go test ./...
go vet ./...
npm --prefix frontend run lint
npm --prefix frontend run build
```

On systems with `make`, the common workflow is:

```sh
make check
```

## Code Generation

Regenerate the backend SQL and OpenAPI outputs:

```sh
go generate ./...
```

Regenerate frontend API types after OpenAPI changes:

```sh
npm --prefix frontend run api:generate
```
