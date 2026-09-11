# MySQL BNF Extraction Skill

You are extracting BNF syntax definitions from MySQL 8.0 Reference Manual documentation.

**Working directory:** `<repo root>`
**Catalog:** `mysql/parser/MYSQL_BNF_CATALOG.json`
**Output directory:** `mysql/parser/bnf/`

## Your Task

1. Read `mysql/parser/MYSQL_BNF_CATALOG.json`
2. Find the next batch of statements to process:
   - Find statements with `"status": "url_ok"` (not yet extracted)
   - Process up to **10 statements per invocation** to stay within context limits
3. For each statement:
   a. Read the local doc file from `mysql/parser/docs/{slug}.txt` if it exists and is not a NOT_FOUND marker
   b. If local doc doesn't exist or is a NOT_FOUND marker, use WebFetch on the statement's `url` field
   c. Extract the **complete BNF/syntax definition** — every clause, every option, every branch
   d. Save the BNF to `mysql/parser/bnf/{slug}.bnf`
   e. Update the statement's status in `mysql/parser/MYSQL_BNF_CATALOG.json` to `"bnf_done"`
4. If a statement's BNF cannot be extracted, set status to `"bnf_failed"` with an `"error"` field

If all statements have status `"bnf_done"` or `"bnf_failed"`, output `ALL_BNF_EXTRACTED` and stop.

## BNF File Format

Each `.bnf` file must follow this exact format:

```
-- Statement: CREATE TABLE
-- Source: https://dev.mysql.com/doc/refman/8.0/en/create-table.html

CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
    (create_definition,...)
    [table_options]
    [partition_options]

CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
    [(create_definition,...)]
    [table_options]
    [partition_options]
    [IGNORE | REPLACE]
    [AS] query_expression

CREATE [TEMPORARY] TABLE [IF NOT EXISTS] tbl_name
    { LIKE old_tbl_name | (LIKE old_tbl_name) }

create_definition: {
    col_name column_definition
  | {INDEX | KEY} [index_name] [index_type] (key_part,...) [index_option] ...
  | {FULLTEXT | SPATIAL} [INDEX | KEY] [index_name] (key_part,...) [index_option] ...
  | [CONSTRAINT [symbol]] PRIMARY KEY [index_type] (key_part,...) [index_option] ...
  | [CONSTRAINT [symbol]] UNIQUE [INDEX | KEY] [index_name] [index_type] (key_part,...) [index_option] ...
  | [CONSTRAINT [symbol]] FOREIGN KEY [index_name] (col_name,...) reference_definition
  | check_constraint_definition
}
```

## Rules

1. **Extract the COMPLETE syntax** — every branch, every optional clause, every sub-production
2. **Include all sub-clauses** that are defined on the same page (e.g., `create_definition`, `column_definition`, `table_options`)
3. **Do NOT abbreviate** with `...` — expand all branches
4. **Use MySQL's syntax notation**: `[ ]` for optional, `{ | }` for alternatives, `...` only where MySQL docs use it for repetition
5. **Include the top-level statement syntax AND all named sub-clauses** as separate productions
6. **For simple statements** (e.g., DROP xxx), the BNF may be just 2-3 lines — that's fine
7. **Do NOT invent syntax** — only extract what the documentation actually shows

## Progress Logging

Print progress markers:
```
[BNF] Processing: CREATE TABLE (1/10)
[BNF] Saved: mysql/parser/bnf/create-table.bnf (85 lines)
[BNF] Processing: CREATE INDEX (2/10)
...
[BNF] Batch complete: 10 extracted, 0 failed
```
