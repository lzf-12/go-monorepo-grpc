# Ops Monorepo

A Go microservices monorepo with user authentication, notifications, inventory management, and order processing.

## Quick Start

```bash
make setup && docker-compose up
```

## Prerequisites

### System Requirements
- **Go**: 1.24.2+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Make**: GNU Make 4.0+

### Tools (for development)
- **Buf**: v1.28+ (protocol buffer management)
- **Mockery**: v2.40+ (mock generation)
- **PostgreSQL Client**: 15+ (for database debugging)

## First Principles & Assumptions

1. **Docker-first development**: All services run in containers
2. **Database isolation**: Each service has its own PostgreSQL database
3. **gRPC for inter-service communication**: Services communicate via gRPC
4. **JWT authentication**: Token-based auth with User service as authority
5. **Environment file hierarchy**: Root .env for shared config, service .env for specifics

## Getting Started

### 1. Clone and Setup
```bash
git clone <repository-url>
cd ops-monorepo
make setup && docker-compose up
```

### 2. Verify Services
```bash
# Check all containers are running
docker-compose ps

# Test user registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "first_name": "Test", "last_name": "User"}'
```

## Services

| Service | gRPC | HTTP | Database | Documentation |
|---------|------|------|----------|---------------|
| svc-user | 50053 | 8080 | 5436 | [README](services/svc-user/README.md) |
| svc-notification | 50052 | - | 5433 | [README](services/svc-notification/README.md) |
| svc-inventory | 50051 | - | 5434 | [README](services/svc-inventory/README.md) |
| svc-order | - | 8081 | 5435 | [README](services/svc-order/readme.md) |

## Development Commands

### Environment Management
```bash
make setup        # Copy .env.example to .env files
make setup-force  # Force overwrite existing .env files
```

### Docker Operations
```bash
make start        # Start all services
make stop         # Stop all services  
make logs         # View all service logs
make clean        # Remove containers and volumes
```

### Development Mode
```bash
make dev          # Start only databases for local development
```

### Code Generation
```bash
# Generate protocol buffers
make protogen
# OR
buf generate

# Generate mocks (per service)
cd services/svc-user && mockery
cd services/svc-notification && mockery  
cd services/svc-inventory && mockery
cd services/svc-order && mockery
```

### Testing
```bash
# All services
make test

# Individual service
cd services/svc-order && make test
```

## Tool Installation

### Install Buf (Protocol Buffers)
```bash
# macOS
brew install bufbuild/buf/buf

# Linux
curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.28.1/buf-$(uname -s)-$(uname -m)" \
  -o "/usr/local/bin/buf" && chmod +x "/usr/local/bin/buf"
```

### Install Mockery (Mock Generation)
```bash
# Go install
go install github.com/vektra/mockery/v2@v2.40.1

# Or download binary from releases
```

## Configuration

### Root .env (shared database config)
```bash
# Database credentials
POSTGRES_USER_DB=user_db
POSTGRES_USER_SERVICE=user_service
POSTGRES_USER_PASSWORD=user_password123
# ... (similar for other services)
```

### Service .env (service-specific)
```bash
# services/svc-user/.env
ENV=local
PORT=50053
HTTP_PORT=8080
DB_URI=postgres://user_service:user_password123@user-db:5432/user_db?sslmode=disable
JWT_SECRET=super-secret-jwt-key-for-development
```

## API Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123", "first_name": "John", "last_name": "Doe"}'
```

### 2. Login and Get Token
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}' \
  | jq -r '.data.access_token')
```

### 3. Create Order (Protected)
```bash
curl -X POST http://localhost:8081/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_items": [{"sku": "OLIVE-OIL-1L", "quantity_per_uom": 2, "uom": "L"}]}'
```

## Troubleshooting

### Port Conflicts
```bash
# Check what's using ports
lsof -i :5432
lsof -i :8080

# Stop conflicting services
docker stop $(docker ps -q --filter ancestor=postgres)
```

### Database Issues
```bash
# Check database connectivity
docker-compose exec user-db pg_isready -U user_service -d user_db

# View service logs
docker-compose logs svc-user
```

### Authentication Issues
```bash
# Verify user service is running
curl http://localhost:8080/health

# Check JWT token validity
echo $TOKEN | cut -d. -f2 | base64 -d
```

## Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌───────────┐
│  svc-user   │    │svc-notification│   │svc-inventory│    │ svc-order │
│   :50053    │    │    :50052      │   │   :50051    │    │   :8081   │
│   :8080     │    │                │   │             │    │           │
└─────────────┘    └──────────────┘    └─────────────┘    └───────────┘
       │                   │                   │                 │
       │                   │                   │                 │
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌───────────┐
│   user-db   │    │notification-db│   │inventory-db │    │ order-db  │
│   :5436     │    │    :5433      │   │   :5434     │    │   :5435   │
└─────────────┘    └──────────────┘    └─────────────┘    └───────────┘
```

## Make Commands Reference

```bash
make help         # Show all commands
make setup        # Setup environment files
make start        # Start all services  
make stop         # Stop all services
make restart      # Restart all services
make logs         # Show logs
make test         # Run tests
make clean        # Clean containers/volumes
make dev          # Development mode (databases only)
make protogen     # Generate protobuf files
make quick-start  # Setup + start (recommended)
```