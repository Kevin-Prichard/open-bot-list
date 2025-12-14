#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")/.."
PATH_BASE="$(pwd)"
PATH_BUILD="${PATH_BASE}/build"

mkdir -p "$PATH_BUILD"

rm -f "$PATH_BUILD"/*

APP_PREFIX="open-bot-list"

function compile() {
    app="${APP_PREFIX}-$1" os="$2" arch="$3"

    cd "$SRC_DIR"
    echo "COMPILING BINARY FOR ${os}-${arch}"
    GOOS="$os" GOARCH="$arch" go build -o "${PATH_BUILD}/${app}-${os}-${arch}" ./cmd/main.go
    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -o "${PATH_BUILD}/${app}-${os}-${arch}-CGO0" ./cmd/main.go

    cd "$PATH_BUILD"
    if [[ "$os" == "windows" ]]
    then
        zip "./${app}-${os}-${arch}.zip" "./${app}-${os}-${arch}"
        zip "./${app}-${os}-${arch}-CGO0.zip" "./${app}-${os}-${arch}-CGO0"
    else
        tar -czf "./${app}-${os}-${arch}.tar.gz" "./${app}-${os}-${arch}"
        tar -czf "./${app}-${os}-${arch}-CGO0.tar.gz" "./${app}-${os}-${arch}-CGO0"
    fi
}

echo ''
echo '### DOWNLOADER ###'
echo ''

SRC_DIR="${PATH_BASE}/apps/downloader"

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

SRC_DIR="${PATH_BASE}/apps/log_flagger"
compile "log_flagger" "linux" "amd64"

echo "COMMAND TO REMOVE ALL NON-ARCHIVES: find ${PATH_BUILD} -type f ! -name '*.*' -delete"