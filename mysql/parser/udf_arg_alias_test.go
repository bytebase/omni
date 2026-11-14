package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// TestUDFArgAlias holds a generic call's arguments to udf_expr, which may carry
// an alias; a native function refuses it with 1583 and a stored one with 1584,
// both past parsing (MySQL 8.0.46), and a qualified call takes none.
func TestUDFArgAlias(t *testing.T) {
	for _, sql := range []string{
		"DO abs(1 x)",
		"DO f(1 AS x, 2)",
		"DO `abs`(1 x)",
		"DO f(1 'x')",
		"DO f(1 AS 'x')",
		"DO f(1 AS '')",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"DO d.f(1 x)",
		"DO f(1 AS)",
		"DO sum(1 x)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): want a syntax error", sql)
		}
	}
	stmts, err := Parse("DO f(1, 2 AS y, 3 '')")
	if err != nil {
		t.Fatal(err)
	}
	fc := stmts.Items[0].(*nodes.DoStmt).Exprs[0].(*nodes.FuncCallExpr)
	if len(fc.ArgAliases) != 3 || fc.ArgAliases[0] != nil || fc.ArgAliases[1] == nil || *fc.ArgAliases[1] != "y" ||
		fc.ArgAliases[2] == nil || *fc.ArgAliases[2] != "" {
		t.Errorf("ArgAliases = %v, want [nil y '']", fc.ArgAliases)
	}
}
