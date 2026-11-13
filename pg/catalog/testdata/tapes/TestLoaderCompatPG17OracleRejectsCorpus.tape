-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestLoaderCompatPG17OracleRejectsCorpus default#1
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_numeric_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_numeric_fk_references_integer_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_numeric_fk_references_integer_pk", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id numeric REFERENCES parent(id)
| 				);
| 			
< error
| 42804 foreign key constraint "child_parent_id_fkey" cannot be implemented
| detail: Key columns "parent_id" and "id" are of incompatible types: numeric and integer.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_numeric_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_composite_fk_rejects_one_incompatible_column_pair" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_composite_fk_rejects_one_incompatible_column_pair"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_composite_fk_rejects_one_incompatible_column_pair", public;
|
| 				CREATE TABLE parent (id integer, tenant_id integer, PRIMARY KEY (id, tenant_id));
| 				CREATE TABLE child (
| 					parent_id numeric,
| 					tenant_id integer,
| 					FOREIGN KEY (parent_id, tenant_id) REFERENCES parent(id, tenant_id)
| 				);
| 			
< error
| 42804 foreign key constraint "child_parent_id_tenant_id_fkey" cannot be implemented
| detail: Key columns "parent_id" and "id" are of incompatible types: numeric and integer.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_composite_fk_rejects_one_incompatible_column_pair" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_wrong_child_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_wrong_child_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_wrong_child_table", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 				CREATE TABLE other (id integer);
| 				CREATE INDEX parent_id_idx ON ONLY parent (id);
| 				CREATE INDEX other_id_idx ON other (id);
| 				ALTER INDEX parent_id_idx ATTACH PARTITION other_id_idx;
| 			
< error
| 55000 cannot attach index "other_id_idx" as a partition of index "parent_id_idx"
| detail: Index "other_id_idx" is not an index on any partition of table "parent".
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_wrong_child_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_nonpartitioned_parent_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_nonpartitioned_parent_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_nonpartitioned_parent_table", public;
|
| 				CREATE TABLE parent (id integer);
| 				CREATE TABLE child (id integer);
| 				CREATE INDEX parent_id_idx ON parent (id);
| 				CREATE INDEX child_id_idx ON child (id);
| 				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
| 			
< error
| 42809 ALTER action ATTACH PARTITION cannot be performed on relation "parent_id_idx"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_alter_index_attach_rejects_nonpartitioned_parent_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_partitioned_zero_column_table_expression_key_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_partitioned_zero_column_table_expression_key_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_partitioned_zero_column_table_expression_key_rejected", public;
| CREATE TABLE zc () PARTITION BY LIST (());
< error
| 42601 syntax error at or near ")"
| position: 164
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_partitioned_zero_column_table_expression_key_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_explicit_variadic_rejects_non_array_argument" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_explicit_variadic_rejects_non_array_argument"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_explicit_variadic_rejects_non_array_argument", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_bad_variadic AS
| 					SELECT count_variadic(VARIADIC 1) AS value;
| 			
< error
| 42883 function count_variadic(integer) does not exist
| position: 295
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_explicit_variadic_rejects_non_array_argument" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_user_variadic_function_zero_args_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_user_variadic_function_zero_args_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_user_variadic_function_zero_args_rejected", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_zero AS
| 					SELECT count_variadic() AS value;
| 			
< error
| 42883 function count_variadic() does not exist
| position: 293
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_user_variadic_function_zero_args_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_recursive_cte_self_reference_without_union_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_recursive_cte_self_reference_without_union_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_recursive_cte_self_reference_without_union_rejected", public;
|
| 				CREATE VIEW v_bad_recursive AS
| 					WITH RECURSIVE r(n) AS (SELECT n + 1 FROM r)
| 					SELECT n FROM r;
| 			
< error
| 42P19 recursive query "r" does not have the form non-recursive-term UNION [ALL] recursive-term
| position: 179
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_recursive_cte_self_reference_without_union_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_cte_column_aliases_more_than_output_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_cte_column_aliases_more_than_output_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_cte_column_aliases_more_than_output_rejected", public;
|
| 				CREATE VIEW v_bad_cte_alias AS
| 					WITH c(a, b) AS (SELECT 1)
| 					SELECT a FROM c;
| 			
< error
| 42P10 WITH query "c" has 1 columns available but 2 columns specified
| position: 162
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_cte_column_aliases_more_than_output_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_non_lateral_subquery_references_prior_from_item_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_non_lateral_subquery_references_prior_from_item_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_non_lateral_subquery_references_prior_from_item_rejected", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_bad_non_lateral AS
| 					SELECT s.id
| 					FROM base_items b,
| 					     (SELECT b.id AS id) s;
| 			
< error
| 42P01 invalid reference to FROM-clause entry for table "b"
| position: 267
| detail: There is an entry for table "b", but it cannot be referenced from this part of the query.
| hint: To reference that table, you must mark this subquery with LATERAL.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_non_lateral_subquery_references_prior_from_item_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_check_constraint_subquery_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_check_constraint_subquery_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_check_constraint_subquery_rejected", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE TABLE t (
| 					id integer,
| 					CHECK (id IN (SELECT id FROM base_items))
| 				);
| 			
< error
| 0A000 cannot use subquery in check constraint
| position: 202
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_check_constraint_subquery_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_generated_column_subquery_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_generated_column_subquery_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_generated_column_subquery_rejected", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE TABLE t (
| 					id integer,
| 					generated integer GENERATED ALWAYS AS ((SELECT max(id) FROM base_items)) STORED
| 				);
| 			
< error
| 0A000 cannot use subquery in column generation expression
| position: 231
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_generated_column_subquery_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_view_column_aliases_more_than_output_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_view_column_aliases_more_than_output_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_view_column_aliases_more_than_output_rejected", public;
| CREATE VIEW v_bad_alias(a, b) AS SELECT 1;
< error
| 42601 CREATE VIEW specifies more column names than columns
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_view_column_aliases_more_than_output_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_view_duplicate_output_column_names_rejected" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_view_duplicate_output_column_names_rejected"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_view_duplicate_output_column_names_rejected", public;
| CREATE VIEW v_dupe_cols AS SELECT 1 AS id, 2 AS id;
< error
| 42701 column "id" specified more than once
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_view_duplicate_output_column_names_rejected" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_drop_primary_key_referenced_by_fk_without_cascade" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_drop_primary_key_referenced_by_fk_without_cascade"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_drop_primary_key_referenced_by_fk_without_cascade", public;
|
| 				CREATE TABLE drop_pk_parent (id integer CONSTRAINT drop_pk_parent_key PRIMARY KEY);
| 				CREATE TABLE drop_pk_child (parent_id integer CONSTRAINT drop_pk_child_fk REFERENCES drop_pk_parent (id));
| 				ALTER TABLE drop_pk_parent DROP CONSTRAINT drop_pk_parent_key;
| 			
< error
| 2BP01 cannot drop constraint drop_pk_parent_key on table drop_pk_parent because other objects depend on it
| detail: constraint drop_pk_child_fk on table drop_pk_child depends on index drop_pk_parent_key
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_drop_primary_key_referenced_by_fk_without_cascade" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_drop_unique_referenced_by_fk_without_cascade" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oraclerejectscorpus_drop_unique_referenced_by_fk_without_cascade"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oraclerejectscorpus_drop_unique_referenced_by_fk_without_cascade", public;
|
| 				CREATE TABLE drop_uq_parent (code text CONSTRAINT drop_uq_parent_key UNIQUE);
| 				CREATE TABLE drop_uq_child (parent_code text CONSTRAINT drop_uq_child_fk REFERENCES drop_uq_parent (code));
| 				ALTER TABLE drop_uq_parent DROP CONSTRAINT drop_uq_parent_key;
| 			
< error
| 2BP01 cannot drop constraint drop_uq_parent_key on table drop_uq_parent because other objects depend on it
| detail: constraint drop_uq_child_fk on table drop_uq_child depends on index drop_uq_parent_key
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oraclerejectscorpus_drop_unique_referenced_by_fk_without_cascade" CASCADE
< ok 0 -
