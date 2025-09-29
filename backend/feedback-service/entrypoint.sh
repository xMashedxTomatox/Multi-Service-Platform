#!/bin/bash
set -e

echo "🚀 Starting auth-service container..."

# Wait until Postgres is ready
echo "⏳ Waiting for Postgres..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$POSTGRES_USER"; do
  sleep 2
done

echo "✅ Postgres is ready!"

# Run migrations
for f in ./migrations/*.sql; do
  echo "📦 Running migration: $f"
  PGPASSWORD=$POSTGRES_PASSWORD psql \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    -f "$f"
done

echo "✅ All migrations applied!"

# Start the feedback-service binary
echo "🚀 Starting feedback-service..."
exec ./feedback-service
