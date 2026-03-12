package pg

import (
	"testing"

	"github.com/bytebase/omni/pg/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantCount int
		checks    func(t *testing.T, stmts []Statement)
	}{
		{
			name:      "single SELECT",
			sql:       "SELECT 1",
			wantCount: 1,
			checks: func(t *testing.T, stmts []Statement) {
				s := stmts[0]
				if s.Text != "SELECT 1" {
					t.Errorf("text = %q, want %q", s.Text, "SELECT 1")
				}
				if _, ok := s.AST.(*ast.SelectStmt); !ok {
					t.Errorf("AST type = %T, want *ast.SelectStmt", s.AST)
				}
				if s.ByteStart != 0 || s.ByteEnd != 8 {
					t.Errorf("byte range = [%d,%d), want [0,8)", s.ByteStart, s.ByteEnd)
				}
				if s.Start.Line != 1 || s.Start.Column != 1 {
					t.Errorf("start = %+v, want {1,1}", s.Start)
				}
			},
		},
		{
			name:      "two statements with semicolons",
			sql:       "SELECT 1; SELECT 2;",
			wantCount: 2,
			checks: func(t *testing.T, stmts []Statement) {
				if stmts[0].Text != "SELECT 1;" {
					t.Errorf("stmt[0].text = %q, want %q", stmts[0].Text, "SELECT 1;")
				}
				if stmts[1].Text != " SELECT 2;" {
					t.Errorf("stmt[1].text = %q, want %q", stmts[1].Text, " SELECT 2;")
				}
				if stmts[1].ByteStart != 9 {
					t.Errorf("stmt[1].ByteStart = %d, want 9", stmts[1].ByteStart)
				}
			},
		},
		{
			name: "multiline",
			sql:  "SELECT\n  1;\nSELECT\n  2;",
			wantCount: 2,
			checks: func(t *testing.T, stmts []Statement) {
				if stmts[0].Start.Line != 1 {
					t.Errorf("stmt[0] start line = %d, want 1", stmts[0].Start.Line)
				}
				if stmts[1].Start.Line != 3 {
					t.Errorf("stmt[1] start line = %d, want 3", stmts[1].Start.Line)
				}
			},
		},
		{
			name:      "CREATE TABLE",
			sql:       "CREATE TABLE t (id int);",
			wantCount: 1,
			checks: func(t *testing.T, stmts []Statement) {
				if _, ok := stmts[0].AST.(*ast.CreateStmt); !ok {
					t.Errorf("AST type = %T, want *ast.CreateStmt", stmts[0].AST)
				}
			},
		},
		{
			name:      "empty input",
			sql:       "",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.sql, err)
			}
			if len(stmts) != tt.wantCount {
				t.Fatalf("Parse(%q) got %d statements, want %d", tt.sql, len(stmts), tt.wantCount)
			}
			if tt.checks != nil {
				tt.checks(t, stmts)
			}
		})
	}
}

func TestParsePosition(t *testing.T) {
	sql := "SELECT 1;\nINSERT INTO t VALUES (1);\nDELETE FROM t;"
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 3 {
		t.Fatalf("got %d stmts, want 3", len(stmts))
	}

	// Line 1: "SELECT 1;"
	if stmts[0].Start != (Position{1, 1}) {
		t.Errorf("stmt[0].Start = %+v, want {1,1}", stmts[0].Start)
	}

	// Line 2: "INSERT INTO t VALUES (1);"
	if stmts[1].Start.Line != 2 {
		t.Errorf("stmt[1].Start.Line = %d, want 2", stmts[1].Start.Line)
	}

	// Line 3: "DELETE FROM t;"
	if stmts[2].Start.Line != 3 {
		t.Errorf("stmt[2].Start.Line = %d, want 3", stmts[2].Start.Line)
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse("SELECTT FROM")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
