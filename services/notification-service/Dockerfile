# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy shared ONLY if needed (notification-service tidak pakai shared)
COPY shared ./shared
COPY services/notification-service ./services/notification-service

WORKDIR /app/services/notification-service

RUN go mod tidy && go build -o /app/notification-service ./cmd/server

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/notification-service .
COPY services/notification-service/.env .env

EXPOSE 8084
CMD ["./notification-service"]
