FROM golang:1.24-alpine AS builder

WORKDIR /app

# ============================
# 1. COPY SHARED FIRST (needed for replace path)
# ============================
COPY shared ./shared

# ============================
# 2. COPY go.mod & go.sum ONLY
# ============================
COPY services/chat-service/go.mod services/chat-service/go.sum ./services/chat-service/

WORKDIR /app/services/chat-service

# ============================
# 3. DOWNLOAD DEPENDENCIES (cached)
# ============================
RUN go mod download

# ============================
# 4. COPY SERVICE SOURCE CODE
# ============================
COPY services/chat-service /app/services/chat-service

# ============================
# 5. BUILD BINARY
# ============================
RUN go build -o /app/service ./cmd/server

# ============================
# FINAL STAGE
# ============================
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/service .
COPY services/chat-service/.env .env

EXPOSE 8083
CMD ["./service"]
