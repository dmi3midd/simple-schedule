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

# Docker: Build containers
docker-build:
	@docker compose build

# Docker: Build and start all containers in background
docker-run:
	@docker compose up --build -d

# Docker: Stop all containers
docker-down:
	@docker compose down

# Docker: Stop all containers and remove volumes (clean database reset)
docker-down-v:
	@docker compose down -v

# Docker: Follow application logs
docker-logs:
	@docker compose logs -f api

# Docker: Restart application container
docker-restart:
	@docker compose restart api

# Generate Swagger documentation
swagger:
	@echo "Generating swagger docs..."
	@swag init -g main.go -d cmd,internal/server/handlers,internal/domain,internal/shared/httputils/apierror --parseInternal


.PHONY: run build swagger clean docker-build docker-run docker-down docker-down-v docker-logs docker-restart swagger
