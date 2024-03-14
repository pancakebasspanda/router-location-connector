#!/bin/bash

# Execute Makefile to build Linux and macOS executables
make build-linux

# Run docker-compose up with detached mode
docker-compose up redis -d

# Run Linux executable
./router-location-connector
