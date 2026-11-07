package deparse

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
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

// TestGroupByDirectionDeparse checks that a DESC direction on a GROUP BY item
// round-trips and that ASC (the default) is dropped.
func TestGroupByDirectionDeparse(t *testing.T) {
	for _, tc := range []struct{ sql, wantContains, wantAbsent string }{
		{"SELECT a FROM t GROUP BY a DESC", "group by `a` desc", ""},
		{"SELECT a FROM t GROUP BY a ASC", "group by `a`", " asc"},
		{"SELECT a, b FROM t GROUP BY a ASC, b DESC", "group by `a`,`b` desc", " asc"},
	} {
		list, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		got := DeparseSelect(sel)
		if !strings.Contains(got, tc.wantContains) {
			t.Errorf("DeparseSelect(%q) = %q, want to contain %q", tc.sql, got, tc.wantContains)
		}
		if tc.wantAbsent != "" && strings.Contains(got, tc.wantAbsent) {
			t.Errorf("DeparseSelect(%q) = %q, should not contain %q", tc.sql, got, tc.wantAbsent)
		}
	}
}
