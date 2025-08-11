# ADR 0006: Backend Language and Runtime

Status: Accepted
Date: 2025-08-11

## Context
We need a consistent backend stack that supports agents and API services with strong performance, observability, and simple deployment.

## Decision
- Language: Go 1.22+ for all backend services and agents.
- Runtime: Cloud Run (Go) with Pub/Sub push subscriptions for agents; Workflows for multi-step orchestration (MVP).
- Data: Firestore as canonical store for the Living Narrative; optional Cloud SQL (PostgreSQL) for auth/billing later.
- Clients: sqlc generates type-safe Go clients for Postgres when used.

## Consequences
- Pros: High performance, low memory footprint, single language across backend; straightforward deployment on Cloud Run.
- Cons: Fewer off-the-shelf LLM SDKs than Node/Python.
- Mitigations: Provider-agnostic Agent Contract; thin adapters for Vertex/OpenAI as needed.

