#!/usr/bin/env bash
# Prints what CI tests for the commits between <base> and HEAD: `go test`
# targets of the root module (./mysql/...) and nested modules with their own
# go.mod (googlesql/spanner-oracle), space-separated. Nothing printed means
# nothing to test.
#
# The engines are independent Go trees, so a change tests its top-level directory
# and every top-level directory whose packages, or their tests, import it: a change
# to mysql tests mysql, a change to review or metadata tests the engines that
# import them. Besides that:
#
#   - go.mod, go.sum and .github/ test everything.
#   - docs/migration/<engine>/ counts as <engine>, Markdown included: googlesql
#     reads its official corpus from .md files there.
#   - Other Markdown tests nothing, wherever it is, and so do root files, the rest
#     of docs/ and the scripts that are not Go.
#   - proto/ counts as metadata, which is generated from it.
#   - A nested module lives in the engine that execs it, so it tests that module
#     and, as any subdirectory would, its engine: googlesql/spanner-oracle tests
#     googlesql. A nested module that imports an affected engine is tested too.
#
# Usage: scripts/affected-packages.sh origin/main [HEAD]
set -euo pipefail

base=${1:?usage: affected-packages.sh <base-ref> [<head-ref>]}
head=${2:-HEAD}
cd "$(git rev-parse --show-toplevel)"

# Tracked files only: a worktree or build tree under an ignored directory
# (.claude/worktrees/) holds whole copies of the repo, go.mod files included.
modules=$(git ls-files '*/go.mod' | xargs -r dirname | sort -u)
dirs=() mods=()
while IFS= read -r f; do
	case "$f" in
	go.mod | go.sum | .github/*)
		echo ./... $modules
		exit 0
		;;
	docs/migration/*)
		e=${f#docs/migration/}
		dirs+=("${e%%/*}")
		;;
	*.md) ;;
	proto/*) dirs+=(metadata) ;;
	*/*)
		dirs+=("${f%%/*}")
		for m in $modules; do
			case "$f" in "$m"/*) mods+=("$m") ;; esac
		done
		;;
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
	# the root module. A module inside the engine it serves shares the root
	# module's import prefix, so its own packages are dropped first: they would
	# otherwise read as an import of their own engine and test it on every change.
	for m in $modules; do
		self=$(cd "$m" && go list -m)
		deps=$(cd "$m" && go list -deps -test ./... 2>/dev/null |
			grep -v "^${self}" |
			sed -n "s#^${mod}\([^/]*\).*#\1#p" | sort -u)
		for d in $affected; do
			grep -qx "$d" <<< "$deps" && mods+=("$m") && break
		done
	done
fi

{
	for d in $affected; do echo "./$d/..."; done
	for m in ${mods[@]+"${mods[@]}"}; do echo "$m"; done
} | sort -u | xargs
