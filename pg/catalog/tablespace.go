package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// checkTablespace refuses a tablespace name the catalog does not hold. A
// tablespace is not part of what a schema stores, but a statement that names
// one PostgreSQL does not have is still refused with 42704.
func (c *Catalog) checkTablespace(name string) error {
	if name == "" || c.tablespaces[name] {
		return nil
	}
	return errUndefinedObject("tablespace", name)
}

// DefineTablespace records a CREATE TABLESPACE, which stores nothing else.
//
// pg: src/backend/commands/tablespace.c — CreateTableSpace
func (c *Catalog) DefineTablespace(stmt *nodes.CreateTableSpaceStmt) error {
	name := stmt.Tablespacename
	if c.tablespaces[name] {
		return errDuplicateObject("tablespace", name)
	}
	c.tablespaces[name] = true
	return nil
}

// RemoveTablespace forgets a DROP TABLESPACE, which stores nothing else.
//
// pg: src/backend/commands/tablespace.c — RemoveTablespace
func (c *Catalog) RemoveTablespace(stmt *nodes.DropTableSpaceStmt) error {
	name := stmt.Tablespacename
	if !c.tablespaces[name] {
		if stmt.MissingOk {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("tablespace %q does not exist, skipping", name))
			return nil
		}
		return errUndefinedObject("tablespace", name)
	}
	delete(c.tablespaces, name)
	return nil
}

// AlterTablespaceOptions checks the tablespace an ALTER TABLESPACE names.
//
// pg: src/backend/commands/tablespace.c — AlterTableSpaceOptions
func (c *Catalog) AlterTablespaceOptions(stmt *nodes.AlterTableSpaceOptionsStmt) error {
	return c.checkTablespace(stmt.Tablespacename)
}

// MoveAllTablespaces checks both tablespaces an ALTER TABLE/INDEX/
// MATERIALIZED VIEW ALL IN TABLESPACE statement names. PostgreSQL validates
// the destination first, then the source.
//
// pg: src/backend/commands/tablecmds.c — AlterTableMoveAll
func (c *Catalog) MoveAllTablespaces(stmt *nodes.AlterTableMoveAllStmt) error {
	if err := c.checkTablespace(stmt.NewTablespacename); err != nil {
		return err
	}
	return c.checkTablespace(stmt.OrigTablespacename)
}
