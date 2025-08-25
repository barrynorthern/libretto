# Make targets for local development and CI-friendly workflows

.PHONY: help build test dev-up dev-smoke dev-down proto lint-proto wails-dev wails-build frontend-install frontend-build sqlc

help:
	@echo "Libretto Make targets"
	@echo "  build            - bazel build //..."
	@echo "  test             - bazel test //... --test_output=errors"
	@echo "  proto            - buf generate"
	@echo "  lint-proto       - buf lint"
	@echo "  wails-dev        - Run Wails dev for apps/desktop (frontend + Go live reload)"
	@echo "  wails-build      - Build Wails desktop app"
	@echo "  frontend-install - Install frontend deps with pnpm (apps/desktop/frontend)"
	@echo "  frontend-build   - Build frontend (apps/desktop/frontend)"
	@echo "  sqlc             - Generate sqlc code (when sqlc.yaml present)"
	@echo "Environment overrides: API_PORT"

build:
	bazel build //...

test:
	bazel test //... --test_output=errors

proto:
	buf generate

lint-proto:
	buf lint

wails-dev:
	cd apps/desktop && wails dev

wails-build:
	cd apps/desktop && wails build

frontend-install:
	pnpm -C apps/desktop/frontend install

frontend-build:
	pnpm -C apps/desktop/frontend run build

sqlc:
	@echo "sqlc.yaml not yet present; will wire after repo scaffolding" && true

dev-up:
	./scripts/dev_up.sh

dev-smoke:
	./scripts/dev_smoke.sh

dev-down:
	pkill -f bazel-bin/services/api/api_/api || true
	pkill -f bazel-bin/services/agents/plotweaver/plotweaver_/plotweaver || true
	pkill -f bazel-bin/services/graphwrite/graphwrite_/graphwrite || true
	@echo "Stopped local services (best effort)."

