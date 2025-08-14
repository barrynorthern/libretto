# 001 – Baton API and Plot Weaver stub vertical slice (retrospective)

Status: Done (merged)
Owner: barrynorthern
Date completed: 2025-08-14

## Context
We delivered a thin end-to-end slice to exercise the narrative event path from API to an agent and back, without persistence.

## Scope delivered
- BatonService (Connect) IssueDirective endpoint
- Event envelope published by API using NopPublisher (local)
- Plot Weaver agent (HTTP handler) consumes trigger and emits SceneProposalReady (stub)
- Protos for baton.v1 and graphwrite.v1; TS/Go codegen in CI
- Health endpoints and basic tests

## Acceptance criteria
- curl to BatonService/IssueDirective returns 200 and correlation id
- Plot Weaver stub returns 200 and logs SceneProposalReady event
- bazel build //... and bazel test //... succeed in CI

## Links
- services/api/*
- services/agents/plotweaver/*
- proto/libretto/baton/v1/baton.proto
- proto/libretto/graph/v1/graphwrite.proto

## Notes / outcomes
- Confirmed ports and env var defaults for local runs
- Deferred real Pub/Sub and Firestore wiring to next iteration

