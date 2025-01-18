#!/bin/sh

set -e

host=$DB_HOST
port=$DB_PORT

echo "Waiting for database at $host:$port..."

until nc -z "$host" "$port"; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is up - executing command"

# Pass all arguments to the command (i.e., your Go application)
exec "$@"