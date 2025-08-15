#!/usr/bin/env bash
set -euo pipefail

API_PORT="${API_PORT:-8080}"
PLOT_PORT="${PLOT_PORT:-8081}"
GRAPHWRITE_PORT="${GRAPHWRITE_PORT:-8082}"

set -x

# API health
curl -sS "http://localhost:${API_PORT}/healthz" | cat

# Baton IssueDirective
curl -sS -X POST -H 'Content-Type: application/json' \
  --data '{"text":"Introduce a betrayal","act":"2","target":"protagonist"}' \
  "http://localhost:${API_PORT}/libretto.baton.v1.BatonService/IssueDirective" | cat

# Plot Weaver stub
curl -sS -X POST "http://localhost:${PLOT_PORT}/" | cat

# GraphWrite Apply
curl -sS -X POST -H 'Content-Type: application/json' \
  --data '{"parentVersionId":"01JROOT","deltas":[{"op":"create","entityType":"Scene","entityId":"sc-1","fields":{"title":"Test"}}]}' \
  "http://localhost:${GRAPHWRITE_PORT}/libretto.graph.v1.GraphWriteService/Apply" | cat

