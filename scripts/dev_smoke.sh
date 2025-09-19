#!/usr/bin/env bash
set -euo pipefail

API_PORT="${API_PORT:-8080}"
DASHBOARD_PORT="${DASHBOARD_PORT:-9000}"

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

echo "🧪 Testing Libretto Cross-Project Functionality"
echo

# Monolith API tests
echo "📡 Testing Monolith API..."
http_check "API health" GET "http://localhost:${API_PORT}/healthz"
http_check "Baton IssueDirective" POST "http://localhost:${API_PORT}/libretto.baton.v1.BatonService/IssueDirective" '{"text":"Elena discovers an ancient artifact","act":"1","target":"protagonist"}'

# Dashboard tests
echo
echo "🎛️  Testing Dashboard..."
http_check "Dashboard home" GET "http://localhost:${DASHBOARD_PORT}/"
http_check "Dashboard demo" GET "http://localhost:${DASHBOARD_PORT}/demo"

# Cross-project functionality test
echo
echo "🌐 Testing Cross-Project Features..."
echo "Running Elena Stormwind cross-project demo..."
if go test -v ./internal/graphwrite -run TestCrossProjectCharacterArcs >/dev/null 2>&1; then
  echo "$(green "✔") $(bold "Cross-project character continuity test")"
  PASS=$((PASS+1))
else
  echo "$(red "✘") $(bold "Cross-project character continuity test")" >&2
  FAIL=$((FAIL+1))
fi

TOTAL=$((PASS+FAIL))

echo
if [ "$FAIL" -eq 0 ]; then
  echo "$(green "🎉 All ${TOTAL} checks passed - Elena must always be Elena!")"
  exit 0
else
  echo "$(red "💥 ${FAIL}/${TOTAL} checks failed")"
  exit 1
fi

