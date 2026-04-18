ifneq (,$(wildcard .env))
  include .env
  export
endif

dev:
	cd api && go run cmd/server/main.go

fmt:
	cd api && gofmt -w .

lint:
	cd api && golangci-lint run ./...


docker-up:
	docker compose up -d

docker-down:
	docker compose down
