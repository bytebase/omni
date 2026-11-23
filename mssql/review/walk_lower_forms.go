package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// The checks here fail a view's or an inline function's query on a form
// SQL Server refuses whatever the schema holds (measured on SQL Server
// 2022, mssql/catalog/testdata/queries). A module's other statements are
// left alone: deferred name resolution may put off compiling them.

// aggregates are the system aggregate functions. A user-defined aggregate
// is always called by a schema-qualified name.
var aggregates = map[string]bool{
	"SUM": true, "AVG": true, "MIN": true, "MAX": true, "COUNT": true, "COUNT_BIG": true,
	"STDEV": true, "STDEVP": true, "VAR": true, "VARP": true, "STRING_AGG": true,
	"CHECKSUM_AGG": true, "GROUPING": true, "GROUPING_ID": true, "APPROX_COUNT_DISTINCT": true,
}

// orderedWindows are the window functions whose OVER clause must order
// the rows (4112).
var orderedWindows = map[string]bool{
	"ROW_NUMBER": true, "RANK": true, "DENSE_RANK": true, "NTILE": true,
	"LAG": true, "LEAD": true, "PERCENT_RANK": true, "CUME_DIST": true,
}

// builtinArity is how many arguments some system functions take (174 for
// a fixed count, 189 for a range).
var builtinArity = map[string][2]int{
	"LEN": {1, 1}, "DATALENGTH": {1, 1}, "SUBSTRING": {3, 3}, "LEFT": {2, 2}, "RIGHT": {2, 2},
	"REPLACE": {3, 3}, "REVERSE": {1, 1}, "REPLICATE": {2, 2}, "STUFF": {4, 4},
	"UPPER": {1, 1}, "LOWER": {1, 1}, "ABS": {1, 1}, "CEILING": {1, 1}, "FLOOR": {1, 1},
	"POWER": {2, 2}, "SQRT": {1, 1}, "SQUARE": {1, 1}, "ROUND": {2, 3}, "CONCAT": {2, 254},
	"YEAR": {1, 1}, "MONTH": {1, 1}, "DAY": {1, 1},
}

// dateParts are the dateparts some date function accepts; a name none
// accepts is no datepart (155).
var dateParts = map[string]bool{
	"year": true, "yy": true, "yyyy": true, "quarter": true, "qq": true, "q": true,
	"month": true, "mm": true, "m": true, "dayofyear": true, "dy": true, "y": true,
	"day": true, "dd": true, "d": true, "week": true, "wk": true, "ww": true,
	"weekday": true, "dw": true, "w": true, "hour": true, "hh": true,
	"minute": true, "mi": true, "n": true, "second": true, "ss": true, "s": true,
	"millisecond": true, "ms": true, "microsecond": true, "mcs": true,
	"nanosecond": true, "ns": true, "tzoffset": true, "tz": true,
	"iso_week": true, "isowk": true, "isoww": true,
}

// builtin returns a function's name in upper case when it is unqualified,
// which means a system function, else "".
func builtin(f *ast.FuncCallExpr) string {
	if f == nil || f.Name == nil || f.Name.Schema != "" || f.Name.Database != "" || f.Name.Server != "" {
		return ""
	}
	return strings.ToUpper(f.Name.Object)
}

// isAggregate reports whether a call is a system aggregate, not windowed.
func isAggregate(n ast.Node) bool {
	f, ok := n.(*ast.FuncCallExpr)
	return ok && f.Over == nil && aggregates[builtin(f)]
}

// inBlock calls f for each node of an expression of one query block, not
// entering the queries nested in it.
func inBlock(n ast.Node, f func(ast.Node) bool) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(m ast.Node) bool {
		switch m.(type) {
		case *ast.SelectStmt, *ast.SubqueryExpr:
			return false
		}
		return f(m)
	})
}

// listNodes returns a list's items, nil for no list.
func listNodes(l *ast.List) []ast.Node {
	if l == nil {
		return nil
	}
	return l.Items
}

// viewForms fails a view's or an inline function's query on a form SQL
// Server refuses before it binds a name: SELECT INTO (156) or an OPTION
// clause (156) in a view; TABLESAMPLE (478) or NEXT VALUE FOR (11719) in
// either; TOP beside OFFSET (10741); and a WITH that names a CTE twice
// (239).
func (b *builder) viewForms(query ast.Node, kind string) bool {
	var rng review.Range
	msg := ""
	set := func(loc ast.Loc, m string) {
		if msg == "" {
			rng, msg = rangeOf(loc), m
		}
	}
	ast.Inspect(query, func(m ast.Node) bool {
		if msg != "" {
			return false
		}
		switch v := m.(type) {
		case *ast.SelectStmt:
			if kind == "view" && v.IntoTable != nil {
				set(v.IntoTable.Loc, "SQL Server doesn't accept `SELECT INTO` in a view.")
			}
			if kind == "view" && v.OptionClause != nil && len(v.OptionClause.Items) > 0 {
				set(v.Loc, "SQL Server doesn't accept an `OPTION` clause in a view.")
			}
			if v.Top != nil && v.OffsetClause != nil {
				set(v.Top.Loc, "SQL Server doesn't accept `TOP` beside `OFFSET`.")
			}
			if v.WithClause != nil {
				var names []string
				for _, item := range listNodes(v.WithClause.CTEs) {
					if cte, ok := item.(*ast.CommonTableExpr); ok {
						names = append(names, cte.Name)
					}
				}
				if n := twice(names); n != "" {
					set(v.WithClause.Loc, msgTwice("names", "CTE", ident(n)))
				}
			}
		case *ast.AliasedTableRef:
			if v.TableSample != nil {
				set(v.TableSample.Loc, "SQL Server doesn't accept `TABLESAMPLE` in a "+kind+".")
			}
		case *ast.NextValueForExpr:
			set(v.Loc, "SQL Server doesn't accept `NEXT VALUE FOR` in a "+kind+".")
		}
		return true
	})
	if msg == "" {
		return false
	}
	b.fail(rng, review.UnsupportedForm, msg)
	return true
}

// queryForms fails a view's or an inline function's query on a form SQL
// Server refuses when it binds the query, whatever the tables hold.
func (b *builder) queryForms(query ast.Node) bool {
	var rng review.Range
	msg := ""
	set := func(loc ast.Loc, m string) {
		if msg == "" {
			rng, msg = rangeOf(loc), m
		}
	}
	roots := make(map[*ast.SelectStmt]bool)
	var root func(n ast.Node)
	root = func(n ast.Node) {
		s, ok := n.(*ast.SelectStmt)
		if !ok || s == nil {
			return
		}
		roots[s] = true
		root(s.Larg)
		root(s.Rarg)
		if s.WithClause != nil {
			for _, item := range listNodes(s.WithClause.CTEs) {
				if cte, ok := item.(*ast.CommonTableExpr); ok {
					root(cte.Query)
				}
			}
		}
	}
	root(query)
	ast.Inspect(query, func(m ast.Node) bool {
		if msg != "" {
			return false
		}
		switch v := m.(type) {
		case *ast.SelectStmt:
			if loc, why := blockForm(v, roots[v]); why != "" {
				set(loc, why)
			}
		case *ast.FuncCallExpr:
			if loc, why := callForm(v); why != "" {
				set(loc, why)
			}
		case *ast.AliasedTableRef:
			if loc, why := derivedForm(v); why != "" {
				set(loc, why)
			}
		}
		return true
	})
	if msg == "" {
		return false
	}
	b.fail(rng, review.UnsupportedForm, msg)
	return true
}

// blockForm checks one query block: its window functions (4112, 5362,
// 4108), its aggregates (147 in the WHERE of a block that has no outer
// query, 130), GROUP BY (164), ORDER BY (145, 108), its subqueries in
// expressions (116), the exposed names of FROM (1011, 1013), and the
// CTEs of its WITH (8158, 8159, 8156, 252).
func blockForm(s *ast.SelectStmt, root bool) (ast.Loc, string) {
	if s.WithClause != nil {
		for _, item := range listNodes(s.WithClause.CTEs) {
			cte, ok := item.(*ast.CommonTableExpr)
			if !ok {
				continue
			}
			q, _ := cte.Query.(*ast.SelectStmt)
			if loc, why := columnsForm(cte.Name, q, cte.Columns, cte.Loc); why != "" {
				return loc, why
			}
			if q != nil && q.Larg == nil && q.Rarg == nil && selfReference(q, cte.Name) {
				return cte.Loc, "SQL Server doesn't accept recursive CTE " + review.Code(ident(cte.Name)) + " without a top-level `UNION ALL`."
			}
		}
	}
	if s.Larg != nil || s.Rarg != nil {
		return ast.Loc{}, ""
	}
	var exprs []ast.Node
	exprs = append(exprs, listNodes(s.TargetList)...)
	for _, item := range listNodes(s.OrderByClause) {
		if o, ok := item.(*ast.OrderByItem); ok {
			exprs = append(exprs, o.Expr)
		}
	}
	var loc ast.Loc
	why := ""
	set := func(l ast.Loc, w string) {
		if why == "" {
			loc, why = l, w
		}
	}
	// Window functions name a window the block defines and order it.
	for _, e := range exprs {
		inBlock(e, func(m ast.Node) bool {
			f, ok := m.(*ast.FuncCallExpr)
			if !ok || f.Over == nil {
				return true
			}
			if l, w := windowForm(s, f); w != "" {
				set(l, w)
			}
			return true
		})
	}
	// A window function belongs only in the select list and ORDER BY,
	// and an aggregate not in WHERE.
	clauses := []ast.Node{s.WhereClause, s.HavingClause}
	clauses = append(clauses, listNodes(s.GroupByClause)...)
	for i, c := range clauses {
		inBlock(c, func(m ast.Node) bool {
			if f, ok := m.(*ast.FuncCallExpr); ok && f.Over != nil {
				set(f.Loc, "SQL Server doesn't accept a window function outside the select list and `ORDER BY`.")
				return false
			}
			if i == 0 && root && isAggregate(m) {
				set(m.(*ast.FuncCallExpr).Loc, "SQL Server doesn't accept an aggregate in `WHERE`.")
				return false
			}
			return true
		})
	}
	// An aggregate of an aggregate or a subquery.
	for _, e := range append(exprs, s.HavingClause) {
		inBlock(e, func(m ast.Node) bool {
			if !isAggregate(m) {
				return true
			}
			f := m.(*ast.FuncCallExpr)
			for _, a := range listNodes(f.Args) {
				ast.Inspect(a, func(n ast.Node) bool {
					switch n.(type) {
					case *ast.SubqueryExpr, *ast.SelectStmt:
						set(f.Loc, "SQL Server doesn't accept an aggregate of an aggregate or a subquery.")
						return false
					}
					if isAggregate(n) {
						set(f.Loc, "SQL Server doesn't accept an aggregate of an aggregate or a subquery.")
						return false
					}
					return true
				})
			}
			return false
		})
	}
	for _, item := range listNodes(s.GroupByClause) {
		if lit, ok := item.(*ast.Literal); ok {
			set(lit.Loc, "SQL Server doesn't accept a constant in `GROUP BY`.")
		}
	}
	if l, w := orderForm(s); w != "" {
		set(l, w)
	}
	// A subquery in an expression returns one column.
	roots := append([]ast.Node{}, exprs...)
	roots = append(roots, s.WhereClause, s.HavingClause)
	roots = append(roots, listNodes(s.GroupByClause)...)
	for _, e := range roots {
		if e == nil {
			continue
		}
		ast.Inspect(e, func(m ast.Node) bool {
			var q ast.Node
			switch v := m.(type) {
			case *ast.ExistsExpr:
				return false
			case *ast.SelectStmt:
				return false
			case *ast.InExpr:
				if sq, ok := v.Subquery.(*ast.SubqueryExpr); ok {
					q = sq.Query
				}
			case *ast.SubqueryComparisonExpr:
				q = v.Subquery
			case *ast.SubqueryExpr:
				q = v.Query
			}
			if q == nil {
				return true
			}
			if sel, ok := q.(*ast.SelectStmt); ok {
				if w, ok := selectWidth(sel); ok && w > 1 {
					set(sel.Loc, "SQL Server doesn't accept a subquery of "+count(w, "column")+" in an expression.")
				}
			}
			return false
		})
	}
	if l, w := exposedTwice(s); w != "" {
		set(l, w)
	}
	return loc, why
}

// windowForm checks a windowed call of a block: a window name the block
// defines (5362), and the ORDER BY a ranking or offset function needs
// (4112, 5366).
func windowForm(s *ast.SelectStmt, f *ast.FuncCallExpr) (ast.Loc, string) {
	over := f.Over
	var def *ast.WindowDef
	if over.WindowName != "" {
		var maybe []*ast.WindowDef
		for _, item := range listNodes(s.WindowClause) {
			if d, ok := item.(*ast.WindowDef); ok && !differ(d.Name, over.WindowName) {
				maybe = append(maybe, d)
			}
		}
		if len(maybe) == 0 {
			return over.Loc, "Window " + review.Code(ident(over.WindowName)) + " is not defined in the query."
		}
		if len(maybe) != 1 || maybe[0].Name != over.WindowName {
			return ast.Loc{}, ""
		}
		def = maybe[0]
	}
	name := builtin(f)
	if !orderedWindows[name] || over.OrderBy != nil && len(over.OrderBy.Items) > 0 {
		return ast.Loc{}, ""
	}
	if def != nil && (def.RefName != "" || def.OrderBy != nil && len(def.OrderBy.Items) > 0 || over.PartitionBy != nil || over.WindowFrame != nil) {
		return ast.Loc{}, ""
	}
	return over.Loc, "SQL Server's " + review.Code(name) + " needs an `OVER` clause with `ORDER BY`."
}

// orderForm checks a block's ORDER BY: under SELECT DISTINCT each item is
// in the select list (145), and a position is one of the select list's
// (108).
func orderForm(s *ast.SelectStmt) (ast.Loc, string) {
	if s.OrderByClause == nil || s.TargetList == nil {
		return ast.Loc{}, ""
	}
	width, known := selectWidth(s)
	var names []string
	for _, item := range s.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok {
			if rt.Name != "" {
				names = append(names, rt.Name)
			}
			if c, ok := rt.Val.(*ast.ColumnRef); ok {
				names = append(names, c.Column)
			}
		}
	}
	for _, item := range s.OrderByClause.Items {
		o, ok := item.(*ast.OrderByItem)
		if !ok {
			continue
		}
		switch e := o.Expr.(type) {
		case *ast.Literal:
			if e.Type == ast.LitInteger && known && (e.Ival < 1 || e.Ival > int64(width)) {
				return e.Loc, "SQL Server doesn't accept `ORDER BY` position " + strconv.FormatInt(e.Ival, 10) + " for a select list of " + count(width, "column") + "."
			}
		case *ast.ColumnRef:
			if s.Distinct && known && !isPseudoCol(e.Column) && !hasName(names, e.Column) {
				return e.Loc, "SQL Server doesn't accept `ORDER BY` column " + review.Code(ident(e.Column)) + ", which `SELECT DISTINCT` doesn't select."
			}
		}
	}
	return ast.Loc{}, ""
}

// exposedTwice finds two sources of a block's FROM that expose one name,
// spelled alike (1011, 1013).
func exposedTwice(s *ast.SelectStmt) (ast.Loc, string) {
	var names []string
	var source func(n ast.Node)
	source = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.TableRef:
			if v.Alias != "" {
				names = append(names, v.Alias)
			} else {
				names = append(names, v.Object)
			}
		case *ast.AliasedTableRef:
			switch {
			case v.Alias != "":
				names = append(names, v.Alias)
			default:
				if tr, ok := v.Table.(*ast.TableRef); ok {
					source(tr)
				}
			}
		case *ast.TableVarRef:
			if v.Alias != "" {
				names = append(names, v.Alias)
			} else {
				names = append(names, v.Name)
			}
		case *ast.JoinClause:
			source(v.Left)
			source(v.Right)
		}
	}
	for _, item := range listNodes(s.FromClause) {
		source(item)
	}
	if n := twice(names); n != "" {
		return s.Loc, "This query exposes " + review.Code(ident(n)) + " twice in its `FROM` clause."
	}
	return ast.Loc{}, ""
}

// derivedForm checks a derived table or a VALUES table of FROM: VALUES
// rows of one length (10709), a column list naming each column once
// (8156) and as many as the table has (8158, 8159), and a query without
// one whose columns have different names (8156).
func derivedForm(v *ast.AliasedTableRef) (ast.Loc, string) {
	switch t := v.Table.(type) {
	case *ast.ValuesClause:
		width := -1
		for _, row := range listNodes(t.Rows) {
			l, ok := row.(*ast.List)
			if !ok {
				return ast.Loc{}, ""
			}
			if width >= 0 && len(l.Items) != width {
				return t.Loc, "SQL Server doesn't accept `VALUES` rows of different lengths."
			}
			width = len(l.Items)
		}
		if width < 0 || v.Columns == nil {
			return ast.Loc{}, ""
		}
		return listForm(v.Alias, width, v.Columns, v.Loc)
	case *ast.SubqueryExpr:
		q, _ := t.Query.(*ast.SelectStmt)
		return columnsForm(v.Alias, q, v.Columns, v.Loc)
	}
	return ast.Loc{}, ""
}

// columnsForm checks the columns of a derived table or a CTE: its column
// list, or else the names its query gives.
func columnsForm(alias string, q *ast.SelectStmt, cols *ast.List, loc ast.Loc) (ast.Loc, string) {
	if q == nil {
		return ast.Loc{}, ""
	}
	if cols != nil {
		width, ok := selectWidth(q)
		if !ok {
			if names, ok := listNames(cols); ok {
				if n := twice(names); n != "" {
					return loc, msgTwice("names", "column", ident(n))
				}
			}
			return ast.Loc{}, ""
		}
		return listForm(alias, width, cols, loc)
	}
	if names, ok := outputNames(q); ok {
		if n := twice(names); n != "" {
			return loc, msgTwice("names", "column", ident(n))
		}
	}
	return ast.Loc{}, ""
}

// listForm checks a column list for a table of width columns.
func listForm(alias string, width int, cols *ast.List, loc ast.Loc) (ast.Loc, string) {
	names, ok := listNames(cols)
	if !ok {
		return ast.Loc{}, ""
	}
	if n := twice(names); n != "" {
		return loc, msgTwice("names", "column", ident(n))
	}
	if len(names) != width {
		return loc, "SQL Server doesn't accept " + review.Code(ident(alias)) + " with a column list of " + count(len(names), "name") + " for " + count(width, "column") + "."
	}
	return ast.Loc{}, ""
}

// listNames returns the names of a column list.
func listNames(cols *ast.List) ([]string, bool) {
	var names []string
	for _, item := range listNodes(cols) {
		s, ok := item.(*ast.String)
		if !ok {
			return nil, false
		}
		names = append(names, s.Str)
	}
	return names, true
}

// outputNames returns the names of the columns a query returns, which its
// first branch gives, false when one has no name the walk knows.
func outputNames(q *ast.SelectStmt) ([]string, bool) {
	for q != nil && q.Larg != nil {
		q = q.Larg
	}
	if q == nil || q.TargetList == nil || q.ForClause != nil {
		return nil, false
	}
	var names []string
	for _, item := range q.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			return nil, false
		}
		switch v := rt.Val.(type) {
		case *ast.StarExpr, *ast.SelectAssign:
			return nil, false
		case *ast.ColumnRef:
			if rt.Name == "" {
				if isPseudoCol(v.Column) {
					return nil, false
				}
				names = append(names, v.Column)
				continue
			}
		}
		if rt.Name == "" {
			return nil, false
		}
		names = append(names, rt.Name)
	}
	return names, true
}

// selfReference reports whether a query reads a table spelled as name,
// which in a CTE's query means the CTE.
func selfReference(q *ast.SelectStmt, name string) bool {
	found := false
	forEachFromTable(q, func(tr *ast.TableRef) {
		if tr.Schema == "" && tr.Database == "" && tr.Server == "" && tr.Object == name {
			found = true
		}
	})
	return found
}

// callForm checks a system function call: how many arguments it takes
// (174, 189), and the datepart it names (155).
func callForm(f *ast.FuncCallExpr) (ast.Loc, string) {
	name := builtin(f)
	if name == "" || f.Star {
		return ast.Loc{}, ""
	}
	args := listNodes(f.Args)
	if r, ok := builtinArity[name]; ok && f.TrimOption == "" && (len(args) < r[0] || len(args) > r[1]) {
		want := count(r[0], "argument")
		if r[0] != r[1] {
			want = strconv.Itoa(r[0]) + " to " + count(r[1], "argument")
		}
		return f.Loc, "SQL Server's " + review.Code(name) + " takes " + want + "."
	}
	for _, a := range args {
		if dp, ok := a.(*ast.DatePart); ok && !strings.HasPrefix(dp.Name, "[") && !dateParts[strings.ToLower(dp.Name)] {
			return dp.Loc, review.Code(dp.Name) + " is not a datepart of " + review.Code(name) + "."
		}
	}
	return ast.Loc{}, ""
}

// sumTypes are the system types SUM and AVG refuse (8117).
var sumTypes = map[string]bool{
	"char": true, "varchar": true, "nchar": true, "nvarchar": true, "text": true, "ntext": true, "sysname": true,
	"date": true, "time": true, "datetime": true, "datetime2": true, "datetimeoffset": true, "smalldatetime": true,
}

// callTypes fails a system function call of a bound view's or inline
// function's query on an argument whose type it refuses: TRIM a number
// (8116), SUM or AVG a string or a date (8117).
func (q *binder) callTypes(f *ast.FuncCallExpr, sc *qscope) {
	typ := func(e ast.ExprNode) string { return q.operandType(e, sc) }
	switch name := builtin(f); name {
	case "TRIM":
		for _, a := range listNodes(f.Args) {
			e, ok := a.(ast.ExprNode)
			if !ok {
				continue
			}
			if t := operandType(e, typ); integerTypes[t] {
				q.b.fail(rangeOf(f.Loc), review.UnsupportedForm, "SQL Server's `TRIM` doesn't accept an argument of type "+review.Code(t)+".")
				return
			}
		}
	case "SUM", "AVG":
		if f.Star || f.Args == nil || len(f.Args.Items) != 1 {
			return
		}
		e, ok := f.Args.Items[0].(ast.ExprNode)
		if !ok {
			return
		}
		if t := operandType(e, typ); sumTypes[t] {
			q.b.fail(rangeOf(f.Loc), review.UnsupportedForm, "SQL Server's "+review.Code(name)+" doesn't accept an argument of type "+review.Code(t)+".")
		}
	}
}

// integerTypes are the exact whole-number system types.
var integerTypes = map[string]bool{"tinyint": true, "smallint": true, "int": true, "bigint": true}

// clrTypes are the system types with properties and methods.
var clrTypes = map[string]bool{"xml": true, "geometry": true, "geography": true, "hierarchyid": true}

// cte returns the one CTE a table reference certainly means: of the WITH
// clauses around it, defined before it, and spelled as it is; nil when
// another CTE or a table may be meant.
func (q *binder) cte(tr *ast.TableRef) *ast.CommonTableExpr {
	if tr.Loc.Start < 0 {
		return nil
	}
	var found *ast.CommonTableExpr
	for _, c := range q.visible {
		if differ(c.Name, tr.Object) {
			continue
		}
		if c.Loc.Start < 0 || c.Loc.Start > tr.Loc.Start {
			// A CTE defined after the reference is not one it can mean.
			return nil
		}
		if c.Name != tr.Object || found != nil {
			return nil
		}
		found = c
	}
	return found
}

// cteNames returns the names of the columns of a derived table or a CTE:
// its column list when it is as long as the query is wide, else the names
// its query gives.
func cteNames(sel *ast.SelectStmt, cols *ast.List) ([]string, bool) {
	if cols != nil {
		names, ok := listNames(cols)
		if !ok {
			return nil, false
		}
		if w, known := selectWidth(sel); !known || w != len(names) {
			return nil, false
		}
		return names, true
	}
	if sel == nil {
		return nil, false
	}
	return outputNames(sel)
}

// derived adds a derived table to a block: one whose column names the
// walk knows lists them, else an opaque one.
func (q *binder) derived(alias string, sel *ast.SelectStmt, cols *ast.List, sc *qscope) {
	if names, ok := cteNames(sel, cols); ok {
		sc.sources = append(sc.sources, qsource{exposed: alias, written: alias, local: true, cols: names})
		return
	}
	sc.opaque = true
	sc.sources = append(sc.sources, qsource{exposed: alias})
}

// hasHint reports whether a table's hints hold one.
func hasHint(hints *ast.List, name string) bool {
	for _, item := range listNodes(hints) {
		if h, ok := item.(*ast.TableHint); ok && strings.EqualFold(h.Name, name) {
			return true
		}
	}
	return false
}

// rowsetFunctions are the system functions FROM calls by a one-part name.
var rowsetFunctions = map[string]bool{
	"OPENJSON": true, "OPENXML": true, "OPENROWSET": true, "OPENQUERY": true, "OPENDATASOURCE": true,
	"STRING_SPLIT": true, "GENERATE_SERIES": true, "CHANGETABLE": true, "CONTAINSTABLE": true,
	"FREETEXTTABLE": true, "SEMANTICKEYPHRASETABLE": true, "SEMANTICSIMILARITYTABLE": true,
	"SEMANTICSIMILARITYDETAILSTABLE": true, "PREDICT": true,
}

// tableFunction checks a table-valued function call of FROM: an object
// of another kind takes no parameters (215). A table with names in
// parentheses is one with old-style table hints instead, such as
// a(NOLOCK), which a name that is no hint fails otherwise (207).
func (b *builder) tableFunction(f *ast.FuncCallExpr) {
	if f.Name == nil || f.Name.Server != "" || f.Name.Database != "" {
		return
	}
	for _, a := range listNodes(f.Args) {
		if _, ok := a.(*ast.ColumnRef); ok {
			return
		}
	}
	if f.Name.Schema == "" {
		up := strings.ToUpper(f.Name.Object)
		if rowsetFunctions[up] || strings.HasPrefix(up, "FN_") || strings.HasPrefix(up, "DM_") {
			return
		}
	}
	r, ok := b.ref(f.Name, false)
	if !ok || r.where != placed || r.temp {
		return
	}
	present := notKind(r.written, "function", kindFunction, kindSynonym)
	b.ofKind(spaceObject, r.parent, r.n, []string{kindFunction, kindSynonym}, note(r.rng, "", present))
	b.ofKind(spaceRelation, r.parent, r.n, []string{kindFunction, kindSynonym}, note(r.rng, "", present))
}

// noSources reports whether no block of a scope has a source.
func noSources(sc *qscope) bool {
	for s := sc; s != nil; s = s.outer {
		if s.opaque || len(s.sources) > 0 {
			return false
		}
	}
	return true
}

// pivot binds a PIVOT on a source of inner: its FOR column and the
// columns it aggregates are the source's (207), and the result has the
// source's other columns and one for each IN value.
func (q *binder) pivot(v *ast.PivotExpr, inner, sc *qscope) {
	var used []*ast.ColumnRef
	if v.ForCol != nil {
		used = append(used, v.ForCol)
	}
	if f, ok := v.AggFunc.(*ast.FuncCallExpr); ok {
		for _, a := range listNodes(f.Args) {
			if c, ok := a.(*ast.ColumnRef); ok {
				used = append(used, c)
			}
		}
	}
	for _, c := range used {
		if c.Table == "" {
			q.column(c, inner)
		}
	}
	src := qsource{exposed: v.Alias, written: v.Alias}
	if len(inner.sources) != 1 || inner.opaque || !inner.sources[0].local {
		sc.opaque = true
		sc.sources = append(sc.sources, src)
		return
	}
	for _, name := range inner.sources[0].cols {
		taken := false
		for _, c := range used {
			if c.Column == name {
				taken = true
			}
		}
		if !taken {
			src.cols = append(src.cols, name)
		}
	}
	for _, item := range listNodes(v.InValues) {
		c, ok := item.(*ast.ColumnRef)
		if !ok {
			sc.opaque = true
			sc.sources = append(sc.sources, qsource{exposed: v.Alias})
			return
		}
		src.cols = append(src.cols, c.Column)
	}
	src.local = true
	sc.sources = append(sc.sources, src)
}

// unpivot binds an UNPIVOT on a source of inner: the columns it turns
// into rows are the source's (207), and the result has the source's
// other columns, the value column, and the FOR column.
func (q *binder) unpivot(v *ast.UnpivotExpr, inner, sc *qscope) {
	var in []*ast.ColumnRef
	for _, item := range listNodes(v.InCols) {
		if c, ok := item.(*ast.ColumnRef); ok {
			in = append(in, c)
			if c.Table == "" {
				q.column(c, inner)
			}
		}
	}
	src := qsource{exposed: v.Alias, written: v.Alias}
	if len(inner.sources) != 1 || inner.opaque || !inner.sources[0].local || v.ValueCol == nil || v.ForCol == nil {
		sc.opaque = true
		sc.sources = append(sc.sources, src)
		return
	}
	for _, name := range inner.sources[0].cols {
		taken := false
		for _, c := range in {
			if c.Column == name {
				taken = true
			}
		}
		if !taken {
			src.cols = append(src.cols, name)
		}
	}
	src.cols = append(src.cols, v.ValueCol.Column, v.ForCol.Column)
	src.local = true
	sc.sources = append(sc.sources, src)
}
