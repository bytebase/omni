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

		// Forms the parser deliberately still rejects.
		{"substring_from", "SELECT SUBSTRING('abcdef' FROM 2)", false},
		{"trim_both_from", "SELECT TRIM(BOTH 'x' FROM 'xax')", false},
		{"position_in", "SELECT POSITION('b' IN 'abc')", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}
