FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy shared modules (proto, utils)
COPY shared ./shared

# Copy chat-service folder
COPY services/chat-service ./services/chat-service

# Move into chat-service
WORKDIR /app/services/chat-service

# Download dependencies + build binary
RUN go mod tidy && go build -o /app/service ./cmd/server

# -------------------------------------------
# Final image
# -------------------------------------------
FROM alpine:3.20
WORKDIR /app

# Copy built binary
COPY --from=builder /app/service .

# Copy .env
COPY services/chat-service/.env .env

EXPOSE 8083
CMD ["./service"]
