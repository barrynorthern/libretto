# Libretto

A multi-agent narrative orchestration engine where the user (the Conductor) directs specialized AI agents to build a Living Narrative graph. This is not a traditional writing app; prose is generated from the model and refined via structured controls.

## Status
- Specification: v1.1 (spec.md) — authoritative source of truth using EARS-style functional requirements
- Planning: plans/execution-plan.md — phased plan, WBS, CI/CD outline, Firestore layout, SLOs
- Schemas: JSON Schemas for events and entities under /schemas

## Core ideas
- The Conductor, not the typist: high-level directives via a Baton (NL command palette)
- The Living Narrative: a versioned, validated graph of scenes, arcs, characters, settings, and relationships
- Multi-agent AI: event-driven agents (Plot Weaver, Thematic Steward, etc.) collaborating over a Narrative Event Bus
- Generated prose: prose is a read-only view generated from the model; refinements happen via Tuners or graph edits

## MVP scope (high level)
- Firestore as the canonical store with denormalized adjacency indexes (Graph Service owns all mutations)
- Google Cloud Workflows for orchestration (Temporal-ready design)
- Firebase Auth (RBAC: Owner/Editor/Viewer), WorkOS later if needed
- Real-time: Firestore listeners as the source of truth; WebSockets for transient streaming only
- Bootstrap: template-first wizard with lightweight paste/upload assist (Markdown/CSV/TXT)
- Export: final copy (Markdown) and compendia (CSV/JSON/Markdown)

## Architecture at a glance
- UI: Next.js (planned) — Canvas, Inspector, Governor
- Event bus: Cloud Pub/Sub with JSON-schema’d events (envelope + versioned payloads)
- Agents: Cloud Functions/Run (starting simple) emitting/consuming events; strictly idempotent
- Durable flows: Google Cloud Workflows for multi-step chains with checkpoints
- RAG (later phases): Vertex AI Vector Search projection of the Living Narrative

## Repository layout
```
/spec.md                    # Specification v1.1 (canonical)
/plans/execution-plan.md    # Phases, WBS, CI/CD, SLOs
/schemas/events/*.json      # Event envelope + MVP events
/schemas/model/*.json       # Minimal entity schemas (MVP)
/docs/graphwrite-api.md     # GraphWrite API contract (draft)
/docs/bootstrap-helper.md   # Paste/upload helper contract (draft)
```

## Spec-Driven Development
- The spec is the single source of truth (see Section 5.0+6.0). All features tie to EARS FRs.
- Changes flow: propose deltas in spec.md → update schemas/contracts → implement → tests/observability → CI.

## Bootstrap (MVP)
- Template-first wizard (Three-Act, Hero’s Journey to start) with archetypes
- Optional paste/upload (Markdown/CSV/TXT) to fill in template fields
- Human-in-the-loop review before applying to a bootstrap branch (versioned)

## Exports (MVP)
- Final copy: Markdown
- Compendia: CSV/JSON/Markdown (story bible, character book, lore book)

## Contributing
- Open issues/PRs against spec.md before implementation changes
- Follow the PR checklist in plans/execution-plan.md (schemas, tests, observability, IaC updates)

## Security & privacy
- Do not commit secrets. Use macOS Keychain or environment variables for local credentials.
- All runtime secrets will be managed via GCP KMS and Terraform.

## License
- TBA

## Getting started (early repo)
- Read spec.md (v1.1) end-to-end
- Review schemas under /schemas
- Review the GraphWrite API and Bootstrap Helper drafts under /docs
- See plans/execution-plan.md for phased delivery and DoD

