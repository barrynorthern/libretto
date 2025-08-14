# 000 – Repo scaffolding and CI baseline (retrospective)

Status: Done (merged)
Owner: barrynorthern
Date completed: 2025-08-14

## Context
We established the monorepo with Bazel, initial CI, ADRs, and baseline docs to enable a spec-driven workflow.

## Scope delivered
- Monorepo with Bazel (rules_go + gazelle); basic BUILD targets
- CI: bazel build/test, buf lint/generate
- ADRs: backend language/runtime, monorepo/bazel
- Docs: spec.md v1.1; execution plan; API/agent drafts

## Acceptance criteria
- Bazel builds and tests pass in CI
- Protos lint and generate cleanly
- Repo layout matches execution plan

## Links
- ADR 0006, ADR 0007
- .github/workflows/ci.yml

## Notes / outcomes
- Decided to keep frontend out of Bazel initially (ADR-0007)
- NPM-based FE to come later; contracts generated with buf for TS when needed

