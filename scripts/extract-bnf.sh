#!/bin/bash
# Extract BNF from documentation using Claude.
#
# Usage: ./scripts/extract-bnf.sh <engine> [--max-iterations N]
#
# Engines: oracle, mysql, mssql
# Each engine needs:
#   - {engine}/parser/{ENGINE}_BNF_CATALOG.json (with url_ok entries)
#   - scripts/prompts/{engine}-extract-bnf-skill.md

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OMNI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

ENGINE=""
MAX_ITERATIONS=20

while [[ $# -gt 0 ]]; do
    case "$1" in
        --max-iterations) MAX_ITERATIONS="$2"; shift 2 ;;
        *) ENGINE="$1"; shift ;;
    esac
done

if [ -z "$ENGINE" ]; then
    echo "Usage: $0 <engine> [--max-iterations N]"
    echo "Engines: oracle, mysql, mssql"
    exit 1
fi

ENGINE_UPPER=$(echo "$ENGINE" | tr '[:lower:]' '[:upper:]')
CATALOG="$OMNI_DIR/$ENGINE/parser/${ENGINE_UPPER}_BNF_CATALOG.json"
BNF_DIR="$OMNI_DIR/$ENGINE/parser/bnf"
LOG_DIR="$OMNI_DIR/$ENGINE/parser/logs"
SKILL_FILE="$SCRIPT_DIR/prompts/${ENGINE}-extract-bnf-skill.md"

if [ ! -f "$CATALOG" ]; then
    echo "Error: $CATALOG not found"
    exit 1
fi
if [ ! -f "$SKILL_FILE" ]; then
    echo "Error: $SKILL_FILE not found"
    exit 1
fi

mkdir -p "$BNF_DIR" "$LOG_DIR"

# Count remaining
count_remaining() {
    python3 -c "
import json
with open('$CATALOG') as f:
    data = json.load(f)
remaining = sum(1 for s in data['statements'] if s['status'] == 'url_ok')
done = sum(1 for s in data['statements'] if s.get('status') in ('bnf_done',))
failed = sum(1 for s in data['statements'] if s.get('status') == 'bnf_failed')
no_doc = sum(1 for s in data['statements'] if s.get('status') == 'no_doc')
total = len(data['statements'])
print(f'{done} done, {remaining} remaining, {failed} failed, {no_doc} no_doc (total: {total})')
if remaining == 0:
    exit(42)
"
}

# Track child PID for signal forwarding
CHILD_PID=""
cleanup() {
    if [ -n "$CHILD_PID" ] && kill -0 "$CHILD_PID" 2>/dev/null; then
        kill -TERM "$CHILD_PID" 2>/dev/null
        wait "$CHILD_PID" 2>/dev/null
    fi
    echo ""
    echo "BNF extraction interrupted."
    exit 130
}
trap cleanup INT TERM

echo "=== BNF Extraction: $ENGINE ==="
echo "Catalog: $CATALOG"
echo "Output: $BNF_DIR"
echo "Skill: $SKILL_FILE"
echo "Max iterations: $MAX_ITERATIONS"
echo ""

for ((i=1; i<=MAX_ITERATIONS; i++)); do
    local_exit=0
    status_line=$(count_remaining) || local_exit=$?

    echo "  [Iteration $i/$MAX_ITERATIONS] $status_line"

    if [ "$local_exit" -eq 42 ]; then
        echo ""
        echo "All BNF extracted for $ENGINE!"
        exit 0
    fi

    timestamp=$(date '+%Y%m%d_%H%M%S')
    log_file="$LOG_DIR/bnf_${i}_${timestamp}.log"

    cd "$OMNI_DIR" && CLAUDECODE="" claude -p "$(cat "$SKILL_FILE")" \
        --dangerously-skip-permissions \
        --output-format stream-json --verbose \
        2>/dev/null | python3 "$SCRIPT_DIR/stream-filter.py" "$log_file" &
    CHILD_PID=$!
    wait $CHILD_PID
    CHILD_PID=""

    sleep 1
done

echo "Reached max iterations for $ENGINE."
