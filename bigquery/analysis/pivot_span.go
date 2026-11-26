package analysis

import (
	"fmt"

	"strings"

	"github.com/bytebase/omni/bigquery/ast"
)

// PIVOT and UNPIVOT replace their input's columns
// (https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#pivot_operator).
// Masking reads results by position, so both need the input column by column.

// pivotColumns is the output of PIVOT over in: the input columns that
// neither an aggregate nor the FOR expression reads, then, for each IN value,
// one column per aggregate, reading the aggregate's sources. Column
// references resolve against the current scope, which holds the input alone.
func (w *spanWalker) pivotColumns(in []projColumn, pv *ast.PivotClause) []projColumn {
	used := map[string]bool{}
	w.markRead(used, pv.For)
	var aggs [][]ColumnRef
	for _, agg := range pv.Aggregates {
		if agg == nil {
			continue
		}
		w.markRead(used, agg.Expr)
		_, sources := w.resolveExprSources(agg.Expr)
		aggs = append(aggs, sources)
	}
	var out []projColumn
	for _, c := range in {
		if !used[foldKey(c.name)] {
			out = append(out, c)
		}
	}
	for _, v := range pv.Values {
		if v == nil {
			continue
		}
		value := v.Alias
		if value == "" {
			value = pivotValueName(v.Expr)
		}
		if value == "" {
			// A reference to the column by its name would match nothing.
			w.failClosed(fmt.Errorf("a PIVOT value without an alias whose column name this walker does not derive (fail closed)"))
		}
		for i, agg := range pv.Aggregates {
			if agg == nil {
				continue
			}
			name := value
			if agg.Alias != "" {
				name = agg.Alias + "_" + value
			}
			out = append(out, projColumn{name: name, sources: aggs[i]})
		}
	}
	return out
}

// markRead adds to used, folded, the input column each column reference of
// expr names: its first part, its second where the first names a range
// variable, the column of a path that names a base table in full
// (`ds.t.c`), and every column of a range variable read whole.
func (w *spanWalker) markRead(used map[string]bool, expr ast.Node) {
	paths, _ := w.collectColumnPaths(expr)
	for _, parts := range paths {
		rel := w.leaves.find(parts[0])
		ref := columnRefFromParts(parts)
		switch {
		case rel == nil && w.qualifierNamesBaseRelation(ref):
			used[foldKey(ref.Column)] = true
		case rel != nil && len(parts) >= 2:
			used[foldKey(parts[1])] = true
		case rel != nil:
			cols, ok := w.pipeConcrete(rel.columns)
			if !ok {
				w.failClosed(fmt.Errorf("a PIVOT reads a range variable whose columns cannot be listed (fail closed)"))
			}
			for _, c := range cols {
				used[foldKey(c.name)] = true
			}
		default:
			used[foldKey(parts[0])] = true
		}
	}
}

// pivotValueName is the name BigQuery gives an IN value without an alias
// (https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax#pivot_operator):
// a string's content, `_1` for 1 and `minus_1` for -1, TRUE, FALSE, NULL,
// and `_2013_11_25` for DATE '2013-11-25'. Another value gets no name here.
func pivotValueName(expr ast.Node) string {
	if u, ok := expr.(*ast.UnaryExpr); ok && u.Op == ast.UnaryMinus {
		if lit, ok := u.Expr.(*ast.Literal); ok && lit.Kind == ast.LitInt && !strings.HasPrefix(strings.ToLower(lit.Value), "0x") {
			return "minus_" + lit.Value
		}
		return ""
	}
	if lit, ok := expr.(*ast.TypedLiteral); ok && strings.EqualFold(lit.TypeKeyword, "DATE") {
		return "_" + strings.ReplaceAll(strings.TrimSpace(lit.Value), "-", "_")
	}
	lit, ok := expr.(*ast.Literal)
	if !ok {
		return ""
	}
	switch lit.Kind {
	case ast.LitString:
		return lit.Value
	case ast.LitInt:
		if strings.HasPrefix(strings.ToLower(lit.Value), "0x") {
			return ""
		}
		return "_" + lit.Value
	case ast.LitBool:
		return strings.ToUpper(lit.Value)
	case ast.LitNull:
		return "NULL"
	}
	return ""
}

// unpivotColumns is the output of UNPIVOT over in: the input columns no IN
// item names, then each value column, reading the IN items' columns at its
// position, then the name column, which holds column names and reads no
// data.
func (w *spanWalker) unpivotColumns(in []projColumn, up *ast.UnpivotClause) []projColumn {
	used := map[string]bool{}
	values := make([]refUnion, len(up.ValueColumns))
	for _, item := range up.Items {
		if item == nil {
			continue
		}
		for i, col := range item.Columns {
			if col == nil || len(col.Parts) == 0 {
				continue
			}
			used[foldKey(col.Parts[len(col.Parts)-1])] = true
			if i < len(values) {
				_, refs := w.resolvePath(col.Parts)
				values[i].add(refs)
			}
		}
	}
	var out []projColumn
	for _, c := range in {
		if !used[foldKey(c.name)] {
			out = append(out, c)
		}
	}
	for i, vc := range up.ValueColumns {
		out = append(out, projColumn{name: lastPart(vc), sources: values[i].refs})
	}
	if up.NameColumn != nil {
		out = append(out, projColumn{name: lastPart(up.NameColumn)})
	}
	return out
}

func lastPart(p *ast.PathExpr) string {
	if p == nil || len(p.Parts) == 0 {
		return ""
	}
	return p.Parts[len(p.Parts)-1]
}

// pivotRelation builds a FROM source that PIVOT or UNPIVOT follows: the
// source is built in a scope of its own, which the clause's expressions
// resolve against, and the clause's output replaces it.
func (w *spanWalker) pivotRelation(n *ast.TableExpr) *relation {
	outer, outerCols := w.leaves, w.pipeColumns
	w.leaves, w.pipeColumns = &leafSet{}, nil
	defer func() { w.leaves, w.pipeColumns = outer, outerCols }()
	in := w.buildTableSource(n)
	if in == nil {
		return nil
	}
	cols, ok := w.pipeConcrete(in.columns)
	if !ok {
		w.failClosed(fmt.Errorf("PIVOT or UNPIVOT over a source whose columns cannot be listed (fail closed)"))
		return nil
	}
	name, alias := in.name, ""
	if n.Pivot != nil {
		cols, alias = w.pivotColumns(cols, n.Pivot), n.Pivot.Alias
	} else {
		cols, alias = w.unpivotColumns(cols, n.Unpivot), n.Unpivot.Alias
	}
	if alias != "" {
		name = strings.ToLower(alias)
	}
	return &relation{name: name, columns: cols}
}
