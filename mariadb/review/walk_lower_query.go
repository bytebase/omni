package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// viewQueryFails checks the query of a view CREATE or ALTER VIEW makes,
// which MariaDB resolves when it stores the view (verified on
// mariadb:11.8.8 by mariadb/catalog/testdata/queries/simulate.tsv):
//
//   - a bare column with no table to read it from (1054), FROM DUAL being
//     none, and a column of the one base table a query block reads that the
//     table lacks (1054);
//   - `*` with no table (1096), and `t.*` with none (1051);
//   - two output columns of one name, when the view names none (1060);
//   - a row as an output column (1241);
//   - a window a function names that its query block does not define
//     (4009), and a frame on a function that takes none (4017);
//   - a non-constant `ESCAPE` or `AGAINST` (1210);
//   - a missing relation it reads (1146), and NEXTVAL, LASTVAL, SETVAL,
//     NEXT VALUE FOR, or PREVIOUS VALUE FOR of a missing relation (1146) or
//     one that is no sequence (4089).
//
// named is set when the view lists its column names.
func (b *builder) viewQueryFails(sel *ast.SelectStmt, named bool) bool {
	if sel == nil {
		return false
	}
	ctes := make(map[string]bool)
	ast.Inspect(sel, func(n ast.Node) bool {
		if c, ok := n.(*ast.CommonTableExpr); ok {
			ctes[strings.ToLower(c.Name)] = true
		}
		return true
	})
	top := ast.UnwrapParenSource(sel)
	if top != nil && top.SetOp == 0 && top.ValuesSource == nil {
		if b.outputFails(top, named) {
			return true
		}
	}
	for _, leaf := range queryLeaves(sel) {
		if b.tablelessFails(leaf) {
			return true
		}
	}
	if top != nil && top.SetOp == 0 && top.ValuesSource == nil && len(ctes) == 0 {
		b.singleTableColumns(top)
	}
	return b.blocksFail(sel, ctes)
}

// queryLeaves returns the query blocks whose rows a query returns: each
// operand of its set operations.
func queryLeaves(sel *ast.SelectStmt) []*ast.SelectStmt {
	sel = ast.UnwrapParenSource(sel)
	if sel == nil {
		return nil
	}
	if sel.SetOp != 0 {
		return append(queryLeaves(sel.Left), queryLeaves(sel.Right)...)
	}
	if sel.ValuesSource != nil {
		return nil
	}
	return []*ast.SelectStmt{sel}
}

// isDual reports whether a FROM clause reads no table: none, or DUAL.
func (b *builder) isDual(from []ast.TableExpr) bool {
	if len(from) == 0 {
		return true
	}
	if len(from) != 1 {
		return false
	}
	tr, ok := from[0].(*ast.TableRef)
	if !ok || tr.Schema != "" || tr.Alias != "" || tr.Loc.Start < 0 || tr.Loc.End > len(b.w.sql) {
		return false
	}
	// DUAL is a keyword; a quoted `dual` is a table.
	return strings.EqualFold(b.w.sql[tr.Loc.Start:tr.Loc.End], "DUAL")
}

// tablelessFails fails a query block that reads no table and names a
// column (1054) or a star (1096, 1051) in its select list.
func (b *builder) tablelessFails(sel *ast.SelectStmt) bool {
	if !b.isDual(sel.From) {
		return false
	}
	for _, target := range sel.TargetList {
		expr := target
		if rt, ok := target.(*ast.ResTarget); ok {
			expr = rt.Val
		}
		switch x := expr.(type) {
		case *ast.StarExpr:
			b.fail(rangeOf(x.Loc), nil, review.UnsupportedForm, "This query reads no table for "+review.Code("*")+".") // 1096
			return true
		case *ast.ColumnRef:
			if x.Star {
				b.fail(rangeOf(x.Loc), nil, review.ObjectNotFound, "Table "+review.Code(ident(x.Table))+" is not in this query.") // 1051
				return true
			}
		}
		if ref := bareColumnRefOutsideFuncs(expr); ref != nil {
			b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
			return true
		}
	}
	return false
}

// outputFails fails a view's query block whose output columns MariaDB
// refuses: a row (1241), or two of one name when the view lists none
// (1060). An item the server names after its text is left out.
func (b *builder) outputFails(sel *ast.SelectStmt, named bool) bool {
	seen := make(map[string]bool)
	for _, target := range sel.TargetList {
		expr, name := target, ""
		if rt, ok := target.(*ast.ResTarget); ok {
			expr = rt.Val
			if rt.NameWritten {
				name = rt.Name
			}
		}
		if row, ok := expr.(*ast.RowExpr); ok && len(row.Items) > 1 {
			b.fail(rangeOf(row.Loc), nil, review.UnsupportedForm, "A view's column cannot be a row.") // 1241
			return true
		}
		if !named && name == "" {
			if ref, ok := expr.(*ast.ColumnRef); ok && !ref.Star {
				name = ref.Column
			}
		}
		if named || name == "" {
			continue
		}
		k, ok := b.key(spaceColumn, name)
		if !ok {
			continue
		}
		if seen[k] {
			b.fail(b.anchor(), nil, review.DuplicateName, statesTwice("names", "column", ident(name))) // 1060
			return true
		}
		seen[k] = true
	}
	return false
}

// singleTableColumns checks the columns a query block that reads one base
// table names in its select list and WHERE clause, outside subqueries and
// function calls: each is the table's (1054). An alias of the select list
// is no column there.
func (b *builder) singleTableColumns(sel *ast.SelectStmt) {
	if len(sel.From) != 1 || b.isDual(sel.From) {
		return
	}
	tr, ok := sel.From[0].(*ast.TableRef)
	if !ok || tr.SystemTime != nil {
		return
	}
	r, st := b.ref(tr)
	if st != refOK {
		return
	}
	f, temp, kind := b.resolve(r)
	if temp || f != walk.Present || kind != relTable || b.w.versioned[relKey(r.db, r.name)] != walk.Absent {
		return
	}
	qualifier := tr.Name
	if tr.Alias != "" {
		qualifier = tr.Alias
	}
	exprs := make([]ast.ExprNode, 0, len(sel.TargetList)+1)
	for _, target := range sel.TargetList {
		if rt, ok := target.(*ast.ResTarget); ok {
			target = rt.Val
		}
		exprs = append(exprs, target)
	}
	exprs = append(exprs, sel.Where)
	for _, expr := range exprs {
		for _, ref := range columnRefsOutsideFuncs(expr) {
			if ref.Schema != "" || ref.Table != "" && !strings.EqualFold(ref.Table, qualifier) {
				continue
			}
			k, ok := b.key(spaceColumn, ref.Column)
			if !ok {
				continue
			}
			b.need(colKey(r.db, r.name, k), walk.Present, nil, rangeOf(ref.Loc), notFound("column", ident(ref.Column), r.written), "") // 1054
		}
	}
}

// columnRefsOutsideFuncs returns the column references of an expression
// outside subqueries and function calls, whose arguments may hold a
// sequence or a temporal unit.
func columnRefsOutsideFuncs(expr ast.ExprNode) []*ast.ColumnRef {
	if expr == nil {
		return nil
	}
	var out []*ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt, *ast.FuncCallExpr, *ast.SubqueryExpr, *ast.ExistsExpr:
			return false
		case *ast.ColumnRef:
			if !x.Star && x.Column != "" {
				out = append(out, x)
			}
			return false
		}
		return true
	})
	return out
}

// frameless lists the window functions that take no frame
// (Item_window_func::is_frame_prohibited).
var frameless = map[string]bool{
	"ROW_NUMBER": true, "RANK": true, "DENSE_RANK": true,
	"PERCENT_RANK": true, "CUME_DIST": true, "NTILE": true,
}

// blocksFail checks each query block of a query for the window, constant,
// and sequence failures viewQueryFails lists.
func (b *builder) blocksFail(root *ast.SelectStmt, ctes map[string]bool) bool {
	queue := []*ast.SelectStmt{root}
	for len(queue) > 0 {
		sel := queue[0]
		queue = queue[1:]
		windows := make(map[string]bool)
		for _, w := range sel.WindowClause {
			if w != nil {
				windows[strings.ToLower(w.Name)] = true
			}
		}
		failed := false
		ast.Inspect(sel, func(n ast.Node) bool {
			if failed {
				return false
			}
			switch x := n.(type) {
			case *ast.SelectStmt:
				if x != sel {
					queue = append(queue, x)
					return false
				}
			case *ast.FuncCallExpr:
				failed = b.callFails(x, windows, ctes)
			case *ast.LikeExpr:
				if x.Escape != nil && len(columnRefsOutsideQueries(x.Escape)) > 0 {
					b.fail(rangeOf(x.Loc), nil, review.UnsupportedForm, "The "+review.Code("ESCAPE")+" of "+review.Code("LIKE")+" must be a constant.") // 1210
					failed = true
				}
			case *ast.MatchExpr:
				if len(columnRefsOutsideQueries(x.Against)) > 0 {
					b.fail(rangeOf(x.Loc), nil, review.UnsupportedForm, "The "+review.Code("AGAINST")+" of "+review.Code("MATCH")+" must be a constant.") // 1210
					failed = true
				}
			case *ast.NextValueForExpr:
				b.sequenceNeeds(x.Sequence, ctes)
				return false
			case *ast.PreviousValueForExpr:
				b.sequenceNeeds(x.Sequence, ctes)
				return false
			case *ast.TableRef:
				b.sourceNeeds(x, ctes)
			}
			return !failed
		})
		if failed {
			return true
		}
	}
	return false
}

// callFails checks a function call: its window (4009, 4017), and the
// sequence a sequence function names.
func (b *builder) callFails(x *ast.FuncCallExpr, windows, ctes map[string]bool) bool {
	if x.Over != nil {
		if name := x.Over.RefName; name != "" && !windows[strings.ToLower(name)] {
			b.fail(rangeOf(x.Over.Loc), nil, review.ObjectNotFound, "Window "+review.Code(ident(name))+" is not defined.") // 4009
			return true
		}
		if x.Over.Frame != nil && x.Schema == "" && !x.NameQuoted && frameless[strings.ToUpper(x.Name)] {
			b.fail(rangeOf(x.Over.Loc), nil, review.UnsupportedForm, "Window function "+review.Code(strings.ToUpper(x.Name))+" takes no frame.") // 4017
			return true
		}
	}
	if x.Schema != "" || x.NameQuoted || len(x.Args) == 0 {
		return false
	}
	switch strings.ToUpper(x.Name) {
	case "NEXTVAL", "LASTVAL", "SETVAL":
		ref, ok := x.Args[0].(*ast.ColumnRef)
		if !ok || ref.Star || ref.Schema != "" || ref.Column == "" {
			return false
		}
		b.sequenceNeeds(&ast.TableRef{Loc: ref.Loc, Schema: ref.Table, Name: ref.Column}, ctes)
	}
	return false
}

// sequenceNeeds checks that a sequence function's relation exists (1146)
// and is a sequence (4089). A temporary table of the name shadows it, and
// a common table expression's name is left alone.
func (b *builder) sequenceNeeds(tr *ast.TableRef, ctes map[string]bool) {
	if tr == nil || tr.Schema == "" && ctes[strings.ToLower(tr.Name)] {
		return
	}
	r, st := b.ref(tr)
	if st != refOK {
		return
	}
	if f, temp, _ := b.resolve(r); temp || f == walk.Unknown {
		return
	}
	b.needKind(relKey(r.db, r.name), []string{relSequence}, r.rng, notFound("table", r.written, ""), r.written, "sequence")
}

// sourceNeeds checks that a relation a query reads exists (1146): a
// table, a view, or a sequence. DUAL and a common table expression's name
// are none, and a temporary table of the name may shadow it.
func (b *builder) sourceNeeds(tr *ast.TableRef, ctes map[string]bool) {
	if tr.Schema == "" && ctes[strings.ToLower(tr.Name)] || b.isDual([]ast.TableExpr{tr}) {
		return
	}
	r, st := b.ref(tr)
	if st != refOK {
		return
	}
	if f, temp, _ := b.resolve(r); temp || f == walk.Unknown {
		return
	}
	b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("table", r.written, ""), "")
}

// columnRefsOutsideQueries returns the column references of an expression
// outside subqueries.
func columnRefsOutsideQueries(expr ast.ExprNode) []*ast.ColumnRef {
	if expr == nil {
		return nil
	}
	var out []*ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt, *ast.SubqueryExpr, *ast.ExistsExpr:
			return false
		case *ast.ColumnRef:
			if !x.Star && x.Column != "" {
				out = append(out, x)
			}
			return false
		}
		return true
	})
	return out
}
