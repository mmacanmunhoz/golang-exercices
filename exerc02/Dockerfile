# Multi-stage build for Go application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o configparser-exerc02 ./main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/configparser-exerc02 .

# Copy example config for testing
COPY example_config.yaml .

# Expose port if needed (adjust based on your app)
# EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./configparser-exerc02"]

# Default command (can be overridden)
CMD ["parse", "example_config.yaml"]