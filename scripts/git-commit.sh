#!/bin/bash
# Safe git commit with lock retry for concurrent pipelines.
# Usage: ./git-commit.sh "commit message" [files...]
#
# Retries up to 10 times with random backoff if git index is locked.

set -euo pipefail

MSG="$1"
shift
FILES=("$@")

REPO_ROOT="$(git rev-parse --show-toplevel)"
MAX_RETRIES=10

for i in $(seq 1 $MAX_RETRIES); do
    # Stage files
    if ! git add "${FILES[@]}" 2>/dev/null; then
        echo "[git-commit] Retry $i/$MAX_RETRIES: git add failed (likely lock), waiting..."
        sleep $(( (RANDOM % 3) + 1 ))
        continue
    fi

    # Check if there's anything to commit
    if git diff --cached --quiet 2>/dev/null; then
        echo "[git-commit] Nothing to commit."
        exit 0
    fi

    # Commit
    if git commit -m "$MSG" 2>/dev/null; then
        echo "[git-commit] Committed: $MSG"
        exit 0
    else
        echo "[git-commit] Retry $i/$MAX_RETRIES: git commit failed, waiting..."
        sleep $(( (RANDOM % 3) + 1 ))
    fi
done

echo "[git-commit] Failed after $MAX_RETRIES retries."
exit 1
