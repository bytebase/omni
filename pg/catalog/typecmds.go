package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// Enum types
// pg: src/backend/commands/typecmds.c
// ---------------------------------------------------------------------------

// EnumValue represents one label of an enum type.
type EnumValue struct {
	OID       uint32
	EnumOID   uint32
	SortOrder float32
	Label     string
}

// EnumType holds the metadata for a user-created enum type.
type EnumType struct {
	TypeOID  uint32
	Values   []*EnumValue
	labelMap map[string]*EnumValue
}

// moveArrayTypeName renames an existing auto-generated array type to avoid collision.
//
// pg: src/backend/catalog/pg_type.c — moveArrayTypeName
func (c *Catalog) moveArrayTypeName(schema *Schema, arrayName string) error {
	existing := c.typeByName[typeKey{ns: schema.OID, name: arrayName}]
	if existing == nil {
		return nil
	}
	// Only rename if it's an auto-generated array type.
	if existing.Category != 'A' || existing.Elem == 0 {
		return errDuplicateObject("type", arrayName)
	}
	// Find a suffix that doesn't conflict.
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s_%d", arrayName, i)
		if c.typeByName[typeKey{ns: schema.OID, name: newName}] == nil {
			delete(c.typeByName, typeKey{ns: schema.OID, name: arrayName})
			existing.TypeName = newName
			c.typeByName[typeKey{ns: schema.OID, name: newName}] = existing
			return nil
		}
	}
}

// DefineEnum creates a new enum type from a parsed CREATE TYPE ... AS ENUM statement.
//
// pg: src/backend/commands/typecmds.c — DefineEnum
func (c *Catalog) DefineEnum(stmt *nodes.CreateEnumStmt) error {
	schemaName, name := qualifiedName(stmt.TypeName)
	values := stringListItems(stmt.Vals)

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Check name conflicts with existing types.
	if c.typeByName[typeKey{ns: schema.OID, name: name}] != nil {
		return errDuplicateObject("type", name)
	}

	// Check for duplicate labels.
	seen := make(map[string]bool, len(values))
	for _, v := range values {
		if seen[v] {
			return errInvalidParameterValue(fmt.Sprintf("enum label %q specified more than once", v))
		}
		seen[v] = true
	}

	// Allocate OIDs.
	typeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	// Create enum values.
	enumValues := make([]*EnumValue, len(values))
	labelMap := make(map[string]*EnumValue, len(values))
	for i, label := range values {
		ev := &EnumValue{
			OID:       c.oidGen.Next(),
			EnumOID:   typeOID,
			SortOrder: float32(i + 1),
			Label:     label,
		}
		enumValues[i] = ev
		labelMap[label] = ev
	}

	// Register the enum type.
	enumType := &BuiltinType{
		OID:       typeOID,
		TypeName:  name,
		Namespace: schema.OID,
		Len:       4,
		ByVal:     true,
		Type:      'e',
		Category:  'E',
		IsDefined: true,
		Delim:     ',',
		Array:     arrayOID,
		Align:     'i',
		Storage:   'p',
		TypeMod:   -1,
	}
	c.typeByOID[typeOID] = enumType
	c.typeByName[typeKey{ns: schema.OID, name: name}] = enumType

	// Move existing auto-generated array type if it collides.
	// pg: src/backend/catalog/pg_type.c — moveArrayTypeName
	arrayTypeName := fmt.Sprintf("_%s", name)
	if err := c.moveArrayTypeName(schema, arrayTypeName); err != nil {
		return err
	}

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  arrayTypeName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      typeOID,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	c.enumTypes[typeOID] = &EnumType{
		TypeOID:  typeOID,
		Values:   enumValues,
		labelMap: labelMap,
	}

	// Enum types use the built-in anyenum operators (OID 3516-3521) which
	// are resolved via IsBinaryCoercible(enumOID, ANYENUMOID) during
	// operator resolution.
	// pg: src/include/catalog/pg_operator.dat — enum_eq/ne/lt/gt/le/ge operators

	return nil
}

// AlterEnumStmt handles ALTER TYPE ... ADD VALUE and ALTER TYPE ... RENAME VALUE.
//
// pg: src/backend/commands/typecmds.c — AlterEnum
func (c *Catalog) AlterEnumStmt(stmt *nodes.AlterEnumStmt) error {
	schemaName, name := qualifiedName(stmt.Typname)

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	bt := c.typeByName[typeKey{ns: schema.OID, name: name}]
	if bt == nil {
		return errUndefinedType(name)
	}
	if bt.Type != 'e' {
		return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an enum type", name)}
	}

	et := c.enumTypes[bt.OID]
	if et == nil {
		return errUndefinedType(name)
	}

	// RENAME VALUE: oldval is set when renaming.
	// pg: src/backend/commands/typecmds.c — RenameEnumLabel
	if stmt.Oldval != "" {
		ev, exists := et.labelMap[stmt.Oldval]
		if !exists {
			return &Error{Code: CodeInvalidParameterValue,
				Message: fmt.Sprintf("%q is not an existing enum label", stmt.Oldval)}
		}
		if _, dup := et.labelMap[stmt.Newval]; dup {
			return errDuplicateObject("enum label", stmt.Newval)
		}
		delete(et.labelMap, stmt.Oldval)
		ev.Label = stmt.Newval
		et.labelMap[stmt.Newval] = ev
		return nil
	}

	// ADD VALUE
	newValue := stmt.Newval
	ifNotExists := stmt.SkipIfNewvalExists

	var before, after string
	if stmt.NewvalIsAfter {
		after = stmt.NewvalNeighbor
	} else if stmt.NewvalNeighbor != "" {
		before = stmt.NewvalNeighbor
	}

	// Check duplicate label.
	if _, exists := et.labelMap[newValue]; exists {
		if ifNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("enum label %q already exists, skipping", newValue))
			return nil
		}
		return errDuplicateObject("enum label", newValue)
	}

	// Determine sort order.
	// pg: src/backend/commands/typecmds.c — AddEnumLabel (neighbor validation)
	var sortOrder float32
	insertIdx := len(et.Values) // default: append

	if before != "" {
		found := false
		for i, ev := range et.Values {
			if ev.Label == before {
				if i == 0 {
					sortOrder = ev.SortOrder - 1
				} else {
					sortOrder = (et.Values[i-1].SortOrder + ev.SortOrder) / 2
				}
				insertIdx = i
				found = true
				break
			}
		}
		if !found {
			return &Error{Code: CodeInvalidParameterValue,
				Message: fmt.Sprintf("%q is not an existing enum label", before)}
		}
	} else if after != "" {
		found := false
		for i, ev := range et.Values {
			if ev.Label == after {
				if i == len(et.Values)-1 {
					sortOrder = ev.SortOrder + 1
				} else {
					sortOrder = (ev.SortOrder + et.Values[i+1].SortOrder) / 2
				}
				insertIdx = i + 1
				found = true
				break
			}
		}
		if !found {
			return &Error{Code: CodeInvalidParameterValue,
				Message: fmt.Sprintf("%q is not an existing enum label", after)}
		}
	} else {
		if len(et.Values) > 0 {
			sortOrder = et.Values[len(et.Values)-1].SortOrder + 1
		} else {
			sortOrder = 1
		}
	}

	ev := &EnumValue{
		OID:       c.oidGen.Next(),
		EnumOID:   bt.OID,
		SortOrder: sortOrder,
		Label:     newValue,
	}

	// Insert at correct position.
	et.Values = append(et.Values, nil)
	copy(et.Values[insertIdx+1:], et.Values[insertIdx:])
	et.Values[insertIdx] = ev
	et.labelMap[newValue] = ev

	return nil
}

// findTypeDependents finds relations that have columns using this type.
func (c *Catalog) findTypeDependents(typeOID uint32) []uint32 {
	var relOIDs []uint32
	for _, rel := range c.relationByOID {
		for _, col := range rel.Columns {
			if col.TypeOID == typeOID {
				relOIDs = append(relOIDs, rel.OID)
				break
			}
		}
	}
	return relOIDs
}

// EnumValues returns the labels of an enum type, or nil if not an enum.
func (c *Catalog) EnumValues(typeOID uint32) []string {
	et := c.enumTypes[typeOID]
	if et == nil {
		return nil
	}
	labels := make([]string, len(et.Values))
	for i, ev := range et.Values {
		labels[i] = ev.Label
	}
	return labels
}

// ---------------------------------------------------------------------------
// Domain types
// pg: src/backend/commands/typecmds.c
// ---------------------------------------------------------------------------

// DomainConstraint represents a CHECK constraint on a domain.
type DomainConstraint struct {
	OID           uint32
	Name          string
	DomainOID     uint32
	CheckExpr     string
	CheckAnalyzed AnalyzedExpr // analyzed form (for Tier 2 deparse)
	ConValidated  bool         // true if constraint has been validated
}

// DomainType holds metadata for a user-created domain type.
type DomainType struct {
	TypeOID     uint32
	BaseTypeOID uint32
	BaseTypMod  int32
	NotNull     bool
	Default     string
	Constraints []*DomainConstraint
}

// DefineDomain creates a new domain type from a parsed CREATE DOMAIN statement.
//
// pg: src/backend/commands/typecmds.c — DefineDomain
func (c *Catalog) DefineDomain(stmt *nodes.CreateDomainStmt) error {
	schemaName, name := qualifiedName(stmt.Domainname)
	baseType := convertTypeNameToInternal(stmt.Typname)

	// Process constraints.
	// pg: src/backend/commands/typecmds.c — DefineDomain (constraint processing)
	var nullDefined bool
	var typNotNull bool
	var sawDefault bool
	var defaultExpr string
	var defaultRawExpr nodes.Node
	type domCheckDef struct {
		name    string
		expr    string
		rawExpr nodes.Node
	}
	var checkDefs []domCheckDef
	var hasCollateClause bool
	if stmt.Constraints != nil {
		for _, item := range stmt.Constraints.Items {
			// pgparser may put COLLATE clause here as *nodes.CollateClause.
			if _, ok := item.(*nodes.CollateClause); ok {
				hasCollateClause = true
				continue
			}
			con, ok := item.(*nodes.Constraint)
			if !ok {
				continue
			}
			switch con.Contype {
			case nodes.CONSTR_NOTNULL:
				// pg: src/backend/commands/typecmds.c — DefineDomain (conflicting NULL/NOT NULL)
				if nullDefined && !typNotNull {
					return &Error{Code: CodeSyntaxError,
						Message: "conflicting NULL/NOT NULL constraints"}
				}
				typNotNull = true
				nullDefined = true
			case nodes.CONSTR_NULL:
				// pg: src/backend/commands/typecmds.c — DefineDomain (conflicting NULL/NOT NULL)
				if nullDefined && typNotNull {
					return &Error{Code: CodeSyntaxError,
						Message: "conflicting NULL/NOT NULL constraints"}
				}
				typNotNull = false
				nullDefined = true
			case nodes.CONSTR_DEFAULT:
				// pg: src/backend/commands/typecmds.c — DefineDomain (multiple default rejection)
				if sawDefault {
					return &Error{Code: CodeSyntaxError,
						Message: "multiple default expressions"}
				}
				sawDefault = true
				if con.CookedExpr != "" {
					defaultExpr = con.CookedExpr
				} else if con.RawExpr != nil {
					defaultRawExpr = con.RawExpr
				}
			case nodes.CONSTR_CHECK:
				// pg: src/backend/commands/typecmds.c — DefineDomain (NO INHERIT rejection)
				if con.IsNoInherit {
					return errInvalidObjectDefinition("check constraints for domains cannot be marked NO INHERIT")
				}
				expr := con.CookedExpr
				checkDefs = append(checkDefs, domCheckDef{name: con.Conname, expr: expr, rawExpr: con.RawExpr})
			case nodes.CONSTR_UNIQUE:
				return &Error{Code: CodeSyntaxError,
					Message: "unique constraints not possible for domains"}
			case nodes.CONSTR_PRIMARY:
				return &Error{Code: CodeSyntaxError,
					Message: "primary key constraints not possible for domains"}
			case nodes.CONSTR_EXCLUSION:
				return &Error{Code: CodeSyntaxError,
					Message: "exclusion constraints not possible for domains"}
			case nodes.CONSTR_FOREIGN:
				return &Error{Code: CodeSyntaxError,
					Message: "foreign key constraints not possible for domains"}
			case nodes.CONSTR_ATTR_DEFERRABLE, nodes.CONSTR_ATTR_NOT_DEFERRABLE,
				nodes.CONSTR_ATTR_DEFERRED, nodes.CONSTR_ATTR_IMMEDIATE:
				// pg: src/backend/commands/typecmds.c — DefineDomain (deferrability rejection)
				return &Error{Code: CodeFeatureNotSupported,
					Message: "specifying constraint deferrability not supported for domains"}
			}
		}
	}
	notNull := typNotNull

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Check name conflicts.
	if c.typeByName[typeKey{ns: schema.OID, name: name}] != nil {
		return errDuplicateObject("type", name)
	}

	// Resolve base type.
	baseOID, baseMod, err := c.ResolveType(baseType)
	if err != nil {
		return err
	}

	// Validate base type kind — reject pseudo-types.
	// pg: src/backend/commands/typecmds.c — DefineDomain (typtype check, line 777-786)
	if bt := c.typeByOID[baseOID]; bt != nil {
		switch bt.Type {
		case 'b', 'c', 'd', 'e', 'r', 'm':
			// OK: base, composite, domain, enum, range, multirange
		default:
			return errDatatypeMismatch(fmt.Sprintf("%q is not a valid base type for a domain", baseType.Name))
		}
	}

	// Validate COLLATE clause — base type must be collatable.
	// pg: src/backend/commands/typecmds.c — DefineDomain (collation check)
	// pgparser may put the COLLATE clause in CollClause or in Constraints as *nodes.CollateClause.
	if stmt.CollClause != nil || hasCollateClause {
		if bt := c.typeByOID[baseOID]; bt != nil && bt.Collation == 0 && bt.Category != 'S' {
			return errDatatypeMismatch(fmt.Sprintf("collations are not supported by type %s", bt.TypeName))
		}
	}

	// Allocate OIDs.
	typeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	// Register the domain type.
	// pg: src/backend/commands/typecmds.c — DefineDomain (inherits typcategory from base type)
	domCategory := byte('D')
	if bt := c.typeByOID[baseOID]; bt != nil {
		domCategory = bt.Category
	}
	domType := &BuiltinType{
		OID:       typeOID,
		TypeName:  name,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'd',
		Category:  domCategory,
		IsDefined: true,
		Delim:     ',',
		Array:     arrayOID,
		BaseType:  baseOID,
		TypeMod:   baseMod,
		Align:     'i',
		Storage:   'x',
		NotNull:   notNull,
	}

	// Copy properties from base type.
	if bt := c.typeByOID[baseOID]; bt != nil {
		domType.Len = bt.Len
		domType.ByVal = bt.ByVal
		domType.Align = bt.Align
		domType.Storage = bt.Storage
		domType.Collation = bt.Collation
	}

	c.typeByOID[typeOID] = domType
	c.typeByName[typeKey{ns: schema.OID, name: name}] = domType

	// Move existing auto-generated array type if it collides.
	// pg: src/backend/catalog/pg_type.c — moveArrayTypeName
	arrayTypeName := fmt.Sprintf("_%s", name)
	if err := c.moveArrayTypeName(schema, arrayTypeName); err != nil {
		return err
	}

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  arrayTypeName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      typeOID,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
		Collation: domType.Collation,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	// Analyze and deparse the DEFAULT expression if raw.
	if defaultRawExpr != nil {
		if analyzed, err := c.AnalyzeDomainExpr(defaultRawExpr, baseOID, baseMod); err == nil && analyzed != nil {
			defaultExpr = c.DeparseExpr(analyzed, nil, false)
		}
	}

	// Create domain metadata.
	dt := &DomainType{
		TypeOID:     typeOID,
		BaseTypeOID: baseOID,
		BaseTypMod:  baseMod,
		NotNull:     notNull,
		Default:     defaultExpr,
	}

	// Add CHECK constraints if provided (supports multiple).
	// pg: src/backend/commands/typecmds.c — DefineDomain (domainAddCheckConstraint)
	for i, chk := range checkDefs {
		conName := chk.name
		if conName == "" {
			if i == 0 {
				conName = name + "_check"
			} else {
				conName = fmt.Sprintf("%s_check%d", name, i+1)
			}
		}
		dc := &DomainConstraint{
			OID:          c.oidGen.Next(),
			Name:         conName,
			DomainOID:    typeOID,
			CheckExpr:    chk.expr,
			ConValidated: true, // pg: CREATE DOMAIN always creates validated constraints
		}
		// Analyze the CHECK expression using Tier 2 pipeline, then deparse.
		if chk.rawExpr != nil {
			if analyzed, err := c.AnalyzeDomainExpr(chk.rawExpr, baseOID, baseMod); err == nil && analyzed != nil {
				dc.CheckAnalyzed = analyzed
				dc.CheckExpr = c.DeparseExpr(analyzed, nil, true)
				// pg: src/backend/catalog/pg_constraint.c:370-380
				// Record dependencies from constraint to referenced functions/operators.
				c.recordDependencyOnExpr('c', dc.OID, analyzed, DepNormal)
			}
		}
		dt.Constraints = append(dt.Constraints, dc)
	}

	c.domainTypes[typeOID] = dt

	// Domain depends on base type.
	c.recordDependency('t', typeOID, 0, 't', baseOID, 0, DepNormal)

	return nil
}

// AlterDomainStmt alters an existing domain from a parsed ALTER DOMAIN statement.
//
// pg: src/backend/commands/typecmds.c — AlterDomainStmt
func (c *Catalog) AlterDomainStmt(stmt *nodes.AlterDomainStmt) error {
	schemaName, name := qualifiedName(stmt.Typname)

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	bt := c.typeByName[typeKey{ns: schema.OID, name: name}]
	if bt == nil {
		return errUndefinedType(name)
	}
	if bt.Type != 'd' {
		return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a domain", name)}
	}

	dt := c.domainTypes[bt.OID]
	if dt == nil {
		return errUndefinedType(name)
	}

	switch stmt.Subtype {
	case 'T': // SET DEFAULT or DROP DEFAULT
		if stmt.Def != nil {
			var setDefault string
			if con, ok := stmt.Def.(*nodes.Constraint); ok && con.CookedExpr != "" {
				setDefault = con.CookedExpr
			} else {
				if analyzed, err := c.AnalyzeDomainExpr(stmt.Def, dt.BaseTypeOID, dt.BaseTypMod); err == nil && analyzed != nil {
					setDefault = c.DeparseExpr(analyzed, nil, false)
				}
			}
			dt.Default = setDefault
		} else {
			dt.Default = ""
		}
	case 'N': // SET NOT NULL
		dt.NotNull = true
		bt.NotNull = true
	case 'O': // DROP NOT NULL
		dt.NotNull = false
		bt.NotNull = false
	case 'C': // ADD CONSTRAINT
		// pg: src/backend/commands/typecmds.c — AlterDomainAddConstraint
		if con, ok := stmt.Def.(*nodes.Constraint); ok {
			switch con.Contype {
			case nodes.CONSTR_CHECK:
				// OK — only CHECK is allowed on domains
			default:
				return &Error{Code: CodeSyntaxError,
					Message: "only CHECK constraints can be added to domains"}
			}
			addCheckExpr := con.CookedExpr
			addCheckName := con.Conname
			conName := addCheckName
			if conName == "" {
				conName = name + "_check"
			}
			// Check for duplicate constraint name.
			for _, dc := range dt.Constraints {
				if dc.Name == conName {
					return errDuplicateObject("constraint", conName)
				}
			}
			dc := &DomainConstraint{
				OID:          c.oidGen.Next(),
				Name:         conName,
				DomainOID:    bt.OID,
				CheckExpr:    addCheckExpr,
				ConValidated: !con.SkipValidation, // pg: default true unless NOT VALID
			}
			if con.RawExpr != nil {
				if analyzed, err := c.AnalyzeDomainExpr(con.RawExpr, dt.BaseTypeOID, dt.BaseTypMod); err == nil && analyzed != nil {
					dc.CheckAnalyzed = analyzed
					dc.CheckExpr = c.DeparseExpr(analyzed, nil, true)
					// pg: src/backend/catalog/pg_constraint.c:370-380
					c.recordDependencyOnExpr('c', dc.OID, analyzed, DepNormal)
				}
			}
			dt.Constraints = append(dt.Constraints, dc)
		}
	case 'X': // DROP CONSTRAINT
		dropConstraint := stmt.Name
		found := false
		for i, dc := range dt.Constraints {
			if dc.Name == dropConstraint {
				dt.Constraints = append(dt.Constraints[:i], dt.Constraints[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("constraint %q does not exist", dropConstraint)}
		}
	case 'V': // VALIDATE CONSTRAINT
		// pg: src/backend/commands/typecmds.c — AlterDomainValidateConstraint
		validateName := stmt.Name
		found := false
		for _, dc := range dt.Constraints {
			if dc.Name == validateName {
				dc.ConValidated = true
				found = true
				break
			}
		}
		if !found {
			return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("constraint %q does not exist", validateName)}
		}
	}

	return nil
}

// DomainInfo returns the domain metadata for the given type OID, or nil.
func (c *Catalog) DomainInfo(typeOID uint32) *DomainType {
	return c.domainTypes[typeOID]
}

// ---------------------------------------------------------------------------
// Composite types
// pg: src/backend/commands/typecmds.c
// ---------------------------------------------------------------------------

// DefineCompositeType creates a new composite type from a parsed CREATE TYPE ... AS statement.
//
// pg: src/backend/commands/typecmds.c — DefineCompositeType
func (c *Catalog) DefineCompositeType(stmt *nodes.CompositeTypeStmt) error {
	// PG: DefineCompositeType converts CompositeTypeStmt to a CreateStmt
	// and calls DefineRelation with relkind='c'.
	rv := stmt.Typevar
	if rv == nil {
		return errInvalidParameterValue("composite type requires a name")
	}

	createStmt := &nodes.CreateStmt{
		Relation:  rv,
		TableElts: stmt.Coldeflist,
	}

	return c.DefineRelation(createStmt, 'c')
}

// ---------------------------------------------------------------------------
// Range types
// pg: src/backend/commands/typecmds.c
// ---------------------------------------------------------------------------

// RangeType holds the metadata for a user-created range type.
type RangeType struct {
	OID           uint32
	Name          string
	Namespace     uint32
	SubTypeOID    uint32
	ArrayOID      uint32
	MultirangeOID uint32 // auto-generated multirange type
}

// DefineRange creates a new range type from a parsed CREATE TYPE ... AS RANGE statement.
//
// pg: src/backend/commands/typecmds.c — DefineRange
func (c *Catalog) DefineRange(stmt *nodes.CreateRangeStmt) error {
	schemaName, name := qualifiedName(stmt.TypeName)

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Check name conflicts with existing types.
	if c.typeByName[typeKey{ns: schema.OID, name: name}] != nil {
		return errDuplicateObject("type", name)
	}

	// Parse params (DefElem list): subtype is required.
	// pg: src/backend/commands/typecmds.c — DefineRange (duplicate param detection)
	var subtypeOID uint32
	var hasCollation bool
	seenParams := make(map[string]bool)
	if stmt.Params != nil {
		for _, item := range stmt.Params.Items {
			def, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			// Check for duplicate parameters.
			if seenParams[def.Defname] {
				return errInvalidParameterValue("conflicting or redundant options")
			}
			seenParams[def.Defname] = true
			switch def.Defname {
			case "subtype":
				// Resolve the subtype.
				if tn, ok := def.Arg.(*nodes.TypeName); ok {
					oid, _, err := c.resolveTypeName(tn)
					if err != nil {
						return err
					}
					subtypeOID = oid
				} else {
					subtypeName := defElemString(def)
					oid, _, err := c.ResolveType(TypeName{Name: resolveAlias(subtypeName), TypeMod: -1})
					if err != nil {
						return err
					}
					subtypeOID = oid
				}
			case "collation":
				hasCollation = true
			case "subtype_opclass", "canonical", "subtype_diff":
				// These are accepted but not needed for catalog-only tracking.
			}
		}
	}

	if subtypeOID == 0 {
		return errInvalidParameterValue(fmt.Sprintf("type %q requires a subtype", name))
	}

	// Validate subtype is not a pseudo-type.
	// pg: src/backend/commands/typecmds.c — DefineRange (line 1469)
	if st := c.typeByOID[subtypeOID]; st != nil && st.Type == 'p' {
		return errDatatypeMismatch(fmt.Sprintf("range subtype cannot be %s", st.TypeName))
	}

	// Validate collation — subtype must be collatable if collation is specified.
	// pg: src/backend/commands/typecmds.c — DefineRange (collation validation)
	if hasCollation {
		if st := c.typeByOID[subtypeOID]; st != nil && st.Collation == 0 && st.Category != 'S' {
			return errDatatypeMismatch(fmt.Sprintf("collations are not supported by type %s", st.TypeName))
		}
	}

	// Compute alignment from subtype.
	// pg: src/backend/commands/typecmds.c — DefineRange (line 1520)
	align := byte('d') // default double alignment for range types
	if st := c.typeByOID[subtypeOID]; st != nil {
		if st.Align == 'd' {
			align = 'd'
		} else {
			align = 'i' // minimum is int alignment for ranges
		}
	}

	// Allocate OIDs.
	typeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	// Register the range type.
	rangeType := &BuiltinType{
		OID:       typeOID,
		TypeName:  name,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'r',
		Category:  'R',
		IsDefined: true,
		Delim:     ',',
		Array:     arrayOID,
		Align:     align,
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[typeOID] = rangeType
	c.typeByName[typeKey{ns: schema.OID, name: name}] = rangeType

	// Move existing auto-generated array type if it collides.
	// pg: src/backend/catalog/pg_type.c — moveArrayTypeName
	arrayTypeName := fmt.Sprintf("_%s", name)
	if err := c.moveArrayTypeName(schema, arrayTypeName); err != nil {
		return err
	}

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  arrayTypeName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      typeOID,
		Align:     align,
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	// Auto-create multirange type and its array type (PG14+).
	// pg: src/backend/commands/typecmds.c — DefineRange (multirange auto-creation)
	multirangeOID := c.oidGen.Next()
	multirangeArrayOID := c.oidGen.Next()
	multirangeName := name + "_multirange"

	// Move conflicting array type name for multirange.
	multirangeArrayName := fmt.Sprintf("_%s", multirangeName)
	if err := c.moveArrayTypeName(schema, multirangeArrayName); err != nil {
		return err
	}

	multirangeType := &BuiltinType{
		OID:       multirangeOID,
		TypeName:  multirangeName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'm', // TYPTYPE_MULTIRANGE
		Category:  'R',
		IsDefined: true,
		Delim:     ',',
		Array:     multirangeArrayOID,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[multirangeOID] = multirangeType
	c.typeByName[typeKey{ns: schema.OID, name: multirangeName}] = multirangeType

	multirangeArrayType := &BuiltinType{
		OID:       multirangeArrayOID,
		TypeName:  multirangeArrayName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      multirangeOID,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[multirangeArrayOID] = multirangeArrayType
	c.typeByName[typeKey{ns: schema.OID, name: multirangeArrayName}] = multirangeArrayType

	// Store range metadata.
	c.rangeTypes[typeOID] = &RangeType{
		OID:           typeOID,
		Name:          name,
		Namespace:     schema.OID,
		SubTypeOID:    subtypeOID,
		ArrayOID:      arrayOID,
		MultirangeOID: multirangeOID,
	}

	// Range depends on subtype.
	c.recordDependency('t', typeOID, 0, 't', subtypeOID, 0, DepNormal)
	// Multirange depends on range.
	c.recordDependency('t', multirangeOID, 0, 't', typeOID, 0, DepNormal)

	return nil
}

// RangeInfo returns the range metadata for the given type OID, or nil.
func (c *Catalog) RangeInfo(typeOID uint32) *RangeType {
	return c.rangeTypes[typeOID]
}

// ---------------------------------------------------------------------------
// Base types (CREATE TYPE name / CREATE TYPE name (...))
// pg: src/backend/commands/typecmds.c — DefineType
// ---------------------------------------------------------------------------

// DefineType creates a new base type.
//
// Two-phase protocol:
//  1. Shell type (no parameters): creates an undefined shell entry.
//  2. Full definition (with parameters): requires existing shell type, parses
//     params, looks up I/O functions, registers the complete type + array type.
//
// pg: src/backend/commands/typecmds.c — DefineType
func (c *Catalog) DefineType(stmt *nodes.DefineStmt) error {
	schemaName, name := qualifiedName(stmt.Defnames)

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Look to see if type already exists.
	// pg: src/backend/commands/typecmds.c — DefineType (line 238-254)
	existing := c.typeByName[typeKey{ns: schema.OID, name: name}]
	if existing != nil && existing.IsDefined {
		// If it's a defined auto-generated array type, rename it out of the way.
		if err := c.moveArrayTypeName(schema, name); err != nil {
			return errDuplicateObject("type", name)
		}
		// Re-check after move.
		existing = c.typeByName[typeKey{ns: schema.OID, name: name}]
	}

	// If parameterless CREATE TYPE, make a shell type.
	// pg: src/backend/commands/typecmds.c — DefineType (line 260-269)
	if stmt.Definition == nil || len(stmt.Definition.Items) == 0 {
		if existing != nil {
			return &Error{
				Code:    CodeDuplicateObject,
				Message: fmt.Sprintf("type \"%s\" already exists", name),
			}
		}
		shellOID := c.oidGen.Next()
		shellType := &BuiltinType{
			OID:       shellOID,
			TypeName:  name,
			Namespace: schema.OID,
			Len:       4,
			ByVal:     true,
			Type:      'p', // pseudo until fully defined
			Category:  'U', // TYPCATEGORY_USER
			IsDefined: false,
			Delim:     ',',
			Align:     'i',
			Storage:   'p',
			TypeMod:   -1,
		}
		c.typeByOID[shellOID] = shellType
		c.typeByName[typeKey{ns: schema.OID, name: name}] = shellType
		return nil
	}

	// Full definition: must already have a shell type.
	// pg: src/backend/commands/typecmds.c — DefineType (line 275-279)
	if existing == nil {
		return &Error{
			Code:    CodeUndefinedObject,
			Message: fmt.Sprintf("type \"%s\" does not exist", name),
		}
	}

	// Parse parameters.
	// pg: src/backend/commands/typecmds.c — DefineType (foreach loop, line 282-339)
	var (
		internalLength int16 = -1    // default: variable-length
		inputName      string        // required
		outputName     string        // required
		receiveName    string
		sendName       string
		typmodinName   string
		typmodoutName  string
		analyzeFnName  string
		category       byte   = 'U'  // TYPCATEGORY_USER
		preferred      bool
		delimiter      byte   = ','  // DEFAULT_TYPDELIM
		elemType       uint32
		defaultValue   string
		byValue        bool
		alignment      byte   = 'i'  // TYPALIGN_INT
		storage        byte   = 'p'  // TYPSTORAGE_PLAIN
		collation      uint32
	)

	seenParams := make(map[string]bool)
	for _, item := range stmt.Definition.Items {
		def, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		paramName := def.Defname

		// Map analyse → analyze
		if paramName == "analyse" {
			paramName = "analyze"
		}

		// Check for duplicate parameters.
		if seenParams[paramName] {
			return errInvalidParameterValue("conflicting or redundant options")
		}
		seenParams[paramName] = true

		switch paramName {
		case "internallength":
			if v, ok := defElemInt(def); ok {
				internalLength = int16(v)
			} else {
				s := strings.ToLower(defElemString(def))
				if s == "variable" {
					internalLength = -1
				}
			}
		case "input":
			inputName = defElemString(def)
		case "output":
			outputName = defElemString(def)
		case "receive":
			receiveName = defElemString(def)
		case "send":
			sendName = defElemString(def)
		case "typmod_in":
			typmodinName = defElemString(def)
		case "typmod_out":
			typmodoutName = defElemString(def)
		case "analyze":
			analyzeFnName = defElemString(def)
		case "subscript":
			// Accepted but not tracked in pgddl.
		case "category":
			s := defElemString(def)
			if len(s) > 0 {
				ch := s[0]
				if ch < 32 || ch > 126 {
					return errInvalidParameterValue(fmt.Sprintf("invalid type category \"%s\": must be simple ASCII", s))
				}
				category = ch
			}
		case "preferred":
			preferred = defElemBool(def)
		case "delimiter":
			s := defElemString(def)
			if len(s) > 0 {
				delimiter = s[0]
			}
		case "element":
			// Resolve element type.
			if tn, ok := def.Arg.(*nodes.TypeName); ok {
				oid, _, err := c.resolveTypeName(tn)
				if err != nil {
					return err
				}
				elemType = oid
			} else {
				elemName := defElemString(def)
				oid, _, err := c.ResolveType(TypeName{Name: resolveAlias(elemName), TypeMod: -1})
				if err != nil {
					return err
				}
				elemType = oid
			}
			// Disallow arrays of pseudotypes.
			// pg: src/backend/commands/typecmds.c — DefineType (line 402-407)
			if elemType != 0 {
				if et := c.typeByOID[elemType]; et != nil && et.Type == 'p' {
					return errDatatypeMismatch(fmt.Sprintf("array element type cannot be %s", et.TypeName))
				}
			}
		case "default":
			defaultValue = defElemString(def)
		case "passedbyvalue":
			byValue = defElemBool(def)
		case "alignment":
			a := strings.ToLower(defElemString(def))
			switch {
			case a == "double" || a == "float8" || a == "pg_catalog.float8":
				alignment = 'd'
			case a == "int4" || a == "pg_catalog.int4":
				alignment = 'i'
			case a == "int2" || a == "pg_catalog.int2":
				alignment = 's'
			case a == "char" || a == "pg_catalog.bpchar":
				alignment = 'c'
			default:
				return errInvalidParameterValue(fmt.Sprintf("alignment \"%s\" not recognized", a))
			}
		case "storage":
			a := strings.ToLower(defElemString(def))
			switch a {
			case "plain":
				storage = 'p'
			case "external":
				storage = 'e'
			case "extended":
				storage = 'x'
			case "main":
				storage = 'm'
			default:
				return errInvalidParameterValue(fmt.Sprintf("storage \"%s\" not recognized", a))
			}
		case "collatable":
			if defElemBool(def) {
				collation = DEFAULT_COLLATION_OID
			}
		case "like":
			// LIKE type — copy properties from an existing type.
			// pg: src/backend/commands/typecmds.c — DefineType (line 346-358)
			likeName := defElemString(def)
			likeOID, _, err := c.ResolveType(TypeName{Name: resolveAlias(likeName), TypeMod: -1})
			if err != nil {
				return err
			}
			if lt := c.typeByOID[likeOID]; lt != nil {
				internalLength = lt.Len
				byValue = lt.ByVal
				alignment = lt.Align
				storage = lt.Storage
			}
		default:
			// PG issues WARNING for unrecognized attributes.
			c.addWarning(CodeSyntaxError, fmt.Sprintf("type attribute \"%s\" not recognized", def.Defname))
		}
	}

	// Validate required I/O functions.
	// pg: src/backend/commands/typecmds.c — DefineType (line 462-474)
	if inputName == "" {
		return errInvalidObjectDefinition("type input function must be specified")
	}
	if outputName == "" {
		return errInvalidObjectDefinition("type output function must be specified")
	}
	if typmodinName == "" && typmodoutName != "" {
		return errInvalidObjectDefinition("type modifier output function is useless without a type modifier input function")
	}

	// Look up I/O function OIDs.
	// pg: src/backend/commands/typecmds.c — findTypeInputFunction etc.
	// For pgddl: verify function exists in procByName and store OID.
	inputOID := c.findProcByName(inputName)
	outputOID := c.findProcByName(outputName)
	var receiveOID, sendOID, typmodinOID, typmodoutOID, analyzeOID uint32
	if receiveName != "" {
		receiveOID = c.findProcByName(receiveName)
	}
	if sendName != "" {
		sendOID = c.findProcByName(sendName)
	}
	if typmodinName != "" {
		typmodinOID = c.findProcByName(typmodinName)
	}
	if typmodoutName != "" {
		typmodoutOID = c.findProcByName(typmodoutName)
	}
	if analyzeFnName != "" {
		analyzeOID = c.findProcByName(analyzeFnName)
	}

	// Allocate array OID.
	arrayOID := c.oidGen.Next()

	// Update shell type to full type.
	// pg: src/backend/commands/typecmds.c — TypeCreate (line 572-604)
	existing.Len = internalLength
	existing.ByVal = byValue
	existing.Type = 'b' // TYPTYPE_BASE
	existing.Category = category
	existing.IsPreferred = preferred
	existing.IsDefined = true
	existing.Delim = delimiter
	existing.Input = inputOID
	existing.Output = outputOID
	existing.Receive = receiveOID
	existing.Send = sendOID
	existing.ModIn = typmodinOID
	existing.ModOut = typmodoutOID
	existing.Analyze = analyzeOID
	existing.Elem = elemType
	existing.Array = arrayOID
	existing.Align = alignment
	existing.Storage = storage
	existing.TypeMod = -1
	existing.Collation = collation

	// Create the array type.
	// pg: src/backend/commands/typecmds.c — DefineType (line 610-646)
	arrayTypeName := fmt.Sprintf("_%s", name)
	if err := c.moveArrayTypeName(schema, arrayTypeName); err != nil {
		return err
	}

	// Array alignment: must be at least TYPALIGN_INT.
	// pg: src/backend/commands/typecmds.c — DefineType (line 613)
	arrayAlign := alignment
	if arrayAlign != 'd' {
		arrayAlign = 'i'
	}

	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  arrayTypeName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     delimiter,
		Elem:      existing.OID,
		Align:     arrayAlign,
		Storage:   'x', // TYPSTORAGE_EXTENDED
		TypeMod:   -1,
		ModIn:     typmodinOID,
		ModOut:    typmodoutOID,
		Collation: collation,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayTypeName}] = arrayType

	_ = defaultValue // Default value is stored on the type but not tracked further in pgddl.
	_ = preferred    // Already set on the type above.

	return nil
}

// findProcByName returns the OID of the first function matching the given name,
// or 0 if not found.
//
// (pgddl helper — PG uses findTypeInputFunction/findTypeOutputFunction etc.)
func (c *Catalog) findProcByName(name string) uint32 {
	if procs := c.procByName[name]; len(procs) > 0 {
		return procs[0].OID
	}
	return 0
}

