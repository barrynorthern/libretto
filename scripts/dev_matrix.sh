#!/usr/bin/env bash
set -euo pipefail

# Runs smoke tests in a small matrix of publisher states.
# Assumes dev_up.sh is already running services locally.

API_PORT="${API_PORT:-8080}"

run_case_env() {
  local name="$1"; shift
  echo ""
  echo "=== Matrix case: $name ==="
  "$@" ./scripts/dev_smoke.sh || true
}

run_case_env "NOP publisher" env -u PUBLISHER -u PUBSUB_ENABLED
run_case_env "Pub/Sub publisher (back-compat)" env PUBSUB_ENABLED=true
run_case_env "DevPush publisher" env PUBLISHER=devpush

echo ""
echo "Matrix complete."

