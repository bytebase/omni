package catalog

import (
	"fmt"
	"slices"

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
			prep, exists := c.preparedStmts[execStmt.Name]
			if !exists || prep == nil {
				return &Error{
					Code:    CodeInvalidSQLStatementName,
					Message: fmt.Sprintf("prepared statement %q does not exist", execStmt.Name),
				}
			}
			paramTypes := slices.Clone(prep.argTypes)
			if len(prep.argTypes) > 0 {
				gotParams := 0
				if execStmt.Params != nil {
					gotParams = len(execStmt.Params.Items)
				}
				if gotParams != len(prep.argTypes) {
					return &Error{
						Code:    CodeSyntaxError,
						Message: fmt.Sprintf("wrong number of parameters for prepared statement %q", execStmt.Name),
					}
				}
			}
			if execStmt.Params != nil {
				for i, item := range execStmt.Params.Items {
					pExpr, perr := c.AnalyzeStandaloneExpr(item, nil)
					if perr != nil {
						return perr
					}
					if i < len(paramTypes) && paramTypes[i] != 0 {
						if pExpr != nil && pExpr.exprType() != UNKNOWNOID {
							if _, cerr := c.coerceToTargetType(pExpr, pExpr.exprType(), paramTypes[i], 'i'); cerr != nil {
								return cerr
							}
						}
					} else {
						for len(paramTypes) <= i {
							paramTypes = append(paramTypes, 0)
						}
						if pExpr != nil && pExpr.exprType() != UNKNOWNOID {
							paramTypes[i] = pExpr.exprType()
						}
					}
				}
			}
			analyzeOpts.ParamTypes = paramTypes
			if prepSel, isSel := prep.query.(*nodes.SelectStmt); isSel {
				selStmt = prepSel
				ok = true
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

	if err := c.validateRelOptionsForCols(into.Options, relkind, "", columnsNeedToastTable(columns)); err != nil {
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
