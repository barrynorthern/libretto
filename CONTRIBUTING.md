# Contributing to Libretto

Thank you for contributing. This repo follows Spec‑Driven Development.

## Ground rules
- The spec (spec.md, v1.1+) is the single source of truth. Update the spec and schemas first, then code.
- All changes must reference one or more EARS Functional Requirements (FRs) and/or NFRs.
- No secrets in the repo. Use macOS Keychain or env vars locally; use KMS and Terraform in cloud.

## Workflow
1) Propose deltas in spec.md (add/modify sections and FRs).
2) Update JSON Schemas under /schemas (events, model). Keep them versioned and backward compatible.
3) Update contracts/docs (e.g., GraphWrite API, agent contracts, templates) under /docs.
4) Implement changes with tests and observability.
5) Submit PR with the checklist below.

## PR checklist
- [ ] Spec updated (sections and/or FRs referenced)
- [ ] JSON Schemas added/updated (schemas/…)
- [ ] Tests updated: unit, integration, and E2E where applicable
- [ ] Observability: logs (correlationId/causationId), metrics, traces
- [ ] Terraform updated (if infra affects)
- [ ] Firestore security rules updated/tests added (if applicable)
- [ ] Export/Bootstrap/GraphWrite docs updated (if applicable)
- [ ] Rollback plan noted

## Commit messages
- Conventional-ish style preferred: feat:, fix:, chore:, docs:, refactor:, test:, infra:
- Reference FR ids, e.g. "feat(FR-6.1): bootstrap wizard template selection UI"

## Running validation locally
- Ensure you have Bazelisk installed (brew install bazelisk)
- Repo pins Bazel via .bazelversion to keep local and CI aligned
- Typical local flow (mirrors CI):
  - bazel version
  - bazel build //...
  - bazel test //... --test_output=errors
- When adding imports or generating code, update BUILD files with Gazelle:
  - bazel run //:gazelle

## Code of Conduct
- Be respectful and constructive. Assume positive intent. Focus on clarity and testability.

