package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftDropSchemaExternalDatabaseParse(t *testing.T) {
	tests := []struct {
		sql       string
		names     []string
		missingOk bool
		behavior  int
	}{
		{
			sql:      "DROP SCHEMA s_spectrum DROP EXTERNAL DATABASE RESTRICT;",
			names:    []string{"s_spectrum"},
			behavior: int(nodes.DROP_RESTRICT),
		},
		{
			sql:      "DROP SCHEMA s_sales, s_profit, s_revenue DROP EXTERNAL DATABASE CASCADE;",
			names:    []string{"s_sales", "s_profit", "s_revenue"},
			behavior: int(nodes.DROP_CASCADE),
		},
		{
			sql:       "DROP SCHEMA IF EXISTS spectrum_schema DROP EXTERNAL DATABASE CASCADE;",
			names:     []string{"spectrum_schema"},
			missingOk: true,
			behavior:  int(nodes.DROP_CASCADE),
		},
		{
			sql:      "DROP SCHEMA external_data DROP EXTERNAL DATABASE;",
			names:    []string{"external_data"},
			behavior: int(nodes.DROP_RESTRICT),
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstDropStmt(t, tt.sql)
			if stmt.RemoveType != int(nodes.OBJECT_SCHEMA) {
				t.Fatalf("expected schema drop, got remove type %d", stmt.RemoveType)
			}
			if stmt.Missing_ok != tt.missingOk {
				t.Fatalf("expected missing_ok=%v, got %v", tt.missingOk, stmt.Missing_ok)
			}
			if stmt.Behavior != tt.behavior {
				t.Fatalf("expected behavior %d, got %d", tt.behavior, stmt.Behavior)
			}
			if stmt.Objects == nil || len(stmt.Objects.Items) != len(tt.names) {
				t.Fatalf("expected %d objects, got %#v", len(tt.names), stmt.Objects)
			}
			for i, expected := range tt.names {
				name := singleNameFromAnyNameListItem(t, stmt.Objects.Items[i])
				if name != expected {
					t.Fatalf("object %d: expected %q, got %q", i, expected, name)
				}
			}
			assertDefElemBool(t, stmt.Options, "drop_external_database", true)
		})
	}
}

func firstDropStmt(t *testing.T, sql string) *nodes.DropStmt {
	t.Helper()
	tree, err := Parse(sql)
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
	stmt, ok := raw.Stmt.(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected DropStmt, got %T", raw.Stmt)
	}
	return stmt
}

func singleNameFromAnyNameListItem(t *testing.T, item nodes.Node) string {
	t.Helper()
	name, ok := item.(*nodes.List)
	if !ok || len(name.Items) != 1 {
		t.Fatalf("expected single-part any_name list, got %#v", item)
	}
	part, ok := name.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected string name part, got %T", name.Items[0])
	}
	return part.Str
}

func assertDefElemBool(t *testing.T, list *nodes.List, name string, want bool) {
	t.Helper()
	elem := findDefElem(list, name)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", name, list)
	}
	arg, ok := elem.Arg.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected option %q boolean arg, got %T", name, elem.Arg)
	}
	if arg.Boolval != want {
		t.Fatalf("expected option %q=%v, got %v", name, want, arg.Boolval)
	}
}
