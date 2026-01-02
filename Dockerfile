# Platform Rocks - Multi-Project Dockerfile
# This Dockerfile can build any of the exercises based on build context

ARG PROJECT=exerc04

# Multi-stage build
FROM golang:1.25-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory based on project
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build based on project type
RUN case "$PROJECT" in \
    "exerc01") CGO_ENABLED=0 GOOS=linux go build -o app ./main.go ;; \
    "exerc02") CGO_ENABLED=0 GOOS=linux go build -o app ./main.go ;; \
    "exerc03") CGO_ENABLED=0 GOOS=linux go build -o app ./main.go ;; \
    "exerc04") CGO_ENABLED=0 GOOS=linux go build -o app ./main.go ;; \
    *) echo "Unknown project: $PROJECT" && exit 1 ;; \
    esac

# Final stage
FROM alpine:latest

ARG PROJECT

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Install Docker CLI for exerc03
RUN if [ "$PROJECT" = "exerc03" ]; then apk add --no-cache docker; fi

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/app .

# Copy configuration files based on project
COPY *.yaml . 2>/dev/null || true

# Make binary executable
RUN chmod +x ./app

# Set entrypoint
ENTRYPOINT ["./app"]

# Default command varies by project
CMD case "$PROJECT" in \
    "exerc01") ["parse", "--help"] ;; \
    "exerc02") ["parse", "--help"] ;; \
    "exerc03") ["--help"] ;; \
    "exerc04") ["k8s", "--help"] ;; \
    *) ["--help"] ;; \
    esac