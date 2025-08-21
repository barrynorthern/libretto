# Libretto User Guide

Welcome to Libretto. This guide helps you get started with the developer MVP: running the services locally, issuing a directive, and seeing agents respond.

## Quick start
1. Prerequisites: Go 1.22+, Make, Docker (for emulators, optional initially).
2. Clone and enter the repo.
3. Run the stack:
   - `make dev-up`
4. Issue a directive (replace text as desired):
   - `make dev-smoke`
5. Inspect logs for the agent flow and responses.

## Concepts (brief)
- Baton: issues directives from API.
- Event bus seam: DevPush simulates Pub/Sub locally.
- Plot Weaver: consumes DirectiveIssued and emits SceneProposalReady.
- GraphWrite: accepts deltas to persist narrative graph (emulator-backed in future steps).

## Next steps
- See plans/tickets for upcoming staircase steps toward the full vertical slice.
- Once Firestore emulator is enabled, you’ll be able to view persisted scenes via a minimal UI.

## Troubleshooting
- Port conflicts: API 8080, Plot Weaver 8081, GraphWrite 8082.
- Publisher selection: set `PUBLISHER=nop|devpush|pubsub`.
- Envelope validation: `ENVELOPE_VALIDATE=1` (default on); set to `0` to bypass during debugging.

