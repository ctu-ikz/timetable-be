#!/bin/bash

set -e

# Initialize Docker Swarm if not already initialized
if ! docker info | grep -q 'Swarm: active'; then
  echo "Docker Swarm is not initialized. Initializing Swarm..."
  docker swarm init || { echo "Error: Failed to initialize Docker Swarm. Please check your Docker setup."; exit 1; }
else
  echo "Docker Swarm is already initialized."
fi

# Check if the .env file exists
if [ ! -f ".env" ]; then
  echo ".env file does not exist. Please create it with the required environment variables."
  exit 1
fi

POSTGRES_USER=$(grep '^DB_USER=' .env | cut -d '=' -f2-)
POSTGRES_PASSWORD=$(grep '^DB_PASSWORD=' .env | cut -d '=' -f2-)
POSTGRES_DB=$(grep '^DB_NAME=' .env | cut -d '=' -f2-)

if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ] || [ -z "$POSTGRES_DB" ]; then
  echo "Error: One or more environment variables (DB_USER, DB_PASSWORD, DB_NAME) are missing in the .env file."
  exit 1
fi

# Ensure secrets directory exists
mkdir -p ./secrets

# Clean up the secrets files
rm -f ./secrets/postgres_user.txt ./secrets/postgres_password.txt ./secrets/postgres_db.txt

# Remove existing Docker secrets, ignoring errors if they don't exist
docker secret rm postgres_user postgres_password postgres_db || true

echo "$POSTGRES_USER" > ./secrets/postgres_user.txt
echo "$POSTGRES_PASSWORD" > ./secrets/postgres_password.txt
echo "$POSTGRES_DB" > ./secrets/postgres_db.txt

echo "Creating Docker secrets..."
docker secret create postgres_user ./secrets/postgres_user.txt
docker secret create postgres_password ./secrets/postgres_password.txt
docker secret create postgres_db ./secrets/postgres_db.txt

# Stop and rebuild the containers
echo "Stopping existing containers..."
docker-compose down

echo "Building and starting the containers..."
docker-compose up --build -d || { echo "Error: Failed to build/start containers."; exit 1; }

echo "Setup completed successfully!"