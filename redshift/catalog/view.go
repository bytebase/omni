package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
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

	// A UNION ALL column whose inputs' collations conflict has none, which
	// a view column of a collatable type cannot take.
	// pg: src/backend/commands/view.c — DefineVirtualRelation
	if analyzedQuery != nil {
		for i, te := range nonJunkTargetEntries(analyzedQuery.TargetList) {
			if i < len(cols) && c.lacksCollation(&queryLevels{query: analyzedQuery}, te.Expr) {
				return &Error{
					Code:    CodeIndeterminateCollation,
					Message: fmt.Sprintf("could not determine which collation to use for view column %q", cols[i].Name),
				}
			}
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
	w := &viewDepWalker{c: c, viewOID: viewOID}
	w.query(q)
}

// viewDepWalker records the dependencies of one view's query. levels holds
// the query levels enclosing the part being walked, innermost last, to
// resolve Vars with LevelsUp > 0. The walk pushes a level on entering a query
// and pops it on leaving, never copying the stack: the analyzer nests each
// set operation of a left-deep UNION one level inside the next, and a copy
// per level made a 30,000-arm UNION ALL allocate 9 GB.
type viewDepWalker struct {
	c       *Catalog
	viewOID uint32
	levels  []*Query
}

// query records the dependencies of one query level of a view: the
// relations in its range table, and what every expression of the level
// references, including the columns its Vars read.
//
// pg: src/backend/catalog/dependency.c — find_expr_references_walker, which
// reaches every part of a Query through query_tree_walker
func (w *viewDepWalker) query(q *Query) {
	if q == nil {
		return
	}
	w.levels = append(w.levels, q)
	defer func() { w.levels = w.levels[:len(w.levels)-1] }()

	// Record deps from range table (tables, subqueries, functions in FROM).
	for _, rte := range q.RangeTable {
		switch rte.Kind {
		case RTERelation:
			if rte.RelOID != 0 {
				w.c.recordDependency('r', w.viewOID, 0, 'r', rte.RelOID, 0, DepNormal)
			}
		case RTESubquery:
			w.query(rte.Subquery)
		case RTEFunction:
			for _, fe := range rte.FuncExprs {
				w.expr(fe)
			}
		}
	}

	// Recurse into set-op branches.
	if q.SetOp != SetOpNone {
		w.query(q.LArg)
		w.query(q.RArg)
	}

	// Recurse into CTEs.
	for _, cte := range q.CTEList {
		w.query(cte.Query)
	}

	// Junk entries carry the ORDER BY, GROUP BY, DISTINCT ON and window
	// expressions that are not in the output.
	for _, te := range q.TargetList {
		w.expr(te.Expr)
	}
	if q.JoinTree != nil {
		for _, jn := range q.JoinTree.FromList {
			w.join(jn)
		}
		w.expr(q.JoinTree.Quals)
	}
	w.expr(q.HavingQual)
	w.expr(q.LimitCount)
	w.expr(q.LimitOffset)
}

// lacksCollation reports whether expr is of a collatable type with no
// collation: an output column of a UNION ALL whose inputs' collations
// conflict, a plain reference to one through a subquery, CTE, or join, or a
// scalar subquery returning one. A column reference with no collation gives
// an expression over it its type's default collation, but a scalar subquery
// passes its missing collation on, as in PostgreSQL, so an expression over
// it lacks one too unless a COLLATE settles it. levels holds the query expr
// belongs to, and through it the queries enclosing it.
//
// pg: src/backend/parser/parse_collate.c — assign_collations_walker (Var, SubLink)
func (c *Catalog) lacksCollation(levels *queryLevels, expr AnalyzedExpr) bool {
	return c.lacksCollationIn(levels, expr, collationMemo{})
}

// queryLevels is a query and, through outer, the queries enclosing it. A
// level shares the levels above it with every level built on them, so
// following a column down through nested subqueries adds one node per
// subquery; copying a slice of all the levels instead made a column of a
// deeply nested view cost time quadratic in the nesting.
type queryLevels struct {
	query *Query
	outer *queryLevels
}

// up returns the level n above l, or nil when there are not that many.
func (l *queryLevels) up(n int) *queryLevels {
	for ; l != nil && n > 0; n-- {
		l = l.outer
	}
	return l
}

// lacksCollationIn is lacksCollation, keeping in memo the answers of
// explicitCollation for expr and the expressions under it.
func (c *Catalog) lacksCollationIn(levels *queryLevels, expr AnalyzedExpr, memo collationMemo) bool {
	switch v := expr.(type) {
	case *SetOpColumnExpr:
		return v.Collation == 0 && c.typeCollation(v.TypeOID) != 0
	case *VarExpr:
		// An outer reference reads a column of an enclosing query.
		if v.LevelsUp < 0 || v.AttNum <= 0 {
			return false
		}
		owner := levels.up(v.LevelsUp)
		if owner == nil {
			return false
		}
		q := owner.query
		if q == nil || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
			return false
		}
		rte := q.RangeTable[v.RangeIdx]
		var sub *Query
		switch rte.Kind {
		case RTESubquery:
			sub = rte.Subquery
		case RTECTE:
			// The CTE may belong to an enclosing query; find it by name,
			// innermost first, as the analyzer resolved it.
			for l := owner; l != nil && sub == nil; l = l.outer {
				for _, cte := range l.query.CTEList {
					if cte.Name == rte.CTEName {
						sub = cte.Query
					}
				}
			}
		case RTEJoin:
			if int(v.AttNum) <= len(rte.JoinAliasVars) {
				return c.lacksCollationIn(owner, rte.JoinAliasVars[v.AttNum-1], memo)
			}
			return false
		}
		if sub == nil {
			return false
		}
		te := nonJunkTargetEntry(sub.TargetList, int(v.AttNum))
		if te == nil {
			return false
		}
		return c.lacksCollationIn(&queryLevels{query: sub, outer: owner}, te.Expr, memo)
	case *SubLinkExpr:
		if v.SubLinkType != SubLinkExprType || v.SubQuery == nil {
			return false
		}
		te := nonJunkTargetEntry(v.SubQuery.TargetList, 1)
		return te != nil && c.lacksCollationIn(&queryLevels{query: v.SubQuery, outer: levels}, te.Expr, memo)
	case *ConstExpr, *ParamExpr, *CollateExprQ:
		return false
	}
	if c.typeCollation(expr.exprType()) == 0 {
		return false
	}
	if _, ok := memo.explicit(expr); ok {
		return false
	}
	for _, child := range AnalyzedExprChildren(expr) {
		if _, ok := child.(*SubLinkExpr); !ok && !isComputedExpr(child) {
			continue
		}
		if c.lacksCollationIn(levels, child, memo) {
			return true
		}
	}
	return false
}

// isComputedExpr reports whether e is computed from other expressions, as
// opposed to a column reference, constant, or parameter.
func isComputedExpr(e AnalyzedExpr) bool {
	switch e.(type) {
	case *VarExpr, *ConstExpr, *ParamExpr, *SetOpColumnExpr:
		return false
	}
	return true
}

// join records the dependencies of the JOIN ON (or USING) conditions in a
// FROM item.
func (w *viewDepWalker) join(jn JoinNode) {
	j, ok := jn.(*JoinExprNode)
	if !ok {
		return
	}
	w.join(j.Left)
	w.join(j.Right)
	w.expr(j.Quals)
}

// expr records what one expression of a view references: the functions and
// types it uses, and the columns it reads.
func (w *viewDepWalker) expr(expr AnalyzedExpr) {
	if expr == nil {
		return
	}
	w.c.recordDependencyOnExpr('r', w.viewOID, expr, DepNormal)
	w.varDeps(expr)
}

// varDeps records a column dependency for each Var in expr, and the
// dependencies of the subqueries it contains.
func (w *viewDepWalker) varDeps(expr AnalyzedExpr) {
	WalkAnalyzedExpr(expr, func(e AnalyzedExpr) bool {
		switch v := e.(type) {
		case *VarExpr:
			w.varDep(v)
		case *SubLinkExpr:
			w.varDeps(v.TestExpr)
			w.query(v.SubQuery)
			return false
		}
		return true
	})
}

// varDep records the view's dependency on the table column a Var reads. A
// Var of a join's output resolves through the join's alias expression, which
// for a FULL JOIN USING column reads both inputs. A Var of a subquery or CTE
// needs nothing here: query walks every target entry of both. A whole-row
// Var (AttNum 0) adds no column dependency.
//
// pg: src/backend/catalog/dependency.c — find_expr_references_walker (Var,
// including the RTE_JOIN case through joinaliasvars)
func (w *viewDepWalker) varDep(v *VarExpr) {
	if v == nil || v.AttNum <= 0 || v.LevelsUp < 0 || v.LevelsUp >= len(w.levels) {
		return
	}
	q := w.levels[len(w.levels)-1-v.LevelsUp]
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	switch rte := q.RangeTable[v.RangeIdx]; rte.Kind {
	case RTERelation:
		if rte.RelOID != 0 {
			w.c.recordDependency('r', w.viewOID, 0, 'r', rte.RelOID, int32(v.AttNum), DepNormal)
		}
	case RTEJoin:
		if int(v.AttNum) > len(rte.JoinAliasVars) {
			return
		}
		WalkAnalyzedExpr(rte.JoinAliasVars[v.AttNum-1], func(e AnalyzedExpr) bool {
			if av, ok := e.(*VarExpr); ok {
				w.varDep(&VarExpr{RangeIdx: av.RangeIdx, AttNum: av.AttNum, LevelsUp: v.LevelsUp + av.LevelsUp})
			}
			return true
		})
	}
}
