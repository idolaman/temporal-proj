# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go module files
COPY temporal-proj/cmd/broker-manager/go.mod temporal-proj/cmd/broker-manager/go.sum ./
# Copy shared packages for replace directives
COPY pkg/models /pkg/models
COPY pkg/utils /pkg/utils
COPY temporal-proj /temporal-proj

RUN go mod download

# Build broker-manager binary
WORKDIR /app/temporal-proj/cmd/broker-manager
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o broker-manager .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/temporal-proj/cmd/broker-manager/broker-manager .

EXPOSE 8080

ENV TEMPORAL_ADDRESS=temporal:7233
ENV TEMPORAL_NAMESPACE=default
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=password
ENV DB_NAME=crawler
ENV GIN_MODE=release

CMD ["./broker-manager"] 