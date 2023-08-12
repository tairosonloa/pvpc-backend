#!/bin/sh

# This script ensures a minimum test coverage is achieved.

set -e -u

go test ./... -coverprofile coverage.out -covermode atomic

perc=`go tool cover -func=coverage.out | tail -n 1 | sed -Ee 's!^[^[:digit:]]+([[:digit:]]+(\.[[:digit:]]+)?)%$!\1!'`
res=`echo "$perc >= 90.0" | bc` # 90% minimum coverage
test "$res" -eq 1 && echo "OK: Coverage of $perc % (threshold requires >= 90.0 %)" && exit 0
echo "Insufficient coverage: $perc" >&2
exit 1