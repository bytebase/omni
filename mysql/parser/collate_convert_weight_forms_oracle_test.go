package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestCollateConvertWeightFormsOracle asserts omni agrees with MySQL 8.0 on
// COLLATE, CONVERT(... USING ...) and WEIGHT_STRING(... AS ...).
//
// MySQL 8.0 takes the reserved word BINARY as a collation name in a column
// option (collation_name has a BINARY_SYM arm) but not in an expression
// postfix, where `x COLLATE BINARY` is a 1064; it takes BINARY as a charset
// name in CONVERT(... USING ...); and WEIGHT_STRING carries an AS CHAR(n) /
// AS CHARACTER(n) / AS BINARY(n) clause whose length must be at least 1,
// and no LEVEL clause.
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

		// WEIGHT_STRING's LEVEL clause: MySQL 5.7 had it and 8.0 removed it,
		// so every form is a 1064, alone or after an AS clause.
		"DO WEIGHT_STRING(x LEVEL 1)",
		"DO WEIGHT_STRING(x LEVEL 0)",
		"DO WEIGHT_STRING(x LEVEL 1 ASC)",
		"DO WEIGHT_STRING(x LEVEL 1 DESC)",
		"DO WEIGHT_STRING(x LEVEL 1 REVERSE)",
		"DO WEIGHT_STRING(x LEVEL 1 DESC REVERSE)",
		"DO WEIGHT_STRING(x LEVEL 1 REVERSE DESC)",
		"DO WEIGHT_STRING(x LEVEL 1, 2)",
		"DO WEIGHT_STRING(x LEVEL 1 DESC, 2 REVERSE)",
		"DO WEIGHT_STRING(x LEVEL 1 - 2)",
		"DO WEIGHT_STRING(x AS CHAR(4) LEVEL 1)",
		"DO WEIGHT_STRING(x AS BINARY(4) LEVEL 1)",
		"DO WEIGHT_STRING(x LEVEL)",
	}
	for _, sql := range stmts {
		t.Run(sql, func(t *testing.T) {
			assertOracleMatch(t, o, sql)
		})
	}
}

// TestWeightStringLevelErrorPositionOracle asserts Parse refuses the LEVEL
// clause where MySQL 8.0 does: the server's 1064 quotes the statement from the
// LEVEL token on, and the ParseError carries that token's offset.
func TestWeightStringLevelErrorPositionOracle(t *testing.T) {
	o := startParserOracle(t)
	stmts := []string{
		"DO WEIGHT_STRING(x LEVEL 1)",
		"DO WEIGHT_STRING(x LEVEL 1 DESC REVERSE)",
		"DO WEIGHT_STRING(x LEVEL 1, 2)",
		"DO WEIGHT_STRING(x LEVEL 1 - 2)",
		"DO WEIGHT_STRING(x AS CHAR(4) LEVEL 1)",
		"DO WEIGHT_STRING(x AS BINARY(4) LEVEL 1)",
	}
	for _, sql := range stmts {
		t.Run(sql, func(t *testing.T) {
			level := strings.Index(sql, "LEVEL")

			_, err := o.db.ExecContext(o.ctx, sql)
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) || myErr.Number != 1064 {
				t.Fatalf("server answered %v, want 1064", err)
			}
			if near := "near '" + sql[level:] + "'"; !strings.Contains(myErr.Message, near) {
				t.Errorf("server's error is %q, want it to say %s", myErr.Message, near)
			}

			_, err = Parse(sql)
			var parseErr *ParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("Parse answered %v, want a ParseError", err)
			}
			if parseErr.Position != level {
				t.Errorf("ParseError.Position = %d, want %d, the LEVEL token", parseErr.Position, level)
			}
		})
	}
}
