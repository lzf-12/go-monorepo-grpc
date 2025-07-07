# Notification Service

A gRPC service for sending emails via SMTP with comprehensive logging capabilities.

## Features

- Send emails via SMTP with configurable settings
- HTML and plain text email support
- Comprehensive email logging to PostgreSQL database
- Built with proven Go libraries (gomail.v2)
- gRPC API for service-to-service communication

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL database
- SMTP server credentials

## Configuration

Set the following environment variables or create a `.env` file:

```bash
# Application Configuration
ENV=development
APPNAME=svc-notification
DEBUG_MODE=true
PORT=50052

# Database Configuration
DB_URI=postgres://user:password@localhost:5432/notification_db
INIT_SEEDS=true

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourcompany.com
```

## Installation

1. Clone the repository and navigate to the service directory:
```bash
cd services/svc-notification
```

2. Install dependencies:
```bash
make deps
```

3. Set up your database and configure environment variables

4. Run the service:
```bash
make run
```

## Usage

### gRPC API

The service exposes one gRPC method:

#### SendEmail

Sends an email and logs the activity to the database.

**Request:**
```protobuf
message SendEmailRequest {
  string to = 1;        // Recipient email address
  string subject = 2;   // Email subject
  string body = 3;      // Email body content
  string from = 4;      // Sender email (optional, uses config default if empty)
  bool is_html = 5;     // Whether the body is HTML content
}
```

**Response:**
```protobuf
message SendEmailResponse {
  bool success = 1;           // Whether the email was sent successfully
  string message = 2;         // Status message
  string email_id = 3;        // Unique identifier for the email log entry
  google.protobuf.Timestamp timestamp = 4;  // Timestamp of the operation
}
```

### Example Usage

#### Go gRPC Client

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    notificationv1 "pb_schemas/notification/v1"
)

func main() {
    // Connect to the notification service
    conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := notificationv1.NewNotificationServiceClient(conn)

    // Send an email
    req := &notificationv1.SendEmailRequest{
        To:      "user@example.com",
        Subject: "Welcome to Our Service",
        Body:    "<h1>Welcome!</h1><p>Thank you for joining our service.</p>",
        From:    "noreply@yourcompany.com",
        IsHtml:  true,
    }

    resp, err := client.SendEmail(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }

    log.Printf("Email sent successfully: %s (ID: %s)", resp.Message, resp.EmailId)
}
```

## Database Schema

The service uses PostgreSQL with the following tables:

### Entity Relationship Diagram

```
┌─────────────────────────────────┐
│           email_logs            │
├─────────────────────────────────┤
│ id (PK)                         │
│ to_email                        │
│ from_email                      │
│ subject                         │
│ body                            │
│ is_html                         │
│ status                          │
│ error                           │
│ created_at                      │
│ updated_at                      │
└─────────────────────────────────┘
```

### Table Details

#### email_logs
- `id`: Unique identifier for each email log entry
- `to_email`: Recipient email address
- `from_email`: Sender email address
- `subject`: Email subject
- `body`: Email body content
- `is_html`: Whether the body is HTML content
- `status`: Email status (pending, sent, failed)
- `error`: Error message if sending failed
- `created_at`: When the email log was created
- `updated_at`: When the email log was last updated

## SMTP Configuration

The service uses the gomail.v2 library for sending emails. Make sure to configure your SMTP settings correctly:

### For Gmail:
- Use App Passwords instead of your regular password
- Enable 2-factor authentication
- Generate an App Password specifically for this application

### For Other Providers:
- Check your email provider's SMTP settings
- Ensure you have the correct host, port, and authentication credentials

## Development

### Running Tests
```bash
make test
```

### Code Formatting
```bash
make fmt
```

### Building
```bash
make build
```

### Docker Support
```bash
make docker-build
make docker-run
```

## Production Deployment

1. Set `ENV=production` in your environment
2. Use a production-grade PostgreSQL database
3. Configure proper SMTP credentials
4. Set `INIT_SEEDS=false` to avoid running seed data in production
5. Use `make build-prod` for production builds

## Security Considerations

- Store SMTP credentials securely (use environment variables or secrets management)
- Use TLS/SSL for SMTP connections
- Implement rate limiting to prevent email spam
- Validate email addresses before sending
- Consider implementing email templates for security

## Troubleshooting

### Common Issues:

1. **SMTP Authentication Failed**
   - Check your SMTP credentials
   - Ensure 2FA is enabled and use App Password for Gmail
   - Verify SMTP host and port settings

2. **Database Connection Issues**
   - Verify database URI format
   - Check database server is running
   - Ensure database exists and user has proper permissions

3. **Email Not Sending**
   - Check the `email_logs` table for error messages
   - Verify SMTP server connectivity
   - Check firewall settings for outbound SMTP ports

## Dependencies

### Service Dependencies
- **Database**: PostgreSQL (port 5432)
- **svc-user**: Required for service startup sequence (no direct communication)

### Docker Dependencies
- **notification-db**: PostgreSQL container for email logs storage
- **svc-user**: Must be started before this service in Docker Compose

### Service Startup Order
1. **user-db** (PostgreSQL database)
2. **svc-user** (User service)
3. **notification-db** (PostgreSQL database)
4. **svc-notification** (This service)

## Docker Deployment

### Building and Running with Docker (from service directory)

```bash
# Navigate to service directory
cd services/svc-notification

# Build Docker image (builds with monorepo context)
make docker-build

# Run with Docker (standalone)
make docker-run

# Or run with environment file
make docker-run-env

# Stop the container
make docker-stop
```

### Manual Docker Commands

```bash
# Build from service directory
docker build -t svc-notification -f Dockerfile ../../

# Run with port mapping
docker run --rm -p 50052:50052 svc-notification

# Run with environment file
docker run --rm --env-file .env -p 50052:50052 svc-notification
```

### Docker Compose Deployment

```bash
# From root directory
docker-compose up svc-notification
```

### Environment Variables for Docker
```bash
ENV=production
APPNAME=svc-notification
DEBUG_MODE=false
PORT=50052
DB_URI=postgres://notification_service:notification_password123@notification-db:5432/notification_db?sslmode=disable
INIT_SEEDS=true
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourcompany.com
```

## Support

For issues and questions, please check the logs in the database and application output for detailed error messages.