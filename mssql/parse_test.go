package mssql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

func TestBuildLineIndex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  lineIndex
	}{
		{"empty string", "", lineIndex{0}},
		{"no newlines", "SELECT 1", lineIndex{0}},
		{"one newline", "SELECT\n1", lineIndex{0, 7}},
		{"multiple newlines", "a\nb\nc", lineIndex{0, 2, 4}},
		{"trailing newline", "SELECT 1\n", lineIndex{0, 9}},
		{"consecutive newlines", "a\n\nb", lineIndex{0, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildLineIndex(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("buildLineIndex(%q) = %v (len %d), want %v (len %d)", tt.input, got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("buildLineIndex(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		stmts, err := Parse("")
		if err != nil {
			t.Fatalf("Parse('') error: %v", err)
		}
		if len(stmts) != 0 {
			t.Errorf("Parse('') returned %d statements, want 0", len(stmts))
		}
	})

	t.Run("single statement", func(t *testing.T) {
		stmts, err := Parse("SELECT 1")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		s := stmts[0]
		if s.AST == nil {
			t.Error("AST is nil")
		}
		if s.Start.Line != 1 || s.Start.Column != 1 {
			t.Errorf("Start = %+v, want {1,1}", s.Start)
		}
	})

	t.Run("single statement with semicolon", func(t *testing.T) {
		stmts, err := Parse("SELECT 1;")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		if !strings.Contains(stmts[0].Text, ";") {
			t.Errorf("Text = %q, want to contain semicolon", stmts[0].Text)
		}
	})

	t.Run("multi-statement", func(t *testing.T) {
		stmts, err := Parse("SELECT 1; SELECT 2; SELECT 3")
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 3 {
			t.Fatalf("got %d statements, want 3", len(stmts))
		}
		for i, s := range stmts {
			if s.AST == nil {
				t.Errorf("stmt[%d].AST is nil", i)
			}
		}
	})

	t.Run("multi-line positions", func(t *testing.T) {
		sql := "SELECT\n1"
		stmts, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		if len(stmts) != 1 {
			t.Fatalf("got %d statements, want 1", len(stmts))
		}
		s := stmts[0]
		if s.Start.Line != 1 {
			t.Errorf("Start.Line = %d, want 1", s.Start.Line)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		_, err := Parse("SELECT FROM WHERE")
		// We just check it doesn't panic — error behavior varies
		_ = err
	})
}

func TestOffsetToPosition(t *testing.T) {
	// "SELECT\n1\nFROM t"
	// Line 1: "SELECT\n"  bytes 0-6 (newline at 6)
	// Line 2: "1\n"       bytes 7-8 (newline at 8)
	// Line 3: "FROM t"    bytes 9-14
	idx := buildLineIndex("SELECT\n1\nFROM t")

	tests := []struct {
		name   string
		offset int
		want   Position
	}{
		{"offset 0 = line 1 col 1", 0, Position{Line: 1, Column: 1}},
		{"S in SELECT", 0, Position{Line: 1, Column: 1}},
		{"T in SELECT", 5, Position{Line: 1, Column: 6}},
		{"start of line 2", 7, Position{Line: 2, Column: 1}},
		{"start of line 3", 9, Position{Line: 3, Column: 1}},
		{"mid line 3", 12, Position{Line: 3, Column: 4}},
		{"past end", 15, Position{Line: 3, Column: 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := offsetToPosition(idx, tt.offset)
			if got != tt.want {
				t.Errorf("offsetToPosition(idx, %d) = %+v, want %+v", tt.offset, got, tt.want)
			}
		})
	}
}

// TestParseWithClauseDML is the end-to-end shape of BYT-10223: a CTE followed
// by INSERT must come back as one statement whose text spans the whole input
// and whose AST can be walked. Before the fix Parse returned a typed-nil
// *ast.SelectStmt holding all the text plus an InsertStmt with none, and
// ast.Inspect on the first one panicked inside every SQL review rule.
func TestParseWithClauseDML(t *testing.T) {
	sql := `-- txn-mode = off

;WITH TargetVenues AS (
    SELECT v.VenueId
    FROM Venue v
    WHERE v.IsActive = 1
)
INSERT INTO VenueMeta (VenueId, Attribute)
SELECT tv.VenueId, 'Roller.Venue.Cap'
FROM TargetVenues tv
WHERE NOT EXISTS (
    SELECT 1 FROM VenueMeta vm WHERE vm.VenueId = tv.VenueId
);`

	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("got %d statements, want 1", len(stmts))
	}
	s := stmts[0]
	ins, ok := s.AST.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("AST is %T, want *ast.InsertStmt", s.AST)
	}
	if ins.WithClause == nil || ins.WithClause.CTEs == nil || len(ins.WithClause.CTEs.Items) != 1 {
		t.Fatalf("WithClause = %+v, want one CTE", ins.WithClause)
	}
	if s.Text != sql {
		t.Errorf("Text = %q, want the whole input", s.Text)
	}
	if s.ByteStart != 0 || s.ByteEnd != len(sql) {
		t.Errorf("ByteStart/ByteEnd = %d/%d, want 0/%d", s.ByteStart, s.ByteEnd, len(sql))
	}
	// Start points at the WITH keyword (line 3, after the leading semicolon).
	if s.Start.Line != 3 || s.Start.Column != 2 {
		t.Errorf("Start = %d:%d, want 3:2", s.Start.Line, s.Start.Column)
	}

	// Every advisor rule walks the AST; this must not panic and must reach
	// both the CTE body and the INSERT source.
	var seen []string
	ast.Inspect(s.AST, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.CommonTableExpr, *ast.SelectStmt, *ast.InsertStmt:
			seen = append(seen, fmt.Sprintf("%T", n))
		}
		return true
	})
	joined := strings.Join(seen, " ")
	for _, want := range []string{"*ast.InsertStmt", "*ast.CommonTableExpr", "*ast.SelectStmt"} {
		if !strings.Contains(joined, want) {
			t.Errorf("Inspect did not visit %s; visited: %s", want, joined)
		}
	}
}
