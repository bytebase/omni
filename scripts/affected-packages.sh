#!/usr/bin/env bash
# Prints what CI tests for the commits between <base> and HEAD: `go test`
# targets of the root module (./mysql/...) and nested modules under harness/
# (harness/conformance), space-separated. Nothing printed means nothing to test.
#
# The engines are independent Go trees, so a change tests its top-level directory
# and every top-level directory whose packages, or their tests, import it: a change
# to mysql tests mysql, a change to review or metadata tests the engines that
# import them. Besides that:
#
#   - go.mod, go.sum and .github/ test everything.
#   - Markdown tests nothing, wherever it is, and so do root files, docs/ and the
#     scripts that are not Go.
#   - proto/ counts as metadata, which is generated from it.
#   - docs/migration/<engine>/ counts as <engine>: googlesql reads a corpus there.
#   - harness/<module> tests that module, and the engine that execs it
#     (mssql-scriptdom: mssql, googlesql-spanner: googlesql). A nested module that
#     imports an affected engine is tested too.
#
# Usage: scripts/affected-packages.sh origin/main [HEAD]
set -euo pipefail

base=${1:?usage: affected-packages.sh <base-ref> [<head-ref>]}
head=${2:-HEAD}
cd "$(git rev-parse --show-toplevel)"

modules=$(find harness -mindepth 2 -maxdepth 2 -name go.mod -printf '%h\n' | sort)
dirs=() mods=()
while IFS= read -r f; do
	case "$f" in
	go.mod | go.sum | .github/*)
		echo ./... $modules
		exit 0
		;;
	*.md) ;;
	proto/*) dirs+=(metadata) ;;
	docs/migration/*)
		e=${f#docs/migration/}
		dirs+=("${e%%/*}")
		;;
	harness/*)
		m=${f#harness/}
		m=${m%%/*}
		[ ! -f "harness/$m/go.mod" ] || mods+=("harness/$m")
		case "$m" in
		mssql-scriptdom) dirs+=(mssql) ;;
		googlesql-spanner) dirs+=(googlesql) ;;
		esac
		;;
	*/*) dirs+=("${f%%/*}") ;;
	esac
done < <(git diff --name-only "$base...$head")

# Keep the directories that hold Go: docs/ and a deleted engine hold none.
changed=$(for d in ${dirs[@]+"${dirs[@]}"}; do
	[ -z "$(find "$d" -name '*.go' -print -quit 2>/dev/null)" ] || echo "$d"
done | sort -u | xargs)

affected=
if [ -n "$changed" ]; then
	mod=$(go list -m)/
	affected=$(go list -f '{{.ImportPath}} {{join .Deps " "}} {{join .TestImports " "}} {{join .XTestImports " "}}' ./... |
		awk -v mod="$mod" -v changed="$changed" '
			function top(p) { sub("^" mod, "", p); sub("/.*", "", p); return p }
			BEGIN { n = split(changed, c, " "); for (i = 1; i <= n; i++) want[c[i]] = 1 }
			{ for (i = 1; i <= NF; i++) if (index($i, mod) == 1 && top($i) in want) { out[top($1)] = 1; break } }
			END { for (d in out) print d }' | sort | xargs)
	# A nested module that imports an affected engine through its replace of
	# the root module.
	for m in $modules; do
		deps=$(cd "$m" && go list -deps -test ./... 2>/dev/null | sed -n "s#^${mod}\([^/]*\).*#\1#p" | sort -u)
		for d in $affected; do
			grep -qx "$d" <<< "$deps" && mods+=("$m") && break
		done
	done
fi

{
	for d in $affected; do echo "./$d/..."; done
	for m in ${mods[@]+"${mods[@]}"}; do echo "$m"; done
} | sort -u | xargs
