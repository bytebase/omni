package parser

import (
	"errors"
	"strings"
	"testing"
)

// collateCases pin where SQL Server 2022 accepts a COLLATE clause, as SET
// PARSEONLY ON reports it. {a}, {b} and {c} stand for three collation names.
//
//   - An expression operand takes one COLLATE: a second one straight after
//     it is Msg 156, in every expression position. Parentheses make a new
//     operand that takes its own, and so does the time zone of AT TIME ZONE.
//   - A column definition takes COLLATE once, directly after the data type:
//     after any other column option it is Msg 156.
//   - The collation name is a regular identifier: bracketed, quoted and
//     string forms are Msg 102.
var collateCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT x COLLATE {a} FROM t", true},
	{"SELECT x COLLATE database_default FROM t", true},
	{"SELECT x COLLATE {a}, y COLLATE {b} FROM t", true},
	{"SELECT (x COLLATE {a}) COLLATE {b} FROM t", true},
	{"SELECT ((x COLLATE {a}) COLLATE {b}) COLLATE {c} FROM t", true},
	{"SELECT -(x COLLATE {a}) COLLATE {b} FROM t", true},
	{"SELECT COALESCE(x COLLATE {a}, y) COLLATE {b} FROM t", true},
	{"SELECT CASE WHEN 1 = 1 THEN 'a' COLLATE {a} END COLLATE {b}", true},
	{"SELECT x COLLATE {a} AT TIME ZONE 'UTC' COLLATE {b} FROM t", true},
	{"SELECT d AT TIME ZONE 'UTC' COLLATE {a} AT TIME ZONE 'UTC' COLLATE {b} FROM t", true},
	{"SELECT 1 FROM t WHERE x COLLATE {a} = y COLLATE {b}", true},
	{"SELECT x FROM t ORDER BY (x COLLATE {a}) COLLATE {b}", true},

	{"SELECT x COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT x COLLATE {a} COLLATE {a} FROM t", false},
	{"SELECT x COLLATE {a} COLLATE {b} COLLATE {c} FROM t", false},
	{"SELECT x COLLATE {a} -- comment\n  COLLATE {b} FROM t", false},
	{"SELECT -x COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT 'a' + 'b' COLLATE {a} COLLATE {b}", false},
	{"SELECT (x + y) COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT CAST(x AS VARCHAR(10)) COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT CASE WHEN 1 = 1 THEN 'a' END COLLATE {a} COLLATE {b}", false},
	{"SELECT (SELECT 'a') COLLATE {a} COLLATE {b}", false},
	{"SELECT d AT TIME ZONE 'UTC' COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT 1 FROM t WHERE x COLLATE {a} COLLATE {b} = y", false},
	{"SELECT x FROM t ORDER BY x COLLATE {a} COLLATE {b}", false},
	{"UPDATE t SET x = y COLLATE {a} COLLATE {b}", false},
	{"DECLARE @v VARCHAR(10) = 'x' COLLATE {a} COLLATE {b}", false},

	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} SPARSE NULL)", true},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} PRIMARY KEY)", true},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} DEFAULT 'x' COLLATE {b})", true},
	{"CREATE TABLE t3 (a VARCHAR(10) DEFAULT ('x' COLLATE {a}) COLLATE {b})", true},
	{"CREATE TABLE t3 (x VARCHAR(10), a AS (x COLLATE {a}) COLLATE {b})", true},
	{"ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE {a} NOT NULL", true},
	{"DECLARE @t TABLE (a VARCHAR(10) COLLATE {a} NOT NULL)", true},
	{"CREATE DATABASE d2 COLLATE {a}", true},
	{"ALTER DATABASE testdb COLLATE {a}", true},

	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE {a} COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) NOT NULL COLLATE {a})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} NOT NULL COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) PRIMARY KEY COLLATE {a})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) CHECK (a <> '') COLLATE {a})", false},
	{"CREATE TABLE t3 (x VARCHAR(10), a AS x COLLATE {a} COLLATE {b})", false},
	{"ALTER TABLE t ADD c VARCHAR(10) NOT NULL COLLATE {a}", false},
	{"DECLARE @t TABLE (a VARCHAR(10) NOT NULL COLLATE {a})", false},
	{"CREATE TYPE tt AS TABLE (a VARCHAR(10) COLLATE {a} COLLATE {b})", false},

	{"SELECT x COLLATE [Latin1_General_CI_AS] FROM t", false},
	{`SELECT x COLLATE "Latin1_General_CI_AS" FROM t`, false},
	{"SELECT x COLLATE [database_default] FROM t", false},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE [Latin1_General_CI_AS])", false},
	{"ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE [Latin1_General_CI_AS]", false},
	{"CREATE DATABASE d2 COLLATE [Latin1_General_CI_AS]", false},
	{"ALTER DATABASE testdb COLLATE [Latin1_General_CI_AS]", false},
}

var collationNames = strings.NewReplacer(
	"{a}", "Latin1_General_CI_AS", "{b}", "Latin1_General_BIN", "{c}", "SQL_Latin1_General_CP1_CI_AS")

// TestCollateClauseOracle: omni accepts exactly the COLLATE clauses SQL
// Server accepts, and reports a rejected one where SqlScriptDOM does.
// omni used to take any number of COLLATE clauses in a row, anywhere in a
// column definition, and a bracketed or quoted collation name.
func TestCollateClauseOracle(t *testing.T) {
	oracle := startParserOracle(t)
	h := getHarness(t)
	for _, tc := range collateCases {
		sql := collationNames.Replace(tc.sql)
		t.Run(tc.sql, func(t *testing.T) {
			serverErr, err := oracle.parseError(sql)
			if err != nil {
				t.Fatal(err)
			}
			if (serverErr == nil) != tc.accept {
				t.Fatalf("SQL Server: %v, want accept = %v", serverErr, tc.accept)
			}
			_, omniErr := Parse(sql)
			if (omniErr == nil) != tc.accept {
				t.Fatalf("omni: %v, want accept = %v", omniErr, tc.accept)
			}
			if tc.accept {
				return
			}
			errs, _ := h.Shape(t, sql)["errors"].([]any)
			if len(errs) == 0 {
				t.Fatal("SqlScriptDOM accepted it")
			}
			first, _ := errs[0].(map[string]any)
			line, _ := first["line"].(float64)
			column, _ := first["column"].(float64)
			var pe *ParseError
			if !errors.As(omniErr, &pe) {
				t.Fatalf("omni error %v is not a *ParseError", omniErr)
			}
			if want := lineColumnOffset(sql, int(line), int(column)); pe.Position != want {
				t.Errorf("omni reports offset %d, SqlScriptDOM line %d column %d (offset %d)",
					pe.Position, int(line), int(column), want)
			}
		})
	}
}

// lineColumnOffset converts SqlScriptDOM's 1-based line and column to a byte
// offset into the ASCII sql.
func lineColumnOffset(sql string, line, column int) int {
	offset := 0
	for ; line > 1; line-- {
		offset += strings.IndexByte(sql[offset:], '\n') + 1
	}
	return offset + column - 1
}
