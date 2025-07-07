# Build stage
FROM golang:1.24.4-alpine AS builder

# Install git for dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy the entire monorepo (when building from root: docker build -f services/svc-notification/Dockerfile .)
COPY . ./

# Set working directory to this service
WORKDIR /app/services/svc-notification

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/svc-notification/main .

# Copy the service's .env file as default
COPY --from=builder /app/services/svc-notification/.env .env

# Expose port
EXPOSE 50052

# Run the binary
CMD ["./main"]