.PHONY: build run test clean proto migrate docker-build docker-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=grpc-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product/product.proto

# Run database migrations
migrate:
	@echo "Running database migrations..."
	$(GOBUILD) -o migrate-tool ./cmd/migrate
	./migrate-tool
	rm -f migrate-tool

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Docker commands
docker-build:
	docker build -t grpc-exmpl -f deployments/docker/Dockerfile .

docker-run:
	docker run -p 8080:8080 --env-file .env grpc-exmpl

# Development setup
dev-setup: deps proto
	@echo "Development setup complete"

# Start development server with hot reload (requires air)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Security scan (requires gosec)
security:
	gosec ./...

# Generate mocks (requires mockgen)
mocks:
	mockgen -source=internal/repository/user_repository.go -destination=tests/mocks/user_repository_mock.go
	mockgen -source=internal/service/user_service.go -destination=tests/mocks/user_service_mock.go

# Run all checks
check: fmt lint test security

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build files"
	@echo "  proto        - Generate protobuf files"
	@echo "  migrate      - Run database migrations"
	@echo "  deps         - Download dependencies"
	@echo "  dev-setup    - Setup development environment"
	@echo "  dev          - Start development server with hot reload"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  security     - Run security scan"
	@echo "  mocks        - Generate mocks"
	@echo "  check        - Run all checks"