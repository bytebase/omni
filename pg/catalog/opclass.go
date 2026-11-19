package catalog

import (
	"fmt"

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
	OID    uint32
	Name   string
	AMOID  uint32
	Schema uint32
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

	// Resolve the handler function by its declared schema and signature: a
	// table access method's handler takes one `internal` argument and an
	// index access method's takes none.
	// pg: src/backend/commands/amcmds.c — CreateAccessMethod (lookup_am_handler)
	var handlerOID uint32
	if stmt.HandlerName != nil {
		schemaName, funcName := qualifiedName(stmt.HandlerName)
		if funcName != "" {
			schema, err := c.resolveTargetSchema(schemaName)
			if err != nil {
				return err
			}
			var want []uint32
			if stmt.Amtype == 't' {
				want = []uint32{INTERNALOID}
			}
			if bp := c.findExactProc(schema, funcName, want); bp != nil {
				handlerOID = bp.OID
			} else if candidates := c.findUserProcsByName(schema, funcName); len(candidates) == 1 {
				handlerOID = candidates[0].OID
			}
		}
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

// CreateOpFamily registers a new operator family.
//
// pg: src/backend/commands/opclasscmds.c — CreateOpFamily
func (c *Catalog) CreateOpFamily(stmt *nodes.CreateOpFamilyStmt) error {
	// Resolve access method.
	am := c.accessMethodByName[stmt.Amname]
	if am == nil {
		return errUndefinedObject("access method", stmt.Amname)
	}
	if am.Type != 'i' {
		return &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("access method %q is not of type INDEX", stmt.Amname)}
	}

	// Resolve schema and name.
	schemaName, famName := qualifiedName(stmt.Opfamilyname)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	fam := &OpFamily{
		OID:    c.oidGen.Next(),
		Name:   famName,
		AMOID:  am.OID,
		Schema: schema.OID,
	}
	c.opFamilies[fam.OID] = fam
	c.recordDependency('Q', fam.OID, 0, 'm', am.OID, 0, DepNormal)

	return nil
}

// DefineOpClass registers a new operator class.
//
// pg: src/backend/commands/opclasscmds.c — DefineOpClass
func (c *Catalog) DefineOpClass(stmt *nodes.CreateOpClassStmt) error {
	// Resolve access method.
	am := c.accessMethodByName[stmt.Amname]
	if am == nil {
		return errUndefinedObject("access method", stmt.Amname)
	}
	if am.Type != 'i' {
		return &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("access method %q is not of type INDEX", stmt.Amname)}
	}

	// Resolve indexed data type.
	tn := convertTypeNameToInternal(stmt.Datatype)
	typeOID, _, err := c.ResolveType(tn)
	if err != nil {
		return err
	}

	// Resolve schema and name.
	schemaName, className := qualifiedName(stmt.Opclassname)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Find or auto-create operator family.
	var familyOID uint32
	if stmt.Opfamilyname != nil {
		// Explicit family specified — find it.
		_, famName := qualifiedName(stmt.Opfamilyname)
		found := false
		for _, fam := range c.opFamilies {
			if fam.Name == famName && fam.AMOID == am.OID {
				familyOID = fam.OID
				found = true
				break
			}
		}
		if !found {
			return errUndefinedObject("operator family",
				fmt.Sprintf("%s for access method %s", famName, am.Name))
		}
	} else {
		// Auto-create family with same name as class.
		fam := &OpFamily{
			OID:    c.oidGen.Next(),
			Name:   className,
			AMOID:  am.OID,
			Schema: schema.OID,
		}
		c.opFamilies[fam.OID] = fam
		c.recordDependency('Q', fam.OID, 0, 'm', am.OID, 0, DepNormal)
		familyOID = fam.OID
	}

	opc := &OpClass{
		OID:       c.oidGen.Next(),
		Name:      className,
		FamilyOID: familyOID,
		AMOID:     am.OID,
		TypeOID:   typeOID,
		IsDefault: stmt.IsDefault,
		Schema:    schema.OID,
	}
	c.opClasses[opc.OID] = opc
	c.recordDependency('q', opc.OID, 0, 'm', am.OID, 0, DepNormal)
	if familyOID != 0 {
		c.recordDependency('q', opc.OID, 0, 'Q', familyOID, 0, DepNormal)
	}
	if typeOID != 0 {
		c.recordDependency('q', opc.OID, 0, 't', typeOID, 0, DepNormal)
	}

	// Register as default for (AM, type) pair.

	if opc.IsDefault {
		c.opClassByKey[opClassKey{amOID: am.OID, typeOID: typeOID}] = opc
	}

	return nil
}

// AlterOpFamily adds or drops members from an operator family.
//
// pg: src/backend/commands/opclasscmds.c — AlterOpFamily
func (c *Catalog) AlterOpFamily(stmt *nodes.AlterOpFamilyStmt) error {
	// Resolve access method.
	am := c.accessMethodByName[stmt.Amname]
	if am == nil {
		return errUndefinedObject("access method", stmt.Amname)
	}
	if am.Type != 'i' {
		return &Error{Code: CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("access method %q is not of type INDEX", stmt.Amname)}
	}

	// Verify family exists.
	_, famName := qualifiedName(stmt.Opfamilyname)
	found := false
	for _, fam := range c.opFamilies {
		if fam.Name == famName && fam.AMOID == am.OID {
			found = true
			break
		}
	}
	if !found {
		return errUndefinedObject("operator family",
			fmt.Sprintf("%s for access method %s", famName, am.Name))
	}

	// ADD/DROP members: simplified no-op — pgddl doesn't track individual members.
	return nil
}

// LookupAccessMethod returns the access method with the given name, or nil.
func (c *Catalog) LookupAccessMethod(name string) *AccessMethod {
	return c.accessMethodByName[name]
}
