package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ALTER TYPE ... ADD/DROP/ALTER ATTRIBUTE reuses AlterTableStmt/AlterTableCmd
// (gram.y AlterCompositeTypeStmt). The statement and every subcommand must
// carry the same byte ranges the ALTER TABLE path sets.
func TestLocAlterTypeAttributeCmds(t *testing.T) {
	cases := []struct {
		sql  string
		cmds []string
	}{
		{"ALTER TYPE t DROP ATTRIBUTE a", []string{"DROP ATTRIBUTE a"}},
		{"ALTER TYPE s.t ADD ATTRIBUTE b int", []string{"ADD ATTRIBUTE b int"}},
		{"ALTER TYPE t DROP ATTRIBUTE IF EXISTS a CASCADE", []string{"DROP ATTRIBUTE IF EXISTS a CASCADE"}},
		{"ALTER TYPE t ALTER ATTRIBUTE a SET DATA TYPE text COLLATE \"C\" RESTRICT", []string{"ALTER ATTRIBUTE a SET DATA TYPE text COLLATE \"C\" RESTRICT"}},
		{"ALTER TYPE t ADD ATTRIBUTE b int, DROP ATTRIBUTE a, ALTER ATTRIBUTE c TYPE text", []string{"ADD ATTRIBUTE b int", "DROP ATTRIBUTE a", "ALTER ATTRIBUTE c TYPE text"}},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			tree, err := parser.Parse(tc.sql)
			require.NoError(t, err)
			raw := tree.Items[0].(*nodes.RawStmt)
			stmt := raw.Stmt.(*nodes.AlterTableStmt)
			assert.Equal(t, int(nodes.OBJECT_TYPE), stmt.ObjType)
			assert.Equal(t, tc.sql, tc.sql[stmt.Loc.Start:stmt.Loc.End])
			require.NotNil(t, stmt.Cmds)
			require.Len(t, stmt.Cmds.Items, len(tc.cmds))
			for i, want := range tc.cmds {
				cmd := stmt.Cmds.Items[i].(*nodes.AlterTableCmd)
				assert.Equal(t, want, tc.sql[cmd.Loc.Start:cmd.Loc.End], "cmd %d", i)
			}
		})
	}
}

// A leading statement shifts every offset: the ranges must be absolute, not
// relative to the statement start.
func TestLocAlterTypeAttributeCmdsSecondStatement(t *testing.T) {
	sql := "SELECT 1; ALTER TYPE t DROP ATTRIBUTE a, ADD ATTRIBUTE b int"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	require.Len(t, tree.Items, 2)
	stmt := tree.Items[1].(*nodes.RawStmt).Stmt.(*nodes.AlterTableStmt)
	assert.Equal(t, "ALTER TYPE t DROP ATTRIBUTE a, ADD ATTRIBUTE b int", sql[stmt.Loc.Start:stmt.Loc.End])
	require.Len(t, stmt.Cmds.Items, 2)
	assert.Equal(t, "DROP ATTRIBUTE a", sql[stmt.Cmds.Items[0].(*nodes.AlterTableCmd).Loc.Start:stmt.Cmds.Items[0].(*nodes.AlterTableCmd).Loc.End])
	assert.Equal(t, "ADD ATTRIBUTE b int", sql[stmt.Cmds.Items[1].(*nodes.AlterTableCmd).Loc.Start:stmt.Cmds.Items[1].(*nodes.AlterTableCmd).Loc.End])
}

func TestLocAlterTypeRenameAttribute(t *testing.T) {
	cases := []string{
		"ALTER TYPE t RENAME ATTRIBUTE a TO b",
		"ALTER TYPE s.t RENAME ATTRIBUTE a TO b CASCADE",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			tree, err := parser.Parse(sql)
			require.NoError(t, err)
			stmt := tree.Items[0].(*nodes.RawStmt).Stmt.(*nodes.RenameStmt)
			assert.Equal(t, nodes.OBJECT_ATTRIBUTE, stmt.RenameType)
			assert.Equal(t, sql, sql[stmt.Loc.Start:stmt.Loc.End])
		})
	}
}

func TestLocAlterTypeRenameTo(t *testing.T) {
	sql := "ALTER TYPE s.t RENAME TO u"
	tree, err := parser.Parse(sql)
	require.NoError(t, err)
	stmt := tree.Items[0].(*nodes.RawStmt).Stmt.(*nodes.RenameStmt)
	assert.Equal(t, nodes.OBJECT_TYPE, stmt.RenameType)
	assert.Equal(t, sql, sql[stmt.Loc.Start:stmt.Loc.End])
}
