#!/usr/bin/env bash

set -e

cd "$(dirname "$0")/.."

ENDPOINT="https://api.cloudflare.com/client/v4/graphql"
OUTPUT="$PWD/cloudflare/cloudflare.schema.graphql"
# CF_API_TOKEN=

if [ -f "$PWD/config.env" ]; then
  # shellcheck source=/dev/null
  source "$PWD/config.env"
fi

npx -p @apollo/rover \
  rover graph introspect "$ENDPOINT" \
  --output "$OUTPUT" \
  --format plain \
  --header "Authorization: Bearer $CF_API_TOKEN"
