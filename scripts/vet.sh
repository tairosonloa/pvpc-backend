#!/bin/sh

# By using the go vet command, it is possible to check for
# errors that are not reported by the compiler.

set -e -u

output=`go vet ./...`
test "$output" = "" && echo "Everything OK" && exit 0
echo "The following issues were found:"
echo "$output"
exit 1