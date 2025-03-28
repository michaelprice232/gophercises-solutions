# Gophercises - Exercise 7 - Task

A simple Task todo CLI. My implementation of https://github.com/gophercises/task/.
Opted to use Postgres instead of BoltDB as it's more relevant to my work-life.

## Running Locally

Pre-reqs:

- [golang-migrate](https://github.com/golang-migrate/migrate) installed (`brew install golang-migrate`)
- [Docker/Docker-Compose](https://www.docker.com/products/docker-desktop/) installed
- [Golang](https://go.dev/dl/) install

```shell
# Start database
make run_db

# Apply DB migrations
migrate -path migrations -database 'postgres://postgres:test@localhost:5432/task_db?sslmode=disable' up

# Install Go modules locally
go mod download                      

# Run CLI
export DATABASE_URL="postgres://postgres:test@localhost:5432/task_db?sslmode=disable"
go run cmd/main.go add [task name]  # Add task
go run cmd/main.go do [task name]   # Complete a task
go run cmd/main.go outstanding      # List the outstanding tasks
go run cmd/main.go completed        # List the completed tasks
```

#### Cleanup

```shell
make down_db
```

## DB Migrations

### Create the migration files

```shell
migrate create -ext sql -dir migrations -seq -digits 4 add_tasks_table
```

### Run the migrations

```shell
# Provision
make db_migrations_up

# Tear down
make db_migrations_down
```

## Useful links

- https://github.com/spf13/cobra
- https://github.com/jackc/pgx
- https://github.com/golang-migrate/migrate
- https://github.com/gophercises/task/