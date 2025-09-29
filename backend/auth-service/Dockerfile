# ---------- Builder ----------
FROM golang:1.25.1-alpine AS builder

WORKDIR /app/auth-service

# Install git (needed for go modules that use replace or private repos)
RUN apk add --no-cache git

# Copy service-specific go.mod/go.sum first for caching
COPY auth-service/go.mod auth-service/go.sum ./
COPY shared/go.mod ../shared/go.mod
COPY shared/go.sum ../shared/go.sum

# Download dependencies
RUN go mod download

# Copy the rest of the source (service + shared)
COPY auth-service/ ./
COPY shared/ ../shared/

# Build binary
RUN go build -o auth-service ./cmd/server

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache bash curl ca-certificates postgresql-client

# Copy binary + migrations + entrypoint
COPY --from=builder /app/auth-service/auth-service ./auth-service
COPY --from=builder /app/auth-service/migrations ./migrations
COPY auth-service/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
