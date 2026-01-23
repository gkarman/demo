ENV_FILE ?= .env.local
STEP ?=

define run_with_env
	@set -a; source $(ENV_FILE); set +a; $(1)
endef

run:
	$(call run_with_env,go run ./cmd/demo)

lint:
	$(call run_with_env,golangci-lint run)

migrate-up:
	$(call run_with_env,go run ./cmd/migrate up $(if $(STEP),--step=$(STEP)))

migrate-down:
	$(call run_with_env,go run ./cmd/migrate down $(if $(STEP),--step=$(STEP)))

migrate-version:
	$(call run_with_env,go run ./cmd/migrate version)
