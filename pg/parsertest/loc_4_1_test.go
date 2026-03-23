package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// --- Section 4.1: Database & Schema Nodes ---

func TestLocCreatedbStmt(t *testing.T) {
	sql := "CREATE DATABASE mydb"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatedbStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreatedbStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreatedbStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreatedbStmtWithOptions(t *testing.T) {
	sql := "CREATE DATABASE mydb WITH OWNER myuser"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatedbStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreatedbStmt text = %q, want %q", got, sql)
	}
}

func TestLocDropdbStmt(t *testing.T) {
	sql := "DROP DATABASE mydb"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DropdbStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("DropdbStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("DropdbStmt text = %q, want %q", got, sql)
	}
}

func TestLocDropdbStmtIfExists(t *testing.T) {
	sql := "DROP DATABASE IF EXISTS mydb"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DropdbStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("DropdbStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterDatabaseStmt(t *testing.T) {
	sql := "ALTER DATABASE mydb CONNECTION LIMIT 10"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterDatabaseStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterDatabaseStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterDatabaseStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterDatabaseSetStmt(t *testing.T) {
	sql := "ALTER DATABASE mydb SET search_path TO public"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterDatabaseSetStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterDatabaseSetStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterDatabaseSetStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterDatabaseResetStmt(t *testing.T) {
	sql := "ALTER DATABASE mydb RESET search_path"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterDatabaseSetStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterDatabaseSetStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateSchemaStmt(t *testing.T) {
	sql := "CREATE SCHEMA myschema"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateSchemaStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CreateSchemaStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateSchemaStmt text = %q, want %q", got, sql)
	}
}

func TestLocCreateSchemaIfNotExists(t *testing.T) {
	sql := "CREATE SCHEMA IF NOT EXISTS myschema"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateSchemaStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CreateSchemaStmt text = %q, want %q", got, sql)
	}
}

func TestLocCommentStmt(t *testing.T) {
	sql := "COMMENT ON TABLE t IS 'a comment'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CommentStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CommentStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CommentStmt text = %q, want %q", got, sql)
	}
}

func TestLocCommentOnColumn(t *testing.T) {
	sql := "COMMENT ON COLUMN t.c IS 'column comment'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CommentStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CommentStmt text = %q, want %q", got, sql)
	}
}
