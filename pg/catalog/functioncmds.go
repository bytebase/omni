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
	OID        uint32
	Name       string
	Schema     *Schema
	OwnerName  string
	Kind       byte // 'f' or 'p'
	ArgTypes   []uint32
	RetType    uint32
	RetSet     bool
	Volatile   byte
	Parallel   byte
	IsStrict   bool
	SecDef     bool
	LeakProof  bool
	IsWindow   bool
	Cost       string
	Rows       string
	SetConfig  []string
	Support    string
	SupportOID uint32
	Language   string
	Body       string
	LinkSymbol string
	HasSqlBody bool // SQL-standard body (RETURN / BEGIN ATOMIC) without AS '...'
	// sqlBody and bodyBindings are a SQL-standard body as parsed and bound,
	// which procBody prints with the names of what it uses as they are now.
	sqlBody      nodes.Node
	bodyBindings *ruleBindings
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
	var language, body, linkSymbol, support string
	var supportOID uint32
	var volatile, parallel byte
	var isStrict, secDef, isLeakProof bool
	var isWindowFunc bool
	var costStr, rowsStr string
	var setConfig []string
	var hasAsClause, rowsSet, seenVolatility, seenStrict, seenSecurity, seenParallel, seenLeakproof, seenWindow, seenCost, seenRows, seenSupport, seenLanguage bool

	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "language":
				if seenLanguage {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenLanguage = true
				language = defElemString(de)
			case "as":
				if hasAsClause {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				hasAsClause = true
				switch v := de.Arg.(type) {
				case *nodes.String:
					body = v.Str
				case *nodes.List:
					if len(v.Items) > 0 {
						body = stringVal(v.Items[0])
					}
					if len(v.Items) > 1 {
						linkSymbol = stringVal(v.Items[1])
					}
				}
			case "volatility":
				// pg: compute_common_attribute — procedures cannot have volatility
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				if seenVolatility {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenVolatility = true
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
				if seenParallel {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenParallel = true
				switch strings.ToLower(defElemString(de)) {
				case "safe":
					parallel = 's'
				case "restricted":
					parallel = 'r'
				case "unsafe":
					parallel = 'u'
				default:
					return &Error{Code: CodeSyntaxError, Message: `parameter "parallel" must be SAFE, RESTRICTED, or UNSAFE`}
				}
			case "strict":
				// pg: compute_common_attribute — procedures cannot have strict
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				if seenStrict {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenStrict = true
				isStrict = defElemBool(de)
			case "called_on_null_input":
				// pg: compute_common_attribute — same as strict (inverse), restricted for procedures
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				if seenStrict {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenStrict = true
				isStrict = !defElemBool(de)
			case "security":
				// pg: SECURITY DEFINER → Boolean(true), SECURITY INVOKER → Boolean(false)
				// Also handle String("definer"/"invoker") for manually constructed AST.
				if seenSecurity {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenSecurity = true
				if s, ok := de.Arg.(*nodes.String); ok {
					secDef = strings.ToLower(s.Str) == "definer"
				} else {
					secDef = defElemBool(de)
				}
			case "leakproof":
				// pg: src/backend/commands/functioncmds.c — compute_common_attribute
				if seenLeakproof {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenLeakproof = true
				isLeakProof = defElemBool(de)
			case "window":
				// pg: compute_function_attributes — procedures cannot have WINDOW
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				if seenWindow {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenWindow = true
				isWindowFunc = defElemBool(de)
			case "cost":
				// pg: src/backend/commands/functioncmds.c — compute_function_attributes
				if seenCost {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenCost = true
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
				if seenRows {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenRows = true
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
					nextSet, err := c.applyFuncSetConfig(setConfig, vs)
					if err != nil {
						return err
					}
					setConfig = nextSet
				}
			case "support":
				if isProcedure {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
				if seenSupport {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenSupport = true
				sup, supOID, err := c.resolveRoutineRef(de, []uint32{INTERNALOID}, INTERNALOID, "internal")
				if err != nil {
					return err
				}
				support = sup
				supportOID = supOID
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
	if language == "" {
		return errInvalidFunctionDefinition("no language specified")
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
	sqlStandardBody := hasSqlBody && !hasAsClause

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

	if retOID == EVENTTRIGGEROID && !langIsSQL && len(allArgTypes) > 0 {
		return errInvalidFunctionDefinition("event trigger functions cannot have declared arguments")
	}

	// A SQL-standard body is parsed and analyzed now, before the routine
	// exists, and records what it uses.
	// pg: src/backend/commands/functioncmds.c — interpret_AS_clause
	fnParams := sqlFnInputParams(name, stmt.Parameters, argOIDs)
	var bodyBindings *ruleBindings
	if sqlStandardBody {
		if !langIsSQL {
			return errInvalidFunctionDefinition("inline SQL function body only valid for language SQL")
		}
		for _, t := range argOIDs {
			if isPolymorphic(t) {
				return errInvalidFunctionDefinition("SQL function with unquoted function body cannot have polymorphic arguments")
			}
		}
		bodyBindings, err = c.bindSQLFunctionBody(sqlBodyStmts(stmt.SqlBody), fnParams)
		if err != nil {
			return err
		}
		body = c.deparseSqlBody(stmt.SqlBody, bodyBindings)
	}

	// A SQL function may not return a pseudo-type but record, void or a
	// polymorphic one.
	// pg: src/backend/executor/functions.c — fmgr_sql_validator
	if langIsSQL && !isPolymorphic(retOID) && retOID != RECORDOID && retOID != VOIDOID {
		if bt := c.typeByOID[retOID]; bt != nil && bt.Type == 'p' {
			return errInvalidFunctionDefinition(fmt.Sprintf("SQL functions cannot return type %s", c.FormatType(retOID, -1)))
		}
	}
	var outTypes []uint32
	for i, mode := range allArgModes {
		if mode == 'o' || mode == 'b' || mode == 't' {
			outTypes = append(outTypes, allArgTypes[i])
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
	} else if isWindowFunc {
		kind = PROKIND_WINDOW
	}

	if volatile == 0 {
		volatile = 'v'
	}
	if parallel == 0 {
		parallel = 'u'
	}

	// validateBody runs the validator over a string body once the routine
	// exists, so the body may call it.
	// pg: src/backend/catalog/pg_proc.c — ProcedureCreate (OidFunctionCall1 of the validator)
	validateBody := func() error {
		if !c.checkFunctionBodies() {
			return nil
		}
		if strings.EqualFold(language, "plpgsql") {
			return c.checkPLpgSQL(body, &plpgsqlRoutine{
				name: name, kind: kind, retType: retOID, retSet: returnSet,
				hasOut: len(outTypes) > 0, params: append(allArgNames, make([]string, len(allArgTypes)-len(allArgNames))...),
			})
		}
		if !langIsSQL {
			return nil
		}
		if sqlStandardBody {
			return c.checkSQLFunctionResult(bodyBindings, sqlBodyStmts(stmt.SqlBody), kind, retOID, outTypes)
		}
		return c.checkSQLFunctionBody(body, fnParams, kind, retOID, outTypes)
	}

	argDefaults := c.paramDefaults(stmt.Parameters, len(argOIDs), int(nArgDefaults))
	defaultTypes := c.paramDefaultTypes(stmt.Parameters, argOIDs, int(nArgDefaults))

	// Check for existing function with same signature.
	existing := c.findExactProc(schema, name, argOIDs)
	if existing != nil {
		if !orReplace {
			return errDuplicateFunction(name)
		}
		if existing.Kind != kind {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf(
				"cannot change routine kind of existing function %q", name)}
		}
		// OrReplace: verify return type and setof flag match.
		if existing.RetType != retOID || existing.RetSet != returnSet {
			return errInvalidFunctionDefinition(fmt.Sprintf(
				"cannot change return type of existing function %q", name))
		}
		// Replace the body and attributes.
		up := c.userProcs[existing.OID]
		if up != nil {
			// A replacement may not rename an input parameter.
			// pg: src/backend/catalog/pg_proc.c — ProcedureCreate (proargnames)
			old := sqlFnInputParams(up.Name, nil, up.ArgTypes)
			old.names = old.names[:0]
			for i, n := range up.ArgNames {
				if up.ArgModes == nil || i < len(up.ArgModes) && up.ArgModes[i] != 'o' && up.ArgModes[i] != 't' {
					old.names = append(old.names, n)
				}
			}
			for i, n := range old.names {
				if n != "" && i < len(fnParams.names) && fnParams.names[i] != n {
					return errInvalidFunctionDefinition(fmt.Sprintf("cannot change name of input parameter %q", n))
				}
			}
		}
		if err := validateBody(); err != nil {
			return err
		}
		if up != nil {
			fnParams.procOID = existing.OID
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
			up.Support = support
			up.SupportOID = supportOID
			up.Language = language
			up.Body = body
			up.LinkSymbol = linkSymbol
			up.HasSqlBody = sqlStandardBody
			up.sqlBody, up.bodyBindings = nil, nil
			if sqlStandardBody {
				up.sqlBody, up.bodyBindings = stmt.SqlBody, bodyBindings
			}
		}
		existing.Volatile = volatile
		existing.Parallel = parallel
		existing.IsStrict = isStrict
		existing.SecDef = secDef
		existing.LeakProof = isLeakProof
		// Rebuild all normal dependencies (arg/return types, support, defaults, SQL body)
		// so dropped helpers/tables from replaced defaults or bodies are detached.
		c.removeDepsOf('f', existing.OID)
		c.recordFunctionDependencies(existing.OID, argOIDs, retOID, supportOID, stmt.Parameters, bodyBindings)
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
		OID:             procOID,
		Name:            name,
		Schema:          schema,
		OwnerName:       c.currentOwner(),
		Kind:            kind,
		ArgTypes:        argOIDs,
		RetType:         retOID,
		RetSet:          returnSet,
		Volatile:        volatile,
		Parallel:        parallel,
		IsStrict:        isStrict,
		SecDef:          secDef,
		LeakProof:       isLeakProof,
		IsWindow:        isWindowFunc,
		Cost:            costStr,
		Rows:            rowsStr,
		SetConfig:       slices.Clone(setConfig),
		Support:         support,
		SupportOID:      supportOID,
		Language:        language,
		Body:            body,
		LinkSymbol:      linkSymbol,
		HasSqlBody:      sqlStandardBody,
		bodyBindings:    bodyBindings,
		ArgNames:        allArgNames,
		ArgModes:        storedArgModes,
		AllArgTypes:     storedAllArgTypes,
		ArgDefaultTypes: defaultTypes,
		ArgDefaults:     argDefaults,
		NArgDefaults:    nArgDefaults,
	}
	if sqlStandardBody {
		up.sqlBody = stmt.SqlBody
	}
	fnParams.procOID = procOID
	c.userProcs[procOID] = up

	c.recordFunctionDependencies(procOID, argOIDs, retOID, supportOID, stmt.Parameters, bodyBindings)

	if err := validateBody(); err != nil {
		c.removeFunction(procOID, name)
		return err
	}
	return nil
}

func (c *Catalog) recordFunctionDependencies(
	procOID uint32,
	argOIDs []uint32,
	retOID uint32,
	supportOID uint32,
	params *nodes.List,
	body *ruleBindings,
) {
	// Record dependencies on arg/return/support types/routines.
	for _, aoid := range argOIDs {
		if aoid >= FirstNormalObjectId {
			c.recordDependency('f', procOID, 0, 't', aoid, 0, DepNormal)
		}
	}
	if retOID >= FirstNormalObjectId {
		c.recordDependency('f', procOID, 0, 't', retOID, 0, DepNormal)
	}
	if supportOID >= FirstNormalObjectId {
		c.recordDependency('f', procOID, 0, 'f', supportOID, 0, DepNormal)
	}

	// Record dependencies from default expressions.
	// pg: src/backend/commands/functioncmds.c — CreateFunction (recordDependencyOnExpr for defaults)
	if params != nil {
		for _, item := range params.Items {
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

	// A SQL-standard body depends on what it uses; a string body on
	// nothing.
	// pg: src/backend/catalog/pg_proc.c — ProcedureCreate
	//   if (languageObjectId == SQLlanguageId && prosqlbody)
	//       recordDependencyOnExpr(&myself, prosqlbody, NIL, DEPENDENCY_NORMAL);
	if body == nil {
		return
	}
	for _, relOID := range body.relDeps {
		c.recordDependency('f', procOID, 0, 'r', relOID, 0, DepNormal)
	}
	for _, cd := range body.colDeps {
		c.recordDependency('f', procOID, 0, 'r', cd.relOID, int32(cd.attNum), DepNormal)
	}
	for _, conOID := range body.conDeps {
		c.recordDependency('f', procOID, 0, 'c', conOID, 0, DepNormal)
	}
	for _, seqOID := range body.seqDeps {
		c.recordDependency('f', procOID, 0, 's', seqOID, 0, DepNormal)
	}
	for _, typeOID := range body.typeDeps {
		c.recordDependency('f', procOID, 0, 't', typeOID, 0, DepNormal)
	}
	for _, collOID := range body.collDeps {
		c.recordDependency('f', procOID, 0, 'L', collOID, 0, DepNormal)
	}
	for _, funcOID := range body.funcDeps {
		if funcOID != procOID {
			c.recordDependency('f', procOID, 0, 'f', funcOID, 0, DepNormal)
		}
	}
	for _, operOID := range body.operDeps {
		c.recordDependency('f', procOID, 0, 'O', operOID, 0, DepNormal)
	}
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

// bodyOrderingDeps returns, for each SQL function with a string body, an
// entry on each table its body names. PostgreSQL records no dependency for
// such a body, so a DROP ignores them; they only order a migration, which
// must create a table before a function check_function_bodies reads it in.
func (c *Catalog) bodyOrderingDeps() []DepEntry {
	var out []DepEntry
	for _, up := range c.userProcs {
		if up.HasSqlBody || !strings.EqualFold(up.Language, "sql") || up.Body == "" {
			continue
		}
		stmts, err := parser.Parse(up.Body)
		if err != nil || stmts == nil {
			continue
		}
		seen := map[uint32]bool{}
		for _, stmt := range stmts.Items {
			nodes.Inspect(stmt, func(n nodes.Node) bool {
				if rv, ok := n.(*nodes.RangeVar); ok && rv != nil {
					if _, rel, err := c.findRelation(rv.Schemaname, rv.Relname); err == nil && rel != nil && !seen[rel.OID] {
						seen[rel.OID] = true
						out = append(out, DepEntry{ObjType: 'f', ObjOID: up.OID, RefType: 'r', RefOID: rel.OID, DepType: DepNormal})
					}
				}
				return true
			})
		}
	}
	return out
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
	return c.resolveRoutineTargetOpts(objType, owa, true)
}

func (c *Catalog) resolveRoutineTargetOpts(objType nodes.ObjectType, owa *nodes.ObjectWithArgs, rejectAggForFunc bool) (*Schema, *BuiltinProc, *UserProc, error) {
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
	matchesObjKind := func(kind byte) bool {
		switch objType {
		case nodes.OBJECT_PROCEDURE:
			return kind == 'p'
		case nodes.OBJECT_FUNCTION:
			return kind != 'p'
		case nodes.OBJECT_AGGREGATE:
			return kind == PROKIND_AGGREGATE
		default:
			return true
		}
	}
	var schema *Schema
	if schemaName == "" {
		for _, nsOID := range c.routineSearchPathWithCatalog() {
			ns := c.schemas[nsOID]
			if ns == nil {
				continue
			}
			if argOIDs == nil {
				foundKind := false
				for _, up := range c.findUserProcsByName(ns, funcName) {
					if matchesObjKind(up.Kind) {
						foundKind = true
						break
					}
				}
				if foundKind {
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
		return nil, nil, nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("could not find a %s named %q", kindStr, funcName)}
	}
	var bp *BuiltinProc
	if argOIDs != nil {
		bp = c.findExactProc(schema, funcName, argOIDs)
	} else {
		var candidates []*UserProc
		for _, up := range c.findUserProcsByName(schema, funcName) {
			if matchesObjKind(up.Kind) {
				candidates = append(candidates, up)
			}
		}
		if len(candidates) > 1 {
			return nil, nil, nil, errAmbiguousFunction(funcName)
		}
		if len(candidates) == 1 {
			bp = c.procByOID[candidates[0].OID]
		}
	}
	if bp == nil {
		return nil, nil, nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("could not find a %s named %q", kindStr, funcName)}
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
		if rejectAggForFunc && bp.Kind == PROKIND_AGGREGATE {
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

	nextVolatile := bp.Volatile
	nextStrict := bp.IsStrict
	nextSecDef := bp.SecDef
	nextParallel := bp.Parallel
	nextLeakProof := bp.LeakProof
	var nextCost, nextRows, nextSupport string
	var nextSupportOID uint32
	var nextSetConfig []string
	var hasSupport bool
	if up != nil {
		nextCost = up.Cost
		nextRows = up.Rows
		nextSetConfig = slices.Clone(up.SetConfig)
		nextSupport = up.Support
		nextSupportOID = up.SupportOID
	}

	// Validate every action before mutating bp, up, or catalog dependencies.
	var seenVolatility, seenStrict, seenSecurity, seenParallel, seenCost, seenRows, seenLeakproof, seenSupport bool
	if stmt.Actions != nil {
		for _, item := range stmt.Actions.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			// pg: src/backend/commands/functioncmds.c — compute_common_attribute (procedure_error)
			switch strings.ToLower(de.Defname) {
			case "volatility", "strict", "called_on_null_input", "leakproof", "cost", "rows", "support", "parallel":
				if objType == nodes.OBJECT_PROCEDURE {
					return errInvalidFunctionDefinition("invalid attribute in procedure definition")
				}
			}
			switch strings.ToLower(de.Defname) {
			case "volatility":
				if seenVolatility {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenVolatility = true
				switch strings.ToLower(defElemString(de)) {
				case "immutable":
					nextVolatile = 'i'
				case "stable":
					nextVolatile = 's'
				case "volatile":
					nextVolatile = 'v'
				}
			case "strict":
				if seenStrict {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenStrict = true
				nextStrict = defElemBool(de)
			case "called_on_null_input":
				if seenStrict {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenStrict = true
				nextStrict = !defElemBool(de)
			case "security":
				if seenSecurity {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenSecurity = true
				// pg: SECURITY DEFINER → Boolean(true), SECURITY INVOKER → Boolean(false)
				if s, ok := de.Arg.(*nodes.String); ok {
					nextSecDef = strings.ToLower(s.Str) == "definer"
				} else {
					nextSecDef = defElemBool(de)
				}
			case "parallel":
				if seenParallel {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenParallel = true
				switch strings.ToLower(defElemString(de)) {
				case "safe":
					nextParallel = 's'
				case "restricted":
					nextParallel = 'r'
				case "unsafe":
					nextParallel = 'u'
				default:
					return &Error{Code: CodeSyntaxError, Message: `parameter "parallel" must be SAFE, RESTRICTED, or UNSAFE`}
				}
			case "cost":
				if seenCost {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenCost = true
				// pg: src/backend/commands/functioncmds.c — AlterFunction
				raw, val, ok := defElemNumeric(de)
				if ok && val <= 0 {
					return &Error{Code: CodeInvalidParameterValue,
						Message: "COST must be positive"}
				}
				if ok {
					nextCost = raw
				}
			case "rows":
				if seenRows {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenRows = true
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
				if ok {
					nextRows = raw
				}
			case "leakproof":
				if seenLeakproof {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenLeakproof = true
				// pg: src/backend/commands/functioncmds.c — AlterFunction (line 1413-1420)
				nextLeakProof = defElemBool(de)
			case "set":
				if up != nil {
					if vs, ok := de.Arg.(*nodes.VariableSetStmt); ok {
						nextSet, err := c.applyFuncSetConfig(nextSetConfig, vs)
						if err != nil {
							return err
						}
						nextSetConfig = nextSet
					}
				}
			case "support":
				if seenSupport {
					return &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
				}
				seenSupport = true
				sup, supOID, err := c.resolveRoutineRef(de, []uint32{INTERNALOID}, INTERNALOID, "internal")
				if err != nil {
					return err
				}
				nextSupport = sup
				nextSupportOID = supOID
				hasSupport = true
			}
		}
	}

	bp.Volatile = nextVolatile
	bp.IsStrict = nextStrict
	bp.SecDef = nextSecDef
	bp.Parallel = nextParallel
	bp.LeakProof = nextLeakProof
	if up != nil {
		up.Volatile = nextVolatile
		up.IsStrict = nextStrict
		up.SecDef = nextSecDef
		up.Parallel = nextParallel
		up.LeakProof = nextLeakProof
		up.Cost = nextCost
		up.Rows = nextRows
		up.SetConfig = nextSetConfig
		if hasSupport {
			oldSupportOID := up.SupportOID
			up.Support = nextSupport
			up.SupportOID = nextSupportOID
			c.updateSupportDep(up.OID, oldSupportOID, nextSupportOID)
		}
	}

	return nil
}

func (c *Catalog) updateSupportDep(procOID, oldSupportOID, newSupportOID uint32) {
	if procOID == 0 || oldSupportOID == newSupportOID {
		return
	}
	if oldSupportOID >= FirstNormalObjectId {
		for i, d := range c.deps {
			if d.ObjType == 'f' && d.ObjOID == procOID && d.RefType == 'f' && d.RefOID == oldSupportOID && d.DepType == DepNormal {
				c.deps = append(c.deps[:i], c.deps[i+1:]...)
				break
			}
		}
	}
	if newSupportOID >= FirstNormalObjectId {
		c.recordDependency('f', procOID, 0, 'f', newSupportOID, 0, DepNormal)
	}
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

func (c *Catalog) applyFuncSetConfig(existing []string, vs *nodes.VariableSetStmt) ([]string, error) {
	if vs == nil {
		return existing, nil
	}
	if vs.Kind == nodes.VAR_RESET_ALL {
		return nil, nil
	}
	key := strings.ToLower(vs.Name)
	if key == "" {
		return existing, nil
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
	case nodes.VAR_SET_DEFAULT:
		return out, nil
	case nodes.VAR_RESET:
		if len(out) == len(existing) && !isKnownPGConfigParam(key) {
			return nil, &Error{
				Code:    CodeUndefinedObject,
				Message: fmt.Sprintf("unrecognized configuration parameter %q", vs.Name),
			}
		}
		return out, nil
	case nodes.VAR_SET_CURRENT:
		resolved, ok := c.resolveCurrentSettingValue(key)
		if !ok {
			return nil, &Error{
				Code:    CodeFeatureNotSupported,
				Message: fmt.Sprintf("untracked configuration parameter %q in SET ... FROM CURRENT", vs.Name),
			}
		}
		return append(out, fmt.Sprintf("SET %s = %s", quoteIdentAlways(vs.Name), resolved)), nil
	case nodes.VAR_SET_VALUE:
		val := formatFuncSetArgs(vs.Args)
		if val == "" {
			return out, nil
		}
		return append(out, fmt.Sprintf("SET %s = %s", quoteIdentAlways(vs.Name), val)), nil
	default:
		return out, nil
	}
}

func isKnownPGConfigParam(key string) bool {
	return strings.Contains(key, ".") || pg17ConfigParams[key]
}

func (c *Catalog) resolveCurrentSettingValue(key string) (string, bool) {
	if c == nil {
		return "", false
	}
	switch key {
	case "search_path":
		if len(c.searchPath) == 0 {
			return "'$user', 'public'", true
		}
		parts := make([]string, 0, len(c.searchPath))
		for _, s := range c.searchPath {
			parts = append(parts, quoteLiteral(s))
		}
		return strings.Join(parts, ", "), true
	case "default_tablespace":
		return quoteLiteral(c.defaultTablespace), true
	case "default_table_access_method":
		if c.defaultTableAccessMethod == "" {
			return "'heap'", true
		}
		return quoteLiteral(c.defaultTableAccessMethod), true
	case "role":
		if c.currentUser == "" {
			return "'none'", true
		}
		return quoteLiteral(c.currentUser), true
	}
	if val, ok := c.sessionSettings[key]; ok {
		return val, true
	}
	return "", false
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

func (c *Catalog) deparseSqlBody(sqlBody nodes.Node, b *ruleBindings) string {
	if sqlBody == nil {
		return ""
	}
	d := &ruleDeparser{cat: c, b: b}
	switch n := sqlBody.(type) {
	case *nodes.ReturnStmt:
		return "RETURN " + d.deparseExpr(n.Returnval)
	case *nodes.List:
		stmts := n.Items
		if len(n.Items) == 1 {
			if inner, ok := n.Items[0].(*nodes.List); ok {
				stmts = inner.Items
			}
		}
		var b strings.Builder
		b.WriteString("BEGIN ATOMIC\n")
		for _, item := range stmts {
			var stmtSQL string
			if ret, ok := item.(*nodes.ReturnStmt); ok {
				stmtSQL = "RETURN " + d.deparseExpr(ret.Returnval)
			} else {
				stmtSQL = d.deparseActionStmt(item)
			}
			if stmtSQL != "" {
				b.WriteByte(' ')
				b.WriteString(stmtSQL)
				b.WriteString(";\n")
			}
		}
		b.WriteString("END")
		return b.String()
	default:
		return d.deparseActionStmt(sqlBody)
	}
}
