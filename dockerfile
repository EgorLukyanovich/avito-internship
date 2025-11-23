# Этап сборки
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -o avito_app ./cmd

# Финальный образ
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/avito_app .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY internal/migrations ./internal/migrations

ENV TZ=UTC
ENV DB_MIGRATION_PATH=/app/internal/migrations

CMD goose -dir "$DB_MIGRATION_PATH" postgres "$DATABASE_URL" up && ./avito_app
