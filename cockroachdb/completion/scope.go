package completion

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// placeholder stands in for the word at the caret. It is a bare identifier
// no keyword spells, so the grammar reads it wherever a name fits.
const placeholder = "omni_completion_caret"

type nameKind int

const (
	kindNone nameKind = iota
	// kindObject is a table-like name: a table, view, sequence or schema.
	kindObject
	// kindSchema is the schema part of a qualified object name.
	kindSchema
	// kindExpr is a column reference or function name in an expression.
	kindExpr
	// kindFunction is the name of a function call.
	kindFunction
	// kindTargetColumn is a column of one known table, as in INSERT INTO t (|).
	kindTargetColumn
	// kindIndex is an index of one known table, as in t@|.
	kindIndex
	// kindType is a type name, as in CAST(x AS |) or a column definition.
	kindType
	// kindColumnList is a column of a known list, as in JOIN ... USING (|).
	kindColumnList
)

// source is a relation a name can reference: a catalog relation or a derived
// one (CTE, subquery, function) with its own columns.
type source struct {
	ref     string // the alias, or the relation's own name
	schema  string // the catalog relation's schema; "" resolves by search path
	name    string // the catalog relation's name; "" for a derived source
	columns []string
	derived bool
	// depth is the nesting level of the query whose FROM holds the source,
	// so a qualified name resolves in the innermost query that defines it.
	depth int
}

// nameContext says what the word at the caret names and which sources and
// CTEs are in scope there.
type nameContext struct {
	kind      nameKind
	qualifier []string // outermost first: [schema] or [table] or [schema, table]
	scope     []source // innermost first
	target    *source
	ctes      []string
	columns   []string // the names kindColumnList offers
	// scoped is set when the analyzer walked the statement to the name, so
	// scope is complete even when empty.
	scoped bool
}

// analyzer walks a statement once, finding placeholder and naming the
// columns of every relation it passes: each CTE, subquery and statement
// source is walked once, so the work is linear in the statement. It stops
// at the placeholder or when ctx is done.
type analyzer struct {
	ctx   context.Context
	cat   *catalog
	ctes  [][]source // a stack of WITH clauses, innermost last
	depth int
	// search is unset when the walk only names columns, as for a view's
	// definition.
	search bool
	found  *nameContext
}

// stopped reports whether the walk is over: the placeholder is found or the
// time is up.
func (a *analyzer) stopped() bool {
	return a.found != nil || a.ctx.Err() != nil
}

// analyzeStatement returns the context of placeholder in stmt, or nil when
// stmt does not hold it as a name completion can resolve.
func analyzeStatement(ctx context.Context, stmt tree.Statement, cat *catalog) *nameContext {
	a := &analyzer{ctx: ctx, cat: cat, search: true}
	a.statement(stmt, nil)
	if ctx.Err() != nil {
		return a.found
	}
	if a.found == nil {
		a.found = typeContext(stmt)
	}
	if a.found == nil {
		a.found = tableNameContext(stmt)
	}
	if a.found == nil {
		a.found = genericExprContext(stmt)
	}
	return a.found
}

// tableNameContext finds placeholder in any table-like name of stmt, through
// the formatter hook CockroachDB calls on every TableName and
// UnresolvedObjectName.
func tableNameContext(stmt tree.Statement) (found *nameContext) {
	defer func() {
		if recover() != nil {
			found = nil
		}
	}()
	fmtCtx := tree.NewFmtCtx(tree.FmtSimple, tree.FmtReformatTableNames(func(ctx *tree.FmtCtx, tn *tree.TableName) {
		ctx.WriteString("t")
		if found != nil {
			return
		}
		switch {
		case string(tn.ObjectName) == placeholder:
			found = &nameContext{kind: kindObject}
			if tn.ExplicitSchema {
				found.qualifier = []string{string(tn.SchemaName)}
			}
		case tn.ExplicitSchema && string(tn.SchemaName) == placeholder:
			found = &nameContext{kind: kindSchema}
		}
	}))
	fmtCtx.FormatNode(stmt)
	return found
}

// genericExprContext finds placeholder as a column reference in any
// expression of a statement the analyzer does not walk.
func genericExprContext(stmt tree.Statement) (found *nameContext) {
	defer func() {
		if recover() != nil {
			found = nil
		}
	}()
	_, _ = tree.SimpleStmtVisit(stmt, func(expr tree.Expr) (bool, tree.Expr, error) {
		if found != nil {
			return false, expr, nil
		}
		if kind, qual, ok := placeholderInExpr(expr); ok {
			found = &nameContext{kind: kind, qualifier: qual}
			return false, expr, nil
		}
		return true, expr, nil
	})
	return found
}

// placeholderInExpr reports whether expr itself is placeholder as a column
// reference or a function name, with the name's qualifier.
func placeholderInExpr(expr tree.Expr) (nameKind, []string, bool) {
	switch e := expr.(type) {
	case *tree.UnresolvedName:
		if qual, ok := placeholderName(e); ok {
			return kindExpr, qual, true
		}
	case *tree.FuncExpr:
		if name, ok := e.Func.FunctionReference.(*tree.UnresolvedName); ok {
			if qual, ok := placeholderName(name); ok {
				return kindFunction, qual, true
			}
		}
	case *tree.CastExpr:
		if qual, ok := placeholderType(e.Type); ok {
			return kindType, qual, true
		}
	case *tree.AnnotateTypeExpr:
		if qual, ok := placeholderType(e.Type); ok {
			return kindType, qual, true
		}
	case *tree.IsOfTypeExpr:
		for _, typ := range e.Types {
			if qual, ok := placeholderType(typ); ok {
				return kindType, qual, true
			}
		}
	}
	return kindNone, nil, false
}

// placeholderType reports whether ref is the placeholder as a type name,
// which the parser leaves unresolved, with the name's qualifier.
func placeholderType(ref tree.ResolvableTypeReference) ([]string, bool) {
	name, ok := ref.(*tree.UnresolvedObjectName)
	if !ok || name.NumParts < 1 || name.Parts[0] != placeholder {
		return nil, false
	}
	var qual []string
	for i := name.NumParts - 1; i >= 1; i-- {
		qual = append(qual, name.Parts[i])
	}
	return qual, true
}

// typeContext finds placeholder as a column's type in DDL, or as a cast's
// type in any expression. It runs before tableNameContext, whose formatter
// hook also sees an unresolved type name.
func typeContext(stmt tree.Statement) *nameContext {
	var defs []*tree.ColumnTableDef
	switch s := stmt.(type) {
	case *tree.CreateTable:
		for _, def := range s.Defs {
			if col, ok := def.(*tree.ColumnTableDef); ok {
				defs = append(defs, col)
			}
		}
	case *tree.AlterTable:
		for _, cmd := range s.Cmds {
			switch c := cmd.(type) {
			case *tree.AlterTableAddColumn:
				defs = append(defs, c.ColumnDef)
			case *tree.AlterTableAlterColumnType:
				if qual, ok := placeholderType(c.ToType); ok {
					return &nameContext{kind: kindType, qualifier: qual}
				}
			}
		}
	}
	for _, def := range defs {
		if def == nil {
			continue
		}
		if qual, ok := placeholderType(def.Type); ok {
			return &nameContext{kind: kindType, qualifier: qual}
		}
	}
	if nc := genericExprContext(stmt); nc != nil && nc.kind == kindType {
		return nc
	}
	return nil
}

func placeholderName(n *tree.UnresolvedName) ([]string, bool) {
	if n.Star || n.NumParts < 1 || n.Parts[0] != placeholder {
		return nil, false
	}
	var qual []string
	for i := n.NumParts - 1; i >= 1; i-- {
		qual = append(qual, n.Parts[i])
	}
	return qual, true
}

// statement walks stmt and returns the names of the columns it returns: a
// query's output, or a data change's RETURNING list.
func (a *analyzer) statement(stmt tree.Statement, outer []source) []string {
	if a.stopped() || stmt == nil {
		return nil
	}
	switch s := stmt.(type) {
	case *tree.Select:
		return a.query(s, outer)
	case *tree.ParenSelect:
		return a.query(s.Select, outer)
	case *tree.Insert:
		return a.insert(s, outer)
	case *tree.Update:
		return a.update(s, outer)
	case *tree.Delete:
		return a.delete(s, outer)
	case *tree.Explain:
		a.statement(s.Statement, outer)
	case *tree.ExplainAnalyze:
		a.statement(s.Statement, outer)
	case *tree.CreateView:
		a.query(s.AsSource, outer)
	case *tree.CreateTable:
		if s.AsSource != nil {
			a.query(s.AsSource, outer)
		}
		for _, e := range createTableExprs(s) {
			a.expr(e, outer)
		}
	case *tree.CreateIndex:
		a.createIndex(s)
	case *tree.AlterTable:
		a.alterTable(s)
	}
	return nil
}

// query walks a query with its WITH, ORDER BY and LIMIT, and returns its
// output names.
func (a *analyzer) query(sel *tree.Select, outer []source) []string {
	if sel == nil || a.stopped() {
		return nil
	}
	popped := a.with(sel.With, outer)
	defer popped()
	scope, cols := a.body(sel.Select, outer)
	if len(sel.OrderBy) > 0 {
		// ORDER BY may also name the query's output columns.
		orderScope := append([]source{{derived: true, columns: cols, depth: a.depth + 1}}, scope...)
		for _, order := range sel.OrderBy {
			a.expr(order.Expr, orderScope)
		}
	}
	if sel.Limit != nil {
		a.expr(sel.Limit.Count, outer)
		a.expr(sel.Limit.Offset, outer)
	}
	return cols
}

// with pushes the CTEs of w, walking each body once, and returns the pop.
// Each CTE sees the ones before it; a recursive one also sees itself, with
// the columns of its non-recursive term.
func (a *analyzer) with(w *tree.With, outer []source) func() {
	if w == nil {
		return func() {}
	}
	level := len(a.ctes)
	a.ctes = append(a.ctes, nil)
	for _, cte := range w.CTEList {
		if a.stopped() {
			break
		}
		self := -1
		if w.Recursive {
			a.ctes[level] = append(a.ctes[level], source{ref: string(cte.Name.Alias), derived: true, columns: aliasColumns(cte.Name)})
			self = len(a.ctes[level]) - 1
		}
		var cols []string
		if sel, ok := cte.Stmt.(*tree.Select); ok && w.Recursive {
			cols = a.recursiveQuery(sel, outer, func(left []string) {
				a.ctes[level][self].columns = overlayAliases(left, cte.Name)
			})
		} else {
			cols = a.statement(cte.Stmt, outer)
		}
		src := source{ref: string(cte.Name.Alias), derived: true, columns: overlayAliases(cols, cte.Name)}
		if self >= 0 {
			a.ctes[level][self] = src
		} else {
			a.ctes[level] = append(a.ctes[level], src)
		}
	}
	return func() { a.ctes = a.ctes[:level] }
}

// recursiveQuery walks a recursive CTE's body: for a UNION, its left term
// first, whose output names seed tells the CTE's own entry before the right
// term reads it.
func (a *analyzer) recursiveQuery(sel *tree.Select, outer []source, seed func([]string)) []string {
	union, ok := sel.Select.(*tree.UnionClause)
	if !ok {
		return a.query(sel, outer)
	}
	popped := a.with(sel.With, outer)
	defer popped()
	left := a.query(union.Left, outer)
	seed(left)
	a.query(union.Right, outer)
	for _, order := range sel.OrderBy {
		a.expr(order.Expr, append([]source{{derived: true, columns: left, depth: a.depth + 1}}, outer...))
	}
	return left
}

// overlayAliases renames the first columns of cols by alias's column list,
// keeping the columns the list does not reach.
func overlayAliases(cols []string, alias tree.AliasClause) []string {
	names := aliasColumns(alias)
	if len(names) < len(cols) {
		names = append(names, cols[len(names):]...)
	}
	return names
}

func aliasColumns(alias tree.AliasClause) []string {
	if len(alias.Cols) == 0 {
		return nil
	}
	cols := make([]string, 0, len(alias.Cols))
	for _, c := range alias.Cols {
		cols = append(cols, string(c.Name))
	}
	return cols
}

// body walks one select body and returns the scope its clauses see and its
// output names.
func (a *analyzer) body(ss tree.SelectStatement, outer []source) ([]source, []string) {
	if a.stopped() {
		return outer, nil
	}
	switch s := ss.(type) {
	case *tree.SelectClause:
		a.depth++
		defer func() { a.depth-- }()
		items := a.tableExprs(s.From.Tables, outer)
		local := itemSources(items)
		scope := slices.Concat(local, outer)
		var cols []string
		for _, e := range s.Exprs {
			a.expr(e.Expr, scope)
			if star, ok := a.starColumns(e.Expr, items); ok {
				cols = append(cols, star...)
			} else {
				cols = append(cols, renderName(e))
			}
		}
		if s.Where != nil {
			a.expr(s.Where.Expr, scope)
		}
		// GROUP BY may also name the output columns, as ORDER BY may.
		groupScope := append([]source{{derived: true, columns: cols, depth: a.depth}}, scope...)
		for _, g := range s.GroupBy {
			a.expr(g, groupScope)
		}
		if s.Having != nil {
			a.expr(s.Having.Expr, scope)
		}
		// DISTINCT ON, as ORDER BY, may name the output columns.
		for _, e := range s.DistinctOn {
			a.expr(e, groupScope)
		}
		for _, def := range s.Window {
			a.windowDef(def, scope)
		}
		return scope, cols
	case *tree.UnionClause:
		left := a.query(s.Left, outer)
		a.query(s.Right, outer)
		return outer, left
	case *tree.ParenSelect:
		return outer, a.query(s.Select, outer)
	case *tree.ValuesClause:
		for _, row := range s.Rows {
			for _, e := range row {
				a.expr(e, outer)
			}
		}
		if len(s.Rows) == 0 {
			return outer, nil
		}
		cols := make([]string, len(s.Rows[0]))
		for i := range cols {
			cols[i] = "column" + strconv.Itoa(i+1)
		}
		return outer, cols
	}
	return outer, nil
}

// windowDef walks a window's partitions, ordering and frame offsets.
func (a *analyzer) windowDef(def *tree.WindowDef, scope []source) {
	for _, e := range def.Partitions {
		a.expr(e, scope)
	}
	for _, order := range def.OrderBy {
		a.expr(order.Expr, scope)
	}
	if def.Frame != nil {
		for _, bound := range []*tree.WindowFrameBound{def.Frame.Bounds.StartBound, def.Frame.Bounds.EndBound} {
			if bound != nil {
				a.expr(bound.OffsetExpr, scope)
			}
		}
	}
}

// fromItem is one item of a FROM list: the relations a qualified name can
// reach in it, and the columns it outputs, which * reads. A join's output
// has each USING or NATURAL column once, while its relations keep theirs.
type fromItem struct {
	sources []source
	columns []string
}

func itemSources(items []fromItem) []source {
	var out []source
	for _, item := range items {
		out = append(out, item.sources...)
	}
	return out
}

// tableExprs returns the items of a FROM list. outer is the scope of the
// enclosing queries.
func (a *analyzer) tableExprs(exprs tree.TableExprs, outer []source) []fromItem {
	var items []fromItem
	for _, te := range exprs {
		if item, ok := a.tableExpr(te, itemSources(items), outer); ok {
			items = append(items, item)
		}
	}
	return items
}

// tableExpr returns one FROM item. preceding are the sources of the items
// before it, which a LATERAL item can read; outer is the scope of the
// enclosing queries, which every item reads.
func (a *analyzer) tableExpr(te tree.TableExpr, preceding, outer []source) (fromItem, bool) {
	if a.stopped() {
		return fromItem{}, false
	}
	switch t := te.(type) {
	case *tree.AliasedTableExpr:
		return a.aliased(t, preceding, outer)
	case *tree.ParenTableExpr:
		return a.tableExpr(t.Expr, preceding, outer)
	case *tree.JoinTableExpr:
		return a.join(t, preceding, outer), true
	case *tree.TableName:
		src := a.tableSource(t, "")
		return fromItem{sources: []source{src}, columns: a.sourceColumns(src, false)}, true
	case *tree.Subquery:
		_, cols := a.body(t.Select, outer)
		return fromItem{sources: []source{{derived: true, columns: cols, depth: a.depth}}, columns: cols}, true
	case *tree.StatementSource:
		cols := a.statement(t.Statement, outer)
		return fromItem{sources: []source{{derived: true, columns: cols, depth: a.depth}}, columns: cols}, true
	}
	return fromItem{}, false
}

// join returns a join's item. Its ON condition sees the two sides and the
// enclosing queries, not the rest of the FROM list; a LATERAL right side
// sees the left. USING and NATURAL merge their columns into one in the
// output.
func (a *analyzer) join(t *tree.JoinTableExpr, preceding, outer []source) fromItem {
	left, _ := a.tableExpr(t.Left, preceding, outer)
	right, _ := a.tableExpr(t.Right, slices.Concat(preceding, left.sources), outer)
	var using []string
	switch cond := t.Cond.(type) {
	case *tree.OnJoinCond:
		a.expr(cond.Expr, slices.Concat(left.sources, right.sources, outer))
	case *tree.UsingJoinCond:
		if a.search && a.found == nil && slices.Contains(cond.Cols, tree.Name(placeholder)) {
			a.found = &nameContext{kind: kindColumnList, columns: commonColumns(left.columns, right.columns), scoped: true}
		}
		for _, name := range cond.Cols {
			using = append(using, string(name))
		}
	case tree.NaturalJoinCond:
		using = commonColumns(left.columns, right.columns)
	}
	item := fromItem{sources: slices.Concat(left.sources, right.sources)}
	merged := map[string]bool{}
	for _, name := range using {
		if slices.Contains(left.columns, name) && slices.Contains(right.columns, name) && !merged[name] {
			item.columns = append(item.columns, name)
			merged[name] = true
		}
	}
	for _, side := range [][]string{left.columns, right.columns} {
		skipped := map[string]bool{}
		for _, c := range side {
			if merged[c] && !skipped[c] {
				skipped[c] = true
				continue
			}
			item.columns = append(item.columns, c)
		}
	}
	return item
}

// commonColumns are the columns both sides of a join output, the names
// JOIN ... USING accepts and NATURAL JOIN merges.
func commonColumns(left, right []string) []string {
	var out []string
	for _, c := range left {
		if slices.Contains(right, c) && !slices.Contains(out, c) {
			out = append(out, c)
		}
	}
	return out
}

func (a *analyzer) aliased(t *tree.AliasedTableExpr, preceding, outer []source) (fromItem, bool) {
	alias := string(t.As.Alias)
	var src source
	switch e := t.Expr.(type) {
	case *tree.TableName:
		src = a.tableSource(e, alias)
		if t.IndexFlags != nil && string(t.IndexFlags.Index) == placeholder && a.search && a.found == nil {
			target := src
			a.found = &nameContext{kind: kindIndex, target: &target}
		}
	case *tree.Subquery:
		inner := outer
		if t.Lateral {
			inner = slices.Concat(preceding, outer)
		}
		_, cols := a.body(e.Select, inner)
		src = source{ref: alias, derived: true, columns: cols}
	case *tree.StatementSource:
		src = source{ref: alias, derived: true, columns: a.statement(e.Statement, outer)}
	case *tree.RowsFromExpr:
		// A function in FROM reads the items before it, as if LATERAL.
		var cols []string
		for _, item := range e.Items {
			a.expr(item, slices.Concat(preceding, outer))
			cols = append(cols, functionColumns(item)...)
		}
		ref := alias
		if len(e.Items) == 1 {
			if ref == "" {
				ref = renderName(tree.SelectExpr{Expr: e.Items[0]})
			} else if len(cols) == 1 && len(t.As.Cols) == 0 {
				// A function's one column takes the alias's name.
				cols[0] = alias
			}
		}
		src = source{ref: ref, derived: true, columns: cols}
	case *tree.ParenTableExpr, *tree.JoinTableExpr:
		inner, ok := a.tableExpr(e, preceding, outer)
		if !ok || (alias == "" && !t.Ordinality) {
			// An unaliased join keeps its tables' own names.
			return inner, ok
		}
		// An aliased join is one relation with the join's output columns.
		src = source{ref: alias, derived: true, columns: inner.columns}
	default:
		return fromItem{}, false
	}
	if t.Ordinality {
		src.columns = append(slices.Clone(a.sourceColumns(src, false)), "ordinality")
		src.derived, src.name = true, ""
	}
	if len(t.As.Cols) > 0 {
		src.columns = overlayAliases(a.sourceColumns(src, false), t.As)
		src.derived, src.name = true, ""
	}
	src.depth = a.depth
	return fromItem{sources: []source{src}, columns: a.sourceColumns(src, false)}, true
}

// tableSource resolves a table name to a CTE in scope or a catalog relation.
func (a *analyzer) tableSource(tn *tree.TableName, alias string) source {
	name := string(tn.ObjectName)
	if a.search && a.found == nil {
		switch {
		case name == placeholder:
			a.found = &nameContext{kind: kindObject, ctes: a.visibleCTEs()}
			if tn.ExplicitSchema {
				a.found.qualifier = []string{string(tn.SchemaName)}
			}
		case tn.ExplicitSchema && string(tn.SchemaName) == placeholder:
			a.found = &nameContext{kind: kindSchema}
		}
	}
	ref := alias
	if ref == "" {
		ref = name
	}
	if !tn.ExplicitSchema {
		for i := len(a.ctes) - 1; i >= 0; i-- {
			for _, cte := range a.ctes[i] {
				if cte.ref == name {
					return source{ref: ref, derived: true, columns: cte.columns, depth: a.depth}
				}
			}
		}
	}
	src := source{ref: ref, name: name, depth: a.depth}
	if tn.ExplicitSchema {
		src.schema = string(tn.SchemaName)
	}
	return src
}

// visibleCTEs names the CTEs the WITH clauses open at this point define.
func (a *analyzer) visibleCTEs() []string {
	var out []string
	for _, level := range a.ctes {
		for _, cte := range level {
			out = append(out, cte.ref)
		}
	}
	return out
}

// expr looks for placeholder in e, and walks e's subqueries with scope as
// their outer scope.
func (a *analyzer) expr(e tree.Expr, scope []source) {
	if e == nil || a.stopped() {
		return
	}
	_, _ = tree.SimpleVisit(e, func(expr tree.Expr) (bool, tree.Expr, error) {
		if a.stopped() {
			return false, expr, nil
		}
		if a.search {
			if kind, qual, ok := placeholderInExpr(expr); ok {
				a.found = &nameContext{kind: kind, qualifier: qual, scope: scope, scoped: true}
				return false, expr, nil
			}
		}
		if sub, ok := expr.(*tree.Subquery); ok {
			a.body(sub.Select, scope)
			return false, expr, nil
		}
		// The expression walk skips a call's OVER clause.
		if fn, ok := expr.(*tree.FuncExpr); ok && fn.WindowDef != nil {
			a.windowDef(fn.WindowDef, scope)
		}
		return true, expr, nil
	})
}

// target resolves the table a DML or DDL statement changes.
func (a *analyzer) target(te tree.TableExpr) *source {
	if aliased, ok := te.(*tree.AliasedTableExpr); ok {
		if tn, ok := aliased.Expr.(*tree.TableName); ok {
			src := a.tableSource(tn, string(aliased.As.Alias))
			return &src
		}
	}
	if tn, ok := te.(*tree.TableName); ok {
		src := a.tableSource(tn, "")
		return &src
	}
	return nil
}

func (a *analyzer) targetColumn(names tree.NameList, target *source) {
	if !a.search || a.found != nil || target == nil {
		return
	}
	if slices.Contains(names, tree.Name(placeholder)) {
		a.found = &nameContext{kind: kindTargetColumn, target: target}
	}
}

// dmlScope opens a data change's query level: its target, then the
// tables its FROM or USING adds.
func (a *analyzer) dmlScope(target *source, from tree.TableExprs, outer []source) []source {
	a.depth++
	sources := itemSources(a.tableExprs(from, outer))
	if target != nil {
		target.depth = a.depth
		sources = append([]source{*target}, sources...)
	}
	return slices.Concat(sources, outer)
}

func (a *analyzer) insert(s *tree.Insert, outer []source) []string {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	a.targetColumn(s.Columns, target)
	a.query(s.Rows, outer)
	scope := a.dmlScope(target, nil, outer)
	defer func() { a.depth-- }()
	if oc := s.OnConflict; oc != nil {
		a.targetColumn(oc.Columns, target)
		conflictScope := scope
		if target != nil {
			excluded := *target
			excluded.ref = "excluded"
			conflictScope = append([]source{*target, excluded}, outer...)
		}
		a.expr(oc.ArbiterPredicate, conflictScope)
		a.updateExprs(oc.Exprs, target, conflictScope)
		if oc.Where != nil {
			a.expr(oc.Where.Expr, conflictScope)
		}
	}
	return a.returning(s.Returning, target, scope)
}

func (a *analyzer) update(s *tree.Update, outer []source) []string {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	scope := a.dmlScope(target, s.From, outer)
	defer func() { a.depth-- }()
	a.updateExprs(s.Exprs, target, scope)
	if s.Where != nil {
		a.expr(s.Where.Expr, scope)
	}
	for _, order := range s.OrderBy {
		a.expr(order.Expr, scope)
	}
	return a.returning(s.Returning, target, scope)
}

func (a *analyzer) delete(s *tree.Delete, outer []source) []string {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	scope := a.dmlScope(target, s.Using, outer)
	defer func() { a.depth-- }()
	if s.Where != nil {
		a.expr(s.Where.Expr, scope)
	}
	for _, order := range s.OrderBy {
		a.expr(order.Expr, scope)
	}
	return a.returning(s.Returning, target, scope)
}

func (a *analyzer) updateExprs(exprs tree.UpdateExprs, target *source, scope []source) {
	for _, ue := range exprs {
		a.targetColumn(ue.Names, target)
		a.expr(ue.Expr, scope)
	}
}

// returning walks a RETURNING list and returns its output names; a star
// stands for the target's visible columns.
func (a *analyzer) returning(r tree.ReturningClause, target *source, scope []source) []string {
	exprs, ok := r.(*tree.ReturningExprs)
	if !ok {
		return nil
	}
	var targets []fromItem
	if target != nil {
		targets = []fromItem{{sources: []source{*target}, columns: a.sourceColumns(*target, false)}}
	}
	var cols []string
	for _, e := range *exprs {
		a.expr(e.Expr, scope)
		if star, ok := a.starColumns(e.Expr, targets); ok {
			cols = append(cols, star...)
		} else {
			cols = append(cols, renderName(tree.SelectExpr(e)))
		}
	}
	return cols
}

func (a *analyzer) createIndex(s *tree.CreateIndex) {
	target := a.tableSource(&s.Table, "")
	scope := []source{target}
	for _, elem := range s.Columns {
		if string(elem.Column) == placeholder {
			a.found = &nameContext{kind: kindTargetColumn, target: &target}
			return
		}
		a.expr(elem.Expr, scope)
	}
	a.targetColumn(s.Storing, &target)
	a.expr(s.Predicate, scope)
}

func (a *analyzer) alterTable(s *tree.AlterTable) {
	if s.Table == nil {
		return
	}
	tn := s.Table.ToTableName()
	target := a.tableSource(&tn, "")
	for _, cmd := range s.Cmds {
		var col tree.Name
		switch c := cmd.(type) {
		case *tree.AlterTableDropColumn:
			col = c.Column
		case *tree.AlterTableAlterColumnType:
			col = c.Column
		case *tree.AlterTableRenameColumn:
			col = c.Column
		case *tree.AlterTableSetDefault:
			col = c.Column
		case *tree.AlterTableSetNotNull:
			col = c.Column
		case *tree.AlterTableDropNotNull:
			col = c.Column
		case *tree.AlterTableDropStored:
			col = c.Column
		case *tree.AlterTableSetVisible:
			col = c.Column
		}
		if string(col) == placeholder {
			a.found = &nameContext{kind: kindTargetColumn, target: &target}
			return
		}
	}
	for _, e := range alterTableExprs(s) {
		a.expr(e, []source{target})
	}
}

// columnDefExprs returns the expressions of a column definition.
func columnDefExprs(def *tree.ColumnTableDef) []tree.Expr {
	if def == nil {
		return nil
	}
	exprs := []tree.Expr{def.DefaultExpr.Expr, def.OnUpdateExpr.Expr, def.Computed.Expr}
	for _, check := range def.CheckExprs {
		exprs = append(exprs, check.Expr)
	}
	return exprs
}

// alterTableExprs returns the expressions an ALTER TABLE evaluates against
// its table: defaults, checks, computed columns and USING conversions.
func alterTableExprs(s *tree.AlterTable) []tree.Expr {
	var exprs []tree.Expr
	for _, cmd := range s.Cmds {
		switch c := cmd.(type) {
		case *tree.AlterTableAddColumn:
			exprs = append(exprs, columnDefExprs(c.ColumnDef)...)
		case *tree.AlterTableAddConstraint:
			if check, ok := c.ConstraintDef.(*tree.CheckConstraintTableDef); ok {
				exprs = append(exprs, check.Expr)
			}
		case *tree.AlterTableSetDefault:
			exprs = append(exprs, c.Default)
		case *tree.AlterTableAlterColumnType:
			exprs = append(exprs, c.Using)
		}
	}
	return exprs
}

// createTableExprs returns the expressions of a CREATE TABLE's column and
// table definitions.
func createTableExprs(s *tree.CreateTable) []tree.Expr {
	var exprs []tree.Expr
	for _, def := range s.Defs {
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			exprs = append(exprs, columnDefExprs(d)...)
		case *tree.CheckConstraintTableDef:
			exprs = append(exprs, d.Expr)
		}
	}
	return exprs
}

// starColumns expands a star select expression: * to the output columns of
// the FROM items, t.* to the visible columns of the relation t names; ok
// is false for any other expression.
func (a *analyzer) starColumns(e tree.Expr, items []fromItem) (cols []string, ok bool) {
	switch n := e.(type) {
	case tree.UnqualifiedStar:
		for _, item := range items {
			cols = append(cols, item.columns...)
		}
		return cols, true
	case *tree.UnresolvedName:
		if !n.Star {
			return nil, false
		}
		if n.NumParts == 1 {
			return a.starColumns(tree.UnqualifiedStar{}, items)
		}
		for _, src := range itemSources(items) {
			if src.ref == n.Parts[1] {
				cols = append(cols, a.sourceColumns(src, false)...)
			}
		}
		return cols, true
	}
	return nil, false
}

// functionColumns names the columns a function in FROM returns: key and
// value for json_each and its variants, one per array for unnest of
// several, and otherwise one named after the function. These are the widths
// cockroachdb/analysis's query span settles on (#1795), copied while that
// package exports no table of them.
func functionColumns(e tree.Expr) []string {
	name := renderName(tree.SelectExpr{Expr: e})
	fn, ok := e.(*tree.FuncExpr)
	switch {
	case ok && name == "unnest" && len(fn.Exprs) > 1:
		cols := make([]string, len(fn.Exprs))
		for i := range cols {
			cols[i] = name
		}
		return cols
	case ok && keyValueFunctions[name]:
		return []string{"key", "value"}
	}
	return []string{name}
}

// keyValueFunctions return a key and a value column.
var keyValueFunctions = map[string]bool{
	"json_each": true, "jsonb_each": true, "json_each_text": true, "jsonb_each_text": true,
}

// renderName is the name CockroachDB gives a select expression's column:
// its alias, or a name derived from the expression (count for count(*),
// ?column? for 1 + 1), by the vendor's own GetRenderColName.
func renderName(e tree.SelectExpr) (name string) {
	defer func() {
		if recover() != nil {
			name = "?column?"
		}
	}()
	name, err := tree.GetRenderColName(context.Background(), tree.EmptySearchPath, e, functionNames{})
	if err != nil {
		return "?column?"
	}
	return name
}

// functionNames resolves a function call to its own name, which is all
// GetRenderColName reads. Resolving stores the result in the call, after
// the analyzer has already looked for the placeholder there.
type functionNames struct{ tree.FunctionReferenceResolver }

func (functionNames) ResolveFunction(_ context.Context, name tree.UnresolvedRoutineName, _ tree.SearchPath) (*tree.ResolvedFunctionDefinition, error) {
	return &tree.ResolvedFunctionDefinition{Name: name.UnresolvedName().Parts[0]}, nil
}

// sourceColumns returns the columns of src; hidden columns are included only
// when withHidden is set.
func (a *analyzer) sourceColumns(src source, withHidden bool) []string {
	if src.derived {
		return src.columns
	}
	rel := a.cat.relation(src.schema, src.name)
	if rel == nil {
		return nil
	}
	if !withHidden {
		return rel.visibleColumns()
	}
	out := make([]string, 0, len(rel.columns))
	for _, col := range rel.columns {
		out = append(out, col.name)
	}
	return out
}

// definitionColumns names the output columns of a view definition.
func definitionColumns(def string, cat *catalog) []string {
	def = strings.TrimSpace(def)
	if def == "" {
		return nil
	}
	stmt, err := parseOne(def)
	if err != nil {
		// A synced definition the server accepted but this parser rejects
		// leaves the view without columns rather than failing completion.
		return nil
	}
	a := &analyzer{ctx: context.Background(), cat: cat}
	if view, ok := stmt.(*tree.CreateView); ok {
		return a.query(view.AsSource, nil)
	}
	return a.statement(stmt, nil)
}
