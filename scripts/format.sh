#!/bin/sh

# This script ensures that code is correctly formated.
# It uses goimport with the -local flag to detect local packages.

set -e -u

go install golang.org/x/tools/cmd/goimports@latest
output=`goimports -e -l -local go-pvpc ./`
test "$output" = "" && echo "Everything OK" && exit 0
echo "The following files have incorrect format:"
echo "$output"
exit 1