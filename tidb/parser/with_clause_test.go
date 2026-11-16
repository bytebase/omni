package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// A WITH clause is kept on the statement it belongs to, every statement form
// that takes one included. V1's prior-backup rule could not tell a reference
// to a CTE from a table without it.

func TestWithClauseOnDML(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
		of   func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr
	}{
		{
			name: "insert",
			sql:  "WITH a AS (SELECT 1), b AS (SELECT 2) INSERT INTO t SELECT * FROM a",
			want: []string{"a", "b"},
			of: func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr {
				stmt, ok := item.(*nodes.InsertStmt)
				if !ok {
					t.Fatalf("expected an INSERT, got %T", item)
				}
				return stmt.CTEs
			},
		},
		{
			name: "replace",
			sql:  "WITH a AS (SELECT 1) REPLACE INTO t SELECT * FROM a",
			want: []string{"a"},
			of: func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr {
				stmt, ok := item.(*nodes.InsertStmt)
				if !ok {
					t.Fatalf("expected a REPLACE, got %T", item)
				}
				return stmt.CTEs
			},
		},
		{
			name: "update",
			sql:  "WITH a AS (SELECT 1), b AS (SELECT 2) UPDATE t SET id = 1",
			want: []string{"a", "b"},
			of: func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr {
				stmt, ok := item.(*nodes.UpdateStmt)
				if !ok {
					t.Fatalf("expected an UPDATE, got %T", item)
				}
				return stmt.CTEs
			},
		},
		{
			name: "delete",
			sql:  "WITH a AS (SELECT 1) DELETE FROM t WHERE id IN (SELECT * FROM a)",
			want: []string{"a"},
			of: func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr {
				stmt, ok := item.(*nodes.DeleteStmt)
				if !ok {
					t.Fatalf("expected a DELETE, got %T", item)
				}
				return stmt.CTEs
			},
		},
		{
			name: "select",
			sql:  "WITH a AS (SELECT 1), b AS (SELECT 2) SELECT * FROM a",
			want: []string{"a", "b"},
			of: func(t *testing.T, item nodes.Node) []*nodes.CommonTableExpr {
				stmt, ok := item.(*nodes.SelectStmt)
				if !ok {
					t.Fatalf("expected a SELECT, got %T", item)
				}
				return stmt.CTEs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("expected one statement, got %d", len(list.Items))
			}
			ctes := tt.of(t, list.Items[0])
			var got []string
			for _, cte := range ctes {
				got = append(got, cte.Name)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("CTE names: want %v, got %v", tt.want, got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("CTE names: want %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func TestWithClauseAbsent(t *testing.T) {
	for _, sql := range []string{"UPDATE t SET id = 1", "DELETE FROM t", "INSERT INTO t VALUES (1)"} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		switch stmt := list.Items[0].(type) {
		case *nodes.UpdateStmt:
			if stmt.CTEs != nil {
				t.Fatalf("%q: unexpected CTEs %v", sql, stmt.CTEs)
			}
		case *nodes.DeleteStmt:
			if stmt.CTEs != nil {
				t.Fatalf("%q: unexpected CTEs %v", sql, stmt.CTEs)
			}
		case *nodes.InsertStmt:
			if stmt.CTEs != nil {
				t.Fatalf("%q: unexpected CTEs %v", sql, stmt.CTEs)
			}
		}
	}
}
