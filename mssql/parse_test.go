package mssql

import (
	"fmt"
	"github.com/bytebase/omni/review"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

func TestPositionAt(t *testing.T) {
	// "SELECT\n1\nFROM t"
	// Line 1: "SELECT\n"  bytes 0-6 (newline at 6)
	// Line 2: "1\n"       bytes 7-8 (newline at 8)
	// Line 3: "FROM t"    bytes 9-14
	idx := review.Index("SELECT\n1\nFROM t")

	tests := []struct {
		name   string
		offset int
		want   Position
	}{
		{"offset 0 = line 1 col 1", 0, Position{Line: 1, Column: 1}},
		{"T in SELECT", 5, Position{Line: 1, Column: 6}},
		{"start of line 2", 7, Position{Line: 2, Column: 1}},
		{"start of line 3", 9, Position{Line: 3, Column: 1}},
		{"mid line 3", 12, Position{Line: 3, Column: 4}},
		{"past end", 15, Position{Line: 3, Column: 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := positionAt(idx, tt.offset); got != tt.want {
				t.Errorf("positionAt(idx, %d) = %+v, want %+v", tt.offset, got, tt.want)
			}
		})
	}

	// Columns count code points, not bytes: 'é' is two bytes and one column.
	multi := review.Index("SELECT 'é', 1")
	if got := positionAt(multi, len("SELECT 'é', ")); got != (Position{Line: 1, Column: 13}) {
		t.Errorf("positionAt after a two-byte rune = %+v, want line 1 column 13", got)
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
