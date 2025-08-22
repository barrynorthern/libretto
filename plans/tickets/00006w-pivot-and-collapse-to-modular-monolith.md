# 00006w – Pivot and collapse to a Modular Monolith (MVP)

Status: Proposed
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
We are over-indexed on distributed seams for a single-user app. To accelerate the MVP and reduce failure modes, we will pivot to a modular monolith: one binary, internal packages for agents and graphwrite, synchronous orchestration.

## Goal
Collapse the current multi-service layout into a single binary while preserving clean interfaces and future extractability. Deliver the thin vertical slice end-to-end with minimal moving parts.

## Scope
- Planning
  - Add ADR 0011 documenting the pivot (done in this branch)
  - Mark tickets 00006/00007/00008 as Won't Do (superseded)
- Code structure
  - Create cmd/libretto/main.go (single binary)
  - Add internal/app/orchestrator driving: Baton -> PlotWeaver -> Narrative -> GraphWrite
  - Move code into internal packages: internal/agents/{plotweaver,narrative}, internal/graphwrite
  - Remove DevPush/push from the happy path; keep Connect handlers for API
- Persistence
  - Keep in-memory store initially; prepare interface for Firestore emulator next ticket
- Tests
  - Adapt existing tests to internal packages and direct calls
- Scripts
  - Simplify dev_up.sh to run single binary

## Acceptance criteria
- `make dev-up` starts one process; issuing a directive produces a scene proposal and applies it via GraphWrite in-process
- `go test ./...` green

## Notes on UI
- UI will be a separate Next.js app (CSR) using an off-the-shelf design system (e.g., Mantine, Chakra, or MUI). Backend serves JSON APIs only. A follow-up ticket will scaffold the UI repo and wire basic calls.

## Out of scope
- Firestore emulator wiring (next ticket)
- Real Pub/Sub or multi-process orchestration

