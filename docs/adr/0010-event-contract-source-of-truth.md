# ADR 0010: Event Contracts — Protobuf as Canonical Source of Truth

Status: Accepted
Date: 2025-08-17
Owners: barrynorthern
Reviewers: Augment Agent
Supersedes: ADR 0005 (Event Schema and Idempotency Conventions)
Related: ADR 0008 (IDL & RPC), ADR 0009 (Agent Orchestration Seams)

## Context
We need a single, authoritative definition for event envelopes and payloads across services and agents. Earlier, ADR 0005 standardized on a JSON Schema envelope. Separately, ADR 0008 established Protobuf as the canonical IDL for services and event payloads. Maintaining hand-authored JSON Schemas alongside protobuf introduces drift risk and duplicated effort.

## Decision
- Protobuf is the canonical source of truth for event contracts:
  - Envelope is defined as a protobuf message (libretto.events.v1.Envelope).
  - Event payloads are protobuf messages (e.g., DirectiveIssued, SceneProposalReady).
- Serialization:
  - Use protojson for JSON encoding/decoding on the bus (Pub/Sub push) and HTTP.
  - Envelope and payload are carried together, using the Envelope plus a payload message appropriate to the topic.
- Validation:
  - Validation is performed via protobuf message decoding and application logic (required/semantic checks), not hand-maintained JSON Schemas.
  - JSON Schema may be generated from protobuf only when required by downstream tooling; generated schemas are not the canonical contract and are not hand-edited.
- Identifiers:
  - UUIDs are acceptable for eventId, correlationId, causationId, idempotencyKey (where applicable). ULIDs may be considered in a future ADR if sortable identifiers are required.

## Consequences
- Pros: Single IDL; less drift; typed clients; strong compatibility checks (buf lint/breaking).
- Cons: Requires code generation and consistent protobuf discipline.
- Mitigations: Retain the ability to generate JSON Schema from protobuf for consumers that require it.

## Migration
- ADR 0005 is superseded; the JSON Schema envelope is no longer authoritative.
- Existing JSON Schemas under /schemas remain for reference during transition but should not be treated as canonical or edited by hand.
- Update documentation and tickets to reflect protobuf-based envelope and payload definitions.

## References
- ADR 0008: docs/adr/0008-idl-and-rpc.md
- ADR 0009: docs/adr/0009-agent-orchestration-seams-and-langgraph-sidecar.md
- Protos: proto/libretto/events/v1/events.proto

