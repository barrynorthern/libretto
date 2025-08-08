# Bootstrap Helper (Draft)

Status: Draft (MVP)

## Purpose
Synchronous, lightweight parsing of pasted/uploaded Markdown/CSV/TXT within the Bootstrap Wizard to propose structured entities for review.

## Input
```json
{
  "projectId": "01J...",
  "mode": "paste" | "upload",
  "files": [
    {"name": "characters.md", "contentType": "text/markdown", "contentBase64": "..."}
  ],
  "hints": {"expect": ["Character", "Setting", "Arc"]},
  "budget": {"maxTokens": 20000, "maxCostUSD": 1.00}
}
```

## Output
```json
{
  "proposals": [
    {"entityType": "Character", "proposedId": "p-01", "name": "Avery", "archetype": "Protagonist", "bio": "...", "confidence": 0.92, "sourceExcerpt": "..."},
    {"entityType": "Setting", "proposedId": "p-02", "name": "Grey Harbor", "description": "...", "confidence": 0.87}
  ],
  "warnings": ["Low confidence for 2 items"],
  "usage": {"tokens": 1540, "estimatedCostUSD": 0.12}
}
```

## Constraints
- Size limits: total input 1MB (MVP). Fail fast with friendly error when exceeded.
- Cost limits enforced by `budget`; helper must short-circuit and return partial results.
- No external connectors; no persistent writes; proposals only.

## Errors
- 400: invalid input or unsupported content type
- 413: payload too large
- 429: budget exceeded
- 500: internal

## Notes
- Parser should be deterministic where possible (regex/heuristics first) and only use LLM for ambiguous cases.
- Provenance: include excerpt or offsets from source to support review.

