# ADR 0007: Monorepo and Bazel Build System

Status: Accepted
Date: 2025-08-11

## Context
We will host the frontend and backend in a monorepo. We want reproducible builds, clear dependency graphs, and consistent tooling for Go services.

## Decision
- Monorepo with Bazel as the primary build system for backend Go services and agents (rules_go + gazelle).
- Frontend (Next.js) initially built outside Bazel using package scripts; consider rules_nodejs later.
- Shared contracts and schemas versioned in-repo; codegen for Go (and TS for FE) as needed.

## Consequences
- Pros: Strong caching, clear build graphs, easy multi-service scaling in CI.
- Cons: Some upfront friction to author BUILD files.
- Mitigations: Use gazelle to generate BUILD files; start with a minimal set of Go targets and expand iteratively.

