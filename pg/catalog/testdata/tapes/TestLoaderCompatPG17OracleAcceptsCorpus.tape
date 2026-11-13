-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestLoaderCompatPG17OracleAcceptsCorpus default#1
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_zero_column_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_zero_column_table", public;
| CREATE TABLE zc ();
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_in_schema" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_zero_column_table_in_schema"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_zero_column_table_in_schema", public;
|
| 				CREATE SCHEMA loader_s1;
| 				CREATE TABLE loader_s1.zc ();
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_in_schema" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_comment" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_zero_column_table_comment"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_zero_column_table_comment", public;
|
| 				CREATE TABLE zc ();
| 				COMMENT ON TABLE zc IS 'zero columns';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_comment" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_grant" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_zero_column_table_grant"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_zero_column_table_grant", public;
|
| 				CREATE TABLE zc ();
| 				GRANT SELECT ON zc TO PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_grant" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_view_star" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_zero_column_table_view_star"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_zero_column_table_view_star", public;
|
| 				CREATE TABLE zc ();
| 				CREATE VIEW v_zc AS SELECT * FROM zc;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_zero_column_table_view_star" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer,
| 					FOREIGN KEY (parent_id) REFERENCES parent(id) MATCH SIMPLE
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_multicolumn" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_multicolumn"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_multicolumn", public;
|
| 				CREATE TABLE parent (a integer, b integer, PRIMARY KEY (a, b));
| 				CREATE TABLE child (
| 					a integer,
| 					b integer,
| 					FOREIGN KEY (a, b) REFERENCES parent(a, b) MATCH SIMPLE
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_multicolumn" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_update_cascade" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_update_cascade"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_update_cascade", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id) MATCH SIMPLE ON UPDATE CASCADE
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_update_cascade" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_delete_set_null" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_delete_set_null"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_delete_set_null", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id) MATCH SIMPLE ON DELETE SET NULL
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_on_delete_set_null" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_deferrable" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_deferrable"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_deferrable", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id) MATCH SIMPLE DEFERRABLE INITIALLY DEFERRED
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_foreign_key_match_simple_deferrable" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_function_argument_names", public;
|
| 				CREATE FUNCTION f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				COMMENT ON FUNCTION f(a integer, b text) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_schema_qualified_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_schema_qualified_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_schema_qualified_function_argument_names", public;
|
| 				CREATE SCHEMA funcs;
| 				CREATE FUNCTION funcs.f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				COMMENT ON FUNCTION funcs.f(a integer, b text) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_schema_qualified_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_quoted_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_function_quoted_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_function_quoted_argument_names", public;
|
| 				CREATE FUNCTION f("select" integer, "from" text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				COMMENT ON FUNCTION f("select" integer, "from" text) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_quoted_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_qualified_argument_types" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_function_qualified_argument_types"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_function_qualified_argument_types", public;
|
| 				CREATE FUNCTION f(a pg_catalog.int4, b pg_catalog.text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				COMMENT ON FUNCTION f(a pg_catalog.int4, b pg_catalog.text) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_qualified_argument_types" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_variadic_function_identity" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_variadic_function_identity"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_variadic_function_identity", public;
|
| 				CREATE FUNCTION f(VARIADIC nums integer[]) RETURNS integer
| 					LANGUAGE sql AS 'SELECT cardinality(nums)';
| 				COMMENT ON FUNCTION f(VARIADIC nums integer[]) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_variadic_function_identity" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_grant_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_grant_function_argument_names", public;
|
| 				CREATE FUNCTION f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				GRANT EXECUTE ON FUNCTION f(a integer, b text) TO PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_revoke_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_revoke_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_revoke_function_argument_names", public;
|
| 				CREATE FUNCTION f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				GRANT EXECUTE ON FUNCTION f(a integer, b text) TO PUBLIC;
| 				REVOKE EXECUTE ON FUNCTION f(a integer, b text) FROM PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_revoke_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_function_argument_names", public;
|
| 				CREATE FUNCTION f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				ALTER FUNCTION f(a integer, b text) RENAME TO f2;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_function_argument_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_function_argument_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_function_argument_names", public;
|
| 				CREATE FUNCTION f(a integer, b text) RETURNS integer
| 					LANGUAGE sql AS 'SELECT 1';
| 				DROP FUNCTION f(a integer, b text);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_function_argument_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_function_identity_in_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_function_identity_in_mode", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				COMMENT ON FUNCTION f(IN a integer) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_identity_inout_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_comment_function_identity_inout_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_comment_function_identity_inout_mode", public;
|
| 				CREATE FUNCTION f(INOUT a integer)
| 					LANGUAGE sql AS 'SELECT $1';
| 				COMMENT ON FUNCTION f(INOUT a integer) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_comment_function_identity_inout_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_grant_function_identity_in_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_grant_function_identity_in_mode", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				GRANT EXECUTE ON FUNCTION f(IN a integer) TO PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_inout_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_grant_function_identity_inout_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_grant_function_identity_inout_mode", public;
|
| 				CREATE FUNCTION f(INOUT a integer)
| 					LANGUAGE sql AS 'SELECT $1';
| 				GRANT EXECUTE ON FUNCTION f(INOUT a integer) TO PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_inout_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_variadic_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_grant_function_identity_variadic_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_grant_function_identity_variadic_mode", public;
|
| 				CREATE FUNCTION f(VARIADIC a integer[]) RETURNS integer
| 					LANGUAGE sql AS 'SELECT cardinality($1)';
| 				GRANT EXECUTE ON FUNCTION f(VARIADIC a integer[]) TO PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_grant_function_identity_variadic_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_revoke_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_revoke_function_identity_in_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_revoke_function_identity_in_mode", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				GRANT EXECUTE ON FUNCTION f(IN a integer) TO PUBLIC;
| 				REVOKE EXECUTE ON FUNCTION f(IN a integer) FROM PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_revoke_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_owner" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_owner"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_owner", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				ALTER FUNCTION f(IN a integer) OWNER TO CURRENT_USER;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_owner" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_options" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_options"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_options", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				ALTER FUNCTION f(IN a integer) IMMUTABLE;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_function_identity_in_mode_options" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_function_identity_in_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_function_identity_in_mode", public;
|
| 				CREATE FUNCTION f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				DROP FUNCTION f(IN a integer);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_function_identity_in_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_procedure_identity_inout_mode" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_procedure_identity_inout_mode"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_procedure_identity_inout_mode", public;
|
| 				CREATE PROCEDURE p(INOUT a integer)
| 					LANGUAGE sql AS 'SELECT $1';
| 				DROP PROCEDURE p(INOUT a integer);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_procedure_identity_inout_mode" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_routine_identity_argument_name" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_routine_identity_argument_name"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_routine_identity_argument_name", public;
|
| 				CREATE FUNCTION r(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				DROP ROUTINE r(a integer);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_routine_identity_argument_name" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_types" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_types"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_types", public;
|
| 				CREATE FUNCTION f(a pg_catalog.int4) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				COMMENT ON FUNCTION f(a pg_catalog.int4) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_types" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_array_type" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_function_identity_array_type"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_function_identity_array_type", public;
|
| 				CREATE FUNCTION f(a integer[]) RETURNS integer
| 					LANGUAGE sql AS 'SELECT cardinality($1)';
| 				COMMENT ON FUNCTION f(a integer[]) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_array_type" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_quoted_function_name" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_function_identity_quoted_function_name"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_function_identity_quoted_function_name", public;
|
| 				CREATE FUNCTION "select"(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				COMMENT ON FUNCTION "select"(a integer) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_quoted_function_name" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_function_name" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_function_name"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_function_name", public;
|
| 				CREATE SCHEMA ident_s;
| 				CREATE FUNCTION ident_s.f(a integer) RETURNS integer
| 					LANGUAGE sql AS 'SELECT $1';
| 				COMMENT ON FUNCTION ident_s.f(a integer) IS 'comment';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_function_identity_schema_qualified_function_name" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_plpgsql" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_returns_table_plpgsql"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_returns_table_plpgsql", public;
|
| 				CREATE FUNCTION f()
| 				RETURNS TABLE(id integer, name text)
| 				LANGUAGE plpgsql
| 				AS $$
| 				BEGIN
| 					RETURN QUERY SELECT 1, 'one'::text;
| 				END
| 				$$;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_plpgsql" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_single_column_plpgsql" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_returns_table_single_column_plpgsql"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_returns_table_single_column_plpgsql", public;
|
| 				CREATE FUNCTION f()
| 				RETURNS TABLE(id integer)
| 				LANGUAGE plpgsql
| 				AS $$
| 				BEGIN
| 					RETURN QUERY SELECT 1;
| 				END
| 				$$;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_single_column_plpgsql" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_qualified_column_types" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_returns_table_qualified_column_types"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_returns_table_qualified_column_types", public;
|
| 				CREATE FUNCTION f()
| 				RETURNS TABLE(id pg_catalog.int4, name pg_catalog.text)
| 				LANGUAGE plpgsql
| 				AS $$
| 				BEGIN
| 					RETURN QUERY SELECT 1, 'one'::text;
| 				END
| 				$$;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_returns_table_qualified_column_types" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id bigint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_integer_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_integer_pk", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id smallint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_integer_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_bigint_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_integer_fk_references_bigint_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_integer_fk_references_bigint_pk", public;
|
| 				CREATE TABLE parent (id bigint PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_bigint_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_smallint_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_integer_fk_references_smallint_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_integer_fk_references_smallint_pk", public;
|
| 				CREATE TABLE parent (id smallint PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_smallint_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_bigint_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_bigint_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_bigint_pk", public;
|
| 				CREATE TABLE parent (id bigint PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id smallint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_smallint_fk_references_bigint_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_smallint_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_smallint_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_smallint_pk", public;
|
| 				CREATE TABLE parent (id smallint PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id bigint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_smallint_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_numeric_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_integer_fk_references_numeric_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_integer_fk_references_numeric_pk", public;
|
| 				CREATE TABLE parent (id numeric PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_integer_fk_references_numeric_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_varchar_fk_references_text_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_varchar_fk_references_text_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_varchar_fk_references_text_pk", public;
|
| 				CREATE TABLE parent (id text PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id varchar REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_varchar_fk_references_text_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_text_fk_references_varchar_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_text_fk_references_varchar_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_text_fk_references_varchar_pk", public;
|
| 				CREATE TABLE parent (id varchar PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id text REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_text_fk_references_varchar_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bpchar_fk_references_text_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_bpchar_fk_references_text_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_bpchar_fk_references_text_pk", public;
|
| 				CREATE TABLE parent (id text PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id character REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bpchar_fk_references_text_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_text_fk_references_bpchar_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_text_fk_references_bpchar_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_text_fk_references_bpchar_pk", public;
|
| 				CREATE TABLE parent (id character PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id text REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_text_fk_references_bpchar_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_domain_fk_references_base_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_domain_fk_references_base_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_domain_fk_references_base_pk", public;
|
| 				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id positive_int REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_domain_fk_references_base_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_base_fk_references_domain_pk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_base_fk_references_domain_pk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_base_fk_references_domain_pk", public;
|
| 				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
| 				CREATE TABLE parent (id positive_int PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_base_fk_references_domain_pk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_same_domain_fk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_same_domain_fk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_same_domain_fk", public;
|
| 				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
| 				CREATE TABLE parent (id positive_int PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id positive_int REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_same_domain_fk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_different_domain_same_base_fk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_different_domain_same_base_fk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_different_domain_same_base_fk", public;
|
| 				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
| 				CREATE DOMAIN other_int AS integer CHECK (VALUE > 0);
| 				CREATE TABLE parent (id positive_int PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id other_int REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_different_domain_same_base_fk" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk_composite" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk_composite"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk_composite", public;
|
| 				CREATE TABLE parent (id integer, tenant_id integer, PRIMARY KEY (id, tenant_id));
| 				CREATE TABLE child (
| 					parent_id bigint,
| 					tenant_id bigint,
| 					FOREIGN KEY (parent_id, tenant_id) REFERENCES parent(id, tenant_id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_bigint_fk_references_integer_pk_composite" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_cross_type" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_cross_type"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_cross_type", public;
|
| 				CREATE TABLE parent (id integer UNIQUE);
| 				CREATE TABLE child (
| 					parent_id bigint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_cross_type" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_with_include_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_with_include_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_with_include_columns", public;
|
| 				CREATE TABLE parent (id integer, extra text);
| 				CREATE UNIQUE INDEX parent_id_idx ON parent (id) INCLUDE (extra);
| 				CREATE TABLE child (
| 					parent_id bigint REFERENCES parent(id)
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_fk_references_unique_index_with_include_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_temporary_zero_column_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_temporary_zero_column_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_temporary_zero_column_table", public;
| CREATE TEMPORARY TABLE zc ();
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_temporary_zero_column_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_unlogged_zero_column_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_unlogged_zero_column_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_unlogged_zero_column_table", public;
| CREATE UNLOGGED TABLE zc ();
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_unlogged_zero_column_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_table_like_zero_column_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_create_table_like_zero_column_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_create_table_like_zero_column_table", public;
|
| 				CREATE TABLE zc ();
| 				CREATE TABLE zc_like (LIKE zc);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_table_like_zero_column_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_table_as_zero_output_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_create_table_as_zero_output_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_create_table_as_zero_output_columns", public;
| CREATE TABLE zc AS SELECT;
< ok 1 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_table_as_zero_output_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_view_zero_output_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_create_view_zero_output_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_create_view_zero_output_columns", public;
| CREATE VIEW v_zc AS SELECT;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_create_view_zero_output_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_inherited_zero_column_parent" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_inherited_zero_column_parent"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_inherited_zero_column_parent", public;
|
| 				CREATE TABLE parent ();
| 				CREATE TABLE child () INHERITS (parent);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_inherited_zero_column_parent" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_zero_column_table_add_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_zero_column_table_add_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_zero_column_table_add_column", public;
|
| 				CREATE TABLE zc ();
| 				ALTER TABLE zc ADD COLUMN id integer;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_zero_column_table_add_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_table_drop_column_to_zero_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_table_drop_column_to_zero_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_table_drop_column_to_zero_columns", public;
|
| 				CREATE TABLE zc (id integer);
| 				ALTER TABLE zc DROP COLUMN id;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_table_drop_column_to_zero_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_wrong_schema_parent_qualified" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_wrong_schema_parent_qualified"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_wrong_schema_parent_qualified", public;
|
| 				CREATE SCHEMA pidx;
| 				CREATE TABLE pidx.parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE pidx.child PARTITION OF pidx.parent FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_id_idx ON ONLY pidx.parent (id);
| 				CREATE INDEX child_id_idx ON pidx.child (id);
| 				ALTER INDEX pidx.parent_id_idx ATTACH PARTITION pidx.child_id_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_wrong_schema_parent_qualified" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_parent_index_first" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_parent_index_first"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_parent_index_first", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE INDEX parent_id_idx ON ONLY parent (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_parent_index_first" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_child_table_first" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_child_table_first"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_child_table_first", public;
|
| 				CREATE TABLE child (id integer);
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				ALTER TABLE parent ATTACH PARTITION child FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_id_idx ON ONLY parent (id);
| 				CREATE INDEX child_id_idx ON child (id);
| 				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_index_attach_child_table_first" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_table_attach_records_dependency" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partition_table_attach_records_dependency"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partition_table_attach_records_dependency", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child (id integer);
| 				ALTER TABLE parent ATTACH PARTITION child FOR VALUES FROM (0) TO (10);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_table_attach_records_dependency" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_unique_index_attach" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_unique_index_attach"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_unique_index_attach", public;
|
| 				CREATE TABLE parent (id integer NOT NULL) PARTITION BY RANGE (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 				CREATE UNIQUE INDEX parent_id_idx ON ONLY parent (id);
| 				CREATE UNIQUE INDEX child_id_idx ON child (id);
| 				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_unique_index_attach" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_expression_index_attach" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_expression_index_attach"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_expression_index_attach", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_expr_idx ON ONLY parent ((id + 1));
| 				CREATE INDEX child_expr_idx ON child ((id + 1));
| 				ALTER INDEX parent_expr_idx ATTACH PARTITION child_expr_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_expression_index_attach" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_partial_index_attach" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partitioned_partial_index_attach"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partitioned_partial_index_attach", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_partial_idx ON ONLY parent (id) WHERE id > 0;
| 				CREATE INDEX child_partial_idx ON child (id) WHERE id > 0;
| 				ALTER INDEX parent_partial_idx ATTACH PARTITION child_partial_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partitioned_partial_index_attach" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_default_bound" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partition_default_bound"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partition_default_bound", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child_default PARTITION OF parent DEFAULT;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_default_bound" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_table_detach" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_partition_table_detach"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_partition_table_detach", public;
|
| 				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
| 				ALTER TABLE parent DETACH PARTITION child;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_partition_table_detach" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_generated_column_expression" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_generated_column_expression"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_generated_column_expression", public;
|
| 				CREATE TABLE t (
| 					a integer,
| 					b integer GENERATED ALWAYS AS (a + 1) STORED
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_generated_column_expression" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_identity_column_dependency" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_identity_column_dependency"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_identity_column_dependency", public;
|
| 				CREATE TABLE t (
| 					id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_identity_column_dependency" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_column_default_function_dependency" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_column_default_function_dependency"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_column_default_function_dependency", public;
|
| 				CREATE FUNCTION next_code() RETURNS integer LANGUAGE sql AS 'SELECT 1';
| 				CREATE TABLE t (
| 					code integer DEFAULT next_code()
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_column_default_function_dependency" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_column_default_sequence_dependency" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_column_default_sequence_dependency"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_column_default_sequence_dependency", public;
|
| 				CREATE SEQUENCE s;
| 				CREATE TABLE t (
| 					id integer DEFAULT nextval('s')
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_column_default_sequence_dependency" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_check_constraint_function_dependency" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_check_constraint_function_dependency"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_check_constraint_function_dependency", public;
|
| 				CREATE FUNCTION is_positive(x integer) RETURNS boolean LANGUAGE sql AS 'SELECT $1 > 0';
| 				CREATE TABLE t (
| 					x integer CHECK (is_positive(x))
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_check_constraint_function_dependency" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_exclusion_constraint_resolves_operator_class" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_exclusion_constraint_resolves_operator_class"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_exclusion_constraint_resolves_operator_class", public;
|
| 				CREATE TABLE t (r int4range);
| 				ALTER TABLE t ADD CONSTRAINT t_r_excl EXCLUDE USING gist (r WITH &&);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_exclusion_constraint_resolves_operator_class" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_unique_nulls_not_distinct" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_unique_nulls_not_distinct"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_unique_nulls_not_distinct", public;
|
| 				CREATE TABLE t (x integer);
| 				ALTER TABLE t ADD CONSTRAINT t_x_key UNIQUE NULLS NOT DISTINCT (x);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_unique_nulls_not_distinct" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_deferrable_unique_constraint" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_deferrable_unique_constraint"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_deferrable_unique_constraint", public;
|
| 				CREATE TABLE t (
| 					x integer UNIQUE DEFERRABLE INITIALLY DEFERRED
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_deferrable_unique_constraint" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_deferrable_foreign_key" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_deferrable_foreign_key"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_deferrable_foreign_key", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (
| 					parent_id integer REFERENCES parent(id) DEFERRABLE INITIALLY DEFERRED
| 				);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_deferrable_foreign_key" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_not_valid_foreign_key" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_not_valid_foreign_key"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_not_valid_foreign_key", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (parent_id integer);
| 				ALTER TABLE child ADD CONSTRAINT child_parent_fk
| 					FOREIGN KEY (parent_id) REFERENCES parent(id) NOT VALID;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_not_valid_foreign_key" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_validate_constraint" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_validate_constraint"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_validate_constraint", public;
|
| 				CREATE TABLE parent (id integer PRIMARY KEY);
| 				CREATE TABLE child (parent_id integer);
| 				ALTER TABLE child ADD CONSTRAINT child_parent_fk
| 					FOREIGN KEY (parent_id) REFERENCES parent(id) NOT VALID;
| 				ALTER TABLE child VALIDATE CONSTRAINT child_parent_fk;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_validate_constraint" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_constraint_comment" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_constraint_comment"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_constraint_comment", public;
|
| 				CREATE TABLE t (x integer CONSTRAINT x_positive CHECK (x > 0));
| 				COMMENT ON CONSTRAINT x_positive ON t IS 'positive';
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_constraint_comment" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_concat_ws_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_concat_ws_variadic_builtin", public;
|
| 				CREATE TABLE items (a text, b integer);
| 				CREATE VIEW v_items AS
| 					SELECT concat_ws('-', a, b) AS label
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_unknown_literals" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_concat_ws_unknown_literals"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_concat_ws_unknown_literals", public;
|
| 				CREATE VIEW v_items AS
| 					SELECT concat_ws('-', 'a', 'b', 'c') AS label;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_unknown_literals" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_mixed_types" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_concat_ws_mixed_types"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_concat_ws_mixed_types", public;
|
| 				CREATE TABLE items (a text, b integer, c uuid, d timestamp);
| 				CREATE VIEW v_items AS
| 					SELECT concat_ws('-', a, b, c, d) AS label
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_ws_mixed_types" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_object_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_object_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_object_variadic_builtin", public;
|
| 				CREATE TABLE items (id integer, name text);
| 				CREATE VIEW v_items AS
| 					SELECT jsonb_build_object('id', id, 'name', name) AS payload
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_object_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_concat_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_concat_variadic_builtin", public;
|
| 				CREATE TABLE items (a text, b integer, c timestamp);
| 				CREATE VIEW v_items AS
| 					SELECT concat(a, b, c) AS label
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_concat_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_format_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_format_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_format_variadic_builtin", public;
|
| 				CREATE TABLE items (a text, b integer);
| 				CREATE VIEW v_items AS
| 					SELECT format('%s %s', a, b) AS label
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_format_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_json_build_object_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_json_build_object_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_json_build_object_variadic_builtin", public;
|
| 				CREATE TABLE items (id integer, name text);
| 				CREATE VIEW v_items AS
| 					SELECT json_build_object('id', id, 'name', name) AS payload
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_json_build_object_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_json_build_array_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_json_build_array_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_json_build_array_variadic_builtin", public;
|
| 				CREATE TABLE items (id integer, name text);
| 				CREATE VIEW v_items AS
| 					SELECT json_build_array(id, name) AS payload
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_json_build_array_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_array_variadic_builtin" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_array_variadic_builtin"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_array_variadic_builtin", public;
|
| 				CREATE TABLE items (id integer, name text);
| 				CREATE VIEW v_items AS
| 					SELECT jsonb_build_array(id, name) AS payload
| 					FROM items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_build_array_variadic_builtin" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_function_default_argument" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_function_default_argument"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_function_default_argument", public;
|
| 				CREATE FUNCTION add_default(a integer, b integer DEFAULT 1)
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT a + b';
| 				CREATE VIEW v_default_arg AS
| 					SELECT add_default(2) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_function_default_argument" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_expanded_arguments" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_expanded_arguments"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_expanded_arguments", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_arg AS
| 					SELECT count_variadic(1, 2, 3) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_expanded_arguments" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_array" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_array"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_array", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_arg AS
| 					SELECT count_variadic(VARIADIC ARRAY[1, 2, 3]::integer[]) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_array" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_one_arg" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_one_arg"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_one_arg", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_one AS
| 					SELECT count_variadic(1) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_one_arg" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_mixed_coercible_args" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_mixed_coercible_args"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_mixed_coercible_args", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums bigint[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_mixed AS
| 					SELECT count_variadic(1::smallint, 2::integer, 3::bigint) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_mixed_coercible_args" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_null_array" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_null_array"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_null_array", public;
|
| 				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
| 				RETURNS integer
| 				LANGUAGE sql
| 				AS 'SELECT cardinality(nums)';
| 				CREATE VIEW v_variadic_null AS
| 					SELECT count_variadic(VARIADIC NULL::integer[]) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_user_variadic_function_explicit_null_array" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_text_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_text_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_text_alias_column", public;
|
| 				CREATE TABLE docs (data jsonb);
| 				CREATE VIEW v_docs AS
| 					SELECT d.data ->> c.col_name AS value
| 					FROM docs d,
| 					     unnest(ARRAY['name', 'status']::text[]) AS c(col_name);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_text_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_varchar_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_varchar_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_varchar_alias_column", public;
|
| 				CREATE TABLE docs (data jsonb);
| 				CREATE VIEW v_docs AS
| 					SELECT d.data ->> c.col_name AS value
| 					FROM docs d,
| 					     unnest(ARRAY['name', 'status']::varchar[]) AS c(col_name);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_unnest_varchar_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_cte_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_cte_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_cte_alias_column", public;
|
| 				CREATE TABLE docs (data jsonb);
| 				CREATE VIEW v_docs AS
| 					WITH c(col_name) AS (SELECT 'name'::text)
| 					SELECT d.data ->> c.col_name AS value
| 					FROM docs d, c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_extract_with_cte_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_integer_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unnest_integer_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unnest_integer_alias_column", public;
|
| 				CREATE VIEW v_nums AS
| 					SELECT n.value + 1 AS next_value
| 					FROM unnest(ARRAY[1, 2, 3]::integer[]) AS n(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_integer_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_bigint_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unnest_bigint_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unnest_bigint_alias_column", public;
|
| 				CREATE VIEW v_nums AS
| 					SELECT n.value + 1 AS next_value
| 					FROM unnest(ARRAY[1, 2, 3]::bigint[]) AS n(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_bigint_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_uuid_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unnest_uuid_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unnest_uuid_alias_column", public;
|
| 				CREATE VIEW v_ids AS
| 					SELECT n.value::text AS id_text
| 					FROM unnest(ARRAY['00000000-0000-0000-0000-000000000001']::uuid[]) AS n(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_uuid_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_jsonb_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unnest_jsonb_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unnest_jsonb_alias_column", public;
|
| 				CREATE VIEW v_payloads AS
| 					SELECT n.value ->> 'name' AS name
| 					FROM unnest(ARRAY['{"name":"one"}'::jsonb]::jsonb[]) AS n(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_jsonb_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_domain_array_alias_column" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unnest_domain_array_alias_column"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unnest_domain_array_alias_column", public;
|
| 				CREATE DOMAIN loader_text_domain AS text;
| 				CREATE VIEW v_domain_values AS
| 					SELECT n.value::text AS value
| 					FROM unnest(ARRAY['one']::loader_text_domain[]) AS n(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unnest_domain_array_alias_column" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_array_polymorphic_builtins" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_array_polymorphic_builtins"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_array_polymorphic_builtins", public;
|
| 				CREATE VIEW v_arrays AS
| 					SELECT
| 						array_length(ARRAY[1, 2, 3]::integer[], 1) AS len,
| 						array_position(ARRAY['a', 'b']::text[], 'b') AS text_pos,
| 						array_position(ARRAY[1, 2, 3]::integer[], 2) AS int_pos,
| 						array_append(ARRAY[1, 2]::integer[], 3) AS appended,
| 						array_prepend(0, ARRAY[1, 2]::integer[]) AS prepended;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_array_polymorphic_builtins" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_anycompatible_common_type_builtins" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_anycompatible_common_type_builtins"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_anycompatible_common_type_builtins", public;
|
| 				CREATE VIEW v_common AS
| 					SELECT
| 						coalesce(NULL, 1, 2::bigint) AS c,
| 						greatest(1, 2::bigint) AS g,
| 						least(1, 2::bigint) AS l;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_anycompatible_common_type_builtins" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_record_returning_builtin_alias_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_record_returning_builtin_alias_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_record_returning_builtin_alias_columns", public;
|
| 				CREATE VIEW v_each AS
| 					SELECT e.key, e.value
| 					FROM jsonb_each_text('{"a":"b"}'::jsonb) AS e(key, value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_record_returning_builtin_alias_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_record_returning_user_function_alias_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_record_returning_user_function_alias_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_record_returning_user_function_alias_columns", public;
|
| 				CREATE FUNCTION loader_pair(OUT id integer, OUT name text)
| 				RETURNS record
| 				LANGUAGE sql
| 				AS 'SELECT 1, ''one''';
| 				CREATE VIEW v_pair AS
| 					SELECT p.id, p.name
| 					FROM loader_pair() AS p(id, name);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_record_returning_user_function_alias_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_returns_table_function_alias_columns" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_returns_table_function_alias_columns"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_returns_table_function_alias_columns", public;
|
| 				CREATE FUNCTION loader_table()
| 				RETURNS TABLE(id integer, name text)
| 				LANGUAGE sql
| 				AS 'SELECT 1, ''one''';
| 				CREATE VIEW v_table AS
| 					SELECT t.id, t.name
| 					FROM loader_table() AS t(id, name);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_returns_table_function_alias_columns" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_srf_alias_column_list" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_srf_alias_column_list"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_srf_alias_column_list", public;
|
| 				CREATE VIEW v_srf AS
| 					SELECT x.value + 1 AS next_value
| 					FROM generate_series(1, 3) AS x(value);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_srf_alias_column_list" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unknown_literal_and_default_resolution" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_unknown_literal_and_default_resolution"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_unknown_literal_and_default_resolution", public;
|
| 				CREATE FUNCTION loader_default_one(a integer, b integer DEFAULT 1)
| 				RETURNS integer LANGUAGE sql AS 'SELECT a + b';
| 				CREATE FUNCTION loader_default_two(a integer, b integer DEFAULT 1, c integer DEFAULT 2)
| 				RETURNS integer LANGUAGE sql AS 'SELECT a + b + c';
| 				CREATE FUNCTION loader_text_arg(a text) RETURNS text LANGUAGE sql AS 'SELECT a';
| 				CREATE FUNCTION loader_int_arg(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
| 				CREATE FUNCTION loader_null_arg(a text, b integer DEFAULT 1) RETURNS text LANGUAGE sql AS 'SELECT a';
| 				CREATE VIEW v_resolve AS
| 					SELECT
| 						loader_default_one(1) AS d1,
| 						loader_default_two(1) AS d2,
| 						loader_text_arg('x') AS text_value,
| 						loader_int_arg(1) AS int_value,
| 						loader_null_arg(NULL) AS null_value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_unknown_literal_and_default_resolution" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_operator_resolution_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_operator_resolution_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_operator_resolution_matrix", public;
|
| 				CREATE TABLE op_items (
| 					j json,
| 					jb jsonb,
| 					t text,
| 					i integer,
| 					b bigint,
| 					n numeric,
| 					d date,
| 					ts timestamp,
| 					ia integer[],
| 					r int4range
| 				);
| 				CREATE VIEW v_ops AS
| 					SELECT
| 						jb -> 'name' AS jb_obj_text,
| 						jb -> 0 AS jb_obj_int,
| 						jb ->> 'name' AS jb_text_text,
| 						jb ->> 0 AS jb_text_int,
| 						j -> 'name' AS j_obj_text,
| 						j ->> 'name' AS j_text_text,
| 						t || 'x' AS text_unknown_right,
| 						'x' || t AS text_unknown_left,
| 						i + b AS int_bigint,
| 						b + i AS bigint_int,
| 						n + i AS numeric_int,
| 						d + 1 AS date_plus_int,
| 						ts + interval '1 day' AS ts_plus_interval,
| 						t LIKE 'a%' AS text_like,
| 						t ILIKE 'a%' AS text_ilike,
| 						ia @> ARRAY[1]::integer[] AS array_contains,
| 						jb @> '{"a":1}'::jsonb AS jsonb_contains,
| 						r && int4range(1, 3) AS range_overlaps
| 					FROM op_items;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_operator_resolution_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_range_resolution" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_range_resolution"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_range_resolution", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_base_items AS
| 					WITH cte1 AS (SELECT id FROM base_items)
| 					SELECT id FROM cte1;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_range_resolution" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_column_alias_resolution" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_column_alias_resolution"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_column_alias_resolution", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_base_items AS
| 					WITH cte1(item_id) AS (SELECT id FROM base_items)
| 					SELECT item_id FROM cte1;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_column_alias_resolution" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_shadows_base_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_shadows_base_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_shadows_base_table", public;
|
| 				CREATE TABLE cte1 (id integer);
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_base_items AS
| 					WITH cte1 AS (SELECT id FROM base_items)
| 					SELECT id FROM cte1;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_shadows_base_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_multiple_ctes_reference_prior_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_multiple_ctes_reference_prior_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_multiple_ctes_reference_prior_cte", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_multi_cte AS
| 					WITH c1 AS (SELECT id FROM base_items),
| 					     c2 AS (SELECT id + 1 AS next_id FROM c1)
| 					SELECT next_id FROM c2;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_multiple_ctes_reference_prior_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_same_column_names_as_base_table" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_same_column_names_as_base_table"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_same_column_names_as_base_table", public;
|
| 				CREATE TABLE base_items (id integer, value text);
| 				CREATE VIEW v_cte_same_cols AS
| 					WITH c AS (SELECT id, value FROM base_items)
| 					SELECT id, value FROM c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_same_column_names_as_base_table" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_subquery_references_outer_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_nested_subquery_references_outer_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_nested_subquery_references_outer_cte", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_nested_outer_cte AS
| 					WITH c AS (SELECT id FROM base_items)
| 					SELECT (SELECT max(id) FROM c) AS max_id;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_subquery_references_outer_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_with_shadows_outer_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_nested_with_shadows_outer_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_nested_with_shadows_outer_cte", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_nested_shadow_cte AS
| 					WITH c AS (SELECT id FROM base_items)
| 					SELECT id
| 					FROM (
| 						WITH c AS (SELECT 2 AS id)
| 						SELECT id FROM c
| 					) s;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_with_shadows_outer_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_with_references_outer_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_nested_with_references_outer_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_nested_with_references_outer_cte", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_nested_outer_ref AS
| 					WITH c AS (SELECT id FROM base_items)
| 					SELECT id
| 					FROM (
| 						WITH d AS (SELECT id + 1 AS id FROM c)
| 						SELECT id FROM d
| 					) s;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_with_references_outer_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_recursive_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_recursive_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_recursive_cte", public;
|
| 				CREATE VIEW v_recursive AS
| 					WITH RECURSIVE r(n) AS (
| 						SELECT 1
| 						UNION ALL
| 						SELECT n + 1 FROM r WHERE n < 3
| 					)
| 					SELECT n FROM r;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_recursive_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_materialized" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_materialized"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_materialized", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_cte_mat AS
| 					WITH c AS MATERIALIZED (SELECT id FROM base_items)
| 					SELECT id FROM c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_materialized" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_not_materialized" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_not_materialized"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_not_materialized", public;
|
| 				CREATE TABLE base_items (id integer);
| 				CREATE VIEW v_cte_not_mat AS
| 					WITH c AS NOT MATERIALIZED (SELECT id FROM base_items)
| 					SELECT id FROM c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_not_materialized" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_column_aliases_fewer_than_output" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cte_column_aliases_fewer_than_output"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cte_column_aliases_fewer_than_output", public;
|
| 				CREATE VIEW v_cte_fewer_aliases AS
| 					WITH c(a) AS (SELECT 1 AS x, 2 AS y)
| 					SELECT a, y FROM c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cte_column_aliases_fewer_than_output" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_later_lateral_references_previous_lateral_alias" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_later_lateral_references_previous_lateral_alias"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_later_lateral_references_previous_lateral_alias", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_lateral_base AS
| 					SELECT s2.next_id
| 					FROM lateral_base b,
| 					     LATERAL (SELECT b.id AS id) s1,
| 				     LATERAL (SELECT s1.id + 1 AS next_id) s2;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_later_lateral_references_previous_lateral_alias" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_cross_lateral AS
| 					SELECT s.id
| 					FROM lateral_base b
| 					CROSS JOIN LATERAL (SELECT b.id AS id) s;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_inner_join_lateral_references_left" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_inner_join_lateral_references_left"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_inner_join_lateral_references_left", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_inner_lateral AS
| 					SELECT s.id
| 					FROM lateral_base b
| 					INNER JOIN LATERAL (SELECT b.id AS id) s ON true;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_inner_join_lateral_references_left" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_references_multiple_prior_items" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_lateral_references_multiple_prior_items"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_lateral_references_multiple_prior_items", public;
|
| 				CREATE TABLE a (id integer);
| 				CREATE TABLE b (id integer);
| 				CREATE VIEW v_multi_lateral AS
| 					SELECT s.sum_id
| 					FROM a, b,
| 					     LATERAL (SELECT a.id + b.id AS sum_id) s;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_references_multiple_prior_items" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_third_lateral_references_two_prior_lateral_aliases" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_third_lateral_references_two_prior_lateral_aliases"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_third_lateral_references_two_prior_lateral_aliases", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_chain_lateral AS
| 					SELECT s3.total
| 					FROM lateral_base b,
| 					     LATERAL (SELECT b.id AS id1) s1,
| 					     LATERAL (SELECT s1.id1 + 1 AS id2) s2,
| 					     LATERAL (SELECT s1.id1 + s2.id2 AS total) s3;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_third_lateral_references_two_prior_lateral_aliases" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_unnest_references_left" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_lateral_unnest_references_left"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_lateral_unnest_references_left", public;
|
| 				CREATE TABLE array_items (id integer, vals integer[]);
| 				CREATE VIEW v_lateral_unnest AS
| 					SELECT u.val + array_items.id AS value
| 					FROM array_items,
| 					     LATERAL unnest(array_items.vals) AS u(val);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_unnest_references_left" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_generate_series_references_left" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_lateral_generate_series_references_left"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_lateral_generate_series_references_left", public;
|
| 				CREATE TABLE series_items (n integer);
| 				CREATE VIEW v_lateral_series AS
| 					SELECT g.x
| 					FROM series_items,
| 					     LATERAL generate_series(1, series_items.n) AS g(x);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_generate_series_references_left" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_rows_from_with_ordinality" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_rows_from_with_ordinality"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_rows_from_with_ordinality", public;
|
| 				CREATE VIEW v_rows_from AS
| 					SELECT x.val, x.ord
| 					FROM ROWS FROM (unnest(ARRAY['a', 'b']::text[])) WITH ORDINALITY AS x(val, ord);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_rows_from_with_ordinality" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_record_coldeflist" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_lateral_record_coldeflist"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_lateral_record_coldeflist", public;
|
| 				CREATE TABLE json_items (payload jsonb);
| 				CREATE VIEW v_lateral_record AS
| 					SELECT r.id, r.name
| 					FROM json_items,
| 					     LATERAL jsonb_to_record(json_items.payload) AS r(id integer, name text);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_record_coldeflist" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_join_scope_nested" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_lateral_join_scope_nested"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_lateral_join_scope_nested", public;
|
| 				CREATE TABLE a (id integer);
| 				CREATE TABLE b (id integer);
| 				CREATE VIEW v_lateral_join_scope AS
| 					SELECT s.sum_id
| 					FROM a
| 					JOIN b ON true
| 					JOIN LATERAL (SELECT a.id + b.id AS sum_id) s ON true;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_lateral_join_scope_nested" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_correlated_subquery_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_correlated_subquery_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_correlated_subquery_matrix", public;
|
| 				CREATE TABLE parents (id integer);
| 				CREATE TABLE children (parent_id integer, value integer);
| 				CREATE VIEW v_correlated AS
| 					SELECT
| 						p.id,
| 						(SELECT max(c.value) FROM children c WHERE c.parent_id = p.id) AS max_value
| 					FROM parents p
| 					WHERE EXISTS (SELECT 1 FROM children c WHERE c.parent_id = p.id)
| 					  AND p.id IN (SELECT c.parent_id FROM children c WHERE c.value > p.id)
| 					  AND (SELECT count(*) FROM children c WHERE c.parent_id = p.id) > 0;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_correlated_subquery_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_correlated_subquery_levels_up_two" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_nested_correlated_subquery_levels_up_two"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_nested_correlated_subquery_levels_up_two", public;
|
| 				CREATE TABLE parents (id integer);
| 				CREATE TABLE children (parent_id integer, value integer);
| 				CREATE VIEW v_nested_correlated AS
| 					SELECT p.id
| 					FROM parents p
| 					WHERE EXISTS (
| 						SELECT 1
| 						FROM children c
| 						WHERE EXISTS (
| 							SELECT 1
| 							WHERE c.parent_id = p.id
| 						)
| 					);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_nested_correlated_subquery_levels_up_two" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_from_unnest_alias" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_star_from_unnest_alias"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_star_from_unnest_alias", public;
|
| 				CREATE VIEW v_star_unnest AS
| 					SELECT *
| 					FROM unnest(ARRAY['a', 'b']::text[]) AS u(val);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_from_unnest_alias" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_from_multi_unnest_alias" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_star_from_multi_unnest_alias"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_star_from_multi_unnest_alias", public;
|
| 				CREATE VIEW v_star_multi_unnest AS
| 					SELECT *
| 					FROM unnest(ARRAY[1, 2]::integer[], ARRAY['a', 'b']::text[]) AS u(id, name);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_from_multi_unnest_alias" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_record_coldeflist" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_record_coldeflist"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_record_coldeflist", public;
|
| 				CREATE VIEW v_jsonb_record AS
| 					SELECT r.id, r.name
| 					FROM jsonb_to_record('{"id":1,"name":"one"}'::jsonb) AS r(id integer, name text);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_record_coldeflist" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_recordset_coldeflist" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_recordset_coldeflist"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_recordset_coldeflist", public;
|
| 				CREATE VIEW v_jsonb_recordset AS
| 					SELECT r.id, r.name
| 					FROM jsonb_to_recordset('[{"id":1,"name":"one"}]'::jsonb) AS r(id integer, name text);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_jsonb_to_recordset_coldeflist" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_srf_in_select_target" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_srf_in_select_target"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_srf_in_select_target", public;
|
| 				CREATE VIEW v_srf_target AS
| 					SELECT generate_series(1, 3) AS value;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_srf_in_select_target" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_cte" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_cte"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_cte", public;
|
| 				CREATE TABLE star_items (id integer, name text);
| 				CREATE VIEW v_star_cte AS
| 					WITH c AS (SELECT id, name FROM star_items)
| 					SELECT * FROM c;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_cte" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_lateral" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_lateral"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_lateral", public;
|
| 				CREATE TABLE star_items (id integer, vals integer[]);
| 				CREATE VIEW v_star_lateral AS
| 					SELECT *
| 					FROM star_items,
| 					     LATERAL unnest(star_items.vals) AS u(val);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_star_expansion_through_lateral" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_column_alias_list_overrides_target_names" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_column_alias_list_overrides_target_names"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_column_alias_list_overrides_target_names", public;
| CREATE VIEW v_alias_override(a, b) AS SELECT 1 AS x, 2 AS y;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_column_alias_list_overrides_target_names" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_column_aliases_fewer_than_output" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_column_aliases_fewer_than_output"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_column_aliases_fewer_than_output", public;
| CREATE VIEW v_alias_fewer(a) AS SELECT 1 AS x, 2 AS y;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_column_aliases_fewer_than_output" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_alter_rename_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_object_namespace_alter_rename_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_object_namespace_alter_rename_matrix", public;
|
| 				CREATE TABLE t (id integer);
| 				CREATE INDEX t_id_idx ON t (id);
| 				ALTER INDEX t_id_idx RENAME TO t_id_idx2;
| 				CREATE SEQUENCE s;
| 				ALTER SEQUENCE s RENAME TO s2;
| 				ALTER SEQUENCE s2 OWNER TO CURRENT_USER;
| 				CREATE VIEW v AS SELECT id FROM t;
| 				ALTER VIEW v RENAME TO v2;
| 				CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
| 				ALTER MATERIALIZED VIEW mv RENAME TO mv2;
| 				CREATE TYPE mood AS ENUM ('sad', 'ok');
| 				ALTER TYPE mood RENAME TO mood2;
| 				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
| 				ALTER DOMAIN positive_int RENAME TO positive_int2;
| 				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
| 				ALTER FUNCTION f(integer) RENAME TO f2;
| 				CREATE PROCEDURE p(a integer) LANGUAGE sql AS 'SELECT 1';
| 				ALTER PROCEDURE p(integer) RENAME TO p2;
| 				ALTER ROUTINE f2(integer) OWNER TO CURRENT_USER;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_alter_rename_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_comment_drop_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_object_namespace_comment_drop_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_object_namespace_comment_drop_matrix", public;
|
| 				CREATE TABLE t (id integer);
| 				CREATE INDEX t_id_idx ON t (id);
| 				COMMENT ON INDEX t_id_idx IS 'idx';
| 				COMMENT ON INDEX t_id_idx IS NULL;
| 				CREATE SEQUENCE s;
| 				COMMENT ON SEQUENCE s IS 'seq';
| 				COMMENT ON SEQUENCE s IS NULL;
| 				CREATE TYPE mood AS ENUM ('sad', 'ok');
| 				COMMENT ON TYPE mood IS 'type';
| 				COMMENT ON TYPE mood IS NULL;
| 				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
| 				COMMENT ON FUNCTION f(a integer) IS 'func';
| 				COMMENT ON FUNCTION f(integer) IS NULL;
| 				DROP FUNCTION f(integer);
| 				DROP INDEX t_id_idx;
| 				DROP SEQUENCE s;
| 				DROP TYPE mood;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_comment_drop_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_grant_revoke_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_object_namespace_grant_revoke_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_object_namespace_grant_revoke_matrix", public;
|
| 				CREATE TABLE t (id integer);
| 				CREATE VIEW v AS SELECT id FROM t;
| 				CREATE SEQUENCE s;
| 				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
| 				GRANT SELECT ON TABLE t TO PUBLIC;
| 				REVOKE SELECT ON TABLE t FROM PUBLIC;
| 				GRANT SELECT ON TABLE v TO PUBLIC;
| 				REVOKE SELECT ON TABLE v FROM PUBLIC;
| 				GRANT USAGE ON SEQUENCE s TO PUBLIC;
| 				REVOKE USAGE ON SEQUENCE s FROM PUBLIC;
| 				GRANT EXECUTE ON FUNCTION f(integer) TO PUBLIC;
| 				REVOKE EXECUTE ON FUNCTION f(integer) FROM PUBLIC;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_grant_revoke_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_set_schema_matrix" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_object_namespace_set_schema_matrix"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_object_namespace_set_schema_matrix", public;
|
| 				CREATE SCHEMA target_schema;
| 				CREATE TABLE t (id integer);
| 				ALTER TABLE t SET SCHEMA target_schema;
| 				CREATE SEQUENCE s;
| 				ALTER SEQUENCE s SET SCHEMA target_schema;
| 				CREATE TYPE mood AS ENUM ('sad', 'ok');
| 				ALTER TYPE mood SET SCHEMA target_schema;
| 				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
| 				ALTER FUNCTION f(integer) SET SCHEMA target_schema;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_object_namespace_set_schema_matrix" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_depends_on_extension" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_index_depends_on_extension"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_index_depends_on_extension", public;
|
| 				CREATE TABLE t (id integer);
| 				CREATE INDEX t_id_idx ON t (id);
| 				ALTER INDEX t_id_idx DEPENDS ON EXTENSION plpgsql;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_depends_on_extension" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_owned_by_loaded_objects" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_owned_by_loaded_objects"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_owned_by_loaded_objects", public;
|
| 				CREATE ROLE loader_owner;
| 				CREATE TABLE t (id integer);
| 				CREATE SEQUENCE s;
| 				ALTER TABLE t OWNER TO loader_owner;
| 				ALTER SEQUENCE s OWNER TO loader_owner;
| 				DROP OWNED BY loader_owner;
| 				DROP ROLE loader_owner;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_owned_by_loaded_objects" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_left_join_lateral_references_left_relation" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_left_join_lateral_references_left_relation"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_left_join_lateral_references_left_relation", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_lateral_left AS
| 					SELECT b.id, s.next_id
| 					FROM lateral_base b
| 					LEFT JOIN LATERAL (SELECT b.id + 1 AS next_id) s ON true;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_left_join_lateral_references_left_relation" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_chained_lateral_references_base_and_prior_alias" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_chained_lateral_references_base_and_prior_alias"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_chained_lateral_references_base_and_prior_alias", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_lateral_base AS
| 					SELECT s3.value
| 					FROM lateral_base b,
| 					     LATERAL (SELECT b.id AS id) s1,
| 					     LATERAL (SELECT s1.id + b.id AS id2) s2,
| 					     LATERAL (SELECT s2.id2 + s1.id AS value) s3;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_chained_lateral_references_base_and_prior_alias" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left_relation" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left_relation"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left_relation", public;
|
| 				CREATE TABLE lateral_base (id integer);
| 				CREATE VIEW v_lateral_cross AS
| 					SELECT b.id, s.next_id
| 					FROM lateral_base b
| 					CROSS JOIN LATERAL (SELECT b.id + 1 AS next_id) s;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_view_cross_join_lateral_references_left_relation" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition", public;
|
| 				CREATE TABLE parent_idx_attach (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE child_idx_attach PARTITION OF parent_idx_attach
| 					FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_idx_attach_id_idx ON ONLY parent_idx_attach (id);
| 				CREATE INDEX child_idx_attach_id_idx ON child_idx_attach (id);
| 				ALTER INDEX parent_idx_attach_id_idx ATTACH PARTITION child_idx_attach_id_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition_schema_qualified" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition_schema_qualified"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition_schema_qualified", public;
|
| 				CREATE SCHEMA part_s;
| 				CREATE TABLE part_s.parent_idx_attach (id integer) PARTITION BY RANGE (id);
| 				CREATE TABLE part_s.child_idx_attach PARTITION OF part_s.parent_idx_attach
| 					FOR VALUES FROM (0) TO (10);
| 				CREATE INDEX parent_idx_attach_id_idx ON ONLY part_s.parent_idx_attach (id);
| 				CREATE INDEX child_idx_attach_id_idx ON part_s.child_idx_attach (id);
| 				ALTER INDEX part_s.parent_idx_attach_id_idx ATTACH PARTITION part_s.child_idx_attach_id_idx;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_alter_index_attach_partition_schema_qualified" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_key_referenced_by_fk_with_cascade" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_key_referenced_by_fk_with_cascade"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_key_referenced_by_fk_with_cascade", public;
|
| 				CREATE TABLE drop_pk_cascade_parent (id integer CONSTRAINT drop_pk_cascade_parent_key PRIMARY KEY);
| 				CREATE TABLE drop_pk_cascade_child (parent_id integer CONSTRAINT drop_pk_cascade_child_fk REFERENCES drop_pk_cascade_parent (id));
| 				ALTER TABLE drop_pk_cascade_parent DROP CONSTRAINT drop_pk_cascade_parent_key CASCADE;
| 				ALTER TABLE drop_pk_cascade_child ADD CONSTRAINT drop_pk_cascade_child_fk CHECK (parent_id > 0);
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_key_referenced_by_fk_with_cascade" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_key_after_dropping_its_fk" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "testloadercompatpg17oracleacceptscorpus_drop_key_after_dropping_its_fk"
< ok 0 -
> exec
| SET search_path TO "testloadercompatpg17oracleacceptscorpus_drop_key_after_dropping_its_fk", public;
|
| 				CREATE TABLE drop_pk_free_parent (id integer CONSTRAINT drop_pk_free_parent_key PRIMARY KEY);
| 				CREATE TABLE drop_pk_free_child (parent_id integer CONSTRAINT drop_pk_free_child_fk REFERENCES drop_pk_free_parent (id));
| 				ALTER TABLE drop_pk_free_child DROP CONSTRAINT drop_pk_free_child_fk;
| 				ALTER TABLE drop_pk_free_parent DROP CONSTRAINT drop_pk_free_parent_key;
| 			
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "testloadercompatpg17oracleacceptscorpus_drop_key_after_dropping_its_fk" CASCADE
< ok 0 -
