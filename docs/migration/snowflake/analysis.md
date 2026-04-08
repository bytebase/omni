# Snowflake Migration Analysis

Migration target: full coverage of the legacy ANTLR4 `bytebase/parser/snowflake` grammar by a hand-written recursive-descent parser under `omni/snowflake`. Bytebase consumption is used to **prioritize** the order of work; it is never used to **scope down** legacy parser features.

Sources analyzed:

- Legacy parser: `/Users/h3n4l/OpenSource/parser/snowflake` (`SnowflakeParser.g4` 5,718 lines, 669 rules; `SnowflakeLexer.g4` 1,275 lines)
- Downstream consumer: `/Users/h3n4l/OpenSource/bytebase` (25 files, ~4,367 LOC consuming the parser)

---

## Grammar Coverage

### Entry rules

| Rule | Parses |
|------|--------|
| `snowflake_file` | Complete SQL file (batch of commands) |
| `batch` | One or more commands separated by `;` |
| `sql_command` | Any single SQL command |
| `dml_command` | INSERT/UPDATE/DELETE/MERGE/SELECT |
| `ddl_command` | CREATE/ALTER/DROP/UNDROP |
| `query_statement` | Outer SELECT wrapper (with WITH/set operators) |

### DDL — CREATE

Covers `CREATE [OR REPLACE] [IF NOT EXISTS]` for:

- Storage objects: `DATABASE`, `SCHEMA`, `TABLE` (TRANSIENT, CTAS, LIKE, CLUSTER BY, COPY GRANTS, WITH TAGS), `VIEW` (SECURE), `MATERIALIZED VIEW`, `DYNAMIC TABLE`, `EXTERNAL TABLE`, `EVENT TABLE`, `SEQUENCE`
- Compute / data movement: `WAREHOUSE`, `STAGE` (internal/external, encryption, directory table, auto-ingest), `FILE FORMAT` (CSV/JSON/AVRO/ORC/PARQUET), `PIPE` (auto-ingest + COPY INTO), `STREAM` (TABLE/EXTERNAL/STAGE, INSERT_ONLY, AT/BEFORE), `TASK` (cron schedule, BEGIN…END Snowflake-Scripting bodies), `ALERT`
- Functions / procedures: `FUNCTION` (SQL/JavaScript/Python/Java, IMMUTABLE/VOLATILE/SECURE, external functions), `PROCEDURE`, language-specific handlers
- Security: `ROLE`, `USER`, `MASKING POLICY`, `ROW ACCESS POLICY`, `PASSWORD POLICY`, `SESSION POLICY`, `NETWORK POLICY`, `SECURITY INTEGRATION` (EXTERNAL_OAUTH/SNOWFLAKE_OAUTH/SAML2/SCIM)
- Integrations: `STORAGE INTEGRATION`, `API INTEGRATION`, `NOTIFICATION INTEGRATION` (AWS SNS / Azure Event Grid / GCP Pub/Sub), `RESOURCE MONITOR`, `SECRET`, `CONNECTION`, `GIT REPOSITORY`
- Replication / sharing: `FAILOVER GROUP`, `REPLICATION GROUP`, `ACCOUNT`, `MANAGED ACCOUNT`, `SHARE`
- Newer surface: `DATASET`, `SEMANTIC VIEW`, `SEMANTIC DIMENSION`, `SEMANTIC METRIC`, `TAG`

### DDL — ALTER

Supports `ALTER` for every object listed above plus `ALTER SESSION`, `ALTER ACCOUNT`. ALTER TABLE alone covers: ADD/DROP/RENAME COLUMN, ALTER COLUMN (SET DEFAULT/NOT NULL/DATA TYPE/MASKING POLICY), RECLUSTER, RENAME TO, SWAP WITH, SET/UNSET (comment, tags, clustering), ADD/DROP CONSTRAINT, CLUSTER BY edits.

### DDL — DROP / UNDROP

DROP for all object types with `[IF EXISTS]` and `[CASCADE|RESTRICT]`. UNDROP for `DATABASE`, `SCHEMA`, `TABLE`, `DYNAMIC TABLE`, `TAG`.

### DML

- `SELECT` — `WITH [RECURSIVE]`, `[DISTINCT|ALL]`, `TOP n`, `EXCLUDE`, GROUP BY (CUBE / GROUPING SETS / ROLLUP / `GROUP BY ALL`), `HAVING`, `QUALIFY`, `ORDER BY` (NULLS FIRST/LAST), `LIMIT/OFFSET/FETCH`
- Joins: INNER, LEFT/RIGHT/FULL OUTER, CROSS, NATURAL, ASOF, DIRECTED, LATERAL, table functions, `MATCH_RECOGNIZE`, `PIVOT`/`UNPIVOT`, `CHANGES`, `AT`/`BEFORE` (time travel), `SAMPLE`/`TABLESAMPLE`
- Set ops: `UNION [ALL]`, `UNION ALL BY NAME`, `EXCEPT`/`MINUS`, `INTERSECT`
- `INSERT [OVERWRITE]` (single + multi-table `INSERT ALL/FIRST WHEN…THEN`), `UPDATE` (USING joins), `DELETE` (USING), `MERGE` (matched/not-matched/not-matched-by-source)
- `COPY INTO` (load + unload), `PUT`, `GET`, `LIST`, `REMOVE`, `CALL`, `EXECUTE IMMEDIATE`, `EXECUTE TASK`, `EXPLAIN [USING TABULAR|DESCRIBE|ANALYZE]`

### DCL / TCL / Utility

- `GRANT`/`REVOKE` (roles, object privileges, FROM SHARE, WITH GRANT OPTION)
- `BEGIN`/`START TRANSACTION`/`COMMIT`/`ROLLBACK` (NAME id), `SAVEPOINT`
- 50+ `SHOW` variants, full `DESCRIBE`/`DESC` set, `USE` (DATABASE/SCHEMA/ROLE/WAREHOUSE/SECONDARY ROLES)
- `COMMENT ON`, `SET`, `UNSET`, `TRUNCATE`, `ABORT*`

### Snowflake-specific surface

- Time travel `AT (TIMESTAMP|OFFSET|STATEMENT|STREAM => …)` and `BEFORE (STATEMENT => …)`
- Zero-copy `CLONE`
- Semi-structured types: `VARIANT`, `OBJECT`, `ARRAY[<type>]`, `GEOGRAPHY`, `GEOMETRY`, `VECTOR(<type>, <dim>)`; JSON path `expr:field`, `expr[idx]`
- Snowflake Scripting (parsed structurally inside TASK/PROC/FUNCTION bodies): `DECLARE`, `BEGIN…END`, `:=` assignment, `IF/CASE/FOR/WHILE/REPEAT`, `EXCEPTION`, `RETURN`, cursor / resultset references
- Window functions, `WITHIN GROUP` for `LISTAGG`/`ARRAY_AGG`, `MATCH_RECOGNIZE`, lambda `param -> expr`, `CONNECT BY`/`START WITH`, `LATERAL FLATTEN`/`SPLIT_TO_TABLE`, `IDENTIFIER()`

### Data types

INT/INTEGER/SMALLINT/TINYINT/BYTEINT/BIGINT, NUMBER(p,s)/NUMERIC/DECIMAL, FLOAT/FLOAT4/FLOAT8/DOUBLE/DOUBLE PRECISION/REAL, BOOLEAN, DATE/TIME(p)/TIMESTAMP[_NTZ|_TZ|_LTZ](p)/DATETIME, VARCHAR(n)/CHAR(n)/NCHAR/CHARACTER/TEXT/STRING, BINARY/VARBINARY, VARIANT/OBJECT/ARRAY[<type>]/GEOGRAPHY/GEOMETRY, VECTOR(elem, dim).

---

## Parse API Surface (legacy)

```go
parser.NewSnowflakeLexer(antlr.CharStream) *SnowflakeLexer
parser.NewSnowflakeParser(antlr.TokenStream) *SnowflakeParser

(*SnowflakeParser).Snowflake_file() ISnowflake_fileContext   // primary entry
(*SnowflakeParser).Batch()          IBatchContext
(*SnowflakeParser).Sql_command()    ISql_commandContext
```

Errors flow through the standard ANTLR `antlr.ErrorListener` (`SyntaxError`, `ReportAmbiguity`, …). Bytebase relies on:

- `AddErrorListener` / `RemoveErrorListeners`
- `BuildParseTrees` flag
- Token constants `SnowflakeLexerSEMI`, `SnowflakeParserSEMI`, `SnowflakeParserEOF`
- `antlr.TokenStreamRewriter` for in-place query rewriting (LIMIT injection)
- `antlr.ParseTreeWalkerDefault.Walk` over the listener interface

---

## AST Surface

The legacy parser exposes:

- Listener interface `SnowflakeParserListener` + `BaseSnowflakeParserListener` (Enter/Exit per rule)
- Visitor interface `SnowflakeParserVisitor` + `BaseSnowflakeParserVisitor`
- One Context type per grammar rule (~669)

**Most-relied-on context types in bytebase** (used by ≥3 features):

- Top-level: `BatchContext`, `Sql_commandContext`, `Ddl_commandContext`, `Dml_commandContext`, `Other_commandContext`, `Show_commandContext`, `Describe_commandContext`
- Query: `Query_statementContext`, `Select_statementContext`, `Select_statement_in_parenthesesContext`, `With_expressionContext`, `Common_table_expressionContext`, `Set_operatorsContext`, `From_clauseContext`, `Where_clauseContext`, `Limit_clauseContext`, `Select_list_topContext`, `Select_list_elemContext`, `Column_elem_starContext`, `Column_elemContext`, `Exclude_clauseContext`, `Object_refContext`, `Object_nameContext`, `Object_name_or_aliasContext`, `Column_nameContext`, `Column_listContext`, `Join_clauseContext`
- DDL: `Create_databaseContext`, `Create_schemaContext`, `Create_tableContext`, `Create_table_as_selectContext`, `Alter_tableContext`, `Drop_tableContext`, `Drop_schemaContext`, `Drop_databaseContext`, `Full_col_declContext`, `Col_declContext`, `Inline_constraintContext`, `Out_of_line_constraintContext`, `Primary_key_clauseContext`, `Foreign_key_clauseContext`, `Null_not_nullContext`, `Data_typeContext`, `Varchar_typeContext`, `Table_column_actionContext`
- Identifiers: `Id_Context` (used by every linter)

---

## Bytebase Import Sites

Import path: `github.com/bytebase/parser/snowflake`. 25 backend files / ~4,367 LOC.

| File | Parser entry | Purpose |
|------|--------------|---------|
| `backend/plugin/parser/snowflake/snowflake.go` | `NewSnowflakeLexer`, `NewSnowflakeParser`, `Snowflake_file` | Common parse helper used by everything else |
| `backend/plugin/parser/snowflake/split.go` | Lexer + `SnowflakeLexerSEMI` | Statement splitting |
| `backend/plugin/parser/snowflake/statement_ranges.go` | Lexer + token constants | Byte ranges of statements (UTF-16) |
| `backend/plugin/parser/snowflake/diagnose.go` | Parser + error listener | Syntax diagnostics for editor |
| `backend/plugin/parser/snowflake/query.go` | Listener walk | Query validation for SQL editor |
| `backend/plugin/parser/snowflake/query_type.go` | Listener walk | DDL/DML/SELECT classification |
| `backend/plugin/parser/snowflake/query_span_extractor.go` (1,589 LOC) | Listener walk over 50+ contexts | Query span / field lineage / masking scope |
| `backend/plugin/parser/snowflake/query_span.go` | Wrapper for above | — |
| `backend/plugin/db/snowflake/query.go` | `TokenStreamRewriter` | Inject `LIMIT` for safe query preview |
| `backend/plugin/advisor/snowflake/generic_checker.go` | Reflection-based listener dispatch | Lint rule infrastructure |
| `backend/plugin/advisor/snowflake/rule_*.go` (14 rules) | Specific contexts (see below) | Lint rules |

### Lint rules (14)

| Rule | AST nodes inspected |
|------|---------------------|
| `column_no_null` | `Create_tableContext`, `Full_col_declContext`, `Null_not_nullContext`, `Inline_constraintContext`, `Alter_tableContext`, `Table_column_actionContext` |
| `column_require` | `Create_tableContext`, `Full_col_declContext`, `Col_declContext` |
| `column_maximum_varchar_length` | `Create_tableContext`, `Full_col_declContext`, `Data_typeContext`, `Varchar_typeContext` |
| `table_require_pk` | `Create_tableContext`, `Alter_tableContext`, `Out_of_line_constraintContext`, `Primary_key_clauseContext`, `Column_listContext` |
| `table_no_foreign_key` | `Create_tableContext`, `Inline_constraintContext`, `Out_of_line_constraintContext`, `Foreign_key_clauseContext` |
| `naming_table` | `Create_tableContext`, `Alter_tableContext`, `Create_schemaContext`, `Create_databaseContext`, `Object_nameContext` |
| `naming_table_no_keyword` | `Create_tableContext`, `Alter_tableContext`, `Id_Context` |
| `naming_identifier_case` | `Id_Context` |
| `naming_identifier_no_keyword` | `Id_Context`, `Column_nameContext`, `Create_table_as_selectContext` |
| `select_no_select_all` | `Select_list_elemContext`, `Column_elem_starContext` |
| `where_require_select` | `Select_statementContext`, `From_clauseContext` |
| `where_require_update_delete` | `Update_statementContext`, `Delete_statementContext`, `Where_clauseContext` |
| `migration_compatibility` | `Create_*`/`Alter_tableContext`/`Object_nameContext` |
| `table_drop_naming_convention` | `Drop_tableContext`, `Drop_schemaContext`, `Drop_databaseContext`, `Object_nameContext` |

### Custom AST walkers

| Walker | File | Role |
|--------|------|------|
| `selectOnlyListener` | `query_span_extractor.go` | Extract result columns + their source tables |
| `accessTablesListener` | `query_span_extractor.go` | Collect every `Object_ref` (touched table) |
| `queryTypeListener` | `query_type.go` | Classify statement type |
| `queryValidateListener` | `query.go` | Validate query for editor |
| `snowsqlRewriter` | `db/snowflake/query.go` | Token-level LIMIT injection |
| `GenericChecker` | `advisor/snowflake/generic_checker.go` | Dispatch lint rules via `EnterEveryRule`/`ExitEveryRule` |

### Normalization helpers shared everywhere

`NormalizeSnowSQLObjectNamePart`, `NormalizeSnowSQLSchemaName`, `NormalizeSnowSQLObjectName`, `normalizeSnowflakeColumnName`, `normalizedObjectNameOrAlias`, `ExtractSnowSQLOrdinaryIdentifier`, `IsSnowflakeKeyword`. The new omni package needs equivalent helpers.

---

## Feature Dependency Map

| Bytebase feature | Parser APIs / contexts | Data extracted | Critical? |
|------------------|------------------------|----------------|-----------|
| Statement parsing & splitting | Lexer + `Snowflake_file`, `SnowflakeLexerSEMI`, error listeners | Token boundaries, parse tree, syntax errors | **CRITICAL** — foundation for everything else |
| Statement byte ranges | Lexer + `EOF`/`SEMI` constants | Statement byte positions | CRITICAL |
| Query type classification | `BatchContext` → DDL/DML/SELECT/SHOW/DESCRIBE/Other | Statement category | HIGH |
| Query validation (editor) | `Sql_commandContext`, `Dml_commandContext` | Validity, EXECUTE presence | MEDIUM |
| Query span / field lineage | Deep walk over 50+ contexts (`Query_statementContext`, `Select_statementContext`, `With_expressionContext`, `From_clauseContext`, `Object_refContext`, …) | Tables touched, result-column lineage, CTE merging, set-op merging | **CRITICAL** — drives masking |
| Result LIMIT injection | `TokenStreamRewriter` + `Query_statementContext`, `Select_statementContext`, `Limit_clauseContext`, `Set_operatorsContext`, `Select_list_topContext` | Where to splice LIMIT | MEDIUM |
| Lint dispatcher | `BaseSnowflakeParserListener.EnterEveryRule/ExitEveryRule`, rule-name reflection | Generic listener dispatch | **CRITICAL** (blocks all 14 rules) |
| 14 lint rules | `Create_tableContext`, `Full_col_declContext`, `Object_nameContext`, `Id_Context`, `Out_of_line_constraintContext`, `Foreign_key_clauseContext`, `Select_list_elemContext`, `Where_clauseContext`, … | Table/column shape, identifiers, statement structure | **CRITICAL** (parallelizable) |
| Syntax diagnostics | Parser + error listener | Errors with positions | MEDIUM |

Notably **absent** from bytebase's snowflake usage (compared to other engines):

- No schema-sync-via-parser (Snowflake schemas come from the live database)
- No standalone masking provider (masking is derived from query span only)
- No deparser / SQL-string round-trip
- No catalog or completion provider yet

---

## Gaps and Limitations of the Legacy Parser

These are documented so the omni rewrite can decide whether to inherit or fix them.

1. **Snowflake Scripting is structural-only.** `BEGIN…END`, `IF`, loops, `DECLARE`, `EXCEPTION` blocks parse but get no scope/type validation.
2. **Window frame specs are simplified** vs the SQL standard (`ROWS BETWEEN … PRECEDING/FOLLOWING` is incomplete).
3. **EXECUTE IMMEDIATE** treats the inner SQL as an opaque string.
4. **Some constraint combinations** (e.g. `INITIALLY DEFERRED` vs `NOT DEFERRABLE`) are unenforced — there is a TODO in the grammar.
5. **Comments** (`--`, `/* */`) are tokenized but not exposed in the AST by default.
6. **Custom error messages** are absent — only stock ANTLR diagnostics.
7. **No UDT (user-defined type)** grammar.
8. The generated Go file (`snowflake_parser.go`) is **211k lines / ~5MB** — a clear motivation for rewriting.

The hand-written omni parser should at minimum match the legacy structural coverage; improving any of the above is a stretch goal, not a requirement.

---

## Priority Ranking

Order in which features should be migrated. Foundational items must land first because every downstream consumer depends on them.

### Tier 0 — Foundation (blocks everything)

1. **Lexer** — keywords, identifiers (incl. quoted), numeric/string/binary literals, JSON path operators (`:`, `[]`), `::` cast, `||` concat, semicolon/EOF tokens, comments. Must keep enough token-position metadata for `statement_ranges.go` and `TokenStreamRewriter`.
2. **AST core types** — node interfaces, position tracking, walker infrastructure (visit/listener-style API parity).
3. **Statement splitter** — uses lexer only; unblocks bytebase's `split.go` and `statement_ranges.go`.
4. **Parser entry** — `ParseSnowflake(sql) → File`; collect syntax errors with positions.

### Tier 1 — Core SQL surface (unblocks query span + classification + most lints)

5. **Identifiers and qualified names** — `Id`, `Object_name` (db.schema.obj), normalization helpers (`NormalizeSnowSQLObjectNamePart` etc.).
2. **Data types** — full type grammar (NUMBER/VARCHAR/VARIANT/ARRAY/OBJECT/VECTOR/GEOGRAPHY/GEOMETRY/TIMESTAMP_*).
3. **Expressions** — operators, function calls, `CASE`, `CAST`/`TRY_CAST`, JSON path, lambda `->`, subqueries, IN/BETWEEN/LIKE/ILIKE/RLIKE.
4. **`SELECT` core** — SELECT list (incl. `*`, `EXCLUDE`), FROM (table refs, joins), WHERE, GROUP BY (incl. CUBE/ROLLUP/GROUPING SETS/ALL), HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET/FETCH/TOP.
5. **CTEs (`WITH [RECURSIVE]`) and set operators** (UNION/UNION BY NAME/EXCEPT/MINUS/INTERSECT).
6. **Statement classification helper** — equivalent of bytebase's `query_type.go`.

### Tier 2 — DDL backbone (unblocks 14 lint rules)

11. **CREATE TABLE** (full) — column declarations, inline + out-of-line constraints (PK/FK/UNIQUE/CHECK/NOT NULL/DEFAULT/IDENTITY/AUTOINCREMENT/MASKING POLICY/COLLATE/WITH TAGS), CTAS, LIKE, CLUSTER BY, COPY GRANTS.
2. **ALTER TABLE** (full set of actions).
3. **CREATE/ALTER/DROP DATABASE / SCHEMA / VIEW / MATERIALIZED VIEW**.
4. **DROP / UNDROP** for the above.
5. **Lint dispatcher port** + the 14 advisor rules (parallelizable once Tier 2 contexts exist).

### Tier 3 — Query span extractor parity

16. **Query span extractor** in omni — full reimplementation of the 1,589-LOC bytebase walker on the new AST. Needs joins, CTE merging, set-op column merging, subquery field resolution, `EXCLUDE`, `Object_ref` collection.
2. **Result LIMIT injection** — needs an omni equivalent of `TokenStreamRewriter` (or AST-level rewrite) over `Query_statementContext`/`Select_statementContext`/`Limit_clauseContext`.

### Tier 4 — Remaining DDL surface (parity with legacy grammar)

18. CREATE/ALTER/DROP for: `STAGE`, `FILE FORMAT`, `PIPE`, `STREAM`, `TASK` (incl. Snowflake-Scripting bodies), `DYNAMIC TABLE`, `EXTERNAL TABLE`, `EVENT TABLE`, `SEQUENCE`, `FUNCTION`, `PROCEDURE`, `ALERT`.
2. Security DDL: `ROLE`, `USER`, `MASKING POLICY`, `ROW ACCESS POLICY`, `PASSWORD/SESSION/NETWORK POLICY`, `SECURITY INTEGRATION`.
3. Integration DDL: `STORAGE/API/NOTIFICATION INTEGRATION`, `RESOURCE MONITOR`, `SECRET`, `CONNECTION`, `GIT REPOSITORY`.
4. Replication: `FAILOVER GROUP`, `REPLICATION GROUP`, `ACCOUNT`, `MANAGED ACCOUNT`, `SHARE`.
5. Newer surface: `DATASET`, `SEMANTIC VIEW`, `SEMANTIC DIMENSION`, `SEMANTIC METRIC`, `TAG`.

### Tier 5 — Remaining DML / Snowflake specifics

23. `INSERT` (single + multi-table `INSERT ALL/FIRST`), `UPDATE` (USING), `DELETE` (USING), `MERGE` (matched / not-matched / not-matched-by-source).
2. `COPY INTO` (load + unload), `PUT`, `GET`, `LIST`, `REMOVE`.
3. `MATCH_RECOGNIZE`, `PIVOT`/`UNPIVOT`, `LATERAL FLATTEN`/`SPLIT_TO_TABLE`, `CONNECT BY`/`START WITH`, `CHANGES`, `AT`/`BEFORE` time travel, `SAMPLE`/`TABLESAMPLE`, `CLONE`.
4. `CALL`, `EXECUTE IMMEDIATE`, `EXECUTE TASK`, `EXPLAIN` variants.

### Tier 6 — Auxiliary commands

27. `GRANT`/`REVOKE` (roles + privileges + FROM SHARE + WITH GRANT OPTION).
2. `BEGIN`/`COMMIT`/`ROLLBACK`/`SAVEPOINT`/`START TRANSACTION`.
3. `SHOW` (50+ variants), `DESCRIBE`/`DESC`, `USE`, `SET`/`UNSET`, `COMMENT ON`, `TRUNCATE`, `ABORT*`.

### Tier 7 — Snowflake Scripting bodies

30. Structural parse of `DECLARE`, `BEGIN…END`, `:=`, `IF`, `CASE`, `FOR`, `WHILE`, `REPEAT`, `EXCEPTION`, `RETURN`, cursor / resultset references inside TASK/PROC/FUNCTION bodies. Match legacy depth — no semantic checking.

---

## Full Coverage Target

Every legacy grammar feature is in scope. Each is tagged P0 (consumed by bytebase today) or P1 (legacy-parity only). P0 must exist before bytebase imports can switch.

| Area | Item | Tier | Tag |
|------|------|------|-----|
| Lexer | Tokens, comments, identifiers, literals, operators, JSON path, `::` cast, semicolon | 0 | P0 |
| Parser entry | `ParseSnowflake`, `Batch`, `Sql_command` | 0 | P0 |
| Statement splitter | Token-driven splitter | 0 | P0 |
| Statement byte ranges | Position metadata | 0 | P0 |
| Identifiers | `Id`, normalization helpers | 1 | P0 |
| Object names | `Object_name` (1/2/3-part), aliases | 1 | P0 |
| Data types | All Snowflake types incl. VARIANT/OBJECT/ARRAY/VECTOR/GEOGRAPHY/GEOMETRY | 1 | P0 |
| Expressions | Arith/cmp/logical, `CASE`, `CAST`/`TRY_CAST`, JSON path, subqueries, lambdas, function calls | 1 | P0 |
| SELECT core | SELECT list, FROM, WHERE, GROUP BY (+CUBE/ROLLUP/GROUPING SETS/ALL), HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET/FETCH/TOP, EXCLUDE | 1 | P0 |
| Joins | INNER/LEFT/RIGHT/FULL OUTER/CROSS/NATURAL/ASOF/DIRECTED/LATERAL | 1 | P0 |
| CTEs | `WITH [RECURSIVE]` | 1 | P0 |
| Set ops | UNION [ALL], UNION ALL BY NAME, EXCEPT/MINUS, INTERSECT | 1 | P0 |
| Statement classification | DDL/DML/SELECT/SHOW/DESCRIBE/Other helper | 1 | P0 |
| CREATE TABLE | Full grammar incl. constraints, CTAS, LIKE, CLUSTER BY, COPY GRANTS, WITH TAGS | 2 | P0 |
| ALTER TABLE | Full action set | 2 | P0 |
| CREATE/ALTER/DROP DATABASE/SCHEMA/VIEW/MATERIALIZED VIEW | — | 2 | P0 |
| DROP / UNDROP | All object types | 2 | P0 (table/schema/db), P1 (others) |
| Lint dispatcher | Generic listener with rule registration | 2 | P0 |
| 14 lint rules | column_no_null, column_require, column_maximum_varchar_length, table_require_pk, table_no_foreign_key, naming_table, naming_table_no_keyword, naming_identifier_case, naming_identifier_no_keyword, select_no_select_all, where_require_select, where_require_update_delete, migration_compatibility, table_drop_naming_convention | 2 | P0 |
| Query span extractor | Result-column lineage, table-access listener, CTE/set-op merging | 3 | P0 |
| Result LIMIT injection | Token-stream-style rewrite or AST splice | 3 | P0 |
| Syntax diagnostics | Error listing for editor | 3 | P0 |
| INSERT (single + multi-table) | — | 5 | P0 |
| UPDATE / DELETE | — | 5 | P0 |
| MERGE | All match clauses | 5 | P0 |
| CREATE/ALTER/DROP STAGE / FILE FORMAT / PIPE / STREAM / TASK / DYNAMIC TABLE / EXTERNAL TABLE / EVENT TABLE / SEQUENCE / FUNCTION / PROCEDURE / ALERT | — | 4 | P1 |
| Security DDL (ROLE/USER/MASKING POLICY/ROW ACCESS POLICY/POLICY/SECURITY INTEGRATION) | — | 4 | P1 |
| Integration DDL (STORAGE/API/NOTIFICATION INTEGRATION, RESOURCE MONITOR, SECRET, CONNECTION, GIT REPOSITORY) | — | 4 | P1 |
| Replication / sharing (FAILOVER/REPLICATION GROUP, ACCOUNT, MANAGED ACCOUNT, SHARE) | — | 4 | P1 |
| Newer surface (DATASET, SEMANTIC VIEW/DIMENSION/METRIC, TAG) | — | 4 | P1 |
| COPY INTO / PUT / GET / LIST / REMOVE | — | 5 | P1 |
| MATCH_RECOGNIZE, PIVOT/UNPIVOT, LATERAL FLATTEN/SPLIT_TO_TABLE, CONNECT BY/START WITH, CHANGES, AT/BEFORE, SAMPLE/TABLESAMPLE, CLONE | — | 5 | P1 |
| CALL, EXECUTE IMMEDIATE, EXECUTE TASK, EXPLAIN | — | 5 | P1 |
| GRANT/REVOKE | — | 6 | P1 |
| Transaction (BEGIN/COMMIT/ROLLBACK/SAVEPOINT) | — | 6 | P1 |
| SHOW / DESCRIBE / USE / SET / UNSET / COMMENT ON / TRUNCATE / ABORT* | — | 6 | P1 |
| Snowflake Scripting bodies (DECLARE/BEGIN…END/IF/CASE/FOR/WHILE/REPEAT/EXCEPTION/RETURN/cursors) | — | 7 | P1 |

P0 = required before bytebase can switch its imports. P1 = required for full parity with the legacy grammar (the migration target).

---

## Test Corpus Strategy (decided)

The omni snowflake parser must be tested against **two independent corpora**, both of which must parse cleanly:

### Corpus A — Legacy ANTLR4 examples (regression baseline)
Lift every `.sql` file from `/Users/h3n4l/OpenSource/parser/snowflake/examples/` (27 files):

```
alerts.sql            create_function.sql       describe.sql        materialized_views.sql  show.sql
alter.sql             create_pipe.sql           drop.sql            merge_statement.sql     task_scripting.sql
at_before.sql         create_procedure.sql      dynamic_table.sql   other.sql               undrop.sql
call.sql              create_table.sql          grant.sql           select.sql              use.sql
comment.sql           create_view.sql           having.sql          create.sql
create_fileformat.sql                                ids.sql                                  iff.sql
```

These are what bytebase has historically been able to parse. The omni parser must match — any regression here is a blocker.

### Corpus B — Official Snowflake SQL reference (authoritative source)
Recursively scrape `https://docs.snowflake.com/en/sql-reference-commands` and its subpages (one page per command, e.g. `/sql/create-table`, `/sql/alter-table`, `/sql/select`, `/sql/merge`, …). For each subpage, extract every code example (typically in `<pre>` / fenced-code blocks) and add it to the corpus. Each command's "Examples" section is the primary target.

These are authoritative — when there is a conflict between the legacy `.g4` grammar and the official docs, **the docs win** (the legacy grammar may have bugs). Cross-reference both when implementing each rule.

### Corpus organization in the omni repo
Suggested layout (to be confirmed in the planning phase):

```
snowflake/parser/testdata/
  legacy/        # mirrored from bytebase/parser/snowflake/examples
    alter.sql
    select.sql
    ...
  official/      # scraped from docs.snowflake.com
    create-table/
      example_01.sql
      example_02.sql
      ...
    select/
      example_01.sql
      ...
```

Each tier of the migration adds the relevant subset of both corpora to the green-test set. By the time the migration completes, every file in both corpora must parse without error.

---

## Architecture Decisions

### Scope: full legacy parity (decided)
The omni snowflake parser must cover **the entire grammar of `bytebase/parser/snowflake`**, including Tier 7 (Snowflake Scripting bodies). Scripting blocks are parsed structurally only — no scope/type checking — matching legacy depth. Nothing from the legacy grammar is dropped.

### Reference engine: `pg/` as primary template, `mongo/analysis` as the model for query span (decided)

Survey of existing omni engines:

| Engine | Packages | Parser LOC | Statement splitter | Query analysis | Notes |
|--------|----------|-----------|---------------------|----------------|-------|
| `mssql` | ast, completion, parser | 34,116 | No | No | Minimal supporting packages |
| `pg` | ast, catalog, completion, parser, parsertest, pgregress, plpgsql, semantic | 30,652 | **Yes** (sophisticated: dollar-quoting, BEGIN ATOMIC, nested comments) | No | Most mature ecosystem; cleanest AST/parser/semantic split |
| `mongo` | analysis, ast, catalog, completion, parser, parsertest | 2,453 | Yes | **Yes** (operation/collection/field extraction, ~600 LOC) | Best example of an `analysis/` package |
| `mysql` | ast, catalog, completion, **deparse**, parser, quality, semantic | 23,860 | Yes | No | Only engine with a `deparse/` package; AST-level rewrite (`deparse/rewrite.go`) for SHOW CREATE VIEW resolver compat |
| `oracle` | ast, parser, quality | 36,540 | Yes | No | Minimal supporting packages |
| `cosmosdb` | analysis, ast, parser | 2,027 | Yes | Yes (~400 LOC) | Smallest analysis example |

**`pg/` wins as the structural template** because:
1. It has the most sophisticated statement splitter — Snowflake needs the same level of robustness for `BEGIN…END` scripting blocks, dollar-style quoting, and nested comments.
2. It has the cleanest separation of `ast/`, `parser/`, `semantic/`, `parsertest/` — Snowflake will grow into a similar shape.
3. Its `parsertest/` is the right home for the two-corpus regression suite (Corpus A + Corpus B).
4. It is actively maintained with ~70 `*_test.go` files.

**`mongo/analysis/`** is the right model for the **query span extractor** (Tier 3, the single biggest bytebase consumer at 1,589 LOC). Its `StatementAnalysis` shape — operation classification, collection set, field set, predicate fields — is structurally close to what bytebase's `selectOnlyListener` + `accessTablesListener` produces.

`mssql` was a tempting match on feature shape (DDL-heavy, no schema sync, no completion), but it lacks both a statement splitter and any analysis package — the two pieces Snowflake needs most.

Target package layout for `omni/snowflake/`:
```
snowflake/
  ast/                     # AST node types (mirror pg/ast structure)
  parser/                  # Lexer + recursive descent parser + statement splitter
    testdata/
      legacy/              # Corpus A — lifted from bytebase/parser/snowflake/examples
      official/            # Corpus B — scraped from docs.snowflake.com
  parsertest/              # Table-driven parser tests (mirror pg/parsertest)
  analysis/                # Query span / field lineage (mirror mongo/analysis)
  advisor/                 # 14 lint rules (new, modeled on bytebase consumer)
  deparse/                 # AST → SQL string + LIMIT injection rewrite
```

### LIMIT injection: AST-level rewrite via a `deparse/` package (decided)

**No engine in omni currently implements token-stream rewriting**, and no engine uses an ANTLR-style `TokenStreamRewriter`. The closest analog is `mysql/deparse/rewrite.go`, which performs AST-to-AST transformations and then deparses back to SQL text (used for SHOW CREATE VIEW resolver compat).

For Snowflake's LIMIT injection:
1. Parse the SELECT into an AST
2. Walk to the deepest top-level `SELECT` (descending into set operators / parentheses, **not** into subqueries)
3. Append or replace a `LIMIT` clause node
4. Deparse the modified AST back to SQL text

This is structurally simpler than the legacy bytebase approach (which uses ANTLR's `TokenStreamRewriter` to splice tokens). The trade-off is that we need a `deparse/` package — but Snowflake will need one anyway for any future use case that round-trips SQL, and `mysql/deparse` is a working template. The bytebase consumer call site (`backend/plugin/db/snowflake/query.go`) will need a thin shim around the new omni API instead of the current ANTLR rewriter.

If we ever need exact-text-preservation rewriting (e.g., to keep formatting/comments intact), we can add a token-stream rewriter later — but the LIMIT-injection use case does not require it.
