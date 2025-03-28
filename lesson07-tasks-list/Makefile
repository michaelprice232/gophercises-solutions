run_db:
	docker-compose up -d --build

down_db:
	docker-compose down --volumes

db_migrations_up:
	migrate -path migrations -database 'postgres://postgres:test@localhost:5432/task_db?sslmode=disable' up

db_migrations_down:
	migrate -path migrations -database 'postgres://postgres:test@localhost:5432/task_db?sslmode=disable' down --all