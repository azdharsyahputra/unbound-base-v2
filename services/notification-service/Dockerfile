# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

# notification-service nggak pakai shared (dari go.mod yg lo kirim)
COPY services/notification-service ./services/notification-service

WORKDIR /app/services/notification-service

RUN go mod tidy && \
    go build -o /app/service ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/service .
COPY services/notification-service/.env .env

EXPOSE 8084
CMD ["./service"]
