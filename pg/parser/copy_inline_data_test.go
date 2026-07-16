package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseCopyInlineData pins parser-level handling of COPY ... FROM STDIN
// inline data (psql / pg_dump plain format): the data lines are not SQL and
// must be captured on the CopyStmt, with parsing resuming cleanly after the
// "\." terminator. Before this fix the parser reported a syntax error on
// the first data line, so the pg_dump shape failed the whole pipeline even
// though the splitter kept the block intact.
func TestParseCopyInlineData(t *testing.T) {
	t.Run("pg_dump shape with following statement", func(t *testing.T) {
		sql := "COPY t (a, b) FROM stdin;\nx;y\t1\n\\N\t2\n\\.\nSELECT 1;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
		raw := list.Items[0].(*nodes.RawStmt)
		cs, ok := raw.Stmt.(*nodes.CopyStmt)
		if !ok {
			t.Fatalf("expected CopyStmt, got %T", raw.Stmt)
		}
		wantData := "\nx;y\t1\n\\N\t2\n\\.\n"
		if cs.InlineData != wantData {
			t.Fatalf("InlineData = %q, want %q", cs.InlineData, wantData)
		}
		if got := sql[raw.Loc.Start:raw.Loc.End]; !strings.HasSuffix(got, "\\.\n") {
			t.Fatalf("RawStmt span %q does not include terminator", got)
		}
		if _, ok := list.Items[1].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt); !ok {
			t.Fatalf("statement after data did not parse as SELECT")
		}
	})

	t.Run("backslash data does not poison lexer", func(t *testing.T) {
		// \N is the pg_dump NULL marker; bare backslash is not legal SQL, so
		// any lookahead into the data region must be discarded.
		sql := "COPY t FROM stdin;\n\\N\n\\.\nSELECT 2;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
	})

	t.Run("two copy blocks", func(t *testing.T) {
		sql := "COPY a FROM stdin;\n1;\n\\.\nCOPY b FROM stdin;\n2;\n\\.\n"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
		for i, item := range list.Items {
			cs := item.(*nodes.RawStmt).Stmt.(*nodes.CopyStmt)
			if cs.InlineData == "" {
				t.Fatalf("copy %d has empty InlineData", i)
			}
		}
	})

	t.Run("missing terminator swallows to EOF like psql", func(t *testing.T) {
		sql := "COPY t FROM stdin;\na\t1\nSELECT 3;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(list.Items))
		}
		cs := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CopyStmt)
		if !strings.Contains(cs.InlineData, "SELECT 3;") {
			t.Fatalf("InlineData %q should contain the swallowed line", cs.InlineData)
		}
	})

	t.Run("no data when copy is last without newline", func(t *testing.T) {
		sql := "COPY t FROM stdin;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		cs := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CopyStmt)
		if cs.InlineData != "" {
			t.Fatalf("InlineData = %q, want empty", cs.InlineData)
		}
	})

	t.Run("copy from file does not consume following statements", func(t *testing.T) {
		sql := "COPY t FROM 'x.csv';\nSELECT 4;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
	})

	t.Run("copy to stdout does not consume following statements", func(t *testing.T) {
		sql := "COPY t TO stdout;\nSELECT 5;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
	})

	t.Run("copy with options captures data", func(t *testing.T) {
		sql := "copy t from stdin with (format text);\na;b\t1\n\\.\nSELECT 6;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 2 {
			t.Fatalf("expected 2 statements, got %d", len(list.Items))
		}
	})
}
