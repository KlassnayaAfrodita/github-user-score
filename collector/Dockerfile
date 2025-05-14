FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o collector ./cmd/collector-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/collector .
EXPOSE 50051
CMD ["./collector"]