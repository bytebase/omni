package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// ---------------------------------------------------------------------------
// Section 4.2: Sequence, Function & Domain Nodes
// ---------------------------------------------------------------------------

func TestLocCreateSeqStmt(t *testing.T) {
	sql := "CREATE SEQUENCE myseq START 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateSeqStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateSeqStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateSeqStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateSeqStmtIfNotExists(t *testing.T) {
	sql := "CREATE SEQUENCE IF NOT EXISTS myseq"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateSeqStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateSeqStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateDomainStmt(t *testing.T) {
	sql := "CREATE DOMAIN posint AS int CHECK (VALUE > 0)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateDomainStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateDomainStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateDomainStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateFunctionStmt(t *testing.T) {
	sql := "CREATE FUNCTION myfunc(int) RETURNS int AS $$ SELECT $1 $$ LANGUAGE sql"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateFunctionStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateFunctionStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateFunctionStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateFunctionStmtOrReplace(t *testing.T) {
	sql := "CREATE OR REPLACE FUNCTION myfunc() RETURNS void AS $$ SELECT 1 $$ LANGUAGE sql"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateFunctionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateFunctionStmt text = %q, want %q", got, sql)
	}
}

func TestLocFunctionParameter(t *testing.T) {
	sql := "CREATE FUNCTION myfunc(x int, y text) RETURNS int AS $$ SELECT $1 $$ LANGUAGE sql"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateFunctionStmt)
	if stmt.Parameters == nil || len(stmt.Parameters.Items) < 2 {
		t.Fatal("expected at least 2 parameters")
	}
	// First parameter: "x int"
	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Loc.Start == -1 || p1.Loc.End == -1 {
		t.Fatalf("FunctionParameter Loc not set: %+v", p1.Loc)
	}
	got1 := sql[p1.Loc.Start:p1.Loc.End]
	if got1 != "x int" {
		t.Errorf("param 1 text = %q, want %q", got1, "x int")
	}
	// Second parameter: "y text"
	p2 := stmt.Parameters.Items[1].(*nodes.FunctionParameter)
	got2 := sql[p2.Loc.Start:p2.Loc.End]
	if got2 != "y text" {
		t.Errorf("param 2 text = %q, want %q", got2, "y text")
	}
}

func TestLocReturnStmt(t *testing.T) {
	sql := "CREATE FUNCTION myfunc() RETURNS int RETURN 42"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateFunctionStmt)
	if stmt.SqlBody == nil {
		t.Fatal("expected SqlBody to be set")
	}
	ret := stmt.SqlBody.(*nodes.ReturnStmt)
	if ret.Loc.Start == -1 || ret.Loc.End == -1 {
		t.Fatalf("ReturnStmt Loc not set: %+v", ret.Loc)
	}
	got := sql[ret.Loc.Start:ret.Loc.End]
	want := "RETURN 42"
	if got != want {
		t.Errorf("ReturnStmt text = %q, want %q", got, want)
	}
}
