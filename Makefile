BINARY := shell-emulator
SRC    := ./src
TESTS  := ./tests

ifeq ($(OS),Windows_NT)
	BINARY := $(BINARY).exe
endif

.PHONY: all build run test test-v fmt vet clean help

all: build

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) $(SRC)

run:
	go run $(SRC)

test:
	go test $(TESTS)

test-v:
	go test -v $(TESTS)

fmt:
	gofmt -l -w src tests

vet:
	go vet ./...

clean:
	rm -rf bin

help:
	@echo "Available targets:"
	@echo "  build   Build the emulator into bin/"
	@echo "  run     Run the emulator from source"
	@echo "  test    Run unit tests"
	@echo "  test-v  Run unit tests with verbose output"
	@echo "  fmt     Format source files with gofmt"
	@echo "  vet     Run go vet"
	@echo "  clean   Remove build artifacts"