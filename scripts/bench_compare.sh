#!/usr/bin/env bash
set -euo pipefail

# Simple benchmark runner and comparator against baseline.
#
# Usage:
#   scripts/bench_compare.sh [output_file]
#
# Env vars:
#   BENCHTIME   - benchtime for go test (default: 1s)
#   BENCH_PKGS  - packages to bench (default: ./bench/...)
#   GOFLAGS     - extra flags forwarded to `go test`
#
# Output:
#   Writes current results to benchmarks/CURRENT_BENCH.txt (or provided output_file)
#   and diffs against benchmarks/BASELINE_BENCH.txt if present.

repo_root_dir() {
  git rev-parse --show-toplevel 2>/dev/null || pwd
}

ROOT="$(repo_root_dir)"
BENCH_DIR="$ROOT/benchmarks"
mkdir -p "$BENCH_DIR"

BASELINE="$BENCH_DIR/BASELINE_BENCH.txt"
CURRENT="${1:-$BENCH_DIR/CURRENT_BENCH.txt}"

BENCHTIME="${BENCHTIME:-1s}"
BENCH_PKGS="${BENCH_PKGS:-./bench/...}"

echo "Running benchmarks: pkgs=$BENCH_PKGS, benchtime=$BENCHTIME" >&2
{
  echo "goos: $(go env GOOS)"
  echo "goarch: $(go env GOARCH)"
} >/dev/null

# Run benchmarks and save output
go test -bench=. -benchmem -benchtime="$BENCHTIME" ${GOFLAGS:-} $BENCH_PKGS | tee "$CURRENT"

echo
if [[ -f "$BASELINE" ]]; then
  echo "Comparing to baseline: $BASELINE" >&2
  if command -v colordiff >/dev/null 2>&1; then
    colordiff -u "$BASELINE" "$CURRENT" || true
  else
    diff -u "$BASELINE" "$CURRENT" || true
  fi
else
  echo "Baseline not found at $BASELINE" >&2
  echo "To set baseline: cp '$CURRENT' '$BASELINE'" >&2
fi

echo "Current results written to: $CURRENT" >&2
