# cypress-parallel-docker-images [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel-docker-images.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel-docker-images)

`cypress-parallel-docker-images` contain cypress docker image and cypress-parallel-docker-cli to run cypress unit testing in parallel.

## Docker image semantics

The docker image tag is defined as below:
- cypress docker image tag
- cypress-parallel-cli version

For cypress docker image `7.2.0` and cli `v0.0.1`, the result will be `docker.pkg.github.com/xxxx/xxxx/xxxx:7.2.0-0.0.1`

## Git hooks

Add githook like so:

```bash
git config core.hooksPath .githooks
```

## Linter
```bash
# https://golangci-lint.run/usage/install/
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```