#!/bin/bash

if [ ! -d "./secrets" ]; then
  echo "Secrets directory does not exist. Creating it now."
  mkdir ./secrets
fi

if [ ! -f ".env" ]; then
  echo ".env file does not exist. Please create it with the required environment variables."
  exit 1
fi

POSTGRES_USER=$(grep -oP '^DB_USER=\K.*' .env)
POSTGRES_PASSWORD=$(grep -oP '^DB_PASSWORD=\K.*' .env)
POSTGRES_DB=$(grep -oP '^DB_NAME=\K.*' .env)

if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ] || [ -z "$POSTGRES_DB" ]; then
  echo "Error: One or more environment variables (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB) are missing in the .env file."
  exit 1
fi

echo "$POSTGRES_USER" > ./secrets/postgres_user.txt
echo "$POSTGRES_PASSWORD" > ./secrets/postgres_password.txt
echo "$POSTGRES_DB" > ./secrets/postgres_db.txt

echo "Creating Docker secrets..."
docker secret create postgres_user ./secrets/postgres_user.txt
docker secret create postgres_password ./secrets/postgres_password.txt
docker secret create postgres_db ./secrets/postgres_db.txt

rm ./secrets/postgres_user.txt
rm ./secrets/postgres_password.txt
rm ./secrets/postgres_db.txt

echo "Stopping existing containers..."
docker-compose down

echo "Building and starting the containers..."
docker-compose up --build -d

echo "Setup completed successfully!"
