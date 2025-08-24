# 00008 – Local vector DB (sqlite-vec), RAG pipeline, and Wails desktop scaffold (MVP)

Status: Proposed
Owner: barrynorthern

## Context
Context management requires local retrieval for project-scoped memory. To minimize dependencies, use SQLite with sqlite-vec for embeddings and similarity search.

## Goal
Implement a minimal RAG pipeline using sqlite-vec:
- Schema: documents(project_id, doc_id, kind, text, embedding BLOB)
- Functions: upsert embedding; search k-NN by cosine or L2
- Integrate Retriever interface from internal/context

## Scope
- Add sqlite-vec dependency and migration
- Create embedding pipeline stub (Ollama by default; optional API-key providers)
- Implement internal/context/sqlitev retriever (Search)
- Scaffold Wails desktop app (React + shadcn/Tailwind) with a Scenes page calling a binding
- Unit tests with small vectors to validate k-NN results

## Acceptance criteria
- Build/test green with sqlite-vec enabled locally
- Retriever.Search returns expected nearest neighbors on fixture data
- Basic upsert and search demonstrated in a test

## Out of scope
- Full document chunker, advanced ranking, or caching
- Cross-project/global memory

