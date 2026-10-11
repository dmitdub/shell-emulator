#!/usr/bin/env bash
# Runs the emulator with a minimal VFS.

set -e
go run ./src/ --vfs testdata/vfs-minimal.csv