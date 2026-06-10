package analysis

import (
	"strings"

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
func resolveDerivedLineage(stmt ast.Node, span *QuerySpan) {
	if span == nil {
		return
	}
	qs, ok := stmt.(*parser.QueryStmt)
	if !ok || qs.Query == nil {
		return
	}
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
	cols := resolveQueryCols(qs.Query, nil)

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
// name; an inner WITH's definitions chain to the outer ones.
type cteDefs struct {
	defs   map[string][]outCol
	parent *cteDefs
}

func (c *cteDefs) lookup(name string) ([]outCol, bool) {
	key := strings.ToLower(name)
	for cur := c; cur != nil; cur = cur.parent {
		if cols, ok := cur.defs[key]; ok {
			return cols, true
		}
	}
	return nil, false
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
	for i := range q.With.CTEs {
		nq := q.With.CTEs[i]
		name := identName(nq.Name)
		if name == "" {
			continue
		}
		cols := resolveQueryCols(nq.Query, d)
		cols = applyColumnAliases(cols, nq.ColumnAliases)
		d.defs[strings.ToLower(name)] = cols
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
		return mergeArms(resolveNodeCols(n.Left, cte), resolveNodeCols(n.Right, cte))
	case *parser.ParenQuery:
		return resolveQueryCols(n.Inner, cte)
	case *parser.TableQuery:
		// TABLE name == SELECT * FROM name. Over an in-scope CTE the projection
		// is the CTE's resolved columns (verified against Trino 481:
		// `WITH w AS (SELECT phone, name …) TABLE w` returns [phone, name]);
		// over a base table the star needs catalog metadata and stays opaque.
		if parts := normalizedParts(n.Name); len(parts) == 1 && cte != nil {
			if cols, ok := cte.lookup(parts[0]); ok && len(cols) > 0 && !hasOpaque(cols) {
				return stampStar(cols, 0, nil)
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
func resolveValuesCols(n *parser.ValuesQuery, cte *cteDefs) []outCol {
	if n == nil {
		return nil
	}
	scope := &rscope{cte: cte}
	var cols []outCol
	for _, row := range n.Rows {
		for i, e := range valuesRowElements(row) {
			for len(cols) <= i {
				cols = append(cols, outCol{item: len(cols)})
			}
			cols[i].sources = unionRefs(cols[i].sources, scope.resolveExprRefs(e))
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

// mergeArms positionally merges the resolved columns of two set-operation arms.
// The result has the left arm's shape (count, names, and item stamps — SQL takes
// the output column names from the first SELECT); each column's sources are the
// union of the two arms' sources at that position, so a sensitive value
// contributed by either arm forces the column to be masked. A right arm with
// fewer columns than the left (only in malformed SQL) contributes nothing beyond
// the positions it has.
//
// With stars expanded inline, both arms of a valid set operation have their true
// width, so the positional union is exact — including a resolved star arm
// (`SELECT * FROM (…) d UNION SELECT phone FROM t` merges phone into d's
// expanded column). When either arm still contains an OPAQUE star its width is
// unknown and positional pairing would be unsound, so the merge is skipped and
// the left arm passes through unchanged (parity with the walk: the consumer's
// own star expansion applies and the other arm's lineage remains a documented
// residual rather than a mis-attributed union).
func mergeArms(left, right []outCol) []outCol {
	if hasOpaque(left) || hasOpaque(right) {
		return left
	}
	out := make([]outCol, len(left))
	copy(out, left)
	for i := range out {
		if i < len(right) {
			out[i].sources = unionRefs(left[i].sources, right[i].sources)
		}
	}
	return out
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
		switch item.Kind {
		case parser.SelectAll:
			if exp := scope.starExpansion(""); exp != nil {
				out = append(out, stampStar(exp, idx, nil)...)
			} else {
				out = append(out, outCol{name: "*", item: idx, opaque: true})
			}
			continue
		case parser.SelectAllFrom:
			if q := starQualifier(item.Expr); q != "" {
				if exp := scope.starExpansion(q); exp != nil {
					out = append(out, stampStar(exp, idx, item.Aliases)...)
					continue
				}
			}
			name := renderExprName(item.Expr)
			if name == "" {
				name = "*"
			} else {
				name += ".*"
			}
			out = append(out, outCol{name: name, sources: scope.resolveExprRefs(item.Expr), item: idx, opaque: true})
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
		// A single-part name that matches an in-scope CTE is a derived relation
		// (the CTE's resolved columns); otherwise it is a base table.
		if parts := normalizedParts(n.Name); len(parts) == 1 && s.cte != nil {
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
//   - every covered relation is a derived relation (subquery, CTE reference, or
//     aliased UNNEST) whose projection is fully resolved: non-empty and free of
//     opaque star placeholders. A base table (width known only to catalog
//     metadata), a lateral/table-function relation, an UNNEST without column
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
		for i := range s.rels {
			if strings.EqualFold(s.rels[i].name, qualifier) {
				match = &s.rels[i]
				count++
			}
		}
		if count != 1 || !match.derived || len(match.cols) == 0 || hasOpaque(match.cols) {
			return nil
		}
		return match.cols
	}
	var out []outCol
	for _, rb := range s.rels {
		if !rb.derived || len(rb.cols) == 0 || hasOpaque(rb.cols) {
			return nil
		}
		out = append(out, rb.cols...)
	}
	return out
}

// resolveRefs resolves a list of result-column refs through the scope,
// deduplicating the result while preserving first-seen order.
func (s *rscope) resolveRefs(refs []ColumnRef) []ColumnRef {
	var out []ColumnRef
	seen := make(map[ColumnRef]bool)
	for _, ref := range refs {
		for _, r := range s.resolveRef(ref) {
			if !seen[r] {
				seen[r] = true
				out = append(out, r)
			}
		}
	}
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
	if ref.Table != "" {
		// Qualified by a relation name: append the recovered sources of every
		// in-scope derived relation of that name exposing the column (a reused
		// alias yields several; unioning them over-includes, which is safe).
		for _, rb := range s.rels {
			if rb.derived && strings.EqualFold(rb.name, ref.Table) {
				if src, ok := lookupOutCol(rb.cols, ref.Column); ok {
					out = append(out, src...)
				}
			}
		}
		return out
	}
	// Bare reference: append the recovered sources of every in-scope derived
	// relation exposing a column of this name. The original bare ref is retained
	// so a base table providing the same column is still covered.
	for _, rb := range s.rels {
		if !rb.derived {
			continue
		}
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
func (s *rscope) resolveExprRefs(expr parser.Expr) []ColumnRef {
	var refs []ColumnRef       // direct column refs
	var subSources []ColumnRef // scalar-subquery output sources
	ew := exprWalk{
		followSub: false,
		onColumn:  func(ref ColumnRef) { refs = append(refs, ref) },
		onSubquery: func(sub *parser.SubqueryExpr) {
			subSources = append(subSources, scalarSubquerySources(sub, s.cte)...)
		},
	}
	ew.walk(expr)
	return unionRefs(s.resolveRefs(refs), s.resolveRefs(subSources))
}

// scalarSubquerySources recovers the resolved base-column sources of a scalar
// subquery used as a value (e.g. SELECT (SELECT phone FROM customer) AS sp). The
// subquery's raw text is re-parsed and its output columns resolved (recursively,
// in the enclosing CTE scope); the union of those columns' sources is the
// lineage feeding the enclosing expression. EXISTS subqueries yield a boolean and
// are skipped. Best-effort: a correlated reference to an outer column is not
// resolved to the outer relation (it resolves within the subquery's own scope).
func scalarSubquerySources(sub *parser.SubqueryExpr, cte *cteDefs) []ColumnRef {
	if sub == nil || sub.Kind != parser.SubqueryScalar {
		return nil
	}
	file, _ := parser.Parse(sub.RawText)
	if file == nil {
		return nil
	}
	var out []ColumnRef
	for _, stmt := range file.Stmts {
		if qs, ok := stmt.(*parser.QueryStmt); ok {
			for _, c := range resolveQueryCols(qs.Query, cte) {
				out = append(out, c.sources...)
			}
		}
	}
	return out
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
