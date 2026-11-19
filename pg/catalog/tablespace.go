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
	if rel != nil && rel.Persistence == 't' {
		return rel.Tablespace
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
	if stmt.Owner != nil {
		if err := c.checkRole(c.ownerRoleName(stmt.Owner)); err != nil {
			return err
		}
	}
	if !strings.HasPrefix(stmt.Location, "/") {
		return errInvalidParameterValue("tablespace location must be an absolute path")
	}
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
	if err := c.checkTablespace(name); err != nil {
		if stmt.MissingOk {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("tablespace %q does not exist, skipping", name))
			return nil
		}
		return err
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
	if err := c.checkTablespace(oldName); err != nil {
		return err
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
	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			d, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch d.Defname {
			case "addroleto", "rolemembers", "adminmembers":
				list, ok := d.Arg.(*nodes.List)
				if !ok {
					continue
				}
				for _, member := range list.Items {
					name := stringVal(member)
					if rs, ok := member.(*nodes.RoleSpec); ok {
						name = c.ownerRoleName(rs)
					}
					if err := c.checkRole(name); err != nil {
						return err
					}
				}
			}
		}
	}
	if strings.HasPrefix(stmt.Role, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("role name %q is reserved", stmt.Role)}
	}
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
		for _, seq := range schema.Sequences {
			if seq.OwnerName == name {
				return true
			}
		}
	}
	for _, proc := range c.userProcs {
		if proc.OwnerName == name {
			return true
		}
	}
	for _, owner := range c.typeOwners {
		if owner == name {
			return true
		}
	}
	return false
}

// DropOwned drops the relations and indexes a DROP OWNED BY names.
func (c *Catalog) DropOwned(stmt *nodes.DropOwnedStmt) error {
	if stmt.Roles == nil {
		return nil
	}
	owners, err := c.checkedRoleSet(stmt.Roles)
	if err != nil {
		return err
	}
	for _, schema := range c.schemas {
		var rels []*Relation
		for _, rel := range schema.Relations {
			if owners[rel.OwnerName] {
				rels = append(rels, rel)
			}
		}
		for _, rel := range rels {
			if hasNormalDependent(c.relationDependents(rel.OID, 0)) {
				if stmt.Behavior != nodes.DROP_CASCADE {
					kind := "table"
					if rel.RelKind == 'v' {
						kind = "view"
					} else if rel.RelKind == 'm' {
						kind = "materialized view"
					}
					return errDependentObjects(kind, rel.Name)
				}
				c.dropDependents('r', rel.OID)
			}
			c.removeRelation(schema, rel.Name, rel)
		}
		var indexes []*Index
		for _, idx := range schema.Indexes {
			if owners[idx.OwnerName] {
				indexes = append(indexes, idx)
			}
		}
		for _, idx := range indexes {
			c.removeIndex(schema, idx.Name, idx)
		}
		var seqs []*Sequence
		for _, seq := range schema.Sequences {
			if owners[seq.OwnerName] {
				seqs = append(seqs, seq)
			}
		}
		for _, seq := range seqs {
			c.removeSequence(schema, seq)
		}
	}
	var procs []*UserProc
	for _, proc := range c.userProcs {
		if owners[proc.OwnerName] {
			procs = append(procs, proc)
		}
	}
	for _, proc := range procs {
		if deps := c.findNormalDependents('f', proc.OID); len(deps) > 0 {
			if stmt.Behavior != nodes.DROP_CASCADE {
				return errDependentObjects("function", proc.Name)
			}
			c.dropDependents('f', proc.OID)
		}
		c.removeFunction(proc.OID, proc.Name)
	}
	var typeObjects []nodes.Node
	var ownedTypes []uint32
	for oid, owner := range c.typeOwners {
		if !owners[owner] {
			continue
		}
		bt := c.typeByOID[oid]
		if bt == nil {
			delete(c.typeOwners, oid)
			continue
		}
		ownedTypes = append(ownedTypes, oid)
		if bt.Type == 'b' {
			continue // the array goes with its element type
		}
		schema := c.schemas[bt.Namespace]
		if schema == nil {
			continue
		}
		typeObjects = append(typeObjects, &nodes.TypeName{Names: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: schema.Name}, &nodes.String{Str: bt.TypeName},
		}}})
	}
	if len(typeObjects) > 0 {
		if err := c.removeTypeObjects(&nodes.DropStmt{
			RemoveType: int(nodes.OBJECT_TYPE),
			Objects:    &nodes.List{Items: typeObjects},
			Behavior:   int(stmt.Behavior),
		}); err != nil {
			return err
		}
	}
	for _, oid := range ownedTypes {
		delete(c.typeOwners, oid)
	}
	kept := c.grants[:0]
	for _, g := range c.grants {
		if !owners[g.Grantee] {
			kept = append(kept, g)
		}
	}
	c.grants = kept
	return nil
}

// ReassignOwned moves the objects a REASSIGN OWNED BY names to another role.
func (c *Catalog) ReassignOwned(stmt *nodes.ReassignOwnedStmt) error {
	if stmt.Newrole == nil {
		return nil
	}
	newOwner := c.ownerRoleName(stmt.Newrole)
	if err := c.checkRole(newOwner); err != nil {
		return err
	}
	owners, err := c.checkedRoleSet(stmt.Roles)
	if err != nil {
		return err
	}
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if owners[rel.OwnerName] {
				rel.OwnerName = newOwner
			}
		}
		for _, idx := range schema.Indexes {
			if owners[idx.OwnerName] {
				idx.OwnerName = newOwner
			}
		}
		for _, seq := range schema.Sequences {
			if owners[seq.OwnerName] {
				seq.OwnerName = newOwner
			}
		}
	}
	for _, proc := range c.userProcs {
		if owners[proc.OwnerName] {
			proc.OwnerName = newOwner
		}
	}
	for oid, owner := range c.typeOwners {
		if owners[owner] {
			c.typeOwners[oid] = newOwner
		}
	}
	return nil
}

// checkedRoleSet reads a RoleSpec list as role names, refusing one the
// catalog does not hold before any object is touched.
func (c *Catalog) checkedRoleSet(roles *nodes.List) (map[string]bool, error) {
	owners := map[string]bool{}
	if roles == nil {
		return owners, nil
	}
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

// RenameRole renames a role and every object it owns.
func (c *Catalog) RenameRole(stmt *nodes.RenameStmt) error {
	oldName := stmt.Subname
	if oldName == c.currentOwner() || oldName == c.sessionUser {
		return &Error{Code: CodeObjectNotInPrerequisiteState, Message: "current user cannot be renamed"}
	}
	if err := c.checkRole(oldName); err != nil {
		return err
	}
	if c.roles[stmt.Newname] {
		return errDuplicateObject("role", stmt.Newname)
	}
	delete(c.roles, oldName)
	c.roles[stmt.Newname] = true
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if rel.OwnerName == oldName {
				rel.OwnerName = stmt.Newname
			}
		}
		for _, idx := range schema.Indexes {
			if idx.OwnerName == oldName {
				idx.OwnerName = stmt.Newname
			}
		}
		for _, seq := range schema.Sequences {
			if seq.OwnerName == oldName {
				seq.OwnerName = stmt.Newname
			}
		}
	}
	for _, proc := range c.userProcs {
		if proc.OwnerName == oldName {
			proc.OwnerName = stmt.Newname
		}
	}
	for oid, owner := range c.typeOwners {
		if owner == oldName {
			c.typeOwners[oid] = stmt.Newname
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
		if err := c.checkRole(name); err != nil {
			if stmt.MissingOk {
				continue
			}
			return err
		}
		if !c.roles[name] {
			continue // a snapshot may know a role the catalog does not
		}
		if name == c.currentOwner() || name == c.sessionUser {
			return &Error{Code: CodeObjectNotInPrerequisiteState, Message: "current user cannot be dropped"}
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
