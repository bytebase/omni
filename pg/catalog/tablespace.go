package catalog

import (
	"cmp"
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// checkTablespace refuses a tablespace name the catalog does not hold. A
// tablespace is not part of what a schema stores, but a statement that names
// one PostgreSQL does not have is still refused with 42704.
func (c *Catalog) checkTablespace(name string) error {
	if name == "" {
		return nil
	}
	if exists, known := c.tablespaces[name]; known {
		if exists {
			return nil
		}
		return errUndefinedObject("tablespace", name)
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
	if c.defaultTablespace == "" {
		return ""
	}
	if exists, known := c.tablespaces[c.defaultTablespace]; known {
		if exists {
			return c.defaultTablespace
		}
		return ""
	}
	if c.snapshot {
		return c.defaultTablespace
	}
	return ""
}

// effectiveIndexTablespace returns the tablespace an index lands in: its own
// TABLESPACE clause when given, "" on a temporary relation, else the session's
// effective default_tablespace (which falls back to the database default, not
// the table's tablespace).
func (c *Catalog) effectiveIndexTablespace(name string, rel *Relation) string {
	if name != "" {
		return name
	}
	if rel != nil && rel.Persistence == 't' {
		return ""
	}
	return c.effectiveTablespace("")
}

// checkIndexTablespace validates an index's target tablespace, including
// PostgreSQL's prohibition on naming pg_default for partitioned indexes.
func (c *Catalog) checkIndexTablespace(name string, rel *Relation) error {
	if rel != nil && rel.RelKind == 'p' && name == "pg_default" {
		return &Error{Code: CodeFeatureNotSupported,
			Message: "cannot specify default tablespace for partitioned relations"}
	}
	return c.checkRelationTablespace(name)
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
	owner := c.currentOwner()
	if stmt.Owner != nil {
		owner = c.ownerRoleName(stmt.Owner)
	}
	c.tablespaceOwners[name] = owner
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
	if c.snapshot {
		c.tablespaces[name] = false
	} else {
		delete(c.tablespaces, name)
	}
	delete(c.tablespaceOwners, name)
	return nil
}

// ExecAlterOwnerStmt dispatches ALTER ... OWNER TO for object classes whose
// ownership is tracked in the catalog.
func (c *Catalog) ExecAlterOwnerStmt(stmt *nodes.AlterOwnerStmt) error {
	switch nodes.ObjectType(stmt.ObjectType) {
	case nodes.OBJECT_TABLESPACE:
		return c.AlterTablespaceOwner(stmt)
	case nodes.OBJECT_DATABASE:
		return c.alterDatabaseOwner(stmt)
	case nodes.OBJECT_SCHEMA:
		return c.alterSchemaOwner(stmt)
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return c.alterTypeOwner(stmt)
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE, nodes.OBJECT_AGGREGATE:
		return c.alterRoutineOwner(stmt)
	case nodes.OBJECT_OPERATOR:
		return c.alterOperatorOwner(stmt)
	case nodes.OBJECT_OPCLASS:
		return c.alterOpClassOwner(stmt)
	case nodes.OBJECT_OPFAMILY:
		return c.alterOpFamilyOwner(stmt)
	case nodes.OBJECT_EVENT_TRIGGER:
		return c.alterEventTriggerOwner(stmt)
	case nodes.OBJECT_EXTENSION:
		return c.alterExtensionOwner(stmt)
	case nodes.OBJECT_COLLATION:
		return c.alterCollationOwner(stmt)
	case nodes.OBJECT_FDW:
		return c.alterFdwOwner(stmt)
	case nodes.OBJECT_FOREIGN_SERVER:
		return c.alterForeignServerOwner(stmt)
	default:
		return nil
	}
}

// alterOperatorOwner changes the owner recorded for a user-defined operator.
func (c *Catalog) alterOperatorOwner(stmt *nodes.AlterOwnerStmt) error {
	owa, ok := stmt.Object.(*nodes.ObjectWithArgs)
	if !ok {
		return nil
	}
	found, err := c.lookupUserOperator(owa)
	if err != nil {
		return err
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	c.operatorOwners[found.OID] = owner
	return nil
}

// alterExtensionOwner changes the owner recorded for an extension.
func (c *Catalog) alterExtensionOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	ext := c.extByName[name]
	if ext == nil {
		return errUndefinedObject("extension", name)
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	ext.OwnerName = owner
	return nil
}

func (c *Catalog) lookupUserCollation(obj nodes.Node) (*Schema, *UserCollation, error) {
	schemaName, name := extractObjectName(obj)
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, nil, errUndefinedSchema(schemaName)
		}
		if uc := s.Collations[name]; uc != nil {
			return s, uc, nil
		}
		if s.OID == PGCatalogNamespace && isBuiltinCollationName(name) {
			return s, nil, nil
		}
		return nil, nil, errUndefinedObject("collation", name)
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if uc := s.Collations[name]; uc != nil {
			return s, uc, nil
		}
		if nsOID == PGCatalogNamespace && isBuiltinCollationName(name) {
			return s, nil, nil
		}
	}
	return nil, nil, errUndefinedObject("collation", name)
}

func (c *Catalog) alterCollationOwner(stmt *nodes.AlterOwnerStmt) error {
	_, uc, err := c.lookupUserCollation(stmt.Object)
	if err != nil {
		if c.snapshot && missingObject(err) {
			return nil
		}
		return err
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	if uc != nil {
		uc.OwnerName = owner
	}
	return nil
}

// alterEventTriggerOwner changes the owner recorded for an event trigger.
func (c *Catalog) alterEventTriggerOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	evt := c.eventTriggerByName[name]
	if evt == nil {
		return errUndefinedObject("event trigger", name)
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	evt.OwnerName = owner
	return nil
}

// AlterTablespaceOwner changes the owner recorded for a tablespace.
func (c *Catalog) AlterTablespaceOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	if err := c.checkTablespace(name); err != nil {
		return err
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	c.tablespaceOwners[name] = owner
	return nil
}

func (c *Catalog) alterDatabaseOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	if name != "" && name == c.databaseName {
		c.databaseOwners[name] = owner
		return nil
	}
	if name != "" && c.databaseOwners[name] != "" {
		c.databaseOwners[name] = owner
		return nil
	}
	if name == "" && c.databaseName == "" && c.databaseOwners["\x00current"] != "" {
		c.databaseOwners["\x00current"] = owner
		return nil
	}
	if c.snapshot || c.databaseName != "" || len(c.databaseOwners) > 0 {
		return errInvalidCatalogName(name)
	}
	return nil
}

func (c *Catalog) alterSchemaOwner(stmt *nodes.AlterOwnerStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	s := c.schemaByName[name]
	if s == nil {
		return errUndefinedSchema(name)
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	oldOwner := s.effectiveOwner()
	s.Owner = owner
	s.OwnerName = owner
	c.transferACLs('n', s.OID, oldOwner, owner)
	return nil
}

func (c *Catalog) alterTypeOwner(stmt *nodes.AlterOwnerStmt) error {
	if stmt.Newowner != nil {
		if err := c.checkRole(c.ownerRoleName(stmt.Newowner)); err != nil {
			return err
		}
	}
	schemaName, typeName := extractObjectName(stmt.Object)
	var bt *BuiltinType
	if schemaName == "" {
		for _, nsOID := range c.searchPathWithCatalog() {
			if found := c.typeByName[typeKey{ns: nsOID, name: typeName}]; found != nil {
				bt = found
				break
			}
		}
	} else {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
		bt = c.typeByName[typeKey{ns: schema.OID, name: typeName}]
	}
	if bt == nil {
		return errUndefinedType(typeName)
	}
	if nodes.ObjectType(stmt.ObjectType) == nodes.OBJECT_DOMAIN && bt.Type != 'd' {
		return errWrongObjectType(typeName, "a domain")
	}
	if bt.Type == 'c' && bt.RelID != 0 {
		rel := c.relationByOID[bt.RelID]
		if rel != nil && rel.RelKind != 'c' {
			return errWrongObjectType(typeName, "a composite type")
		}
		if stmt.Newowner == nil {
			return nil
		}
		owner := c.ownerRoleName(stmt.Newowner)
		if err := c.checkRole(owner); err != nil {
			return err
		}
		if rel != nil {
			oldOwner := rel.OwnerName
			rel.OwnerName = owner
			c.transferACLs('r', rel.OID, oldOwner, owner)
			if rel.RowTypeOID != 0 {
				if _, ok := c.typeOwners[rel.RowTypeOID]; ok {
					c.typeOwners[rel.RowTypeOID] = owner
				}
				c.transferACLs('T', rel.RowTypeOID, oldOwner, owner)
			}
			if rel.ArrayOID != 0 {
				if _, ok := c.typeOwners[rel.ArrayOID]; ok {
					c.typeOwners[rel.ArrayOID] = owner
				}
				c.transferACLs('T', rel.ArrayOID, oldOwner, owner)
			}
		}
		return nil
	}
	if bt.Category == 'A' && bt.Elem != 0 {
		return errWrongObjectType(typeName, "a non-array type")
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	for _, oid := range c.typeFamilyOIDs(bt) {
		oldOwner := c.typeOwners[oid]
		c.typeOwners[oid] = owner
		c.transferACLs('T', oid, oldOwner, owner)
	}
	return nil
}

func (c *Catalog) alterRoutineOwner(stmt *nodes.AlterOwnerStmt) error {
	if stmt.Newowner != nil {
		if err := c.checkRole(c.ownerRoleName(stmt.Newowner)); err != nil {
			return err
		}
	}
	owa, ok := stmt.Object.(*nodes.ObjectWithArgs)
	if !ok {
		return nil
	}
	objType := nodes.ObjectType(stmt.ObjectType)
	kindStr := "function"
	if objType == nodes.OBJECT_PROCEDURE {
		kindStr = "procedure"
	} else if objType == nodes.OBJECT_ROUTINE {
		kindStr = "routine"
	}
	schemaName, funcName := qualifiedName(owa.Objname)
	var schema *Schema
	if schemaName == "" {
		// PostgreSQL matches the full signature first and then reports the
		// first visible schema that holds it, so follow that order here.
		argOIDs, argErr := c.routineArgOIDs(owa)
		if argErr != nil {
			return argErr
		}
		for _, nsOID := range c.routineSearchPathWithCatalog() {
			ns := c.schemas[nsOID]
			if ns == nil {
				continue
			}
			if argOIDs == nil {
				if len(c.findUserProcsByName(ns, funcName)) > 0 {
					schema = ns
					break
				}
				continue
			}
			if c.findExactProc(ns, funcName, argOIDs) != nil {
				schema = ns
				break
			}
		}
	} else {
		var err error
		schema, err = c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
	}
	if schema == nil {
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
	}
	var bp *BuiltinProc
	if !owa.ArgsUnspecified {
		var argOIDs []uint32
		if owa.Objargs != nil {
			argOIDs = make([]uint32, 0, len(owa.Objargs.Items))
			for _, item := range owa.Objargs.Items {
				var tn *nodes.TypeName
				switch n := item.(type) {
				case *nodes.TypeName:
					tn = n
				case *nodes.FunctionParameter:
					tn = n.ArgType
				}
				if tn == nil {
					continue
				}
				oid, _, err := c.resolveTypeName(tn)
				if err != nil {
					return err
				}
				argOIDs = append(argOIDs, oid)
			}
		}
		bp = c.findExactProc(schema, funcName, argOIDs)
	} else {
		candidates := c.findUserProcsByName(schema, funcName)
		if len(candidates) > 1 {
			return errAmbiguousFunction(funcName)
		}
		if len(candidates) == 1 {
			bp = c.procByOID[candidates[0].OID]
		}
	}
	if bp == nil {
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
	}
	up := c.userProcs[bp.OID]
	if up == nil {
		return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%s %q is not user-defined", kindStr, funcName)}
	}
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		if bp.Kind != 'p' {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a procedure", funcName)}
		}
	case nodes.OBJECT_FUNCTION:
		if bp.Kind == 'p' {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a function", funcName)}
		}
		if bp.Kind == PROKIND_AGGREGATE {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is an aggregate function", funcName)}
		}
	case nodes.OBJECT_AGGREGATE:
		if bp.Kind != PROKIND_AGGREGATE {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an aggregate", funcName)}
		}
	}
	if stmt.Newowner == nil {
		return nil
	}
	owner := c.ownerRoleName(stmt.Newowner)
	if err := c.checkRole(owner); err != nil {
		return err
	}
	oldOwner := up.OwnerName
	up.OwnerName = owner
	c.transferACLs('f', up.OID, oldOwner, owner)
	return nil
}

// routineArgOIDs resolves an ObjectWithArgs' argument types, or nil when the
// statement named no signature.
func (c *Catalog) routineArgOIDs(owa *nodes.ObjectWithArgs) ([]uint32, error) {
	if owa == nil || owa.ArgsUnspecified {
		return nil, nil
	}
	if owa.Objargs == nil {
		return []uint32{}, nil
	}
	argOIDs := make([]uint32, 0, len(owa.Objargs.Items))
	for _, item := range owa.Objargs.Items {
		var tn *nodes.TypeName
		switch n := item.(type) {
		case *nodes.TypeName:
			tn = n
		case *nodes.FunctionParameter:
			tn = n.ArgType
		}
		if tn == nil {
			continue
		}
		oid, _, err := c.resolveTypeName(tn)
		if err != nil {
			return nil, err
		}
		argOIDs = append(argOIDs, oid)
	}
	return argOIDs, nil
}

func (c *Catalog) execAlterSequenceTableCmds(seq *Sequence, relName string, cmds *nodes.List) error {
	if cmds == nil {
		return nil
	}
	var newOwner string
	hasOwner := false
	// Validate every subcommand before mutating seq.
	// pg: src/backend/commands/tablecmds.c — ATSimplePermissions (ATT_SEQUENCE)
	for _, item := range cmds.Items {
		atc, ok := item.(*nodes.AlterTableCmd)
		if !ok {
			continue
		}
		switch nodes.AlterTableType(atc.Subtype) {
		case nodes.AT_SetLogged, nodes.AT_SetUnLogged:
			if seq.Persistence == 't' {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("cannot change logged status of table %q because it is temporary", seq.Name),
				}
			}
		case nodes.AT_ChangeOwner:
			if atc.Newowner == nil {
				continue
			}
			owner := c.ownerRoleName(atc.Newowner)
			if err := c.checkRole(owner); err != nil {
				return err
			}
			if seq.OwnerRelOID != 0 && owner != seq.OwnerName {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: fmt.Sprintf("cannot change owner of sequence %q", seq.Name),
				}
			}
			newOwner = owner
			hasOwner = true
		default:
			return errWrongObjectType(relName, "a table")
		}
	}
	for _, item := range cmds.Items {
		atc, ok := item.(*nodes.AlterTableCmd)
		if !ok {
			continue
		}
		switch nodes.AlterTableType(atc.Subtype) {
		case nodes.AT_SetLogged:
			seq.Persistence = 'p'
		case nodes.AT_SetUnLogged:
			seq.Persistence = 'u'
		}
	}
	if hasOwner {
		oldOwner := seq.OwnerName
		seq.OwnerName = newOwner
		c.transferACLs('s', seq.OID, oldOwner, newOwner)
	}
	return nil
}

// RenameTablespace renames a tablespace and every object stored in it.
//
// pg: src/backend/commands/tablespace.c — RenameTableSpace
func (c *Catalog) RenameTablespace(stmt *nodes.RenameStmt) error {
	oldName := stmt.Subname
	if err := c.checkTablespace(oldName); err != nil {
		return err
	}
	if strings.HasPrefix(stmt.Newname, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("unacceptable tablespace name %q", stmt.Newname)}
	}
	if oldName == "pg_default" || oldName == "pg_global" {
		return &Error{Code: CodeInsufficientPrivilege,
			Message: fmt.Sprintf("permission denied for tablespace %s", oldName)}
	}
	if c.tablespaces[stmt.Newname] || oldName == stmt.Newname {
		return errDuplicateObject("tablespace", stmt.Newname)
	}
	if c.snapshot {
		c.tablespaces[oldName] = false
	} else {
		delete(c.tablespaces, oldName)
	}
	c.tablespaces[stmt.Newname] = true
	if owner, ok := c.tablespaceOwners[oldName]; ok {
		delete(c.tablespaceOwners, oldName)
		c.tablespaceOwners[stmt.Newname] = owner
	}
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
// session user is the server's bootstrap superuser, postgres; during snapshot
// installation no default owner is invented for unattributed objects.
func (c *Catalog) currentOwner() string {
	if c.installingSnapshot {
		return ""
	}
	if c.currentUser != "" {
		return c.currentUser
	}
	return "postgres"
}

func (c *Catalog) recordTypeOwner(oid uint32) {
	if owner := c.currentOwner(); owner != "" {
		c.typeOwners[oid] = owner
	}
}

// checkRole refuses a role name the catalog does not hold. A snapshot
// carries no roles, so an unobserved one may still exist on its server.
func (c *Catalog) checkRole(name string) error {
	if name == "" {
		return nil
	}
	if exists, known := c.roles[name]; known {
		if exists {
			return nil
		}
		return errUndefinedObject("role", name)
	}
	if c.snapshot {
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

// roleInUse reports whether any catalog object or ACL/policy entry references
// the named role.
func (c *Catalog) roleInUse(name string) bool {
	for _, schema := range c.schemas {
		if schema.OID != PGCatalogNamespace && schema.OID != PGToastNamespace && schema.effectiveOwner() == name {
			return true
		}
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
		for _, uc := range schema.Collations {
			if uc != nil && uc.OwnerName == name {
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
	for _, owner := range c.operatorOwners {
		if owner == name {
			return true
		}
	}
	for _, g := range c.grants {
		if g.Grantee == name || g.Grantor == name {
			return true
		}
	}
	for _, p := range c.policies {
		for _, r := range p.Roles {
			if r != "" && r != "public" && r == name {
				return true
			}
		}
	}
	for _, owner := range c.tablespaceOwners {
		if owner == name {
			return true
		}
	}
	for _, owner := range c.databaseOwners {
		if owner == name {
			return true
		}
	}
	for _, ext := range c.extensions {
		if ext.OwnerName == name {
			return true
		}
	}
	for _, evt := range c.eventTriggers {
		if evt.OwnerName == name {
			return true
		}
	}
	for _, opc := range c.opClasses {
		if opc.OwnerName == name {
			return true
		}
	}
	for _, fam := range c.opFamilies {
		if fam.OwnerName == name {
			return true
		}
	}
	return false
}

// ownerAMUser is an operator class or family that uses an access method or
// family of an owned extension.
type ownerAMUser struct {
	dep     DepEntry
	extName string
}

type ownerRelExt struct {
	rel     *Relation
	extName string
}

// ownerIndexExt pairs an index with the extension whose operator class or
// access method it uses, for the DROP OWNED report.
type ownerIndexExt struct {
	idx     *Index
	extName string
}

func (c *Catalog) objectOwnerName(objType byte, objOID uint32) string {
	switch objType {
	case 'r':
		if rel := c.relationByOID[objOID]; rel != nil {
			return rel.OwnerName
		}
	case 's':
		if seq := c.sequenceByOID[objOID]; seq != nil {
			return seq.OwnerName
		}
	case 'f':
		if up := c.userProcs[objOID]; up != nil {
			return up.OwnerName
		}
	case 'n':
		if s := c.schemas[objOID]; s != nil {
			return s.effectiveOwner()
		}
	case 'T':
		if owner := c.typeOwners[objOID]; owner != "" {
			return owner
		}
		if bt := c.typeByOID[objOID]; bt != nil && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil {
				return rel.OwnerName
			}
		}
	}
	return ""
}

// DropOwned drops the objects owned by the named roles and revokes their
// privileges, validating the complete deletion set before mutating the catalog.
func (c *Catalog) DropOwned(stmt *nodes.DropOwnedStmt) error {
	if stmt.Roles == nil {
		return nil
	}
	owners, err := c.checkedRoleSet(stmt.Roles)
	if err != nil {
		return err
	}
	cascade := stmt.Behavior == nodes.DROP_CASCADE

	var ownedSchemas []*Schema
	ownedSchemaOIDs := map[uint32]bool{}
	for _, schema := range c.schemas {
		if schema.OID == PGCatalogNamespace || schema.OID == PGToastNamespace {
			continue
		}
		if so := schema.effectiveOwner(); so != "" && owners[so] {
			ownedSchemas = append(ownedSchemas, schema)
			ownedSchemaOIDs[schema.OID] = true
		}
	}
	sort.Slice(ownedSchemas, func(i, j int) bool { return ownedSchemas[i].OID < ownedSchemas[j].OID })

	var ownedExts []*Extension
	ownedExtNames := map[string]bool{}
	for _, ext := range c.extensions {
		if (ext.OwnerName != "" && owners[ext.OwnerName]) || (cascade && ownedSchemaOIDs[ext.SchemaOID]) {
			ownedExts = append(ownedExts, ext)
			ownedExtNames[ext.Name] = true
		}
	}
	sort.Slice(ownedExts, func(i, j int) bool { return ownedExts[i].OID < ownedExts[j].OID })
	inOwnedExt := func(oid uint32) bool {
		extName := c.extensionMembers[oid]
		return extName != "" && ownedExtNames[extName]
	}
	ownedExtOIDs := map[uint32]bool{}
	for _, ext := range ownedExts {
		ownedExtOIDs[ext.OID] = true
	}
	autoExtOIDs := map[uint32]bool{}
	for changed := true; changed; {
		changed = false
		for _, d := range c.deps {
			if d.RefType != 'e' || !ownedExtOIDs[d.RefOID] || d.DepType != DepAutoExtension {
				continue
			}
			if d.ObjOID < FirstNormalObjectId {
				return &Error{
					Code:    CodeDependentObjects,
					Message: "cannot drop system object because it is required by the database system",
				}
			}
			autoExtOIDs[d.ObjOID] = true
			if ownerExtName := c.extensionMembers[d.ObjOID]; ownerExtName != "" && !ownedExtNames[ownerExtName] {
				if !cascade {
					return errDependentObjects("extension", ownerExtName)
				}
				if ownerExt := c.extByName[ownerExtName]; ownerExt != nil {
					ownedExtNames[ownerExtName] = true
					ownedExtOIDs[ownerExt.OID] = true
					ownedExts = append(ownedExts, ownerExt)
					changed = true
				}
			}
		}
	}
	sort.Slice(ownedExts, func(i, j int) bool { return ownedExts[i].OID < ownedExts[j].OID })
	for _, schema := range c.schemas {
		if schema.OID == PGCatalogNamespace || schema.OID == PGToastNamespace || ownedSchemaOIDs[schema.OID] {
			continue
		}
		if inOwnedExt(schema.OID) {
			ownedSchemas = append(ownedSchemas, schema)
			ownedSchemaOIDs[schema.OID] = true
		}
	}
	sort.Slice(ownedSchemas, func(i, j int) bool { return ownedSchemas[i].OID < ownedSchemas[j].OID })
	// opclassIndexes holds the indexes that use an operator class of an owned
	// extension: a RESTRICT drop reports them, CASCADE takes them with it.
	var opclassIndexes []ownerIndexExt
	opclassIndexOIDs := map[uint32]bool{}
	// amRelations holds relations that use an access method of an owned
	// extension; CASCADE drops them with it.
	var amRelations []ownerRelExt
	amRelOIDs := map[uint32]bool{}
	var amIndexes []ownerIndexExt
	var amUsers []ownerAMUser
	amUserSeen := map[[2]uint32]bool{}
	for oid, extName := range c.extensionMembers {
		if !ownedExtNames[extName] {
			continue
		}
		for _, dep := range c.findNormalDependents('o', oid) {
			if dep.ObjType != 'i' || opclassIndexOIDs[dep.ObjOID] {
				continue
			}
			if idx := c.indexes[dep.ObjOID]; idx != nil {
				opclassIndexes = append(opclassIndexes, ownerIndexExt{idx: idx, extName: extName})
				opclassIndexOIDs[idx.OID] = true
			}
		}
		if c.accessMethods[oid] != nil {
			for _, dep := range c.findNormalDependents('m', oid) {
				switch dep.ObjType {
				case 'r':
					if rel := c.relationByOID[dep.ObjOID]; rel != nil && !amRelOIDs[rel.OID] {
						amRelOIDs[rel.OID] = true
						amRelations = append(amRelations, ownerRelExt{rel: rel, extName: extName})
					}
				case 'i':
					if idx := c.indexes[dep.ObjOID]; idx != nil && !opclassIndexOIDs[idx.OID] {
						opclassIndexOIDs[idx.OID] = true
						amIndexes = append(amIndexes, ownerIndexExt{idx: idx, extName: extName})
					}
				case 'q', 'Q', 'm':
					if !inOwnedExt(dep.ObjOID) {
						key := [2]uint32{uint32(dep.ObjType), dep.ObjOID}
						if !amUserSeen[key] {
							amUserSeen[key] = true
							amUsers = append(amUsers, ownerAMUser{dep: dep, extName: extName})
						}
					}
				}
			}
		}
		if c.opFamilies[oid] != nil {
			for _, dep := range c.findNormalDependents('Q', oid) {
				if dep.ObjType == 'q' && !inOwnedExt(dep.ObjOID) {
					key := [2]uint32{uint32(dep.ObjType), dep.ObjOID}
					if !amUserSeen[key] {
						amUserSeen[key] = true
						amUsers = append(amUsers, ownerAMUser{dep: dep, extName: extName})
					}
				}
			}
		}
	}
	sort.Slice(amRelations, func(i, j int) bool { return amRelations[i].rel.OID < amRelations[j].rel.OID })
	sort.Slice(amIndexes, func(i, j int) bool { return amIndexes[i].idx.OID < amIndexes[j].idx.OID })
	sort.Slice(amUsers, func(i, j int) bool {
		if amUsers[i].dep.ObjType != amUsers[j].dep.ObjType {
			return amUsers[i].dep.ObjType < amUsers[j].dep.ObjType
		}
		return amUsers[i].dep.ObjOID < amUsers[j].dep.ObjOID
	})
	// externalCasts holds user-created casts that reference an owned
	// extension's type or function.
	type castExt struct {
		key     castKey
		extName string
	}
	var externalCasts []castExt
	extTypeOIDs := map[uint32]string{}
	extFuncOIDs := map[uint32]string{}
	for oid, extName := range c.extensionMembers {
		if !ownedExtNames[extName] {
			continue
		}
		if c.typeByOID[oid] != nil {
			extTypeOIDs[oid] = extName
		}
		if c.userProcs[oid] != nil {
			extFuncOIDs[oid] = extName
		}
	}
	for key, cast := range c.castIndex {
		if owner := c.extensionCasts[key]; owner != "" && ownedExtNames[owner] {
			continue // the cast goes with the extension that created it
		}
		if name := cmp.Or(extTypeOIDs[cast.Source], extTypeOIDs[cast.Target], extFuncOIDs[cast.Func]); name != "" {
			externalCasts = append(externalCasts, castExt{key: key, extName: name})
		}
	}
	sort.Slice(externalCasts, func(i, j int) bool {
		if externalCasts[i].key.source != externalCasts[j].key.source {
			return externalCasts[i].key.source < externalCasts[j].key.source
		}
		return externalCasts[i].key.target < externalCasts[j].key.target
	})
	sort.Slice(opclassIndexes, func(i, j int) bool { return opclassIndexes[i].idx.OID < opclassIndexes[j].idx.OID })

	type ownedRel struct {
		schema *Schema
		rel    *Relation
	}
	var rels []ownedRel
	ownedRelOIDs := map[uint32]bool{}
	for _, schema := range c.schemas {
		for _, rel := range schema.Relations {
			if owners[rel.OwnerName] || inOwnedExt(rel.OID) || autoExtOIDs[rel.OID] || (cascade && ownedSchemaOIDs[schema.OID]) {
				rels = append(rels, ownedRel{schema: schema, rel: rel})
				ownedRelOIDs[rel.OID] = true
			}
		}
	}
	sort.Slice(rels, func(i, j int) bool { return rels[i].rel.OID < rels[j].rel.OID })

	// Auto/internal relation dependents (such as partitions of an owned
	// partitioned table) are dropped automatically with their parent.
	for changed := true; changed; {
		changed = false
		for _, d := range c.deps {
			if d.ObjType == 'r' && d.RefType == 'r' && (d.DepType == DepAuto || d.DepType == DepInternal) &&
				ownedRelOIDs[d.RefOID] && !ownedRelOIDs[d.ObjOID] {
				if child := c.relationByOID[d.ObjOID]; child != nil {
					ownedRelOIDs[child.OID] = true
					rels = append(rels, ownedRel{schema: child.Schema, rel: child})
					changed = true
				}
			}
		}
	}
	sort.Slice(rels, func(i, j int) bool { return rels[i].rel.OID < rels[j].rel.OID })

	type ownedIdx struct {
		schema *Schema
		idx    *Index
	}
	var indexes []ownedIdx
	ownedIdxOIDs := map[uint32]bool{}
	type ownedSeq struct {
		schema *Schema
		seq    *Sequence
	}
	var seqs []ownedSeq
	ownedSeqOIDs := map[uint32]bool{}
	for _, schema := range c.schemas {
		for _, idx := range schema.Indexes {
			if ownedRelOIDs[idx.RelOID] {
				ownedIdxOIDs[idx.OID] = true
				continue
			}
			if owners[idx.OwnerName] || inOwnedExt(idx.OID) || autoExtOIDs[idx.OID] || (cascade && ownedSchemaOIDs[schema.OID]) {
				indexes = append(indexes, ownedIdx{schema: schema, idx: idx})
				ownedIdxOIDs[idx.OID] = true
			}
		}
		for _, seq := range schema.Sequences {
			if ownedRelOIDs[seq.OwnerRelOID] {
				ownedSeqOIDs[seq.OID] = true
				continue
			}
			if owners[seq.OwnerName] || inOwnedExt(seq.OID) || (cascade && ownedSchemaOIDs[schema.OID]) {
				seqs = append(seqs, ownedSeq{schema: schema, seq: seq})
				ownedSeqOIDs[seq.OID] = true
			}
		}
	}
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].idx.OID < indexes[j].idx.OID })
	sort.Slice(seqs, func(i, j int) bool { return seqs[i].seq.OID < seqs[j].seq.OID })
	for changed := true; changed; {
		changed = false
		for oid, child := range c.indexes {
			if child != nil && child.Schema != nil && child.ParentIndexOID != 0 && ownedIdxOIDs[child.ParentIndexOID] && !ownedIdxOIDs[oid] {
				ownedIdxOIDs[oid] = true
				if !ownedRelOIDs[child.RelOID] {
					indexes = append(indexes, ownedIdx{schema: child.Schema, idx: child})
				}
				changed = true
			}
		}
	}
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].idx.OID < indexes[j].idx.OID })
	for _, x := range indexes {
		if x.idx.ParentIndexOID != 0 && !ownedIdxOIDs[x.idx.ParentIndexOID] && !ownedRelOIDs[x.idx.RelOID] {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because parent index requires it", x.idx.Name),
			}
		}
	}

	var autoExtTriggers []*Trigger
	ownedTrigOIDs := map[uint32]bool{}
	for oid, trig := range c.triggers {
		if trig == nil {
			continue
		}
		if ownedRelOIDs[trig.RelOID] {
			ownedTrigOIDs[oid] = true
			continue
		}
		if inOwnedExt(oid) || autoExtOIDs[oid] {
			autoExtTriggers = append(autoExtTriggers, trig)
			ownedTrigOIDs[oid] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for oid, child := range c.triggers {
			if child != nil && child.ParentTriggerOID != 0 && ownedTrigOIDs[child.ParentTriggerOID] && !ownedTrigOIDs[oid] {
				ownedTrigOIDs[oid] = true
				if !ownedRelOIDs[child.RelOID] {
					autoExtTriggers = append(autoExtTriggers, child)
				}
				changed = true
			}
		}
	}
	sort.Slice(autoExtTriggers, func(i, j int) bool { return autoExtTriggers[i].OID < autoExtTriggers[j].OID })
	for _, trig := range autoExtTriggers {
		if trig.ParentTriggerOID != 0 && !ownedTrigOIDs[trig.ParentTriggerOID] && !ownedRelOIDs[trig.RelOID] {
			return errDependentObjects("trigger", trig.Name)
		}
	}

	var procs []*UserProc
	ownedProcOIDs := map[uint32]bool{}
	for _, proc := range c.userProcs {
		if owners[proc.OwnerName] || inOwnedExt(proc.OID) || autoExtOIDs[proc.OID] || (cascade && proc.Schema != nil && ownedSchemaOIDs[proc.Schema.OID]) {
			procs = append(procs, proc)
			ownedProcOIDs[proc.OID] = true
		}
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].OID < procs[j].OID })

	var types []*BuiltinType
	seenPrimaryType := map[uint32]bool{}
	ownedTypeOIDs := map[uint32]bool{}
	var staleTypeOwnerOIDs []uint32
	for oid, owner := range c.typeOwners {
		if !owners[owner] {
			continue
		}
		bt := c.typeByOID[oid]
		if bt == nil {
			staleTypeOwnerOIDs = append(staleTypeOwnerOIDs, oid)
			continue
		}
		ownedTypeOIDs[oid] = true
		if bt.Category == 'A' && bt.Elem != 0 || bt.Type == 'm' || (bt.Type == 'c' && bt.RelID != 0) {
			continue
		}
		if !seenPrimaryType[bt.OID] {
			seenPrimaryType[bt.OID] = true
			types = append(types, bt)
			for _, famOID := range c.typeFamilyOIDs(bt) {
				ownedTypeOIDs[famOID] = true
			}
		}
	}
	for _, bt := range c.typeByOID {
		if bt.OID < FirstNormalObjectId || bt.Elem != 0 || bt.Type == 'c' {
			continue
		}
		if !inOwnedExt(bt.OID) && !(cascade && ownedSchemaOIDs[bt.Namespace]) {
			continue
		}
		primary := bt
		if bt.Type == 'm' {
			if rng := c.owningRangeType(bt.OID); rng != nil {
				primary = rng
			}
		}
		if !seenPrimaryType[primary.OID] {
			seenPrimaryType[primary.OID] = true
			types = append(types, primary)
			for _, famOID := range c.typeFamilyOIDs(primary) {
				ownedTypeOIDs[famOID] = true
			}
		}
	}
	sort.Slice(types, func(i, j int) bool { return types[i].OID < types[j].OID })

	for _, x := range rels {
		if x.rel.RowTypeOID != 0 {
			ownedTypeOIDs[x.rel.RowTypeOID] = true
		}
		if x.rel.ArrayOID != 0 {
			ownedTypeOIDs[x.rel.ArrayOID] = true
		}
	}
	ownedDomainConOIDs := map[uint32]bool{}
	for oid := range ownedTypeOIDs {
		if dt := c.domainTypes[oid]; dt != nil {
			for _, dc := range dt.Constraints {
				ownedDomainConOIDs[dc.OID] = true
			}
		}
	}

	var ownedOperators []uint32
	seenOperator := map[uint32]bool{}
	for oid, owner := range c.operatorOwners {
		if owners[owner] && c.operByOID[oid] != nil && oid >= FirstNormalObjectId && !seenOperator[oid] {
			seenOperator[oid] = true
			ownedOperators = append(ownedOperators, oid)
		}
	}
	for oid := range c.operByOID {
		if seenOperator[oid] || oid < FirstNormalObjectId {
			continue
		}
		if inOwnedExt(oid) || (cascade && ownedSchemaOIDs[c.operatorSchemas[oid]]) {
			seenOperator[oid] = true
			ownedOperators = append(ownedOperators, oid)
		}
	}

	var ownedOpClasses []*OpClass
	ownedOpClassOIDs := map[uint32]bool{}
	for _, opc := range c.opClasses {
		if (opc.OwnerName != "" && owners[opc.OwnerName]) || inOwnedExt(opc.OID) || (cascade && ownedSchemaOIDs[opc.Schema]) {
			ownedOpClasses = append(ownedOpClasses, opc)
			ownedOpClassOIDs[opc.OID] = true
		}
	}
	sort.Slice(ownedOpClasses, func(i, j int) bool { return ownedOpClasses[i].OID < ownedOpClasses[j].OID })

	var ownedOpFamilies []*OpFamily
	ownedOpFamilyOIDs := map[uint32]bool{}
	for _, fam := range c.opFamilies {
		if (fam.OwnerName != "" && owners[fam.OwnerName]) || inOwnedExt(fam.OID) || (cascade && ownedSchemaOIDs[fam.Schema]) {
			ownedOpFamilies = append(ownedOpFamilies, fam)
			ownedOpFamilyOIDs[fam.OID] = true
		}
	}
	sort.Slice(ownedOpFamilies, func(i, j int) bool { return ownedOpFamilies[i].OID < ownedOpFamilies[j].OID })

	// Event triggers are collected before preflight so a function they call
	// and they themselves both count as part of one DROP OWNED set.
	var ownedEventTriggers []*EventTrigger
	ownedEventTriggerOIDs := map[uint32]bool{}
	for _, evt := range c.eventTriggers {
		if (evt.OwnerName != "" && owners[evt.OwnerName]) || inOwnedExt(evt.OID) {
			ownedEventTriggers = append(ownedEventTriggers, evt)
			ownedEventTriggerOIDs[evt.OID] = true
		}
	}
	sort.Slice(ownedEventTriggers, func(i, j int) bool { return ownedEventTriggers[i].OID < ownedEventTriggers[j].OID })

	type ownedColl struct {
		schema *Schema
		uc     *UserCollation
	}
	var ownedCollations []ownedColl
	ownedCollationOIDs := map[uint32]bool{}
	for _, schema := range c.schemas {
		for _, uc := range schema.Collations {
			if uc == nil {
				continue
			}
			if (uc.OwnerName != "" && owners[uc.OwnerName]) || inOwnedExt(uc.OID) || (cascade && ownedSchemaOIDs[schema.OID]) {
				ownedCollations = append(ownedCollations, ownedColl{schema: schema, uc: uc})
				ownedCollationOIDs[uc.OID] = true
			}
		}
	}
	sort.Slice(ownedCollations, func(i, j int) bool { return ownedCollations[i].uc.OID < ownedCollations[j].uc.OID })

	// Refuse dropping any object pinned to an un-dropped extension (DEPENDENCY_EXTENSION).
	checkExtPin := func(oid uint32) error {
		if c.isExtensionMember(oid) && !inOwnedExt(oid) {
			return errDependentObjects("extension", c.extensionMembers[oid])
		}
		return nil
	}
	for _, x := range rels {
		if err := checkExtPin(x.rel.OID); err != nil {
			return err
		}
	}
	for _, x := range indexes {
		if err := checkExtPin(x.idx.OID); err != nil {
			return err
		}
	}
	for _, x := range seqs {
		if err := checkExtPin(x.seq.OID); err != nil {
			return err
		}
	}
	for _, proc := range procs {
		if err := checkExtPin(proc.OID); err != nil {
			return err
		}
	}
	for _, bt := range types {
		if err := checkExtPin(bt.OID); err != nil {
			return err
		}
	}
	for _, oid := range ownedOperators {
		if err := checkExtPin(oid); err != nil {
			return err
		}
	}
	for _, opc := range ownedOpClasses {
		if err := checkExtPin(opc.OID); err != nil {
			return err
		}
	}
	for _, fam := range ownedOpFamilies {
		if err := checkExtPin(fam.OID); err != nil {
			return err
		}
	}
	for _, evt := range ownedEventTriggers {
		if err := checkExtPin(evt.OID); err != nil {
			return err
		}
	}
	for _, x := range ownedCollations {
		if err := checkExtPin(x.uc.OID); err != nil {
			return err
		}
	}
	for _, s := range ownedSchemas {
		if err := checkExtPin(s.OID); err != nil {
			return err
		}
	}

	// Collect non-extension casts referencing any owned type or function.
	var dependentCasts []castKey
	for key, cast := range c.castIndex {
		if owner := c.extensionCasts[key]; owner != "" && ownedExtNames[owner] {
			continue
		}
		if ownedTypeOIDs[cast.Source] || ownedTypeOIDs[cast.Target] || (cast.Func != 0 && ownedProcOIDs[cast.Func]) {
			dependentCasts = append(dependentCasts, key)
		}
	}
	sort.Slice(dependentCasts, func(i, j int) bool {
		if dependentCasts[i].source != dependentCasts[j].source {
			return dependentCasts[i].source < dependentCasts[j].source
		}
		return dependentCasts[i].target < dependentCasts[j].target
	})

	inDropSet := func(objType byte, objOID uint32) bool {
		switch objType {
		case 'r':
			return ownedRelOIDs[objOID]
		case 'i':
			if ownedIdxOIDs[objOID] {
				return true
			}
			idx := c.indexes[objOID]
			return idx != nil && ownedRelOIDs[idx.RelOID]
		case 's':
			if ownedSeqOIDs[objOID] {
				return true
			}
			seq := c.sequenceByOID[objOID]
			return seq != nil && ownedRelOIDs[seq.OwnerRelOID]
		case 'f':
			return ownedProcOIDs[objOID]
		case 't':
			return ownedTypeOIDs[objOID]
		case 'c':
			if ownedDomainConOIDs[objOID] {
				return true
			}
			con := c.constraints[objOID]
			return con != nil && ownedRelOIDs[con.RelOID]
		case 'g':
			if ownedTrigOIDs[objOID] {
				return true
			}
			trig := c.triggers[objOID]
			return trig != nil && ownedRelOIDs[trig.RelOID]
		case 'p':
			pol := c.policies[objOID]
			return pol != nil && ownedRelOIDs[pol.RelOID]
		case 'R':
			rule := c.rules[objOID]
			return rule != nil && ownedRelOIDs[rule.OwnerRelOID]
		case 'E':
			return ownedEventTriggerOIDs[objOID]
		case 'O':
			return seenOperator[objOID] || inOwnedExt(objOID)
		case 'm':
			return inOwnedExt(objOID)
		case 'q':
			return ownedOpClassOIDs[objOID]
		case 'Q':
			return ownedOpFamilyOIDs[objOID]
		case 'L':
			return ownedCollationOIDs[objOID]
		}
		return false
	}

	// Preflight all external dependencies before mutating any catalog state.
	if !cascade {
		for _, s := range ownedSchemas {
			for _, rel := range s.Relations {
				if !ownedRelOIDs[rel.OID] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, idx := range s.Indexes {
				if !inDropSet('i', idx.OID) {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, seq := range s.Sequences {
				if !inDropSet('s', seq.OID) {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, up := range c.userProcs {
				if up.Schema != nil && up.Schema.OID == s.OID && !ownedProcOIDs[up.OID] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, bt := range c.typeByOID {
				if bt.Namespace == s.OID && bt.OID >= FirstNormalObjectId && bt.Elem == 0 && bt.Type != 'c' && !ownedTypeOIDs[bt.OID] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, ext := range c.extensions {
				if ext.SchemaOID == s.OID && !ownedExtNames[ext.Name] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for oid, nsOID := range c.operatorSchemas {
				if nsOID == s.OID && !inDropSet('O', oid) && c.operByOID[oid] != nil {
					return errDependentObjects("schema", s.Name)
				}
			}
			for oid, opc := range c.opClasses {
				if opc.Schema == s.OID && !inDropSet('q', oid) {
					return errDependentObjects("schema", s.Name)
				}
			}
			for oid, fam := range c.opFamilies {
				if fam.Schema == s.OID && !inDropSet('Q', oid) {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, uc := range s.Collations {
				if uc != nil && !inDropSet('L', uc.OID) {
					return errDependentObjects("schema", s.Name)
				}
			}
		}
		for _, x := range rels {
			for _, dep := range c.relationDependents(x.rel.OID, 0) {
				if !dep.Normal || inDropSet(dep.ObjType, dep.ObjOID) {
					continue
				}
				return errDependentObjects(dropRelKindLabel(x.rel.RelKind), x.rel.Name)
			}
			if x.rel.RowTypeOID != 0 {
				if bt := c.typeByOID[x.rel.RowTypeOID]; bt != nil {
					for _, depRelOID := range c.typeRelationDependents(bt) {
						if !ownedRelOIDs[depRelOID] {
							return errDependentObjects(dropRelKindLabel(x.rel.RelKind), x.rel.Name)
						}
					}
					for _, dep := range c.typeNormalDependents(bt) {
						if !inDropSet(dep.ObjType, dep.ObjOID) {
							return errDependentObjects(dropRelKindLabel(x.rel.RelKind), x.rel.Name)
						}
					}
				}
			}
		}
		for _, x := range indexes {
			if x.idx.ConstraintOID != 0 && !inDropSet('c', x.idx.ConstraintOID) {
				return &Error{
					Code:    CodeDependentObjects,
					Message: fmt.Sprintf("cannot drop index %q because constraint requires it", x.idx.Name),
				}
			}
			for _, refClass := range []byte{'r', 'i'} {
				for _, dep := range c.findNormalDependents(refClass, x.idx.OID) {
					if !inDropSet(dep.ObjType, dep.ObjOID) {
						return errDependentObjects("index", x.idx.Name)
					}
				}
			}
		}
		for _, trig := range autoExtTriggers {
			for _, dep := range c.findNormalDependents('g', trig.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("trigger", trig.Name)
				}
			}
		}
		for _, x := range opclassIndexes {
			if !ownedRelOIDs[x.idx.RelOID] && !inDropSet('i', x.idx.OID) {
				return errDependentObjects("extension", x.extName)
			}
		}
		for _, x := range amRelations {
			if !ownedRelOIDs[x.rel.OID] {
				return errDependentObjects("extension", x.extName)
			}
		}
		for _, x := range amIndexes {
			if !ownedRelOIDs[x.idx.RelOID] && !inDropSet('i', x.idx.OID) {
				return errDependentObjects("extension", x.extName)
			}
		}
		for _, x := range amUsers {
			if inDropSet(x.dep.ObjType, x.dep.ObjOID) {
				continue
			}
			return errDependentObjects("extension", x.extName)
		}
		for _, x := range externalCasts {
			return errDependentObjects("extension", x.extName)
		}
		if len(dependentCasts) > 0 {
			return errDependentObjects("cast", "")
		}
		for _, x := range seqs {
			for _, dep := range c.findNormalDependents('s', x.seq.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("sequence", x.seq.Name)
				}
			}
		}
		for _, proc := range procs {
			for _, dep := range c.findNormalDependents('f', proc.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					kind := "function"
					if proc.Kind == 'p' {
						kind = "procedure"
					}
					return errDependentObjects(kind, proc.Name)
				}
			}
		}
		for _, opOID := range ownedOperators {
			for _, dep := range c.findNormalDependents('O', opOID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					opName := ""
					if op := c.operByOID[opOID]; op != nil {
						opName = op.Name
					}
					return errDependentObjects("operator", opName)
				}
			}
		}
		for _, opc := range ownedOpClasses {
			for _, refClass := range []byte{'o', 'q'} {
				for _, dep := range c.findNormalDependents(refClass, opc.OID) {
					if !inDropSet(dep.ObjType, dep.ObjOID) {
						return errDependentObjects("operator class", opc.Name)
					}
				}
			}
		}
		for _, fam := range ownedOpFamilies {
			for _, dep := range c.findNormalDependents('Q', fam.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("operator family", fam.Name)
				}
			}
		}
		for _, x := range ownedCollations {
			for _, dep := range c.findNormalDependents('L', x.uc.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("collation", x.uc.Name)
				}
			}
		}
		// Revoking a grant option the dropped role holds abandons the grants it
		// passed on: RESTRICT refuses that, CASCADE revokes the chain.
		grantPrivCovers := func(held, passed string) bool {
			h := strings.ToUpper(held)
			p := strings.ToUpper(passed)
			return h == "ALL" || p == "ALL" || h == p
		}
		for _, g := range c.grants {
			if !owners[g.Grantee] || !g.WithGrant {
				continue
			}
			if inDropSet(g.ObjType, g.ObjOID) || (g.ObjType == 'n' && ownedSchemaOIDs[g.ObjOID]) || (g.ObjType == 'T' && ownedTypeOIDs[g.ObjOID]) {
				continue
			}
			for _, dependent := range c.grants {
				if dependent.ObjType == g.ObjType && dependent.ObjOID == g.ObjOID &&
					dependent.Grantor == g.Grantee && !owners[dependent.Grantee] &&
					grantPrivCovers(g.Privilege, dependent.Privilege) {
					return errDependentObjects("grant", g.Grantee)
				}
			}
		}
		for _, bt := range types {
			for _, depRelOID := range c.typeRelationDependents(bt) {
				if !ownedRelOIDs[depRelOID] {
					return errDependentObjects("type", bt.TypeName)
				}
			}
			for _, dep := range c.typeNormalDependents(bt) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("type", bt.TypeName)
				}
			}
		}
	}

	restore := func() {}
	if cascade && len(ownedCollations) > 0 {
		restore = c.snapshotForAlterRollback()
	}

	// Execute the drop set now that all checks passed.
	for _, x := range rels {
		if c.relationByOID[x.rel.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropRelation(x.rel)
			continue
		}
		if hasNormalDependent(c.relationDependents(x.rel.OID, 0)) {
			c.dropDependents('r', x.rel.OID)
		}
		if c.relationByOID[x.rel.OID] != nil {
			c.removeRelation(x.schema, x.rel.Name, x.rel)
		}
	}
	for _, x := range indexes {
		if x.schema.Indexes[x.idx.Name] == x.idx {
			if cascade {
				c.dropDependentObject('i', x.idx.OID)
			} else {
				c.removeIndex(x.schema, x.idx.Name, x.idx)
			}
		}
	}
	for _, x := range seqs {
		if c.sequenceByOID[x.seq.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('s', x.seq.OID)) > 0 {
			c.dropDependents('s', x.seq.OID)
		}
		if c.sequenceByOID[x.seq.OID] != nil {
			c.removeSequence(x.schema, x.seq)
		}
	}
	for _, proc := range procs {
		if c.userProcs[proc.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('f', proc.OID)) > 0 {
			c.dropDependents('f', proc.OID)
		}
		if c.userProcs[proc.OID] != nil {
			c.removeFunction(proc.OID, proc.Name)
		}
	}
	for _, bt := range types {
		if c.typeByOID[bt.OID] == nil {
			continue
		}
		if cascade {
			for _, depRelOID := range c.typeRelationDependents(bt) {
				if depRel := c.relationByOID[depRelOID]; depRel != nil {
					c.cascadeDropRelation(depRel)
				}
			}
			c.dropTypeDependents(bt)
		}
		if c.typeByOID[bt.OID] != nil {
			c.deleteUserType(bt)
		}
	}
	for _, oid := range staleTypeOwnerOIDs {
		delete(c.typeOwners, oid)
	}
	for oid, owner := range c.typeOwners {
		if owners[owner] && c.typeByOID[oid] == nil {
			delete(c.typeOwners, oid)
		}
	}
	for _, trig := range autoExtTriggers {
		if c.triggers[trig.OID] == nil {
			continue
		}
		if cascade {
			c.dropDependents('g', trig.OID)
		}
		c.removeTrigger(trig)
	}
	for _, ext := range ownedExts {
		if c.extensions[ext.OID] != nil {
			c.removeExtensionMetadata(ext)
		}
	}
	for _, x := range opclassIndexes {
		if c.indexes[x.idx.OID] == x.idx {
			c.dropDependentObject('i', x.idx.OID)
		}
	}
	if cascade {
		for _, x := range amRelations {
			if c.relationByOID[x.rel.OID] != nil {
				c.cascadeDropRelation(x.rel)
			}
		}
		for _, x := range amIndexes {
			if c.indexes[x.idx.OID] == x.idx {
				c.dropDependentObject('i', x.idx.OID)
			}
		}
		for _, x := range amUsers {
			switch x.dep.ObjType {
			case 'm':
				c.cascadeDropAccessMethod(x.dep.ObjOID)
			case 'q':
				c.cascadeDropOpClass(x.dep.ObjOID)
			case 'Q':
				c.cascadeDropOpFamily(x.dep.ObjOID)
			}
		}
		for _, x := range externalCasts {
			delete(c.castIndex, x.key)
			delete(c.extensionCasts, x.key)
		}
		for _, key := range dependentCasts {
			delete(c.castIndex, key)
			delete(c.extensionCasts, key)
		}
	}
	sort.Slice(ownedOperators, func(i, j int) bool { return ownedOperators[i] < ownedOperators[j] })
	for _, oid := range ownedOperators {
		if cascade {
			c.dropDependents('O', oid)
		}
		c.removeOperator(oid)
	}
	for _, opc := range ownedOpClasses {
		if c.opClasses[opc.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropOpClass(opc.OID)
		} else {
			c.removeOpClass(opc.OID)
		}
	}
	for _, fam := range ownedOpFamilies {
		if c.opFamilies[fam.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropOpFamily(fam.OID)
		} else {
			c.removeOpFamily(fam.OID)
		}
	}
	for _, s := range ownedSchemas {
		opcOIDs := make([]uint32, 0, len(c.opClasses))
		for oid, opc := range c.opClasses {
			if opc.Schema == s.OID {
				opcOIDs = append(opcOIDs, oid)
			}
		}
		sort.Slice(opcOIDs, func(i, j int) bool { return opcOIDs[i] < opcOIDs[j] })
		for _, oid := range opcOIDs {
			c.cascadeDropOpClass(oid)
		}
		famOIDs := make([]uint32, 0, len(c.opFamilies))
		for oid, fam := range c.opFamilies {
			if fam.Schema == s.OID {
				famOIDs = append(famOIDs, oid)
			}
		}
		sort.Slice(famOIDs, func(i, j int) bool { return famOIDs[i] < famOIDs[j] })
		for _, oid := range famOIDs {
			c.cascadeDropOpFamily(oid)
		}
	}
	for _, evt := range ownedEventTriggers {
		if c.eventTriggers[evt.OID] == evt {
			c.removeEventTrigger(evt)
		}
	}
	for _, x := range ownedCollations {
		if x.schema.Collations[x.uc.Name] != x.uc {
			continue
		}
		if cascade {
			if err := c.cascadeDropCollation(x.schema, x.uc.Name, x.uc.OID); err != nil {
				restore()
				return err
			}
		} else {
			delete(x.schema.Collations, x.uc.Name)
			c.removeComments('L', x.uc.OID)
			c.removeDepsOn('L', x.uc.OID)
		}
	}
	for _, s := range ownedSchemas {
		if c.schemas[s.OID] == nil {
			continue
		}
		c.removeGrantsForObject('n', s.OID)
		c.removeComments('n', s.OID)
		c.removeDepsOf('n', s.OID)
		c.removeDepsOn('n', s.OID)
		delete(c.schemas, s.OID)
		delete(c.schemaByName, s.Name)
	}

	if len(c.grants) > 0 {
		grantPrivCovers := func(held, passed string) bool {
			h := strings.ToUpper(held)
			p := strings.ToUpper(passed)
			return h == "ALL" || p == "ALL" || h == p
		}
		removed := make([]bool, len(c.grants))
		for i, g := range c.grants {
			if owners[g.Grantee] || owners[g.Grantor] {
				removed[i] = true
			}
		}
		for changed := true; changed; {
			changed = false
			for i, g := range c.grants {
				if removed[i] || g.Grantor == "" || c.objectOwnerName(g.ObjType, g.ObjOID) == g.Grantor {
					continue
				}
				lostOption := false
				stillHasOption := false
				for j, src := range c.grants {
					if src.ObjType != g.ObjType || src.ObjOID != g.ObjOID || src.Grantee != g.Grantor || !src.WithGrant {
						continue
					}
					if !grantPrivCovers(src.Privilege, g.Privilege) {
						continue
					}
					if removed[j] {
						lostOption = true
					} else {
						stillHasOption = true
					}
				}
				if lostOption && !stillHasOption {
					removed[i] = true
					changed = true
				}
			}
		}
		kept := c.grants[:0]
		for i, g := range c.grants {
			if !removed[i] {
				kept = append(kept, g)
			}
		}
		c.grants = kept
	}

	var policiesToRemove []*Policy
	for _, p := range c.policies {
		var remaining []string
		removedAny := false
		for _, r := range p.Roles {
			if r != "" && r != "public" && owners[r] {
				removedAny = true
				continue
			}
			remaining = append(remaining, r)
		}
		if !removedAny {
			continue
		}
		if len(remaining) == 0 {
			policiesToRemove = append(policiesToRemove, p)
		} else {
			p.Roles = remaining
		}
	}
	for _, p := range policiesToRemove {
		c.removePolicy(p)
	}
	return nil
}

func dropRelKindLabel(relkind byte) string {
	switch relkind {
	case 'v':
		return "view"
	case 'm':
		return "materialized view"
	case 'c':
		return "type"
	case 'f':
		return "foreign table"
	default:
		return "table"
	}
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
		if schema.OID != PGCatalogNamespace && schema.OID != PGToastNamespace {
			if so := schema.effectiveOwner(); so != "" && owners[so] {
				schema.OwnerName = newOwner
				if schema.Owner != "" {
					schema.Owner = newOwner
				}
				c.transferACLs('n', schema.OID, so, newOwner)
			}
		}
		for _, rel := range schema.Relations {
			if owners[rel.OwnerName] {
				oldOwner := rel.OwnerName
				rel.OwnerName = newOwner
				c.transferACLs('r', rel.OID, oldOwner, newOwner)
				if rel.RowTypeOID != 0 {
					c.transferACLs('T', rel.RowTypeOID, oldOwner, newOwner)
				}
				if rel.ArrayOID != 0 {
					c.transferACLs('T', rel.ArrayOID, oldOwner, newOwner)
				}
			}
		}
		for _, idx := range schema.Indexes {
			if owners[idx.OwnerName] {
				idx.OwnerName = newOwner
			}
		}
		for _, seq := range schema.Sequences {
			if owners[seq.OwnerName] {
				oldOwner := seq.OwnerName
				seq.OwnerName = newOwner
				c.transferACLs('s', seq.OID, oldOwner, newOwner)
			}
		}
		for _, uc := range schema.Collations {
			if uc != nil && owners[uc.OwnerName] {
				uc.OwnerName = newOwner
			}
		}
	}
	for _, proc := range c.userProcs {
		if owners[proc.OwnerName] {
			oldOwner := proc.OwnerName
			proc.OwnerName = newOwner
			c.transferACLs('f', proc.OID, oldOwner, newOwner)
		}
	}
	for oid, owner := range c.typeOwners {
		if owners[owner] {
			c.typeOwners[oid] = newOwner
			c.transferACLs('T', oid, owner, newOwner)
		}
	}
	for oid, owner := range c.operatorOwners {
		if owners[owner] {
			c.operatorOwners[oid] = newOwner
		}
	}
	for _, opc := range c.opClasses {
		if owners[opc.OwnerName] {
			opc.OwnerName = newOwner
		}
	}
	for _, fam := range c.opFamilies {
		if owners[fam.OwnerName] {
			fam.OwnerName = newOwner
		}
	}
	for name, owner := range c.tablespaceOwners {
		if owners[owner] {
			c.tablespaceOwners[name] = newOwner
		}
	}
	for name, owner := range c.databaseOwners {
		if owners[owner] {
			c.databaseOwners[name] = newOwner
		}
	}
	for _, ext := range c.extensions {
		if owners[ext.OwnerName] {
			ext.OwnerName = newOwner
		}
	}
	for _, evt := range c.eventTriggers {
		if owners[evt.OwnerName] {
			evt.OwnerName = newOwner
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

// RenameRole renames a role and every object, grant, or policy referencing it.
//
// pg: src/backend/commands/user.c — RenameRole
func (c *Catalog) RenameRole(stmt *nodes.RenameStmt) error {
	oldName := stmt.Subname
	if err := c.checkRole(oldName); err != nil {
		return err
	}
	if oldName == c.currentOwner() || oldName == c.currentSessionUser() {
		return &Error{Code: CodeFeatureNotSupported, Message: "current user cannot be renamed"}
	}
	if strings.HasPrefix(oldName, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("role name %q is reserved", oldName)}
	}
	if strings.HasPrefix(stmt.Newname, "pg_") {
		return &Error{Code: CodeReservedName,
			Message: fmt.Sprintf("role name %q is reserved", stmt.Newname)}
	}
	if c.roles[stmt.Newname] || oldName == stmt.Newname {
		return errDuplicateObject("role", stmt.Newname)
	}
	if c.snapshot {
		c.roles[oldName] = false
	} else {
		delete(c.roles, oldName)
	}
	c.roles[stmt.Newname] = true
	for _, schema := range c.schemas {
		if schema.OID != PGCatalogNamespace && schema.OID != PGToastNamespace {
			if schema.Owner == oldName {
				schema.Owner = stmt.Newname
			}
			if schema.OwnerName == oldName {
				schema.OwnerName = stmt.Newname
			}
		}
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
		for _, uc := range schema.Collations {
			if uc != nil && uc.OwnerName == oldName {
				uc.OwnerName = stmt.Newname
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
	for oid, owner := range c.operatorOwners {
		if owner == oldName {
			c.operatorOwners[oid] = stmt.Newname
		}
	}
	for name, owner := range c.tablespaceOwners {
		if owner == oldName {
			c.tablespaceOwners[name] = stmt.Newname
		}
	}
	for name, owner := range c.databaseOwners {
		if owner == oldName {
			c.databaseOwners[name] = stmt.Newname
		}
	}
	for _, ext := range c.extensions {
		if ext.OwnerName == oldName {
			ext.OwnerName = stmt.Newname
		}
	}
	for _, evt := range c.eventTriggers {
		if evt.OwnerName == oldName {
			evt.OwnerName = stmt.Newname
		}
	}
	for _, opc := range c.opClasses {
		if opc.OwnerName == oldName {
			opc.OwnerName = stmt.Newname
		}
	}
	for _, fam := range c.opFamilies {
		if fam.OwnerName == oldName {
			fam.OwnerName = stmt.Newname
		}
	}
	for i := range c.grants {
		if c.grants[i].Grantee == oldName {
			c.grants[i].Grantee = stmt.Newname
		}
		if c.grants[i].Grantor == oldName {
			c.grants[i].Grantor = stmt.Newname
		}
	}
	for _, p := range c.policies {
		for i, r := range p.Roles {
			if r == oldName {
				p.Roles[i] = stmt.Newname
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
	var names, missingNames []string
	seen := map[string]bool{}
	for _, item := range stmt.Roles.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		name := c.ownerRoleName(rs)
		if seen[name] {
			if stmt.MissingOk {
				continue
			}
			return errUndefinedObject("role", name)
		}
		if err := c.checkRole(name); err != nil {
			if stmt.MissingOk {
				seen[name] = true
				missingNames = append(missingNames, name)
				continue
			}
			return err
		}
		if name == c.currentOwner() || name == c.currentSessionUser() {
			return &Error{Code: CodeObjectNotInPrerequisiteState, Message: "current user cannot be dropped"}
		}
		if c.roleInUse(name) {
			return errDependentObjects("role", name)
		}
		seen[name] = true
		names = append(names, name)
	}
	for _, name := range missingNames {
		if c.snapshot {
			c.roles[name] = false
		}
	}
	for _, name := range names {
		if c.snapshot {
			c.roles[name] = false
		} else {
			delete(c.roles, name)
		}
	}
	return nil
}

func (c *Catalog) currentSessionUser() string {
	if c.sessionUser != "" {
		return c.sessionUser
	}
	return "postgres"
}

// ownerRoleName resolves a role spec to the role name an object's owner is
// recorded under.
func (c *Catalog) ownerRoleName(rs *nodes.RoleSpec) string {
	name := roleSpecName(rs)
	switch name {
	case "CURRENT_USER", "CURRENT_ROLE":
		return c.currentOwner()
	case "SESSION_USER":
		return c.currentSessionUser()
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
