#!/bin/bash

set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
MIGRATIONS_DIR="$PROJECT_DIR/migrations"
MIGRATIONS_DIR=$(echo "$MIGRATIONS_DIR" | sed 's|^/\([a-zA-Z]\)/|\U\1:/|')
echo "Running database migrations..."
docker run --rm -v "$MIGRATIONS_DIR:/migrations" --network host migrate/migrate \
  -path=/migrations/ \
  -database "postgres://user:password@localhost:5432/orders?sslmode=disable" up

echo "Migrations completed."