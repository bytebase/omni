package deparse

import (
	"sort"
	"strings"

	ast "github.com/bytebase/omni/mysql/ast"
	scopepkg "github.com/bytebase/omni/mysql/scope"
)

// fromScope is the scope of one query block: the tables of its FROM clause,
// and the columns its NATURAL and USING joins join on.
//
// MySQL lays out the result columns of a NATURAL or USING join
// (store_natural_using_join_columns, sql/sql_base.cc) as the columns it joins
// on, in the order its left input lists them, then the other columns of its
// left input, then the other columns of its right input. Within the inputs of
// such a join, every RIGHT join is a LEFT join of its inputs swapped, so it
// lists its right input first. Elsewhere a join lists its left input's
// columns, then its right input's, as written. A column a join joins on is the
// same column in every place that names it: in the join's ON condition, in
// SELECT *, and in an unqualified reference.
//
// Verified on MySQL 8.0 with SHOW CREATE VIEW: SELECT * FROM t1 RIGHT JOIN t2
// USING (k1) lists t2.k1 first and writes on((t2.k1 = t1.k1)); USING (b, k2,
// k1) over t5(a, k1, b, k2) compares k1, b, k2 in that order; and in
// u JOIN ((q RIGHT JOIN p ON 1=1) JOIN s ON 1=1) USING (k), p's columns come
// before q's.
type fromScope struct {
	*scopepkg.Scope
	from []ast.TableExpr

	// tables holds the scope entries of the FROM clause's tables and derived
	// tables in the order they were registered; entries maps each such FROM
	// item to its index there.
	tables  []scopepkg.Entry
	entries map[ast.TableExpr]int

	// spans maps each FROM item to the range [lo, hi) of tables under it.
	spans map[ast.TableExpr][2]int

	// joined maps each NATURAL or USING join resolved here to the columns it
	// joins on, taken from its left input, in that input's order.
	joined map[*ast.JoinClause][]joinColumn

	// hidden holds the right input's copy of each column a join joins on.
	hidden map[colKey]bool

	// coalesced maps a lowered column name to the column an unqualified
	// reference to that name means, when a NATURAL or USING join joins on it.
	coalesced map[string]joinColumn

	// firsts caches, for an input of a NATURAL or USING join, each lowered
	// column name to the first column of that name in the input's result.
	firsts map[ast.TableExpr]map[string]joinColumn
}

// joinColumn is one column of a FROM item's result.
type joinColumn struct {
	table     int    // index of its table in fromScope.tables
	qualifier string // the name the scope knows its table by
	name      string // the column's name in its table
	pos       int    // the column's position in its table
}

// colKey identifies a column of a table.
type colKey struct {
	table int
	name  string
}

func (c joinColumn) key() colKey {
	return colKey{table: c.table, name: strings.ToLower(c.name)}
}

func newFromScope(from []ast.TableExpr) *fromScope {
	return &fromScope{
		Scope:   scopepkg.New(),
		from:    from,
		entries: make(map[ast.TableExpr]int),
		spans:   make(map[ast.TableExpr][2]int),
	}
}

// add registers a table or derived table of the FROM clause.
func (sc *fromScope) add(te ast.TableExpr, name string, table *ResolverTable) {
	sc.Add(name, table)
	sc.entries[te] = len(sc.tables)
	sc.tables = append(sc.tables, scopepkg.Entry{Name: name, Table: table})
}

// span records the tables registered for a FROM item since lo.
func (sc *fromScope) span(te ast.TableExpr, lo int) {
	sc.spans[te] = [2]int{lo, len(sc.tables)}
}

// within reports whether c belongs to a table under te.
func (sc *fromScope) within(te ast.TableExpr, c joinColumn) bool {
	s := sc.spans[te]
	return s[0] <= c.table && c.table < s[1]
}

func isRightJoin(t ast.JoinType) bool {
	return t == ast.JoinRight || t == ast.JoinNaturalRight
}

// inputs returns a join's inputs in the order it lists their columns. A
// NATURAL or USING join, and a RIGHT join within the inputs of one
// (natural), lists a RIGHT join's right input first.
func (sc *fromScope) inputs(j *ast.JoinClause, natural bool) (ast.TableExpr, ast.TableExpr) {
	if _, ok := sc.joined[j]; ok {
		natural = true
	}
	if natural && isRightJoin(j.Type) {
		return j.Right, j.Left
	}
	return j.Left, j.Right
}

// tableColumns returns the columns of a registered table, in table order.
func (sc *fromScope) tableColumns(te ast.TableExpr) ([]joinColumn, bool) {
	i, ok := sc.entries[te]
	if !ok {
		return nil, false
	}
	e := sc.tables[i]
	cols := sortedResolverColumns(e.Table)
	out := make([]joinColumn, len(cols))
	for k, c := range cols {
		out[k] = joinColumn{table: i, qualifier: e.Name, name: c.Name, pos: c.Position}
	}
	return out, true
}

// firstColumns returns, for each lowered column name in the result of an
// input of a NATURAL or USING join, the first column of that name. A join's
// map is built from its inputs' maps, adding the smaller one's names to the
// larger, so a chain of joins builds its maps in linear time; an input's map
// is handed to its join and dropped from the cache, and rebuilt if asked for
// again.
func (sc *fromScope) firstColumns(te ast.TableExpr) (map[string]joinColumn, bool) {
	if m, ok := sc.firsts[te]; ok {
		return m, true
	}
	var m map[string]joinColumn
	if j, ok := te.(*ast.JoinClause); ok {
		left, right := sc.inputs(j, true)
		lm, ok := sc.firstColumns(left)
		if !ok {
			return nil, false
		}
		rm, ok := sc.firstColumns(right)
		if !ok {
			return nil, false
		}
		delete(sc.firsts, left)
		delete(sc.firsts, right)
		if len(lm) >= len(rm) {
			for k, c := range rm {
				if _, ok := lm[k]; !ok {
					lm[k] = c
				}
			}
			m = lm
		} else {
			for k, c := range lm {
				rm[k] = c
			}
			m = rm
		}
	} else {
		cols, ok := sc.tableColumns(te)
		if !ok {
			return nil, false
		}
		m = make(map[string]joinColumn, len(cols))
		for _, c := range cols {
			if k := strings.ToLower(c.name); !hasColumn(m, k) {
				m[k] = c
			}
		}
	}
	if sc.firsts == nil {
		sc.firsts = make(map[ast.TableExpr]map[string]joinColumn)
	}
	sc.firsts[te] = m
	return m, true
}

func hasColumn(m map[string]joinColumn, k string) bool {
	_, ok := m[k]
	return ok
}

// columns appends a FROM item's result columns, in MySQL's order, to out.
// natural is set within the inputs of a NATURAL or USING join. A column that
// joins join on is listed once, where its outermost such join lists it;
// listed records the columns already placed.
func (sc *fromScope) columns(te ast.TableExpr, natural bool, listed map[colKey]bool, out []joinColumn) []joinColumn {
	j, ok := te.(*ast.JoinClause)
	if !ok {
		cols, _ := sc.tableColumns(te)
		for _, c := range cols {
			if k := c.key(); !listed[k] && !sc.hidden[k] {
				out = append(out, c)
			}
		}
		return out
	}
	for _, c := range sc.joined[j] {
		// An outer join that joins on this column too may have listed it
		// already, or hidden it as the copy from its right input.
		if k := c.key(); !listed[k] && !sc.hidden[k] {
			listed[k] = true
			out = append(out, c)
		}
	}
	if _, ok := sc.joined[j]; ok {
		natural = true
	}
	left, right := sc.inputs(j, natural)
	out = sc.columns(left, natural, listed, out)
	return sc.columns(right, natural, listed, out)
}

// precedes reports whether column x comes before column y in the result of
// te, an input of a NATURAL or USING join. It walks down from te only until
// the two columns part, so ordering a join's few columns does not list its
// whole input.
func (sc *fromScope) precedes(te ast.TableExpr, x, y joinColumn) bool {
	for {
		j, ok := te.(*ast.JoinClause)
		if !ok {
			return x.pos < y.pos
		}
		if pulled := sc.joined[j]; len(pulled) > 0 {
			ix, iy := indexOfColumn(pulled, x), indexOfColumn(pulled, y)
			switch {
			case ix >= 0 && iy >= 0:
				return ix < iy
			case ix >= 0:
				return true
			case iy >= 0:
				return false
			}
		}
		left, right := sc.inputs(j, true)
		xl, yl := sc.within(left, x), sc.within(left, y)
		if xl != yl {
			return xl
		}
		if xl {
			te = left
		} else {
			te = right
		}
	}
}

func indexOfColumn(cols []joinColumn, c joinColumn) int {
	k := c.key()
	for i, x := range cols {
		if x.key() == k {
			return i
		}
	}
	return -1
}

// joinOn records that j joins on the named columns, and returns the ON
// condition MySQL writes for it: each column of its first input, in that
// input's order, equal to the same-named column of its second input. ok is
// false when a column cannot be found in both inputs, as when a table is
// unknown.
func (sc *fromScope) joinOn(j *ast.JoinClause, names []string) (ast.ExprNode, bool) {
	left, right := j.Left, j.Right
	if isRightJoin(j.Type) {
		left, right = right, left
	}
	lm, ok := sc.firstColumns(left)
	if !ok {
		return nil, false
	}
	rm, ok := sc.firstColumns(right)
	if !ok {
		return nil, false
	}
	type pair struct{ l, r joinColumn }
	pairs := make([]pair, 0, len(names))
	for _, name := range names {
		k := strings.ToLower(name)
		lc, lok := lm[k]
		rc, rok := rm[k]
		if !lok || !rok {
			return nil, false
		}
		pairs = append(pairs, pair{lc, rc})
	}
	sort.SliceStable(pairs, func(a, b int) bool {
		return sc.precedes(left, pairs[a].l, pairs[b].l)
	})
	if sc.joined == nil {
		sc.joined = make(map[*ast.JoinClause][]joinColumn)
		sc.hidden = make(map[colKey]bool)
		sc.coalesced = make(map[string]joinColumn)
	}
	lefts := make([]joinColumn, len(pairs))
	rights := make([]joinColumn, len(pairs))
	for i, p := range pairs {
		lefts[i], rights[i] = p.l, p.r
		sc.hidden[p.r.key()] = true
		sc.coalesced[strings.ToLower(p.l.name)] = p.l
	}
	sc.joined[j] = lefts
	return columnEqualities(lefts, rights), true
}

// joinOnNothing records a NATURAL join whose inputs share no column.
func (sc *fromScope) joinOnNothing(j *ast.JoinClause) {
	if sc.joined == nil {
		sc.joined = make(map[*ast.JoinClause][]joinColumn)
		sc.hidden = make(map[colKey]bool)
		sc.coalesced = make(map[string]joinColumn)
	}
	sc.joined[j] = []joinColumn{}
}

// noColumnsCondition is the condition MySQL writes for an outer NATURAL join
// whose inputs share no column: on(true).
func noColumnsCondition() ast.ExprNode {
	return &ast.BoolLit{Value: true}
}

// naturalColumns returns the names of the columns a NATURAL join joins on:
// those its inputs share.
func (sc *fromScope) naturalColumns(j *ast.JoinClause) ([]string, bool) {
	left, right := j.Left, j.Right
	if isRightJoin(j.Type) {
		left, right = right, left
	}
	lm, ok := sc.firstColumns(left)
	if !ok {
		return nil, false
	}
	rm, ok := sc.firstColumns(right)
	if !ok {
		return nil, false
	}
	small, large := lm, rm
	if len(small) > len(large) {
		small, large = large, small
	}
	var names []string
	for k := range small {
		if hasColumn(large, k) {
			names = append(names, lm[k].name)
		}
	}
	return names, true
}

// columnEqualities builds left[0] = right[0] AND left[1] = right[1] ...,
// chained to the left.
func columnEqualities(lefts, rights []joinColumn) ast.ExprNode {
	var result ast.ExprNode
	for i := range lefts {
		eq := &ast.BinaryExpr{
			Op:    ast.BinOpEq,
			Left:  &ast.ColumnRef{Table: lefts[i].qualifier, Column: lefts[i].name},
			Right: &ast.ColumnRef{Table: rights[i].qualifier, Column: rights[i].name},
		}
		if result == nil {
			result = eq
		} else {
			result = &ast.BinaryExpr{Op: ast.BinOpAnd, Left: result, Right: eq}
		}
	}
	return result
}
