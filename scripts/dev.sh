#!/usr/bin/env bash
set -euo pipefail

API_PORT=8080
API_URL="http://localhost:${API_PORT}/api/exercises?pageSize=1"
GO_BINARY=".dev-server"

PIDS=()

cleanup() {
  echo ""
  echo "Shutting down..."
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  # Wait briefly, then force-kill stragglers
  sleep 0.3
  for pid in "${PIDS[@]}"; do
    kill -9 "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
  done
  rm -f "$GO_BINARY"
  echo "Done."
}

trap cleanup EXIT

# Build Go binary (avoids go run's grandchild problem)
echo "Building Go API..."
go build -o "$GO_BINARY" .

# Start API
DEV=true AUTH_FALLBACK_USER=anon "./$GO_BINARY" &
PIDS+=($!)
echo "API server starting (pid ${PIDS[0]})..."

# Wait for API readiness
for i in $(seq 1 30); do
  if curl -sf "$API_URL" > /dev/null 2>&1; then
    break
  fi
  if ! kill -0 "${PIDS[0]}" 2>/dev/null; then
    echo "API process died."
    exit 1
  fi
  sleep 0.5
done

if ! curl -sf "$API_URL" > /dev/null 2>&1; then
  echo "API did not become ready in time."
  exit 1
fi
echo "API ready."

# Start Angular dev server
echo "Starting Angular dev server..."
cd web
CI=true npx ng serve "$@" &
PIDS+=($!)
cd ..

# Wait for any child to exit — that means something crashed
wait -n "${PIDS[@]}" 2>/dev/null || true
echo "A child process exited unexpectedly."
exit 1
