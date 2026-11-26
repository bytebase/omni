package parser

import (
	"math"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestStringLengthTakesANumber pins the length a decimal string length reads as.
// string_length_i takes a number_literal (sql_parser_mysql_mode.y:7322), so
// OceanBase 4.3.5.6 parses CHAR(1.5) and answers 1074 for CHAR(1e5), a length
// past CHAR's limit; a length past int's range stays too long, as a run of
// digits too wide for an integer does. Verdicts: compat.tsv, suite:string-length.
func TestStringLengthTakesANumber(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want int
	}{
		{"ALTER TABLE a ADD a CHAR(1.5)", 1},
		{"ALTER TABLE a ADD a VARCHAR(.5)", 0},
		{"ALTER TABLE a ADD a BINARY(2.)", 2},
		{"ALTER TABLE a ADD a CHAR(1e5)", 100000},
		{"ALTER TABLE a ADD a CHAR(1e10)", 10000000000},
		{"ALTER TABLE a ADD a VARCHAR(2e400)", math.MaxInt},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		stmt := list.Items[0].(*nodes.AlterTableStmt)
		if got := stmt.Commands[0].Column.TypeName.Length; got != tc.want {
			t.Errorf("%s: length %d, want %d", tc.sql, got, tc.want)
		}
	}
}
