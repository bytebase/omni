package catalog

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// UserProc represents a user-defined function or procedure.
type UserProc struct {
	OID          uint32
	Name         string
	Schema       *Schema
	OwnerName    string
	Kind         byte // 'f' or 'p'
	ArgTypes     []uint32
	RetType      uint32
	RetSet       bool
	Volatile     byte
	Parallel     byte
	IsStrict     bool
	SecDef       bool
	LeakProof    bool
	IsWindow     bool
	Cost         string
	Rows         string
	SetConfig    []string
	Language     string
	Body         string
	ArgNames     []string // all parameter names (IN+OUT+INOUT+TABLE)
	ArgModes     []byte   // 'i','o','b','v','t' per parameter; nil if all IN
	AllArgTypes  []uint32 // all parameter types (IN+OUT+INOUT); nil if no OUT params
	NArgDefaults int16    // count of input parameters with defaults
	AggKind      byte     // for an aggregate, pg_aggregate.aggkind: 'n' normal, 'o' ordered-set, 'h' hypothetical-set
	// AggNumDirectArgs is pg_aggregate.aggnumdirectargs: the number of
	// direct arguments of an ordered-set or hypothetical-set aggregate.
	AggNumDirectArgs int16
	// ArgDefaultTypes is the type of each parameter default, for the last
	// NArgDefaults input parameters; nil when one could not be analyzed.
	ArgDefaultTypes []uint32
	// ArgDefaults stores deparsed default expressions parallel to ArgTypes.
	ArgDefaults []string
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

	isProcedure := createFunctionStmtIsProcedure(stmt)

	// ---------------------------------------------------------------
	// Process options.
	// pg: src/backend/commands/functioncmds.c — compute_function_attributes
	// ---------------------------------------------------------------
	var language, body string
	var volatile, parallel byte
	var isStrict, secDef, isLeakProof bool
	var isWindowFunc bool
	var costStr, rowsStr string
	var setConfig []string
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
				raw, val, ok := defElemNumeric(de)
				if ok && val <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "COST must be positive"}
				}
				if ok {
					costStr = raw
				}
			case "rows":
				// pg: src/backend/commands/functioncmds.c — compute_function_attributes
				raw, val, ok := defElemNumeric(de)
				if ok && val <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS must be positive"}
				}
				if ok {
					rowsStr = raw
				}
				rowsSet = true
			case "set":
				if vs, ok := de.Arg.(*nodes.VariableSetStmt); ok {
					setConfig = applyFuncSetConfig(setConfig, vs)
				}
			}
		}
	}

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
	var variadicOID uint32
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
					variadicOID = ANYOID
				default:
					bt := c.typeByOID[toid]
					if bt == nil || bt.Elem == 0 {
						return errInvalidFunctionDefinition("VARIADIC parameter must be an array")
					}
					variadicOID = bt.Elem
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
		if hasTableParams && outCount == 1 && requiredResultType != 0 && retOID == RECORDOID {
			retOID = requiredResultType
		}
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

	if retOID == EVENTTRIGGEROID {
		if strings.EqualFold(language, "sql") {
			return errInvalidFunctionDefinition("SQL functions cannot return type event_trigger")
		}
		if len(allArgTypes) > 0 {
			return errInvalidFunctionDefinition("event trigger functions cannot have declared arguments")
		}
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

	argDefaults := c.paramDefaults(stmt.Parameters, len(argOIDs), int(nArgDefaults))
	defaultTypes := c.paramDefaultTypes(stmt.Parameters, argOIDs, int(nArgDefaults))

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
			// A replacement may add parameter defaults, but not remove one
			// or change one's type.
			// pg: src/backend/catalog/pg_proc.c — ProcedureCreate (pronargdefaults)
			if old := int(up.NArgDefaults); old > 0 {
				if int(nArgDefaults) < old {
					return errInvalidFunctionDefinition("cannot remove parameter defaults from existing function")
				}
				if up.ArgDefaultTypes != nil && defaultTypes != nil {
					added := len(defaultTypes) - old
					for i, typ := range up.ArgDefaultTypes {
						if defaultTypes[added+i] != typ {
							return errInvalidFunctionDefinition("cannot change data type of existing parameter default value")
						}
					}
				}
			}
			up.NArgDefaults = nArgDefaults
			up.ArgDefaultTypes = defaultTypes
			up.ArgDefaults = argDefaults
			up.ArgDefaults = argDefaults
			existing.NArgDefaults = nArgDefaults
			up.Volatile = volatile
			up.Parallel = parallel
			up.IsStrict = isStrict
			up.SecDef = secDef
			up.LeakProof = isLeakProof
			up.IsWindow = isWindowFunc
			up.Cost = costStr
			up.Rows = rowsStr
			up.SetConfig = slices.Clone(setConfig)
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
		OID:          procOID,
		Name:         name,
		Kind:         kind,
		SecDef:       secDef,
		LeakProof:    isLeakProof,
		IsStrict:     isStrict,
		RetSet:       returnSet,
		Volatile:     volatile,
		Parallel:     parallel,
		Variadic:     variadicOID,
		NArgs:        int16(len(argOIDs)),
		NArgDefaults: nArgDefaults,
		RetType:      retOID,
		ArgTypes:     argOIDs,
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
	bp.ArgNames = allArgNames
	bp.ArgModes = storedArgModes
	bp.AllArgTypes = storedAllArgTypes

	// Create UserProc metadata.
	up := &UserProc{
		OID:          procOID,
		Name:         name,
		Schema:       schema,
		OwnerName:    c.currentOwner(),
		Kind:         kind,
		ArgTypes:     argOIDs,
		RetType:      retOID,
		RetSet:       returnSet,
		Volatile:     volatile,
		Parallel:     parallel,
		IsStrict:     isStrict,
		SecDef:       secDef,
		LeakProof:    isLeakProof,
		IsWindow:     isWindowFunc,
		Cost:         costStr,
		Rows:         rowsStr,
		SetConfig:    slices.Clone(setConfig),
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
	up.ArgDefaultTypes = defaultTypes
	up.ArgDefaults = argDefaults

	// Record dependencies from SQL function body.
	// pg: src/backend/catalog/pg_proc.c — ProcedureCreate
	//   if (languageObjectId == SQLlanguageId && prosqlbody)
	//       recordDependencyOnExpr(&myself, prosqlbody, NIL, DEPENDENCY_NORMAL);
	//
	// For LANGUAGE sql functions, parse the body and extract table/function
	// references so the migration toposort orders tables before functions
	// that reference them.
	if strings.EqualFold(language, "sql") && body != "" {
		c.recordFuncBodyDeps(procOID, body)
	}

	return nil
}

// paramDefaultTypes returns the type of each parameter default, in order,
// as PostgreSQL stores it: the parameter's type, or for a polymorphic
// parameter the default expression's own type. Defaults belong to the last
// nDefaults input parameters. It returns nil when a default cannot be
// analyzed.
//
// pg: src/backend/commands/functioncmds.c — interpret_function_parameter_list (coerce_to_specific_type)
func (c *Catalog) paramDefaultTypes(params *nodes.List, argOIDs []uint32, nDefaults int) []uint32 {
	if params == nil || nDefaults == 0 || nDefaults > len(argOIDs) {
		return nil
	}
	var types []uint32
	for _, item := range params.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if !ok || fp.Defexpr == nil {
			continue
		}
		analyzed, err := c.AnalyzeExprNoContext(fp.Defexpr)
		if err != nil || analyzed == nil {
			return nil
		}
		types = append(types, analyzed.exprType())
	}
	if len(types) != nDefaults {
		return nil
	}
	first := len(argOIDs) - nDefaults
	for i := range types {
		if pt := argOIDs[first+i]; !isPolymorphic(pt) {
			types[i] = pt
		}
	}
	return types
}

func (c *Catalog) paramDefaults(params *nodes.List, nInputArgs, nDefaults int) []string {
	if params == nil || nDefaults == 0 || nDefaults > nInputArgs {
		return nil
	}
	defaults := make([]string, nInputArgs)
	inputIdx := 0
	for _, item := range params.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if !ok {
			continue
		}
		if fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
			continue
		}
		if fp.Defexpr != nil && inputIdx < len(defaults) {
			analyzed, err := c.AnalyzeExprNoContext(fp.Defexpr)
			if err == nil && analyzed != nil {
				defaults[inputIdx] = c.DeparseExpr(analyzed, nil, false)
			}
		}
		inputIdx++
	}
	return defaults
}

func createFunctionStmtIsProcedure(stmt *nodes.CreateFunctionStmt) bool {
	if stmt == nil || stmt.Options == nil {
		return false
	}
	for _, item := range stmt.Options.Items {
		de, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		if strings.EqualFold(de.Defname, "isProcedure") {
			return defElemBool(de)
		}
	}
	return false
}

// recordFuncBodyDeps parses a SQL function body and records dependencies on
// tables referenced within it. Extracts relation names from the raw AST
// rather than doing full semantic analysis, because function bodies may
// reference function parameters that the analyzer doesn't know about.
// Silently ignores parse errors to avoid blocking catalog loading.
func (c *Catalog) recordFuncBodyDeps(funcOID uint32, body string) {
	stmts, err := parser.Parse(body)
	if err != nil || stmts == nil {
		return
	}
	for _, stmt := range stmts.Items {
		raw, ok := stmt.(*nodes.RawStmt)
		if !ok {
			continue
		}
		c.recordDepsFromNode('f', funcOID, raw.Stmt)
	}
}

// recordDepsFromNode walks an AST node tree and records dependencies for
// any RangeVar (table reference) found. This is a lightweight alternative
// to full semantic analysis that works even when the query contains
// unresolvable references (e.g., function parameters).
func (c *Catalog) recordDepsFromNode(objType byte, objOID uint32, n nodes.Node) {
	if n == nil {
		return
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		c.recordDepsOnRangeVar(objType, objOID, v)
	case *nodes.SelectStmt:
		if v.FromClause != nil {
			for _, from := range v.FromClause.Items {
				c.recordDepsFromNode(objType, objOID, from)
			}
		}
		c.recordDepsFromNode(objType, objOID, v.WhereClause)
		if v.Larg != nil {
			c.recordDepsFromNode(objType, objOID, v.Larg)
		}
		if v.Rarg != nil {
			c.recordDepsFromNode(objType, objOID, v.Rarg)
		}
		if v.WithClause != nil {
			for _, cte := range v.WithClause.Ctes.Items {
				c.recordDepsFromNode(objType, objOID, cte)
			}
		}
	case *nodes.CommonTableExpr:
		c.recordDepsFromNode(objType, objOID, v.Ctequery)
	case *nodes.JoinExpr:
		c.recordDepsFromNode(objType, objOID, v.Larg)
		c.recordDepsFromNode(objType, objOID, v.Rarg)
	case *nodes.RangeSubselect:
		c.recordDepsFromNode(objType, objOID, v.Subquery)
	case *nodes.SubLink:
		c.recordDepsFromNode(objType, objOID, v.Subselect)
	case *nodes.InsertStmt:
		c.recordDepsFromNode(objType, objOID, v.Relation)
		c.recordDepsFromNode(objType, objOID, v.SelectStmt)
	case *nodes.UpdateStmt:
		c.recordDepsFromNode(objType, objOID, v.Relation)
		if v.FromClause != nil {
			for _, from := range v.FromClause.Items {
				c.recordDepsFromNode(objType, objOID, from)
			}
		}
	case *nodes.DeleteStmt:
		c.recordDepsFromNode(objType, objOID, v.Relation)
	}
}

// recordDepsOnRangeVar resolves a RangeVar to a catalog relation and records
// a dependency if found.
func (c *Catalog) recordDepsOnRangeVar(objType byte, objOID uint32, rv *nodes.RangeVar) {
	if rv == nil {
		return
	}
	_, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
	if err == nil && rel != nil {
		c.recordDependency(objType, objOID, 0, 'r', rel.OID, 0, DepNormal)
	}
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
	c.recordTypeOwner(oid)
	return oid
}

// removeFunction removes a function from all catalog maps.
func (c *Catalog) removeFunction(procOID uint32, name string) {
	delete(c.userProcs, procOID)
	delete(c.procByOID, procOID)
	delete(c.extensionMembers, procOID)

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

	for key, cast := range c.castIndex {
		if cast.Func == procOID {
			delete(c.castIndex, key)
			delete(c.extensionCasts, key)
		}
	}

	c.removeGrantsForObject('f', procOID)
	c.removeComments('f', procOID)
	c.removeDepsOf('f', procOID)
	c.removeDepsOn('f', procOID)
}

// resolveRoutineTarget resolves an ObjectWithArgs to a user-defined routine,
// searching search_path when unqualified and enforcing routine-kind rules.
func (c *Catalog) resolveRoutineTarget(objType nodes.ObjectType, owa *nodes.ObjectWithArgs) (*Schema, *BuiltinProc, *UserProc, error) {
	kindStr := "function"
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		kindStr = "procedure"
	case nodes.OBJECT_AGGREGATE:
		kindStr = "aggregate"
	case nodes.OBJECT_ROUTINE:
		kindStr = "routine"
	}
	schemaName, funcName := qualifiedName(owa.Objname)
	argOIDs, argErr := c.routineArgOIDs(owa)
	if argErr != nil {
		return nil, nil, nil, argErr
	}
	var schema *Schema
	if schemaName == "" {
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
			return nil, nil, nil, err
		}
	}
	if schema == nil {
		return nil, nil, nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
	}
	var bp *BuiltinProc
	if argOIDs != nil {
		bp = c.findExactProc(schema, funcName, argOIDs)
	} else {
		candidates := c.findUserProcsByName(schema, funcName)
		if len(candidates) > 1 {
			return nil, nil, nil, errAmbiguousFunction(funcName)
		}
		if len(candidates) == 1 {
			bp = c.procByOID[candidates[0].OID]
		}
	}
	if bp == nil {
		return nil, nil, nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
	}
	up := c.userProcs[bp.OID]
	if up == nil {
		return nil, nil, nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%s %q is not user-defined", kindStr, funcName)}
	}
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		if bp.Kind != 'p' {
			return nil, nil, nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a procedure", funcName)}
		}
	case nodes.OBJECT_FUNCTION:
		if bp.Kind == 'p' {
			return nil, nil, nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a function", funcName)}
		}
		if bp.Kind == PROKIND_AGGREGATE {
			return nil, nil, nil, &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is an aggregate function\nHint: Use DROP AGGREGATE to drop aggregate functions.", funcName)}
		}
	case nodes.OBJECT_AGGREGATE:
		if bp.Kind != PROKIND_AGGREGATE {
			return nil, nil, nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an aggregate", funcName)}
		}
	}
	return schema, bp, up, nil
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

	objType := stmt.Objtype
	if objType == 0 {
		objType = nodes.OBJECT_ROUTINE
	}
	_, bp, up, err := c.resolveRoutineTarget(objType, stmt.Func)
	if err != nil {
		return err
	}

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
				raw, val, ok := defElemNumeric(de)
				if ok && val <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "COST must be positive"}
				}
				if ok && up != nil {
					up.Cost = raw
				}
			case "rows":
				// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1429-1439)
				raw, val, ok := defElemNumeric(de)
				if ok && val <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS must be positive"}
				}
				if !bp.RetSet {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "ROWS is not applicable when function does not return a set"}
				}
				if ok && up != nil {
					up.Rows = raw
				}
			case "leakproof":
				// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1413-1420)
				lp := defElemBool(de)
				bp.LeakProof = lp
				if up != nil {
					up.LeakProof = lp
				}
			case "set":
				if up != nil {
					if vs, ok := de.Arg.(*nodes.VariableSetStmt); ok {
						up.SetConfig = applyFuncSetConfig(up.SetConfig, vs)
					}
				}
			case "support":
				// No-op: not tracked in pgddl.
			}
		}
	}

	return nil
}

func defElemNumeric(d *nodes.DefElem) (string, float64, bool) {
	if d == nil || d.Arg == nil {
		return "", 0, false
	}
	extract := func(n nodes.Node) (string, float64, bool) {
		switch v := n.(type) {
		case *nodes.Integer:
			return strconv.FormatInt(v.Ival, 10), float64(v.Ival), true
		case *nodes.Float:
			f, err := strconv.ParseFloat(v.Fval, 64)
			if err != nil {
				return "", 0, false
			}
			return v.Fval, f, true
		}
		return "", 0, false
	}
	if ac, ok := d.Arg.(*nodes.A_Const); ok {
		return extract(ac.Val)
	}
	return extract(d.Arg)
}

func applyFuncSetConfig(existing []string, vs *nodes.VariableSetStmt) []string {
	if vs == nil {
		return existing
	}
	if vs.Kind == nodes.VAR_RESET_ALL {
		return nil
	}
	key := strings.ToLower(vs.Name)
	if key == "" {
		return existing
	}
	removeVar := func(items []string) []string {
		out := make([]string, 0, len(items))
		for _, item := range items {
			if name, ok := funcSetConfigVarName(item); ok && name == key {
				continue
			}
			out = append(out, item)
		}
		return out
	}
	out := removeVar(existing)
	switch vs.Kind {
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		return out
	case nodes.VAR_SET_CURRENT:
		return append(out, fmt.Sprintf("SET %s FROM CURRENT", quoteIdentAlways(vs.Name)))
	case nodes.VAR_SET_VALUE:
		val := formatFuncSetArgs(vs.Args)
		if val == "" {
			return out
		}
		return append(out, fmt.Sprintf("SET %s = %s", quoteIdentAlways(vs.Name), val))
	default:
		return out
	}
}

func funcSetConfigVarName(clause string) (string, bool) {
	rest, ok := strings.CutPrefix(clause, "SET ")
	if !ok {
		return "", false
	}
	var raw string
	if strings.HasPrefix(rest, `"`) {
		var b strings.Builder
		for i := 1; i < len(rest); i++ {
			if rest[i] == '"' {
				if i+1 < len(rest) && rest[i+1] == '"' {
					b.WriteByte('"')
					i++
					continue
				}
				return strings.ToLower(b.String()), true
			}
			b.WriteByte(rest[i])
		}
		return "", false
	}
	if idx := strings.IndexAny(rest, " ="); idx >= 0 {
		raw = rest[:idx]
	} else {
		raw = rest
	}
	return strings.ToLower(raw), raw != ""
}

func formatFuncSetArgs(args *nodes.List) string {
	if args == nil || len(args.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args.Items))
	for _, item := range args.Items {
		switch v := item.(type) {
		case *nodes.A_Const:
			switch val := v.Val.(type) {
			case *nodes.String:
				parts = append(parts, quoteLiteral(val.Str))
			case *nodes.Integer:
				parts = append(parts, strconv.FormatInt(val.Ival, 10))
			case *nodes.Float:
				parts = append(parts, val.Fval)
			case *nodes.Boolean:
				if val.Boolval {
					parts = append(parts, "true")
				} else {
					parts = append(parts, "false")
				}
			}
		case *nodes.String:
			parts = append(parts, quoteLiteral(v.Str))
		case *nodes.Integer:
			parts = append(parts, strconv.FormatInt(v.Ival, 10))
		case *nodes.Float:
			parts = append(parts, v.Fval)
		}
	}
	return strings.Join(parts, ", ")
}
