.PHONY: all
all:
	bazel clean //...
	gofmt -w -s pkg/ cmd/
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel build //...

.PHONY: all
build:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
	bazel build //...

.PHONY: test
test:
	bazel test //...

.PHONY: clean
clean:
	bazel clean

.PHONY: run
run:
	bazel run //cmd/server:server

.PHONY: fmt
fmt:
	gofmt -w -s pkg/ cmd/

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

.PHONY: dependency
dependency:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies