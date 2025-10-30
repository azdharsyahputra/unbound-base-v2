# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/service .
COPY .env .env
EXPOSE 8080
CMD ["./service"]
