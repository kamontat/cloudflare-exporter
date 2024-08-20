#!/usr/bin/env bash

set -e

printf '>-------------------------<\n'
printf 'current path: %s\n' "$PWD"
ls -la

printf '>-------------------------<\n'
path="$1"
printf 'path: %s\n' "$path"
[ -f "$path" ] && path="$(dirname "$1")"
cd "$path" && ls -la

printf '>-------------------------<\n'
path="$2"
printf 'module path: %s\n' "$path"
[ -f "$path" ] && path="$(dirname "$2")"
cd "$path" && ls -la
