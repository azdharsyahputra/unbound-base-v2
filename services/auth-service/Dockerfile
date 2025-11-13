# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

# Root kerja dalam container
WORKDIR /app

# Modul shared (untuk proto / common)
COPY shared ./shared

# Source auth-service
COPY services/auth-service ./services/auth-service

# Masuk ke modul auth-service
WORKDIR /app/services/auth-service

# Download dep + build
RUN go mod tidy && \
    go build -o /app/service ./cmd/server

# ============================
# Final image
# ============================
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/service .
COPY services/auth-service/.env .env

EXPOSE 8081
CMD ["./service"]
