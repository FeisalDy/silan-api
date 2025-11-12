# Dockerfile
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o simple-go ./cmd/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/simple-go .
EXPOSE 8080
ENTRYPOINT ["./simple-go"]
