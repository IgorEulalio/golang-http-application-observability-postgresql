#!/bin/sh

# Wait for Redis to initialize
echo "Waiting for Redis to initialize..."

# In this version, we'll run Redis in the background,
# then wait for it to become available.
redis-server &
REDIS_PID=$!

until redis-cli ping; do
  echo "Redis not ready, sleeping..."
  sleep 1
done

# Populate data
echo "Populating Redis data..."
redis-cli set 10 terraform
redis-cli set 11 golang
redis-cli set 12 aws

# Keep Redis running in the foreground
wait $REDIS_PID
