run:
	@set -a; source .env.local; set +a; go run ./cmd/demo

lint:
	@golangci-lint run