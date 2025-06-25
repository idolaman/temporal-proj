# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy main module files (for replace directive)
COPY go.mod go.sum ./

# Copy broker-manager module files first for caching
COPY cmd/broker-manager/go.mod ./cmd/broker-manager/

# Copy the entire source tree
COPY . .

# Download dependencies for the broker-manager sub-module
WORKDIR /app/cmd/broker-manager
RUN go mod download

# Build broker-manager binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /broker-manager .

# -------- Runtime image --------
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /broker-manager ./

EXPOSE 8080

CMD ["./broker-manager"] 