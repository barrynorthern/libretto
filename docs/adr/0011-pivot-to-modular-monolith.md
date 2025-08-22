# ADR 0011: Pivot to a Modular Monolith for the MVP

Status: Proposed
Date: 2025-08-20

## Context
The initial design pursued a microservice-per-agent approach (API, PlotWeaver, NarrativeIngest, GraphWrite) with event envelopes and simulated Pub/Sub. This introduced operational complexity and slowed delivery of the MVP vertical slice. For a single-user application, independent scaling and heterogenous tech stacks are not required.

## Decision
Adopt a modular monolith for the MVP:
- Single binary (cmd/libretto) hosting API handlers and internal modules.
- Agent modules implemented as Go packages with clear interfaces; invoked synchronously by an orchestrator.
- Preserve extraction seams via interfaces and typed contracts; protobuf stays for long-lived contracts but internal calls prefer native Go structs.
- Defer distributed messaging and service decomposition until justified by real needs.

## Consequences
- Simpler dev loop: one process; less flakiness; faster iteration and testing.
- Reduced surface area: fewer network/serialization failures and schema drifts.
- Clear upgrade path: any module can be extracted later behind the same interface.

## Implementation Plan (high level)
1. Collapse services into internal packages and a single main:
   - internal/app/orchestrator (drives flow)
   - internal/agents/plotweaver
   - internal/agents/narrative
   - internal/graphwrite (Apply + in-memory store)
2. Replace DevPush/push with direct calls in the happy path.
3. Keep protobuf messages for externalization; use Go types internally where simpler.
4. UI: separate Next.js app (CSR only) in a sibling project/repo; backend exposes HTTP APIs for it to consume.

## Alternatives Considered
- In-process pub/sub: preserves event semantics but still adds abstraction cost; can be added later if needed.
- Staying distributed: rejected for MVP due to complexity vs benefit mismatch.

## References
- Prior ADRs 0009/0010 on seams and contracts; this ADR supersedes distribution for MVP.

