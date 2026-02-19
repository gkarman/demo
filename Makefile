ENV_FILE ?= .env

define run_with_env
	@set -a; source $(ENV_FILE); set +a; $(1)
endef

up:
	docker compose up -d db

down:
	docker compose down

run:
	$(call run_with_env,go run ./cmd/api)

lint:
	$(call run_with_env,golangci-lint run)

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1

migrate-version:
	docker compose run --rm migrate version

migrate-create:
	docker compose run --rm migrate create -ext sql -dir /migrations $(name)
