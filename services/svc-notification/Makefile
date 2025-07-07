BINARY_NAME=svc-notification
MAIN_PATH=./main.go

.PHONY: build clean run test deps

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

# Docker build
docker-build:
	docker build -t ${BINARY_NAME} .

# Run with docker-compose
docker-run:
	docker-compose up --build

# Stop docker-compose
docker-stop:
	docker-compose down