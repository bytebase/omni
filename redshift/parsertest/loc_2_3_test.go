package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
)

// ---------------------------------------------------------------------------
// Section 2.3: Foreign data wrapper nodes
// ---------------------------------------------------------------------------

func TestLocCreateFdwStmt(t *testing.T) {
	sql := "CREATE FOREIGN DATA WRAPPER myfdw"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.CreateFdwStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "CREATE FOREIGN DATA WRAPPER myfdw"
	if got != want {
		t.Errorf("CreateFdwStmt text = %q, want %q", got, want)
	}
}

func TestLocCreateForeignServerStmt(t *testing.T) {
	sql := "CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.CreateForeignServerStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw"
	if got != want {
		t.Errorf("CreateForeignServerStmt text = %q, want %q", got, want)
	}
}

func TestLocCreateForeignTableStmt(t *testing.T) {
	sql := "CREATE FOREIGN TABLE ft (a int) SERVER myserver"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.CreateForeignTableStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "CREATE FOREIGN TABLE ft (a int) SERVER myserver"
	if got != want {
		t.Errorf("CreateForeignTableStmt text = %q, want %q", got, want)
	}
}

func TestLocCreatePLangStmt(t *testing.T) {
	sql := "CREATE LANGUAGE plmylang"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.CreatePLangStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "CREATE LANGUAGE plmylang"
	if got != want {
		t.Errorf("CreatePLangStmt text = %q, want %q", got, want)
	}
}

func TestLocCreateUserMappingStmt(t *testing.T) {
	sql := "CREATE USER MAPPING FOR current_user SERVER myserver"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.CreateUserMappingStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "CREATE USER MAPPING FOR current_user SERVER myserver"
	if got != want {
		t.Errorf("CreateUserMappingStmt text = %q, want %q", got, want)
	}
}

func TestLocAlterFdwStmt(t *testing.T) {
	sql := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD host 'foo')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.AlterFdwStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD host 'foo')"
	if got != want {
		t.Errorf("AlterFdwStmt text = %q, want %q", got, want)
	}
}

func TestLocAlterForeignServerStmt(t *testing.T) {
	sql := "ALTER SERVER myserver OPTIONS (SET host 'bar')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.AlterForeignServerStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "ALTER SERVER myserver OPTIONS (SET host 'bar')"
	if got != want {
		t.Errorf("AlterForeignServerStmt text = %q, want %q", got, want)
	}
}

func TestLocAlterUserMappingStmt(t *testing.T) {
	sql := "ALTER USER MAPPING FOR current_user SERVER myserver OPTIONS (SET password 'x')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.AlterUserMappingStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "ALTER USER MAPPING FOR current_user SERVER myserver OPTIONS (SET password 'x')"
	if got != want {
		t.Errorf("AlterUserMappingStmt text = %q, want %q", got, want)
	}
}

func TestLocDropUserMappingStmt(t *testing.T) {
	sql := "DROP USER MAPPING FOR current_user SERVER myserver"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.DropUserMappingStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "DROP USER MAPPING FOR current_user SERVER myserver"
	if got != want {
		t.Errorf("DropUserMappingStmt text = %q, want %q", got, want)
	}
}

func TestLocImportForeignSchemaStmt(t *testing.T) {
	sql := "IMPORT FOREIGN SCHEMA public FROM SERVER myserver INTO local_schema"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	n := raw.Stmt.(*nodes.ImportForeignSchemaStmt)
	got := sql[n.Loc.Start:n.Loc.End]
	want := "IMPORT FOREIGN SCHEMA public FROM SERVER myserver INTO local_schema"
	if got != want {
		t.Errorf("ImportForeignSchemaStmt text = %q, want %q", got, want)
	}
}
