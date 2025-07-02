# gRPC Example with PostgreSQL

This repository contains a minimal gRPC server written in Go. It stores data in PostgreSQL and demonstrates typical middleware such as authentication and logging. The protocol definitions live in the `proto/` directory and generated Go files are kept next to them.

## Setup

1. **Clone the repository** and download dependencies:

   ```bash
   git clone <repo>
   cd grpc-go-postgre
   go mod download
   ```

2. **Generate gRPC code** from the Protocol Buffer definitions:

   ```bash
   ./scripts/gen-proto.sh
   ```

   The script invokes `protoc` for all files in `proto/` and writes the generated code next to the `.proto` files.

3. **Configure the application** by editing `configs/app.yaml`. An example environment file is provided in `.env.example`.

## Running the server

Run the gRPC server directly with Go:

```bash
go run cmd/server/main.go
```

The server listens on the port defined in `configs/app.yaml` (default `8080`).

You can also use the `Makefile` for common tasks such as `make run`, `make proto` or `make migrate`.
