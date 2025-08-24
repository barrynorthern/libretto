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
    App["internal/app Orchestrator"]
    Plot["internal/agents/plotweaver"]
    Narr["internal/agents/narrative"]
    Store["internal/graphwrite (Store + Repos)"]
    Ctx["internal/context (Context Manager: memory, RAG)"]
  end
  UI["Desktop UI (Wails React)"] <--> |Bindings + DTOs| App
  App --> Plot
  App --> Narr
  App --> Ctx
  Narr --> Store
```

## Domain Model and Seams (high level)
```mermaid
flowchart LR
  subgraph "Application Layer"
    Orchestrator["Orchestrator (Baton service)"]
    GraphWriteSvc["GraphWrite App Service"]
  end
  subgraph "Domain Layer"
    ProjectAgg["Aggregate: Project"]
    SceneEnt["Entity: Scene"]
    Versioning["Value Obj: GraphVersion"]
  end
  subgraph "Infrastructure Layer"
    Repos["Repositories (sqlc)"]
    DB["SQLite (canonical store)"]
    VectorDB["sqlite-vec (local vector DB)"]
    Models["Model Providers (Ollama default; API keys optional)"]
  end
  UI["Wails React UI (DTOs via proto)"] -->|DTOs| Orchestrator
  Orchestrator -->|domain calls| GraphWriteSvc
  GraphWriteSvc -->|repos| Repos
  Repos --> DB
  Orchestrator -->|ContextBundle| VectorDB
  Orchestrator -->|ModelSpec| Models
  ProjectAgg --> SceneEnt
  Versioning --> GraphWriteSvc
```

### Key seams
- UI ↔ Application: protobuf DTOs for requests/responses (TS/Go generated)
- Application ↔ Domain: method calls with domain types (no proto inside core)
- Domain ↔ Persistence: repository interfaces (sqlc-backed) against SQLite
- Application ↔ Context: Context Manager provides ContextBundle (RAG via sqlite-vec)
- Application ↔ Models: ModelSelector chooses Ollama or API provider based on task/constraints

## Directive → Persisted Scene (Sequence)

```mermaid
sequenceDiagram
  participant UI
  participant App as Orchestrator
  participant Plot as PlotWeaver
  participant Ctx as ContextMgr
  participant Narr as Narrative
  participant Repo as "Store (sqlc)"
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

### Seams-first sequence (Context + Model selection in loop)

```mermaid
sequenceDiagram
  participant UI
  participant App as Orchestrator
  participant Ctx as ContextMgr
  participant Sel as ModelSelector
  participant Plot as PlotWeaver
  participant GW as GraphWriteSvc
  participant Repo as "Store (sqlc)"
  UI->>App: IssueDirective(DirectiveDTO)
  App->>Ctx: BuildContext(project, directive)
  Ctx-->>App: ContextBundle
  App->>Sel: ChooseModel(task, complexity, budget)
  Sel-->>App: ModelSpec (Ollama default)
  App->>Plot: ProcessDirective(directive, ContextBundle, ModelSpec)
  Plot-->>App: SceneProposal
  App->>GW: ApplySceneProposal(proposal)
  GW->>Repo: Persist(scene)
  Repo-->>GW: ok
  GW-->>App: ok
  App-->>UI: ProposalApplied(sceneId)
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

