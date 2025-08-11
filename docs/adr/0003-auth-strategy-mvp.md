# ADR 0003: Auth Strategy for MVP

Status: Accepted
Date: 2025-08-11

## Context
We need authentication and RBAC quickly for a single-tenant MVP while preserving a path to enterprise SSO.

## Decision
- Use Firebase Authentication for MVP.
- Implement RBAC (Owner, Editor, Viewer) in-app; authorize on GraphWrite and Baton.
- Abstract auth to allow future WorkOS integration without UI changes.

## Consequences
- Pros: Fast, low cost, integrates with Firestore rules.
- Cons: Not enterprise SSO out of the box.
- Mitigations: Add WorkOS when needed; keep user model and sessions decoupled.

