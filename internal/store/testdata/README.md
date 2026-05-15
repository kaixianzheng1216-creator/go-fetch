# Store Test Data

This directory is reserved for future PostgreSQL integration test fixtures.

Conventions:

- Keep default `go test ./...` fast and self-contained.
- Put database-backed tests behind an explicit integration setup or build tag.
- Prefer small SQL fixtures that describe one scenario per file.
- Do not store generated sqlc output or production data here.
