package analysis

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// resolveDerivedLineage deepens result-column lineage through derived relations.
//
// GetQuerySpan's primary walk records each result column's source refs as the
// refs written in the select item (e.g. `d.x`, or a CTE column `pp`). When the
// referenced relation is a derived relation — a subquery in FROM or a CTE — that
// ref names the derived relation's output column, not a base table, so a
// downstream name-match against the FROM tables finds nothing and the column is
// left without lineage. Bytebase then has no source column to attach a masker
// to and returns the (possibly sensitive) value unmasked.
//
// This pass resolves those refs: for the outermost query it builds the set of
// derived relations visible in FROM, computes each derived relation's output
// columns and their underlying base-column refs (recursively, so a derived
// relation over another derived relation resolves), and ADDS the recovered base
// refs to every result ref that points at a derived relation. Resolution is
// additive — the original ref is always retained — so the resolved source set is
// a superset of the original and the pass can only ever deepen lineage (mask
// more), never drop a previously-correct ref (which would under-mask and leak).
// A ref that names a real base table (or an unrecognised relation) is left as-is.
//
// Scope: derived-relation projection only (subqueries in FROM and CTE
// references). Scalar subqueries in the select list, set-operation arm merging,
// and UNNEST are deliberately out of scope here (each tracked separately) and
// pass through unchanged.
func resolveDerivedLineage(stmt ast.Node, span *QuerySpan) {
	if span == nil {
		return
	}
	qs, ok := stmt.(*parser.QueryStmt)
	if !ok || qs.Query == nil {
		return
	}
	spec, cte := outermostScope(qs.Query, nil)
	if spec == nil {
		return
	}
	scope := newRScope(spec.From, cte)
	if !scope.hasDerived() {
		// No derived relations in the outermost FROM: there is nothing to
		// resolve, so leave Results byte-for-byte unchanged.
		return
	}
	for i := range span.Results {
		span.Results[i].SourceColumns = scope.resolveRefs(span.Results[i].SourceColumns)
	}
}

// ---------------------------------------------------------------------------
// Resolved output columns and CTE definitions
// ---------------------------------------------------------------------------

// outCol is one output column of a (sub)query or relation: its name and the
// base-column refs that feed it, recovered through any derived relations.
type outCol struct {
	name    string
	sources []ColumnRef
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
		return resolveNodeCols(n.Left, cte)
	case *parser.ParenQuery:
		return resolveQueryCols(n.Inner, cte)
	case *parser.TableQuery:
		// TABLE name == SELECT * FROM name.
		return []outCol{{name: "*"}}
	default:
		return nil
	}
}

// resolveSpecCols computes the resolved output columns of one SELECT block: it
// builds the FROM scope, then resolves each select item's direct column refs
// through that scope.
func resolveSpecCols(spec *parser.QuerySpec, cte *cteDefs) []outCol {
	if spec == nil {
		return nil
	}
	scope := newRScope(spec.From, cte)
	out := make([]outCol, 0, len(spec.Items))
	for _, item := range spec.Items {
		switch item.Kind {
		case parser.SelectAll:
			out = append(out, outCol{name: "*"})
			continue
		case parser.SelectAllFrom:
			name := renderExprName(item.Expr)
			if name == "" {
				name = "*"
			} else {
				name += ".*"
			}
			out = append(out, outCol{name: name, sources: scope.resolveExprRefs(item.Expr)})
			continue
		}
		name := identName(item.Alias)
		if name == "" {
			name = renderExprName(item.Expr)
		}
		out = append(out, outCol{name: name, sources: scope.resolveExprRefs(item.Expr)})
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

// hasDerived reports whether the scope contains any derived relation (so the
// caller can skip resolution entirely for base-table-only queries).
func (s *rscope) hasDerived() bool {
	for _, rb := range s.rels {
		if rb.derived {
			return true
		}
	}
	return false
}

// add binds one FROM relation (recursing through joins/parentheses). alias and
// colAliases carry an enclosing AliasedRelation's `AS name (c1, …)` down to the
// primary relation it names.
func (s *rscope) add(rel parser.Relation, alias string, colAliases []*ast.Identifier) {
	switch n := rel.(type) {
	case *parser.AliasedRelation:
		s.add(n.Inner, identName(n.Alias), n.ColumnAliases)
	case *parser.Join:
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
	default:
		// Lateral / UNNEST / table-function relations are not resolved here
		// (tracked separately). Bind the alias as a non-derived (opaque)
		// relation so references through it pass through unchanged.
		s.rels = append(s.rels, rbind{name: alias, derived: false})
	}
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

// resolveExprRefs collects the direct column refs of a select-item expression
// (not crossing subquery boundaries, matching the primary walk's
// collectDirectColumns) and resolves each through the scope.
func (s *rscope) resolveExprRefs(expr parser.Expr) []ColumnRef {
	var refs []ColumnRef
	ew := exprWalk{
		followSub: false,
		onColumn:  func(ref ColumnRef) { refs = append(refs, ref) },
	}
	ew.walk(expr)
	return s.resolveRefs(refs)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// outermostScope unwraps a query to the SELECT block whose select list produced
// the span's Results — descending through a set operation's left arm and
// parenthesised queries — and returns it together with the CTE scope visible
// there (accumulated across any WITH clauses encountered on the way down).
func outermostScope(q *parser.Query, parent *cteDefs) (*parser.QuerySpec, *cteDefs) {
	if q == nil {
		return nil, parent
	}
	cte := buildCTEDefs(q, parent)
	return specOf(q.Body, cte)
}

func specOf(node parser.QueryNode, cte *cteDefs) (*parser.QuerySpec, *cteDefs) {
	switch n := node.(type) {
	case *parser.QuerySpec:
		return n, cte
	case *parser.SetOperation:
		return specOf(n.Left, cte)
	case *parser.ParenQuery:
		return outermostScope(n.Inner, cte)
	default:
		return nil, cte
	}
}

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
