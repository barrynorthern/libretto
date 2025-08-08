# GraphWrite API (Draft)

Status: Draft (MVP)

## Principles
- Single entry point for all graph mutations.
- Validates schema, referential integrity, and produces immutable GraphVersions with deltas.
- Idempotent on `idempotencyKey`.

## Entities (reference)
- See `schemas/model/entities.schema.json`.

## Endpoints

### POST /graph/versions/{projectId}/apply
Apply a batch of deltas as a new GraphVersion.

Request
```json
{
  "parentVersionId": "01J...",
  "idempotencyKey": "bootstrap-...",
  "deltas": [
    {"op": "create", "entityType": "Character", "entity": {"id": "01J...", "name": "Protagonist"}},
    {"op": "create", "entityType": "Arc", "entity": {"id": "01J...", "name": "Main Plot"}},
    {"op": "create", "entityType": "Relationship", "entity": {"from": "scene-1", "to": "arc-1", "type": "advances"}}
  ]
}
```

Response
```json
{
  "graphVersionId": "01J...",
  "applied": 3,
  "warnings": []
}
```

### POST /graph/bootstrap/{projectId}/template
Create a bootstrap branch and skeleton from a template.

Request
```json
{
  "templateId": "three-act-v1",
  "parameters": {"genre": "mystery", "tone": "noir"},
  "idempotencyKey": "tpl-01J..."
}
```

Response
```json
{
  "graphVersionId": "01J...",
  "branch": "bootstrap/01J..."
}
```

### POST /graph/bootstrap/{projectId}/apply-review
Apply reviewed proposals to the bootstrap branch.

Request
```json
{
  "bootstrapBranchVersionId": "01J...",
  "idempotencyKey": "review-01J...",
  "decisions": [
    {"proposedId": "char-01", "action": "accept"},
    {"proposedId": "char-02", "action": "merge", "mergeWithId": "char-01"}
  ]
}
```

Response
```json
{
  "graphVersionId": "01J...",
  "appliedDeltaCount": 12
}
```

## Errors
- 400: schema/validation error
- 404: parentVersionId or project not found
- 409: precondition failed (optimistic concurrency)
- 409: idempotencyKey replay (returns prior result)
- 500: internal

## Preconditions & Concurrency
- `parentVersionId` must exist and be current for linear histories; merges use a dedicated merge endpoint (future).
- Firestore write preconditions enforced per document.

