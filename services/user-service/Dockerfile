# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/auth-service
RUN go mod tidy && go build -o /app/service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/service .
COPY services/user-service/.env .env
EXPOSE 8082
CMD ["./service"]
