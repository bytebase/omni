#!/usr/bin/env bash
# Prints the `go test` targets affected by the commits between <base> and HEAD.
#
# The engines are independent Go trees (only metadata is shared), so a change
# maps to its top-level directory plus every directory whose packages, or
# their tests, import it. A change outside any Go directory (go.mod, proto,
# workflows, root files) or under scripts/, which holds the CI drivers
# themselves, means ./... . Markdown under docs/ or at the root cannot affect a
# test and is ignored.
#
# Usage: scripts/affected-packages.sh origin/main
set -euo pipefail

base=${1:?usage: affected-packages.sh <base-ref>}
changed=$(git diff --name-only "$base...HEAD" | { grep -vE '^(docs/|[^/]+\.md$)' || true; } |
	awk -F/ '{ print (NF > 1 ? $1 : ".") }' | sort -u | xargs)
[ -n "$changed" ] || exit 0
for d in $changed; do
	[ "$d" != . ] && [ "$d" != scripts ] && [ -n "$(find "$d" -name '*.go' -print -quit)" ] || { echo ./...; exit 0; }
done

mod=$(go list -m)/
go list -f '{{.ImportPath}} {{join .Deps " "}} {{join .TestImports " "}} {{join .XTestImports " "}}' ./... |
	awk -v mod="$mod" -v changed="$changed" '
		function top(p) { sub("^" mod, "", p); sub("/.*", "", p); return p }
		BEGIN { n = split(changed, c, " "); for (i = 1; i <= n; i++) want[c[i]] = 1 }
		{ for (i = 1; i <= NF; i++) if (index($i, mod) == 1 && top($i) in want) { out[top($1)] = 1; break } }
		END { for (d in out) print "./" d "/..." }' | sort | xargs
