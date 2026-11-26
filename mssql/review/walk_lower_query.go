package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// A view, an inline table-valued function, and each statement of a
// procedure, function, or trigger body are bound when the module is
// created: SQL Server resolves their tables and columns against the
// database as it is then (verified). The checks here are the ones that
// certainly fail on the ledger; a reference the walk cannot resolve says
// nothing.

// qsource is one table source of a query block.
type qsource struct {
	// exposed is the name a column's qualifier matches: the alias, or
	// the table's own name; "" for a source whose name the walk does not
	// know, which any qualifier may mean.
	exposed string
	written string
	// path is the table's path when the ledger lists its columns.
	path string
	// local is set for a table variable, whose columns are cols.
	local bool
	cols  []string
}

// qscope is a query block's sources, inside the blocks it is nested in.
type qscope struct {
	sources []qsource
	// opaque is set when a source's columns are not the ledger's: a view,
	// a derived table, a CTE, a function, or one the walk cannot place.
	opaque bool
	outer  *qscope
	// aliases are the block's select-list names, which its ORDER BY may
	// name instead of a column.
	aliases []string
	// order is set while the block's ORDER BY is bound.
	order bool
}

// binder binds the queries of one module or statement on the ledger.
type binder struct {
	b *builder
	// ctes are every CTE the statement defines, which a table reference
	// after it may mean.
	ctes []*ast.CommonTableExpr
	// visible are the CTEs of the WITH clauses around the query being
	// bound.
	visible []*ast.CommonTableExpr
	// noFrom reports a column with no source at all as missing: in a
	// view, an inline function, or a body's SELECT (verified).
	noFrom bool
	// pseudo is the table path a DML trigger's inserted and deleted
	// stand for, and pseudoTable the trigger's table as written;
	// pseudoCols are the columns of a view's INSTEAD OF trigger's.
	pseudo, pseudoTable string
	pseudoCols          []string
	// tableVars are the table variables the body declares.
	tableVars []*ast.VariableDecl
	// bound is set while a statement SQL Server certainly binds is bound:
	// every table it names exists, so a COLLATE in it is checked (448).
	bound bool
	// vars are the declared types of the module's parameters and
	// variables, for a COLLATE operand.
	vars []typedVar
	// types is set for a view's or an inline function's query, whose
	// operand types SQL Server checks when it creates the module (447,
	// 8116, 8117; measured).
	types bool
}

func newBinder(b *builder, root ast.Node, noFrom bool) *binder {
	q := &binder{b: b, noFrom: noFrom}
	ast.Inspect(root, func(n ast.Node) bool {
		if w, ok := n.(*ast.WithClause); ok && w.CTEs != nil {
			for _, item := range w.CTEs.Items {
				if cte, ok := item.(*ast.CommonTableExpr); ok {
					q.ctes = append(q.ctes, cte)
				}
			}
		}
		return true
	})
	return q
}

// isCTE reports whether an unqualified table name may mean a CTE: one
// defined before it, or itself, as a recursive CTE's query names it. A
// CTE defined after the name is not one it can mean (208, verified).
func (q *binder) isCTE(tr *ast.TableRef) bool {
	if tr.Schema != "" || tr.Database != "" || tr.Server != "" {
		return false
	}
	for _, c := range q.ctes {
		if !differ(c.Name, tr.Object) && (c.Loc.Start < 0 || tr.Loc.Start < 0 || c.Loc.Start < tr.Loc.Start) {
			return true
		}
	}
	return false
}

// query binds a SELECT and its set operation branches, CTEs, and
// subqueries.
func (q *binder) query(n ast.Node, outer *qscope) {
	s, ok := n.(*ast.SelectStmt)
	if !ok || s == nil {
		return
	}
	if s.WithClause != nil && s.WithClause.CTEs != nil {
		saved := q.visible
		q.visible = append([]*ast.CommonTableExpr{}, saved...)
		for _, item := range s.WithClause.CTEs.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				q.visible = append(q.visible, cte)
			}
		}
		defer func() { q.visible = saved }()
		for _, item := range s.WithClause.CTEs.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				q.query(cte.Query, outer)
			}
		}
	}
	if s.Larg != nil || s.Rarg != nil {
		q.query(s.Larg, outer)
		q.query(s.Rarg, outer)
		return
	}
	sc := &qscope{outer: outer}
	var conds []ast.ExprNode
	if s.FromClause != nil {
		for _, item := range s.FromClause.Items {
			q.source(item, sc, &conds, false)
		}
	}
	for _, c := range conds {
		q.expr(c, sc)
	}
	if s.TargetList != nil {
		for _, item := range s.TargetList.Items {
			q.expr(item, sc)
			if rt, ok := item.(*ast.ResTarget); ok {
				if rt.Name != "" {
					sc.aliases = append(sc.aliases, rt.Name)
				} else if c, ok := rt.Val.(*ast.ColumnRef); ok {
					sc.aliases = append(sc.aliases, c.Column)
				}
			}
		}
	}
	q.expr(s.WhereClause, sc)
	if outer == nil && q.ungrouped(s) {
		return
	}
	if s.GroupByClause != nil {
		for _, item := range s.GroupByClause.Items {
			q.expr(item, sc)
		}
	}
	q.expr(s.HavingClause, sc)
	if s.OrderByClause != nil {
		sc.order = true
		for _, item := range s.OrderByClause.Items {
			q.expr(item, sc)
		}
		sc.order = false
	}
}

// source adds a FROM item's sources to a block, binding the derived
// tables in it, and collects the join conditions to bind once every
// source is in. apply is set for the right side of APPLY.
func (q *binder) source(n ast.Node, sc *qscope, conds *[]ast.ExprNode, apply bool) {
	switch v := n.(type) {
	case *ast.TableRef:
		q.table(v, v.Alias, sc)
	case *ast.AliasedTableRef:
		switch t := v.Table.(type) {
		case *ast.TableRef:
			if v.Columns != nil {
				sc.opaque = true
				sc.sources = append(sc.sources, qsource{exposed: v.Alias})
				return
			}
			alias := v.Alias
			if alias == "" {
				alias = t.Alias
			}
			q.table(t, alias, sc)
		case *ast.SubqueryExpr:
			// A derived table may refer to the block's other sources
			// under APPLY, and else only to the blocks around it
			// (4104, verified).
			outer := sc
			if !apply {
				outer = sc.outer
			}
			q.query(t.Query, outer)
			sel, _ := t.Query.(*ast.SelectStmt)
			q.derived(v.Alias, sel, v.Columns, sc)
		case *ast.ValuesClause:
			names, ok := listNames(v.Columns)
			if !ok || v.Columns == nil {
				sc.opaque = true
				sc.sources = append(sc.sources, qsource{exposed: v.Alias})
				return
			}
			sc.sources = append(sc.sources, qsource{exposed: v.Alias, written: v.Alias, local: true, cols: names})
		case *ast.FuncCallExpr:
			if v.Columns == nil {
				// A column list fails first on the call (317).
				q.b.tableFunction(t)
				if fn := q.b.syncedFunction(t.Name); fn != nil && fn.kind != funcScalar && fn.cols != nil {
					src := qsource{exposed: v.Alias, written: v.Alias, local: true, cols: fn.cols}
					if v.Alias == "" {
						src.exposed, src.written = t.Name.Object, object(t.Name)
					}
					sc.sources = append(sc.sources, src)
					return
				}
			}
			sc.opaque = true
			sc.sources = append(sc.sources, qsource{exposed: v.Alias})
		default:
			sc.opaque = true
			sc.sources = append(sc.sources, qsource{exposed: v.Alias})
		}
	case *ast.TableVarRef:
		exposed := v.Alias
		if exposed == "" {
			exposed = v.Name
		}
		d := q.tableVar(v.Name)
		if d == nil || d.TableDef == nil {
			sc.opaque = true
			sc.sources = append(sc.sources, qsource{exposed: exposed})
			return
		}
		src := qsource{exposed: exposed, written: v.Name, local: true}
		for _, item := range d.TableDef.Items {
			if cd, ok := item.(*ast.ColumnDef); ok {
				src.cols = append(src.cols, cd.Name)
			}
		}
		sc.sources = append(sc.sources, src)
	case *ast.PivotExpr:
		inner := &qscope{outer: sc.outer}
		q.source(v.Source, inner, conds, false)
		q.pivot(v, inner, sc)
	case *ast.UnpivotExpr:
		inner := &qscope{outer: sc.outer}
		q.source(v.Source, inner, conds, false)
		q.unpivot(v, inner, sc)
	case *ast.JoinClause:
		q.source(v.Left, sc, conds, apply)
		// A join under APPLY's right side still sees what APPLY's left
		// side holds.
		q.source(v.Right, sc, conds, apply || v.Type == ast.JoinCrossApply || v.Type == ast.JoinOuterApply)
		if v.Condition != nil {
			*conds = append(*conds, v.Condition)
		}
		if v.Using != nil {
			sc.opaque = true
		}
	default:
		sc.opaque = true
		sc.sources = append(sc.sources, qsource{})
	}
}

// table adds a named table source: one whose columns the ledger lists
// when it is a table the walk knows, else an opaque one.
func (q *binder) table(tr *ast.TableRef, alias string, sc *qscope) {
	exposed := alias
	if exposed == "" {
		exposed = tr.Object
	}
	src := qsource{exposed: exposed, written: object(tr)}
	if q.isPseudo(tr) {
		src.path, src.local, src.cols = q.pseudo, q.pseudoCols != nil, q.pseudoCols
		sc.sources = append(sc.sources, src)
		return
	}
	if q.isCTE(tr) && !tr.ForPath && tr.Temporal == nil {
		if cte := q.cte(tr); cte != nil {
			sel, _ := cte.Query.(*ast.SelectStmt)
			if names, ok := cteNames(sel, cte.Columns); ok {
				src.local, src.cols = true, names
				sc.sources = append(sc.sources, src)
				return
			}
		}
	}
	if tr.Temporal != nil && !q.isCTE(tr) && !strings.HasPrefix(tr.Object, "#") {
		// FOR SYSTEM_TIME reads a system-versioned table (13544).
		if r, ok := q.b.ref(tr, false); ok && q.b.periodless(r, rangeOf(tr.Loc)) {
			sc.opaque = true
			sc.sources = append(sc.sources, src)
			return
		}
	}
	if q.isCTE(tr) || strings.HasPrefix(tr.Object, "#") || tr.Temporal != nil || tr.ForPath {
		sc.opaque = true
		sc.sources = append(sc.sources, src)
		return
	}
	r, ok := q.b.ref(tr, false)
	if ok && r.where == placed && !r.temp && hasHint(tr.Hints, "NOEXPAND") {
		// NOEXPAND reads an indexed view (8171).
		present := notKind(r.written, "view", kindView, kindSynonym)
		q.b.ofKind(spaceRelation, r.parent, r.n, []string{kindView, kindSynonym}, note(r.rng, "", present))
	}
	if ok {
		if cols := q.b.viewColumns(r); cols != nil {
			src.local, src.cols, src.written = true, cols, r.written
			sc.sources = append(sc.sources, src)
			return
		}
	}
	if !ok || !q.b.knownTable(r) {
		sc.opaque = true
		sc.sources = append(sc.sources, src)
		return
	}
	src.path = r.tablePath()
	src.written = r.written
	sc.sources = append(sc.sources, src)
}

// expr binds the columns of an expression and the queries nested in it.
func (q *binder) expr(n ast.Node, sc *qscope) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(m ast.Node) bool {
		switch v := m.(type) {
		case *ast.SelectStmt:
			q.query(v, sc)
			return false
		case *ast.ColumnRef:
			q.column(v, sc)
			return false
		case *ast.CollateExpr:
			// The operand is bound before the operator.
			q.expr(v.Expr, sc)
			if q.bound {
				q.b.collateFail(v, func(e ast.ExprNode) string { return q.operandType(e, sc) }, q.types)
			}
			return false
		case *ast.FuncCallExpr:
			if q.bound && q.types {
				q.callTypes(v, sc)
			}
			if q.call(v) {
				return false
			}
			if v.Name != nil && v.Name.Schema == "" && strings.EqualFold(v.Name.Object, "UPDATE") && v.Args != nil && len(v.Args.Items) == 1 {
				// A DML trigger's UPDATE(column) names a column of its
				// table, whatever the query around it reads.
				if c, ok := v.Args.Items[0].(*ast.ColumnRef); ok {
					if q.pseudo != "" || q.pseudoCols != nil {
						src := qsource{exposed: q.pseudoTable, written: q.pseudoTable, path: q.pseudo, local: q.pseudoCols != nil, cols: q.pseudoCols}
						q.column(c, &qscope{sources: []qsource{src}})
					}
					return false
				}
			}
		}
		return true
	})
}

// column checks that a column reference names a column of a source. An
// unqualified name looks in the block's sources and then in those of the
// blocks around it; a qualified one in the source its qualifier names.
func (q *binder) column(c *ast.ColumnRef, sc *qscope) {
	if c.Schema != "" || c.Database != "" || c.Server != "" {
		return
	}
	if c.Table != "" {
		q.qualified(c, sc)
		return
	}
	if sc != nil && sc.order {
		for _, a := range sc.aliases {
			if !differ(a, c.Column) {
				return
			}
		}
	}
	if !q.b.isDelimitedAt(c.Loc) {
		switch {
		case strings.EqualFold(c.Column, "$IDENTITY") && q.pseudoColumn(c, sc, func(st *snapTable) bool { return len(st.identity) == 0 }):
			return
		case strings.EqualFold(c.Column, "$ROWGUID") && q.pseudoColumn(c, sc, func(st *snapTable) bool { return st.noRowguid }):
			return
		}
	}
	cn, ok := keyed(c.Column)
	if !ok {
		return
	}
	pseudo := isPseudoCol(c.Column)
	var first *sight
	var firstSrc qsource
	any := false
	for s := sc; s != nil; s = s.outer {
		if s.opaque {
			return
		}
		if !pseudo && q.ambiguous(c, cn, s) {
			return
		}
		for _, src := range s.sources {
			if pseudo {
				return
			}
			if src.local {
				if hasName(src.cols, c.Column) {
					return
				}
				if !any {
					firstSrc = src
				}
				any = true
				continue
			}
			got := q.b.see(spaceColumn, src.path, cn)
			if got.fact != walk.Absent {
				return
			}
			if first == nil {
				first = &got
				if !any {
					firstSrc = src
				}
			}
			any = true
		}
	}
	rng := rangeOf(c.Loc)
	switch {
	case !any:
		if q.noFrom {
			q.b.fail(rng, review.ObjectNotFound, msgMissing("column", ident(c.Column)))
		}
	case firstSrc.local || first == nil:
		q.b.fail(rng, review.ObjectNotFound, msgMissingIn("column", ident(c.Column), firstSrc.written))
	default:
		q.b.Need(first.key, walk.Present, nil, note(rng, msgMissingIn("column", ident(c.Column), firstSrc.written), nil))
	}
}

// ambiguous fails a column that two of a block's tables certainly hold
// (209), and reports whether it did.
func (q *binder) ambiguous(c *ast.ColumnRef, cn name, s *qscope) bool {
	var keys []walk.Key
	for _, src := range s.sources {
		if src.local {
			if hasName(src.cols, c.Column) {
				// A table variable's column may be spelled otherwise.
				return false
			}
			continue
		}
		if got := q.b.see(spaceColumn, src.path, cn); got.fact == walk.Present {
			keys = append(keys, got.key)
		}
	}
	if len(keys) < 2 {
		return false
	}
	msg := "Column " + review.Code(ident(c.Column)) + " is in more than one table of the query."
	q.b.Fail(keys[:2], walk.Present, &walkNote{rng: rangeOf(c.Loc), absent: msg, present: text(msg), fixed: review.UnsupportedForm})
	return true
}

// qualified checks a column whose qualifier names one source the walk
// knows the columns of.
func (q *binder) qualified(c *ast.ColumnRef, sc *qscope) {
	if isPseudoCol(c.Column) {
		// A view's query with no source at all has no table to qualify
		// a column with (107).
		if q.noFrom && noSources(sc) {
			q.unbound(c, sc)
		}
		return
	}
	for s := sc; s != nil; s = s.outer {
		var match *qsource
		for i, src := range s.sources {
			if src.exposed == "" {
				return
			}
			if differ(src.exposed, c.Table) {
				continue
			}
			if match != nil {
				return
			}
			match = &s.sources[i]
		}
		if match == nil {
			continue
		}
		if match.local {
			if !hasName(match.cols, c.Column) {
				q.b.fail(rangeOf(c.Loc), review.ObjectNotFound, msgMissingIn("column", ident(c.Column), match.written))
			}
			return
		}
		if match.path == "" {
			return
		}
		cn, ok := keyed(c.Column)
		if !ok {
			return
		}
		got := q.b.see(spaceColumn, match.path, cn)
		if got.fact == walk.Absent {
			q.b.Need(got.key, walk.Present, nil, note(rangeOf(c.Loc), msgMissingIn("column", ident(c.Column), match.written), nil))
		}
		return
	}
	q.unbound(c, sc)
}

// unbound fails a column whose qualifier names no source of the query
// (4104, 107): none is named so, and none has a column so named, whose
// property the reference may be instead.
func (q *binder) unbound(c *ast.ColumnRef, sc *qscope) {
	qn, ok := keyed(c.Table)
	if !ok {
		return
	}
	var keys []walk.Key
	any := false
	for s := sc; s != nil; s = s.outer {
		if s.opaque {
			return
		}
		for _, src := range s.sources {
			any = true
			if src.local {
				if hasName(src.cols, c.Table) {
					return
				}
				continue
			}
			got := q.b.see(spaceColumn, src.path, qn)
			if t := q.b.pathColType(src.path, c.Table); got.fact == walk.Present && isBuiltinSystemType(t) && !clrTypes[t] {
				// A column of a system type that is no CLR type has
				// no property to read (4104, verified).
				continue
			}
			if got.fact != walk.Absent {
				return
			}
			keys = append(keys, got.key)
		}
	}
	if !any && !q.noFrom {
		return
	}
	msg := review.Code(ident(c.Table)) + " is not a table or alias of the query."
	q.b.Fail(keys, walk.Absent, &walkNote{rng: rangeOf(c.Loc), absent: msg, present: text(msg), fixed: review.ObjectNotFound})
}

// unnamedDerived finds a derived table or a CTE without a column list
// whose query leaves a column without a name, used or not (8155,
// verified).
func unnamedDerived(n ast.Node) (review.Range, bool) {
	var rng review.Range
	found := false
	unnamed := func(q ast.Node) bool {
		s, ok := q.(*ast.SelectStmt)
		if !ok || s == nil {
			return false
		}
		for s.Larg != nil {
			s = s.Larg
		}
		if s.TargetList == nil || s.ForClause != nil {
			return false
		}
		for _, item := range s.TargetList.Items {
			if rt, ok := item.(*ast.ResTarget); ok && rt.Name == "" && unnamedExpr(rt.Val) {
				if _, assign := rt.Val.(*ast.SelectAssign); !assign {
					return true
				}
			}
		}
		return false
	}
	ast.Inspect(n, func(m ast.Node) bool {
		if found {
			return false
		}
		switch v := m.(type) {
		case *ast.AliasedTableRef:
			if v.Columns != nil {
				break
			}
			switch t := v.Table.(type) {
			case *ast.SubqueryExpr:
				found = unnamed(t.Query)
			case *ast.ValuesClause:
				// VALUES names none of its columns.
				found = true
			}
			if found {
				rng = rangeOf(v.Loc)
			}
		case *ast.CommonTableExpr:
			if v.Columns == nil && unnamed(v.Query) {
				rng, found = rangeOf(v.Loc), true
			}
		}
		return !found
	})
	return rng, found
}

// moduleSchema returns the schema a module statement's name resolves
// to, as createModule resolves it: the default schema for a creation,
// and for ALTER the one a lookup finds; "" when the walk cannot tell.
func (b *builder) moduleSchema(name *ast.TableRef, loc ast.Loc, orAlter bool) string {
	r, ok := b.ref(name, b.moduleVerb(loc, orAlter) != "ALTER")
	if !ok || r.where != placed || r.temp {
		return ""
	}
	return r.schema.spelling
}

// enterModule makes unqualified names resolve as they do in the body of
// a module of a schema, nowhere for "", until the returned function runs.
func (b *builder) enterModule(schema string) func() {
	saved, savedIn := b.modSchema, b.inModule
	b.inModule, b.modSchema = true, schema
	return func() { b.modSchema, b.inModule = saved, savedIn }
}

// viewQuery checks what SQL Server checks of a view's or an inline
// function's query when it creates the module, which nothing defers: its
// variables (137, 1087), unless params declares them; its temporary
// tables (4508); its tables (208); its column names (4511, 4506, 8156,
// 8158, 8159); SELECT * without FROM anywhere in it (263); and its
// columns (207).
func (b *builder) viewQuery(query ast.Node, cols *ast.List, kind string, params []string, bound bool) {
	if query == nil {
		return
	}
	if b.undeclared(query, params) {
		return
	}
	if rng, ok := tempTableIn(query); ok {
		b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept a temporary table in a "+kind+".")
		return
	}
	if b.viewForms(query, kind) {
		return
	}
	if rng, bad := anySelectStarWithoutFrom(query); bad {
		b.fail(rng, review.UnsupportedForm, "`SELECT *` requires a `FROM` clause.")
		return
	}
	q := newBinder(b, query, true)
	if bound && b.schemaBound(q, query, kind, true) || b.setWidths(query) {
		return
	}
	if rng, bad := orderWithoutTop(query); bad {
		b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept `ORDER BY` without `TOP`, `OFFSET`, or `FOR XML` in a "+kind+".")
		return
	}
	if b.columnNames(query, cols, kind) {
		return
	}
	if rng, bad := unnamedDerived(query); bad {
		b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept a derived table or CTE column without a name.")
		return
	}
	if b.queryForms(query) {
		return
	}
	b.queryTables(q, query)
	q.bound = q.resolvable(query)
	q.types = true
	q.query(query, nil)
}

// undeclared fails a module whose body names a variable before declaring
// it: a parameter, or a DECLARE before it in the text (137, 1087). A
// system variable (@@) is always there.
func (b *builder) undeclared(body ast.Node, params []string) bool {
	type decl struct {
		name string
		at   int
	}
	var decls []decl
	for _, p := range params {
		decls = append(decls, decl{p, -1})
	}
	type use struct {
		name  string
		loc   ast.Loc
		table bool
	}
	var uses []use
	ref := func(name string, loc ast.Loc, table bool) {
		if strings.HasPrefix(name, "@") && !strings.HasPrefix(name, "@@") {
			uses = append(uses, use{name, loc, table})
		}
	}
	ast.Inspect(body, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.VariableDecl:
			// The variable is declared from the end of its declaration.
			decls = append(decls, decl{v.Name, v.Loc.End})
		case *ast.VariableRef:
			ref(v.Name, v.Loc, false)
		case *ast.TableVarRef:
			ref(v.Name, v.Loc, true)
		case *ast.TableVarMethodCallRef:
			ref(v.Var, v.Loc, false)
		case *ast.SelectAssign:
			ref(v.Variable, v.Loc, false)
		case *ast.SetStmt:
			ref(v.Variable, v.Loc, false)
		case *ast.SetExpr:
			ref(v.Variable, v.Loc, false)
		}
		return true
	})
	for _, u := range uses {
		if u.loc.Start < 0 {
			return false
		}
	}
	for _, d := range decls {
		if d.at == 0 {
			// A declaration whose place the walk does not know.
			return false
		}
	}
	var bad *use
	for i, u := range uses {
		ok := false
		for _, d := range decls {
			if d.at <= u.loc.Start && !differ(strings.TrimPrefix(d.name, "@"), strings.TrimPrefix(u.name, "@")) {
				ok = true
				break
			}
		}
		if !ok && (bad == nil || u.loc.Start < bad.loc.Start) {
			bad = &uses[i]
		}
	}
	if bad == nil {
		return false
	}
	if bad.table {
		b.fail(rangeOf(bad.loc), review.ObjectNotFound, "Table variable "+review.Code(bad.name)+" is not declared.")
	} else {
		b.fail(rangeOf(bad.loc), review.ObjectNotFound, "Variable "+review.Code(bad.name)+" is not declared.")
	}
	return true
}

// tempTableIn returns a temporary table a query reads.
func tempTableIn(n ast.Node) (review.Range, bool) {
	var rng review.Range
	found := false
	forEachFromTable(n, func(tr *ast.TableRef) {
		if !found && strings.HasPrefix(tr.Object, "#") && tr.Server == "" {
			rng, found = rangeOf(tr.Loc), true
		}
	})
	return rng, found
}

// forEachFromTable calls f for each named table of every FROM list in a
// query, at any depth.
func forEachFromTable(n ast.Node, f func(*ast.TableRef)) {
	var source func(ast.Node)
	source = func(m ast.Node) {
		switch v := m.(type) {
		case *ast.TableRef:
			f(v)
		case *ast.AliasedTableRef:
			if tr, ok := v.Table.(*ast.TableRef); ok {
				f(tr)
			}
		case *ast.JoinClause:
			source(v.Left)
			source(v.Right)
		}
	}
	ast.Inspect(n, func(m ast.Node) bool {
		if s, ok := m.(*ast.SelectStmt); ok && s.FromClause != nil {
			for _, item := range s.FromClause.Items {
				source(item)
			}
		}
		return true
	})
}

// queryTables checks that each table a view's query reads exists (208).
func (b *builder) queryTables(q *binder, n ast.Node) {
	forEachFromTable(n, func(tr *ast.TableRef) {
		// A system table's compatibility view, such as sysobjects, is
		// found without its schema.
		if q.isCTE(tr) || strings.HasPrefix(tr.Object, "#") || len(tr.Object) >= 3 && strings.EqualFold(tr.Object[:3], "sys") {
			return
		}
		r, ok := b.ref(tr, false)
		if !ok || r.where != placed || r.temp {
			return
		}
		if b.see(spaceRelation, r.parent, r.n).fact == walk.Absent {
			b.needs(spaceRelation, r.parent, r.n, nil, note(r.rng, msgMissing("", r.written), nil))
		}
	})
}

// anySelectStarWithoutFrom finds SELECT * without a FROM clause in any
// query block of a statement (263). EXISTS takes one such query as it
// is, but not a set operation of them (verified).
func anySelectStarWithoutFrom(n ast.Node) (review.Range, bool) {
	var rng review.Range
	found := false
	exempt := make(map[*ast.SelectStmt]bool)
	ast.Inspect(n, func(m ast.Node) bool {
		if found {
			return false
		}
		switch v := m.(type) {
		case *ast.ExistsExpr:
			if s, ok := v.Subquery.(*ast.SelectStmt); ok && s.Larg == nil && s.Rarg == nil {
				exempt[s] = true
			}
		case *ast.SelectStmt:
			if v.Larg == nil && v.Rarg == nil && !exempt[v] {
				if r, bad := selectStarWithoutFrom(v); bad {
					rng, found = r, true
					return false
				}
			}
		}
		return true
	})
	return rng, found
}

// columnNames checks the names a view or an inline function gives its
// columns: every column needs one (4511), and no two the same (4506);
// a column list names each once (8156) and as many columns as the query
// returns (8158, 8159). Two names are the same here only when spelled
// alike, which every collation equates.
func (b *builder) columnNames(query ast.Node, cols *ast.List, kind string) bool {
	s, ok := query.(*ast.SelectStmt)
	if !ok || s == nil {
		return false
	}
	for s.Larg != nil {
		s = s.Larg
	}
	if s.TargetList == nil || s.IntoTable != nil {
		return false
	}
	if s.ForClause != nil {
		// FOR XML and FOR JSON return one column, which has no name.
		if (s.ForClause.Mode == ast.ForXML || s.ForClause.Mode == ast.ForJSON) && cols == nil {
			b.fail(rangeOf(s.ForClause.Loc), review.UnsupportedForm, "SQL Server doesn't accept a "+kind+" column without a name.")
			return true
		}
		return false
	}
	var names []string
	star := false
	// A * stands for the columns of the sources it names, when the walk
	// knows them.
	sources, known := b.starSources(query)
	expand := func(st *ast.StarExpr) {
		star = true
		for _, src := range sources {
			if known && (st.Qualifier == "" || st.Qualifier == src.exposed) {
				names = append(names, src.cols...)
			}
		}
	}
	for _, item := range s.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			if st, ok := item.(*ast.StarExpr); ok {
				expand(st)
			}
			continue
		}
		switch v := rt.Val.(type) {
		case *ast.StarExpr:
			expand(v)
			continue
		case *ast.SelectAssign:
			return false
		case *ast.ColumnRef:
			if rt.Name == "" {
				names = append(names, v.Column)
				continue
			}
		}
		if rt.Name != "" {
			names = append(names, rt.Name)
			continue
		}
		if unnamedExpr(rt.Val) && cols == nil {
			b.fail(rangeOf(rt.Loc), review.UnsupportedForm, "SQL Server doesn't accept a "+kind+" column without a name.")
			return true
		}
		names = append(names, "")
	}
	if cols != nil {
		var list []string
		for _, item := range cols.Items {
			str, ok := item.(*ast.String)
			if !ok {
				return false
			}
			list = append(list, str.Str)
		}
		if n := twice(list); n != "" {
			b.fail(b.anchor(), review.DuplicateName, msgTwice("names", "column", ident(n)))
			return true
		}
		if !star && len(list) != len(s.TargetList.Items) {
			b.fail(b.anchor(), review.UnsupportedForm, "SQL Server doesn't accept a "+kind+" column list of "+count(len(list), "name")+" for a query of "+count(len(s.TargetList.Items), "column")+".")
			return true
		}
		return false
	}
	if n := twice(names); n != "" {
		b.fail(b.anchor(), review.DuplicateName, msgTwice("names", "column", ident(n)))
		return true
	}
	return false
}

// starSources returns the sources of a view's query, a single block, with
// the columns a * gives of each: false unless each is a synced table not
// changed since. A column that may be hidden from * (a period column, a
// ledger table's, a graph table's) is left out, so two names a * gives
// are the columns of two tables.
func (b *builder) starSources(query ast.Node) ([]qsource, bool) {
	s, ok := query.(*ast.SelectStmt)
	if !ok || s == nil || s.Larg != nil || s.Rarg != nil || s.WithClause != nil || s.FromClause == nil {
		return nil, false
	}
	var out []qsource
	known := true
	var add func(n ast.Node)
	add = func(n ast.Node) {
		var tr *ast.TableRef
		alias := ""
		switch v := n.(type) {
		case *ast.TableRef:
			tr, alias = v, v.Alias
		case *ast.AliasedTableRef:
			t, ok := v.Table.(*ast.TableRef)
			if !ok || v.Columns != nil {
				known = false
				return
			}
			tr, alias = t, v.Alias
			if alias == "" {
				alias = t.Alias
			}
		case *ast.JoinClause:
			if v.Using != nil {
				known = false
				return
			}
			add(v.Left)
			add(v.Right)
			return
		default:
			known = false
			return
		}
		if strings.HasPrefix(tr.Object, "#") || tr.Temporal != nil || tr.ForPath {
			known = false
			return
		}
		r, ok := b.ref(tr, false)
		if !ok {
			known = false
			return
		}
		st := b.syncedTable(r)
		if st == nil || st.ledger || st.edge {
			known = false
			return
		}
		src := qsource{exposed: alias}
		if alias == "" {
			src.exposed = tr.Object
		}
		for _, c := range st.names {
			if t := st.colTypes[c]; t != "" && t != "datetime2" && !strings.HasPrefix(c, "$") {
				src.cols = append(src.cols, c)
			}
		}
		out = append(out, src)
	}
	for _, item := range s.FromClause.Items {
		add(item)
	}
	return out, known
}

// unnamedExpr reports whether a select-list expression without an alias
// certainly has no name: anything but a column, or a column in
// parentheses or under COLLATE, which the walk leaves alone.
func unnamedExpr(e ast.ExprNode) bool {
	switch v := e.(type) {
	case *ast.ParenExpr:
		return v.Expr != nil && unnamedExpr(v.Expr)
	case *ast.ColumnRef, *ast.CollateExpr, *ast.StarExpr, nil:
		return false
	}
	return true
}

// twice returns a name spelled alike twice in a list, "" when none is.
func twice(names []string) string {
	seen := make(map[string]bool, len(names))
	for _, n := range names {
		if n == "" {
			continue
		}
		if seen[n] {
			return n
		}
		seen[n] = true
	}
	return ""
}

// moduleBody checks what SQL Server checks of a procedure's, a
// function's, or a trigger's body when it creates the module: each
// variable is declared before it is named (137, 1087), and each
// statement whose tables all exist binds its columns (207); one that
// names a table that may not exist, or a temporary table, is bound when
// it runs instead (deferred name resolution). A body that changes a
// table may bind against what it changed, and binds nothing here.
// pseudo is the path of the table whose rows a DML trigger's inserted
// and deleted hold, "" when the walk does not know them, and table that
// table as written; pseudoCols are a view's columns instead. vars are the
// module's typed parameters.
func (b *builder) moduleBody(body ast.Node, params []string, vars []typedVar, tableVars []*ast.VariableDecl, pseudo, table string, pseudoCols []string) {
	if body == nil || b.undeclared(body, params) {
		return
	}
	var stmts []ast.Node
	changes := flatBody(body, &stmts)
	for _, s := range stmts {
		if ins, ok := s.(*ast.InsertStmt); ok && b.insertCounts(ins) || b.setWidths(s) {
			return
		}
	}
	if changes {
		return
	}
	q := newBinder(b, body, false)
	q.pseudo, q.pseudoTable, q.pseudoCols = pseudo, table, pseudoCols
	q.tableVars = tableVars
	ast.Inspect(body, func(n ast.Node) bool {
		if v, ok := n.(*ast.VariableDecl); ok && v.IsTable {
			q.tableVars = append(q.tableVars, v)
		}
		return true
	})
	q.vars = vars
	ast.Inspect(body, func(n ast.Node) bool {
		if v, ok := n.(*ast.VariableDecl); ok && !v.IsTable {
			q.vars = append(q.vars, typedVar{v.Name, v.DataType})
		}
		return true
	})
	for _, s := range stmts {
		if q.resolvable(s) {
			q.bound = true
			q.statement(s)
			q.bound = false
		}
	}
}

// flatBody lists the statements of a body that bind on their own, and
// reports whether the body holds a statement that may change a table.
func flatBody(n ast.Node, out *[]ast.Node) (changes bool) {
	switch v := n.(type) {
	case nil:
	case *ast.BeginEndStmt:
		if v.Stmts != nil {
			for _, item := range v.Stmts.Items {
				changes = flatBody(item, out) || changes
			}
		}
	case *ast.TryCatchStmt:
		for _, l := range []*ast.List{v.TryBlock, v.CatchBlock} {
			if l == nil {
				continue
			}
			for _, item := range l.Items {
				changes = flatBody(item, out) || changes
			}
		}
	case *ast.IfStmt:
		*out = append(*out, v.Condition)
		changes = flatBody(v.Then, out)
		changes = flatBody(v.Else, out) || changes
	case *ast.WhileStmt:
		*out = append(*out, v.Condition)
		changes = flatBody(v.Body, out)
	case *ast.SelectStmt:
		*out = append(*out, n)
		changes = v.IntoTable != nil
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.DeclareStmt, *ast.SetStmt, *ast.ReturnStmt, *ast.DeclareCursorStmt, *ast.PrintStmt:
		*out = append(*out, n)
	case *ast.OpenCursorStmt, *ast.FetchCursorStmt, *ast.CloseCursorStmt, *ast.DeallocateCursorStmt,
		*ast.BreakStmt, *ast.ContinueStmt, *ast.RaiseErrorStmt, *ast.ThrowStmt,
		*ast.BeginTransStmt, *ast.CommitTransStmt, *ast.RollbackTransStmt, *ast.SaveTransStmt,
		*ast.SetOptionStmt, *ast.GotoStmt, *ast.LabelStmt, *ast.WaitForStmt:
	default:
		changes = true
	}
	return changes
}

// resolvable reports whether every table a statement names certainly
// exists, so SQL Server binds it when the module is created.
func (q *binder) resolvable(n ast.Node) bool {
	if n == nil {
		return false
	}
	ok := true
	item := func(m ast.Node) {
		switch v := m.(type) {
		case *ast.TableRef:
			ok = ok && q.exists(v)
		case *ast.AliasedTableRef:
			switch t := v.Table.(type) {
			case *ast.TableRef:
				ok = ok && q.exists(t)
			case *ast.SubqueryExpr:
			case *ast.FuncCallExpr:
				// A table-valued function the walk knows is bound too.
				fn := q.b.syncedFunction(t.Name)
				ok = ok && fn != nil && fn.kind != funcScalar
			default:
				ok = false
			}
		case *ast.TableVarRef, *ast.JoinClause, nil:
		default:
			ok = false
		}
	}
	var from func(m ast.Node)
	from = func(m ast.Node) {
		if j, isJoin := m.(*ast.JoinClause); isJoin {
			from(j.Left)
			from(j.Right)
			return
		}
		item(m)
	}
	list := func(l *ast.List) {
		if l != nil {
			for _, it := range l.Items {
				from(it)
			}
		}
	}
	ast.Inspect(n, func(m ast.Node) bool {
		switch v := m.(type) {
		case *ast.SelectStmt:
			list(v.FromClause)
			if v.IntoTable != nil {
				ok = false
			}
		case *ast.InsertStmt:
			item(v.Relation)
		case *ast.UpdateStmt:
			item(v.Relation)
			list(v.FromClause)
		case *ast.DeleteStmt:
			item(v.Relation)
			list(v.FromClause)
		case *ast.MergeStmt:
			from(v.Target)
			from(v.Source)
		case *ast.ExecStmt:
			ok = false
		}
		return ok
	})
	return ok
}

// exists reports whether a named table certainly exists: a CTE, a DML
// trigger's inserted or deleted, or a placed relation of the ledger.
func (q *binder) exists(tr *ast.TableRef) bool {
	if q.isCTE(tr) || q.isPseudo(tr) {
		return true
	}
	if strings.HasPrefix(tr.Object, "#") {
		return false
	}
	r, ok := q.b.ref(tr, false)
	if !ok || r.where != placed || r.temp {
		return false
	}
	return q.b.see(spaceRelation, r.parent, r.n).fact == walk.Present
}

// isPseudo reports whether a name means inserted or deleted in a DML
// trigger's body whose table the walk knows.
func (q *binder) isPseudo(tr *ast.TableRef) bool {
	if q.pseudo == "" && q.pseudoCols == nil || tr.Schema != "" || tr.Database != "" || tr.Server != "" || q.isCTE(tr) {
		return false
	}
	return strings.EqualFold(tr.Object, "inserted") || strings.EqualFold(tr.Object, "deleted")
}

// tableVar returns the declaration of the table variable a name means,
// nil when the walk cannot tell which.
func (q *binder) tableVar(name string) *ast.VariableDecl {
	var found *ast.VariableDecl
	for _, d := range q.tableVars {
		if d.Name == name {
			found = d
			continue
		}
		if !differ(strings.TrimPrefix(d.Name, "@"), strings.TrimPrefix(name, "@")) {
			return nil
		}
	}
	return found
}

// statement binds one statement of a module body.
func (q *binder) statement(n ast.Node) {
	switch v := n.(type) {
	case *ast.SelectStmt:
		// A column of a query with no table is refused as well (verified).
		q.noFrom = true
		q.query(v, nil)
		q.noFrom = false
	case *ast.InsertStmt:
		if v.WithClause != nil {
			return
		}
		sc := &qscope{}
		var conds []ast.ExprNode
		q.source(v.Relation, sc, &conds, false)
		if v.Cols != nil && !sc.opaque && len(sc.sources) == 1 {
			for _, item := range v.Cols.Items {
				if s, ok := item.(*ast.String); ok {
					q.column(&ast.ColumnRef{Column: s.Str, Loc: v.Loc}, sc)
				}
			}
		}
		if src, ok := v.Source.(*ast.SelectStmt); ok {
			q.query(src, nil)
		}
	case *ast.UpdateStmt:
		if v.WithClause != nil || v.FromClause != nil && len(v.FromClause.Items) > 0 {
			return
		}
		sc := &qscope{}
		var conds []ast.ExprNode
		q.source(v.Relation, sc, &conds, false)
		if v.SetClause != nil {
			for _, item := range v.SetClause.Items {
				se, ok := item.(*ast.SetExpr)
				if !ok {
					continue
				}
				if se.Column != nil && !se.WriteMethod {
					q.column(se.Column, sc)
				}
				q.expr(se.Value, sc)
			}
		}
		q.expr(v.WhereClause, sc)
	case *ast.DeleteStmt:
		if v.WithClause != nil || v.FromClause != nil && len(v.FromClause.Items) > 0 {
			return
		}
		sc := &qscope{}
		var conds []ast.ExprNode
		q.source(v.Relation, sc, &conds, false)
		q.expr(v.WhereClause, sc)
	case *ast.MergeStmt:
		if v.WithClause != nil {
			return
		}
		sc := &qscope{}
		var conds []ast.ExprNode
		q.source(v.Target, sc, &conds, false)
		q.source(v.Source, sc, &conds, false)
		for _, c := range conds {
			q.expr(c, sc)
		}
		q.expr(v.OnCondition, sc)
	case *ast.DeclareCursorStmt:
		q.query(v.Query, nil)
	default:
		// A condition, an assignment, a declaration's default, a
		// RETURN: the queries in it stand alone.
		q.expr(n, nil)
	}
}

// hasName reports whether a list may hold a name under some collation.
func hasName(names []string, n string) bool {
	for _, x := range names {
		if !differ(x, n) {
			return true
		}
	}
	return false
}

// insertCounts fails an INSERT whose column list names fewer or more
// columns than its VALUES rows or its query give (109, 110, 120, 121).
func (b *builder) insertCounts(v *ast.InsertStmt) bool {
	if v.Cols == nil || len(v.Cols.Items) == 0 || v.DefaultValues {
		return false
	}
	n := len(v.Cols.Items)
	got := -1
	switch src := v.Source.(type) {
	case *ast.ValuesClause:
		if src.Rows == nil {
			return false
		}
		for _, row := range src.Rows.Items {
			l, ok := row.(*ast.List)
			if !ok {
				return false
			}
			if len(l.Items) != n {
				got = len(l.Items)
				break
			}
		}
	case *ast.SelectStmt:
		if c, ok := selectWidth(src); ok && c != n {
			got = c
		}
	}
	if got < 0 {
		return false
	}
	b.fail(rangeOf(v.Loc), review.UnsupportedForm, "SQL Server doesn't accept an `INSERT` that lists "+count(n, "column")+" and "+count(got, "value")+".")
	return true
}

// selectWidth returns how many columns a query returns, false when the
// walk cannot count them: a *, or branches that differ.
func selectWidth(s *ast.SelectStmt) (int, bool) {
	if s == nil {
		return 0, false
	}
	if s.Larg != nil || s.Rarg != nil {
		l, ok1 := selectWidth(s.Larg)
		r, ok2 := selectWidth(s.Rarg)
		if !ok1 || !ok2 || l != r {
			return 0, false
		}
		return l, true
	}
	if s.TargetList == nil || s.ForClause != nil {
		return 0, false
	}
	for _, item := range s.TargetList.Items {
		if _, ok := item.(*ast.StarExpr); ok {
			return 0, false
		}
		if rt, ok := item.(*ast.ResTarget); ok {
			switch rt.Val.(type) {
			case *ast.StarExpr, *ast.SelectAssign:
				return 0, false
			}
		}
	}
	return len(s.TargetList.Items), true
}

// setWidths fails a statement whose set operation joins queries of
// different widths (205), at any depth. Each query's width is counted
// once, so a long chain of branches costs no more than its length.
func (b *builder) setWidths(n ast.Node) bool {
	type width struct {
		n  int
		ok bool
	}
	memo := make(map[*ast.SelectStmt]width)
	var widthOf func(s *ast.SelectStmt) (int, bool)
	widthOf = func(s *ast.SelectStmt) (int, bool) {
		if s == nil || s.Larg == nil && s.Rarg == nil {
			return selectWidth(s)
		}
		if w, seen := memo[s]; seen {
			return w.n, w.ok
		}
		l, ok1 := widthOf(s.Larg)
		r, ok2 := widthOf(s.Rarg)
		w := width{l, ok1 && ok2 && l == r}
		memo[s] = w
		return w.n, w.ok
	}
	var bad *ast.SelectStmt
	ast.Inspect(n, func(m ast.Node) bool {
		if bad != nil {
			return false
		}
		s, ok := m.(*ast.SelectStmt)
		if !ok || s.Larg == nil || s.Rarg == nil {
			return true
		}
		l, ok1 := widthOf(s.Larg)
		r, ok2 := widthOf(s.Rarg)
		if ok1 && ok2 && l != r {
			bad = s
			return false
		}
		return true
	})
	if bad == nil {
		return false
	}
	b.fail(rangeOf(bad.Loc), review.UnsupportedForm, "SQL Server doesn't accept a `UNION`, `INTERSECT`, or `EXCEPT` of queries with different numbers of columns.")
	return true
}

// orderWithoutTop finds an ORDER BY that SQL Server refuses in a view
// or an inline function, at any depth (1033): one without TOP, OFFSET,
// or FOR XML or JSON.
func orderWithoutTop(n ast.Node) (review.Range, bool) {
	var rng review.Range
	found := false
	memo := make(map[*ast.SelectStmt]bool)
	var anyTop func(s *ast.SelectStmt) bool
	anyTop = func(s *ast.SelectStmt) bool {
		if s == nil {
			return false
		}
		if s.Top != nil {
			return true
		}
		if v, seen := memo[s]; seen {
			return v
		}
		v := anyTop(s.Larg) || anyTop(s.Rarg)
		memo[s] = v
		return v
	}
	ast.Inspect(n, func(m ast.Node) bool {
		if found {
			return false
		}
		s, ok := m.(*ast.SelectStmt)
		if !ok || s.OrderByClause == nil || len(s.OrderByClause.Items) == 0 || s.OffsetClause != nil || s.ForClause != nil || anyTop(s) {
			return true
		}
		rng, found = rangeOf(s.Loc), true
		return false
	})
	return rng, found
}

// schemaBound fails a schema-bound view's or function's query that names
// a table by one part (4512), a synced view or function that isn't
// schema-bound (4513), or, with star, selects * (1054).
func (b *builder) schemaBound(q *binder, n ast.Node, kind string, star bool) bool {
	var rng review.Range
	msg := ""
	unbound := func(loc ast.Loc, written string) {
		rng, msg = rangeOf(loc), "SQL Server doesn't accept "+review.Code(written)+", which isn't schema-bound, in a schema-bound "+kind+"."
	}
	forEachFromTable(n, func(tr *ast.TableRef) {
		if msg == "" && tr.Schema == "" && tr.Database == "" && tr.Server == "" && !q.isCTE(tr) {
			rng, msg = rangeOf(tr.Loc), "SQL Server doesn't accept the one-part name "+review.Code(object(tr))+" in a schema-bound "+kind+"."
		}
		if msg == "" && tr.Database == "" && tr.Server == "" && (strings.EqualFold(tr.Schema, "sys") || strings.EqualFold(tr.Schema, "INFORMATION_SCHEMA")) {
			rng, msg = rangeOf(tr.Loc), "SQL Server doesn't accept the system object "+review.Code(object(tr))+" in a schema-bound "+kind+"."
		}
		if msg == "" && tr.Schema != "" {
			if r, ok := b.ref(tr, false); ok && b.unboundView(r) {
				unbound(tr.Loc, object(tr))
			}
		}
	})
	if msg == "" {
		ast.Inspect(n, func(m ast.Node) bool {
			if f, ok := m.(*ast.FuncCallExpr); ok && msg == "" && f.Name != nil && f.Name.Schema != "" {
				if fn := b.syncedFunction(f.Name); fn != nil && !fn.bound {
					unbound(f.Loc, object(f.Name))
				}
			}
			return msg == ""
		})
	}
	if msg == "" && star {
		ast.Inspect(n, func(m ast.Node) bool {
			if msg != "" {
				return false
			}
			if st, ok := m.(*ast.StarExpr); ok {
				rng, msg = rangeOf(st.Loc), "SQL Server doesn't accept `*` in a schema-bound "+kind+"."
				return false
			}
			switch v := m.(type) {
			case *ast.FuncCallExpr:
				if v.Star {
					// COUNT(*) is no column list.
					return false
				}
			case *ast.ExistsExpr:
				return false
			}
			return true
		})
	}
	if msg == "" {
		return false
	}
	b.fail(rng, review.UnsupportedForm, msg)
	return true
}

// hasOption reports whether a module's WITH options hold one.
func hasOption(opts *ast.List, name string) bool {
	if opts == nil {
		return false
	}
	for _, item := range opts.Items {
		if s, ok := item.(*ast.String); ok && strings.EqualFold(s.Str, name) {
			return true
		}
	}
	return false
}

const lastReturn = "SQL Server doesn't accept a function whose last statement isn't `RETURN`."

// count writes n of a noun.
func count(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(n) + " " + noun + "s"
}

// functionBody fails a scalar or multi-statement function whose body
// changes a table other than its own table variables (443), calls
// NEWID (443), or does not end with RETURN (455).
func (b *builder) functionBody(body ast.Node) bool {
	be, ok := body.(*ast.BeginEndStmt)
	if !ok || be.Stmts == nil {
		return false
	}
	var rng review.Range
	msg := ""
	ast.Inspect(be, func(m ast.Node) bool {
		if msg != "" {
			return false
		}
		var target ast.TableExpr
		verb := ""
		switch v := m.(type) {
		case *ast.InsertStmt:
			if v.WithClause != nil {
				return true
			}
			target, verb = v.Relation, "INSERT"
		case *ast.UpdateStmt:
			if v.WithClause != nil || v.FromClause != nil {
				return true
			}
			target, verb = v.Relation, "UPDATE"
		case *ast.DeleteStmt:
			if v.WithClause != nil || v.FromClause != nil {
				return true
			}
			target, verb = v.Relation, "DELETE"
		case *ast.MergeStmt:
			if v.WithClause != nil {
				return true
			}
			target, verb = v.Target, "MERGE"
		case *ast.VariableDecl:
			// A table variable's default is not the function's call.
			return false
		case *ast.FuncCallExpr:
			if v.Name != nil && v.Name.Schema == "" && strings.EqualFold(v.Name.Object, "NEWID") {
				rng, msg = rangeOf(v.Loc), "SQL Server doesn't accept `NEWID` in a function."
			}
			return true
		default:
			return true
		}
		if tr, ok := target.(*ast.TableRef); ok && !strings.HasPrefix(tr.Object, "@") {
			rng, msg = rangeOf(ast.NodeLoc(m)), "SQL Server doesn't accept `"+verb+"` of a table in a function."
		}
		return true
	})
	if msg == "" && len(be.Stmts.Items) > 0 {
		// IF ... RETURN ELSE RETURN fails too (verified); a block that
		// ends with RETURN does not.
		switch last := be.Stmts.Items[len(be.Stmts.Items)-1].(type) {
		case *ast.ReturnStmt, *ast.BeginEndStmt:
		default:
			rng, msg = rangeOf(ast.NodeLoc(last)), lastReturn
		}
	}
	if msg == "" {
		return false
	}
	b.fail(rng, review.UnsupportedForm, msg)
	return true
}

// ungrouped fails a query block whose select list names a column that
// its GROUP BY, a list of columns, does not (8120). Only a block no other
// encloses is held to it, since an outer block's column is one value per
// group.
func (q *binder) ungrouped(s *ast.SelectStmt) bool {
	if s.GroupByClause == nil || len(s.GroupByClause.Items) == 0 || s.GroupByAll || s.TargetList == nil {
		return false
	}
	var grouped []string
	for _, item := range s.GroupByClause.Items {
		c, ok := item.(*ast.ColumnRef)
		if !ok {
			return false
		}
		grouped = append(grouped, c.Column)
	}
	for _, item := range s.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		c, ok := rt.Val.(*ast.ColumnRef)
		if !ok || isPseudoCol(c.Column) || hasName(grouped, c.Column) {
			continue
		}
		q.b.fail(rangeOf(c.Loc), review.UnsupportedForm, "SQL Server doesn't accept column "+review.Code(ident(c.Column))+", which is neither grouped nor aggregated.")
		return true
	}
	// HAVING (8121) and ORDER BY (8127) hold to it too, outside an
	// aggregate; ORDER BY may name the select list's columns instead.
	var bad *ast.ColumnRef
	clause := ""
	inBlock(s.HavingClause, func(m ast.Node) bool {
		if bad != nil || isAggregate(m) {
			return false
		}
		if c, ok := m.(*ast.ColumnRef); ok && c.Schema == "" && !isPseudoCol(c.Column) && !hasName(grouped, c.Column) {
			bad, clause = c, "HAVING"
		}
		return true
	})
	var selected []string
	for _, item := range s.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok {
			if rt.Name != "" {
				selected = append(selected, rt.Name)
			} else if c, ok := unparen(rt.Val).(*ast.ColumnRef); ok {
				selected = append(selected, c.Column)
			}
		}
	}
	for _, item := range listNodes(s.OrderByClause) {
		o, ok := item.(*ast.OrderByItem)
		if !ok || bad != nil {
			continue
		}
		if c, ok := o.Expr.(*ast.ColumnRef); ok && c.Schema == "" && !isPseudoCol(c.Column) && !hasName(grouped, c.Column) && !hasName(selected, c.Column) {
			bad, clause = c, "ORDER BY"
		}
	}
	if bad == nil {
		return false
	}
	q.b.fail(rangeOf(bad.Loc), review.UnsupportedForm, "SQL Server doesn't accept column "+review.Code(ident(bad.Column))+" in `"+clause+"`, which is neither grouped nor aggregated.")
	return true
}

// call checks a function call of a bound statement, and reports whether
// it failed: a synced scalar function takes no more arguments than it
// declares (8144).
func (q *binder) call(f *ast.FuncCallExpr) bool {
	if !q.bound || f.Star || f.Name == nil || f.Name.Schema == "" {
		return false
	}
	fn := q.b.syncedFunction(f.Name)
	if fn == nil || fn.kind != funcScalar || len(listNodes(f.Args)) <= fn.params {
		return false
	}
	q.b.fail(rangeOf(f.Loc), review.UnsupportedForm, msgArgs(object(f.Name), fn.params))
	return true
}

// msgArgs says how many arguments a function takes.
func msgArgs(fn string, n int) string {
	return "Function " + review.Code(fn) + " takes " + count(n, "argument") + "."
}

// pseudoColumn checks $IDENTITY or $ROWGUID in a block with a source, and
// reports whether it did: it fails when every source, and every source of
// the blocks around it, is a synced table that lacks the column the name
// means (207), as lacks says: an identity column, or one that can be the
// ROWGUIDCOL column.
func (q *binder) pseudoColumn(c *ast.ColumnRef, sc *qscope, lacks func(*snapTable) bool) bool {
	written := ""
	for s := sc; s != nil; s = s.outer {
		if s.opaque {
			return true
		}
		for _, src := range s.sources {
			if src.local || src.path == "" {
				return true
			}
			if st := q.b.pathTable(src.path); st == nil || !lacks(st) {
				return true
			}
			if written == "" {
				written = src.written
			}
		}
	}
	if written == "" {
		return false
	}
	q.b.fail(rangeOf(c.Loc), review.ObjectNotFound, msgMissingIn("column", strings.ToUpper(c.Column), written))
	return true
}
