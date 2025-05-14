FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN apk update && apk add git
RUN go mod download
RUN go build -o scoring-manager ./cmd/scoring-manager-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/scoring-manager .
EXPOSE 50052
CMD ["./scoring-manager"]
