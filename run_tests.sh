#!/bin/bash

# Run docker-compose up with detached mode
docker-compose up -d

# Run unit-tests
make unit-test

# Run integration-tests
make integration-test
