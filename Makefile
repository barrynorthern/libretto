# Simple dev targets for local workflow
# Make targets for local development and CI-friendly workflows

help:
	@echo "Available targets:"
	@echo "  build       - bazel build //..."
	@echo "  test        - bazel test //... --test_output=errors"
	@echo "  dev-up      - start monolith (API + agents + store)"
	@echo "  dev-smoke   - run smoke checks against monolith"
	@echo "  dev-down    - stop local services (best-effort)"
	@echo "Environment overrides: API_PORT"
	@echo "Examples: API_PORT=8090 make dev-up"



.PHONY: build test dev-up dev-smoke dev-down

build:
	bazel build //...

test:
	bazel test //... --test_output=errors

dev-up:
	./scripts/dev_up.sh

dev-smoke:
	./scripts/dev_smoke.sh



dev-down:
	pkill -f bazel-bin/services/api/api_/api || true
	pkill -f bazel-bin/services/agents/plotweaver/plotweaver_/plotweaver || true
	pkill -f bazel-bin/services/graphwrite/graphwrite_/graphwrite || true
	@echo "Stopped local services (best effort)."

