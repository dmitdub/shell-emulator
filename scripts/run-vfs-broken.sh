#!/usr/bin/env bash
# Runs the emulator with a malformed VFS to check error handling.

set -e
go run ./src/ --vfs testdata/vfs-broken.csv