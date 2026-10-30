package analysis

import (
	"strconv"

	"github.com/bytebase/omni/googlesql/ast"
)

// recursiveFixpoint gives what the fixpoint loop of resolveCTEColumns gives,
// projection and side effects alike, without resolving the recursive arm once
// per pass.
//
// The loop resolves the arm against the CTE's projection so far, unions each
// arm column's sources into the same position, and stops after a pass that
// adds none. When the anchor and the arm are concrete, an arm column's sources
// are the ordered union of a fixed list of terms, each either a source read
// from elsewhere or the sources of CTE column j so far: resolving the arm only
// passes source lists along and takes ordered, deduplicated unions of them,
// and apart from a fixpoint of its own nothing it decides depends on what a
// list holds. An arm that rotates its columns (`SELECT c2, c3, ..., c1 FROM
// r`) moves a source one column per pass, so it needs a pass per column, and
// each pass resolved the whole arm and unioned every column's whole source
// list again: O(n^3) for n columns.
//
// Here the arm is resolved with each CTE column's sources replaced by one
// placeholder, which leaves the term lists as the arm's output, and the
// passes are replayed on the lists. The first pass reads every term. After
// it, a column holds every source a term read in an earlier pass, so a later
// pass can add only what the columns it reads gained in the pass before; only
// those gains are read, and only for the columns that read a column that
// gained, in the order the loop would add them. A second resolution with
// different placeholders confirms that no source read from elsewhere was
// taken for a placeholder.
//
// ok is false, with no side effect left that the loop does not repeat, when
// the replay would not be exact: the anchor or the arm is not concrete, their
// widths differ, an anchor column lists a source twice (the loop would drop
// the copy only in a pass that adds to that column), or the arm starts a
// fixpoint of its own (whose pass count depends on the sources it reads).
// An arm that holds a WITH RECURSIVE clause is left to the loop before it is
// resolved at all: resolving it here twice and then again in the loop would
// triple the work at every level of nested recursive CTEs, exponential in
// their depth, where the loop alone resolves an arm that converges at once a
// single time.
func (w *spanWalker) recursiveFixpoint(arm ast.Node, anchor []projColumn, scope *cteScope, lower string) ([]projColumn, bool) {
	n := len(anchor)
	if n == 0 || !allConcrete(anchor) || w.holdsRecursiveWith(arm) {
		return nil, false
	}
	cur := make([]refUnion, n) // each column's sources so far
	for i, c := range anchor {
		cur[i].add(c.sources)
		if len(cur[i].refs) != len(c.sources) {
			return nil, false
		}
	}

	// Resolving the arm records tables and predicate columns, which are
	// deduplicated, and may fail closed, which keeps the first failure; none
	// of that depends on the CTE's sources, so it is the same on every pass.
	// The names of CTEs the arm defines are appended on every pass: those of
	// the second resolution are removed, and the loop's other passes add them
	// below.
	ctes, fixpoints := len(w.span.CTEs), w.fixpoints
	projA, phA := placeholderProjection(anchor, "a")
	scope.columns[lower] = projA
	armA := w.visitBody(arm)
	perPass := len(w.span.CTEs) - ctes
	projB, phB := placeholderProjection(anchor, "b")
	scope.columns[lower] = projB
	armB := w.visitBody(arm)
	w.span.CTEs = w.span.CTEs[:ctes+perPass]
	scope.columns[lower] = anchor
	terms, ok := armTerms(armA, armB, phA, phB, n)
	if !ok || w.fixpoints != fixpoints {
		w.span.CTEs = w.span.CTEs[:ctes]
		return nil, false
	}

	// The first pass: every column reads every term in full.
	type gain struct{ col, lo, hi int }
	var last []gain // what the last pass added, by column
	for i := range terms {
		start := len(cur[i].refs)
		for _, t := range terms[i] {
			if t.col < 0 {
				cur[i].add([]ColumnRef{t.ref})
			} else {
				cur[i].add(anchor[t.col].sources)
			}
		}
		if end := len(cur[i].refs); end > start {
			last = append(last, gain{i, start, end})
		}
	}
	if len(last) == 0 {
		return anchor, true
	}

	// The later passes, up to the loop's bound.
	readers := make([][]int, n) // column j → the columns with a term reading j
	for i, ts := range terms {
		for _, t := range ts {
			if t.col >= 0 {
				readers[t.col] = append(readers[t.col], i)
			}
		}
	}
	gained := make([]gain, n) // gained[j]: what column j gained in the last pass, if it did
	inLast := make([]int, n)  // inLast[j] == pass when it did
	queued := make([]int, n)  // queued[i] == pass when column i is to be recomputed
	passes := 1
	for pass := 1; pass <= n+1 && len(last) > 0; pass++ {
		passes++
		var todo []int
		for _, g := range last {
			gained[g.col], inLast[g.col] = g, pass
			for _, i := range readers[g.col] {
				if queued[i] != pass {
					queued[i] = pass
					todo = append(todo, i)
				}
			}
		}
		last = last[:0]
		for _, i := range todo {
			start := len(cur[i].refs)
			for _, t := range terms[i] {
				if t.col >= 0 && inLast[t.col] == pass {
					g := gained[t.col]
					cur[i].add(cur[t.col].refs[g.lo:g.hi])
				}
			}
			if end := len(cur[i].refs); end > start {
				last = append(last, gain{i, start, end})
			}
		}
	}
	for p := 1; p < passes && perPass > 0; p++ {
		w.span.CTEs = append(w.span.CTEs, w.span.CTEs[ctes:ctes+perPass]...)
	}

	cols := make([]projColumn, n)
	for i, c := range anchor {
		sources := cur[i].refs
		if len(sources) == len(c.sources) {
			sources = c.sources
		}
		cols[i] = projColumn{name: c.name, sources: sources}
	}
	scope.columns[lower] = cols
	return cols, true
}

// holdsRecursiveWith reports whether n's subtree holds a WITH RECURSIVE
// clause, the only way resolving n can start a fixpoint: a CTE's columns are
// resolved where its WITH clause is, not where it is read. The answer for
// every node of the scanned subtree is recorded, so the arms of recursive
// CTEs nested in one another are scanned once in all, not once per level.
func (w *spanWalker) holdsRecursiveWith(n ast.Node) bool {
	if held, ok := w.recursiveWiths[n]; ok {
		return held
	}
	if w.recursiveWiths == nil {
		w.recursiveWiths = make(map[ast.Node]bool)
	}
	ast.Walk(&recursiveWithScan{seen: w.recursiveWiths}, n)
	return w.recursiveWiths[n]
}

// recursiveWithScan records in seen, for each node it walks, whether the
// node's subtree holds a WITH RECURSIVE clause. stack holds the nodes whose
// subtrees are being walked.
type recursiveWithScan struct {
	seen  map[ast.Node]bool
	stack []ast.Node
}

func (s *recursiveWithScan) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		top := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		s.markParent(s.seen[top])
		return nil
	}
	if held, ok := s.seen[n]; ok {
		s.markParent(held)
		return nil
	}
	with, ok := n.(*ast.WithClause)
	s.seen[n] = ok && with.Recursive
	s.stack = append(s.stack, n)
	return s
}

// markParent records that the node being walked holds a WITH RECURSIVE
// clause when its child, just walked, does.
func (s *recursiveWithScan) markParent(held bool) {
	if held && len(s.stack) > 0 {
		s.seen[s.stack[len(s.stack)-1]] = true
	}
}

// armTerm is one term of an arm column's sources: the sources of CTE column
// col, or, when col is -1, the source ref.
type armTerm struct {
	col int
	ref ColumnRef
}

// placeholderProjection returns anchor with each column's sources replaced by
// one placeholder ref of the given family, and the column of each
// placeholder. A placeholder has an empty Column, which no source has (the
// parser rejects an empty identifier); armTerms checks that anyway.
func placeholderProjection(anchor []projColumn, family string) ([]projColumn, map[ColumnRef]int) {
	cols := make([]projColumn, len(anchor))
	index := make(map[ColumnRef]int, len(anchor))
	for j, c := range anchor {
		ref := ColumnRef{Catalog: family, Table: strconv.Itoa(j)}
		c.sources = []ColumnRef{ref}
		cols[j] = c
		index[ref] = j
	}
	return cols, index
}

// armTerms reads each arm column's term list off two resolutions of the arm
// against placeholder projections of different families (phA and phB). The
// resolutions must agree term by term: a placeholder of column j in one where
// the other has column j's, or the same ref with a non-empty Column in both.
// ok is false when they do not, or when an arm is not n concrete columns.
func armTerms(armA, armB []projColumn, phA, phB map[ColumnRef]int, n int) ([][]armTerm, bool) {
	if len(armA) != n || len(armB) != n || !allConcrete(armA) || !allConcrete(armB) {
		return nil, false
	}
	terms := make([][]armTerm, n)
	for i := range armA {
		a, b := armA[i].sources, armB[i].sources
		if len(a) != len(b) {
			return nil, false
		}
		ts := make([]armTerm, len(a))
		for p := range a {
			if j, ok := phA[a[p]]; ok {
				if k, ok := phB[b[p]]; !ok || k != j {
					return nil, false
				}
				ts[p] = armTerm{col: j}
				continue
			}
			if a[p] != b[p] || a[p].Column == "" {
				return nil, false
			}
			ts[p] = armTerm{col: -1, ref: a[p]}
		}
		terms[i] = ts
	}
	return terms, true
}
