#!/usr/bin/env bash

set -e

printf '-------------------------\n'
printf 'current path: %s\n' "$PWD"
ls -la

printf '-------------------------\n'
printf 'module path: %s\n' "$1"
ls -la
