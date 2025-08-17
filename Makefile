# Simple dev targets for local workflow
# Make targets for local development and CI-friendly workflows

help:
	@echo "Available targets:"
	@echo "  build       - bazel build //..."
	@echo "  test        - bazel test //... --test_output=errors"
	@echo "  dev-up      - start API, Plot Weaver, GraphWrite (readiness checks)"
	@echo "  dev-smoke   - run smoke checks against local services"
	@echo "  matrix      - run smoke in NOP and PUBSUB modes"
	@echo "  dev-down    - stop local services (best-effort)"
	@echo "Environment overrides: API_PORT, PLOT_PORT, GRAPHWRITE_PORT, PUBSUB_ENABLED"
	@echo "Examples: API_PORT=8090 PLOT_PORT=8091 GRAPHWRITE_PORT=8092 make dev-up"
	@echo "          PUBSUB_ENABLED=true make dev-smoke"



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

