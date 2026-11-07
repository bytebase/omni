package deparse

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestLimitParamMarkerBindingOrder guards that deparsing a LIMIT with `?`
// parameter markers preserves their positional binding order. The comma form
// `LIMIT offset,count` and the keyword form `LIMIT count OFFSET offset` place
// the markers in opposite positions, so the original form must be reproduced.
func TestLimitParamMarkerBindingOrder(t *testing.T) {
	cases := []struct {
		sql  string
		want string // trailing LIMIT clause expected in the deparsed output
	}{
		{"SELECT a FROM t LIMIT ? OFFSET ?", "limit ? offset ?"},
		{"SELECT a FROM t LIMIT ?, ?", "limit ?,?"},
		{"SELECT a FROM t LIMIT ?", "limit ?"},
		// A literal offset can still normalize to the comma form: the operand
		// swap is value-preserving when the operands are not markers.
		{"SELECT a FROM t LIMIT 10 OFFSET 20", "limit 20,10"},
	}
	for _, c := range cases {
		list, err := parser.Parse(c.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.sql, err)
			continue
		}
		out := DeparseSelect(list.Items[0].(*ast.SelectStmt))
		if !hasSuffix(out, c.want) {
			t.Errorf("Deparse(%q) = %q, want it to end with %q", c.sql, out, c.want)
		}
	}
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
