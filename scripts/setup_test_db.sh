#!/bin/bash

# Exit on error
set -e

echo "Setting up test database..."

# Check if PostgreSQL is running in Docker
if docker ps | grep -q meeting-scheduler-db; then
  echo "Using PostgreSQL from Docker..."
  
  # Create test database
  docker exec meeting-scheduler-db psql -U postgres -c "DROP DATABASE IF EXISTS scheduler_test;" || true
  docker exec meeting-scheduler-db psql -U postgres -c "CREATE DATABASE scheduler_test;"
  
  echo "Test database created successfully in Docker."
  
  # Set environment variables for tests
  export TEST_DB_HOST=localhost
  export TEST_DB_PORT=5432
  export TEST_DB_USER=postgres
  export TEST_DB_PASSWORD=postgres
  export TEST_DB_NAME=scheduler_test
  
  echo "Environment variables set for Docker PostgreSQL."
else
  echo "Using local PostgreSQL..."
  
  # Check if psql is available
  if ! command -v psql &> /dev/null; then
    echo "Error: PostgreSQL client (psql) not found. Please install PostgreSQL."
    exit 1
  fi
  
  # Create test database
  PGPASSWORD=postgres psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS scheduler_test;" || true
  PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE scheduler_test;"
  
  echo "Test database created successfully locally."
  
  # Set environment variables for tests
  export TEST_DB_HOST=localhost
  export TEST_DB_PORT=5432
  export TEST_DB_USER=postgres
  export TEST_DB_PASSWORD=postgres
  export TEST_DB_NAME=scheduler_test
  
  echo "Environment variables set for local PostgreSQL."
fi

echo "Test database setup complete." 