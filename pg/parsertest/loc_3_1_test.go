package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

func TestLocAlterCollationStmt(t *testing.T) {
	sql := `ALTER COLLATION "C" REFRESH VERSION`
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterCollationStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterCollationStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterCollationStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterDomainStmt(t *testing.T) {
	sql := "ALTER DOMAIN mydom SET NOT NULL"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterDomainStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterDomainStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterDomainStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterEnumStmt(t *testing.T) {
	sql := "ALTER TYPE mood ADD VALUE 'great'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterEnumStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterEnumStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterEnumStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterEventTrigStmt(t *testing.T) {
	sql := "ALTER EVENT TRIGGER mytrig DISABLE"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterEventTrigStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterEventTrigStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterEventTrigStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterFunctionStmt(t *testing.T) {
	sql := "ALTER FUNCTION myfunc(int) STABLE"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterFunctionStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterFunctionStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterFunctionStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterObjectDependsStmt(t *testing.T) {
	sql := "ALTER FUNCTION myfunc DEPENDS ON EXTENSION hstore"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterObjectDependsStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterObjectDependsStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterObjectDependsStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterObjectSchemaStmt(t *testing.T) {
	sql := "ALTER TABLE t SET SCHEMA newschema"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterObjectSchemaStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterObjectSchemaStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterObjectSchemaStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterOwnerStmt(t *testing.T) {
	// ALTER TABLE ... OWNER TO produces AlterTableStmt with AT_ChangeOwner cmd,
	// not AlterOwnerStmt. Use ALTER SCHEMA which produces AlterOwnerStmt.
	sql := "ALTER SCHEMA myschema OWNER TO newowner"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterOwnerStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterOwnerStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterOwnerStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterTableSpaceOptionsStmt(t *testing.T) {
	sql := "ALTER TABLESPACE myts SET (seq_page_cost=2)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTableSpaceOptionsStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterTableSpaceOptionsStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterTableSpaceOptionsStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterTSConfigurationStmt(t *testing.T) {
	sql := "ALTER TEXT SEARCH CONFIGURATION myconfig ADD MAPPING FOR word WITH simple"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTSConfigurationStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterTSConfigurationStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterTSConfigurationStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterTSDictionaryStmt(t *testing.T) {
	sql := "ALTER TEXT SEARCH DICTIONARY mydict (STOPWORDS = 'english')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTSDictionaryStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("AlterTSDictionaryStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("AlterTSDictionaryStmt text = %q, want %q", got, sql)
	}
}

func TestLocAlterTypeStmt(t *testing.T) {
	sql := "ALTER TYPE comptype ADD ATTRIBUTE c int"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	// ALTER TYPE ... ADD ATTRIBUTE produces AlterTableStmt (composite type path)
	// so verify it parses and the AlterTableStmt is returned
	_, ok := raw.Stmt.(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt for composite type ALTER, got %T", raw.Stmt)
	}
}
