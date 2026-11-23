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
//     one that is no sequence (4089);
//   - what it holds as MariaDB parses it (viewTextFault): a variable
//     (1351), an INTO or PROCEDURE clause or a parameter marker (1064), and
//     VALUES rows of different lengths (4099) or of none (4141);
//   - a column of its top query block qualified by a table the block
//     doesn't read (qualifiersFail, 1054), and a column of the block's one
//     base table the table lacks in a function's arguments too (1054);
//   - a geometry constructor given a value that is no geometry (4079);
//   - windows windowsFault refuses (4010, 4009, 4014, 4022).
//
// named is set when the view lists its column names.
func (b *builder) viewQueryFails(sel *ast.SelectStmt, named bool) bool {
	if sel == nil {
		return false
	}
	if msg := viewTextFault(sel); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
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
		if b.qualifiersFail(top) {
			return true
		}
		b.singleTableColumns(top)
		if b.geometryArgsFail(top) {
			return true
		}
	}
	if b.blocksFail(sel, ctes) {
		return true
	}
	if msg := windowsFault(sel); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
	}
	return false
}

// viewTextFault returns the fault MariaDB meets in a view's query as it
// parses it, "" for none (each verified on mariadb:11.8.8): a variable,
// user or system (1351); an INTO or PROCEDURE clause, or a parameter
// marker (1064); and a VALUES query whose rows differ in length (4099) or
// one with no column (4141).
func viewTextFault(sel *ast.SelectStmt) string {
	msg := ""
	ast.Inspect(sel, func(n ast.Node) bool {
		if msg != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectStmt:
			switch {
			case x.Into != nil:
				msg = "A view's query can't hold an " + review.Code("INTO") + " clause."
			case x.Procedure != nil:
				msg = "A view's query can't hold a " + review.Code("PROCEDURE") + " clause."
			}
		case *ast.VariableRef:
			msg = "A view's query can't read a variable."
		case *ast.ParamMarkerExpr:
			msg = "A view's query can't hold a parameter marker " + review.Code("?") + "."
		case *ast.ValuesStmt:
			for i, row := range x.Rows {
				switch {
				case len(row) == 0:
					msg = "Each row of a " + review.Code("VALUES") + " query in a view needs at least one column."
				case i > 0 && len(row) != len(x.Rows[0]):
					msg = "Each row of a " + review.Code("VALUES") + " query needs the same number of columns."
				}
				if msg != "" {
					break
				}
			}
		}
		return msg == ""
	})
	return msg
}

// blockQualifiers returns the names a query block's FROM clause gives its
// tables, folded, and the keys of the databases each base table is in by
// name; false when it reads anything the walk does not name.
func (b *builder) blockQualifiers(from []ast.TableExpr) (map[string][]string, bool) {
	out := make(map[string][]string)
	var visit func(te ast.TableExpr) bool
	visit = func(te ast.TableExpr) bool {
		switch x := te.(type) {
		case *ast.TableRef:
			if x.Alias != "" {
				out[strings.ToLower(x.Alias)] = append(out[strings.ToLower(x.Alias)], "")
				return true
			}
			db := x.Schema
			if db == "" {
				if b.w.lost || b.w.noDatabase {
					return false
				}
				db = b.w.current
			} else {
				k, ok := b.key(spaceDatabase, db)
				if !ok {
					return false
				}
				db = k
			}
			out[strings.ToLower(x.Name)] = append(out[strings.ToLower(x.Name)], db)
			return true
		case *ast.JoinClause:
			return visit(x.Left) && visit(x.Right)
		case *ast.SubqueryExpr:
			if x.Alias == "" {
				return false
			}
			out[strings.ToLower(x.Alias)] = append(out[strings.ToLower(x.Alias)], "")
			return true
		case *ast.JsonTableExpr:
			if x.Alias == "" {
				return false
			}
			out[strings.ToLower(x.Alias)] = append(out[strings.ToLower(x.Alias)], "")
			return true
		}
		return false
	}
	for _, te := range from {
		if !visit(te) {
			return nil, false
		}
	}
	return out, true
}

// qualifiersFail fails a view's top query block whose select list or
// WHERE clause reads a column by a table name its FROM clause doesn't give
// (1054), or by database and table where the table isn't in that
// database (1054), outside subqueries (verified on mariadb:11.8.8). A
// qualifier that differs from a table's name only in case is taken for
// it. MariaDB matches db.t.c against an alias too, and a database the
// walk can't compare leaves the column alone.
func (b *builder) qualifiersFail(sel *ast.SelectStmt) bool {
	if b.isDual(sel.From) {
		return false
	}
	quals, ok := b.blockQualifiers(sel.From)
	if !ok {
		return false
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
		for _, ref := range columnRefsInCalls(expr) {
			if ref.Table == "" || outsideASCII(ref.Table) || outsideASCII(ref.Schema) {
				continue
			}
			dbs, named := quals[strings.ToLower(ref.Table)]
			miss := !named
			if named && ref.Schema != "" {
				db, ok := b.key(spaceDatabase, ref.Schema)
				if !ok {
					continue
				}
				miss = true
				for _, d := range dbs {
					// An alias takes any database; one that differs only
					// in case may be the table's.
					if d == "" || strings.EqualFold(d, db) {
						miss = false
					}
				}
			}
			if !miss {
				continue
			}
			name := ident(ref.Table) + "." + ident(ref.Column)
			if ref.Schema != "" {
				name = ident(ref.Schema) + "." + name
			}
			b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, "Column "+review.Code(name)+" doesn't exist in the query.") // 1054
			return true
		}
	}
	return false
}

// geometryConstructors are the functions that build a geometry of the
// geometries they take; any other value is 4079.
var geometryConstructors = map[string]bool{
	"linestring": true, "polygon": true, "multipoint": true, "multilinestring": true,
	"multipolygon": true, "geometrycollection": true,
}

// geometryTypes are the base types of MariaDB's geometry columns.
var geometryTypes = map[string]bool{
	"geometry": true, "point": true, "linestring": true, "polygon": true, "multipoint": true,
	"multilinestring": true, "multipolygon": true, "geometrycollection": true,
}

// geometryArgsFail fails a view's top query block whose select list or
// WHERE clause gives a geometry constructor a literal, or a column of the
// block's one base table whose type is no geometry (4079, verified on
// mariadb:11.8.8), outside subqueries.
func (b *builder) geometryArgsFail(sel *ast.SelectStmt) bool {
	var r *tableRef
	qualifier := ""
	if len(sel.From) == 1 {
		if tr, ok := sel.From[0].(*ast.TableRef); ok && !b.isDual(sel.From) && tr.SystemTime == nil {
			if t, st := b.ref(tr); st == refOK {
				if f, temp, kind := b.resolve(t); !temp && f == walk.Present && kind == relTable && b.w.versioned[relKey(t.db, t.name)] == walk.Absent {
					r, qualifier = &t, tr.Name
					if tr.Alias != "" {
						qualifier = tr.Alias
					}
				}
			}
		}
	}
	argType := func(e ast.ExprNode) string {
		switch x := e.(type) {
		case *ast.IntLit, *ast.FloatLit, *ast.StringLit:
			return "literal"
		case *ast.ColumnRef:
			if r == nil || x.Star || x.Schema != "" || x.Table != "" && !strings.EqualFold(x.Table, qualifier) {
				return ""
			}
			k, ok := b.key(spaceColumn, x.Column)
			if !ok {
				return ""
			}
			return b.ledgerType(*r, k)
		}
		return ""
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
		failed := false
		ast.Inspect(expr, func(n ast.Node) bool {
			if failed {
				return false
			}
			switch x := n.(type) {
			case *ast.SelectStmt, *ast.SubqueryExpr, *ast.ExistsExpr:
				return false
			case *ast.FuncCallExpr:
				if x.Schema != "" || x.NameQuoted || !geometryConstructors[strings.ToLower(x.Name)] {
					return true
				}
				for _, a := range x.Args {
					t := argType(a)
					if t == "literal" || t != "" && fixedTypes[t] && !geometryTypes[t] {
						b.fail(rangeOf(x.Loc), nil, review.UnsupportedForm, review.Code(strings.ToUpper(x.Name))+" takes geometry values, and is given one that is no geometry.") // 4079
						failed = true
						return false
					}
				}
			}
			return true
		})
		if failed {
			return true
		}
	}
	return false
}

// windowsFault returns the fault of a view's query's windows MariaDB
// refuses, outside common table expressions, "" for none (each verified on
// mariadb:11.8.8): a WINDOW clause that defines a name twice (4010) or
// builds on a name it doesn't define (4009); a frame bound combination
// Window_frame::check_frame_bounds refuses (4014): a start at UNBOUNDED
// FOLLOWING, an end at UNBOUNDED PRECEDING, a start at CURRENT ROW and an
// end PRECEDING, or an end at CURRENT ROW and a start FOLLOWING; and a
// frame's EXCLUDE (4022).
func windowsFault(sel *ast.SelectStmt) string {
	msg := ""
	ast.Inspect(sel, func(n ast.Node) bool {
		if msg != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.CommonTableExpr:
			return false
		case *ast.SelectStmt:
			seen := make(map[string]bool, len(x.WindowClause))
			for _, w := range x.WindowClause {
				if w == nil || w.Name == "" {
					continue
				}
				k := strings.ToLower(w.Name)
				if seen[k] {
					msg = "Window " + review.Code(ident(w.Name)) + " is defined twice."
					return false
				}
				seen[k] = true
			}
			for _, w := range x.WindowClause {
				if w != nil && w.RefName != "" && !seen[strings.ToLower(w.RefName)] {
					msg = "Window " + review.Code(ident(w.RefName)) + " is not defined."
					return false
				}
			}
		case *ast.WindowFrame:
			msg = frameFault(x)
		}
		return msg == ""
	})
	return msg
}

// frameFault returns the fault of a window frame windowsFault lists.
func frameFault(f *ast.WindowFrame) string {
	if f.Exclude != "" {
		return engineName + " doesn't support " + review.Code("EXCLUDE") + " in a window frame."
	}
	if f.Start == nil || f.End == nil {
		return ""
	}
	s, e := f.Start.Type, f.End.Type
	if s == ast.BoundUnboundedFollowing || e == ast.BoundUnboundedPreceding ||
		s == ast.BoundCurrentRow && (e == ast.BoundPreceding || e == ast.BoundUnboundedPreceding) ||
		e == ast.BoundCurrentRow && (s == ast.BoundFollowing || s == ast.BoundUnboundedFollowing) {
		return "This window frame's bounds are a combination " + engineName + " doesn't take."
	}
	return ""
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
		for _, ref := range columnRefsInCalls(expr) {
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

// columnRefsInCalls returns the column references of an expression outside
// subqueries and window specifications, in the arguments of function calls
// too, but those of a sequence function, which names a sequence. A name
// MariaDB reads as a function when written bare is no column.
func columnRefsInCalls(expr ast.ExprNode) []*ast.ColumnRef {
	if expr == nil {
		return nil
	}
	var out []*ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt, *ast.SubqueryExpr, *ast.ExistsExpr, *ast.WindowDef,
			*ast.NextValueForExpr, *ast.PreviousValueForExpr:
			return false
		case *ast.FuncCallExpr:
			switch strings.ToUpper(x.Name) {
			case "NEXTVAL", "LASTVAL", "SETVAL":
				return false
			}
		case *ast.ColumnRef:
			if !x.Star && x.Column != "" && !niladicNames[strings.ToLower(x.Column)] {
				out = append(out, x)
			}
			return false
		}
		return true
	})
	return out
}

// niladicNames are names MariaDB reads as a function call when written
// bare; the walk never takes one for a column.
var niladicNames = map[string]bool{
	"current_date": true, "current_time": true, "current_timestamp": true, "current_user": true,
	"current_role": true, "localtime": true, "localtimestamp": true, "utc_date": true,
	"utc_time": true, "utc_timestamp": true, "sysdate": true,
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
