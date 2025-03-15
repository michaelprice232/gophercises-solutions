# WIP - Gophercises - Exercise 7 - Task

My implementation of https://github.com/gophercises/task/

Opted to use Postgres instead of BoltDB as it's more relevant to my work-life. 
Also implemented DB migrations and unit tests as well to further my knowledge.

## Useful links

- https://github.com/spf13/cobra
- https://github.com/jackc/pgx
- https://github.com/golang-migrate/migrate



## DB Migrations

### Create the migration files

```shell
migrate create -ext sql -dir migrations -seq -digits 4 add_tasks_table
```

### Run the migrations

```shell
# Provision
migrate -path migrations -database 'postgres://postgres:test@localhost:5432/task_db?sslmode=disable' up

# Tear down
migrate -path migrations -database 'postgres://postgres:test@localhost:5432/task_db?sslmode=disable' down --all
```