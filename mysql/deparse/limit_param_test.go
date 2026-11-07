package deparse

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// TestLimitParamMarkerDeparse checks that a `?` LIMIT operand round-trips. The
// comma form `limit offset,count` swaps the operands, which is value-preserving
// for literals but would reorder the positional bindings of `?` markers, so the
// OFFSET keyword form is reproduced when a marker is present.
func TestLimitParamMarkerDeparse(t *testing.T) {
	for _, tc := range []struct{ sql, wantSuffix string }{
		{"SELECT * FROM t LIMIT ?", "limit ?"},
		{"SELECT * FROM t LIMIT ?, ?", "limit ?,?"},
		{"SELECT * FROM t LIMIT ? OFFSET ?", "limit ? offset ?"},
		// Literals keep the comma normalization: LIMIT 2 OFFSET 3 -> limit 3,2.
		{"SELECT * FROM t LIMIT 2 OFFSET 3", "limit 3,2"},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		got := DeparseSelect(sel)
		if !strings.HasSuffix(got, tc.wantSuffix) {
			t.Errorf("DeparseSelect(%q) = %q, want suffix %q", tc.sql, got, tc.wantSuffix)
		}
	}
}
