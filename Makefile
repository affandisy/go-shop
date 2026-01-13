.PHONY: help run build test clean install-deps

help: ## Tampilkan bantuan
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

run: ## Jalankan aplikasi
	@echo "Running application..."
	go run cmd/api/main.go

build: ## Build aplikasi
	@echo "Building application..."
	go build -o bin/goshop cmd/api/main.go

test: ## Jalankan test
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Jalankan test dengan coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Bersihkan build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up: ## Start docker compose
	docker-compose up -d

docker-down: ## Stop docker compose
	docker-compose down

docker-logs: ## Lihat docker logs
	docker-compose logs -f