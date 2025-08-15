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

http_check() {
  local name="$1"; shift
  local method="$1"; shift
  local url="$1"; shift
  local data="${1:-}"
  local tmp
  tmp=$(mktemp)
  local http
  if [ -n "$data" ]; then
    http=$(curl -sS -o "$tmp" -w "%{http_code}" -X "$method" -H 'Content-Type: application/json' --data "$data" "$url" || true)
  else
    http=$(curl -sS -o "$tmp" -w "%{http_code}" -X "$method" "$url" || true)
  fi
  if [ "$http" = "200" ]; then
    echo "$(green "✔") $(bold "$name")"
    PASS=$((PASS+1))
  else
    echo "$(red "✘") $(bold "$name")" >&2
    echo "Status: $http" >&2
    echo "Body:" >&2
    sed -e 's/.*/  &/' "$tmp" >&2 || true
    FAIL=$((FAIL+1))
  fi
  rm -f "$tmp"
}

# API health
http_check "API health" GET "http://localhost:${API_PORT}/healthz"

# Baton IssueDirective (expect 200)
http_check "Baton IssueDirective" POST "http://localhost:${API_PORT}/libretto.baton.v1.BatonService/IssueDirective" '{"text":"Introduce a betrayal","act":"2","target":"protagonist"}'

# Plot Weaver stub (expect 200)
http_check "Plot Weaver stub" POST "http://localhost:${PLOT_PORT}/"

# GraphWrite Apply (expect 200)
http_check "GraphWrite Apply" POST "http://localhost:${GRAPHWRITE_PORT}/libretto.graph.v1.GraphWriteService/Apply" '{"parentVersionId":"01JROOT","deltas":[{"op":"create","entityType":"Scene","entityId":"sc-1","fields":{"title":"Test"}}]}'

TOTAL=$((PASS+FAIL))

echo
if [ "$FAIL" -eq 0 ]; then
  echo "$(green "All ${TOTAL} checks passed")"
  exit 0
else
  echo "$(red "${FAIL}/${TOTAL} checks failed")"
  exit 1
fi

