# ADR 0002: Orchestration Choice for MVP

Status: Accepted
Date: 2025-08-11

## Context
We need durable multi-step orchestration. Temporal is a strong fit but adds cost and infra. We need speed for MVP.

## Decision
- Use Google Cloud Workflows for MVP orchestration.
- Design event schemas and agent contracts to be Temporal-ready (correlation/idempotency/checkpoints).

## Consequences
- Pros: Minimal ops, native GCP integration, fast to ship.
- Cons: Less feature-rich than Temporal for long-lived workflows.
- Mitigations: Re-evaluate after 2–3 complex flows; migrate when necessary without changing surface contracts.

