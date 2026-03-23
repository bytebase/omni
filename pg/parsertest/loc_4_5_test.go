package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocDeclareCursorStmt(t *testing.T) {
	sql := "DECLARE mycur CURSOR FOR SELECT 1"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DeclareCursorStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocFetchStmt(t *testing.T) {
	sql := "FETCH NEXT FROM mycur"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.FetchStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocClosePortalStmt(t *testing.T) {
	sql := "CLOSE mycur"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ClosePortalStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocPrepareStmt(t *testing.T) {
	sql := "PREPARE myplan AS SELECT 1"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.PrepareStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocExecuteStmt(t *testing.T) {
	sql := "EXECUTE myplan"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ExecuteStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCopyStmt(t *testing.T) {
	sql := "COPY t FROM STDIN"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CopyStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocLockStmt(t *testing.T) {
	sql := "LOCK TABLE t IN ACCESS EXCLUSIVE MODE"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.LockStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}
