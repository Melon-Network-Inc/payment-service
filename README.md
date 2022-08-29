# MelonWallet Payment Service

[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username=anuraghazra&repo=github-readme-stats)](https://github.com/anuraghazra/github-readme-stats)

The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

## Compile and build

---------------------

```bash
bazel build //...
```

## Start payment server

---------------------

```bash
bazel run cmd/server:server
```

## Update dependencies for Bazel build

---------------------

```bash
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
```

## Run tests

---------------------

```bash
bazel test //...
```
