package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseForeignTableColumnOptions(t *testing.T) {
	stmts, err := Parse(`CREATE FOREIGN TABLE ft1 (
  c1 integer OPTIONS ("param 1" 'val1') NOT NULL,
  c2 text OPTIONS (param2 'val2', param3 'val3') CHECK (c2 <> '')
) SERVER s0 OPTIONS (delimiter ',', quote '"')`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	stmt := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateForeignTableStmt)
	col := stmt.Base.TableElts.Items[0].(*nodes.ColumnDef)
	if col.Fdwoptions == nil || len(col.Fdwoptions.Items) != 1 {
		t.Fatalf("expected one FDW option on c1, got %#v", col.Fdwoptions)
	}
	if col.Constraints == nil || len(col.Constraints.Items) == 0 {
		t.Fatalf("expected NOT NULL constraint after column OPTIONS")
	}
}

func TestParseAlterForeignTableAddColumnOptions(t *testing.T) {
	stmts, err := Parse(`ALTER FOREIGN TABLE ft1 ADD COLUMN c10 integer OPTIONS (p1 'v1')`)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	stmt := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	col := cmd.Def.(*nodes.ColumnDef)
	if col.Fdwoptions == nil || len(col.Fdwoptions.Items) != 1 {
		t.Fatalf("expected one FDW option on added column, got %#v", col.Fdwoptions)
	}
}
