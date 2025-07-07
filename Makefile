# =============================================================================
# Ops Monorepo Makefile
# =============================================================================

.PHONY: help setup setup-force clean start stop restart logs build test test-user test-notification test-inventory test-order install buf protogen quick-start dev test-docker test-docker-clean

# Default target
help:
	@echo "==================================================================="
	@echo "Ops Monorepo - Available Commands"
	@echo "==================================================================="
	@echo "Setup:"
	@echo "  setup        Copy .env.example files to .env (safe - won't overwrite)"
	@echo "  setup-force  Force copy .env.example files (overwrites existing .env)"
	@echo ""
	@echo "Docker Operations:"
	@echo "  start        Start all services with docker-compose up"
	@echo "  stop         Stop all services"
	@echo "  restart      Restart all services"
	@echo "  logs         Show logs for all services"
	@echo "  build        Build all service images"
	@echo ""
	@echo "Development:"
	@echo "  test         Run tests for all services"
	@echo "  test-<svc>   Run tests for specific service (user, notification, inventory, order)"
	@echo "  test-docker  Test Docker builds for all services"
	@echo "  clean        Clean up docker containers and volumes"
	@echo "  dev          Start only databases for local development"
	@echo ""
	@echo "Tools:"
	@echo "  install      Install required tools (buf, grpcurl)"
	@echo "  protogen     Generate protobuf files"
	@echo ""
	@echo "Quick Start:"
	@echo "  quick-start  Setup environment and start all services"
	@echo ""
	@echo "Examples:"
	@echo "  make setup && make start    # Setup environment and start services"
	@echo "  make logs                   # View logs"
	@echo "  make stop                   # Stop all services"

# Setup environment files
setup:
	@echo "Setting up environment files..."
	@chmod +x scripts/setup-env.sh
	@scripts/setup-env.sh

setup-force:
	@echo "Force setting up environment files..."
	@chmod +x scripts/setup-env.sh
	@scripts/setup-env.sh --force

# Docker operations
start:
	@echo "Starting all services..."
	@docker-compose up -d
	@echo "Services started! Check status with: make logs"

stop:
	@echo "Stopping all services..."
	@docker-compose down

restart:
	@echo "Restarting all services..."
	@docker-compose down
	@docker-compose up -d

logs:
	@echo "Showing logs for all services..."
	@docker-compose logs -f

build:
	@echo "Building all service images..."
	@docker-compose build

# Development
test:
	@echo "Running tests for all services..."
	@failed_services=""; \
	for service in svc-user svc-notification svc-inventory svc-order; do \
		echo "=== Testing $$service ==="; \
		if [ -d "services/$$service" ]; then \
			if cd services/$$service && make test && cd ../..; then \
				echo "✓ $$service tests passed"; \
			else \
				echo "✗ $$service tests failed"; \
				failed_services="$$failed_services $$service"; \
			fi; \
		else \
			echo "Warning: Service $$service not found"; \
		fi; \
	done; \
	if [ -n "$$failed_services" ]; then \
		echo "❌ Tests failed for:$$failed_services"; \
		exit 1; \
	else \
		echo "✅ All service tests completed successfully!"; \
	fi

# Test individual services
test-user:
	@echo "Testing svc-user..."
	@cd services/svc-user && make test

test-notification:
	@echo "Testing svc-notification..."
	@cd services/svc-notification && make test

test-inventory:
	@echo "Testing svc-inventory..."
	@cd services/svc-inventory && make test

test-order:
	@echo "Testing svc-order..."
	@cd services/svc-order && make test

# Cleanup
clean:
	@echo "Cleaning up docker containers and volumes..."
	@docker-compose down -v
	@docker system prune -f

# Tools installation
install:
	@echo "Installing required tools..."
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Protobuf operations (alias for existing buf target)
buf:
	buf lint
	buf dep update
	buf generate

protogen: buf

# Quick start (setup + start)
quick-start: setup start
	@echo "Quick start completed!"
	@echo "Services are running. Use 'make logs' to view logs."

# Development mode (start only databases)
dev:
	@echo "Starting development mode (databases only)..."
	@docker-compose up -d user-db notification-db inventory-db order-db
	@echo "Databases started. Run services locally with 'make run' in each service directory."

# Test Docker builds for all services
test-docker:
	@echo "Testing Docker builds for all services..."
	@chmod +x scripts/test-docker-builds.sh
	@scripts/test-docker-builds.sh

# Test Docker builds with cleanup
test-docker-clean:
	@echo "Testing Docker builds with cleanup..."
	@chmod +x scripts/test-docker-builds.sh
	@scripts/test-docker-builds.sh --clean


