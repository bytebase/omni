package completion

import (
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
)

// source is a relation a name can reference: a catalog relation or a derived
// one (CTE, subquery) with its own columns.
type source struct {
	ref     string // the alias, or the relation's own name
	schema  string // the catalog relation's schema; "" resolves by search path
	name    string // the catalog relation's name; "" for a derived source
	columns []string
	derived bool
}

// nameContext says what the word at the caret names and which sources and
// CTEs are in scope there.
type nameContext struct {
	kind      nameKind
	qualifier []string // outermost first: [schema] or [table] or [schema, table]
	scope     []source // innermost first
	target    *source
	ctes      []string
}

// analyzer finds placeholder in a statement and records its context.
type analyzer struct {
	cat     *catalog
	ctes    [][]source // a stack of WITH clauses, innermost last
	allCTEs []string   // every CTE the statement defines
	found   *nameContext
}

// analyzeStatement returns the context of placeholder in stmt, or nil when
// stmt does not hold it as a name completion can resolve.
func analyzeStatement(stmt tree.Statement, cat *catalog) *nameContext {
	a := &analyzer{cat: cat}
	a.statement(stmt, nil)
	if a.found == nil {
		a.found = tableNameContext(stmt)
	}
	if a.found == nil {
		a.found = genericExprContext(stmt)
	}
	if a.found != nil && a.found.kind == kindObject {
		a.found.ctes = a.allCTEs
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
	}
	return kindNone, nil, false
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

func (a *analyzer) statement(stmt tree.Statement, outer []source) {
	if a.found != nil || stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *tree.Select:
		a.selectQuery(s, outer)
	case *tree.ParenSelect:
		a.selectQuery(s.Select, outer)
	case *tree.Insert:
		a.insert(s, outer)
	case *tree.Update:
		a.update(s, outer)
	case *tree.Delete:
		a.delete(s, outer)
	case *tree.Explain:
		a.statement(s.Statement, outer)
	case *tree.ExplainAnalyze:
		a.statement(s.Statement, outer)
	case *tree.CreateView:
		a.selectQuery(s.AsSource, outer)
	case *tree.CreateTable:
		if s.AsSource != nil {
			a.selectQuery(s.AsSource, outer)
		}
	case *tree.CreateIndex:
		a.createIndex(s)
	case *tree.AlterTable:
		a.alterTable(s)
	}
}

func (a *analyzer) selectQuery(sel *tree.Select, outer []source) []source {
	if sel == nil || a.found != nil {
		return outer
	}
	popped := a.with(sel.With, outer)
	defer popped()
	scope := a.selectStatement(sel.Select, outer)
	for _, order := range sel.OrderBy {
		a.expr(order.Expr, scope)
	}
	if sel.Limit != nil {
		a.expr(sel.Limit.Count, outer)
		a.expr(sel.Limit.Offset, outer)
	}
	return scope
}

// with pushes the CTEs of w, walks their bodies, and returns the pop. A
// recursive CTE is visible inside its own body.
func (a *analyzer) with(w *tree.With, outer []source) func() {
	if w == nil {
		return func() {}
	}
	level := len(a.ctes)
	a.ctes = append(a.ctes, nil)
	for _, cte := range w.CTEList {
		src := source{ref: string(cte.Name.Alias), derived: true, columns: aliasColumns(cte.Name)}
		a.allCTEs = append(a.allCTEs, src.ref)
		if w.Recursive {
			a.ctes[level] = append(a.ctes[level], src)
		}
		a.statement(cte.Stmt, outer)
		if src.columns == nil {
			if sel, ok := cte.Stmt.(*tree.Select); ok {
				src.columns = a.outputColumns(sel)
			}
		}
		if w.Recursive {
			a.ctes[level][len(a.ctes[level])-1] = src
		} else {
			a.ctes[level] = append(a.ctes[level], src)
		}
	}
	return func() { a.ctes = a.ctes[:level] }
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

// selectStatement walks one select body and returns the scope its clauses
// see.
func (a *analyzer) selectStatement(ss tree.SelectStatement, outer []source) []source {
	if a.found != nil {
		return outer
	}
	switch s := ss.(type) {
	case *tree.SelectClause:
		local, conds := a.tableExprs(s.From.Tables, outer)
		scope := slices.Concat(local, outer)
		for _, cond := range conds {
			a.expr(cond, scope)
		}
		for _, e := range s.Exprs {
			a.expr(e.Expr, scope)
		}
		if s.Where != nil {
			a.expr(s.Where.Expr, scope)
		}
		for _, g := range s.GroupBy {
			a.expr(g, scope)
		}
		if s.Having != nil {
			a.expr(s.Having.Expr, scope)
		}
		for _, e := range s.DistinctOn {
			a.expr(e, scope)
		}
		return scope
	case *tree.UnionClause:
		left := a.selectQuery(s.Left, outer)
		a.selectQuery(s.Right, outer)
		return left
	case *tree.ParenSelect:
		return a.selectQuery(s.Select, outer)
	case *tree.ValuesClause:
		for _, row := range s.Rows {
			for _, e := range row {
				a.expr(e, outer)
			}
		}
	}
	return outer
}

// tableExprs returns the sources of a FROM list and the join conditions to
// check once every source is known.
func (a *analyzer) tableExprs(exprs tree.TableExprs, outer []source) ([]source, []tree.Expr) {
	var sources []source
	var conds []tree.Expr
	for _, te := range exprs {
		s, c := a.tableExpr(te, outer)
		sources = append(sources, s...)
		conds = append(conds, c...)
	}
	return sources, conds
}

func (a *analyzer) tableExpr(te tree.TableExpr, outer []source) ([]source, []tree.Expr) {
	switch t := te.(type) {
	case *tree.AliasedTableExpr:
		src, ok := a.aliased(t, outer)
		if !ok {
			return nil, nil
		}
		return []source{src}, nil
	case *tree.ParenTableExpr:
		return a.tableExpr(t.Expr, outer)
	case *tree.JoinTableExpr:
		left, lc := a.tableExpr(t.Left, outer)
		right, rc := a.tableExpr(t.Right, outer)
		conds := slices.Concat(lc, rc)
		if on, ok := t.Cond.(*tree.OnJoinCond); ok {
			conds = append(conds, on.Expr)
		}
		return slices.Concat(left, right), conds
	case *tree.TableName:
		return []source{a.tableSource(t, "")}, nil
	case *tree.Subquery:
		a.selectStatement(t.Select, outer)
	case *tree.StatementSource:
		a.statement(t.Statement, outer)
	}
	return nil, nil
}

func (a *analyzer) aliased(t *tree.AliasedTableExpr, outer []source) (source, bool) {
	alias := string(t.As.Alias)
	var src source
	switch e := t.Expr.(type) {
	case *tree.TableName:
		src = a.tableSource(e, alias)
		if t.IndexFlags != nil && string(t.IndexFlags.Index) == placeholder && a.found == nil {
			target := src
			a.found = &nameContext{kind: kindIndex, target: &target}
		}
	case *tree.Subquery:
		var inner []source
		if t.Lateral {
			inner = outer
		}
		a.selectStatement(e.Select, inner)
		src = source{ref: alias, derived: true, columns: a.selectColumns(e.Select)}
	case *tree.StatementSource:
		a.statement(e.Statement, outer)
		src = source{ref: alias, derived: true}
	case *tree.RowsFromExpr:
		for _, item := range e.Items {
			a.expr(item, outer)
		}
		src = source{ref: alias, derived: true}
	default:
		return source{}, false
	}
	if cols := aliasColumns(t.As); cols != nil {
		src.columns = cols
		src.derived = true
		src.name = ""
	}
	return src, true
}

// tableSource resolves a table name to a CTE in scope or a catalog relation.
func (a *analyzer) tableSource(tn *tree.TableName, alias string) source {
	name := string(tn.ObjectName)
	ref := alias
	if ref == "" {
		ref = name
	}
	if !tn.ExplicitSchema {
		for i := len(a.ctes) - 1; i >= 0; i-- {
			for _, cte := range a.ctes[i] {
				if cte.ref == name {
					return source{ref: ref, derived: true, columns: cte.columns}
				}
			}
		}
	}
	src := source{ref: ref, name: name}
	if tn.ExplicitSchema {
		src.schema = string(tn.SchemaName)
	}
	return src
}

// expr looks for placeholder in e, and walks e's subqueries with scope as
// their outer scope.
func (a *analyzer) expr(e tree.Expr, scope []source) {
	if e == nil || a.found != nil {
		return
	}
	_, _ = tree.SimpleVisit(e, func(expr tree.Expr) (bool, tree.Expr, error) {
		if a.found != nil {
			return false, expr, nil
		}
		if kind, qual, ok := placeholderInExpr(expr); ok {
			a.found = &nameContext{kind: kind, qualifier: qual, scope: scope}
			return false, expr, nil
		}
		if sub, ok := expr.(*tree.Subquery); ok {
			a.selectStatement(sub.Select, scope)
			return false, expr, nil
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
	if a.found != nil || target == nil {
		return
	}
	if slices.Contains(names, tree.Name(placeholder)) {
		a.found = &nameContext{kind: kindTargetColumn, target: target}
	}
}

func (a *analyzer) insert(s *tree.Insert, outer []source) {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	a.targetColumn(s.Columns, target)
	a.selectQuery(s.Rows, outer)
	scope := outer
	if target != nil {
		scope = append([]source{*target}, outer...)
	}
	if oc := s.OnConflict; oc != nil {
		a.targetColumn(oc.Columns, target)
		conflictScope := scope
		if target != nil {
			excluded := *target
			excluded.ref = "excluded"
			conflictScope = append([]source{*target, excluded}, outer...)
		}
		a.updateExprs(oc.Exprs, target, conflictScope)
		if oc.Where != nil {
			a.expr(oc.Where.Expr, conflictScope)
		}
	}
	a.returning(s.Returning, scope)
}

func (a *analyzer) update(s *tree.Update, outer []source) {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	from, conds := a.tableExprs(s.From, outer)
	scope := slices.Concat(from, outer)
	if target != nil {
		scope = append([]source{*target}, scope...)
	}
	for _, cond := range conds {
		a.expr(cond, scope)
	}
	a.updateExprs(s.Exprs, target, scope)
	if s.Where != nil {
		a.expr(s.Where.Expr, scope)
	}
	for _, order := range s.OrderBy {
		a.expr(order.Expr, scope)
	}
	a.returning(s.Returning, scope)
}

func (a *analyzer) delete(s *tree.Delete, outer []source) {
	popped := a.with(s.With, outer)
	defer popped()
	target := a.target(s.Table)
	using, conds := a.tableExprs(s.Using, outer)
	scope := slices.Concat(using, outer)
	if target != nil {
		scope = append([]source{*target}, scope...)
	}
	for _, cond := range conds {
		a.expr(cond, scope)
	}
	if s.Where != nil {
		a.expr(s.Where.Expr, scope)
	}
	for _, order := range s.OrderBy {
		a.expr(order.Expr, scope)
	}
	a.returning(s.Returning, scope)
}

func (a *analyzer) updateExprs(exprs tree.UpdateExprs, target *source, scope []source) {
	for _, ue := range exprs {
		a.targetColumn(ue.Names, target)
		a.expr(ue.Expr, scope)
	}
}

func (a *analyzer) returning(r tree.ReturningClause, scope []source) {
	if exprs, ok := r.(*tree.ReturningExprs); ok {
		for _, e := range *exprs {
			a.expr(e.Expr, scope)
		}
	}
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
}

// selectColumns names the output columns of a select body.
func (a *analyzer) selectColumns(ss tree.SelectStatement) []string {
	switch s := ss.(type) {
	case *tree.SelectClause:
		var out []string
		for _, e := range s.Exprs {
			switch {
			case e.As != "":
				out = append(out, string(e.As))
			default:
				out = append(out, a.exprColumns(e.Expr, s.From.Tables)...)
			}
		}
		return out
	case *tree.UnionClause:
		return a.outputColumns(s.Left)
	case *tree.ParenSelect:
		return a.outputColumns(s.Select)
	case *tree.ValuesClause:
		if len(s.Rows) == 0 {
			return nil
		}
		out := make([]string, len(s.Rows[0]))
		for i := range out {
			out[i] = "column" + strconv.Itoa(i+1)
		}
		return out
	}
	return nil
}

func (a *analyzer) outputColumns(sel *tree.Select) []string {
	if sel == nil {
		return nil
	}
	popped := a.with(sel.With, nil)
	defer popped()
	return a.selectColumns(sel.Select)
}

// exprColumns names the output columns one select expression yields: a
// column reference keeps its name, and a star expands to the visible columns
// of the sources it covers.
func (a *analyzer) exprColumns(e tree.Expr, from tree.TableExprs) []string {
	switch n := e.(type) {
	case *tree.UnresolvedName:
		if !n.Star {
			return []string{n.Parts[0]}
		}
		var qual string
		if n.NumParts > 1 {
			qual = n.Parts[1]
		}
		sources, _ := a.tableExprs(from, nil)
		var out []string
		for _, src := range sources {
			if qual == "" || src.ref == qual {
				out = append(out, a.sourceColumns(src, false)...)
			}
		}
		return out
	case tree.UnqualifiedStar:
		return a.exprColumns(&tree.UnresolvedName{NumParts: 1, Star: true}, from)
	}
	return nil
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
	// found is preset so the walk does not look for the placeholder.
	a := &analyzer{cat: cat, found: &nameContext{}}
	switch s := stmt.(type) {
	case *tree.Select:
		return a.outputColumns(s)
	case *tree.CreateView:
		return a.outputColumns(s.AsSource)
	}
	return nil
}
