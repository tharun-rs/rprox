# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Environment variables for Go build
ENV CGO_ENABLED=0 GOOS=linux

# Create working directory inside container
WORKDIR /app

# Install git (for private module access if needed)
RUN apk add --no-cache git

# Copy go.mod, go.sum, and vendor directory
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy source code
COPY . .

# Build the application
RUN go build -mod=vendor -o rprox ./cmd/rprox

# Stage 2: Create a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/rprox .

# Expose application port
EXPOSE 8080

# Start the application
ENTRYPOINT ["./rprox"]
