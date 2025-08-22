# 00007 – GraphWrite persistence via Firestore emulator (persistence step 2)

Status: Proposed
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
We have an in-memory GraphWrite server and a plan to call Apply from a consumer of SceneProposalReady. To move closer to the true vertical slice, we need GraphWrite to persist deltas into a Firestore emulator with a minimal schema.

## Goal
Implement Firestore emulator-backed persistence for GraphWrite.Apply, mapping incoming Deltas into documents.

## Scope
- Infra/Tooling
  - Add Firestore emulator to dev stack (Make: start/stop)
  - Configure emulator env for GraphWrite (`FIRESTORE_EMULATOR_HOST`, project)
- GraphWrite service
  - Add a Firestore-backed Store implementation alongside InMemoryStore
  - Map Delta(op=create, entity_type=Scene) to collection `nodes_scene` with fields {title, summary}
  - Return a new graph_version_id (e.g., ULID or UUID) and applied count
- Tests
  - Unit test mapping function(s)
  - Integration-style test behind an env flag to run against emulator locally

## Acceptance criteria
- `make dev-up` starts emulator and GraphWrite uses it when configured
- Apply requests with a Scene create delta persist a document to nodes_scene
- Logs include graph_version_id and applied count

## Non-functional requirements
- Keep emulator boot + test under a minute on CI
- No production credentials required; emulator only

## Out of scope
- Adjacency indexes, edges, or multi-entity transactions
- Security rules beyond minimal local defaults

## References
- Google Firestore emulator docs
- proto/libretto/graph/v1/graphwrite.proto

