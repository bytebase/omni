package catalog

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestLimitParamMarkerBindingOrderQ guards that deparsing a LIMIT with `?`
// parameter markers through the catalog IR preserves their positional binding
// order: the comma form and the OFFSET-keyword form place the markers in
// opposite positions, so the original form must be reproduced.
func TestLimitParamMarkerBindingOrderQ(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT)")

	cases := []struct {
		sql  string
		want string // trailing LIMIT clause expected in the deparsed IR
	}{
		{"SELECT a FROM t LIMIT ? OFFSET ?", "limit ? offset ?"},
		{"SELECT a FROM t LIMIT ?, ?", "limit ?,?"},
		{"SELECT a FROM t LIMIT ?", "limit ?"},
	}
	for _, c2 := range cases {
		list, err := parser.Parse(c2.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", c2.sql, err)
			continue
		}
		q, err := c.AnalyzeSelectStmt(list.Items[0].(*ast.SelectStmt))
		if err != nil {
			t.Errorf("Analyze(%q): %v", c2.sql, err)
			continue
		}
		out := DeparseQuery(q)
		if !strings.HasSuffix(out, c2.want) {
			t.Errorf("DeparseQuery(%q) = %q, want it to end with %q", c2.sql, out, c2.want)
		}
	}
}
