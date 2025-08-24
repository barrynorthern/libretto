# Libretto

A multi-agent narrative orchestration engine where the user (the Conductor) directs specialized AI agents to build a Living Narrative graph. This is not a traditional writing app; prose is generated from the model and refined via structured controls.

## Strategy & Stack
- Product strategy: plans/libretto-product-strategy.md
- Tech stack guidelines: plans/tech-stack-guidelines.md

## Core ideas
- The Conductor, not the typist: high-level directives via a Baton (NL command palette)
- The Living Narrative: a versioned, validated graph of scenes, arcs, characters, settings, and relationships
- Multi-agent AI: event-driven agents (Plot Weaver, Thematic Steward, etc.) collaborating over a Narrative Event Bus
- Generated prose: prose is a read-only view generated from the model; refinements happen via Tuners or graph edits

## MVP at a glance
- Desktop-first; no cloud infra required to use
- SQLite persistence; sqlc repositories
- Context Manager + RAG for narrative-aware prompts
- Simple Baton flow: Directive → Proposal → Persisted Scene → UI list/detail

## Architecture (current)
- Single process with internal modules: Orchestrator, PlotWeaver, Narrative, GraphWrite store, Context Manager
- Wails React UI calls Go bindings; protobuf DTOs at the boundary; TS client generated via buf

## Repository layout
```
/plans/libretto-product-strategy.md   # Product strategy and staircase
/plans/tech-stack-guidelines.md       # Tech stack and workflow
/docs/ddd/overview.md                 # DDD overview (mermaid diagrams)
/cmd, /internal                       # Backend monolith (Go)
/proto, /gen                          # Shared protobufs and generated code (Go/TS)
```

## Development workflow (short)
- Keep domain logic in Go modules; expose Wails bindings
- Author SQL and run sqlc for type-safe repositories
- Maintain proto DTOs for UI boundary; generate TS clients via buf

## Getting started
- Build: `make build`  •  Test: `make test`
- Run monolith: `make dev-up` (serves API + bindings)
- UI (Wails): see ticket 00009 for scaffold plan

## Contributing
- Discuss changes against strategy and tech stack docs first
- Keep proto contracts and sqlc schema changes small and reviewed

## Security & privacy
- Do not commit secrets. Use macOS Keychain or environment variables for local credentials.
- Local-first; no external services required. Future cloud features will be opt-in.

## License
- TBA

## Getting started (early repo)
- Read spec.md (v1.1) end-to-end
- Review schemas under /schemas
- Review the GraphWrite API and Bootstrap Helper drafts under /docs
- See plans/execution-plan.md for phased delivery and DoD
