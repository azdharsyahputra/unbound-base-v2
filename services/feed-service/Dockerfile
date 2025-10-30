# syntax=docker/dockerfile:1
FROM rust:1.81 AS builder
WORKDIR /usr/src/app
COPY . .
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /usr/src/app/target/release/* .
COPY .env .env
EXPOSE 8080
CMD ["./post-service"]
