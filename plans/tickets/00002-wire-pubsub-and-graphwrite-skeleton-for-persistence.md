# 002 – Wire Pub/Sub and GraphWrite skeleton for persistence (next)

Status: Done (merged)
Owner: barrynorthern
Start: TBC
Date completed: 2025-08-17

## Goal
Turn the stubbed event path into a real pipeline by introducing Pub/Sub publisher for API, Pub/Sub push handler for Plot Weaver, and a GraphWrite service skeleton that persists versions/deltas in-memory (or to Firestore if ready).

## Scope
- API: Replace NopPublisher with interface + Pub/Sub implementation behind env flag; keep NOP for local
- Agents: Plot Weaver HTTP handler to accept Pub/Sub push format; verify envelope schema
- GraphWrite: Implement Connect service with Apply method (proto exists); stub persistence with in-memory store and tests
- CI: Add unit tests for publisher + handler; basic contract tests validating event envelope against schema

## Acceptance criteria
- API publishes to configured Pub/Sub topic when enabled; falls back to NOP locally
- Plot Weaver accepts Pub/Sub push JSON and emits SceneProposalReady via publisher
- GraphWrite Apply returns new version id and applied count for simple create deltas
- Developer UX: Make-first workflow (dev-up, dev-smoke, matrix); CI smoke-matrix running NOP and PUBSUB paths

- All tests green (bazel build/test), smoke checks valid via Make targets

## Risks / mitigations
- Firestore readiness: start with in-memory, add Firestore in a follow-up
- Schema drift: add schema validation step for event envelope in CI

## References
- docs/graphwrite-api.md
- schemas/events/envelope.schema.json
- proto/libretto/graph/v1/graphwrite.proto

