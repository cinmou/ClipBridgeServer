#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/dist"
APP_NAME="clipbridge-server"
GOCACHE_DIR="${GOCACHE:-${ROOT_DIR}/.gocache}"

mkdir -p "${OUTPUT_DIR}"
mkdir -p "${GOCACHE_DIR}"
rm -rf "${OUTPUT_DIR:?}/"*

TARGETS=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

for target in "${TARGETS[@]}"; do
  GOOS="${target%% *}"
  GOARCH="${target##* }"
  EXT=""
  if [[ "${GOOS}" == "windows" ]]; then
    EXT=".exe"
  fi

  OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}${EXT}"
  echo "Building ${OUTPUT_NAME}"
  env GOCACHE="${GOCACHE_DIR}" CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go build -trimpath -ldflags="-s -w" -o "${OUTPUT_DIR}/${OUTPUT_NAME}" ./cmd/server
done

echo "Release binaries written to ${OUTPUT_DIR}"
