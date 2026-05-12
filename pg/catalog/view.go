package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// DefineView creates a new view in the catalog.
//
// pg: src/backend/commands/view.c — DefineView
func (c *Catalog) DefineView(stmt *nodes.ViewStmt) error {
	rv := stmt.View
	schemaName := ""
	viewName := ""
	if rv != nil {
		schemaName = rv.Schemaname
		viewName = rv.Relname
	}

	// Reject UNLOGGED views.
	// pg: src/backend/commands/view.c — DefineView (persistence check)
	if rv != nil && rv.Relpersistence == 'u' {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "views cannot be unlogged",
		}
	}

	// Convert column aliases.
	columnNames := stringListItems(stmt.Aliases)

	selStmt, _ := stmt.Query.(*nodes.SelectStmt)

	// pg: src/backend/commands/view.c — DefineView (SELECT INTO check)
	if selStmt != nil && selStmt.IntoClause != nil {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "views must not contain SELECT INTO",
		}
	}

	// pg: src/backend/commands/view.c — DefineView (hasModifyingCTE check)
	// Views must not contain data-modifying statements (INSERT/UPDATE/DELETE) in WITH clauses.
	if selStmt != nil && selStmt.WithClause != nil && selStmt.WithClause.Ctes != nil {
		for _, cteNode := range selStmt.WithClause.Ctes.Items {
			if cte, ok := cteNode.(*nodes.CommonTableExpr); ok && cte.Ctequery != nil {
				switch cte.Ctequery.(type) {
				case *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt:
					return &Error{
						Code:    CodeFeatureNotSupported,
						Message: "views must not contain data-modifying statements in WITH",
					}
				}
			}
		}
	}

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	// Analyze the query — single pipeline.
	// pg: src/backend/commands/view.c — DefineView → DefineVirtualRelation
	// PG loops over Query.targetList calling exprType/exprTypmod/exprCollation.
	var analyzedQuery *Query
	if selStmt != nil {
		analyzedQuery, err = c.analyzeSelectStmt(selStmt)
		if err != nil {
			return err
		}
	}

	// pg: src/backend/commands/view.c — DefineVirtualRelation
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

	// Apply explicit column names if provided.
	if len(columnNames) > 0 {
		if len(columnNames) > len(cols) {
			return &Error{
				Code:    CodeDatatypeMismatch,
				Message: fmt.Sprintf("CREATE VIEW specifies %d column names, but query produces %d columns", len(columnNames), len(cols)),
			}
		}
		for i, name := range columnNames {
			cols[i].Name = name
		}
	}

	// Check for existing relation with the same name.
	if existing, exists := schema.Relations[viewName]; exists {
		if !stmt.Replace {
			return errDuplicateTable(viewName)
		}
		// OR REPLACE: must be a view.
		if existing.RelKind != 'v' {
			return errWrongObjectType(viewName, "a view")
		}
		// checkViewColumns: validate column compatibility.
		if err := c.checkViewColumns(cols, existing); err != nil {
			return err
		}
		// Update existing view in place.
		return c.replaceView(schema, existing, cols, analyzedQuery)
	}

	// Check conflict with existing type name.
	if c.typeByName[typeKey{ns: schema.OID, name: viewName}] != nil {
		return errDuplicateTable(viewName)
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

	rel := &Relation{
		OID:           relOID,
		Name:          viewName,
		Schema:        schema,
		RelKind:       'v',
		Columns:       columns,
		colByName:     colByName,
		RowTypeOID:    rowTypeOID,
		ArrayOID:      arrayOID,
		AnalyzedQuery: analyzedQuery,
	}

	// Store CHECK OPTION.
	// pg: src/backend/commands/view.c — DefineView (checkOption)
	// pgparser may set WithCheckOption (1=LOCAL, 2=CASCADED) or store it in Options.
	switch stmt.WithCheckOption {
	case 1: // LOCAL
		rel.CheckOption = 'l'
	case 2: // CASCADED
		rel.CheckOption = 'c'
	}
	// Also check Options list for check_option (pgparser stores it there for WITH clause).
	if rel.CheckOption == 0 && stmt.Options != nil {
		for _, opt := range stmt.Options.Items {
			de, ok := opt.(*nodes.DefElem)
			if !ok || de.Defname != "check_option" {
				continue
			}
			val := defElemString(de)
			switch strings.ToLower(val) {
			case "local":
				rel.CheckOption = 'l'
			case "cascaded":
				rel.CheckOption = 'c'
			}
		}
	}

	// Views with CHECK OPTION must be auto-updatable.
	// pg: src/backend/commands/view.c — DefineView (line 510-520)
	// A simple view is auto-updatable if it has exactly one table in FROM,
	// no aggregates, no DISTINCT, no GROUP BY, no HAVING, no LIMIT, no OFFSET,
	// no set operations, no CTEs, no window functions.
	// For pgddl, we only do a basic check: reject CHECK OPTION on set-op views.
	if rel.CheckOption != 0 && analyzedQuery != nil && analyzedQuery.SetOp != SetOpNone {
		return &Error{Code: CodeFeatureNotSupported,
			Message: "WITH CHECK OPTION is not supported on recursive views"}
	}

	schema.Relations[viewName] = rel
	c.relationByOID[relOID] = rel

	// Register composite row type.
	rowType := &BuiltinType{
		OID:       rowTypeOID,
		TypeName:  viewName,
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
	c.typeByName[typeKey{ns: schema.OID, name: viewName}] = rowType

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  fmt.Sprintf("_%s", viewName),
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

	// Record dependencies on source tables/views.
	// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
	c.recordViewDepsQ(relOID, analyzedQuery)

	return nil
}

// checkViewColumns validates that new view columns are compatible with old view columns.
// PG rules (view.c:checkViewColumns):
// 1. New column count >= old column count (cannot drop columns).
// 2. For each existing column: name, type OID, typmod must match exactly.
// 3. Collation must match (PG also checks this).
//
// pg: src/backend/commands/view.c — checkViewColumns
func (c *Catalog) checkViewColumns(newCols []*ResultColumn, existing *Relation) error {
	if len(newCols) < len(existing.Columns) {
		return &Error{
			Code:    CodeInvalidObjectDefinition,
			Message: "cannot drop columns from view",
		}
	}
	for i, oldCol := range existing.Columns {
		nc := newCols[i]
		if nc.Name != oldCol.Name {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: fmt.Sprintf("cannot change name of view column %q to %q", oldCol.Name, nc.Name),
			}
		}
		if nc.TypeOID != oldCol.TypeOID {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: fmt.Sprintf("cannot change data type of view column %q", oldCol.Name),
			}
		}
		if nc.TypeMod != oldCol.TypeMod {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: fmt.Sprintf("cannot change data type of view column %q", oldCol.Name),
			}
		}
		// pg: src/backend/commands/view.c — checkViewColumns (collation check)
		if nc.Collation != oldCol.Collation {
			return &Error{
				Code: CodeInvalidObjectDefinition,
				Message: fmt.Sprintf("cannot change collation of view column %q from %q to %q",
					oldCol.Name, c.collationName(oldCol.Collation), c.collationName(nc.Collation)),
			}
		}
	}
	return nil
}

// replaceView updates an existing view with new columns and query.
//
// (pgddl helper — PG does this inline in DefineVirtualRelation)
func (c *Catalog) replaceView(schema *Schema, rel *Relation, newCols []*ResultColumn, analyzedQuery *Query) error {
	// Remove old view dependencies.
	c.removeDepsOf('r', rel.OID)

	// Add new columns if any.
	for i := len(rel.Columns); i < len(newCols); i++ {
		rc := newCols[i]
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
		rel.Columns = append(rel.Columns, col)
		rel.colByName[col.Name] = i
	}

	// Update the query.
	rel.AnalyzedQuery = analyzedQuery

	// Re-record dependencies on source tables/views.
	// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
	c.recordViewDepsQ(rel.OID, analyzedQuery)

	return nil
}

// dropViewDependents drops all views that depend on the given relation.
func (c *Catalog) dropViewDependents(refType byte, refOID uint32) {
	deps := c.findNormalDependents(refType, refOID)
	for _, dep := range deps {
		switch dep.ObjType {
		case 'r':
			// Dependent view.
			rel := c.relationByOID[dep.ObjOID]
			if rel == nil || rel.RelKind != 'v' {
				continue
			}
			// Recursively drop dependents of this view first.
			c.dropViewDependents('r', rel.OID)
			c.removeRelation(rel.Schema, rel.Name, rel)
		case 'c':
			// Dependent constraint (e.g., FK).
			con := c.constraints[dep.ObjOID]
			if con == nil {
				continue
			}
			if ownerRel := c.relationByOID[con.RelOID]; ownerRel != nil {
				c.removeConstraint(ownerRel.Schema, con)
			}
		}
	}
}

// recordViewDepsQ records dependencies from a view to the tables referenced in its query.
// Walks the analyzed Query tree (range table entries) instead of the legacy SelectStmt.
//
// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
func (c *Catalog) recordViewDepsQ(viewOID uint32, q *Query) {
	if q == nil {
		return
	}

	// Record deps from range table (tables, subqueries).
	for _, rte := range q.RangeTable {
		switch rte.Kind {
		case RTERelation:
			if rte.RelOID != 0 {
				c.recordDependency('r', viewOID, 0, 'r', rte.RelOID, 0, DepNormal)
			}
		case RTESubquery:
			c.recordViewDepsQ(viewOID, rte.Subquery)
		}
	}

	// Recurse into set-op branches.
	if q.SetOp != SetOpNone {
		c.recordViewDepsQ(viewOID, q.LArg)
		c.recordViewDepsQ(viewOID, q.RArg)
	}

	// Recurse into CTEs.
	for _, cte := range q.CTEList {
		c.recordViewDepsQ(viewOID, cte.Query)
	}

	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		c.recordDependencyOnExpr('r', viewOID, te.Expr, DepNormal)
		c.recordViewVarDepsQ(viewOID, q, te.Expr)
	}
	if q.JoinTree != nil && q.JoinTree.Quals != nil {
		c.recordDependencyOnExpr('r', viewOID, q.JoinTree.Quals, DepNormal)
		c.recordViewVarDepsQ(viewOID, q, q.JoinTree.Quals)
	}
}

func (c *Catalog) recordViewVarDepsQ(viewOID uint32, q *Query, expr AnalyzedExpr) {
	if q == nil || expr == nil {
		return
	}
	switch v := expr.(type) {
	case *VarExpr:
		c.recordViewVarDep(viewOID, q, v)
	case *FuncCallExpr:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *AggExpr:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *OpExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Left)
		c.recordViewVarDepsQ(viewOID, q, v.Right)
	case *DistinctExprQ:
		c.recordViewVarDepsQ(viewOID, q, v.Left)
		c.recordViewVarDepsQ(viewOID, q, v.Right)
	case *ScalarArrayOpExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Left)
		c.recordViewVarDepsQ(viewOID, q, v.Right)
	case *NullIfExprQ:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *CoalesceExprQ:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *MinMaxExprQ:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *BoolExprQ:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *NullTestExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
	case *BooleanTestExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
	case *CaseExprQ:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
		for _, w := range v.When {
			c.recordViewVarDepsQ(viewOID, q, w.Condition)
			c.recordViewVarDepsQ(viewOID, q, w.Result)
		}
		c.recordViewVarDepsQ(viewOID, q, v.Default)
	case *ArrayExprQ:
		for _, elem := range v.Elements {
			c.recordViewVarDepsQ(viewOID, q, elem)
		}
	case *RowExprQ:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
	case *RelabelExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
	case *CoerceViaIOExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
	case *CoerceToDomainExpr:
		c.recordViewVarDepsQ(viewOID, q, v.Arg)
	case *SubscriptingRefExpr:
		c.recordViewVarDepsQ(viewOID, q, v.ContainerExpr)
		for _, idx := range v.SubscriptExprs {
			c.recordViewVarDepsQ(viewOID, q, idx)
		}
		for _, idx := range v.LowerExprs {
			c.recordViewVarDepsQ(viewOID, q, idx)
		}
	case *SubLinkExpr:
		c.recordViewDepsQ(viewOID, v.SubQuery)
	case *WindowFuncExpr:
		for _, arg := range v.Args {
			c.recordViewVarDepsQ(viewOID, q, arg)
		}
		c.recordViewVarDepsQ(viewOID, q, v.AggFilter)
	}
}

func (c *Catalog) recordViewVarDep(viewOID uint32, q *Query, v *VarExpr) {
	if v == nil || v.LevelsUp != 0 || v.AttNum <= 0 || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colIdx := int(v.AttNum - 1)
	switch rte.Kind {
	case RTERelation:
		c.recordDependency('r', viewOID, 0, 'r', rte.RelOID, int32(v.AttNum), DepNormal)
	case RTESubquery:
		if rte.Subquery != nil && colIdx < len(rte.Subquery.TargetList) {
			c.recordViewVarDepsQ(viewOID, rte.Subquery, rte.Subquery.TargetList[colIdx].Expr)
		}
	case RTECTE:
		if rte.CTEIndex >= 0 && rte.CTEIndex < len(q.CTEList) {
			cte := q.CTEList[rte.CTEIndex]
			if cte.Query != nil && colIdx < len(cte.Query.TargetList) {
				c.recordViewVarDepsQ(viewOID, cte.Query, cte.Query.TargetList[colIdx].Expr)
			}
		}
	}
}
