FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git

# Копируем go.mod и go.sum ДО запуска go mod download
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN go build -o scoring-manager ./cmd/scoring-manager-service

# Финальный образ
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/scoring-manager ./scoring-manager

EXPOSE 50052
CMD ["./scoring-manager"]
