# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go module files
COPY service2/go.mod service2/go.sum ./
# Copy shared models for replace directive
COPY pkg/models /pkg/models
# Copy shared utils package
COPY pkg/utils /pkg/utils

# Download dependencies
RUN go mod download

# Copy service source code
COPY service2/. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service2 .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/service2 .

# Expose port for REST API
EXPOSE 8080

# Set environment variables with defaults
ENV TEMPORAL_ADDRESS=temporal:7233
ENV TEMPORAL_NAMESPACE=default
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=password
ENV DB_NAME=crawler
ENV GIN_MODE=release

# Run the binary
CMD ["./service2"] 