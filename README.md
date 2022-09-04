# MelonWallet Payment Service

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/569/badge)](https://bestpractices.coreinfrastructure.org/projects/569)
[![Build status](https://badge.buildkite.com/ae3cfe307da872f79cdb77b54afb3dfdb52bee3edccf6b367b.svg)](https://buildkite.com/melon-network-inc/payment-service-pipeline)

<img src="https://avatars.githubusercontent.com/u/104064333?s=400&u=fe08053ed0a72719e2ea4bb0229766ef9b4fdfee&v=4" width="100">

---------------------

The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

---------------------

## Project Setup

### Compile and build

```bash
bazel build //...
```

### Start payment server

```bash
bazel run cmd/server:server
```

### Update dependencies for Bazel build

```bash
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
```

### Run tests

```bash
bazel test //...
```
