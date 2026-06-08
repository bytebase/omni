package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftDeleteWithoutFromParse(t *testing.T) {
	tree, err := Parse("DELETE category;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected DeleteStmt, got %T", raw.Stmt)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "category" {
		t.Fatalf("Relation = %#v, want category", stmt.Relation)
	}
}
