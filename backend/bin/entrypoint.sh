#!/bin/bash

export DIR=$(pwd)
cd /app/src

# Generate a swagger json
go install github.com/swaggo/swag/cmd/swag@v1.16.4
swag fmt
swag init

# Run the API
go run main.go
