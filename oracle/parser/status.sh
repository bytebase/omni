#!/bin/bash
# Live progress monitor for the Oracle PL/SQL recursive descent parser implementation.
# Usage: ./status.sh          (live watch, refreshes every 3s)
#        ./status.sh --once   (print once and exit)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROGRESS_FILE="$SCRIPT_DIR/PROGRESS.json"
WATCH=true
[ "${1:-}" = "--once" ] && WATCH=false

show_status() {
python3 -c "
import json, datetime

with open('$PROGRESS_FILE') as f:
    data = json.load(f)

batches = data['batches']
done = [b for b in batches if b['status'] == 'done']
pending = [b for b in batches if b['status'] == 'pending']
failed = [b for b in batches if b['status'] == 'failed']
in_prog = [b for b in batches if b['status'] == 'in_progress']

total = len(batches)
pct = len(done) / total * 100

# Progress bar
bar_len = 40
filled = int(bar_len * len(done) / total)
bar = '#' * filled + '.' * (bar_len - filled)

now = datetime.datetime.now().strftime('%H:%M:%S')
print(f'[{now}] Progress: [{bar}] {pct:.0f}% ({len(done)}/{total})')
print(f'  Done: {len(done)}  Pending: {len(pending)}  In Progress: {len(in_prog)}  Failed: {len(failed)}')
print()

# Ready to run (dependencies met)
done_ids = {b['id'] for b in done}
ready = [b for b in pending if all(d in done_ids for d in b['depends_on'])]
if in_prog:
    print('Currently working on:')
    for b in in_prog:
        print(f'  >> batch {b[\"id\"]}: {b[\"name\"]} ({b[\"file\"]})')
    print()
if ready:
    print('Next up (dependencies met):')
    for b in ready:
        print(f'  -> batch {b[\"id\"]}: {b[\"name\"]} ({b[\"file\"]})')
    print()

# Show all batches
for b in batches:
    icon = {'done': '[x]', 'pending': '[ ]', 'in_progress': '[>]', 'failed': '[!]'}.get(b['status'], '[?]')
    tests = ', '.join(b['tests']) if b['tests'] else '-'
    err = f' ERROR: {b[\"error\"]}' if 'error' in b else ''
    print(f'  {icon} {b[\"id\"]:2d}. {b[\"name\"]:<20s} {b[\"file\"]:<25s} tests: {tests}{err}')
"
}

if [ "$WATCH" = true ]; then
    trap 'echo ""; exit 0' INT
    while true; do
        clear
        echo "=== Oracle PL/SQL Recursive Descent Parser Progress (Ctrl+C to quit) ==="
        echo ""
        show_status
        sleep 3
    done
else
    show_status
fi
