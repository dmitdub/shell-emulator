#!/usr/bin/env bash
# Runs the emulator with a startup script that exercises cp.

set -e
go run ./src/ --vfs testdata/vfs.csv --script testdata/startup.txt