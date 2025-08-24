# 00007 – DDD overview and Context Manager foundation

Status: Proposed
Owner: barrynorthern

## Context
We’ve pivoted to a modular monolith with Wails+React and a local-first strategy. The primary differentiator is effective context management (task-specific memory + RAG + model selection). We need a DDD doc and an initial Context Manager seam.

## Goal
- Author a DDD document (mermaid diagrams) capturing bounded contexts, modules, domain model, and the directive→scene sequence with context.
- Introduce an internal/context package with interfaces for Context Builder, Retrieval, and Model Selection (no heavy implementation yet).

## Scope
- Docs
  - docs/ddd/overview.md with mermaid: container/module, class diagram, sequence
- Code
  - internal/context: package skeleton with interfaces:
    - Builder: Build(projectID, directive) -> ContextBundle
    - Retriever: Search(projectID, query, k) -> []Result
    - ModelSelector: Choose(task, complexity, budget) -> ModelSpec
  - Choose vector DB: sqlite-vec (preferred for simplicity). Add TODO placeholder; implement in follow-up.

## Acceptance criteria
- DDD doc present with the three diagrams
- internal/context package compiles with interfaces and TODOs

## Out of scope
- Full RAG implementation
- Embedding pipeline and chunking
- Model registry and selection policies (stubs only)

