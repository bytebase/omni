// Package deparse — resolver.go implements schema-aware column qualification.
// The resolver takes a TableLookup + SelectStmt and returns a new SelectStmt where
// all column references are fully qualified with their table name/alias.
package deparse

import (
	"fmt"
	"sort"
	"strings"

	ast "github.com/bytebase/omni/tidb/ast"
	scopepkg "github.com/bytebase/omni/tidb/scope"
)

// ResolverColumn is a type alias for scope.Column, preserving backward compatibility.
type ResolverColumn = scopepkg.Column

// ResolverTable is a type alias for scope.Table, preserving backward compatibility.
type ResolverTable = scopepkg.Table

// TableLookup is a function that looks up a table by name in the catalog.
// Returns nil if the table is not found.
type TableLookup func(tableName string) *ResolverTable

// Resolver resolves column references in a SelectStmt using catalog metadata.
type Resolver struct {
	Lookup TableLookup
	// DefaultCharset is the database's default character set (e.g., "utf8mb4", "latin1").
	// Used to populate CAST(... AS CHAR) charset when not explicitly specified.
	// If empty, defaults to "utf8mb4".
	DefaultCharset string
	// Source is the text the statement was parsed from, which the AST's Loc
	// offsets index into. When it is set, an unaliased select item is named
	// after its source text, as MySQL names it; otherwise the name is rebuilt
	// from the AST, which matches MySQL only for text written the way the
	// rebuilt name spells it.
	Source string
	// Restore prepares the query for RestoreSelect the way TiDB's planner
	// prepares a view body before storing it (planner/core
	// logical_plan_builder.go, unfoldWildStar and addAliasName): every *
	// becomes the columns it stands for, qualified by schema and table;
	// every select field gets a name; a GROUP BY position becomes the field
	// it points at; and a table that is not a CTE gets its schema. Column
	// references, joins, and ORDER BY positions stay as written.
	Restore bool
	// Database is the schema Restore gives a table written without one.
	Database string
	// LookupIn, when set, looks up a table written with its schema
	// (db.t) in that schema; Lookup answers for every other name. Without
	// it a schema-qualified table is looked up by name alone.
	LookupIn func(schema, name string) *ResolverTable

	// visible holds the CTEs in scope where the resolver is working, so a
	// subquery sees the CTEs of the query around it.
	visible map[string]*ResolverTable
	// fieldSchema holds, for each select field that names a column, the
	// schema TiDB carries in its name (scope.Column.Schema).
	fieldSchema map[*ast.ResTarget]string
	// anonymous holds the select fields that were written with no alias
	// and are not a column reference: the ones TiDB renames when a view's
	// column names collide.
	anonymous map[*ast.ResTarget]bool
}

// child returns a resolver for a CTE body, with this one's settings and the
// given lookup.
func (r *Resolver) child(lookup TableLookup) *Resolver {
	return &Resolver{
		Lookup:         lookup,
		DefaultCharset: r.DefaultCharset,
		Source:         r.Source,
		Restore:        r.Restore,
		Database:       r.Database,
		LookupIn:       r.LookupIn,
		visible:        r.visible,
		anonymous:      r.anonymousSet(),
		fieldSchema:    r.fieldSchemas(),
	}
}

func (r *Resolver) fieldSchemas() map[*ast.ResTarget]string {
	if r.fieldSchema == nil {
		r.fieldSchema = make(map[*ast.ResTarget]string)
	}
	return r.fieldSchema
}

// annotate gives the columns of a derived table or CTE built from sel the
// schemas their select fields carry, in Restore mode, and the names TiDB
// gives them when no column list renames them: the field's name, empty for
// an empty string literal (SELECT * FROM (SELECT ”) d is stored as
// SELECT `d`.“ FROM ...).
func (r *Resolver) annotate(vt *ResolverTable, sel *ast.SelectStmt, named bool) *ResolverTable {
	if !r.Restore || vt == nil {
		return vt
	}
	leaf := ast.LeftmostQueryLeaf(sel)
	if leaf == nil {
		return vt
	}
	if leaf.TableSource != nil && leaf.TableSource.Table != nil {
		// TABLE t is SELECT * FROM t: each column carries t's schema.
		for i := range vt.Columns {
			vt.Columns[i].Schema = leaf.TableSource.Table.Schema
		}
		return vt
	}
	for i, target := range leaf.TargetList {
		if i >= len(vt.Columns) {
			break
		}
		if rt, ok := target.(*ast.ResTarget); ok {
			vt.Columns[i].Schema = r.fieldSchema[rt]
			if !named {
				vt.Columns[i].Name = rt.Name
			}
		}
	}
	return vt
}

// columnSchema returns the schema TiDB carries in the name of a column of the
// table at index in the scope (scope.Column.Schema).
func (r *Resolver) columnSchema(sc *fromScope, index int, column string) string {
	if index < 0 {
		return ""
	}
	if schema := sc.schemaOf(index); schema != "" {
		return schema
	}
	if t := sc.tables[index].Table; t != nil {
		if c := t.GetColumn(column); c != nil {
			return c.Schema
		}
	}
	return ""
}

func (r *Resolver) anonymousSet() map[*ast.ResTarget]bool {
	if r.anonymous == nil {
		r.anonymous = make(map[*ast.ResTarget]bool)
	}
	return r.anonymous
}

// Resolve takes a SelectStmt and returns a new SelectStmt with all column
// references fully qualified. The original AST is modified in-place.
func (r *Resolver) Resolve(stmt *ast.SelectStmt) *ast.SelectStmt {
	r.resolveWithCTEs(stmt, nil)
	if r.Restore {
		r.renameDuplicateFields(ast.LeftmostQueryLeaf(stmt))
	}
	return stmt
}

// resolveNested resolves a query nested in the one being resolved -- a
// subquery or a derived table -- with the CTEs around it in scope.
func (r *Resolver) resolveNested(stmt *ast.SelectStmt) {
	r.resolveWithCTEs(stmt, r.visible)
}

// renameDuplicateFields renames the unaliased fields of the view's first
// SELECT whose names collide with an earlier field's, as TiDB does for the
// SELECT that names a view's columns (addAliasName, AsViewSchema): the first
// repeat of a name n becomes Name_exp_n, the next Name_exp_1_n, and so on.
// Names compare without case, and an aliased field or column reference keeps
// its name and is counted first.
func (r *Resolver) renameDuplicateFields(sel *ast.SelectStmt) {
	if sel == nil {
		return
	}
	dedup := make(map[string]int)
	for _, target := range sel.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok && !r.anonymous[rt] {
			dedup[strings.ToLower(rt.Name)] = 0
		}
	}
	for _, target := range sel.TargetList {
		rt, ok := target.(*ast.ResTarget)
		if !ok || !r.anonymous[rt] {
			continue
		}
		old := strings.ToLower(rt.Name)
		if dup, ok := dedup[old]; ok {
			if dup == 0 {
				rt.Name = "Name_exp_" + rt.Name
			} else {
				rt.Name = fmt.Sprintf("Name_exp_%d_%s", dup, rt.Name)
			}
			dedup[old] = dup + 1
		} else {
			dedup[old] = 0
		}
	}
}

// resolveWithCTEs resolves a SelectStmt with optional CTE virtual tables
// available in scope. cteTables maps CTE names to virtual ResolverTables
// built from their SELECT target lists.
// resolveOrderByOrdinals rewrites positional ORDER BY items (ORDER BY 1) to the
// matching column alias from the leftmost SELECT's target list, matching MySQL
// 8.0 behavior.
func resolveOrderByOrdinals(orderBy []*ast.OrderByItem, leftmost *ast.SelectStmt) {
	if leftmost == nil {
		return
	}
	for _, item := range orderBy {
		lit, ok := item.Expr.(*ast.IntLit)
		if !ok {
			continue
		}
		idx := int(lit.Value) - 1 // 1-based to 0-based
		if idx < 0 || idx >= len(leftmost.TargetList) {
			continue
		}
		rt, ok := leftmost.TargetList[idx].(*ast.ResTarget)
		if !ok {
			continue
		}
		alias := rt.Name
		if alias == "" {
			// Derive alias from a bare column ref.
			if cr, ok := rt.Val.(*ast.ColumnRef); ok {
				alias = cr.Column
			}
		}
		if alias != "" {
			item.Expr = &ast.ColumnRef{Column: alias}
		}
	}
}

// resolveParenSource resolves a parenthesized query expression: the inner query
// (with the wrapper's CTEs visible), then any OUTER ORDER BY ordinals against
// the inner's leftmost SELECT target list.
func (r *Resolver) resolveParenSource(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable) {
	mergedCTETables := cteTables
	if len(stmt.CTEs) > 0 {
		mergedCTETables = make(map[string]*ResolverTable)
		for k, v := range cteTables {
			mergedCTETables[k] = v
		}
		for _, cte := range stmt.CTEs {
			cteResolver := r.child(r.withCTELookup(mergedCTETables))
			cteResolver.resolveWithCTEs(cte.Select, mergedCTETables)
			if vt := r.annotate(buildCTEVirtualTable(cte, cteResolver.Lookup), cte.Select, len(cte.Columns) > 0); vt != nil {
				mergedCTETables[strings.ToLower(cte.Name)] = vt
			}
		}
	}

	r.resolveWithCTEs(stmt.ParenSource, mergedCTETables)

	if len(stmt.OrderBy) == 0 {
		return
	}
	// Outer ORDER BY ordinals resolve against the inner's leftmost output SELECT.
	if !r.Restore {
		resolveOrderByOrdinals(stmt.OrderBy, ast.LeftmostQueryLeaf(stmt.ParenSource))
	}
}

func (r *Resolver) resolveWithCTEs(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable) *ast.SelectStmt {
	if stmt == nil {
		return nil
	}
	// Handle set operations recursively
	if stmt.SetOp != ast.SetOpNone {
		// The query-expression WITH lives on the set-op root node; register it
		// so it is visible to both sides of the set operation (matching MySQL
		// semantics where WITH ... applies to the entire UNION).
		mergedCTETables := cteTables
		if leftCTEs := collectSetOpCTEs(stmt); len(leftCTEs) > 0 {
			mergedCTETables = make(map[string]*ResolverTable)
			if cteTables != nil {
				for k, v := range cteTables {
					mergedCTETables[k] = v
				}
			}
			for _, cte := range leftCTEs {
				cteResolver := r.child(r.withCTELookup(mergedCTETables))
				cteResolver.resolveWithCTEs(cte.Select, mergedCTETables)
				vt := r.annotate(buildCTEVirtualTable(cte, cteResolver.Lookup), cte.Select, len(cte.Columns) > 0)
				if vt != nil {
					mergedCTETables[strings.ToLower(cte.Name)] = vt
				}
			}
		}
		if stmt.Left != nil {
			r.resolveWithCTEs(stmt.Left, mergedCTETables)
		}
		if stmt.Right != nil {
			r.resolveWithCTEs(stmt.Right, mergedCTETables)
		}
		// Resolve ORDER BY ordinals (e.g., ORDER BY 1) to column aliases
		// from the leftmost SELECT's target list, matching MySQL 8.0 behavior.
		if len(stmt.OrderBy) > 0 && !r.Restore {
			resolveOrderByOrdinals(stmt.OrderBy, ast.LeftmostQueryLeaf(stmt))
		}
		return stmt
	}

	// Parenthesized query expression: resolve the inner query (with the
	// wrapper's CTEs visible) and the outer ORDER BY ordinals.
	if stmt.ParenSource != nil {
		r.resolveParenSource(stmt, cteTables)
		return stmt
	}

	// Process CTEs: resolve each CTE's SELECT, then build virtual tables.
	// CTEs are resolved in order; later CTEs can reference earlier ones.
	localCTETables := make(map[string]*ResolverTable)
	if cteTables != nil {
		for k, v := range cteTables {
			localCTETables[k] = v
		}
	}
	for _, cte := range stmt.CTEs {
		if cte.Recursive && cte.Select != nil && ast.UnwrapParenSource(cte.Select).SetOp != ast.SetOpNone {
			// Recursive CTE: resolve left (non-recursive) branch first to get column info,
			// then add CTE to scope, then resolve right (recursive) branch. The body
			// may be a parenthesized set operation, so unwrap ParenSource first.
			sel := ast.UnwrapParenSource(cte.Select)

			// Step 1: Resolve the non-recursive (left) branch
			leftResolver := r.child(r.withCTELookup(localCTETables))
			leftResolver.resolveWithCTEs(sel.Left, localCTETables)

			// Step 2: Build virtual table from left branch's target list and add CTE to scope
			vt := r.annotate(buildCTEVirtualTableFromSelect(cte.Name, cte.Columns, sel.Left, leftResolver.Lookup), sel.Left, len(cte.Columns) > 0)
			if vt != nil {
				localCTETables[strings.ToLower(cte.Name)] = vt
			}

			// Step 3: Resolve the recursive (right) branch — CTE is now in scope
			rightResolver := r.child(r.withCTELookup(localCTETables))
			rightResolver.resolveWithCTEs(sel.Right, localCTETables)
		} else {
			// Non-recursive CTE: resolve entire SELECT, then build virtual table
			cteResolver := r.child(r.withCTELookup(localCTETables))
			cteResolver.resolveWithCTEs(cte.Select, localCTETables)

			vt := r.annotate(buildCTEVirtualTable(cte, cteResolver.Lookup), cte.Select, len(cte.Columns) > 0)
			if vt != nil {
				localCTETables[strings.ToLower(cte.Name)] = vt
			}
		}
	}

	// FROM DUAL is no table: TiDB's parser leaves the FROM clause empty and
	// writes FROM DUAL back only before a WHERE (ast.SelectStmt.Restore).
	if r.Restore && len(stmt.From) == 1 {
		if t, ok := stmt.From[0].(*ast.TableRef); ok && t.Schema == "" && t.Alias == "" && strings.EqualFold(t.Name, "dual") {
			stmt.From = nil
		}
	}

	// Subqueries of this block see its CTEs.
	origVisible := r.visible
	r.visible = localCTETables
	defer func() { r.visible = origVisible }()

	// Build scope from FROM clause, using a lookup that includes CTE virtual tables
	origLookup := r.Lookup
	if len(localCTETables) > 0 {
		r.Lookup = r.withCTELookup(localCTETables)
	}
	sc := r.buildScope(stmt.From)
	r.Lookup = origLookup

	// Resolve JOIN ON conditions (walk FROM clause) BEFORE target list resolution
	// so that NATURAL JOIN expansion can mark coalesced columns for star expansion.
	r.resolveFromExprs(stmt.From, sc)

	// Resolve target list (may expand stars)
	stmt.TargetList = r.resolveTargetList(stmt.TargetList, sc)

	// A TABLE query primary names a table TiDB qualifies like any other.
	if r.Restore && stmt.TableSource != nil && stmt.TableSource.Table != nil {
		if t := stmt.TableSource.Table; t.Schema == "" && !r.isCTE(t.Name) {
			t.Schema = r.Database
		}
	}

	// Resolve the rows of a VALUES query primary. A VALUES primary has no
	// target list and no FROM, so nothing above reaches its rows -- and a row
	// may hold a subquery, whose own SELECT needs resolving before deparse
	// writes it: `VALUES ROW((SELECT * FROM t))` otherwise came out as
	// `values row((select  from `t`))`, which is not SQL.
	if stmt.ValuesSource != nil {
		for _, row := range stmt.ValuesSource.Rows {
			for i, expr := range row {
				row[i] = r.resolveExpr(expr, sc)
			}
		}
	}

	// Resolve WHERE
	if stmt.Where != nil {
		stmt.Where = r.resolveExpr(stmt.Where, sc)
	}

	// Resolve GROUP BY. TiDB's planner replaces a position with the field
	// it points at, which the view then stores (GROUP BY 2 over b+1 is
	// GROUP BY `b`+1; verified on v8.5.5).
	for _, item := range stmt.GroupBy {
		if r.Restore {
			if field := groupByPosition(item.Expr, stmt.TargetList); field != nil {
				item.Expr = field
				continue
			}
		}
		item.Expr = r.resolveExpr(item.Expr, sc)
	}

	// Resolve HAVING
	if stmt.Having != nil {
		stmt.Having = r.resolveExpr(stmt.Having, sc)
		if r.Restore {
			stmt.Having = r.qualifyAggregateArgs(stmt.Having, stmt.TargetList, sc)
		}
	}

	// Resolve WINDOW clause
	for _, wd := range stmt.WindowClause {
		for i, expr := range wd.PartitionBy {
			wd.PartitionBy[i] = r.resolveExpr(expr, sc)
		}
		for _, item := range wd.OrderBy {
			item.Expr = r.resolveExpr(item.Expr, sc)
		}
	}

	// Resolve ORDER BY
	for _, item := range stmt.OrderBy {
		item.Expr = r.resolveExpr(item.Expr, sc)
		if r.Restore {
			item.Expr = r.qualifyAggregateArgs(item.Expr, stmt.TargetList, sc)
		}
	}

	return stmt
}

// buildScope constructs a scope from the FROM clause table expressions.
func (r *Resolver) buildScope(from []ast.TableExpr) *fromScope {
	sc := newFromScope(from)
	for _, tbl := range from {
		r.addTableExprToScope(tbl, sc)
	}
	return sc
}

// addTableExprToScope recursively adds table references from a table expression to the scope.
func (r *Resolver) addTableExprToScope(tbl ast.TableExpr, sc *fromScope) {
	defer sc.span(tbl, len(sc.tables))
	switch t := tbl.(type) {
	case *ast.TableRef:
		written := t.Schema
		if r.Restore && t.Schema == "" && !r.isCTE(t.Name) {
			t.Schema = r.Database
		}
		var table *ResolverTable
		if written != "" && r.LookupIn != nil {
			// A schema-qualified name is never a CTE, and names that
			// schema's table, not the session's.
			table = r.LookupIn(written, t.Name)
		} else {
			table = r.Lookup(t.Name)
		}
		if table == nil {
			return
		}
		effectiveName := t.Name
		if t.Alias != "" {
			effectiveName = t.Alias
		}
		sc.add(t, effectiveName, table)
		sc.setSchema(t.Schema)
	case *ast.JoinClause:
		r.addTableExprToScope(t.Left, sc)
		r.addTableExprToScope(t.Right, sc)
	case *ast.SubqueryExpr:
		// Derived table: resolve the inner SELECT, then build a virtual table
		// from its target list so outer queries can reference derived columns.
		if t.Select != nil {
			r.resolveNested(t.Select)
			vt := r.annotate(buildDerivedVirtualTable(t, r.Lookup), t.Select, len(t.Columns) > 0)
			if vt != nil && t.Alias != "" {
				sc.add(t, t.Alias, vt)
			}
		}
	}
}

// resolveTargetList resolves all target list entries, expanding qualified stars.
func (r *Resolver) resolveTargetList(targets []ast.ExprNode, sc *fromScope) []ast.ExprNode {
	var result []ast.ExprNode
	for i, target := range targets {
		expanded := r.resolveTarget(target, sc, i+1)
		result = append(result, expanded...)
	}
	return result
}

// resolveTarget resolves a single target list entry. Returns a slice because
// star expansion can produce multiple entries.
func (r *Resolver) resolveTarget(target ast.ExprNode, sc *fromScope, position int) []ast.ExprNode {
	rt, isRT := target.(*ast.ResTarget)

	var expr ast.ExprNode
	var explicitAlias string
	if isRT {
		expr = rt.Val
		explicitAlias = rt.Name
	} else {
		expr = target
	}

	// Check for qualified star: t1.*
	if col, ok := expr.(*ast.ColumnRef); ok && col.Star {
		return r.expandQualifiedStar(col.Table, sc)
	}

	// Check for unqualified star: *
	if _, ok := expr.(*ast.StarExpr); ok {
		return r.expandStar(sc)
	}

	// Apply CAST/CONVERT charset resolution before computing auto-alias.
	// This ensures the auto-alias includes the database charset (e.g., "charset latin1").
	// TiDB keeps the cast as written.
	if !r.Restore {
		r.resolveCastCharsets(expr)
	}

	// Compute auto-alias from the pre-resolution expression when no explicit alias.
	// MySQL 8.0 uses the original (unqualified) expression text for auto-aliases,
	// so we must derive it before column qualification changes the expression.
	anonymous := false
	if explicitAlias == "" {
		explicitAlias = r.autoName(expr, position)
		_, isCol := expr.(*ast.ColumnRef)
		anonymous = !isCol
	}

	// Resolve the expression (column qualification, etc.)
	resolved := r.resolveExpr(expr, sc)

	if !isRT {
		rt = &ast.ResTarget{}
	}
	rt.Val = resolved
	rt.Name = explicitAlias
	if anonymous && r.Restore {
		r.anonymousSet()[rt] = true
	}
	if col, ok := resolved.(*ast.ColumnRef); ok && r.Restore {
		schema := col.Schema
		if schema == "" {
			index, _ := sc.findColumnTable(col)
			schema = r.columnSchema(sc, index, col.Column)
		}
		r.fieldSchemas()[rt] = schema
	}
	return []ast.ExprNode{rt}
}

// autoName names an unaliased select item: from its source text when the
// resolver has it, otherwise rebuilt from the AST.
func (r *Resolver) autoName(expr ast.ExprNode, position int) string {
	if r.Source != "" {
		if name, ok := sourceAutoName(r.Source, expr, position, !r.Restore); ok {
			return name
		}
	}
	return autoAlias(expr, "", position)
}

// expandStar expands * to the columns of the FROM clause in MySQL's order: a
// NATURAL or USING join lists the columns it joins on once, first (see
// fromScope).
func (r *Resolver) expandStar(sc *fromScope) []ast.ExprNode {
	listed := make(map[colKey]bool)
	var cols []joinColumn
	for _, te := range sc.from {
		cols = sc.columns(te, false, listed, cols)
	}
	result := make([]ast.ExprNode, 0, len(cols))
	for _, c := range cols {
		rt := &ast.ResTarget{
			Name: c.name,
			Val:  r.starColumn(sc, c.table, c.qualifier, c.name),
		}
		if r.Restore {
			r.fieldSchemas()[rt] = r.columnSchema(sc, c.table, c.name)
		}
		result = append(result, rt)
	}
	return result
}

// expandQualifiedStar expands t1.* to all columns of table t1.
func (r *Resolver) expandQualifiedStar(tableName string, sc *fromScope) []ast.ExprNode {
	table := sc.GetTable(tableName)
	if table == nil {
		// Table not found in scope; return as-is
		return []ast.ExprNode{&ast.ResTarget{
			Val: &ast.ColumnRef{Table: tableName, Star: true},
		}}
	}

	// Find the effective name from scope entries (preserves alias casing)
	effectiveName := tableName
	if entry, ok := sc.FirstEqualFold(tableName); ok {
		effectiveName = entry.Name
	}

	index := sc.indexOf(tableName)
	var result []ast.ExprNode
	cols := sortedResolverColumns(table)
	for _, col := range cols {
		rt := &ast.ResTarget{
			Name: col.Name,
			Val:  r.starColumn(sc, index, effectiveName, col.Name),
		}
		if r.Restore {
			r.fieldSchemas()[rt] = r.columnSchema(sc, index, col.Name)
		}
		result = append(result, rt)
	}
	return result
}

// sortedResolverColumns returns columns sorted by Position.
func sortedResolverColumns(table *ResolverTable) []ResolverColumn {
	cols := make([]ResolverColumn, len(table.Columns))
	copy(cols, table.Columns)
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].Position < cols[j].Position
	})
	return cols
}

// resolveExpr resolves column references in an expression.
func (r *Resolver) resolveExpr(node ast.ExprNode, sc *fromScope) ast.ExprNode {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.ColumnRef:
		return r.resolveColumnRef(n, sc)
	case *ast.BinaryExpr:
		n.Left = r.resolveExpr(n.Left, sc)
		n.Right = r.resolveExpr(n.Right, sc)
		return n
	case *ast.UnaryExpr:
		n.Operand = r.resolveExpr(n.Operand, sc)
		return n
	case *ast.ParenExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.InExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		for i, item := range n.List {
			n.List[i] = r.resolveExpr(item, sc)
		}
		if n.Select != nil {
			r.resolveNested(n.Select)
		}
		return n
	case *ast.BetweenExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		n.Low = r.resolveExpr(n.Low, sc)
		n.High = r.resolveExpr(n.High, sc)
		return n
	case *ast.LikeExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		n.Pattern = r.resolveExpr(n.Pattern, sc)
		if n.Escape != nil {
			n.Escape = r.resolveExpr(n.Escape, sc)
		}
		return n
	case *ast.IsExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.CaseExpr:
		if n.Operand != nil {
			n.Operand = r.resolveExpr(n.Operand, sc)
		}
		for _, w := range n.Whens {
			w.Cond = r.resolveExpr(w.Cond, sc)
			w.Result = r.resolveExpr(w.Result, sc)
		}
		if n.Default != nil {
			n.Default = r.resolveExpr(n.Default, sc)
		}
		return n
	case *ast.FuncCallExpr:
		for i, arg := range n.Args {
			n.Args[i] = r.resolveExpr(arg, sc)
		}
		// Resolve ORDER BY in aggregate functions (e.g., GROUP_CONCAT)
		for _, item := range n.OrderBy {
			item.Expr = r.resolveExpr(item.Expr, sc)
		}
		// Resolve window function OVER clause
		if n.Over != nil {
			r.resolveWindowDef(n.Over, sc)
		}
		return n
	case *ast.CastExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		if r.Restore {
			r.fillCastLength(n.TypeName, n.Expr, sc)
		} else {
			r.resolveCastCharset(n.TypeName)
		}
		return n
	case *ast.ConvertExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		if r.Restore {
			r.fillCastLength(n.TypeName, n.Expr, sc)
		} else {
			r.resolveCastCharset(n.TypeName)
		}
		return n
	case *ast.PositionExpr:
		n.Substr = r.resolveExpr(n.Substr, sc)
		n.Str = r.resolveExpr(n.Str, sc)
		return n
	case *ast.ExtractExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.CollateExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.IntervalExpr:
		n.Value = r.resolveExpr(n.Value, sc)
		return n
	case *ast.RowExpr:
		for i, item := range n.Items {
			n.Items[i] = r.resolveExpr(item, sc)
		}
		return n
	case *ast.ExistsExpr:
		if n.Select != nil {
			r.resolveNested(n.Select)
		}
		return n
	case *ast.SubqueryExpr:
		if n.Select != nil {
			r.resolveNested(n.Select)
		}
		return n
	case *ast.ResTarget:
		n.Val = r.resolveExpr(n.Val, sc)
		return n
	default:
		// Leaf nodes (literals, etc.) — no resolution needed
		return node
	}
}

// starColumn returns the reference a * expands to for one column of the
// table at index in the scope. TiDB qualifies it by the table's schema too,
// when the table has one -- a CTE or derived table has none
// (unfoldWildStar).
func (r *Resolver) starColumn(sc *fromScope, index int, qualifier, name string) *ast.ColumnRef {
	col := &ast.ColumnRef{Table: qualifier, Column: name}
	if r.Restore {
		col.Schema = r.qualifierSchema(qualifier, r.columnSchema(sc, index, name))
	}
	return col
}

// qualifierSchema returns the schema TiDB writes before table.column: none
// when the table is named like a CTE in scope (ast.ColumnName.Restore,
// IsCTETableName).
func (r *Resolver) qualifierSchema(qualifier, schema string) string {
	if r.isCTE(qualifier) {
		return ""
	}
	return schema
}

// isCTE reports whether an unqualified table name denotes a CTE in scope.
func (r *Resolver) isCTE(name string) bool {
	_, ok := r.visible[strings.ToLower(name)]
	return ok
}

// groupByPosition returns the select field a GROUP BY position points at, or
// nil when expr is not a position in range.
func groupByPosition(expr ast.ExprNode, targets []ast.ExprNode) ast.ExprNode {
	lit, ok := expr.(*ast.IntLit)
	if !ok || lit.Value < 1 || lit.Value > int64(len(targets)) {
		return nil
	}
	if rt, ok := targets[lit.Value-1].(*ast.ResTarget); ok {
		return rt.Val
	}
	return nil
}

// resolveColumnRef qualifies an unqualified column reference by finding which
// table in scope contains the column. Restore leaves it as written.
func (r *Resolver) resolveColumnRef(col *ast.ColumnRef, sc *fromScope) ast.ExprNode {
	if r.Restore {
		return col
	}
	// Already qualified — just validate the table name maps to an alias
	if col.Table != "" {
		// Check if this table name is in scope (might be an alias)
		if sc.GetTable(col.Table) != nil {
			// Find the effective name from scope (preserves case)
			if entry, ok := sc.FirstEqualFold(col.Table); ok {
				col.Table = entry.Name
			}
		}
		return col
	}

	// A column a NATURAL or USING join joins on is the one copy the join
	// keeps, which comes from its right input in a RIGHT join.
	if c, ok := sc.coalesced[strings.ToLower(col.Column)]; ok {
		col.Table = c.qualifier
		return col
	}

	// Unqualified — qualify with the first table in scope holding the column.
	// When more than one does, MySQL would raise ERROR 1052: Column 'x' in
	// field list is ambiguous. TODO: return error
	entry, ok := sc.FirstWithColumn(col.Column)
	if !ok {
		// Column not found — return as-is (could be a literal alias or error)
		return col
	}
	col.Table = entry.Name
	return col
}

// resolveWindowDef resolves column references in a window definition.
func (r *Resolver) resolveWindowDef(wd *ast.WindowDef, sc *fromScope) {
	for i, expr := range wd.PartitionBy {
		wd.PartitionBy[i] = r.resolveExpr(expr, sc)
	}
	for _, item := range wd.OrderBy {
		item.Expr = r.resolveExpr(item.Expr, sc)
	}
}

// resolveFromExprs walks FROM clause table expressions and resolves
// ON condition expressions in JoinClauses.
func (r *Resolver) resolveFromExprs(from []ast.TableExpr, sc *fromScope) {
	for _, tbl := range from {
		r.resolveTableExpr(tbl, sc)
	}
}

// resolveTableExpr resolves expressions within a table expression (e.g., ON conditions).
// A NATURAL or USING join becomes the ON condition MySQL writes for it (see
// fromScope), comparing the columns it joins on.
func (r *Resolver) resolveTableExpr(tbl ast.TableExpr, sc *fromScope) {
	switch t := tbl.(type) {
	case *ast.JoinClause:
		r.resolveTableExpr(t.Left, sc)
		r.resolveTableExpr(t.Right, sc)

		if r.Restore {
			// TiDB keeps NATURAL and USING as written; the scope still
			// records the columns they join on, which * lists first.
			r.recordJoinColumns(t, sc)
			if on, ok := t.Condition.(*ast.OnCondition); ok {
				on.Expr = r.resolveExpr(on.Expr, sc)
			}
			return
		}

		// Expand NATURAL JOIN → find common columns → build ON condition
		if t.Type == ast.JoinNatural || t.Type == ast.JoinNaturalLeft || t.Type == ast.JoinNaturalRight {
			r.expandNaturalJoin(t, sc)
		}

		// Expand USING → build ON condition with qualified column refs
		if t.Condition != nil {
			if using, ok := t.Condition.(*ast.UsingCondition); ok {
				r.expandUsingCondition(t, using, sc)
			}
		}

		if t.Condition != nil {
			if on, ok := t.Condition.(*ast.OnCondition); ok {
				on.Expr = r.resolveExpr(on.Expr, sc)
			}
		}
	}
}

// recordJoinColumns records the columns a NATURAL or USING join joins on in
// the scope, without rewriting the join.
func (r *Resolver) recordJoinColumns(j *ast.JoinClause, sc *fromScope) {
	switch j.Type {
	case ast.JoinNatural, ast.JoinNaturalLeft, ast.JoinNaturalRight:
		if names, ok := sc.naturalColumns(j); ok {
			if len(names) > 0 {
				sc.joinOn(j, names)
			} else {
				sc.joinOnNothing(j)
			}
		}
		return
	}
	if using, ok := j.Condition.(*ast.UsingCondition); ok && len(using.Columns) > 0 {
		sc.joinOn(j, using.Columns)
	}
}

// expandNaturalJoin joins a NATURAL join on the columns its inputs share and
// builds its ON condition. It also changes the join type:
//   - NATURAL JOIN → JoinInner
//   - NATURAL LEFT JOIN → JoinLeft
//   - NATURAL RIGHT JOIN → JoinRight (deparse will then swap to LEFT)
//
// When an input's columns are unknown, the join keeps no condition.
func (r *Resolver) expandNaturalJoin(j *ast.JoinClause, sc *fromScope) {
	if names, ok := sc.naturalColumns(j); ok {
		if len(names) > 0 {
			if on, ok := sc.joinOn(j, names); ok {
				j.Condition = &ast.OnCondition{Expr: on}
			}
		} else {
			// Inputs that share no column still make a NATURAL join, which
			// lists a RIGHT join's right input first; an outer join needs a
			// condition, written on(true) as MySQL does (TiDB keeps NATURAL
			// as written).
			sc.joinOnNothing(j)
			if j.Type != ast.JoinNatural {
				j.Condition = &ast.OnCondition{Expr: noColumnsCondition()}
			}
		}
	}

	// Change join type from NATURAL variant to standard variant
	switch j.Type {
	case ast.JoinNatural:
		j.Type = ast.JoinInner
	case ast.JoinNaturalLeft:
		j.Type = ast.JoinLeft
	case ast.JoinNaturalRight:
		j.Type = ast.JoinRight
	}
}

// expandUsingCondition converts a USING condition into an ON condition with
// fully qualified column references, then replaces the condition on the join.
// When an input's columns are unknown, each side is qualified by the name of
// its input if that is a table.
func (r *Resolver) expandUsingCondition(j *ast.JoinClause, using *ast.UsingCondition, sc *fromScope) {
	if len(using.Columns) == 0 {
		return
	}
	if on, ok := sc.joinOn(j, using.Columns); ok {
		j.Condition = &ast.OnCondition{Expr: on}
		return
	}

	leftName := tableExprEffectiveName(j.Left)
	rightName := tableExprEffectiveName(j.Right)
	if leftName != "" && rightName != "" {
		j.Condition = &ast.OnCondition{
			Expr: buildColumnEqualityChain(using.Columns, leftName, rightName),
		}
	}
}

// tableExprEffectiveName returns the effective name (alias or table name) of a table expression.
func tableExprEffectiveName(tbl ast.TableExpr) string {
	switch t := tbl.(type) {
	case *ast.TableRef:
		if t.Alias != "" {
			return t.Alias
		}
		return t.Name
	default:
		return ""
	}
}

// buildColumnEqualityChain builds an AND-chained equality expression for column pairs.
// e.g., columns [a, b] with left=t1, right=t2 →
//
//	((`t1`.`a` = `t2`.`a`) and (`t1`.`b` = `t2`.`b`))
func buildColumnEqualityChain(columns []string, leftName, rightName string) ast.ExprNode {
	if len(columns) == 0 {
		return nil
	}

	// Build individual equality expressions
	equalities := make([]ast.ExprNode, len(columns))
	for i, col := range columns {
		equalities[i] = &ast.BinaryExpr{
			Op:    ast.BinOpEq,
			Left:  &ast.ColumnRef{Table: leftName, Column: col},
			Right: &ast.ColumnRef{Table: rightName, Column: col},
		}
	}

	// Chain with AND if multiple
	if len(equalities) == 1 {
		return equalities[0]
	}
	result := equalities[0]
	for i := 1; i < len(equalities); i++ {
		result = &ast.BinaryExpr{
			Op:    ast.BinOpAnd,
			Left:  result,
			Right: equalities[i],
		}
	}
	return result
}

// resolveCastCharsets walks the expression tree and sets charset on CAST/CONVERT DataTypes.
func (r *Resolver) resolveCastCharsets(node ast.ExprNode) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.CastExpr:
		r.resolveCastCharset(n.TypeName)
		r.resolveCastCharsets(n.Expr)
	case *ast.ConvertExpr:
		r.resolveCastCharset(n.TypeName)
		r.resolveCastCharsets(n.Expr)
	case *ast.BinaryExpr:
		r.resolveCastCharsets(n.Left)
		r.resolveCastCharsets(n.Right)
	case *ast.UnaryExpr:
		r.resolveCastCharsets(n.Operand)
	case *ast.ParenExpr:
		r.resolveCastCharsets(n.Expr)
	case *ast.FuncCallExpr:
		for _, arg := range n.Args {
			r.resolveCastCharsets(arg)
		}
	case *ast.CaseExpr:
		r.resolveCastCharsets(n.Operand)
		for _, w := range n.Whens {
			r.resolveCastCharsets(w.Cond)
			r.resolveCastCharsets(w.Result)
		}
		r.resolveCastCharsets(n.Default)
	}
}

// resolveCastCharset sets the charset on a CAST/CONVERT DataType for CHAR types.
// MySQL adds "charset <db_default>" when no charset is explicitly specified.
// The resolver uses DefaultCharset (from the catalog's database charset).
func (r *Resolver) resolveCastCharset(dt *ast.DataType) {
	if dt == nil {
		return
	}
	name := strings.ToLower(dt.Name)
	if name == "char" && dt.Charset == "" {
		charset := r.DefaultCharset
		if charset == "" {
			charset = "utf8mb4"
		}
		dt.Charset = charset
	}
}

// withCTELookup returns a TableLookup function that first checks CTE virtual tables,
// then falls back to the given fallback lookup.
// We capture fallback by value (not r.Lookup) to avoid infinite recursion when
// r.Lookup is later overwritten to point at the returned function itself.
func (r *Resolver) withCTELookup(cteTables map[string]*ResolverTable) TableLookup {
	fallback := r.Lookup // capture current value, not a reference to r.Lookup
	return func(tableName string) *ResolverTable {
		key := strings.ToLower(tableName)
		if vt, ok := cteTables[key]; ok {
			return vt
		}
		if fallback != nil {
			return fallback(tableName)
		}
		return nil
	}
}

// collectSetOpCTEs returns the WITH clause of a set-operation query expression.
// The parser attaches a query-expression-level WITH to the set-op ROOT node, so
// both operands can see it (matching MySQL, where WITH applies to the whole
// UNION). A parenthesized operand's own WITH lives inside its parens and is
// resolved when that operand is recursed into.
func collectSetOpCTEs(stmt *ast.SelectStmt) []*ast.CommonTableExpr {
	if len(stmt.CTEs) > 0 {
		return stmt.CTEs
	}
	return nil
}

// buildCTEVirtualTableFromSelect constructs a ResolverTable from a specific SELECT branch.
// Used for recursive CTEs where we need to build the virtual table from the non-recursive
// (left) branch before resolving the recursive (right) branch.
func buildCTEVirtualTableFromSelect(cteName string, cteColumns []string, sel *ast.SelectStmt, lookup TableLookup) *ResolverTable {
	if sel == nil {
		return nil
	}

	// If the CTE has an explicit column list, use those names
	if len(cteColumns) > 0 {
		cols := make([]ResolverColumn, len(cteColumns))
		for i, name := range cteColumns {
			cols[i] = ResolverColumn{Name: name, Position: i + 1}
		}
		return &ResolverTable{Name: cteName, Columns: cols}
	}

	cols := queryOutputColumns(sel, lookup)
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: cteName, Columns: cols}
}

// buildCTEVirtualTable constructs a ResolverTable from a CTE's SELECT target list.
// This allows the main query to resolve column references to CTE columns.
func buildCTEVirtualTable(cte *ast.CommonTableExpr, lookup TableLookup) *ResolverTable {
	if cte.Select == nil {
		return nil
	}

	// If the CTE has an explicit column list, use those names
	if len(cte.Columns) > 0 {
		cols := make([]ResolverColumn, len(cte.Columns))
		for i, name := range cte.Columns {
			cols[i] = ResolverColumn{Name: name, Position: i + 1}
		}
		return &ResolverTable{Name: cte.Name, Columns: cols}
	}

	// Otherwise, derive columns from the CTE's query expression
	cols := queryOutputColumns(cte.Select, lookup)
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: cte.Name, Columns: cols}
}

// PrimaryOutputColumns returns the columns a query expression outputs when its
// leftmost output query is a TABLE or VALUES primary, and nil when it is a
// SELECT, whose target list the caller reads with its own naming rules.
//
// It is the catalog's way in to the derivation the resolver already uses for a
// CTE or derived table over a primary, so a view body that is one gets the
// columns the server gives it. A database-only lookup is not enough on its
// own: the primary may name a CTE of its own query expression, and that name
// may shadow a real table -- `WITH t AS (SELECT 9 AS z) TABLE t` outputs z on
// both v8.5.5 and mysql:8.0, where resolving t to the stored table answered
// with the table's columns instead (container-verified).
func PrimaryOutputColumns(sel *ast.SelectStmt, lookup TableLookup) []ResolverColumn {
	leaf := ast.LeftmostQueryLeaf(sel)
	if leaf == nil || (leaf.TableSource == nil && leaf.ValuesSource == nil) {
		return nil
	}
	return queryOutputColumns(sel, lookup)
}

// queryOutputColumns returns the columns a query expression outputs, for the
// virtual table a CTE or derived table over it stands for. The names come from
// its leftmost output query, which is one of three shapes:
//
//   - a SELECT: its target list, named as cteColumnName reads it;
//   - a TABLE query primary: the primary is SELECT * FROM t, so the columns are
//     the table's, in table order, which lookup resolves -- unless the name is
//     a CTE of the query itself;
//   - a VALUES query primary: a row carries no names, so the columns are
//     ValuesColumnName's.
//
// A primary has no target list, so deriving from that alone left a CTE or
// derived table over one with no columns at all -- and SELECT * over it then
// expanded to nothing, which deparsed to an invalid "select  from `c`".
// Container-verified: a view over `WITH c AS (TABLE t) SELECT * FROM c` or
// `SELECT * FROM (TABLE t) d` has the columns of t.
func queryOutputColumns(sel *ast.SelectStmt, lookup TableLookup) []ResolverColumn {
	return queryOutputColumnsIn(sel, nil, lookup, nil)
}

// queryOutputColumnsIn is queryOutputColumns with the WITH scopes that enclose
// sel from outside it (innermost first) and the set of CTE declarations whose
// bodies are already being read. Both exist for one reason: a TABLE primary may name
// a CTE rather than a table, that CTE's body may be a primary naming another
// CTE, and a body may name itself.
func queryOutputColumnsIn(sel *ast.SelectStmt, outer []cteScope, lookup TableLookup, seen map[*ast.CommonTableExpr]bool) []ResolverColumn {
	scopes := append(cteScopesOf(sel), outer...)
	sel = ast.LeftmostQueryLeaf(sel)
	if sel == nil {
		return nil
	}
	switch {
	case sel.TableSource != nil:
		if sel.TableSource.Table == nil {
			return nil
		}
		// The primary may name a CTE in scope rather than a table. Only an
		// unqualified name can: a schema-qualified one always denotes a real
		// object, so `WITH t AS (SELECT 9 AS z) TABLE mydb.t` is the stored
		// table even where a CTE t shadows the unqualified name
		// (container-verified on both servers; the same rule the catalog's own
		// CTE-name check states).
		if sel.TableSource.Table.Schema == "" {
			if cols := cteOutputColumns(scopes, sel.TableSource.Table.Name, lookup, seen); cols != nil {
				return cols
			}
		}
		if lookup == nil {
			return nil
		}
		table := lookup(sel.TableSource.Table.Name)
		if table == nil {
			return nil
		}
		sorted := sortedResolverColumns(table)
		cols := make([]ResolverColumn, 0, len(sorted))
		for i, c := range sorted {
			cols = append(cols, ResolverColumn{Name: c.Name, Position: i + 1})
		}
		return cols
	case sel.ValuesSource != nil:
		if len(sel.ValuesSource.Rows) == 0 {
			return nil
		}
		cols := make([]ResolverColumn, 0, len(sel.ValuesSource.Rows[0]))
		for i := range sel.ValuesSource.Rows[0] {
			cols = append(cols, ResolverColumn{Name: ast.ValuesColumnName(i), Position: i + 1})
		}
		return cols
	}
	var cols []ResolverColumn
	for i, target := range sel.TargetList {
		cols = append(cols, ResolverColumn{Name: cteColumnName(target, i+1), Position: i + 1})
	}
	return cols
}

// cteScope is one WITH clause's CTEs.
type cteScope []*ast.CommonTableExpr

// cteScopesOf returns the WITH scopes that enclose root's leftmost output
// query, innermost first: each level's own CTEs, from root down through the
// parenthesized wrappers and set-operation left arms between root and that
// query.
//
// A set operation's RIGHT operand is not on that path: its CTEs are local to
// it, and a name the left operand uses resolves outside it. Searching them
// made `TABLE t UNION (WITH t AS (...) SELECT * FROM t)` read the left
// operand's t as the sibling's CTE, where both servers read the stored table
// (container-verified). A set operation's own WITH lives on the set-op node
// and is visible to both arms, so it is on the path.
func cteScopesOf(root *ast.SelectStmt) []cteScope {
	var scopes []cteScope
	for root != nil {
		if len(root.CTEs) > 0 {
			scopes = append(scopes, root.CTEs)
		}
		switch {
		case root.ParenSource != nil:
			root = root.ParenSource
		case root.SetOp != ast.SetOpNone:
			root = root.Left
		default:
			root = nil
		}
	}
	// Collected outermost first; the search wants innermost first.
	for i, j := 0, len(scopes)-1; i < j; i, j = i+1, j-1 {
		scopes[i], scopes[j] = scopes[j], scopes[i]
	}
	return scopes
}

// cteOutputColumns returns the columns of the CTE named name declared in
// scopes, or nil when none is. A CTE's explicit column list answers directly;
// otherwise its body is read as a query expression of its own, with the scopes
// outside it still in scope -- so `WITH x AS (SELECT 1 AS a), y AS (TABLE x)
// TABLE y` follows y into x and answers a, as both servers do.
//
// seen carries the DECLARATIONS already being read, not their names: a body
// may declare a CTE of the same name as the one it belongs to, and that inner
// declaration is a different one. Keying by name read
// `WITH x AS (WITH x AS (SELECT 1 AS a) TABLE x) TABLE x` as a cycle and
// derived nothing, where both servers give that view a. Keying by declaration
// still stops a body that names itself or a pair that name each other
// (`WITH RECURSIVE x AS (TABLE x)`, `WITH x AS (TABLE y), y AS (TABLE x)`,
// neither of which the servers will run), since those come back to the same
// declaration.
func cteOutputColumns(scopes []cteScope, name string, lookup TableLookup, seen map[*ast.CommonTableExpr]bool) []ResolverColumn {
	cte, rest := findCTE(scopes, name)
	if cte == nil {
		return nil
	}
	if len(cte.Columns) > 0 {
		cols := make([]ResolverColumn, len(cte.Columns))
		for i, n := range cte.Columns {
			cols[i] = ResolverColumn{Name: n, Position: i + 1}
		}
		return cols
	}
	if cte.Select == nil {
		return nil
	}
	if seen[cte] {
		return nil
	}
	next := make(map[*ast.CommonTableExpr]bool, len(seen)+1)
	for k := range seen {
		next[k] = true
	}
	next[cte] = true
	return queryOutputColumnsIn(cte.Select, rest, lookup, next)
}

// bodyScopes returns the WITH scopes in scope inside the body of the CTE found
// at scopes[i][j]: the siblings declared before it in its own WITH, then
// everything outside that WITH.
//
// A body sees an earlier sibling -- `WITH x AS (SELECT 1 AS a), y AS (TABLE x)
// TABLE y` outputs a on both servers -- and not a later one, which both refuse
// as a missing table (error 1146).
func bodyScopes(scopes []cteScope, i, j int) []cteScope {
	rest := make([]cteScope, 0, len(scopes)-i)
	if j > 0 {
		rest = append(rest, scopes[i][:j])
	}
	return append(rest, scopes[i+1:]...)
}

// findCTE finds the CTE named name in scopes and returns it with the scopes
// that remain in scope inside its body (see bodyScopes).
//
// The scopes are searched OUTERMOST first, which is what v8.5.5 does with a
// nested WITH that shadows an outer declaration of the same name:
// `WITH t AS (SELECT 1 AS z) (WITH t AS (SELECT 2 AS q) TABLE t)` answers z
// there, and its SHOW CREATE VIEW drops the inner WITH from the body
// altogether. mysql:8.0 answers q, the innermost, which is what the MySQL tree
// does; both container-verified, and this is the delta between them.
func findCTE(scopes []cteScope, name string) (*ast.CommonTableExpr, []cteScope) {
	for i := len(scopes) - 1; i >= 0; i-- {
		for j, cte := range scopes[i] {
			if cte != nil && strings.EqualFold(cte.Name, name) {
				return cte, bodyScopes(scopes, i, j)
			}
		}
	}
	return nil, nil
}

// cteColumnName extracts the column name from a target list entry.
// Uses alias if present, otherwise column ref name, otherwise positional name.
func cteColumnName(target ast.ExprNode, position int) string {
	if rt, ok := target.(*ast.ResTarget); ok {
		if rt.Name != "" {
			return rt.Name
		}
		if col, ok := rt.Val.(*ast.ColumnRef); ok {
			return col.Column
		}
		return fmt.Sprintf("Name_exp_%d", position)
	}
	if col, ok := target.(*ast.ColumnRef); ok {
		return col.Column
	}
	return fmt.Sprintf("Name_exp_%d", position)
}

// buildDerivedVirtualTable constructs a ResolverTable from a derived table's SELECT target list.
// This allows the outer query to resolve column references to derived table columns.
func buildDerivedVirtualTable(sub *ast.SubqueryExpr, lookup TableLookup) *ResolverTable {
	if sub.Select == nil || sub.Alias == "" {
		return nil
	}

	cols := queryOutputColumns(sub.Select, lookup)
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: sub.Alias, Columns: cols}
}

// AmbiguousColumnError is returned when a column reference matches multiple tables.
type AmbiguousColumnError struct {
	Column string
	Tables []string
}

func (e *AmbiguousColumnError) Error() string {
	return fmt.Sprintf("column %q is ambiguous, found in tables: %s", e.Column, strings.Join(e.Tables, ", "))
}

// ---------------------------------------------------------------------------
// Restore: what TiDB's planner changes in a view body
// ---------------------------------------------------------------------------

// aggregateFuncs are the aggregate functions TiDB's parser builds an
// AggregateFuncExpr for.
var aggregateFuncs = map[string]bool{
	"COUNT": true, "SUM": true, "AVG": true, "MIN": true, "MAX": true,
	"GROUP_CONCAT": true, "BIT_AND": true, "BIT_OR": true, "BIT_XOR": true,
	"STD": true, "STDDEV": true, "STDDEV_POP": true, "STDDEV_SAMP": true,
	"VARIANCE": true, "VAR_POP": true, "VAR_SAMP": true,
	"JSON_ARRAYAGG": true, "JSON_OBJECTAGG": true, "APPROX_COUNT_DISTINCT": true,
	"APPROX_PERCENTILE": true, "ANY_VALUE": true,
}

// qualifyAggregateArgs rewrites the column references inside the aggregate
// functions of a HAVING or ORDER BY expression as TiDB's planner does
// (logical_plan_builder.go, havingWindowAndOrderbyExprResolver): a column a
// plain column field of the select list names becomes that field as written,
// and any other becomes schema.table.column, with the schema the column's
// name carries (scope.Column.Schema) and none for a CTE. The view stores the
// rewritten expression: SELECT COUNT(*) FROM t GROUP BY a HAVING SUM(a) > 10
// is stored with HAVING SUM(`test`.`t`.`a`)>10 (verified on v8.5.5).
func (r *Resolver) qualifyAggregateArgs(node ast.ExprNode, targets []ast.ExprNode, sc *fromScope) ast.ExprNode {
	return rewriteExprs(node, func(n ast.ExprNode) (ast.ExprNode, bool) {
		fc, ok := n.(*ast.FuncCallExpr)
		if !ok || fc.Over != nil || !aggregateFuncs[strings.ToUpper(fc.Name)] {
			return n, false
		}
		for i, arg := range fc.Args {
			fc.Args[i] = rewriteExprs(arg, func(m ast.ExprNode) (ast.ExprNode, bool) {
				col, ok := m.(*ast.ColumnRef)
				if !ok || col.Star {
					return m, false
				}
				return r.aggregateColumn(col, targets, sc), true
			})
		}
		for _, item := range fc.OrderBy {
			item.Expr = rewriteExprs(item.Expr, func(m ast.ExprNode) (ast.ExprNode, bool) {
				col, ok := m.(*ast.ColumnRef)
				if !ok || col.Star {
					return m, false
				}
				return r.aggregateColumn(col, targets, sc), true
			})
		}
		return fc, true
	})
}

// aggregateColumn returns what a column reference inside a HAVING or ORDER BY
// aggregate becomes (see qualifyAggregateArgs). The schema is the one the
// column's name carries: a derived table's column over a stored column has
// that column's schema, so HAVING SUM(a) over (SELECT a FROM t) d is stored
// as SUM(`test`.`d`.`a`).
func (r *Resolver) aggregateColumn(col *ast.ColumnRef, targets []ast.ExprNode, sc *fromScope) ast.ExprNode {
	index, table := sc.findColumnTable(col)
	if table == nil {
		return col
	}
	c := table.GetColumn(col.Column)
	if c == nil {
		return col
	}
	entry := sc.tables[index]
	schema := r.columnSchema(sc, index, c.Name)
	// A plain column field naming the same column answers first.
	for _, target := range targets {
		rt, ok := target.(*ast.ResTarget)
		if !ok {
			continue
		}
		f, ok := rt.Val.(*ast.ColumnRef)
		if !ok || f.Star || !strings.EqualFold(f.Column, c.Name) {
			continue
		}
		if f.Table != "" && !strings.EqualFold(f.Table, entry.Name) {
			continue
		}
		if f.Schema != "" && !strings.EqualFold(f.Schema, schema) {
			continue
		}
		return &ast.ColumnRef{Schema: f.Schema, Table: f.Table, Column: f.Column}
	}
	return &ast.ColumnRef{Schema: r.qualifierSchema(entry.Name, schema), Table: entry.Name, Column: c.Name}
}

// findColumnTable returns the index and table of the scope entry a column
// reference names, or nil when none holds it.
func (sc *fromScope) findColumnTable(col *ast.ColumnRef) (int, *ResolverTable) {
	if col.Schema != "" && col.Table != "" {
		// db.t.c names the table of that schema: a.t and b.t may both be
		// in scope.
		for i, e := range sc.tables {
			if strings.EqualFold(e.Name, col.Table) && strings.EqualFold(sc.schemaOf(i), col.Schema) &&
				e.Table != nil && e.Table.GetColumn(col.Column) != nil {
				return i, e.Table
			}
		}
		return -1, nil
	}
	name := col.Table
	if name == "" {
		// A column a NATURAL or USING join joins on is the copy the join
		// keeps, the right input's in a RIGHT join.
		if c, ok := sc.coalesced[strings.ToLower(col.Column)]; ok {
			return c.table, sc.tables[c.table].Table
		}
		e, ok := sc.FirstWithColumn(col.Column)
		if !ok {
			return -1, nil
		}
		name = e.Name
	}
	i := sc.indexOf(name)
	if i < 0 || sc.tables[i].Table == nil || sc.tables[i].Table.GetColumn(col.Column) == nil {
		return -1, nil
	}
	return i, sc.tables[i].Table
}

// rewriteExprs walks an expression, outside its subqueries, and replaces
// each node f answers for; f's second result stops the walk below a node.
func rewriteExprs(node ast.ExprNode, f func(ast.ExprNode) (ast.ExprNode, bool)) ast.ExprNode {
	if node == nil {
		return nil
	}
	if out, done := f(node); done {
		return out
	}
	walk := func(e ast.ExprNode) ast.ExprNode { return rewriteExprs(e, f) }
	switch n := node.(type) {
	case *ast.BinaryExpr:
		n.Left, n.Right = walk(n.Left), walk(n.Right)
	case *ast.UnaryExpr:
		n.Operand = walk(n.Operand)
	case *ast.ParenExpr:
		n.Expr = walk(n.Expr)
	case *ast.InExpr:
		n.Expr = walk(n.Expr)
		for i, item := range n.List {
			n.List[i] = walk(item)
		}
	case *ast.BetweenExpr:
		n.Expr, n.Low, n.High = walk(n.Expr), walk(n.Low), walk(n.High)
	case *ast.LikeExpr:
		n.Expr, n.Pattern = walk(n.Expr), walk(n.Pattern)
	case *ast.IsExpr:
		n.Expr = walk(n.Expr)
	case *ast.CaseExpr:
		n.Operand = walk(n.Operand)
		for _, w := range n.Whens {
			w.Cond, w.Result = walk(w.Cond), walk(w.Result)
		}
		n.Default = walk(n.Default)
	case *ast.FuncCallExpr:
		for i, arg := range n.Args {
			n.Args[i] = walk(arg)
		}
	case *ast.CastExpr:
		n.Expr = walk(n.Expr)
	case *ast.ConvertExpr:
		n.Expr = walk(n.Expr)
	case *ast.CollateExpr:
		n.Expr = walk(n.Expr)
	case *ast.RowExpr:
		for i, item := range n.Items {
			n.Items[i] = walk(item)
		}
	}
	return node
}

// fillCastLength gives CAST(expr AS CHAR) and CAST(expr AS BINARY), written
// with no length, the length TiDB's planner gives them, which the view
// stores: the display width of expr's type. CAST(a AS CHAR) over an INT is
// stored as CAST(`a` AS CHAR(11)) (verified on v8.5.5). The width is known
// here only for the forms castLength reads; for any other the cast is kept as
// written.
func (r *Resolver) fillCastLength(dt *ast.DataType, expr ast.ExprNode, sc *fromScope) {
	if dt == nil || dt.Length > 0 {
		return
	}
	switch strings.ToUpper(dt.Name) {
	case "CHAR", "BINARY":
	default:
		return
	}
	if n, _, ok := castLength(expr, sc); ok {
		dt.Length = n
	}
}

// castLength returns the display width TiDB gives an expression's type and
// whether that type is an integer, for the forms whose width was measured on
// v8.5.5: a column of a stored table, an integer or string literal, integer
// arithmetic (BIGINT, 20), CONCAT (the sum of its arguments), and NOW() (19).
func castLength(expr ast.ExprNode, sc *fromScope) (int, bool, bool) {
	switch n := expr.(type) {
	case *ast.ParenExpr:
		return castLength(n.Expr, sc)
	case *ast.ColumnRef:
		_, table := sc.findColumnTable(n)
		if table == nil {
			return 0, false, false
		}
		c := table.GetColumn(n.Column)
		if c == nil || c.CharLength == 0 {
			return 0, false, false
		}
		return c.CharLength, c.Integer, true
	case *ast.IntLit:
		return 20, true, true
	case *ast.UnaryExpr:
		if n.Op == ast.UnaryMinus {
			if _, isInt, ok := castLength(n.Operand, sc); ok && isInt {
				return 20, true, true
			}
		}
	case *ast.StringLit:
		return len([]rune(n.Value)), false, true
	case *ast.BinaryExpr:
		switch n.Op {
		case ast.BinOpAdd, ast.BinOpSub, ast.BinOpMul:
			_, li, lok := castLength(n.Left, sc)
			_, ri, rok := castLength(n.Right, sc)
			if lok && rok && li && ri {
				return 20, true, true
			}
		}
	case *ast.FuncCallExpr:
		switch strings.ToUpper(n.Name) {
		case "NOW", "CURRENT_TIMESTAMP", "LOCALTIME", "LOCALTIMESTAMP":
			if len(n.Args) == 0 {
				return 19, false, true
			}
		case "CONCAT":
			total := 0
			for _, arg := range n.Args {
				l, _, ok := castLength(arg, sc)
				if !ok {
					return 0, false, false
				}
				total += l
			}
			return total, false, true
		}
	}
	return 0, false, false
}
