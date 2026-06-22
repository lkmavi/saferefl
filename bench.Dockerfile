# Benchmark runner image. Build once, reuse across runs.
# Build: docker build -f bench.Dockerfile --build-arg GO_VERSION=1.22 -t saferefl-bench:go1.22 .
# Run:   docker run --rm saferefl-bench:go1.22
# Quick: docker run --rm -e BENCH_COUNT=1 saferefl-bench:go1.22

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm

WORKDIR /workspace

# Copy full source (benchmarks/go.mod has a replace directive pointing to ../)
COPY . .

# Download dependencies for both modules
RUN go mod download && cd benchmarks && go mod download

# Pre-compile all test packages so the benchmark run starts immediately.
# (-run and -bench both match nothing → compile only, no tests execute)
RUN cd benchmarks && go test -run='^$' -bench='^$' ./... >/dev/null

WORKDIR /workspace/benchmarks

ENV BENCH_COUNT=5
# Shell form so $BENCH_COUNT expands; override at runtime with -e BENCH_COUNT=1
CMD go test ./... -bench=. -benchmem -count=$BENCH_COUNT -run='^$' -v
