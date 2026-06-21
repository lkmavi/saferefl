#!/usr/bin/env bash
# Run all benchmarks and auto-generate benchmarks/RESULTS.md + inject into README.md.
#
# Usage:
#   ./scripts/bench.sh              # count=5, benchtime=1s (default)
#   ./scripts/bench.sh 10 2s        # count=10, benchtime=2s
set -euo pipefail

cd "$(dirname "$0")/.."

COUNT="${1:-5}"
BENCHTIME="${2:-1s}"

echo "▶ benchmarks (count=${COUNT}, benchtime=${BENCHTIME})..."

go test ./benchmarks/... \
    -bench=. \
    -benchmem \
    -count="${COUNT}" \
    -benchtime="${BENCHTIME}" \
    | go run ./scripts/benchfmt/main.go \
        -output benchmarks/RESULTS.md \
        -readme README.md

echo ""
echo "✓ benchmarks/RESULTS.md"
echo "✓ README.md"
