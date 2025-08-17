# 00003 – DevPush publisher and envelope validation (next)

Status: Proposed
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
We can issue directives (API) and exercise GraphWrite, but the async “API → bus → Plot Weaver” path is not yet truly end‑to‑end. ADR 0009 emphasizes clear seams and keeping orchestration concerns separate. Rather than introduce a heavy Pub/Sub emulator now, we can add a dev‑only publisher that simulates the bus by HTTP POSTing to Plot Weaver, while validating the event envelope on both sides.

## Goal
Introduce a development‑only DevPush publisher to simulate the bus locally and enforce JSON envelope validation pre‑publish (API) and post‑decode (Plot Weaver). Validate both NOP and DevPush modes via Make-based smoke/matrix.

## Scope
- API
  - Publisher selection via env enum: `PUBLISHER` ∈ {`nop`, `devpush`, `pubsub`}. Backward compatibility: if `PUBSUB_ENABLED=true`, treat as `PUBLISHER=pubsub`.
  - Implement `DevPushPublisher` that POSTs the JSON envelope to Plot Weaver’s push endpoint.
    - Config: `PLOT_WEAVER_URL` (default `http://localhost:${PLOT_PORT:-8081}/push`).
  - Validate the envelope against `schemas/events/envelope.schema.json` before publishing; return 400 on invalid payloads.
  - Log selection: `publisher=devpush|nop|pubsub topic=...` for visibility.
- Plot Weaver
  - `/push` handler: after base64 decode, validate the envelope against the same schema; return 400 if invalid, 200 otherwise.
  - Keep existing stub handler for local flows.
- Tooling / Make targets
  - Keep Make as the interface (scripts are implementation details).
  - Extend `make matrix` to include `PUBLISHER=devpush` case (in addition to default NOP).
  - Unit tests: API and Plot Weaver validation (valid/invalid envelopes).
  - Optional env gate `ENVELOPE_VALIDATE=1` (default on) to allow bypass if needed during debugging.
- CI
  - Extend smoke‑matrix job to add a DevPush case; retain log artifact.

## Acceptance criteria
- `make dev-up` brings the stack up; API logs show the selected publisher.
- `make matrix` runs at least NOP and DevPush cases; both pass.
  - In DevPush mode, issuing a directive results in Plot Weaver `/push` receiving the envelope (visible in logs).
- Invalid envelopes are rejected with 400 by API (pre‑publish) and by Plot Weaver (post‑decode).
- Unit tests cover envelope validation pass/fail paths.
- CI smoke‑matrix runs NOP and DevPush and uploads logs.

## Non‑functional requirements
- Matrix step completes in < 60s on CI runners under typical load.
- Validation errors include actionable messages (which field failed) in logs (not necessarily user‑facing).

## Risks / mitigations
- DevPush diverges from real Pub/Sub.
  - Mitigation: keep `pubsub` path available; plan follow‑up ticket for emulator or cloud‑based tests.
- Schema churn causing false negatives.
  - Mitigation: version envelope schema; pin version referenced by both API and Plot Weaver; add CI check to ensure schema exists.

## Out of scope (deferred)
- Real Pub/Sub client wiring/emulator
- Firestore persistence changes
- Orchestration wiring per ADR 0009 (AgentClient RPCs)

## References
- ADR 0009: docs/adr/0009-agent-orchestration-seams-and-langgraph-sidecar.md
- Envelope schema: schemas/events/envelope.schema.json
- Current publisher selection/logging: services/api
- Make targets: Makefile (dev-up, dev-smoke, matrix)

