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

// TestParsePsqlClientMetacommands pins the 8 client-side psql backslash
// metacommand forms that PostgreSQL's server raw parser rejects but pg/parser
// intentionally accepts as comment trivia (lexer.go, internal/metacmd).
// Server-rejected non-metacommand forms (\;, \, \1, \\, :a, :'a', :"a") stay
// in testdata/compat.tsv under suite:psql so both server and omni agree.
func TestParsePsqlClientMetacommands(t *testing.T) {
	t.Run("standalone metacommands produce no statements", func(t *testing.T) {
		standalone := []struct {
			name  string
			input string
		}{
			{"describe", `\d`},
			{"describe table", `\dt a`},
			{"copy", `\copy a FROM STDIN`},
			{"restrict", `\restrict abc`},
			{"echo", `\echo a`},
		}
		for _, tt := range standalone {
			t.Run(tt.name, func(t *testing.T) {
				list, err := Parse(tt.input)
				if err != nil {
					t.Fatalf("Parse(%q) failed: %v", tt.input, err)
				}
				if list != nil && len(list.Items) != 0 {
					t.Fatalf("Parse(%q) produced %d items, want 0", tt.input, len(list.Items))
				}
			})
		}
	})

	assertSelectInt := func(t *testing.T, item nodes.Node, want int64) {
		t.Helper()
		raw, ok := item.(*nodes.RawStmt)
		if !ok {
			t.Fatalf("expected *nodes.RawStmt, got %T", item)
		}
		sel, ok := raw.Stmt.(*nodes.SelectStmt)
		if !ok {
			t.Fatalf("expected *nodes.SelectStmt, got %T", raw.Stmt)
		}
		if sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
			t.Fatalf("expected 1 target list item, got %v", sel.TargetList)
		}
		rt, ok := sel.TargetList.Items[0].(*nodes.ResTarget)
		if !ok {
			t.Fatalf("expected *nodes.ResTarget, got %T", sel.TargetList.Items[0])
		}
		ac, ok := rt.Val.(*nodes.A_Const)
		if !ok {
			t.Fatalf("expected *nodes.A_Const, got %T", rt.Val)
		}
		iv, ok := ac.Val.(*nodes.Integer)
		if !ok || iv.Ival != want {
			t.Fatalf("expected integer %d, got %#v", want, ac.Val)
		}
	}

	t.Run("trailing metacommands are skipped as line trivia", func(t *testing.T) {
		trailing := []struct {
			name  string
			input string
		}{
			{"gset", `SELECT 1 \gset`},
			{"g", `SELECT 1 \g`},
			{"gexec", `SELECT 1 \gexec`},
		}
		for _, tt := range trailing {
			t.Run(tt.name, func(t *testing.T) {
				list, err := Parse(tt.input)
				if err != nil {
					t.Fatalf("Parse(%q) failed: %v", tt.input, err)
				}
				if list == nil || len(list.Items) != 1 {
					t.Fatalf("Parse(%q) produced %d statements, want 1", tt.input, len(list.Items))
				}
				assertSelectInt(t, list.Items[0], 1)
			})
		}
	})

	t.Run("multi-statement boundaries around client metacommands", func(t *testing.T) {
		cases := []struct {
			name  string
			input string
			wants []int64
		}{
			{
				name:  "standalone metacommands between statements",
				input: "\\d\n\\dt a\nSELECT 1;\n\\copy a FROM STDIN\n\\restrict abc\n\\echo a\nSELECT 2;",
				wants: []int64{1, 2},
			},
			{
				name:  "trailing gset before semicolon and next statement",
				input: "SELECT 1 \\gset\n; SELECT 2;",
				wants: []int64{1, 2},
			},
			{
				name:  "trailing g before semicolon and next statement",
				input: "SELECT 1 \\g\n; SELECT 2;",
				wants: []int64{1, 2},
			},
			{
				name:  "trailing gexec before semicolon and next statement",
				input: "SELECT 1 \\gexec\n; SELECT 2;",
				wants: []int64{1, 2},
			},
		}
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				list, err := Parse(tt.input)
				if err != nil {
					t.Fatalf("Parse(%q) failed: %v", tt.input, err)
				}
				if list == nil || len(list.Items) != len(tt.wants) {
					got := 0
					if list != nil {
						got = len(list.Items)
					}
					t.Fatalf("Parse(%q) produced %d statements, want %d", tt.input, got, len(tt.wants))
				}
				for i, want := range tt.wants {
					assertSelectInt(t, list.Items[i], want)
				}
			})
		}

		// \g is skipped as a line comment and does not act as a statement
		// terminator on its own, so a second SELECT without a separating
		// semicolon fails to parse.
		if _, err := Parse("SELECT 1 \\g\nSELECT 2;"); err == nil {
			t.Fatal("expected parse error for missing semicolon after \\g, got nil")
		}
	})
}
