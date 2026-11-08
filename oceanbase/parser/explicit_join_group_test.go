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
	// A conditionless join before a conditioned chain reassociates through the
	// chain's own left spine but still stops at the group: `a JOIN b RIGHT JOIN
	// c ON ... JOIN d ON ...` is `((a JOIN b) RIGHT JOIN c) JOIN d`, while the
	// same with (b RIGHT JOIN c) parenthesized keeps the group under `a JOIN`.
	st3, err := Parse("SELECT 1 FROM a JOIN b RIGHT JOIN c ON b.id = c.id JOIN d ON c.id = d.id")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	outer := st3.Items[0].(*ast.SelectStmt).From[0].(*ast.JoinClause)
	rj, ok := outer.Left.(*ast.JoinClause)
	if !ok || rj.Type != ast.JoinRight {
		t.Fatalf("outer.Left = %T, want the RIGHT join", outer.Left)
	}
	if _, ok := rj.Left.(*ast.JoinClause); !ok {
		t.Errorf("RIGHT join's left = %T, want the (a JOIN b) join", rj.Left)
	}

	st4, err := Parse("SELECT 1 FROM a JOIN (b RIGHT JOIN c ON b.id = c.id) JOIN d ON c.id = d.id")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	outer4 := st4.Items[0].(*ast.SelectStmt).From[0].(*ast.JoinClause)
	inner4 := outer4.Left.(*ast.JoinClause)
	if g, ok := inner4.Right.(*ast.JoinClause); !ok || g.Type != ast.JoinRight {
		t.Errorf("explicit group not preserved: inner.Right = %T", inner4.Right)
	}
	if l, ok := inner4.Left.(*ast.TableRef); !ok || l.Name != "a" {
		t.Errorf("inner.Left = %T, want table a", inner4.Left)
	}
}
