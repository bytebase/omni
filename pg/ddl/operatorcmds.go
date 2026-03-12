package ddl

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// DefineOperator creates a user-defined operator.
//
// pg: src/backend/commands/operatorcmds.c — DefineOperator
func (c *Catalog) DefineOperator(stmt *nodes.DefineStmt) error {
	_, oprName := qualifiedName(stmt.Defnames)
	if oprName == "" {
		return errInvalidParameterValue("no operator name specified")
	}

	// Parse parameters from the DefElem list.
	// pg: src/backend/commands/operatorcmds.c — DefineOperator (line 98-170)
	var (
		functionName  string
		leftTypeName  *nodes.TypeName
		rightTypeName *nodes.TypeName
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
	funcOID := c.findOperatorFunc(functionName, typeId1, typeId2)
	if funcOID == 0 {
		var argOIDs []uint32
		if typeId1 != 0 {
			argOIDs = append(argOIDs, typeId1)
		}
		argOIDs = append(argOIDs, typeId2)
		return errUndefinedFunction(functionName, argOIDs)
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
	key := operKey{name: oprName, left: typeId1, right: typeId2}
	c.operByKey[key] = append(c.operByKey[key], op)

	// Cross-link commutator if found.
	if commutatorOID != 0 {
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
	if negatorOID != 0 {
		if negOp := c.operByOID[negatorOID]; negOp != nil && negOp.Negate == 0 {
			negOp.Negate = op.OID
		}
	}

	return nil
}

// findOperatorFunc finds a function suitable for an operator implementation.
// It looks for a function with the given name that accepts the operator's arg types.
//
// (pgddl helper)
func (c *Catalog) findOperatorFunc(name string, leftType, rightType uint32) uint32 {
	procs := c.procByName[name]
	if leftType == 0 {
		// Unary (prefix) operator: 1 arg matching rightType.
		for _, p := range procs {
			if p.NArgs == 1 && p.ArgTypes[0] == rightType {
				return p.OID
			}
		}
	} else {
		// Binary operator: 2 args matching leftType, rightType.
		for _, p := range procs {
			if p.NArgs == 2 && p.ArgTypes[0] == leftType && p.ArgTypes[1] == rightType {
				return p.OID
			}
		}
	}
	return 0
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
