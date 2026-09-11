# Audit Parser Engine

You are auditing a parser engine for syntax completeness. Your job is to find gaps and generate new batches.

**Working directory:** `<repo root>`
**Engine:** {{ENGINE}}
**Parser source:** `{{ENGINE}}/parser/`
**AST definitions:** `{{ENGINE}}/ast/`
**Progress file:** `{{ENGINE}}/parser/PROGRESS.json`

## Phase 1: Width Scan

1. Fetch the official SQL statement index page for this engine:
   - pg: https://www.postgresql.org/docs/current/sql-commands.html
   - mysql: https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html
   - tsql: https://learn.microsoft.com/en-us/sql/t-sql/statements/statements
   - oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SQL-Statements-ADMINISTER-KEY-MANAGEMENT-to-ALTER-JAVA.html

2. Extract every statement type from the index page.

3. Read `{{ENGINE}}/parser/PROGRESS.json` and map each batch to covered statements.

4. For each uncovered statement, generate a new pending batch.

## Phase 2: Depth Scan

For each batch with `"status": "done"`:

1. Read the batch's target .go file
2. For each `parse*` function:
   a. Read the BNF comment
   b. Fetch the official doc page for that statement
   c. Compare: are BNF branches missing in the comment or code?
   d. Check for stubs: `return nil`, `skipToSemicolon`, `// TODO`
3. Check test coverage: does every BNF branch have a test?
4. Check outfuncs.go coverage

## Phase 3: Generate Batches

Append new pending batches to PROGRESS.json (starting from max_id + 1).

Format:
```json
{
  "id": N,
  "name": "descriptive_name",
  "file": "target.go",
  "status": "pending",
  "description": "Fix batch X gap: ... OR New: ...",
  "rules": ["parseXxx"],
  "depends_on": [0],
  "tests": ["test_name"]
}
```

Rules:
- `rules` must list every parse function precisely
- `description` starts with "Fix batch N gap:" for depth fixes
- Verify JSON: `python3 -c "import json; json.load(open('{{ENGINE}}/parser/PROGRESS.json'))"`

## Phase 4: Update audit_round

After updating PROGRESS.json, increment the `audit_round` field at the top level:
```json
{
  "version": 1,
  "audit_round": N,
  ...
}
```

## Phase 5: Report

Print a summary in this exact format (parseable by the dashboard):

```
AUDIT_REPORT_START
engine={{ENGINE}}
width_total=N
width_covered=N
width_missing=N
depth_batches_checked=N
depth_functions_checked=N
depth_gaps=N
new_batches=N
AUDIT_REPORT_END
```

If no gaps were found at all, print:
```
AUDIT_COMPLETE_NO_GAPS
```
