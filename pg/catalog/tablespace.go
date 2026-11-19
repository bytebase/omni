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
	if c.defaultTablespace == "" || c.snapshot || c.tablespaces[c.defaultTablespace] {
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
	delete(c.tablespaces, name)
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
	case nodes.OBJECT_EVENT_TRIGGER:
		return c.alterEventTriggerOwner(stmt)
	case nodes.OBJECT_EXTENSION:
		return c.alterExtensionOwner(stmt)
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
	schemaName, oprName := qualifiedName(owa.Objname)
	var opSchema *Schema
	if schemaName != "" {
		var err error
		opSchema, err = c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
	}
	var left, right uint32
	if owa.Objargs != nil {
		var oids []uint32
		for _, item := range owa.Objargs.Items {
			tn, ok := item.(*nodes.TypeName)
			if !ok {
				continue
			}
			oid, _, err := c.resolveTypeName(tn)
			if err != nil {
				return err
			}
			oids = append(oids, oid)
		}
		if len(oids) == 1 {
			right = oids[0]
		} else if len(oids) >= 2 {
			left, right = oids[0], oids[1]
		}
	}
	var found *BuiltinOperator
	rank := func(op *BuiltinOperator) int {
		if op.OID < FirstNormalObjectId {
			return -1
		}
		nsOID, ok := c.operatorSchemas[op.OID]
		if !ok {
			return -1
		}
		if opSchema != nil {
			if opSchema.OID == nsOID {
				return 0
			}
			return -1
		}
		for i, visible := range c.searchPathWithCatalog() {
			if visible == nsOID {
				return i
			}
		}
		return -1
	}
	bestRank := -1
	for _, op := range c.operByKey[operKey{name: oprName, left: left, right: right}] {
		if r := rank(op); r >= 0 && (bestRank < 0 || r < bestRank) {
			bestRank = r
			found = op
		}
	}
	if found == nil {
		return errUndefinedObject("operator", oprName)
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
	s.Owner = owner
	s.OwnerName = owner
	return nil
}

func (c *Catalog) alterTypeOwner(stmt *nodes.AlterOwnerStmt) error {
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
			rel.OwnerName = owner
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
		c.typeOwners[oid] = owner
	}
	return nil
}

func (c *Catalog) alterRoutineOwner(stmt *nodes.AlterOwnerStmt) error {
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
		for _, nsOID := range c.searchPathWithCatalog() {
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
	up.OwnerName = owner
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
	for _, item := range cmds.Items {
		atc, ok := item.(*nodes.AlterTableCmd)
		if !ok {
			continue
		}
		switch nodes.AlterTableType(atc.Subtype) {
		case nodes.AT_SetTableSpace:
			if err := c.checkRelationTablespace(atc.Name); err != nil {
				return err
			}
			return errWrongObjectType(relName, "a table, materialized view, or index")
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
		}
	}
	if hasOwner {
		seq.OwnerName = newOwner
	}
	return nil
}

func allAlterCmdsChangeOwner(cmds *nodes.List) bool {
	if cmds == nil || len(cmds.Items) == 0 {
		return false
	}
	for _, item := range cmds.Items {
		atc, ok := item.(*nodes.AlterTableCmd)
		if !ok || nodes.AlterTableType(atc.Subtype) != nodes.AT_ChangeOwner {
			return false
		}
	}
	return true
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
	if c.tablespaces[stmt.Newname] {
		return errDuplicateObject("tablespace", stmt.Newname)
	}
	delete(c.tablespaces, oldName)
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
	return false
}

// ownerAMUser is an operator class or family that uses an access method of an
// owned extension.
type ownerAMUser struct {
	dep     DepEntry
	extName string
}

// removeAMUser removes the operator class or family a dependency entry names.
func (c *Catalog) removeAMUser(dep DepEntry) {
	switch dep.ObjType {
	case 'q':
		if opc := c.opClasses[dep.ObjOID]; opc != nil {
			if opc.IsDefault {
				delete(c.opClassByKey, opClassKey{amOID: opc.AMOID, typeOID: opc.TypeOID})
			}
			delete(c.opClasses, opc.OID)
		}
	case 'Q':
		delete(c.opFamilies, dep.ObjOID)
	}
}

// ownerIndexExt pairs an index with the extension whose operator class it
// uses, for the DROP OWNED report.
type ownerIndexExt struct {
	idx     *Index
	extName string
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
	// opclassIndexes holds the indexes that use an operator class of an owned
	// extension: a RESTRICT drop reports them, CASCADE takes them with it.
	var opclassIndexes []ownerIndexExt
	opclassIndexOIDs := map[uint32]bool{}
	// amRelations holds relations that use an access method of an owned
	// extension; CASCADE drops them with it.
	var amRelations []*Relation
	amRelOIDs := map[uint32]bool{}
	var amIndexes []ownerIndexExt
	var amUsers []ownerAMUser
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
		if c.accessMethods[oid] == nil {
			continue
		}
		for _, dep := range c.findNormalDependents('m', oid) {
			switch dep.ObjType {
			case 'r':
				if rel := c.relationByOID[dep.ObjOID]; rel != nil && !amRelOIDs[rel.OID] {
					amRelOIDs[rel.OID] = true
					amRelations = append(amRelations, rel)
				}
			case 'i':
				if idx := c.indexes[dep.ObjOID]; idx != nil && !opclassIndexOIDs[idx.OID] {
					opclassIndexOIDs[idx.OID] = true
					amIndexes = append(amIndexes, ownerIndexExt{idx: idx, extName: extName})
				}
			}
		}
		for _, dep := range c.findNormalDependents('m', oid) {
			if dep.ObjType == 'q' || dep.ObjType == 'Q' {
				amUsers = append(amUsers, ownerAMUser{dep: dep, extName: extName})
			}
		}
	}
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
			if owners[rel.OwnerName] || inOwnedExt(rel.OID) || (cascade && ownedSchemaOIDs[schema.OID]) {
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
			if owners[idx.OwnerName] || inOwnedExt(idx.OID) || (cascade && ownedSchemaOIDs[schema.OID]) {
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

	var procs []*UserProc
	ownedProcOIDs := map[uint32]bool{}
	for _, proc := range c.userProcs {
		if owners[proc.OwnerName] || inOwnedExt(proc.OID) || (cascade && proc.Schema != nil && ownedSchemaOIDs[proc.Schema.OID]) {
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
		if bt.OID < FirstNormalObjectId || bt.Elem != 0 || bt.Type == 'm' || bt.Type == 'c' {
			continue
		}
		if !inOwnedExt(bt.OID) && !(cascade && ownedSchemaOIDs[bt.Namespace]) {
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
	if cascade {
		for oid := range c.operByOID {
			if seenOperator[oid] || oid < FirstNormalObjectId {
				continue
			}
			if nsOID, ok := c.operatorSchemas[oid]; ok && ownedSchemaOIDs[nsOID] {
				seenOperator[oid] = true
				ownedOperators = append(ownedOperators, oid)
			}
		}
	}

	// Event triggers are collected before preflight so a function they call
	// and they themselves both count as part of one DROP OWNED set.
	var ownedEventTriggers []*EventTrigger
	ownedEventTriggerOIDs := map[uint32]bool{}
	for _, evt := range c.eventTriggers {
		if owners[evt.OwnerName] {
			ownedEventTriggers = append(ownedEventTriggers, evt)
			ownedEventTriggerOIDs[evt.OID] = true
		}
	}
	sort.Slice(ownedEventTriggers, func(i, j int) bool { return ownedEventTriggers[i].OID < ownedEventTriggers[j].OID })

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
			trig := c.triggers[objOID]
			return trig != nil && ownedRelOIDs[trig.RelOID]
		case 'p':
			pol := c.policies[objOID]
			return pol != nil && ownedRelOIDs[pol.RelOID]
		case 'E':
			return ownedEventTriggerOIDs[objOID]
		case 'O':
			return inOwnedExt(objOID) || c.operatorOwners[objOID] != "" && owners[c.operatorOwners[objOID]]
		case 'm':
			return inOwnedExt(objOID)
		case 'q':
			if inOwnedExt(objOID) {
				return true
			}
			if opc := c.opClasses[objOID]; opc != nil && inOwnedExt(opc.AMOID) {
				return true
			}
			return false
		case 'Q':
			if inOwnedExt(objOID) {
				return true
			}
			if fam := c.opFamilies[objOID]; fam != nil && inOwnedExt(fam.AMOID) {
				return true
			}
			return false
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
				if bt.Namespace == s.OID && bt.OID >= FirstNormalObjectId && bt.Elem == 0 && bt.Type != 'm' && bt.Type != 'c' && !ownedTypeOIDs[bt.OID] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for _, ext := range c.extensions {
				if ext.SchemaOID == s.OID && !ownedExtNames[ext.Name] {
					return errDependentObjects("schema", s.Name)
				}
			}
			for oid, nsOID := range c.operatorSchemas {
				if nsOID == s.OID && !seenOperator[oid] && c.operByOID[oid] != nil {
					return errDependentObjects("schema", s.Name)
				}
			}
			if !cascade {
				for _, opc := range c.opClasses {
					if opc.Schema == s.OID {
						return errDependentObjects("schema", s.Name)
					}
				}
				for _, fam := range c.opFamilies {
					if fam.Schema == s.OID {
						return errDependentObjects("schema", s.Name)
					}
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
		}
		for _, x := range opclassIndexes {
			if !ownedRelOIDs[x.idx.RelOID] && !inDropSet('i', x.idx.OID) {
				return errDependentObjects("extension", x.extName)
			}
		}
		for _, rel := range amRelations {
			if !ownedRelOIDs[rel.OID] {
				return errDependentObjects("extension", c.extensionMembers[rel.OID])
			}
		}
		for _, x := range amIndexes {
			if !ownedRelOIDs[x.idx.RelOID] && !inDropSet('i', x.idx.OID) {
				return errDependentObjects("extension", x.extName)
			}
		}
		for _, x := range amUsers {
			if inOwnedExt(x.dep.ObjOID) {
				continue // part of the same extension's drop set
			}
			return errDependentObjects("extension", x.extName)
		}
		for _, x := range externalCasts {
			return errDependentObjects("extension", x.extName)
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
			c.removeIndex(x.schema, x.idx.Name, x.idx)
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
	for _, ext := range ownedExts {
		if c.extensions[ext.OID] != nil {
			c.removeExtensionMetadata(ext)
		}
	}
	for _, x := range opclassIndexes {
		if c.indexes[x.idx.OID] == x.idx {
			c.removeIndex(x.idx.Schema, x.idx.Name, x.idx)
		}
	}
	if cascade {
		for _, rel := range amRelations {
			if c.relationByOID[rel.OID] != nil {
				c.cascadeDropRelation(rel)
			}
		}
		for _, x := range amIndexes {
			if c.indexes[x.idx.OID] == x.idx {
				c.removeIndex(x.idx.Schema, x.idx.Name, x.idx)
			}
		}
		for _, x := range amUsers {
			c.removeAMUser(x.dep)
		}
		for _, x := range externalCasts {
			delete(c.castIndex, x.key)
		}
	}
	sort.Slice(ownedOperators, func(i, j int) bool { return ownedOperators[i] < ownedOperators[j] })
	for _, oid := range ownedOperators {
		c.removeOperator(oid)
	}
	for _, s := range ownedSchemas {
		for oid, opc := range c.opClasses {
			if opc.Schema != s.OID {
				continue
			}
			if opc.IsDefault {
				delete(c.opClassByKey, opClassKey{amOID: opc.AMOID, typeOID: opc.TypeOID})
			}
			delete(c.opClasses, oid)
		}
		for oid, fam := range c.opFamilies {
			if fam.Schema == s.OID {
				delete(c.opFamilies, oid)
			}
		}
	}
	for _, evt := range ownedEventTriggers {
		if c.eventTriggers[evt.OID] == evt {
			c.removeEventTrigger(evt)
		}
	}
	for _, s := range ownedSchemas {
		if c.schemas[s.OID] == nil {
			continue
		}
		c.removeGrantsForObject('n', s.OID)
		c.removeComments('n', s.OID)
		delete(c.schemas, s.OID)
		delete(c.schemaByName, s.Name)
	}

	kept := c.grants[:0]
	for _, g := range c.grants {
		if owners[g.Grantee] || owners[g.Grantor] {
			continue
		}
		kept = append(kept, g)
	}
	c.grants = kept

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
		delete(c.policies, p.OID)
		c.removeDepsOf('p', p.OID)
		c.removeComments('p', p.OID)
		plist := c.policiesByRel[p.RelOID]
		for i, item := range plist {
			if item.OID == p.OID {
				c.policiesByRel[p.RelOID] = append(plist[:i], plist[i+1:]...)
				break
			}
		}
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
			}
		}
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
	for oid, owner := range c.operatorOwners {
		if owners[owner] {
			c.operatorOwners[oid] = newOwner
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
	if c.roles[stmt.Newname] {
		return errDuplicateObject("role", stmt.Newname)
	}
	delete(c.roles, oldName)
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
			if c.roleInUse(name) {
				return errDependentObjects("role", name)
			}
			continue // a snapshot may know a role the catalog does not
		}
		if name == c.currentOwner() || name == c.currentSessionUser() {
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
