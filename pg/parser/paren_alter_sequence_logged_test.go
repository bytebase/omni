package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// KB-2b: PG 15+ added SET LOGGED / SET UNLOGGED to the ALTER SEQUENCE action
// list (SeqOptElem). Prior to this fix, parseAlterSequence stopped at an
// earlier match and left `SET LOGGED` / `SET UNLOGGED` as residual tail,
// which (combined with Parse()'s silent multi-stmt fallback) caused 4
// pgregress statements to silently split. These tests pin the parser to
// accept both forms and emit an AlterSeqStmt with a single DefElem option
// mirroring the bare LOGGED/UNLOGGED SeqOptElem AST shape.

func firstSeqOptDefElem(t *testing.T, stmt nodes.Node) *nodes.DefElem {
	t.Helper()
	as, ok := stmt.(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected AlterSeqStmt, got %T", stmt)
	}
	if as.Options == nil || len(as.Options.Items) == 0 {
		t.Fatalf("AlterSeqStmt.Options empty")
	}
	de, ok := as.Options.Items[0].(*nodes.DefElem)
	if !ok {
		t.Fatalf("expected DefElem, got %T", as.Options.Items[0])
	}
	return de
}

func TestAlterSequenceSetLogged(t *testing.T) {
	stmt := singleStmt(t, "ALTER SEQUENCE foo SET LOGGED")
	as, ok := stmt.(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected AlterSeqStmt, got %T", stmt)
	}
	if as.Sequence == nil || as.Sequence.Relname != "foo" {
		t.Fatalf("expected sequence name 'foo', got %+v", as.Sequence)
	}
	de := firstSeqOptDefElem(t, stmt)
	if de.Defname != "logged" {
		t.Fatalf("expected Defname 'logged', got %q", de.Defname)
	}
	b, ok := de.Arg.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected Boolean arg, got %T", de.Arg)
	}
	if !b.Boolval {
		t.Fatalf("expected Boolean{true} for SET LOGGED, got %+v", b)
	}
}

func TestAlterSequenceSetUnlogged(t *testing.T) {
	stmt := singleStmt(t, "ALTER SEQUENCE foo SET UNLOGGED")
	as, ok := stmt.(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected AlterSeqStmt, got %T", stmt)
	}
	if as.Sequence == nil || as.Sequence.Relname != "foo" {
		t.Fatalf("expected sequence name 'foo', got %+v", as.Sequence)
	}
	de := firstSeqOptDefElem(t, stmt)
	if de.Defname != "logged" {
		t.Fatalf("expected Defname 'logged', got %q", de.Defname)
	}
	b, ok := de.Arg.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected Boolean arg, got %T", de.Arg)
	}
	if b.Boolval {
		t.Fatalf("expected Boolean{false} for SET UNLOGGED, got %+v", b)
	}
}

// Guard that we didn't break the existing SET SCHEMA path.
func TestAlterSequenceSetSchemaStillWorks(t *testing.T) {
	stmt := singleStmt(t, "ALTER SEQUENCE foo SET SCHEMA myschema")
	if _, ok := stmt.(*nodes.AlterObjectSchemaStmt); !ok {
		t.Fatalf("expected AlterObjectSchemaStmt, got %T", stmt)
	}
}

// Guard: IF EXISTS + SET LOGGED/UNLOGGED propagates MissingOk.
func TestAlterSequenceIfExistsSetLogged(t *testing.T) {
	stmt := singleStmt(t, "ALTER SEQUENCE IF EXISTS foo SET UNLOGGED")
	as, ok := stmt.(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected AlterSeqStmt, got %T", stmt)
	}
	if !as.MissingOk {
		t.Fatalf("expected MissingOk=true for IF EXISTS variant")
	}
	de := firstSeqOptDefElem(t, stmt)
	b, ok := de.Arg.(*nodes.Boolean)
	if !ok || b.Boolval {
		t.Fatalf("expected Boolean{false} for SET UNLOGGED, got %+v", de.Arg)
	}
}

// Sanity: verify the statement list has exactly one entry — i.e. the parser
// consumed the entire SET LOGGED / SET UNLOGGED phrase rather than leaving
// residual tokens that Parse()'s multi-stmt fallback would split into a
// second (bogus) statement. This is the direct failure mode KB-2b caused
// upstream in the pgregress harness.
func TestAlterSequenceSetLoggedSingleStatement(t *testing.T) {
	for _, sql := range []string{
		"ALTER SEQUENCE foo SET LOGGED",
		"ALTER SEQUENCE foo SET UNLOGGED",
	} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if list == nil || len(list.Items) != 1 {
			t.Fatalf("Parse(%q): expected 1 stmt, got %d", sql, len(list.Items))
		}
	}
}
