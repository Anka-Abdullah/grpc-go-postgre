# Architecture Documentation

## Overview

This document describes the architecture of the gRPC Example Application, including its components, patterns, and design decisions.

## System Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │    │   gRPC Client   │    │   gRPC Client   │
│   (Mobile App)  │    │   (Web App)     │    │   (Service)     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │                         │
                    │    Load Balancer        │
                    │                         │
                    └────────────┬────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
┌─────────▼───────┐    ┌─────────▼───────┐    ┌─────────▼───────┐
│                 │    │                 │    │                 │
│  gRPC Server    │    │  gRPC Server    │    │  gRPC Server    │
│   Instance 1    │    │   Instance 2    │    │   Instance 3    │
│                 │    │                 │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │                         │
                    │    PostgreSQL DB        │
                    │                         │
                    └─────────────────────────┘
```

### Component Architecture

```
┌─────────────────────────────────────────────────────┐
│                 gRPC Server                         │
├─────────────────────────────────────────────────────┤
│                  Middleware Layer                   │
│  ┌─────────────┐ ┌──────────────┐ ┌──────────────┐  │
│  │    Auth     │ │   Logging    │ │   Recovery   │  │
│  │ Middleware  │ │  Middleware  │ │  Middleware  │  │
│  └─────────────┘ └──────────────┘ └──────────────┘  │
├─────────────────────────────────────────────────────┤
│                  Handler Layer                      │
│  ┌─────────────┐ ┌──────────────┐ ┌──────────────┐  │
│  │    User     │ │   Product    │ │   Future     │  │
│  │   Handler   │ │   Handler    │ │   Handler    │  │
│  └─────────────┘ └──────────────┘ └──────────────┘  │
├─────────────────────────────────────────────────────┤
│                  Service Layer                      │
│  ┌─────────────┐ ┌──────────────┐ ┌──────────────┐  │
│  │    User     │ │   Product    │ │   Future     │  │
│  │   Service   │ │   Service    │ │   Service    │  │
│  └─────────────┘ └──────────────┘ └──────────────┘  │
├─────────────────────────────────────────────────────┤
│                Repository Layer                     │
│  ┌─────────────┐ ┌──────────────┐ ┌──────────────┐  │
│  │    User     │ │   Product    │ │   Future     │  │
│  │ Repository  │ │ Repository   │ │ Repository   │  │
│  └─────────────┘ └──────────────┘ └──────────────┘  │
├─────────────────────────────────────────────────────┤
│                  Database Layer                     │
│              ┌──────────────────┐                   │
│              │   PostgreSQL     │                   │
│              └──────────────────┘                   │
└─────────────────────────────────────────────────────┘
```

## Design Patterns

### 1. Clean Architecture

The application follows Clean Architecture principles:

- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Separation of Concerns**: Each layer has a specific responsibility
- **Independence**: Business logic is independent of frameworks and databases

### 2. Repository Pattern

- Abstracts data access logic
- Provides a uniform interface for data operations
- Enables easy testing with mock implementations
- Supports multiple data sources

### 3. Service Layer Pattern

- Encapsulates business logic
- Coordinates between different repositories
- Handles transaction management
- Provides a clean API for handlers

### 4. Middleware Pattern

- Implements cross-cutting concerns
- Provides modular and reusable components
- Enables request/response processing pipeline
- Supports authentication, logging, and error handling

## Component Details

### gRPC Server

**Location**: `api/grpc/server.go`

**Responsibilities**:
- Initialize gRPC server with middleware
- Register service handlers
- Handle graceful shutdown
- Configure server options

**Key Features**:
- Unary and streaming interceptors
- Reflection support for development
- Middleware chain configuration
- Health check integration

### Middleware Layer

#### Authentication Middleware
**Location**: `internal/middleware/auth.go`

**Responsibilities**:
- JWT token validation
- User context injection
- Public endpoint filtering
- Authorization header processing

#### Logging Middleware
**Location**: Uses `grpc-ecosystem/go-grpc-middleware/logging/logrus`

**Responsibilities**:
- Request/response logging
- Performance metrics
- Error tracking
- Structured logging

#### Recovery Middleware
**Location**: Uses `grpc-ecosystem/go-grpc-middleware/recovery`

**Responsibilities**:
- Panic recovery
- Error response generation
- Graceful error handling
- System stability

### Handler Layer

**Location**: `internal/handler/grpc/`

**Responsibilities**:
- gRPC request/response handling
- Input validation
- Protocol buffer conversion
- Error status code mapping

**Pattern**:
```go
func (h *Handler) Method(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // 1. Convert proto request to domain model
    // 2. Call service method
    // 3. Handle errors
    // 4. Convert domain model to proto response
    // 5. Return response
}
```

### Service Layer

**Location**: `internal/service/`

**Responsibilities**:
- Business logic implementation
- Data validation
- Transaction coordination
- Integration with external services

**Pattern**:
```go
type Service interface {
    Method(req *model.Request) (*model.Response, error)
}

type serviceImpl struct {
    repo Repository
    // other dependencies
}
```

### Repository Layer

**Location**: `internal/repository/`

**Responsibilities**:
- Data persistence
- Database queries
- Transaction management
- Data mapping

**Pattern**:
```go
type Repository interface {