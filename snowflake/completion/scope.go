package completion

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// maxOutputDepth bounds how deep the output columns of nested derived
// tables and CTEs are read.
const maxOutputDepth = 8

// flattenColumns are the columns FLATTEN returns
// (https://docs.snowflake.com/en/sql-reference/functions/flatten#output).
var flattenColumns = []string{"SEQ", "KEY", "PATH", "INDEX", "VALUE", "THIS"}

// rangeItem is a relation in scope: a table, a derived table, a CTE
// reference or a table function.
type rangeItem struct {
	// names are the resolved qualifiers that address the item: its alias,
	// or its table name with each of its qualifications.
	names [][]string
	// display is the qualifier offered for it, as written.
	display string
	columns func() []colInfo
}

type colInfo struct {
	name, definition, comment string
}

// addressedBy reports whether qualifier names the item.
func (r *rangeItem) addressedBy(qualifier []string) bool {
	return slices.ContainsFunc(r.names, func(n []string) bool { return slices.Equal(n, qualifier) })
}

// scope is what is visible at the caret.
type scope struct {
	// levels hold the relations in scope, the innermost query's first, so an
	// inner alias hides an outer one of the same name.
	levels [][]*rangeItem
	// aliases are the select-list aliases the caret's clause may use.
	aliases []string
	// ctes are the CTE names visible at the caret.
	ctes []string
}

// cteMap maps a visible CTE's resolved name to its definition. An inner WITH
// shadows an outer one.
type cteMap map[string]*cteDef

// cteDef is a CTE with what its body sees: the CTEs visible where its WITH
// is, and the CTEs of the WITH before it, or through it when recursive.
type cteDef struct {
	cte      *ast.CTE
	outer    cteMap
	siblings []*ast.CTE
	index    int
}

func (d *cteDef) bodyCTEs() cteMap {
	n := d.index
	if d.cte.Recursive {
		n++
	}
	return d.outer.with(d.siblings[:n])
}

func (m cteMap) with(ctes []*ast.CTE) cteMap {
	if len(ctes) == 0 {
		return m
	}
	out := make(cteMap, len(m)+len(ctes))
	for k, v := range m {
		out[k] = v
	}
	for i, c := range ctes {
		out[c.Name.Normalize()] = &cteDef{cte: c, outer: m, siblings: ctes, index: i}
	}
	return out
}

// scopeBuilder builds rangeItems against a catalog. It reads each CTE's
// columns once, however often the statement references it. query is set in
// a query, where unqualified names resolve through the search path.
type scopeBuilder struct {
	cat     *catalog
	query   bool
	cteCols map[*ast.CTE][]colInfo
}

func newScopeBuilder(cat *catalog, query bool) *scopeBuilder {
	return &scopeBuilder{cat: cat, query: query, cteCols: map[*ast.CTE][]colInfo{}}
}

// ctesAlong returns the CTEs visible at each node of path, the path from
// the statement to the node holding the placeholder. Inside the body of a
// WITH's k-th CTE only the CTEs before it are visible, and the k-th itself
// when it is recursive.
func ctesAlong(path []ast.Node) []cteMap {
	ctesAt := make([]cteMap, len(path))
	visible := cteMap{}
	for i, n := range path {
		// The parser keeps a set operation's WITH on its leftmost operand,
		// and every operand sees it.
		if setOp, ok := n.(*ast.SetOperationStmt); ok {
			if left := leftmostSelect(setOp); left != nil && !slices.Contains(path[i+1:], ast.Node(left)) {
				visible = visible.with(left.With)
			}
		}
		if sel, ok := n.(*ast.SelectStmt); ok && len(sel.With) > 0 {
			withVisible := sel.With
			if i+1 < len(path) {
				for k, c := range sel.With {
					if c.Query == path[i+1] {
						withVisible = sel.With[:k]
						if c.Recursive {
							withVisible = sel.With[:k+1]
						}
						break
					}
				}
			}
			visible = visible.with(withVisible)
		}
		ctesAt[i] = visible
	}
	return ctesAt
}

// leftmostSelect returns the first operand of a set operation.
func leftmostSelect(n ast.Node) *ast.SelectStmt {
	for {
		switch q := n.(type) {
		case *ast.SelectStmt:
			return q
		case *ast.SetOperationStmt:
			n = q.Left
		default:
			return nil
		}
	}
}

// scopeAt returns the scope at the end of path, given the CTEs visible
// along it.
func (b *scopeBuilder) scopeAt(path []ast.Node, ctesAt []cteMap) scope {
	var sc scope
	for name := range ctesAt[len(path)-1] {
		sc.ctes = append(sc.ctes, name)
	}
	slices.Sort(sc.ctes)
	for i := len(path) - 1; i >= 0; i-- {
		var child ast.Node
		if i+1 < len(path) {
			child = path[i+1]
		}
		ctes := ctesAt[i]
		switch n := path[i].(type) {
		case *ast.SelectStmt:
			if slices.ContainsFunc(n.With, func(c *ast.CTE) bool { return c.Query == child }) {
				// A CTE body sees no FROM item of the query it is defined in.
				continue
			}
			if slices.Contains(n.From, child) {
				sc.levels = append(sc.levels, b.fromVisible(n.From, path[i+1:], ctes))
				continue
			}
			sc.levels = append(sc.levels, b.leavesOf(n.From, ctes))
			if !slices.ContainsFunc(n.Targets, func(t *ast.SelectTarget) bool { return t.Expr == child }) {
				for _, t := range n.Targets {
					if !t.Alias.IsEmpty() {
						sc.aliases = append(sc.aliases, t.Alias.Normalize())
					}
				}
			}
		case *ast.UpdateStmt:
			level := []*rangeItem{b.relation(n.Target, n.Alias, ctes, 0)}
			sc.levels = append(sc.levels, append(level, b.leavesOf(n.From, ctes)...))
		case *ast.DeleteStmt:
			level := []*rangeItem{b.relation(n.Target, n.Alias, ctes, 0)}
			sc.levels = append(sc.levels, append(level, b.leavesOf(n.Using, ctes)...))
		case *ast.MergeStmt:
			if child == n.Source {
				continue
			}
			level := []*rangeItem{b.relation(n.Target, n.TargetAlias, ctes, 0)}
			if src := b.leaves(n.Source, ctes); len(src) == 1 && !n.SourceAlias.IsEmpty() {
				src[0].names = [][]string{{n.SourceAlias.Normalize()}}
				src[0].display = n.SourceAlias.String()
				level = append(level, src...)
			} else {
				level = append(level, src...)
			}
			sc.levels = append(sc.levels, level)
		}
	}
	return sc
}

// fromVisible returns the FROM items visible inside the FROM item that
// rest[0] is, rest being the path from that item down to the placeholder.
func (b *scopeBuilder) fromVisible(from []ast.Node, rest []ast.Node, ctes cteMap) []*rangeItem {
	var preceding []*rangeItem
	for _, item := range from {
		if item == rest[0] {
			return b.visibleIn(rest, preceding, ctes)
		}
		preceding = append(preceding, b.leaves(item, ctes)...)
	}
	return nil
}

// visibleIn returns what is visible inside rest[0], a FROM item, from the
// path below it: an ON condition sees its join's operands, a LATERAL
// subquery or table function sees the items before it, and a clause
// attached to a table (PIVOT, MATCH_RECOGNIZE, ...) sees that table.
func (b *scopeBuilder) visibleIn(rest []ast.Node, preceding []*rangeItem, ctes cteMap) []*rangeItem {
	if len(rest) < 2 {
		return nil
	}
	next := rest[1]
	switch it := rest[0].(type) {
	case *ast.JoinExpr:
		switch next {
		case it.Left:
			return b.visibleIn(rest[1:], preceding, ctes)
		case it.Right:
			return b.visibleIn(rest[1:], append(preceding, b.leaves(it.Left, ctes)...), ctes)
		}
		return b.leaves(it, ctes)
	case *ast.TableRef:
		switch {
		case it.Name != nil && next == ast.Node(it.Name):
			return nil
		case next == it.Subquery:
			if it.Lateral {
				return preceding
			}
			return nil
		case it.FuncCall != nil && next == ast.Node(it.FuncCall):
			return preceding
		case it.Nested != nil && next == ast.Node(it.Nested):
			return b.visibleIn(rest[1:], preceding, ctes)
		}
		return []*rangeItem{b.clauseSource(it, ctes, 0)}
	}
	return nil
}

func (b *scopeBuilder) leavesOf(items []ast.Node, ctes cteMap) []*rangeItem {
	var out []*rangeItem
	for _, item := range items {
		out = append(out, b.leaves(item, ctes)...)
	}
	return out
}

func (b *scopeBuilder) leaves(item ast.Node, ctes cteMap) []*rangeItem {
	switch it := item.(type) {
	case *ast.JoinExpr:
		return append(b.leaves(it.Left, ctes), b.leaves(it.Right, ctes)...)
	case *ast.TableRef:
		return []*rangeItem{b.tableRef(it, ctes, 0)}
	}
	return nil
}

// relation returns the item for a table name, addressed by its alias when it
// has one, and otherwise by its name with each qualification Snowflake
// accepts: orders in FROM orders answers to orders.c, public.orders.c and
// db.public.orders.c when public is the session's schema.
func (b *scopeBuilder) relation(name *ast.ObjectName, alias ast.Ident, ctes cteMap, depth int) *rangeItem {
	item := &rangeItem{columns: func() []colInfo { return nil }}
	if name == nil {
		return item
	}
	parts := normalizeParts(name.Parts())
	def, isCTE := ctes[parts[0]]
	isCTE = isCTE && len(parts) == 1
	switch {
	case !alias.IsEmpty():
		item.names = [][]string{{alias.Normalize()}}
		item.display = alias.String()
	case isCTE:
		item.names = [][]string{parts}
		item.display = name.Name.String()
	default:
		full := b.cat.qualify(parts, b.query)
		for i := range full {
			item.names = append(item.names, full[i:])
		}
		item.display = name.Name.String()
	}
	if isCTE {
		item.columns = func() []colInfo { return b.cteColumns(def) }
		return item
	}
	cat, full := b.cat, b.cat.qualify(parts, b.query)
	item.columns = memo(func() []colInfo {
		r := cat.relationAt(full)
		if r == nil {
			return nil
		}
		out := make([]colInfo, 0, len(r.columns))
		for _, col := range r.columns {
			out = append(out, colInfo{name: col.name, definition: r.qualifiedName() + " | " + col.typ, comment: col.comment})
		}
		return out
	})
	return item
}

// tableRef returns the item for a FROM table reference.
func (b *scopeBuilder) tableRef(t *ast.TableRef, ctes cteMap, depth int) *rangeItem {
	if t.Pivot != nil || t.Unpivot != nil || t.MatchRecognize != nil {
		return b.reshaped(t, ctes, depth)
	}
	if t.Name != nil {
		item := b.relation(t.Name, t.Alias, ctes, depth)
		if len(t.Columns) > 0 {
			item.columns = namedColumns(t.Columns)
		}
		return item
	}
	item := &rangeItem{columns: func() []colInfo { return nil }}
	if !t.Alias.IsEmpty() {
		item.names = [][]string{{t.Alias.Normalize()}}
		item.display = t.Alias.String()
	}
	switch {
	case len(t.Columns) > 0:
		item.columns = namedColumns(t.Columns)
	case t.Subquery != nil:
		item.columns = memo(func() []colInfo { return b.outputs(t.Subquery, ctes, depth+1) })
	case t.FuncCall != nil && t.FuncCall.Name.Normalize() == "FLATTEN":
		item.columns = func() []colInfo {
			out := make([]colInfo, len(flattenColumns))
			for i, c := range flattenColumns {
				out[i] = colInfo{name: c, definition: "FLATTEN"}
			}
			return out
		}
	}
	return item
}

// clauseSource returns the relation a table reference's PIVOT, UNPIVOT or
// MATCH_RECOGNIZE reads: the previous link of a chain, or the reference
// without those clauses.
func (b *scopeBuilder) clauseSource(t *ast.TableRef, ctes cteMap, depth int) *rangeItem {
	if t.Nested != nil {
		return b.tableRef(t.Nested, ctes, depth+1)
	}
	stripped := *t
	stripped.Pivot, stripped.Unpivot, stripped.MatchRecognize = nil, nil, nil
	return b.tableRef(&stripped, ctes, depth)
}

// reshaped returns the item for a table reference a MATCH_RECOGNIZE, PIVOT
// or UNPIVOT reshapes. The parser keeps the trailing alias on the last
// clause, which addresses the result; without one the source's names do.
func (b *scopeBuilder) reshaped(t *ast.TableRef, ctes cteMap, depth int) *rangeItem {
	source := b.clauseSource(t, ctes, depth)
	item := &rangeItem{names: source.names, display: source.display}
	var alias ast.Ident
	switch {
	case t.Pivot != nil && !t.Pivot.Alias.IsEmpty():
		alias = t.Pivot.Alias
	case t.Unpivot != nil && !t.Unpivot.Alias.IsEmpty():
		alias = t.Unpivot.Alias
	case t.MatchRecognize != nil && !t.MatchRecognize.Alias.IsEmpty():
		alias = t.MatchRecognize.Alias
	case !t.Alias.IsEmpty():
		alias = t.Alias
	}
	if !alias.IsEmpty() {
		item.names = [][]string{{alias.Normalize()}}
		item.display = alias.String()
	}
	item.columns = memo(func() []colInfo {
		cols := source.columns()
		if t.MatchRecognize != nil {
			cols = matchRecognizeColumns(t.MatchRecognize, cols)
		}
		if t.Pivot != nil {
			cols = pivotColumns(t.Pivot, cols)
		}
		if t.Unpivot != nil {
			cols = unpivotColumns(t.Unpivot, cols)
		}
		return cols
	})
	return item
}

// matchRecognizeColumns returns MATCH_RECOGNIZE's output: the partition
// columns and the measures for ONE ROW PER MATCH, every input column and the
// measures for ALL ROWS PER MATCH.
func matchRecognizeColumns(mr *ast.MatchRecognizeClause, in []colInfo) []colInfo {
	var out []colInfo
	if mr.RowsPerMatch != nil && mr.RowsPerMatch.Kind == ast.AllRowsPerMatch {
		out = append(out, in...)
	} else {
		for _, p := range mr.PartitionBy {
			if ref, ok := p.(*ast.ColumnRef); ok && len(ref.Parts) > 0 {
				out = append(out, colInfo{name: ref.Parts[len(ref.Parts)-1].Normalize()})
			}
		}
	}
	for _, m := range mr.Measures {
		if !m.Alias.IsEmpty() {
			out = append(out, colInfo{name: m.Alias.Normalize()})
		}
	}
	return out
}

// pivotColumns returns PIVOT's output: the input columns but the pivot
// column and the aggregated ones, then a column per IN value, named by its
// alias or, unaliased, by its literal as written ('Q1'). An IN of ANY or a
// subquery names its columns only at run time.
func pivotColumns(p *ast.PivotClause, in []colInfo) []colInfo {
	consumed := map[string]bool{}
	if p.ForColumn != nil && len(p.ForColumn.Parts) > 0 {
		consumed[p.ForColumn.Parts[len(p.ForColumn.Parts)-1].Normalize()] = true
	}
	if p.Agg != nil {
		for _, arg := range p.Agg.Args {
			if ref, ok := arg.(*ast.ColumnRef); ok && len(ref.Parts) > 0 {
				consumed[ref.Parts[len(ref.Parts)-1].Normalize()] = true
			}
		}
	}
	out := slices.DeleteFunc(slices.Clone(in), func(c colInfo) bool { return consumed[c.name] })
	if p.In != nil && p.In.Kind == ast.PivotInValues {
		for _, v := range p.In.Values {
			switch {
			case !v.Alias.IsEmpty():
				out = append(out, colInfo{name: v.Alias.Normalize()})
			default:
				if lit, ok := v.Value.(*ast.Literal); ok {
					name := lit.Value
					if lit.Kind == ast.LitString {
						name = "'" + name + "'"
					}
					out = append(out, colInfo{name: name})
				}
			}
		}
	}
	return out
}

// unpivotColumns returns UNPIVOT's output: the input columns but the
// unpivoted ones, then the name and value columns.
func unpivotColumns(u *ast.UnpivotClause, in []colInfo) []colInfo {
	out := slices.DeleteFunc(slices.Clone(in), func(c colInfo) bool {
		return slices.ContainsFunc(u.Columns, func(uc *ast.UnpivotColumn) bool { return uc.Column.Normalize() == c.name })
	})
	return append(out, colInfo{name: u.NameColumn.Normalize()}, colInfo{name: u.ValueColumn.Normalize()})
}

func namedColumns(ids []ast.Ident) func() []colInfo {
	return func() []colInfo {
		out := make([]colInfo, len(ids))
		for i, id := range ids {
			out[i] = colInfo{name: id.Normalize()}
		}
		return out
	}
}

// cteColumns returns a CTE's columns: its column list, or its body's
// outputs read in the scope of its definition. Each CTE is read once, so the
// depth bound starts again at its body; a recursive CTE read while its own
// body is read has no columns yet.
func (b *scopeBuilder) cteColumns(def *cteDef) []colInfo {
	c := def.cte
	if cols, ok := b.cteCols[c]; ok {
		return cols
	}
	if len(c.Columns) > 0 {
		b.cteCols[c] = namedColumns(c.Columns)()
		return b.cteCols[c]
	}
	b.cteCols[c] = nil
	cols := b.outputs(c.Query, def.bodyCTEs(), 0)
	b.cteCols[c] = cols
	return cols
}

// outputs returns the output columns of a query: each select-list item's
// alias or column name, with * expanded from the query's FROM.
func (b *scopeBuilder) outputs(n ast.Node, ctes cteMap, depth int) []colInfo {
	if depth > maxOutputDepth {
		return nil
	}
	switch q := n.(type) {
	case *ast.SetOperationStmt:
		left := b.outputs(q.Left, ctes, depth+1)
		if !q.ByName {
			return left
		}
		// UNION [ALL] BY NAME aligns the operands by name: the left's
		// columns, then the right's it lacks.
		out := slices.Clone(left)
		for _, c := range b.outputs(q.Right, ctes, depth+1) {
			if !slices.ContainsFunc(out, func(o colInfo) bool { return o.name == c.name }) {
				out = append(out, c)
			}
		}
		return out
	case *ast.ValuesClause:
		if len(q.Rows) == 0 {
			return nil
		}
		out := make([]colInfo, len(q.Rows[0]))
		for i := range out {
			out[i] = colInfo{name: "COLUMN" + strconv.Itoa(i+1)}
		}
		return out
	case *ast.SelectStmt:
		ctes = ctes.with(q.With)
		var out []colInfo
		for _, t := range q.Targets {
			switch {
			case !t.Alias.IsEmpty():
				out = append(out, colInfo{name: t.Alias.Normalize()})
			case t.Star:
				out = append(out, b.starColumns(q, t, ctes, depth)...)
			default:
				if ref, ok := t.Expr.(*ast.ColumnRef); ok && len(ref.Parts) > 0 {
					out = append(out, colInfo{name: ref.Parts[len(ref.Parts)-1].Normalize()})
				}
			}
		}
		return out
	}
	return nil
}

// starColumns expands * or t.* of a select list.
func (b *scopeBuilder) starColumns(q *ast.SelectStmt, t *ast.SelectTarget, ctes cteMap, depth int) []colInfo {
	var qualifier []string
	if star, ok := t.Expr.(*ast.StarExpr); ok && star.Qualifier != nil {
		qualifier = normalizeParts(star.Qualifier.Parts())
	}
	var out []colInfo
	match := func(string) bool { return true }
	if t.Ilike != nil {
		match = ilikeMatcher(t.Ilike.Value)
	}
	for _, item := range b.leavesOfDepth(q.From, ctes, depth) {
		if qualifier != nil && !item.addressedBy(qualifier) {
			continue
		}
		for _, c := range item.columns() {
			if !match(c.name) {
				continue
			}
			if slices.ContainsFunc(t.Exclude, func(id ast.Ident) bool { return id.Normalize() == c.name }) {
				continue
			}
			for _, r := range t.Rename {
				if r.Col.Normalize() == c.name {
					c.name = r.Alias.Normalize()
				}
			}
			out = append(out, colInfo{name: c.name})
		}
	}
	return out
}

func (b *scopeBuilder) leavesOfDepth(items []ast.Node, ctes cteMap, depth int) []*rangeItem {
	var out []*rangeItem
	var walk func(ast.Node)
	walk = func(item ast.Node) {
		switch it := item.(type) {
		case *ast.JoinExpr:
			walk(it.Left)
			walk(it.Right)
		case *ast.TableRef:
			out = append(out, b.tableRef(it, ctes, depth))
		}
	}
	for _, item := range items {
		walk(item)
	}
	return out
}

// ilikeMatcher compiles a Snowflake ILIKE pattern into a case-insensitive
// match of a whole name: % matches any run, _ any one character. A pattern
// that fails to compile matches nothing.
func ilikeMatcher(pattern string) func(string) bool {
	var sb strings.Builder
	sb.WriteString(`(?is)\A`)
	for _, r := range pattern {
		switch r {
		case '%':
			sb.WriteString(`.*`)
		case '_':
			sb.WriteString(`.`)
		default:
			sb.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	sb.WriteString(`\z`)
	re, err := regexp.Compile(sb.String())
	if err != nil {
		return func(string) bool { return false }
	}
	return re.MatchString
}

// memo returns f, calling it at most once.
func memo(f func() []colInfo) func() []colInfo {
	var done bool
	var cols []colInfo
	return func() []colInfo {
		if !done {
			cols, done = f(), true
		}
		return cols
	}
}
