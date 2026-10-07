#!/usr/bin/env bash
# Runs the emulator with the --script flag.

set -e
go run ./src/ --script testdata/startup.txt