# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install build dependencies if needed (e.g. for CGO, though we usually disable it)
# RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Final Stage
FROM alpine:3.19
WORKDIR /app

# Install basic certificates for HTTPS calls if needed
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/worker .

# Copy configuration files
COPY --from=builder /app/env ./env

# Copy migrations (required for auto-migration on startup)
COPY --from=builder /app/internal/db/migration ./internal/db/migration

# Expose the application port (defaults to 3001 in dev)
EXPOSE 3001

# Command to run the application
CMD ["/app/main"]
