# Getting Started

This walkthrough shows how to run Libretto locally and exercise the current event-driven flow.

## Install prerequisites
- Go 1.22+
- Make
- Docker (for emulators in later steps)

## Run the services
- `make dev-up`
  - Starts API (8080), Plot Weaver (8081), GraphWrite (8082)
  - Select publisher via env: `PUBLISHER=nop|devpush|pubsub` (default nop)

## Issue a directive
- `make dev-smoke`
  - In DevPush mode, the API publishes a DirectiveIssued; Plot Weaver consumes it and emits SceneProposalReady.

## Verify behavior
- Look for logs like:
  - "plotweaver: consumed=DirectiveIssued published=SceneProposalReady correlationId=..."
  - "plotweaver publisher=devpush|nop|pubsub"

## What’s next
- Persistence: GraphWrite to Firestore emulator
- Minimal UI: Canvas/Inspector reads scenes

## Troubleshooting
- If `go test ./...` fails, see bug ticket 00005-bug for known issues to resolve in a separate branch.

