#!/usr/bin/env bash

set -euo pipefail

bazel run //:gazelle -- update-repos -from_file=aip-controller/go.mod -to_macro=bazel/deps.bzl%go_dependencies
bazel run //:gazelle -- update-repos -from_file=aip-wiki/go.mod -to_macro=bazel/deps.bzl%go_dependencies
bazel run //:gazelle -- update-repos -from_file=aip-psi/go.mod -to_macro=bazel/deps.bzl%go_dependencies
bazel run //:gazelle -- update-repos -from_file=aip-sdk/go.mod -to_macro=bazel/deps.bzl%go_dependencies
bazel run //:gazelle -- update-repos -from_file=aip-forddb/go.mod -to_macro=bazel/deps.bzl%go_dependencies
