# Use Go base image
FROM golang:1.25-alpine

# Install required packages
RUN apk add --no-cache git bash

# Install Air (hot reload tool)
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy Go module files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Expose port
EXPOSE 8080

# Seed and run app with Air
CMD /bin/bash -c "go run ./cmd/seed/main.go && air"
