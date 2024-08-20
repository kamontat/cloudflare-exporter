#!/usr/bin/env bash

set -e

BUILD_ID="$1"
## https://github.com/goreleaser/goreleaser/blob/3e6d825c80268b1b795971e3bdc0bc7b5a769062/internal/pipe/gomod/gomod_proxy.go#L144
TARGET_DIST="$PWD/dist/proxy/$BUILD_ID"

if [[ -d "$TARGET_DIST" ]]; then
  cd "$TARGET_DIST"
  go generate ./...
else
  printf 'Cannot find proxy directory to run go generate'
fi
