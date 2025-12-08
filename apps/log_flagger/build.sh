#!/bin/bash

set -euo pipefail
cd "$(dirname "$0")"

go build -o ../../../build/testdata-log-flagger cmd/main.go
