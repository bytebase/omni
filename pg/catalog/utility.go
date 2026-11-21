package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ProcessUtility dispatches a utility (DDL) statement to the appropriate handler.
//
// pg: src/backend/tcop/utility.c — standard_ProcessUtility
func (c *Catalog) ProcessUtility(stmt nodes.Node) error {
	switch s := stmt.(type) {
	case *nodes.CreateSchemaStmt:
		return c.CreateSchemaCommand(s)
	case *nodes.CreateStmt:
		return c.DefineRelation(s, 'r')
	case *nodes.ViewStmt:
		return c.DefineView(s)
	case *nodes.IndexStmt:
		return c.DefineIndex(s)
	case *nodes.CreateSeqStmt:
		return c.DefineSequence(s)
	case *nodes.AlterSeqStmt:
		return c.AlterSequenceStmt(s)
	case *nodes.CreateEnumStmt:
		return c.DefineEnum(s)
	case *nodes.AlterEnumStmt:
		return c.AlterEnumStmt(s)
	case *nodes.CreateDomainStmt:
		return c.DefineDomain(s)
	case *nodes.AlterDomainStmt:
		return c.AlterDomainStmt(s)
	case *nodes.CreateFunctionStmt:
		return c.CreateFunctionStmt(s)
	case *nodes.CreateTrigStmt:
		return c.CreateTriggerStmt(s)
	case *nodes.CommentStmt:
		return c.CommentObject(s)
	case *nodes.CompositeTypeStmt:
		return c.DefineCompositeType(s)
	case *nodes.CreateRangeStmt:
		return c.DefineRange(s)
	case *nodes.CreateTableAsStmt:
		objType := nodes.ObjectType(s.Objtype)
		if objType == nodes.OBJECT_MATVIEW || objType == nodes.OBJECT_TABLE {
			return c.ExecCreateTableAs(s)
		}
		return errInvalidParameterValue(fmt.Sprintf("unsupported CreateTableAsStmt objtype: %d", s.Objtype))
	case *nodes.RefreshMatViewStmt:
		return c.ExecRefreshMatView(s)
	case *nodes.DropStmt:
		return c.RemoveObjects(s)
	case *nodes.AlterTableStmt:
		return c.AlterTableStmt(s)
	case *nodes.RenameStmt:
		return c.ExecRenameStmt(s)
	case *nodes.GrantStmt:
		return c.ExecGrantStmt(s)
	case *nodes.CreatePolicyStmt:
		return c.CreatePolicy(s)
	case *nodes.AlterPolicyStmt:
		return c.AlterPolicy(s)
	case *nodes.AlterFunctionStmt:
		return c.AlterFunction(s)
	case *nodes.AlterObjectSchemaStmt:
		return c.ExecAlterObjectSchemaStmt(s)
	case *nodes.AlterOwnerStmt:
		return c.ExecAlterOwnerStmt(s)
	case *nodes.TruncateStmt:
		return c.ExecuteTruncate(s)
	case *nodes.AlterTypeStmt:
		// ALTER TYPE ... SET (...) / DROP ATTRIBUTE: no-op for pgddl
		// pg: src/backend/commands/typecmds.c — AlterType
		return nil
	case *nodes.DefineStmt:
		// pg: src/backend/tcop/utility.c — ProcessUtilitySlow (DefineStmt cases)
		switch s.Kind {
		case nodes.OBJECT_TYPE:
			return c.DefineType(s)
		case nodes.OBJECT_OPERATOR:
			return c.DefineOperator(s)
		case nodes.OBJECT_AGGREGATE:
			return c.DefineAggregate(s)
		case nodes.OBJECT_COLLATION:
			schemaName, name := qualifiedName(s.Defnames)
			schema, err := c.resolveTargetSchema(schemaName)
			if err != nil {
				return err
			}
			defClause, err := c.validateAndFormatCollationDef(s.Definition)
			if err != nil {
				return err
			}
			if name != "" {
				if (schema.Collations != nil && schema.Collations[name] != nil) || (schema.OID == PGCatalogNamespace && isBuiltinCollationName(name)) {
					if s.IfNotExists {
						c.addWarning(CodeWarningSkip, fmt.Sprintf("collation %q already exists, skipping", name))
						return nil
					}
					return errDuplicateObject("collation", name)
				}
				if schema.Collations == nil {
					schema.Collations = make(map[string]*UserCollation)
				}
				schema.Collations[name] = &UserCollation{
					OID:       c.oidGen.Next(),
					Name:      name,
					DefClause: defClause,
					OwnerName: c.currentOwner(),
				}
			}
			return nil
		default:
			// CREATE TEXT SEARCH/etc: no-op for pgddl
			return nil
		}
	case *nodes.RuleStmt:
		return c.CreateRuleStmt(s)
	case *nodes.CreateCastStmt:
		return c.CreateCast(s)
	case *nodes.CreateForeignTableStmt:
		// CREATE FOREIGN TABLE: delegate to DefineRelation with relkind='f'
		// pg: src/backend/commands/tablecmds.c — DefineRelation (foreign table)
		var schemaName, relName string
		if s.Base.Relation != nil {
			schemaName, relName = s.Base.Relation.Schemaname, s.Base.Relation.Relname
		}
		targetSchema, _ := c.resolveTargetSchema(schemaName)
		existed := s.Base.IfNotExists && targetSchema != nil && targetSchema.Relations[relName] != nil
		if !existed && !c.allowImplicitForeignServer {
			if s.Servername == "" || c.foreignServers[s.Servername] == nil {
				return errUndefinedObject("server", s.Servername)
			}
		}
		if err := c.DefineRelation(&s.Base, 'f'); err != nil {
			return err
		}
		if !existed {
			var rel *Relation
			if targetSchema != nil {
				rel = targetSchema.Relations[relName]
			} else {
				rel = c.GetRelation(schemaName, relName)
			}
			if rel != nil && rel.RelKind == 'f' {
				rel.ForeignServer = s.Servername
				if s.Options != nil {
					for _, item := range s.Options.Items {
						if de, ok := item.(*nodes.DefElem); ok {
							rel.ForeignOptions = append(rel.ForeignOptions, de.Defname+"="+defElemString(de))
						}
					}
				}
			}
		}
		return nil
	case *nodes.CreateStatsStmt:
		// CREATE STATISTICS: no-op for pgddl (planner hint)
		return nil
	case *nodes.AlterStatsStmt:
		// ALTER STATISTICS: no-op for pgddl (planner hint)
		return nil
	case *nodes.AlterDefaultPrivilegesStmt:
		// ALTER DEFAULT PRIVILEGES: no-op for pgddl (no ACL tracking)
		return nil
	case *nodes.AlterTSDictionaryStmt:
		// ALTER TEXT SEARCH DICTIONARY: no-op for pgddl
		return nil
	case *nodes.AlterTSConfigurationStmt:
		// ALTER TEXT SEARCH CONFIGURATION: no-op for pgddl
		return nil
	case *nodes.CreateExtensionStmt:
		return c.CreateExtension(s)
	case *nodes.ReindexStmt:
		// REINDEX: no-op for pgddl (physical index rebuild)
		// pg: src/backend/commands/indexcmds.c — ReindexTable etc.
		return nil
	case *nodes.CreateOpClassStmt:
		return c.DefineOpClass(s)
	case *nodes.CreateOpFamilyStmt:
		return c.CreateOpFamily(s)
	case *nodes.AlterOpFamilyStmt:
		return c.AlterOpFamily(s)
	case *nodes.CreateConversionStmt:
		// CREATE CONVERSION: no-op for pgddl
		return nil
	case *nodes.CreateAmStmt:
		return c.CreateAccessMethod(s)
	case *nodes.CreatePublicationStmt:
		// CREATE PUBLICATION: no-op for pgddl
		return nil
	case *nodes.AlterPublicationStmt:
		// ALTER PUBLICATION: no-op for pgddl
		return nil
	case *nodes.CreateSubscriptionStmt:
		// CREATE SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.AlterSubscriptionStmt:
		// ALTER SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.DropSubscriptionStmt:
		// DROP SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.CreateTransformStmt:
		// CREATE TRANSFORM: no-op for pgddl
		return nil
	case *nodes.CreateEventTrigStmt:
		return c.CreateEventTriggerStmt(s)
	case *nodes.AlterEventTrigStmt:
		return c.AlterEventTrigger(s)
	case *nodes.VariableSetStmt:
		// pg: src/backend/tcop/utility.c — standard_ProcessUtility (VariableSetStmt)
		return c.processVariableSet(s)
	case *nodes.TransactionStmt:
		// BEGIN/COMMIT/ROLLBACK: no-op for pgddl
		return nil
	// ---- Foreign data wrappers ----
	case *nodes.AlterExtensionStmt:
		// ALTER EXTENSION: no-op for pgddl
		return nil
	case *nodes.AlterExtensionContentsStmt:
		// ALTER EXTENSION ... ADD/DROP: no-op for pgddl
		return nil
	case *nodes.CreateFdwStmt:
		return c.CreateFdwStmt(s)
	case *nodes.AlterFdwStmt:
		return c.AlterFdwStmt(s)
	case *nodes.CreateForeignServerStmt:
		return c.CreateForeignServerStmt(s)
	case *nodes.AlterForeignServerStmt:
		return c.AlterForeignServerStmt(s)
	case *nodes.CreateUserMappingStmt:
		// CREATE USER MAPPING: no-op for pgddl
		return nil
	case *nodes.AlterUserMappingStmt:
		// ALTER USER MAPPING: no-op for pgddl
		return nil
	case *nodes.DropUserMappingStmt:
		// DROP USER MAPPING: no-op for pgddl
		return nil
	case *nodes.ImportForeignSchemaStmt:
		// IMPORT FOREIGN SCHEMA: no-op for pgddl
		return nil

	// ---- Operators/Collations ----
	case *nodes.AlterCollationStmt:
		if _, _, err := c.resolveCollationList(s.Collname); err != nil {
			return err
		}
		return nil
	case *nodes.AlterOperatorStmt:
		// ALTER OPERATOR: no-op for pgddl
		return nil

	// ---- Tablespaces ----
	case *nodes.CreateTableSpaceStmt:
		return c.DefineTablespace(s)
	case *nodes.DropTableSpaceStmt:
		return c.RemoveTablespace(s)
	case *nodes.AlterTableSpaceOptionsStmt:
		return c.AlterTablespaceOptions(s)
	case *nodes.AlterTableMoveAllStmt:
		return c.MoveAllTablespaces(s)

	// ---- Roles/Databases ----
	case *nodes.CreateRoleStmt:
		return c.DefineRole(s)
	case *nodes.AlterRoleStmt:
		// ALTER ROLE: no-op for pgddl
		return nil
	case *nodes.AlterRoleSetStmt:
		// ALTER ROLE ... SET: no-op for pgddl
		return nil
	case *nodes.DropRoleStmt:
		return c.RemoveRoles(s)
	case *nodes.GrantRoleStmt:
		// GRANT role TO role: no-op for pgddl
		return nil
	case *nodes.CreatedbStmt:
		// CREATE DATABASE: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseStmt:
		// ALTER DATABASE: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseSetStmt:
		// ALTER DATABASE ... SET: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseRefreshCollStmt:
		// ALTER DATABASE ... REFRESH COLLATION VERSION: no-op for pgddl
		return nil
	case *nodes.DropdbStmt:
		// DROP DATABASE: no-op for pgddl
		return nil

	// ---- Misc DDL ----
	case *nodes.AlterObjectDependsStmt:
		// ALTER ... DEPENDS ON EXTENSION: no-op for pgddl
		return nil
	case *nodes.CreatePLangStmt:
		// CREATE LANGUAGE: no-op for pgddl
		return nil
	case *nodes.AlterSystemStmt:
		// ALTER SYSTEM: no-op for pgddl
		return nil
	case *nodes.DropOwnedStmt:
		return c.DropOwned(s)
	case *nodes.ReassignOwnedStmt:
		return c.ReassignOwned(s)
	case *nodes.SecLabelStmt:
		// SECURITY LABEL: no-op for pgddl
		return nil

	// ---- Non-DDL utility statements ----
	case *nodes.CopyStmt:
		// COPY: no-op for pgddl
		return nil
	case *nodes.DoStmt:
		// DO: no-op for pgddl
		return nil
	case *nodes.ExplainStmt:
		// EXPLAIN: no-op for pgddl
		return nil
	case *nodes.VacuumStmt:
		// VACUUM/ANALYZE: no-op for pgddl
		return nil
	case *nodes.ClusterStmt:
		// CLUSTER: no-op for pgddl
		return nil
	case *nodes.CheckPointStmt:
		// CHECKPOINT: no-op for pgddl
		return nil
	case *nodes.DiscardStmt:
		switch s.Target {
		case nodes.DISCARD_TEMP:
			c.discardTempObjects()
		case nodes.DISCARD_ALL:
			clear(c.preparedStmts)
			c.SetSearchPath([]string{"$user", "public"})
			c.defaultTablespace = ""
			c.defaultTableAccessMethod = ""
			c.ResetRole()
			c.discardTempObjects()
		}
		return nil
	case *nodes.LockStmt:
		// LOCK TABLE: no-op for pgddl
		return nil
	case *nodes.ListenStmt:
		// LISTEN: no-op for pgddl
		return nil
	case *nodes.UnlistenStmt:
		// UNLISTEN: no-op for pgddl
		return nil
	case *nodes.NotifyStmt:
		// NOTIFY: no-op for pgddl
		return nil
	case *nodes.LoadStmt:
		// LOAD: no-op for pgddl
		return nil
	case *nodes.ConstraintsSetStmt:
		// SET CONSTRAINTS: no-op for pgddl
		return nil
	case *nodes.VariableShowStmt:
		// SHOW: no-op for pgddl
		return nil
	case *nodes.CallStmt:
		// CALL: no-op for pgddl
		return nil
	case *nodes.PrepareStmt:
		if s.Name != "" && c.preparedStmts[s.Name] != nil {
			return &Error{
				Code:    CodeDuplicatePstatement,
				Message: fmt.Sprintf("prepared statement %q already exists", s.Name),
			}
		}
		var argTypes []uint32
		fixedParams := s.Argtypes != nil && len(s.Argtypes.Items) > 0
		if s.Argtypes != nil {
			argTypes = make([]uint32, 0, len(s.Argtypes.Items))
			for _, item := range s.Argtypes.Items {
				tn, ok := item.(*nodes.TypeName)
				if !ok {
					continue
				}
				oid, _, err := c.resolveTypeName(tn)
				if err != nil {
					return err
				}
				argTypes = append(argTypes, oid)
			}
		}
		numParams := len(argTypes)
		maxParam := maxParamRefInNode(s.Query)
		const maxPreparedParams = 65535
		if fixedParams && maxParam > numParams && !c.snapshot {
			return &Error{
				Code:    CodeUndefinedParameter,
				Message: fmt.Sprintf("there is no parameter $%d", maxParam),
			}
		}
		if maxParam > maxPreparedParams {
			return &Error{
				Code:    CodeUndefinedParameter,
				Message: fmt.Sprintf("there is no parameter $%d", maxParam),
			}
		}
		if !fixedParams && maxParam > numParams {
			numParams = maxParam
		}
		paramTypes := make([]uint32, numParams)
		copy(paramTypes, argTypes)
		var resultCols []preparedResultCol
		var hasResultCols bool
		if sel, isSel := s.Query.(*nodes.SelectStmt); isSel {
			q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: paramTypes})
			if err != nil {
				return err
			}
			if err := collectQueryParamTypes(q, paramTypes, fixedParams); err != nil && !c.snapshot {
				return err
			}
			for i := range paramTypes {
				if (paramTypes[i] == 0 || paramTypes[i] == UNKNOWNOID) && !c.snapshot {
					return &Error{
						Code:    CodeIndeterminateDatatype,
						Message: fmt.Sprintf("could not determine data type of parameter $%d", i+1),
					}
				}
			}
			if q != nil {
				hasResultCols = true
				for _, te := range nonJunkTargetEntries(q.TargetList) {
					resultCols = append(resultCols, preparedResultCol{
						TypeOID: te.Expr.exprType(),
						TypeMod: te.Expr.exprTypMod(),
					})
				}
			}
		} else if s.Query != nil {
			if err := c.validatePreparedDMLStmt(s.Query, paramTypes, fixedParams); err != nil {
				return err
			}
			for i := range paramTypes {
				if (paramTypes[i] == 0 || paramTypes[i] == UNKNOWNOID) && !c.snapshot {
					return &Error{
						Code:    CodeIndeterminateDatatype,
						Message: fmt.Sprintf("could not determine data type of parameter $%d", i+1),
					}
				}
			}
		}
		if s.Name != "" {
			if c.preparedStmts == nil {
				c.preparedStmts = make(map[string]*preparedStmt)
			}
			c.preparedStmts[s.Name] = &preparedStmt{
				query:         s.Query,
				argTypes:      paramTypes,
				numParams:     numParams,
				hasResultCols: hasResultCols,
				resultCols:    resultCols,
			}
		}
		return nil
	case *nodes.ExecuteStmt:
		// EXECUTE changes no schema, but PostgreSQL still refuses one that
		// names no prepared statement or passes unfit arguments.
		if c.snapshot {
			return nil
		}
		_, err := c.checkExecuteStmt(s)
		return err
	case *nodes.DeallocateStmt:
		if s.IsAll || s.Name == "" {
			clear(c.preparedStmts)
			return nil
		}
		if c.preparedStmts[s.Name] == nil {
			return &Error{
				Code:    CodeInvalidSQLStatementName,
				Message: fmt.Sprintf("prepared statement %q does not exist", s.Name),
			}
		}
		delete(c.preparedStmts, s.Name)
		return nil

	default:
		return errInvalidParameterValue(fmt.Sprintf("unsupported utility statement: %T", stmt))
	}
}

func isTempSchemaName(name string) bool {
	return name == "pg_temp" || strings.HasPrefix(name, "pg_temp_")
}

func (c *Catalog) discardTempObjects() {
	var tempRels []*Relation
	for _, rel := range c.relationByOID {
		if rel == nil {
			continue
		}
		if rel.Persistence == 't' || (rel.Schema != nil && isTempSchemaName(rel.Schema.Name)) {
			tempRels = append(tempRels, rel)
		}
	}
	for _, rel := range tempRels {
		if c.relationByOID[rel.OID] != nil {
			c.cascadeDropRelation(rel)
		}
	}

	var tempSeqs []*Sequence
	for _, seq := range c.sequenceByOID {
		if seq == nil {
			continue
		}
		if seq.Persistence == 't' || (seq.Schema != nil && isTempSchemaName(seq.Schema.Name)) {
			tempSeqs = append(tempSeqs, seq)
		}
	}
	for _, seq := range tempSeqs {
		if c.sequenceByOID[seq.OID] == nil {
			continue
		}
		if len(c.findNormalDependents('s', seq.OID)) > 0 {
			c.dropDependents('s', seq.OID)
		}
		if c.sequenceByOID[seq.OID] != nil {
			c.removeSequence(seq.Schema, seq)
		}
	}

	for _, sc := range c.schemas {
		if sc == nil || !isTempSchemaName(sc.Name) {
			continue
		}
		c.removeFunctionsInSchema(sc)
		c.removeTypesInSchema(sc)
		for name, idx := range sc.Indexes {
			if idx != nil {
				c.removeIndex(sc, name, idx)
			}
		}
		for name, uc := range sc.Collations {
			if uc != nil {
				delete(sc.Collations, name)
				c.removeDepsOn('L', uc.OID)
			}
		}
	}
}

type maxParamVisitor struct {
	max int
}

func (v *maxParamVisitor) Visit(node nodes.Node) nodes.Visitor {
	if node == nil {
		return nil
	}
	if pr, ok := node.(*nodes.ParamRef); ok && pr.Number > v.max {
		v.max = pr.Number
	}
	return v
}

func maxParamRefInNode(node nodes.Node) int {
	if node == nil {
		return 0
	}
	var v maxParamVisitor
	nodes.Walk(&v, node)
	return v.max
}

func collectQueryParamTypes(q *Query, paramTypes []uint32, fixedParams ...bool) error {
	if q == nil || len(paramTypes) == 0 {
		return nil
	}
	isFixed := len(fixedParams) > 0 && fixedParams[0]
	var firstErr error
	setParamType := func(num int, typ uint32) {
		if firstErr != nil || num <= 0 || num > len(paramTypes) || typ == 0 || typ == UNKNOWNOID {
			return
		}
		if paramTypes[num-1] == 0 || paramTypes[num-1] == UNKNOWNOID {
			paramTypes[num-1] = typ
		} else if !isFixed && paramTypes[num-1] != typ {
			firstErr = &Error{
				Code:    CodeAmbiguousParameter,
				Message: fmt.Sprintf("could not determine data type of parameter $%d", num),
			}
		}
	}
	var scanExpr func(e AnalyzedExpr)
	scanExpr = func(e AnalyzedExpr) {
		if e == nil {
			return
		}
		WalkAnalyzedExpr(e, func(n AnalyzedExpr) bool {
			switch v := n.(type) {
			case *ParamExpr:
				setParamType(v.ParamNum, v.ParamType)
			case *BoolExprQ:
				for _, a := range v.Args {
					if pe, ok := a.(*ParamExpr); ok {
						setParamType(pe.ParamNum, BOOLOID)
					}
				}
			case *BooleanTestExpr:
				if pe, ok := v.Arg.(*ParamExpr); ok {
					setParamType(pe.ParamNum, BOOLOID)
				}
			case *SubLinkExpr:
				if pe, ok := v.TestExpr.(*ParamExpr); ok && v.SubQuery != nil {
					for _, te := range v.SubQuery.TargetList {
						if te != nil && !te.ResJunk && te.Expr != nil {
							setParamType(pe.ParamNum, te.Expr.exprType())
							break
						}
					}
				}
				scanExpr(v.TestExpr)
				if err := collectQueryParamTypes(v.SubQuery, paramTypes, isFixed); err != nil && firstErr == nil {
					firstErr = err
				}
				return false
			}
			return true
		})
	}
	scanBoolQual := func(e AnalyzedExpr) {
		if pe, ok := e.(*ParamExpr); ok {
			setParamType(pe.ParamNum, BOOLOID)
		}
		scanExpr(e)
	}
	scanLimitExpr := func(e AnalyzedExpr) {
		if pe, ok := e.(*ParamExpr); ok {
			setParamType(pe.ParamNum, INT8OID)
		}
		scanExpr(e)
	}
	var scanJoin func(jn JoinNode)
	scanJoin = func(jn JoinNode) {
		switch j := jn.(type) {
		case *JoinExprNode:
			scanJoin(j.Left)
			scanJoin(j.Right)
			scanBoolQual(j.Quals)
		}
	}
	for _, cte := range q.CTEList {
		if err := collectQueryParamTypes(cte.Query, paramTypes, isFixed); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	for _, rte := range q.RangeTable {
		if rte.Subquery != nil {
			if err := collectQueryParamTypes(rte.Subquery, paramTypes, isFixed); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		for _, fe := range rte.FuncExprs {
			scanExpr(fe)
		}
		if rte.Tablesample != nil {
			for _, a := range rte.Tablesample.Args {
				scanExpr(a)
			}
			scanExpr(rte.Tablesample.Repeatable)
		}
	}
	for _, te := range q.TargetList {
		scanExpr(te.Expr)
	}
	if q.JoinTree != nil {
		for _, fn := range q.JoinTree.FromList {
			scanJoin(fn)
		}
		scanBoolQual(q.JoinTree.Quals)
	}
	scanBoolQual(q.HavingQual)
	scanLimitExpr(q.LimitCount)
	scanLimitExpr(q.LimitOffset)
	for _, wc := range q.WindowClause {
		if wc != nil {
			scanExpr(wc.StartOffset)
			scanExpr(wc.EndOffset)
		}
	}
	if err := collectQueryParamTypes(q.LArg, paramTypes, isFixed); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := collectQueryParamTypes(q.RArg, paramTypes, isFixed); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

func (c *Catalog) validateAndFormatCollationDef(def *nodes.List) (string, error) {
	if def == nil || len(def.Items) == 0 {
		return "", errInvalidObjectDefinition("parameter \"locale\" must be specified")
	}
	seen := make(map[string]bool, len(def.Items))
	var fromList *nodes.List
	var hasLocale, hasLCCollate, hasLCCtype, hasRules bool
	deterministic := true
	var provider, localeVal string
	var parts []string

	for _, item := range def.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		opt := strings.ToLower(d.Defname)
		if seen[opt] {
			return "", &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
		}
		seen[opt] = true

		switch opt {
		case "from":
			l, ok := d.Arg.(*nodes.List)
			if !ok || len(l.Items) == 0 {
				return "", &Error{Code: CodeSyntaxError, Message: "collation FROM requires a collation name"}
			}
			fromList = l
		case "deterministic":
			b, ok := parseDefElemBoolStrict(d)
			if !ok {
				return "", &Error{Code: CodeSyntaxError, Message: "deterministic requires a Boolean value"}
			}
			deterministic = b
			if b {
				parts = append(parts, "DETERMINISTIC = true")
			} else {
				parts = append(parts, "DETERMINISTIC = false")
			}
		case "locale", "lc_collate", "lc_ctype", "provider", "rules", "version":
			val, ok := parseDefElemStringStrict(d)
			if !ok {
				return "", &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("%s requires a string value", opt)}
			}
			switch opt {
			case "locale":
				hasLocale = true
				localeVal = val
			case "lc_collate":
				hasLCCollate = true
			case "lc_ctype":
				hasLCCtype = true
			case "provider":
				provider = strings.ToLower(val)
				val = provider
			case "rules":
				hasRules = true
			}
			parts = append(parts, strings.ToUpper(opt)+" = "+singleQuote(val))
		default:
			return "", &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("collation attribute %q not recognized", d.Defname)}
		}
	}

	if fromList != nil {
		if len(seen) > 1 {
			return "", &Error{Code: CodeSyntaxError, Message: "conflicting or redundant options"}
		}
		srcName, srcUserOID, err := c.resolveCollationList(fromList)
		if err != nil {
			return "", err
		}
		if srcUserOID == 0 && srcName == "default" {
			return "", errInvalidObjectDefinition("collation \"default\" cannot be copied")
		}
		if srcUserOID != 0 {
			for _, sc := range c.schemas {
				for _, uc := range sc.Collations {
					if uc != nil && uc.OID == srcUserOID {
						return uc.DefClause, nil
					}
				}
			}
		}
		qParts := make([]string, 0, len(fromList.Items))
		for _, item := range fromList.Items {
			qParts = append(qParts, quoteIdentifier(stringVal(item)))
		}
		return "FROM " + strings.Join(qParts, "."), nil
	}

	switch provider {
	case "", "libc", "icu", "builtin":
	default:
		return "", errInvalidObjectDefinition(fmt.Sprintf("unrecognized collation provider: %s", provider))
	}
	if !deterministic && provider != "icu" {
		return "", &Error{
			Code:    CodeFeatureNotSupported,
			Message: "nondeterministic collations are not supported with this provider",
		}
	}
	if hasRules && provider != "icu" {
		return "", errInvalidObjectDefinition("ICU rules cannot be specified unless locale provider is ICU")
	}
	if hasLocale && (hasLCCollate || hasLCCtype) {
		return "", errInvalidObjectDefinition("conflicting or redundant options")
	}
	if !hasLocale && !(hasLCCollate && hasLCCtype) {
		if hasLCCollate || hasLCCtype {
			return "", errInvalidObjectDefinition("parameters \"lc_collate\" and \"lc_ctype\" must both be specified")
		}
		return "", errInvalidObjectDefinition("parameter \"locale\" must be specified")
	}
	if provider == "builtin" {
		if !hasLocale {
			return "", errInvalidObjectDefinition("parameter \"locale\" must be specified")
		}
		switch localeVal {
		case "C", "C.UTF-8", "PG_UNICODE_FAST":
		default:
			return "", errInvalidParameterValue(fmt.Sprintf("invalid locale name %q for builtin provider", localeVal))
		}
	}
	sort.Strings(parts)
	return "(" + strings.Join(parts, ", ") + ")", nil
}

// processVariableSet handles SET/RESET variable statements.
//
// pg: src/backend/utils/misc/guc.c — ExecSetVariableStmt
func (c *Catalog) processVariableSet(s *nodes.VariableSetStmt) error {
	if s.Kind == nodes.VAR_RESET_ALL {
		c.SetSearchPath([]string{"$user", "public"})
		c.defaultTablespace = ""
		c.defaultTableAccessMethod = ""
		clear(c.sessionSettings)
		return nil
	}
	switch strings.ToLower(s.Name) {
	case "search_path":
		return c.processSetSearchPath(s)
	case "default_tablespace":
		return c.processSetDefaultTablespace(s)
	case "default_table_access_method":
		return c.processSetDefaultTableAccessMethod(s)
	case "role":
		return c.processSetRole(s)
	}
	key := strings.ToLower(s.Name)
	if key != "" {
		if c.sessionSettings == nil {
			c.sessionSettings = make(map[string]string)
		}
		switch s.Kind {
		case nodes.VAR_SET_VALUE:
			if val := formatFuncSetArgs(s.Args); val != "" {
				c.sessionSettings[key] = val
			}
		case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
			delete(c.sessionSettings, key)
		}
	}
	return nil
}

// processSetSearchPath handles SET search_path = ...
//
// pg: src/backend/catalog/namespace.c — preprocessNamespacePath
func (c *Catalog) processSetSearchPath(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			return nil
		}
		schemas := make([]string, 0, len(s.Args.Items))
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			schemas = append(schemas, sv.Str)
		}
		c.SetSearchPath(schemas)
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		// RESET search_path / SET search_path TO DEFAULT
		// pg: GUC default for search_path is '"$user", public'
		c.SetSearchPath([]string{"$user", "public"})
	}
	return nil
}

// processSetDefaultTablespace handles SET default_tablespace = ...
func (c *Catalog) processSetDefaultTablespace(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			c.defaultTablespace = ""
			return nil
		}
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			if sv.Str != "" && c.checkTablespace(sv.Str) != nil {
				return errInvalidParameterValue(fmt.Sprintf("invalid value for parameter \"default_tablespace\": %q", sv.Str))
			}
			c.defaultTablespace = sv.Str
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		c.defaultTablespace = ""
	}
	return nil
}

// processSetDefaultTableAccessMethod handles SET default_table_access_method.
func (c *Catalog) processSetDefaultTableAccessMethod(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			c.defaultTableAccessMethod = ""
			return nil
		}
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			if sv.Str != "" && !c.snapshot {
				am := c.accessMethodByName[sv.Str]
				if am == nil || am.Type != 't' {
					return errInvalidParameterValue(fmt.Sprintf("invalid value for parameter \"default_table_access_method\": %q", sv.Str))
				}
			}
			c.defaultTableAccessMethod = sv.Str
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		c.defaultTableAccessMethod = ""
	}
	return nil
}

// processSetRole handles SET ROLE, RESET ROLE, SET ROLE NONE.
//
// pg: src/backend/commands/variable.c — check_role, assign_role
// pg: src/backend/utils/init/miscinit.c — SetCurrentRoleId
func (c *Catalog) processSetRole(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			return nil
		}
		// Extract the role name from the first argument.
		// SET ROLE rolename is parsed as SET role = 'rolename'.
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			roleName := sv.Str
			// pg: check_role — "none" maps to InvalidOid → reset
			if roleName == "none" {
				c.ResetRole()
			} else {
				if err := c.checkRole(roleName); err != nil {
					return err
				}
				c.SetRole(roleName)
			}
			return nil
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		// RESET ROLE / SET ROLE TO DEFAULT → revert to session user.
		// pg: SetCurrentRoleId(InvalidOid, ...) → falls back to SessionUserId
		c.ResetRole()
	}
	return nil
}
