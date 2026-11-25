# BigQuery GoogleSQL truth1 — Master Index

version: bigquery-current  
corpus_date: 2026-06-02  
source_root: https://cloud.google.com/bigquery/docs/reference/standard-sql/

---

## query.md — 40 forms

| form_id | source_url | note |
|---------|-----------|------|
| QUERY-001 | .../query-syntax | top-level query_statement / query_expr grammar |
| QUERY-002 | .../query-syntax#select_statement | SELECT with all modifiers |
| QUERY-003 | .../query-syntax#select_except | SELECT * EXCEPT (col, ...) |
| QUERY-004 | .../query-syntax#select_replace | SELECT * REPLACE (expr AS col, ...) |
| QUERY-005 | .../query-syntax#select_distinct | SELECT DISTINCT |
| QUERY-006 | .../query-syntax#select_all | SELECT ALL |
| QUERY-007 | .../query-syntax#select_as_struct | SELECT AS STRUCT |
| QUERY-008 | .../query-syntax#select_as_value | SELECT AS VALUE |
| QUERY-009 | .../query-syntax#select_expression_star | SELECT expression.* |
| QUERY-010 | .../query-syntax#from_clause | FROM clause full grammar |
| QUERY-011 | .../query-syntax#for_system_time_as_of | FOR SYSTEM_TIME AS OF (time travel) |
| QUERY-012 | .../query-syntax#unnest_operator | UNNEST operator + WITH OFFSET |
| QUERY-013 | .../query-syntax#pivot_operator | PIVOT operator |
| QUERY-014 | .../query-syntax#unpivot_operator | UNPIVOT operator (single & multi-column) |
| QUERY-015 | .../query-syntax#tablesample_operator | TABLESAMPLE SYSTEM (n PERCENT) |
| QUERY-016 | .../query-syntax#match_recognize_clause | MATCH_RECOGNIZE clause |
| QUERY-017 | .../query-syntax#inner_join | [INNER] JOIN ... ON / USING |
| QUERY-018 | .../query-syntax#cross_join | CROSS JOIN / comma join |
| QUERY-019 | .../query-syntax#full_outer_join | FULL [OUTER] JOIN |
| QUERY-020 | .../query-syntax#left_outer_join | LEFT [OUTER] JOIN |
| QUERY-021 | .../query-syntax#right_outer_join | RIGHT [OUTER] JOIN |
| QUERY-022 | .../query-syntax#join_operation | join_operation full grammar |
| QUERY-023 | .../query-syntax#where_clause | WHERE clause |
| QUERY-024 | .../query-syntax#group_by_clause | GROUP BY full grammar |
| QUERY-025 | .../query-syntax#group_by_all | GROUP BY ALL |
| QUERY-026 | .../query-syntax#group_by_grouping_sets | GROUP BY GROUPING SETS |
| QUERY-027 | .../query-syntax#group_by_rollup | GROUP BY ROLLUP |
| QUERY-028 | .../query-syntax#group_by_cube | GROUP BY CUBE |
| QUERY-029 | .../query-syntax#having_clause | HAVING clause |
| QUERY-030 | .../query-syntax#order_by_clause | ORDER BY ... ASC/DESC NULLS FIRST/LAST |
| QUERY-031 | .../query-syntax#qualify_clause | QUALIFY clause (window filter) |
| QUERY-032 | .../query-syntax#window_clause | WINDOW clause (named windows) |
| QUERY-033 | .../query-syntax#set_operators | UNION/INTERSECT/EXCEPT ALL/DISTINCT + BY NAME/CORRESPONDING |
| QUERY-034 | .../query-syntax#limit_and_offset_clause | LIMIT count OFFSET skip |
| QUERY-035 | .../query-syntax#with_clause | WITH (non-recursive CTEs) |
| QUERY-036 | .../query-syntax#recursive_keyword | WITH RECURSIVE (recursive CTEs) |
| QUERY-037 | .../query-syntax#correlated_join_operation | correlated join / correlated subquery |
| QUERY-038 | .../query-syntax#field_path | field_path (implicit UNNEST in FROM) |
| QUERY-039 | .../query-syntax#aggregation_threshold_clause | SELECT WITH AGGREGATION_THRESHOLD |
| QUERY-040 | .../query-syntax#table_function_calls | table function call in FROM |

---

## ddl.md — 54 forms

| form_id | source_url | note |
|---------|-----------|------|
| DDL-001 | .../data-definition-language#create_schema_statement | CREATE SCHEMA |
| DDL-002 | .../data-definition-language#create_table_statement | CREATE [OR REPLACE] [TEMP] TABLE + AS SELECT |
| DDL-003 | .../data-definition-language#column | column_schema (column definition with NOT NULL, DEFAULT, GENERATED) |
| DDL-004 | .../data-definition-language#partition_expression | PARTITION BY (all variants) |
| DDL-005 | .../data-definition-language#clustering_column_list | CLUSTER BY |
| DDL-006 | .../data-definition-language#create_table_like_statement | CREATE TABLE LIKE |
| DDL-007 | .../data-definition-language#create_table_copy_statement | CREATE TABLE COPY |
| DDL-008 | .../data-definition-language#create_snapshot_table_statement | CREATE SNAPSHOT TABLE |
| DDL-009 | .../data-definition-language#create_table_clone_statement | CREATE TABLE CLONE |
| DDL-010 | .../data-definition-language#create_view_statement | CREATE VIEW |
| DDL-011 | .../data-definition-language#create_materialized_view_statement | CREATE MATERIALIZED VIEW |
| DDL-012 | .../data-definition-language#create_materialized_view_as_replica_of_statement | CREATE MATERIALIZED VIEW AS REPLICA OF |
| DDL-013 | .../data-definition-language#create_external_schema_statement | CREATE EXTERNAL SCHEMA |
| DDL-014 | .../data-definition-language#create_external_table_statement | CREATE EXTERNAL TABLE |
| DDL-015 | .../data-definition-language#create_function_statement | CREATE FUNCTION (SQL) |
| DDL-016 | .../data-definition-language#create_function_statement | CREATE FUNCTION (JavaScript) |
| DDL-017 | .../data-definition-language#create_function_statement | CREATE FUNCTION (Python / remote) |
| DDL-018 | .../data-definition-language#create_aggregate_function_statement_sql | CREATE AGGREGATE FUNCTION |
| DDL-019 | .../data-definition-language#create_table_function_statement | CREATE TABLE FUNCTION (TVF) |
| DDL-020 | .../data-definition-language#create_procedure_statement | CREATE PROCEDURE (SQL + Spark variants) |
| DDL-021 | .../data-definition-language#create_row_access_policy_statement | CREATE ROW ACCESS POLICY |
| DDL-022 | .../data-definition-language#create_search_index_statement | CREATE SEARCH INDEX |
| DDL-023 | .../data-definition-language#create_vector_index_statement | CREATE VECTOR INDEX |
| DDL-024 | .../data-definition-language#create_capacity_statement | CREATE CAPACITY |
| DDL-025 | .../data-definition-language#create_reservation_statement | CREATE RESERVATION |
| DDL-026 | .../data-definition-language#create_assignment_statement | CREATE ASSIGNMENT |
| DDL-027 | .../data-definition-language#alter_table_set_options_statement | ALTER TABLE SET OPTIONS |
| DDL-028 | .../data-definition-language#alter_table_add_column_statement | ALTER TABLE ADD COLUMN |
| DDL-029 | .../data-definition-language#alter_table_add_foreign_key_statement | ALTER TABLE ADD FOREIGN KEY |
| DDL-030 | .../data-definition-language#alter_table_add_primary_key_statement | ALTER TABLE ADD PRIMARY KEY |
| DDL-031 | .../data-definition-language#alter_table_rename_to_statement | ALTER TABLE RENAME TO |
| DDL-032 | .../data-definition-language#alter_table_rename_column_statement | ALTER TABLE RENAME COLUMN |
| DDL-033 | .../data-definition-language#alter_table_drop_column_statement | ALTER TABLE DROP COLUMN |
| DDL-034 | .../data-definition-language#alter_table_drop_constraint_statement | ALTER TABLE DROP CONSTRAINT / DROP PRIMARY KEY |
| DDL-035 | .../data-definition-language#alter_table_set_default_collate_statement | ALTER TABLE SET DEFAULT COLLATE |
| DDL-036 | .../data-definition-language#alter_column_set_options_statement | ALTER COLUMN (SET OPTIONS / DROP NOT NULL / SET DATA TYPE / SET DEFAULT / DROP DEFAULT) |
| DDL-037 | .../data-definition-language#alter_schema_set_default_collate_statement | ALTER SCHEMA SET DEFAULT COLLATE |
| DDL-038 | .../data-definition-language#alter_schema_set_options_statement | ALTER SCHEMA SET OPTIONS |
| DDL-039 | .../data-definition-language#alter_view_set_options_statement | ALTER VIEW SET OPTIONS |
| DDL-040 | .../data-definition-language#alter_materialized_view_set_options_statement | ALTER MATERIALIZED VIEW SET OPTIONS |
| DDL-041 | .../data-definition-language#alter_vector_index_rebuild_statement | ALTER VECTOR INDEX REBUILD |
| DDL-042 | .../data-definition-language#drop_schema_statement | DROP [EXTERNAL] SCHEMA [CASCADE|RESTRICT] |
| DDL-043 | .../data-definition-language#undrop_schema_statement | UNDROP SCHEMA |
| DDL-044 | .../data-definition-language#drop_table_statement | DROP TABLE |
| DDL-045 | .../data-definition-language#drop_snapshot_table_statement | DROP SNAPSHOT TABLE |
| DDL-046 | .../data-definition-language#drop_external_table_statement | DROP EXTERNAL TABLE |
| DDL-047 | .../data-definition-language#drop_view_statement | DROP VIEW |
| DDL-048 | .../data-definition-language#drop_materialized_view_statement | DROP MATERIALIZED VIEW |
| DDL-049 | .../data-definition-language#drop_function_statement | DROP FUNCTION |
| DDL-050 | .../data-definition-language#drop_table_function_statement | DROP TABLE FUNCTION |
| DDL-051 | .../data-definition-language#drop_procedure_statement | DROP PROCEDURE |
| DDL-052 | .../data-definition-language#drop_row_access_policy_statement | DROP ROW ACCESS POLICY / DROP ALL ROW ACCESS POLICIES |
| DDL-053 | .../data-definition-language#drop_capacity_statement | DROP CAPACITY / RESERVATION / ASSIGNMENT |
| DDL-054 | .../data-definition-language#drop_search_index_statement | DROP SEARCH INDEX / DROP VECTOR INDEX |

---

## dml.md — 5 forms

| form_id | source_url | note |
|---------|-----------|------|
| DML-001 | .../dml-syntax#insert_statement | INSERT INTO ... VALUES / SELECT + DEFAULT |
| DML-002 | .../dml-syntax#delete_statement | DELETE [FROM] ... WHERE |
| DML-003 | .../dml-syntax#truncate_table_statement | TRUNCATE TABLE |
| DML-004 | .../dml-syntax#update_statement | UPDATE ... SET ... [FROM] WHERE |
| DML-005 | .../dml-syntax#merge_statement | MERGE ... USING ... ON ... WHEN [NOT] MATCHED |

---

## datatypes.md — 6 forms

| form_id | source_url | note |
|---------|-----------|------|
| DT-001 | .../data-types | All scalar types: INT64, NUMERIC, BIGNUMERIC, FLOAT64, BOOL, STRING, BYTES, DATE, DATETIME, TIME, TIMESTAMP, INTERVAL, JSON, GEOGRAPHY |
| DT-002 | .../data-types#parameterized_data_types | Parameterized types: STRING(L), BYTES(L), NUMERIC(P,S), BIGNUMERIC(P,S) |
| DT-003 | .../data-types#array_type | ARRAY<T> declaration and construction |
| DT-004 | .../data-types#struct_type | STRUCT<...> declaration and all construction syntaxes (tuple / typeless / typed) |
| DT-005 | .../data-types#range_type | RANGE<DATE|DATETIME|TIMESTAMP> and range literals |
| DT-006 | .../data-types#data_type_list | Type aliases (INT=INT64, DECIMAL=NUMERIC, etc.) |

---

## lexical.md — 20 forms

| form_id | source_url | note |
|---------|-----------|------|
| LEX-001 | .../lexical#identifiers | Unquoted identifier rules |
| LEX-002 | .../lexical#identifiers | Backtick-quoted identifier rules |
| LEX-003 | .../lexical#path_expressions | Path expressions (dot-separated with / : - allowed) |
| LEX-004 | .../lexical#table_names | Table name with dashes (project ID only) |
| LEX-005 | .../lexical#string_and_bytes_literals | All string/bytes literal forms: quoted, triple-quoted, raw (r""), bytes (b""), raw bytes (rb""), concatenation |
| LEX-006 | .../lexical#escape_sequences | All escape sequences \a \b \n \t \x \u \U \ooo etc |
| LEX-007 | .../lexical#integer_literals | Decimal and hex integer literals |
| LEX-008 | .../lexical#numeric_literals | NUMERIC/DECIMAL and BIGNUMERIC/BIGDECIMAL literals |
| LEX-009 | .../lexical#floating_point_literals | Floating point literal syntax |
| LEX-010 | .../lexical#array_literals | Array literal syntax |
| LEX-011 | .../lexical#struct_literals | Struct literal syntax (tuple/typeless/typed) |
| LEX-012 | .../lexical#date_literals | DATE/TIME/DATETIME/TIMESTAMP typed literals |
| LEX-013 | .../lexical#interval_literals | INTERVAL literals (single part and datetime range) |
| LEX-014 | .../lexical#range_literals | RANGE<T> literal syntax |
| LEX-015 | .../lexical#json_literals | JSON literal syntax |
| LEX-016 | .../lexical#named_query_parameters | Named query parameters @name |
| LEX-017 | .../lexical#positional_query_parameters | Positional query parameters ? |
| LEX-018 | .../lexical#comments | Comments: # -- /* */ |
| LEX-019 | .../lexical#reserved_keywords | Full reserved keyword list (80 keywords) |
| LEX-020 | .../lexical#terminating_semicolons | Trailing commas and semicolons |

---

## scripting.md — 19 forms

| form_id | source_url | note |
|---------|-----------|------|
| SCRIPT-001 | .../procedural-language#declare | DECLARE variable [DEFAULT expr] |
| SCRIPT-002 | .../procedural-language#set | SET single and tuple assignment |
| SCRIPT-003 | .../procedural-language#execute_immediate | EXECUTE IMMEDIATE sql [INTO] [USING] |
| SCRIPT-004 | .../procedural-language#beginend | BEGIN...END block |
| SCRIPT-005 | .../procedural-language#beginexceptionend | BEGIN...EXCEPTION WHEN ERROR THEN...END |
| SCRIPT-006 | .../procedural-language#case | CASE (searched) statement |
| SCRIPT-007 | .../procedural-language#case_search_expression | CASE expr WHEN ... END CASE (simple) |
| SCRIPT-008 | .../procedural-language#if | IF/ELSEIF/ELSE/END IF |
| SCRIPT-009 | .../procedural-language#labels | Label syntax for blocks and loops |
| SCRIPT-010 | .../procedural-language#loop | LOOP...END LOOP |
| SCRIPT-011 | .../procedural-language#repeat | REPEAT...UNTIL...END REPEAT |
| SCRIPT-012 | .../procedural-language#while | WHILE...DO...END WHILE |
| SCRIPT-013 | .../procedural-language#break | BREAK / LEAVE |
| SCRIPT-014 | .../procedural-language#continue | CONTINUE / ITERATE |
| SCRIPT-015 | .../procedural-language#for_in | FOR...IN...DO...END FOR |
| SCRIPT-016 | .../procedural-language#begin_transaction | BEGIN/COMMIT/ROLLBACK TRANSACTION |
| SCRIPT-017 | .../procedural-language#raise | RAISE [USING MESSAGE = msg] |
| SCRIPT-018 | .../procedural-language#return | RETURN |
| SCRIPT-019 | .../procedural-language#call | CALL procedure_name(args...) |

---

## dcl.md — 2 forms

| form_id | source_url | note |
|---------|-----------|------|
| DCL-001 | .../data-control-language#grant_statement | GRANT role_list ON resource_type TO user_list |
| DCL-002 | .../data-control-language#revoke_statement | REVOKE role_list ON resource_type FROM user_list |

---

## expressions.md — 16 forms

| form_id | source_url | note |
|---------|-----------|------|
| EXPR-001 | .../conditional_expressions#case_expr | CASE WHEN ... THEN ... END (searched) |
| EXPR-002 | .../conditional_expressions#case_expr | CASE expr WHEN ... END (simple) |
| EXPR-003 | .../conversion_functions#cast | CAST(expr AS type [FORMAT fmt]) / SAFE_CAST |
| EXPR-004 | .../date_functions#extract | EXTRACT(part FROM datetime [AT TIME ZONE]) |
| EXPR-005 | .../aggregate_functions | Aggregate call: [DISTINCT], IGNORE/RESPECT NULLS, ORDER BY, LIMIT inside aggregate |
| EXPR-006 | .../window-function-calls | Window function OVER (window_spec) full grammar with ROWS/RANGE frame |
| EXPR-007 | .../arrays#accessing_array_elements | Array element access: [OFFSET(n)], [ORDINAL(n)], [SAFE_OFFSET(n)], [SAFE_ORDINAL(n)] |
| EXPR-008 | .../data-types#struct_type | Struct field access: struct_expr.field_name |
| EXPR-009 | .../timestamp_functions#at_time_zone | AT TIME ZONE in expressions |
| EXPR-010 | .../expression_subqueries#scalar_subqueries | Scalar subquery: (subquery) |
| EXPR-011 | .../expression_subqueries#array_subqueries | Array subquery: ARRAY(subquery) |
| EXPR-012 | .../expression_subqueries#in_subqueries | IN / NOT IN subquery |
| EXPR-013 | .../expression_subqueries#exists_subqueries | EXISTS(subquery) |
| EXPR-014 | .../expression_subqueries#table_subqueries | Table subquery: FROM (subquery) AS alias |
| EXPR-015 | .../lexical#interval_literals | INTERVAL expression in arithmetic contexts |
| EXPR-016 | .../operators | Full operator set and precedence |

---

## other.md — 3 forms

| form_id | source_url | note |
|---------|-----------|------|
| OTHER-001 | .../other-statements#export_data_statement | EXPORT DATA ... OPTIONS(...) AS query |
| OTHER-002 | .../debugging-statements#assert | ASSERT expression [AS description] |
| OTHER-003 | .../load-data | LOAD DATA {OVERWRITE\|INTO} table FROM FILES (options) [WITH PARTITION COLUMNS] — page was 404 at corpus time; syntax from known docs |

---

## Totals

| file | forms |
|------|-------|
| query.md | 40 |
| ddl.md | 54 |
| dml.md | 5 |
| datatypes.md | 6 |
| lexical.md | 20 |
| scripting.md | 19 |
| dcl.md | 2 |
| expressions.md | 16 |
| other.md | 3 |
| **total** | **165** |

---

## Coverage notes

- **MATCH_RECOGNIZE**: documented as BigQuery-specific clause (QUERY-016); not standard SQL.
- **Differential privacy clause**: documented but omitted from this corpus as it is a BigQuery-specific extension with minimal parser surface (just `WITH differential_privacy_clause` inside SELECT).
- **AGGREGATION_THRESHOLD clause**: captured as QUERY-039.
- **LOAD DATA**: the dedicated reference page (cloud.google.com/bigquery/docs/reference/standard-sql/load-data) returned 404 at corpus build time (2026-06-02). Syntax in OTHER-003 is reconstructed from cross-references in other pages. Revisit this entry when the page becomes available.
- **ALTER ORGANIZATION / ALTER PROJECT / ALTER BI_CAPACITY**: These exist in the DDL page but are very narrow admin-only forms; omitted to keep the corpus focused on SQL syntax forms. Add them if parser coverage is needed.
- **CREATE DATA_POLICY / CREATE CONNECTION / ALTER DATA_POLICY / ALTER CONNECTION / DROP DATA_POLICY / DROP CONNECTION**: Admin-plane statements omitted similarly.
- **GRAPH_TABLE operator**: Mentioned in FROM clause grammar (QUERY-010) but the full GQL property-graph syntax is a separate dialect not in scope of this corpus.
