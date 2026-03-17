# Build configuration
BINARY_NAME=logger-app
BUILD_DIR=.build

.PHONY: all build clean run seed

all: build

build:
	@echo "Building application and scripts..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@go build -o $(BUILD_DIR)/user-manager scripts/user_manager/main.go
	@go build -o $(BUILD_DIR)/seeder scripts/seeder/main.go
	@cp -r static $(BUILD_DIR)/ 2>/dev/null || true
	@echo "Build complete. Binaries and assets are in $(BUILD_DIR)/"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

manage-user:
	@if [ ! -f $(BUILD_DIR)/user-manager ]; then $(MAKE) build; fi
	@./$(BUILD_DIR)/user-manager $(ARGS)

seed:
	@if [ ! -f $(BUILD_DIR)/seeder ]; then $(MAKE) build; fi
	@./$(BUILD_DIR)/seeder $(ARGS)
