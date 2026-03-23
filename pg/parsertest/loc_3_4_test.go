package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Section 3.4: Publication & subscription nodes ---

func TestLocCreatePublicationStmt(t *testing.T) {
	sql := "CREATE PUBLICATION mypub FOR ALL TABLES"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatePublicationStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreatePublicationStmtWithTable(t *testing.T) {
	sql := "CREATE PUBLICATION mypub FOR TABLE t1, t2"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatePublicationStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreatePublicationStmtNameOnly(t *testing.T) {
	sql := "CREATE PUBLICATION mypub"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatePublicationStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterPublicationStmt(t *testing.T) {
	sql := "ALTER PUBLICATION mypub ADD TABLE t"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterPublicationStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterPublicationStmtDrop(t *testing.T) {
	sql := "ALTER PUBLICATION mypub DROP TABLE t"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterPublicationStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateSubscriptionStmt(t *testing.T) {
	sql := "CREATE SUBSCRIPTION mysub CONNECTION 'conninfo' PUBLICATION mypub"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateSubscriptionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterSubscriptionStmt(t *testing.T) {
	sql := "ALTER SUBSCRIPTION mysub SET PUBLICATION mypub"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterSubscriptionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterSubscriptionStmtConnection(t *testing.T) {
	sql := "ALTER SUBSCRIPTION mysub CONNECTION 'newconn'"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterSubscriptionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterSubscriptionStmtEnable(t *testing.T) {
	sql := "ALTER SUBSCRIPTION mysub ENABLE"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterSubscriptionStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocPublicationTable(t *testing.T) {
	sql := "CREATE PUBLICATION mypub FOR TABLE t1"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatePublicationStmt)
	require.NotNil(t, stmt.Pubobjects)
	require.Greater(t, len(stmt.Pubobjects.Items), 0)
	spec := stmt.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
	require.NotNil(t, spec.Pubtable)
	pt := spec.Pubtable
	got := sql[pt.Loc.Start:pt.Loc.End]
	assert.Equal(t, "t1", got)
}

func TestLocRuleStmt(t *testing.T) {
	sql := "CREATE RULE myrule AS ON INSERT TO t DO INSTEAD NOTHING"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.RuleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocRuleStmtOrReplace(t *testing.T) {
	sql := "CREATE OR REPLACE RULE myrule AS ON INSERT TO t DO INSTEAD NOTHING"
	list, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.RuleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}
