#!/bin/sh

# Staticcheck is a state of the art linter for the Go programming language.
# Using static analysis, it finds bugs and performance issues, offers
# simplifications, and enforces style rules.

set -e -u

go install honnef.co/go/tools/cmd/staticcheck@latest
output=`staticcheck ./...`
test "$output" = "" && echo "Everything OK" && exit 0
echo "The following issues were found:"
echo "$output"
exit 1