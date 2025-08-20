# ADR 0009: Agent Orchestration Strategy — API Seams + LangGraph Sidecar for MVP

Status: Accepted
Date: 2025-08-16
Owners: barrynorthern
Reviewers: Augment Agent
Supersedes: ADR 0002 (Orchestration Choice for MVP)
Related: ADR 0001 (Data Store and Graph Abstraction), ADR 0005 (Events and Idempotency), ADR 0006 (Backend Language and Runtime), ADR 0008 (IDL and RPC)

## Context
- We are building an agentic AI product with a Go backend and Bazel builds, prioritizing a thin, reliable vertical slice.
- The agentic runtime landscape is evolving. Mature agent runtimes (e.g., LangGraph) are Python/TS-first; Go-native options are lighter-weight.
- We want to:
  - Establish durable API seams so orchestration is swappable.
  - Use a sensible off-the-shelf solution for MVP orchestration to reduce time-to-value.
  - Prefer Go for application code and plumbing; accept ring-fenced Python for specialized runtime pieces where justified.
  - Keep the door open to adopt Temporal for durable execution in later iterations.
- Previous decision (ADR 0002) selected Google Cloud Workflows for MVP. We are revising based on updated requirements and evaluation of agentic runtimes.

## Decision
- Define clear, versioned internal interfaces for agent orchestration and keep them language- and vendor-neutral.
- For the MVP, implement orchestration in a ring-fenced Python sidecar using LangGraph (production-ready agent runtime with stateful graph semantics and checkpointing).
- Keep Go as the system-of-record and external API surface; the Python sidecar is invoked via a strict RPC boundary (gRPC preferred; HTTP/JSON acceptable), with versioned schemas and strong isolation.
- Defer integrating Temporal until later; design the seams (IDs, idempotency, step boundaries, cancellation semantics) to enable a straightforward refactor to Temporal when warranted.

## Decision Drivers
- Speed to value: deliver a reliable, debuggable agent graph quickly.
- Reliability and observability: benefit from LangGraph's graph semantics and checkpointing now, with OTEL-friendly tracing from day one.
- Optionality: neutral seams let us swap/augment runtime (e.g., Temporal later) without rewriting product code.
- Team focus: keep Go app code and product features first; "framework-y" complexity is isolated.

## Architecture and API Seams

### High-level topology
- Go application services (APIs, tools, business logic)
- Agent Orchestration Service (Python/LangGraph sidecar)
- Shared data plane (Postgres, vector store such as pgvector/Qdrant)
- Observability: OpenTelemetry traces/metrics/logs correlated across services

### RPC boundary (stable, versioned)
- Transport: gRPC (per ADR 0008); fallback HTTP/JSON allowed when necessary
- Versioning: v1 namespace; additive changes only; explicit deprecation policy
- Contracts live under: idl/agent/v1 (protobuf and/or OpenAPI); codegen for Go and Python stubs

### Core messages (shape/semantics)
- StartAgentRequest
  - workflow_id (externally stable), session_id, actor/user_id
  - input (structured), tool_caps, budget_limits, metadata (tenant, PII flags)
- StartAgentResponse
  - run_id (opaque), initial_state_digest, accepted=true/false, errors[]
- Streamed events (server stream or webhook callback)
  - Event types: LlmStep, ToolCallRequested, ToolCallResult, HumanApprovalRequested, StateCheckpointed, RunCompleted, RunFailed, Cancelled
  - Every event carries: run_id, sequence_no, idempotency_key, trace_id/span_id, timestamp
- Control plane
  - CancelRun(run_id), GetRun(run_id), ResumeRun(run_id, payload), ListRuns(query)
- Error semantics
  - Clear retryability contract (retryable vs terminal); idempotency per step (see ADR 0005)

### Go-side interfaces (keep narrow)
- AgentClient: Start, Stream, Cancel, Resume, Get
- ToolServer: Implements Tools invoked by LangGraph via RPC callbacks (strict JSON schema, timeouts)
- CheckpointReader: Go reads snapshots for product UI/analytics (no direct writes into sidecar state)
- Tracing: OTEL everywhere; shared trace propagation (traceparent)

### Python sidecar contracts
- Deterministic graph nodes with explicit step boundaries and checkpointing
- Idempotent tool calls (idempotency_key provided; sidecar must pass it through)
- Configurable policies (timeouts, budgets, safety filters)
- Persistent store for checkpoints (Postgres) and run metadata

## Rationale and Trade-offs
- Not Go-only: Building graph semantics, checkpointing, retries, and human-in-the-loop to production reliability would delay the MVP.
- LangGraph choice: Currently the most production-ready agent runtime for stateful graphs and recovery.
- Ring-fenced Python: Keeps the rest of the stack in Go; orchestration is isolated, replaceable, and operationally constrained.
- Temporal later: Temporal adds infra and learning overhead; we retain a clean migration path via our seams.

## Temporal Future-proofing (explicit design hooks)
- Stable identifiers: workflow_id (external), run_id (internal), step_id (per node execution)
- Idempotency and replay: idempotency_key per step; steps designed to be side-effect free or compensable
- Cancellation and heartbeats: cooperative cancellation; periodic heartbeats from long steps for health checks
- State modeling: event-sourced or snapshot+event hybrid; compact "cursor" at each edge
- Activity boundaries: tool invocations and LLM calls treated as activities with clear retry/backoff policies

These choices map onto Temporal concepts (Workflow ID, Activities, Signals, Queries) if/when we migrate.

## Observability and Ops
- Tracing: OpenTelemetry across Go and Python with shared context; redact PII as needed
- Logging: Structured logs with correlation IDs; prompt/response redaction
- Metrics: Latency, cost, token usage, step retry counts, tool errors, completion rates
- Testing: Contract tests at the RPC boundary; golden traces for key flows; ephemeral env tests in CI

## Security and Compliance
- Boundary hardening: sidecar runs with least privilege; network ACLs; resource limits
- Data handling: PII flags in metadata; redaction at sources and in logs; encryption at rest/in transit
- Multi-tenancy: tenant ID propagated end-to-end; no cross-tenant state leakage (row-level security where feasible)

## Alternatives Considered
- Go-only orchestration (langchaingo + custom state)
  - Pros: monolingual; simpler ops
  - Cons: slower to reliable graph semantics; reinvents mature features
- Semantic Kernel / AutoGen / CrewAI
  - Pros: quick starts
  - Cons: weaker production graph semantics vs LangGraph; Python/TS; less aligned with our needs
- Temporal from day one
  - Pros: rock-solid durability
  - Cons: infra and learning overhead now; can defer while keeping migration easy
- Google Cloud Workflows (ADR 0002)
  - Pros: minimal ops; native GCP integration
  - Cons: insufficient for stateful agent graphs and rich checkpoint semantics; superseded by this ADR

## Consequences
- Positive: faster path to robust agentic behavior and checkpointing; clear seams reduce lock-in; migration path to Temporal preserved
- Negative: polyglot ops (Go + Python) and CI/build complexity; need to maintain RPC contracts and versioning discipline

## Implementation Outline (MVP)
- Define v1 protobuf/JSON schemas for AgentClient and event stream (per ADR 0008)
- Implement Go AgentClient and ToolServer (proto-defined contracts; JSON Schema generated only if required by downstream tooling)
- Implement Python LangGraph sidecar with:
  - Planner node, tool-executing node(s), checkpointing in Postgres
  - RPC adaptors to call Go tools and emit events
  - OTEL instrumentation and structured logging
- Wire tracing/metrics/logging end-to-end
- Add contract tests and a golden-path E2E in CI

## Acceptance Criteria
- Start/stream/cancel/resume flows work end-to-end with correlated traces
- Checkpoints persist; runs are resumable after process restarts
- Tool calls are idempotent and time-bounded
- Dashboards show latency, cost, and error breakdown per step
- No cross-tenant leakage in queries or logs

## Rollback Plan
- If the sidecar proves too costly/complex, replace with a minimal Go orchestrator behind the same AgentClient interface. RPC contracts remain stable.

## Open Questions
- Which vector store for MVP (pgvector vs Qdrant) and which Go client libraries?
- Hosting model for the sidecar (Bazel-built container vs rules_python vs managed runtime)?
- Do we need LangSmith or stick to OTEL + our own dashboards?

