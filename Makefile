BINARY_NAME=gitwise
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build install clean test lint all

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe

test:
	go test ./...

lint:
	golangci-lint run

all: clean build
