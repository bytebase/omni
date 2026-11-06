package deparse

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestODBCEscapeRoundTrip: a temporal literal with a non-string argument is
// written back as the ODBC escape it came from, by every writer here, so the text
// parses again to the same literal. DATE 123 would not parse.
func TestODBCEscapeRoundTrip(t *testing.T) {
	for _, sql := range []string{
		"SELECT {d 123}",
		"SELECT {t now()}",
		"SELECT {ts 1 + 1}",
		"SELECT {ts {abc {d 1}}}",
	} {
		t.Run(sql, func(t *testing.T) {
			list, err := parser.Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
			sel := list.Items[0].(*ast.SelectStmt)
			want := sel.TargetList[0].(*ast.TemporalLit)
			for name, out := range map[string]string{
				"Deparse":       "SELECT " + Deparse(sel.TargetList[0]),
				"Restore":       "SELECT " + Restore(sel.TargetList[0]),
				"DeparseSelect": DeparseSelect(sel),
				"RestoreSelect": RestoreSelect(sel),
			} {
				if !strings.Contains(out, "{") {
					t.Errorf("%s wrote %q, want the escape", name, out)
					continue
				}
				again, err := parser.Parse(out)
				if err != nil {
					t.Errorf("%s wrote %q, which does not parse: %v", name, out, err)
					continue
				}
				item := again.Items[0].(*ast.SelectStmt).TargetList[0]
				if rt, ok := item.(*ast.ResTarget); ok {
					item = rt.Val
				}
				got, ok := item.(*ast.TemporalLit)
				if !ok || got.Type != want.Type || got.Expr == nil {
					t.Errorf("%s wrote %q, which parses to %#v, want a %s literal with an argument", name, out, item, want.Type)
				}
			}
		})
	}
}
