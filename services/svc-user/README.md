# User Service

simple user authentication service providing both gRPC and HTTP APIs for user management, JWT authentication, and token validation.

## Features

- User registration and authentication
- JWT-based access and refresh tokens
- Password hashing with bcrypt
- Role-based access control
- Token validation for other services
- Both gRPC and HTTP APIs
- PostgreSQL database with comprehensive user management

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL database

## Configuration

Set the following environment variables or create a `.env` file:

```bash
# Example Application Configuration
ENV=development
APPNAME=svc-user
DEBUG_MODE=true
PORT=50053           # gRPC port
HTTP_PORT=8080       # HTTP port

# Database Configuration
DB_URI=postgres://user:password@localhost:5432/user_db
INIT_SEEDS=true

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_DURATION=15m
JWT_REFRESH_DURATION=24h
```

## Installation

1. Clone the repository and navigate to the service directory:
```bash
cd services/svc-user
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

The service will start both gRPC server (port 50053) and HTTP server (port 8080).

## APIs

### gRPC API

#### ValidateToken

Validates JWT tokens for other services.

**Request:**
```protobuf
message ValidateTokenRequest {
  string token = 1;  // JWT token to validate
}
```

**Response:**
```protobuf
message ValidateTokenResponse {
  bool valid = 1;              // Whether the token is valid
  string user_email = 2;       // Email of the authenticated user
  repeated string roles = 3;   // User roles
}
```

### HTTP API

#### POST /api/v1/auth/register

Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "success": true,
  "message": "user registered successfully",
  "data": {
    "id": "user-uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "roles": ["user"],
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### POST /api/v1/auth/login

Authenticate user and get tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "uuid-refresh-token",
    "user": {
      "id": "user-uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "roles": ["user"],
      "is_active": true
    }
  }
}
```

#### POST /api/v1/auth/refresh

Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "uuid-refresh-token"
}
```

**Response:**
```json
{
  "success": true,
  "message": "token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "new-uuid-refresh-token"
  }
}
```

## Usage Examples

### HTTP API Usage

#### Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

#### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

#### Refresh token:
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token"
  }'
```

### gRPC API Usage

#### Go gRPC Client for Token Validation:

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    userv1 "pb_schemas/user/v1"
)

func main() {
    // Connect to the user service
    conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := userv1.NewUserServiceClient(conn)

    // Validate a token
    req := &userv1.ValidateTokenRequest{
        Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    }

    resp, err := client.ValidateToken(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to validate token: %v", err)
    }

    if resp.Valid {
        log.Printf("Token is valid for user: %s with roles: %v", 
                   resp.UserEmail, resp.Roles)
    } else {
        log.Printf("Token is invalid")
    }
}
```

## Database Schema

The service uses PostgreSQL with the following tables:

### Entity Relationship Diagram

```
┌─────────────────────────────────┐
│              users              │
├─────────────────────────────────┤
│ id (PK)                         │
│ email                           │
│ password                        │
│ first_name                      │
│ last_name                       │
│ roles                           │
│ is_active                       │
│ created_at                      │
│ updated_at                      │
└─────────────────────────────────┘
                │
                │ 1:N
                │
                ▼
┌─────────────────────────────────┐
│         refresh_tokens          │
├─────────────────────────────────┤
│ id (PK)                         │
│ user_id (FK)                    │
│ token                           │
│ expires_at                      │
│ created_at                      │
│ is_revoked                      │
└─────────────────────────────────┘
```

### Table Details

#### users
- `id`: Unique identifier for each user
- `email`: User's email address (unique)
- `password`: Bcrypt hashed password
- `first_name`: User's first name
- `last_name`: User's last name
- `roles`: Array of user roles (e.g., ["user", "admin"])
- `is_active`: Whether the user account is active
- `created_at`: When the user was created
- `updated_at`: When the user was last updated

#### refresh_tokens
- `id`: Unique identifier for each refresh token
- `user_id`: Reference to the user
- `token`: The refresh token value
- `expires_at`: When the token expires
- `created_at`: When the token was created
- `is_revoked`: Whether the token has been revoked

## Default Users

When `INIT_SEEDS=true`, the following test users are created:

1. **Admin User**
   - Email: `admin@example.com`
   - Password: `password123`
   - Roles: `["admin", "user"]`

2. **Regular User**
   - Email: `user@example.com`
   - Password: `password123`
   - Roles: `["user"]`

3. **Test User**
   - Email: `test@example.com`
   - Password: `password123`
   - Roles: `["user"]`

## Security Features

1. **Password Hashing**: Uses bcrypt with default cost for secure password storage
2. **JWT Tokens**: Stateless authentication with configurable expiration
3. **Refresh Tokens**: Secure token refresh mechanism with database storage
4. **Role-Based Access**: Support for multiple user roles
5. **Token Validation**: Comprehensive token validation for other services

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
2. Use a strong, unique JWT secret key
3. Configure appropriate token expiration times
4. Set `INIT_SEEDS=false` to avoid creating test users in production
5. Use a production-grade PostgreSQL database
6. Consider implementing rate limiting for authentication endpoints

## Integration with Other Services

Other services can validate user tokens by calling the gRPC `ValidateToken` method:

```go
// In your middleware or authentication handler
func validateUserToken(tokenString string) (*User, error) {
    conn, err := grpc.Dial("user-service:50053", grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := userv1.NewUserServiceClient(conn)
    
    resp, err := client.ValidateToken(context.Background(), &userv1.ValidateTokenRequest{
        Token: tokenString,
    })
    if err != nil {
        return nil, err
    }

    if !resp.Valid {
        return nil, errors.New("invalid token")
    }

    return &User{
        Email: resp.UserEmail,
        Roles: resp.Roles,
    }, nil
}
```

## Troubleshooting

### Common Issues:

1. **Database Connection Issues**
   - Verify database URI format
   - Check database server is running
   - Ensure database exists and user has proper permissions

2. **JWT Token Issues**
   - Verify JWT secret is consistent across restarts
   - Check token expiration settings
   - Ensure tokens are being passed correctly in requests

3. **Authentication Failures**
   - Check password hashing/verification
   - Verify user exists and is active
   - Check for case sensitivity in email addresses

## Dependencies

### Service Dependencies
- **Database**: PostgreSQL (port 5432)
- **No service dependencies** - This is a foundational service

### Docker Dependencies
- **user-db**: PostgreSQL container for user data storage

### Service Startup Order
1. **user-db** (PostgreSQL database)
2. **svc-user** (This service)

## Docker Deployment

### Building and Running with Docker (from service directory)

```bash
# Navigate to service directory
cd services/svc-user

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
docker build -t svc-user -f Dockerfile ../../

# Run with port mapping
docker run --rm -p 50053:50053 -p 8080:8080 svc-user

# Run with environment file
docker run --rm --env-file .env -p 50053:50053 -p 8080:8080 svc-user
```

### Docker Compose Deployment

```bash
# From root directory
docker-compose up svc-user
```

### Environment Variables for Docker
```bash
ENV=production
APPNAME=svc-user
DEBUG_MODE=false
PORT=50053
HTTP_PORT=8080
DB_URI=postgres://user_service:user_password123@user-db:5432/user_db?sslmode=disable
INIT_SEEDS=true
JWT_SECRET=super-secret-jwt-key-for-production
JWT_ACCESS_DURATION=15m
JWT_REFRESH_DURATION=24h
```

## Support

For issues and questions, check the application logs and database for detailed error information.