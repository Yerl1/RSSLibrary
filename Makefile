include .env

export $(shell sed 's/=.*//' .env)

MIGRATIONS_DIR=./migrations
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
DOCKER_MIGRATE=docker run --rm -v $(shell pwd)/$(MIGRATIONS_DIR):/migrations:Z --network rsslibrary_default migrate/migrate

# ---------------------
# MIGRATIONS
# ---------------------

.PHONY: migrate-up migrate-down migrate-force migrate-version

migrate-up:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_URL)" up

migrate-down:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_URL)" down 1

migrate-force:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_URL)" force

migrate-version:
	$(DOCKER_MIGRATE) -path=/migrations -database "$(DB_URL)" version

# ---------------------
# DOCKER COMPOSE
# ---------------------

.PHONY: up down restart logs rebuild

up:
	docker compose up --build

down:
	docker compose down

restart:
	docker compose down && docker compose up --build

logs:
	docker compose logs -f --tail=100

rebuild:
	docker compose build app

db:
	docker compose exec db psql -U $(DB_USER) -d $(DB_NAME)


# ---------------------
# SHORTCUTS
# ---------------------

ps:
	docker compose ps

bash:
	docker compose exec app /bin/sh
