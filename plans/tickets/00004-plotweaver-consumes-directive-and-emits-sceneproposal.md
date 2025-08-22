# 00004 – Plot Weaver consumes DirectiveIssued and emits SceneProposalReady (oneof Event)

Status: Completed
Owner: barrynorthern
Start: 2025-08-19
Date completed: 2025-08-20

## Context
We now publish typed Events (Envelope + oneof payload) and validate envelopes. Plot Weaver’s /push validates and logs incoming Events but does not yet act on DirectiveIssued nor publish a SceneProposalReady. Completing this agent loop will prove the end‑to‑end async path using DevPush locally.

## Goal
Implement a minimal Plot Weaver consumer that decodes DirectiveIssued Events and publishes a typed SceneProposalReady Event with a valid Envelope. Validate envelopes pre‑publish (Plot Weaver) and post‑decode as per the current approach and keep the ENVELOPE_VALIDATE gate.

## Scope
- Protobuf/contracts
  - Reuse libretto.events.v1.Event (Envelope + oneof) with DirectiveIssued and SceneProposalReady
- Plot Weaver
  - Add a minimal publisher (reuse PUBLISHER env: nop|devpush|pubsub; default nop)
  - /push: base64 decode, unmarshal Event (protojson), validate Envelope when ENVELOPE_VALIDATE != "0"
  - On DirectiveIssued: synthesize a minimal SceneProposalReady and publish as typed Event with fresh Envelope (UUIDs; semver; producer=plotweaver; correlationId propagated)
  - Validate Envelope pre‑publish when ENVELOPE_VALIDATE != "0"; log consumed/published with IDs
  - Keep existing stub root handler
- Tooling
  - Keep Make as interface; rely on logs for smoke verification
- Tests
  - Plot Weaver unit tests:
    - Valid DirectiveIssued → publish SceneProposalReady (fake publisher)
    - Invalid envelope → 400; nothing published
- CI
  - Smoke‑matrix remains NOP, PUBSUB back‑compat, and DEVPUSH; verify via logs that SceneProposalReady is emitted in DevPush case

## Acceptance criteria
- In DevPush mode: issuing a directive leads to Plot Weaver logging that it consumed DirectiveIssued and published SceneProposalReady (with correlationId continuity)
- Invalid envelopes rejected with 400 on /push; pre‑publish validation in Plot Weaver rejects invalid outbound Envelopes when enabled
- Unit tests cover publish/no‑publish paths and validation failures
- CI smoke‑matrix runs all three cases and uploads logs

## Non‑functional requirements
- Matrix step completes in < 60s on CI
- Validation errors log actionable messages

## Risks / mitigations
- Divergence from real Pub/Sub push semantics → retain nop/devpush; emulator/cloud‑based tests later
- Overfitting to dev behavior → keep publisher interface narrow; clear dev‑only boundaries

## Out of scope (deferred)
- Real Pub/Sub client/emulator
- Thematic Steward behavior
- Persistence downstream of SceneProposalReady

## Outcomes / Notes
- Plot Weaver `/push` now base64-decodes, protojson-unmarshals libretto.events.v1.Event, and switches on oneof payload.
- On DirectiveIssued, it emits a typed SceneProposalReady with a fresh Envelope, propagating correlationId and setting causationId to the incoming eventId.
- Envelope validation is enforced when ENVELOPE_VALIDATE != "0" for both inbound decode and outbound publish using packages/contracts/events.
- Publisher selection added (nop|devpush|pubsub) via PUBLISHER env; current implementations log distinctly per mode.
- Unit tests cover invalid/valid push payloads; Plot Weaver package tests pass locally.
- Dev smoke verifies logs: consumed DirectiveIssued and published SceneProposalReady with correlationId continuity.
- Follow-up: align DevPush semantics across services (API vs agents) and expand CI matrix to include pubsub placeholder once available.

## References
- ADR 0009: docs/adr/0009-agent-orchestration-seams-and-langgraph-sidecar.md
- ADR 0010: docs/adr/0010-event-contract-source-of-truth.md
- Protos: proto/libretto/events/v1/events.proto
- Make targets: Makefile (dev-up, dev-smoke, matrix)

