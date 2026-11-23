package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// viewBodyFails fails a CREATE VIEW whose query TiDB refuses while it plans
// it, before it looks at the view's own name (each verified on v8.5.5, in
// tidb/catalog/testdata/queries/simulate.tsv):
//
//   - a user or system variable (1351) or a parameter marker (1149);
//   - a builtin called with no argument where it takes one (1582);
//   - a window name no WINDOW clause of the query block defines (3579), or
//     one the clause defines twice (3591);
//   - a frame TiDB cannot run: GROUPS (1235), a start of UNBOUNDED
//     FOLLOWING (3584), an end before its start (3585), or an INTERVAL
//     bound on ROWS (3596);
//   - NTH_VALUE or NTILE reading a column where they take a constant (1210);
//   - two columns of one name, with no column list to rename them (1060);
//   - a sequence function on a relation that is no sequence (1347);
//   - a column the query block reads that its one table does not have, or
//     any column of a block with no table (1054).
func (b *builder) viewBodyFails(v *ast.CreateViewStmt) bool {
	sel := v.Select
	if sel == nil {
		return false
	}
	rng := b.anchor()
	if msg := viewTextFault(sel); msg != "" {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	if len(v.Columns) == 0 {
		if name := duplicateViewColumn(sel); name != "" {
			b.fail(rng, nil, review.DuplicateName, statesTwice("names", "view column", ident(name)))
			return true
		}
	}
	if b.viewSequenceFails(sel, rng) {
		return true
	}
	return b.viewColumnsFail(sel, rng)
}

// viewTextFault returns the text of a fault the view's query shows by
// itself, "" for none.
func viewTextFault(sel *ast.SelectStmt) string {
	// TiDB reads a variable or a parameter marker before it plans anything
	// (1149 for ROWS ? PRECEDING, not a frame error).
	msg := ""
	ast.Inspect(sel, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.VariableRef:
			msg = "A view's query can't read a variable."
		case *ast.ParamMarkerExpr:
			msg = "A view's query can't hold a parameter marker " + review.Code("?") + "."
		}
		return msg == ""
	})
	if msg != "" {
		return msg
	}
	ast.Inspect(sel, func(n ast.Node) bool {
		if msg != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectStmt:
			msg = windowNamesFault(x)
			if msg == "" {
				msg = rowSelectItemFault(x)
			}
		case *ast.FuncCallExpr:
			msg = funcCallFault(x)
		case *ast.WindowFrame:
			msg = windowFrameFault(x)
		case *ast.MatchExpr:
			// No table has a FULLTEXT index to search: v8.5.5 refuses to
			// create one (8200), and MATCH ... AGAINST is 1105.
			msg = engineName + " doesn't support " + review.Code("MATCH ... AGAINST") + "."
		case *ast.CastExpr:
			if x.Array {
				msg = engineName + " doesn't support " + review.Code("CAST(... AS ... ARRAY)") + " outside an index."
			}
		case *ast.TemporalLit:
			if x.Expr == nil && !strings.ContainsAny(x.Value, "0123456789") {
				msg = review.Code(x.Type+" '"+x.Value+"'") + " isn't a valid " + strings.ToLower(x.Type) + " literal."
			}
		}
		return msg == ""
	})
	return msg
}

// requiresArguments lists builtins that take at least one argument, each
// of which TiDB refuses with 1582 when called with none (verified on
// v8.5.5 for every name here).
var requiresArguments = map[string]bool{
	"ascii": true, "charset": true, "coalesce": true, "collation": true, "date": true, "day": true,
	"format": true, "hour": true, "if": true, "insert": true, "interval": true, "left": true,
	"log": true, "microsecond": true, "minute": true, "month": true, "password": true,
	"quarter": true, "repeat": true, "replace": true, "reverse": true, "right": true,
	"second": true, "time": true, "timestamp": true, "truncate": true, "week": true, "year": true,
}

func funcCallFault(f *ast.FuncCallExpr) string {
	name := strings.ToLower(f.Name)
	// A qualified name is a stored function's, and TiDB has none (1305).
	if f.Schema != "" {
		return "Function " + review.Code(ident(f.Schema)+"."+ident(strings.ToLower(f.Name))) + " doesn't exist; " + engineName + " has no stored functions."
	}
	if f.Schema == "" && f.HasParens && !f.Star && len(f.Args) == 0 && requiresArguments[name] {
		return review.Code(strings.ToUpper(f.Name)+"()") + " needs at least one argument."
	}
	// NTH_VALUE's N and NTILE's N are constants (1210 for a column).
	arg := -1
	switch name {
	case "nth_value":
		arg = 1
	case "ntile":
		arg = 0
	}
	if arg >= 0 && f.Schema == "" && len(f.Args) > arg {
		if _, ok := unwrapParens(f.Args[arg]).(*ast.ColumnRef); ok {
			return review.Code(strings.ToUpper(f.Name)) + " takes a constant, not a column, as its count."
		}
	}
	return ""
}

// windowNamesFault checks one query block's window names: every name an
// OVER or a WINDOW definition refers to has to be defined in the block's
// own WINDOW clause, and that clause defines a name once. A window name is
// compared case-sensitively for a duplicate and without case for a
// reference, so neither check claims more than TiDB refuses.
func windowNamesFault(sel *ast.SelectStmt) string {
	defined := make(map[string]bool, len(sel.WindowClause))
	folded := make(map[string]bool, len(sel.WindowClause))
	for _, w := range sel.WindowClause {
		if w == nil || w.Name == "" {
			continue
		}
		if defined[w.Name] {
			return "Window " + review.Code(ident(w.Name)) + " is defined twice."
		}
		defined[w.Name] = true
		folded[strings.ToLower(w.Name)] = true
	}
	undefined := func(name string) string {
		if name == "" || folded[strings.ToLower(name)] {
			return ""
		}
		return "Window " + review.Code(ident(name)) + " isn't defined in the query's " + review.Code("WINDOW") + " clause."
	}
	for _, w := range sel.WindowClause {
		if w != nil {
			if msg := undefined(w.RefName); msg != "" {
				return msg
			}
		}
	}
	msg := ""
	inBlock(sel, func(n ast.Node) {
		if f, ok := n.(*ast.FuncCallExpr); ok && f.Over != nil && msg == "" {
			msg = undefined(f.Over.RefName)
		}
	})
	return msg
}

// inBlock calls fn on every node of a query block, its WINDOW clause
// aside, without entering the query blocks it holds.
func inBlock(sel *ast.SelectStmt, fn func(ast.Node)) {
	visit := func(root ast.Node) {
		if root == nil {
			return
		}
		ast.Inspect(root, func(n ast.Node) bool {
			if _, ok := n.(*ast.SelectStmt); ok {
				return false
			}
			fn(n)
			return true
		})
	}
	for _, e := range sel.TargetList {
		visit(e)
	}
	for _, te := range sel.From {
		visit(te)
	}
	visit(sel.Where)
	for _, g := range sel.GroupBy {
		visit(g)
	}
	visit(sel.Having)
	for _, o := range sel.OrderBy {
		visit(o)
	}
}

func windowFrameFault(f *ast.WindowFrame) string {
	if f.Type == ast.FrameGroups {
		return engineName + " doesn't support a " + review.Code("GROUPS") + " window frame."
	}
	if f.Start == nil {
		return ""
	}
	if f.Start.Type == ast.BoundUnboundedFollowing {
		return "A window frame can't start at " + review.Code("UNBOUNDED FOLLOWING") + "."
	}
	if f.End != nil {
		if f.End.Type == ast.BoundUnboundedPreceding {
			return "A window frame can't end at " + review.Code("UNBOUNDED PRECEDING") + "."
		}
		if frameRank(f.End.Type) < frameRank(f.Start.Type) {
			return "A window frame can't end before it starts."
		}
	}
	if f.Type == ast.FrameRows {
		for _, bound := range []*ast.WindowFrameBound{f.Start, f.End} {
			if bound == nil {
				continue
			}
			if _, ok := bound.Offset.(*ast.IntervalExpr); ok {
				return "A " + review.Code("ROWS") + " window frame can't take an " + review.Code("INTERVAL") + " bound."
			}
		}
	}
	return ""
}

// frameRank orders frame bounds from the first row to the last; a
// PRECEDING or FOLLOWING bound with an offset ranks by its side alone,
// which leaves the order of two bounds of one side unjudged.
func frameRank(t ast.WindowBoundType) int {
	switch t {
	case ast.BoundUnboundedPreceding:
		return 0
	case ast.BoundPreceding:
		return 1
	case ast.BoundCurrentRow:
		return 2
	case ast.BoundFollowing:
		return 3
	}
	return 4
}

// duplicateViewColumn returns a name two columns of the view's first query
// block would take, "" for none. Only a column read by name or an alias
// names one: any other expression's name is its text, which the walk does
// not derive, and a star's columns are the table's.
func duplicateViewColumn(sel *ast.SelectStmt) string {
	for sel.ParenSource != nil || sel.Left != nil {
		if sel.ParenSource != nil {
			sel = sel.ParenSource
		} else {
			sel = sel.Left
		}
	}
	if sel.TableSource != nil || sel.ValuesSource != nil {
		return ""
	}
	seen := make(map[string]bool, len(sel.TargetList))
	for _, field := range sel.TargetList {
		name := ""
		expr := field
		if rt, ok := field.(*ast.ResTarget); ok {
			name, expr = rt.Name, rt.Val
		}
		if c, ok := expr.(*ast.ColumnRef); ok && c.Star {
			return ""
		}
		if _, ok := expr.(*ast.StarExpr); ok {
			return ""
		}
		if name == "" {
			c, ok := expr.(*ast.ColumnRef)
			if !ok {
				continue
			}
			name = c.Column
		}
		k := strings.ToLower(name)
		if outsideASCII(k) {
			continue
		}
		if seen[k] {
			return name
		}
		seen[k] = true
	}
	return ""
}

// sequenceFuncs are the functions that read a sequence by name, the
// sequence written as a column reference: NEXT VALUE FOR reads as NEXTVAL.
var sequenceFuncs = map[string]bool{"nextval": true, "lastval": true, "setval": true}

// viewSequenceFails fails a view whose query calls a sequence function on
// a relation that exists and is no sequence (1347).
func (b *builder) viewSequenceFails(sel *ast.SelectStmt, rng review.Range) bool {
	failed := false
	ast.Inspect(sel, func(n ast.Node) bool {
		if failed {
			return false
		}
		f, ok := n.(*ast.FuncCallExpr)
		if !ok || f.Schema != "" || !sequenceFuncs[strings.ToLower(f.Name)] || len(f.Args) == 0 {
			return true
		}
		c, ok := f.Args[0].(*ast.ColumnRef)
		if !ok || c.Star || c.Schema != "" || c.Column == "" {
			return true
		}
		r, st := b.ref(&ast.TableRef{Schema: c.Table, Name: c.Column, Loc: c.Loc})
		if st != refOK {
			failed = st == refFailed
			return !failed
		}
		k := relKey(r.db, r.name)
		if tf, _ := b.l.Lookup(tempKey(r.db, r.name)); tf != walk.Absent {
			return true
		}
		if fk, kind := b.l.Lookup(k); fk == walk.Present && kind != "" && kind != relSequence {
			b.needKind(k, []string{relSequence}, rng, notFound("sequence", r.written, ""), r.written, "sequence")
			failed = true
		}
		return !failed
	})
	return failed
}

// niladicNames are names TiDB reads as a function call when written bare;
// the parser reads most as calls already, and the walk never takes one of
// these for a column.
var niladicNames = map[string]bool{
	"current_date": true, "current_time": true, "current_timestamp": true, "current_user": true,
	"current_role": true, "localtime": true, "localtimestamp": true, "utc_date": true,
	"utc_time": true, "utc_timestamp": true,
}

// viewColumnsFail checks the columns each top-level query block of the
// view reads in its select list and WHERE clause, which no alias of the
// block can name (1054). A block with no table reads no column at all; a
// block reading one table, by its name or alias, reads that table's. A
// block that reads more than one table, a derived table or a CTE, and every
// nested query, are left alone, as are names qualified by a database.
func (b *builder) viewColumnsFail(sel *ast.SelectStmt, rng review.Range) bool {
	if len(sel.CTEs) > 0 {
		return false
	}
	var blocks []*ast.SelectStmt
	var collect func(s *ast.SelectStmt)
	collect = func(s *ast.SelectStmt) {
		switch {
		case s == nil:
		case len(s.CTEs) > 0:
			blocks = append(blocks, nil)
		case s.ParenSource != nil:
			collect(s.ParenSource)
		case s.Left != nil || s.Right != nil:
			collect(s.Left)
			collect(s.Right)
		default:
			blocks = append(blocks, s)
		}
	}
	collect(sel)
	for _, s := range blocks {
		if s == nil {
			return false
		}
	}
	for _, s := range blocks {
		if s.TableSource != nil || s.ValuesSource != nil {
			continue
		}
		if b.blockColumnsFail(s, rng) {
			return true
		}
	}
	return false
}

func (b *builder) blockColumnsFail(s *ast.SelectStmt, rng review.Range) bool {
	var table *tableRef
	qualifier := ""
	switch {
	case viewBlockReadsOnlyNoTable(s, b.w.sql):
	case len(s.From) == 1:
		tr, ok := s.From[0].(*ast.TableRef)
		if !ok {
			return false
		}
		r, st := b.ref(tr)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			return false
		}
		f, temp, kind := b.resolve(r)
		if temp || f != walk.Present || kind != relTable {
			return false
		}
		table = &r
		qualifier = tr.Name
		if tr.Alias != "" {
			qualifier = tr.Alias
		}
	default:
		return false
	}
	var refs []*ast.ColumnRef
	collect := func(root ast.Node) {
		if root == nil {
			return
		}
		ast.Inspect(root, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectStmt, *ast.WindowDef, *ast.OrderByItem:
				return false
			case *ast.FuncCallExpr:
				if sequenceFuncs[strings.ToLower(x.Name)] {
					return false
				}
			case *ast.ColumnRef:
				refs = append(refs, x)
			}
			return true
		})
	}
	for _, e := range s.TargetList {
		if rt, ok := e.(*ast.ResTarget); ok {
			collect(rt.Val)
			continue
		}
		collect(e)
	}
	collect(s.Where)
	for _, c := range refs {
		if c.Star || c.Schema != "" || c.Column == "" || niladicNames[strings.ToLower(c.Column)] {
			continue
		}
		if table == nil {
			if c.Table != "" {
				continue
			}
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(c.Column))+" doesn't exist in a query that reads no table.")
			return true
		}
		if c.Table != "" && !strings.EqualFold(c.Table, qualifier) {
			continue
		}
		k, ok := b.key(spaceColumn, c.Column)
		if !ok {
			continue
		}
		ck := colKey(table.db, table.name, k)
		if f, _ := b.l.Lookup(ck); f == walk.Absent {
			b.need(ck, walk.Present, nil, rng, notFound("column", ident(c.Column), table.written), "") // 1054
			return true
		}
	}
	return false
}

// rowSelectItemFault returns the fault of a select item that is a row of
// more than one column, which a query's result can't hold (1241).
func rowSelectItemFault(sel *ast.SelectStmt) string {
	for _, field := range sel.TargetList {
		expr := field
		if rt, ok := field.(*ast.ResTarget); ok {
			expr = rt.Val
		}
		if r, ok := expr.(*ast.RowExpr); ok && len(r.Items) > 1 {
			return "A select item can't be a row of " + strconv.Itoa(len(r.Items)) + " columns."
		}
	}
	return ""
}
