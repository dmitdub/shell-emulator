BINARY := shell-emulator
SRC    := ./src
TESTS  := ./tests

ifeq ($(OS),Windows_NT)
	BINARY := $(BINARY).exe
endif

.PHONY: all build run run-vfs run-script run-flags \
        run-vfs-minimal run-vfs-deep run-vfs-broken \
        run-commands run-cp \
        demo test test-v fmt vet clean help

all: build

build:
	@mkdir -p bin
	go build -o bin/$(BINARY) $(SRC)

run:
	go run $(SRC)

run-vfs:
	go run $(SRC) --vfs testdata/vfs.csv

run-script:
	go run $(SRC) --script testdata/startup.txt

run-flags:
	go run $(SRC) --vfs testdata/vfs.csv --script testdata/startup.txt

run-vfs-minimal:
	go run $(SRC) --vfs testdata/vfs-minimal.csv

run-vfs-deep:
	go run $(SRC) --vfs testdata/vfs-deep.csv

run-vfs-broken:
	go run $(SRC) --vfs testdata/vfs-broken.csv

run-commands:
	go run $(SRC) --vfs testdata/vfs.csv --script testdata/startup.txt

run-cp:
	go run $(SRC) --vfs testdata/vfs.csv --script testdata/startup.txt

demo:
	bash scripts/run-all.sh

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
	@echo "  build           Build the emulator into bin/"
	@echo "  run             Run the emulator from source"
	@echo "  run-vfs         Run with --vfs flag"
	@echo "  run-script      Run with --script flag"
	@echo "  run-flags       Run with both flags"
	@echo "  run-vfs-minimal Run with minimal VFS"
	@echo "  run-vfs-deep    Run with deeply nested VFS"
	@echo "  run-vfs-broken  Run with malformed VFS (error)"
	@echo "  run-commands    Run all stage 4 commands"
	@echo "  run-cp          Run stage 5 cp tests"
	@echo "  demo            Run all shell scripts in scripts/"
	@echo "  test            Run unit tests"
	@echo "  test-v          Run unit tests with verbose output"
	@echo "  fmt             Format source files with gofmt"
	@echo "  vet             Run go vet"
	@echo "  clean           Remove build artifacts"