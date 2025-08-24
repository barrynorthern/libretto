# Libretto: A Multi-Agent Narrative Orchestration Engine
## Specification Document v1.1

### 1.0 Core Philosophy: The Conductor and the Living Narrative

This system is founded on a radical departure from traditional writing software. The user does not "write" in the conventional sense. Instead, they **conduct** a team of specialized AI agents who collaboratively build, refine, and render a story based on high-level creative directives.

*   **The Conductor, Not the Typist:** The user's role is strategic and creative, akin to a film director or orchestra conductor. They guide the narrative's shape, tone, and direction through a dynamic, visual interface, leaving the granular generation of prose and detail to their AI orchestra.
*   **The Living Narrative, Not a Static Document:** The story is not a linear sequence of text. It is a complex, interconnected data model—a "narrative graph"—where characters, plot points, locations, dialogue, and thematic elements are distinct but related objects. This allows agents to reason about the story holistically and enables the narrative to be explored and rendered in non-linear ways.
*   **Emergent Storytelling:** The process is one of discovery. The Conductor provides the creative impetus, and the AI orchestra explores the narrative space, generating possibilities, highlighting consequences, and weaving together threads that the Conductor can then shape, approve, or reject.

### 2.0 The Conductor's Console: The New Experience

The primary user interface is not an editor; it is a dynamic, multi-faceted **Orchestration Canvas**. This is the Conductor's podium, providing a god's-eye view of the Living Narrative and the tools to direct the AI Orchestra.

#### 2.1 The Narrative Canvas & Inspector
This is the central view, replacing the "editor window." It is a non-linear, interactive visualization of the story's structure, always accompanied by a context-aware Inspector panel. [1, 2]

| ID | Type | Requirement |
| :--- | :--- | :--- |
| **FR-2.1.1** | Ubiquitous | The system **shall** represent the Living Narrative as an interactive graph on the Narrative Canvas, where nodes represent story elements (scenes, character introductions, plot twists) and edges represent their relationships. |
| **FR-2.1.2** | Ubiquitous | The system **shall** permanently display an "Inspector" panel, which **shall** provide contextual information and controls for any element currently selected on the Narrative Canvas or within the agent management interface. |
| **FR-2.1.3** | Event-Driven | When the Conductor selects a node on the canvas, the Inspector **shall** display the generated prose, associated characters, location data, and agent-generated metadata (e.g., emotional analysis, thematic relevance score). |
| **FR-2.1.4** | State-Driven | While the Conductor is interacting with the canvas, the system **shall** allow them to drag, drop, connect, and re-sequence nodes to visually restructure the narrative flow. |
| **FR-2.1.5** | Event-Driven | When the Conductor issues a structural change on the canvas (e.g., moving a scene from Act 1 to Act 3), the system **shall** automatically post a "NarrativeRestructured" event to the Narrative Event Bus, triggering the Continuity Steward agent. |

#### 2.2 The Baton: The Command & Control Interface
The Conductor interacts with the AI Orchestra through a high-level, intention-driven command interface.

| ID | Type | Requirement |
| :--- | :--- | :--- |
| **FR-2.2.1** | Ubiquitous | The system **shall** provide a universal, natural language command palette ("The Baton") as the primary method for issuing directives to the AI Orchestra. |
| **FR-2.2.2** | Event-Driven | When the Conductor issues a high-level directive (e.g., "Create a scene where the protagonist confronts their mentor"), the system **shall** post a "DirectiveIssued" event to the Narrative Event Bus. |
| **FR-2.2.3** | State-Driven | While an agent's proposal is under review in the Inspector, the system **shall** provide the Conductor with "Tuner" controls to refine the output (e.g., sliders for "Tone: More Ominous," "Pacing: Faster," or "Dialogue: More Formal"). [3] |
| **FR-2.2.4** | State-Driven | While agents are active, the system **shall** provide a "Governor" view, allowing the Conductor to see the active events on the Narrative Event Bus and monitor the progress of each specialist agent. [3] |
| **FR-2.2.5** | Ubiquitous | The system **shall** offer contextual "Nudges" and "Suggestions" to the Conductor, with agents proactively identifying opportunities for plot development, character arcs, or thematic resonance. [3] |
| **FR-2.2.6** | Event-Driven | When the Conductor selects an agent in the Governor view, the Inspector **shall** display the agent's details, including its defined role and the specific event "hooks" it is subscribed to on the Narrative Event Bus. |

### 3.0 The AI Orchestra: The Multi-Agent System

The creative engine of the application is a collaborative, multi-agent system. The agents operate in a decentralized, event-driven manner, subscribing to a central "Narrative Event Bus" to find tasks that match their expertise. [4, 5]

#### 3.1 Agent Roles & Responsibilities

| Agent Name | Role | Core Function |
| :--- | :--- | :--- |
| **The Plot Weaver** | Specialist | Generates and refines narrative structure, causality, and plot points. Subscribes to "DirectiveIssued" events related to plot. Can be tasked with creating inciting incidents, rising action, climaxes, and resolutions. |
| **The World Architect** | Specialist (RAG) | Responsible for creating and maintaining the story's setting and lore. Subscribes to events requiring world-building details. Uses the Living Narrative as its knowledge base to ensure all world-building is internally consistent. |
| **The Character Troupe** | Specialist | A collective of autonomous sub-agents, one for each character in the story. Each agent is a "method actor," an expert in their assigned character's voice, motivations, and knowledge. They subscribe to events mentioning their character and will "fight" for their voice to be represented authentically, flagging dialogue or actions that feel out of character. |
| **The Continuity Steward** | Specialist (RAG) | The guardian of narrative consistency. Subscribes to all events that modify the Living Narrative (e.g., "SceneAdded," "CharacterUpdated," "NarrativeRestructured"). It continuously monitors for contradictions in timelines, character knowledge, or physical descriptions. |
| **The Dramaturg** | Specialist | The ruthless editor. Its purpose is to "kill your darlings." It subscribes to "ActCompleted" or "DraftFinished" events and analyzes the narrative for superfluous scenes, characters, or dialogue that do not serve the core theme or advance the plot, proposing cuts to improve focus and impact. |
| **The Empath** | Specialist | The steward of the story's emotional core. Subscribes to "SceneGenerated" events. It analyzes scenes for emotional impact using sentiment analysis and narrative theory, tracking the emotional arc of characters and the story as a whole. It provides feedback on how to heighten specific feelings (tension, joy, sorrow) to achieve the Conductor's desired emotional response from the audience. |
| **The Thematic Steward** | Specialist | The guardian of the story's central idea. It is initialized with the story's core thematic question (e.g., "Can love conquer greed?"). It subscribes to all major narrative events and evaluates new content for its relevance and contribution to this theme, flagging anything that diverges or weakens the story's scope. |

#### 3.2 Agent Collaboration & Workflow
The system rejects a rigid, top-down hierarchy in favor of a dynamic, event-driven architecture. This promotes emergent collaboration and reduces bottlenecks. [4]

| ID | Type | Requirement |
| :--- | :--- | :--- |
| **FR-3.2.1** | Ubiquitous | The system **shall** use a publish-subscribe message bus ("The Narrative Event Bus") as the primary mechanism for agent orchestration. |
| **FR-3.2.2** | Event-Driven | When a Conductor issues a directive, a "publisher" service **shall** post a structured event to the Narrative Event Bus. |
| **FR-3.2.3** | State-Driven | While active, each specialist agent **shall** subscribe to specific event types ("hooks") on the Narrative Event Bus that match its expertise. |
| **FR-3.2.4** | Event-Driven | When an agent consumes an event, it **shall** execute its core logic and, upon completion, **shall** publish its output as a new event (e.g., "SceneProposalReady") back to the Narrative Event Bus. |
| **FR-3.2.5** | Complex | Where direct, targeted communication is required, when an agent needs a specific piece of information from another, it **shall** be able to use the A2A protocol for a point-to-point request, complementing the primary event-driven workflow. [6, 7] |

### 4.0 Architecture & Implementation (Local‑First Desktop)

This experience is built on a modern, real-time, serverless architecture designed for massive scalability and creative flexibility. The architecture adopts a hybrid database approach to use the best tool for each job.

#### 4.1 Technology Stack & Services

| Component | Technology/Service | Rationale |
| :--- | :--- | :--- |
| **Desktop UI** | Wails (React + TypeScript) | Local-first desktop app with a fast feedback loop and direct bindings into Go. |
| **State Persistence** | SQLite + sqlc (Go) | Zero-config local DB; type-safe queries via sqlc; portable single-file storage. |
| **Agent Orchestration** | In-process orchestrator (publish-subscribe semantics) | Decoupled modules with event semantics inside one process; extractable later if needed. |
| **Context & RAG** | sqlite-vec (local vector DB) | Project-scoped embeddings and similarity search without new infra. |
| **Embeddings/Models** | Ollama (default) + optional API-key providers | Offline/private by default; opt-in to hosted models for convenience. |
| **Build** | Monorepo with Bazel (Go) + pnpm (UI) + buf (codegen) | Reproducible builds and shared contracts; TS/Go clients from protos. |

#### 4.2 Architectural Flow: A Conductor's Directive (Local‑First)

1.  **Directive Issued:** The Conductor issues a command via the Baton: *"Introduce a betrayal that complicates the protagonist's main goal."*
2.  **Context Built:** The Orchestrator asks the Context Manager to assemble a ContextBundle (relevant scenes, characters, beats) via sqlite‑vec and domain rules.
3.  **Plot Proposal:** The Orchestrator calls PlotWeaver with the directive + ContextBundle. PlotWeaver uses the selected model (Ollama by default) to produce a SceneProposal.
4.  **Apply:** The Orchestrator invokes the Narrative module to persist the scene via repositories (SQLite/sqlc).
5.  **UI Update:** The UI reads from the local store and displays the new scene for review in the Inspector.
6.  **Refine:** Tuners issue a refinement directive; the same path repeats with an updated ContextBundle.

#### 4.3 Durable Execution and State Management
The loss of creative work or state is a catastrophic failure. The system is architected to prevent this at all costs.

| ID | Type | Requirement |
| :--- | :--- | :--- |
| **FR-4.3.1** | Ubiquitous | The system **shall** ensure durable execution for all long-running, multi-step agent workflows. [8, 9] |
| **FR-4.3.2** | Event-Driven | When a complex directive is issued that requires a sequence of steps, the system **shall** orchestrate the chain with a durable mechanism (in‑process for MVP; pluggable durable engine later), ensuring that the overall process state is persisted at every step. |
| **FR-4.3.3** | Ubiquitous | The system **shall** use Firestore's transactional capabilities to ensure that all changes to the Living Narrative graph and the Conductor's view state are saved atomically and durably. [10] |
| **FR-4.3.4** | Ubiquitous | The system **shall** implement a robust checkpointing mechanism for agent workflows, allowing a process to be resumed from the last successful step in the event of a transient failure in a downstream service or agent. [9, 11] |


#### 4.4 Build & Repo Strategy (Monorepo)

#### 4.4.1 Local build and tooling
- Single Go binary (monolith) built with Bazel (rules_go, gazelle).
- Desktop UI via Wails (React + TypeScript) using pnpm; Bazel wrappers optional later.
- Codegen: buf for protobufs; generate Go and TS clients used by app and UI.
- Monorepo: Shared contracts and schemas versioned in-repo for Go and TS.

### 5.0 v1.1 Addendum: Canonical Specifications for Implementation

This section is authoritative for implementation and augments Sections 1–4 without altering the product vision. Where Functional Requirements (FRs) are restated below in EARS format, these supersede earlier phrasing.

#### 5.1 Glossary
- Living Narrative: A versioned, directed narrative graph (nodes: scenes, beats, arcs, characters, locations, items, themes; edges: typed relationships) plus analyses and annotations.
- Narrative Event Bus: Pub/Sub topics carrying versioned events that conform to the Event Schema Registry.
- Agent: A compute unit that implements the Agent Contract to consume events and emit outputs.
- Baton: Natural-language command palette that emits structured directives.
- Inspector/Governor/Canvas: UI surfaces for context, orchestration visibility, and graph interaction.

#### 5.2 Scope and Assumptions
- In scope (MVP): Single-project orchestration; offline-first desktop; versioning/snapshots; undo/redo; basic branching; cost budgets; export.
- Out of scope (MVP): Rich free-form text editing; third-party agent marketplace; multi-user collaboration.
- Assumptions: Local-first monolith; SQLite as canonical store; sqlite‑vec for RAG; Ollama for embeddings/models by default; optional API providers via user keys.

#### 5.3 Prose Policy and Export/Import
- Prose is a generated view of the Living Narrative, not a primary editing surface. Refinement occurs via structured Tuners or graph edits that trigger regeneration.
- The Inspector may support prose highlighting/selection to attach precise feedback (annotations) that map back to narrative elements; this is akin to PR review comments and remains structured.
- Export: The system shall export (a) final copy (e.g., Markdown, DOCX, PDF), and (b) reference compendia (e.g., story bible, character book, lore book) derived from the model; exports are frozen, read-only artifacts.
- Import (later): If supported, imports must be transformed into structured deltas that map to the narrative model; arbitrary text blobs are not imported without structure.

#### 5.4 Data Model (Living Narrative)
- Entities (ULID ids): Project, GraphVersion, Scene, Beat, Arc, Character, Location, Item, Theme, Relationship (edge), Annotation, Analysis.
- Versioning: Each write produces a new GraphVersion with parentVersionId; snapshots are immutable; a working set pointer indicates the latest for editing.
- Edges: Directed and typed (e.g., contains, advances, features, occursAt). Integrity and acyclicity enforced where applicable.
- Graph Service: All mutations pass through an application service (GraphWrite) that validates schema, referential integrity, and produces versioned deltas. SQLite stores nodes, edges, and supporting indexes to enable efficient neighbor and sequence queries.

#### 5.5 Event Schema Registry
- Envelope (JSON Schema v2020-12): eventName, eventVersion (semver), eventId (ULID), occurredAt, correlationId, causationId, idempotencyKey, producer, tenantId, payload.
- Core events v1: DirectiveIssued, NarrativeRestructured, SceneAdded, CharacterUpdated, CharacterReactionRequested, CharacterReactionProvided, SceneProposalReady, RefinementDirective, AgentError, AgentHeartbeat.
- Delivery: At-least-once; consumers must be idempotent on idempotencyKey; ordering not guaranteed; per-topic retry and DLQ policies defined.

#### 5.6 Agent Contract
- Inputs: Subscribed event schemas; tools (RAG query, read-only graph queries, GraphWrite mutations via API).
- Outputs: Emitted events conforming to registry; optional GraphWrite requests (never direct storage writes).
- Result envelope: status (success|retriable|fatal), emittedEvents[], annotations[], costUsage, logsRef.
- Idempotency: Agents must detect duplicate idempotencyKey and avoid duplicate side effects.
- Security: Service accounts with least privilege; all writes through GraphWrite API.

#### 5.7 UI Semantics and Real-time
- UI state: Local SQLite is the source of truth; streaming used only for transient previews. On completion, canonical results are committed to SQLite and reflected in the UI.
- Canvas: Structural edits show preview; on confirm, emit NarrativeRestructured with a structured diff.
- Inspector: Shows prose (read-only), linked entities, analyses, and allows highlight-based annotations that become structured feedback tied to nodes/edges.
- Governor: Displays in-flight events by correlationId and per-agent state (queued/running/succeeded/failed) with recent logsRef.

#### 5.8 Collaboration, Locking, and Undo/Redo
- Collaboration: Optimistic concurrency via version preconditions; on conflict, present merge UI; resulting decision emits NarrativeRestructured.
- Locking: Coarse scene/sequence locks during restructuring; visible owner and TTL; locks are advisory.
- Undo/Redo: Every user write is a reversible delta; per-user undo/redo replays authored deltas but persists as new versions; history is immutable.

#### 5.9 Security, Privacy, and Auth
- Auth: Firebase Authentication (MVP). RBAC: Owner, Editor, Viewer. Baton directives require Editor or above. Abstracted to allow future WorkOS SSO.
- Data: TLS in transit; KMS for secrets; audit logs for all writes and emissions; retention policies configurable; hard-delete SLA configurable per tenant.
- Safety: Generated content is passed through safety filters; unsafe content yields safetyReview annotations and is blocked unless Owner overrides.

#### 5.10 Observability and SLOs
- Metrics: Event latency per agent; queue depth; DLQ rate; UI update latency; RAG query latency; per-project LLM token/$$ usage.
- Tracing: Propagate correlationId/causationId UI→API→Event→Agent→GraphWrite; logs include these ids.
- Initial SLOs: UI live update p95 < 500ms post-write; Baton→first agent start p95 < 2s; Directive→SceneProposalReady p95 < 20s baseline; DLQ < 0.5%/day.

#### 5.11 Non-Functional Requirements
- Performance (10k scenes/project, 50 concurrent edits), accessibility (WCAG 2.1 AA), i18n-ready, cost budgets and throttles per project, Terraform IaC for all infra.

#### 5.12 Functional Requirements (EARS v1.1, canonical)
The following FRs restate and extend earlier requirements in EARS form and are authoritative.

- FR-2.1.1 (Ubiquitous): The system shall represent the Living Narrative as an interactive graph with nodes for story elements and typed edges for relationships.
- FR-2.1.2 (Ubiquitous): The system shall permanently display an Inspector panel providing contextual information and controls for the selected element.
- FR-2.1.3 (Event-driven): When the Conductor selects a node, the system shall display prose, linked entities, and agent-generated metadata in the Inspector.
- FR-2.1.4 (State-driven): While the Conductor is interacting with the Canvas, the system shall allow drag/connect/resequence operations to restructure narrative flow.
- FR-2.1.5 (Event-driven): When the Conductor confirms a structural change, the system shall emit NarrativeRestructured with a structured diff.
- FR-2.1.6 (Event-driven): When the Conductor highlights prose in the Inspector and submits feedback, the system shall create a structured annotation linked to the underlying narrative elements.

- FR-2.2.1 (Ubiquitous): The system shall provide a natural-language Baton as the primary method for issuing directives.
- FR-2.2.2 (Event-driven): When a directive is issued, the system shall emit DirectiveIssued with correlation metadata and idempotencyKey.
- FR-2.2.3 (State-driven): While a proposal is under review, the system shall expose Tuners (tone, pacing, style) and shall emit RefinementDirective when applied.
- FR-2.2.4 (State-driven): While agents are active, the system shall provide a Governor view of active events and per-agent progress.
- FR-2.2.5 (Ubiquitous): The system shall present contextual Nudges for plot, characters, and themes.
- FR-2.2.6 (Event-driven): When an agent is selected, the system shall display role, subscriptions, and recent activity.

- FR-3.2.1 (Ubiquitous): The system shall use a publish-subscribe Event Bus as the primary agent orchestration mechanism.
- FR-3.2.2 (Event-driven): When a directive is issued, the publisher shall post a structured event with idempotency and correlation metadata.
- FR-3.2.3 (State-driven): While active, each agent shall subscribe only to declared event types matching its expertise.
- FR-3.2.4 (Event-driven): When an agent consumes an event, it shall emit outputs as new events conforming to the registry.
- FR-3.2.5 (Event-driven): When targeted information is required, agents shall use A2A request/response events correlated by correlationId.
- FR-3.2.6 (Unwanted behavior): Where a duplicate idempotencyKey is detected, the system shall not apply a duplicate side effect and shall return the prior result.

- FR-4.3.1 (Ubiquitous): The system shall ensure durable execution for multi-step workflows with persisted state and resumability.
- FR-4.3.2 (Event-driven): When a directive requires a sequence of steps, the system shall orchestrate via Google Cloud Workflows and persist checkpoints at each step.
- FR-4.3.3 (Ubiquitous): The system shall store narrative changes atomically with validation against the canonical schema.
- FR-4.3.4 (Ubiquitous): The system shall implement checkpointing to resume from the last successful step after transient failures.
- FR-4.3.5 (Event-driven): When a processing failure occurs, the system shall emit AgentError with diagnostics and route the message per retry/DLQ policy.

- FR-5.1 (Ubiquitous): The system shall maintain versioned snapshots and support branch/fork and merge with conflict resolution.
- FR-5.2 (Event-driven): When a user performs an edit, the system shall capture an undoable delta and enable undo/redo within the session.
- FR-5.3 (Ubiquitous): The system shall enforce RBAC roles (Owner, Editor, Viewer); Baton directives require Editor.
- FR-5.4 (Event-driven): When writes occur, the system shall use optimistic concurrency (version preconditions) and surface a merge UI on conflict.
- FR-5.5 (Ubiquitous): The system shall enforce per-project LLM budget thresholds and throttle processing once reached, notifying the user.
- FR-5.6 (Ubiquitous): The system shall record audit logs for all writes and emissions correlated by correlationId.
- FR-5.7 (Event-driven): When generated content violates safety policy, the system shall attach a safetyReview annotation and block display unless overridden by Owner.
- FR-5.8 (Ubiquitous): The system shall expose operational metrics and traces across UI, API, events, agents, and graph writes.
- FR-5.9 (Event-driven): When the Conductor requests export, the system shall generate frozen artifacts (final copy and compendia) derived from the current GraphVersion.
- FR-5.10 (Event-driven): When import is performed (if enabled), the system shall transform inputs into structured deltas mapped to the narrative model and reject unstructured blobs.


### 6.0 Bootstrap and Templates (MVP)

#### 6.1 Bootstrap Modes (MVP scope)
- Template-first (primary): The Conductor selects a plot blueprint (e.g., Three-Act, Hero's Journey) and character archetypes, parameterizes theme/genre/tone, and the system seeds a graph skeleton (arcs, beats, placeholder scenes, character slots).
- Paste/Import assist (lightweight): Within the Bootstrap Wizard, the Conductor may paste Markdown/CSV/TXT snippets to populate template-required fields (e.g., character bios, setting notes). The system proposes structured fields from pasted content for review.
- Uploads-only (MVP): File uploads limited to Markdown, CSV, and plain text. No external connectors (e.g., Notion) in MVP.

#### 6.2 Bootstrap Wizard UX
- Steps: (1) Choose template and parameters; (2) Optional paste/upload; (3) Review structured proposals; (4) Apply as bootstrap branch (new GraphVersion).
- Review: Each proposed entity displays confidence and source excerpt. The Conductor can accept/edit/reject items; merge duplicates.

#### 6.3 Agents and Orchestration (MVP)
- No dedicated ingestion agents in MVP. Extraction/classification executed synchronously within the API tier or a single lightweight “Bootstrap Helper” function with strict cost limits.
- All outputs go through GraphWrite API to produce a versioned bootstrap branch.

#### 6.4 EARS Functional Requirements (Bootstrap)
- FR-6.1 (Ubiquitous): The system shall provide a template-first Bootstrap Wizard that seeds a narrative skeleton with parameterized templates and archetypes.
- FR-6.2 (Event-driven): When the Conductor confirms a template selection, the system shall emit TemplateApplied and create the corresponding graph skeleton as a new GraphVersion (bootstrap branch).
- FR-6.3 (Ubiquitous): The system shall allow the Conductor to paste or upload Markdown/CSV/TXT during bootstrap and shall propose structured fields for characters, settings, and arcs.
- FR-6.4 (State-driven): While reviewing proposals, the system shall allow accept/edit/reject and duplicate merges; on submit, it shall emit ImportReviewSubmitted.
- FR-6.5 (Event-driven): When the review is approved, the system shall apply changes via GraphWrite and shall emit ImportApplied and BootstrapGraphReady.
- FR-6.6 (Unwanted behavior): Where proposed entities match existing canonical entities above a threshold, the system shall not create duplicates without explicit merge.
- FR-6.7 (Ubiquitous): The system shall surface estimated cost and processing time for bootstrap parsing and enforce per-project LLM budgets.

#### 6.5 Export Formats (MVP)
- Final copy: Markdown (MVP). Optional: HTML later.
- Compendia: CSV/JSON/Markdown for story bible, character book, and lore book.
