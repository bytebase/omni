package deparse

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestWindowNullRoundTrip: IGNORE NULLS and FROM LAST are written back by every
// writer and parse again to the same flags. TiDB v8.5.5 parses these but refuses
// them at execution (1235), so a view never stores them; the text must still be
// omni's own valid SQL.
func TestWindowNullRoundTrip(t *testing.T) {
	for _, sql := range []string{
		"SELECT FIRST_VALUE(v) IGNORE NULLS OVER w",
		"SELECT NTH_VALUE(v, 2) FROM LAST OVER w",
		"SELECT NTH_VALUE(v, 2) FROM LAST IGNORE NULLS OVER w",
	} {
		t.Run(sql, func(t *testing.T) {
			list, err := parser.Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
			sel := list.Items[0].(*ast.SelectStmt)
			want := sel.TargetList[0].(*ast.FuncCallExpr)
			for name, out := range map[string]string{
				"DeparseSelect": DeparseSelect(sel),
				"RestoreSelect": RestoreSelect(sel),
			} {
				again, err := parser.Parse(out)
				if err != nil {
					t.Errorf("%s wrote %q, which does not parse: %v", name, out, err)
					continue
				}
				item := again.Items[0].(*ast.SelectStmt).TargetList[0]
				if rt, ok := item.(*ast.ResTarget); ok {
					item = rt.Val
				}
				got, ok := item.(*ast.FuncCallExpr)
				if !ok || got.IgnoreNull != want.IgnoreNull || got.FromLast != want.FromLast {
					t.Errorf("%s wrote %q, which parses to IgnoreNull/FromLast != %v/%v", name, out, want.IgnoreNull, want.FromLast)
				}
				if !strings.Contains(out, "OVER") {
					t.Errorf("%s wrote %q, missing OVER", name, out)
				}
			}
		})
	}
}
