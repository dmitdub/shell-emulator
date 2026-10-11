#!/usr/bin/env bash
# Runs the emulator and invokes vfs-init to reset the VFS.

set -e
cd "$(dirname "$0")/.."
cp testdata/vfs.csv testdata/vfs-tmp.csv

cat testdata/vfs-tmp.csv
echo "---"
echo "vfs-init" | go run ./src/ --vfs testdata/vfs-tmp.csv
echo "---"
cat testdata/vfs-tmp.csv

rm testdata/vfs-tmp.csv