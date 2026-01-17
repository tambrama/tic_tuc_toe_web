# ===== build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ["../go.mod", "../go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o app ./cmd/app

FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache ca-certificates

# 🔑 ВОТ ЗДЕСЬ ГЛАВНОЕ ИЗМЕНЕНИЕ
COPY --from=builder /app/app .

EXPOSE 8081

CMD ["./app"]
