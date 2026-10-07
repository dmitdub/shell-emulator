#!/usr/bin/env bash
# Runs the emulator with a missing startup script
# to check error handling.

set -e
go run ./src/ --script /no/such/file.txt