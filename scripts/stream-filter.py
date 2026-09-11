#!/usr/bin/env python3
"""Filter stream-json output from claude -p to show real-time progress.

Reads JSON lines from stdin (claude --output-format stream-json --verbose),
extracts meaningful progress info, and displays it on the terminal.
Also writes the full final text to the log file specified as arg.

Usage:
    claude -p "..." --output-format stream-json --verbose | \
        python3 stream-filter.py <logfile>
"""

import json
import sys
import re
import time

log_file = sys.argv[1] if len(sys.argv) > 1 else None
full_text_parts = []
start_time = time.time()
last_progress_line = None
current_tool = None

# ANSI colors
DIM = "\033[2m"
CYAN = "\033[36m"
GREEN = "\033[32m"
YELLOW = "\033[33m"
RED = "\033[31m"
BOLD = "\033[1m"
RESET = "\033[0m"


def elapsed():
    secs = int(time.time() - start_time)
    m, s = divmod(secs, 60)
    return f"{m}:{s:02d}"


def show(msg, color=""):
    print(f"    {DIM}{elapsed()}{RESET} {color}{msg}{RESET}", flush=True)


def check_batch_markers(text):
    """Extract and display [BATCH N] progress markers."""
    global last_progress_line
    for m in re.finditer(r'\[BATCH \d+\].*', text):
        line = m.group(0).strip()
        if line == last_progress_line:
            continue
        last_progress_line = line
        if "FAIL" in line:
            show(line, RED)
        elif "DONE" in line:
            show(line, GREEN)
        elif "STARTED" in line:
            show(line, BOLD + CYAN)
        elif "RETRY" in line:
            show(line, YELLOW)
        else:
            show(line, CYAN)


for raw_line in sys.stdin:
    raw_line = raw_line.strip()
    if not raw_line:
        continue
    try:
        event = json.loads(raw_line)
    except json.JSONDecodeError:
        continue

    etype = event.get("type", "")

    if etype == "system" and event.get("subtype") == "init":
        show("Session started", DIM)
        continue

    if etype == "assistant":
        msg = event.get("message", {})
        for part in msg.get("content", []):
            if part.get("type") == "text":
                text = part["text"]
                full_text_parts.append(text)
                check_batch_markers(text)
            elif part.get("type") == "tool_use":
                tool_name = part.get("name", "?")
                inp = part.get("input", {})
                if tool_name == "Bash":
                    cmd = inp.get("command", "")
                    if any(k in cmd for k in ["go build", "go test", "git commit", "git-commit"]):
                        show(f"$ {cmd[:120]}", DIM)
                elif tool_name == "Write":
                    fp = inp.get("file_path", "")
                    if fp:
                        show(f"Write: {fp.split('/')[-1]}", DIM)
                elif tool_name == "Edit":
                    fp = inp.get("file_path", "")
                    if fp:
                        show(f"Edit: {fp.split('/')[-1]}", DIM)
                elif tool_name == "Read":
                    fp = inp.get("file_path", "")
                    if fp:
                        show(f"Read: {fp.split('/')[-1]}", DIM)
                elif tool_name == "WebFetch":
                    url = inp.get("url", "")
                    if url:
                        show(f"Fetch: {url[:80]}", DIM)
        continue

    if etype == "result":
        sub = event.get("subtype", "")
        dur = event.get("duration_ms", 0)
        cost = event.get("total_cost_usd", 0)
        dur_s = dur / 1000
        # Use result text as the canonical output (replaces assistant text parts)
        result_text = event.get("result", "")
        if result_text:
            full_text_parts = [result_text]
            check_batch_markers(result_text)
        if sub == "success":
            show(f"Done ({dur_s:.0f}s, ${cost:.2f})", GREEN)
        elif sub == "error":
            err = event.get("error", "unknown")
            show(f"Error: {err} ({dur_s:.0f}s)", RED)
        continue

# Write full output to log
if log_file:
    with open(log_file, "w") as f:
        f.write("\n".join(full_text_parts))
