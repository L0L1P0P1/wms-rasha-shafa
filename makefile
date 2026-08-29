-include .env
export

DRIVER ?= postgres
PGPORT ?= 5432
DB_STRING ?= postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=disable
MIGRATIONS_DIR ?= db/migrations

GOOSE_CMD = go tool goose
SQLC_CMD = go tool sqlc

.PHONY: migrate-status migrate-up migrate-down create-migration generate

# goose migration commands

migrate-status:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_STRING)" status

migrate-up:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_STRING)" up

migrate-down:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DB_STRING)" down

create-migration:
	@read -p "Enter migration name (e.g., add_phone_number): " name; \
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) create $$name sql

# sqlc code generation

generate:
	@echo "Generating Go code from SQL queries..."
	$(SQLC_CMD) generate
	@echo "Done!"
