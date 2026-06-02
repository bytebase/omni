# GoogleSQL Legacy ANTLR4 Grammar Catalog (`antlr_rules.md`)

> **Status:** truth2 — a **HINT corpus**, NOT an oracle. This document catalogs the legacy
> ANTLR4 grammar that the new hand-written omni parser must eventually match for coverage.
> It defines the full coverage target; it does not by itself validate the new parser.
>
> **Sources (authoritative legacy grammar, read in full):**
> - `/Users/h3n4l/OpenSource/parser/googlesql/GoogleSQLParser.g4` (2798 lines)
> - `/Users/h3n4l/OpenSource/parser/googlesql/GoogleSQLLexer.g4` (497 lines)
>
> **Secondary cross-reference (HINT only, NOT the target):**
> - `/Users/h3n4l/OpenSource/bq-parser/BigQueryParser.g4` (343 lines, `BigQueryParser`)
> - `/Users/h3n4l/OpenSource/bq-parser/BigQueryLexer.g4` (192 lines, `BigQueryLexer`)
>
> **Consumed in bytebase as** `github.com/bytebase/parser/googlesql`, serving **BOTH BigQuery and
> Spanner** dialects from one grammar. The grammar is a hand-port of Google ZetaSQL's
> `bison_parser.y` / `flex_tokenizer.l` to ANTLR4 (see inline source links in the `.g4`),
> with embedded Go semantic predicates (`p.NotifyErrorListeners(...)`) used to emit ZetaSQL-style
> syntax errors for known-bad constructs.

---

## Headline numbers

| Metric | Value |
| --- | --- |
| Parser rules (unique named productions) | **625** |
| Lexer tokens (keyword `*_SYMBOL`) | 307 word-keywords |
| Lexer tokens (operator/punctuation `*_SYMBOL` / `*_OPERATOR`) | 43 (16 of them `*_OPERATOR`) |
| Reserved keywords (NOT usable as bare identifier) | **97** (= 307 word-keywords − 211 non-reserved; verified by set-difference) |
| Non-reserved keywords (usable as identifier via `common_keyword_as_identifier` + `SIMPLE`) | **211** |
| Top-level statement kinds (`sql_statement_body` alternatives) | **53** SQL bodies + **12** procedural-script statement forms (incl. BEGIN…END block) |
| `opt_*` optional-helper rules | 101 |
| GQL / graph rules (`graph_*`, `gql_*`) | 40 |
| TODO / FIXME / XXX markers | 10 |

---

## 1. Entry rules

| Rule | Parses |
| --- | --- |
| `root` | `stmts EOF` — the grammar start symbol. A complete input is a statement list followed by EOF. |
| `stmts` | One-or-more `unterminated_sql_statement` separated by `;` (`SEMI_SYMBOL`), with an optional trailing `;`. This is the multi-statement list. |
| `unterminated_sql_statement` | A single statement: an optional `statement_level_hint` followed by `sql_statement_body`. Also has two error-only alternatives that reject `DEFINE MACRO` (cannot be composed from expansions; hints not allowed on it). |
| `sql_statement_body` | The **53-way dispatch** to every SQL statement kind (query, DDL, DML, DCL, utility, transaction, batch, GQL, etc.). This is the real "single statement" rule. See §4 for the full enumeration. |
| `query_statement` | `query` — the top-level query entry (BigQuery query-syntax root). Delegates to `query` → `query_without_pipe_operators`. |
| `query` | `query_without_pipe_operators` — the SELECT/set-operation/CTE query root (see §5). |
| `statement_list` | `unterminated_non_empty_statement_list SEMI_SYMBOL` — the procedural body used inside `BEGIN…END`, `LOOP`, `WHILE`, `IF`, `CASE`, etc. Mixes SQL statements and script statements. |
| `unterminated_statement` | `unterminated_sql_statement \| unterminated_script_statement` — one element of a procedural body. |

**Two statement "universes":**
- **SQL statements** (`unterminated_sql_statement` → `sql_statement_body`): the normal top-level surface.
- **Procedural / scripting statements** (`unterminated_script_statement`, `unterminated_unlabeled_script_statement`): only legal inside `statement_list` (i.e. inside a stored-procedure/script `BEGIN…END` block). These add `IF`, `CASE`, `WHILE`, `LOOP`, `REPEAT`, `FOR…IN`, `DECLARE`, `BREAK/LEAVE`, `CONTINUE/ITERATE`, `RETURN`, `RAISE`, and labels.

---

## 2. Full rule catalog

> Categories: **entry / statement / ddl / dml / query / expression / identifier / datatype / literal / clause / helper**.
> Trivial helper/`opt_*` rules are grouped into compact rows; every statement / ddl / dml / query / expression
> rule appears individually with its real rule name. All 625 rules are accounted for.

### 2.1 Entry & top-level dispatch

| rule | purpose | category |
| --- | --- | --- |
| `root` | start symbol: `stmts EOF` | entry |
| `stmts` | `;`-separated statement list | entry |
| `unterminated_sql_statement` | one statement (+ optional hint); rejects composed `DEFINE MACRO` | entry |
| `sql_statement_body` | 53-way dispatch to all SQL statement kinds | entry |
| `query_statement` | top-level query | entry |
| `statement_list` | procedural body (terminated) | entry |
| `unterminated_non_empty_statement_list` | `;`-separated procedural statements | helper |
| `unterminated_statement` | SQL-or-script statement element | entry |
| `statement_level_hint` | `hint` applied at statement scope | clause |

### 2.2 DDL — CREATE statements

| rule | purpose | category |
| --- | --- | --- |
| `create_table_statement` | `CREATE … TABLE` with elements, Spanner options, LIKE/CLONE/COPY, default-collate, PARTITION BY, CLUSTER BY, TTL row-deletion policy, WITH CONNECTION, OPTIONS, `AS query` | ddl |
| `create_table_function_statement` | `CREATE … TABLE FUNCTION` (TVF) with params, RETURNS, language/options, `AS query`-or-string; predicate requires `(` | ddl |
| `create_external_table_statement` | `CREATE … EXTERNAL TABLE` with elements, LIKE, default-collate, external WITH clauses, OPTIONS | ddl |
| `create_external_table_function_statement` | error-only: emits "CREATE EXTERNAL TABLE FUNCTION is not supported" | ddl |
| `create_view_statement` | `CREATE [MATERIALIZED\|APPROX] [RECURSIVE] VIEW` (+ column-options, SQL SECURITY, PARTITION/CLUSTER for materialized, OPTIONS, `AS query`/replica-source) | ddl |
| `create_index_statement` | `CREATE [UNIQUE] [NULL_FILTERED] [SEARCH\|VECTOR] INDEX … ON …` (+ unnest exprs, ordering+options, STORING, partition/options/Spanner-interleave suffix) | ddl |
| `create_schema_statement` | `CREATE SCHEMA` (+ default-collate, OPTIONS) | ddl |
| `create_external_schema_statement` | `CREATE … EXTERNAL SCHEMA` (+ WITH CONNECTION, OPTIONS) | ddl |
| `create_function_statement` | `CREATE [AGGREGATE] FUNCTION` (+ returns, SQL SECURITY, determinism, LANGUAGE/REMOTE WITH CONNECTION, options/body) | ddl |
| `create_procedure_statement` | `CREATE PROCEDURE` (+ params, EXTERNAL SECURITY, WITH CONNECTION, OPTIONS, BEGIN…END or LANGUAGE…AS code) | ddl |
| `create_database_statement` | `CREATE DATABASE … [OPTIONS]` | ddl |
| `create_connection_statement` | `CREATE … CONNECTION … [OPTIONS]` | ddl |
| `create_constant_statement` | `CREATE … CONSTANT … = expression` | ddl |
| `create_model_statement` | `CREATE … MODEL` (+ INPUT/OUTPUT, TRANSFORM, REMOTE WITH CONNECTION, OPTIONS, `AS query`/aliased-query-list) | ddl |
| `create_snapshot_statement` | `CREATE … SNAPSHOT TABLE\|<obj> … CLONE <source> [OPTIONS]` | ddl |
| `create_property_graph_statement` | `CREATE … PROPERTY GRAPH … NODE TABLES <list> [EDGE TABLES <list>]` | ddl |
| `create_row_access_policy_statement` | `CREATE … ROW ACCESS POLICY … ON … [GRANT TO] FILTER USING (...)` | ddl |
| `create_privilege_restriction_statement` | `CREATE … PRIVILEGE RESTRICTION ON <privs> ON <type> <path> [RESTRICT TO]` | ddl |
| `create_entity_statement` | `CREATE … <generic_entity_type> … [OPTIONS] [AS <body>]` (generic/extensible object) | ddl |

### 2.3 DDL — ALTER / DROP / RENAME / UNDROP / TRUNCATE

| rule | purpose | category |
| --- | --- | --- |
| `alter_statement` | `ALTER` for TABLE/[TABLE]FUNCTION, schema-object, generic entity, PRIVILEGE RESTRICTION, ROW ACCESS POLICY, ALL ROW ACCESS POLICIES | ddl |
| `drop_statement` | `DROP` for PRIVILEGE RESTRICTION, ROW ACCESS POLICY, INDEX, TABLE/[TABLE]FUNCTION, SNAPSHOT TABLE, generic entity, generic schema-object (+ RESTRICT/CASCADE) | ddl |
| `drop_all_row_access_policies_statement` | `DROP ALL ROW ACCESS POLICIES ON <path>` | ddl |
| `rename_statement` | `RENAME <object-kind> <path> TO <path>` | ddl |
| `undrop_statement` | `UNDROP <schema_object_kind> [IF NOT EXISTS] <path> [AT SYSTEM TIME] [OPTIONS]` | ddl |
| `truncate_statement` | `TRUNCATE TABLE <path> [WHERE expr]` | ddl |
| `define_table_statement` | `DEFINE TABLE <path> [OPTIONS]` | ddl |

### 2.4 DDL — table-element / column / constraint structure

| rule | purpose | category |
| --- | --- | --- |
| `table_element_list` | parenthesized `( table_element , … )` | clause |
| `table_element` | `table_column_definition \| table_constraint_definition` | helper |
| `table_column_definition` | `identifier table_column_schema [attributes] [OPTIONS]` | ddl |
| `table_column_schema` | column schema + collate + column-info, or generated-column-info | datatype |
| `column_schema_inner` | raw column schema + optional type parameters | datatype |
| `raw_column_schema_inner` | simple/array/struct/range column schema | datatype |
| `simple_column_schema_inner` | `path_expression \| INTERVAL` | datatype |
| `array_column_schema_inner` | `ARRAY< field_schema >` | datatype |
| `struct_column_schema_inner` | `STRUCT< struct_column_field,… >` | datatype |
| `struct_column_field` | struct field (schema+collate+attrs, or `identifier field_schema`) | datatype |
| `range_column_schema_inner` | `RANGE< field_schema >` | datatype |
| `field_schema` | column-schema-inner + collate + attrs + OPTIONS | datatype |
| `opt_field_attributes` | `not_null_column_attribute` | helper |
| `not_null_column_attribute` | `NOT NULL` | clause |
| `column_attributes` | `column_attribute+` + optional enforcement | clause |
| `column_attribute` | primary-key / foreign-key / hidden / not-null attribute | clause |
| `primary_key_column_attribute` | `PRIMARY KEY` (inline) | clause |
| `foreign_key_column_attribute` | `[CONSTRAINT id] foreign_key_reference` | clause |
| `hidden_column_attribute` | `HIDDEN` | clause |
| `opt_constraint_identity` | `CONSTRAINT identifier` | helper |
| `table_constraint_definition` | primary-key-spec / table-constraint-spec / `id id table_constraint_spec` | ddl |
| `primary_key_or_table_constraint_spec` | `primary_key_spec \| table_constraint_spec` | helper |
| `primary_key_spec` | `PRIMARY KEY ( elements ) [enforcement] [OPTIONS]` | ddl |
| `primary_key_element_list` | parenthesized PK element list | clause |
| `primary_key_element` | `identifier [ASC\|DESC] [null_order]` | helper |
| `table_constraint_spec` | `CHECK(expr)` or `FOREIGN KEY (cols) <reference>` (+ enforcement/options) | ddl |
| `foreign_key_reference` | `REFERENCES path (cols) [MATCH] [actions]` | clause |
| `opt_foreign_key_action` | ON UPDATE/ON DELETE action pair | helper |
| `foreign_key_on_update` | `ON UPDATE <action>` | clause |
| `foreign_key_on_delete` | `ON DELETE <action>` | clause |
| `foreign_key_action` | `NO ACTION \| RESTRICT \| CASCADE \| SET NULL` | helper |
| `opt_foreign_key_match` | `MATCH <mode>` | helper |
| `foreign_key_match_mode` | `SIMPLE \| FULL \| NOT DISTINCT` | helper |
| `constraint_enforcement` | `[NOT] ENFORCED` | clause |
| `column_list` | parenthesized `( identifier , … )` | clause |
| `column_position` | `PRECEDING id \| FOLLOWING id` (ADD COLUMN position) | clause |
| `fill_using_expression` | `FILL USING expr` (ADD COLUMN backfill) | clause |
| `default_column_info` | `DEFAULT expr` | clause |
| `opt_column_info` | generated-or-default column info (+ predicate rejecting both) | helper |
| `invalid_generated_column` / `invalid_default_column` | predicate markers for the "both DEFAULT and GENERATED" error | helper |
| `generated_column_info` | `<generated_mode> (expr) [stored_mode]` or identity column | clause |
| `generated_mode` | `GENERATED [ALWAYS\|BY DEFAULT] AS \| AS` | helper |
| `stored_mode` | `STORED [VOLATILE]` | helper |
| `identity_column_info` | `IDENTITY( [START WITH][INCREMENT BY][MAXVALUE][MINVALUE][CYCLE] )` | clause |
| `opt_start_with` / `opt_increment_by` / `opt_maxvalue` / `opt_minvalue` / `opt_cycle` | identity-column sub-clauses | helper |
| `signed_numeric_literal` | int/numeric/bignumeric/float, optionally `-` prefixed | literal |

### 2.5 DDL — ALTER actions & generic-entity machinery

| rule | purpose | category |
| --- | --- | --- |
| `alter_action_list` | `,`-separated `alter_action` | clause |
| `alter_action` | the big ALTER-action union: SET OPTIONS / SET AS body / ADD constraint\|PK\|column\|row-deletion-policy / DROP constraint\|PK\|column\|row-deletion-policy / ALTER constraint\|column / RENAME column\|TO / SET DEFAULT collate / generic-sub-entity add/drop/alter / Spanner column actions | ddl |
| `spanner_alter_column_action` | Spanner `ALTER COLUMN <schema> [NOT NULL] [generated/default] [OPTIONS]` | ddl |
| `spanner_set_on_delete_action` | Spanner `SET ON DELETE <foreign_key_action>` | ddl |
| `spanner_generated_or_default` | `AS (expr) STORED` (Spanner) | clause |
| `row_access_policy_alter_action_list` | `,`-separated row-access-policy alter actions | clause |
| `row_access_policy_alter_action` | GRANT TO / FILTER USING / REVOKE FROM (...)/ALL / RENAME TO | clause |
| `generic_entity_type` / `generic_entity_type_unchecked` | extensible entity type token (`IDENTIFIER \| PROJECT`) | identifier |
| `generic_sub_entity_type` / `sub_entity_type_identifier` | sub-entity type (`IDENTIFIER \| REPLICA`) | identifier |
| `generic_entity_body` | `json_literal \| string_literal` | helper |
| `schema_object_kind` | the object-kind keyword set (AGGREGATE FUNCTION, APPROX VIEW, CONNECTION, CONSTANT, DATABASE, EXTERNAL TABLE/SCHEMA, FUNCTION, INDEX, MATERIALIZED VIEW, MODEL, PROCEDURE, SCHEMA, VIEW) | helper |
| `table_or_table_function` | `TABLE [FUNCTION]` | helper |

### 2.6 DDL — create-statement option/sub-clauses

| rule | purpose | category |
| --- | --- | --- |
| `opt_or_replace` | `OR REPLACE` | helper |
| `opt_create_scope` | `TEMP \| TEMPORARY \| PUBLIC \| PRIVATE` | helper |
| `opt_if_not_exists` | `IF NOT EXISTS` | helper |
| `opt_if_exists` | `IF EXISTS` | helper |
| `opt_drop_mode` | `RESTRICT \| CASCADE` | helper |
| `opt_options_list` | `OPTIONS <options_list>` | clause |
| `opt_default_collate_clause` | `DEFAULT <collate_clause>` | clause |
| `opt_like_path_expression` | `LIKE <path>` | clause |
| `opt_clone_table` / `opt_copy_table` | `CLONE`/`COPY <data source>` | clause |
| `copy_data_source` / `clone_data_source` | `<path> [AT SYSTEM TIME] [WHERE]` | clause |
| `clone_data_source_list` | `UNION ALL`-joined clone sources | clause |
| `opt_ttl_clause` | `ROW DELETION POLICY (expr)` | clause |
| `opt_spanner_table_options` | Spanner `PRIMARY KEY` + INTERLEAVE-in-parent | clause |
| `opt_spanner_interleave_in_parent_clause` | `, INTERLEAVE IN PARENT <path> <on-delete>` | clause |
| `spanner_primary_key` | Spanner `PRIMARY KEY <elements>` | clause |
| `opt_spanner_null_filtered` | `NULL_FILTERED` (Spanner index) | helper |
| `spanner_index_interleave_clause` | `, INTERLEAVE IN <path>` (Spanner index) | clause |
| `opt_create_index_statement_suffix` | partition/options/Spanner-interleave tail of CREATE INDEX | clause |
| `index_type` | `SEARCH \| VECTOR` | helper |
| `index_order_by_and_options` | parenthesized column-ordering list or `index_all_columns` | clause |
| `index_all_columns` | `( ALL COLUMNS [WITH COLUMN OPTIONS …] )` | clause |
| `opt_with_column_options` | `WITH COLUMN OPTIONS <all_column_column_options>` | clause |
| `all_column_column_options` | parenthesized column-ordering+options list | clause |
| `column_ordering_and_options_expr` | `expr [collate] [ASC/DESC] [null_order] [OPTIONS]` | clause |
| `index_storing_list` | `STORING <expr-list>` | clause |
| `index_storing_expression_list` | parenthesized expression list | clause |
| `index_unnest_expression_list` | `unnest_expression_with_opt_alias_and_offset+` | clause |
| `unnest_expression_with_opt_alias_and_offset` | `unnest_expression [alias] [WITH OFFSET]` | clause |
| `on_path_expression` | `ON <path>` | clause |
| `column_with_options_list` / `column_with_options` | view column list with per-column OPTIONS | clause |
| `query_or_replica_source` | `query \| REPLICA OF <path>` (materialized view source) | clause |
| `opt_sql_security_clause` | `SQL SECURITY INVOKER\|DEFINER` | clause |
| `sql_security_clause_kind` | `INVOKER \| DEFINER` | helper |
| `opt_external_security_clause` | `EXTERNAL SECURITY INVOKER\|DEFINER` | clause |
| `external_security_clause_kind` | `INVOKER \| DEFINER` | helper |
| `opt_determinism_level` | `DETERMINISTIC \| NOT DETERMINISTIC \| IMMUTABLE \| STABLE \| VOLATILE` | helper |
| `opt_aggregate` | `AGGREGATE` | helper |
| `opt_returns` / `opt_function_returns` | `RETURNS <type_or_tvf_schema>` | clause |
| `opt_input_output_clause` | model `INPUT <elems> OUTPUT <elems>` | clause |
| `opt_transform_clause` | model `TRANSFORM (select_list)` | clause |
| `opt_as_query_or_aliased_query_list` | model `AS query` or `AS (aliased_query_list)` | clause |
| `aliased_query_list` | `,`-separated `aliased_query` | clause |
| `as_query` | `AS query` | clause |
| `opt_as_query_or_string` | `AS query \| AS string` (TVF body) | clause |
| `opt_generic_entity_body` | `AS <generic_entity_body>` | helper |
| `opt_edge_table_clause` | `EDGE TABLES <element_table_list>` | clause |
| `element_table_list` | parenthesized element-table definitions | clause |
| `element_table_definition` | property-graph node/edge table def (key, source/dest node, labels) | clause |
| `opt_label_and_properties_clause` | properties-clause or label-and-properties list | helper |
| `label_and_properties_list` / `label_and_properties` | `[DEFAULT] LABEL id [properties]` | clause |
| `properties_clause` | `NO PROPERTIES \| PROPERTIES ARE ALL COLUMNS [EXCEPT] \| PROPERTIES(derived…)` | clause |
| `properties_all_columns` | `PROPERTIES [ARE] ALL COLUMNS` | helper |
| `derived_property_list` / `derived_property` | property-graph derived properties | clause |
| `opt_except_column_list` | `EXCEPT <column_list>` | helper |
| `opt_key_clause` | `KEY <column_list>` | helper |
| `opt_source_node_table_clause` | `SOURCE KEY (cols) REFERENCES id (cols)` | clause |
| `opt_dest_node_table_clause` | `DESTINATION KEY (cols) REFERENCES id (cols)` | clause |
| `filter_using_clause` | `[FILTER] USING (expr)` (row access policy) | clause |
| `create_row_access_policy_grant_to_clause` | `grant_to_clause \| TO grantee_list` | helper |
| `restrict_to_clause` | `RESTRICT TO <possibly-empty grantee list>` | clause |
| `possibly_empty_grantee_list` | parenthesized optional string/param list | clause |
| `with_connection_clause` | `WITH <connection_clause>` | clause |
| `remote_with_connection_clause` | `REMOTE [WITH CONNECTION] \| WITH CONNECTION` | clause |
| `opt_language_or_remote_with_connection` | `LANGUAGE id [REMOTE WITH] \| REMOTE WITH [LANGUAGE]` | clause |
| `language` | `LANGUAGE identifier` | clause |
| `unordered_language_options` | language+options in either order (TVF) | clause |
| `unordered_options_body` | options+body in either order (function) | clause |
| `as_sql_function_body_or_string` | `AS <sql_function_body> \| AS string` | clause |
| `sql_function_body` | `(expr)`; error alt rejects `(SELECT …)` | clause |
| `begin_end_block_or_language_as_code` | `begin_end_block \| LANGUAGE id [AS code]` | clause |
| `opt_as_code` | `AS string` | helper |
| `function_declaration` | `path_expression function_parameters` | helper |
| `function_parameters` / `opt_function_parameters` | parenthesized `function_parameter` list | clause |
| `function_parameter` | `[id] type_or_tvf_schema [AS alias] [DEFAULT expr] [NOT AGGREGATE]` | clause |
| `opt_not_aggregate` | `NOT AGGREGATE` | helper |
| `opt_default_expression` | `DEFAULT expr` | helper |
| `type_or_tvf_schema` | `type \| templated_parameter_type \| tvf_schema` | datatype |
| `tvf_schema` | `TABLE< tvf_schema_column,… >` | datatype |
| `tvf_schema_column` | `[identifier] type` | datatype |
| `templated_parameter_type` | `ANY <kind>` | datatype |
| `templated_parameter_kind` | `PROTO \| ENUM \| STRUCT \| ARRAY \| identifier` | helper |
| `procedure_parameters` | parenthesized `procedure_parameter` list | clause |
| `procedure_parameter` | `[mode] id type_or_tvf_schema`; error alt for missing type | clause |
| `procedure_parameter_termination` | `) \| ,` (error-path marker) | helper |
| `opt_procedure_parameter_mode` | `IN \| OUT \| INOUT` | helper |

### 2.7 DML

| rule | purpose | category |
| --- | --- | --- |
| `dml_statement` | `insert_statement \| delete_statement \| update_statement` | dml |
| `insert_statement` | `INSERT [OR …] [INTO] <path> [hint] [cols] <values-or-query> [ON CONFLICT] [ASSERT ROWS MODIFIED] [THEN RETURN]` | dml |
| `update_statement` | `UPDATE <path> [hint] [alias] [WITH OFFSET] SET <items> [FROM] [WHERE] [ASSERT ROWS MODIFIED] [THEN RETURN]` | dml |
| `delete_statement` | `DELETE [FROM] <path> [hint] [alias] [WITH OFFSET] [WHERE] [ASSERT ROWS MODIFIED] [THEN RETURN]` | dml |
| `merge_statement` | `MERGE [INTO] <path> [alias] USING <source> ON <expr> <when-clauses>+` | dml |
| `merge_source` | `table_path_expression \| table_subquery` | helper |
| `merge_when_clause` | `WHEN [NOT] MATCHED [BY TARGET\|BY SOURCE] [AND expr] THEN <action>` | clause |
| `merge_action` | `INSERT [cols] <values/row> \| UPDATE SET <items> \| DELETE` | clause |
| `merge_insert_value_list_or_source_row` | `VALUES <row> \| ROW` | helper |
| `by_target` | `BY TARGET` | helper |
| `opt_and_expression` | `AND expr` | helper |
| `insert_statement_prefix` | `INSERT [OR IGNORE/REPLACE/UPDATE] [INTO] <path> [hint]` | helper |
| `opt_or_ignore_replace_update` | `[OR] IGNORE \| [OR] REPLACE \| [OR] UPDATE` | helper |
| `opt_into` | `INTO` | helper |
| `insert_values_or_query` | `insert_values_list \| query` | helper |
| `insert_values_list` | `VALUES <row> [, <row>]…` | clause |
| `insert_values_list_or_table_clause` | `insert_values_list \| table_clause_unreversed` | helper |
| `insert_values_row` | parenthesized `expression_or_default` list | clause |
| `expression_or_default` | `expression \| DEFAULT` | helper |
| `on_conflict_clause` | `ON CONFLICT [target] DO NOTHING \| DO UPDATE SET …` | clause |
| `opt_conflict_target` | `column_list \| ON UNIQUE CONSTRAINT id` | helper |
| `update_item_list` | `,`-separated `update_item` | clause |
| `update_item` | `update_set_value \| nested_dml_statement` | helper |
| `update_set_value` | `<generalized_path> = expression_or_default` | clause |
| `nested_dml_statement` | `( dml_statement )` | dml |
| `table_clause_unreversed` | `TABLE <table_clause_no_keyword>` | helper |
| `table_clause_no_keyword` | `path [WHERE] \| tvf [WHERE]` | helper |
| `opt_returning_clause` | `THEN RETURN [WITH ACTION [AS id]] <select_list>` | clause |
| `opt_assert_rows_modified` | `ASSERT_ROWS_MODIFIED <int>` | clause |
| `opt_where_expression` | `WHERE expr` | helper |
| `maybe_dashed_generalized_path_expression` | `generalized_path_expression \| dashed_path_expression` (DML target) | identifier |

### 2.8 DCL — privileges, roles, grants

| rule | purpose | category |
| --- | --- | --- |
| `grant_statement` | `GRANT <privileges> ON [type] <path> TO <grantee_list>` | ddl (DCL) |
| `revoke_statement` | `REVOKE <privileges> ON [type] <path> FROM <grantee_list>` | ddl (DCL) |
| `privileges` | `ALL [PRIVILEGES] \| privilege_list` | helper |
| `privilege_list` | `,`-separated `privilege` | helper |
| `privilege` | `privilege_name [ (path-list) ]` | clause |
| `privilege_name` | `identifier \| SELECT` | helper |
| `path_expression_list_with_parens` | parenthesized path-expression list | helper |
| `grant_to_clause` | `GRANT TO ( grantee_list )` | clause |
| `grantee_list` | `,`-separated string/param grantees | helper |

### 2.9 Transaction / batch / session

| rule | purpose | category |
| --- | --- | --- |
| `begin_statement` | `BEGIN [TRANSACTION] \| START TRANSACTION` + mode list | statement |
| `begin_transaction_keywords` | `START TRANSACTION \| BEGIN [TRANSACTION]` | helper |
| `commit_statement` | `COMMIT [TRANSACTION]` | statement |
| `rollback_statement` | `ROLLBACK [TRANSACTION]` | statement |
| `set_statement` | `SET TRANSACTION <modes> \| SET id = expr \| SET @param = expr \| SET @@sysvar = expr \| SET (id-list) = expr` (+ error alt for unparenthesized multi-var) | statement |
| `transaction_mode_list` | `,`-separated `transaction_mode` | helper |
| `transaction_mode` | `READ ONLY \| READ WRITE \| ISOLATION LEVEL id [id]` | helper |
| `start_batch_statement` | `START BATCH [id]` | statement |
| `run_batch_statement` | `RUN BATCH` | statement |
| `abort_batch_statement` | `ABORT BATCH` | statement |
| `identifier_list` | `,`-separated identifiers | helper |

### 2.10 Utility / metadata / misc statements

| rule | purpose | category |
| --- | --- | --- |
| `explain_statement` | `EXPLAIN <unterminated_sql_statement>` | statement |
| `describe_statement` | `DESCRIBE\|DESC <info>` | statement |
| `describe_keyword` | `DESCRIBE \| DESC` | helper |
| `describe_info` | `[id] <slashed/dashed path> [FROM <path>]` | helper |
| `opt_from_path_expression` | `FROM <slashed/dashed path>` | helper |
| `show_statement` | `SHOW <target> [FROM <path>] [LIKE string]` | statement |
| `show_target` | `MATERIALIZED VIEWS \| identifier` | helper |
| `opt_like_string_literal` | `LIKE string` | helper |
| `analyze_statement` | `ANALYZE [OPTIONS] [table-and-column-info-list]` | statement |
| `table_and_column_info_list` / `table_and_column_info` | `<path> [column_list]` targets | clause |
| `assert_statement` | `ASSERT expr [AS description]` | statement |
| `opt_description` | `AS string` | helper |
| `call_statement` | `CALL <path> ( [tvf_argument,…] )` | statement |
| `execute_immediate` | `EXECUTE IMMEDIATE expr [INTO <ids>] [USING <args>]` | statement |
| `opt_execute_into_clause` | `INTO <identifier_list>` | helper |
| `opt_execute_using_clause` | `USING <execute_using_argument_list>` | helper |
| `execute_using_argument_list` / `execute_using_argument` | `expr [AS id]` args | helper |
| `import_statement` | `IMPORT <MODULE\|PROTO> <path/string> [AS/INTO id] [OPTIONS]` | statement |
| `import_type` | `MODULE \| PROTO` | helper |
| `opt_as_or_into_alias` | `(AS\|INTO) identifier` | helper |
| `path_expression_or_string` | `path_expression \| string_literal` | helper |
| `module_statement` | `MODULE <path> [OPTIONS]` | statement |
| `export_data_statement` | `EXPORT DATA [WITH CONNECTION] [OPTIONS] AS query` | statement |
| `export_data_no_query` | `EXPORT DATA [WITH CONNECTION] [OPTIONS]` | helper |
| `export_model_statement` | `EXPORT MODEL <path> [WITH CONNECTION] [OPTIONS]` | statement |
| `export_metadata_statement` | `EXPORT TABLE\|TABLE FUNCTION METADATA FROM <path> [WITH CONNECTION] [OPTIONS]` | statement |
| `aux_load_data_statement` | `LOAD DATA INTO\|OVERWRITE <scoped path> [elems] [PARTITIONS] [collate] [PARTITION BY] [CLUSTER BY] [OPTIONS] FROM FILES <opts> [external WITH]` | statement |
| `maybe_dashed_path_expression_with_scope` | `[TEMP\|TEMPORARY] TABLE <path>` for LOAD DATA | helper |
| `append_or_overwrite` | `INTO \| OVERWRITE` | helper |
| `load_data_partitions_clause` | `[OVERWRITE] PARTITIONS (expr)` | clause |
| `aux_load_data_from_files_options_list` | `FROM FILES <options_list>` | clause |
| `clone_data_statement` | `CLONE DATA INTO <path> FROM <clone source list>` | statement |
| `opt_external_table_with_clauses` | `WITH PARTITION COLUMNS [elems]` and/or `WITH CONNECTION` | clause |
| `with_partition_columns_clause` | `WITH PARTITION COLUMNS [elems]` | clause |
| `cluster_by_clause_prefix_no_hint` | `CLUSTER BY <expr-list>` | clause |

### 2.11 Procedural / scripting statements

| rule | purpose | category |
| --- | --- | --- |
| `unterminated_script_statement` | dispatch to script statements (if/case/declare/break/continue/return/raise/labeled) | statement |
| `unterminated_unlabeled_script_statement` | `begin_end_block \| while \| loop \| repeat \| for_in` | statement |
| `begin_end_block` | `BEGIN [statement_list] [exception handler] END` | clause |
| `opt_exception_handler` | `EXCEPTION WHEN ERROR THEN <statement_list>` | clause |
| `if_statement` | `IF expr THEN … [ELSEIF…] [ELSE…] END IF` | statement |
| `elseif_clauses` | `(ELSEIF expr THEN …)+` | clause |
| `opt_else` | `ELSE [statement_list]` | clause |
| `case_statement` | `CASE [expr] <when_then>… [ELSE] END CASE` | statement |
| `when_then_clauses` | `(WHEN expr THEN …)+` | clause |
| `while_statement` | `WHILE expr DO … END WHILE` | statement |
| `loop_statement` | `LOOP … END LOOP` | statement |
| `repeat_statement` | `REPEAT … UNTIL expr END REPEAT` | statement |
| `until_clause` | `UNTIL expr` | clause |
| `for_in_statement` | `FOR id IN (query) DO … END FOR` | statement |
| `variable_declaration` | `DECLARE <ids> type [DEFAULT expr] \| DECLARE <ids> DEFAULT expr` | statement |
| `break_statement` | `BREAK \| LEAVE [id]` | statement |
| `continue_statement` | `CONTINUE \| ITERATE [id]` | statement |
| `return_statement` | `RETURN` | statement |
| `raise_statement` | `RAISE \| RAISE USING MESSAGE = expr` | statement |
| `label` | `identifier` (script label; `TODO: refine label`) | helper |

### 2.12 GQL — graph query language (`GRAPH … <ops>`)

| rule | purpose | category |
| --- | --- | --- |
| `gql_statement` | `GRAPH <path> <graph_operation_block>` | statement |
| `graph_operation_block` | `NEXT`-separated composite query blocks | query |
| `graph_composite_query_block` | linear op or composite prefix | query |
| `graph_composite_query_prefix` | set-operation-joined linear ops | query |
| `graph_set_operation_metadata` | `<set-op-type> <ALL\|DISTINCT>` | clause |
| `graph_linear_query_operation` | `[operator list] <return op>` | query |
| `graph_linear_operator_list` / `graph_linear_operator` | sequence of MATCH/OPTIONAL MATCH/LET/FILTER/ORDER BY/PAGE/WITH/FOR/SAMPLE | query |
| `graph_match_operator` / `graph_optional_match_operator` | `[OPTIONAL] MATCH [hint] <pattern>` | clause |
| `graph_let_operator` | `LET <var defs>` | clause |
| `graph_let_variable_definition_list` / `graph_let_variable_definition` | `id = expr` | clause |
| `graph_filter_operator` | `FILTER <where>\|<expr>` | clause |
| `graph_order_by_operator` / `graph_order_by_clause` | `ORDER [hint] BY <ordering exprs>` | clause |
| `graph_ordering_expression` | `expr [collate] [asc/desc] [null_order]` | clause |
| `opt_graph_asc_or_desc` | `ASC\|DESC \| ASCENDING \| DESCENDING` | helper |
| `graph_page_operator` / `graph_page_clause` | `OFFSET\|SKIP <n> [LIMIT <n>] \| LIMIT <n>` | clause |
| `graph_with_operator` | `WITH [ALL\|DISTINCT] [hint] <return items> [GROUP BY]` | clause |
| `graph_for_operator` | `FOR id IN expr [WITH OFFSET [alias]]` | clause |
| `opt_with_offset_and_alias_with_required_as` | `WITH OFFSET [AS alias]` | helper |
| `graph_sample_clause` | `TABLESAMPLE id ( size ) [suffix]` | clause |
| `opt_graph_sample_clause_suffix` | repeatable / `WITH WEIGHT [AS id]` | helper |
| `graph_return_operator` | `RETURN [hint] [ALL\|DISTINCT] <items> [GROUP BY] [ORDER BY] [PAGE]` | clause |
| `graph_return_item_list` / `graph_return_item` | `expr [AS id] \| *` | clause |
| `graph_pattern` | `<path pattern list> [WHERE]` | clause |
| `graph_path_pattern_list` / `graph_path_pattern` | path patterns (+ hints, search/mode prefixes) | clause |
| `graph_path_pattern_expr` | `graph_path_factor (hint? factor)*` | clause |
| `graph_path_factor` | primary or quantified primary | helper |
| `graph_quantified_path_primary` | `primary { [n,] m }` quantifier | clause |
| `graph_path_primary` | element pattern or parenthesized path | helper |
| `graph_parenthesized_path_pattern` | `( [hint] <path> [WHERE] )` | clause |
| `graph_element_pattern` | node or edge pattern | helper |
| `graph_node_pattern` | `( <filler> )` | clause |
| `graph_edge_pattern` | `-[ <filler> ]-`, `<-`, `->`, etc. | clause |
| `graph_element_pattern_filler` | `[hint] [id] [label-expr] [property-spec] [WHERE]` (`TODO`: empty-production note) | clause |
| `graph_property_specification` / `graph_property_name_and_value` | `{ id : expr , … }` | clause |
| `opt_is_label_expression` | `IS <label_expr> \| : <label_expr>` | helper |
| `label_expression` / `label_primary` | label algebra (`&`, `\|`, `!`, `%`, parens, id) | expression |
| `parenthesized_label_expression` | `( label_expression )` | helper |
| `opt_graph_element_identifier` | `graph_identifier` | helper |
| `opt_graph_path_mode_prefix` / `opt_graph_path_mode` | `WALK\|TRAIL\|SIMPLE\|ACYCLIC [PATH(S)]` | helper |
| `path_or_paths` | `PATH \| PATHS` | helper |
| `opt_graph_search_prefix` | `(ANY\|ALL) [SHORTEST]` | helper |
| `opt_path_variable_assignment` | `graph_identifier =` | helper |
| `graph_identifier` | `token_identifier \| common_keyword_as_identifier` | identifier |

### 2.13 Query / SELECT core

| rule | purpose | category |
| --- | --- | --- |
| `query` | `query_without_pipe_operators` | query |
| `query_without_pipe_operators` | `[WITH] <primary-or-set-op> [ORDER BY] [LIMIT/OFFSET]`; several error alts (trailing comma after WITH, pipe `\|>` after WITH, unexpected/bad keyword after FROM-query) | query |
| `bad_keyword_after_from_query` | error marker: `WHERE\|SELECT\|GROUP` after a FROM-query | helper |
| `bad_keyword_after_from_query_allows_parens` | error marker: `ORDER\|UNION\|INTERSECT\|EXCEPT\|LIMIT` | helper |
| `with_clause_with_trailing_comma` | error marker | helper |
| `select_or_from_keyword` | `SELECT \| FROM` (error-path) | helper |
| `query_primary_or_set_operation` | `query_primary \| query_set_operation` | query |
| `query_primary` | `select \| ( parenthesized_query ) [AS alias]` | query |
| `query_set_operation` / `query_set_operation_prefix` | left-recursive set-operation chain (+ FROM-after-set-op error alts) | query |
| `query_set_operation_item` | `<set-op metadata> query_primary` | query |
| `set_operation_metadata` | `[corresponding outer] <UNION\|EXCEPT\|INTERSECT> [hint] <ALL\|DISTINCT> [STRICT] [CORRESPONDING [BY]]` | clause |
| `query_set_operation_type` | `UNION \| EXCEPT \| INTERSECT` | helper |
| `all_or_distinct` | `ALL \| DISTINCT` | helper |
| `opt_corresponding_outer_mode` | `FULL [OUTER] \| OUTER \| LEFT [OUTER]` | helper |
| `opt_outer` | `OUTER` | helper |
| `opt_strict` | `STRICT` | helper |
| `opt_column_match_suffix` | `CORRESPONDING [BY]` | helper |
| `with_clause` | `WITH [RECURSIVE] <aliased_query>,…` (CTE) | clause |
| `aliased_query` | `id AS ( query ) [recursion-depth modifiers]` | clause |
| `opt_aliased_query_modifiers` | `recursion_depth_modifier` | helper |
| `recursion_depth_modifier` | `WITH DEPTH [AS alias] [BETWEEN n AND n \| MAX n]` | clause |
| `possibly_unbounded_int_literal_or_parameter` | `<int/param> \| UNBOUNDED` | helper |
| `int_literal_or_parameter` | `integer_literal \| @param \| @@sysvar` | helper |
| `order_by_clause` / `order_by_clause_prefix` | `ORDER [hint] BY <ordering exprs>` | clause |
| `ordering_expression` | `expr [collate] [ASC/DESC] [null_order]` | clause |
| `select` | `select_clause [from_clause] [clauses-following-from]` | query |
| `select_clause` | `SELECT [hint] [WITH …] [ALL\|DISTINCT] [AS struct/path] <select_list>`; error alt for empty list before FROM | query |
| `opt_select_with` | `WITH id [OPTIONS list]` | clause |
| `opt_select_as_clause` | `AS STRUCT \| AS <path>` | clause |
| `opt_clauses_following_from` | WHERE→GROUP BY→HAVING→QUALIFY→WINDOW ordering chain (from FROM) | clause |
| `opt_clauses_following_where` | GROUP BY→HAVING→QUALIFY→WINDOW chain | clause |
| `opt_clauses_following_group_by` | HAVING→QUALIFY→WINDOW chain | clause |
| `where_clause` | `WHERE expr` | clause |
| `group_by_clause` | `group_by_all \| group_by_clause_prefix` | clause |
| `group_by_all` | `GROUP [hint] [AND ORDER] BY ALL` | clause |
| `group_by_clause_prefix` | `GROUP [hint] [AND ORDER] BY <grouping items>` | clause |
| `group_by_preamble` | `GROUP [hint] [AND ORDER] BY` | helper |
| `opt_and_order` | `AND ORDER` | helper |
| `grouping_item` | `() \| expr [AS alias] [order] \| ROLLUP(...) \| CUBE(...) \| GROUPING SETS(...)` | clause |
| `grouping_set_list` | `GROUPING SETS ( <set>,… )` | clause |
| `grouping_set` | `() \| expr \| ROLLUP(...) \| CUBE(...)` | clause |
| `cube_list` | `CUBE ( <exprs> )` | clause |
| `rollup_list` | `ROLLUP ( <exprs> )` | clause |
| `having_clause` | `HAVING expr` | clause |
| `qualify_clause_nonreserved` | `QUALIFY expr` | clause |
| `window_clause` / `window_clause_prefix` | `WINDOW <window_definition>,…` | clause |
| `window_definition` | `id AS <window_specification>` | clause |
| `limit_offset_clause` | `LIMIT expr [OFFSET expr]` | clause |
| `opt_grouping_item_order` | `opt_selection_item_order \| null_order` | helper |
| `opt_selection_item_order` | `<ASC/DESC> [null_order]` | helper |
| `asc_or_desc` | `ASC \| DESC` | helper |
| `null_order` | `NULLS FIRST \| NULLS LAST` | helper |

### 2.14 Query — FROM clause / table sources / joins

| rule | purpose | category |
| --- | --- | --- |
| `from_clause` | `FROM <from_clause_contents>` | clause |
| `from_clause_contents` | `table_primary <suffix>*`; error alts for `@`/`?`/`@@` used as table names | clause |
| `from_clause_contents_suffix` | `, table_primary` or `[NATURAL] [join_type] [join_hint] JOIN [hint] table_primary [on/using]` | clause |
| `table_primary` | TVF / path / `( join )` / subquery / match-recognize / sample | clause |
| `table_path_expression` | `<base> [hint] [pivot/unpivot/alias] [WITH OFFSET] [AT SYSTEM TIME]` | clause |
| `table_path_expression_base` | `unnest_expression \| <slashed/dashed path>`; error alts for array/field access without UNNEST | clause |
| `table_subquery` | `( query ) [pivot/unpivot/alias]` | clause |
| `join` | `table_primary join_item*` | query |
| `join_item` | `[NATURAL] [join_type] [join_hint] JOIN [hint] table_primary [on/using]` | clause |
| `join_type` | `CROSS \| FULL [OUTER] \| INNER \| LEFT [OUTER] \| RIGHT [OUTER]` | helper |
| `opt_natural` | `NATURAL` | helper |
| `join_hint` | `HASH \| LOOKUP` | helper |
| `on_or_using_clause_list` / `on_or_using_clause` | `on_clause+`/`using_clause+` | clause |
| `on_clause` | `ON expr` | clause |
| `using_clause` | `USING ( id [. id]* )` | clause |
| `tvf_with_suffixes` | TVF call `( … )` with hint + pivot/unpivot/alias | clause |
| `tvf_prefix` / `tvf_prefix_no_args` | TVF name + argument list start (`path \| IF (`) | clause |
| `tvf_argument` | `expr \| descriptor \| table \| model \| connection \| named-arg`; many parenthesization error alts | clause |
| `descriptor_argument` | `DESCRIPTOR ( <columns> )` | clause |
| `descriptor_column_list` / `descriptor_column` | descriptor column identifiers | helper |
| `table_clause` | `TABLE <tvf> \| TABLE <path>` | clause |
| `model_clause` | `MODEL <path>` | clause |
| `connection_clause` | `CONNECTION <path-or-DEFAULT>` | clause |
| `path_expression_or_default` | `path_expression \| DEFAULT` | helper |
| `sequence_arg` | `SEQUENCE <path>` | clause |
| `sample_clause` | `TABLESAMPLE id ( size ) <suffix>` | clause |
| `opt_sample_clause_suffix` | repeatable / `WITH WEIGHT [id\|AS id]` | helper |
| `repeatable_clause` | `REPEATABLE ( <int> )` | clause |
| `sample_size` | `<value> <unit> [PARTITION BY]` | clause |
| `sample_size_value` | `<int/param> \| float` | helper |
| `sample_size_unit` | `ROWS \| PERCENT` | helper |
| `possibly_cast_int_literal_or_parameter` | `cast_int… \| int_literal_or_parameter` | helper |
| `cast_int_literal_or_parameter` | `CAST ( <int> AS type [FORMAT] )` | helper |
| `partition_by_clause_prefix_no_hint` | `PARTITION BY <exprs>` (no hint) | clause |
| `opt_at_system_time` | `FOR SYSTEM [TIME] AS OF expr` | clause |
| `opt_with_offset_and_alias` | `WITH OFFSET [alias]` | helper |
| `as_alias` | `[AS] identifier` | helper |
| `opt_as_alias_with_required_as` | `AS identifier` | helper |
| `maybe_slashed_or_dashed_path_expression` | `maybe_dashed \| slashed` path | identifier |
| `maybe_dashed_path_expression` | `path_expression \| dashed_path_expression` | identifier |

### 2.15 Query — PIVOT / UNPIVOT / UNNEST / MATCH_RECOGNIZE

| rule | purpose | category |
| --- | --- | --- |
| `pivot_or_unpivot_clause_and_aliases` | post-TVF pivot/unpivot + alias combinations (+ QUALIFY error alts) | clause |
| `opt_pivot_or_unpivot_clause_and_alias` | post-table pivot/unpivot + alias combinations (+ QUALIFY error alts) | clause |
| `pivot_clause` | `PIVOT ( <exprs> FOR <expr> IN ( <values> ) )` | clause |
| `pivot_expression_list` / `pivot_expression` | `expr [alias]` (pivot aggregates) | clause |
| `pivot_value_list` / `pivot_value` | `expr [alias]` (pivot column values) | clause |
| `unpivot_clause` | `UNPIVOT [nulls-filter] ( <cols> FOR <path> IN ( <items> ) )` | clause |
| `unpivot_nulls_filter` | `EXCLUDE NULLS \| INCLUDE NULLS` | helper |
| `unpivot_in_item_list` / `unpivot_in_item_list_prefix` / `unpivot_in_item` | UNPIVOT IN-items | clause |
| `opt_as_string_or_integer` | `[AS] string \| [AS] integer` | helper |
| `path_expression_list_with_opt_parens` | `( path-list ) \| path-list` | helper |
| `path_expression_list` | `,`-separated path expressions | helper |
| `unnest_expression` | `UNNEST ( <exprs-with-alias> [, zip-mode] )`; error alt for `(SELECT…)` | clause |
| `unnest_expression_prefix` | `UNNEST ( <expression_with_opt_alias>,… ` | helper |
| `opt_array_zip_mode` | `, <named_argument>` | helper |
| `expression_with_opt_alias` | `expr [AS alias]` | helper |
| `match_recognize_clause` | `MATCH_RECOGNIZE ( [PARTITION BY] <ORDER BY> MEASURES <aliased exprs> PATTERN ( <pat> ) DEFINE <var defs> ) [alias]` | clause |
| `row_pattern_expr` | row-pattern alternation (`\|`) | expression |
| `row_pattern_concatenation` | row-pattern concatenation | expression |
| `row_pattern_factor` | `id \| ( <pattern> )` | helper |
| `select_list_prefix_with_as_aliases` / `select_column_expr_with_as_alias` | MEASURES list (`expr AS id`) | clause |

### 2.16 Query — SELECT list & star modifiers

| rule | purpose | category |
| --- | --- | --- |
| `select_list` | `,`-separated `select_list_item` (trailing comma allowed) | clause |
| `select_list_item` | `select_column_expr \| select_column_dot_star \| select_column_star` | clause |
| `select_column_expr` | `expr \| expr AS id \| expr id` | clause |
| `select_column_star` | `* [star_modifiers]` | clause |
| `select_column_dot_star` | `<expr>.* [star_modifiers]` | clause |
| `star_modifiers` | `EXCEPT(...) [REPLACE(...)]` | clause |
| `star_except_list` | `EXCEPT ( id,… )` | clause |
| `star_replace_list` | `REPLACE ( <items> )` | clause |
| `star_replace_item` | `expr AS id` | helper |

### 2.17 Expressions — precedence chain & operators

| rule | purpose | category |
| --- | --- | --- |
| `expression` | top of precedence chain: `<higher-prec> \| and_expression \| expression OR expression` | expression |
| `expression_higher_prec_than_and` | the big primary+operator expression rule (literals, constructors, CASE/CAST/EXTRACT, function calls, field/array access, NOT, LIKE/IN/BETWEEN/IS, comparisons, bitwise, shift, additive, multiplicative, unary, parenthesized) — port of ZetaSQL `unparenthesized_expression_higher_prec_than_and` | expression |
| `expression_maybe_parenthesized_not_a_query` | same chain but allowing parenthesized-not-a-query at head (RHS contexts) | expression |
| `and_expression` | `<higher-prec> AND <higher-prec> (AND <higher-prec>)*` | expression |
| `parenthesized_query` | `( query )` | expression |
| `parenthesized_expression_not_a_query` | `( <expr-not-a-query> )` | expression |
| `parenthesized_in_rhs` | RHS of IN: `( query ) \| ( expr ) \| <2+ list>` | expression |
| `parenthesized_anysomeall_list_in_rhs` | RHS of ANY/SOME/ALL list | expression |
| `in_list_two_or_more_prefix` | `( expr , expr [, expr]* ` | helper |
| `unary_operator` | `+ \| - \| ~` | helper |
| `comparative_operator` | `= != <> < <= > >=` | helper |
| `shift_operator` | `<< \| >>` | helper |
| `additive_operator` | `+ \| -` | helper |
| `multiplicative_operator` | `* \| /` | helper |
| `is_operator` | `IS [NOT]` | helper |
| `between_operator` | `[NOT] BETWEEN` | helper |
| `in_operator` | `[NOT] IN` | helper |
| `distinct_operator` | `IS [NOT] DISTINCT FROM` | helper |
| `like_operator` | `LIKE \| NOT LIKE` | helper |
| `any_some_all` | `ANY \| SOME \| ALL` | helper |
| `expression_subquery_with_keyword` | `ARRAY ( query ) \| EXISTS [hint] ( query )` | expression |
| `interval_expression` | `INTERVAL expr id [TO id]` | expression |

### 2.18 Expressions — constructors, CASE/CAST/EXTRACT, function calls

| rule | purpose | category |
| --- | --- | --- |
| `case_expression` | `<prefix> [ELSE expr] END` | expression |
| `case_expression_prefix` | `case_value… \| case_no_value…` | expression |
| `case_value_expression_prefix` | `CASE expr (WHEN expr THEN expr)+` | expression |
| `case_no_value_expression_prefix` | `CASE (WHEN expr THEN expr)+` | expression |
| `cast_expression` | `CAST ( expr AS type [FORMAT] )` / `SAFE_CAST(...)`; error alts for `CAST(CAST…` | expression |
| `extract_expression` | `EXTRACT ( expr FROM expr ) [AT TIME ZONE expr]` | expression |
| `extract_expression_base` | `EXTRACT ( expr FROM expr` | helper |
| `replace_fields_expression` / `replace_fields_prefix` | `REPLACE_FIELDS ( expr, <args> )` | expression |
| `replace_fields_arg` | `expr AS <generalized path/extension>` | helper |
| `with_expression` | `WITH ( <var defs>, expr )` (inline WITH expression) | expression |
| `with_expression_variable_prefix` / `with_expression_variable` | `id AS expr` | helper |
| `array_constructor` | `[ARRAY\|array_type] [ <exprs> ]` | expression |
| `array_constructor_prefix` / `array_constructor_prefix_no_expressions` | array-constructor start | helper |
| `struct_constructor` | parenthesized struct constructor (keyword/typed/bare forms) | expression |
| `struct_constructor_prefix_with_keyword` / `…_no_arg` / `…_without_keyword` | struct-constructor variants | helper |
| `struct_constructor_arg` | `expr [AS alias]` | helper |
| `new_constructor` | `NEW <type> ( <args> )` (proto) | expression |
| `new_constructor_prefix` / `…_no_arg` | new-constructor start | helper |
| `new_constructor_arg` | `expr [AS id] [AS (path)]` | helper |
| `braced_new_constructor` | `NEW <type> <braced_constructor>` | expression |
| `braced_constructor` | `{ <fields> }` (proto braced) | expression |
| `braced_constructor_start` / `…_prefix` / `…_field` / `…_lhs` / `…_field_value` / `…_extension` | braced-constructor internals | helper |
| `struct_braced_constructor` | `[struct_type \| STRUCT] <braced_constructor>` | expression |
| `function_call_expression_with_clauses` | `<path> ( [DISTINCT] … ) \| <keyword-fn> ( … )` | expression |
| `function_call_expression_with_clauses_suffix` | argument list + null-handling/having/clamped/with-report/order/limit modifiers + hint + with-group-rows + over | expression |
| `function_call_argument` | `expr [AS alias] \| named-arg \| lambda \| sequence-arg`; error alt for `SELECT` | expression |
| `named_argument` | `id => expr \| id => lambda` | expression |
| `lambda_argument` | `<arg list> -> expr` | expression |
| `lambda_argument_list` | `expr \| ()` | helper |
| `function_name_from_keyword` | `IF \| GROUPING \| LEFT \| RIGHT \| COLLATE \| RANGE` (keywords usable as function names) | identifier |
| `opt_null_handling_modifier` | `IGNORE NULLS \| RESPECT NULLS` | helper |
| `opt_having_or_group_by_modifier` | `HAVING MAX expr \| HAVING MIN expr <group-by>` | clause |
| `clamped_between_modifier` | `CLAMPED BETWEEN … AND …` (differential privacy) | clause |
| `with_report_modifier` | `WITH REPORT <format>` | clause |
| `with_report_format` | `options_list` | helper |
| `with_group_rows` | `WITH GROUP ROWS` | clause |
| `over_clause` | `OVER <window_specification>` | clause |
| `window_specification` | `id \| ( [id] [PARTITION BY] [ORDER BY] [frame] )` | clause |
| `opt_window_frame_clause` | `<unit> BETWEEN <bound> AND <bound> \| <unit> <bound>` | clause |
| `window_frame_bound` | `UNBOUNDED <prec/foll> \| CURRENT ROW \| expr <prec/foll>` | clause |
| `preceding_or_following` | `PRECEDING \| FOLLOWING` | helper |
| `frame_unit` | `ROWS \| RANGE` | helper |
| `partition_by_clause` / `partition_by_clause_prefix` | `PARTITION [hint] BY <exprs>` | clause |
| `generalized_path_expression` | path with field/array access (`. id`, `.(path)`, `[expr]`) | identifier |
| `generalized_extension_path` | proto extension path `(path)` chains | identifier |

### 2.19 Hints & OPTIONS

| rule | purpose | category |
| --- | --- | --- |
| `hint` | `@<int> \| <hint_with_body>` (statement/table/operator hints) | clause |
| `hint_with_body` / `hint_with_body_prefix` | `@[<int>@]{ <entries> }` | clause |
| `hint_entry` | `<id> = expr \| <id>.<id> = expr` | clause |
| `identifier_in_hints` | `identifier \| extra_identifier_in_hints_name` | identifier |
| `extra_identifier_in_hints_name` | `HASH \| PROTO \| PARTITION` | helper |
| `options_list` | `( <entries> ) \| ()` | clause |
| `options_list_prefix` | `( <options_entry>,… ` | helper |
| `options_entry` | `<id-in-hints> <op> <expr-or-PROTO>` | clause |
| `expression_or_proto` | `PROTO \| expression` | helper |
| `options_assignment_operator` | `= \| += \| -=` | helper |

### 2.20 Data types

| rule | purpose | category |
| --- | --- | --- |
| `type` | `raw_type [type_parameters] [collate]` | datatype |
| `raw_type` | `array_type \| struct_type \| type_name \| range_type \| function_type \| map_type` | datatype |
| `type_name` | `path_expression \| INTERVAL` | datatype |
| `array_type` | `ARRAY< type >` | datatype |
| `struct_type` / `struct_type_prefix` | `STRUCT< [fields] >` | datatype |
| `struct_field` | `[identifier] type` | datatype |
| `range_type` | `RANGE< type >` | datatype |
| `map_type` | `MAP< keyType, valueType >` | datatype |
| `function_type` / `function_type_prefix` | `FUNCTION< (args) -> returnType >` | datatype |
| `template_type_open` / `template_type_close` | `<` / `>` (type-template brackets) | helper |
| `opt_type_parameters` / `type_parameters_prefix` | `( <params> )`; error alt for trailing comma | datatype |
| `type_parameter` | `int \| bool \| string \| bytes \| float \| MAX` | datatype |
| `collate_clause` | `COLLATE <string-or-param>` | clause |

### 2.21 Literals

| rule | purpose | category |
| --- | --- | --- |
| `string_literal` | one-or-more `string_literal_component` (adjacent concatenation; predicate requires whitespace between, rejects mixing with bytes) | literal |
| `string_literal_component` | `STRING_LITERAL` token | literal |
| `bytes_literal` | one-or-more `bytes_literal_component` (predicate requires whitespace between, rejects mixing with string) | literal |
| `bytes_literal_component` | `BYTES_LITERAL` token | literal |
| `integer_literal` | `INTEGER_LITERAL` | literal |
| `floating_point_literal` | `FLOATING_POINT_LITERAL` | literal |
| `numeric_literal` / `numeric_literal_prefix` | `NUMERIC\|DECIMAL '…'` | literal |
| `bignumeric_literal` / `bignumeric_literal_prefix` | `BIGNUMERIC\|BIGDECIMAL '…'` | literal |
| `json_literal` | `JSON '…'` | literal |
| `range_literal` | `RANGE<type> '…'` | literal |
| `date_or_time_literal` / `date_or_time_literal_kind` | `DATE\|TIME\|DATETIME\|TIMESTAMP '…'` | literal |
| `null_literal` | `NULL` | literal |
| `boolean_literal` | `TRUE \| FALSE` | literal |
| `signed_numeric_literal` | (listed in §2.4) | literal |

### 2.22 Identifiers, paths, parameters

| rule | purpose | category |
| --- | --- | --- |
| `identifier` | `token_identifier \| keyword_as_identifier` | identifier |
| `token_identifier` | `IDENTIFIER` token (unquoted or backtick-quoted) | identifier |
| `keyword_as_identifier` | `common_keyword_as_identifier \| SIMPLE` | identifier |
| `common_keyword_as_identifier` | the **211**-keyword non-reserved set that may be used as a bare identifier | identifier |
| `path_expression` | `identifier (. identifier)*` | identifier |
| `dashed_path_expression` | dashed path (`a-b.c`) for BigQuery project/dataset names | identifier |
| `dashed_identifier` | hyphen-joined identifier/integer/float sequences | identifier |
| `slashed_path_expression` / `slashed_identifier` | slash/colon-joined path (BigQuery resource paths) | identifier |
| `slashed_identifier_separator` | `- / :` separator sequence | helper |
| `identifier_or_integer` | `identifier \| INTEGER_LITERAL` | helper |
| `parameter_expression` | `named_parameter_expression \| ?` | identifier |
| `named_parameter_expression` | `@ identifier` | identifier |
| `system_variable_expression` | `@@ path_expression` | identifier |
| `string_literal_or_parameter` | `string \| @param \| @@sysvar` | helper |

---

## 3. Lexer summary (`GoogleSQLLexer.g4`)

- **Case sensitivity:** `options { caseInsensitive = true; }` — keywords match case-insensitively via per-letter `fragment A..Z` and the global flag. The new parser must treat keywords/identifiers case-insensitively.

### 3.1 Token categories

| Category | Count / notes |
| --- | --- |
| Word keywords (`*_SYMBOL: 'WORD';`) | **307** |
| Operator tokens (`*_OPERATOR`) | 16 |
| Punctuation/bracket tokens (`*_SYMBOL` with non-alpha value) | 27 |
| Literal tokens | `STRING_LITERAL`, `BYTES_LITERAL`, `INTEGER_LITERAL`, `FLOATING_POINT_LITERAL` (+ several `UNCLOSED_*` error tokens) |
| Identifier tokens | `IDENTIFIER`, `UNCLOSED_ESCAPED_IDENTIFIER` |
| Hidden-channel | `WHITESPACE`, `COMMENT` |

### 3.2 Operator & punctuation tokens

`=` `!=` `<>` `<` `<=` `>` `>=` `<<` `>>` `+` `-` `*` `/` `~` `!` `%`
(`*_OPERATOR` group, 16 total: `EQUAL_OPERATOR NOT_EQUAL_OPERATOR NOT_EQUAL2_OPERATOR LT_OPERATOR LE_OPERATOR GT_OPERATOR GE_OPERATOR KL_OPERATOR KR_OPERATOR PLUS_OPERATOR MINUS_OPERATOR MULTIPLY_OPERATOR DIVIDE_OPERATOR BITWISE_NOT_OPERATOR EXCLAMATION_OPERATOR MODULO_OPERATOR`)

Punctuation `*_SYMBOL`: `, . { } ( ) [ ] | : ; ' ''' " """ \` ? @ @@ => -> += -= |> ^ & ||`
(`COMMA DOT LC_BRACKET RC_BRACKET LR_BRACKET RR_BRACKET LS_BRACKET RS_BRACKET STROKE COLON SEMI SINGLE_QUOTE SINGLE_QUOTE_3 DOUBLE_QUOTE DOUBLE_QUOTE_3 BACKQUOTE QUESTION AT ATAT EQUAL_GT_BRACKET SUB_GT_BRACKET PLUS_EQUAL SUB_EQUAL PIPE CIRCUMFLEX BIT_AND BOOL_OR`).
Notable: `|>` (`PIPE_SYMBOL`) is tokenized (pipe-query operator) but the **parser only emits errors** for it — pipe-syntax is not actually parsed (see §6). `=>` for named args, `->` for lambdas/function-type, `+=`/`-=` for OPTIONS assignment, `@`/`@@` for parameters/system-variables.

### 3.3 Identifier rules

- `fragment UNQUOTED_IDENTIFIER: [A-Z_][A-Z0-9_]*;` (case-insensitive ⇒ `[A-Za-z_][A-Za-z0-9_]*`).
- `fragment BQTEXT: \`…\`` — **backtick-quoted** identifier; body allows escapes (`ANY_ESCAPE`) and any non-backtick/backslash/newline char.
- `IDENTIFIER: UNQUOTED_IDENTIFIER | BQTEXT;`
- `UNCLOSED_ESCAPED_IDENTIFIER: BQTEXT_0;` — error token for an unterminated backtick identifier.
- **Dashed/slashed identifiers** are NOT lexer tokens; they are assembled in the **parser** (`dashed_identifier`, `slashed_identifier`) from `IDENTIFIER`, `MINUS_OPERATOR`, `SLASH_SYMBOL`, `COLON_SYMBOL`, and integer/float tokens — supporting BigQuery `project-id.dataset.table` and resource-path forms.

### 3.4 Literal forms

- **String:** `STRING_LITERAL: R? ( SQTEXT | DQTEXT | SQ3TEXT | DQ3TEXT )` — single/double quoted and triple-quoted, with optional leading `R` (raw). Escapes via `ANY_ESCAPE` (`\` + any char, incl. line continuations).
- **Bytes:** `BYTES_LITERAL: (B | R B | B R) ( … )` — `b'…'`, raw-bytes `rb'…'`/`br'…'`, all four quote forms.
- **Raw** is the `R` prefix on strings/bytes; **triple-quoted** forms (`'''…'''`, `"""…"""`) supported for both.
- **Unclosed error tokens:** `UNCLOSED_STRING_LITERAL`, `UNCLOSED_TRIPLE_QUOTED_STRING_LITERAL`, `UNCLOSED_RAW_STRING_LITERAL`, `UNCLOSED_TRIPLE_QUOTED_RAW_STRING_LITERAL`, `UNCLOSED_BYTES_LITERAL`, `UNCLOSED_TRIPLE_QUOTED_BYTES_LITERAL`, `UNCLOSED_RAW_BYTES_LITERAL`, `UNCLOSED_TRIPLE_QUOTED_RAW_BYTES_LITERAL`.
- **Integer:** `INTEGER_LITERAL: DECIMAL_DIGITS | HEX_DIGITS` (`HEX_DIGITS: '0x' HEX_DIGIT+`).
- **Floating-point:** `FLOATING_POINT_LITERAL` — `digits . [digits] [E[±]digits]`, `[digits] . digits [E…]`, or `digits E [±] digits`; optional leading `+`/`-`.
- **Typed/prefixed literals** (DATE/TIME/DATETIME/TIMESTAMP/NUMERIC/DECIMAL/BIGNUMERIC/BIGDECIMAL/JSON/RANGE + string) are assembled in the **parser** (§2.21), not the lexer.

### 3.5 Comment styles & hidden channels

- `WHITESPACE: [ \t\f\r\n] -> channel(HIDDEN);`
- `COMMENT: (BLOCK_COMMENT | DASH_COMMENT | POUND_COMMENT) -> channel(HIDDEN);`
  - `BLOCK_COMMENT: '/**/' | '/*' ~[!] .*? '*/'` (note: `/*!…` is deliberately excluded — leaves room for hint-comment style, though not used).
  - `DASH_COMMENT: '--' …` to end-of-line.
  - `POUND_COMMENT: '#' …` to end-of-line (BigQuery-style `#` comments).

### 3.6 Keyword list — reserved vs non-reserved

The grammar distinguishes them via the `common_keyword_as_identifier` rule: any word-keyword **listed there is non-reserved** (usable as a bare identifier); word-keywords **absent from it are reserved**.

**Non-reserved keywords (211)** — the full `common_keyword_as_identifier` set (+ `SIMPLE` via `keyword_as_identifier`):
`ABORT, ACCESS, ACTION, AGGREGATE, ADD, ALTER, ALWAYS, ANALYZE, APPROX, ARE, ASSERT, BATCH, BEGIN, BIGDECIMAL, BIGNUMERIC, BREAK, CALL, CASCADE, CHECK, CLAMPED, CONFLICT, CLONE, COPY, CLUSTER, COLUMN, COLUMNS, COMMIT, CONNECTION, CONSTANT, CONSTRAINT, CONTINUE, CORRESPONDING, CYCLE, DATA, DATABASE, DATE, DATETIME, DECIMAL, DECLARE, DEFINER, DELETE, DELETION, DEPTH, DESCRIBE, DETERMINISTIC, DO, DROP, ELSEIF, ENFORCED, ERROR, EXCEPTION, EXECUTE, EXPLAIN, EXPORT, EXTEND, EXTERNAL, FILES, FILTER, FILL, FIRST, FOREIGN, FORMAT, FUNCTION, GENERATED, GRANT, GROUP_ROWS, HIDDEN, IDENTITY, IMMEDIATE, IMMUTABLE, IMPORT, INCLUDE, INCREMENT, INDEX, INOUT, INPUT, INSERT, INVOKER, ISOLATION, ITERATE, JSON, KEY, LANGUAGE, LAST, LEAVE, LEVEL, LOAD, LOOP, MACRO, MAP, MATCH, KW_MATCH_RECOGNIZE_NONRESERVED, MATCHED, MATERIALIZED, MAX, MAXVALUE, MEASURES, MESSAGE, METADATA, MIN, MINVALUE, MODEL, MODULE, NUMERIC, OFFSET, ONLY, OPTIONS, OUT, OUTPUT, OVERWRITE, PARTITIONS, PATTERN, PERCENT, PIVOT, POLICIES, POLICY, PRIMARY, PRIVATE, PRIVILEGE, PRIVILEGES, PROCEDURE, PROJECT, PUBLIC, RAISE, READ, REFERENCES, REMOTE, REMOVE, RENAME, REPEAT, REPEATABLE, REPLACE, REPLACE_FIELDS, REPLICA, REPORT, RESTRICT, RESTRICTION, RETURNS, RETURN, REVOKE, ROLLBACK, ROW, RUN, SAFE_CAST, SCHEMA, SEARCH, SECURITY, SEQUENCE, SETS, SHOW, SNAPSHOT, SOURCE, SQL, STABLE, START, STATIC_DESCRIBE, STORED, STORING, STRICT, SYSTEM, SYSTEM_TIME, TABLE, TABLES, TARGET, TEMP, TEMPORARY, TIME, TIMESTAMP, TRANSACTION, TRANSFORM, TRUNCATE, TYPE, UNDROP, UNIQUE, UNKNOWN, UNPIVOT, UNTIL, UPDATE, VALUE, VALUES, VECTOR, VIEW, VIEWS, VOLATILE, WEIGHT, WHILE, WRITE, ZONE, DESCRIPTOR, INTERLEAVE, NULL_FILTERED, PARENT, DESTINATION, PROPERTY, GRAPH, NODE, PROPERTIES, LABEL, EDGE, NEXT, ASCENDING, DESCENDING, SKIP, PATH, PATHS, WALK, TRAIL, ACYCLIC, OPTIONAL, LET` (+ `SIMPLE`).

**Reserved keywords (97)** — word-keywords NOT in `common_keyword_as_identifier`; these cannot be bare identifiers. Exact set (computed by set-difference of all 307 word tokens minus the 211 non-reserved):
`ALL, AND, ANY, ARRAY, AS, ASC, ASSERT_ROWS_MODIFIED, BETWEEN, BY, CASE, CAST, COLLATE, CREATE, CROSS, CUBE, CURRENT, DEFAULT, DEFINE, DELTA, DESC, DIFFERENTIAL_PRIVACY, DISTINCT, ELSE, END, ENUM, EPSILON, EXCEPT, EXCLUDE, EXISTS, EXTRACT, FALSE, FOLLOWING, FOR, FROM, FULL, GROUP, GROUPING, HASH, HAVING, IF, IGNORE, IN, INNER, INTERSECT, INTERVAL, INTO, IS, JOIN, LEFT, LIKE, LIMIT, LOOKUP, MATCH_RECOGNIZE, MAX_GROUPS_CONTRIBUTED, MERGE, NATURAL, NEW, NO, NOT, NOTHING, NULL, NULLS, OF, ON, OR, ORDER, OUTER, OVER, PARTITION, PRECEDING, PRIVACY_UNIT_COLUMN, PROTO, QUALIFY, RANGE, RECURSIVE, RESPECT, RIGHT, ROLLUP, ROWS, SELECT, SET, SHORTEST, SLASH, SOME, STRUCT, TABLESAMPLE, THEN, TO, TRUE, UNBOUNDED, UNION, UNNEST, USING, WHEN, WHERE, WINDOW, WITH`.
(Note: `SIMPLE`, `SYSTEM`, `VALUE`, `CONFLICT` are **non-reserved** — they appear in `common_keyword_as_identifier`/`keyword_as_identifier`. The precise reserved/non-reserved boundary is exactly "present in `common_keyword_as_identifier`?". Treat that rule as the source of truth.)

> **Migration note:** the new parser should drive its reserved/non-reserved split from a single
> table equivalent to `common_keyword_as_identifier`, because dozens of statement rules rely on
> non-reserved keywords being usable as object/column names.

---

## 4. Statement inventory

`sql_statement_body` has **53 active alternatives** (top-level SQL statements; a 54th, `define_macro_statement`, is commented out). Procedural-script statements (legal only inside `statement_list`) add **12** more forms. Organized:

### 4.1 DDL — CREATE (one rule per object kind unless noted)

| Statement | Rule |
| --- | --- |
| CREATE TABLE | `create_table_statement` |
| CREATE TABLE FUNCTION (TVF) | `create_table_function_statement` |
| CREATE EXTERNAL TABLE | `create_external_table_statement` |
| CREATE EXTERNAL TABLE FUNCTION (rejected) | `create_external_table_function_statement` |
| CREATE [MATERIALIZED/APPROX] VIEW | `create_view_statement` |
| CREATE [SEARCH/VECTOR] INDEX | `create_index_statement` |
| CREATE SCHEMA | `create_schema_statement` |
| CREATE EXTERNAL SCHEMA | `create_external_schema_statement` |
| CREATE [AGGREGATE] FUNCTION | `create_function_statement` |
| CREATE PROCEDURE | `create_procedure_statement` |
| CREATE DATABASE | `create_database_statement` |
| CREATE CONNECTION | `create_connection_statement` |
| CREATE CONSTANT | `create_constant_statement` |
| CREATE MODEL | `create_model_statement` |
| CREATE SNAPSHOT TABLE | `create_snapshot_statement` |
| CREATE PROPERTY GRAPH | `create_property_graph_statement` |
| CREATE ROW ACCESS POLICY | `create_row_access_policy_statement` |
| CREATE PRIVILEGE RESTRICTION | `create_privilege_restriction_statement` |
| CREATE \<generic entity\> | `create_entity_statement` |

### 4.2 DDL — ALTER / DROP / RENAME / UNDROP / TRUNCATE / DEFINE TABLE

| Statement | Rule |
| --- | --- |
| ALTER (table/view/schema/function/entity/privilege-restriction/row-access-policy/all-row-access-policies) | `alter_statement` |
| DROP (table/function/index/view/schema/model/snapshot/row-access-policy/privilege-restriction/entity) | `drop_statement` |
| DROP ALL ROW ACCESS POLICIES | `drop_all_row_access_policies_statement` |
| RENAME … TO … | `rename_statement` |
| UNDROP | `undrop_statement` |
| TRUNCATE TABLE | `truncate_statement` |
| DEFINE TABLE | `define_table_statement` |

### 4.3 DML

| Statement | Rule |
| --- | --- |
| INSERT / DELETE / UPDATE | `dml_statement` → `insert_statement` / `delete_statement` / `update_statement` |
| MERGE | `merge_statement` |

### 4.4 Query

| Statement | Rule |
| --- | --- |
| SELECT / set-operation / WITH-CTE query | `query_statement` → `query` |
| GQL graph query (`GRAPH …`) | `gql_statement` |

### 4.5 DCL

| Statement | Rule |
| --- | --- |
| GRANT | `grant_statement` |
| REVOKE | `revoke_statement` |

### 4.6 Transaction / batch / session

| Statement | Rule |
| --- | --- |
| BEGIN / START TRANSACTION | `begin_statement` |
| COMMIT | `commit_statement` |
| ROLLBACK | `rollback_statement` |
| SET (transaction / variable / system-var / variable-list) | `set_statement` |
| START BATCH | `start_batch_statement` |
| RUN BATCH | `run_batch_statement` |
| ABORT BATCH | `abort_batch_statement` |

### 4.7 Utility / metadata / scripting-entry / other

| Statement | Rule |
| --- | --- |
| EXPLAIN | `explain_statement` |
| DESCRIBE / DESC | `describe_statement` |
| SHOW | `show_statement` |
| ANALYZE | `analyze_statement` |
| ASSERT | `assert_statement` |
| CALL | `call_statement` |
| EXECUTE IMMEDIATE | `execute_immediate` |
| IMPORT (MODULE/PROTO) | `import_statement` |
| MODULE | `module_statement` |
| EXPORT DATA | `export_data_statement` |
| EXPORT MODEL | `export_model_statement` |
| EXPORT … METADATA | `export_metadata_statement` |
| LOAD DATA (FROM FILES) | `aux_load_data_statement` |
| CLONE DATA | `clone_data_statement` |

### 4.8 Procedural / scripting (only inside `statement_list`)

| Statement | Rule |
| --- | --- |
| BEGIN…END block | `begin_end_block` |
| IF … THEN … END IF | `if_statement` |
| CASE … END CASE | `case_statement` |
| WHILE … END WHILE | `while_statement` |
| LOOP … END LOOP | `loop_statement` |
| REPEAT … UNTIL … END REPEAT | `repeat_statement` |
| FOR … IN (query) DO … END FOR | `for_in_statement` |
| DECLARE | `variable_declaration` |
| BREAK / LEAVE | `break_statement` |
| CONTINUE / ITERATE | `continue_statement` |
| RETURN | `return_statement` |
| RAISE | `raise_statement` |

> **DEFINE MACRO** is intentionally **not implemented** — `sql_statement_body` references it only in a
> commented-out `// | define_macro_statement`, and `unterminated_sql_statement` has two error
> alternatives that explicitly reject `DEFINE MACRO`.

---

## 5. Notable structures

### 5.1 Expression precedence chain
ZetaSQL-faithful, encoded as left-recursive ANTLR rules (resolved to avoid mutual left-recursion):

```
expression
 ├─ expression_higher_prec_than_and      (everything tighter than AND)
 ├─ and_expression                       (chained AND)
 └─ expression OR expression             (lowest)

expression_higher_prec_than_and  (single rule, ordered alternatives ≈ precedence):
   literals / constructors / CASE / CAST / EXTRACT / WITH-expr / REPLACE_FIELDS /
   function calls / INTERVAL / identifier / subquery-with-keyword
   →  [expr]  /  .(path)  /  .id            (postfix field / array / member access)
   →  NOT expr
   →  expr LIKE/IN/BETWEEN/IS …  (incl. ANY/SOME/ALL, NOT DISTINCT FROM)
   →  comparative (= != <> < <= > >=)
   →  |  (STROKE / pipe-as-concat)  ^  &  ||
   →  shift (<< >>)
   →  additive (+ -)
   →  multiplicative (* /)
   →  unary (+ - ~) expr
   →  ( parenthesized-expr-not-a-query )  |  ( query )
```
`expression_maybe_parenthesized_not_a_query` mirrors the chain for RHS contexts (IN-lists, ANY/SOME/ALL).

### 5.2 Joins
- Types: `join_type` = `CROSS | FULL [OUTER] | INNER | LEFT [OUTER] | RIGHT [OUTER]`; `NATURAL` prefix via `opt_natural`; join hints `HASH | LOOKUP` via `join_hint`; per-join ANTLR `hint` (`@{…}`) allowed after `JOIN`.
- Criteria: `on_or_using_clause_list` = one-or-more of `on_clause` (`ON expr`) / `using_clause` (`USING (id…)`).
- Structure: FROM is `table_primary from_clause_contents_suffix*` (comma-cross-join or JOIN). Parenthesized joins via `table_primary: ( join )`; `join`/`join_item` resolve the mutual left-recursion of join-input.

### 5.3 Set operations
`query_set_operation_prefix` = `query_primary query_set_operation_item+`, where `set_operation_metadata` carries `[corresponding-outer] (UNION|EXCEPT|INTERSECT) [hint] (ALL|DISTINCT) [STRICT] [CORRESPONDING [BY]]`. FROM-queries directly following a set-op are rejected (must be parenthesized).

### 5.4 CTEs / WITH
`with_clause` = `WITH [RECURSIVE] aliased_query (, aliased_query)*`. `aliased_query` = `id AS ( query ) [recursion_depth_modifier]`. Recursion depth: `WITH DEPTH [AS alias] [BETWEEN n AND n | MAX n]` (`possibly_unbounded_int_literal_or_parameter` allows `UNBOUNDED`). Trailing comma after WITH and pipe after WITH are error alts.

### 5.5 Subqueries
- Scalar/derived-table: `parenthesized_query` = `( query )`.
- In expressions: `expression_subquery_with_keyword` = `ARRAY ( query )` / `EXISTS [hint] ( query )`; plus bare `( query )` and `( query ) [AS alias]` as a `query_primary`.
- As table: `table_subquery` = `( query ) [pivot/unpivot/alias]`.

### 5.6 Window / OVER
`over_clause` = `OVER window_specification`; `window_specification` = named ref or `( [name] [PARTITION BY] [ORDER BY] [frame] )`. Frame: `opt_window_frame_clause` over `frame_unit` (`ROWS|RANGE`) with `window_frame_bound` = `UNBOUNDED PRECEDING/FOLLOWING | CURRENT ROW | expr PRECEDING/FOLLOWING`. Named windows declared in `window_clause` (`WINDOW id AS spec, …`).

### 5.7 PIVOT / UNPIVOT / UNNEST / MATCH_RECOGNIZE
- **PIVOT:** `pivot_clause` = `PIVOT ( <agg exprs> FOR <expr> IN ( <values> ) )`.
- **UNPIVOT:** `unpivot_clause` = `UNPIVOT [EXCLUDE/INCLUDE NULLS] ( <cols> FOR <path> IN ( <items> ) )`.
- Applied to tables/TVFs/subqueries via `opt_pivot_or_unpivot_clause_and_alias` / `pivot_or_unpivot_clause_and_aliases` (with QUALIFY-misuse error alts).
- **UNNEST:** `unnest_expression` = `UNNEST ( <exprs-with-alias> [, array-zip-mode] )`; usable as a FROM table source (`table_path_expression_base`) and as an IN/LIKE RHS. Array element / generalized field access in FROM without UNNEST is rejected.
- **MATCH_RECOGNIZE:** present (`match_recognize_clause`, `row_pattern_expr/concatenation/factor`) with PARTITION BY / ORDER BY / MEASURES / PATTERN / DEFINE.

### 5.8 Data-type grammar
`type` = `raw_type [type_parameters] [collate]`; `raw_type` ∈ {`array_type` `ARRAY<…>`, `struct_type` `STRUCT<…>`, `range_type` `RANGE<…>`, `map_type` `MAP<k,v>`, `function_type` `FUNCTION<(args)->ret>`, `type_name` (`path_expression | INTERVAL`)}. Template brackets `< >` via `template_type_open/close` (`LT_OPERATOR`/`GT_OPERATOR`). Type parameters `( … )` with `MAX` allowed; trailing comma rejected.

### 5.9 OPTIONS / hints
- **OPTIONS:** `options_list` = `( options_entry,… ) | ()`; entry = `id <op> (expr|PROTO)` with op ∈ `= | += | -=`. Used pervasively in DDL.
- **Hints:** `hint` = `@<int>` (simple) or `@[<int>@]{ entry,… }` (body). Entries `id = expr` or `id.id = expr`. Hints attach at statement level (`statement_level_hint`), on joins, on operators within expressions (IN/EXISTS/etc.), and inside SELECT/GROUP BY/ORDER BY/PARTITION BY.

---

## 6. TODOs / FIXMEs / gaps

| Line | Marker | Note |
| --- | --- | --- |
| 58 | TODO(zp) | `define_macro_statement` commented out — DEFINE MACRO statement **not implemented** (only rejected). |
| 200 | TODO(zp) | `graph_element_pattern_filler` uses an empty production "which confused listener user" — flagged for refactor. |
| 313 | TODO(zp) | `drop_statement` TABLE/FUNCTION alt: "Refine syntax error". |
| 671 | TODO(zp) | `label: identifier` — "refine label" (script labels are just identifiers for now). |
| 1393 | FIXME(zp) | `query_without_pipe_operators`: error messages say `<KEYWORD>` literally — "Inject the keyword from original input" not done. |
| 1737 | TODO(zp) | `identifier_or_integer`: `SCRIPT_LABEL` token not yet handled. |
| 2173 | XXX(zp) | `with_group_rows`: `WITH GROUP ROWS` does not actually parse the `(query)` body — placeholder. |
| 2222 | XXX(zp) | `lambda_argument_list`: expr-kind check not enforced (any expression accepted as a lambda param list). |
| 2241 | XXX(zp) | `hint`: `@<int>` form lacks `ABORT_CHECK` handling. |
| 2333 | XXX(zp) | `with_expression`: implemented directly instead of via ZetaSQL's lookahead-transformer. |

**Visibly incomplete / placeholder behavior:**
- **Pipe syntax (`|>`)** is lexed (`PIPE_SYMBOL`) but **not parsed** — every occurrence after WITH/FROM only triggers a "Consider using pipe operator" or "pipe cannot follow WITH" **error** (`query_without_pipe_operators`). The new parser will need real pipe-operator support if/when coverage requires it; the rule name itself (`query_without_pipe_operators`) signals the omission.
- `with_group_rows` accepts no body (XXX above).
- Many alternatives exist **only to emit ZetaSQL-style syntax errors** via `p.NotifyErrorListeners(...)` (e.g. `bad_keyword_after_from_query`, the parenthesization errors in `tvf_argument`, `(SELECT…)` rejections in `cast_expression`/`sql_function_body`/`unnest_expression`, multi-variable `SET` without parens, trailing comma in type params, string/bytes concatenation rules). These are not real grammar coverage but encode error UX that the new parser may want to reproduce.

---

## 7. BigQuery-only vs Spanner-only vs Shared

GoogleSQL is a single grammar serving both dialects; classification is by which product the construct targets (BigQuery is the dominant surface; Spanner-specific constructs are clearly marked in the grammar).

### 7.1 Spanner-specific
| Construct | Rule(s) |
| --- | --- |
| INTERLEAVE IN PARENT (table) | `opt_spanner_interleave_in_parent_clause`, `opt_spanner_table_options`, `spanner_primary_key` |
| INTERLEAVE IN (index) | `spanner_index_interleave_clause` |
| NULL_FILTERED index | `opt_spanner_null_filtered` |
| Spanner ALTER COLUMN (`<schema> [NOT NULL] [generated/default] [OPTIONS]`) | `spanner_alter_column_action`, `spanner_generated_or_default` |
| Spanner SET ON DELETE | `spanner_set_on_delete_action` |
| Generated-column `AS (expr) STORED` (Spanner form) | `spanner_generated_or_default` |
| Sequences (`SEQUENCE` keyword, `sequence_arg`) | `sequence_arg` (+ `SEQUENCE_SYMBOL`) — Spanner sequences |
| Role/grantee-based GRANT/REVOKE & ROW ACCESS POLICY | `grant_statement`, `revoke_statement`, `create_row_access_policy_statement`, `grant_to_clause` (Spanner fine-grained access control; also BigQuery row-access policies — shared-ish but Spanner-leaning) |
| `REPLICA` / `CHANGE STREAM`-style generic sub-entities | `sub_entity_type_identifier` (`REPLICA`), generic-entity machinery (CHANGE STREAM is modeled as a generic entity, not a dedicated rule) |
| Statement hints `@{...}` | `hint`, `statement_level_hint` (Spanner query hints; BigQuery also uses `@{}` hints, so partly shared) |

> Note: Spanner CHANGE STREAM, named schemas, and similar are not first-class rules — they ride on
> the **generic entity** mechanism (`create_entity_statement` / `alter_statement` generic-entity alts /
> `generic_entity_type`), which is the grammar's extensibility hook.

### 7.2 BigQuery-specific
| Construct | Rule(s) |
| --- | --- |
| Procedural scripting (IF/CASE/WHILE/LOOP/REPEAT/FOR-IN/DECLARE/BREAK/CONTINUE/RETURN/RAISE, BEGIN…END, EXCEPTION) | `if_statement`, `case_statement`, `while_statement`, `loop_statement`, `repeat_statement`, `for_in_statement`, `variable_declaration`, `break_statement`, `continue_statement`, `return_statement`, `raise_statement`, `begin_end_block`, `opt_exception_handler` |
| EXPORT DATA / MODEL / METADATA | `export_data_statement`, `export_model_statement`, `export_metadata_statement` |
| LOAD DATA / CLONE DATA | `aux_load_data_statement`, `clone_data_statement` |
| ML / model | `create_model_statement`, `opt_input_output_clause`, `opt_transform_clause`, `model_clause`, `export_model_statement` |
| BigQuery ML/remote functions (REMOTE, CONNECTION, LANGUAGE) | `remote_with_connection_clause`, `connection_clause`, `language` |
| Differential privacy (DIFFERENTIAL_PRIVACY, EPSILON, CLAMPED, PRIVACY_UNIT_COLUMN, MAX_GROUPS_CONTRIBUTED, WITH REPORT) | `clamped_between_modifier`, `with_report_modifier`, tokens `DIFFERENTIAL_PRIVACY/EPSILON/PRIVACY_UNIT_COLUMN/MAX_GROUPS_CONTRIBUTED` |
| Snapshots | `create_snapshot_statement` |
| Search/Vector indexes | `index_type` (`SEARCH`/`VECTOR`) |
| Materialized / Approx / Recursive views; REPLICA OF | `create_view_statement`, `query_or_replica_source` |
| Dashed & slashed paths (`project-id.dataset.table`, resource paths) | `dashed_path_expression`, `slashed_path_expression` |
| BigQuery `#` comments | lexer `POUND_COMMENT` |
| Property graphs + GQL (`GRAPH …`, CREATE PROPERTY GRAPH) | all `graph_*`/`gql_*` rules, `create_property_graph_statement` (GoogleSQL graph extension — BigQuery/Spanner Graph) |
| MERGE, PIVOT/UNPIVOT, QUALIFY, MATCH_RECOGNIZE, INTERVAL/RANGE/JSON/NUMERIC literals, TABLESAMPLE | predominantly BigQuery query-syntax features (some now also in Spanner) |
| Capacity/reservation/assignment, DDL on connections | `create_connection_statement` (BigQuery connections) |
| AT SYSTEM TIME (time travel) | `opt_at_system_time` (BigQuery `FOR SYSTEM_TIME AS OF`) |

### 7.3 Shared core
SELECT and the full query stack (`query`, `select`, `from_clause`, joins, set operations, CTEs, ORDER/LIMIT/GROUP/HAVING/WINDOW), DML (`insert/update/delete`), the entire expression grammar (§2.17–2.18), data-type grammar (§2.20), literals (§2.21), identifiers/paths/parameters (§2.22), basic DDL (CREATE/ALTER/DROP TABLE/VIEW/SCHEMA/FUNCTION/PROCEDURE/INDEX, constraints, primary/foreign keys, OPTIONS), transactions (`begin/commit/rollback/set`), `CALL`, `EXECUTE IMMEDIATE`, hints (`@{}`/`@n`), and the generic-entity extensibility mechanism. These constitute the bulk of the 625 rules and are dialect-neutral.

---

## Appendix — cross-reference vs `bq-parser/BigQueryParser.g4` (HINT only)

The secondary `BigQueryParser.g4` is a **small, WIP** grammar (~25 parser rules) and is **not** the migration target. Coverage comparison:

**Covered by GoogleSQL but absent/stubbed in BigQueryParser (i.e. BigQueryParser gaps):**
- All DDL (CREATE/ALTER/DROP/RENAME/UNDROP/TRUNCATE), all DML except SELECT, MERGE, all DCL, all scripting, EXPORT/LOAD/CLONE, GQL/property-graph, PIVOT/UNPIVOT/MATCH_RECOGNIZE, hints, OPTIONS, full type grammar, typed literals, parameters/system-variables, dashed/slashed paths, set-operation metadata (CORRESPONDING/STRICT), window frames, CTE recursion depth.
- BigQueryParser's `field_path`, `array_path`, `window_name`, `window_definition`, `bool_expression` are **empty/`WIP` stubs**; its expression rule is a flat single rule (no real precedence) and it explicitly notes STRUCT/ARRAY handling is incomplete.

**Present in BigQueryParser but worth noting vs GoogleSQL:**
- `ORDINAL` / `SAFE_OFFSET` / `SAFE_ORDINAL` array-access keywords and `TREAT`, `FETCH`, `WITHIN`, `LATERAL`, `CONTAINS`, `GROUPS`, `ESCAPE` appear in BigQueryParser's keyword list. In GoogleSQL these are **not** dedicated tokens — `OFFSET`/`ORDINAL` array indexing is handled as ordinary function-style/path access, and `TREAT`/`FETCH`/`LATERAL`/`WITHIN`/`CONTAINS`/`GROUPS`/`ESCAPE` are **not** modeled. If full BigQuery coverage is required these are potential gaps in **both** grammars relative to BigQuery's published lexical spec, but GoogleSQL is far more complete overall.
- BigQueryLexer keeps `RAW_BYTE_STRING` (`rb`/`br`) and triple-quoted strings — both also covered (more rigorously, incl. unclosed-token error states) by GoogleSQLLexer.

**Conclusion:** the GoogleSQL grammar is the strict superset and the only viable coverage target;
`bq-parser` is useful only as a sanity hint for a handful of BigQuery array-access / lexical keywords
(`OFFSET/ORDINAL/SAFE_OFFSET/SAFE_ORDINAL`, `TREAT`, `WITHIN`, `LATERAL`) that may warrant explicit
test cases when validating the new parser against real BigQuery.
