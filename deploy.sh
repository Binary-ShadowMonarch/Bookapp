#!/usr/bin/env bash
# deploy.sh - Pull latest code and rebuild Docker containers

# Exit immediately if a command exits with a non-zero status,
# treat unset variables as an error, and propagate errors in pipelines
set -euo pipefail

# Determine script directory and cd into it
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "\n=>> Pulling latest changes from 'origin/master'..."
git pull origin master

echo "\n=>> Stopping containers..."
sudo docker compose down 

echo "\n=>> Pruning build cache..."

sudo docker builder prune -f
sudo docker buildx prune -f

echo "\n=>> Building and starting containers..."
sudo docker compose --env-file ./.env build --no-cache --pull && sudo docker compose --env-file ./.env up --force-recreate --build -d
echo "\n=>> Deployment complete!"

