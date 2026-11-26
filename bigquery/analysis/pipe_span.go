package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/bigquery/ast"
)

// Pipe operators that change the projection.
//
// A pipe operator reads one table, the output of the query or operator before
// it, and returns another. pipeTable is that table: its columns in order and
// the range variables (table aliases) still visible. BigQuery keeps an input's
// range variables through WHERE, EXTEND, SET, DROP, RENAME, JOIN and the
// other operators that keep its rows' identity, and a range variable still
// reads the original row: after `|> SET x = y`, `t.x` is the old x. SELECT,
// AGGREGATE and the set operations return a new table with none, and AS
// replaces them all with one.
//
// Masking reads results by position, so every operator here needs its input
// as concrete columns: a base table's star is expanded from metadata
// (spanWalker.columnsOf), a star's EXCEPT and REPLACE are applied, and an
// input omni cannot list column by column (a deferred set-operation merge, a
// table without metadata) fails closed. So do the operators whose output omni does not
// derive: PIVOT, UNPIVOT, CALL, MATCH_RECOGNIZE, a recursive union, and the
// control and write operators.
//
// The answers below were checked against BigQuery's dry run on inline tables
// (`FROM (SELECT 1 AS x, 2 AS y) AS t |> ...`); testdata/span.tsv holds the
// recorded cases.

// pipeTable is the table between two pipe operators.
type pipeTable struct {
	cols   []projColumn
	ranges []*relation // range variables, each a named FROM leaf
}

// pipeOperator applies one pipe operator to in and returns its output, or nil
// when the operator fails closed. Either way the operator's tables are
// recorded.
func (w *spanWalker) pipeOperator(op ast.Node, in *pipeTable) *pipeTable {
	switch n := op.(type) {
	case *ast.PipeWhere, *ast.PipeOrderBy, *ast.PipeLimitOffset, *ast.PipeDistinct,
		*ast.PipeTablesample, *ast.PipeAssert:
		// Rows change, columns and range variables do not.
		restore := w.enterPipeScope(in)
		w.visitPipeOp(op)
		restore()
		return in
	case *ast.PipeLog:
		w.visitPipeOp(op)
		if !pipeKeepsProjection(n.Pipeline) {
			return w.pipeFailClosed("a pipe LOG whose subpipeline changes the projection")
		}
		return in
	}

	cols, ok := w.pipeConcrete(in.cols)
	if !ok {
		w.visitPipeOp(op)
		return w.pipeFailClosed("a pipe operator over an input whose columns cannot be listed")
	}
	if n, ok := op.(*ast.PipeSetOperation); ok {
		// The operands are queries of their own: the input's range variables
		// are not visible to them.
		return w.pipeSetOperation(cols, n)
	}
	cur := &pipeTable{cols: cols, ranges: in.ranges}
	restore := w.enterPipeScope(cur)
	defer restore()

	switch n := op.(type) {
	case *ast.PipeAs:
		alias := &relation{name: strings.ToLower(n.Alias), columns: cols}
		return &pipeTable{cols: cols, ranges: []*relation{alias}}
	case *ast.PipeSelect:
		w.visitPipeOp(op)
		if n.Select == nil || n.Select.As != ast.SelectAsNone {
			return w.pipeFailClosed("a pipe SELECT AS STRUCT or AS VALUE")
		}
		if w.starOverValueElem(n.Select.Items, cols) {
			return w.pipeFailClosed("a pipe SELECT star over an UNNEST element, which expands a struct element's fields,")
		}
		return &pipeTable{cols: w.resolveSelectProjection(n.Select, []*relation{{columns: cols}})}
	case *ast.PipeExtend:
		w.visitPipeOp(op)
		if n.Select == nil {
			return w.pipeFailClosed("an empty pipe EXTEND")
		}
		return w.pipeExtend(cur, n.Select.Items)
	case *ast.PipeWindow:
		w.visitPipeOp(op)
		return w.pipeExtend(cur, n.Items)
	case *ast.PipeAggregate:
		w.visitPipeOp(op)
		return w.pipeAggregate(n.Select)
	case *ast.PipeSet:
		w.visitPipeOp(op)
		return w.pipeSet(cur, n.Items)
	case *ast.PipeDrop:
		return w.pipeDrop(cur, n.Columns)
	case *ast.PipeRename:
		return w.pipeRename(cur, n.Items)
	case *ast.PipeJoin:
		return w.pipeJoin(cur, n.Join)
	}
	w.visitPipeOp(op)
	return w.pipeFailClosed(fmt.Sprintf("the pipe operator %T", op))
}

// starOverValueElem reports whether a star item of items expands an UNNEST
// element column: a bare star over cols, or `rel.*` over a range variable
// holding one.
func (w *spanWalker) starOverValueElem(items []*ast.SelectItem, cols []projColumn) bool {
	for _, item := range items {
		if item == nil || !item.Star {
			continue
		}
		over := cols
		if item.Expr != nil {
			over = nil
			if parts := exprToParts(item.Expr); len(parts) == 1 {
				if rel := w.leaves.find(parts[0]); rel != nil {
					over = rel.columns
				}
			}
		}
		for _, c := range over {
			if c.valueElem {
				return true
			}
		}
	}
	return false
}

// pipeFailClosed fails the span closed for a pipe operator omni does not
// resolve and returns the nil table that stops the resolution.
func (w *spanWalker) pipeFailClosed(what string) *pipeTable {
	w.failClosed(fmt.Errorf("%s changes the projection in a way this walker does not resolve (fail closed)", what))
	return nil
}

// enterPipeScope makes t the table column references resolve against: an
// unqualified name reads t's columns, and a qualified one a range variable. It
// returns the function that restores the previous scope.
func (w *spanWalker) enterPipeScope(t *pipeTable) func() {
	prevLeaves, prevFrom, prevUnresolved := w.leaves, w.fromAliases, w.fromUnresolved
	prevCols := w.pipeColumns
	w.leaves = &leafSet{}
	w.fromAliases = map[string]bool{}
	w.fromUnresolved = false
	w.pipeColumns = map[string]bool{}
	for _, c := range t.cols {
		w.pipeColumns[foldKey(c.name)] = true
	}
	w.leaves.add(&relation{columns: t.cols})
	for _, rel := range t.ranges {
		w.leaves.add(rel)
		w.registerFromAlias(rel.name)
	}
	return func() {
		w.leaves, w.fromAliases, w.fromUnresolved = prevLeaves, prevFrom, prevUnresolved
		w.pipeColumns = prevCols
	}
}

// pipeConcrete lists cols column by column, expanding each base table's star
// from metadata. It reports false when a column stands for columns omni
// cannot list.
func (w *spanWalker) pipeConcrete(cols []projColumn) ([]projColumn, bool) {
	out := make([]projColumn, 0, len(cols))
	for _, c := range cols {
		switch {
		case c.baseStar != nil:
			if w.columnsOf == nil {
				return nil, false
			}
			names, ok := w.columnsOf(*c.baseStar)
			if !ok {
				return nil, false
			}
			for _, name := range names {
				if containsFold(c.baseStarExcept, name) {
					continue
				}
				ref := *c.baseStar
				ref.Column = name
				out = append(out, projColumn{name: name, sources: []ColumnRef{ref}, plain: c.plain})
			}
		case c.starGroup != nil:
			expanded, ok := w.pipeStarGroup(c.starGroup)
			if !ok {
				return nil, false
			}
			out = append(out, expanded...)
		case c.starMerge != nil || c.setOpMerge != nil:
			return nil, false
		default:
			out = append(out, c)
		}
	}
	return out, true
}

// pipeStarGroup lists a star with EXCEPT or REPLACE column by column: its
// segments expanded, each EXCEPT name dropped, each REPLACE name re-pointed to
// its value. It reports false when two segments share a name, where the
// modifiers' effect depends on which one they meet.
func (w *spanWalker) pipeStarGroup(g *starGroup) ([]projColumn, bool) {
	segs, ok := w.pipeConcrete(g.segments)
	if !ok {
		return nil, false
	}
	seen := map[string]bool{}
	for _, c := range segs {
		key := strings.ToLower(c.name)
		if c.name == "" || seen[key] {
			return nil, false
		}
		seen[key] = true
	}
	out := make([]projColumn, 0, len(segs))
	for _, c := range segs {
		if containsFold(g.except, c.name) {
			continue
		}
		for _, r := range g.replace {
			if strings.EqualFold(r.Name, c.name) {
				c = projColumn{name: r.Name, sources: r.Sources}
			}
		}
		out = append(out, c)
	}
	return out, true
}

// pipeExtend appends one column per item (EXTEND, WINDOW). The items read the
// input, not each other, and the range variables stay.
func (w *spanWalker) pipeExtend(in *pipeTable, items []*ast.SelectItem) *pipeTable {
	cols := append([]projColumn{}, in.cols...)
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.Star {
			return w.pipeFailClosed("a star in a pipe EXTEND")
		}
		cols = append(cols, w.resolveExprColumn(item))
	}
	return &pipeTable{cols: cols, ranges: in.ranges}
}

// pipeAggregate returns the grouping columns, then the aggregate columns: the
// order BigQuery's AGGREGATE gives them. A ROLLUP, CUBE or GROUPING SETS
// contributes each distinct column it groups by, in order of first appearance.
func (w *spanWalker) pipeAggregate(sel *ast.SelectStmt) *pipeTable {
	if sel == nil {
		return w.pipeFailClosed("an empty pipe AGGREGATE")
	}
	var cols []projColumn
	if sel.GroupBy != nil {
		if sel.GroupBy.Kind != ast.GroupByItems {
			return w.pipeFailClosed("a pipe AGGREGATE with GROUP BY ALL")
		}
		for _, gi := range sel.GroupBy.Items {
			if gi == nil {
				continue
			}
			switch gi.Kind {
			case ast.GroupingExpr:
				_, sources := w.resolveExprSources(gi.Expr)
				name := implicitAlias(gi.Expr)
				if gi.Alias != "" {
					name = gi.Alias
				}
				cols = append(cols, projColumn{name: name, sources: sources})
			case ast.GroupingEmpty:
			case ast.GroupingRollup, ast.GroupingCube, ast.GroupingSets:
				seen := map[string]bool{}
				for _, e := range gi.Items {
					parts := exprToParts(e)
					if len(parts) == 0 {
						return w.pipeFailClosed("a pipe AGGREGATE grouping set of expressions")
					}
					key := strings.ToLower(strings.Join(parts, "."))
					if seen[key] {
						continue
					}
					seen[key] = true
					_, sources := w.resolvePath(parts)
					cols = append(cols, projColumn{name: parts[len(parts)-1], sources: sources})
				}
			default:
				return w.pipeFailClosed("a pipe AGGREGATE grouping item")
			}
		}
	}
	for _, item := range sel.Items {
		if item == nil {
			continue
		}
		if item.Star {
			return w.pipeFailClosed("a star in a pipe AGGREGATE")
		}
		cols = append(cols, w.resolveExprColumn(item))
	}
	return &pipeTable{cols: cols}
}

// pipeColumnIndex returns the position of the one input column named name
// (case-insensitive), or -1 when none or several are. SET and RENAME require
// exactly one.
func pipeColumnIndex(cols []projColumn, name string) int {
	at := -1
	for i, c := range cols {
		if strings.EqualFold(c.name, name) {
			if at >= 0 {
				return -1
			}
			at = i
		}
	}
	return at
}

// pipeSet replaces each named column's value in place. Every value reads the
// input, and the range variables keep reading the original columns.
func (w *spanWalker) pipeSet(in *pipeTable, items []*ast.PipeSetItem) *pipeTable {
	cols := append([]projColumn{}, in.cols...)
	for _, item := range items {
		if item == nil {
			continue
		}
		i := pipeColumnIndex(in.cols, item.Column)
		if i < 0 {
			return w.pipeFailClosed("a pipe SET of a column that is not exactly one input column")
		}
		_, sources := w.resolveExprSources(item.Value)
		cols[i] = projColumn{name: in.cols[i].name, sources: sources}
	}
	return &pipeTable{cols: cols, ranges: in.ranges}
}

// pipeDrop removes every column with each name; the range variables keep them.
func (w *spanWalker) pipeDrop(in *pipeTable, names []string) *pipeTable {
	var cols []projColumn
	for _, name := range names {
		if !containsColumnFold(in.cols, name) {
			return w.pipeFailClosed("a pipe DROP of a column the input lacks")
		}
	}
	for _, c := range in.cols {
		if !containsFold(names, c.name) {
			cols = append(cols, c)
		}
	}
	return &pipeTable{cols: cols, ranges: in.ranges}
}

// pipeRename renames each column in place; the range variables keep the old
// names.
func (w *spanWalker) pipeRename(in *pipeTable, items []*ast.PipeRenameItem) *pipeTable {
	cols := append([]projColumn{}, in.cols...)
	for _, item := range items {
		if item == nil {
			continue
		}
		i := pipeColumnIndex(in.cols, item.Old)
		if i < 0 {
			return w.pipeFailClosed("a pipe RENAME of a column that is not exactly one input column")
		}
		cols[i].name = item.New
		cols[i].baseFieldName = false
	}
	return &pipeTable{cols: cols, ranges: in.ranges}
}

// pipeJoin joins the input, as the left side, to the operator's table. The
// output is the left columns then the right ones; USING puts each key first,
// once, reading both sides, as a FROM join does. The input's range variables
// stay and the right side's join them.
func (w *spanWalker) pipeJoin(in *pipeTable, join *ast.JoinExpr) *pipeTable {
	if join == nil {
		return w.pipeFailClosed("an empty pipe JOIN")
	}
	known := len(w.leaves.rels)
	rel := w.buildFromItem(join.Right)
	if te, ok := join.Right.(*ast.TableExpr); ok && te != nil && te.Subquery != nil && isValueTableQuery(te.Subquery) {
		w.recordSourceTVFArgs(join.Right)
		w.walkPredicate(join.On)
		return w.pipeFailClosed("a pipe JOIN to a SELECT AS STRUCT or AS VALUE subquery")
	}
	if sourceReshapes(join.Right) {
		w.recordSourceTVFArgs(join.Right)
		w.walkPredicate(join.On)
		return w.pipeFailClosed("a pipe JOIN to a PIVOT or UNPIVOT source")
	}
	w.recordSourceTVFArgs(join.Right)
	w.walkPredicate(join.On)
	for _, col := range join.Using {
		if col != "" {
			w.addPredicateColumn(ColumnRef{Column: col})
		}
	}
	if join.Natural {
		return w.pipeFailClosed("a pipe NATURAL JOIN")
	}
	if join.MatchRecognize != nil {
		w.walkMatchRecognize(join.MatchRecognize)
		return w.pipeFailClosed("MATCH_RECOGNIZE on a pipe JOIN")
	}
	if rel == nil || w.fromUnresolved {
		return w.pipeFailClosed("a pipe JOIN to a source without resolvable columns")
	}
	right, ok := w.pipeConcrete(rel.columns)
	if !ok {
		return w.pipeFailClosed("a pipe JOIN to a source whose columns cannot be listed")
	}
	var cols []projColumn
	if len(join.Using) > 0 {
		cols, _, ok = coalesceUsingColumns(in.cols, right, join.Using)
		if !ok {
			return w.pipeFailClosed("a pipe JOIN USING")
		}
	} else {
		cols = make([]projColumn, 0, len(in.cols)+len(right))
		for _, c := range in.cols {
			cols = append(cols, rewrapNonPlain(c))
		}
		cols = append(cols, right...)
	}
	ranges := append([]*relation{}, in.ranges...)
	for _, leaf := range w.leaves.rels[known:] {
		if leaf.name != "" {
			ranges = append(ranges, leaf)
		}
	}
	return &pipeTable{cols: cols, ranges: ranges}
}

// pipeSetOperation merges the input with each operand, left to right, as a set
// operation does: by position, or by name for BY NAME and CORRESPONDING. The
// output has no range variables.
func (w *spanWalker) pipeSetOperation(cols []projColumn, n *ast.PipeSetOperation) *pipeTable {
	set := n.Set
	if set == nil {
		for _, rhs := range n.Right {
			w.visitPipeOperand(rhs)
		}
		return w.pipeFailClosed("a pipe set operation")
	}
	// An outer or INNER mode, and a CORRESPONDING without STRICT, keep a
	// subset or a union of the inputs' names that mergeByName does not
	// compute (`|> INNER UNION ALL BY NAME` keeps only the shared columns).
	if set.OuterMode != "" || (set.Corresponding && !set.Strict) {
		for _, rhs := range n.Right {
			w.visitPipeOperand(rhs)
		}
		return w.pipeFailClosed("a pipe set operation that matches a subset of the inputs' names")
	}
	m := setOpMerger{cols: cols}
	for _, rhs := range n.Right {
		if isValueTableQuery(rhs) {
			w.visitPipeOperand(rhs)
			return w.pipeFailClosed("a pipe set operation with a SELECT AS STRUCT or AS VALUE operand")
		}
		right := w.visitBody(rhs)
		if set.ByName || set.Corresponding {
			m.mergeByName(right, set.MatchColumns)
		} else {
			m.merge(right)
		}
	}
	return &pipeTable{cols: m.result()}
}

// containsColumnFold reports whether some column of cols is named name,
// case-insensitively.
func containsColumnFold(cols []projColumn, name string) bool {
	for _, c := range cols {
		if strings.EqualFold(c.name, name) {
			return true
		}
	}
	return false
}

// sourceReshapes reports whether a FROM source carries a PIVOT or UNPIVOT,
// whose output columns buildFromItem does not derive: it returns the
// source's columns before the operator.
func sourceReshapes(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.TableExpr:
		return n != nil && (n.Pivot != nil || n.Unpivot != nil)
	case *ast.UnnestExpr:
		return n != nil && (n.Pivot != nil || n.Unpivot != nil)
	case *ast.JoinExpr:
		return n != nil && (sourceReshapes(n.Left) || sourceReshapes(n.Right))
	}
	return false
}

// isValueTableQuery reports whether a query body returns a value table: a
// SELECT AS STRUCT or AS VALUE (through parentheses, and a set operation's
// first operand). Its rows are one value, not the select list's columns that
// visitBody returns.
func isValueTableQuery(node ast.Node) bool {
	for {
		switch n := node.(type) {
		case *ast.SelectStmt:
			return n != nil && n.As != ast.SelectAsNone
		case *ast.QueryStmt:
			if n == nil {
				return false
			}
			node = n.Body
		case *ast.SetOperation:
			if n == nil {
				return false
			}
			node = n.Left
		default:
			return false
		}
	}
}
