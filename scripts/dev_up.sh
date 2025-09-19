#!/usr/bin/env bash
set -euo pipefail

# Starts the local monolith binary with cross-project functionality
# Stop with Ctrl+C; all child processes are cleaned up.

API_PORT="${API_PORT:-8080}"
DASHBOARD_PORT="${DASHBOARD_PORT:-9000}"

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
for name in API:${API_PORT} DASHBOARD:${DASHBOARD_PORT}; do
  svc="${name%%:*}"; port="${name##*:}"
  if port_in_use "${port}"; then
    echo "Error: Port ${port} for ${svc} appears to be in use. Set ${svc}_PORT to a free port or stop the other process." >&2
    exit 1
  fi
done

echo "Starting Libretto services..."

echo "Starting monolith on :${API_PORT}"
PORT="${API_PORT}" bazel run //:libretto &
pids+=("$!")

echo "Starting dashboard on :${DASHBOARD_PORT}"
go run cmd/dashboard/main.go -port="${DASHBOARD_PORT}" &
pids+=("$!")

# Simple readiness wait
sleep 3

echo "\nServices started:"
echo "- Monolith API: http://localhost:${API_PORT}"
echo "- Dashboard: http://localhost:${DASHBOARD_PORT}"
echo "- Cross-project demo: http://localhost:${DASHBOARD_PORT}/demo"

echo "\nPress Ctrl+C to stop..."

# Wait on background jobs (portable; macOS bash lacks 'wait -n')
wait || true

