# Build configuration
BINARY_NAME=logger-app
BUILD_DIR=.build

.PHONY: all build clean run seed

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

seed:
	@echo "Seeding database..."
	@go run scripts/seed.go
