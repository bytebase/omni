package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// DefineOperator creates a user-defined operator.
//
// pg: src/backend/commands/operatorcmds.c — DefineOperator
func (c *Catalog) DefineOperator(stmt *nodes.DefineStmt) error {
	schemaName, oprName := qualifiedName(stmt.Defnames)
	opSchema, schemaErr := c.resolveTargetSchema(schemaName)
	if schemaErr != nil {
		return schemaErr
	}
	if oprName == "" {
		return errInvalidParameterValue("no operator name specified")
	}

	// Parse parameters from the DefElem list.
	// pg: src/backend/commands/operatorcmds.c — DefineOperator (line 98-170)
	var (
		functionName   string
		leftTypeName   *nodes.TypeName
		rightTypeName  *nodes.TypeName
		commutatorName string
		negatorName    string
		restrictName   string
		joinName       string
		canHash        bool
		canMerge       bool
	)

	if stmt.Definition != nil {
		for _, item := range stmt.Definition.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "leftarg":
				if tn, ok := de.Arg.(*nodes.TypeName); ok {
					leftTypeName = tn
				}
			case "rightarg":
				if tn, ok := de.Arg.(*nodes.TypeName); ok {
					rightTypeName = tn
				}
			case "function", "procedure":
				functionName = defElemQualifiedName(de)
			case "commutator":
				commutatorName = defElemOperatorName(de)
			case "negator":
				negatorName = defElemOperatorName(de)
			case "restrict":
				restrictName = defElemQualifiedName(de)
			case "join":
				joinName = defElemQualifiedName(de)
			case "hashes":
				canHash = defElemBool(de)
			case "merges":
				canMerge = defElemBool(de)
			case "sort1", "sort2", "ltcmp", "gtcmp":
				// pg: obsolete merge-join options, treat as canMerge=true
				canMerge = true
			default:
				c.addWarning(CodeSyntaxError,
					fmt.Sprintf("operator attribute \"%s\" not recognized", de.Defname))
			}
		}
	}

	// Validation.
	// pg: DefineOperator (line 174-184)
	if functionName == "" {
		return errInvalidObjectDefinition("operator function must be specified")
	}
	if leftTypeName == nil && rightTypeName == nil {
		return errInvalidObjectDefinition("operator argument types must be specified")
	}
	// pg: DefineOperator (line 186-189) — postfix operators removed in PG14
	if leftTypeName != nil && rightTypeName == nil {
		return errInvalidObjectDefinition("operator right argument type must be specified")
	}

	// Resolve operand types.
	var typeId1, typeId2 uint32
	if leftTypeName != nil {
		lt := convertTypeNameToInternal(leftTypeName)
		oid, _, err := c.ResolveType(lt)
		if err != nil {
			return err
		}
		typeId1 = oid
	}
	if rightTypeName != nil {
		rt := convertTypeNameToInternal(rightTypeName)
		oid, _, err := c.ResolveType(rt)
		if err != nil {
			return err
		}
		typeId2 = oid
	}

	// Look up the implementing function.
	// pg: DefineOperator → OperatorCreate → OperatorLookup (function resolution)
	var funcSchema string
	if stmt.Definition != nil {
		for _, item := range stmt.Definition.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			if low := strings.ToLower(de.Defname); low == "function" || low == "procedure" {
				switch v := de.Arg.(type) {
				case *nodes.List:
					funcSchema, _ = qualifiedName(v)
				case *nodes.TypeName:
					funcSchema, _ = typeNameParts(v)
				}
			}
		}
	}
	funcOID, funcErr := c.findOperatorFuncInSchema(funcSchema, functionName, typeId1, typeId2)
	if funcErr != nil {
		return funcErr
	}
	if funcOID == 0 {
		var argOIDs []uint32
		if typeId1 != 0 {
			argOIDs = append(argOIDs, typeId1)
		}
		argOIDs = append(argOIDs, typeId2)
		return errUndefinedFunction(functionName, argOIDs)
	}

	// Duplicate check in target schema.
	// pg: src/backend/catalog/pg_operator.c — OperatorCreate (42723 ERRCODE_DUPLICATE_FUNCTION)
	key := operKey{name: oprName, left: typeId1, right: typeId2}
	for _, existing := range c.operByKey[key] {
		existingNS := uint32(PGCatalogNamespace)
		if existing.OID >= FirstNormalObjectId {
			existingNS = c.operatorSchemas[existing.OID]
		}
		if existingNS == opSchema.OID {
			return errDuplicateFunction(oprName)
		}
	}

	// Determine result type from function return type.
	funcProc := c.procByOID[funcOID]
	rettype := funcProc.RetType

	// Determine operator kind.
	var kind byte
	if typeId1 != 0 {
		kind = 'b' // binary
	} else {
		kind = 'l' // left (prefix)
	}

	// Look up commutator/negator if specified.
	var commutatorOID, negatorOID uint32
	if commutatorName != "" {
		commutatorOID = c.lookupOperatorByName(commutatorName, typeId2, typeId1)
	}
	if negatorName != "" {
		negatorOID = c.lookupOperatorByName(negatorName, typeId1, typeId2)
	}

	// Look up restriction/join estimators (best effort).
	var restrictOID, joinOID uint32
	if restrictName != "" {
		restrictOID = c.findProcByName(restrictName)
	}
	if joinName != "" {
		joinOID = c.findProcByName(joinName)
	}

	// Register the operator.
	op := &BuiltinOperator{
		OID:      c.oidGen.Next(),
		Name:     oprName,
		Kind:     kind,
		CanMerge: canMerge,
		CanHash:  canHash,
		Left:     typeId1,
		Right:    typeId2,
		Result:   rettype,
		Com:      commutatorOID,
		Negate:   negatorOID,
		Code:     funcOID,
		Rest:     restrictOID,
		Join:     joinOID,
	}
	c.operByOID[op.OID] = op
	c.operatorOwners[op.OID] = c.currentOwner()
	c.operatorSchemas[op.OID] = opSchema.OID
	for _, fnOID := range []uint32{funcOID, restrictOID, joinOID} {
		if fnOID != 0 {
			c.recordDependency('O', op.OID, 0, 'f', fnOID, 0, DepNormal)
		}
	}
	for _, typeOID := range []uint32{typeId1, typeId2, rettype} {
		if typeOID != 0 {
			c.recordDependency('O', op.OID, 0, 't', typeOID, 0, DepNormal)
		}
	}
	c.operByKey[key] = append(c.operByKey[key], op)

	// Cross-link commutator if found (only on user-defined operators so shared
	// built-in operator entries are never mutated).
	if commutatorOID >= FirstNormalObjectId {
		if comOp := c.operByOID[commutatorOID]; comOp != nil && comOp.Com == 0 {
			comOp.Com = op.OID
		}
	}
	// Self-commutator: if commutator is the same operator name with reversed types
	// and both types match, point to self.
	if commutatorName == oprName && typeId1 == typeId2 {
		op.Com = op.OID
	}

	// Cross-link negator if found.
	if negatorOID >= FirstNormalObjectId {
		if negOp := c.operByOID[negatorOID]; negOp != nil && negOp.Negate == 0 {
			negOp.Negate = op.OID
		}
	}

	return nil
}

func (c *Catalog) findOperatorFuncInSchema(schemaName, name string, leftType, rightType uint32) (uint32, error) {
	var argOIDs []uint32
	if leftType != 0 {
		argOIDs = []uint32{leftType, rightType}
	} else {
		argOIDs = []uint32{rightType}
	}
	matchBuiltin := func() uint32 {
		for _, p := range c.procByName[name] {
			if p.OID >= FirstNormalObjectId || int(p.NArgs) != len(argOIDs) {
				continue
			}
			ok := true
			for i, a := range argOIDs {
				if p.ArgTypes[i] != a {
					ok = false
					break
				}
			}
			if ok {
				return p.OID
			}
		}
		return 0
	}
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return 0, err
		}
		if bp := c.findExactProc(schema, name, argOIDs); bp != nil {
			return bp.OID, nil
		}
		if schema.OID == PGCatalogNamespace {
			return matchBuiltin(), nil
		}
		return 0, nil
	}
	for _, nsOID := range c.routineSearchPathWithCatalog() {
		if ns := c.schemas[nsOID]; ns != nil {
			if bp := c.findExactProc(ns, name, argOIDs); bp != nil {
				return bp.OID, nil
			}
		}
		if nsOID == PGCatalogNamespace {
			if oid := matchBuiltin(); oid != 0 {
				return oid, nil
			}
		}
	}
	return 0, nil
}

// lookupOperatorByName finds an existing operator by name and arg types.
//
// (pgddl helper)
func (c *Catalog) lookupOperatorByName(name string, left, right uint32) uint32 {
	key := operKey{name: name, left: left, right: right}
	if ops := c.operByKey[key]; len(ops) > 0 {
		return ops[0].OID
	}
	return 0
}

// lookupUserOperator resolves an ObjectWithArgs to a user-defined operator.
func (c *Catalog) lookupUserOperator(owa *nodes.ObjectWithArgs) (*BuiltinOperator, error) {
	if owa == nil {
		return nil, errUndefinedObject("operator", "")
	}
	schemaName, oprName := qualifiedName(owa.Objname)
	var opSchema *Schema
	if schemaName != "" {
		var err error
		opSchema, err = c.resolveTargetSchema(schemaName)
		if err != nil {
			return nil, err
		}
	}
	var left, right uint32
	if owa.Objargs != nil {
		var oids []uint32
		for _, item := range owa.Objargs.Items {
			tn, ok := item.(*nodes.TypeName)
			if !ok {
				oids = append(oids, 0)
				continue
			}
			oid, _, err := c.resolveTypeName(tn)
			if err != nil {
				return nil, err
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
		for i, visible := range c.routineSearchPathWithCatalog() {
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
		return nil, errUndefinedObject("operator", oprName)
	}
	return found, nil
}

func (c *Catalog) alterOperatorNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	owa, ok := stmt.Object.(*nodes.ObjectWithArgs)
	if !ok {
		return nil
	}
	op, err := c.lookupUserOperator(owa)
	if err != nil {
		return err
	}
	oldNS := c.operatorSchemas[op.OID]
	if oldNS == newSchema.OID {
		return nil
	}
	for _, other := range c.operByKey[operKey{name: op.Name, left: op.Left, right: op.Right}] {
		otherNS := uint32(PGCatalogNamespace)
		if other.OID >= FirstNormalObjectId {
			otherNS = c.operatorSchemas[other.OID]
		}
		if other.OID != op.OID && otherNS == newSchema.OID {
			return errDuplicateObject("operator", op.Name)
		}
	}
	c.operatorSchemas[op.OID] = newSchema.OID
	return nil
}

// removeOperator forgets a user-defined operator, clearing the commutator and
// negator links other operators hold on it.
func (c *Catalog) removeOperator(oid uint32) {
	op := c.operByOID[oid]
	if op == nil || oid < FirstNormalObjectId {
		return
	}
	for otherOID, other := range c.operByOID {
		if otherOID < FirstNormalObjectId {
			continue
		}
		if other.Com == oid {
			other.Com = 0
		}
		if other.Negate == oid {
			other.Negate = 0
		}
	}
	delete(c.operByOID, oid)
	delete(c.operatorOwners, oid)
	delete(c.operatorSchemas, oid)
	delete(c.extensionMembers, oid)
	c.removeComments('O', oid)
	c.removeDepsOf('O', oid)
	c.removeDepsOn('O', oid)
	for key, ops := range c.operByKey {
		n := ops[:0]
		for _, item := range ops {
			if item.OID != oid {
				n = append(n, item)
			}
		}
		if len(n) == 0 {
			delete(c.operByKey, key)
			continue
		}
		c.operByKey[key] = n
	}
}

// defElemQualifiedName extracts a qualified name from a DefElem as a simple string.
// Handles both String arg and List of Strings (schema.name).
//
// (pgddl helper)
func defElemQualifiedName(de *nodes.DefElem) string {
	if de.Arg == nil {
		return ""
	}
	switch v := de.Arg.(type) {
	case *nodes.String:
		return v.Str
	case *nodes.List:
		// Take the last element as the function name.
		if len(v.Items) > 0 {
			return stringVal(v.Items[len(v.Items)-1])
		}
	case *nodes.TypeName:
		_, n := typeNameParts(v)
		return n
	}
	return fmt.Sprintf("%v", de.Arg)
}

// defElemOperatorName extracts an operator name from a DefElem.
// The parser may encode COMMUTATOR/NEGATOR as a List of Strings.
//
// (pgddl helper)
func defElemOperatorName(de *nodes.DefElem) string {
	if de.Arg == nil {
		return ""
	}
	switch v := de.Arg.(type) {
	case *nodes.String:
		return v.Str
	case *nodes.List:
		// Last element is the operator name.
		if len(v.Items) > 0 {
			return stringVal(v.Items[len(v.Items)-1])
		}
	}
	return ""
}
