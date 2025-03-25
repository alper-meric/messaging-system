.PHONY: build run test clean docker-build docker-run mockery

# Variables
APP_NAME = messaging-system
BUILD_DIR = build
MAIN_PATH = ./cmd/server

# Build
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)"

# Run
run: build
	@echo "Running application..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Test
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean
clean:
	@echo "Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaning completed."

# Generate mocks using mockery
mockery:
	@echo "Checking if mockery is installed..."
	@which mockery > /dev/null || (echo "Installing mockery..." && go install github.com/vektra/mockery/v2@latest)
	@echo "Generating mocks..."
	@mockery --dir=repository --name=MessageRepository --output=./mocks/repository --outpkg=repository
	@mockery --dir=repository --name=CacheRepository --output=./mocks/repository --outpkg=repository
	@mockery --dir=services --name=MessageServiceInterface --output=./mocks/services --outpkg=services
	@echo "Mock generation completed successfully."

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .
	@echo "Docker image created: $(APP_NAME)"

# Run with Docker
docker-run: docker-build
	@echo "Running application with Docker..."
	@docker run -p 8080:8080 --name $(APP_NAME)-container $(APP_NAME)

# Run with Docker Compose
docker-compose-up:
	@echo "Starting applications with Docker Compose..."
	@docker-compose up -d
	@echo "Application started, you can access it at http://localhost:8080"

# Stop Docker Compose
docker-compose-down:
	@echo "Stopping applications with Docker Compose..."
	@docker-compose down 