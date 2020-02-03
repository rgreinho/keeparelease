#!/bin/bash
set -euo pipefail

BINARY_NAME=fetch_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64
: ${BINARY_VERSION:=0.3.7}
FETCH_RELEASE_URL=https://github.com/gruntwork-io/fetch/releases/download
curl -sLO "${FETCH_RELEASE_URL}/v${BINARY_VERSION}/${BINARY_NAME}"
chmod +x "${BINARY_NAME}"
mv "./${BINARY_NAME}" /usr/local/bin/fetch
