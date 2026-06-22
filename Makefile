.PHONY: bench bench-realistic bench-stat bench-local bench-docker test lint

bench: bench-realistic

bench-realistic:
	cd benchmarks && go test ./realistic/... -bench=. -benchmem -count=5

# Capture realistic results to JSON for benchstat comparison.
# Usage: make bench-json OUT=bench.json
bench-json:
	cd benchmarks && go test ./realistic/... -bench=. -benchmem -count=5 -json > ../$(OUT)

# Compare two saved JSON runs with benchstat.
# Usage: make bench-stat OLD=bench-main.json NEW=bench-pr.json
bench-stat:
	go run golang.org/x/perf/cmd/benchstat@latest $(OLD) $(NEW)

# Run benchmarks on local hardware, regenerate benchmarks/RESULTS.md.
# Usage: make bench-local           — results file only
#        make bench-local readme=1  — results + inject into README.md
#        COUNT=1 make bench-local   — quick single-pass
bench-local:
	@if [ "$(readme)" = "1" ]; then \
		./scripts/bench.sh -readme; \
	else \
		./scripts/bench.sh; \
	fi

# Cross-version benchmarks via Docker (builds first, then runs sequentially).
# Usage: make bench-docker
#        REBUILD=0 make bench-docker      — skip rebuild if images exist
#        COUNT=1 make bench-docker        — quick single-pass
#        PARALLEL=2 make bench-docker     — run 2 containers at once
bench-docker:
	./scripts/bench-versions.sh

# Single-version Docker benchmark.
# Usage: make bench-docker-1.24
#        COUNT=1 make bench-docker-1.24
bench-docker-%:
	docker build -f bench.Dockerfile --build-arg GO_VERSION=$* -t saferefl-bench:go$* .
	docker run --rm -e BENCH_COUNT=${COUNT:-5} saferefl-bench:go$* \
	  | go run ./scripts/benchfmt/main.go \
	      -output benchmarks/results/go$*.md \
	      -go-version $*
	docker image rm saferefl-bench:go$* --force

test:
	go test ./... -race -count=1

lint:
	golangci-lint run
