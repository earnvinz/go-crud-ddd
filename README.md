# Backend Go API with Database Migration

This is a sample backend API project written in Go, featuring database schema management with golang-migrate/migrate and using Docker Compose for running the database and migrations.

---

## Environment Configuration

Create a `.env` file in the project root followed by .example.env

---

## Docker Compose Setup

Navigate to the project root folder then command `docker compose up -d`

---

## Running Migrations

1. Navigate to the `cmd` directory inside your project folder:

2. Run the migration with command `go run migrate.go up` to apply migrations to your database using the environment variables from your .env file

[optional] 3. To rollback (migrate down) the last migration, run: `go run migrate.go down`

---

## Writing Migrations

Create a new migration file with a sequential number and a descriptive name using:
`migrate create -ext sql -dir ./migrations -seq {action_name} ex.add_users_table`

---

## How to start the server

1. Navigate to the project root folder containing `main.go`.
2. Run the command: `go run main.go`

---

## How to run API integration tests

1. Generate Swagger documentation by running the command: `swag init`
2. Open your browser and navigate to: `http://localhost:3000/swagger/index.html` to testing API

---

## How to run unit test

1. Navigate to the project root folder containing `main.go`.
2. command `go test./...`

```

```
