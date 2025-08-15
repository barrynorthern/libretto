#!/usr/bin/env bash
set -euo pipefail

# Starts the local stack (API, Plot Weaver, GraphWrite) with Bazel on dev ports.
# Stop with Ctrl+C; all child processes are cleaned up.

API_PORT="${API_PORT:-8080}"
PLOT_PORT="${PLOT_PORT:-8081}"
GRAPHWRITE_PORT="${GRAPHWRITE_PORT:-8082}"

pids=()

cleanup() {
  echo "\nShutting down services..."
  for pid in "${pids[@]:-}"; do
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done
}
trap cleanup INT TERM EXIT

echo "Starting API on :${API_PORT}"
PORT="${API_PORT}" bazel run //services/api:api &
pids+=("$!")

sleep 0.5

echo "Starting Plot Weaver on :${PLOT_PORT}"
PORT="${PLOT_PORT}" bazel run //services/agents/plotweaver:plotweaver &
pids+=("$!")

sleep 0.5

echo "Starting GraphWrite on :${GRAPHWRITE_PORT}"
PORT="${GRAPHWRITE_PORT}" bazel run //services/graphwrite:graphwrite &
pids+=("$!")

# Simple readiness wait
sleep 2

echo "\nServices started:"
echo "- API:          http://localhost:${API_PORT}"
echo "- Plot Weaver:  http://localhost:${PLOT_PORT}"
echo "- GraphWrite:   http://localhost:${GRAPHWRITE_PORT}"

echo "\nTip: In a new terminal, run ./scripts/dev_smoke.sh to smoke test endpoints."

echo "\nPress Ctrl+C to stop..."

# Wait on background jobs
wait -n || true

