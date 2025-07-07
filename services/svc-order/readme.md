# Order Service

A REST API service for managing orders with JWT authentication, inventory integration, and comprehensive order processing.

## Features

- RESTful API for order management
- JWT authentication via middleware
- Integration with inventory service for stock validation
- PostgreSQL database for order persistence
- Gin framework for HTTP routing
- Docker containerization support

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL database
- Access to user service for authentication
- Access to inventory service for stock management

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Application Configuration
ENV=development
APPNAME=svc-order
DEBUG_MODE=true
PORT=8080

# Database Configuration
DB_URI=postgres://user:password@localhost:5432/order_db
INIT_SEEDS=true

# Service Dependencies
USER_SERVICE_URL=localhost:50053
INVENTORY_SERVICE_URL=localhost:50051
NOTIFICATION_SERVICE_URL=localhost:50052
```

## Installation

1. Clone the repository and navigate to the service directory:
```bash
cd services/svc-order
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Adjust environment values as needed

4. Install dependencies:
```bash
make deps
```

5. Run the service:
```bash
make run
```

## API Endpoints

### Authentication Required

All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### POST /api/v1/orders

Create a new order.

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "customer_id": "customer-123",
  "items": [
    {
      "sku": "WIDGET-001",
      "quantity": 2,
      "unit_price": 29.99
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "order_id": "order-456",
    "status": "pending",
    "total_amount": 59.98,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

## Authentication

The service uses JWT authentication middleware that validates tokens with the user service.

### Getting a JWT Token

1. Register/Login via user service:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com", 
    "password": "password123"
  }'
```

2. Use the returned `access_token` in subsequent requests:
```bash
curl -X POST http://localhost:8081/api/v1/orders \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "123", "items": [...]}'
```

## Database Schema

The service uses PostgreSQL for order management and tracking.

### Entity Relationship Diagram

```
┌─────────────────────────────────┐
│             orders              │
├─────────────────────────────────┤
│ id (PK)                         │
│ user_id                         │
│ user_email                      │
│ status                          │
│ total_amount                    │
│ currency                        │
│ created_at                      │
│ updated_at                      │
└─────────────────────────────────┘
                │
                │ 1:N
                │
                ▼
┌─────────────────────────────────┐
│          order_items            │
├─────────────────────────────────┤
│ id (PK)                         │
│ order_id (FK)                   │
│ sku                             │
│ quantity_per_uom                │
│ price_per_uom                   │
│ uom_code                        │
└─────────────────────────────────┘
```

### Table Details

#### orders
- `id`: Unique identifier for each order (UUID)
- `user_id`: Reference to the user who placed the order
- `user_email`: Email address of the user
- `status`: Order status (PENDING, CONFIRMED, CANCELLED)
- `total_amount`: Total order amount
- `currency`: Currency code (default: USD)
- `created_at`: When the order was created
- `updated_at`: When the order was last updated

#### order_items
- `id`: Unique identifier for each order item (UUID)
- `order_id`: Reference to the parent order
- `sku`: Stock keeping unit identifier
- `quantity_per_uom`: Quantity per unit of measure
- `price_per_uom`: Price per unit of measure
- `uom_code`: Unit of measure code

### Key Relationships

- **orders** can have multiple **order_items** (one-to-many)
- **order_items** reference inventory SKUs but don't enforce foreign key constraints (loose coupling)
- Unique constraint on (order_id, sku) prevents duplicate items in the same order

## Dependencies

### Service Dependencies
- **svc-user**: JWT token validation (gRPC on port 50053)
- **svc-inventory**: Stock availability checking (gRPC on port 50051)
- **svc-notification**: Email notifications (gRPC on port 50052)
- **Database**: PostgreSQL (port 5432)

### Docker Dependencies
- **order-db**: PostgreSQL container for order data storage
- **svc-user**: Required for JWT authentication
- **svc-inventory**: Required for stock validation
- **svc-notification**: Required for order notifications

### Service Startup Order
1. **user-db** (PostgreSQL database)
2. **svc-user** (User service)
3. **notification-db** (PostgreSQL database)
4. **svc-notification** (Notification service)
5. **inventory-db** (PostgreSQL database)
6. **svc-inventory** (Inventory service)
7. **order-db** (PostgreSQL database)
8. **svc-order** (This service)

## Docker Deployment

### Docker Build
```bash
# Build the service
docker build -f services/svc-order/Dockerfile -t svc-order .

# Run with Docker Compose (recommended)
docker-compose up svc-order
```

### Environment Variables for Docker
```bash
ENV=production
APPNAME=svc-order
DEBUG_MODE=false
PORT=8080
DB_URI=postgres://order_service:order_password123@order-db:5432/order_db?sslmode=disable
INIT_SEEDS=true
USER_SERVICE_URL=svc-user:50053
INVENTORY_SERVICE_URL=svc-inventory:50051
NOTIFICATION_SERVICE_URL=svc-notification:50052
```

## Development

### Running Tests
```bash
make test
```

### Generate Mocks
```bash
mockery --all --dir "./internal" --output "./mocks/internal" --keeptree && mockery --all --dir "./validator" --output "./mocks/validator" --keeptree
```

### Code Formatting
```bash
make fmt
```

### Building
```bash
make build
```

## Integration Testing

### End-to-End Flow

1. **Start all services:**
```bash
docker-compose up
```

2. **Get JWT token:**
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}' \
  | jq -r '.data.access_token')
```

3. **Create order:**
```bash
curl -X POST http://localhost:8081/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-123",
    "items": [
      {
        "sku": "WIDGET-001",
        "quantity": 2,
        "unit_price": 29.99
      }
    ]
  }'
```

## Error Handling

### Authentication Errors
- **401 Unauthorized**: Missing or invalid JWT token
- **403 Forbidden**: Insufficient permissions

### Business Logic Errors
- **400 Bad Request**: Invalid order data
- **409 Conflict**: Insufficient inventory
- **500 Internal Server Error**: Service communication failures

## Troubleshooting

### Common Issues:

1. **Authentication Failures**
   - Verify user service is running and accessible
   - Check JWT token format and expiration
   - Ensure proper Authorization header format

2. **Service Communication Errors**
   - Verify all dependent services are running
   - Check service URLs in environment variables
   - Test gRPC connectivity manually

3. **Database Issues**
   - Verify database connection string
   - Check database permissions
   - Ensure migrations/seeds have run

## Support

For issues and questions, check the application logs and verify all dependent services are running and accessible.