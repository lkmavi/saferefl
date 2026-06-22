#!/usr/bin/env bash
# Run all benchmarks on local hardware and regenerate benchmarks/RESULTS.md + README.md.
#
# Usage:
#   ./scripts/bench.sh              # count=5, benchtime=1s (default)
#   ./scripts/bench.sh 10 2s        # count=10, benchtime=2s
#   COUNT=1 ./scripts/bench.sh      # env-var override (quick smoke check)
set -euo pipefail

cd "$(dirname "$0")/.."

COUNT="${COUNT:-${1:-5}}"
BENCHTIME="${2:-1s}"
GOVERSION=$(go env GOVERSION | sed 's/^go//' | cut -d. -f1-2)

echo "▶ benchmarks (go=${GOVERSION}, count=${COUNT}, benchtime=${BENCHTIME})..."

(cd benchmarks && go test ./... \
    -bench=. \
    -benchmem \
    -count="${COUNT}" \
    -benchtime="${BENCHTIME}" \
    -run='^$') \
    | go run ./scripts/benchfmt/main.go \
        -output benchmarks/RESULTS.md \
        -go-version "${GOVERSION}" \
        -readme README.md

echo ""
echo "✓ benchmarks/RESULTS.md"
echo "✓ README.md"
