#!/usr/bin/env bash
# Drift guard: regenerate the bindings and fail if the committed files are stale.
# Run in CI (and locally) so a schema change without a regenerate can't land.
set -euo pipefail
cd "$(dirname "$0")/.."

scripts/generate.sh >/dev/null

if ! git diff --quiet -- sdui/contract/contract.gen.go ts/contract.gen.ts; then
  echo "ERROR: generated bindings are stale. Run scripts/generate.sh and commit." >&2
  git --no-pager diff --stat -- sdui/contract/contract.gen.go ts/contract.gen.ts >&2
  exit 1
fi
echo "generated bindings are up to date."
