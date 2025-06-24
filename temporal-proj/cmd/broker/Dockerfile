# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY temporal-proj/cmd/broker/go.mod ./
COPY temporal-proj /temporal-proj

RUN go mod download

WORKDIR /app/temporal-proj/cmd/broker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o broker .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/temporal-proj/cmd/broker/broker .

EXPOSE 8080
CMD ["./broker"] 