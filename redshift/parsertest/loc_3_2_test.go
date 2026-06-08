package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocAlterTableStmt(t *testing.T) {
	sql := "ALTER TABLE t ADD COLUMN x int"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTableStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterTableCmd(t *testing.T) {
	sql := "ALTER TABLE t ADD COLUMN x int"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTableStmt)
	require.NotNil(t, stmt.Cmds)
	require.Len(t, stmt.Cmds.Items, 1)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	got := sql[cmd.Loc.Start:cmd.Loc.End]
	assert.Equal(t, "ADD COLUMN x int", got)
}

func TestLocAlterSeqStmt(t *testing.T) {
	sql := "ALTER SEQUENCE myseq RESTART WITH 1"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterSeqStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterTableMoveAllStmt(t *testing.T) {
	sql := "ALTER TABLE ALL IN TABLESPACE ts1 SET TABLESPACE ts2"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTableMoveAllStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocPartitionCmd(t *testing.T) {
	sql := "ALTER TABLE t ATTACH PARTITION p FOR VALUES FROM (1) TO (10)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterTableStmt)
	require.NotNil(t, stmt.Cmds)
	require.Len(t, stmt.Cmds.Items, 1)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	pcmd := cmd.Def.(*nodes.PartitionCmd)
	got := sql[pcmd.Loc.Start:pcmd.Loc.End]
	assert.Equal(t, "ATTACH PARTITION p FOR VALUES FROM (1) TO (10)", got)
}

func TestLocRenameStmt(t *testing.T) {
	sql := "ALTER TABLE t RENAME COLUMN old TO new"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.RenameStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}
