package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// checkTablespace refuses a tablespace name the catalog does not hold. A
// tablespace is not part of what a schema stores, but a statement that names
// one PostgreSQL does not have is still refused with 42704.
func (c *Catalog) checkTablespace(name string) error {
	if name == "" || c.tablespaces[name] {
		return nil
	}
	if c.snapshot {
		return nil
	}
	return errUndefinedObject("tablespace", name)
}

// effectiveTablespace returns the tablespace an object with an optional
// TABLESPACE clause lands in: the clause, else the session default.
func (c *Catalog) effectiveTablespace(name string) string {
	if name != "" {
		return name
	}
	return c.defaultTablespace
}

// effectiveIndexTablespace is effectiveTablespace for an index: with no
// clause and no session default, an index takes its table's tablespace.
func (c *Catalog) effectiveIndexTablespace(name string, rel *Relation) string {
	if name != "" {
		return name
	}
	if c.defaultTablespace != "" {
		return c.defaultTablespace
	}
	if rel != nil {
		return rel.Tablespace
	}
	return ""
}

// normalizeTablespace stores the default tablespace as "", the way a
// relation created without a TABLESPACE clause holds it.
func normalizeTablespace(name string) string {
	if name == "pg_default" {
		return ""
	}
	return name
}

// checkRelationTablespace refuses a tablespace an ordinary relation cannot be
// placed in. pg_global holds shared system relations only, even though the
// tablespace itself exists.
func (c *Catalog) checkRelationTablespace(name string) error {
	if err := c.checkTablespace(name); err != nil {
		return err
	}
	if name == "pg_global" {
		return errInvalidParameterValue("only shared relations can be placed in pg_global tablespace")
	}
	return nil
}

// tablespaceMatches reports whether an object's stored tablespace is the one
// a bulk move names. The empty value is the default, pg_default.
func tablespaceMatches(stored, name string) bool {
	if name == "pg_default" {
		return stored == "" || stored == "pg_default"
	}
	return stored == name
}

// checkTableAccessMethod refuses a table access method the catalog does not
// hold, or one registered for indexes: CREATE TABLE ... USING takes only a
// table access method.
func (c *Catalog) checkTableAccessMethod(name string) error {
	if name == "" {
		return nil
	}
	am := c.accessMethodByName[name]
	if am == nil {
		if c.snapshot {
			return nil
		}
		return errUndefinedObject("access method", name)
	}
	if am.Type != 't' {
		return &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("access method %q is not of type TABLE", name)}
	}
	return nil
}

// DefineTablespace records a CREATE TABLESPACE, which stores nothing else.
//
// pg: src/backend/commands/tablespace.c — CreateTableSpace
func (c *Catalog) DefineTablespace(stmt *nodes.CreateTableSpaceStmt) error {
	name := stmt.Tablespacename
	if strings.HasPrefix(name, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("unacceptable tablespace name %q", name)}
	}
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
	// The two a server starts with cannot be dropped, even by a superuser.
	if name == "pg_default" || name == "pg_global" {
		return &Error{Code: CodeInsufficientPrivilege,
			Message: fmt.Sprintf("permission denied for tablespace %s", name)}
	}
	if !c.tablespaces[name] {
		if stmt.MissingOk {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("tablespace %q does not exist, skipping", name))
			return nil
		}
		return errUndefinedObject("tablespace", name)
	}
	if c.tablespaceInUse(name) {
		return &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("tablespace %q is not empty", name)}
	}
	delete(c.tablespaces, name)
	return nil
}

// RenameTablespace renames a tablespace and every object stored in it.
//
// pg: src/backend/commands/tablespace.c — RenameTableSpace
func (c *Catalog) RenameTablespace(stmt *nodes.RenameStmt) error {
	oldName := stmt.Subname
	if strings.HasPrefix(stmt.Newname, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("unacceptable tablespace name %q", stmt.Newname)}
	}
	if oldName == "pg_default" || oldName == "pg_global" {
		return &Error{Code: CodeInsufficientPrivilege,
			Message: fmt.Sprintf("permission denied for tablespace %s", oldName)}
	}
	if !c.tablespaces[oldName] {
		return errUndefinedObject("tablespace", oldName)
	}
	if c.tablespaces[stmt.Newname] {
		return errDuplicateObject("tablespace", stmt.Newname)
	}
	delete(c.tablespaces, oldName)
	c.tablespaces[stmt.Newname] = true
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if rel.Tablespace == oldName {
				rel.Tablespace = stmt.Newname
			}
		}
		for _, idx := range schema.Indexes {
			if idx.Tablespace == oldName {
				idx.Tablespace = stmt.Newname
			}
		}
	}
	return nil
}

// AlterTablespaceOptions checks the tablespace an ALTER TABLESPACE names.
//
// pg: src/backend/commands/tablespace.c — AlterTableSpaceOptions
func (c *Catalog) AlterTablespaceOptions(stmt *nodes.AlterTableSpaceOptionsStmt) error {
	return c.checkTablespace(stmt.Tablespacename)
}

// currentOwner is the role a new object belongs to. A catalog with no
// session user is the server's bootstrap superuser, postgres.
func (c *Catalog) currentOwner() string {
	if c.currentUser != "" {
		return c.currentUser
	}
	return "postgres"
}

// checkRole refuses a role name the catalog does not hold. A snapshot
// carries no roles, so an unobserved one may still exist on its server.
func (c *Catalog) checkRole(name string) error {
	if name == "" || c.roles[name] || c.snapshot {
		return nil
	}
	return errUndefinedObject("role", name)
}

// DefineRole records a CREATE ROLE.
func (c *Catalog) DefineRole(stmt *nodes.CreateRoleStmt) error {
	if c.roles[stmt.Role] {
		return errDuplicateObject("role", stmt.Role)
	}
	c.roles[stmt.Role] = true
	return nil
}

// roleInUse reports a relation or index the current schemas hold for a role.
func (c *Catalog) roleInUse(name string) bool {
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if rel.OwnerName == name {
				return true
			}
		}
		for _, idx := range schema.Indexes {
			if idx.OwnerName == name {
				return true
			}
		}
	}
	return false
}

// DropOwned releases the ownership a DROP OWNED BY names. The catalog keeps
// the objects, but a later DROP ROLE no longer sees them as dependents.
func (c *Catalog) DropOwned(stmt *nodes.DropOwnedStmt) error {
	if stmt.Roles == nil {
		return nil
	}
	owners := map[string]bool{}
	for _, item := range stmt.Roles.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		owners[c.ownerRoleName(rs)] = true
	}
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if owners[rel.OwnerName] {
				rel.OwnerName = ""
			}
		}
		for _, idx := range schema.Indexes {
			if owners[idx.OwnerName] {
				idx.OwnerName = ""
			}
		}
	}
	return nil
}

// RemoveRoles forgets DROP ROLE.
func (c *Catalog) RemoveRoles(stmt *nodes.DropRoleStmt) error {
	if stmt.Roles == nil {
		return nil
	}
	var names []string
	for _, item := range stmt.Roles.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		name := c.ownerRoleName(rs)
		if !c.roles[name] {
			if stmt.MissingOk {
				continue
			}
			return errUndefinedObject("role", name)
		}
		if c.roleInUse(name) {
			return errDependentObjects("role", name)
		}
		names = append(names, name)
	}
	for _, name := range names {
		delete(c.roles, name)
	}
	return nil
}

// ownerRoleName resolves a role spec to the role name an object's owner is
// recorded under.
func (c *Catalog) ownerRoleName(rs *nodes.RoleSpec) string {
	name := roleSpecName(rs)
	switch name {
	case "CURRENT_USER", "CURRENT_ROLE":
		return c.currentOwner()
	case "SESSION_USER":
		if c.sessionUser != "" {
			return c.sessionUser
		}
		return c.currentOwner()
	}
	return name
}

// moveAllOwners reads the OWNED BY list of a bulk move. Nil means every
// owner.
func (c *Catalog) moveAllOwners(roles *nodes.List) (map[string]bool, error) {
	if roles == nil || len(roles.Items) == 0 {
		return nil, nil
	}
	owners := make(map[string]bool, len(roles.Items))
	for _, item := range roles.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		name := c.ownerRoleName(rs)
		if err := c.checkRole(name); err != nil {
			return nil, err
		}
		owners[name] = true
	}
	return owners, nil
}

// ownerMatches reports whether an object's owner is in a bulk move's OWNED
// BY list. A nil list takes every owner.
func ownerMatches(owner string, owners map[string]bool) bool {
	return owners == nil || owners[owner]
}

// MoveAllTablespaces checks both tablespaces an ALTER TABLE/INDEX/
// MATERIALIZED VIEW ALL IN TABLESPACE statement names, then moves the
// relations and indexes the source holds. PostgreSQL validates the
// destination first, then the source.
//
// pg: src/backend/commands/tablecmds.c — AlterTableMoveAll
func (c *Catalog) MoveAllTablespaces(stmt *nodes.AlterTableMoveAllStmt) error {
	if err := c.checkTablespace(stmt.NewTablespacename); err != nil {
		return err
	}
	if err := c.checkTablespace(stmt.OrigTablespacename); err != nil {
		return err
	}
	if stmt.NewTablespacename == "pg_global" || stmt.OrigTablespacename == "pg_global" {
		return errInvalidParameterValue("cannot move relations in to or out of pg_global tablespace")
	}
	target := normalizeTablespace(stmt.NewTablespacename)
	owners, err := c.moveAllOwners(stmt.Roles)
	if err != nil {
		return err
	}
	for _, schema := range c.schemas {
		switch nodes.ObjectType(stmt.ObjType) {
		case nodes.OBJECT_TABLE:
			for _, rel := range schema.Relations {
				if rel.Persistence != 't' && (rel.RelKind == 'r' || rel.RelKind == 'p') && tablespaceMatches(rel.Tablespace, stmt.OrigTablespacename) && ownerMatches(rel.OwnerName, owners) {
					rel.Tablespace = target
				}
			}
		case nodes.OBJECT_MATVIEW:
			for _, rel := range schema.Relations {
				if rel.Persistence != 't' && rel.RelKind == 'm' && tablespaceMatches(rel.Tablespace, stmt.OrigTablespacename) && ownerMatches(rel.OwnerName, owners) {
					rel.Tablespace = target
				}
			}
		case nodes.OBJECT_INDEX:
			for _, idx := range schema.Indexes {
				rel := c.relationByOID[idx.RelOID]
				if rel == nil || rel.Persistence == 't' {
					continue
				}
				if tablespaceMatches(idx.Tablespace, stmt.OrigTablespacename) && ownerMatches(idx.OwnerName, owners) {
					idx.Tablespace = target
				}
			}
		}
	}
	return nil
}

// tablespaceInUse reports a relation or index the current schemas hold in a
// tablespace.
func (c *Catalog) tablespaceInUse(name string) bool {
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if rel.Tablespace == name {
				return true
			}
		}
		for _, idx := range schema.Indexes {
			if idx.Tablespace == name {
				return true
			}
		}
	}
	return false
}
