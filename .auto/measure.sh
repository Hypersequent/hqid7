#!/usr/bin/env bash
set -euo pipefail

export GOMAXPROCS=${GOMAXPROCS:-1}
TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

go test -run '^$' -bench . -benchmem ./... -count=5 -benchtime=200ms >"$TMP"

python3 - "$TMP" <<'PY'
import math
import re
import statistics
import sys
from collections import defaultdict

path = sys.argv[1]
# Example: BenchmarkNewString          12345   678.9 ns/op   24 B/op   2 allocs/op
line_re = re.compile(r'^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op(?:\s+([0-9.]+)\s+B/op)?(?:\s+([0-9.]+)\s+allocs/op)?')
vals = defaultdict(lambda: {"ns": [], "b": [], "allocs": []})
with open(path, 'r', encoding='utf-8') as f:
    for line in f:
        m = line_re.match(line.strip())
        if not m:
            continue
        raw, ns, b, allocs = m.groups()
        # Strip CPU suffix and make a stable metric-safe name.
        name = re.sub(r'-\d+$', '', raw)
        name = name.replace('/', '_').replace('-', '_')
        vals[name]["ns"].append(float(ns))
        if b is not None:
            vals[name]["b"].append(float(b))
        if allocs is not None:
            vals[name]["allocs"].append(float(allocs))

if not vals:
    print(open(path, 'r', encoding='utf-8').read())
    raise SystemExit('no benchmark lines parsed')

med_ns = {}
for name in sorted(vals):
    ns = statistics.median(vals[name]["ns"])
    med_ns[name] = ns
    print(f"METRIC {name}_ns_op={ns}")
    if vals[name]["b"]:
        print(f"METRIC {name}_B_op={statistics.median(vals[name]['b'])}")
    if vals[name]["allocs"]:
        print(f"METRIC {name}_allocs_op={statistics.median(vals[name]['allocs'])}")

geo = math.exp(sum(math.log(v) for v in med_ns.values()) / len(med_ns))
print(f"METRIC geomean_ns_op={geo}")
PY
