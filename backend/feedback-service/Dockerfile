# ---------- Builder ----------
FROM golang:1.25.1-alpine AS builder

WORKDIR /app/feedback-service

RUN apk add --no-cache git

# Copy service-specific go.mod/go.sum first
COPY feedback-service/go.mod feedback-service/go.sum ./
COPY shared/go.mod ../shared/go.mod
COPY shared/go.sum ../shared/go.sum

RUN go mod download

# Copy source
COPY feedback-service/ ./
COPY shared/ ../shared/

# Build binary
RUN go build -o feedback-service ./cmd/server

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash curl ca-certificates postgresql-client

COPY --from=builder /app/feedback-service/feedback-service ./feedback-service
COPY --from=builder /app/feedback-service/migrations ./migrations
COPY feedback-service/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
