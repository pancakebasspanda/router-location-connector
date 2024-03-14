#!/bin/bash

# Execute Makefile to build Linux and macOS executables
make build-macos

# Run docker-compose up with detached mode
docker-compose up redis -d

./router-location-connector