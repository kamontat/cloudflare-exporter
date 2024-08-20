# CONTRIBUTING

Contributions are always welcome, no matter how large or small.

## Set up

- Install golang
- Install dependencies by run `go mod tidy`
- If you change cloudflare/queries/*.gql files:
  - run: `go generate ./...`
- Configure githooks directory by run `git config core.hooksPath .githooks`
- Good to go:
  - run app: `go run $PWD`
