# MSSQL BNF Extraction Skill

You are extracting BNF syntax definitions for T-SQL statements. You have TWO sources:
1. **Primary**: `mssql/parser/bnf/TSql170.g` — the official ANTLR grammar from Microsoft's sqlscriptdom (35000 lines)
2. **Secondary**: Microsoft T-SQL documentation for additional context and human-readable BNF

**Working directory:** `<repo root>`
**Catalog:** `mssql/parser/MSSQL_BNF_CATALOG.json`
**Grammar file:** `mssql/parser/bnf/TSql170.g`
**Output directory:** `mssql/parser/bnf/`

## Your Task

1. Read `mssql/parser/MSSQL_BNF_CATALOG.json`
2. Find the next batch of statements to process:
   - Find statements with `"status": "url_ok"` (not yet extracted)
   - Process up to **8 statements per invocation** (grammar extraction is more work)
3. For each statement:
   a. Find the grammar rule in `TSql170.g` using the `rule` field from the catalog (e.g., `createTableStatement`)
   b. Extract the complete rule and all referenced sub-rules
   c. If the statement has a `url` field, also WebFetch the Microsoft docs page to get the human-readable BNF
   d. Combine both into a `.bnf` file
   e. Save to `mssql/parser/bnf/{slug}.bnf`
   f. Update the statement's status in `mssql/parser/MSSQL_BNF_CATALOG.json` to `"bnf_done"`
4. For statements with `"status": "no_doc"`, extract BNF solely from the grammar file
5. If a statement's BNF cannot be extracted, set status to `"bnf_failed"` with an `"error"` field

If all statements have status `"bnf_done"`, `"bnf_failed"`, or `"no_doc"` with extracted grammar, output `ALL_BNF_EXTRACTED` and stop.

## BNF File Format

Each `.bnf` file must follow this format:

```
-- Statement: CREATE TABLE
-- Grammar rule: createTableStatement (TSql170.g)
-- Doc: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql

-- From TSql170.g:
createTableStatement:
    CREATE TABLE tableName
    LeftParen columnDefinition (Comma columnDefinition)* RightParen
    (ON fileGroup)?
    (TEXTIMAGE_ON fileGroup)?
    (WITH LeftParen tableOption (Comma tableOption)* RightParen)?
    ;

-- From Microsoft docs:
CREATE TABLE
    { database_name.schema_name.table_name | schema_name.table_name | table_name }
    ( { <column_definition> | <computed_column_definition> | <column_set_definition>
      | <table_constraint> | <table_index> } [ ,...n ]
    )
    [ ON { partition_scheme_name ( partition_column_name )
         | filegroup | "default" } ]
    [ TEXTIMAGE_ON { filegroup | "default" } ]
    [ FILESTREAM_ON { partition_scheme_name | filegroup | "default" } ]
    [ WITH ( <table_option> [ ,...n ] ) ]
[ ; ]
```

## Rules

1. **Always extract from TSql170.g first** — this is the authoritative source
2. **Include the ANTLR rule AND all directly referenced sub-rules**
3. **Then add the human-readable BNF from docs** if available
4. **For grammar rules**: search for the rule name (e.g., `createTableStatement`) and extract the complete rule body including all alternatives
5. **Follow rule references**: if a rule references `columnDefinition`, include that sub-rule too (one level deep)
6. **Do NOT extract the entire grammar** — only rules relevant to the specific statement

## Searching the Grammar File

The grammar file is large (35000 lines). To find rules efficiently:
- Use Grep to search for the rule name: `grep -n "createTableStatement" mssql/parser/bnf/TSql170.g`
- Read the surrounding lines to get the full rule
- Rules end at the next rule definition (a line starting with a lowercase letter followed by colon or `returns`)

## Progress Logging

Print progress markers:
```
[BNF] Processing: CREATE TABLE (1/8)
[BNF] Grammar rule found at line 1234
[BNF] Saved: mssql/parser/bnf/create-table-transact-sql.bnf (60 lines)
[BNF] Batch complete: 8 extracted, 0 failed
```
