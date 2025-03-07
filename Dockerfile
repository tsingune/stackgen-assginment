# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata wget

# Set timezone
ENV TZ=UTC

# Copy the binary from builder
COPY --from=builder /app/main .

# Create a non-root user
RUN adduser -D appuser
USER appuser

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 