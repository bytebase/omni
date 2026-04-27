#!/bin/bash
# Local MySQL test driver.
#
# Usage:
#   scripts/test-mysql.sh quick
#   scripts/test-mysql.sh full
#   scripts/test-mysql.sh list-shards
#   scripts/test-mysql.sh container-shards [shard-name...]

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

declare -a SHARD_NAMES=(
    "smoke-and-misc"
    "section-1-core"
    "section-1-indexes-options"
    "section-2-ddl"
    "sections-3-through-6"
    "deparse-container"
    "implicit-scenarios-c"
    "implicit-scenarios-extra"
)

shard_regex() {
    case "$1" in
        smoke-and-misc)
            printf '%s\n' '^(TestDDLWorkflow_Container|TestShowCreateTable_ContainerComparison|TestContainerSmoke|TestContainer_ReservedKeywordAcceptance|TestSpotCheck_CatalogVerification|TestMySQL_DeparseRules)$'
            ;;
        section-1-core)
            printf '%s\n' '^TestContainer_Section_1_([1-9]|10|11|12)_'
            ;;
        section-1-indexes-options)
            printf '%s\n' '^TestContainer_Section_1_(13|14|15|16|17|18|19|20)_'
            ;;
        section-2-ddl)
            printf '%s\n' '^TestContainer_Section_2_'
            ;;
        sections-3-through-6)
            printf '%s\n' '^TestContainer_Section_[3-6]_'
            ;;
        deparse-container)
            printf '%s\n' '^(TestDeparse_.*_Container|TestDeparseContainer_)'
            ;;
        implicit-scenarios-c)
            printf '%s\n' '^TestScenario_C'
            ;;
        implicit-scenarios-extra)
            printf '%s\n' '^TestScenario_(AX|PS)$'
            ;;
        *)
            return 1
            ;;
    esac
}

usage() {
    sed -n '2,9p' "$0" | sed 's/^# \{0,1\}//'
}

run_go_test() {
    echo "+ go test $*"
    (cd "$ROOT_DIR" && go test "$@")
}

list_shards() {
    for name in "${SHARD_NAMES[@]}"; do
        printf '%s\t%s\n' "$name" "$(shard_regex "$name")"
    done
}

selected_shards() {
    if [[ $# -eq 0 ]]; then
        printf '%s\n' "${SHARD_NAMES[@]}"
        return
    fi

    local name
    for name in "$@"; do
        if ! shard_regex "$name" >/dev/null; then
            echo "unknown shard: $name" >&2
            echo "known shards:" >&2
            list_shards >&2
            exit 2
        fi
        printf '%s\n' "$name"
    done
}

run_container_shards() {
    local failed=0
    local name
    local regex
    local -a failures=()

    while IFS= read -r name; do
        echo
        echo "=== mysql catalog container shard: $name ==="
        regex="$(shard_regex "$name")"
        set +e
        run_go_test -timeout 15m ./mysql/catalog/ -run "$regex" -count=1 -v
        local status=$?
        set -e
        if [[ $status -ne 0 ]]; then
            failures+=("$name")
            failed=1
        fi
    done < <(selected_shards "$@")

    if [[ $failed -ne 0 ]]; then
        echo
        echo "failed mysql container shards:"
        printf '  %s\n' "${failures[@]}"
        return 1
    fi
}

mode="${1:-}"
if [[ -z "$mode" ]]; then
    usage
    exit 2
fi
shift

case "$mode" in
    quick)
        run_go_test -short ./mysql/... -count=1
        ;;
    full)
        run_go_test ./mysql/... -count=1 -timeout=20m
        ;;
    list-shards)
        list_shards
        ;;
    container-shards)
        run_container_shards "$@"
        ;;
    *)
        echo "unknown mode: $mode" >&2
        usage >&2
        exit 2
        ;;
esac
