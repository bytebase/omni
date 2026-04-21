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
}

// Resolve takes a SelectStmt and returns a new SelectStmt with all column
// references fully qualified. The original AST is modified in-place.
func (r *Resolver) Resolve(stmt *ast.SelectStmt) *ast.SelectStmt {
	return r.resolveWithCTEs(stmt, nil)
}

// resolveWithCTEs resolves a SelectStmt with optional CTE virtual tables
// available in scope. cteTables maps CTE names to virtual ResolverTables
// built from their SELECT target lists.
func (r *Resolver) resolveWithCTEs(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable) *ast.SelectStmt {
	if stmt == nil {
		return nil
	}
	// Handle set operations recursively
	if stmt.SetOp != ast.SetOpNone {
		// Before recursing, hoist CTEs from the leftmost leaf so they are
		// visible to both sides of the set operation (matching MySQL semantics
		// where WITH ... applies to the entire UNION).
		mergedCTETables := cteTables
		if leftCTEs := collectLeftmostCTEs(stmt); len(leftCTEs) > 0 {
			mergedCTETables = make(map[string]*ResolverTable)
			if cteTables != nil {
				for k, v := range cteTables {
					mergedCTETables[k] = v
				}
			}
			for _, cte := range leftCTEs {
				cteResolver := &Resolver{
					Lookup:         r.withCTELookup(mergedCTETables),
					DefaultCharset: r.DefaultCharset,
				}
				cteResolver.resolveWithCTEs(cte.Select, mergedCTETables)
				vt := buildCTEVirtualTable(cte)
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
		if len(stmt.OrderBy) > 0 {
			leftmost := stmt
			for leftmost.SetOp != ast.SetOpNone && leftmost.Left != nil {
				leftmost = leftmost.Left
			}
			for _, item := range stmt.OrderBy {
				if lit, ok := item.Expr.(*ast.IntLit); ok {
					idx := int(lit.Value) - 1 // 1-based to 0-based
					if idx >= 0 && idx < len(leftmost.TargetList) {
						if rt, ok := leftmost.TargetList[idx].(*ast.ResTarget); ok {
							alias := rt.Name
							if alias == "" {
								// Derive alias from column ref
								if cr, ok := rt.Val.(*ast.ColumnRef); ok {
									alias = cr.Column
								}
							}
							if alias != "" {
								item.Expr = &ast.ColumnRef{Column: alias}
							}
						}
					}
				}
			}
		}
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
		if cte.Recursive && cte.Select != nil && cte.Select.SetOp != ast.SetOpNone {
			// Recursive CTE: resolve left (non-recursive) branch first to get column info,
			// then add CTE to scope, then resolve right (recursive) branch.
			sel := cte.Select

			// Step 1: Resolve the non-recursive (left) branch
			leftResolver := &Resolver{
				Lookup:         r.withCTELookup(localCTETables),
				DefaultCharset: r.DefaultCharset,
			}
			leftResolver.resolveWithCTEs(sel.Left, localCTETables)

			// Step 2: Build virtual table from left branch's target list and add CTE to scope
			vt := buildCTEVirtualTableFromSelect(cte.Name, cte.Columns, sel.Left)
			if vt != nil {
				localCTETables[strings.ToLower(cte.Name)] = vt
			}

			// Step 3: Resolve the recursive (right) branch — CTE is now in scope
			rightResolver := &Resolver{
				Lookup:         r.withCTELookup(localCTETables),
				DefaultCharset: r.DefaultCharset,
			}
			rightResolver.resolveWithCTEs(sel.Right, localCTETables)
		} else {
			// Non-recursive CTE: resolve entire SELECT, then build virtual table
			cteResolver := &Resolver{
				Lookup:         r.withCTELookup(localCTETables),
				DefaultCharset: r.DefaultCharset,
			}
			cteResolver.resolveWithCTEs(cte.Select, localCTETables)

			vt := buildCTEVirtualTable(cte)
			if vt != nil {
				localCTETables[strings.ToLower(cte.Name)] = vt
			}
		}
	}

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

	// Resolve WHERE
	if stmt.Where != nil {
		stmt.Where = r.resolveExpr(stmt.Where, sc)
	}

	// Resolve GROUP BY
	for i, expr := range stmt.GroupBy {
		stmt.GroupBy[i] = r.resolveExpr(expr, sc)
	}

	// Resolve HAVING
	if stmt.Having != nil {
		stmt.Having = r.resolveExpr(stmt.Having, sc)
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
	}

	return stmt
}

// buildScope constructs a scope from the FROM clause table expressions.
func (r *Resolver) buildScope(from []ast.TableExpr) *scopepkg.Scope {
	sc := scopepkg.New()
	for _, tbl := range from {
		r.addTableExprToScope(tbl, sc)
	}
	return sc
}

// addTableExprToScope recursively adds table references from a table expression to the scope.
func (r *Resolver) addTableExprToScope(tbl ast.TableExpr, sc *scopepkg.Scope) {
	switch t := tbl.(type) {
	case *ast.TableRef:
		table := r.Lookup(t.Name)
		if table == nil {
			return
		}
		effectiveName := t.Name
		if t.Alias != "" {
			effectiveName = t.Alias
		}
		sc.Add(effectiveName, table)
	case *ast.JoinClause:
		r.addTableExprToScope(t.Left, sc)
		r.addTableExprToScope(t.Right, sc)
	case *ast.SubqueryExpr:
		// Derived table: resolve the inner SELECT, then build a virtual table
		// from its target list so outer queries can reference derived columns.
		if t.Select != nil {
			r.Resolve(t.Select)
			vt := buildDerivedVirtualTable(t)
			if vt != nil && t.Alias != "" {
				sc.Add(t.Alias, vt)
			}
		}
	}
}

// resolveTargetList resolves all target list entries, expanding qualified stars.
func (r *Resolver) resolveTargetList(targets []ast.ExprNode, sc *scopepkg.Scope) []ast.ExprNode {
	var result []ast.ExprNode
	for i, target := range targets {
		expanded := r.resolveTarget(target, sc, i+1)
		result = append(result, expanded...)
	}
	return result
}

// resolveTarget resolves a single target list entry. Returns a slice because
// star expansion can produce multiple entries.
func (r *Resolver) resolveTarget(target ast.ExprNode, sc *scopepkg.Scope, position int) []ast.ExprNode {
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
	r.resolveCastCharsets(expr)

	// Compute auto-alias from the pre-resolution expression when no explicit alias.
	// MySQL 8.0 uses the original (unqualified) expression text for auto-aliases,
	// so we must derive it before column qualification changes the expression.
	if explicitAlias == "" {
		exprStr := deparseExpr(expr)
		explicitAlias = autoAlias(expr, exprStr, position)
	}

	// Resolve the expression (column qualification, etc.)
	resolved := r.resolveExpr(expr, sc)

	if isRT {
		rt.Val = resolved
		rt.Name = explicitAlias
		return []ast.ExprNode{rt}
	}
	return []ast.ExprNode{&ast.ResTarget{Name: explicitAlias, Val: resolved}}
}

// expandStar expands * to all columns from all tables in scope order.
// Columns marked as coalesced (from NATURAL JOIN / USING) are excluded.
func (r *Resolver) expandStar(sc *scopepkg.Scope) []ast.ExprNode {
	var result []ast.ExprNode
	for _, entry := range sc.AllEntries() {
		cols := sortedResolverColumns(entry.Table)
		for _, col := range cols {
			// Skip columns coalesced by NATURAL JOIN or USING
			if sc.IsCoalesced(entry.Name, col.Name) {
				continue
			}
			result = append(result, &ast.ResTarget{
				Name: col.Name,
				Val: &ast.ColumnRef{
					Table:  entry.Name,
					Column: col.Name,
				},
			})
		}
	}
	return result
}

// expandQualifiedStar expands t1.* to all columns of table t1.
func (r *Resolver) expandQualifiedStar(tableName string, sc *scopepkg.Scope) []ast.ExprNode {
	table := sc.GetTable(tableName)
	if table == nil {
		// Table not found in scope; return as-is
		return []ast.ExprNode{&ast.ResTarget{
			Val: &ast.ColumnRef{Table: tableName, Star: true},
		}}
	}

	// Find the effective name from scope entries (preserves alias casing)
	effectiveName := tableName
	for _, entry := range sc.AllEntries() {
		if strings.EqualFold(entry.Name, tableName) {
			effectiveName = entry.Name
			break
		}
	}

	var result []ast.ExprNode
	cols := sortedResolverColumns(table)
	for _, col := range cols {
		result = append(result, &ast.ResTarget{
			Name: col.Name,
			Val: &ast.ColumnRef{
				Table:  effectiveName,
				Column: col.Name,
			},
		})
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
func (r *Resolver) resolveExpr(node ast.ExprNode, sc *scopepkg.Scope) ast.ExprNode {
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
			r.Resolve(n.Select)
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
		r.resolveCastCharset(n.TypeName)
		return n
	case *ast.ConvertExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		r.resolveCastCharset(n.TypeName)
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
			r.Resolve(n.Select)
		}
		return n
	case *ast.SubqueryExpr:
		if n.Select != nil {
			r.Resolve(n.Select)
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

// resolveColumnRef qualifies an unqualified column reference by finding which
// table in scope contains the column.
func (r *Resolver) resolveColumnRef(col *ast.ColumnRef, sc *scopepkg.Scope) ast.ExprNode {
	// Already qualified — just validate the table name maps to an alias
	if col.Table != "" {
		// Check if this table name is in scope (might be an alias)
		if sc.GetTable(col.Table) != nil {
			// Find the effective name from scope (preserves case)
			for _, entry := range sc.AllEntries() {
				if strings.EqualFold(entry.Name, col.Table) {
					col.Table = entry.Name
					break
				}
			}
		}
		return col
	}

	// Unqualified — search all tables in scope
	var matchTable string
	var matchCount int
	for _, entry := range sc.AllEntries() {
		if entry.Table.GetColumn(col.Column) != nil {
			if matchCount == 0 {
				matchTable = entry.Name
			}
			matchCount++
		}
	}

	if matchCount == 0 {
		// Column not found — return as-is (could be a literal alias or error)
		return col
	}
	if matchCount > 1 {
		// Ambiguous — for now, qualify with first match
		// MySQL would raise ERROR 1052: Column 'x' in field list is ambiguous
		// TODO: return error
	}

	col.Table = matchTable
	return col
}

// resolveWindowDef resolves column references in a window definition.
func (r *Resolver) resolveWindowDef(wd *ast.WindowDef, sc *scopepkg.Scope) {
	for i, expr := range wd.PartitionBy {
		wd.PartitionBy[i] = r.resolveExpr(expr, sc)
	}
	for _, item := range wd.OrderBy {
		item.Expr = r.resolveExpr(item.Expr, sc)
	}
}

// resolveFromExprs walks FROM clause table expressions and resolves
// ON condition expressions in JoinClauses.
func (r *Resolver) resolveFromExprs(from []ast.TableExpr, sc *scopepkg.Scope) {
	for _, tbl := range from {
		r.resolveTableExpr(tbl, sc)
	}
}

// resolveTableExpr resolves expressions within a table expression (e.g., ON conditions).
// For NATURAL JOINs, it expands the join by finding common columns between both tables
// and building an ON condition. For USING clauses, it resolves column references.
func (r *Resolver) resolveTableExpr(tbl ast.TableExpr, sc *scopepkg.Scope) {
	switch t := tbl.(type) {
	case *ast.JoinClause:
		r.resolveTableExpr(t.Left, sc)
		r.resolveTableExpr(t.Right, sc)

		// Expand NATURAL JOIN → find common columns → build ON condition
		if t.Type == ast.JoinNatural || t.Type == ast.JoinNaturalLeft || t.Type == ast.JoinNaturalRight {
			r.expandNaturalJoin(t, sc)
		}

		// Expand USING → build ON condition with qualified column refs
		if t.Condition != nil {
			if using, ok := t.Condition.(*ast.UsingCondition); ok {
				r.expandUsingCondition(t, using)
			}
		}

		if t.Condition != nil {
			if on, ok := t.Condition.(*ast.OnCondition); ok {
				on.Expr = r.resolveExpr(on.Expr, sc)
			}
		}
	}
}

// expandNaturalJoin finds common columns between the left and right tables of a
// NATURAL JOIN and builds an ON condition. It also changes the join type:
//   - NATURAL JOIN → JoinInner
//   - NATURAL LEFT JOIN → JoinLeft
//   - NATURAL RIGHT JOIN → JoinRight (deparse will then swap to LEFT)
func (r *Resolver) expandNaturalJoin(j *ast.JoinClause, sc *scopepkg.Scope) {
	leftTable := r.lookupTableExpr(j.Left)
	rightTable := r.lookupTableExpr(j.Right)
	if leftTable == nil || rightTable == nil {
		// Can't expand without schema info; leave as-is
		switch j.Type {
		case ast.JoinNatural:
			j.Type = ast.JoinInner
		case ast.JoinNaturalLeft:
			j.Type = ast.JoinLeft
		case ast.JoinNaturalRight:
			j.Type = ast.JoinRight
		}
		return
	}

	// Find common columns (columns with matching names, case-insensitive)
	// Use left table column order for deterministic output
	leftCols := sortedResolverColumns(leftTable)
	var commonCols []string
	for _, lc := range leftCols {
		if rightTable.GetColumn(lc.Name) != nil {
			commonCols = append(commonCols, lc.Name)
		}
	}

	// Get effective table names for qualified column refs
	leftName := tableExprEffectiveName(j.Left)
	rightName := tableExprEffectiveName(j.Right)

	// Build ON condition from common columns
	if len(commonCols) > 0 {
		j.Condition = &ast.OnCondition{
			Expr: buildColumnEqualityChain(commonCols, leftName, rightName),
		}
	}

	// Mark common columns from the right table as coalesced for star expansion.
	// In NATURAL JOIN, common columns appear only once (from the left table).
	if sc != nil && len(commonCols) > 0 {
		for _, col := range commonCols {
			sc.MarkCoalesced(rightName, col)
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
func (r *Resolver) expandUsingCondition(j *ast.JoinClause, using *ast.UsingCondition) {
	leftName := tableExprEffectiveName(j.Left)
	rightName := tableExprEffectiveName(j.Right)

	if len(using.Columns) > 0 && leftName != "" && rightName != "" {
		j.Condition = &ast.OnCondition{
			Expr: buildColumnEqualityChain(using.Columns, leftName, rightName),
		}
	}
}

// lookupTableExpr returns the ResolverTable for a table expression.
// Only works for simple TableRef nodes.
func (r *Resolver) lookupTableExpr(tbl ast.TableExpr) *ResolverTable {
	switch t := tbl.(type) {
	case *ast.TableRef:
		return r.Lookup(t.Name)
	default:
		return nil
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

// collectLeftmostCTEs walks down the left spine of a set operation tree and
// returns CTEs from the leftmost leaf. Unlike extractCTEs in the deparser,
// this does NOT clear the CTEs — they must remain in the AST for the deparser
// to emit the WITH clause later. Instead, it marks them as already-resolved
// by removing them from a copy so the resolver doesn't double-process.
func collectLeftmostCTEs(stmt *ast.SelectStmt) []*ast.CommonTableExpr {
	cur := stmt
	for cur.SetOp != ast.SetOpNone && cur.Left != nil {
		cur = cur.Left
	}
	if len(cur.CTEs) > 0 {
		return cur.CTEs
	}
	return nil
}

// buildCTEVirtualTableFromSelect constructs a ResolverTable from a specific SELECT branch.
// Used for recursive CTEs where we need to build the virtual table from the non-recursive
// (left) branch before resolving the recursive (right) branch.
func buildCTEVirtualTableFromSelect(cteName string, cteColumns []string, sel *ast.SelectStmt) *ResolverTable {
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

	// Walk down to the leftmost leaf for set operations
	for sel.SetOp != ast.SetOpNone && sel.Left != nil {
		sel = sel.Left
	}

	var cols []ResolverColumn
	for i, target := range sel.TargetList {
		name := cteColumnName(target, i+1)
		cols = append(cols, ResolverColumn{Name: name, Position: i + 1})
	}
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: cteName, Columns: cols}
}

// buildCTEVirtualTable constructs a ResolverTable from a CTE's SELECT target list.
// This allows the main query to resolve column references to CTE columns.
func buildCTEVirtualTable(cte *ast.CommonTableExpr) *ResolverTable {
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

	// Otherwise, derive columns from the CTE's SELECT target list
	sel := cte.Select
	// For set operations, use the left side's target list
	for sel.SetOp != ast.SetOpNone && sel.Left != nil {
		sel = sel.Left
	}

	var cols []ResolverColumn
	for i, target := range sel.TargetList {
		name := cteColumnName(target, i+1)
		cols = append(cols, ResolverColumn{Name: name, Position: i + 1})
	}
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: cte.Name, Columns: cols}
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
func buildDerivedVirtualTable(sub *ast.SubqueryExpr) *ResolverTable {
	if sub.Select == nil || sub.Alias == "" {
		return nil
	}

	sel := sub.Select
	// For set operations, use the left side's target list
	for sel.SetOp != ast.SetOpNone && sel.Left != nil {
		sel = sel.Left
	}

	var cols []ResolverColumn
	for i, target := range sel.TargetList {
		name := cteColumnName(target, i+1)
		cols = append(cols, ResolverColumn{Name: name, Position: i + 1})
	}
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
