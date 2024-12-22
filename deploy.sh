#!/bin/bash

# Exit on error
set -e

# Log output to a file
exec >> /var/log/deploy.log 2>&1

# Variables
PROJECT_DIR="/var/www/go-chat"
GO_BINARY="go-chat"          # Name of your Go binary
SERVICE_PATH="/etc/systemd/system/go-chat.service"

# 1. Pull latest changes from Git repository
cd $PROJECT_DIR
git pull origin main

# 2. Download dependencies (if Go modules are used)
echo "Downloading dependencies..."
go mod tidy   # Ensures all dependencies are fetched

# 3. Build the Go application
echo "Building the Go application..."
go build -o $GO_BINARY main.go

# 4. Reload systemd and restart the service
echo "Restarting service..."
systemctl daemon-reload
systemctl restart go-chat

# 5. Restart nginx
echo "Restarting nginx..."
systemctl restart nginx

echo "Deployment completed successfully!"