package deparse

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// TestWindowOptionsKept: what a window carries past its bounds is written
// back, so a deparse keeps the query's meaning.
func TestWindowOptionsKept(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"SELECT SUM(1) OVER (ROWS UNBOUNDED PRECEDING EXCLUDE CURRENT ROW)", "exclude current row"},
		{"SELECT SUM(1) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING EXCLUDE NO OTHERS)", "exclude no others"},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		got := DeparseSelect(list.Items[0].(*ast.SelectStmt))
		if !strings.Contains(strings.ToLower(got), strings.ToLower(tc.want)) {
			t.Errorf("DeparseSelect(%q) = %q, want it to hold %q", tc.sql, got, tc.want)
		}
	}
}
