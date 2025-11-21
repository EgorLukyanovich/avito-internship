# Этап сборки
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o avito_app ./cmd

# Финальный образ
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/avito_app .
COPY internal/migrations ./internal/migrations

ENV TZ=UTC

CMD ["./avito_app"]
