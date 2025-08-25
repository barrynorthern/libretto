# Getting Started

This guide helps you run the Libretto desktop scaffold and core tooling.

## Prerequisites
- Go 1.22+
- Node 18+ with corepack (pnpm)
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## Common tasks (Make)
Run `make help` to see available targets:

- proto: Run buf generate (protos -> gen/go, gen/ts)
- lint-proto: Lint protobufs
- build: Bazel build //...
- test: go/bazel tests
- wails-dev: Start Wails dev server for apps/desktop
- wails-build: Build the desktop app
- frontend-install: Install frontend deps with pnpm
- frontend-build: Build the frontend
- sqlc: Generate repository code once sqlc.yaml is added

## Desktop app (Wails)
1. Install frontend deps: `make frontend-install`
2. Start dev server: `make wails-dev`
3. Build release binary: `make wails-build`

## Notes
- DTOs across the UI boundary are protobuf-generated (baton.v1, scene.v1). The app currently uses a simple DTO struct for ListScenes() and will switch to generated types when repositories are wired.
- Offline-first: No external services are required to run the dev app.
