#!/usr/bin/env bash

if [[ -z "$MODE_TEST" ]]
then
  MODE_TEST=0
fi

set -euo pipefail

cd "$(dirname "$0")/.."
BASE_DIR="$(pwd)"

PATH_OUT="${BASE_DIR}/build"
mkdir -p "$PATH_OUT"

echo ''
echo '### DOWNLOADER ###'
echo ''

cd "${BASE_DIR}/apps/downloader"
go build -o "${PATH_OUT}/downloader" ./cmd/main.go

echo ''
echo '### LOG-FLAGGER ###'
echo ''

cd "${BASE_DIR}/apps/log_flagger"
go build -o "${PATH_OUT}/log_flagger" ./cmd/main.go

echo ''
echo '### DONE ###'
echo ''

ls -l "$PATH_OUT"

echo ''
echo ''
