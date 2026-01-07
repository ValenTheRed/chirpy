# chirpy
Project to get my feet wet with BE development

## Requirement
1. Database used: postgreSQL.
2. Go tools used
    1. `sqlc` as a postgreSQL ORM and
    2. `goose` for database migration.
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Installation
1. Create an `.env` file to hold these secrets
    1. `DB_URL`: connection string for your postgreSQL instance. It would look something like this: `postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable`.
    2. `PLATFORM`: `DEV`, `PROD`, etc.
    3. `TOKEN_SECRET`: 256 bit secret for signing and verifying of JWTs. Below can be used to generate it.
    ```bash
    openssl rand -base64 64
    ```
    4. `POLKA_KEY`: API key for authenticate a webhook/an external caller of our server. Provided in the course.
2. `goose` migrate all schemas
```bash
cd ./sql/schema
goose $DB_URL up # as many times as there are files
```
3. Generate go bindings for all SQL queries in `./sql/queries`.
```bash
sqlc generate
```

## Run
```
go run .
```
