package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocCreateTrigStmt(t *testing.T) {
	sql := "CREATE TRIGGER mytrig BEFORE INSERT ON t EXECUTE FUNCTION myfunc()"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateTrigStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateEventTrigStmt(t *testing.T) {
	sql := "CREATE EVENT TRIGGER myevt ON ddl_command_start EXECUTE FUNCTION myfunc()"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateEventTrigStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocTriggerTransition(t *testing.T) {
	sql := "CREATE TRIGGER mytrig AFTER INSERT ON t REFERENCING NEW TABLE AS newtab FOR EACH ROW EXECUTE FUNCTION myfunc()"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateTrigStmt)
	require.NotNil(t, stmt.TransitionRels)
	require.Len(t, stmt.TransitionRels.Items, 1)
	tt := stmt.TransitionRels.Items[0].(*nodes.TriggerTransition)
	got := sql[tt.Loc.Start:tt.Loc.End]
	assert.Equal(t, "NEW TABLE AS newtab", got)
}

func TestLocIndexStmt(t *testing.T) {
	sql := "CREATE INDEX myidx ON t (col)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.IndexStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocIndexElem(t *testing.T) {
	sql := "CREATE INDEX myidx ON t (col)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.IndexStmt)
	require.NotNil(t, stmt.IndexParams)
	require.Len(t, stmt.IndexParams.Items, 1)
	elem := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	got := sql[elem.Loc.Start:elem.Loc.End]
	assert.Equal(t, "col", got)
}

func TestLocViewStmt(t *testing.T) {
	sql := "CREATE VIEW myview AS SELECT * FROM t"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ViewStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateTableAsStmt(t *testing.T) {
	sql := "CREATE TABLE t AS SELECT 1"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateTableAsStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocRefreshMatViewStmt(t *testing.T) {
	sql := "REFRESH MATERIALIZED VIEW myview"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.RefreshMatViewStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateStmt(t *testing.T) {
	sql := "CREATE TABLE t (a int, b text)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}
