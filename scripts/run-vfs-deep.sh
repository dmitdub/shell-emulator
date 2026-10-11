#!/usr/bin/env bash
# Runs the emulator with a deeply nested VFS.

set -e
go run ./src/ --vfs testdata/vfs-deep.csv