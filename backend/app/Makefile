.PHONY: run build test test-cover docker-up docker-down docker-logs lint tidy

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o dose-certa-backend ./cmd/server

test:
	go test ./... -v

test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | tail -n1

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

tidy:
	go mod tidy

lint:
	golangci-lint run ./...
