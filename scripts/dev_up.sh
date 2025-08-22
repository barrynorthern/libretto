#!/usr/bin/env bash
set -euo pipefail

# Starts the local stack (API, Plot Weaver, Narrative Ingest, GraphWrite) with Bazel on dev ports.
# Stop with Ctrl+C; all child processes are cleaned up.

API_PORT="${API_PORT:-8080}"
PLOT_PORT="${PLOT_PORT:-8081}"
NARRATIVE_PORT="${NARRATIVE_PORT:-8083}"
GRAPHWRITE_PORT="${GRAPHWRITE_PORT:-8082}"

port_in_use() {
  local port="$1"
  if command -v lsof >/dev/null 2>&1; then
    lsof -i tcp:"${port}" -sTCP:LISTEN -Pn >/dev/null 2>&1
  else
    # Fallback: try connecting quickly; if success, assume in use
    (exec 3<>/dev/tcp/127.0.0.1/"${port}") >/dev/null 2>&1 || return 1
  fi
}

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

# Preflight: ensure ports are free
for name in API:${API_PORT} PLOT:${PLOT_PORT} NARRATIVE:${NARRATIVE_PORT} GRAPHWRITE:${GRAPHWRITE_PORT}; do
  svc="${name%%:*}"; port="${name##*:}"
  if port_in_use "${port}"; then
    echo "Error: Port ${port} for ${svc} appears to be in use. Set ${svc}_PORT to a free port or stop the other process." >&2
    exit 1
  fi
done

echo "Starting API on :${API_PORT}"
PORT="${API_PORT}" bazel run //services/api:api &
pids+=("$!")

sleep 0.5

echo "Starting Plot Weaver on :${PLOT_PORT}"
PORT="${PLOT_PORT}" bazel run //services/agents/plotweaver:plotweaver &
pids+=("$!")

sleep 0.5

echo "Starting Narrative Ingest on :${NARRATIVE_PORT}"
PORT="${NARRATIVE_PORT}" bazel run //services/agents/narrativeingest:narrativeingest &
pids+=("$!")

sleep 0.5

echo "Starting GraphWrite on :${GRAPHWRITE_PORT}"
PORT="${GRAPHWRITE_PORT}" bazel run //services/graphwrite:graphwrite &
pids+=("$!")

# Simple readiness wait
sleep 2

echo "\nServices started:"
echo "- API:              http://localhost:${API_PORT}"
echo "- Plot Weaver:      http://localhost:${PLOT_PORT}"
echo "- Narrative Ingest: http://localhost:${NARRATIVE_PORT}"
echo "- GraphWrite:       http://localhost:${GRAPHWRITE_PORT}"

echo "\nTip: In a new terminal, run ./scripts/dev_smoke.sh or ./scripts/dev_matrix.sh (matrix runs NOP, PUBSUB back-compat, and DevPush)"

echo "\nPress Ctrl+C to stop..."

# Wait on background jobs (portable; macOS bash lacks 'wait -n')
wait || true

