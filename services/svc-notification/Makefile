BINARY_NAME=svc-notification
MAIN_PATH=./main.go

.PHONY: build clean run test deps fmt lint build-prod docker-build docker-run docker-run-env docker-stop

# Build the application
build:
	go build -o ${BINARY_NAME} ${MAIN_PATH}

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}

# Run the application
run:
	go run ${MAIN_PATH}

# Run tests
test:
	go test -v ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_NAME} ${MAIN_PATH}

# Docker build (from service directory with monorepo context)
docker-build:
	cd ../../ && docker build -t ${BINARY_NAME} -f services/svc-notification/Dockerfile .

# Run with docker (standalone)
docker-run:
	docker run --rm -p 50052:50052 ${BINARY_NAME}

# Run with docker (with environment file)
docker-run-env:
	docker run --rm --env-file .env -p 50052:50052 ${BINARY_NAME}

# Stop docker containers
docker-stop:
	docker stop $$(docker ps -q --filter ancestor=${BINARY_NAME})