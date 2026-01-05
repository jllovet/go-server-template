# Variables
BINARY_NAME=server
BUILD_DIR=bin
MAIN_PATH=cmd/main.go

.PHONY: all build run run-silent test clean vet tidy benchmark

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: build
	@echo "Running..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

run-silent: build
	@echo "Running silently..."
	@DISABLE_LOGGING=true ./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Testing..."
	@go test -race ./...

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

vet:
	@go vet ./...

tidy:
	@go mod tidy