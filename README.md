# Meeting Scheduler API

A RESTful API built with Go for scheduling meetings across different time zones. This service helps distributed teams find the optimal meeting time by collecting participant availability and suggesting the best possible time slots.

## Features

- Create events with multiple time slot options
- Collect participant availability
- Calculate optimal meeting times
- RESTful API design
- PostgreSQL for data persistence
- Docker support

## Prerequisites

- Go 1.21+ (for local development)
- Docker and Docker Compose (for containerized deployment)
- Make (optional, for convenience commands)

## Running the Application

### Using Docker (Recommended)

The easiest way to run the application is using Docker Compose:

```bash
# Start the application and database
make docker-up

# View logs
make docker-logs

# Stop the application and database
make docker-down
```

This will start:
- The Meeting Scheduler API on http://localhost:8080
- PostgreSQL database on port 5432
- pgAdmin (PostgreSQL admin interface) on http://localhost:5050
  - Email: admin@example.com
  - Password: admin

### Local Development

For local development without Docker:

1. Make sure PostgreSQL is running and create a database:
```bash
createdb scheduler
```

2. Run the application:
```bash
make run
```

## Testing the API

### Using curl

You can test the API endpoints using curl commands:

```bash
# Health check
curl http://localhost:8080/health

# Create a participant
curl -X POST http://localhost:8080/api/v1/participants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'

# Create an event
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "organizer_id": 1,
    "duration": 60
  }'

# Add a time slot
curl -X POST http://localhost:8080/api/v1/events/1/timeslots \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 1,
    "start_time": "2023-06-01T10:00:00Z",
    "end_time": "2023-06-01T11:00:00Z"
  }'

# Submit availability
curl -X POST http://localhost:8080/api/v1/events/1/availability \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": 1,
    "time_slot_id": 1,
    "is_available": true
  }'

# Get recommendations
curl -X GET http://localhost:8080/api/v1/events/1/recommendations
```

### Using the Test Script

A test script is provided to test all API endpoints in sequence:

```bash
# Make the script executable
chmod +x scripts/test_api.sh

# Run the test script
./scripts/test_api.sh
```

## Running Tests

### Running Tests with Docker

If you're using Docker, you can run the tests inside the Docker container:

```bash
# Start the application (if not already running)
make docker-up

# Run tests inside the Docker container
make test-in-docker
```

### Running Tests Locally

To run tests locally:

```bash
# Set up the test database
make setup-test-db

# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage report
make test-coverage
```

## API Endpoints

The API provides the following endpoints:

### Events
- `POST /api/v1/events` - Create a new event
- `GET /api/v1/events/{id}` - Get event details
- `PUT /api/v1/events/{id}` - Update an event
- `DELETE /api/v1/events/{id}` - Delete an event

### Time Slots
- `POST /api/v1/events/{id}/timeslots` - Add time slots to an event
- `GET /api/v1/events/{id}/timeslots` - Get time slots for an event

### Availability
- `POST /api/v1/events/{id}/availability` - Submit availability
- `GET /api/v1/events/{id}/recommendations` - Get time slot recommendations

### Participants
- `POST /api/v1/participants` - Create a participant
- `GET /api/v1/participants/{id}` - Get participant details

### Debug and Health
- `GET /health` - Health check endpoint
- `GET /debug/db` - Database connection check

## Useful Commands

```bash
# Build the application
make build

# Run the application locally
make run

# Start Docker environment
make docker-up

# Stop Docker environment
make docker-down

# View Docker logs
make docker-logs

# Connect to PostgreSQL CLI
make docker-exec-db
```

For a complete list of available commands:
```bash
make help
```

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entrypoint
├── internal/
│   ├── api/             # API handlers
│   ├── config/          # Configuration
│   ├── logger/          # Logging
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   └── repository/      # Database operations
├── docs/                # Documentation
├── scripts/             # Utility scripts
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── Makefile             # Build and run commands
└── README.md            # This file
```

## License

MIT License - see LICENSE file for details

## API Documentation

The API is documented using OpenAPI/Swagger. You can access the interactive API documentation at:

```
http://localhost:8080/docs
```

or directly at:

```
http://localhost:8080/swagger/index.html
```

This provides a user-friendly interface to:
- View all available endpoints
- See request and response schemas
- Test API calls directly from the browser

### Installing Swagger Dependencies

If you're developing the API and need to install the Swagger dependencies:

```bash
make swagger-deps
```