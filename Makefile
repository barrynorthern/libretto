# Make targets for local development and CI-friendly workflows

.PHONY: help build test dev-up dev-smoke dev-down proto lint-proto wails-dev wails-build frontend-install frontend-build sqlc
.PHONY: test-unit test-integration test-all test-coverage test-watch
.PHONY: db-seed db-inspect db-clean dashboard
.PHONY: tools-build tools-install monitoring-start monitoring-stop
.PHONY: docs-serve docs-build

# Default database for development
DB_FILE ?= libretto-dev.db
DB_PRESET ?= fantasy
DASHBOARD_PORT ?= 8080

help:
	@echo "Libretto Make targets"
	@echo ""
	@echo "Build & Development:"
	@echo "  build            - bazel build //..."
	@echo "  proto            - buf generate"
	@echo "  lint-proto       - buf lint"
	@echo "  sqlc             - Generate sqlc code"
	@echo "  wails-dev        - Run Wails dev for apps/desktop"
	@echo "  wails-build      - Build Wails desktop app"
	@echo "  frontend-install - Install frontend deps with pnpm"
	@echo "  frontend-build   - Build frontend"
	@echo ""
	@echo "Testing:"
	@echo "  test             - Run all tests (unit + integration)"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration test suite"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-watch       - Run tests in watch mode"
	@echo ""
	@echo "Database & Tools:"
	@echo "  db-seed          - Seed database with test data"
	@echo "  db-inspect       - Launch database inspection CLI"
	@echo "  db-clean         - Clean and reseed database"
	@echo "  dashboard        - Launch web dashboard"
	@echo "  tools-build      - Build all CLI tools"
	@echo "  tools-install    - Install CLI tools to PATH"
	@echo ""
	@echo "Monitoring & Docs:"
	@echo "  monitoring-start - Start monitoring dashboard"
	@echo "  monitoring-stop  - Stop monitoring services"
	@echo "  docs-serve       - Serve documentation locally"
	@echo "  docs-build       - Build documentation"
	@echo ""
	@echo "Environment variables:"
	@echo "  DB_FILE          - Database file path (default: libretto-dev.db)"
	@echo "  DB_PRESET        - Database preset (default: fantasy)"
	@echo "  DASHBOARD_PORT   - Dashboard port (default: 8080)"

build:
	bazel build //...

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
	sqlc generate

dev-up:
	./scripts/dev_up.sh

dev-smoke:
	./scripts/dev_smoke.sh

dev-down:
	pkill -f bazel-bin/services/api/api_/api || true
	pkill -f bazel-bin/services/agents/plotweaver/plotweaver_/plotweaver || true
	pkill -f bazel-bin/services/graphwrite/graphwrite_/graphwrite || true
	@echo "Stopped local services (best effort)."

# Testing targets
test: test-unit test-integration
	@echo "All tests completed successfully"

test-unit:
	@echo "Running unit tests..."
	go test ./internal/... -v -race

test-integration:
	@echo "Running integration tests..."
	go run cmd/integration-test/main.go -v -output test-results.json
	@echo "Integration test results saved to test-results.json"

test-all: test-unit test-integration
	@echo "Complete test suite finished"

test-coverage:
	@echo "Running tests with coverage..."
	go test ./internal/... -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-watch:
	@echo "Running tests in watch mode (requires entr)..."
	find . -name "*.go" | entr -c make test-unit

# Database and tools targets
db-seed:
	@echo "Seeding database: $(DB_FILE) with preset: $(DB_PRESET)"
	go run cmd/dbseed/main.go -db $(DB_FILE) -preset $(DB_PRESET) -clean

db-inspect:
	@echo "Launching database inspector for: $(DB_FILE)"
	@echo "Available commands: projects, entities, relationships, annotations, graph, stats"
	@echo "Example: make db-inspect-projects"
	go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd projects

db-inspect-projects:
	go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd projects -v

db-inspect-schema:
	go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd schema

db-inspect-stats:
	@if [ ! -f $(DB_FILE) ]; then echo "Database $(DB_FILE) not found. Run 'make db-seed' first."; exit 1; fi
	@echo "Getting project ID..."
	$(eval PROJECT_ID := $(shell go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd projects | tail -n +3 | head -n 1 | cut -d' ' -f1))
	@if [ -z "$(PROJECT_ID)" ]; then echo "No projects found. Run 'make db-seed' first."; exit 1; fi
	go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd stats -project $(PROJECT_ID)

db-inspect-graph:
	@if [ ! -f $(DB_FILE) ]; then echo "Database $(DB_FILE) not found. Run 'make db-seed' first."; exit 1; fi
	@echo "Getting project ID..."
	$(eval PROJECT_ID := $(shell go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd projects | tail -n +3 | head -n 1 | cut -d' ' -f1))
	@if [ -z "$(PROJECT_ID)" ]; then echo "No projects found. Run 'make db-seed' first."; exit 1; fi
	go run cmd/dbinspect/main.go -db $(DB_FILE) -cmd graph -project $(PROJECT_ID)

db-clean: 
	@echo "Cleaning and reseeding database..."
	rm -f $(DB_FILE)
	$(MAKE) db-seed

dashboard:
	@if [ ! -f $(DB_FILE) ]; then echo "Database $(DB_FILE) not found. Run 'make db-seed' first."; exit 1; fi
	@echo "Starting web dashboard on http://localhost:$(DASHBOARD_PORT)"
	@echo "Database: $(DB_FILE)"
	go run cmd/dashboard/main.go -db $(DB_FILE) -port $(DASHBOARD_PORT)

# Tools targets
tools-build:
	@echo "Building all CLI tools..."
	go build -o bin/dbinspect cmd/dbinspect/main.go
	go build -o bin/dbseed cmd/dbseed/main.go
	go build -o bin/dashboard cmd/dashboard/main.go
	go build -o bin/integration-test cmd/integration-test/main.go
	@echo "Tools built in ./bin/"

tools-install: tools-build
	@echo "Installing tools to PATH..."
	cp bin/* $(GOPATH)/bin/ 2>/dev/null || cp bin/* ~/go/bin/ 2>/dev/null || echo "Could not find Go bin directory"
	@echo "Tools installed. You can now use: dbinspect, dbseed, dashboard, integration-test"

# Monitoring targets
monitoring-start:
	@echo "Starting monitoring dashboard..."
	$(MAKE) dashboard &
	@echo "Dashboard started in background on port $(DASHBOARD_PORT)"

monitoring-stop:
	@echo "Stopping monitoring services..."
	pkill -f "cmd/dashboard/main.go" || true
	@echo "Monitoring services stopped"

# Documentation targets
docs-serve:
	@echo "Serving documentation on http://localhost:8000"
	@echo "Press Ctrl+C to stop"
	python3 -m http.server 8000 -d docs/ || python -m SimpleHTTPServer 8000

docs-build:
	@echo "Documentation is in markdown format in docs/"
	@echo "Main documentation files:"
	@ls -la docs/

# Development workflow targets
dev-setup: db-clean tools-build
	@echo "Development environment setup complete!"
	@echo "- Database seeded: $(DB_FILE)"
	@echo "- Tools built in ./bin/"
	@echo ""
	@echo "Quick start:"
	@echo "  make dashboard          # Launch web interface"
	@echo "  make test              # Run all tests"
	@echo "  make db-inspect-graph  # View narrative graph"

dev-test: test-unit db-clean test-integration
	@echo "Development test cycle complete"

# CI/CD targets
ci-test: test-unit test-integration
	@echo "CI test suite completed"

ci-build: build tools-build
	@echo "CI build completed"

