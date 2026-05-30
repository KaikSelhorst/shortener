ifneq (,$(wildcard .env))
  include .env
  export
endif

dev:
	cd api && go run cmd/server/main.go

seed:
	cd api && go run ./cmd/seed

fmt:
	cd api && gofmt -w .

lint:
	cd api && golangci-lint run ./...


test:
	cd api && go test ./internal/service/... ./internal/handler/...

test-integration:
	cd api && DATABASE_URL="$(DATABASE_URL)" go test -tags integration ./internal/repository/...

test-all: test test-integration

docker-up:
	docker compose up -d

docker-down:
	docker compose down
