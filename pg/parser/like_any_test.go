package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseLikeIlikeAnyAllArray(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantKind nodes.A_Expr_Kind
		wantOp   string
	}{
		{
			name:     "like any",
			sql:      `SELECT 'foo' LIKE ANY (ARRAY['%a', '%o'])`,
			wantKind: nodes.AEXPR_OP_ANY,
			wantOp:   "~~",
		},
		{
			name:     "like all",
			sql:      `SELECT 'foo' LIKE ALL (ARRAY['f%', '%o'])`,
			wantKind: nodes.AEXPR_OP_ALL,
			wantOp:   "~~",
		},
		{
			name:     "not like any",
			sql:      `SELECT 'foo' NOT LIKE ANY (ARRAY['%a', '%b'])`,
			wantKind: nodes.AEXPR_OP_ANY,
			wantOp:   "!~~",
		},
		{
			name:     "not like all",
			sql:      `SELECT 'foo' NOT LIKE ALL (ARRAY['%a', '%o'])`,
			wantKind: nodes.AEXPR_OP_ALL,
			wantOp:   "!~~",
		},
		{
			name:     "ilike any",
			sql:      `SELECT 'foo' ILIKE ANY (ARRAY['%A', '%O'])`,
			wantKind: nodes.AEXPR_OP_ANY,
			wantOp:   "~~*",
		},
		{
			name:     "ilike all",
			sql:      `SELECT 'foo' ILIKE ALL (ARRAY['F%', '%O'])`,
			wantKind: nodes.AEXPR_OP_ALL,
			wantOp:   "~~*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
			target := sel.TargetList.Items[0].(*nodes.ResTarget)
			expr, ok := target.Val.(*nodes.A_Expr)
			if !ok {
				t.Fatalf("expected target A_Expr, got %T", target.Val)
			}
			if expr.Kind != tt.wantKind {
				t.Fatalf("A_Expr.Kind = %v, want %v", expr.Kind, tt.wantKind)
			}
			if got := firstString(expr.Name); got != tt.wantOp {
				t.Fatalf("operator = %q, want %q", got, tt.wantOp)
			}
		})
	}
}

func TestParseLikeAnySubquery(t *testing.T) {
	stmts, err := Parse(`SELECT 'foo' LIKE ANY (SELECT pattern FROM patterns)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	sel := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	target := sel.TargetList.Items[0].(*nodes.ResTarget)
	sub, ok := target.Val.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected target SubLink, got %T", target.Val)
	}
	if sub.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Fatalf("SubLinkType = %d, want %d", sub.SubLinkType, int(nodes.ANY_SUBLINK))
	}
	if got := firstString(sub.OperName); got != "~~" {
		t.Fatalf("operator = %q, want %q", got, "~~")
	}
}

func firstString(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	str, ok := list.Items[0].(*nodes.String)
	if !ok {
		return ""
	}
	return str.Str
}
