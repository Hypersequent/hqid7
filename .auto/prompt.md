# Autoresearch: hqid7 performance

## Objective
Improve runtime performance and allocation behavior of the hqid7 Go library and CLI without changing any externally observable output, encoding format, decoding behavior, UUIDv7 bit layout, errors, or public API behavior. Do not overfit to benchmark constants; changes must be generally correct for all valid hqid7 values.

## Metrics
- **Primary**: geomean_ns_op (ns/op, lower is better) — geometric mean of median ns/op across all repository Go benchmarks.
- **Secondary**: per-benchmark median ns/op and B/op/allocs where reported.

## How to Run
`./.auto/measure.sh` — runs `go test -run '^$' -bench . -benchmem ./... -count=5`, computes per-benchmark medians, and emits `METRIC name=value` lines.

## Files in Scope
- `encode.go` — Base58 encode/decode helpers and public string conversion API.
- `generate.go` — UUIDv7 generation from current/supplied time.
- `cmd/hqid7/tool.go` — CLI command behavior and parse/generate formatting.
- Tests/benchmarks may be read but should not be weakened to improve metrics.

## Off Limits
- Do not change output strings, separator position, Base58 alphabet, UUID layout, CLI text, or documented behavior.
- Do not remove or weaken tests/benchmarks/checks.
- Do not add new dependencies unless clearly justified by a large general improvement.
- Do not introduce benchmark-specific special cases or constants.

## Constraints
- Correctness checks must pass: `go test ./...`.
- Preserve public API compatibility.
- Prefer simple, maintainable optimizations. Primary metric is king, but catastrophic allocation or correctness regressions are not acceptable.

## What's Been Tried
- Session initialized. Baseline pending.
