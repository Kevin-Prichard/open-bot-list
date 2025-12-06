#!/usr/bin/env bash

if [[ -z "$MODE_TEST" ]]
then
  MODE_TEST=0
fi

set -euo pipefail

cd "$(dirname "$0")/.."
BASE_DIR="$(pwd)"

PATH_OUT="$(pwd)/build"
mkdir -p "$PATH_OUT"

cd "${BASE_DIR}/downloader"
go build -o "${PATH_OUT}/downloader" ./cmd/main.go

echo ''
echo '### DONE ###'
echo ''

ls -l "$PATH_OUT"

echo ''
echo ''
