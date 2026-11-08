package parser

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// TestExplicitJoinGroupPreserved pins that an explicitly parenthesized join
// group is kept whole when a conditionless join precedes it, since the two
// nestings have different results: `a JOIN (b RIGHT JOIN c ON ...)` cross-joins
// a with the right join (no rows when a is empty), while `a JOIN b RIGHT JOIN c
// ON ...` binds left-deep to `(a JOIN b) RIGHT JOIN c` (which preserves c's
// rows). parseTableFactor returns a group as a bare JoinClause, so the chain
// re-associates only the joins it built itself.
func TestExplicitJoinGroupPreserved(t *testing.T) {
	// a JOIN (b RIGHT JOIN c ON ...): top is the JOIN, its right the whole group.
	st, err := Parse("SELECT 1 FROM a JOIN (b RIGHT JOIN c ON b.id = c.id)")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	top := st.Items[0].(*ast.SelectStmt).From[0].(*ast.JoinClause)
	if top.Type != ast.JoinInner {
		t.Errorf("top type = %v, want JOIN", top.Type)
	}
	if l, ok := top.Left.(*ast.TableRef); !ok || l.Name != "a" {
		t.Errorf("top.Left = %T, want table a", top.Left)
	}
	if r, ok := top.Right.(*ast.JoinClause); !ok || r.Type != ast.JoinRight {
		t.Errorf("top.Right = %T, want the RIGHT JOIN group kept whole", top.Right)
	}

	// a JOIN b RIGHT JOIN c ON ...: binds left-deep to (a JOIN b) RIGHT JOIN c.
	st2, err := Parse("SELECT 1 FROM a JOIN b RIGHT JOIN c ON b.id = c.id")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	top2 := st2.Items[0].(*ast.SelectStmt).From[0].(*ast.JoinClause)
	if top2.Type != ast.JoinRight {
		t.Errorf("top2 type = %v, want RIGHT JOIN", top2.Type)
	}
	if _, ok := top2.Left.(*ast.JoinClause); !ok {
		t.Errorf("top2.Left = %T, want the (a JOIN b) join", top2.Left)
	}
}
