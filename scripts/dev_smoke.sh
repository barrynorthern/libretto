#!/usr/bin/env bash
set -euo pipefail

API_PORT="${API_PORT:-8080}"
PLOT_PORT="${PLOT_PORT:-8081}"
GRAPHWRITE_PORT="${GRAPHWRITE_PORT:-8082}"

PASS=0
FAIL=0

bold() { printf "\033[1m%s\033[0m" "$1"; }
green() { printf "\033[32m%s\033[0m" "$1"; }
red() { printf "\033[31m%s\033[0m" "$1"; }

check() {
  local name="$1"; shift
  local cmd=("$@")
  if output="$(${cmd[@]} 2>/dev/null)"; then
    echo "$(green "✔") $(bold "$name")"
  else
    echo "$(red "✘") $(bold "$name")" >&2
    echo "$output" >&2 || true
    FAIL=$((FAIL+1))
    return 1
  fi
  PASS=$((PASS+1))
}

# API health
check "API health" curl -sS "http://localhost:${API_PORT}/healthz"

# Baton IssueDirective (expect correlation_id)
check "Baton IssueDirective" curl -sS -X POST -H 'Content-Type: application/json' \
  --data '{"text":"Introduce a betrayal","act":"2","target":"protagonist"}' \
  "http://localhost:${API_PORT}/libretto.baton.v1.BatonService/IssueDirective"

# Plot Weaver stub
check "Plot Weaver stub" curl -sS -X POST "http://localhost:${PLOT_PORT}/"

# GraphWrite Apply (expect graphVersionId)
check "GraphWrite Apply" curl -sS -X POST -H 'Content-Type: application/json' \
  --data '{"parentVersionId":"01JROOT","deltas":[{"op":"create","entityType":"Scene","entityId":"sc-1","fields":{"title":"Test"}}]}' \
  "http://localhost:${GRAPHWRITE_PORT}/libretto.graph.v1.GraphWriteService/Apply"

TOTAL=$((PASS+FAIL))

echo
if [ "$FAIL" -eq 0 ]; then
  echo "$(green "All ${TOTAL} checks passed")"
  exit 0
else
  echo "$(red "${FAIL}/${TOTAL} checks failed")"
  exit 1
fi

