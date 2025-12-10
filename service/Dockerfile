FROM golang:1.25.0-alpine AS builder
    WORKDIR /app

    COPY ./service/ ./

    RUN go build -o service ./cmd/app/main.go

FROM debian:bookworm

    WORKDIR /app
    RUN apt update && apt install make
    COPY --from=builder /app/service .
    COPY --from=builder /app/migrations ./migrations

    CMD ["./service"]

