#!/usr/bin/env bash
# Runs the emulator with both flags.

set -e
go run ./src/ --vfs testdata/vfs.csv --script testdata/startup.txt