package catalog

import "testing"

type loaderCompatCase struct {
	name string
	sql  string
}

func loaderCompatRejectCases() []loaderCompatCase {
	return []loaderCompatCase{
		{
			name: "numeric_fk_references_integer_pk",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id numeric REFERENCES parent(id)
				);
			`,
		},
		{
			name: "composite_fk_rejects_one_incompatible_column_pair",
			sql: `
				CREATE TABLE parent (id integer, tenant_id integer, PRIMARY KEY (id, tenant_id));
				CREATE TABLE child (
					parent_id numeric,
					tenant_id integer,
					FOREIGN KEY (parent_id, tenant_id) REFERENCES parent(id, tenant_id)
				);
			`,
		},
		{
			name: "alter_index_attach_rejects_wrong_child_table",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
				CREATE TABLE other (id integer);
				CREATE INDEX parent_id_idx ON ONLY parent (id);
				CREATE INDEX other_id_idx ON other (id);
				ALTER INDEX parent_id_idx ATTACH PARTITION other_id_idx;
			`,
		},
		{
			name: "alter_index_attach_rejects_nonpartitioned_parent_table",
			sql: `
				CREATE TABLE parent (id integer);
				CREATE TABLE child (id integer);
				CREATE INDEX parent_id_idx ON parent (id);
				CREATE INDEX child_id_idx ON child (id);
				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
			`,
		},
		{
			name: "partitioned_zero_column_table_expression_key_rejected",
			sql:  `CREATE TABLE zc () PARTITION BY LIST (());`,
		},
		{
			name: "explicit_variadic_rejects_non_array_argument",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_bad_variadic AS
					SELECT count_variadic(VARIADIC 1) AS value;
			`,
		},
		{
			name: "user_variadic_function_zero_args_rejected",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_zero AS
					SELECT count_variadic() AS value;
			`,
		},
		{
			name: "recursive_cte_self_reference_without_union_rejected",
			sql: `
				CREATE VIEW v_bad_recursive AS
					WITH RECURSIVE r(n) AS (SELECT n + 1 FROM r)
					SELECT n FROM r;
			`,
		},
		{
			name: "cte_column_aliases_more_than_output_rejected",
			sql: `
				CREATE VIEW v_bad_cte_alias AS
					WITH c(a, b) AS (SELECT 1)
					SELECT a FROM c;
			`,
		},
		{
			name: "non_lateral_subquery_references_prior_from_item_rejected",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_bad_non_lateral AS
					SELECT s.id
					FROM base_items b,
					     (SELECT b.id AS id) s;
			`,
		},
		{
			name: "check_constraint_subquery_rejected",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE TABLE t (
					id integer,
					CHECK (id IN (SELECT id FROM base_items))
				);
			`,
		},
		{
			name: "generated_column_subquery_rejected",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE TABLE t (
					id integer,
					generated integer GENERATED ALWAYS AS ((SELECT max(id) FROM base_items)) STORED
				);
			`,
		},
		{
			name: "view_column_aliases_more_than_output_rejected",
			sql:  `CREATE VIEW v_bad_alias(a, b) AS SELECT 1;`,
		},
		{
			name: "view_duplicate_output_column_names_rejected",
			sql:  `CREATE VIEW v_dupe_cols AS SELECT 1 AS id, 2 AS id;`,
		},
	}
}

func loaderCompatAcceptCases() []loaderCompatCase {
	return []loaderCompatCase{
		{
			name: "zero_column_table",
			sql:  `CREATE TABLE zc ();`,
		},
		{
			name: "zero_column_table_in_schema",
			sql: `
				CREATE SCHEMA loader_s1;
				CREATE TABLE loader_s1.zc ();
			`,
		},
		{
			name: "zero_column_table_comment",
			sql: `
				CREATE TABLE zc ();
				COMMENT ON TABLE zc IS 'zero columns';
			`,
		},
		{
			name: "zero_column_table_grant",
			sql: `
				CREATE TABLE zc ();
				GRANT SELECT ON zc TO PUBLIC;
			`,
		},
		{
			name: "zero_column_table_view_star",
			sql: `
				CREATE TABLE zc ();
				CREATE VIEW v_zc AS SELECT * FROM zc;
			`,
		},
		{
			name: "foreign_key_match_simple",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer,
					FOREIGN KEY (parent_id) REFERENCES parent(id) MATCH SIMPLE
				);
			`,
		},
		{
			name: "foreign_key_match_simple_multicolumn",
			sql: `
				CREATE TABLE parent (a integer, b integer, PRIMARY KEY (a, b));
				CREATE TABLE child (
					a integer,
					b integer,
					FOREIGN KEY (a, b) REFERENCES parent(a, b) MATCH SIMPLE
				);
			`,
		},
		{
			name: "foreign_key_match_simple_on_update_cascade",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id) MATCH SIMPLE ON UPDATE CASCADE
				);
			`,
		},
		{
			name: "foreign_key_match_simple_on_delete_set_null",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id) MATCH SIMPLE ON DELETE SET NULL
				);
			`,
		},
		{
			name: "foreign_key_match_simple_deferrable",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id) MATCH SIMPLE DEFERRABLE INITIALLY DEFERRED
				);
			`,
		},
		{
			name: "comment_function_argument_names",
			sql: `
				CREATE FUNCTION f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				COMMENT ON FUNCTION f(a integer, b text) IS 'comment';
			`,
		},
		{
			name: "comment_schema_qualified_function_argument_names",
			sql: `
				CREATE SCHEMA funcs;
				CREATE FUNCTION funcs.f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				COMMENT ON FUNCTION funcs.f(a integer, b text) IS 'comment';
			`,
		},
		{
			name: "comment_function_quoted_argument_names",
			sql: `
				CREATE FUNCTION f("select" integer, "from" text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				COMMENT ON FUNCTION f("select" integer, "from" text) IS 'comment';
			`,
		},
		{
			name: "comment_function_qualified_argument_types",
			sql: `
				CREATE FUNCTION f(a pg_catalog.int4, b pg_catalog.text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				COMMENT ON FUNCTION f(a pg_catalog.int4, b pg_catalog.text) IS 'comment';
			`,
		},
		{
			name: "comment_variadic_function_identity",
			sql: `
				CREATE FUNCTION f(VARIADIC nums integer[]) RETURNS integer
					LANGUAGE sql AS 'SELECT cardinality(nums)';
				COMMENT ON FUNCTION f(VARIADIC nums integer[]) IS 'comment';
			`,
		},
		{
			name: "grant_function_argument_names",
			sql: `
				CREATE FUNCTION f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				GRANT EXECUTE ON FUNCTION f(a integer, b text) TO PUBLIC;
			`,
		},
		{
			name: "revoke_function_argument_names",
			sql: `
				CREATE FUNCTION f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				GRANT EXECUTE ON FUNCTION f(a integer, b text) TO PUBLIC;
				REVOKE EXECUTE ON FUNCTION f(a integer, b text) FROM PUBLIC;
			`,
		},
		{
			name: "alter_function_argument_names",
			sql: `
				CREATE FUNCTION f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				ALTER FUNCTION f(a integer, b text) RENAME TO f2;
			`,
		},
		{
			name: "drop_function_argument_names",
			sql: `
				CREATE FUNCTION f(a integer, b text) RETURNS integer
					LANGUAGE sql AS 'SELECT 1';
				DROP FUNCTION f(a integer, b text);
			`,
		},
		{
			name: "comment_function_identity_in_mode",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				COMMENT ON FUNCTION f(IN a integer) IS 'comment';
			`,
		},
		{
			name: "comment_function_identity_inout_mode",
			sql: `
				CREATE FUNCTION f(INOUT a integer)
					LANGUAGE sql AS 'SELECT $1';
				COMMENT ON FUNCTION f(INOUT a integer) IS 'comment';
			`,
		},
		{
			name: "grant_function_identity_in_mode",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				GRANT EXECUTE ON FUNCTION f(IN a integer) TO PUBLIC;
			`,
		},
		{
			name: "grant_function_identity_inout_mode",
			sql: `
				CREATE FUNCTION f(INOUT a integer)
					LANGUAGE sql AS 'SELECT $1';
				GRANT EXECUTE ON FUNCTION f(INOUT a integer) TO PUBLIC;
			`,
		},
		{
			name: "grant_function_identity_variadic_mode",
			sql: `
				CREATE FUNCTION f(VARIADIC a integer[]) RETURNS integer
					LANGUAGE sql AS 'SELECT cardinality($1)';
				GRANT EXECUTE ON FUNCTION f(VARIADIC a integer[]) TO PUBLIC;
			`,
		},
		{
			name: "revoke_function_identity_in_mode",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				GRANT EXECUTE ON FUNCTION f(IN a integer) TO PUBLIC;
				REVOKE EXECUTE ON FUNCTION f(IN a integer) FROM PUBLIC;
			`,
		},
		{
			name: "alter_function_identity_in_mode_owner",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				ALTER FUNCTION f(IN a integer) OWNER TO CURRENT_USER;
			`,
		},
		{
			name: "alter_function_identity_in_mode_options",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				ALTER FUNCTION f(IN a integer) IMMUTABLE;
			`,
		},
		{
			name: "drop_function_identity_in_mode",
			sql: `
				CREATE FUNCTION f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				DROP FUNCTION f(IN a integer);
			`,
		},
		{
			name: "drop_procedure_identity_inout_mode",
			sql: `
				CREATE PROCEDURE p(INOUT a integer)
					LANGUAGE sql AS 'SELECT $1';
				DROP PROCEDURE p(INOUT a integer);
			`,
		},
		{
			name: "drop_routine_identity_argument_name",
			sql: `
				CREATE FUNCTION r(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				DROP ROUTINE r(a integer);
			`,
		},
		{
			name: "function_identity_schema_qualified_types",
			sql: `
				CREATE FUNCTION f(a pg_catalog.int4) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				COMMENT ON FUNCTION f(a pg_catalog.int4) IS 'comment';
			`,
		},
		{
			name: "function_identity_array_type",
			sql: `
				CREATE FUNCTION f(a integer[]) RETURNS integer
					LANGUAGE sql AS 'SELECT cardinality($1)';
				COMMENT ON FUNCTION f(a integer[]) IS 'comment';
			`,
		},
		{
			name: "function_identity_quoted_function_name",
			sql: `
				CREATE FUNCTION "select"(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				COMMENT ON FUNCTION "select"(a integer) IS 'comment';
			`,
		},
		{
			name: "function_identity_schema_qualified_function_name",
			sql: `
				CREATE SCHEMA ident_s;
				CREATE FUNCTION ident_s.f(a integer) RETURNS integer
					LANGUAGE sql AS 'SELECT $1';
				COMMENT ON FUNCTION ident_s.f(a integer) IS 'comment';
			`,
		},
		{
			name: "returns_table_plpgsql",
			sql: `
				CREATE FUNCTION f()
				RETURNS TABLE(id integer, name text)
				LANGUAGE plpgsql
				AS $$
				BEGIN
					RETURN QUERY SELECT 1, 'one'::text;
				END
				$$;
			`,
		},
		{
			name: "returns_table_single_column_plpgsql",
			sql: `
				CREATE FUNCTION f()
				RETURNS TABLE(id integer)
				LANGUAGE plpgsql
				AS $$
				BEGIN
					RETURN QUERY SELECT 1;
				END
				$$;
			`,
		},
		{
			name: "returns_table_qualified_column_types",
			sql: `
				CREATE FUNCTION f()
				RETURNS TABLE(id pg_catalog.int4, name pg_catalog.text)
				LANGUAGE plpgsql
				AS $$
				BEGIN
					RETURN QUERY SELECT 1, 'one'::text;
				END
				$$;
			`,
		},
		{
			name: "bigint_fk_references_integer_pk",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id bigint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "smallint_fk_references_integer_pk",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id smallint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "integer_fk_references_bigint_pk",
			sql: `
				CREATE TABLE parent (id bigint PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id)
				);
			`,
		},
		{
			name: "integer_fk_references_smallint_pk",
			sql: `
				CREATE TABLE parent (id smallint PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id)
				);
			`,
		},
		{
			name: "smallint_fk_references_bigint_pk",
			sql: `
				CREATE TABLE parent (id bigint PRIMARY KEY);
				CREATE TABLE child (
					parent_id smallint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "bigint_fk_references_smallint_pk",
			sql: `
				CREATE TABLE parent (id smallint PRIMARY KEY);
				CREATE TABLE child (
					parent_id bigint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "integer_fk_references_numeric_pk",
			sql: `
				CREATE TABLE parent (id numeric PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id)
				);
			`,
		},
		{
			name: "varchar_fk_references_text_pk",
			sql: `
				CREATE TABLE parent (id text PRIMARY KEY);
				CREATE TABLE child (
					parent_id varchar REFERENCES parent(id)
				);
			`,
		},
		{
			name: "text_fk_references_varchar_pk",
			sql: `
				CREATE TABLE parent (id varchar PRIMARY KEY);
				CREATE TABLE child (
					parent_id text REFERENCES parent(id)
				);
			`,
		},
		{
			name: "bpchar_fk_references_text_pk",
			sql: `
				CREATE TABLE parent (id text PRIMARY KEY);
				CREATE TABLE child (
					parent_id character REFERENCES parent(id)
				);
			`,
		},
		{
			name: "text_fk_references_bpchar_pk",
			sql: `
				CREATE TABLE parent (id character PRIMARY KEY);
				CREATE TABLE child (
					parent_id text REFERENCES parent(id)
				);
			`,
		},
		{
			name: "domain_fk_references_base_pk",
			sql: `
				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id positive_int REFERENCES parent(id)
				);
			`,
		},
		{
			name: "base_fk_references_domain_pk",
			sql: `
				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
				CREATE TABLE parent (id positive_int PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id)
				);
			`,
		},
		{
			name: "same_domain_fk",
			sql: `
				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
				CREATE TABLE parent (id positive_int PRIMARY KEY);
				CREATE TABLE child (
					parent_id positive_int REFERENCES parent(id)
				);
			`,
		},
		{
			name: "different_domain_same_base_fk",
			sql: `
				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
				CREATE DOMAIN other_int AS integer CHECK (VALUE > 0);
				CREATE TABLE parent (id positive_int PRIMARY KEY);
				CREATE TABLE child (
					parent_id other_int REFERENCES parent(id)
				);
			`,
		},
		{
			name: "bigint_fk_references_integer_pk_composite",
			sql: `
				CREATE TABLE parent (id integer, tenant_id integer, PRIMARY KEY (id, tenant_id));
				CREATE TABLE child (
					parent_id bigint,
					tenant_id bigint,
					FOREIGN KEY (parent_id, tenant_id) REFERENCES parent(id, tenant_id)
				);
			`,
		},
		{
			name: "fk_references_unique_index_cross_type",
			sql: `
				CREATE TABLE parent (id integer UNIQUE);
				CREATE TABLE child (
					parent_id bigint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "fk_references_unique_index_with_include_columns",
			sql: `
				CREATE TABLE parent (id integer, extra text);
				CREATE UNIQUE INDEX parent_id_idx ON parent (id) INCLUDE (extra);
				CREATE TABLE child (
					parent_id bigint REFERENCES parent(id)
				);
			`,
		},
		{
			name: "temporary_zero_column_table",
			sql:  `CREATE TEMPORARY TABLE zc ();`,
		},
		{
			name: "unlogged_zero_column_table",
			sql:  `CREATE UNLOGGED TABLE zc ();`,
		},
		{
			name: "create_table_like_zero_column_table",
			sql: `
				CREATE TABLE zc ();
				CREATE TABLE zc_like (LIKE zc);
			`,
		},
		{
			name: "create_table_as_zero_output_columns",
			sql:  `CREATE TABLE zc AS SELECT;`,
		},
		{
			name: "create_view_zero_output_columns",
			sql:  `CREATE VIEW v_zc AS SELECT;`,
		},
		{
			name: "inherited_zero_column_parent",
			sql: `
				CREATE TABLE parent ();
				CREATE TABLE child () INHERITS (parent);
			`,
		},
		{
			name: "alter_zero_column_table_add_column",
			sql: `
				CREATE TABLE zc ();
				ALTER TABLE zc ADD COLUMN id integer;
			`,
		},
		{
			name: "alter_table_drop_column_to_zero_columns",
			sql: `
				CREATE TABLE zc (id integer);
				ALTER TABLE zc DROP COLUMN id;
			`,
		},
		{
			name: "partitioned_index_attach_wrong_schema_parent_qualified",
			sql: `
				CREATE SCHEMA pidx;
				CREATE TABLE pidx.parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE pidx.child PARTITION OF pidx.parent FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_id_idx ON ONLY pidx.parent (id);
				CREATE INDEX child_id_idx ON pidx.child (id);
				ALTER INDEX pidx.parent_id_idx ATTACH PARTITION pidx.child_id_idx;
			`,
		},
		{
			name: "partitioned_index_attach_parent_index_first",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE INDEX parent_id_idx ON ONLY parent (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
			`,
		},
		{
			name: "partitioned_index_attach_child_table_first",
			sql: `
				CREATE TABLE child (id integer);
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				ALTER TABLE parent ATTACH PARTITION child FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_id_idx ON ONLY parent (id);
				CREATE INDEX child_id_idx ON child (id);
				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
			`,
		},
		{
			name: "partition_table_attach_records_dependency",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child (id integer);
				ALTER TABLE parent ATTACH PARTITION child FOR VALUES FROM (0) TO (10);
			`,
		},
		{
			name: "partitioned_unique_index_attach",
			sql: `
				CREATE TABLE parent (id integer NOT NULL) PARTITION BY RANGE (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
				CREATE UNIQUE INDEX parent_id_idx ON ONLY parent (id);
				CREATE UNIQUE INDEX child_id_idx ON child (id);
				ALTER INDEX parent_id_idx ATTACH PARTITION child_id_idx;
			`,
		},
		{
			name: "partitioned_expression_index_attach",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_expr_idx ON ONLY parent ((id + 1));
				CREATE INDEX child_expr_idx ON child ((id + 1));
				ALTER INDEX parent_expr_idx ATTACH PARTITION child_expr_idx;
			`,
		},
		{
			name: "partitioned_partial_index_attach",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_partial_idx ON ONLY parent (id) WHERE id > 0;
				CREATE INDEX child_partial_idx ON child (id) WHERE id > 0;
				ALTER INDEX parent_partial_idx ATTACH PARTITION child_partial_idx;
			`,
		},
		{
			name: "partition_default_bound",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child_default PARTITION OF parent DEFAULT;
			`,
		},
		{
			name: "partition_table_detach",
			sql: `
				CREATE TABLE parent (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (0) TO (10);
				ALTER TABLE parent DETACH PARTITION child;
			`,
		},
		{
			name: "generated_column_expression",
			sql: `
				CREATE TABLE t (
					a integer,
					b integer GENERATED ALWAYS AS (a + 1) STORED
				);
			`,
		},
		{
			name: "identity_column_dependency",
			sql: `
				CREATE TABLE t (
					id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY
				);
			`,
		},
		{
			name: "column_default_function_dependency",
			sql: `
				CREATE FUNCTION next_code() RETURNS integer LANGUAGE sql AS 'SELECT 1';
				CREATE TABLE t (
					code integer DEFAULT next_code()
				);
			`,
		},
		{
			name: "column_default_sequence_dependency",
			sql: `
				CREATE SEQUENCE s;
				CREATE TABLE t (
					id integer DEFAULT nextval('s')
				);
			`,
		},
		{
			name: "check_constraint_function_dependency",
			sql: `
				CREATE FUNCTION is_positive(x integer) RETURNS boolean LANGUAGE sql AS 'SELECT $1 > 0';
				CREATE TABLE t (
					x integer CHECK (is_positive(x))
				);
			`,
		},
		{
			name: "exclusion_constraint_resolves_operator_class",
			sql: `
				CREATE TABLE t (r int4range);
				ALTER TABLE t ADD CONSTRAINT t_r_excl EXCLUDE USING gist (r WITH &&);
			`,
		},
		{
			name: "unique_nulls_not_distinct",
			sql: `
				CREATE TABLE t (x integer);
				ALTER TABLE t ADD CONSTRAINT t_x_key UNIQUE NULLS NOT DISTINCT (x);
			`,
		},
		{
			name: "deferrable_unique_constraint",
			sql: `
				CREATE TABLE t (
					x integer UNIQUE DEFERRABLE INITIALLY DEFERRED
				);
			`,
		},
		{
			name: "deferrable_foreign_key",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (
					parent_id integer REFERENCES parent(id) DEFERRABLE INITIALLY DEFERRED
				);
			`,
		},
		{
			name: "not_valid_foreign_key",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (parent_id integer);
				ALTER TABLE child ADD CONSTRAINT child_parent_fk
					FOREIGN KEY (parent_id) REFERENCES parent(id) NOT VALID;
			`,
		},
		{
			name: "validate_constraint",
			sql: `
				CREATE TABLE parent (id integer PRIMARY KEY);
				CREATE TABLE child (parent_id integer);
				ALTER TABLE child ADD CONSTRAINT child_parent_fk
					FOREIGN KEY (parent_id) REFERENCES parent(id) NOT VALID;
				ALTER TABLE child VALIDATE CONSTRAINT child_parent_fk;
			`,
		},
		{
			name: "constraint_comment",
			sql: `
				CREATE TABLE t (x integer CONSTRAINT x_positive CHECK (x > 0));
				COMMENT ON CONSTRAINT x_positive ON t IS 'positive';
			`,
		},
		{
			name: "view_concat_ws_variadic_builtin",
			sql: `
				CREATE TABLE items (a text, b integer);
				CREATE VIEW v_items AS
					SELECT concat_ws('-', a, b) AS label
					FROM items;
			`,
		},
		{
			name: "view_concat_ws_unknown_literals",
			sql: `
				CREATE VIEW v_items AS
					SELECT concat_ws('-', 'a', 'b', 'c') AS label;
			`,
		},
		{
			name: "view_concat_ws_mixed_types",
			sql: `
				CREATE TABLE items (a text, b integer, c uuid, d timestamp);
				CREATE VIEW v_items AS
					SELECT concat_ws('-', a, b, c, d) AS label
					FROM items;
			`,
		},
		{
			name: "view_jsonb_build_object_variadic_builtin",
			sql: `
				CREATE TABLE items (id integer, name text);
				CREATE VIEW v_items AS
					SELECT jsonb_build_object('id', id, 'name', name) AS payload
					FROM items;
			`,
		},
		{
			name: "view_concat_variadic_builtin",
			sql: `
				CREATE TABLE items (a text, b integer, c timestamp);
				CREATE VIEW v_items AS
					SELECT concat(a, b, c) AS label
					FROM items;
			`,
		},
		{
			name: "view_format_variadic_builtin",
			sql: `
				CREATE TABLE items (a text, b integer);
				CREATE VIEW v_items AS
					SELECT format('%s %s', a, b) AS label
					FROM items;
			`,
		},
		{
			name: "view_json_build_object_variadic_builtin",
			sql: `
				CREATE TABLE items (id integer, name text);
				CREATE VIEW v_items AS
					SELECT json_build_object('id', id, 'name', name) AS payload
					FROM items;
			`,
		},
		{
			name: "view_json_build_array_variadic_builtin",
			sql: `
				CREATE TABLE items (id integer, name text);
				CREATE VIEW v_items AS
					SELECT json_build_array(id, name) AS payload
					FROM items;
			`,
		},
		{
			name: "view_jsonb_build_array_variadic_builtin",
			sql: `
				CREATE TABLE items (id integer, name text);
				CREATE VIEW v_items AS
					SELECT jsonb_build_array(id, name) AS payload
					FROM items;
			`,
		},
		{
			name: "view_user_function_default_argument",
			sql: `
				CREATE FUNCTION add_default(a integer, b integer DEFAULT 1)
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT a + b';
				CREATE VIEW v_default_arg AS
					SELECT add_default(2) AS value;
			`,
		},
		{
			name: "view_user_variadic_function_expanded_arguments",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_arg AS
					SELECT count_variadic(1, 2, 3) AS value;
			`,
		},
		{
			name: "view_user_variadic_function_explicit_array",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_arg AS
					SELECT count_variadic(VARIADIC ARRAY[1, 2, 3]::integer[]) AS value;
			`,
		},
		{
			name: "view_user_variadic_function_one_arg",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_one AS
					SELECT count_variadic(1) AS value;
			`,
		},
		{
			name: "view_user_variadic_function_mixed_coercible_args",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums bigint[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_mixed AS
					SELECT count_variadic(1::smallint, 2::integer, 3::bigint) AS value;
			`,
		},
		{
			name: "view_user_variadic_function_explicit_null_array",
			sql: `
				CREATE FUNCTION count_variadic(VARIADIC nums integer[])
				RETURNS integer
				LANGUAGE sql
				AS 'SELECT cardinality(nums)';
				CREATE VIEW v_variadic_null AS
					SELECT count_variadic(VARIADIC NULL::integer[]) AS value;
			`,
		},
		{
			name: "view_jsonb_extract_with_unnest_text_alias_column",
			sql: `
				CREATE TABLE docs (data jsonb);
				CREATE VIEW v_docs AS
					SELECT d.data ->> c.col_name AS value
					FROM docs d,
					     unnest(ARRAY['name', 'status']::text[]) AS c(col_name);
			`,
		},
		{
			name: "view_jsonb_extract_with_unnest_varchar_alias_column",
			sql: `
				CREATE TABLE docs (data jsonb);
				CREATE VIEW v_docs AS
					SELECT d.data ->> c.col_name AS value
					FROM docs d,
					     unnest(ARRAY['name', 'status']::varchar[]) AS c(col_name);
			`,
		},
		{
			name: "view_jsonb_extract_with_cte_alias_column",
			sql: `
				CREATE TABLE docs (data jsonb);
				CREATE VIEW v_docs AS
					WITH c(col_name) AS (SELECT 'name'::text)
					SELECT d.data ->> c.col_name AS value
					FROM docs d, c;
			`,
		},
		{
			name: "view_unnest_integer_alias_column",
			sql: `
				CREATE VIEW v_nums AS
					SELECT n.value + 1 AS next_value
					FROM unnest(ARRAY[1, 2, 3]::integer[]) AS n(value);
			`,
		},
		{
			name: "view_unnest_bigint_alias_column",
			sql: `
				CREATE VIEW v_nums AS
					SELECT n.value + 1 AS next_value
					FROM unnest(ARRAY[1, 2, 3]::bigint[]) AS n(value);
			`,
		},
		{
			name: "view_unnest_uuid_alias_column",
			sql: `
				CREATE VIEW v_ids AS
					SELECT n.value::text AS id_text
					FROM unnest(ARRAY['00000000-0000-0000-0000-000000000001']::uuid[]) AS n(value);
			`,
		},
		{
			name: "view_unnest_jsonb_alias_column",
			sql: `
				CREATE VIEW v_payloads AS
					SELECT n.value ->> 'name' AS name
					FROM unnest(ARRAY['{"name":"one"}'::jsonb]::jsonb[]) AS n(value);
			`,
		},
		{
			name: "view_unnest_domain_array_alias_column",
			sql: `
				CREATE DOMAIN loader_text_domain AS text;
				CREATE VIEW v_domain_values AS
					SELECT n.value::text AS value
					FROM unnest(ARRAY['one']::loader_text_domain[]) AS n(value);
			`,
		},
		{
			name: "view_array_polymorphic_builtins",
			sql: `
				CREATE VIEW v_arrays AS
					SELECT
						array_length(ARRAY[1, 2, 3]::integer[], 1) AS len,
						array_position(ARRAY['a', 'b']::text[], 'b') AS text_pos,
						array_position(ARRAY[1, 2, 3]::integer[], 2) AS int_pos,
						array_append(ARRAY[1, 2]::integer[], 3) AS appended,
						array_prepend(0, ARRAY[1, 2]::integer[]) AS prepended;
			`,
		},
		{
			name: "view_anycompatible_common_type_builtins",
			sql: `
				CREATE VIEW v_common AS
					SELECT
						coalesce(NULL, 1, 2::bigint) AS c,
						greatest(1, 2::bigint) AS g,
						least(1, 2::bigint) AS l;
			`,
		},
		{
			name: "view_record_returning_builtin_alias_columns",
			sql: `
				CREATE VIEW v_each AS
					SELECT e.key, e.value
					FROM jsonb_each_text('{"a":"b"}'::jsonb) AS e(key, value);
			`,
		},
		{
			name: "view_record_returning_user_function_alias_columns",
			sql: `
				CREATE FUNCTION loader_pair(OUT id integer, OUT name text)
				RETURNS record
				LANGUAGE sql
				AS 'SELECT 1, ''one''';
				CREATE VIEW v_pair AS
					SELECT p.id, p.name
					FROM loader_pair() AS p(id, name);
			`,
		},
		{
			name: "view_returns_table_function_alias_columns",
			sql: `
				CREATE FUNCTION loader_table()
				RETURNS TABLE(id integer, name text)
				LANGUAGE sql
				AS 'SELECT 1, ''one''';
				CREATE VIEW v_table AS
					SELECT t.id, t.name
					FROM loader_table() AS t(id, name);
			`,
		},
		{
			name: "view_srf_alias_column_list",
			sql: `
				CREATE VIEW v_srf AS
					SELECT x.value + 1 AS next_value
					FROM generate_series(1, 3) AS x(value);
			`,
		},
		{
			name: "view_unknown_literal_and_default_resolution",
			sql: `
				CREATE FUNCTION loader_default_one(a integer, b integer DEFAULT 1)
				RETURNS integer LANGUAGE sql AS 'SELECT a + b';
				CREATE FUNCTION loader_default_two(a integer, b integer DEFAULT 1, c integer DEFAULT 2)
				RETURNS integer LANGUAGE sql AS 'SELECT a + b + c';
				CREATE FUNCTION loader_text_arg(a text) RETURNS text LANGUAGE sql AS 'SELECT a';
				CREATE FUNCTION loader_int_arg(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
				CREATE FUNCTION loader_null_arg(a text, b integer DEFAULT 1) RETURNS text LANGUAGE sql AS 'SELECT a';
				CREATE VIEW v_resolve AS
					SELECT
						loader_default_one(1) AS d1,
						loader_default_two(1) AS d2,
						loader_text_arg('x') AS text_value,
						loader_int_arg(1) AS int_value,
						loader_null_arg(NULL) AS null_value;
			`,
		},
		{
			name: "view_operator_resolution_matrix",
			sql: `
				CREATE TABLE op_items (
					j json,
					jb jsonb,
					t text,
					i integer,
					b bigint,
					n numeric,
					d date,
					ts timestamp,
					ia integer[],
					r int4range
				);
				CREATE VIEW v_ops AS
					SELECT
						jb -> 'name' AS jb_obj_text,
						jb -> 0 AS jb_obj_int,
						jb ->> 'name' AS jb_text_text,
						jb ->> 0 AS jb_text_int,
						j -> 'name' AS j_obj_text,
						j ->> 'name' AS j_text_text,
						t || 'x' AS text_unknown_right,
						'x' || t AS text_unknown_left,
						i + b AS int_bigint,
						b + i AS bigint_int,
						n + i AS numeric_int,
						d + 1 AS date_plus_int,
						ts + interval '1 day' AS ts_plus_interval,
						t LIKE 'a%' AS text_like,
						t ILIKE 'a%' AS text_ilike,
						ia @> ARRAY[1]::integer[] AS array_contains,
						jb @> '{"a":1}'::jsonb AS jsonb_contains,
						r && int4range(1, 3) AS range_overlaps
					FROM op_items;
			`,
		},
		{
			name: "view_cte_range_resolution",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_base_items AS
					WITH cte1 AS (SELECT id FROM base_items)
					SELECT id FROM cte1;
			`,
		},
		{
			name: "view_cte_column_alias_resolution",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_base_items AS
					WITH cte1(item_id) AS (SELECT id FROM base_items)
					SELECT item_id FROM cte1;
			`,
		},
		{
			name: "view_cte_shadows_base_table",
			sql: `
				CREATE TABLE cte1 (id integer);
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_base_items AS
					WITH cte1 AS (SELECT id FROM base_items)
					SELECT id FROM cte1;
			`,
		},
		{
			name: "view_multiple_ctes_reference_prior_cte",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_multi_cte AS
					WITH c1 AS (SELECT id FROM base_items),
					     c2 AS (SELECT id + 1 AS next_id FROM c1)
					SELECT next_id FROM c2;
			`,
		},
		{
			name: "view_cte_same_column_names_as_base_table",
			sql: `
				CREATE TABLE base_items (id integer, value text);
				CREATE VIEW v_cte_same_cols AS
					WITH c AS (SELECT id, value FROM base_items)
					SELECT id, value FROM c;
			`,
		},
		{
			name: "view_nested_subquery_references_outer_cte",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_nested_outer_cte AS
					WITH c AS (SELECT id FROM base_items)
					SELECT (SELECT max(id) FROM c) AS max_id;
			`,
		},
		{
			name: "view_nested_with_shadows_outer_cte",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_nested_shadow_cte AS
					WITH c AS (SELECT id FROM base_items)
					SELECT id
					FROM (
						WITH c AS (SELECT 2 AS id)
						SELECT id FROM c
					) s;
			`,
		},
		{
			name: "view_nested_with_references_outer_cte",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_nested_outer_ref AS
					WITH c AS (SELECT id FROM base_items)
					SELECT id
					FROM (
						WITH d AS (SELECT id + 1 AS id FROM c)
						SELECT id FROM d
					) s;
			`,
		},
		{
			name: "view_recursive_cte",
			sql: `
				CREATE VIEW v_recursive AS
					WITH RECURSIVE r(n) AS (
						SELECT 1
						UNION ALL
						SELECT n + 1 FROM r WHERE n < 3
					)
					SELECT n FROM r;
			`,
		},
		{
			name: "view_cte_materialized",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_cte_mat AS
					WITH c AS MATERIALIZED (SELECT id FROM base_items)
					SELECT id FROM c;
			`,
		},
		{
			name: "view_cte_not_materialized",
			sql: `
				CREATE TABLE base_items (id integer);
				CREATE VIEW v_cte_not_mat AS
					WITH c AS NOT MATERIALIZED (SELECT id FROM base_items)
					SELECT id FROM c;
			`,
		},
		{
			name: "view_cte_column_aliases_fewer_than_output",
			sql: `
				CREATE VIEW v_cte_fewer_aliases AS
					WITH c(a) AS (SELECT 1 AS x, 2 AS y)
					SELECT a, y FROM c;
			`,
		},
		{
			name: "view_later_lateral_references_previous_lateral_alias",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_lateral_base AS
					SELECT s2.next_id
					FROM lateral_base b,
					     LATERAL (SELECT b.id AS id) s1,
				     LATERAL (SELECT s1.id + 1 AS next_id) s2;
			`,
		},
		{
			name: "view_cross_join_lateral_references_left",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_cross_lateral AS
					SELECT s.id
					FROM lateral_base b
					CROSS JOIN LATERAL (SELECT b.id AS id) s;
			`,
		},
		{
			name: "view_inner_join_lateral_references_left",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_inner_lateral AS
					SELECT s.id
					FROM lateral_base b
					INNER JOIN LATERAL (SELECT b.id AS id) s ON true;
			`,
		},
		{
			name: "view_lateral_references_multiple_prior_items",
			sql: `
				CREATE TABLE a (id integer);
				CREATE TABLE b (id integer);
				CREATE VIEW v_multi_lateral AS
					SELECT s.sum_id
					FROM a, b,
					     LATERAL (SELECT a.id + b.id AS sum_id) s;
			`,
		},
		{
			name: "view_third_lateral_references_two_prior_lateral_aliases",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_chain_lateral AS
					SELECT s3.total
					FROM lateral_base b,
					     LATERAL (SELECT b.id AS id1) s1,
					     LATERAL (SELECT s1.id1 + 1 AS id2) s2,
					     LATERAL (SELECT s1.id1 + s2.id2 AS total) s3;
			`,
		},
		{
			name: "view_lateral_unnest_references_left",
			sql: `
				CREATE TABLE array_items (id integer, vals integer[]);
				CREATE VIEW v_lateral_unnest AS
					SELECT u.val + array_items.id AS value
					FROM array_items,
					     LATERAL unnest(array_items.vals) AS u(val);
			`,
		},
		{
			name: "view_lateral_generate_series_references_left",
			sql: `
				CREATE TABLE series_items (n integer);
				CREATE VIEW v_lateral_series AS
					SELECT g.x
					FROM series_items,
					     LATERAL generate_series(1, series_items.n) AS g(x);
			`,
		},
		{
			name: "view_rows_from_with_ordinality",
			sql: `
				CREATE VIEW v_rows_from AS
					SELECT x.val, x.ord
					FROM ROWS FROM (unnest(ARRAY['a', 'b']::text[])) WITH ORDINALITY AS x(val, ord);
			`,
		},
		{
			name: "view_lateral_record_coldeflist",
			sql: `
				CREATE TABLE json_items (payload jsonb);
				CREATE VIEW v_lateral_record AS
					SELECT r.id, r.name
					FROM json_items,
					     LATERAL jsonb_to_record(json_items.payload) AS r(id integer, name text);
			`,
		},
		{
			name: "view_lateral_join_scope_nested",
			sql: `
				CREATE TABLE a (id integer);
				CREATE TABLE b (id integer);
				CREATE VIEW v_lateral_join_scope AS
					SELECT s.sum_id
					FROM a
					JOIN b ON true
					JOIN LATERAL (SELECT a.id + b.id AS sum_id) s ON true;
			`,
		},
		{
			name: "view_correlated_subquery_matrix",
			sql: `
				CREATE TABLE parents (id integer);
				CREATE TABLE children (parent_id integer, value integer);
				CREATE VIEW v_correlated AS
					SELECT
						p.id,
						(SELECT max(c.value) FROM children c WHERE c.parent_id = p.id) AS max_value
					FROM parents p
					WHERE EXISTS (SELECT 1 FROM children c WHERE c.parent_id = p.id)
					  AND p.id IN (SELECT c.parent_id FROM children c WHERE c.value > p.id)
					  AND (SELECT count(*) FROM children c WHERE c.parent_id = p.id) > 0;
			`,
		},
		{
			name: "view_nested_correlated_subquery_levels_up_two",
			sql: `
				CREATE TABLE parents (id integer);
				CREATE TABLE children (parent_id integer, value integer);
				CREATE VIEW v_nested_correlated AS
					SELECT p.id
					FROM parents p
					WHERE EXISTS (
						SELECT 1
						FROM children c
						WHERE EXISTS (
							SELECT 1
							WHERE c.parent_id = p.id
						)
					);
			`,
		},
		{
			name: "view_star_from_unnest_alias",
			sql: `
				CREATE VIEW v_star_unnest AS
					SELECT *
					FROM unnest(ARRAY['a', 'b']::text[]) AS u(val);
			`,
		},
		{
			name: "view_star_from_multi_unnest_alias",
			sql: `
				CREATE VIEW v_star_multi_unnest AS
					SELECT *
					FROM unnest(ARRAY[1, 2]::integer[], ARRAY['a', 'b']::text[]) AS u(id, name);
			`,
		},
		{
			name: "view_jsonb_to_record_coldeflist",
			sql: `
				CREATE VIEW v_jsonb_record AS
					SELECT r.id, r.name
					FROM jsonb_to_record('{"id":1,"name":"one"}'::jsonb) AS r(id integer, name text);
			`,
		},
		{
			name: "view_jsonb_to_recordset_coldeflist",
			sql: `
				CREATE VIEW v_jsonb_recordset AS
					SELECT r.id, r.name
					FROM jsonb_to_recordset('[{"id":1,"name":"one"}]'::jsonb) AS r(id integer, name text);
			`,
		},
		{
			name: "view_srf_in_select_target",
			sql: `
				CREATE VIEW v_srf_target AS
					SELECT generate_series(1, 3) AS value;
			`,
		},
		{
			name: "view_star_expansion_through_cte",
			sql: `
				CREATE TABLE star_items (id integer, name text);
				CREATE VIEW v_star_cte AS
					WITH c AS (SELECT id, name FROM star_items)
					SELECT * FROM c;
			`,
		},
		{
			name: "view_star_expansion_through_lateral",
			sql: `
				CREATE TABLE star_items (id integer, vals integer[]);
				CREATE VIEW v_star_lateral AS
					SELECT *
					FROM star_items,
					     LATERAL unnest(star_items.vals) AS u(val);
			`,
		},
		{
			name: "view_column_alias_list_overrides_target_names",
			sql:  `CREATE VIEW v_alias_override(a, b) AS SELECT 1 AS x, 2 AS y;`,
		},
		{
			name: "view_column_aliases_fewer_than_output",
			sql:  `CREATE VIEW v_alias_fewer(a) AS SELECT 1 AS x, 2 AS y;`,
		},
		{
			name: "object_namespace_alter_rename_matrix",
			sql: `
				CREATE TABLE t (id integer);
				CREATE INDEX t_id_idx ON t (id);
				ALTER INDEX t_id_idx RENAME TO t_id_idx2;
				CREATE SEQUENCE s;
				ALTER SEQUENCE s RENAME TO s2;
				ALTER SEQUENCE s2 OWNER TO CURRENT_USER;
				CREATE VIEW v AS SELECT id FROM t;
				ALTER VIEW v RENAME TO v2;
				CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
				ALTER MATERIALIZED VIEW mv RENAME TO mv2;
				CREATE TYPE mood AS ENUM ('sad', 'ok');
				ALTER TYPE mood RENAME TO mood2;
				CREATE DOMAIN positive_int AS integer CHECK (VALUE > 0);
				ALTER DOMAIN positive_int RENAME TO positive_int2;
				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
				ALTER FUNCTION f(integer) RENAME TO f2;
				CREATE PROCEDURE p(a integer) LANGUAGE sql AS 'SELECT 1';
				ALTER PROCEDURE p(integer) RENAME TO p2;
				ALTER ROUTINE f2(integer) OWNER TO CURRENT_USER;
			`,
		},
		{
			name: "object_namespace_comment_drop_matrix",
			sql: `
				CREATE TABLE t (id integer);
				CREATE INDEX t_id_idx ON t (id);
				COMMENT ON INDEX t_id_idx IS 'idx';
				COMMENT ON INDEX t_id_idx IS NULL;
				CREATE SEQUENCE s;
				COMMENT ON SEQUENCE s IS 'seq';
				COMMENT ON SEQUENCE s IS NULL;
				CREATE TYPE mood AS ENUM ('sad', 'ok');
				COMMENT ON TYPE mood IS 'type';
				COMMENT ON TYPE mood IS NULL;
				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
				COMMENT ON FUNCTION f(a integer) IS 'func';
				COMMENT ON FUNCTION f(integer) IS NULL;
				DROP FUNCTION f(integer);
				DROP INDEX t_id_idx;
				DROP SEQUENCE s;
				DROP TYPE mood;
			`,
		},
		{
			name: "object_namespace_grant_revoke_matrix",
			sql: `
				CREATE TABLE t (id integer);
				CREATE VIEW v AS SELECT id FROM t;
				CREATE SEQUENCE s;
				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
				GRANT SELECT ON TABLE t TO PUBLIC;
				REVOKE SELECT ON TABLE t FROM PUBLIC;
				GRANT SELECT ON TABLE v TO PUBLIC;
				REVOKE SELECT ON TABLE v FROM PUBLIC;
				GRANT USAGE ON SEQUENCE s TO PUBLIC;
				REVOKE USAGE ON SEQUENCE s FROM PUBLIC;
				GRANT EXECUTE ON FUNCTION f(integer) TO PUBLIC;
				REVOKE EXECUTE ON FUNCTION f(integer) FROM PUBLIC;
			`,
		},
		{
			name: "object_namespace_set_schema_matrix",
			sql: `
				CREATE SCHEMA target_schema;
				CREATE TABLE t (id integer);
				ALTER TABLE t SET SCHEMA target_schema;
				CREATE SEQUENCE s;
				ALTER SEQUENCE s SET SCHEMA target_schema;
				CREATE TYPE mood AS ENUM ('sad', 'ok');
				ALTER TYPE mood SET SCHEMA target_schema;
				CREATE FUNCTION f(a integer) RETURNS integer LANGUAGE sql AS 'SELECT a';
				ALTER FUNCTION f(integer) SET SCHEMA target_schema;
			`,
		},
		{
			name: "alter_index_depends_on_extension",
			sql: `
				CREATE TABLE t (id integer);
				CREATE INDEX t_id_idx ON t (id);
				ALTER INDEX t_id_idx DEPENDS ON EXTENSION plpgsql;
			`,
		},
		{
			name: "drop_owned_by_loaded_objects",
			sql: `
				CREATE ROLE loader_owner;
				CREATE TABLE t (id integer);
				CREATE SEQUENCE s;
				ALTER TABLE t OWNER TO loader_owner;
				ALTER SEQUENCE s OWNER TO loader_owner;
				DROP OWNED BY loader_owner;
				DROP ROLE loader_owner;
			`,
		},
		{
			name: "view_left_join_lateral_references_left_relation",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_lateral_left AS
					SELECT b.id, s.next_id
					FROM lateral_base b
					LEFT JOIN LATERAL (SELECT b.id + 1 AS next_id) s ON true;
			`,
		},
		{
			name: "view_chained_lateral_references_base_and_prior_alias",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_lateral_base AS
					SELECT s3.value
					FROM lateral_base b,
					     LATERAL (SELECT b.id AS id) s1,
					     LATERAL (SELECT s1.id + b.id AS id2) s2,
					     LATERAL (SELECT s2.id2 + s1.id AS value) s3;
			`,
		},
		{
			name: "view_cross_join_lateral_references_left_relation",
			sql: `
				CREATE TABLE lateral_base (id integer);
				CREATE VIEW v_lateral_cross AS
					SELECT b.id, s.next_id
					FROM lateral_base b
					CROSS JOIN LATERAL (SELECT b.id + 1 AS next_id) s;
			`,
		},
		{
			name: "alter_index_attach_partition",
			sql: `
				CREATE TABLE parent_idx_attach (id integer) PARTITION BY RANGE (id);
				CREATE TABLE child_idx_attach PARTITION OF parent_idx_attach
					FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_idx_attach_id_idx ON ONLY parent_idx_attach (id);
				CREATE INDEX child_idx_attach_id_idx ON child_idx_attach (id);
				ALTER INDEX parent_idx_attach_id_idx ATTACH PARTITION child_idx_attach_id_idx;
			`,
		},
		{
			name: "alter_index_attach_partition_schema_qualified",
			sql: `
				CREATE SCHEMA part_s;
				CREATE TABLE part_s.parent_idx_attach (id integer) PARTITION BY RANGE (id);
				CREATE TABLE part_s.child_idx_attach PARTITION OF part_s.parent_idx_attach
					FOR VALUES FROM (0) TO (10);
				CREATE INDEX parent_idx_attach_id_idx ON ONLY part_s.parent_idx_attach (id);
				CREATE INDEX child_idx_attach_id_idx ON part_s.child_idx_attach (id);
				ALTER INDEX part_s.parent_idx_attach_id_idx ATTACH PARTITION part_s.child_idx_attach_id_idx;
			`,
		},
	}
}

func TestLoaderCompatZeroColumnTable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE zc ();`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	_, rel, err := c.findRelation("", "zc")
	if err != nil {
		t.Fatalf("findRelation failed: %v", err)
	}
	if got := len(rel.Columns); got != 0 {
		t.Fatalf("columns: got %d, want 0", got)
	}
}

func TestLoaderCompatForeignKeyMatchSimple(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE parent (id integer PRIMARY KEY);
		CREATE TABLE child (
			parent_id integer,
			FOREIGN KEY (parent_id) REFERENCES parent(id) MATCH SIMPLE
		);
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatCommentOnFunctionWithArgumentNames(t *testing.T) {
	_, err := LoadSQL(`
		CREATE FUNCTION f(a integer, b text) RETURNS integer
			LANGUAGE sql AS 'SELECT 1';
		COMMENT ON FUNCTION f(a integer, b text) IS 'comment';
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatReturnsTablePlpgsql(t *testing.T) {
	_, err := LoadSQL(`
		CREATE FUNCTION f()
		RETURNS TABLE(id integer, name text)
		LANGUAGE plpgsql
		AS $$
		BEGIN
			RETURN QUERY SELECT 1, 'one'::text;
		END
		$$;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatBigintForeignKeyReferencesIntegerPrimaryKey(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE parent (id integer PRIMARY KEY);
		CREATE TABLE child (
			parent_id bigint REFERENCES parent(id)
		);
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatViewConcatWSVariadicBuiltin(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE items (a text, b integer);
		CREATE VIEW v_items AS
			SELECT concat_ws('-', a, b) AS label
			FROM items;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatViewJsonbExtractWithUnnestAliasColumn(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE docs (data jsonb);
		CREATE VIEW v_docs AS
			SELECT d.data ->> c.col_name AS value
			FROM docs d,
			     unnest(ARRAY['name', 'status']::text[]) AS c(col_name);
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatViewCTERangeResolution(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE base_items (id integer);
		CREATE VIEW v_base_items AS
			WITH cte1 AS (SELECT id FROM base_items)
			SELECT id FROM cte1;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatViewLaterLateralReferencesPreviousLateralAlias(t *testing.T) {
	_, err := LoadSQL(`
		CREATE TABLE lateral_base (id integer);
		CREATE VIEW v_lateral_base AS
			SELECT s2.next_id
			FROM lateral_base b,
			     LATERAL (SELECT b.id AS id) s1,
			     LATERAL (SELECT s1.id + 1 AS next_id) s2;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}
}

func TestLoaderCompatAlterIndexAttachPartition(t *testing.T) {
	c, err := LoadSQL(`
		CREATE TABLE parent_idx_attach (id integer) PARTITION BY RANGE (id);
		CREATE TABLE child_idx_attach PARTITION OF parent_idx_attach
			FOR VALUES FROM (0) TO (10);
		CREATE INDEX parent_idx_attach_id_idx ON ONLY parent_idx_attach (id);
		CREATE INDEX child_idx_attach_id_idx ON child_idx_attach (id);
		ALTER INDEX parent_idx_attach_id_idx ATTACH PARTITION child_idx_attach_id_idx;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	schema, err := c.resolveTargetSchema("")
	if err != nil {
		t.Fatalf("resolve schema: %v", err)
	}
	parentIdx := schema.Indexes["parent_idx_attach_id_idx"]
	childIdx := schema.Indexes["child_idx_attach_id_idx"]
	if parentIdx == nil || childIdx == nil {
		t.Fatalf("expected parent and child indexes to be registered")
	}
	for _, dep := range c.deps {
		if dep.ObjType == 'i' && dep.ObjOID == childIdx.OID &&
			dep.RefType == 'i' && dep.RefOID == parentIdx.OID &&
			dep.DepType == DepInternal {
			return
		}
	}
	t.Fatalf("missing internal child-index to parent-index dependency")
}

func TestLoaderCompatPhase3PartitionTableAttachDependency(t *testing.T) {
	c, err := LoadSQL(`
		CREATE TABLE parent_attach_dep (id integer) PARTITION BY RANGE (id);
		CREATE TABLE child_attach_dep (id integer);
		ALTER TABLE parent_attach_dep ATTACH PARTITION child_attach_dep
			FOR VALUES FROM (0) TO (10);
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	schema, err := c.resolveTargetSchema("")
	if err != nil {
		t.Fatalf("resolve schema: %v", err)
	}
	parent := schema.Relations["parent_attach_dep"]
	child := schema.Relations["child_attach_dep"]
	if parent == nil || child == nil {
		t.Fatalf("expected parent and child relations to be registered")
	}
	if !hasDependency(c, 'r', child.OID, 0, 'r', parent.OID, 0, DepAuto) {
		t.Fatalf("missing auto child-table to parent-table dependency")
	}
}

func TestLoaderCompatPhase3ExpressionDependencies(t *testing.T) {
	c, err := LoadSQL(`
		CREATE SEQUENCE loader_dep_seq;
		CREATE FUNCTION loader_dep_next_code() RETURNS integer LANGUAGE sql AS 'SELECT 1';
		CREATE FUNCTION loader_dep_is_positive(x integer) RETURNS boolean LANGUAGE sql AS 'SELECT $1 > 0';
		CREATE TABLE loader_dep_t (
			base integer,
			generated integer GENERATED ALWAYS AS (base + 1) STORED,
			code integer DEFAULT loader_dep_next_code(),
			seq_id integer DEFAULT nextval('loader_dep_seq'),
			CONSTRAINT loader_dep_positive CHECK (loader_dep_is_positive(base))
		);
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	schema, err := c.resolveTargetSchema("")
	if err != nil {
		t.Fatalf("resolve schema: %v", err)
	}
	rel := schema.Relations["loader_dep_t"]
	if rel == nil {
		t.Fatalf("expected table to be registered")
	}
	nextCode := c.findUserProcsByName(schema, "loader_dep_next_code")
	if len(nextCode) != 1 {
		t.Fatalf("expected loader_dep_next_code function, got %d", len(nextCode))
	}
	isPositive := c.findUserProcsByName(schema, "loader_dep_is_positive")
	if len(isPositive) != 1 {
		t.Fatalf("expected loader_dep_is_positive function, got %d", len(isPositive))
	}
	seq := schema.Sequences["loader_dep_seq"]
	if seq == nil {
		t.Fatalf("expected sequence to be registered")
	}

	generatedCol := rel.Columns[1]
	codeCol := rel.Columns[2]
	seqCol := rel.Columns[3]

	if !hasDependency(c, 'r', rel.OID, int32(generatedCol.AttNum), 'r', rel.OID, int32(rel.Columns[0].AttNum), DepNormal) {
		t.Fatalf("missing generated column dependency on referenced base column")
	}
	if !hasDependency(c, 'r', rel.OID, int32(codeCol.AttNum), 'f', nextCode[0].OID, 0, DepNormal) {
		t.Fatalf("missing column default dependency on user function")
	}
	if !hasDependency(c, 'r', rel.OID, int32(seqCol.AttNum), 's', seq.OID, 0, DepNormal) {
		t.Fatalf("missing column default dependency on referenced sequence")
	}

	var checkOID uint32
	for _, con := range c.ConstraintsOf(rel.OID) {
		if con.Name == "loader_dep_positive" {
			checkOID = con.OID
			break
		}
	}
	if checkOID == 0 {
		t.Fatalf("expected check constraint to be registered")
	}
	if !hasDependency(c, 'c', checkOID, 0, 'f', isPositive[0].OID, 0, DepNormal) {
		t.Fatalf("missing check constraint dependency on user function")
	}
}

func TestLoaderCompatPhase5CTEViewDependencies(t *testing.T) {
	c, err := LoadSQL(`
		CREATE TABLE phase5_base (id integer, value integer);
		CREATE FUNCTION phase5_inc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT $1 + 1';
		CREATE VIEW phase5_view AS
			WITH c AS (
				SELECT id, phase5_inc(value) AS next_value
				FROM phase5_base
			)
			SELECT id, next_value FROM c;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	schema, err := c.resolveTargetSchema("")
	if err != nil {
		t.Fatalf("resolve schema: %v", err)
	}
	base := schema.Relations["phase5_base"]
	view := schema.Relations["phase5_view"]
	if base == nil || view == nil {
		t.Fatalf("expected base table and view to be registered")
	}
	procs := c.findUserProcsByName(schema, "phase5_inc")
	if len(procs) != 1 {
		t.Fatalf("expected phase5_inc function, got %d", len(procs))
	}
	if !hasDependency(c, 'r', view.OID, 0, 'r', base.OID, 0, DepNormal) {
		t.Fatalf("missing view dependency on base table through CTE")
	}
	if !hasDependency(c, 'r', view.OID, 0, 'r', base.OID, int32(base.Columns[0].AttNum), DepNormal) {
		t.Fatalf("missing view dependency on base column through CTE")
	}
	if !hasDependency(c, 'r', view.OID, 0, 'f', procs[0].OID, 0, DepNormal) {
		t.Fatalf("missing view dependency on function used inside CTE")
	}
}

func hasDependency(c *Catalog, objType byte, objOID uint32, objSubID int32, refType byte, refOID uint32, refSubID int32, depType DepType) bool {
	for _, dep := range c.deps {
		if dep.ObjType == objType && dep.ObjOID == objOID && dep.ObjSubID == objSubID &&
			dep.RefType == refType && dep.RefOID == refOID && dep.RefSubID == refSubID &&
			dep.DepType == depType {
			return true
		}
	}
	return false
}

func TestLoaderCompatAcceptCorpus(t *testing.T) {
	for _, tc := range loaderCompatAcceptCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadSQL(tc.sql); err != nil {
				t.Fatalf("LoadSQL failed: %v\nSQL:\n%s", err, tc.sql)
			}
		})
	}
}

func TestLoaderCompatRejectCorpus(t *testing.T) {
	for _, tc := range loaderCompatRejectCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadSQL(tc.sql); err == nil {
				t.Fatalf("LoadSQL unexpectedly succeeded\nSQL:\n%s", tc.sql)
			}
		})
	}
}
