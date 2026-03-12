package ddl

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// UserProc represents a user-defined function or procedure.
type UserProc struct {
	OID          uint32
	Name         string
	Schema       *Schema
	Kind         byte // 'f' or 'p'
	ArgTypes     []uint32
	RetType      uint32
	RetSet       bool
	Volatile     byte
	Parallel     byte
	IsStrict     bool
	SecDef       bool
	LeakProof    bool
	Language     string
	Body         string
	ArgNames     []string // all parameter names (IN+OUT+INOUT+TABLE)
	ArgModes     []byte   // 'i','o','b','v','t' per parameter; nil if all IN
	AllArgTypes  []uint32 // all parameter types (IN+OUT+INOUT); nil if no OUT params
	NArgDefaults int16    // count of input parameters with defaults
}

// CreateFunctionStmt creates a user-defined function or procedure from a pgparser AST.
//
// pg: src/backend/commands/functioncmds.c — CreateFunction
func (c *Catalog) CreateFunctionStmt(stmt *nodes.CreateFunctionStmt) error {
	schemaName, name := qualifiedName(stmt.Funcname)
	orReplace := stmt.IsOrReplace

	// Resolve schema.
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Return type: nil means procedure.
	isProcedure := stmt.ReturnType == nil

	// ---------------------------------------------------------------
	// Process options.
	// pg: src/backend/commands/functioncmds.c — compute_function_attributes
	// ---------------------------------------------------------------
	var language, body string
	var volatile, parallel byte
	var isStrict, secDef, isLeakProof bool
	var isWindowFunc bool
	hasAsClause := false
	rowsSet := false

	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "language":
				language = defElemString(de)
			case "as":
				hasAsClause = true
				// Arg is a String or a List of Strings; use first string.
				switch v := de.Arg.(type) {
				case *nodes.String:
					body = v.Str
				case *nodes.List:
					if len(v.Items) > 0 {
						body = stringVal(v.Items[0])
					}
				}
			case "volatility":
				// pg: compute_common_attribute — procedures cannot have volatility
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				switch strings.ToLower(defElemString(de)) {
				case "immutable":
					volatile = 'i'
				case "stable":
					volatile = 's'
				case "volatile":
					volatile = 'v'
				}
			case "parallel":
				// pg: compute_common_attribute — procedures cannot have parallel
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				switch strings.ToLower(defElemString(de)) {
				case "safe":
					parallel = 's'
				case "restricted":
					parallel = 'r'
				case "unsafe":
					parallel = 'u'
				}
			case "strict":
				// pg: compute_common_attribute — procedures cannot have strict
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				isStrict = defElemBool(de)
			case "called_on_null_input":
				// pg: compute_common_attribute — same as strict (inverse), restricted for procedures
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				isStrict = !defElemBool(de)
			case "security":
				// pg: SECURITY DEFINER → Boolean(true), SECURITY INVOKER → Boolean(false)
				// Also handle String("definer"/"invoker") for manually constructed AST.
				if s, ok := de.Arg.(*nodes.String); ok {
					secDef = strings.ToLower(s.Str) == "definer"
				} else {
					secDef = defElemBool(de)
				}
			case "leakproof":
				// pg: src/backend/commands/functioncmds.c — compute_common_attribute
				isLeakProof = defElemBool(de)
			case "window":
				// pg: compute_function_attributes — procedures cannot have WINDOW
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				isWindowFunc = defElemBool(de)
			case "cost":
				// pg: src/backend/commands/functioncmds.c — compute_function_attributes
				if v, ok := defElemInt(de); ok && v <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "COST must be positive"}
				}
			case "rows":
				// pg: src/backend/commands/functioncmds.c — compute_function_attributes
				if v, ok := defElemInt(de); ok && v <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS must be positive"}
				}
				rowsSet = true
			case "set":
				// Ignored.
			}
		}
	}

	// Suppress unused variable warning — isWindowFunc is validated above but not stored.
	_ = isWindowFunc

	// ---------------------------------------------------------------
	// SQL body default language.
	// pg: src/backend/commands/functioncmds.c — CreateFunction (line 1084-1087)
	// If sql_body is present and no language specified, default to "sql".
	// ---------------------------------------------------------------
	hasSqlBody := stmt.SqlBody != nil
	if hasSqlBody && language == "" {
		language = "sql"
	}

	// Validate language name.
	// pg: src/backend/commands/functioncmds.c — CreateFunction (line 1095-1101)
	validLanguages := map[string]bool{
		"sql": true, "plpgsql": true, "c": true, "internal": true,
		"plperl": true, "plperlu": true, "pltcl": true, "pltclu": true,
		"plpythonu": true, "plpython2u": true, "plpython3u": true,
	}
	if language != "" && !validLanguages[strings.ToLower(language)] {
		return &Error{Code: CodeUndefinedObject,
			Message: fmt.Sprintf("language \"%s\" does not exist", language)}
	}

	// ---------------------------------------------------------------
	// Function body requirement.
	// pg: src/backend/commands/functioncmds.c — interpret_AS_clause
	// Either "as" or sql_body must be present.
	// ---------------------------------------------------------------
	if !hasAsClause && !hasSqlBody {
		return errInvalidFunctionDefinition("no function body specified")
	}

	// ---------------------------------------------------------------
	// Resolve parameters and validate.
	// pg: src/backend/commands/functioncmds.c — interpret_function_parameter_list
	// ---------------------------------------------------------------
	var argOIDs []uint32
	var varCount, outCount int
	var requiredResultType uint32 // 0 = InvalidOid
	var hasTableParams bool
	haveDefaults := false
	langIsSQL := strings.ToLower(language) == "sql"

	// Collectors for proargnames, proargmodes, proallargtypes, pronargdefaults.
	var allArgNames []string
	var allArgModes []byte
	var allArgTypes []uint32
	var nArgDefaults int16

	if stmt.Parameters != nil {
		for i, item := range stmt.Parameters.Items {
			fp, ok := item.(*nodes.FunctionParameter)
			if !ok {
				continue
			}

			fpmode := fp.Mode
			// For our purposes, a defaulted mode spec is identical to IN.
			if fpmode == nodes.FUNC_PARAM_DEFAULT {
				fpmode = nodes.FUNC_PARAM_IN
			}

			// Resolve the parameter type.
			argType := convertTypeNameToInternal(fp.ArgType)
			toid, _, err := c.ResolveType(argType)
			if err != nil {
				return err
			}

			// ---------------------------------------------------------------
			// Shell type detection for parameters.
			// pg: interpret_function_parameter_list (line 235-256)
			// SQL functions cannot accept shell types.
			// ---------------------------------------------------------------
			if bt := c.typeByOID[toid]; bt != nil && !bt.IsDefined {
				if langIsSQL {
					return errInvalidFunctionDefinition(
						fmt.Sprintf("SQL function cannot accept shell type %s", bt.TypeName))
				}
			}

			// ---------------------------------------------------------------
			// SETOF parameter rejection.
			// pg: interpret_function_parameter_list — "functions cannot accept set arguments"
			// ---------------------------------------------------------------
			if fp.ArgType != nil && fp.ArgType.Setof {
				if isProcedure {
					return errInvalidFunctionDefinition("procedures cannot accept set arguments")
				}
				return errInvalidFunctionDefinition("functions cannot accept set arguments")
			}

			// Collect parameter metadata.
			// pg: interpret_function_parameter_list — argnames, argmodes, allargtypes
			allArgNames = append(allArgNames, fp.Name)
			allArgTypes = append(allArgTypes, toid)
			switch fpmode {
			case nodes.FUNC_PARAM_IN:
				allArgModes = append(allArgModes, 'i')
			case nodes.FUNC_PARAM_OUT:
				allArgModes = append(allArgModes, 'o')
			case nodes.FUNC_PARAM_INOUT:
				allArgModes = append(allArgModes, 'b')
			case nodes.FUNC_PARAM_VARIADIC:
				allArgModes = append(allArgModes, 'v')
			case nodes.FUNC_PARAM_TABLE:
				allArgModes = append(allArgModes, 't')
			default:
				allArgModes = append(allArgModes, 'i')
			}

			isInput := false

			// Handle input parameters.
			if fpmode != nodes.FUNC_PARAM_OUT && fpmode != nodes.FUNC_PARAM_TABLE {
				// ---------------------------------------------------------------
				// VARIADIC ordering: other input parameters can't follow a VARIADIC parameter.
				// pg: interpret_function_parameter_list
				// ---------------------------------------------------------------
				if varCount > 0 {
					return errInvalidFunctionDefinition("VARIADIC parameter must be the last input parameter")
				}
				argOIDs = append(argOIDs, toid)
				isInput = true
			}

			// Track TABLE parameters for RETURNS TABLE detection.
			if fpmode == nodes.FUNC_PARAM_TABLE {
				hasTableParams = true
			}

			// Handle output parameters.
			// pg: interpret_function_parameter_list (line 304-324)
			if fpmode != nodes.FUNC_PARAM_IN && fpmode != nodes.FUNC_PARAM_VARIADIC {
				if isProcedure {
					if varCount > 0 {
						return errInvalidFunctionDefinition("VARIADIC parameter must be the last parameter")
					}
					// Procedures with output parameters always return RECORD.
					requiredResultType = RECORDOID
				} else if outCount == 0 {
					// Save first output param's type.
					requiredResultType = toid
				}
				outCount++
			}

			// ---------------------------------------------------------------
			// VARIADIC type validation.
			// pg: interpret_function_parameter_list — "VARIADIC parameter must be an array"
			// ---------------------------------------------------------------
			if fpmode == nodes.FUNC_PARAM_VARIADIC {
				varCount++
				switch toid {
				case ANYARRAYOID, ANYCOMPATIBLEARRAYOID, ANYOID:
					// These pseudo-types are okay for VARIADIC.
				default:
					bt := c.typeByOID[toid]
					if bt == nil || bt.Elem == 0 {
						return errInvalidFunctionDefinition("VARIADIC parameter must be an array")
					}
				}
			}

			// ---------------------------------------------------------------
			// Parameter name duplicate detection.
			// pg: interpret_function_parameter_list — "parameter name used more than once"
			// ---------------------------------------------------------------
			if fp.Name != "" {
				for j := 0; j < i; j++ {
					prevItem := stmt.Parameters.Items[j]
					prevFP, ok := prevItem.(*nodes.FunctionParameter)
					if !ok {
						continue
					}
					if prevFP.Name == "" {
						continue
					}
					if prevFP.Name != fp.Name {
						continue
					}

					prevMode := prevFP.Mode
					if prevMode == nodes.FUNC_PARAM_DEFAULT {
						prevMode = nodes.FUNC_PARAM_IN
					}

					// Pure IN doesn't conflict with pure OUT.
					if (fpmode == nodes.FUNC_PARAM_IN || fpmode == nodes.FUNC_PARAM_VARIADIC) &&
						(prevMode == nodes.FUNC_PARAM_OUT || prevMode == nodes.FUNC_PARAM_TABLE) {
						continue
					}
					if (prevMode == nodes.FUNC_PARAM_IN || prevMode == nodes.FUNC_PARAM_VARIADIC) &&
						(fpmode == nodes.FUNC_PARAM_OUT || fpmode == nodes.FUNC_PARAM_TABLE) {
						continue
					}

					return errInvalidFunctionDefinition(
						fmt.Sprintf("parameter name %q used more than once", fp.Name))
				}
			}

			// ---------------------------------------------------------------
			// Default parameter ordering.
			// pg: interpret_function_parameter_list — "input parameters after one with a
			// default value must also have defaults"
			// ---------------------------------------------------------------
			if fp.Defexpr != nil {
				haveDefaults = true
				nArgDefaults++
			} else {
				if isInput && haveDefaults {
					return errInvalidFunctionDefinition(
						"input parameters after one with a default value must also have defaults")
				}
				// For procedures, OUT parameters after one with a default are also rejected.
				if isProcedure && haveDefaults {
					return errInvalidFunctionDefinition(
						"procedure OUT parameters cannot appear after one with a default value")
				}
			}
		}
	}

	// pg: interpret_function_parameter_list (line 467-468)
	// If multiple OUT params for a function, requiredResultType = RECORDOID.
	if outCount > 1 && !isProcedure {
		requiredResultType = RECORDOID
	}

	// ---------------------------------------------------------------
	// Resolve return type.
	// pg: src/backend/commands/functioncmds.c — CreateFunction (line 1171-1202)
	// ---------------------------------------------------------------
	var retOID uint32
	var returnSet bool
	if isProcedure {
		// pg: line 1174 — procedure returns requiredResultType or VOIDOID
		if requiredResultType != 0 {
			retOID = requiredResultType
		} else {
			retOID = VOIDOID
		}
	} else if stmt.ReturnType != nil {
		// Explicit RETURNS clause.
		retType := convertTypeNameToInternal(stmt.ReturnType)
		returnSet = stmt.ReturnType.Setof
		retOID, _, err = c.ResolveType(retType)
		if err != nil {
			// pg: compute_return_type (line 113-153)
			// Only C and internal language functions can auto-create shell types.
			langLower := strings.ToLower(language)
			if langLower != "c" && langLower != "internal" {
				return err // preserve original "type does not exist" error
			}
			// Auto-create shell type.
			// pg: compute_return_type — TypeShellMake
			retOID = c.typeShellMake(retType.Name, schema.OID)
			c.addWarning(CodeUndefinedObject,
				fmt.Sprintf("type \"%s\" is not yet defined", retType.Name))
			err = nil
		}

		// ---------------------------------------------------------------
		// Shell type detection for return type.
		// pg: compute_return_type (line 88-164)
		// SQL functions cannot return shell types.
		// ---------------------------------------------------------------
		if bt := c.typeByOID[retOID]; bt != nil && !bt.IsDefined {
			if langIsSQL {
				return errInvalidFunctionDefinition(
					fmt.Sprintf("SQL function cannot return shell type %s", bt.TypeName))
			}
		}

		// pg: line 1182-1186 — validate against requiredResultType from OUT params
		if requiredResultType != 0 && retOID != requiredResultType {
			return errInvalidFunctionDefinition(
				fmt.Sprintf("function result type must be %s because of OUT parameters",
					c.formatType(requiredResultType, -1)))
		}
	} else if requiredResultType != 0 {
		// pg: line 1188-1192 — default RETURNS clause from OUT parameters
		retOID = requiredResultType
		returnSet = false
	} else {
		// pg: line 1194-1201 — no return type specified and no OUT params
		return errInvalidFunctionDefinition("function result type must be specified")
	}

	// pgparser quirk: RETURNS TABLE(...) doesn't set ReturnType.Setof = true.
	// PG grammar's TableFuncTypeName() always sets setof = true.
	// pg: src/backend/parser/gram.y — TableFuncTypeName
	if !returnSet && hasTableParams {
		returnSet = true
	}

	// ---------------------------------------------------------------
	// ROWS vs returnsSet validation.
	// pg: src/backend/commands/functioncmds.c (line 1248-1251)
	// ROWS is not applicable when function does not return a set.
	// ---------------------------------------------------------------
	if rowsSet && !returnSet {
		return &Error{Code: CodeInvalidParameterValue,
			Message: "ROWS is not applicable when function does not return a set"}
	}

	kind := byte('f')
	if isProcedure {
		kind = 'p'
	}

	if volatile == 0 {
		volatile = 'v'
	}
	if parallel == 0 {
		parallel = 'u'
	}

	// Check for existing function with same signature.
	existing := c.findExactProc(schema, name, argOIDs)
	if existing != nil {
		if !orReplace {
			return errDuplicateFunction(name)
		}
		// OrReplace: verify return type and setof flag match.
		if existing.RetType != retOID || existing.RetSet != returnSet {
			return errInvalidObjectDefinition(fmt.Sprintf(
				"cannot change return type of existing function %q", name))
		}
		// Replace the body and attributes.
		up := c.userProcs[existing.OID]
		if up != nil {
			up.Volatile = volatile
			up.Parallel = parallel
			up.IsStrict = isStrict
			up.SecDef = secDef
			up.LeakProof = isLeakProof
			up.Language = language
			up.Body = body
		}
		existing.Volatile = volatile
		existing.Parallel = parallel
		existing.IsStrict = isStrict
		existing.SecDef = secDef
		existing.LeakProof = isLeakProof
		return nil
	}

	// Allocate OID.
	procOID := c.oidGen.Next()

	// Register as BuiltinProc so resolveFunc works.
	bp := &BuiltinProc{
		OID:       procOID,
		Name:      name,
		Kind:      kind,
		SecDef:    secDef,
		LeakProof: isLeakProof,
		IsStrict:  isStrict,
		RetSet:    returnSet,
		Volatile:  volatile,
		Parallel:  parallel,
		NArgs:     int16(len(argOIDs)),
		RetType:   retOID,
		ArgTypes:  argOIDs,
	}
	c.procByOID[procOID] = bp
	c.procByName[name] = append(c.procByName[name], bp)

	// PG convention: proargnames is NULL when no parameter has a name.
	// pg: interpret_function_parameter_list — if (!have_names) *parameterNames = NIL
	hasNames := false
	for _, n := range allArgNames {
		if n != "" {
			hasNames = true
			break
		}
	}
	if !hasNames {
		allArgNames = nil
	}

	// PG convention: proargmodes and proallargtypes are NULL when all params are IN.
	// pg: interpret_function_parameter_list (line 467)
	hasNonInParam := outCount > 0 || varCount > 0 || hasTableParams
	var storedArgModes []byte
	var storedAllArgTypes []uint32
	if hasNonInParam {
		storedArgModes = allArgModes
		storedAllArgTypes = allArgTypes
	}

	// Create UserProc metadata.
	up := &UserProc{
		OID:          procOID,
		Name:         name,
		Schema:       schema,
		Kind:         kind,
		ArgTypes:     argOIDs,
		RetType:      retOID,
		RetSet:       returnSet,
		Volatile:     volatile,
		Parallel:     parallel,
		IsStrict:     isStrict,
		SecDef:       secDef,
		LeakProof:    isLeakProof,
		Language:     language,
		Body:         body,
		ArgNames:     allArgNames,
		ArgModes:     storedArgModes,
		AllArgTypes:  storedAllArgTypes,
		NArgDefaults: nArgDefaults,
	}
	c.userProcs[procOID] = up

	// Record dependencies on arg/return types.
	for _, aoid := range argOIDs {
		if aoid >= FirstNormalObjectId {
			c.recordDependency('f', procOID, 0, 't', aoid, 0, DepNormal)
		}
	}
	if retOID >= FirstNormalObjectId {
		c.recordDependency('f', procOID, 0, 't', retOID, 0, DepNormal)
	}

	// Record dependencies from default expressions.
	// pg: src/backend/commands/functioncmds.c — CreateFunction (recordDependencyOnExpr for defaults)
	if stmt.Parameters != nil {
		for _, item := range stmt.Parameters.Items {
			fp, ok := item.(*nodes.FunctionParameter)
			if !ok || fp.Defexpr == nil {
				continue
			}
			analyzed, err := c.AnalyzeExprNoContext(fp.Defexpr)
			if err == nil && analyzed != nil {
				c.recordDependencyOnExpr('f', procOID, analyzed, DepNormal)
			}
		}
	}

	return nil
}

// CreateCast registers a user-defined cast.
//
// pg: src/backend/commands/functioncmds.c — CreateCast
func (c *Catalog) CreateCast(stmt *nodes.CreateCastStmt) error {
	// Resolve source type.
	srcType := convertTypeNameToInternal(stmt.Sourcetype)
	srcOID, _, err := c.ResolveType(srcType)
	if err != nil {
		return err
	}
	srcBT := c.typeByOID[srcOID]

	// Resolve target type.
	tgtType := convertTypeNameToInternal(stmt.Targettype)
	tgtOID, _, err := c.ResolveType(tgtType)
	if err != nil {
		return err
	}
	tgtBT := c.typeByOID[tgtOID]

	// pg: CreateCast (line 1591-1599) — pseudo-type check
	if srcBT != nil && srcBT.Type == 'p' {
		return &Error{Code: CodeWrongObjectType,
			Message: fmt.Sprintf("source data type %s is a pseudo-type", srcBT.TypeName)}
	}
	if tgtBT != nil && tgtBT.Type == 'p' {
		return &Error{Code: CodeWrongObjectType,
			Message: fmt.Sprintf("target data type %s is a pseudo-type", tgtBT.TypeName)}
	}

	// pg: CreateCast (line 1603-1612) — domain warning
	if srcBT != nil && srcBT.Type == 'd' {
		c.addWarning(CodeWarning,
			"cast will be ignored because the source data type is a domain")
	}
	if tgtBT != nil && tgtBT.Type == 'd' {
		c.addWarning(CodeWarning,
			"cast will be ignored because the target data type is a domain")
	}

	// Determine cast method.
	// pg: CreateCast (line 1617-1690)
	var method byte
	var funcOID uint32
	if stmt.Func != nil {
		method = 'f'
		funcOID, err = c.resolveObjectWithArgs(stmt.Func)
		if err != nil {
			return err
		}
	} else if stmt.Inout {
		method = 'i'
	} else {
		method = 'b'
	}

	// pg: CreateCast (line 1697-1701) — same-type check
	if srcOID == tgtOID && method != 'f' {
		return &Error{Code: CodeInvalidObjectDefinition,
			Message: fmt.Sprintf("source data type and target data type are the same")}
	}

	// Map coercion context.
	// pg: CreateCast (line 1711-1717)
	var ctx byte
	switch stmt.Context {
	case nodes.COERCION_IMPLICIT:
		ctx = 'i'
	case nodes.COERCION_ASSIGNMENT:
		ctx = 'a'
	default:
		ctx = 'e'
	}

	// Duplicate check.
	// pg: CreateCast (line 1720-1726)
	key := castKey{source: srcOID, target: tgtOID}
	if c.castIndex[key] != nil {
		return &Error{Code: CodeDuplicateObject,
			Message: fmt.Sprintf("cast from type %s to type %s already exists",
				c.formatType(srcOID, -1), c.formatType(tgtOID, -1))}
	}

	// Register the cast.
	cast := &BuiltinCast{
		Source:  srcOID,
		Target:  tgtOID,
		Func:    funcOID,
		Context: ctx,
		Method:  method,
	}
	c.castIndex[key] = cast

	return nil
}

// resolveObjectWithArgs resolves a function from an ObjectWithArgs node.
//
// (pgddl helper — simplified version of PG's LookupFuncWithArgs)
func (c *Catalog) resolveObjectWithArgs(owa *nodes.ObjectWithArgs) (uint32, error) {
	_, funcName := qualifiedName(owa.Objname)
	if funcName == "" {
		return 0, errUndefinedFunction("", nil)
	}

	// Resolve argument types.
	var argOIDs []uint32
	if owa.Objargs != nil {
		for _, item := range owa.Objargs.Items {
			if tn, ok := item.(*nodes.TypeName); ok {
				argType := convertTypeNameToInternal(tn)
				oid, _, err := c.ResolveType(argType)
				if err != nil {
					return 0, err
				}
				argOIDs = append(argOIDs, oid)
			}
		}
	}

	// Search procByName for matching function.
	for _, p := range c.procByName[funcName] {
		if int(p.NArgs) != len(argOIDs) {
			continue
		}
		match := true
		for i, a := range argOIDs {
			if p.ArgTypes[i] != a {
				match = false
				break
			}
		}
		if match {
			return p.OID, nil
		}
	}

	return 0, errUndefinedFunction(funcName, argOIDs)
}

// typeShellMake creates a shell type entry (IsDefined=false).
//
// pg: src/backend/catalog/pg_type.c — TypeShellMake
func (c *Catalog) typeShellMake(name string, namespace uint32) uint32 {
	oid := c.oidGen.Next()
	bt := &BuiltinType{
		OID:       oid,
		TypeName:  name,
		Namespace: namespace,
		IsDefined: false,
		Type:      'p', // pseudo until defined
		Category:  'U', // user-defined
		Delim:     ',',
		Align:     'i',
		Storage:   'p',
		TypeMod:   -1,
	}
	c.typeByOID[oid] = bt
	c.typeByName[typeKey{ns: namespace, name: name}] = bt
	return oid
}

// removeFunction removes a function from all catalog maps.
func (c *Catalog) removeFunction(procOID uint32, name string) {
	delete(c.userProcs, procOID)
	delete(c.procByOID, procOID)

	// Remove from procByName.
	list := c.procByName[name]
	for i, p := range list {
		if p.OID == procOID {
			c.procByName[name] = append(list[:i], list[i+1:]...)
			break
		}
	}
	if len(c.procByName[name]) == 0 {
		delete(c.procByName, name)
	}

	c.removeComments('f', procOID)
	c.removeDepsOf('f', procOID)
	c.removeDepsOn('f', procOID)
}

// findExactProc finds a proc in a specific schema with exact arg type match.
func (c *Catalog) findExactProc(schema *Schema, name string, argOIDs []uint32) *BuiltinProc {
	for _, p := range c.procByName[name] {
		up := c.userProcs[p.OID]
		if up == nil || up.Schema.OID != schema.OID {
			continue
		}
		if int(p.NArgs) != len(argOIDs) {
			continue
		}
		match := true
		for i, a := range argOIDs {
			if p.ArgTypes[i] != a {
				match = false
				break
			}
		}
		if match {
			return p
		}
	}
	return nil
}

// findUserProcsByName finds all user procs with the given name in a schema.
func (c *Catalog) findUserProcsByName(schema *Schema, name string) []*UserProc {
	var result []*UserProc
	for _, up := range c.userProcs {
		if up.Name == name && up.Schema.OID == schema.OID {
			result = append(result, up)
		}
	}
	return result
}

// AlterFunction alters a function's attributes.
//
// pg: src/backend/commands/functioncmds.c — AlterFunction
func (c *Catalog) AlterFunction(stmt *nodes.AlterFunctionStmt) error {
	// Resolve function from ObjectWithArgs.
	if stmt.Func == nil {
		return errInvalidParameterValue("ALTER FUNCTION requires a function name")
	}

	schemaName, funcName := qualifiedName(stmt.Func.Objname)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Resolve argument types.
	var argOIDs []uint32
	if stmt.Func.Objargs != nil {
		for _, item := range stmt.Func.Objargs.Items {
			if tn, ok := item.(*nodes.TypeName); ok {
				argType := convertTypeNameToInternal(tn)
				oid, _, err := c.ResolveType(argType)
				if err != nil {
					return err
				}
				argOIDs = append(argOIDs, oid)
			}
		}
	}

	bp := c.findExactProc(schema, funcName, argOIDs)
	if bp == nil {
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("function %s does not exist", funcName)}
	}

	up := c.userProcs[bp.OID]

	// Process actions.
	if stmt.Actions != nil {
		for _, item := range stmt.Actions.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "volatility":
				v := strings.ToLower(defElemString(de))
				switch v {
				case "immutable":
					bp.Volatile = 'i'
				case "stable":
					bp.Volatile = 's'
				case "volatile":
					bp.Volatile = 'v'
				}
				if up != nil {
					up.Volatile = bp.Volatile
				}
			case "strict":
				bp.IsStrict = defElemBool(de)
				if up != nil {
					up.IsStrict = bp.IsStrict
				}
			case "called_on_null_input":
				bp.IsStrict = !defElemBool(de)
				if up != nil {
					up.IsStrict = bp.IsStrict
				}
			case "security":
				// pg: SECURITY DEFINER → Boolean(true), SECURITY INVOKER → Boolean(false)
				var secDef bool
				if s, ok := de.Arg.(*nodes.String); ok {
					secDef = strings.ToLower(s.Str) == "definer"
				} else {
					secDef = defElemBool(de)
				}
				bp.SecDef = secDef
				if up != nil {
					up.SecDef = secDef
				}
			case "parallel":
				v := strings.ToLower(defElemString(de))
				switch v {
				case "safe":
					bp.Parallel = 's'
				case "restricted":
					bp.Parallel = 'r'
				case "unsafe":
					bp.Parallel = 'u'
				}
				if up != nil {
					up.Parallel = bp.Parallel
				}
			case "cost":
				// pg: src/backend/commands/functioncmds.c — AlterFunction
				if v, ok := defElemInt(de); ok && v <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "COST must be positive"}
				}
			case "rows":
				// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1429-1439)
				if v, ok := defElemInt(de); ok && v <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS must be positive"}
				}
				if !bp.RetSet {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS is not applicable when function does not return a set"}
				}
			case "leakproof":
				// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1413-1420)
				lp := defElemBool(de)
				bp.LeakProof = lp
				if up != nil {
					up.LeakProof = lp
				}
			case "support", "set":
				// No-op: not tracked in pgddl.
			}
		}
	}

	return nil
}
