# Domain-Driven Design Overview (MVP)

Status: Proposed

## Ubiquitous Language
- Project: container for scenes and related narrative assets
- Scene: atomic narrative unit with title, summary, content
- Directive: user instruction to agents (text, act, target)
- Scene Proposal: candidate scene derived from a directive

## Bounded Contexts
- Creation: directive -> proposal
- Narrative Graph: storage, indexing, retrieval of scenes (later: acts/beats)
- Rendering/Export: formats (later)

## Module/Container View (Monolith)
```mermaid
flowchart TB
  subgraph Monolith
    App[internal/app Orchestrator]
    Plot[internal/agents/plotweaver]
    Narr[internal/agents/narrative]
    Store[internal/graphwrite (Store + Repos)]
    Ctx[internal/context (Context Manager: memory, RAG)]
  end
  UI[Desktop UI (Wails React)] <---> |Bindings + DTOs| App
  App --> Plot
  App --> Narr
  App --> Ctx
  Narr --> Store
```

## Domain Model (initial)
```mermaid
classDiagram
  class Project {
    UUID id
    string name
    +AddScene(Scene)
    +ListScenes()
  }
  class Scene {
    UUID id
    string title
    string summary
    string content
  }
  Project "1" o-- "*" Scene
```

## Directive → Persisted Scene (Sequence)
```mermaid
sequenceDiagram
  participant UI
  participant App as Orchestrator
  participant Plot as PlotWeaver
  participant Ctx as ContextMgr
  participant Narr as Narrative
  participant Repo as Store(sqlc)
  UI->>App: IssueDirective(DirectiveDTO)
  App->>Ctx: BuildContext(project, directive)
  Ctx-->>App: ContextBundle
  App->>Plot: ProcessDirective(text, act, target, ContextBundle)
  Plot-->>App: SceneProposal(domain)
  App->>Narr: ApplySceneProposal(Store, proposal)
  Narr->>Repo: CreateScene(id, title, summary, content)
  Repo-->>Narr: ok
  Narr-->>App: ok
  App-->>UI: IssueDirectiveResponse(correlationId)
```

## Proto Seams (DTOs)
- baton.v1: IssueDirectiveRequest/Response
- scene.v1: Scene, SceneList
- context.v1 (later): ContextBundle summary for debugging/telemetry

## Notes on Context Management
- Local Context Manager responsible for:
  - Token budget planning per task
  - Prompt assembly with domain-aware sections (beats, characters, prior scenes)
  - RAG over local vector DB (scoped to project)
  - Model selection policy (Ollama vs API providers)
- Local vector DB options: SQLite+Vec (pgvecto.rs alt), Qdrant (embedded), or Milvus lite; start with sqlite-vec for fewer moving parts.

