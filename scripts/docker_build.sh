#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/../docker"

docker build -f Dockerfile_downloader -t oxlorg/open-bot-list-downloader:latest --no-cache --network=host ..
