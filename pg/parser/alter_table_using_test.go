package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestAlterColumnTypeUsing holds ALTER COLUMN TYPE's USING expression in the
// ColumnDef's RawDefault, where gram.y's alter_table_cmd puts it.
func TestAlterColumnTypeUsing(t *testing.T) {
	const sql = "ALTER TABLE t ALTER COLUMN a TYPE mood USING a::mood"
	cmd := firstAlterCmd(t, singleStmt(t, sql))
	if cmd.Subtype != int(nodes.AT_AlterColumnType) {
		t.Fatalf("expected AT_AlterColumnType, got %v", cmd.Subtype)
	}
	cd, ok := cmd.Def.(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected ColumnDef, got %T", cmd.Def)
	}
	tc, ok := cd.RawDefault.(*nodes.TypeCast)
	if !ok {
		t.Fatalf("expected the USING TypeCast in RawDefault, got %T", cd.RawDefault)
	}
	if got := sql[tc.Loc.Start:tc.Loc.End]; got != "a::mood" {
		t.Errorf("USING expression spans %q, want %q", got, "a::mood")
	}

	cmd = firstAlterCmd(t, singleStmt(t, "ALTER TABLE t ALTER COLUMN a TYPE mood"))
	if cd := cmd.Def.(*nodes.ColumnDef); cd.RawDefault != nil {
		t.Errorf("without USING, RawDefault is %T, want nil", cd.RawDefault)
	}
}
