package catalog

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ExecCreateTableAs creates a table or materialized view from a query
// (CREATE TABLE ... AS SELECT or CREATE MATERIALIZED VIEW ... AS SELECT).
//
// pg: src/backend/commands/createas.c — ExecCreateTableAs
func (c *Catalog) ExecCreateTableAs(stmt *nodes.CreateTableAsStmt) error {
	into := stmt.Into
	if into == nil || into.Rel == nil {
		return errInvalidParameterValue("CREATE MATERIALIZED VIEW requires a target name")
	}

	rv := into.Rel
	schemaName := rv.Schemaname
	mvName := rv.Relname

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Analyze the query to infer column types (same approach as DefineView).
	var analyzeOpts AnalyzeOptions
	var execPrep *preparedStmt
	selStmt, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		// A prepared statement is run only after the target name is checked.
		if _, exists := schema.Relations[mvName]; exists {
			if stmt.IfNotExists {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
				return nil
			}
			return errDuplicateTable(mvName)
		}
		if c.typeByName[typeKey{ns: schema.OID, name: mvName}] != nil {
			if stmt.IfNotExists {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
				return nil
			}
			return errDuplicateTable(mvName)
		}
		if _, held := relationNamespaceKind(schema, mvName); held {
			if stmt.IfNotExists {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
				return nil
			}
			return errDuplicateTable(mvName)
		}
		if execStmt, isExec := stmt.Query.(*nodes.ExecuteStmt); isExec {
			prep, err := c.checkExecuteStmt(execStmt)
			if err != nil {
				return err
			}
			execPrep = prep
			paramTypes := slices.Clone(prep.argTypes)
			analyzeOpts.ParamTypes = paramTypes
			if prepSel, isSel := prep.query.(*nodes.SelectStmt); isSel {
				selStmt = prepSel
				ok = true
			} else if !c.snapshot {
				// pg: src/backend/commands/prepare.c — ExecuteQuery refuses a
				// prepared INSERT/UPDATE/DELETE/MERGE as a CTAS source, even
				// with RETURNING.
				return &Error{Code: CodeWrongObjectType, Message: "prepared statement is not a SELECT"}
			}
		}
		if !ok {
			return errInvalidParameterValue("CREATE MATERIALIZED VIEW requires a SELECT query")
		}
	}

	// Analyze the query — single pipeline.
	// pg: src/backend/commands/createas.c — ExecCreateTableAs
	analyzedQuery, err := c.analyzeSelectStmtWithOptions(selStmt, analyzeOpts)
	if err != nil {
		return err
	}
	if execPrep != nil && execPrep.hasResultCols && analyzedQuery != nil && !c.snapshot {
		outEntries := nonJunkTargetEntries(analyzedQuery.TargetList)
		if len(outEntries) != len(execPrep.resultCols) {
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: "cached plan must not change result type",
			}
		}
		for i, te := range outEntries {
			want := execPrep.resultCols[i]
			if te.Expr.exprType() != want.TypeOID || te.Expr.exprTypMod() != want.TypeMod {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: "cached plan must not change result type",
				}
			}
		}
	}

	// The query is read before the target name is checked: CREATE TABLE a AS
	// SELECT zzz names a column that does not exist even when a does. Only
	// then is a name another relation holds a duplicate (42P07), whether that
	// relation is a table or an index.
	if _, exists := schema.Relations[mvName]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
			return nil
		}
		return errDuplicateTable(mvName)
	}
	if c.typeByName[typeKey{ns: schema.OID, name: mvName}] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
			return nil
		}
		return errDuplicateTable(mvName)
	}
	if _, ok := relationNamespaceKind(schema, mvName); ok {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", mvName))
			return nil
		}
		return errDuplicateTable(mvName)
	}

	relkind := byte('m')
	if nodes.ObjectType(stmt.Objtype) == nodes.OBJECT_TABLE {
		relkind = 'r'
	}
	if into.OnCommit != nodes.ONCOMMIT_NOOP && rv.Relpersistence != 't' {
		return &Error{
			Code:    CodeInvalidTableDefinition,
			Message: "ON COMMIT can only be used on temporary tables",
		}
	}

	effectiveTablespace := c.effectiveTablespace(into.TableSpaceName)
	if into.TableSpaceName == "" && rv.Relpersistence == 't' {
		effectiveTablespace = ""
	}
	if err := c.checkRelationTablespace(effectiveTablespace); err != nil {
		return err
	}
	ctasAccessMethod := into.AccessMethod
	if ctasAccessMethod == "" {
		ctasAccessMethod = c.defaultTableAccessMethod
	}
	if err := c.checkTableAccessMethod(ctasAccessMethod); err != nil {
		return err
	}

	// Extract column types from analyzed query target entries.
	var cols []*ResultColumn
	if analyzedQuery != nil {
		for _, te := range analyzedQuery.TargetList {
			if te.ResJunk {
				continue
			}
			cols = append(cols, &ResultColumn{
				Name:      te.ResName,
				TypeOID:   te.Expr.exprType(),
				TypeMod:   te.Expr.exprTypMod(),
				Collation: te.Expr.exprCollation(),
			})
		}
	}

	// Apply explicit column names from INTO clause if provided.
	if into.ColNames != nil {
		columnNames := stringListItems(into.ColNames)
		if len(columnNames) != len(cols) {
			// pg: src/backend/commands/createas.c — ExecCreateTableAs (column count mismatch)
			return &Error{
				Code:    CodeSyntaxError,
				Message: fmt.Sprintf("CREATE MATERIALIZED VIEW specifies %d column names, but query produces %d columns", len(columnNames), len(cols)),
			}
		}
		for i, name := range columnNames {
			cols[i].Name = name
		}
	}

	// Build Column objects.
	columns := make([]*Column, len(cols))
	colByName := make(map[string]int, len(cols))
	for i, rc := range cols {
		if _, dup := colByName[rc.Name]; dup {
			return errDuplicateColumn(rc.Name)
		}
		colByName[rc.Name] = i

		typ := c.typeByOID[rc.TypeOID]
		coll := rc.Collation
		if coll == 0 {
			coll = c.typeCollation(rc.TypeOID)
		}
		col := &Column{
			AttNum:    int16(i + 1),
			Name:      rc.Name,
			TypeOID:   rc.TypeOID,
			TypeMod:   rc.TypeMod,
			IsLocal:   true,
			Collation: coll,
		}
		if typ != nil {
			col.Len = typ.Len
			col.ByVal = typ.ByVal
			col.Align = typ.Align
			col.Storage = typ.Storage
		}
		columns[i] = col
	}

	if err := c.validateRelOptions(into.Options, relkind, ""); err != nil {
		return err
	}

	// Allocate OIDs.
	relOID := c.oidGen.Next()
	rowTypeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	// Apply column names to analyzed query target list.
	if analyzedQuery != nil {
		for i, te := range analyzedQuery.TargetList {
			if i < len(cols) {
				te.ResName = cols[i].Name
			}
		}
	}

	persistence := rv.Relpersistence
	if persistence == 0 {
		persistence = 'p'
	}
	if ctasAccessMethod != "" {
		if am := c.accessMethodByName[ctasAccessMethod]; am != nil {
			c.recordDependency('r', relOID, 0, 'm', am.OID, 0, DepNormal)
		}
	}
	rel := &Relation{
		OID:           relOID,
		Name:          mvName,
		Schema:        schema,
		RelKind:       relkind,
		Persistence:   persistence,
		Tablespace:    normalizeTablespace(effectiveTablespace),
		AccessMethod:  ctasAccessMethod,
		OwnerName:     c.currentOwner(),
		Columns:       columns,
		colByName:     colByName,
		RowTypeOID:    rowTypeOID,
		ArrayOID:      arrayOID,
		AnalyzedQuery: analyzedQuery,
		IsPopulated:   relkind == 'm' && !into.SkipData,
	}

	schema.Relations[mvName] = rel
	c.relationByOID[relOID] = rel

	// Register composite row type.
	rowType := &BuiltinType{
		OID:       rowTypeOID,
		TypeName:  mvName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'c',
		Category:  'C',
		IsDefined: true,
		Delim:     ',',
		RelID:     relOID,
		Array:     arrayOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[rowTypeOID] = rowType
	c.typeByName[typeKey{ns: schema.OID, name: mvName}] = rowType

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  fmt.Sprintf("_%s", mvName),
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      rowTypeOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	// Record dependencies on source tables.
	// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
	c.recordViewDepsQ(relOID, analyzedQuery)

	return nil
}

// ExecRefreshMatView refreshes a materialized view.
// For pgddl, this is a no-op (no physical data to refresh) — just verifies the matview exists.
//
// pg: src/backend/commands/matview.c — ExecRefreshMatView
func (c *Catalog) ExecRefreshMatView(stmt *nodes.RefreshMatViewStmt) error {
	if stmt.Relation == nil {
		return errInvalidParameterValue("REFRESH MATERIALIZED VIEW requires a target name")
	}
	_, rel, err := c.findRelation(stmt.Relation.Schemaname, stmt.Relation.Relname)
	if err != nil {
		return err
	}
	if rel.RelKind == 'm' {
		if stmt.Concurrent {
			if stmt.SkipData {
				return &Error{
					Code:    CodeSyntaxError,
					Message: "CONCURRENTLY and WITH NO DATA options cannot be used together",
				}
			}
			if !rel.IsPopulated {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: "CONCURRENTLY cannot be used when the materialized view is not populated",
				}
			}
			hasUniqueIdx := false
			for _, idx := range c.indexesByRel[rel.OID] {
				if !idx.IsUnique || idx.WhereClause != "" || len(idx.Exprs) > 0 {
					continue
				}
				nKey := idx.NKeyColumns
				if nKey <= 0 || nKey > len(idx.Columns) {
					nKey = len(idx.Columns)
				}
				if nKey == 0 {
					continue
				}
				allCols := true
				for k := 0; k < nKey; k++ {
					if idx.Columns[k] <= 0 {
						allCols = false
						break
					}
				}
				if allCols {
					hasUniqueIdx = true
					break
				}
			}
			if !hasUniqueIdx {
				return &Error{
					Code:    CodeObjectNotInPrerequisiteState,
					Message: fmt.Sprintf("cannot refresh materialized view %q concurrently", rel.Name),
				}
			}
		}
		rel.IsPopulated = !stmt.SkipData
		return nil
	}
	if rel.RelKind == 'r' || rel.RelKind == 'p' {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: fmt.Sprintf("%q is not a materialized view", rel.Name),
		}
	}
	return errWrongObjectType(rel.Name, "a materialized view")
}

func (c *Catalog) validateLiteralForType(val string, targetType uint32) error {
	for dt := c.domainTypes[targetType]; dt != nil; dt = c.domainTypes[targetType] {
		targetType = dt.BaseTypeOID
	}
	trimmed := strings.TrimSpace(val)
	switch targetType {
	case INT2OID:
		v, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type smallint", val)}
			}
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type smallint: %q", val)}
		}
		if v < -32768 || v > 32767 {
			return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type smallint", val)}
		}
	case INT4OID, OIDOID, XIDOID, CIDOID:
		v, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type %s", val, c.typeName(targetType))}
			}
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type %s: %q", c.typeName(targetType), val)}
		}
		if targetType == INT4OID {
			if v < -2147483648 || v > 2147483647 {
				return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type integer", val)}
			}
		} else if v < 0 || v > 4294967295 {
			return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type %s", val, c.typeName(targetType))}
		}
	case INT8OID:
		if _, err := strconv.ParseInt(trimmed, 10, 64); err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type bigint", val)}
			}
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type bigint: %q", val)}
		}
	case FLOAT4OID, FLOAT8OID:
		switch strings.ToLower(trimmed) {
		case "nan", "infinity", "+infinity", "-infinity", "inf", "+inf", "-inf":
			return nil
		}
		bitSize := 64
		if targetType == FLOAT4OID {
			bitSize = 32
		}
		if _, err := strconv.ParseFloat(trimmed, bitSize); err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return &Error{Code: CodeNumericValueOutOfRange, Message: fmt.Sprintf("value %q is out of range for type %s", val, c.typeName(targetType))}
			}
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type %s: %q", c.typeName(targetType), val)}
		}
	case NUMERICOID:
		return validateNumericLiteral(val, trimmed)
	case BOOLOID:
		switch strings.ToLower(trimmed) {
		case "t", "tr", "tru", "true", "f", "fa", "fal", "fals", "false",
			"y", "ye", "yes", "n", "no", "on", "of", "off", "1", "0":
			return nil
		default:
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type boolean: %q", val)}
		}
	default:
		if et := c.enumTypes[targetType]; et != nil {
			if et.labelMap[val] == nil {
				return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input value for enum %s: %q", c.typeName(targetType), val)}
			}
		}
	}
	return nil
}

func validateNumericLiteral(raw, trimmed string) error {
	switch strings.ToLower(trimmed) {
	case "nan", "infinity", "+infinity", "-infinity", "inf", "+inf", "-inf":
		return nil
	}
	s := trimmed
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}
	if s == "" {
		return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
	}
	mantissa := s
	expStr := ""
	if idx := strings.IndexAny(s, "eE"); idx >= 0 {
		mantissa = s[:idx]
		expStr = s[idx+1:]
		if expStr == "" {
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
		}
	}
	digitsBefore := 0
	digitsAfter := 0
	seenDot := false
	for i := 0; i < len(mantissa); i++ {
		ch := mantissa[i]
		if ch == '.' {
			if seenDot {
				return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
			}
			seenDot = true
			continue
		}
		if ch < '0' || ch > '9' {
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
		}
		if seenDot {
			digitsAfter++
		} else {
			digitsBefore++
		}
	}
	if digitsBefore+digitsAfter == 0 {
		return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
	}
	if expStr != "" {
		exp, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return &Error{Code: CodeNumericValueOutOfRange, Message: "value overflows numeric format"}
			}
			return &Error{Code: CodeInvalidTextRepresentation, Message: fmt.Sprintf("invalid input syntax for type numeric: %q", raw)}
		}
		if exp > 131072 || exp < -16383 {
			return &Error{Code: CodeNumericValueOutOfRange, Message: "value overflows numeric format"}
		}
	}
	return nil
}

// checkExecuteStmt validates an EXECUTE against the session's prepared
// statements: the statement exists, and its arguments match the parameter
// count and coerce to the parameter types.
func (c *Catalog) checkExecuteStmt(execStmt *nodes.ExecuteStmt) (*preparedStmt, error) {
	prep, exists := c.preparedStmts[execStmt.Name]
	if !exists || prep == nil {
		return nil, &Error{
			Code:    CodeInvalidSQLStatementName,
			Message: fmt.Sprintf("prepared statement %q does not exist", execStmt.Name),
		}
	}
	// pg: src/backend/commands/prepare.c — ExecuteQuery wraps EvaluateParams
	// inside `if (entry->plansource->num_params > 0)`, so PostgreSQL 17 ignores
	// trailing EXECUTE arguments when a prepared statement has 0 parameters
	// (verified against PostgreSQL 17 in TestSimulationCausesOracle.tape:1959-1961).
	if prep.numParams == 0 {
		return prep, nil
	}
	gotParams := 0
	if execStmt.Params != nil {
		gotParams = len(execStmt.Params.Items)
	}
	if gotParams != prep.numParams {
		return nil, &Error{
			Code:    CodeSyntaxError,
			Message: fmt.Sprintf("wrong number of parameters for prepared statement %q", execStmt.Name),
		}
	}
	if execStmt.Params == nil {
		return prep, nil
	}
	for i, item := range execStmt.Params.Items {
		pExpr, perr := c.AnalyzeStandaloneExpr(item, nil)
		if perr != nil {
			return nil, perr
		}
		if i < len(prep.argTypes) && prep.argTypes[i] != 0 && pExpr != nil {
			if pExpr.exprType() == UNKNOWNOID {
				if ce, isConst := pExpr.(*ConstExpr); isConst && !ce.IsNull {
					if lerr := c.validateLiteralForType(ce.Value, prep.argTypes[i]); lerr != nil {
						return nil, lerr
					}
				}
			}
			if _, cerr := c.coerceToTargetType(pExpr, pExpr.exprType(), prep.argTypes[i], 'a'); cerr != nil {
				return nil, cerr
			}
		}
	}
	return prep, nil
}
