// Package deparse — resolver.go implements schema-aware column qualification.
// The resolver takes a TableLookup + SelectStmt and returns a new SelectStmt where
// all column references are fully qualified with their table name/alias.
package deparse

import (
	"fmt"
	"sort"
	"strings"

	ast "github.com/bytebase/omni/mysql/ast"
)

// ResolverColumn represents a column in a table, used by the resolver.
type ResolverColumn struct {
	Name     string
	Position int
}

// ResolverTable represents a table in the catalog, used by the resolver.
type ResolverTable struct {
	Name    string
	Columns []ResolverColumn
}

// GetColumn returns a column by name (case-insensitive), or nil if not found.
func (t *ResolverTable) GetColumn(name string) *ResolverColumn {
	lower := strings.ToLower(name)
	for i := range t.Columns {
		if strings.ToLower(t.Columns[i].Name) == lower {
			return &t.Columns[i]
		}
	}
	return nil
}

// TableLookup is a function that looks up a table by name in the catalog.
// Returns nil if the table is not found.
type TableLookup func(tableName string) *ResolverTable

// Resolver resolves column references in a SelectStmt using catalog metadata.
type Resolver struct {
	Lookup TableLookup
}

// scope maps table alias/name → *ResolverTable for the current FROM clause.
type scope struct {
	// tables maps effective name (alias if present, else table name) → *ResolverTable.
	tables map[string]*ResolverTable
	// order preserves insertion order for deterministic star expansion.
	order []scopeEntry
}

type scopeEntry struct {
	name  string // effective name (alias or table name)
	table *ResolverTable
}

// Resolve takes a SelectStmt and returns a new SelectStmt with all column
// references fully qualified. The original AST is modified in-place.
func (r *Resolver) Resolve(stmt *ast.SelectStmt) *ast.SelectStmt {
	if stmt == nil {
		return nil
	}
	// Handle set operations recursively
	if stmt.SetOp != ast.SetOpNone {
		if stmt.Left != nil {
			r.Resolve(stmt.Left)
		}
		if stmt.Right != nil {
			r.Resolve(stmt.Right)
		}
		return stmt
	}

	// Build scope from FROM clause
	sc := r.buildScope(stmt.From)

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

	// Resolve ORDER BY
	for _, item := range stmt.OrderBy {
		item.Expr = r.resolveExpr(item.Expr, sc)
	}

	// Resolve JOIN ON conditions (walk FROM clause)
	r.resolveFromExprs(stmt.From, sc)

	return stmt
}

// buildScope constructs a scope from the FROM clause table expressions.
func (r *Resolver) buildScope(from []ast.TableExpr) *scope {
	sc := &scope{
		tables: make(map[string]*ResolverTable),
	}
	for _, tbl := range from {
		r.addTableExprToScope(tbl, sc)
	}
	return sc
}

// addTableExprToScope recursively adds table references from a table expression to the scope.
func (r *Resolver) addTableExprToScope(tbl ast.TableExpr, sc *scope) {
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
		key := strings.ToLower(effectiveName)
		sc.tables[key] = table
		sc.order = append(sc.order, scopeEntry{name: effectiveName, table: table})
	case *ast.JoinClause:
		r.addTableExprToScope(t.Left, sc)
		r.addTableExprToScope(t.Right, sc)
	}
}

// resolveTargetList resolves all target list entries, expanding qualified stars.
func (r *Resolver) resolveTargetList(targets []ast.ExprNode, sc *scope) []ast.ExprNode {
	var result []ast.ExprNode
	for i, target := range targets {
		expanded := r.resolveTarget(target, sc, i+1)
		result = append(result, expanded...)
	}
	return result
}

// resolveTarget resolves a single target list entry. Returns a slice because
// star expansion can produce multiple entries.
func (r *Resolver) resolveTarget(target ast.ExprNode, sc *scope, position int) []ast.ExprNode {
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

	// Compute auto-alias from the pre-resolution expression when no explicit alias.
	// MySQL 8.0 uses the original (unqualified) expression text for auto-aliases,
	// so we must derive it before column qualification changes the expression.
	if explicitAlias == "" {
		exprStr := deparseExpr(expr)
		explicitAlias = autoAlias(expr, exprStr, position)
	}

	// Resolve the expression
	resolved := r.resolveExpr(expr, sc)
	if isRT {
		rt.Val = resolved
		rt.Name = explicitAlias
		return []ast.ExprNode{rt}
	}
	return []ast.ExprNode{&ast.ResTarget{Name: explicitAlias, Val: resolved}}
}

// expandStar expands * to all columns from all tables in scope order.
func (r *Resolver) expandStar(sc *scope) []ast.ExprNode {
	var result []ast.ExprNode
	for _, entry := range sc.order {
		cols := sortedResolverColumns(entry.table)
		for _, col := range cols {
			result = append(result, &ast.ResTarget{
				Name: col.Name,
				Val: &ast.ColumnRef{
					Table:  entry.name,
					Column: col.Name,
				},
			})
		}
	}
	return result
}

// expandQualifiedStar expands t1.* to all columns of table t1.
func (r *Resolver) expandQualifiedStar(tableName string, sc *scope) []ast.ExprNode {
	key := strings.ToLower(tableName)
	table, ok := sc.tables[key]
	if !ok {
		// Table not found in scope; return as-is
		return []ast.ExprNode{&ast.ResTarget{
			Val: &ast.ColumnRef{Table: tableName, Star: true},
		}}
	}

	// Find the effective name from scope order (preserves alias casing)
	effectiveName := tableName
	for _, entry := range sc.order {
		if strings.EqualFold(entry.name, tableName) {
			effectiveName = entry.name
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
func (r *Resolver) resolveExpr(node ast.ExprNode, sc *scope) ast.ExprNode {
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
		return n
	case *ast.ConvertExpr:
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
func (r *Resolver) resolveColumnRef(col *ast.ColumnRef, sc *scope) ast.ExprNode {
	// Already qualified — just validate the table name maps to an alias
	if col.Table != "" {
		// Check if this table name is in scope (might be an alias)
		key := strings.ToLower(col.Table)
		if _, ok := sc.tables[key]; ok {
			// Find the effective name from scope (preserves case)
			for _, entry := range sc.order {
				if strings.EqualFold(entry.name, col.Table) {
					col.Table = entry.name
					break
				}
			}
		}
		return col
	}

	// Unqualified — search all tables in scope
	var matchTable string
	var matchCount int
	for _, entry := range sc.order {
		if entry.table.GetColumn(col.Column) != nil {
			if matchCount == 0 {
				matchTable = entry.name
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
func (r *Resolver) resolveWindowDef(wd *ast.WindowDef, sc *scope) {
	for i, expr := range wd.PartitionBy {
		wd.PartitionBy[i] = r.resolveExpr(expr, sc)
	}
	for _, item := range wd.OrderBy {
		item.Expr = r.resolveExpr(item.Expr, sc)
	}
}

// resolveFromExprs walks FROM clause table expressions and resolves
// ON condition expressions in JoinClauses.
func (r *Resolver) resolveFromExprs(from []ast.TableExpr, sc *scope) {
	for _, tbl := range from {
		r.resolveTableExpr(tbl, sc)
	}
}

// resolveTableExpr resolves expressions within a table expression (e.g., ON conditions).
// For NATURAL JOINs, it expands the join by finding common columns between both tables
// and building an ON condition. For USING clauses, it resolves column references.
func (r *Resolver) resolveTableExpr(tbl ast.TableExpr, sc *scope) {
	switch t := tbl.(type) {
	case *ast.JoinClause:
		r.resolveTableExpr(t.Left, sc)
		r.resolveTableExpr(t.Right, sc)

		// Expand NATURAL JOIN → find common columns → build ON condition
		if t.Type == ast.JoinNatural || t.Type == ast.JoinNaturalLeft || t.Type == ast.JoinNaturalRight {
			r.expandNaturalJoin(t)
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
func (r *Resolver) expandNaturalJoin(j *ast.JoinClause) {
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

// AmbiguousColumnError is returned when a column reference matches multiple tables.
type AmbiguousColumnError struct {
	Column string
	Tables []string
}

func (e *AmbiguousColumnError) Error() string {
	return fmt.Sprintf("column %q is ambiguous, found in tables: %s", e.Column, strings.Join(e.Tables, ", "))
}
