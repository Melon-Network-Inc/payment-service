all:
	bazel clean //...
	gofmt -w -s pkg/ cmd/
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel build //...
.PHONY: all

build:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel build //...
.PHONY: all

test:
	bazel test //...
.PHONY: test

clean:
	bazel clean
.PHONY: clean

run:
	bazel run //cmd/server:server
.PHONY: run

fmt:
	gofmt -w -s pkg/ cmd/
.PHONY: fmt

gazelle:
	bazel run //:gazelle
.PHONY: gazelle

dependency:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
.PHONY: dependency