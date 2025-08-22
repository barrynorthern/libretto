#!/usr/bin/env bash
set -euo pipefail

# Starts the local monolith binary (API + agents + store) with Bazel on dev port.
# Stop with Ctrl+C; all child processes are cleaned up.

API_PORT="${API_PORT:-8080}"


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
for name in API:${API_PORT}; do
  svc="${name%%:*}"; port="${name##*:}"
  if port_in_use "${port}"; then
    echo "Error: Port ${port} for ${svc} appears to be in use. Set ${svc}_PORT to a free port or stop the other process." >&2
    exit 1
  fi
done

echo "Starting monolith on :${API_PORT}"
PORT="${API_PORT}" bazel run //:libretto &
pids+=("$!")

# Simple readiness wait
sleep 2

echo "\nService started:"
echo "- Monolith: http://localhost:${API_PORT}"

echo "\nPress Ctrl+C to stop..."

# Wait on background jobs (portable; macOS bash lacks 'wait -n')
wait || true

