#!/bin/bash

export DIR=$(pwd)
cd /app/src

# Generate a swagger json
go install github.com/swaggo/swag/cmd/swag@v1.8.12
echo "Generating swagger json"
swag init

# Run the API
go run main.go
