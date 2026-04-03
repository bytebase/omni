# Snowflake Migration Analysis

## Grammar Coverage

The legacy ANTLR4 Snowflake parser is sourced from [antlr/grammars-v4](https://github.com/antlr/grammars-v4/tree/master/sql/snowflake) and covers a large surface area. Grammar files: `SnowflakeParser.g4` (141KB), `SnowflakeLexer.g4` (56KB). Generated `snowflake_parser.go` is 5.7MB.

### DDL Statements

#### CREATE (46 rules)

- **Database objects**: DATABASE, SCHEMA, TABLE (+ AS SELECT, LIKE, CLONE), VIEW, MATERIALIZED VIEW, DYNAMIC TABLE, EXTERNAL TABLE, EVENT TABLE
- **Functions & procedures**: FUNCTION, PROCEDURE, EXTERNAL FUNCTION
- **Data integration**: STAGE, FILE FORMAT, PIPE, STREAM, SEQUENCE, TAG
- **Integrations**: API INTEGRATION, SECURITY INTEGRATION (EXTERNAL_OAUTH, SNOWFLAKE_OAUTH, SAML2, SCIM), NOTIFICATION INTEGRATION, STORAGE INTEGRATION, NETWORK POLICY
- **Governance**: MASKING POLICY, ROW ACCESS POLICY, SESSION POLICY, PASSWORD POLICY
- **Semantic/analytics**: SEMANTIC VIEW, DATASET
- **Users & roles**: ROLE, USER, MANAGED ACCOUNT
- **Failover/replication**: FAILOVER GROUP, REPLICATION GROUP, GIT REPOSITORY
- **Misc**: ALERT, CONNECTION, RESOURCE MONITOR, SHARE, ACCOUNT, OBJECT CLONE

#### ALTER (54+ rules)

Full ALTER support for all CREATE-able objects. Notable specifics:
- ALTER TABLE: ADD/DROP/RENAME/MODIFY COLUMN, constraints
- ALTER DATABASE: RENAME, SWAP, PRIMARY/FAILOVER/REPLICATION
- ALTER DYNAMIC TABLE: RESUME/SUSPEND, REFRESH, SET params
- ALTER EXTERNAL TABLE: REFRESH, ADD/REMOVE FILES, PARTITIONS
- ALTER FUNCTION/PROCEDURE: signatures, set/unset properties
- ALTER TASK, STREAM, PIPE, SEQUENCE, STAGE, FILE FORMAT
- ALTER ROLE, USER, WAREHOUSE, ACCOUNT
- ALTER all integration types, all policy types
- ALTER FAILOVER/REPLICATION GROUPS, GIT REPOSITORY

#### DROP (37 rules)

All core object types with IF EXISTS, CASCADE/RESTRICT, signature support for functions/procedures.

#### UNDROP (6 rules)

DATABASE, SCHEMA, TABLE, DYNAMIC TABLE, TAG (time-travel recovery).

### DML Statements

- **SELECT**: Full query support with FROM, WHERE, GROUP BY, HAVING, QUALIFY, ORDER BY, LIMIT, TOP, CTEs (recursive), set operations (UNION ALL/BY NAME, EXCEPT, MINUS, INTERSECT)
- **INSERT**: Standard, OVERWRITE, multi-table (FIRST/ALL, WHEN...THEN), INSERT...SELECT, VALUES
- **UPDATE**: SET, FROM, WHERE
- **DELETE**: FROM, USING, WHERE
- **MERGE**: INTO with USING, ON, WHEN MATCHED/NOT MATCHED (UPDATE/DELETE/INSERT)

### Query Features

- **JOINs**: INNER, LEFT, RIGHT, FULL OUTER, NATURAL, CROSS, ASOF (with match conditions), DIRECTED
- **Table functions**: FLATTEN, SPLIT_TO_TABLE, LATERAL derived tables
- **PIVOT/UNPIVOT**: Full PIVOT with aggregation, IN clause; UNPIVOT with INCLUDE/EXCLUDE NULLS
- **MATCH_RECOGNIZE**: PARTITION BY, ORDER BY, measures, row match, after match skip
- **Time travel**: AT (TIMESTAMP/OFFSET/STATEMENT/STREAM), BEFORE, CHANGES
- **Window functions**: OVER, PARTITION BY, ORDER BY, aggregate/ranking with ROWS/RANGE frames
- **SAMPLE**: Sampling clause on tables

### DCL Statements

#### GRANT/REVOKE (8 rules)

- Global privileges (ON ACCOUNT)
- Object privileges (DATABASE, WAREHOUSE, INTEGRATION, RESOURCE MONITOR, USER)
- Schema privileges (CREATE TABLE/VIEW/MATERIALIZED VIEW/MASKING POLICY/ROW ACCESS POLICY, etc.)
- Schema object privileges (SELECT, INSERT, UPDATE, DELETE, TRUNCATE, REFERENCES, USAGE, READ, WRITE, MONITOR, OPERATE, APPLY)
- Future grants (ON FUTURE OBJECTS IN SCHEMA/DATABASE)
- Grant options (WITH GRANT OPTION, COPY/REVOKE CURRENT GRANTS)
- Role hierarchy (GRANT ROLE TO ROLE/USER)
- Share grants (GRANT TO SHARE, GRANT ON SHARE)

### SHOW Commands (60+ variants)

ALERTS, CHANNELS, COLUMNS, CONNECTIONS, DATABASES, DATASETS, DELEGATED AUTHORIZATIONS, DYNAMIC TABLES, EVENT TABLES, EXTERNAL FUNCTIONS, EXTERNAL TABLES, FAILOVER GROUPS, FILE FORMATS, FUNCTIONS, GIT BRANCHES, GIT REPOSITORIES, GIT TAGS, GLOBAL ACCOUNTS, GRANTS, INTEGRATIONS, LOCKS, MANAGED ACCOUNTS, MASKING POLICIES, MATERIALIZED VIEWS, NETWORK POLICIES, OBJECTS, ORGANIZATION ACCOUNTS, PARAMETERS, PIPES, PRIMARY KEYS, PROCEDURES, REGIONS, REPLICATION ACCOUNTS/DATABASES/GROUPS, RESOURCE MONITORS, ROLES, ROW ACCESS POLICIES, SCHEMAS, SECRETS, SEMANTIC VIEWS/DIMENSIONS/METRICS, SEQUENCES, SESSION POLICIES, SHARES, STREAMS, TABLES, TAGS, TASKS, TRANSACTIONS, USER FUNCTIONS, USERS, VARIABLES, VERSIONS, VIEWS, WAREHOUSES.

Filters: LIKE, IN (ACCOUNT/DATABASE/SCHEMA), STARTS WITH, LIMIT. Options: TERSE, HISTORY.

### DESCRIBE Commands (27 variants)

ALERT, DATABASE, DYNAMIC TABLE, EVENT TABLE, EXTERNAL TABLE, FILE FORMAT, FUNCTION, GIT REPOSITORY, INTEGRATION, MASKING POLICY, MATERIALIZED VIEW, NETWORK POLICY, PIPE, PROCEDURE, RESULT, ROW ACCESS POLICY, SCHEMA, SEARCH OPTIMIZATION, SEMANTIC VIEW, SEQUENCE, SESSION POLICY, PASSWORD POLICY, SHARE, STAGE, STREAM, TABLE, TASK, TRANSACTION, USER, VIEW, WAREHOUSE.

### USE Commands (5 variants)

DATABASE, SCHEMA, ROLE, WAREHOUSE, SECONDARY ROLES (ALL/NONE).

### Other Utility Statements

- **COMMENT**: ON object or COLUMN
- **COPY INTO**: TABLE/LOCATION with stages, file format, transformation queries, validation modes
- **GET/PUT**: File staging operations
- **TRANSACTIONS**: BEGIN/START, COMMIT, ROLLBACK (with WORK)
- **EXPLAIN**: TABULAR, JSON, TEXT formats
- **CALL**: Stored procedure invocation
- **EXECUTE**: IMMEDIATE (dynamic SQL), TASK
- **TRUNCATE**: TABLE (IF EXISTS), MATERIALIZED VIEW
- **SET/UNSET**: Session/statement parameters
- **LIST/REMOVE**: Stage operations

### Data Types (30+)

Integer types (INT, SMALLINT, TINYINT, BIGINT), numeric (NUMBER, DECIMAL, FLOAT, DOUBLE, REAL), string (VARCHAR, STRING, CHAR, TEXT, BINARY, VARBINARY), date/time (DATE, TIME, TIMESTAMP, TIMESTAMP_LTZ/NTZ/TZ, DATETIME), semi-structured (VARIANT, OBJECT, ARRAY), geospatial (GEOGRAPHY, GEOMETRY), VECTOR, BOOLEAN.

### Expression Support

Arithmetic, comparison, logical (AND, OR, NOT), string (LIKE, ILIKE, RLIKE), NULL checks (IS NULL, IS NOT NULL), BETWEEN, IN, EXISTS, CASE WHEN, CAST, TRY_CAST, IFF, COALESCE, NULLIF, subqueries, array/object access (bracket notation, colon notation, dot notation).

---

## Parse API Surface

Pure ANTLR4 with no wrapper layer. Public API:

### Functions

| Function | Signature | Purpose |
|----------|-----------|---------|
| `NewSnowflakeLexer` | `(input *antlr.InputStream) *SnowflakeLexer` | Create lexer (case-insensitive) |
| `NewSnowflakeParser` | `(stream *antlr.TokenStream) *SnowflakeParser` | Create parser from token stream |
| `Snowflake_file()` | `() Snowflake_fileContext` | Root parse rule (multi-statement) |

### Traversal

- **Listener**: `SnowflakeParserListener` / `SnowflakeParserBaseListener` (top-down walk)
- **Visitor**: `SnowflakeParserVisitor` / `SnowflakeParserBaseVisitor` (bottom-up traverse)
- **Walker**: `antlr.ParseTreeWalkerDefault.Walk(listener, tree)`

### Configuration

- `parser.BuildParseTrees = true`
- `parser.RemoveErrorListeners()` / `parser.AddErrorListener()`
- `lexer.RemoveErrorListeners()` / `lexer.AddErrorListener()`

---

## AST Types

**No AST layer exists.** The parser produces raw ANTLR4 parse tree contexts. Each grammar rule generates a context type (e.g., `Create_tableContext`, `Select_statementContext`). Consumers must traverse the parse tree directly.

---

## Gaps and Limitations

1. **Constraint validation**: No semantic validation of mutually exclusive constraints (e.g., INITIALLY DEFERRED + NOT DEFERRABLE). Grammar TODO at line 1314.
2. **Element ordering**: Incomplete validation of element order in constraint definitions. TODOs at lines 1422, 1425.
3. **Function arity**: Built-in functions not split by arity (no-param vs unary vs binary). TODO at line 4087.
4. **Comment-only SQL**: Grammar requires at least one statement; files with only comments fail.
5. **Dynamic SQL**: EXECUTE IMMEDIATE accepts arbitrary strings; no validation of dynamic SQL content.
6. **Geospatial functions**: Incomplete coverage of GIS built-in functions.
7. **Vector type**: Recently added; limited built-in vector function support.
8. **Reserved keyword collisions**: Some context-sensitive keywords may collide with identifiers (managed via build script).

---

## Bytebase Import Sites

### A. SQL Parsing & Statement Splitting

**Files:**
- `backend/plugin/parser/snowflake/snowflake.go`
- `backend/plugin/parser/snowflake/split.go`

**APIs called:**
- `ParseSnowSQL()` -- full parsing, returns `[]*base.ANTLRAST`
- `SplitSQL()` -- statement splitting via lexer with `SnowflakeLexerSEMI`
- `NewSnowflakeLexer()`, `NewSnowflakeParser()`, `Snowflake_file()`

**Registered as:**
```go
base.RegisterParseStatementsFunc(storepb.Engine_SNOWFLAKE, parseSnowflakeStatements)
base.RegisterSplitterFunc(storepb.Engine_SNOWFLAKE, SplitSQL)
```

### B. Query Validation & Execution

**Files:**
- `backend/plugin/parser/snowflake/query.go`

**APIs called:**
- `ParseSnowSQL()`, tree walk with listener
- Context: `Dml_command()`, `Other_command()`, `Describe_command()`, `Show_command()`, `Query_statement()`

**Data extracted:** Query validity, whether statement contains EXECUTE (affects result availability), valid statement type detection (DML queries only).

**Registered as:**
```go
base.RegisterQueryValidator(storepb.Engine_SNOWFLAKE, validateQuery)
```

### C. Query Type Classification

**Files:**
- `backend/plugin/parser/snowflake/query_type.go`

**Data extracted:** QueryType enum (DDL, DML, DQL/SELECT, system info queries). Used for statement classification.

### D. Query Span Extraction (Data Masking / Sensitive Data)

**Files:**
- `backend/plugin/parser/snowflake/query_span_extractor.go`
- `backend/plugin/parser/snowflake/query_span.go`

**APIs called:**
- `ParseSnowSQL()`, tree walk
- Context: `With_expression()`, `Select_statement_in_parentheses()`, `AllSet_operators()`, `Object_ref()`, `Full_column_name()`, `Exclude_clause()`, `Trim_expression()`, `Try_cast_expr()`

**Data extracted:** Table sources (database.schema.table), column sources, column destinations, CTE definitions, recursive CTE handling, mixed system/user table detection.

**Registered as:**
```go
base.RegisterGetQuerySpan(storepb.Engine_SNOWFLAKE, GetQuerySpan)
```

### E. Diagnostics (Syntax Checking)

**Files:**
- `backend/plugin/parser/snowflake/diagnose.go`

**APIs called:**
- `NewSnowflakeLexer()`, `NewSnowflakeParser()`, `Snowflake_file()`
- Custom error listeners for lexer and parser

**Data extracted:** Syntax errors with line/column positions.

**Registered as:**
```go
base.RegisterDiagnoseFunc(storepb.Engine_SNOWFLAKE, Diagnose)
```

### F. Statement Position Ranges

**Files:**
- `backend/plugin/parser/snowflake/statement_ranges.go`

**APIs called:**
- `NewSnowflakeLexer()`, token constants (`SnowflakeParserEOF`, `SnowflakeParserSEMI`)
- `base.GetANTLRStatementRangesUTF16Position()`

**Data extracted:** Statement position ranges (UTF-16 aware) for IDE highlighting.

**Registered as:**
```go
base.RegisterStatementRangesFunc(storepb.Engine_SNOWFLAKE, GetStatementRanges)
```

### G. SQL Review & Linting (15 advisor rules)

**Files:** All in `backend/plugin/advisor/snowflake/`

**Advisor rules:**

| Rule | Context Types Used |
|------|--------------------|
| `rule_naming_table` | `Create_tableContext`, `Alter_tableContext` |
| `rule_column_require` | `Create_tableContext`, `Column_decl_item_listContext`, `Alter_tableContext` |
| `rule_column_no_null` | `Full_col_declContext`, `Out_of_line_constraintContext`, `Alter_column_declContext` |
| `rule_column_maximum_varchar_length` | `Data_typeContext` |
| `rule_table_require_pk` | `Create_tableContext`, `Out_of_line_constraintContext` |
| `rule_table_no_foreign_key` | `Out_of_line_constraintContext` |
| `rule_where_require_select` | `Select_statementContext` |
| `rule_where_require_update_delete` | `Update_statementContext`, `Delete_statementContext` |
| `rule_table_drop_naming_convention` | `Drop_tableContext` |
| `rule_naming_identifier_case` | `Id_Context` |
| `rule_naming_identifier_no_keyword` | `Id_Context` |
| `rule_naming_table_no_keyword` | `Table_nameContext` |
| `rule_select_no_select_all` | `Select_list_elemContext` |
| `rule_migration_compatibility` | `Create_tableContext`, `Create_table_as_selectContext`, `Create_schemaContext`, `Create_databaseContext` |
| `generic_checker` | All DDL/DML/column contexts (base framework) |

### H. Query Limit Rewriting

**Files:**
- `backend/plugin/db/snowflake/query.go`

**APIs called:**
- `ParseSnowSQL()`, `antlr.NewTokenStreamRewriter()`, tree walk
- Context: `AllSet_operators()`, `Select_statement_in_parentheses()`, `Limit_clause()`, `LIMIT()`, `OFFSET()`, `AllNum()`

**Data extracted/modified:** Rewrites SELECT queries to add/modify LIMIT for safety.

---

## Feature Dependency Map

| Bytebase Feature | Parser APIs Used | Data Extracted |
|------------------|-----------------|----------------|
| SQL parsing | `ParseSnowSQL()` | ANTLR parse tree per statement |
| Statement splitting | `SplitSQL()` via lexer | Statement list with positions |
| Query type classification | `ParseSnowSQL()` + walk | QueryType enum (DDL/DML/DQL) |
| Query validation | `ParseSnowSQL()` + walk | Valid/invalid + has_execute flag |
| Query span (data masking) | `ParseSnowSQL()` + walk | Table/column dependencies, lineage |
| Syntax diagnostics | Lexer + parser + error listeners | Errors with line/column positions |
| Statement ranges | Lexer tokens | UTF-16 position ranges |
| SQL review (15 rules) | `ParseSnowSQL()` + walk | Advice list (violations) |
| Query limit rewrite | `ParseSnowSQL()` + token rewriter | Modified SQL text |
| Schema sync | N/A (uses `GET_DDL()` DB function) | N/A |

---

## Priority Ranking

### P0: Actively consumed by bytebase (must exist before import migration)

1. **Lexer / tokenization** -- foundation for splitting, ranges, diagnostics
2. **Parser / parse tree** -- foundation for all tree-walking consumers
3. **Statement splitting** -- used by query execution, migration pipeline
4. **Query type classification** -- DDL/DML/DQL routing
5. **Query validation** -- SQL editor safety gate
6. **Syntax diagnostics** -- IDE error display
7. **Statement ranges** -- IDE highlighting
8. **Query span extraction** -- data masking, sensitive data detection, column lineage
9. **Query limit rewriting** -- query execution safety
10. **SQL review framework** -- generic checker + 14 specific rules

### P1: Supported by legacy parser but not currently consumed by bytebase

11. **DDL: CREATE TABLE** (full syntax including constraints, AS SELECT, LIKE, CLONE)
12. **DDL: ALTER TABLE** (ADD/DROP/RENAME/MODIFY COLUMN, constraints)
13. **DDL: DROP TABLE** (IF EXISTS, CASCADE/RESTRICT)
14. **DDL: CREATE/ALTER/DROP VIEW, MATERIALIZED VIEW**
15. **DDL: CREATE/ALTER/DROP SCHEMA, DATABASE**
16. **DDL: CREATE/ALTER/DROP FUNCTION, PROCEDURE, EXTERNAL FUNCTION**
17. **DDL: CREATE/ALTER/DROP DYNAMIC TABLE, EXTERNAL TABLE, EVENT TABLE**
18. **DDL: CREATE/ALTER/DROP STAGE, FILE FORMAT, PIPE, STREAM, SEQUENCE, TAG**
19. **DDL: CREATE/ALTER/DROP all integration types** (API, SECURITY, NOTIFICATION, STORAGE)
20. **DDL: CREATE/ALTER/DROP all policy types** (MASKING, ROW ACCESS, SESSION, PASSWORD, NETWORK)
21. **DDL: CREATE/ALTER/DROP ROLE, USER, MANAGED ACCOUNT, WAREHOUSE**
22. **DDL: CREATE/ALTER/DROP FAILOVER/REPLICATION GROUP, GIT REPOSITORY**
23. **DDL: UNDROP** (DATABASE, SCHEMA, TABLE, DYNAMIC TABLE, TAG)
24. **DDL: Remaining CREATE/ALTER/DROP** (ALERT, CONNECTION, RESOURCE MONITOR, SHARE, ACCOUNT, SEMANTIC VIEW, DATASET)
25. **DML: INSERT** (standard, OVERWRITE, multi-table, INSERT...SELECT)
26. **DML: UPDATE, DELETE, MERGE**
27. **DML: SELECT** (full query features: CTEs, JOINs, PIVOT/UNPIVOT, MATCH_RECOGNIZE, time travel, window functions, SAMPLE)
28. **DCL: GRANT/REVOKE** (all variants)
29. **SHOW commands** (60+ variants)
30. **DESCRIBE commands** (27 variants)
31. **USE commands** (5 variants)
32. **Utility: COPY INTO, GET, PUT** (file staging)
33. **Utility: COMMENT, TRUNCATE, EXPLAIN, CALL, EXECUTE**
34. **Utility: SET/UNSET, LIST/REMOVE, BEGIN/COMMIT/ROLLBACK**
35. **Expression system** (arithmetic, comparison, logical, string, NULL, BETWEEN, IN, EXISTS, CASE, CAST, subqueries, semi-structured access)
36. **Data types** (30+ types including semi-structured, geospatial, vector)

---

## Full Coverage Target

Everything listed above. The omni Snowflake parser must implement all P0 and P1 features for full parity with the legacy ANTLR4 parser. P0 items are prioritized for implementation because they unblock the bytebase import migration. P1 items are implemented for coverage parity even if bytebase does not currently consume them directly.
