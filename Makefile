.PHONY: bench bench-spike bench-realistic bench-stat test lint

bench: bench-spike bench-realistic

bench-spike:
	cd benchmarks && go test ./spike/... -bench=. -benchmem -count=5

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

test:
	go test ./... -race -count=1

lint:
	golangci-lint run
