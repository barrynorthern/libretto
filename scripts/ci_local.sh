#!/usr/bin/env bash
set -euo pipefail

# Ensure Bazelisk is used and Bazel version matches repo pin
bazel version

# Keep BUILD files in sync with imports when developing
# Uncomment locally if you want Gazelle to update BUILD files automatically
# bazel run //:gazelle

# Build everything like CI
bazel build //...

# Run tests like CI (show errors)
bazel test //... --test_output=errors

