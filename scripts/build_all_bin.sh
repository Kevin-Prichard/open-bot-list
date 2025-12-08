#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")/.."
BASE_DIR="$(pwd)"
BUILD_DIR="${BASE_DIR}/build"

mkdir -p "$BUILD_DIR"

rm -f "$BUILD_DIR"/*

APP_PREFIX="open-bot-list"

function compile() {
    app="${APP_PREFIX}-$1" os="$2" arch="$3"
    echo "COMPILING BINARY FOR ${os}-${arch}"
    GOOS="$os" GOARCH="$arch" go build -o "${BUILD_DIR}/${app}-${os}-${arch}" ./cmd/main.go
    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -o "${BUILD_DIR}/${app}-${os}-${arch}-CGO0" ./cmd/main.go
    if [[ "$os" == "windows" ]]
    then
        zip "${BUILD_DIR}/${app}-${os}-${arch}.zip" "${BUILD_DIR}/${app}-${os}-${arch}"
        zip "${BUILD_DIR}/${app}-${os}-${arch}-CGO0.zip" "${BUILD_DIR}/${app}-${os}-${arch}-CGO0"
    else
        tar -czf "${BUILD_DIR}/${app}-${os}-${arch}.tar.gz" "${BUILD_DIR}/${app}-${os}-${arch}"
        tar -czf "${BUILD_DIR}/${app}-${os}-${arch}-CGO0.tar.gz" "${BUILD_DIR}/${app}-${os}-${arch}-CGO0"
    fi
}

echo ''
echo '### DOWNLOADER ###'
echo ''

cd "${BASE_DIR}/apps/downloader"

#compile "linux" "386"
compile "downloader" "linux" "amd64"
#compile "linux" "arm"
#compile "linux" "arm64"

# untested
#compile "freebsd" "386"
#compile "freebsd" "amd64"
#compile "freebsd" "arm"

#compile "openbsd" "386"
#compile "openbsd" "amd64"
#compile "openbsd" "arm"

#compile "darwin" "amd64"
#compile "darwin" "arm64"

#compile "windows" "386"
#compile "windows" "amd64"

echo ''
echo '### LOG-FLAGGER ###'
echo ''

cd "${BASE_DIR}/apps/log_flagger"
compile "log_flagger" "linux" "amd64"
