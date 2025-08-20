# ADR 0008: IDL & RPC Strategy

Status: Accepted
Date: 2025-08-11

## Context
We need a single source of truth for service APIs and event payloads, typed clients for Go/TS, and browser-friendly transport without extra gateways.

## Decision
- Use Protobuf as the canonical IDL for services, event payloads, and the event envelope.
- Use Connect RPC for services (connect-go/connect-web). Same handlers can serve Connect/gRPC/gRPC-web if needed.
- Serialize event payloads and envelope with protojson on Pub/Sub.
- Deprecate hand-maintained JSON Schemas; generate JSON Schema from proto only when required by downstream tooling (not committed).

## Consequences
- Pros: One IDL; typed clients; strong compatibility checks (buf lint/breaking); web-friendly transport.
- Cons: Adds codegen to the workflow.
- Mitigations: Check in generated code for now; CI enforces that it is up-to-date.

