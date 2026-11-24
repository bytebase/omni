package review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	obcatalog "github.com/bytebase/omni/oceanbase/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// queryScope is one query block of a view's or CREATE TABLE ... AS's query
// as the walk reads it: whether a column there may name something, the CTE
// names it sees, and the names its select list gives a subquery inside it.
type queryScope struct {
	outer *queryScope
	// tables is set when the block reads a table or another row source, or,
	// for a derived table or a table function, when another item of the FROM
	// it stands in may lend it columns (LATERAL).
	tables bool
	ctes   []string
	// names are the select list's names, lower-cased, which a subquery's
	// column may refer to as an outer reference; anyName is set when an item
	// has a name the walk does not derive.
	names   []string
	anyName bool
}

// sees reports whether a column of the name may resolve in the scope: a
// row source of the block or of an enclosing one, or a name an enclosing
// block's select list gives. The block's own names are not in reach: a
// select list item or WHERE refers to no alias of its list (1054).
func (s *queryScope) sees(col string) bool {
	if s.tables {
		return true
	}
	for p := s.outer; p != nil; p = p.outer {
		if p.tables || p.anyName {
			return true
		}
		for _, n := range p.names {
			if strings.EqualFold(n, col) {
				return true
			}
		}
	}
	return false
}

// readsTables reports whether the scope or an enclosing one reads a row
// source.
func (s *queryScope) readsTables() bool {
	for p := s; p != nil; p = p.outer {
		if p.tables {
			return true
		}
	}
	return false
}

// cte reports whether a table name may be a CTE the scope sees. CTE names
// are compared without case, so a name that may be one is never taken for a
// table.
func (s *queryScope) cte(name string) bool {
	for p := s; p != nil; p = p.outer {
		for _, c := range p.ctes {
			if strings.EqualFold(c, name) {
				return true
			}
		}
	}
	return false
}

// queryCheck reads a query for what OceanBase refuses in it, given the
// tables the ledger holds: a column with no row source in reach (1054), * or
// t.* with none (1096, 1051), a table named twice in one FROM (1066), a table
// that does not exist (1146), a window no WINDOW clause defines (5929), a
// VALUES row with no column (3942) or a bare DEFAULT (3943), and in a view a
// variable (1351) or a parameter marker (1064). It stops at the first.
type queryCheck struct {
	b    *builder
	view bool
	done bool
}

// queryFails reports whether the query certainly fails, having recorded why.
func (b *builder) queryFails(sel *ast.SelectStmt, view bool) bool {
	q := &queryCheck{b: b, view: view}
	q.query(sel, &queryScope{})
	if !q.done && view {
		return b.viewFunctionsFail(sel)
	}
	return q.done
}

// viewFunctionsFail fails a view whose query calls, unqualified, a builtin
// function with no argument that takes some (1582), or a function that is
// neither a builtin nor a stored function of the current database (1305),
// as OceanBase 4.3.5 does once the query resolves (oceanbase/catalog's
// validateViewFuncExists). The walk does not resolve every column, so it
// reads only a query that names no table and no column.
func (b *builder) viewFunctionsFail(sel *ast.SelectStmt) bool {
	plain := true
	ast.Inspect(sel, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ColumnRef, *ast.TableRef, *ast.StarExpr, *ast.DefaultExpr:
			plain = false
		}
		return plain
	})
	if !plain || b.w.noDatabase || b.w.lost {
		return false
	}
	failed, definite := false, false
	ast.Inspect(sel, func(n ast.Node) bool {
		f, ok := n.(*ast.FuncCallExpr)
		if failed || !ok || f.Schema != "" || f.Name == "" {
			return !failed
		}
		name := strings.ToLower(f.Name)
		if obcatalog.IsBuiltinFunction(name) {
			if len(f.Args) == 0 && obcatalog.BuiltinRequiresArgs(name) {
				b.fail(rangeOf(f.Loc), nil, review.UnsupportedForm, review.Code(strings.ToUpper(name)+"()")+" needs at least one argument.")
				failed, definite = true, true
			}
			return !failed
		}
		if k, ok := b.key(spaceFunction, f.Name); ok {
			b.need(objKey(spaceFunction, b.w.current, k), walk.Present, nil, rangeOf(f.Loc), notFound("function", ident(f.Name), ""), "") // 1305
		}
		failed = true
		return false
	})
	return definite
}

func (q *queryCheck) fail(loc ast.Loc, code review.FindingCode, message string) {
	if q.done {
		return
	}
	q.done = true
	q.b.fail(rangeOf(loc), nil, code, message)
}

// query reads one query expression whose blocks see outer.
func (q *queryCheck) query(sel *ast.SelectStmt, outer *queryScope) {
	if sel == nil || q.done {
		return
	}
	if len(sel.CTEs) > 0 {
		with := &queryScope{outer: outer}
		for _, c := range sel.CTEs {
			with.ctes = append(with.ctes, c.Name)
		}
		// A CTE's query may refer to the columns of enclosing blocks the walk
		// does not resolve; it is read for everything but its columns.
		body := &queryScope{outer: with, tables: true}
		for _, c := range sel.CTEs {
			q.query(c.Select, body)
		}
		outer = with
	}
	switch {
	case sel.ParenSource != nil:
		q.query(sel.ParenSource, outer)
		return
	case sel.SetOp != ast.SetOpNone:
		q.query(sel.Left, outer)
		q.query(sel.Right, outer)
		return
	case sel.ValuesSource != nil:
		q.values(sel.ValuesSource, &queryScope{outer: outer})
		return
	case sel.TableSource != nil:
		return
	}
	s := &queryScope{outer: outer}
	q.from(sel.From, s)
	if q.done {
		return
	}
	q.ambiguous(sel, s)
	if q.done {
		return
	}
	for _, t := range sel.TargetList {
		switch v := t.(type) {
		case *ast.ResTarget:
			s.names = append(s.names, v.Name)
		case *ast.ColumnRef:
			if v.Star {
				s.anyName = true
			} else {
				s.names = append(s.names, v.Column)
			}
		case *ast.IntLit, *ast.SubqueryExpr, *ast.ExistsExpr:
			// Named after its text, which names no column a subquery
			// refers to by a plain name: digits alone are no column.
		default:
			s.anyName = true
		}
	}
	windows := make([]string, 0, len(sel.WindowClause))
	for _, w := range sel.WindowClause {
		if w != nil {
			windows = append(windows, w.Name)
		}
	}
	for _, t := range sel.TargetList {
		q.expr(t, s, true, windows)
	}
	q.expr(sel.Where, s, true, nil)
	for _, g := range sel.GroupBy {
		if g != nil {
			q.expr(g.Expr, s, false, nil)
		}
	}
	q.expr(sel.Having, s, false, windows)
	for _, o := range sel.OrderBy {
		if o != nil {
			q.expr(o.Expr, s, false, windows)
		}
	}
	for _, w := range sel.WindowClause {
		if w != nil {
			q.expr(w, s, false, nil)
		}
	}
}

// from reads a FROM clause into s: its tables must exist and have names of
// their own, and the queries and table functions in it are read against
// what they may see.
func (q *queryCheck) from(items []ast.TableExpr, s *queryScope) {
	var leaves []ast.TableExpr
	var conds []ast.Node
	var collect func(te ast.TableExpr)
	collect = func(te ast.TableExpr) {
		switch v := te.(type) {
		case *ast.JoinClause:
			collect(v.Left)
			collect(v.Right)
			if v.Condition != nil {
				conds = append(conds, v.Condition)
			}
		case nil:
		default:
			leaves = append(leaves, te)
		}
	}
	for _, te := range items {
		collect(te)
	}
	// DUAL reads no table, unless a table of the name exists.
	dual := false
	for _, te := range leaves {
		tr, ok := te.(*ast.TableRef)
		if !ok || !q.isDual(tr) {
			s.tables = true
			continue
		}
		dual = true
	}
	if dual && len(leaves) > 1 {
		s.tables = true
	}
	type exposed struct{ db, name string }
	seen := make(map[exposed]bool, len(leaves))
	for _, te := range leaves {
		if q.done {
			return
		}
		switch v := te.(type) {
		case *ast.TableRef:
			if q.isDual(v) {
				continue
			}
			q.table(v, s)
			if q.done {
				return
			}
			// A CTE's reference by its own name is left out.
			if v.Alias == "" && v.Schema == "" && s.cte(v.Name) {
				continue
			}
			r, ok := q.resolved(v)
			if !ok {
				continue
			}
			name := r.name
			if v.Alias != "" {
				k, ok := q.b.key(spaceRelation, v.Alias)
				if !ok {
					continue
				}
				name = k
			}
			if strings.HasPrefix(name, nonASCII) {
				continue
			}
			e := exposed{r.db, name}
			if seen[e] {
				written := v.Alias
				if written == "" {
					written = v.Name
				}
				q.fail(v.Loc, review.DuplicateName, "This query names table "+review.Code(ident(written))+" twice in one FROM clause.")
				return
			}
			seen[e] = true
		case *ast.SubqueryExpr:
			// A derived table sees the blocks around its own, and, LATERAL,
			// the other items of its FROM.
			q.query(v.Select, &queryScope{outer: s.outer, tables: len(leaves) > 1})
		case *ast.JsonTableExpr:
			q.expr(v.Expr, &queryScope{outer: s.outer, tables: len(leaves) > 1}, true, nil)
		case *ast.RbIterateExpr:
			q.expr(v.Arg, &queryScope{outer: s.outer, tables: len(leaves) > 1}, true, nil)
		case *ast.UnnestExpr:
			for _, a := range v.Args {
				q.expr(a, &queryScope{outer: s.outer, tables: len(leaves) > 1}, true, nil)
			}
		}
	}
	for _, c := range conds {
		q.expr(c, s, false, nil)
	}
}

// isDual reports whether a FROM item is DUAL, which reads no table: an
// unqualified dual with no alias, where no table of the name exists.
func (q *queryCheck) isDual(tr *ast.TableRef) bool {
	if tr.Schema != "" || tr.Catalog != "" || tr.Dblink != "" || tr.Alias != "" || !strings.EqualFold(tr.Name, "dual") {
		return false
	}
	r, ok := q.resolved(tr)
	if !ok {
		return false
	}
	f, temp, _ := q.b.resolve(r)
	return !temp && f == walk.Absent
}

// resolved resolves a FROM item's table without failing the statement for a
// missing current database.
func (q *queryCheck) resolved(tr *ast.TableRef) (tableRef, bool) {
	if tr.Schema == "" && (q.b.w.lost || q.b.w.noDatabase) {
		return tableRef{}, false
	}
	r, st := q.b.ref(tr)
	return r, st == refOK
}

// table checks that a FROM item's table exists, unless it may be a CTE.
func (q *queryCheck) table(tr *ast.TableRef, s *queryScope) {
	if tr.Schema == "" && s.cte(tr.Name) {
		return
	}
	r, ok := q.resolved(tr)
	if !ok {
		return
	}
	f, temp, _ := q.b.resolve(r)
	if temp || f != walk.Absent {
		return
	}
	q.done = true
	q.b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("", r.written, ""), "") // 1146
}

// values reads a VALUES query's rows.
func (q *queryCheck) values(v *ast.ValuesStmt, s *queryScope) {
	if v.LeadingComma {
		q.fail(v.Loc, review.UnsupportedForm, "Each row of a VALUES query must have a value for every column.")
		return
	}
	for _, row := range v.Rows {
		if len(row) == 0 {
			q.fail(v.Loc, review.UnsupportedForm, "Each row of a VALUES query must have at least one column.")
			return
		}
		for _, e := range row {
			if d, ok := e.(*ast.DefaultExpr); ok && d.Column == "" {
				q.fail(d.Loc, review.UnsupportedForm, "A VALUES query cannot use `DEFAULT`.")
				return
			}
			q.expr(e, s, true, nil)
		}
	}
}

// expr reads an expression of a block: its subqueries against the block,
// and, where cols is set, each column it names against what the block sees.
// windows are the block's WINDOW names, which an OVER clause there must
// name; nil where the expression cannot hold a window function of the block.
func (q *queryCheck) expr(n ast.Node, s *queryScope, cols bool, windows []string) {
	if n == nil || q.done {
		return
	}
	ast.Inspect(n, func(m ast.Node) bool {
		if q.done || m == nil {
			return false
		}
		switch v := m.(type) {
		case *ast.SubqueryExpr:
			q.query(v.Select, s)
			return false
		case *ast.ExistsExpr:
			q.query(v.Select, s)
			return false
		case *ast.SelectStmt:
			q.query(v, s)
			return false
		case *ast.LambdaExpr:
			return false
		case *ast.VariableRef:
			if q.view && !v.Bare {
				name := "@" + v.Name
				if v.System {
					name = "@@" + v.Name
				}
				q.fail(v.Loc, review.UnsupportedForm, "A view's query cannot use variable "+review.Code(name)+".")
			}
			return false
		case *ast.ParamMarkerExpr:
			if q.view {
				q.fail(v.Loc, review.UnsupportedForm, "A view's query cannot use a parameter marker.")
			}
			return false
		case *ast.StarExpr:
			if !s.readsTables() {
				q.fail(v.Loc, review.UnsupportedForm, "This query reads `*` but names no table.")
			}
			return false
		case *ast.ColumnRef:
			switch {
			case v.Star && v.Table != "" && !s.readsTables():
				q.fail(v.Loc, review.ObjectNotFound, notFound("table", ident(v.Table), ""))
			case cols && !v.Star && v.Column != "" && (v.Table == "" && !s.sees(v.Column) || v.Table != "" && !s.readsTables()):
				q.fail(v.Loc, review.ObjectNotFound, notFound("column", ident(v.Column), ""))
			}
			return false
		case *ast.DefaultExpr:
			if cols && v.Column != "" && !s.sees(v.Column) {
				q.fail(v.Loc, review.ObjectNotFound, notFound("column", ident(v.Column), ""))
			}
			return false
		case *ast.FuncCallExpr:
			if w := v.Over; w != nil && windows != nil && w.RefName != "" && !hasFold(windows, w.RefName) {
				q.fail(w.Loc, review.ObjectNotFound, "Window "+review.Code(ident(w.RefName))+" is not defined.")
				return false
			}
		}
		return true
	})
}

func hasFold(list []string, name string) bool {
	for _, s := range list {
		if strings.EqualFold(s, name) {
			return true
		}
	}
	return false
}

// ambiguous fails a block whose select list or WHERE names, unqualified, a
// column two of its tables have (1052). A NATURAL or USING join merges
// columns of a name, so a block with one is left out, and only a table whose
// column the ledger holds counts.
func (q *queryCheck) ambiguous(sel *ast.SelectStmt, s *queryScope) {
	var tables []tableRef
	merges := false
	var collect func(te ast.TableExpr)
	collect = func(te ast.TableExpr) {
		switch v := te.(type) {
		case *ast.JoinClause:
			switch v.Type {
			case ast.JoinNatural, ast.JoinNaturalLeft, ast.JoinNaturalRight, ast.JoinNaturalFull:
				merges = true
			}
			if _, ok := v.Condition.(*ast.UsingCondition); ok {
				merges = true
			}
			collect(v.Left)
			collect(v.Right)
		case *ast.TableRef:
			if v.Schema == "" && s.cte(v.Name) {
				return
			}
			r, ok := q.resolved(v)
			if !ok {
				return
			}
			if f, temp, _ := q.b.resolve(r); temp || f != walk.Present {
				return
			}
			tables = append(tables, r)
		}
	}
	for _, te := range sel.From {
		collect(te)
	}
	if merges || len(tables) < 2 {
		return
	}
	check := func(n ast.Node) {
		ast.Inspect(n, func(m ast.Node) bool {
			if q.done || m == nil {
				return false
			}
			switch v := m.(type) {
			case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SelectStmt, *ast.LambdaExpr:
				return false
			case *ast.ColumnRef:
				if v.Star || v.Table != "" || v.Column == "" {
					return false
				}
				k, ok := q.b.key(spaceColumn, v.Column)
				if !ok || strings.HasPrefix(k, nonASCII) {
					return false
				}
				n := 0
				for _, r := range tables {
					if f, _ := q.b.l.Lookup(colKey(r.db, r.name, k)); f == walk.Present {
						n++
					}
				}
				if n >= 2 {
					q.fail(v.Loc, review.DuplicateName, "Column "+review.Code(ident(v.Column))+" is in more than one table of this query.")
				}
				return false
			}
			return true
		})
	}
	for _, t := range sel.TargetList {
		check(t)
	}
	check(sel.Where)
}
