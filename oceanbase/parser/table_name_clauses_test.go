package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestExternalPartitionColumnLoc holds the column of an external table's
// PARTITION (col = value) to the text that names it, quoted or not: a quoted
// name is longer in the source than the name it holds.
func TestExternalPartitionColumnLoc(t *testing.T) {
	for _, tc := range []struct{ sql, column, pair string }{
		{"SELECT * FROM t PARTITION (a = 1)", "a", "a = 1"},
		{"SELECT * FROM t PARTITION (`a` = 1)", "`a`", "`a` = 1"},
		{"SELECT * FROM t PARTITION (`a``b` = 's', c = 2)", "`a``b`", "`a``b` = 's'"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		ref := list.Items[0].(*nodes.SelectStmt).From[0].(*nodes.TableRef)
		pair := ref.PartitionValues[0]
		if got := tc.sql[pair.Column.Loc.Start:pair.Column.Loc.End]; got != tc.column {
			t.Errorf("Parse(%q): column spans %q, want %q", tc.sql, got, tc.column)
		}
		if got := tc.sql[pair.Loc.Start:pair.Loc.End]; got != tc.pair {
			t.Errorf("Parse(%q): pair spans %q, want %q", tc.sql, got, tc.pair)
		}
	}
}
