//go:build ignore

// benchfmt reads raw `go test -bench` output from stdin and writes a Markdown
// benchmark report. It also injects the report into README.md when -readme is set.
//
// Usage:
//
//	go test ./benchmarks/... -bench=. -benchmem -count=5 | \
//	    go run ./scripts/benchfmt/main.go -output benchmarks/RESULTS.md -readme README.md
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	outputPath := flag.String("output", "benchmarks/RESULTS.md", "path for benchmark results Markdown")
	readmePath := flag.String("readme", "", "README.md to inject results into (optional)")
	flag.Parse()

	env, groups := parseBenchOutput(os.Stdin)
	md := renderMarkdown(env, groups)

	if err := os.WriteFile(*outputPath, []byte(md), 0o644); err != nil {
		fatalf("write %s: %v", *outputPath, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", *outputPath)

	if *readmePath != "" {
		if err := injectReadme(*readmePath, md); err != nil {
			fatalf("inject readme: %v", err)
		}
		fmt.Fprintf(os.Stderr, "updated %s\n", *readmePath)
	}
}

// --- Domain types ---

type benchEnv struct {
	date, goos, goarch, cpu string
}

type group struct {
	name     string
	variants []*variant
}

type variant struct {
	name        string
	nsPerOp     float64
	bytesPerOp  float64
	allocsPerOp float64
}

// --- Parser ---

type rawSample struct {
	nsOp, bOp, allocsOp float64
}

type sampleKey struct{ group, variant string }

func parseBenchOutput(r io.Reader) (benchEnv, []group) {
	env := benchEnv{date: time.Now().UTC().Format("2006-01-02 15:04 UTC")}

	samples := map[sampleKey][]rawSample{}
	var keyOrder []sampleKey
	keySet := map[sampleKey]bool{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "goos: "):
			env.goos = strings.TrimPrefix(line, "goos: ")
		case strings.HasPrefix(line, "goarch: "):
			env.goarch = strings.TrimPrefix(line, "goarch: ")
		case strings.HasPrefix(line, "cpu: "):
			env.cpu = strings.TrimPrefix(line, "cpu: ")
		case strings.HasPrefix(line, "Benchmark"):
			g, v, s, ok := parseBenchLine(line)
			if !ok {
				continue
			}
			k := sampleKey{g, v}
			if !keySet[k] {
				keySet[k] = true
				keyOrder = append(keyOrder, k)
			}
			samples[k] = append(samples[k], s)
		}
	}

	// Build ordered groups
	var (
		groupOrder []string
		groupSet   = map[string]bool{}
		groupVars  = map[string][]string{}
		groupVarSet = map[string]map[string]bool{}
	)
	for _, k := range keyOrder {
		if !groupSet[k.group] {
			groupSet[k.group] = true
			groupOrder = append(groupOrder, k.group)
			groupVarSet[k.group] = map[string]bool{}
		}
		if !groupVarSet[k.group][k.variant] {
			groupVarSet[k.group][k.variant] = true
			groupVars[k.group] = append(groupVars[k.group], k.variant)
		}
	}

	var groups []group
	for _, gName := range groupOrder {
		g := group{name: gName}
		for _, vName := range groupVars[gName] {
			ss := samples[sampleKey{gName, vName}]
			g.variants = append(g.variants, &variant{
				name:        vName,
				nsPerOp:     meanOf(ss, func(s rawSample) float64 { return s.nsOp }),
				bytesPerOp:  meanOf(ss, func(s rawSample) float64 { return s.bOp }),
				allocsPerOp: meanOf(ss, func(s rawSample) float64 { return s.allocsOp }),
			})
		}
		groups = append(groups, g)
	}
	return env, groups
}

// parseBenchLine extracts group name, variant name, and metrics from one output line.
// Benchmark names must follow the convention: BenchmarkGroup_Variant-N
func parseBenchLine(line string) (grp, vrt string, s rawSample, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return
	}

	name := fields[0]
	// Strip -N goroutine count suffix
	if i := strings.LastIndex(name, "-"); i > 0 {
		name = name[:i]
	}
	name = strings.TrimPrefix(name, "Benchmark")

	// Split at last underscore → group / variant
	i := strings.LastIndex(name, "_")
	if i < 0 {
		return // skip benchmarks without a variant suffix
	}
	grp = name[:i]
	vrt = name[i+1:]

	// Parse metric pairs: value unit value unit …
	for i := 2; i+1 < len(fields); i += 2 {
		val, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			continue
		}
		switch fields[i+1] {
		case "ns/op":
			s.nsOp = val
		case "B/op":
			s.bOp = val
		case "allocs/op":
			s.allocsOp = val
		}
	}
	ok = true
	return
}

// --- Renderer ---

func renderMarkdown(env benchEnv, groups []group) string {
	var sb strings.Builder

	sb.WriteString("# Benchmark Results\n\n")
	fmt.Fprintf(&sb, "**Generated:** %s  \n", env.date)
	if env.goos != "" {
		fmt.Fprintf(&sb, "**Platform:** %s/%s  \n", env.goos, env.goarch)
	}
	if env.cpu != "" {
		fmt.Fprintf(&sb, "**CPU:** %s  \n", env.cpu)
	}
	sb.WriteString("\n")

	for _, g := range groups {
		fmt.Fprintf(&sb, "## %s\n\n", insertSpaces(g.name))
		sb.WriteString("| Variant | ns/op | B/op | allocs/op | vs first |\n")
		sb.WriteString("|---|---|---|---|---|\n")

		baseline := g.variants[0].nsPerOp
		for i, v := range g.variants {
			ratio := ratioLabel(i, baseline, v.nsPerOp)
			fmt.Fprintf(&sb, "| %s | %s | %.0f | %.0f | %s |\n",
				v.name, fmtNs(v.nsPerOp), v.bytesPerOp, v.allocsPerOp, ratio)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ratioLabel produces "—", "Nx faster", or "Nx slower" relative to baseline.
func ratioLabel(i int, baseline, current float64) string {
	if i == 0 || baseline <= 0 || current <= 0 {
		return "—"
	}
	r := baseline / current
	if r >= 1 {
		return fmt.Sprintf("%.1f× faster", r)
	}
	return fmt.Sprintf("%.1f× slower", 1/r)
}

// --- README injection ---

const (
	startMarker = "<!-- bench:start -->"
	endMarker   = "<!-- bench:end -->"
)

func injectReadme(path, content string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(data)

	start := strings.Index(s, startMarker)
	end := strings.Index(s, endMarker)
	if start < 0 || end < 0 || end <= start {
		return fmt.Errorf("markers %q / %q not found in %s — add them to enable injection",
			startMarker, endMarker, path)
	}

	updated := s[:start+len(startMarker)] + "\n" + content + s[end:]
	return os.WriteFile(path, []byte(updated), 0o644)
}

// --- Helpers ---

func meanOf(ss []rawSample, f func(rawSample) float64) float64 {
	if len(ss) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, s := range ss {
		sum += f(s)
	}
	return sum / float64(len(ss))
}

func fmtNs(ns float64) string {
	switch {
	case math.IsNaN(ns):
		return "—"
	case ns < 1:
		return fmt.Sprintf("%.3f", ns)
	case ns < 10:
		return fmt.Sprintf("%.2f", ns)
	case ns < 100:
		return fmt.Sprintf("%.1f", ns)
	default:
		return fmt.Sprintf("%.0f", ns)
	}
}

// insertSpaces turns "ReadInt64" → "Read Int64", "JSONLike" → "JSON Like".
func insertSpaces(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 {
			prev := runes[i-1]
			isUpper := r >= 'A' && r <= 'Z'
			prevLower := prev >= 'a' && prev <= 'z'
			prevUpper := prev >= 'A' && prev <= 'Z'
			nextLower := i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z'
			// Space before Aa (camelCase) or before Aa in AAAa (acronym end)
			if isUpper && (prevLower || (prevUpper && nextLower)) {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "benchfmt: "+format+"\n", args...)
	os.Exit(1)
}
