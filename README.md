# MelonWallet Payment Service

---------------------
The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

## Compile and build

---------------------

```bash
bazel build //...
```

## Start payment server

---------------------

```bash
bazel run cmd/server:start
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
