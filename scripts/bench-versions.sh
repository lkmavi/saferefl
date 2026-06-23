#!/usr/bin/env bash
# Cross-version benchmark runner via Docker.
#
# Usage:
#   ./scripts/bench-versions.sh                    # all versions, sequential run
#   ./scripts/bench-versions.sh 1.24 1.25          # specific versions
#   REBUILD=0 ./scripts/bench-versions.sh          # skip image rebuild
#   COUNT=1 ./scripts/bench-versions.sh            # quick single-pass
#   PARALLEL=2 ./scripts/bench-versions.sh         # 2 bench containers at once
#
# Env vars:
#   COUNT=N       go test -count=N (default 5)
#   REBUILD=0|1   skip docker build if image exists (default 1 = always rebuild)
#   PARALLEL=N    benchmark run concurrency (default 1 = sequential for accuracy)
#
# Requirements: Docker, Go (for benchfmt)
set -euo pipefail

if [ $# -eq 0 ]; then
  VERSIONS=(1.22 1.24)
else
  VERSIONS=("$@")
fi

COUNT=${COUNT:-5}
REBUILD=${REBUILD:-1}
MAX_PARALLEL=${PARALLEL:-1}

RESULTS_DIR="benchmarks/results"
mkdir -p "$RESULTS_DIR"

# Pre-compile benchfmt once to avoid concurrent go run cache conflicts.
BENCHFMT_BIN=$(mktemp -t benchfmt-XXXXXX)
trap 'rm -f "$BENCHFMT_BIN"' EXIT
go build -o "$BENCHFMT_BIN" ./scripts/benchfmt/main.go

# ── Phase 1: build all images in parallel (I/O-bound, no CPU competition) ──
echo "=== Phase 1: building images ===" >&2
build_pids=()
for v in "${VERSIONS[@]}"; do
  (
    image="saferefl-bench:go${v}"
    if [ "$REBUILD" = "0" ] && docker image inspect "$image" >/dev/null 2>&1; then
      echo "[go${v}] image exists, skipping build" >&2
    else
      echo "[go${v}] building..." >&2
      docker build -f bench.Dockerfile --build-arg "GO_VERSION=${v}" -t "$image" . -q >&2
      echo "[go${v}] image ready" >&2
    fi
  ) &
  build_pids+=($!)
done

for i in "${!build_pids[@]}"; do
  if ! wait "${build_pids[$i]}"; then
    echo "[go${VERSIONS[$i]}] build FAILED" >&2
    exit 1
  fi
done

# ── Phase 2: run benchmarks (CPU-bound, controlled parallelism) ─────────────
echo "=== Phase 2: running benchmarks (COUNT=${COUNT}, PARALLEL=${MAX_PARALLEL}) ===" >&2

run_one() {
  local v="$1"
  local image="saferefl-bench:go${v}"
  echo "[go${v}] benchmarking..." >&2
  docker run --rm -e "BENCH_COUNT=${COUNT}" "$image" \
    | "$BENCHFMT_BIN" \
        -output "${RESULTS_DIR}/go${v}.md" \
        -go-version "${v}"
  docker image rm "$image" --force >/dev/null 2>&1 || true
  echo "[go${v}] done -> ${RESULTS_DIR}/go${v}.md" >&2
}

failed=0
n=${#VERSIONS[@]}
batch_start=0

while [ "$batch_start" -lt "$n" ]; do
  batch_end=$((batch_start + MAX_PARALLEL))
  [ "$batch_end" -gt "$n" ] && batch_end=$n

  batch_pids=()
  i=$batch_start
  while [ "$i" -lt "$batch_end" ]; do
    run_one "${VERSIONS[$i]}" &
    batch_pids+=($!)
    i=$((i + 1))
  done

  j=0
  for pid in "${batch_pids[@]}"; do
    vi=$((batch_start + j))
    if ! wait "$pid"; then
      echo "[go${VERSIONS[$vi]}] FAILED" >&2
      failed=1
    fi
    j=$((j + 1))
  done

  batch_start=$batch_end
done

[ "$failed" -eq 0 ] || { echo "One or more versions failed." >&2; exit 1; }
echo "All done. Results in ${RESULTS_DIR}/" >&2
