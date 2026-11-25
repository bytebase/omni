# Cloud Spanner GoogleSQL truth1 — Index

All documented syntax forms for Cloud Spanner GoogleSQL dialect.
Source corpus: `spanner/` subdirectory.
version: spanner-current

MERGE support: **NOT SUPPORTED** in Cloud Spanner GoogleSQL.

---

## DDL forms (`ddl.md`) — 57 forms

| form_id | source_url | note |
|---------|-----------|------|
| DDL-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_database | CREATE DATABASE database_id |
| DDL-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_database | ALTER DATABASE SET OPTIONS (default_leader, optimizer_version, etc.) |
| DDL-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_table | CREATE TABLE with column defs, PRIMARY KEY, INTERLEAVE IN PARENT, ROW DELETION POLICY, OPTIONS |
| DDL-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#column_default | Column DEFAULT (expression) |
| DDL-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#generated_column | Generated column AS (expr) STORED / VIRTUAL |
| DDL-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#allow_commit_timestamp | TIMESTAMP OPTIONS (allow_commit_timestamp = true) |
| DDL-007 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#interleave_in_parent | INTERLEAVE IN PARENT table ON DELETE CASCADE/NO ACTION |
| DDL-008 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#foreign_key | FOREIGN KEY with REFERENCES and ON DELETE |
| DDL-009 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#check_constraint | CHECK (bool_expression) with optional CONSTRAINT name |
| DDL-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#row_deletion_policy | ROW DELETION POLICY (OLDER_THAN (col, INTERVAL n DAY)) |
| DDL-011 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE ADD COLUMN [IF NOT EXISTS] |
| DDL-012 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE DROP COLUMN [IF EXISTS] |
| DDL-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE ALTER COLUMN (type change, SET/DROP DEFAULT, SET OPTIONS) |
| DDL-014 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE ADD/DROP CONSTRAINT |
| DDL-015 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE SET OPTIONS |
| DDL-016 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE ADD/ALTER/DROP ROW DELETION POLICY |
| DDL-017 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table | ALTER TABLE RENAME TO |
| DDL-018 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_table | DROP TABLE [IF EXISTS] |
| DDL-019 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_index | CREATE [UNIQUE] [NULL_FILTERED] INDEX ... STORING (...) [, INTERLEAVE IN ...] |
| DDL-020 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_index | ALTER INDEX ADD/DROP STORED COLUMN |
| DDL-021 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_index | DROP INDEX [IF EXISTS] |
| DDL-022 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_view | CREATE [OR REPLACE] VIEW SQL SECURITY INVOKER/DEFINER AS query |
| DDL-023 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_view | DROP VIEW [IF EXISTS] |
| DDL-024 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_change_stream | CREATE CHANGE STREAM FOR ALL / FOR table (...) OPTIONS |
| DDL-025 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_change_stream | ALTER CHANGE STREAM SET FOR / DROP FOR ALL / SET OPTIONS |
| DDL-026 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_change_stream | DROP CHANGE STREAM |
| DDL-027 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_sequence | CREATE SEQUENCE OPTIONS (sequence_kind='bit_reversed_positive', skip_range, start_with_counter) |
| DDL-028 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_sequence | ALTER SEQUENCE SET OPTIONS |
| DDL-029 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_sequence | DROP SEQUENCE [IF EXISTS] |
| DDL-030 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_schema | CREATE SCHEMA [IF NOT EXISTS] |
| DDL-031 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_schema | DROP SCHEMA [IF EXISTS] |
| DDL-032 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_role | CREATE ROLE role_name |
| DDL-033 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_role | DROP ROLE role_name |
| DDL-034 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#grant | GRANT SELECT/INSERT/UPDATE/DELETE/EXECUTE ON TABLE/VIEW/CHANGE STREAM TO ROLE |
| DDL-035 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#grant_role | GRANT ROLE role TO ROLE target |
| DDL-036 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#revoke | REVOKE privileges ON object FROM ROLE |
| DDL-037 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#revoke_role | REVOKE ROLE role FROM ROLE target |
| DDL-038 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_model | CREATE [OR REPLACE] MODEL INPUT/OUTPUT REMOTE OPTIONS (endpoints) |
| DDL-039 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_model | ALTER MODEL SET OPTIONS |
| DDL-040 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_model | DROP MODEL [IF EXISTS] |
| DDL-041 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_locality_group | CREATE LOCALITY GROUP OPTIONS (storage='ssd'/'hdd', ssd_to_hdd_spill_timespan) |
| DDL-042 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_locality_group | ALTER LOCALITY GROUP SET OPTIONS |
| DDL-043 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_locality_group | DROP LOCALITY GROUP |
| DDL-044 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_placement | CREATE PLACEMENT OPTIONS (instance_partition, default_leader) |
| DDL-045 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_placement | DROP PLACEMENT |
| DDL-046 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_proto_bundle | CREATE PROTO BUNDLE (proto_type_name, ...) |
| DDL-047 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_proto_bundle | ALTER PROTO BUNDLE INSERT/UPDATE/DELETE (proto_type_name) |
| DDL-048 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_search_index | CREATE SEARCH INDEX ... ON ... (tokenlist_col) STORING PARTITION BY ORDER BY WHERE INTERLEAVE IN OPTIONS |
| DDL-049 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_search_index | ALTER SEARCH INDEX ADD/DROP STORED COLUMN |
| DDL-050 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_search_index | DROP SEARCH INDEX [IF EXISTS] |
| DDL-051 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#tokenlist_column | TOKENLIST column type with TOKENIZE_FULLTEXT/TOKENIZE_NGRAMS/etc. AS STORED |
| DDL-052 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_vector_index | CREATE VECTOR INDEX ON (embedding_col) OPTIONS (distance_type, num_leaves); ARRAY<FLOAT32>(vector_length=>) |
| DDL-053 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_vector_index | DROP VECTOR INDEX [IF EXISTS] |
| DDL-054 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#sequences | NEXT VALUE FOR sequence_name (in DEFAULT expression) |
| DDL-055 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#sequences | GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq) (in DEFAULT expression) |
| DDL-056 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#named_schema | Schema-qualified table names in DDL (schema_name.table_name) |
| DDL-057 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#locality_group | Column-level OPTIONS (locality_group = 'name') for columnar storage |

---

## DML forms (`dml.md`) — 10 forms

| form_id | source_url | note |
|---------|-----------|------|
| DML-001 | https://cloud.google.com/spanner/docs/dml-syntax#insert-syntax | INSERT INTO ... (cols) VALUES (...) [, ...] [THEN RETURN] |
| DML-002 | https://cloud.google.com/spanner/docs/dml-syntax#insert-or-ignore | INSERT OR IGNORE — skip on PK conflict |
| DML-003 | https://cloud.google.com/spanner/docs/dml-syntax#insert-or-update | INSERT OR UPDATE — upsert on PK conflict |
| DML-004 | https://cloud.google.com/spanner/docs/dml-syntax#insert-select | INSERT ... SELECT query_statement |
| DML-005 | https://cloud.google.com/spanner/docs/dml-syntax#then-return | THEN RETURN [WITH ACTION [AS alias]] select_list |
| DML-006 | https://cloud.google.com/spanner/docs/dml-syntax#update-syntax | UPDATE ... SET ... [FROM] WHERE ... [THEN RETURN] |
| DML-007 | https://cloud.google.com/spanner/docs/dml-syntax#update-from | UPDATE ... SET ... FROM from_item WHERE ... |
| DML-008 | https://cloud.google.com/spanner/docs/dml-syntax#delete-syntax | DELETE [FROM] table WHERE ... [THEN RETURN] |
| DML-009 | https://cloud.google.com/spanner/docs/dml-syntax#statement-hints | @{PDML_MAX_PARALLELISM=n} statement hint on DML |
| DML-010 | https://cloud.google.com/spanner/docs/dml-syntax | MERGE — NOT SUPPORTED in Cloud Spanner GoogleSQL |

---

## Query forms (`query.md`) — 26 forms

| form_id | source_url | note |
|---------|-----------|------|
| QUERY-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#sql_syntax | query_statement top-level grammar with optional statement hint, ORDER BY, LIMIT/OFFSET, FOR UPDATE |
| QUERY-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_list | SELECT ALL/DISTINCT, AS STRUCT/VALUE, select_list, * EXCEPT, * REPLACE |
| QUERY-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#from_clause | FROM clause: table, join, subquery, field_path, unnest, cte, tablesample |
| QUERY-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_types | JOIN: CROSS, INNER, LEFT OUTER, RIGHT OUTER, FULL OUTER, HASH; ON/USING |
| QUERY-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#where_clause | WHERE bool_expression |
| QUERY-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#group_by_clause | GROUP BY [group_hint] expr / ROLLUP |
| QUERY-007 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#having_clause | HAVING bool_expression |
| QUERY-008 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#order_by_clause | ORDER BY expr [ASC/DESC] [NULLS FIRST/LAST] |
| QUERY-009 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#limit_and_offset_clause | LIMIT count [OFFSET skip] |
| QUERY-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#for_update | FOR UPDATE clause |
| QUERY-011 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators | UNION {ALL|DISTINCT} |
| QUERY-012 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators | INTERSECT {ALL|DISTINCT} |
| QUERY-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#set_operators | EXCEPT {ALL|DISTINCT} |
| QUERY-014 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#with_clause | WITH cte AS (...), multiple CTEs |
| QUERY-015 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#with_clause | WITH RECURSIVE — NOT SUPPORTED |
| QUERY-016 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#subqueries | Subqueries: scalar, EXISTS, IN, FROM subquery |
| QUERY-017 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#unnest_operator | UNNEST(array) [WITH OFFSET] |
| QUERY-018 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#field_path | Field path in FROM clause for nested struct/array fields |
| QUERY-019 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#tablesample_operator | TABLESAMPLE BERNOULLI(n PERCENT) / RESERVOIR(n ROWS) |
| QUERY-020 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#implicit_alias | Implicit comma-join (FROM t1, t2) |
| QUERY-021 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_as_struct | SELECT AS STRUCT |
| QUERY-022 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#select_as_value | SELECT AS VALUE |
| QUERY-023 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#qualify_clause | QUALIFY clause (post-window filter) |
| QUERY-024 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#window_clause | WINDOW clause (named window definitions) |
| QUERY-025 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#parenthesized_query_expressions | Parenthesized query expression (query_expr) |
| QUERY-026 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#graph_table_operator | GRAPH_TABLE operator (Spanner Graph — noted, not detailed) |

---

## Data type forms (`datatypes.md`) — 18 forms

| form_id | source_url | note |
|---------|-----------|------|
| TYPE-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#boolean_type | BOOL |
| TYPE-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#integer_type | INT64 |
| TYPE-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#floating_point_type | FLOAT32 |
| TYPE-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#floating_point_type | FLOAT64 |
| TYPE-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#numeric_type | NUMERIC (38 digits, 9 decimal) |
| TYPE-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#string_type | STRING(length) / STRING(MAX) |
| TYPE-007 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#bytes_type | BYTES(length) / BYTES(MAX) |
| TYPE-008 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#json_type | JSON |
| TYPE-009 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#date_type | DATE |
| TYPE-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#timestamp_type | TIMESTAMP (with allow_commit_timestamp, PENDING_COMMIT_TIMESTAMP) |
| TYPE-011 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#array_type | ARRAY<T> and ARRAY<FLOAT32>(vector_length=>) |
| TYPE-012 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#struct_type | STRUCT<field type, ...> with named/anonymous fields |
| TYPE-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#tokenlist_type | TOKENLIST — for full-text/vector search |
| TYPE-014 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#proto_type | Proto buffer type (backtick-quoted FQTN, via PROTO BUNDLE) |
| TYPE-015 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#interval_type | INTERVAL (used in date arithmetic and ROW DELETION POLICY) |
| TYPE-016 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#enum_type | Proto ENUM type (via PROTO BUNDLE) |
| TYPE-017 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#casting | CAST(expr AS type) / SAFE_CAST(expr AS type) |
| TYPE-018 | https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#null_type | NULL type (untyped NULL literal) |

---

## Hints forms (`hints.md`) — 13 forms

| form_id | source_url | note |
|---------|-----------|------|
| HINT-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#statement_hints | Statement hint @{key=val, ...} before SQL statement |
| HINT-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#table_hints | Table hint table_name @{key=val} |
| HINT-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_hints | Join hint JOIN @{key=val} |
| HINT-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#force_index_hint | FORCE_INDEX = index_name / _BASE_TABLE |
| HINT-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#groupby_scan_optimization_hint | GROUPBY_SCAN_OPTIMIZATION = TRUE/FALSE |
| HINT-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#optimizer_version_hint | OPTIMIZER_VERSION = integer / 'latest' |
| HINT-007 | https://cloud.google.com/spanner/docs/query-execution-hints#optimizer_statistics_package_hint | OPTIMIZER_STATISTICS_PACKAGE = 'package_name' |
| HINT-008 | https://cloud.google.com/spanner/docs/query-execution-hints#use_additional_parallelism | USE_ADDITIONAL_PARALLELISM = TRUE/FALSE |
| HINT-009 | https://cloud.google.com/spanner/docs/dml-syntax#statement-hints | PDML_MAX_PARALLELISM = integer (DML/PDML hint) |
| HINT-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#join_hints | JOIN_METHOD = HASH_JOIN / APPLY_JOIN / LOOP_JOIN / MERGE_JOIN |
| HINT-011 | https://cloud.google.com/spanner/docs/query-execution-hints#lock_scanned_ranges | LOCK_SCANNED_RANGES = exclusive / shared |
| HINT-012 | https://cloud.google.com/spanner/docs/query-execution-hints#scan_method | SCAN_METHOD = BATCH / ROW |
| HINT-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/query-syntax#hints | Multiple hints combined in one @{...} block |

---

## Lexical forms (`lexical.md`) — 21 forms

| form_id | source_url | note |
|---------|-----------|------|
| LEX-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#identifiers | Unquoted identifiers: letter/_ start, letter/digit/_ body, max 128 chars, case-insensitive |
| LEX-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#identifiers | Backtick-quoted identifiers for keywords, spaces, special chars |
| LEX-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Single-quoted string literals |
| LEX-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Double-quoted string literals |
| LEX-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Triple-quoted string literals (''' or """) |
| LEX-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Raw string literals r'' r"" r'''''' r"""""" |
| LEX-007 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Bytes literals b'' B'' b"""...""" |
| LEX-008 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals | Raw bytes literals rb'' br'' RB"" BR"" |
| LEX-009 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#escape_sequences | Escape sequences \\, \', \", \n, \t, \xhh, \uhhhh, \Uhhhhhhhh |
| LEX-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#integer_literals | Integer literals: decimal and 0x hex |
| LEX-011 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#floating_point_literals | Float literals: 3.14, .5, 1.23e10 |
| LEX-012 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#numeric_literals | NUMERIC 'decimal_string' typed literal |
| LEX-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#special_floating_point_values | Special float values '+inf', '-inf', 'nan' (as cast-from-string) |
| LEX-014 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#boolean_literals | TRUE / FALSE literals |
| LEX-015 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#null_literal | NULL literal |
| LEX-016 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#date_literals | DATE/TIMESTAMP/JSON typed literals |
| LEX-017 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#query_parameters | Query parameters @param_name (named only, no positional) |
| LEX-018 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#comments | Comments -- single-line and /* */ multi-line |
| LEX-019 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#reserved_keywords | Reserved keywords list (ALL, AND, ARRAY, ... WITH, WITHIN) |
| LEX-020 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#case_sensitivity | Case sensitivity: keywords/identifiers case-insensitive, strings case-sensitive |
| LEX-021 | https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#terminating_semicolons | Semicolon as statement separator in DDL batches |

---

## Expression forms (`expressions.md`) — 40 forms

| form_id | source_url | note |
|---------|-----------|------|
| EXPR-001 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#case_expressions | CASE expr WHEN val THEN result ... END (simple form) |
| EXPR-002 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#case_expressions | CASE WHEN condition THEN result ... END (searched form) |
| EXPR-003 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#cast_expressions | CAST(expr AS type) — error on failure |
| EXPR-004 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#safe_cast_expressions | SAFE_CAST(expr AS type) — NULL on failure |
| EXPR-005 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#extract_expressions | EXTRACT(part FROM datetime [AT TIME ZONE tz]) |
| EXPR-006 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#function_calls | Function call fn(arg, ...) |
| EXPR-007 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#function_calls | Named argument function call fn(name => val, ...) |
| EXPR-008 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls | Aggregate with DISTINCT: COUNT(DISTINCT col) |
| EXPR-009 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls | Aggregate with IGNORE NULLS / RESPECT NULLS |
| EXPR-010 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls | Aggregate with ORDER BY: ARRAY_AGG(x ORDER BY y) |
| EXPR-011 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls | Aggregate with LIMIT: ARRAY_AGG(x ORDER BY y LIMIT n) |
| EXPR-012 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#window_function_calls | Window fn OVER (PARTITION BY ... ORDER BY ... ROWS/RANGE BETWEEN ...) |
| EXPR-013 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator | Array subscript [OFFSET(n)] — 0-based |
| EXPR-014 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator | Array subscript [ORDINAL(n)] — 1-based |
| EXPR-015 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator | Array subscript [SAFE_OFFSET(n)] — NULL on OOB |
| EXPR-016 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator | Array subscript [SAFE_ORDINAL(n)] — NULL on OOB |
| EXPR-017 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_constructors | Array constructors: [1,2,3], ARRAY<T>[...], ARRAY(subquery) |
| EXPR-018 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#struct_constructors | STRUCT constructors: STRUCT(...), STRUCT<T>(...), (a, b, c) |
| EXPR-019 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#field_access_operator | Struct field access: struct_expr.field_name |
| EXPR-020 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#interval_expressions | INTERVAL integer PART for date/timestamp arithmetic |
| EXPR-021 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#conditional_expressions | IF, IFNULL, NULLIF, COALESCE, IIF conditional functions |
| EXPR-022 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#in_operators | IN (...), NOT IN (...), IN UNNEST(array) |
| EXPR-023 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#exists_expression | EXISTS (subquery) |
| EXPR-024 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#like_operator | LIKE / NOT LIKE pattern matching |
| EXPR-025 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#is_expressions | IS NULL / IS NOT NULL / IS TRUE / IS FALSE |
| EXPR-026 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#between_operator | BETWEEN lower AND upper / NOT BETWEEN |
| EXPR-027 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#arithmetic_operators | Arithmetic: +, -, *, /, unary -/+ |
| EXPR-028 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#comparison_operators | Comparison: =, !=, <>, <, <=, >, >= |
| EXPR-029 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#logical_operators | Logical: AND, OR, NOT |
| EXPR-030 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#bitwise_operators | Bitwise: &, |, ^, ~, <<, >> |
| EXPR-031 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#concatenation_operator | String/array concatenation: \|\| |
| EXPR-032 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#scalar_subqueries | Scalar subquery (SELECT expr FROM ...) |
| EXPR-033 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subqueries | ARRAY(SELECT ...) subquery |
| EXPR-034 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#quantified_subqueries | Quantified: expr op ALL/SOME/ANY (subquery) |
| EXPR-035 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#at_time_zone | AT TIME ZONE timezone for timestamp conversion |
| EXPR-036 | https://cloud.google.com/spanner/docs/commit-timestamp | PENDING_COMMIT_TIMESTAMP() — special commit timestamp function |
| EXPR-037 | https://cloud.google.com/spanner/docs/sequences | NEXT VALUE FOR seq / GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq) |
| EXPR-038 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#in_operators | IN UNNEST(array_param) for array parameter membership |
| EXPR-039 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#proto_constructors | NEW proto_type (field: val, ...) proto message constructor |
| EXPR-040 | https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#treat_expression | TREAT(expr AS proto_type) proto type conversion |

---

## Summary

| file | form count |
|------|-----------|
| ddl.md | 57 |
| dml.md | 10 |
| query.md | 26 |
| datatypes.md | 18 |
| hints.md | 13 |
| lexical.md | 21 |
| expressions.md | 40 |
| **TOTAL** | **185** |
