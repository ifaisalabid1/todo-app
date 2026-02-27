.PHONY: build run clean migrate help

include .env

# Variables
APP_NAME=todo-api
BUILD_DIR=./bin
MAIN_PATH=./cmd/api

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete!"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)/main.go

# Run tests

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Run database migrations
migrate:
	@echo "Running migrations..."
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path ./migrations up
	@echo "Migrations complete!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete!"

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run
	@echo "Linting complete!"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied!"

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make migrate        - Run database migrations"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"
	@echo "  make tidy           - Tidy dependencies"
	@echo "  make help           - Show this help message"