<!-- /docs/README.md -->
# gRPC Example Application

A comprehensive gRPC-based application built with Go, featuring user authentication, database integration, and microservices architecture.

## Features

- **gRPC API** with Protocol Buffers
- **User Authentication** with JWT tokens
- **PostgreSQL Database** integration
- **Middleware** for authentication, logging, and recovery
- **Docker** containerization
- **Kubernetes** deployment manifests
- **Database Migration** system
- **Structured Logging** with Logrus
- **Configuration Management** with Viper

## Architecture

```
├── api/grpc/           # gRPC server setup
├── cmd/server/         # Application entry point
├── configs/            # Configuration files
├── deployments/        # Docker and Kubernetes manifests
├── internal/           # Application internal packages
│   ├── config/         # Configuration handling
│   ├── handler/grpc/   # gRPC handlers
│   ├── middleware/     # Middleware components
│   ├── model/          # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic layer
├── pkg/                # Shared packages
│   ├── database/       # Database connection and migration
│   ├── logger/         # Logging configuration
│   └── utils/          # Utility functions
├── proto/              # Protocol Buffer definitions
└── scripts/            # Utility scripts
```

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Protocol Buffers compiler (protoc)
- Docker (optional)
- Kubernetes (optional)

## Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd grpc-exmpl
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Generate Protocol Buffer code:**
```bash
./scripts/gen-proto.sh
```

4. **Set up PostgreSQL database:**
```bash
createdb grpc_exmpl
```

5. **Configure the application:**
   - Copy `configs/app.yaml.example` to `configs/app.yaml`
   - Update database credentials and JWT secret

## Running the Application

### Local Development

1. **Start the server:**
```bash
go run cmd/server/main.go
```

2. **The server will start on port 8080 by default**

### Using Docker

1. **Build the image:**
```bash
docker build -f deployments/docker/Dockerfile -t grpc-exmpl .
```

2. **Run with docker-compose:**
```bash
docker-compose up -d
```

### Kubernetes Deployment

1. **Apply the manifests:**
```bash
kubectl apply -f deployments/k8s/
```

## API Usage

### User Registration

```bash
grpcurl -plaintext -d '{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}' localhost:8080 user.UserService/Register
```

### User Login

```bash
grpcurl -plaintext -d '{
  "email": "john@example.com",
  "password": "password123"
}' localhost:8080 user.UserService/Login
```

### Get User Profile

```bash
grpcurl -plaintext -H "authorization: Bearer <JWT_TOKEN>" -d '{
  "token": "<JWT_TOKEN>"
}' localhost:8080 user.UserService/GetProfile
```

## Configuration

The application uses YAML configuration files located in the `configs/` directory:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  shutdown_timeout: "5s"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  database: "grpc_exmpl"
  ssl_mode: "disable"

jwt:
  secret: "your-secret-key"
  expiration: "24h"

log:
  level: "info"
  format: "json"
```

Configuration can be overridden using environment variables with uppercase and underscore format (e.g., `DATABASE_HOST`).

## Database Schema

### Users Table
- `id` - Primary key
- `username` - Unique username
- `email` - Unique email address
- `password` - Hashed password
- `full_name` - User's full name
- `created_at` - Record creation timestamp
- `updated_at` - Record update timestamp

## Development

### Project Structure

The project follows Go's standard project layout and clean architecture principles:

- **Handler Layer**: Handles gRPC requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Data access and persistence
- **Model Layer**: Data structures and validation

### Adding New Services

1. Define the service in Protocol Buffers (`.proto` files)
2. Generate Go code using `protoc`
3. Implement the service handler
4. Create service logic and repository
5. Register the service in the gRPC server

### Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Input validation and sanitization
- SQL injection prevention with parameterized queries
- HTTPS/TLS support (configure in production)

## Monitoring and Logging

- Structured logging with Logrus
- Request/response logging middleware
- Panic recovery middleware
- Health check endpoints
- Metrics collection ready

## Production Deployment

### Environment Variables

Set the following environment variables in production:

```bash
DATABASE_HOST=your-db-host
DATABASE_PASSWORD=secure-password
JWT_SECRET=super-secure-jwt-secret
LOG_LEVEL=warn
```

### Security Considerations

1. Change default JWT secret
2. Use strong database passwords
3. Enable SSL/TLS for database connections
4. Configure proper firewall rules
5. Use non-root user in containers
6. Regularly update dependencies

## Troubleshooting

### Common Issues

1. **Database connection failed**: Check PostgreSQL is running and credentials are correct
2. **Proto compilation errors**: Ensure `protoc` is installed and PATH is set correctly
3. **Port already in use**: Change the server port in configuration
4. **JWT validation errors**: Verify JWT secret matches between services

### Logs

Check application logs for detailed error information:

```bash
# Local development
tail -f logs/app.log

# Docker
docker logs grpc-exmpl-container

# Kubernetes
kubectl logs -f deployment/grpc-exmpl-app
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions or issues, please:
1. Check the documentation
2. Search existing issues
3. Create a new issue with detailed information