#!/usr/bin/env bash
# Generate the per-language contract bindings from the single source of truth,
# schema/sdui.schema.json. Do NOT hand-edit the generated files.
#
#   Go   -> sdui/contract/contract.gen.go   (package contract)
#   TS   -> ts/contract.gen.ts
#   Dart -> add when a Dart client lands (quicktype --lang dart)
#
# Run: scripts/generate.sh   (requires npx + gofmt on PATH)
set -euo pipefail
cd "$(dirname "$0")/.."

SCHEMA="schema/sdui.schema.json"
QT=(npx --yes quicktype -s schema --top-level MosaicSDUI "$SCHEMA")

echo "generating Go -> sdui/contract/contract.gen.go"
"${QT[@]}" --lang go --package contract --just-types-and-package -o sdui/contract/contract.gen.go
# The consolidated schema has a wrapper root (so the generator reaches every
# top-level type). Strip that unused wrapper struct, then mark as generated.
sed -i '/^type MosaicSDUI struct {/,/^}/d' sdui/contract/contract.gen.go
sed -i '1s|^|// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.\n|' sdui/contract/contract.gen.go
gofmt -w sdui/contract/contract.gen.go

echo "generating TypeScript -> ts/contract.gen.ts"
"${QT[@]}" --lang typescript --just-types -o ts/contract.gen.ts
sed -i '/^export interface MosaicSDUI {/,/^}/d' ts/contract.gen.ts
sed -i '1s|^|// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.\n|' ts/contract.gen.ts

echo "done."
