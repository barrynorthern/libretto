# Libretto Execution Plan (MVP to v1)

## Objectives
- Deliver an end-to-end slice: Baton -> Event Bus -> Agents -> GraphWrite -> Firestore -> Canvas/Inspector.
- Establish platform primitives: canonical data model, Event Schema Registry, Agent Contract, observability, safety rails, and export.

## Architectural decisions
- Canonical store: Firestore + Graph Service with denormalized adjacency indexes (Path A).
- Orchestration: Google Cloud Workflows for MVP; design interfaces Temporal-ready.
- Auth: Firebase Auth MVP; abstraction layer to allow WorkOS later.
- Real-time: Firestore listeners as source of truth; WebSockets only for streaming previews.
- IaC: Terraform for all infra.

## Milestones

### Phase 0 – Decisions & scaffolding (1–2 weeks)
- Author ADRs: data model, event schema registry, agent contract, orchestration choice.
- Repo structure: /spec, /plans, /schemas/events, /schemas/model, /services/{api,agents}, /infra/terraform.
- CI/CD: build/test/lint; preview envs; secrets via KMS.
- Observability baseline: correlationId/causationId propagation; Cloud Logging/Monitoring/Trace; DLQ alert policy.

### Phase 1 – MVP vertical slice (2–3 weeks)
- GraphWrite API: schema validation, referential integrity, versioned deltas, Firestore persistence with adjacency indexes.
- Baton -> DirectiveIssued publisher with idempotencyKey.
- Agents: Plot Weaver (minimal), Thematic Steward (validator). Consume DirectiveIssued; produce SceneProposalReady.
- Canvas/Inspector: render scenes from Firestore; show analyses stub; allow highlight-based annotations.
- Export: final copy (Markdown) from current GraphVersion.
- Tests: E2E for directive flow; contract tests for events; API unit tests; UI integration test for Canvas load and selection.

### Phase 1.5 – Bootstrap (simple) (1 week)
- Bootstrap Wizard (template-first): choose template/archetypes; parameterize; seed graph skeleton (bootstrap branch).
- Paste/Import assist: paste Markdown/CSV/TXT into wizard; synchronous lightweight parsing in API/Bootstrap Helper; review and apply via GraphWrite.
- Export compendia (MVP): CSV/JSON/Markdown for character book and story bible.

### Phase 2 – Collaboration & refinement (2–3 weeks)
- Character Troupe agent for protagonist; CharacterReactionRequested/Provided.
- Tuners & RefinementDirective loop; Governor view with in-flight event tracking.
- Undo/redo and snapshots; optimistic concurrency with merge UI.
- Cost budgets and throttling for LLM usage; user notifications.


## Work breakdown structure (WBS)

1. Foundations
   - Terraform modules: pubsub_topic, pubsub_subscription, cloud_run, cloud_workflows, firestore_rules, kms, monitoring_alerts
   - CI/CD: lint/test/build; terraform plan/apply (manual approve); schema validation step
   - Repo structure: /spec, /plans, /schemas/{events,model}, /docs, /services/{api,agents}, /infra/terraform
2. Data & contracts
   - Finalize model schemas; event schemas; Agent Contract doc; GraphWrite API
   - Firestore layout: collections nodes_{type}, edges, adjacency_out, adjacency_in, versions, deltas
   - Security rules: read/write by role; GraphWrite-only writes; rule tests
3. Vertical slice
   - Baton API + DirectiveIssued; Plot Weaver + Thematic Steward (stubs); Canvas/Inspector basic
   - Observability: tracing with correlationId; dashboards for latency, DLQ, queue depth
4. Bootstrap (Phase 1.5)
   - Template library v1 (Three-Act, Hero's Journey) + archetypes
   - Bootstrap Wizard + paste/upload parsing (Markdown/CSV/TXT)
   - Export: Markdown (final), CSV/JSON/Markdown (compendia)
5. Collaboration & refinement
   - Governor view; Tuners loop; Character Troupe (protagonist); undo/redo & merge UI
6. Reliability & safety
   - DLQs/retries/idempotency; safetyReview flow; SLO validation & load tests

## Firestore structure (draft)
- projects/{projectId}
- graphs/{projectId}/versions/{versionId}
- graphs/{projectId}/deltas/{deltaId}
- nodes_{type}/{nodeId} (scene, character, arc, setting, item, theme)
- edges/{edgeId} { from, to, type, versionId }
- adjacency_out/{nodeId} { edges: [...] }
- adjacency_in/{nodeId} { edges: [...] }
- annotations/{annotationId}
- analyses/{analysisId}
- exports/{exportId}

Notes: Writes go through GraphWrite; preconditions enforce version; triggers maintain adjacency.

## Pub/Sub topics (naming)
- libretto.dev.directive.issued.v1
- libretto.dev.narrative.restructured.v1
- libretto.dev.scene.proposal.ready.v1
- libretto.dev.template.applied.v1
- libretto.dev.import.review.submitted.v1
- libretto.dev.import.applied.v1
- libretto.dev.bootstrap.graph.ready.v1
- libretto.dev.agent.error.v1

## CI/CD pipeline (draft)
- PR: lint + unit tests + schema validation (ajv) + firestore rules tests
- Main: build, deploy API/agents (staging), terraform plan; manual approve to apply; smoke tests; promote
- Release tags: trigger prod deploy; migrations (if any) behind feature flags

## Dashboards & alerts
- Dashboards: event latency by agent; queue depth; DLQ rate; UI update latency; token usage per project
- Alerts: DLQ rate > 0.5%/day; p95 latency breaches; budget threshold exceeded; push failures

## Definition of Done (per feature)
- Spec references updated; schemas validated; tests (unit+integration+E2E) passing; dashboards updated; IaC applied; docs updated; PR checklist complete

## PR checklist
- [ ] Linked FRs/NFRs
- [ ] JSON Schemas added/updated
- [ ] Tests added/updated
- [ ] Observability added (logs, metrics, traces)
- [ ] Terraform updated
- [ ] Security rules updated
- [ ] Rollback plan noted

## Timeline (indicative)
- Phase 0: 1–2 weeks; Phase 1: 2–3 weeks; Phase 1.5: 1 week; Phase 2: 2–3 weeks; Phase 3: 2 weeks

### Phase 3 – Reliability & safety (2 weeks)
- DLQ, retries, idempotency enforcement; AgentError surfacing in Governor.
- SafetyReview annotations and Owner override path.
- SLO dashboards; load tests and failure injection.

### Phase 4 – Expansion (ongoing)
- Continuity Steward with RAG (Vertex Vector Search) + CDC embedding pipeline.
- Optional: Neo4j Aura read-model projection if query complexity rises.
- Presence indicators; scene/sequence locks with TTL.
- Accessibility sweep (WCAG 2.1 AA) and i18n readiness.

## Deliverables per phase
- Updated spec (v1.1 already integrated), schemas, ADRs, CI/CD pipelines, Terraform modules, dashboards.

## Success metrics
- p95 DirectiveIssued→SceneProposalReady < 20s; UI update post-write p95 < 500ms; DLQ < 0.5%/day; zero data loss under failure injection; undo/redo correct over 20 ops.

## Risks & mitigations
- Firestore graph complexity: hidden behind Graph Service; option to add Aura read-model.
- Cost spikes: budgets, throttles, alerts.
- Vendor lock-in: provider-agnostic agent contract; abstracted Auth.

