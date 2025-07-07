BINARY_NAME=svc-inventory
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
	cd ../../ && docker build -t ${BINARY_NAME} -f services/svc-inventory/Dockerfile .

# Run with docker (standalone)
docker-run:
	docker run --rm -p 50051:50051 ${BINARY_NAME}

# Run with docker (with environment file)
docker-run-env:
	docker run --rm --env-file .env -p 50051:50051 ${BINARY_NAME}

# Stop docker containers
docker-stop:
	docker stop $$(docker ps -q --filter ancestor=${BINARY_NAME})

# Legacy targets (keep for compatibility)
install:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen

check-specs:
	@SPEC_FILES=$$(find ./specs -maxdepth 1 -type f \( -name '*.yaml' -o -name '*.yml' \)); \
	if [ -z "$$SPEC_FILES" ]; then \
		echo "Error: No OpenAPI spec files (.yaml or .yml) found in ./specs directory"; \
		exit 1; \
	fi; \
	echo "found specs .yaml or .yml files:"; \
	echo "$$SPEC_FILES" | sed 's/ /\n/g'

gen-oapi-config:
	# generate or refresh oapi config file
	mkdir -p internal/delivery/types
	rm -f internal/delivery/types/cfg.yaml;\
	[ -f internal/delivery/types/cfg.yaml ] || echo "package: types\ngenerate:\n models: true" > internal/delivery/types/cfg.yaml;\
	
gen-go-directive:	
	rm -f internal/delivery/types/gen.go;\
	echo "package types" > internal/delivery/types/gen.go
	echo "" >> internal/delivery/types/gen.go
	find ./specs -maxdepth 1 -type f \( -name '*.yaml' -o -name '*.yml' \) | \
	while read -r spec_file; do \
		filename=$$(basename "$$spec_file"); \
		basename="$${filename%.*}"; \
		echo "//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml --o=$$basename.gen.go ../../../specs/$$filename" >> internal/delivery/types/gen.go; \
	done
	echo "" >> internal/delivery/types/gen.go
	echo "Generated directives for:" $$(find ./specs -maxdepth 1 -type f \( -name '*.yaml' -o -name '*.yml' \) -printf "%f ")


gen-types: check-specs gen-oapi-config gen-go-directive
	go generate ./internal/delivery/types/...

# run order service
run:
	go run main.go

	