# MySQL Implicit-Behavior Scenario Bug Queue

This directory collects bugs discovered while executing the
`mysql-implicit-behavior` starmap. Each of the 25 sections in
`SCENARIOS-mysql-implicit-behavior.md` gets its own file,
`<section-id>.md`, listing every scenario whose dual-assertion test
(real MySQL 8.0 container vs. the omni catalog) diverged.

The queue is **append-only during a batch**: workers add entries as
they hit diffs and keep the test using `t.Error` so the whole section
runs. A follow-up fix pass drains entries back to zero by either
patching the omni catalog or updating the recorded expectation.

## File layout

```
scenarios_bug_queue/
  README.md          <- this file
  S01.md             <- bugs for section 1
  S02.md
  ...
  S25.md
```

Create a file lazily the first time a section records a bug; empty
sections have no file.

## Entry format

Each bug is a level-2 heading followed by a short metadata block and
free-form notes. Fields in square brackets are required:

```
## [scenario id] — one-line summary

- Discovered: YYYY-MM-DD
- Section: [section id + short name, e.g. "S04 — implicit indexes"]
- Scenario: [scenario id from SCENARIOS-mysql-implicit-behavior.md]
- Severity: [blocker | high | medium | low]
- Expected (MySQL 8.0): <what the container produced>
- Actual (omni): <what the omni catalog produced>
- Fix hint: <pointer to the suspected file / function / missing rule>

<free-form notes, reproduction SQL, information_schema dumps, etc.>
```

## Severity guide

- **blocker** — breaks a load path entirely (parse error, panic,
  wrong DDL accepted silently).
- **high** — observable semantic diff on common DDL (missing column
  metadata, wrong default, wrong index shape).
- **medium** — diff on edge-case DDL that real users occasionally hit.
- **low** — cosmetic / ordering / representation-only diff that does
  not affect downstream consumers.

## Workflow

1. Worker runs the section's scenario tests.
2. For every `t.Error` that fires, the worker adds one entry to
   `S<NN>.md` with enough detail to reproduce.
3. Worker reports counts (pass / fail / new bug entries) in its
   return summary.
4. Fix pass (separate batch) picks up entries, implements fixes,
   removes the entry, and re-runs to confirm green.
