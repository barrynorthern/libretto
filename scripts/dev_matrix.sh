#!/usr/bin/env bash
set -euo pipefail

# Runs smoke tests in a small matrix of PUBSUB_ENABLED states.
# Assumes dev_up.sh is already running services locally.

API_PORT="${API_PORT:-8080}"

run_case() {
  local name="$1"; shift
  echo "\n=== Matrix case: $name ==="
  PUBSUB_ENABLED="$1" ./scripts/dev_smoke.sh || true
}

run_case "NOP publisher" ""
run_case "Pub/Sub publisher" "true"

echo "\nMatrix complete."

