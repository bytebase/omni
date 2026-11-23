// Package deparse — resolver.go implements schema-aware column qualification.
// The resolver takes a TableLookup + SelectStmt and returns a new SelectStmt where
// all column references are fully qualified with their table name/alias.
package deparse

import (
	"fmt"
	"sort"
	"strings"

	ast "github.com/bytebase/omni/mysql/ast"
	scopepkg "github.com/bytebase/omni/mysql/scope"
)

// ResolverColumn is a type alias for scope.Column, preserving backward compatibility.
type ResolverColumn = scopepkg.Column

// ResolverTable is a type alias for scope.Table, preserving backward compatibility.
type ResolverTable = scopepkg.Table

// TableLookup is a function that looks up a table by name in the catalog.
// Returns nil if the table is not found.
type TableLookup func(tableName string) *ResolverTable

// SchemaTableLookup looks up a table by schema and name in the catalog.
// Returns nil if the schema or table is not found.
type SchemaTableLookup func(schema, tableName string) *ResolverTable

// Resolver resolves column references in a SelectStmt using catalog metadata.
type Resolver struct {
	Lookup TableLookup
	// SchemaLookup optionally resolves schema-qualified table references
	// (e.g., db2.t). When nil, schema-qualified TABLE query primaries are
	// preserved unexpanded rather than guessing from the single-schema Lookup.
	SchemaLookup SchemaTableLookup
	// DefaultCharset is the database's default character set (e.g., "utf8mb4", "latin1").
	// Used to populate CAST(... AS CHAR) charset when not explicitly specified.
	// If empty, defaults to "utf8mb4".
	DefaultCharset string
	// DefaultWeekFormat is the effective session default_week_format (0..7) used
	// to expand 1-arg WEEK(x) -> week(x, mode). A negative value means unknown,
	// preserving the 1-arg call unexpanded.
	DefaultWeekFormat int
	// NormalizeCastAsYear rewrites MySQL 8.0 SHOW CREATE VIEW's internal
	// `cast_as_year`(x) spelling to `cast(x as year)` so engine-readback views
	// converge with executable CAST(... AS YEAR). Left false when a stored
	// function named `cast_as_year` is in scope or for direct Resolver callers.
	NormalizeCastAsYear bool
	// Source is the text the statement was parsed from, which the AST's Loc
	// offsets index into. When it is set, an unaliased select item is named
	// after its source text, as MySQL names it; otherwise the name is rebuilt
	// from the AST, which matches MySQL only for text written the way the
	// rebuilt name spells it.
	Source string
	// DefaultSchema is the schema of the view the statement defines, or ""
	// outside a view. When it is set, table references are written as MySQL
	// stores a view (mysql:8.0.46 SHOW CREATE VIEW): a table of DefaultSchema
	// loses its schema (FROM a.t in a view of a is written from `t`), an
	// alias equal to its table's name is dropped (FROM t AS t is from `t`),
	// and a column of an unaliased table of another schema is qualified by
	// that schema as well (`o`.`x`.`id`).
	DefaultSchema string

	// shadowing holds the CTEs in scope while a FROM clause is read, keyed
	// by lowered name.
	shadowing map[string]*ResolverTable
}

// Resolve takes a SelectStmt and returns a new SelectStmt with all column
// references fully qualified. The original AST is modified in-place.
func (r *Resolver) Resolve(stmt *ast.SelectStmt) *ast.SelectStmt {
	return r.resolveWithCTEs(stmt, nil)
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
func (r *Resolver) resolveParenSource(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable, parentScope *fromScope) {
	mergedCTETables := cteTables
	if len(stmt.CTEs) > 0 {
		mergedCTETables = make(map[string]*ResolverTable)
		for k, v := range cteTables {
			mergedCTETables[k] = v
		}
		for _, cte := range stmt.CTEs {
			cteResolver := &Resolver{
				Lookup:              r.withCTELookup(mergedCTETables),
				SchemaLookup:        r.SchemaLookup,
				DefaultCharset:      r.DefaultCharset,
				DefaultWeekFormat:   r.DefaultWeekFormat,
				NormalizeCastAsYear: r.NormalizeCastAsYear,
				Source:              r.Source,
				DefaultSchema:       r.DefaultSchema,
			}
			cteResolver.resolveWithCTEs(cte.Select, mergedCTETables)
			if vt := buildCTEVirtualTable(cte, cteResolver.Lookup, r.SchemaLookup); vt != nil {
				mergedCTETables[strings.ToLower(cte.Name)] = vt
			}
		}
	}

	r.resolveWithScope(stmt.ParenSource, mergedCTETables, parentScope)

	if len(stmt.OrderBy) == 0 {
		return
	}
	// Outer ORDER BY ordinals resolve against the inner's leftmost output SELECT.
	resolveOrderByOrdinals(stmt.OrderBy, ast.LeftmostQueryLeaf(stmt.ParenSource))
}

func (r *Resolver) resolveWithCTEs(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable) *ast.SelectStmt {
	return r.resolveWithScope(stmt, cteTables, nil)
}

func (r *Resolver) resolveWithScope(stmt *ast.SelectStmt, cteTables map[string]*ResolverTable, parentScope *fromScope) *ast.SelectStmt {
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
				cteResolver := &Resolver{
					Lookup:              r.withCTELookup(mergedCTETables),
					SchemaLookup:        r.SchemaLookup,
					DefaultCharset:      r.DefaultCharset,
					DefaultWeekFormat:   r.DefaultWeekFormat,
					NormalizeCastAsYear: r.NormalizeCastAsYear,
					Source:              r.Source,
					DefaultSchema:       r.DefaultSchema,
				}
				cteResolver.resolveWithCTEs(cte.Select, mergedCTETables)
				vt := buildCTEVirtualTable(cte, cteResolver.Lookup, r.SchemaLookup)
				if vt != nil {
					mergedCTETables[strings.ToLower(cte.Name)] = vt
				}
			}
		}
		if stmt.Left != nil {
			r.resolveWithScope(stmt.Left, mergedCTETables, parentScope)
		}
		if stmt.Right != nil {
			r.resolveWithScope(stmt.Right, mergedCTETables, parentScope)
		}
		// Resolve ORDER BY ordinals (e.g., ORDER BY 1) to column aliases
		// from the leftmost SELECT's target list, matching MySQL 8.0 behavior.
		if len(stmt.OrderBy) > 0 {
			resolveOrderByOrdinals(stmt.OrderBy, ast.LeftmostQueryLeaf(stmt))
		}
		return stmt
	}

	// Parenthesized query expression: resolve the inner query (with the
	// wrapper's CTEs visible) and the outer ORDER BY ordinals.
	if stmt.ParenSource != nil {
		r.resolveParenSource(stmt, cteTables, parentScope)
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
			leftResolver := &Resolver{
				Lookup:              r.withCTELookup(localCTETables),
				SchemaLookup:        r.SchemaLookup,
				DefaultCharset:      r.DefaultCharset,
				DefaultWeekFormat:   r.DefaultWeekFormat,
				NormalizeCastAsYear: r.NormalizeCastAsYear,
				Source:              r.Source,
				DefaultSchema:       r.DefaultSchema,
			}
			leftResolver.resolveWithCTEs(sel.Left, localCTETables)

			// Step 2: Build virtual table from left branch's target list and add CTE to scope
			vt := buildCTEVirtualTableFromSelect(cte.Name, cte.Columns, sel.Left, leftResolver.Lookup, r.SchemaLookup)
			if vt != nil {
				localCTETables[strings.ToLower(cte.Name)] = vt
			}

			// Step 3: Resolve the recursive (right) branch — CTE is now in scope
			rightResolver := &Resolver{
				Lookup:              r.withCTELookup(localCTETables),
				SchemaLookup:        r.SchemaLookup,
				DefaultCharset:      r.DefaultCharset,
				DefaultWeekFormat:   r.DefaultWeekFormat,
				NormalizeCastAsYear: r.NormalizeCastAsYear,
				Source:              r.Source,
				DefaultSchema:       r.DefaultSchema,
			}
			rightResolver.resolveWithCTEs(sel.Right, localCTETables)
		} else {
			// Non-recursive CTE: resolve entire SELECT, then build virtual table
			cteResolver := &Resolver{
				Lookup:              r.withCTELookup(localCTETables),
				SchemaLookup:        r.SchemaLookup,
				DefaultCharset:      r.DefaultCharset,
				DefaultWeekFormat:   r.DefaultWeekFormat,
				NormalizeCastAsYear: r.NormalizeCastAsYear,
				Source:              r.Source,
				DefaultSchema:       r.DefaultSchema,
			}
			cteResolver.resolveWithCTEs(cte.Select, localCTETables)

			vt := buildCTEVirtualTable(cte, cteResolver.Lookup, r.SchemaLookup)
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
	// Expand an explicit TABLE query primary into SELECT * FROM t so the
	// resolved view body matches MySQL 8.0's SHOW CREATE VIEW form
	// (`select t.c1 AS c1, ... from t`). Direct DeparseSelect without Resolver
	// keeps `table t`. For a schema-qualified TABLE db.t, only expand when a
	// schema-aware SchemaLookup can resolve the requested schema so we never
	// substitute columns from another schema's same-named table.
	if stmt.TableSource != nil && stmt.TableSource.Table != nil {
		tbl := stmt.TableSource.Table
		if r.lookupTableRef(tbl, origLookup, false) != nil {
			stmt.From = []ast.TableExpr{tbl}
			stmt.TargetList = []ast.ExprNode{&ast.StarExpr{Loc: stmt.TableSource.Loc}}
			stmt.TableSource = nil
		}
	}
	// A CTE of a table's name keeps that table's schema (viewTableSchema).
	r.shadowing = localCTETables
	sc := r.buildScope(stmt.From, origLookup)
	r.shadowing = nil
	sc.parent = parentScope
	sc.ctes = localCTETables
	r.Lookup = origLookup

	// Resolve JOIN ON conditions (walk FROM clause) BEFORE target list resolution
	// so that NATURAL JOIN expansion can mark coalesced columns for star expansion.
	r.resolveFromExprs(stmt.From, sc)

	// Resolve target list (may expand stars)
	stmt.TargetList = r.resolveTargetList(stmt.TargetList, sc)

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
// baseLookup (when non-nil) resolves schema-qualified table references without
// shadowing by local CTEs.
func (r *Resolver) buildScope(from []ast.TableExpr, baseLookup TableLookup) *fromScope {
	sc := newFromScope(from)
	for _, tbl := range from {
		r.addTableExprToScope(tbl, sc, baseLookup)
	}
	return sc
}

// lookupTableRef resolves a TableRef using SchemaLookup when schema-qualified
// (or optionally falling back to baseLookup when no SchemaLookup is configured)
// and r.Lookup when unqualified.
func (r *Resolver) lookupTableRef(t *ast.TableRef, baseLookup TableLookup, allowSingleSchemaFallback bool) *ResolverTable {
	if t == nil {
		return nil
	}
	if t.Schema != "" {
		if r.SchemaLookup != nil {
			return r.SchemaLookup(t.Schema, t.Name)
		}
		if allowSingleSchemaFallback && baseLookup != nil {
			return baseLookup(t.Name)
		}
		return nil
	}
	if r.Lookup != nil {
		return r.Lookup(t.Name)
	}
	return nil
}

// addTableExprToScope recursively adds table references from a table expression to the scope.
func (r *Resolver) addTableExprToScope(tbl ast.TableExpr, sc *fromScope, baseLookup TableLookup) {
	defer sc.span(tbl, len(sc.tables))
	switch t := tbl.(type) {
	case *ast.TableRef:
		table := r.lookupTableRef(t, baseLookup, true)
		if table == nil {
			return
		}
		schema := r.viewTableSchema(t)
		effectiveName := t.Name
		if t.Alias != "" {
			effectiveName = t.Alias
		}
		sc.add(t, effectiveName, table)
		if schema != "" {
			sc.qualify(len(sc.tables)-1, schema)
		}
	case *ast.JoinClause:
		r.addTableExprToScope(t.Left, sc, baseLookup)
		r.addTableExprToScope(t.Right, sc, baseLookup)
	case *ast.SubqueryExpr:
		// Derived table: resolve the inner SELECT, then build a virtual table
		// from its target list so outer queries can reference derived columns.
		if t.Select != nil {
			r.Resolve(t.Select)
			vt := buildDerivedVirtualTable(t, r.Lookup, r.SchemaLookup)
			if vt != nil && t.Alias != "" {
				sc.add(t, t.Alias, vt)
			}
		}
	}
}

// viewTableSchema rewrites a table reference of a view as MySQL stores it
// (see DefaultSchema) and returns the schema its columns are qualified by, or
// "" when they are qualified by the table's name or alias alone.
func (r *Resolver) viewTableSchema(t *ast.TableRef) string {
	if r.DefaultSchema == "" {
		return ""
	}
	// The schema is kept where a CTE of the table's name would otherwise
	// take the reference: MySQL drops an unreferenced CTE from a stored view
	// (WITH t AS (SELECT 2 AS id) SELECT a.t.id FROM a.t is stored as
	// select `t`.`id` AS `id` from `t`), which the catalog keeps.
	if _, shadowed := r.shadowing[strings.ToLower(t.Name)]; t.Schema != "" && !shadowed && strings.EqualFold(t.Schema, r.DefaultSchema) {
		t.Schema = ""
	}
	if t.Alias == t.Name {
		t.Alias = ""
	}
	if t.Alias != "" {
		return ""
	}
	return t.Schema
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
	r.resolveCastCharsets(expr)

	// Compute auto-alias from the pre-resolution expression when no explicit alias.
	// MySQL 8.0 uses the original (unqualified) expression text for auto-aliases,
	// so we must derive it before column qualification changes the expression.
	if explicitAlias == "" {
		explicitAlias = r.autoName(expr, position)
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

// autoName names an unaliased select item: from its source text when the
// resolver has it, otherwise rebuilt from the AST.
func (r *Resolver) autoName(expr ast.ExprNode, position int) string {
	return AutoName(r.Source, expr, position)
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
		result = append(result, &ast.ResTarget{
			Name: c.name,
			Val: &ast.ColumnRef{
				Schema: c.schema,
				Table:  c.qualifier,
				Column: c.name,
			},
		})
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

	var result []ast.ExprNode
	cols := sortedVisibleResolverColumns(table)
	for _, col := range cols {
		result = append(result, &ast.ResTarget{
			Name: col.Name,
			Val: &ast.ColumnRef{
				Schema: sc.schemaOf(effectiveName),
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

// sortedVisibleResolverColumns returns non-invisible columns sorted by Position
// (for SELECT *, t.*, TABLE t, and NATURAL JOIN expansion).
func sortedVisibleResolverColumns(table *ResolverTable) []ResolverColumn {
	cols := make([]ResolverColumn, 0, len(table.Columns))
	for _, c := range table.Columns {
		if !c.Invisible {
			cols = append(cols, c)
		}
	}
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
			r.resolveWithScope(n.Select, sc.ctes, sc)
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
		// MySQL 8.0 rewrites 1-arg TIME(x) -> cast(x as time[(fsp)]), 1-arg
		// TIMESTAMP(x) -> cast(x as datetime[(fsp)]) using the argument's type,
		// and 1-arg WEEK(x) -> week(x, default_week_format) using the session mode.
		if n.Schema == "" && !n.SpacedParen && len(n.Args) == 1 && n.Over == nil {
			upper := strings.ToUpper(n.Name)
			switch upper {
			case "TIME", "TIMESTAMP":
				if fsp, ok := exprTemporalFSP(n.Args[0], sc); ok {
					targetType := "time"
					if upper == "TIMESTAMP" {
						targetType = "datetime"
					}
					return &ast.CastExpr{
						Loc:      n.Loc,
						Expr:     n.Args[0],
						TypeName: &ast.DataType{Name: targetType, Length: fsp},
					}
				}
			case "WEEK":
				if r.DefaultWeekFormat >= 0 {
					n.Args = append(n.Args, &ast.IntLit{Value: int64(r.DefaultWeekFormat)})
				}
			case "CAST_AS_YEAR":
				if r.NormalizeCastAsYear {
					return &ast.CastExpr{
						Loc:      n.Loc,
						Expr:     n.Args[0],
						TypeName: &ast.DataType{Name: "year"},
					}
				}
			}
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
	case *ast.PositionExpr:
		n.Substr = r.resolveExpr(n.Substr, sc)
		n.Str = r.resolveExpr(n.Str, sc)
		return n
	case *ast.CollateExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.IntervalExpr:
		n.Value = r.resolveExpr(n.Value, sc)
		return n
	case *ast.ExtractExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.WeightStringExpr:
		n.Expr = r.resolveExpr(n.Expr, sc)
		return n
	case *ast.RowExpr:
		for i, item := range n.Items {
			n.Items[i] = r.resolveExpr(item, sc)
		}
		return n
	case *ast.ExistsExpr:
		if n.Select != nil {
			r.resolveWithScope(n.Select, sc.ctes, sc)
		}
		return n
	case *ast.SubqueryExpr:
		if n.Select != nil {
			r.resolveWithScope(n.Select, sc.ctes, sc)
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
// table in scope (or an enclosing correlated scope) contains the column.
func (r *Resolver) resolveColumnRef(col *ast.ColumnRef, sc *fromScope) ast.ExprNode {
	// Already qualified — just validate the table name maps to an alias
	if col.Table != "" {
		for cur := sc; cur != nil; cur = cur.parent {
			if cur.GetTable(col.Table) != nil {
				// Find the effective name from scope (preserves case)
				if entry, ok := cur.FirstEqualFold(col.Table); ok {
					col.Table = entry.Name
					if r.DefaultSchema != "" {
						col.Schema = cur.schemaOf(entry.Name)
					}
				}
				return col
			}
		}
		return col
	}

	for cur := sc; cur != nil; cur = cur.parent {
		// A column a NATURAL or USING join joins on is the one copy the join
		// keeps, which comes from its right input in a RIGHT join.
		if c, ok := cur.coalesced[strings.ToLower(col.Column)]; ok {
			col.Schema, col.Table = c.schema, c.qualifier
			return col
		}

		// Unqualified — qualify with the first table in scope holding the column.
		// When more than one does, MySQL would raise ERROR 1052: Column 'x' in
		// field list is ambiguous. TODO: return error
		if entry, ok := cur.FirstWithColumn(col.Column); ok {
			col.Schema, col.Table = cur.schemaOf(entry.Name), entry.Name
			return col
		}
	}
	// Column not found — return as-is (could be a literal alias or error)
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
			if on, ok := sc.joinOn(j, names, false); ok {
				j.Condition = &ast.OnCondition{Expr: on}
			}
		} else {
			// Inputs that share no column still make a NATURAL join, which
			// lists a RIGHT join's right input first; an outer join needs a
			// condition, and MySQL writes on(true) (SELECT * FROM p NATURAL
			// RIGHT JOIN q).
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
	if on, ok := sc.joinOn(j, using.Columns, true); ok {
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
	case *ast.ExtractExpr:
		r.resolveCastCharsets(n.Expr)
	case *ast.WeightStringExpr:
		r.resolveCastCharsets(n.Expr)
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
func buildCTEVirtualTableFromSelect(cteName string, cteColumns []string, sel *ast.SelectStmt, lookup TableLookup, schemaLookup SchemaTableLookup) *ResolverTable {
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

	cols := queryOutputColumns(sel, lookup, schemaLookup)
	if len(cols) == 0 {
		return nil
	}
	return &ResolverTable{Name: cteName, Columns: cols}
}

// buildCTEVirtualTable constructs a ResolverTable from a CTE's SELECT target list.
// This allows the main query to resolve column references to CTE columns.
func buildCTEVirtualTable(cte *ast.CommonTableExpr, lookup TableLookup, schemaLookup SchemaTableLookup) *ResolverTable {
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
	cols := queryOutputColumns(cte.Select, lookup, schemaLookup)
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
	return PrimaryOutputColumnsWithSchema(sel, lookup, nil)
}

// PrimaryOutputColumnsWithSchema is PrimaryOutputColumns with an optional
// schema-aware lookup for schema-qualified TABLE primaries (e.g., TABLE db2.t).
func PrimaryOutputColumnsWithSchema(sel *ast.SelectStmt, lookup TableLookup, schemaLookup SchemaTableLookup) []ResolverColumn {
	leaf := ast.LeftmostQueryLeaf(sel)
	if leaf == nil || (leaf.TableSource == nil && leaf.ValuesSource == nil) {
		return nil
	}
	return queryOutputColumns(sel, lookup, schemaLookup)
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
func queryOutputColumns(sel *ast.SelectStmt, lookup TableLookup, schemaLookup SchemaTableLookup) []ResolverColumn {
	return queryOutputColumnsIn(sel, nil, lookup, schemaLookup, nil)
}

// queryOutputColumnsIn is queryOutputColumns with the WITH scopes that enclose
// sel from outside it (innermost first) and the set of CTE declarations whose
// bodies are already being read. Both exist for one reason: a TABLE primary may name
// a CTE rather than a table, that CTE's body may be a primary naming another
// CTE, and a body may name itself.
func queryOutputColumnsIn(sel *ast.SelectStmt, outer []cteScope, lookup TableLookup, schemaLookup SchemaTableLookup, seen map[*ast.CommonTableExpr]bool) []ResolverColumn {
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
			if cols := cteOutputColumns(scopes, sel.TableSource.Table.Name, lookup, schemaLookup, seen); cols != nil {
				return cols
			}
		}
		var table *ResolverTable
		if sel.TableSource.Table.Schema != "" {
			if schemaLookup != nil {
				table = schemaLookup(sel.TableSource.Table.Schema, sel.TableSource.Table.Name)
			}
		} else if lookup != nil {
			table = lookup(sel.TableSource.Table.Name)
		}
		if table == nil {
			return nil
		}
		sorted := sortedVisibleResolverColumns(table)
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
func cteOutputColumns(scopes []cteScope, name string, lookup TableLookup, schemaLookup SchemaTableLookup, seen map[*ast.CommonTableExpr]bool) []ResolverColumn {
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
	return queryOutputColumnsIn(cte.Select, rest, lookup, schemaLookup, next)
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
// The scopes are searched innermost first, so a nested WITH shadows an outer
// declaration of the same name:
// `WITH t AS (SELECT 1 AS z) (WITH t AS (SELECT 2 AS q) TABLE t)` outputs q on
// mysql:8.0 (container-verified). v8.5.5 answers z for the same query and
// drops the inner WITH from its own canonical body, so the TiDB tree keeps the
// outer declaration there; see its copy of this function.
func findCTE(scopes []cteScope, name string) (*ast.CommonTableExpr, []cteScope) {
	for i, scope := range scopes {
		for j, cte := range scope {
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
func buildDerivedVirtualTable(sub *ast.SubqueryExpr, lookup TableLookup, schemaLookup SchemaTableLookup) *ResolverTable {
	if sub.Select == nil || sub.Alias == "" {
		return nil
	}

	cols := queryOutputColumns(sub.Select, lookup, schemaLookup)
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

// exprTemporalFSP returns the fractional seconds precision MySQL 8.0 uses when
// rewriting 1-arg TIME(expr) -> cast(expr as time[(fsp)]) or 1-arg
// TIMESTAMP(expr) -> cast(expr as datetime[(fsp)]). Returns (0, false) when
// the argument's type is not known.
func exprTemporalFSP(expr ast.ExprNode, sc *fromScope) (int, bool) {
	expr = unwrapParen(expr)
	switch n := expr.(type) {
	case *ast.IntLit, *ast.BoolLit, *ast.NullLit:
		return 0, true
	case *ast.StringLit, *ast.HexLit, *ast.BitLit:
		return 6, true
	case *ast.FloatLit:
		if strings.ContainsAny(n.Value, "eE") {
			return 6, true
		}
		if _, frac, ok := strings.Cut(n.Value, "."); ok {
			return clampFSP(len(frac)), true
		}
		return 0, true
	case *ast.TemporalLit:
		return temporalLitFSP(n.Type, n.Value), true
	case *ast.CastExpr:
		if n.TypeName == nil {
			return 0, false
		}
		switch strings.ToLower(n.TypeName.Name) {
		case "signed", "signed integer", "unsigned", "unsigned integer", "date", "year":
			return 0, true
		case "datetime", "time":
			return clampFSP(n.TypeName.Length), true
		case "decimal":
			return clampFSP(n.TypeName.Scale), true
		case "char", "binary", "float", "double", "json":
			return 6, true
		}
		return 0, false
	case *ast.ColumnRef:
		col := lookupScopeColumn(n, sc)
		if col == nil || col.DataType == "" {
			return 0, false
		}
		switch strings.ToLower(col.DataType) {
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "year", "date", "bit":
			return 0, true
		case "datetime", "timestamp", "time", "decimal", "numeric":
			return clampFSP(col.Scale), true
		case "float", "double",
			"char", "varchar", "tinytext", "text", "mediumtext", "longtext",
			"binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob",
			"enum", "set", "json":
			return 6, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func lookupScopeColumn(ref *ast.ColumnRef, sc *fromScope) *ResolverColumn {
	if ref == nil || ref.Star {
		return nil
	}
	for cur := sc; cur != nil; cur = cur.parent {
		if ref.Table != "" {
			if tbl := cur.GetTable(ref.Table); tbl != nil {
				return tbl.GetColumn(ref.Column)
			}
			continue
		}
		if entry, ok := cur.FirstWithColumn(ref.Column); ok && entry.Table != nil {
			return entry.Table.GetColumn(ref.Column)
		}
	}
	return nil
}

func clampFSP(n int) int {
	if n < 0 {
		return 0
	}
	if n > 6 {
		return 6
	}
	return n
}

// temporalLitFSP returns the fractional seconds precision of a MySQL DATE,
// TIME, or TIMESTAMP literal. Unlike splitting at the first '.', this
// distinguishes date/time field delimiters (such as '.' in
// TIMESTAMP '2020.01.01 12:34:56.1' or '2020.01.01 12.34.56') from the
// fractional-seconds suffix and stops before any trailing timezone suffix.
func temporalLitFSP(typ, val string) int {
	if strings.EqualFold(typ, "DATE") {
		return 0
	}
	s := strings.TrimSpace(val)
	if s == "" {
		return 0
	}
	if strings.EqualFold(typ, "TIME") {
		t := s
		if t[0] == '+' || t[0] == '-' {
			t = strings.TrimSpace(t[1:])
		}
		plain := true
		dotIdx := -1
		for i := 0; i < len(t); i++ {
			ch := t[i]
			switch {
			case ch >= '0' && ch <= '9', ch == ':', ch == ' ':
			case ch == '.':
				if dotIdx >= 0 {
					plain = false
				} else {
					dotIdx = i
				}
			default:
				plain = false
			}
		}
		if plain {
			if dotIdx < 0 {
				return 0
			}
			n := 0
			for i := dotIdx + 1; i < len(t) && t[i] >= '0' && t[i] <= '9'; i++ {
				n++
			}
			return clampFSP(n)
		}
	}
	return datetimeLitFSP(s)
}

func datetimeLitFSP(s string) int {
	var parts []string
	var seps []string
	i := 0
	for i < len(s) {
		sepStart := i
		for i < len(s) && (s[i] < '0' || s[i] > '9') {
			i++
		}
		sep := s[sepStart:i]
		if i >= len(s) {
			break
		}
		digStart := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		seps = append(seps, sep)
		parts = append(parts, s[digStart:i])
	}
	if len(parts) == 0 {
		return 0
	}
	n0 := len(parts[0])
	if n0 >= 5 {
		filled := 3
		switch {
		case n0 >= 11:
			filled = 6
		case n0 >= 9:
			filled = 5
		case n0 == 7:
			filled = 4
		}
		for idx := 1; idx < len(parts); idx++ {
			if filled < 6 {
				if filled == 3 && len(parts[idx]) >= 5 {
					filled = 6
				} else {
					filled++
				}
			} else {
				if seps[idx] == "." {
					return clampFSP(len(parts[idx]))
				}
				return 0
			}
		}
		return 0
	}
	if len(parts) >= 7 && seps[6] == "." {
		return clampFSP(len(parts[6]))
	}
	return 0
}
