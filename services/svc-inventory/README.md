# Inventory Service

A gRPC service for managing inventory stock operations including checking availability, reserving stock, and releasing reservations.

## Features

- Check stock availability for multiple SKUs
- Reserve stock for orders with transaction safety
- Release stock reservations
- Historical tracking of reservations
- PostgreSQL database with ACID compliance
- gRPC API for service-to-service communication

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL database

## Configuration

Set the following environment variables or create a `.env` file:

```bash
# Application Configuration
ENV=development
APPNAME=svc-inventory
DEBUG_MODE=true
PORT=50051

# Database Configuration
DB_URI=postgres://user:password@localhost:5432/inventory_db
INIT_SEEDS=true
```

## Installation

1. Clone the repository and navigate to the service directory:
```bash
cd services/svc-inventory
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

## gRPC API

The service exposes three gRPC methods:

### CheckStock

Check stock availability for items without modifying inventory.

**Request:**
```protobuf
message StandardInventoryRequest {
  string order_id = 1;
  repeated InventoryItem items = 2;
}

message InventoryItem {
  string sku = 1;
  double req_qty_per_uom = 2;
  string uom = 3;
}
```

**Response:**
```protobuf
message InventoryStatusResponse {
  repeated InventoryStatus items = 1;
  google.protobuf.Timestamp timestamp = 2;
}
```

### ReserveStock

Reserve inventory items for an order.

**Request:** Same as CheckStock

**Response:**
```protobuf
message InventoryReservationResponse {
  string order_id = 1;
  SuccessProcessedItems success_processed_items = 2;
  FailedProcessedItems failed_processed_items = 3;
  google.protobuf.Timestamp timestamp = 4;
}
```

### ReleaseStock

Release previously reserved inventory items.

**Request:** Same as CheckStock

**Response:** Same as ReserveStock

## Usage Examples

### Go gRPC Client

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    inventoryv1 "pb_schemas/inventory/v1"
)

func main() {
    // Connect to the inventory service
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := inventoryv1.NewInventoryServiceClient(conn)

    // Check stock
    req := &inventoryv1.StandardInventoryRequest{
        OrderId: "order-123",
        Items: []*inventoryv1.InventoryItem{
            {
                Sku: "WIDGET-001",
                ReqQtyPerUom: 5.0,
                Uom: "pieces",
            },
        },
    }

    resp, err := client.CheckStock(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to check stock: %v", err)
    }

    log.Printf("Stock check result: %+v", resp)
}
```

## Database Schema

The service uses PostgreSQL with inventory-related tables for tracking stock levels, reservations, and historical data.

## Dependencies

### Service Dependencies
- **Database**: PostgreSQL (port 5432)
- **svc-user**: Required for service startup sequence (no direct communication)
- **svc-notification**: Required for service startup sequence (no direct communication)

### Docker Dependencies
- **inventory-db**: PostgreSQL container for inventory data storage
- **svc-notification**: Must be started before this service in Docker Compose

### Service Startup Order
1. **user-db** (PostgreSQL database)
2. **svc-user** (User service)
3. **notification-db** (PostgreSQL database)
4. **svc-notification** (Notification service)
5. **inventory-db** (PostgreSQL database)
6. **svc-inventory** (This service)

## Docker Deployment

### Docker Build
```bash
# Build the service
docker build -f services/svc-inventory/Dockerfile -t svc-inventory .

# Run with Docker Compose
docker-compose up svc-inventory
```

### Environment Variables for Docker
```bash
ENV=production
APPNAME=svc-inventory
DEBUG_MODE=false
PORT=50051
DB_URI=postgres://inventory_service:inventory_password123@inventory-db:5432/inventory_db?sslmode=disable
INIT_SEEDS=true
```

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
3. Set `INIT_SEEDS=false` to avoid running seed data in production
4. Use `make build-prod` for production builds

## Error Handling

The service includes comprehensive error handling for:
- SKU not found
- Insufficient quantity
- Database transaction errors
- Invalid UOM pairs

## Troubleshooting

### Common Issues:

1. **Database Connection Issues**
   - Verify database URI format
   - Check database server is running
   - Ensure database exists and user has proper permissions

2. **Stock Reservation Failures**
   - Check for sufficient available quantity
   - Verify SKU exists in inventory
   - Check for concurrent reservation conflicts

3. **gRPC Connection Issues**
   - Verify service is running on correct port
   - Check network connectivity
   - Ensure proper gRPC client configuration

## Support

For issues and questions, check the application logs and database for detailed error information.