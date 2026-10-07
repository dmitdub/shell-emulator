#!/usr/bin/env bash
# Runs the non-interactive scripts in this directory in sequence.

set -e
cd "$(dirname "$0")/.."

for s in run-with-script.sh run-all-flags.sh run-missing-script.sh; do
  echo "--- $s ---"
  bash "scripts/$s" || true
  echo
done