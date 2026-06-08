package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftSelectTopAndMinusParse(t *testing.T) {
	tree, err := Parse(`
SELECT TOP 10 * FROM sales ORDER BY amount DESC;
SELECT product_id FROM all_products MINUS SELECT product_id FROM discontinued_products;
SELECT TOP 100 * INTO top_sales FROM sales ORDER BY pricepaid DESC;
`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 3 {
		t.Fatalf("expected three statements, got %d", len(tree.Items))
	}

	top := rawSelectAt(t, tree, 0)
	if top.LimitCount == nil {
		t.Fatal("SELECT TOP should populate LimitCount")
	}
	if top.IntoClause != nil {
		t.Fatal("plain SELECT TOP should not populate IntoClause")
	}

	minus := rawSelectAt(t, tree, 1)
	if minus.Op != nodes.SETOP_EXCEPT {
		t.Fatalf("MINUS should parse as SETOP_EXCEPT, got %v", minus.Op)
	}
	if minus.Larg == nil || minus.Rarg == nil {
		t.Fatalf("MINUS should populate left and right set-op branches: %#v", minus)
	}

	topInto := rawSelectAt(t, tree, 2)
	if topInto.LimitCount == nil {
		t.Fatal("SELECT TOP ... INTO should populate LimitCount")
	}
	if topInto.IntoClause == nil || topInto.IntoClause.Rel == nil || topInto.IntoClause.Rel.Relname != "top_sales" {
		t.Fatalf("SELECT TOP ... INTO should populate IntoClause target, got %#v", topInto.IntoClause)
	}
}

func TestRedshiftCreateViewNoSchemaBindingParse(t *testing.T) {
	tree, err := Parse("CREATE OR REPLACE VIEW public.sales_v AS SELECT amount FROM public.sales WITH NO SCHEMA BINDING;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.ViewStmt)
	if !ok {
		t.Fatalf("expected ViewStmt, got %T", raw.Stmt)
	}
	if !stmt.Replace {
		t.Fatal("expected CREATE OR REPLACE view")
	}
	if stmt.Options == nil {
		t.Fatal("expected WITH NO SCHEMA BINDING option")
	}
	elem := findDefElem(stmt.Options, "no_schema_binding")
	if elem == nil {
		t.Fatalf("missing no_schema_binding option: %#v", stmt.Options.Items)
	}
}

func rawSelectAt(t *testing.T, tree *nodes.List, index int) *nodes.SelectStmt {
	t.Helper()
	raw, ok := tree.Items[index].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("statement %d: expected RawStmt, got %T", index, tree.Items[index])
	}
	stmt, ok := raw.Stmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("statement %d: expected SelectStmt, got %T", index, raw.Stmt)
	}
	return stmt
}
