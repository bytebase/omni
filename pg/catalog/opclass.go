package catalog

import (
	"fmt"
	"sort"

	nodes "github.com/bytebase/omni/pg/ast"
)

// AccessMethod represents a registered index/table access method.
//
// pg: src/backend/commands/amcmds.c — CreateAccessMethod
type AccessMethod struct {
	OID     uint32
	Name    string
	Type    byte   // 'i' = index, 't' = table
	Handler uint32 // handler function OID (informational only)
}

// OpFamily represents an operator family (groups related operator classes).
//
// pg: src/include/catalog/pg_opfamily.h
type OpFamily struct {
	OID       uint32
	Name      string
	AMOID     uint32
	Schema    uint32
	OwnerName string
	Operators map[string]bool
}

// OpClass represents an operator class for an access method.
//
// pg: src/include/catalog/pg_opclass.h
type OpClass struct {
	OID       uint32
	Name      string
	FamilyOID uint32
	AMOID     uint32
	TypeOID   uint32
	IsDefault bool
	Schema    uint32
	OwnerName string
	Operators map[string]bool
}

// opClassKey identifies a default operator class by (amOID, typeOID).
type opClassKey struct {
	amOID   uint32
	typeOID uint32
}

// CreateAccessMethod registers a new access method.
//
// pg: src/backend/commands/amcmds.c — CreateAccessMethod
func (c *Catalog) CreateAccessMethod(stmt *nodes.CreateAmStmt) error {
	// Duplicate check.
	if c.accessMethodByName[stmt.Amname] != nil {
		return errDuplicateObject("access method", stmt.Amname)
	}

	// Both index and table access-method handlers take one `internal` argument.
	// pg: src/backend/commands/amcmds.c — lookup_am_handler_func
	handlerOID, err := c.resolveAMHandlerOID(stmt.HandlerName, stmt.Amtype)
	if err != nil {
		return err
	}

	am := &AccessMethod{
		OID:     c.oidGen.Next(),
		Name:    stmt.Amname,
		Type:    stmt.Amtype,
		Handler: handlerOID,
	}
	c.accessMethods[am.OID] = am
	c.accessMethodByName[am.Name] = am
	if handlerOID != 0 {
		c.recordDependency('m', am.OID, 0, 'f', handlerOID, 0, DepNormal)
	}

	return nil
}

func (c *Catalog) resolveAMHandlerOID(handlerName *nodes.List, amType byte) (uint32, error) {
	if handlerName == nil {
		return 0, errUndefinedObject("handler function", "")
	}
	schemaName, funcName := qualifiedName(handlerName)
	want := []uint32{INTERNALOID}
	if funcName == "" {
		return 0, errUndefinedFunction("", want)
	}
	findInNamespace := func(nsOID uint32) *BuiltinProc {
		if ns := c.schemas[nsOID]; ns != nil {
			if bp := c.findExactProc(ns, funcName, want); bp != nil {
				return bp
			}
			if nsOID == PGCatalogNamespace {
				for _, p := range c.procByName[funcName] {
					if p.OID < FirstNormalObjectId && int(p.NArgs) == 1 && p.ArgTypes[0] == INTERNALOID {
						return p
					}
				}
			}
		}
		return nil
	}
	var bp *BuiltinProc
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return 0, err
		}
		bp = findInNamespace(schema.OID)
	} else {
		for _, nsOID := range c.routineSearchPathWithCatalog() {
			if found := findInNamespace(nsOID); found != nil {
				bp = found
				break
			}
		}
	}
	if bp == nil {
		return 0, errUndefinedFunction(funcName, want)
	}
	switch amType {
	case 't':
		if bp.RetType != TABLEAMHANDLEROID {
			return 0, &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("function %s must return type table_am_handler", funcName)}
		}
	case 'i':
		if bp.RetType != INDEXAMHANDLEROID {
			return 0, &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("function %s must return type index_am_handler", funcName)}
		}
	default:
		return 0, &Error{Code: CodeInternalError,
			Message: fmt.Sprintf("unrecognized access method type: %c", amType)}
	}
	return bp.OID, nil
}

func (c *Catalog) lookupIndexAccessMethod(name string) (*AccessMethod, error) {
	am := c.accessMethodByName[name]
	if am == nil {
		return nil, errUndefinedObject("access method", name)
	}
	if am.Type != 'i' {
		return nil, &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("access method %q is not of type INDEX", name)}
	}
	return am, nil
}

func (c *Catalog) lookupOpFamily(am *AccessMethod, schemaName, famName string) (*OpFamily, error) {
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return nil, err
		}
		for _, fam := range c.opFamilies {
			if fam.Schema == schema.OID && fam.AMOID == am.OID && fam.Name == famName {
				return fam, nil
			}
		}
		return nil, errUndefinedObject("operator family",
			fmt.Sprintf("%s for access method %s", famName, am.Name))
	}
	var best *OpFamily
	for _, nsOID := range c.searchPathWithCatalog() {
		for _, fam := range c.opFamilies {
			if fam.Schema == nsOID && fam.AMOID == am.OID && fam.Name == famName {
				if best == nil || fam.OID < best.OID {
					best = fam
				}
			}
		}
		if best != nil {
			return best, nil
		}
	}
	return nil, errUndefinedObject("operator family",
		fmt.Sprintf("%s for access method %s", famName, am.Name))
}

func (c *Catalog) lookupOpClass(am *AccessMethod, schemaName, className string) (*OpClass, error) {
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return nil, err
		}
		for _, opc := range c.opClasses {
			if opc.Schema == schema.OID && opc.AMOID == am.OID && opc.Name == className {
				return opc, nil
			}
		}
		return nil, errUndefinedObject("operator class",
			fmt.Sprintf("%s for access method %s", className, am.Name))
	}
	var best *OpClass
	for _, nsOID := range c.searchPathWithCatalog() {
		for _, opc := range c.opClasses {
			if opc.Schema == nsOID && opc.AMOID == am.OID && opc.Name == className {
				if best == nil || opc.OID < best.OID {
					best = opc
				}
			}
		}
		if best != nil {
			return best, nil
		}
	}
	return nil, errUndefinedObject("operator class",
		fmt.Sprintf("%s for access method %s", className, am.Name))
}

// CreateOpFamily registers a new operator family.
//
// pg: src/backend/commands/opclasscmds.c — CreateOpFamily
func (c *Catalog) CreateOpFamily(stmt *nodes.CreateOpFamilyStmt) error {
	am, err := c.lookupIndexAccessMethod(stmt.Amname)
	if err != nil {
		return err
	}

	schemaName, famName := qualifiedName(stmt.Opfamilyname)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}
	for _, existing := range c.opFamilies {
		if existing.Schema == schema.OID && existing.AMOID == am.OID && existing.Name == famName {
			return errDuplicateObject("operator family", famName)
		}
	}

	fam := &OpFamily{
		OID:       c.oidGen.Next(),
		Name:      famName,
		AMOID:     am.OID,
		Schema:    schema.OID,
		OwnerName: c.currentOwner(),
	}
	c.opFamilies[fam.OID] = fam
	c.recordDependency('Q', fam.OID, 0, 'm', am.OID, 0, DepNormal)

	return nil
}

// DefineOpClass registers a new operator class.
//
// pg: src/backend/commands/opclasscmds.c — DefineOpClass
func (c *Catalog) DefineOpClass(stmt *nodes.CreateOpClassStmt) error {
	am, err := c.lookupIndexAccessMethod(stmt.Amname)
	if err != nil {
		return err
	}

	tn := convertTypeNameToInternal(stmt.Datatype)
	typeOID, _, err := c.ResolveType(tn)
	if err != nil {
		return err
	}

	schemaName, className := qualifiedName(stmt.Opclassname)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}
	// pg: OpclassAmNameNspIndexId — (opcmethod, opcname, opcnamespace) is unique.
	for _, existing := range c.opClasses {
		if existing.Schema == schema.OID && existing.AMOID == am.OID && existing.Name == className {
			return errDuplicateObject("operator class", className)
		}
	}
	if stmt.IsDefault && c.opClassByKey[opClassKey{amOID: am.OID, typeOID: typeOID}] != nil {
		return errDuplicateObject("operator class", className)
	}

	var familyOID uint32
	if stmt.Opfamilyname != nil {
		famSchemaName, famName := qualifiedName(stmt.Opfamilyname)
		fam, err := c.lookupOpFamily(am, famSchemaName, famName)
		if err != nil {
			return err
		}
		familyOID = fam.OID
	} else {
		// Reuse an existing family of the same name in the target schema, or create one.
		// pg: src/backend/commands/opclasscmds.c — DefineOpClass
		for _, fam := range c.opFamilies {
			if fam.Schema == schema.OID && fam.AMOID == am.OID && fam.Name == className {
				familyOID = fam.OID
				break
			}
		}
		if familyOID == 0 {
			fam := &OpFamily{
				OID:       c.oidGen.Next(),
				Name:      className,
				AMOID:     am.OID,
				Schema:    schema.OID,
				OwnerName: c.currentOwner(),
			}
			c.opFamilies[fam.OID] = fam
			c.recordDependency('Q', fam.OID, 0, 'm', am.OID, 0, DepNormal)
			familyOID = fam.OID
		}
	}

	var declaredOps map[string]bool
	if stmt.Items != nil {
		for _, raw := range stmt.Items.Items {
			ci, ok := raw.(*nodes.CreateOpClassItem)
			if !ok || ci.Name == nil {
				continue
			}
			if ci.Itemtype == 1 { // OPCLASS_ITEM_OPERATOR
				_, opName := qualifiedName(ci.Name.Objname)
				if opName != "" {
					if declaredOps == nil {
						declaredOps = make(map[string]bool)
					}
					declaredOps[opName] = true
				}
			}
		}
	}

	opc := &OpClass{
		OID:       c.oidGen.Next(),
		Name:      className,
		FamilyOID: familyOID,
		AMOID:     am.OID,
		TypeOID:   typeOID,
		IsDefault: stmt.IsDefault,
		Schema:    schema.OID,
		OwnerName: c.currentOwner(),
		Operators: declaredOps,
	}
	if familyOID != 0 && len(declaredOps) > 0 {
		if fam := c.opFamilies[familyOID]; fam != nil {
			if fam.Operators == nil {
				fam.Operators = make(map[string]bool, len(declaredOps))
			}
			for opName := range declaredOps {
				fam.Operators[opName] = true
			}
		}
	}
	c.opClasses[opc.OID] = opc
	c.recordDependency('q', opc.OID, 0, 'm', am.OID, 0, DepNormal)
	if familyOID != 0 {
		c.recordDependency('q', opc.OID, 0, 'Q', familyOID, 0, DepNormal)
	}
	if typeOID != 0 {
		c.recordDependency('q', opc.OID, 0, 't', typeOID, 0, DepNormal)
	}

	if opc.IsDefault {
		c.opClassByKey[opClassKey{amOID: am.OID, typeOID: typeOID}] = opc
	}

	return nil
}

// AlterOpFamily adds or drops members from an operator family.
//
// pg: src/backend/commands/opclasscmds.c — AlterOpFamily
func (c *Catalog) AlterOpFamily(stmt *nodes.AlterOpFamilyStmt) error {
	am, err := c.lookupIndexAccessMethod(stmt.Amname)
	if err != nil {
		return err
	}

	famSchemaName, famName := qualifiedName(stmt.Opfamilyname)
	fam, err := c.lookupOpFamily(am, famSchemaName, famName)
	if err != nil {
		return err
	}

	if stmt.Items != nil {
		for _, raw := range stmt.Items.Items {
			ci, ok := raw.(*nodes.CreateOpClassItem)
			if !ok || ci == nil || ci.Name == nil || ci.Itemtype != 1 {
				continue
			}
			_, opName := qualifiedName(ci.Name.Objname)
			if opName == "" {
				continue
			}
			if stmt.IsDrop {
				if fam.Operators != nil {
					delete(fam.Operators, opName)
				}
			} else {
				if fam.Operators == nil {
					fam.Operators = make(map[string]bool)
				}
				fam.Operators[opName] = true
			}
		}
	}
	return nil
}

// LookupAccessMethod returns the access method with the given name, or nil.
func (c *Catalog) LookupAccessMethod(name string) *AccessMethod {
	return c.accessMethodByName[name]
}

// extractOpClassObject reads the [amName, name] or [amName, schema, name] list
// that ALTER/RENAME/DROP OPERATOR CLASS and OPERATOR FAMILY carry.
func extractOpClassObject(obj nodes.Node) (amName, schemaName, name string) {
	list, ok := obj.(*nodes.List)
	if !ok || len(list.Items) < 2 {
		return "", "", ""
	}
	amName = stringVal(list.Items[0])
	switch len(list.Items) {
	case 2:
		return amName, "", stringVal(list.Items[1])
	case 3:
		return amName, stringVal(list.Items[1]), stringVal(list.Items[2])
	default:
		return amName, stringVal(list.Items[len(list.Items)-2]), stringVal(list.Items[len(list.Items)-1])
	}
}

func (c *Catalog) alterOpClassOwner(stmt *nodes.AlterOwnerStmt) error {
	amName, schemaName, className := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	opc, err := c.lookupOpClass(am, schemaName, className)
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
	opc.OwnerName = owner
	return nil
}

func (c *Catalog) alterOpFamilyOwner(stmt *nodes.AlterOwnerStmt) error {
	amName, schemaName, famName := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	fam, err := c.lookupOpFamily(am, schemaName, famName)
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
	fam.OwnerName = owner
	return nil
}

func (c *Catalog) renameOpClass(stmt *nodes.RenameStmt) error {
	amName, schemaName, className := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	opc, err := c.lookupOpClass(am, schemaName, className)
	if err != nil {
		return err
	}
	for _, other := range c.opClasses {
		if other.Schema == opc.Schema && other.AMOID == opc.AMOID && other.Name == stmt.Newname {
			return errDuplicateObject("operator class", stmt.Newname)
		}
	}
	opc.Name = stmt.Newname
	return nil
}

func (c *Catalog) renameOpFamily(stmt *nodes.RenameStmt) error {
	amName, schemaName, famName := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	fam, err := c.lookupOpFamily(am, schemaName, famName)
	if err != nil {
		return err
	}
	for _, other := range c.opFamilies {
		if other.Schema == fam.Schema && other.AMOID == fam.AMOID && other.Name == stmt.Newname {
			return errDuplicateObject("operator family", stmt.Newname)
		}
	}
	fam.Name = stmt.Newname
	return nil
}

func (c *Catalog) alterOpClassNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	amName, schemaName, className := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	opc, err := c.lookupOpClass(am, schemaName, className)
	if err != nil {
		return err
	}
	if opc.Schema == newSchema.OID {
		return nil
	}
	for _, other := range c.opClasses {
		if other.OID != opc.OID && other.Schema == newSchema.OID && other.AMOID == opc.AMOID && other.Name == opc.Name {
			return errDuplicateObject("operator class", opc.Name)
		}
	}
	opc.Schema = newSchema.OID
	return nil
}

func (c *Catalog) alterOpFamilyNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	amName, schemaName, famName := extractOpClassObject(stmt.Object)
	am, err := c.lookupIndexAccessMethod(amName)
	if err != nil {
		return err
	}
	fam, err := c.lookupOpFamily(am, schemaName, famName)
	if err != nil {
		return err
	}
	if fam.Schema == newSchema.OID {
		return nil
	}
	for _, other := range c.opFamilies {
		if other.OID != fam.OID && other.Schema == newSchema.OID && other.AMOID == fam.AMOID && other.Name == fam.Name {
			return errDuplicateObject("operator family", fam.Name)
		}
	}
	fam.Schema = newSchema.OID
	return nil
}

// removeOpClass unregisters an operator class and clears its dependency and
// comment rows.
func (c *Catalog) removeOpClass(oid uint32) {
	if opc := c.opClasses[oid]; opc != nil {
		if opc.IsDefault {
			key := opClassKey{amOID: opc.AMOID, typeOID: opc.TypeOID}
			if c.opClassByKey[key] == opc {
				delete(c.opClassByKey, key)
			}
		}
		delete(c.opClasses, oid)
	}
	delete(c.extensionMembers, oid)
	c.removeComments('q', oid)
	c.removeDepsOf('q', oid)
	c.removeDepsOn('q', oid)
	c.removeDepsOn('o', oid)
}

// cascadeDropOpClass removes an operator class after dropping all objects
// (such as indexes and their backing constraints) that depend on it.
func (c *Catalog) cascadeDropOpClass(oid uint32) {
	for _, dep := range c.findNormalDependents('o', oid) {
		c.dropDependentObject(dep.ObjType, dep.ObjOID)
	}
	for _, dep := range c.findNormalDependents('q', oid) {
		c.dropDependentObject(dep.ObjType, dep.ObjOID)
	}
	c.removeOpClass(oid)
}

// removeOpFamily unregisters an operator family and clears its dependency and
// comment rows.
func (c *Catalog) removeOpFamily(oid uint32) {
	delete(c.opFamilies, oid)
	delete(c.extensionMembers, oid)
	c.removeComments('Q', oid)
	c.removeDepsOf('Q', oid)
	c.removeDepsOn('Q', oid)
}

// cascadeDropOpFamily removes an operator family after cascading through all
// operator classes (and their indexes) that belong to it.
func (c *Catalog) cascadeDropOpFamily(oid uint32) {
	deps := c.findNormalDependents('Q', oid)
	sort.Slice(deps, func(i, j int) bool { return deps[i].ObjOID < deps[j].ObjOID })
	for _, dep := range deps {
		if dep.ObjType == 'q' {
			c.cascadeDropOpClass(dep.ObjOID)
			continue
		}
		c.dropDependentObject(dep.ObjType, dep.ObjOID)
	}
	c.removeOpFamily(oid)
}

// removeAccessMethod unregisters a user-defined access method and clears its
// dependency and comment rows.
func (c *Catalog) removeAccessMethod(oid uint32) {
	if oid < FirstNormalObjectId {
		return
	}
	if am := c.accessMethods[oid]; am != nil {
		delete(c.accessMethods, oid)
		if c.accessMethodByName[am.Name] == am {
			delete(c.accessMethodByName, am.Name)
		}
	}
	delete(c.extensionMembers, oid)
	c.removeComments('m', oid)
	c.removeDepsOf('m', oid)
	c.removeDepsOn('m', oid)
}

// accessMethodHasOpClass reports whether any operator class belongs to the
// access method.
func (c *Catalog) accessMethodHasOpClass(amOID uint32) bool {
	for _, opc := range c.opClasses {
		if opc.AMOID == amOID {
			return true
		}
	}
	return false
}

// accessMethodNeedsCascade reports whether DROP ACCESS METHOD needs CASCADE.
// Operator families and classes depend on their access method automatically
// (opclasscmds.c records DEPENDENCY_AUTO), so they alone do not block the
// drop; an object that uses one of them, such as an index, does.
func (c *Catalog) accessMethodNeedsCascade(amOID uint32) bool {
	opClassUsed := func(oid uint32) bool {
		return len(c.findNormalDependents('q', oid)) > 0 || len(c.findNormalDependents('o', oid)) > 0
	}
	for _, dep := range c.findNormalDependents('m', amOID) {
		switch dep.ObjType {
		case 'q':
			if opClassUsed(dep.ObjOID) {
				return true
			}
		case 'Q':
			for _, member := range c.findNormalDependents('Q', dep.ObjOID) {
				if member.ObjType != 'q' || opClassUsed(member.ObjOID) {
					return true
				}
			}
		default:
			return true
		}
	}
	return false
}

// cascadeDropAccessMethod drops all relations, indexes, operator classes, and
// operator families that use an access method, then unregisters it.
func (c *Catalog) cascadeDropAccessMethod(oid uint32) {
	deps := c.findNormalDependents('m', oid)
	sort.Slice(deps, func(i, j int) bool { return deps[i].ObjOID < deps[j].ObjOID })
	for _, user := range deps {
		switch user.ObjType {
		case 'r':
			if rel := c.relationByOID[user.ObjOID]; rel != nil {
				c.cascadeDropRelation(rel)
			}
		case 'i':
			c.dropDependentObject('i', user.ObjOID)
		case 'q':
			c.cascadeDropOpClass(user.ObjOID)
		case 'Q':
			c.cascadeDropOpFamily(user.ObjOID)
		default:
			c.dropDependentObject(user.ObjType, user.ObjOID)
		}
	}
	c.removeAccessMethod(oid)
}
