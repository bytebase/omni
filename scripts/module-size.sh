#!/usr/bin/env bash
# Fails when the github.com/bytebase/omni module zip grows past its budget. Go
# refuses to download a module whose files pass 500 MB (524288000 bytes), so a
# module over that cannot be pinned by anyone. The zip holds every tracked file
# outside a directory with its own go.mod (the tapes directories, nested
# modules); the budget leaves room below Go's cap.
#
# Usage: scripts/module-size.sh [<ref>]
set -euo pipefail

ref=${1:-HEAD}
budget=$((400 * 1024 * 1024))
cd "$(git rev-parse --show-toplevel)"

nested=$(git ls-tree -r --name-only "$ref" | grep '/go.mod$' | xargs -r dirname | sort -u)
size=$(git ls-tree -r -l "$ref" | awk -v nested="$nested" '
	BEGIN { n = split(nested, dirs, "\n") }
	{
		path = $5
		for (i = 1; i <= n; i++) {
			if (dirs[i] != "" && index(path, dirs[i] "/") == 1) next
		}
		total += $4
	}
	END { print total + 0 }')

echo "module files: $size bytes (budget $budget, Go's cap 524288000)"
if [ "$size" -gt "$budget" ]; then
	echo "the module is over budget: move large test data into a directory with its own go.mod (as testdata/tapes is)" >&2
	exit 1
fi
