#!/usr/bin/env bash

# ./scripts/migrate.sh
# Build and run the standalone migration command

set -euo pipefail

# Compile migration utility
go build -o migrate-tool ./cmd/migrate

# Execute migrations
./migrate-tool

# Clean up binary
rm -f migrate-tool
