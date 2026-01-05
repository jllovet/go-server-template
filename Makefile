# Variables
BINARY_NAME=server
BUILD_DIR=bin
MAIN_PATH=cmd/main.go
TEST_DB_URL=postgres://user:password@localhost:5432/todo_db?sslmode=disable

.PHONY: all build run run-silent test clean vet tidy benchmark

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-db:
	@echo "Building database..."
	@docker compose up -d

run: build build-db
	@echo "Running..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

run-silent: build
	@echo "Running silently..."
	@DISABLE_LOGGING=true ./$(BUILD_DIR)/$(BINARY_NAME)

test: build-db
	@echo "Testing..."
	@TEST_DATABASE_URL="$(TEST_DB_URL)" go test -race ./...

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@docker compose down -v

vet:
	@go vet ./...

tidy:
	@go mod tidy