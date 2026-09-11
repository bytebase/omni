#!/bin/bash
# Unified driver script for parser pipeline with Build → Audit loop.
#
# Usage: ./driver.sh <engine> [--max-iterations N] [--max-audit-rounds N]
#
# The pipeline runs in a loop:
#   Build phase: run SKILL.md until all batches are done
#   Audit phase: run audit to find gaps, generate new batches
#   Repeat until audit finds no gaps
#
# Examples:
#   ./driver.sh pg
#   ./driver.sh tsql --max-iterations 100 --max-audit-rounds 5
#   ./driver.sh oracle --audit-only   # skip build, just audit

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OMNI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Parse arguments
ENGINE=""
MAX_ITERATIONS=100
MAX_AUDIT_ROUNDS=10
AUDIT_ONLY=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --max-iterations) MAX_ITERATIONS="$2"; shift 2 ;;
        --max-audit-rounds) MAX_AUDIT_ROUNDS="$2"; shift 2 ;;
        --audit-only) AUDIT_ONLY=true; shift ;;
        *) ENGINE="$1"; shift ;;
    esac
done

if [ -z "$ENGINE" ]; then
    echo "Usage: $0 <engine> [--max-iterations N] [--max-audit-rounds N] [--audit-only]"
    echo "Engines: pg, mysql, tsql, oracle"
    exit 1
fi

PARSER_DIR="$OMNI_DIR/$ENGINE/parser"
PROGRESS_FILE="$PARSER_DIR/PROGRESS.json"
SKILL_FILE="$PARSER_DIR/SKILL.md"
AUDIT_SKILL="$SCRIPT_DIR/prompts/audit-skill.md"
LOG_DIR="$PARSER_DIR/logs"

if [ ! -f "$PROGRESS_FILE" ]; then
    echo "Error: $PROGRESS_FILE not found"
    exit 1
fi

mkdir -p "$LOG_DIR"

echo "=== Parser Pipeline: $ENGINE ==="
echo "Working directory: $OMNI_DIR"
echo "Max build iterations: $MAX_ITERATIONS"
echo "Max audit rounds: $MAX_AUDIT_ROUNDS"
echo ""

# Get current audit round from PROGRESS.json
get_audit_round() {
    python3 -c "
import json
with open('$PROGRESS_FILE') as f:
    data = json.load(f)
print(data.get('audit_round', 0))
"
}

# Check if all batches are done
all_done() {
    python3 -c "
import json, sys
with open('$PROGRESS_FILE') as f:
    data = json.load(f)
pending = [b for b in data['batches'] if b['status'] == 'pending']
failed = [b for b in data['batches'] if b['status'] == 'failed']
in_prog = [b for b in data['batches'] if b['status'] == 'in_progress']
done = [b for b in data['batches'] if b['status'] == 'done']
actionable = len(pending) + len(failed) + len(in_prog)
print(f'{len(done)} done, {len(pending)} pending, {len(in_prog)} in_progress, {len(failed)} failed')
if actionable == 0:
    sys.exit(42)
" 2>&1
}

# Run build phase
run_build() {
    local iteration=0
    while [ $iteration -lt $MAX_ITERATIONS ]; do
        iteration=$((iteration + 1))
        local timestamp=$(date '+%Y%m%d_%H%M%S')
        local log_file="$LOG_DIR/build_${iteration}_${timestamp}.log"

        local exit_code=0
        local status_line
        status_line=$(all_done) || exit_code=$?

        echo "  [Build $iteration/$MAX_ITERATIONS] $status_line"

        if [ "$exit_code" -eq 42 ]; then
            echo "  All batches complete."
            return 0
        fi

        cd "$OMNI_DIR" && claude -p "$(cat "$SKILL_FILE")" \
            --dangerously-skip-permissions \
            2>&1 | tee "$log_file"

        sleep 2
    done
    echo "  Reached max build iterations."
    return 1
}

# Run audit phase
run_audit() {
    local audit_round=$(get_audit_round)
    local new_round=$((audit_round + 1))
    local timestamp=$(date '+%Y%m%d_%H%M%S')
    local log_file="$LOG_DIR/audit_${new_round}_${timestamp}.log"

    echo "  [Audit round $new_round] Starting..."

    # Generate audit prompt with engine substituted
    local prompt
    prompt=$(sed "s/{{ENGINE}}/$ENGINE/g" "$AUDIT_SKILL")

    local output
    output=$(cd "$OMNI_DIR" && claude -p "$prompt" \
        --dangerously-skip-permissions \
        2>&1 | tee "$log_file")

    # Check if audit found no gaps
    if echo "$output" | grep -q "AUDIT_COMPLETE_NO_GAPS"; then
        echo "  Audit found no gaps. Parser is complete!"
        return 42
    fi

    # Extract report if present
    if echo "$output" | grep -q "AUDIT_REPORT_START"; then
        echo "$output" | sed -n '/AUDIT_REPORT_START/,/AUDIT_REPORT_END/p'
    fi

    echo "  Audit round $new_round complete. New batches may have been added."
    return 0
}

# Main loop: Build → Audit → Build → ...
AUDIT_ROUND=0

while [ $AUDIT_ROUND -lt $MAX_AUDIT_ROUNDS ]; do
    CURRENT_AUDIT=$(get_audit_round)
    echo ""
    echo "========================================="
    echo "  Phase: BUILD (audit round $CURRENT_AUDIT)"
    echo "========================================="

    if [ "$AUDIT_ONLY" = false ]; then
        run_build || true
    fi

    # Verify build
    echo ""
    echo "Running verification..."
    cd "$OMNI_DIR" && go build ./$ENGINE/... 2>&1
    cd "$OMNI_DIR" && go test ./$ENGINE/... 2>&1 | tee "$LOG_DIR/verify_$(date '+%Y%m%d_%H%M%S').log"

    echo ""
    echo "========================================="
    echo "  Phase: AUDIT (round $((CURRENT_AUDIT + 1)))"
    echo "========================================="

    audit_exit=0
    run_audit || audit_exit=$?

    if [ "$audit_exit" -eq 42 ]; then
        echo ""
        echo "========================================="
        echo "  PIPELINE COMPLETE: $ENGINE"
        echo "  All syntax covered. No gaps found."
        echo "========================================="
        exit 0
    fi

    AUDIT_ROUND=$((AUDIT_ROUND + 1))
done

echo ""
echo "Reached max audit rounds ($MAX_AUDIT_ROUNDS)."
echo "Run again to continue, or check PROGRESS.json."
exit 1
