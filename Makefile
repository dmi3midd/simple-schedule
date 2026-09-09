# Run the application locally
run:
	@echo "Running..."
	@go run ./cmd

# Build the binary locally
build:
	@echo "Building..."
	@go build -o main ./cmd

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f main

.PHONY: run build clean
