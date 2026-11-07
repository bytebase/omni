package parser

import "testing"

// TestCollateConvertWeightFormsOracle asserts omni agrees with MySQL 8.0 on
// COLLATE, CONVERT(... USING ...) and WEIGHT_STRING(... AS ...).
//
// MySQL 8.0 takes the reserved word BINARY as a collation name in a column
// option (collation_name has a BINARY_SYM arm) but not in an expression
// postfix, where `x COLLATE BINARY` is a 1064; it takes BINARY as a charset
// name in CONVERT(... USING ...); and WEIGHT_STRING carries an AS CHAR(n) /
// AS CHARACTER(n) / AS BINARY(n) clause whose length must be at least 1.
// assertOracleMatch parses each with omni and runs it on the live server,
// asserting both accept or both reject.
//
// Ref: collation_name, charset_name and the weight_string rules in
// mysql-server sql/sql_yacc.yy.
func TestCollateConvertWeightFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	stmts := []string{
		// COLLATE: BINARY is a column-option collation name, not an
		// expression one.
		"ALTER TABLE t ADD c JSON COLLATE BINARY",
		"ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin",
		"DO x COLLATE BINARY",
		"DO x COLLATE utf8mb4_bin",

		// CONVERT(expr USING charset): BINARY is a charset name here.
		"DO CONVERT(x USING BINARY)",
		"DO CONVERT(x USING utf8mb4)",

		// WEIGHT_STRING's AS clause: CHAR and CHARACTER are one keyword, the
		// length is required and must be >= 1.
		"DO WEIGHT_STRING(x)",
		"DO WEIGHT_STRING(x AS CHAR(1))",
		"DO WEIGHT_STRING(x AS CHARACTER(1))",
		"DO WEIGHT_STRING(x AS BINARY(1))",
		"DO WEIGHT_STRING(x AS CHAR(0))",
		"DO WEIGHT_STRING(x AS CHAR)",
		"DO WEIGHT_STRING(x AS DECIMAL(1))",
	}
	for _, sql := range stmts {
		t.Run(sql, func(t *testing.T) {
			assertOracleMatch(t, o, sql)
		})
	}
}
