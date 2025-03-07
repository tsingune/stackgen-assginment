#!/bin/bash

# Make script exit on any error
set -e

echo "Starting Meeting Scheduler..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker is not running. Please start Docker and try again."
  exit 1
fi

# Check if port 8080 is already in use
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
  echo "Warning: Port 8080 is already in use. Stopping any containers using this port..."
  docker ps --format "{{.ID}}" --filter "publish=8080" | xargs -r docker stop
fi

# Clean up any existing containers
echo "Stopping any existing containers..."
docker-compose down 2>/dev/null || true

# Rebuild the application
echo "Building the application..."
docker-compose build app

# Start the database first
echo "Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
  if docker-compose exec postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "PostgreSQL is ready!"
    break
  fi
  echo "Waiting for PostgreSQL to start... ($i/30)"
  sleep 2
  if [ $i -eq 30 ]; then
    echo "Error: PostgreSQL failed to start in time."
    docker-compose logs postgres
    exit 1
  fi
done

# Start the application
echo "Starting application..."
docker-compose up -d app pgadmin

# Wait for application to be ready
echo "Waiting for application to be ready..."
for i in {1..30}; do
  if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "Application is ready!"
    echo "API is available at: http://localhost:8080"
    echo "pgAdmin is available at: http://localhost:5050"
    echo "  - Email: admin@example.com"
    echo "  - Password: admin"
    exit 0
  fi
  echo "Waiting for application to start... ($i/30)"
  sleep 2
  
  # Check if the container is still running
  if ! docker ps | grep -q meeting-scheduler-app; then
    echo "Error: Application container stopped unexpectedly."
    echo "Checking logs for errors:"
    docker-compose logs app
    exit 1
  fi
done

echo "Error: Application failed to start in time."
echo "Checking logs for errors:"
docker-compose logs app

# Print database connection info for debugging
echo "Database connection info:"
docker-compose exec postgres psql -U postgres -c "SELECT version();"
docker-compose exec postgres psql -U postgres -c "\\l"
docker-compose exec postgres psql -U postgres -c "\\du"

exit 1 