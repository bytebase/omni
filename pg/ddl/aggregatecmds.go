package ddl

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// DefineAggregate creates a user-defined aggregate function.
//
// pg: src/backend/commands/aggregatecmds.c — DefineAggregate
func (c *Catalog) DefineAggregate(stmt *nodes.DefineStmt) error {
	schemaName, aggName := qualifiedName(stmt.Defnames)
	if aggName == "" {
		return errInvalidParameterValue("no aggregate name specified")
	}

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Parse parameters from the DefElem list.
	// pg: src/backend/commands/aggregatecmds.c — DefineAggregate (line 95-230)
	var (
		transfuncName  string
		finalfuncName  string
		combinefuncName string
		serialfuncName  string
		deserialfuncName string
		mtransfuncName  string
		minvtransfuncName string
		mfinalfuncName  string
		transTypeName   *nodes.TypeName
		msTransTypeName *nodes.TypeName
		baseTypeName    *nodes.TypeName
		initCond        string
		mInitCond       string
		sortOperatorName string
		parallelStr     string
	)

	if stmt.Definition != nil {
		for _, item := range stmt.Definition.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "sfunc", "sfunc1":
				transfuncName = defElemQualifiedName(de)
			case "finalfunc":
				finalfuncName = defElemQualifiedName(de)
			case "combinefunc":
				combinefuncName = defElemQualifiedName(de)
			case "serialfunc":
				serialfuncName = defElemQualifiedName(de)
			case "deserialfunc":
				deserialfuncName = defElemQualifiedName(de)
			case "msfunc":
				mtransfuncName = defElemQualifiedName(de)
			case "minvfunc":
				minvtransfuncName = defElemQualifiedName(de)
			case "mfinalfunc":
				mfinalfuncName = defElemQualifiedName(de)
			case "finalfunc_extra", "mfinalfunc_extra":
				// Accepted, not tracked.
			case "finalfunc_modify", "mfinalfunc_modify":
				// Accepted, not tracked.
			case "sortop":
				sortOperatorName = defElemQualifiedName(de)
			case "basetype":
				if tn, ok := de.Arg.(*nodes.TypeName); ok {
					baseTypeName = tn
				}
			case "stype", "stype1":
				if tn, ok := de.Arg.(*nodes.TypeName); ok {
					transTypeName = tn
				}
			case "sspace":
				// Accepted, not tracked.
			case "mstype":
				if tn, ok := de.Arg.(*nodes.TypeName); ok {
					msTransTypeName = tn
				}
			case "msspace":
				// Accepted, not tracked.
			case "initcond", "initcond1":
				initCond = defElemString(de)
			case "minitcond":
				mInitCond = defElemString(de)
			case "parallel":
				parallelStr = strings.ToLower(defElemString(de))
			case "hypothetical":
				// Accepted, not tracked.
			default:
				c.addWarning(CodeSyntaxError,
					fmt.Sprintf("aggregate attribute \"%s\" not recognized", de.Defname))
			}
		}
	}

	// Suppress unused warnings for stored-but-not-tracked parameters.
	_ = combinefuncName
	_ = serialfuncName
	_ = deserialfuncName
	_ = mtransfuncName
	_ = minvtransfuncName
	_ = mfinalfuncName
	_ = msTransTypeName
	_ = initCond
	_ = mInitCond
	_ = sortOperatorName

	// ---------------------------------------------------------------
	// Resolve argument types.
	// pg: src/backend/commands/aggregatecmds.c — DefineAggregate (line 235-320)
	// ---------------------------------------------------------------
	var argOIDs []uint32

	if stmt.Oldstyle || stmt.Args == nil {
		// Old-style syntax: use basetype DefElem.
		// pg: line 241-267
		if baseTypeName != nil {
			_, baseName := typeNameParts(baseTypeName)
			if strings.ToLower(baseName) != "any" {
				bt := convertTypeNameToInternal(baseTypeName)
				oid, _, err := c.ResolveType(bt)
				if err != nil {
					return err
				}
				argOIDs = append(argOIDs, oid)
			}
			// If basetype is "ANY", argOIDs stays empty → agg(*).
		}
	} else {
		// New-style syntax: resolve arg types from stmt.Args.
		// pg: line 270-320
		// pgparser encodes Args as [*List(of FunctionParameter), *Integer(-1 sentinel)].
		// The first element is the actual arg list.
		for _, item := range stmt.Args.Items {
			switch v := item.(type) {
			case *nodes.List:
				for _, inner := range v.Items {
					oid, err := c.resolveAggArg(inner)
					if err != nil {
						return err
					}
					if oid != 0 {
						argOIDs = append(argOIDs, oid)
					}
				}
			case *nodes.FunctionParameter:
				oid, err := c.resolveAggArg(v)
				if err != nil {
					return err
				}
				if oid != 0 {
					argOIDs = append(argOIDs, oid)
				}
			case *nodes.TypeName:
				argType := convertTypeNameToInternal(v)
				oid, _, err := c.ResolveType(argType)
				if err != nil {
					return err
				}
				argOIDs = append(argOIDs, oid)
			case *nodes.Integer:
				// ORDER BY sentinel (-1), skip.
			}
		}
	}

	// ---------------------------------------------------------------
	// Validation.
	// pg: src/backend/commands/aggregatecmds.c — DefineAggregate (line 322-340)
	// ---------------------------------------------------------------
	if transTypeName == nil {
		return errInvalidObjectDefinition("aggregate stype must be specified")
	}
	if transfuncName == "" {
		return errInvalidObjectDefinition("aggregate sfunc must be specified")
	}

	// Resolve state type.
	transType := convertTypeNameToInternal(transTypeName)
	transTypeId, _, err := c.ResolveType(transType)
	if err != nil {
		return err
	}

	// ---------------------------------------------------------------
	// Determine return type.
	// pg: src/backend/commands/aggregatecmds.c — AggregateCreate (line 392-410)
	// If finalfunc is specified, use its return type. Otherwise, retType = transType.
	// ---------------------------------------------------------------
	retType := transTypeId
	if finalfuncName != "" {
		// Try to find the finalfunc by name and infer its return type.
		if procs := c.procByName[finalfuncName]; len(procs) > 0 {
			retType = procs[0].RetType
		}
	}

	// Map parallel safety.
	var parallel byte
	switch parallelStr {
	case "safe":
		parallel = 's'
	case "restricted":
		parallel = 'r'
	default:
		parallel = 'u'
	}

	// ---------------------------------------------------------------
	// Register as BuiltinProc with Kind='a'.
	// pg: src/backend/catalog/pg_aggregate.c — AggregateCreate → ProcedureCreate
	// ---------------------------------------------------------------
	procOID := c.oidGen.Next()
	bp := &BuiltinProc{
		OID:      procOID,
		Name:     aggName,
		Kind:     PROKIND_AGGREGATE,
		NArgs:    int16(len(argOIDs)),
		RetType:  retType,
		ArgTypes: argOIDs,
		IsStrict: false, // aggregates are not strict by default
		Volatile: 'i',   // aggregates are immutable
		Parallel: parallel,
	}
	c.procByOID[procOID] = bp
	c.procByName[aggName] = append(c.procByName[aggName], bp)

	// Create UserProc metadata for dependency tracking.
	// pg: src/backend/catalog/pg_aggregate.c — AggregateCreate → ProcedureCreate
	// PG stores prosrc = "aggregate_dummy" for all aggregates.
	up := &UserProc{
		OID:      procOID,
		Name:     aggName,
		Schema:   schema,
		Kind:     PROKIND_AGGREGATE,
		ArgTypes: argOIDs,
		RetType:  retType,
		Volatile: 'i',
		Parallel: parallel,
		Language: "internal",
		Body:     "aggregate_dummy",
	}
	c.userProcs[procOID] = up

	// Record dependencies.
	if transTypeId >= FirstNormalObjectId {
		c.recordDependency('f', procOID, 0, 't', transTypeId, 0, DepNormal)
	}
	for _, aoid := range argOIDs {
		if aoid >= FirstNormalObjectId {
			c.recordDependency('f', procOID, 0, 't', aoid, 0, DepNormal)
		}
	}

	// Record dependency on transition function.
	// pg: src/backend/catalog/pg_aggregate.c — AggregateCreate (recordDependencyOn for transfn)
	if transfuncName != "" {
		transfnArgs := append([]uint32{transTypeId}, argOIDs...)
		if bp := c.findExactProc(schema, transfuncName, transfnArgs); bp != nil {
			if bp.OID >= FirstNormalObjectId {
				c.recordDependency('f', procOID, 0, 'f', bp.OID, 0, DepNormal)
			}
		}
	}

	return nil
}

// resolveAggArg resolves a single aggregate argument (FunctionParameter or TypeName) to a type OID.
//
// (pgddl helper)
func (c *Catalog) resolveAggArg(item nodes.Node) (uint32, error) {
	switch v := item.(type) {
	case *nodes.FunctionParameter:
		argType := convertTypeNameToInternal(v.ArgType)
		oid, _, err := c.ResolveType(argType)
		return oid, err
	case *nodes.TypeName:
		argType := convertTypeNameToInternal(v)
		oid, _, err := c.ResolveType(argType)
		return oid, err
	}
	return 0, nil
}
