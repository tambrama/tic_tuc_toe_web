# ===== build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN go build -o app ./cmd/app

# ===== runtime stage =====
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY src/migrations ./migrations

# goose
RUN apk add --no-cache curl \
    && curl -fsSL https://github.com/pressly/goose/releases/download/v3.21.1/goose_linux_x86_64 \
    -o /usr/local/bin/goose \
    && chmod +x /usr/local/bin/goose

EXPOSE 8081

CMD ["./app"]
