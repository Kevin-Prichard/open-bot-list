#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

export MODE_TEST=1

BASE_DIR="$(pwd)"

echo ''
echo '### DOWNLOADER ###'
echo ''

cd "${BASE_DIR}/apps/downloader"
go run gotest.tools/gotestsum@latest --format pkgname ./...

echo ''
echo '### LOG-FLAGGER ###'
echo ''

cd "${BASE_DIR}/apps/log_flagger"
go run gotest.tools/gotestsum@latest --format pkgname ./...
