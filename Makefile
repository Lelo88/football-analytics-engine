.PHONY: all ci fmt test build tidy help db-up db-down migrate-up migrate-down migrate-validate

COMPOSE := docker compose --env-file ./.env -f deploy/docker-compose.yml
VALIDATION_DB := migration_validation

.DEFAULT_GOAL := all

all: fmt test build

ci: fmt test build migrate-validate

fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build ./...

tidy:
	go mod tidy

db-up:
	$(COMPOSE) up -d postgres

db-down:
	$(COMPOSE) down

migrate-up: db-up
	@echo "Esperando PostgreSQL..."
	@until docker exec football-analytics-postgres pg_isready -U postgres -d postgres >/dev/null 2>&1; do \
		sleep 1; \
	done
	docker exec football-analytics-postgres psql -U postgres -d postgres -c "DROP DATABASE IF EXISTS $(VALIDATION_DB);"
	docker exec football-analytics-postgres psql -U postgres -d postgres -c "CREATE DATABASE $(VALIDATION_DB);"
	docker exec -i football-analytics-postgres psql -v ON_ERROR_STOP=1 -U postgres -d $(VALIDATION_DB) < migrations/0001_initial_schema.up.sql

migrate-down: db-up
	docker exec -i football-analytics-postgres psql -v ON_ERROR_STOP=1 -U postgres -d $(VALIDATION_DB) < migrations/0001_initial_schema.down.sql

migrate-validate: migrate-up migrate-down
	docker exec football-analytics-postgres psql -U postgres -d $(VALIDATION_DB) -c "SELECT COUNT(*) AS remaining_tables FROM pg_tables WHERE schemaname='public';"

help:
	@echo "Targets disponibles:"
	@echo "  make        -> fmt + test + build"
	@echo "  make ci     -> fmt + test + build + migrate-validate"
	@echo "  make fmt    -> formatea el código"
	@echo "  make test   -> ejecuta tests"
	@echo "  make build  -> compila el proyecto"
	@echo "  make tidy   -> limpia dependencias go.mod/go.sum"
	@echo "  make db-up  -> levanta PostgreSQL"
	@echo "  make db-down -> apaga stack docker"
	@echo "  make migrate-up -> aplica migración up en base temporal"
	@echo "  make migrate-down -> aplica migración down en base temporal"
	@echo "  make migrate-validate -> ejecuta up/down y valida esquema limpio"
