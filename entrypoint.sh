#!/bin/bash
set -e

# Check if the SSH root password secret exists and update the root password
if [ -f /run/secrets/ssh_root_password ]; then
  SSH_ROOT_PASSWORD=$(cat /run/secrets/ssh_root_password)
  echo "Updating root password..."
  echo "root:${SSH_ROOT_PASSWORD}" | chpasswd
else
  echo "Warning: SSH root password secret not found. Using existing root password."
fi

# Execute the command passed to the container (typically starting the SSH daemon)
exec "$@"
