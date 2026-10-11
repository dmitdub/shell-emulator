#!/usr/bin/env bash
# Runs the emulator with a startup script that exercises
# all stage 4 commands.

set -e
go run ./src/ --vfs testdata/vfs.csv --script testdata/startup.txt