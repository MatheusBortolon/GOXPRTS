# Go Orders Service

This project implements a service for managing orders using three different interfaces: REST, gRPC, and GraphQL. It includes a PostgreSQL database for data persistence and utilizes Docker for containerization.

## Project Structure

```
go-orders-service
├── cmd
│   ├── rest          # REST API entry point
│   ├── grpc          # gRPC service entry point
│   └── graphql       # GraphQL service entry point
├── api
│   ├── proto         # Protocol Buffers definitions for gRPC
│   └── graphql       # GraphQL schema definitions
├── internal
│   ├── orders        # Business logic and data access for orders
│   ├── transport     # HTTP, gRPC, and GraphQL transport layers
│   └── db           # Database connection logic
├── migrations        # Database migration files
├── scripts           # Scripts for running migrations
├── docker            # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── Makefile          # Build and run commands
├── .env.example      # Example environment variables
├── go.mod            # Go module definition
└── README.md         # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-orders-service
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Set up the environment variables:**
   Copy `.env.example` to `.env` and fill in the required values.

4. **Start all services (database + applications):**
   ```
   docker-compose up --build
   ```

5. **Run database migrations (in a separate terminal):**
   ```powershell
   # PowerShell (Windows)
   .\scripts\migrate.ps1
   
   # Bash (Linux/Mac/Git Bash)
   ./scripts/migrate.sh
   ```

6. **Run tests:**
   ```bash
   # Run all tests
   go test ./...
   
   # Run tests with coverage
   go test -cover ./...
   
   # Run tests with detailed coverage report
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   
   # Run specific test package
   go test ./internal/orders
   go test ./internal/transport/rest
   go test ./internal/transport/grpc
   go test ./internal/transport/graphql
   ```

## Usage

- **REST API:** Access the REST API at `http://localhost:8080/orders`.
- **gRPC Service:** Connect to the gRPC service at `localhost:50051`.
- **GraphQL API:** Access the GraphQL API at `http://localhost:8081/graphql`.

## Testing

Unit tests provide comprehensive coverage across all components:

### Coverage Summary

- **Orders Service:** 84.2% coverage - Tests for business logic, order creation, retrieval, and listing
- **REST Transport:** 100% coverage - HTTP handler tests including error scenarios and edge cases
- **gRPC Transport:** 75% coverage - gRPC server tests with context handling and error propagation
- **GraphQL Transport:** 100% coverage - GraphQL resolver tests for queries and mutations

### Test Files Structure

- **Service Layer:**
  - `internal/orders/service_test.go` - Core business logic tests
  - `internal/orders/service_extended_test.go` - Extended tests for edge cases
  - `internal/orders/service_edge_cases_test.go` - Additional edge case tests
  - `internal/orders/repository_test.go` - Repository implementation tests

- **REST Transport:**
  - `internal/transport/rest/handlers_test.go` - Basic HTTP handler tests
  - `internal/transport/rest/handlers_extended_test.go` - Extended handler tests
  - `internal/transport/rest/handlers_edge_cases_test.go` - Edge case and error handling tests

- **gRPC Transport:**
  - `internal/transport/grpc/server_test.go` - Basic gRPC server tests
  - `internal/transport/grpc/server_extended_test.go` - Extended server tests including context and error handling

- **GraphQL Transport:**
  - `internal/transport/graphql/resolver_test.go` - Complete resolver tests

- **Mock Repository:** `internal/orders/mock.go` - Shared mock for all test packages

### Run Tests

First, navigate to the project directory:

```bash
cd go-orders-service
```

Then run tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report for all packages
go test ./internal/orders ./internal/transport/rest ./internal/transport/grpc ./internal/transport/graphql -cover

# Generate HTML coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run tests for specific package
go test ./internal/orders -v -cover
go test ./internal/transport/rest -v -cover
go test ./internal/transport/grpc -v -cover
go test ./internal/transport/graphql -v -cover

# Clean test cache and re-run
go clean -testcache
go test ./internal/... -cover
```

### Run Tests (PowerShell)

First, navigate to the project directory:

```powershell
cd .\go-orders-service
```

Then run tests:

```powershell
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report for all packages
go test ./internal/orders ./internal/transport/rest ./internal/transport/grpc ./internal/transport/graphql -cover

# Generate coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run tests for specific package
go test ./internal/orders -v -cover
go test ./internal/transport/rest -v -cover
go test ./internal/transport/grpc -v -cover
go test ./internal/transport/graphql -v -cover

# Clean test cache and re-run
go clean -testcache
go test ./internal/... -cover
```

### Coverage Analysis

To view detailed coverage by function:

```bash
# Generate coverage profile
go test ./internal/... -coverprofile=coverage.out

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.