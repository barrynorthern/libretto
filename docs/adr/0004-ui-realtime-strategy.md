# ADR 0004: UI Real-time Strategy

Status: Accepted
Date: 2025-08-11

## Context
The UI must reflect the Living Narrative with low latency while avoiding divergent sources of truth. We plan incremental streaming for generation previews but want durable state to be canonical and consistent.

## Decision
- Firestore listeners are the single source of truth for persistent UI state (Canvas, Inspector, Governor).
- WebSockets are used only for transient streaming (e.g., token-by-token previews during generation). On completion, the canonical result is written to Firestore; the UI reconciles to Firestore.
- Presence and locks are modeled in Firestore (presence collection with TTL; advisory scene/sequence locks with owner and expiry).
- Canvas structural edits show previews locally; upon confirmation, emit NarrativeRestructured and write via GraphWrite, then update via Firestore listeners.

## Consequences
- Pros: Simpler consistency model, less duplication of real-time layers, straightforward offline behavior.
- Cons: Slightly higher latency for committing results compared to full WebSocket state mirroring.
- Mitigations: Use streaming only where it matters (generation previews) and keep deltas small to meet SLOs.

