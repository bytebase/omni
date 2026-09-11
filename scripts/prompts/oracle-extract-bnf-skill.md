# Oracle BNF Extraction Skill

You are extracting BNF syntax definitions from Oracle SQL Reference documentation.

**Working directory:** `<repo root>`
**Catalog:** `oracle/parser/ORACLE_BNF_CATALOG.json`
**Output directory:** `oracle/parser/bnf/`

## Your Task

1. Read `oracle/parser/ORACLE_BNF_CATALOG.json`
2. Find the next batch of statements to process:
   - Find statements with `"status": "url_ok"` (not yet extracted)
   - Process up to **10 statements per invocation** to stay within context limits
3. For each statement:
   a. Read the local doc file from `oracle/parser/docs/{slug}.txt` if it exists
   b. If local doc doesn't exist or is a NOT_FOUND marker, use WebFetch on the statement's `url` field
   c. Extract the **complete BNF/syntax definition** — every clause, every option, every branch
   d. Save the BNF to `oracle/parser/bnf/{slug}.bnf`
   e. Update the statement's status in `ORACLE_BNF_CATALOG.json` to `"bnf_done"`
4. If a statement's BNF cannot be extracted, set status to `"bnf_failed"` with an `"error"` field

If all statements have status `"bnf_done"` or `"bnf_failed"`, output `ALL_BNF_EXTRACTED` and stop.

## BNF File Format

Each `.bnf` file must follow this exact format:

```
-- Statement: CREATE TABLE
-- Source: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html

CREATE [ { GLOBAL TEMPORARY | PRIVATE TEMPORARY | SHARDED | DUPLICATED |
           IMMUTABLE | BLOCKCHAIN | DEFAULT COLLATION collation_name } ]
    TABLE [ IF NOT EXISTS ] [ schema. ] table_name
    [ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
    { relational_table | object_table | XMLType_table }
    [ MEMOPTIMIZE FOR READ ] [ MEMOPTIMIZE FOR WRITE ]
    [ PARENT [ schema. ] table ] ;

relational_table:
    [ ( relational_properties ) ]
    ...

relational_properties:
    { column_definition
    | virtual_column_definition
    | period_definition
    | out_of_line_constraint
    | out_of_line_ref_constraint
    | supplemental_logging_props
    } [, ...]
```

## Rules

1. **Extract the COMPLETE syntax** — every branch, every optional clause, every sub-production
2. **Include all sub-clauses** that are defined on the same page (e.g., `relational_properties`, `column_definition`)
3. **Do NOT abbreviate** with `...` — expand all branches
4. **Use the exact syntax notation** from Oracle docs: `[ ]` for optional, `{ | }` for alternatives, `[, ...]` for repetition
5. **Include the top-level statement syntax AND all named sub-clauses** as separate productions
6. **For simple statements** (e.g., DROP xxx), the BNF may be just 2-3 lines — that's fine
7. **Do NOT invent syntax** — only extract what the documentation actually shows

## Progress Logging

Print progress markers:
```
[BNF] Processing: CREATE TABLE (1/10)
[BNF] Saved: oracle/parser/bnf/CREATE-TABLE.bnf (45 lines)
[BNF] Processing: CREATE INDEX (2/10)
...
[BNF] Batch complete: 10 extracted, 0 failed
```
