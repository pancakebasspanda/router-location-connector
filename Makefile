NAME := router-location-connector
SRC := ./cmd/router-location-connector
TEST_PKG := ./integration-test

.PHONY: build-linux
build-linux:
	@echo "Building Linux executable"
	env GOOS=linux GOARCH=amd64 go build -o $(NAME) $(SRC)

.PHONY: build-macos
build-macos:
	@echo "Building macOS executable"
	env GOOS=darwin GOARCH=amd64 go build -o $(NAME) $(SRC)

.PHONY: integration-test
integration-test:
	@echo "Running integration tests"
	go test $(TEST_PKG) -v

.PHONY:
unit-test:
	@echo "Running unit tests"
	go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /integration-test) -race

.PHONY: all-tests
all-tests: integration-test unit-test

.PHONY: all
all: build-linux build-macos integration-test unit-test
