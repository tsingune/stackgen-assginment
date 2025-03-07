.PHONY: build run test test-unit test-integration lint clean docker-build docker-run migrate help docker-up docker-down docker-logs docker-ps docker-exec-db start setup-test-db test-in-docker swagger-deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=meeting-scheduler
MAIN_PATH=./cmd/api

# Docker parameters
DOCKER_IMAGE=meeting-scheduler
DOCKER_TAG=latest

# Default target
all: test build

# Build the application
build:
	@echo "Building..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running..."
	$(GORUN) $(MAIN_PATH)

# Set up test database
setup-test-db:
	@echo "Setting up test database..."
	@chmod +x scripts/setup_test_db.sh
	@./scripts/setup_test_db.sh

# Run tests
test: setup-test-db
	@echo "Running all tests..."
	$(GOTEST) -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./...

# Run integration tests only
test-integration: setup-test-db
	@echo "Running integration tests..."
	$(GOTEST) -v -run "TestBasic|TestTimeSlot|TestError" ./internal/api

# Run tests with coverage
test-coverage: setup-test-db
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start Docker Compose environment
docker-up:
	@echo "Starting Docker Compose environment..."
	docker-compose up -d

# Stop Docker Compose environment
docker-down:
	@echo "Stopping Docker Compose environment..."
	docker-compose down

# View Docker Compose logs
docker-logs:
	@echo "Viewing Docker Compose logs..."
	docker-compose logs -f

# List running Docker containers
docker-ps:
	@echo "Listing Docker containers..."
	docker-compose ps

# Execute PostgreSQL command line
docker-exec-db:
	@echo "Connecting to PostgreSQL..."
	docker exec -it meeting-scheduler-db psql -U postgres -d scheduler

# Run database migrations
migrate:
	@echo "Running database migrations..."
	$(GORUN) ./cmd/migrate

# Generate OpenAPI documentation
docs:
	@echo "Generating API documentation..."
	swagger generate spec -o ./docs/swagger.json

# Deploy to Kubernetes
k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/

# Start the application with the startup script
start:
	@echo "Starting application with the startup script..."
	@chmod +x scripts/start.sh
	@./scripts/start.sh

# Run tests inside Docker
test-in-docker:
	@echo "Setting up test database in Docker..."
	docker exec meeting-scheduler-db psql -U postgres -c "DROP DATABASE IF EXISTS scheduler_test;" || true
	docker exec meeting-scheduler-db psql -U postgres -c "CREATE DATABASE scheduler_test;"
	@echo "Running tests inside Docker container..."
	docker exec -e TEST_DB_HOST=postgres -e TEST_DB_PORT=5432 -e TEST_DB_USER=postgres -e TEST_DB_PASSWORD=postgres -e TEST_DB_NAME=scheduler_test meeting-scheduler-app go test -v ./...

# Install Swagger dependencies
swagger-deps:
	@echo "Installing Swagger dependencies..."
	go get -u github.com/swaggo/http-swagger
	go get -u github.com/swaggo/swag/cmd/swag
	go mod tidy

# Help command
help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run all tests"
	@echo "  make test-unit          - Run unit tests"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make lint               - Run linter"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make deps               - Download dependencies"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make docker-up          - Start Docker Compose environment"
	@echo "  make docker-down        - Stop Docker Compose environment"
	@echo "  make docker-logs        - View Docker Compose logs"
	@echo "  make docker-ps          - List running Docker containers"
	@echo "  make docker-exec-db      - Connect to PostgreSQL command line"
	@echo "  make migrate            - Run database migrations"
	@echo "  make docs               - Generate OpenAPI documentation"
	@echo "  make k8s-deploy         - Deploy to Kubernetes"
	@echo "  make start              - Start application with diagnostic checks"
	@echo "  make setup-test-db      - Set up test database"
	@echo "  make test-in-docker     - Run tests inside Docker container"
	@echo "  make swagger-deps         - Install Swagger dependencies"
	@echo "  make help               - Show this help message" 