package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func firstAlterCmd(t *testing.T, stmt nodes.Node) *nodes.AlterTableCmd {
	t.Helper()
	at, ok := stmt.(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", stmt)
	}
	if at.Cmds == nil || len(at.Cmds.Items) == 0 {
		t.Fatalf("no cmds")
	}
	cmd, ok := at.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", at.Cmds.Items[0])
	}
	return cmd
}

func TestAlterTableSetCompressionDefault(t *testing.T) {
	t.Run("SET COMPRESSION default", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER TABLE t ALTER COLUMN c SET COMPRESSION default")
		cmd := firstAlterCmd(t, stmt)
		if cmd.Subtype != int(nodes.AT_SetCompression) {
			t.Fatalf("expected AT_SetCompression, got %v", cmd.Subtype)
		}
		s, ok := cmd.Def.(*nodes.String)
		if !ok || s.Str != "default" {
			t.Fatalf("expected String{default}, got %+v", cmd.Def)
		}
	})

	t.Run("SET COMPRESSION pglz (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER TABLE t ALTER COLUMN c SET COMPRESSION pglz")
		cmd := firstAlterCmd(t, stmt)
		s, ok := cmd.Def.(*nodes.String)
		if !ok || s.Str != "pglz" {
			t.Fatalf("expected String{pglz}, got %+v", cmd.Def)
		}
	})

	t.Run("SET COMPRESSION lz4 (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER TABLE t ALTER COLUMN c SET COMPRESSION lz4")
		cmd := firstAlterCmd(t, stmt)
		s, ok := cmd.Def.(*nodes.String)
		if !ok || s.Str != "lz4" {
			t.Fatalf("expected String{lz4}, got %+v", cmd.Def)
		}
	})
}
