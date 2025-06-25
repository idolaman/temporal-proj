# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy main module files (for replace directive)
COPY go.mod go.sum ./

# Copy broker module files first to leverage docker layer caching
COPY cmd/broker/go.mod ./cmd/broker/

# Copy the entire source tree
COPY . .

# Download dependencies for the broker sub-module
WORKDIR /app/cmd/broker
RUN go mod download

# Build the broker binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /broker .

# -------- Runtime image --------
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /broker ./

EXPOSE 8080
CMD ["./broker"] 