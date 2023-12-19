.PHONY: create-migration up-migrate down-migrate

DSN := sqlite3.db
MIGRATIONS_DIR := sqlite/migrations

create-migration:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

up-migrate:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DSN) up

down-migrate:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DSN) down

run:
	go run cmd/fileserve/main.go
