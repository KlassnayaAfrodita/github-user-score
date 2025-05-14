FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o score-service ./cmd/scoring-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/score-service .
CMD ["./score-service"]
