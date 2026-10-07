#!/usr/bin/env bash
# Runs the emulator with the --vfs flag.

set -e
go run ./src/ --vfs testdata/vfs.csv