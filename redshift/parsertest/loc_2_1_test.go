package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Section 2.1: Type & operator definitions ---

func TestLocDefineStmtType(t *testing.T) {
	sql := "CREATE TYPE mytype AS (a int, b text)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CompositeTypeStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocDefineStmtOperator(t *testing.T) {
	sql := "CREATE OPERATOR === (LEFTARG=int, RIGHTARG=int, FUNCTION=int4eq)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DefineStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocDefineStmtAggregate(t *testing.T) {
	sql := "CREATE AGGREGATE myagg(int) (SFUNC=int4pl, STYPE=int)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DefineStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCompositeTypeStmt(t *testing.T) {
	sql := "CREATE TYPE comptype AS (x int, y int)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CompositeTypeStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateEnumStmt(t *testing.T) {
	sql := "CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateEnumStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateRangeStmt(t *testing.T) {
	sql := "CREATE TYPE floatrange AS RANGE (SUBTYPE = float8)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateRangeStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateOpClassStmt(t *testing.T) {
	sql := "CREATE OPERATOR CLASS myclass FOR TYPE int4 USING btree AS OPERATOR 1 <"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateOpClassStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateOpFamilyStmt(t *testing.T) {
	sql := "CREATE OPERATOR FAMILY myfam USING btree"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateOpFamilyStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateOpClassItem(t *testing.T) {
	sql := "CREATE OPERATOR CLASS myclass FOR TYPE int4 USING btree AS OPERATOR 1 <"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateOpClassStmt)
	require.NotNil(t, stmt.Items)
	require.True(t, len(stmt.Items.Items) > 0)
	item := stmt.Items.Items[0].(*nodes.CreateOpClassItem)
	got := sql[item.Loc.Start:item.Loc.End]
	assert.Equal(t, "OPERATOR 1 <", got)
}

func TestLocCreateConversionStmt(t *testing.T) {
	sql := "CREATE CONVERSION myconv FOR 'UTF8' TO 'LATIN1' FROM utf8_to_iso8859_1"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateConversionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateStatsStmt(t *testing.T) {
	sql := "CREATE STATISTICS mystats (dependencies) ON a, b FROM t"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateStatsStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterDefaultPrivilegesStmt(t *testing.T) {
	sql := "ALTER DEFAULT PRIVILEGES GRANT SELECT ON TABLES TO public"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterDefaultPrivilegesStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterOpFamilyStmt(t *testing.T) {
	sql := "ALTER OPERATOR FAMILY myfam USING btree ADD OPERATOR 1 <(int4,int4)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterOpFamilyStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterOperatorStmt(t *testing.T) {
	sql := "ALTER OPERATOR ===(int4,int4) SET (RESTRICT=eqsel)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterOperatorStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterStatsStmt(t *testing.T) {
	sql := "ALTER STATISTICS mystats SET STATISTICS 100"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterStatsStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocStatsElem(t *testing.T) {
	sql := "CREATE STATISTICS mystats (dependencies) ON a, b FROM t"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateStatsStmt)
	require.NotNil(t, stmt.Exprs)
	require.True(t, len(stmt.Exprs.Items) >= 2)
	elem0 := stmt.Exprs.Items[0].(*nodes.StatsElem)
	got0 := sql[elem0.Loc.Start:elem0.Loc.End]
	assert.Equal(t, "a", got0)
	elem1 := stmt.Exprs.Items[1].(*nodes.StatsElem)
	got1 := sql[elem1.Loc.Start:elem1.Loc.End]
	assert.Equal(t, "b", got1)
}
