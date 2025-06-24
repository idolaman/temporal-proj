# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY service1/go.mod service1/go.sum ./
# Copy shared models so replace works
COPY pkg/models /pkg/models
# Copy shared utils package
COPY pkg/utils /pkg/utils

# Download dependencies
RUN go mod download

# Copy source code
COPY service1/. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service1 .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/service1 .

# Set environment variables with defaults
ENV TEMPORAL_ADDRESS=localhost:8233
ENV TEMPORAL_NAMESPACE=default

# Run the binary
CMD ["./service1"] 