package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParsePsqlMetacommandLines pins parser-level handling of psql
// metacommand lines (\restrict / \unrestrict / \connect ...): they are not
// SQL, produce no AST node, and must not fail the parse — the backslash is
// not a legal SQL token, so skipping also discards the lexer error it
// produced. pg_dump plain format emits \restrict bracket lines since the
// CVE-2025-8714 point releases.
func TestParsePsqlMetacommandLines(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStmts int
	}{
		{"restrict header", "\\restrict Abc123\nSELECT 1;", 1},
		{"unrestrict trailer", "SELECT 1;\n\\unrestrict Abc123\n", 1},
		{"consecutive metacommands", "\\restrict A\n\\connect mydb --create\nSELECT 1;", 1},
		{"between statements", "SELECT 1;\n\\connect other\nSELECT 2;", 2},
		{"metacommand only", "\\restrict Abc123\n", 0},
		{"copy data terminator unaffected", "COPY t FROM stdin;\n\\N\n\\.\nSELECT 2;", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			n := 0
			if list != nil {
				for _, item := range list.Items {
					if _, ok := item.(*nodes.RawStmt); ok {
						n++
					}
				}
			}
			if n != tt.wantStmts {
				t.Fatalf("got %d statements, want %d", n, tt.wantStmts)
			}
		})
	}

	t.Run("backslash inside string stays literal", func(t *testing.T) {
		list, err := Parse("SELECT E'a\n\\restrict b';")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(list.Items))
		}
	})
}
