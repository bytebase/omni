package analysis

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// resolveDerivedLineage deepens result-column lineage through relation
// indirection (derived tables, CTEs, UNNEST) and scalar subqueries.
//
// GetQuerySpan's primary walk records each result column's source refs as the
// refs written in the select item (e.g. `d.x`, a CTE column `pp`, an UNNEST
// output `t.p`), or records nothing for a scalar subquery used as a value. When
// the reference names a derived relation's output column rather than a base
// table, a downstream name-match against the FROM tables finds nothing and the
// column is left without lineage; Bytebase then has no source column to attach a
// masker to and returns the (possibly sensitive) value unmasked.
//
// This pass recomputes the outermost select block's lineage with relation-scope
// resolution: it builds the relations visible in FROM (resolving each derived
// relation's output columns and their underlying base-column refs recursively,
// so derived-over-derived and UNNEST-of-a-base-column resolve) and re-resolves
// each select item — following scalar subqueries to their output sources — then
// unions the result into the primary walk's Results. The union is additive: the
// walk's refs are always retained, so the resolved source set is a superset of
// the original and the pass can only ever deepen lineage (mask more), never drop
// a previously-correct ref (which would under-mask and leak). A ref naming a real
// base table (or an unrecognised relation) is left as-is.
//
// Scope: derived-relation projection (subqueries in FROM, CTE references),
// UNNEST output columns, scalar subqueries in the select list, set-operation
// arm merging (each output column's lineage is the union of that column across
// all arms), and star expansion — a `*` / `<rel>.*` covering only derived
// relations with fully-resolved projections is replaced by its exact output
// columns (width and order), so `SELECT *` over a derived table or CTE masks
// positionally; a star that needs catalog metadata (base table, coalescing
// join, unresolved relation) stays a single opaque "*" result, as before.
func resolveDerivedLineage(stmt ast.Node, span *QuerySpan, vs *viewState) {
	if span == nil {
		return
	}
	qs, ok := stmt.(*parser.QueryStmt)
	if !ok || qs.Query == nil {
		return
	}
	// The root CTE scope carries the catalog-aware resolution state (nil
	// disables catalog resolution); buildCTEDefs propagates it to nested
	// scopes.
	root := &cteDefs{views: vs, defined: make(map[string]int)}
	// Recompute the outermost query's column lineage with relation-scope
	// resolution (derived tables, CTEs, UNNEST, scalar subqueries, set-operation
	// arm merging, and star expansion) and fold it into the primary walk's
	// Results. The walk produces exactly one Result per select item; the resolved
	// columns carry their select-item ordinal, so each walk Result aligns with
	// its resolved group:
	//   - an ordinary item (one resolved column) keeps its walk Result with the
	//     resolved sources unioned in — additive, so lineage only deepens;
	//   - a resolved star is SPLICED: its walk Result (the single "*" entry,
	//     which had no real source columns) is replaced by the expanded output
	//     columns, giving the projection its true width and order so the
	//     positional masker downstream stays aligned;
	//   - an opaque star (not expandable without metadata) is left byte-for-byte
	//     untouched, preserving the exact result shape consumers key on to apply
	//     their own metadata-based expansion.
	cols := resolveQueryCols(qs.Query, root)

	// A top-level VALUES produces no Results in the primary walk (it has no
	// select items), leaving the consumer with zero positional maskers — a
	// sensitive value in a VALUES row (e.g. a scalar subquery) would pass
	// through unmasked. When the walk produced nothing but the resolver knows
	// the exact projection (VALUES columns are always resolved), synthesize the
	// Results from it. Other statement shapes always have walk Results, so this
	// only fires for VALUES (including set operations whose first arm is one).
	if len(span.Results) == 0 {
		if len(cols) > 0 && !hasOpaque(cols) {
			for _, c := range cols {
				span.Results = append(span.Results, ColumnInfo{Name: c.name, SourceColumns: c.sources})
			}
		}
		return
	}

	groups := groupByItem(cols)
	newResults := make([]ColumnInfo, 0, len(span.Results))
	for i := range span.Results {
		r := span.Results[i]
		var g []outCol
		if i < len(groups) {
			g = groups[i]
		}
		switch {
		case len(g) == 0:
			newResults = append(newResults, r)
		case g[0].opaque:
			// Unresolved star: leave the walk's result untouched (including its
			// source refs — a qualified star's single relation-name ref is shape
			// consumers detect; unioning anything in could break that detection).
			newResults = append(newResults, r)
		case len(g) == 1 && !g[0].fromStar:
			r.SourceColumns = unionRefs(r.SourceColumns, g[0].sources)
			newResults = append(newResults, r)
		default:
			// Resolved star: splice the expanded projection in place of the "*".
			for _, c := range g {
				newResults = append(newResults, ColumnInfo{Name: c.name, SourceColumns: c.sources})
			}
		}
	}
	span.Results = newResults
}

// groupByItem buckets resolved columns by their select-item ordinal, preserving
// order within each bucket. Ordinals are dense (assigned per select item), so a
// slice indexed by ordinal suffices.
func groupByItem(cols []outCol) [][]outCol {
	var groups [][]outCol
	for _, c := range cols {
		for len(groups) <= c.item {
			groups = append(groups, nil)
		}
		groups[c.item] = append(groups[c.item], c)
	}
	return groups
}

// ---------------------------------------------------------------------------
// Resolved output columns and CTE definitions
// ---------------------------------------------------------------------------

// outCol is one output column of a (sub)query or relation: its name and the
// base-column refs that feed it, recovered through any derived relations.
//
// A star select item is represented in one of two ways:
//   - resolved: the star is expanded inline into the covered relations' output
//     columns (each marked fromStar), so the projection's width and order are
//     exact — this is what lets `SELECT *` over a derived table or CTE be
//     masked positionally;
//   - opaque: the star cannot be expanded without catalog metadata (it covers a
//     base table, an unresolved relation, or a coalescing join), so a single
//     placeholder entry (opaque=true) holds its position and the projection's
//     true width is unknown.
type outCol struct {
	name    string
	sources []ColumnRef
	// item is the select-item ordinal this column came from, used by the
	// top-level splice to align resolved columns with the primary walk's
	// Results (a resolved star yields several columns with the same item).
	item int
	// fromStar marks a column produced by expanding a resolved star.
	fromStar bool
	// opaque marks an unresolved star placeholder; a projection containing one
	// has unknown width and blocks star expansion through it.
	opaque bool
}

// hasOpaque reports whether cols contains an unresolved star placeholder (the
// projection's width is then unknown).
func hasOpaque(cols []outCol) bool {
	for _, c := range cols {
		if c.opaque {
			return true
		}
	}
	return false
}

// cteDefs maps a CTE name (lower-cased) to its resolved output columns, with a
// parent link for nested WITH scopes. A CTE shadows an outer one of the same
// name; an inner WITH's definitions chain to the outer ones. The chain also
// carries the catalog-aware resolution state (views) from the root scope.
type cteDefs struct {
	defs   map[string][]outCol
	parent *cteDefs
	views  *viewState
	// defined counts, across every scope chained to the same root, the
	// scopes defining each name. A name no scope defines, such as a base
	// table's, is answered without walking the chain: walking it for every
	// table made n nested WITH clauses O(n²).
	defined map[string]int
}

func (c *cteDefs) lookup(name string) ([]outCol, bool) {
	key := strings.ToLower(name)
	if c == nil || c.defined[key] == 0 {
		return nil, false
	}
	for cur := c; cur != nil; cur = cur.parent {
		if cols, ok := cur.defs[key]; ok {
			return cols, true
		}
	}
	return nil, false
}

// viewState returns the catalog-aware resolution state carried by this scope
// chain, or nil when there is none (catalog-less analysis). Every scope copies
// its parent's, so the root's state is on each scope.
func (c *cteDefs) viewState() *viewState {
	if c == nil {
		return nil
	}
	return c.views
}

// buildCTEDefs resolves a query's WITH clause into a cteDefs scope chained to
// parent. Each CTE body is resolved in a scope that already includes the earlier
// siblings (sequential visibility), matching standard non-recursive CTE scoping;
// a CTE is not visible to itself, so a recursive self-reference resolves as a
// base table (best-effort) rather than recursing forever.
func buildCTEDefs(q *parser.Query, parent *cteDefs) *cteDefs {
	if q == nil || q.With == nil || len(q.With.CTEs) == 0 {
		return parent
	}
	d := &cteDefs{defs: make(map[string][]outCol), parent: parent}
	if parent != nil {
		d.views = parent.views
		d.defined = parent.defined
	}
	if d.defined == nil {
		d.defined = make(map[string]int)
	}
	for i := range q.With.CTEs {
		nq := q.With.CTEs[i]
		name := identName(nq.Name)
		if name == "" {
			continue
		}
		cols := resolveQueryCols(nq.Query, d)
		cols = applyColumnAliases(cols, nq.ColumnAliases)
		key := strings.ToLower(name)
		if _, ok := d.defs[key]; !ok {
			d.defined[key]++
		}
		d.defs[key] = cols
	}
	return d
}

// ---------------------------------------------------------------------------
// Output-column resolution
// ---------------------------------------------------------------------------

// resolveQueryCols computes the resolved output columns of a query (its WITH
// scope plus body).
func resolveQueryCols(q *parser.Query, parent *cteDefs) []outCol {
	if q == nil {
		return nil
	}
	cte := buildCTEDefs(q, parent)
	return resolveNodeCols(q.Body, cte)
}

// resolveNodeCols computes the resolved output columns of a query body / set-op
// operand / query primary. For a set operation the left arm's columns are used
// (SQL's "column names come from the first select" rule), matching the primary
// walk; merging the arms' lineage is tracked separately.
func resolveNodeCols(node parser.QueryNode, cte *cteDefs) []outCol {
	switch n := node.(type) {
	case *parser.QuerySpec:
		return resolveSpecCols(n, cte)
	case *parser.SetOperation:
		// A set operation's output columns take their names from the first (left)
		// arm, but a value in ANY arm surfaces in the corresponding output column,
		// so the lineage of column i is the union of column i across all arms. The
		// right operand may itself be a set operation, so recursion accumulates
		// every arm's positional sources.
		return resolveSetOpCols(n, cte)
	case *parser.ParenQuery:
		return resolveQueryCols(n.Inner, cte)
	case *parser.TableQuery:
		// TABLE name == SELECT * FROM name. Over an in-scope CTE the projection
		// is the CTE's resolved columns (verified against Trino 481:
		// `WITH w AS (SELECT phone, name …) TABLE w` returns [phone, name]);
		// the same holds for a catalog-resolved view (its definition's
		// projection) or table (its catalog columns). Otherwise the star needs
		// metadata the analysis does not have and stays opaque.
		parts := normalizedParts(n.Name)
		if len(parts) == 1 && cte != nil {
			if cols, ok := cte.lookup(parts[0]); ok {
				if len(cols) > 0 && !hasOpaque(cols) {
					return stampStar(cols, 0, nil)
				}
				return []outCol{{name: "*", opaque: true}}
			}
		}
		if vs := cte.viewState(); vs != nil {
			if key, view, table, found := vs.lookupRelation(parts); found {
				var cols []outCol
				if view != nil {
					if proj := vs.projectionFor(key, view); proj != nil {
						vs.collectTables(proj.tables)
						cols = proj.cols
					}
				} else {
					cols = tableCols(key, table)
				}
				if len(cols) > 0 && !hasOpaque(cols) {
					return stampStar(cols, 0, nil)
				}
			}
		}
		return []outCol{{name: "*", opaque: true}}
	case *parser.ValuesQuery:
		return resolveValuesCols(n, cte)
	default:
		return nil
	}
}

// resolveValuesCols computes the output columns of a VALUES clause. Each row is a
// row expression (a RowConstructor for the multi-column form, a bare expression
// for a single column); output column i's lineage is the union of column i's
// expression sources across all rows, so a sensitive value (e.g. a scalar
// subquery) anywhere in a VALUES arm of a set operation is still masked. VALUES
// has no FROM relations, so column refs resolve to themselves; scalar subqueries
// within the row expressions are followed.
//
// Each row's sources are appended to its column's, and seen holds each
// column's sources as a set. Uniting them into a copy of the column's sources,
// as each row did, copied lineage that grows by a row each time, so n rows of
// distinct refs cost O(n²). A column's sources start empty, not nil, as the
// copy did: a column of literals has no sources, but not nil ones.
func resolveValuesCols(n *parser.ValuesQuery, cte *cteDefs) []outCol {
	if n == nil {
		return nil
	}
	scope := &rscope{cte: cte}
	var cols []outCol
	var seen []map[ColumnRef]bool
	for _, row := range n.Rows {
		for i, e := range valuesRowElements(row) {
			for len(cols) <= i {
				cols = append(cols, outCol{item: len(cols), sources: []ColumnRef{}})
				seen = append(seen, make(map[ColumnRef]bool))
			}
			for _, r := range scope.resolveExprRefs(e) {
				if !seen[i][r] {
					seen[i][r] = true
					cols[i].sources = append(cols[i].sources, r)
				}
			}
		}
	}
	return cols
}

// valuesRowElements returns the per-column expressions of one VALUES row: the
// elements of a RowConstructor, or the row itself as a single column.
func valuesRowElements(row parser.Expr) []parser.Expr {
	if rc, ok := row.(*parser.RowConstructor); ok {
		return rc.Elements
	}
	return []parser.Expr{row}
}

// resolveSetOpCols computes the resolved output columns of a set operation:
// its arms merged positionally (see armMerger.merge), from its leftmost arm
// on, through the chain of set operations down its left side (`a UNION b
// UNION c` is (a UNION b) UNION c). The arms are merged into one accumulator:
// merging each arm into a copy of the lineage of the arms before it copied
// lineage that grows by an arm's sources each time, so n arms of distinct
// columns cost O(n²).
func resolveSetOpCols(n *parser.SetOperation, cte *cteDefs) []outCol {
	// rights lists the chain's right operands, outermost first.
	var rights []parser.QueryNode
	var node parser.QueryNode = n
	for {
		op, ok := node.(*parser.SetOperation)
		if !ok {
			break
		}
		rights = append(rights, op.Right)
		node = op.Left
	}
	m := newArmMerger(resolveNodeCols(node, cte))
	for i := len(rights) - 1; i >= 0; i-- {
		m.merge(resolveNodeCols(rights[i], cte))
	}
	return m.cols
}

// armMerger accumulates the merged columns of set-operation arms. The first
// merge copies the columns, and a column keeps the set of its sources from
// its first merge on, so merging an arm costs that arm's size.
type armMerger struct {
	cols   []outCol
	opaque bool // hasOpaque(cols), which no merge changes
	copied bool // cols is the merger's own copy
	seen   []map[ColumnRef]bool
}

func newArmMerger(first []outCol) *armMerger {
	return &armMerger{cols: first, opaque: hasOpaque(first)}
}

// merge positionally merges the resolved columns of the next set-operation
// arm into cols. The result keeps the shape of cols (count, names, and item
// stamps — SQL takes the output column names from the first SELECT); each
// column's sources become the union (unionRefs) of its sources and the arm's
// at that position, so a sensitive value contributed by any arm forces the
// column to be masked. An arm with fewer columns (only in malformed SQL)
// contributes nothing beyond the positions it has.
//
// With stars expanded inline, the arms of a valid set operation have their
// true width, so the positional union is exact — including a resolved star
// arm (`SELECT * FROM (…) d UNION SELECT phone FROM t` merges phone into d's
// expanded column). When either side still contains an OPAQUE star its width
// is unknown and positional pairing would be unsound, so the merge is skipped
// and cols passes through unchanged (parity with the walk: the consumer's own
// star expansion applies and the other arm's lineage remains a documented
// residual rather than a mis-attributed union).
func (m *armMerger) merge(right []outCol) {
	if m.opaque || hasOpaque(right) {
		return
	}
	if !m.copied {
		m.cols = append([]outCol(nil), m.cols...)
		m.seen = make([]map[ColumnRef]bool, len(m.cols))
		m.copied = true
	}
	for i := 0; i < len(m.cols) && i < len(right); i++ {
		seen := m.seen[i]
		if seen == nil {
			// The column's first merge: the union drops a repeated source of
			// the column's own too.
			left := m.cols[i].sources
			seen = make(map[ColumnRef]bool, len(left)+len(right[i].sources))
			out := make([]ColumnRef, 0, len(left)+len(right[i].sources))
			for _, r := range left {
				if !seen[r] {
					seen[r] = true
					out = append(out, r)
				}
			}
			m.cols[i].sources = out
			m.seen[i] = seen
		}
		for _, r := range right[i].sources {
			if !seen[r] {
				seen[r] = true
				m.cols[i].sources = append(m.cols[i].sources, r)
			}
		}
	}
}

// resolveSpecCols computes the resolved output columns of one SELECT block: it
// builds the FROM scope, then resolves each select item through that scope. A
// star item is expanded inline when every relation it covers is a derived
// relation with a fully-resolved projection (see starExpansion); otherwise it
// stays a single opaque placeholder, exactly as before.
func resolveSpecCols(spec *parser.QuerySpec, cte *cteDefs) []outCol {
	if spec == nil {
		return nil
	}
	scope := newRScope(spec.From, cte)
	out := make([]outCol, 0, len(spec.Items))
	for idx, item := range spec.Items {
		if cols, ok := scope.starItemCols(item, idx); ok {
			out = append(out, cols...)
			continue
		}
		name := identName(item.Alias)
		if name == "" {
			name = renderExprName(item.Expr)
		}
		out = append(out, outCol{name: name, sources: scope.resolveExprRefs(item.Expr), item: idx})
	}
	return out
}

// starItemCols returns the output columns of select item idx when it is a
// star, and ok false when it is not.
func (s *rscope) starItemCols(item parser.SelectItem, idx int) (cols []outCol, ok bool) {
	switch item.Kind {
	case parser.SelectAll:
		if exp := s.starExpansion(""); exp != nil {
			return stampStar(exp, idx, nil), true
		}
		return []outCol{{name: "*", item: idx, opaque: true}}, true
	case parser.SelectAllFrom:
		if q := starQualifier(item.Expr); q != "" {
			if exp := s.starExpansion(q); exp != nil {
				return stampStar(exp, idx, item.Aliases), true
			}
		}
		name := renderExprName(item.Expr)
		if name == "" {
			name = "*"
		} else {
			name += ".*"
		}
		return []outCol{{name: name, sources: s.resolveExprRefs(item.Expr), item: idx, opaque: true}}, true
	}
	return nil, false
}

// starQualifier returns the relation qualifier a `<rel>.*` select item names:
// the rightmost identifier of the row expression (the relation/alias for `d.*`,
// the table for `sch.t.*`). Empty when the expression is not a plain name chain
// (e.g. a row-valued function), in which case the star is not expandable.
func starQualifier(expr parser.Expr) string {
	switch e := expr.(type) {
	case *parser.ColumnRef:
		return identName(e.Name)
	case *parser.Dereference:
		return identName(e.FieldName)
	case *parser.ParenExpr:
		return starQualifier(e.Expr)
	}
	return ""
}

// stampStar returns a copy of a star's expanded columns stamped with the star's
// select-item ordinal and the fromStar marker, with any `AS (a, b, …)` column
// aliases applied positionally. The input (a relation's shared projection) is
// never mutated.
func stampStar(cols []outCol, item int, aliases []*ast.Identifier) []outCol {
	out := make([]outCol, len(cols))
	copy(out, cols)
	out = applyColumnAliases(out, aliases)
	for i := range out {
		out[i].item = item
		out[i].fromStar = true
	}
	return out
}

// ---------------------------------------------------------------------------
// Relation scope
// ---------------------------------------------------------------------------

// rscope is the set of relations visible in a query body, used to resolve a
// column reference to base-column refs. A derived relation (subquery / CTE
// reference) carries its resolved output columns; a base table carries no
// columns (its references resolve to themselves).
type rscope struct {
	rels []rbind
	cte  *cteDefs
	// coalesced is set when the FROM tree contains a USING or NATURAL join,
	// which coalesces the join columns into single output columns — `SELECT *`
	// then has a different width/order than the relations' concatenation, so
	// star expansion is blocked in this scope.
	coalesced bool

	// idx indexes rels once the scope is large (see index); starAll holds
	// starExpansion("") for the first starRels relations, once starDone.
	idx      *relIndex
	starAll  []outCol
	starRels int
	starDone bool

	// recovering is set once a relation of rels[:recoverChecked] has a
	// column with sources (see recovers).
	recovering     bool
	recoverChecked int
}

// relIndex indexes a scope's relations and their columns by name, by
// foldKey, so a reference finds the relations and the column it names
// without comparing it against every one: comparing made a select list that
// names the relations of a long join chain, or the columns of a wide derived
// table, cost the product of the two. The scope only appends relations, and
// the index catches up with them as it is used.
type relIndex struct {
	indexed int              // the relations indexed: rels[:indexed]
	byName  map[string][]int // a relation name: the relations of that name, in order
	// byCol maps a column name to the relations exposing a column of that
	// name, in order, each with the first such column.
	byCol map[string][]colHit
}

type colHit struct{ rel, col int }

// scopeIndexMin is how many relations and columns a scope holds before its
// lookups go through a relIndex; below it, comparing names is as fast and
// allocates nothing. It is a variable so that tests can index every scope.
var scopeIndexMin = 32

// index returns the scope's relIndex, brought up to date, or nil while the
// scope is small.
func (s *rscope) index() *relIndex {
	if s.idx == nil {
		size := len(s.rels)
		for i := 0; i < len(s.rels) && size < scopeIndexMin; i++ {
			size += len(s.rels[i].cols)
		}
		if size < scopeIndexMin {
			return nil
		}
		s.idx = &relIndex{byName: make(map[string][]int), byCol: make(map[string][]colHit)}
	}
	x := s.idx
	for ; x.indexed < len(s.rels); x.indexed++ {
		ri := x.indexed
		name := foldKey(s.rels[ri].name)
		x.byName[name] = append(x.byName[name], ri)
		for ci, c := range s.rels[ri].cols {
			col := foldKey(c.name)
			hits := x.byCol[col]
			if n := len(hits); n > 0 && hits[n-1].rel == ri {
				continue // not the relation's first column of that name
			}
			x.byCol[col] = append(hits, colHit{rel: ri, col: ci})
		}
	}
	return x
}

// lookupOutCol is the package function lookupOutCol on relation ri's
// columns, for a column name whose foldKey is col.
func (x *relIndex) lookupOutCol(s *rscope, ri int, col string) ([]ColumnRef, bool) {
	hits := x.byCol[col]
	i := sort.Search(len(hits), func(i int) bool { return hits[i].rel >= ri })
	if i < len(hits) && hits[i].rel == ri {
		return s.rels[ri].cols[hits[i].col].sources, true
	}
	return nil, false
}

// foldKey returns a key that two strings share exactly when strings.EqualFold
// holds for them: each rune, as EqualFold reads it (an invalid byte is
// utf8.RuneError), replaced by the smallest rune of its case-folding orbit
// (unicode.SimpleFold).
func foldKey(s string) string {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		// The smallest rune of an ASCII letter's orbit is its upper case;
		// the orbits of k and s also hold U+212A and U+017F, which are larger.
		return strings.ToUpper(s)
	}
	var b strings.Builder
	for _, r := range s {
		least := r
		for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
			least = min(least, f)
		}
		b.WriteRune(least)
	}
	return b.String()
}

// rbind binds one relation in a FROM scope to a name. derived marks a subquery
// or CTE reference (cols then holds its resolved output columns); a base table
// has derived=false and nil cols.
type rbind struct {
	name    string
	derived bool
	cols    []outCol
}

// newRScope builds a relation scope from a FROM list.
func newRScope(from []parser.Relation, cte *cteDefs) *rscope {
	s := &rscope{cte: cte}
	for _, rel := range from {
		s.add(rel, "", nil)
	}
	return s
}

// add binds one FROM relation (recursing through joins/parentheses). alias and
// colAliases carry an enclosing AliasedRelation's `AS name (c1, …)` down to the
// primary relation it names.
func (s *rscope) add(rel parser.Relation, alias string, colAliases []*ast.Identifier) {
	switch n := rel.(type) {
	case *parser.AliasedRelation:
		s.add(n.Inner, identName(n.Alias), n.ColumnAliases)
	case *parser.Join:
		if alias != "" {
			// An aliased (parenthesized) join — ((…) a JOIN (…) b ON …) AS j —
			// is one relation named j whose projection is the join's column
			// concatenation; Trino hides the inner aliases behind it (verified
			// against Trino 481: `a.phone` does not resolve through `… AS j`).
			// Bind a single synthetic relation: derived when the join subtree
			// is fully resolvable (all derived, no coalescing), opaque-blocking
			// otherwise.
			sub := &rscope{cte: s.cte}
			if len(n.Using) > 0 || n.Natural {
				sub.coalesced = true
			}
			sub.add(n.Left, "", nil)
			sub.add(n.Right, "", nil)
			cols := sub.starExpansion("")
			s.rels = append(s.rels, rbind{name: alias, derived: cols != nil, cols: applyColumnAliases(cols, colAliases)})
			return
		}
		if len(n.Using) > 0 || n.Natural {
			// USING/NATURAL coalesces the join columns, changing the output
			// width and order: an unqualified `*` lists the coalesced columns
			// once, first; and a QUALIFIED star EXCLUDES the using columns from
			// its relation's projection (verified against Trino 481: with
			// a(k, p), `SELECT a.* FROM a JOIN b USING (k)` returns only [p]).
			// Star expansion in this scope would therefore be width-wrong in
			// both forms; block it. Named-column resolution is unaffected.
			s.coalesced = true
		}
		s.add(n.Left, "", nil)
		s.add(n.Right, "", nil)
	case *parser.ParenRelation:
		s.add(n.Inner, alias, colAliases)
	case *parser.TableRelation:
		parts := normalizedParts(n.Name)
		// A single-part name that matches an in-scope CTE is a derived relation
		// (the CTE's resolved columns); a CTE shadows any same-named catalog
		// object.
		if len(parts) == 1 && s.cte != nil {
			if cols, ok := s.cte.lookup(parts[0]); ok {
				name := alias
				if name == "" {
					name = parts[0]
				}
				s.rels = append(s.rels, rbind{name: name, derived: true, cols: applyColumnAliases(cols, colAliases)})
				return
			}
		}
		name := alias
		if name == "" {
			name = lastPart(n.Name)
		}
		// Resolve against the catalog when one was supplied. A VIEW binds as a
		// derived relation carrying its definition's resolved projection, so
		// references and stars through it reach the underlying base columns; a
		// TABLE binds with its catalog columns, which star expansion uses (named
		// references through a base table keep their written form — additive
		// resolution needs no help there, see resolveRef).
		if vs := s.cte.viewState(); vs != nil {
			if key, view, table, found := vs.lookupRelation(parts); found {
				if view != nil {
					if proj := vs.projectionFor(key, view); proj != nil {
						vs.collectTables(proj.tables)
						s.rels = append(s.rels, rbind{name: name, derived: true, cols: applyColumnAliases(proj.cols, colAliases)})
						return
					}
					// Unresolvable view (no/unanalyzable definition, cycle):
					// bind opaquely so stars through it stay unexpanded.
					s.rels = append(s.rels, rbind{name: name, derived: false})
					return
				}
				s.rels = append(s.rels, rbind{name: name, derived: false, cols: applyColumnAliases(tableCols(key, table), colAliases)})
				return
			}
		}
		s.rels = append(s.rels, rbind{name: name, derived: false})
	case *parser.SubqueryRelation:
		cols := resolveQueryCols(n.Query, s.cte)
		s.rels = append(s.rels, rbind{name: alias, derived: true, cols: applyColumnAliases(cols, colAliases)})
	case *parser.UnnestRelation:
		s.rels = append(s.rels, rbind{name: alias, derived: true, cols: s.unnestColumns(n, colAliases)})
	default:
		// Lateral and table-function relations are not resolved here (tracked
		// separately). Bind the alias as a non-derived (opaque) relation so
		// references through it pass through unchanged.
		s.rels = append(s.rels, rbind{name: alias, derived: false})
	}
}

// unnestColumns computes the output columns of an UNNEST relation. Each named
// output column (from the enclosing `AS alias (c1, …)` list) is an element of an
// unnested collection; its lineage is the column(s) referenced by the unnested
// expressions, resolved against the relations already in scope — UNNEST is
// implicitly lateral over the left side of its join, so those expressions may
// reference the left relations. The number of output columns produced per
// expression depends on the element type (array → 1, map → 2, array<row> → N),
// which is unknown without type metadata, so every named column conservatively
// carries the union of all unnested expressions' sources (over-inclusion is safe
// for masking). A trailing WITH ORDINALITY column is an ordinal and carries no
// lineage. Without column aliases the output columns are unnamed and cannot be
// referenced by name, so none are produced (best-effort, unchanged behaviour).
//
// Boundary: a scalar subquery as the unnested expression — UNNEST((SELECT
// array_agg(x) FROM t)) — is not followed here (resolveExprRefs collects direct
// columns only), so its inner columns are not yet recovered. Because resolution
// is additive this never regresses below the prior behaviour; the inner lineage
// is picked up once scalar-subquery resolution lands (BYT-9676).
func (s *rscope) unnestColumns(n *parser.UnnestRelation, colAliases []*ast.Identifier) []outCol {
	var srcs []ColumnRef
	for _, e := range n.Exprs {
		srcs = append(srcs, s.resolveExprRefs(e)...)
	}
	cols := make([]outCol, 0, len(colAliases))
	for i, a := range colAliases {
		name := identName(a)
		if n.WithOrdinality && i == len(colAliases)-1 {
			cols = append(cols, outCol{name: name})
			continue
		}
		cols = append(cols, outCol{name: name, sources: srcs})
	}
	return cols
}

// starExpansion returns the exact output columns a star select item covers, or
// nil when expansion is not provably width- and order-correct. An empty
// qualifier is an unqualified `*` (covering every FROM relation, in order); a
// non-empty qualifier is `<rel>.*` (covering exactly that relation).
//
// Expansion requires certainty — a wrong width or order would misalign the
// positional masker downstream, which is precisely the bug this resolves — so
// it bails (nil) unless:
//   - no USING/NATURAL join coalesces columns in this scope, and
//   - every covered relation has a fully-resolved projection: non-empty and
//     free of opaque star placeholders. That is a derived relation (subquery,
//     CTE reference, aliased UNNEST, catalog-resolved view) or a base table
//     whose columns the supplied catalog knows. A base table without catalog
//     metadata, a lateral/table-function relation, an UNNEST without column
//     aliases, or a qualifier matching zero or several relations all bail.
//
// A nil return leaves the star opaque — the consumer's metadata-based expansion
// applies, exactly as before this resolver existed.
func (s *rscope) starExpansion(qualifier string) []outCol {
	if s.coalesced || len(s.rels) == 0 {
		return nil
	}
	if qualifier != "" {
		var match *rbind
		count := 0
		if x := s.index(); x != nil {
			named := x.byName[foldKey(qualifier)]
			if count = len(named); count > 0 {
				match = &s.rels[named[count-1]]
			}
		} else {
			for i := range s.rels {
				if strings.EqualFold(s.rels[i].name, qualifier) {
					match = &s.rels[i]
					count++
				}
			}
		}
		if count != 1 || len(match.cols) == 0 || hasOpaque(match.cols) {
			return nil
		}
		return match.cols
	}
	// Computed once per scope: every `*` item of a select list expands to the
	// same columns, and expanding them per item cost the product of the two.
	if s.starDone && s.starRels == len(s.rels) {
		return s.starAll
	}
	var out []outCol
	for _, rb := range s.rels {
		if len(rb.cols) == 0 || hasOpaque(rb.cols) {
			out = nil
			break
		}
		out = append(out, rb.cols...)
	}
	s.starAll, s.starRels, s.starDone = out, len(s.rels), true
	return out
}

// resolveRef resolves a single column ref through the scope, ADDITIVELY: the
// original ref is always kept, and any base refs recovered through a derived
// relation are appended. This guarantees the resolved source set is a superset
// of the original, so the pass can only ever deepen lineage (mask more), never
// drop a previously-correct ref (which would under-mask and leak).
//
// Additivity is what makes the pass safe in the cases where a derived relation's
// column has incomplete or unknown lineage but a same-named column is also valid
// against a base table:
//   - a set operation inside the derived relation (only the left arm's lineage
//     is computed, so the right arm's sensitive column would otherwise be lost);
//   - a derived column sourced from `SELECT *` (unknown lineage — empty here);
//   - `JOIN ... USING (c)`, where bare `c` is a valid output coalescing a base
//     table and a derived relation.
//
// In every such case the kept original bare/base ref still masks the base
// column, while the recovered derived sources broaden coverage. A ref that names
// no derived relation (a base table or an unknown relation) is returned
// unchanged.
func (s *rscope) resolveRef(ref ColumnRef) []ColumnRef {
	out := []ColumnRef{ref}
	if x := s.index(); x != nil {
		// The loops below, through the index.
		col := foldKey(ref.Column)
		if ref.Table != "" {
			for _, ri := range x.byName[foldKey(ref.Table)] {
				if src, ok := x.lookupOutCol(s, ri, col); ok {
					out = append(out, src...)
				}
			}
			return out
		}
		for _, h := range x.byCol[col] {
			out = append(out, s.rels[h.rel].cols[h.col].sources...)
		}
		return out
	}
	if ref.Table != "" {
		// Qualified by a relation name: append the recovered sources of every
		// in-scope relation of that name exposing the column (a reused alias
		// yields several; unioning them over-includes, which is safe). A
		// catalog-known base table's cols matter when relation column aliases
		// rename its columns (FROM customer AS c(i, p, …): c.p must reach
		// customer.phone); without aliases they only restate the written ref
		// in qualified form, which is harmless.
		for _, rb := range s.rels {
			if strings.EqualFold(rb.name, ref.Table) {
				if src, ok := lookupOutCol(rb.cols, ref.Column); ok {
					out = append(out, src...)
				}
			}
		}
		return out
	}
	// Bare reference: append the recovered sources of every in-scope relation
	// exposing a column of this name. The original bare ref is retained so a
	// base table providing the same column is still covered.
	for _, rb := range s.rels {
		if src, ok := lookupOutCol(rb.cols, ref.Column); ok {
			out = append(out, src...)
		}
	}
	return out
}

// resolveExprRefs collects a select-item expression's lineage: its direct column
// refs (resolved through the scope) plus the recovered output sources of any
// scalar subquery used as a value within it (SELECT (SELECT phone …) AS sp). It
// does not cross subquery boundaries for direct columns (matching the primary
// walk's collectDirectColumns); scalar subqueries are handled via onSubquery.
//
// A scalar subquery's recovered sources are themselves resolved through this
// (outer) scope, so a correlated reference to an outer relation — e.g. the `d`
// in SELECT (SELECT d.phone) AS sp FROM (SELECT phone FROM customer) d — is
// resolved to its base column rather than left as the unmaskable alias ref.
//
// The result is, without repeats, the direct refs resolved through the scope,
// then the subqueries' sources resolved through it, where resolving a list
// replaces each ref by what resolveRef returns for it. See lineage for how
// it is computed.
func (s *rscope) resolveExprRefs(expr parser.Expr) []ColumnRef {
	l := &lineage{
		out:    refSet{refs: []ColumnRef{}, seen: make(map[ColumnRef]bool)},
		levels: []lineageLevel{{scope: s}},
	}
	l.expr(expr, 0)
	return l.out.refs
}

// lineage computes the result of resolveExprRefs for one expression.
//
// Collecting the sources of each scalar subquery, resolving them through the
// scope, and returning them, as each level of nested subqueries did, copied
// at each level the lineage of every level inside it: n levels that each add
// a column, SELECT a1 + (SELECT a2 + (…)), cost O(n²). Instead, each ref is
// passed up the levels of subqueries as it is found, and only a level whose
// scope may recover sources for it resolves it; any other passes it on
// unchanged, as resolving it there would. A level is found through byKey, so
// a ref passes the levels that do not recover anything for it in one step.
// A ref resolved at a level once is not resolved there again: its first
// pass already brought everything it resolves to.
//
// The result is what resolving the levels' lists one inside the other gives,
// because resolving a list replaces each ref by a list that depends on that
// ref alone, and the result keeps the first of each repeated ref: resolving
// a list that has lost repeats of a ref, or that holds more, gives the same.
type lineage struct {
	out refSet
	// levels[0] is the scope of the expression; levels[i] that of the scalar
	// subquery at levels[i-1] whose select items are being resolved.
	levels []lineageLevel
	// byKey lists, by refKey and in increasing order, the levels from 1 on
	// whose scope may recover sources for a ref of that key (see push).
	// The scope of levels[0] is asked directly: registering it for every
	// expression of a select list would make the list cost the product of
	// its length and its scope's width.
	byKey map[refKey][]int
}

type lineageLevel struct {
	scope *rscope
	seen  map[ColumnRef]bool // the refs resolved at this level already
	keys  []refKey           // the keys registered for this level
}

// refKey is what resolveRef matches a ref by: its table, or "" for a bare
// ref, and its column, each by foldKey.
type refKey struct{ table, column string }

func keyOf(ref ColumnRef) refKey {
	return refKey{table: foldKey(ref.Table), column: foldKey(ref.Column)}
}

// refSet collects refs in the order they first come, each once.
type refSet struct {
	refs []ColumnRef
	seen map[ColumnRef]bool
}

func (a *refSet) add(ref ColumnRef) {
	if !a.seen[ref] {
		a.seen[ref] = true
		a.refs = append(a.refs, ref)
	}
}

// expr passes on the lineage of expr, at level i: its direct refs resolved
// through the level's scope, then its scalar subqueries' sources.
func (l *lineage) expr(expr parser.Expr, i int) {
	var refs []ColumnRef
	var subs []*parser.SubqueryExpr
	exprWalk{
		onColumn:   func(ref ColumnRef) { refs = append(refs, ref) },
		onSubquery: func(sub *parser.SubqueryExpr) { subs = append(subs, sub) },
	}.walk(expr)
	s := l.levels[i].scope
	for _, ref := range refs {
		if !s.recovers() {
			l.up(ref, i)
			continue
		}
		for _, r := range s.resolveRef(ref) {
			l.up(r, i)
		}
	}
	for _, sub := range subs {
		l.scalarSubquery(sub, s.cte, i+1)
	}
}

// up passes ref, a ref of level i, on through the levels below i.
func (l *lineage) up(ref ColumnRef, i int) {
	if i == 0 {
		l.out.add(ref)
		return
	}
	j := 0
	if l.byKey != nil {
		stack := l.byKey[keyOf(ref)]
		if k := sort.SearchInts(stack, i) - 1; k >= 0 {
			j = stack[k]
		}
	}
	// levels[j] is the first level below i that may resolve ref to more than
	// itself; the levels between pass it on unchanged.
	lv := &l.levels[j]
	if j == 0 && !lv.scope.recovers() {
		l.out.add(ref)
		return
	}
	if lv.seen[ref] {
		return
	}
	if lv.seen == nil {
		lv.seen = make(map[ColumnRef]bool)
	}
	lv.seen[ref] = true
	for _, r := range lv.scope.resolveRef(ref) {
		l.up(r, j)
	}
}

// push makes s, the scope of a scalar subquery of level i-1, level i.
func (l *lineage) push(s *rscope, i int) {
	l.levels = append(l.levels[:i], lineageLevel{scope: s})
	if !s.recovers() {
		return
	}
	lv := &l.levels[i]
	if l.byKey == nil {
		l.byKey = make(map[refKey][]int)
	}
	register := func(k refKey) {
		stack := l.byKey[k]
		if n := len(stack); n > 0 && stack[n-1] == i {
			return
		}
		l.byKey[k] = append(stack, i)
		lv.keys = append(lv.keys, k)
	}
	// resolveRef recovers sources for a bare ref from a column of that name
	// in any relation, and for a qualified one from a column of that name in
	// a relation of that name.
	for _, rb := range s.rels {
		rel := foldKey(rb.name)
		for _, c := range rb.cols {
			if len(c.sources) == 0 {
				continue
			}
			col := foldKey(c.name)
			register(refKey{column: col})
			register(refKey{table: rel, column: col})
		}
	}
}

// pop ends level i, the last.
func (l *lineage) pop(i int) {
	for _, k := range l.levels[i].keys {
		stack := l.byKey[k]
		l.byKey[k] = stack[:len(stack)-1]
	}
	l.levels = l.levels[:i]
}

// recovers reports whether a relation in scope has a column with sources.
// When none has, resolveRef resolves every ref to itself alone.
func (s *rscope) recovers() bool {
	for ; !s.recovering && s.recoverChecked < len(s.rels); s.recoverChecked++ {
		for _, c := range s.rels[s.recoverChecked].cols {
			if len(c.sources) > 0 {
				s.recovering = true
				break
			}
		}
	}
	return s.recovering
}

// scalarSubquery passes on the resolved base-column sources of a scalar
// subquery used as a value (e.g. SELECT (SELECT phone FROM customer) AS sp),
// as refs of level i. The subquery's query (the tree the strict parse built
// in place, or else its raw text parsed again) has its output columns
// resolved (recursively, in the enclosing CTE scope); the sources of those
// columns are the lineage feeding the enclosing expression. EXISTS
// subqueries yield a boolean and are skipped. Best-effort: a correlated
// reference to an outer column is not resolved to the outer relation (it
// resolves within the subquery's own scope).
func (l *lineage) scalarSubquery(sub *parser.SubqueryExpr, cte *cteDefs, i int) {
	if sub == nil || sub.Kind != parser.SubqueryScalar {
		return
	}
	// The walker already failed the span closed on any subquery that does not
	// parse (analyzeSubqueryText), so by the time lineage resolution runs every
	// placeholder body parses; an error here is unreachable and yields no
	// sources rather than a partial tree's.
	queries := []*parser.Query{sub.Query}
	if sub.Query == nil {
		file, err := parser.Parse(sub.RawText)
		if err != nil || file == nil {
			return
		}
		queries = queries[:0]
		for _, stmt := range file.Stmts {
			if qs, ok := stmt.(*parser.QueryStmt); ok {
				queries = append(queries, qs.Query)
			}
		}
	}
	for _, q := range queries {
		l.query(q, cte, i)
	}
}

// query passes on the sources of the output columns resolveQueryCols
// resolves for q, in order, as refs of level i. A select block's scope
// becomes level i while its items are resolved.
func (l *lineage) query(q *parser.Query, parent *cteDefs, i int) {
	if q == nil {
		return
	}
	cte := buildCTEDefs(q, parent)
	switch n := q.Body.(type) {
	case *parser.QuerySpec:
		if n == nil {
			return
		}
		scope := newRScope(n.From, cte)
		l.push(scope, i)
		for idx, item := range n.Items {
			if cols, ok := scope.starItemCols(item, idx); ok {
				l.cols(cols, i)
				continue
			}
			l.expr(item.Expr, i)
		}
		l.pop(i)
	case *parser.ParenQuery:
		l.query(n.Inner, cte, i)
	default:
		l.cols(resolveNodeCols(q.Body, cte), i)
	}
}

// cols passes on the sources of cols, in order, as refs of level i.
func (l *lineage) cols(cols []outCol, i int) {
	for _, c := range cols {
		for _, ref := range c.sources {
			l.up(ref, i)
		}
	}
}

// unionRefs returns the deduplicated concatenation of a and b, preserving a's
// order first and then b's new entries.
func unionRefs(a, b []ColumnRef) []ColumnRef {
	out := make([]ColumnRef, 0, len(a)+len(b))
	seen := make(map[ColumnRef]bool, len(a)+len(b))
	for _, r := range a {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	for _, r := range b {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// lookupOutCol returns the sources of the output column named name (case-
// insensitively), and whether such a column exists. found can be true with
// empty sources (e.g. a literal projection), which correctly yields no lineage.
func lookupOutCol(cols []outCol, name string) ([]ColumnRef, bool) {
	for _, c := range cols {
		if strings.EqualFold(c.name, name) {
			return c.sources, true
		}
	}
	return nil, false
}

// applyColumnAliases renames the leading output columns to the given relation /
// CTE column aliases (`AS t (a, b)` or `WITH w (a, b) AS …`), preserving each
// column's recovered sources. Columns beyond the alias list keep their names;
// surplus aliases are ignored. The input is not mutated.
func applyColumnAliases(cols []outCol, aliases []*ast.Identifier) []outCol {
	if len(aliases) == 0 {
		return cols
	}
	out := make([]outCol, len(cols))
	copy(out, cols)
	for i := range out {
		if i >= len(aliases) {
			break
		}
		if name := identName(aliases[i]); name != "" {
			out[i].name = name
		}
	}
	return out
}

// lastPart returns the rightmost component (the table name) of a qualified name.
func lastPart(name *ast.QualifiedName) string {
	parts := normalizedParts(name)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
