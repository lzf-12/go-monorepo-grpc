# Docker Setup Guide

This document explains how to build and run each service individually using Docker from their respective service directories.

## Overview

Each service can be built and run independently from its own directory using Docker. The Dockerfiles are designed to copy the entire monorepo context, ensuring access to shared libraries and dependencies.

## Quick Start

### Building All Services

From the root directory:
```bash
make test-docker
```

This will test build all services and verify they can start properly.

### Building Individual Services

From any service directory:
```bash
cd services/svc-user
make docker-build
make docker-run
```

## Service-Specific Instructions

### User Service (svc-user)

```bash
cd services/svc-user

# Build
make docker-build
# OR manually: docker build -t svc-user -f Dockerfile ../../

# Run (requires database)
make docker-run-env
# OR manually: docker run --rm --env-file .env -p 50053:50053 -p 8080:8080 svc-user

# Ports: 50053 (gRPC), 8080 (HTTP)
```

### Notification Service (svc-notification)

```bash
cd services/svc-notification

# Build
make docker-build
# OR manually: docker build -t svc-notification -f Dockerfile ../../

# Run (requires database and SMTP config)
make docker-run-env
# OR manually: docker run --rm --env-file .env -p 50052:50052 svc-notification

# Ports: 50052 (gRPC)
```

### Inventory Service (svc-inventory)

```bash
cd services/svc-inventory

# Build
make docker-build
# OR manually: docker build -t svc-inventory -f Dockerfile ../../

# Run (requires database)
make docker-run-env
# OR manually: docker run --rm --env-file .env -p 50051:50051 svc-inventory

# Ports: 50051 (gRPC)
```

### Order Service (svc-order)

```bash
cd services/svc-order

# Build
make docker-build
# OR manually: docker build -t svc-order -f Dockerfile ../../

# Run (requires database and other services)
make docker-run-env
# OR manually: docker run --rm --env-file .env -p 8080:8080 svc-order

# Ports: 8080 (HTTP)
```

## How It Works

### Dockerfile Structure

Each Dockerfile follows this pattern:

1. **Build Stage**: Uses `golang:1.24.4-alpine`
   - Copies entire monorepo context (`COPY ../../ ./`)
   - Changes to service directory
   - Downloads dependencies and builds binary

2. **Runtime Stage**: Uses `alpine:latest`
   - Copies binary from build stage
   - Exposes appropriate ports
   - Runs the service

### Context Directory

When building from a service directory, the Docker context is set to `../../` (the monorepo root), ensuring access to:
- Shared libraries (`shared-libs/`)
- Protocol buffer definitions (`pb_schemas/`)
- Go workspace files (`go.work`, `go.work.sum`)

### Environment Variables

Each service requires specific environment variables. Use `.env` files in each service directory or pass them via Docker:

```bash
# With environment file
docker run --env-file .env -p <port>:<port> <service-name>

# With individual variables
docker run -e DB_URI="..." -e JWT_SECRET="..." -p <port>:<port> <service-name>
```

## Testing Docker Builds

### Automated Testing

```bash
# Test all services
make test-docker

# Test with cleanup
make test-docker-clean
```

### Manual Testing

```bash
# Build service
cd services/svc-user
make docker-build

# Verify image exists
docker images | grep svc-user

# Test basic functionality (will fail without database, which is expected)
docker run --rm -p 50053:50053 -p 8080:8080 svc-user
```

## Production Considerations

### Environment Variables

- Set proper database connection strings
- Use strong JWT secrets
- Configure SMTP credentials for notification service
- Set `ENV=production` and `DEBUG_MODE=false`

### Resource Limits

```bash
docker run --memory=512m --cpus=0.5 --env-file .env -p <port>:<port> <service>
```

### Health Checks

Services don't include built-in health endpoints yet, but you can test connectivity:

```bash
# For gRPC services
grpcurl -plaintext localhost:50053 list

# For HTTP services
curl http://localhost:8080/health
```

## Troubleshooting

### Common Issues

1. **Build Fails with Missing Dependencies**
   - Ensure you're building from the service directory
   - Verify the Docker context is set to `../../`
   - Check that shared-libs exists in the monorepo root

2. **Service Starts but Crashes**
   - Check environment variables are set correctly
   - Verify database is accessible (for services that need it)
   - Check Docker logs: `docker logs <container-name>`

3. **Port Conflicts**
   - Ensure ports aren't already in use
   - Use different host ports: `-p 8081:8080` instead of `-p 8080:8080`

### Debug Commands

```bash
# Check if image was built
docker images | grep <service-name>

# Check running containers
docker ps

# View logs
docker logs <container-name>

# Interactive shell in container
docker run -it --entrypoint /bin/sh <service-name>
```

## Integration with Docker Compose

Individual Docker builds are compatible with the root-level docker-compose.yml:

```bash
# Build individual service
cd services/svc-user && make docker-build

# Use in compose (will use local image if available)
cd ../../ && docker-compose up svc-user
```