# 00006 – Consume SceneProposalReady → GraphWrite.Apply (persistence step 1)

Status: Won't Do
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
We now have Plot Weaver consuming DirectiveIssued and emitting a typed SceneProposalReady. To complete the MVP vertical slice, we need to persist a Scene node via GraphWrite when a SceneProposalReady arrives, proving the bus→agent→GraphWrite path. We'll do this with a minimal consumer that translates SceneProposalReady into one or more GraphWrite deltas and calls GraphWrite.Apply.

## Goal
Implement a minimal consumer for SceneProposalReady that constructs a GraphWrite ApplyRequest to create a Scene, then invokes GraphWrite.Apply using the existing connect-go client. Keep this local-only (DevPush/NOP) without real Pub/Sub yet.

## Scope
- Contracts
  - Reuse libretto.events.v1.Event (Envelope + oneof) and SceneProposalReady payload
  - No schema changes in this ticket
- New service: Narrative Ingest (temporary name) or extend Plot Weaver for MVP
  - Option A (preferred): Add a tiny agent service `narrative-ingest` with `/push` that listens for SceneProposalReady and calls GraphWrite
  - Option B: Temporarily add the consumer to Plot Weaver to shorten the path; extract later
- Consumer behavior
  - Base64 decode and protojson unmarshal Event
  - On SceneProposalReady, build a GraphWrite ApplyRequest with a single Delta:
    - op: "create"
    - entity_type: "Scene"
    - entity_id: scene_id from event
    - fields: { "title": title, "summary": summary }
  - Call GraphWrite.Apply via connect-go client to http://localhost:${GRAPHWRITE_PORT:-8082}
  - Log correlationId, causationId, graphVersionId returned
- Configuration
  - GRAPHWRITE_URL env (default http://localhost:8082)
- Tests
  - Unit test: valid SceneProposalReady → client.Apply called with expected request
  - Unit test: invalid payload → 400; no Apply call
- Tooling/CI
  - Extend make matrix/dev smoke to exercise the flow end-to-end locally (DevPush)

## Acceptance criteria
- Running dev smoke with a DirectiveIssued results in:
  - Plot Weaver consumes and emits SceneProposalReady
  - Narrative Ingest (or Plot Weaver extended) consumes SceneProposalReady and calls GraphWrite.Apply
  - Logs show correlationId continuity and a non-empty graphVersionId in response
- Unit tests cover apply/no-apply paths and validation failures

## Non-functional requirements
- Keep < 60s matrix run time
- Clear logs for tracing correlationId/causationId

## Risks / mitigations
- Service sprawl: start with Option B (extend Plot Weaver) if needed, extract to `narrative-ingest` in a follow-up
- Schema leakage: keep the mapper from event → deltas isolated in a small function with tests

## Out of scope (deferred)
- Real Pub/Sub client/emulator
- Firestore persistence implementation (GraphWrite remains in-memory)
- Thematic Steward participation

