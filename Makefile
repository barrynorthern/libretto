# Simple dev targets for local workflow

.PHONY: build test dev-up dev-smoke dev-down matrix

build:
	bazel build //...

test:
	bazel test //... --test_output=errors

dev-up:
	./scripts/dev_up.sh

dev-smoke:
	./scripts/dev_smoke.sh


matrix:
	./scripts/dev_matrix.sh

dev-down:
	pkill -f bazel-bin/services/api/api_/api || true
	pkill -f bazel-bin/services/agents/plotweaver/plotweaver_/plotweaver || true
	pkill -f bazel-bin/services/graphwrite/graphwrite_/graphwrite || true
	@echo "Stopped local services (best effort)."

