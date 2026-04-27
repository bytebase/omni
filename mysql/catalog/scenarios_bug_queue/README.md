# MySQL Implicit-Behavior Scenario Bug Queue

This directory is a historical archive of bugs discovered while
executing the `mysql-implicit-behavior` starmap. It is no longer the
active source of truth for MySQL implicit-behavior compatibility.

The current source of truth is:

- `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md`
- the corresponding `TestScenario_*` tests under `mysql/catalog`

As of the 2026-04-27 reconciliation, the tracker is fully verified:
239 total scenarios, 239 verified, 0 remaining. Historical entries in
this directory are preserved for context and may describe behavior that
has already been fixed or reclassified. Treat an entry as active only if
it is explicitly re-opened with a newer date and linked to a failing
test.

## File layout

```
scenarios_bug_queue/
  README.md          <- this file
  c1.md              <- archived C1 notes
  c2.md
  ...
  c25.md
  ax.md              <- archived ALTER TABLE path notes
  ps.md              <- archived CREATE vs ALTER path-split notes
```

These files are retained as archival notes. Do not add new compatibility
findings here unless they are also represented in the tracker and in a
failing or skipped `TestScenario_*` test.

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

## Current workflow

1. Add or update the scenario in
   `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md`.
2. Encode the observed MySQL behavior in the matching `TestScenario_*`
   test, preferably with a container-backed assertion when session state
   or server-version behavior matters.
3. Fix the omni catalog behavior or explicitly mark the scenario out of
   scope in the tracker.
4. Use this archive only as supporting historical context.
