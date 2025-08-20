# ADR 0005: Event Schema and Idempotency Conventions

Status: Superseded by ADR 0010
Date: 2025-08-11

## Context
Agents communicate via Pub/Sub with at-least-once delivery and potential reordering. We need consistent event envelopes, correlation, and idempotency to avoid duplicate side effects.

## Decision (Superseded)
- Use a versioned JSON Schema envelope: eventName, eventVersion, eventId, occurredAt, correlationId, causationId, idempotencyKey, producer, tenantId, payload.
- All side-effecting operations must be idempotent keyed by idempotencyKey; consumers return prior results on duplicates.
- Use correlationId/causationId to link traces across UI→API→Event→Agent→GraphWrite.
- Per-topic retry and DLQ policies are defined; AgentError is emitted with diagnostics on failure.

## Consequences
- Pros: Strong debugging/traceability, safer retries, simpler reasoning about duplicates.
- Cons: Additional bookkeeping for keys and result caching.
- Mitigations: Provide shared libraries/utilities for key generation and dedup checks.

