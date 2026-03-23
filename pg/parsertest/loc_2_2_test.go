package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

func TestLocCreateExtensionStmt(t *testing.T) {
	sql := "CREATE EXTENSION hstore"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateExtensionStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateExtensionStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateExtensionStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateAmStmt(t *testing.T) {
	sql := "CREATE ACCESS METHOD myam TYPE INDEX HANDLER myhandler"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateAmStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateAmStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateAmStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateCastStmt(t *testing.T) {
	sql := "CREATE CAST (int AS text) WITH FUNCTION int4_to_text(int)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateCastStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateCastStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateCastStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateTransformStmt(t *testing.T) {
	sql := "CREATE TRANSFORM FOR int LANGUAGE plpgsql (FROM SQL WITH FUNCTION int4_to_plpgsql(internal), TO SQL WITH FUNCTION plpgsql_to_int4(internal))"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateTransformStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateTransformStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateTransformStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterExtensionStmt(t *testing.T) {
	sql := "ALTER EXTENSION hstore UPDATE TO '2.0'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterExtensionStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterExtensionStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterExtensionStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterExtensionContentsStmt(t *testing.T) {
	sql := "ALTER EXTENSION hstore ADD TABLE mytable"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterExtensionContentsStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterExtensionContentsStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterExtensionContentsStmt text = %q, want %q", got, sql)
	}
}
