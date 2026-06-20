#!/usr/bin/env bash
set -euo pipefail

# Cap noisy failing test output so pi's exec stdout accumulator cannot grow
# without bound (for example, a broken encoder can make TestEncodeDecode emit
# gigabytes of repeated t.Errorf lines before exiting).
go test ./... 2>&1 | head -c "${CHECK_OUTPUT_LIMIT_BYTES:-200000}"
