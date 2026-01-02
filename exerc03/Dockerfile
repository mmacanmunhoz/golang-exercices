# Multi-stage build for Docker CLI application
FROM golang:1.24-alpine AS builder

# Install git (may be needed for Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-cli ./main.go

# Final stage - use Docker-in-Docker base for Docker operations
FROM docker:25-dind

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/docker-cli .

# Copy configuration files
COPY deploy.yaml .

# Make binary executable
RUN chmod +x ./docker-cli

# Set entrypoint
ENTRYPOINT ["./docker-cli"]

# Default command
CMD ["--help"]