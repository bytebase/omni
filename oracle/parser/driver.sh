#!/bin/bash
# Driver script for autonomous recursive descent parser implementation.
# Launches Claude Code repeatedly until all batches in PROGRESS.json are done.
#
# Usage: ./driver.sh [--max-iterations N]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OMNI_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
PROGRESS_FILE="$SCRIPT_DIR/PROGRESS.json"
SKILL_FILE="$SCRIPT_DIR/SKILL.md"
LOG_DIR="$SCRIPT_DIR/logs"

MAX_ITERATIONS=${1:-100}
ITERATION=0

mkdir -p "$LOG_DIR"

echo "=== Oracle PL/SQL Recursive Descent Parser Builder ==="
echo "Working directory: $OMNI_DIR"
echo "Max iterations: $MAX_ITERATIONS"
echo ""

while [ $ITERATION -lt $MAX_ITERATIONS ]; do
    ITERATION=$((ITERATION + 1))
    TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
    LOG_FILE="$LOG_DIR/iteration_${ITERATION}_${TIMESTAMP}.log"

    # Check if all batches are done
    EXIT_CODE=0
    PENDING=$(python3 -c "
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
" 2>&1) || EXIT_CODE=$?

    echo "[$ITERATION/$MAX_ITERATIONS] Status: $PENDING"

    if [ "$EXIT_CODE" -eq 42 ]; then
        echo ""
        echo "=== ALL BATCHES COMPLETE ==="
        echo "Running final verification..."
        cd "$OMNI_DIR" && go test ./oracle/... 2>&1 | tee "$LOG_DIR/final_verification.log"
        exit 0
    fi

    echo "Starting iteration $ITERATION..."

    # Build the prompt from SKILL.md
    PROMPT=$(cat "$SKILL_FILE")

    # Run Claude Code with the skill prompt
    cd "$OMNI_DIR" && claude -p "$PROMPT" \
        --dangerously-skip-permissions \
        2>&1 | tee "$LOG_FILE"

    echo ""
    echo "Iteration $ITERATION complete. Log: $LOG_FILE"
    echo "---"

    # Brief pause between iterations
    sleep 2
done

echo "Reached max iterations ($MAX_ITERATIONS). Check PROGRESS.json for status."
exit 1
