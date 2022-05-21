all:
	bazel build //...

test:
	bazel test //...

server:
	bazel run //cmd/server:server

fmt:
	gofmt -w -s pkg/ cmd/

gazelle:
	bazel run //:gazelle

dependency:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
