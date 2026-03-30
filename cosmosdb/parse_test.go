package cosmosdb_test

import (
	"testing"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/cosmosdb/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
		check     func(t *testing.T, stmts []cosmosdb.Statement)
	}{
		{
			name:      "simple select star",
			sql:       "SELECT * FROM c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Text != "SELECT * FROM c" {
					t.Errorf("Text = %q, want %q", s.Text, "SELECT * FROM c")
				}
				if _, ok := s.AST.(*ast.SelectStmt); !ok {
					t.Errorf("AST type = %T, want *ast.SelectStmt", s.AST)
				}
				if s.ByteStart != 0 {
					t.Errorf("ByteStart = %d, want 0", s.ByteStart)
				}
				if s.ByteEnd != 15 {
					t.Errorf("ByteEnd = %d, want 15", s.ByteEnd)
				}
				if s.Start.Line != 1 || s.Start.Column != 1 {
					t.Errorf("Start = %v, want {1 1}", s.Start)
				}
			},
		},
		{
			name:      "multiline",
			sql:       "SELECT\n  *\nFROM\n  c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Start.Line != 1 {
					t.Errorf("Start.Line = %d, want 1", s.Start.Line)
				}
				if s.End.Line != 4 {
					t.Errorf("End.Line = %d, want 4", s.End.Line)
				}
			},
		},
		{
			name:      "with trailing whitespace",
			sql:       "SELECT * FROM c   ",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				s := stmts[0]
				if s.Text != "SELECT * FROM c" {
					t.Errorf("Text = %q, want %q", s.Text, "SELECT * FROM c")
				}
				if s.ByteEnd != 15 {
					t.Errorf("ByteEnd = %d, want 15", s.ByteEnd)
				}
			},
		},
		{
			name:      "empty input",
			sql:       "",
			wantCount: 0,
		},
		{
			name:      "whitespace only",
			sql:       "   \n\t  ",
			wantCount: 0,
		},
		{
			name:      "select with value",
			sql:       "SELECT VALUE c.name FROM c",
			wantCount: 1,
			check: func(t *testing.T, stmts []cosmosdb.Statement) {
				sel := stmts[0].AST.(*ast.SelectStmt)
				if !sel.Value {
					t.Error("Value should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := cosmosdb.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if len(stmts) != tt.wantCount {
				t.Fatalf("got %d statements, want %d", len(stmts), tt.wantCount)
			}
			if tt.check != nil {
				tt.check(t, stmts)
			}
		})
	}
}
