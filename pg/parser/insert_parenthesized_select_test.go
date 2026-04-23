package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseInsertParenthesizedSelectRest(t *testing.T) {
	tests := []string{
		`INSERT INTO dst (SELECT x FROM src)`,
		`INSERT INTO dst ((SELECT x FROM src))`,
		`INSERT INTO dst (VALUES (1), (2))`,
		`EXPLAIN (COSTS OFF) INSERT INTO dst (SELECT * FROM src)`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}

func TestParseInsertParenthesizedSelectDoesNotBecomeColumnList(t *testing.T) {
	stmts, err := Parse(`INSERT INTO dst (SELECT x FROM src)`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	insert := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.InsertStmt)
	if insert.Cols != nil && len(insert.Cols.Items) > 0 {
		t.Fatalf("expected no column list, got %d columns", len(insert.Cols.Items))
	}
	if _, ok := insert.SelectStmt.(*nodes.SelectStmt); !ok {
		t.Fatalf("expected SelectStmt rest, got %T", insert.SelectStmt)
	}
}
