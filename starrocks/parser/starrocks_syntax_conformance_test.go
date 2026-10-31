package parser

import "testing"

// TestStarRocksSyntaxConformance is the container-grounded conformance matrix
// for the expression and clause forms added alongside the Doris work. The
// StarRocks fork carries the same parser changes, but StarRocks is a different
// engine with its own grammar, so its acceptance has to be proven against
// StarRocks rather than inferred from the Doris run.
func TestStarRocksSyntaxConformance(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		// EXTRACT.
		{"extract_year", "SELECT EXTRACT(YEAR FROM NOW())", true},
		{"extract_month", "SELECT EXTRACT(MONTH FROM NOW())", true},
		{"extract_in_cte", "WITH c AS (SELECT EXTRACT(YEAR FROM NOW()) y) SELECT * FROM c", true},

		// Reserved keywords usable as function names.
		{"trim", "SELECT TRIM('  x  ')", true},
		{"left", "SELECT LEFT('abc', 2)", true},
		{"right", "SELECT RIGHT('abc', 2)", true},
		{"if", "SELECT IF(1, 2, 3)", true},
		{"database", "SELECT DATABASE()", true},
		{"schema", "SELECT SCHEMA()", true},
		{"user", "SELECT USER()", true},
		{"session_user", "SELECT SESSION_USER()", true},
		{"current_user", "SELECT CURRENT_USER()", true},
		{"connection_id", "SELECT CONNECTION_ID()", true},

		// Nilary datetime functions. Unlike Doris, only CURRENT_TIMESTAMP takes
		// a fractional-seconds precision here; the rest are strictly nilary.
		{"current_timestamp_precision", "SELECT CURRENT_TIMESTAMP(3)", true},
		{"current_time_precision", "SELECT CURRENT_TIME(3)", false},
		{"localtime_precision", "SELECT LOCALTIME(3)", false},
		{"localtimestamp_precision", "SELECT LOCALTIMESTAMP(3)", false},
		{"current_date_arg", "SELECT CURRENT_DATE(1)", false},
		{"current_user_args", "SELECT CURRENT_USER(1, 2)", false},
		{"current_timestamp_junk_args", "SELECT CURRENT_TIMESTAMP(foo, bar)", false},
		{"current_date_parens", "SELECT CURRENT_DATE()", true},
		{"current_time_parens", "SELECT CURRENT_TIME()", true},

		// Literals.
		{"hex_literal", "SELECT x'1f'", true},
		{"bit_literal", "SELECT b'101'", true},
		{"array_literal", "SELECT [1, 2, 3]", true},
		// A map constructor requires the MAP prefix here; bare braces are not a
		// literal.
		{"map_literal_prefixed", "SELECT map{'a': 1}", true},
		{"map_literal_typed", "SELECT map<varchar,int>{'a': 1}", true},
		{"map_literal_bare", "SELECT {'a': 1}", false},

		// Session and user variables.
		{"system_var", "SELECT @@version_comment", true},
		{"system_var_session_scope", "SELECT @@session.query_timeout", true},
		{"user_var", "SELECT @uservar", true},

		// Lambdas.
		{"lambda_array_map", "SELECT array_map(x -> x + 1, [1, 2, 3])", true},
		{"lambda_multi_param", "SELECT array_map((x, y) -> x + y, [1, 2], [3, 4])", true},
		// Unlike Doris, the parenthesized single-parameter form is accepted.
		{"lambda_single_param_parenthesized", "SELECT array_map((x) -> x + 1, [1, 2])", true},

		// The -> JSON path operator (BYT-10088). Unlike Doris, which has no
		// `->` at all, StarRocks extracts a JSON path with it. The engine's
		// rule, pinned by these cases: `X -> Y` is accepted when Y is a bare
		// string literal (the arrow, X any primary expression) or when X is an
		// identifier or identifier list (a lambda, Y any expression).
		//
		// omni implements the arrow branch in full. It keeps the lambda branch
		// scoped to function-argument position, where the analyzer is the only
		// consumer that can use one, so engine-parseable oddities such as
		// `SELECT x -> 1` stay rejected here. That under-acceptance predates
		// this operator and is not pinned below.
		{"json_arrow", "SELECT j->'$.a' FROM t", true},
		{"json_arrow_spaced", "SELECT j -> '$.a' FROM t", true},
		{"json_arrow_double_quoted_path", `SELECT j->"$.a" FROM t`, true},
		{"json_arrow_in_function_call", "SELECT length(j->'$.a') FROM t", true},
		{"json_arrow_chained", "SELECT j->'$.b'->'$.c' FROM t", true},
		{"json_arrow_aliased", "SELECT j->'$.a' AS v FROM t", true},
		{"json_arrow_in_where", "SELECT id FROM t WHERE j->'$.a' = 1", true},
		// A left operand of the wrong type is a semantic error on the engine,
		// not a syntax one, so each of these still parses.
		{"json_arrow_qualified_left", "SELECT t.j->'$.a' FROM t", true},
		{"json_arrow_paren_left", "SELECT (j)->'$.a' FROM t", true},
		{"json_arrow_func_left", `SELECT PARSE_JSON('{"a":1}')->'$.a'`, true},
		{"json_arrow_int_left", "SELECT 1->'$.a'", true},
		{"json_arrow_null_left", "SELECT NULL->'$.a'", true},
		// Suffix chaining in either order, and tighter binding than the
		// operators either side.
		{"json_arrow_then_subscript", "SELECT j->'$.b'['c'] FROM t", true},
		{"json_arrow_after_subscript", "SELECT j['b']->'$.c' FROM t", true},
		{"json_arrow_tighter_than_eq", "SELECT j->'$.a' = 1 FROM t", true},
		{"json_arrow_tighter_than_not", "SELECT NOT j->'$.a' FROM t", true},
		{"json_arrow_tighter_than_unary_minus", "SELECT -id->'$.a' FROM t", true},
		// The ambiguous corner: in higher-order argument position the engine
		// reads `-> '<string>'` as the arrow, then rejects the call in its
		// analyzer with "There are no lambda functions in high-order function
		// array_map" — so it parses, and omni takes the same reading, which is
		// also the one that keeps the left column in lineage.
		{"json_arrow_wins_over_lambda", "SELECT array_map(j -> '$.a', [1, 2])", true},
		{"json_arrow_wins_over_paren_lambda", "SELECT array_map((x) -> '$.a', [1, 2])", true},
		{"json_arrow_wins_then_concat", "SELECT array_map(x -> 'a' || 'b', [1, 2])", true},
		// A parenthesized parameter LIST is not a primary expression, so the
		// arrow branch cannot apply and the lambda reading stands.
		{"lambda_multi_param_string_body", "SELECT array_map((x, y) -> '$.a', [1, 2], [3, 4])", true},
		{"lambda_body_contains_arrow", `SELECT array_map(x -> PARSE_JSON('{}')->'$.a', [1, 2])`, true},
		// Satisfying neither branch is a syntax error on the engine.
		{"json_arrow_int_rhs_rejected", "SELECT 1 -> 1", false},
		{"json_arrow_string_lhs_int_rhs_rejected", "SELECT 'a' -> 1", false},
		// StarRocks has no unquoting arrow at all.
		{"json_double_arrow_rejected", "SELECT j->>'$.a' FROM t", false},

		// Optimizer hints.
		{"set_var_hint", "SELECT /*+ SET_VAR(query_timeout = 100) */ 1", true},

		// CONVERT / USING / BINARY. This engine has no USING charset clause at
		// all — not on CONVERT, not on CHAR — though it keeps CONVERT(x, type).
		{"convert_to_type", "SELECT CONVERT('1', SIGNED)", true},
		{"convert_using", "SELECT CONVERT('abc' USING utf8)", false},
		{"char_plain", "SELECT CHAR(65)", true},
		{"char_using", "SELECT CHAR(65 USING utf8)", false},
		{"using_on_aggregate", "SELECT SUM(a USING utf8) FROM t", false},
		{"binary_operator", "SELECT BINARY 'abc'", true},
		// Unlike Doris, the general prefix form is engine-valid here.
		{"binary_int", "SELECT BINARY 1", true},
		{"binary_func", "SELECT BINARY now()", true},
		// Unlike Doris, array literal elements take full expressions.
		{"array_expr_element", "SELECT [1+1]", true},
		{"array_column_element", "SELECT [a] FROM t", true},
		{"user_var_string", "SELECT @'quoted'", true},
		{"paren_select_nested", "((SELECT 1))", true},
		{"paren_setop_after_nested", "((SELECT 1) UNION SELECT 2)", true},
		{"paren_setop_after_with", "(WITH c AS (SELECT 1) SELECT 1 UNION SELECT 2)", true},
		{"paren_trailing_clauses", "(SELECT 1) ORDER BY 1 LIMIT 5", true},
		{"paren_nested_trailing_limit", "((SELECT 1)) LIMIT 5", true},
		{"paren_double_order_by", "(SELECT 1 ORDER BY 1) ORDER BY 1", true},
		{"paren_with_dml_rejected", "(WITH c AS (SELECT 1) DELETE FROM t)", false},
		{"union_paren_rhs_limit", "SELECT 1 UNION (SELECT 2) LIMIT 5", true},
		{"paren_inner_trailing_limit", "(SELECT 1 UNION (SELECT 2) LIMIT 5)", true},
		{"paren_double_limit_nested", "((SELECT 1 LIMIT 1)) LIMIT 2", true},
		{"paren_double_order_nested", "((SELECT 1 ORDER BY 1)) ORDER BY 2", true},
		{"paren_order_by_subquery", "(SELECT 1) ORDER BY (SELECT 1)", true},
		// Unlike Doris, clause order is strict here: a trailing ORDER BY
		// after the set operation's LIMIT is engine-rejected.
		{"union_limit_then_order_rejected", "SELECT 1 UNION SELECT 2 LIMIT 5 ORDER BY 1", false},
		{"select_limit_then_order_rejected", "SELECT 1 LIMIT 5 ORDER BY 1", false},
		// Unlike Doris, StarRocks has no Hive-style LATERAL VIEW at all —
		// its lateral form is `, [LATERAL] unnest(...)` (engine-verified).
		{"lateral_view_rejected", "SELECT * FROM person LATERAL VIEW EXPLODE(ARRAY(30, 60)) tableName AS c_age", false},
		{"lateral_view_chained_rejected", "SELECT * FROM t LATERAL VIEW EXPLODE([1]) a AS x LATERAL VIEW EXPLODE([2]) b AS y", false},
		{"build_index_rejected", "BUILD INDEX index1 ON table1", false},

		// Forms the parser deliberately still rejects.
		{"substring_from", "SELECT SUBSTRING('abcdef' FROM 2)", false},
		{"trim_both_from", "SELECT TRIM(BOTH 'x' FROM 'xax')", false},
		{"position_in", "SELECT POSITION('b' IN 'abc')", false},

		// --- Constructs previously hidden by the trailing-token swallow
		// (BYT-10084): each parsed as a valid prefix and silently dropped the
		// rest of the statement.
		{"limit_offset_comma", "SELECT * FROM t LIMIT 5, 10", true},
		{"string_alias", `SELECT (1 + 1) AS "20%"`, true},
		{"string_alias_empty", "SELECT 1 AS ''", true},
		{"element_at_literal", "SELECT [1, 2, 3][1]", true},
		{"element_at_column", "SELECT c1[1] FROM t", true},
		{"element_at_qualified", "SELECT t.c1[1] FROM t", true},
		{"qualified_column_arith", "SELECT t.a + 1 FROM t", true},
		// Unlike Doris, the slice form is engine-rejected here.
		{"array_slice_rejected", "SELECT [1, 2, 3][1:2]", false},
		{"grouping_sets", "SELECT a, SUM(b) FROM t GROUP BY GROUPING SETS ((a), ())", true},
		{"grouping_sets_trailing_item", "SELECT a FROM t GROUP BY GROUPING SETS ((a)), b", false},
		// TABLET is engine-valid; TABLESAMPLE and GROUP BY ... WITH ROLLUP are
		// not StarRocks syntax (engine-verified) and, with the strict
		// trailing-token check (BYT-10085), stay rejected instead of being
		// swallowed as trailing junk.
		{"tablet", "SELECT * FROM t TABLET(10001)", true},
		{"tablesample_rejected", "SELECT * FROM t TABLESAMPLE(1000 ROWS)", false},
		{"group_by_with_rollup_rejected", "SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP", false},
		{"show_databases_from", "SHOW DATABASES FROM default_catalog", true},
		{"stray_comment_close_rejected", "SELECT 1 */ 2", false},
		{"explain_incomplete_rejected", "EXPLAIN SELECT * FROM", false},
		{"explain_bare_rejected", "EXPLAIN", false},
		{"set_expr_incomplete_rejected", "SET x = 1 +", false},
		{"set_valueless_rejected", "SET x", false},
		{"show_where_incomplete_rejected", "SHOW TABLES WHERE (", false},
		{"show_from_missing_db_rejected", "SHOW TABLES FROM", false},
		{"set_names_missing_charset_rejected", "SET NAMES", false},
		{"set_names_dangling_collate_rejected", "SET NAMES utf8 COLLATE", false},
		{"set_names", "SET NAMES utf8", true},
		{"set_names_default", "SET NAMES DEFAULT", true},
		{"set_role", "SET ROLE admin_role", true},
		{"set_role_default", "SET ROLE DEFAULT", true},
		{"set_default_role_missing_to_rejected", "SET DEFAULT ROLE r1", false},
		{"set_default_role_qualified_user", "SET DEFAULT ROLE r1 TO 'alice'@'%'", true},
		{"set_role_none_combined", "SET ROLE NONE, r2", true},
		{"set_names_collate", "SET NAMES utf8 COLLATE utf8_general_ci", true},
		{"show_like_missing_pattern_rejected", "SHOW TABLES LIKE", false},
		{"paren_select", "(SELECT 1)", true},
		{"from_table_named_selected", "SELECT * FROM selected", true},
		{"create_table_primary_key", "CREATE TABLE conf_pk (id BIGINT NOT NULL, v VARCHAR(64)) PRIMARY KEY(id) DISTRIBUTED BY HASH(id)", true},

		// --- Property lists. The grammar gives every one of them as
		// property (COMMA property)* -- so at least one item,
		// comma-separated, not ending on a comma. Each clause below reaches a
		// different one of the parsers that used to carry its own copy of the
		// list loop with the comma optional, so each accepted all three
		// departures the engine refuses.
		{"props_create_table", `CREATE TABLE zzp1 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_create_table_missing_comma_rejected", `CREATE TABLE zzp2 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_create_table_trailing_comma_rejected", `CREATE TABLE zzp3 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v",)`, false},
		{"props_create_table_empty_rejected", `CREATE TABLE zzp4 (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES()`, false},

		{"props_create_database", `CREATE DATABASE zzpdb PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_create_database_missing_comma_rejected", `CREATE DATABASE zzpdb2 PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_create_database_trailing_comma_rejected", `CREATE DATABASE zzpdb3 PROPERTIES("k"="v",)`, false},
		{"props_create_database_empty_rejected", `CREATE DATABASE zzpdb4 PROPERTIES()`, false},

		{"props_create_external_catalog", `CREATE EXTERNAL CATALOG zzpc PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_create_external_catalog_missing_comma_rejected", `CREATE EXTERNAL CATALOG zzpc2 PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_create_external_catalog_trailing_comma_rejected", `CREATE EXTERNAL CATALOG zzpc3 PROPERTIES("k"="v",)`, false},
		{"props_create_external_catalog_empty_rejected", `CREATE EXTERNAL CATALOG zzpc4 PROPERTIES()`, false},

		// The bare list, with no PROPERTIES keyword in front of it.
		{"props_using_gin", `CREATE TABLE zzg1 (a INT, INDEX idx (a) USING GIN ("imp_type"="clucene","parser"="english")) DISTRIBUTED BY HASH(a) BUCKETS 1`, true},
		{"props_using_gin_missing_comma_rejected", `CREATE TABLE zzg2 (a INT, INDEX idx (a) USING GIN ("a"="b" "c"="d")) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_using_gin_trailing_comma_rejected", `CREATE TABLE zzg3 (a INT, INDEX idx (a) USING GIN ("a"="b",)) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},
		{"props_using_gin_empty_rejected", `CREATE TABLE zzg4 (a INT, INDEX idx (a) USING GIN ()) DISTRIBUTED BY HASH(a) BUCKETS 1`, false},

		{"props_alter_table_set", `ALTER TABLE t SET ("k"="v","k2"="v2")`, true},
		{"props_alter_table_set_missing_comma_rejected", `ALTER TABLE t SET ("k"="v" "k2"="v2")`, false},
		{"props_alter_table_set_trailing_comma_rejected", `ALTER TABLE t SET ("k"="v",)`, false},
		{"props_alter_table_set_empty_rejected", `ALTER TABLE t SET ()`, false},

		// One config, because that is all the engine takes here -- the
		// multi-item form is a reject, where Doris accepts it. That arity is
		// a difference of its own and this parser still accepts several, so
		// only the single-item form is pinned.
		{"props_admin_frontend_config", `ADMIN SET FRONTEND CONFIG ("k"="v")`, true},
		{"props_admin_frontend_config_missing_comma_rejected", `ADMIN SET FRONTEND CONFIG ("k"="v" "k2"="v2")`, false},
		{"props_admin_frontend_config_trailing_comma_rejected", `ADMIN SET FRONTEND CONFIG ("k"="v",)`, false},
		{"props_admin_frontend_config_empty_rejected", `ADMIN SET FRONTEND CONFIG ()`, false},

		{"props_create_mv", `CREATE MATERIALIZED VIEW zzpmv DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v","k2"="v2") AS SELECT 1 AS a`, true},
		{"props_create_mv_missing_comma_rejected", `CREATE MATERIALIZED VIEW zzpmv2 DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v" "k2"="v2") AS SELECT 1 AS a`, false},
		{"props_create_mv_trailing_comma_rejected", `CREATE MATERIALIZED VIEW zzpmv3 DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES("k"="v",) AS SELECT 1 AS a`, false},
		{"props_create_mv_empty_rejected", `CREATE MATERIALIZED VIEW zzpmv4 DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES() AS SELECT 1 AS a`, false},

		{"props_routine_load_job", `CREATE ROUTINE LOAD zzrl ON t PROPERTIES("k"="v","k2"="v2") FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, true},
		{"props_routine_load_job_missing_comma_rejected", `CREATE ROUTINE LOAD zzrl2 ON t PROPERTIES("k"="v" "k2"="v2") FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, false},
		{"props_routine_load_job_trailing_comma_rejected", `CREATE ROUTINE LOAD zzrl3 ON t PROPERTIES("k"="v",) FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`, false},
		{"props_routine_load_source_missing_comma_rejected", `CREATE ROUTINE LOAD zzrl4 ON t PROPERTIES("desired_concurrent_number"="3") FROM KAFKA("k"="v" "k2"="v2")`, false},
		{"props_routine_load_source_empty_rejected", `CREATE ROUTINE LOAD zzrl5 ON t PROPERTIES("desired_concurrent_number"="3") FROM KAFKA()`, false},

		{"props_outfile", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v","k2"="v2")`, true},
		{"props_outfile_missing_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v" "k2"="v2")`, false},
		{"props_outfile_trailing_comma_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES("k"="v",)`, false},
		{"props_outfile_empty_rejected", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES()`, false},

		// ALTER CATALOG ... PROPERTY (...) reaches the sixth list parser, but
		// it gets no rows here: the engine rejects that statement form
		// outright, well-formed list and all, so a parity row would pin the
		// form rather than the list. omni accepting a form the engine does not
		// have -- and rejecting ALTER CATALOG c SET ("k"="v"), which it does --
		// is a separate, pre-existing divergence; the list shape for that
		// parser is covered by TestPropertyListSeparators.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
