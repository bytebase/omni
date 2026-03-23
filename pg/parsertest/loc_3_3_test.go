package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocGrantStmt(t *testing.T) {
	sql := "GRANT SELECT ON t TO myrole"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.GrantStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocGrantRoleStmt(t *testing.T) {
	sql := "GRANT myrole TO otheruser"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.GrantRoleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreateRoleStmt(t *testing.T) {
	sql := "CREATE ROLE myrole WITH LOGIN"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreateRoleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterRoleStmt(t *testing.T) {
	sql := "ALTER ROLE myrole WITH SUPERUSER"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterRoleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterRoleSetStmt(t *testing.T) {
	sql := "ALTER ROLE myrole SET search_path TO public"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterRoleSetStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAlterPolicyStmt(t *testing.T) {
	sql := "ALTER POLICY mypol ON t USING (true)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.AlterPolicyStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocCreatePolicyStmt(t *testing.T) {
	sql := "CREATE POLICY mypol ON t USING (true)"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CreatePolicyStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocDropRoleStmt(t *testing.T) {
	sql := "DROP ROLE myrole"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DropRoleStmt)
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	assert.Equal(t, sql, got)
}

func TestLocAccessPriv(t *testing.T) {
	sql := "GRANT SELECT ON t TO myrole"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	raw := tree.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.GrantStmt)
	require.NotNil(t, stmt.Privileges)
	require.Len(t, stmt.Privileges.Items, 1)
	priv := stmt.Privileges.Items[0].(*nodes.AccessPriv)
	got := sql[priv.Loc.Start:priv.Loc.End]
	assert.Equal(t, "SELECT", got)
}
