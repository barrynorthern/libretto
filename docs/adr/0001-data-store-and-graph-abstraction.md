# ADR 0001: Canonical Store and Graph Abstraction

Status: Accepted
Date: 2025-08-11

## Context
We need a reliable, low-ops canonical data store with real-time sync for the Living Narrative graph. Pure graph DBs add ops/learning overhead; GCP lacks a first-party graph DB.

## Decision
- Use Firestore as canonical store.
- Provide a Graph Service (GraphWrite API + derived adjacency indexes) to enforce schema, referential integrity, and versioned deltas.
- Maintain denormalized adjacency collections for O(1) neighbor fetches. All mutations go through GraphWrite, not direct writes.

## Consequences
- Pros: Real-time listeners, low ops, simple scaling, good developer velocity.
- Cons: No native graph traversals; must maintain derived indexes and validations.
- Mitigations: Encapsulate complexity in Graph Service; consider optional read-model projection to Neo4j Aura if analytic queries become complex.

