package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// viewQueryFails fails a CREATE VIEW or ALTER VIEW whose query MySQL
// refuses, in the order MySQL 8.0.46 meets each fault (each verified on the
// pinned server, most in mysql/catalog/testdata/queries/simulate.tsv):
//
//   - while it parses: see parsedFault;
//   - while it opens the tables: a table that doesn't exist (1146);
//   - while it resolves each query block: a locking clause naming a table
//     the block doesn't read (3568), and then item by item and WHERE last, a
//     column the block's one table doesn't have, or any column of a block
//     with no table (1054), and the faults of the items themselves (see
//     itemFault and selectItemFault);
//   - while it sets up the windows: a window name defined twice (3591), and
//     a frame MySQL can't run (see frameFault);
//   - last, two view columns of one name, with no column list to rename them
//     (1060).
//
// A fault the walk can't place in that order with certainty is still a
// fault MySQL refuses; only which one it names first may differ.
func (b *builder) viewQueryFails(sel *ast.SelectStmt, columns []string, rng review.Range) bool {
	if sel == nil {
		return false
	}
	if msg := parsedFault(sel); msg != "" {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	if b.queryTablesFail(sel, rng) {
		return true
	}
	ctes := cteNames(sel)
	for _, s := range queryLeaves(sel) {
		if b.blockFails(s, ctes, rng) {
			return true
		}
	}
	var msg string
	ast.Inspect(sel, func(n ast.Node) bool {
		if msg != "" {
			return false
		}
		switch x := n.(type) {
		case *ast.CommonTableExpr:
			// MySQL sets up the windows of a CTE only where the query
			// reads it.
			return false
		case *ast.SelectStmt:
			msg = windowClauseFault(x)
		case *ast.WindowFrame:
			msg = frameFault(x)
		}
		return msg == ""
	})
	if msg != "" {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	if len(columns) == 0 {
		if name := duplicateViewColumn(sel); name != "" {
			b.fail(rng, nil, review.DuplicateName, statesTwice("names", "column", ident(name)))
			return true
		}
	}
	return false
}

// parsedFault returns the fault MySQL meets in the query's text while it
// parses it, "" for none: an INTO clause (1350), a temporal literal with no
// digit (1525), a variable (1351) or a parameter marker (1064), and a name
// where the grammar takes a number or a variable, which MySQL reads as a
// variable no routine declares (1327).
func parsedFault(sel *ast.SelectStmt) string {
	msg := ""
	ast.Inspect(sel, func(n ast.Node) bool {
		if s, ok := n.(*ast.SelectStmt); ok && s.Into != nil {
			msg = "A view's query can't hold an " + review.Code("INTO") + " clause."
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
			if x.Limit != nil {
				for _, e := range []ast.ExprNode{x.Limit.Count, x.Limit.Offset} {
					if name := bareName(e); name != "" {
						msg = review.Code("LIMIT") + " takes a number, not " + review.Code(ident(name)) + "."
					}
				}
			}
		case *ast.FuncCallExpr:
			arg := -1
			switch strings.ToLower(x.Name) {
			case "ntile":
				arg = 0
			case "lead", "lag":
				arg = 1
			}
			if x.Schema == "" && arg >= 0 && len(x.Args) > arg && x.Over != nil {
				if name := bareName(x.Args[arg]); name != "" {
					msg = review.Code(strings.ToUpper(x.Name)) + " takes a number, not " + review.Code(ident(name)) + "."
				}
			}
		case *ast.TemporalLit:
			// MySQL reads the literal as it parses it (1525 before even a
			// table that doesn't exist); one with no digit is no value of
			// any temporal type.
			if !strings.ContainsAny(x.Value, "0123456789") {
				msg = review.Code(x.Type+" '"+x.Value+"'") + " isn't a valid " + strings.ToLower(x.Type) + " literal."
			}
		case *ast.VariableRef:
			msg = "A view's query can't read a variable."
		case *ast.ParamMarkerExpr:
			msg = "A view's query can't hold a parameter marker " + review.Code("?") + "."
		}
		return msg == ""
	})
	return msg
}

// bareName returns the name an expression reads unqualified, "" when it is
// anything else.
func bareName(e ast.ExprNode) string {
	if c, ok := e.(*ast.ColumnRef); ok && !c.Star && c.Table == "" && c.Schema == "" {
		return c.Column
	}
	return ""
}

// queryLeaves returns the query blocks a query expression combines, left
// to right: the operands of its set operations, inside any parentheses.
func queryLeaves(sel *ast.SelectStmt) []*ast.SelectStmt {
	switch {
	case sel == nil:
		return nil
	case sel.ParenSource != nil:
		return queryLeaves(sel.ParenSource)
	case sel.Left != nil || sel.Right != nil:
		return append(queryLeaves(sel.Left), queryLeaves(sel.Right)...)
	}
	return []*ast.SelectStmt{sel}
}

// queryTablesFail fails a query that reads a table the ledger shows doesn't
// exist (1146): MySQL opens every table a view's query reads before it
// resolves anything in it. A name some WITH clause of the query defines may
// be that CTE, and is left alone, as are the tables of a CTE's own query,
// which MySQL doesn't open when nothing reads it.
func (b *builder) queryTablesFail(sel *ast.SelectStmt, rng review.Range) bool {
	ctes := cteNames(sel)
	failed := false
	ast.Inspect(sel, func(n ast.Node) bool {
		if failed {
			return false
		}
		switch n.(type) {
		case *ast.ForUpdate:
			// A locking clause names the query's tables by alias.
			return false
		case *ast.CommonTableExpr:
			// MySQL opens a CTE's tables only where the query reads it.
			return false
		}
		tr, ok := n.(*ast.TableRef)
		if !ok {
			return true
		}
		if tr.Schema == "" && (ctes[strings.ToLower(tr.Name)] || b.isDual(tr)) {
			return true
		}
		r, st := b.ref(tr)
		switch st {
		case refFailed:
			failed = true
			return false
		case refUnknown:
			return true
		}
		if tf, _ := b.l.Lookup(tempKey(r.db, r.name)); tf != walk.Absent {
			return true
		}
		k := relKey(r.db, r.name)
		if f, _ := b.l.Lookup(k); f == walk.Absent {
			b.fail(rangeOf(tr.Loc), []walk.Key{k}, review.ObjectNotFound, notFound("table", r.written, ""))
			failed = true
		}
		return !failed
	})
	return failed
}

// cteNames returns the names, folded, every WITH clause of a query defines.
func cteNames(sel *ast.SelectStmt) map[string]bool {
	ctes := make(map[string]bool)
	ast.Inspect(sel, func(n ast.Node) bool {
		if c, ok := n.(*ast.CommonTableExpr); ok {
			ctes[strings.ToLower(c.Name)] = true
		}
		return true
	})
	return ctes
}

// isDual reports whether a table reference is the keyword DUAL, which reads
// no table; a quoted `dual` names one.
func (b *builder) isDual(tr *ast.TableRef) bool {
	if tr == nil || tr.Schema != "" || !strings.EqualFold(tr.Name, "dual") {
		return false
	}
	start := tr.Loc.Start
	return start >= 0 && start < len(b.w.sql) && b.w.sql[start] != '`' && b.w.sql[start] != '"'
}

// blockScope is what a query block reads columns from: no table at all, or
// one table whose columns the ledger knows. A block the walk can't scope
// leaves its columns alone.
type blockScope struct {
	known     bool
	table     *tableRef
	qualifier string
}

func (b *builder) scopeOf(s *ast.SelectStmt, ctes map[string]bool) blockScope {
	if s.TableSource != nil {
		return blockScope{}
	}
	if s.ValuesSource != nil || len(s.From) == 0 {
		return blockScope{known: true}
	}
	if len(s.From) != 1 {
		return blockScope{}
	}
	tr, ok := s.From[0].(*ast.TableRef)
	if !ok {
		return blockScope{}
	}
	if b.isDual(tr) {
		return blockScope{known: true}
	}
	if tr.Schema == "" && ctes[strings.ToLower(tr.Name)] {
		// A CTE of the name shadows the table.
		return blockScope{}
	}
	r, st := b.ref(tr)
	if st != refOK {
		return blockScope{}
	}
	f, temp, kind := b.resolve(r)
	if temp || f != walk.Present || kind != relTable {
		return blockScope{}
	}
	sc := blockScope{known: true, table: &r, qualifier: tr.Name}
	if tr.Alias != "" {
		sc.qualifier = tr.Alias
	}
	return sc
}

// blockFails checks one query block: each select item in order, its
// columns and then its own faults, and then WHERE.
func (b *builder) blockFails(s *ast.SelectStmt, ctes map[string]bool, rng review.Range) bool {
	sc := b.scopeOf(s, ctes)
	if s.ValuesSource != nil {
		if msg := valuesFault(s.ValuesSource); msg != "" {
			b.fail(rng, nil, review.UnsupportedForm, msg)
			return true
		}
		for _, row := range s.ValuesSource.Rows {
			for _, e := range row {
				if b.columnsFail(sc, e, rng) {
					return true
				}
			}
		}
		return false
	}
	if msg := lockingFault(sc, s); msg != "" {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	windows := definedWindows(s)
	for _, item := range s.TargetList {
		expr := item
		if rt, ok := item.(*ast.ResTarget); ok {
			expr = rt.Val
		}
		if b.columnsFail(sc, expr, rng) {
			return true
		}
		if msg := b.itemFault(sc, s, expr, windows); msg != "" {
			b.fail(rng, nil, review.UnsupportedForm, msg)
			return true
		}
		if msg := b.selectItemFault(sc, expr); msg != "" {
			b.fail(rng, nil, review.UnsupportedForm, msg)
			return true
		}
	}
	if s.Where != nil {
		if b.columnsFail(sc, s.Where, rng) {
			return true
		}
		if msg := b.itemFault(sc, s, s.Where, windows); msg != "" {
			b.fail(rng, nil, review.UnsupportedForm, msg)
			return true
		}
	}
	return false
}

// lockingFault returns the fault of a scoped block's locking clause naming
// a table the block doesn't read (3568), which MySQL meets before the
// block's columns.
func lockingFault(sc blockScope, s *ast.SelectStmt) string {
	if !sc.known {
		return ""
	}
	locks := s.MoreLocking
	if s.ForUpdate != nil {
		locks = append([]*ast.ForUpdate{s.ForUpdate}, locks...)
	}
	for _, l := range locks {
		for _, t := range l.Tables {
			if t == nil || t.Schema != "" {
				continue
			}
			if sc.table == nil || !strings.EqualFold(t.Name, sc.qualifier) {
				return "The locking clause names " + review.Code(ident(t.Name)) + ", which its query block doesn't read."
			}
		}
	}
	return ""
}

// niladicNames are names MySQL reads as a function call when written bare;
// the walk never takes one for a column.
var niladicNames = map[string]bool{
	"current_date": true, "current_time": true, "current_timestamp": true, "current_user": true,
	"current_role": true, "localtime": true, "localtimestamp": true, "utc_date": true,
	"utc_time": true, "utc_timestamp": true,
}

// columnsFail fails an expression of a scoped block that reads a column the
// block can't (1054): any column of a block with no table, and of a block
// reading one table a column that table doesn't have, or one qualified by a
// name that isn't the table's. Columns in nested queries, window
// specifications and ORDER BY items are left alone.
func (b *builder) columnsFail(sc blockScope, expr ast.Node, rng review.Range) bool {
	if !sc.known || expr == nil {
		return false
	}
	var refs []*ast.ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt, *ast.SubqueryExpr, *ast.ExistsExpr, *ast.WindowDef, *ast.OrderByItem:
			return false
		case *ast.ColumnRef:
			refs = append(refs, x)
		}
		return true
	})
	for _, c := range refs {
		if c.Star || c.Column == "" || niladicNames[strings.ToLower(c.Column)] || outsideASCII(c.Column) {
			continue
		}
		if sc.table == nil {
			if c.Table != "" || c.Schema != "" {
				continue
			}
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(c.Column))+" doesn't exist in the query.")
			return true
		}
		t := sc.table
		switch b.qualifies(sc, c) {
		case qualUnsure:
			continue
		case qualMiss:
			name := ident(c.Table) + "." + ident(c.Column)
			if c.Schema != "" {
				name = ident(c.Schema) + "." + name
			}
			b.fail(rng, nil, review.ObjectNotFound, "Column "+review.Code(name)+" doesn't exist in the query.")
			return true
		}
		k, ok := b.key(spaceColumn, c.Column)
		if !ok {
			continue
		}
		ck := colKey(t.db, t.name, k)
		if f, _ := b.l.Lookup(ck); f == walk.Absent {
			b.fail(rng, []walk.Key{ck}, review.ObjectNotFound, notFound("column", ident(c.Column), t.written))
			return true
		}
	}
	return false
}

type qualification int

const (
	qualMatch qualification = iota
	qualMiss
	qualUnsure
)

// qualifies tells whether a column reference names the block's one table:
// unqualified, or by its alias or its name, and, when written, its
// database. A qualifier that differs from the table's only in case is taken
// for a match, which leaves the column checked against the table; a
// database the walk can't compare leaves it unsure.
func (b *builder) qualifies(sc blockScope, c *ast.ColumnRef) qualification {
	if c.Table == "" {
		return qualMatch
	}
	if !strings.EqualFold(c.Table, sc.qualifier) {
		return qualMiss
	}
	if c.Schema == "" {
		return qualMatch
	}
	// MySQL matches the table qualifier against the alias, and the
	// database against the table's, so db.alias.c reads the table too.
	db, ok := b.key(spaceDatabase, c.Schema)
	switch {
	case !ok:
		return qualUnsure
	case db == sc.table.db:
		return qualMatch
	case strings.EqualFold(db, sc.table.db):
		return qualUnsure
	}
	return qualMiss
}

// colType returns the base type of a column of the block's table that an
// expression reads by name, "" when it is anything else or not known.
func (b *builder) colType(sc blockScope, e ast.ExprNode) string {
	c, ok := unparen(e).(*ast.ColumnRef)
	if !ok || sc.table == nil || c.Star || c.Schema != "" || (c.Table != "" && !strings.EqualFold(c.Table, sc.qualifier)) {
		return ""
	}
	k, ok := b.key(spaceColumn, c.Column)
	if !ok {
		return ""
	}
	ck := colKey(sc.table.db, sc.table.name, k)
	if f, _ := b.l.Lookup(ck); f != walk.Present {
		return ""
	}
	return b.w.colTypes[ck]
}

func unparen(e ast.ExprNode) ast.ExprNode {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.Expr
	}
}

// geometryConstructors are the functions that build a geometry of the
// geometries they take; any other value is 1367 (POINT takes numbers).
var geometryConstructors = map[string]bool{
	"linestring": true, "polygon": true, "multipoint": true, "multilinestring": true,
	"multipolygon": true, "geometrycollection": true, "geomcollection": true,
}

// isLiteral reports a literal value, which no geometry and no column is.
func isLiteral(e ast.ExprNode) bool {
	switch unparen(e).(type) {
	case *ast.IntLit, *ast.FloatLit, *ast.StringLit, *ast.NullLit, *ast.BoolLit:
		return true
	}
	return false
}

// isNumericLiteral reports a number written as a literal.
func isNumericLiteral(e ast.ExprNode) bool {
	switch unparen(e).(type) {
	case *ast.IntLit, *ast.FloatLit:
		return true
	}
	return false
}

// notJSONType reports a column type JSON_VALUE refuses as its document
// (3146): a number or a temporal value, neither a JSON value nor a string.
func notJSONType(t string) bool {
	switch t {
	case "tinyint", "smallint", "mediumint", "int", "bigint", "decimal", "float", "double",
		"date", "datetime", "timestamp", "time", "year":
		return true
	}
	return false
}

// itemFault returns the fault of an expression of a query block that MySQL
// meets while it resolves it, "" for none, nested queries included:
//
//   - a window name the block's WINDOW clause doesn't define (3579);
//   - MATCH ... AGAINST or LIKE ... ESCAPE reading a column where it takes a
//     constant (1210);
//   - GROUPING() in a block with no GROUP BY ... WITH ROLLUP (1111);
//   - a geometry constructor given a value that is no geometry (1367);
//   - JSON_VALUE reading a number or a temporal value (3146);
//   - CAST(... AT TIME ZONE ...) of a column that is no TIMESTAMP (3998), or
//     to a zone other than 'UTC' or '+00:00' (1298).
func (b *builder) itemFault(sc blockScope, s *ast.SelectStmt, expr ast.ExprNode, windows map[string]bool) string {
	msg := ""
	block := s
	var visit func(n ast.Node, block *ast.SelectStmt, windows map[string]bool, top bool)
	visit = func(root ast.Node, block *ast.SelectStmt, windows map[string]bool, top bool) {
		ast.Inspect(root, func(n ast.Node) bool {
			if msg != "" {
				return false
			}
			switch x := n.(type) {
			case *ast.SelectStmt:
				if x == block {
					return true
				}
				for _, leaf := range queryLeaves(x) {
					w := definedWindows(leaf)
					for _, item := range leaf.TargetList {
						visit(item, leaf, w, false)
					}
					for _, te := range leaf.From {
						visit(te, leaf, w, false)
					}
					if leaf.Where != nil {
						visit(leaf.Where, leaf, w, false)
					}
				}
				return false
			case *ast.FuncCallExpr:
				msg = b.callFault(sc, x, block, windows, top)
			case *ast.MatchExpr:
				if _, ok := unparen(x.Against).(*ast.ColumnRef); ok {
					msg = review.Code("AGAINST") + " takes a constant, not a column."
				}
			case *ast.LikeExpr:
				if _, ok := unparen(x.Escape).(*ast.ColumnRef); ok {
					msg = review.Code("ESCAPE") + " takes a constant, not a column."
				}
			case *ast.CastExpr:
				msg = b.castFault(sc, x, top)
			}
			return msg == ""
		})
	}
	visit(expr, block, windows, true)
	return msg
}

func (b *builder) callFault(sc blockScope, f *ast.FuncCallExpr, block *ast.SelectStmt, windows map[string]bool, top bool) string {
	if f.Over != nil && f.Over.RefName != "" && !windows[strings.ToLower(f.Over.RefName)] {
		return "Window " + review.Code(ident(f.Over.RefName)) + " isn't defined in the query's " + review.Code("WINDOW") + " clause."
	}
	if f.Schema != "" {
		return ""
	}
	name := strings.ToLower(f.Name)
	switch {
	case name == "grouping" && block != nil && !block.WithRollup:
		return review.Code("GROUPING()") + " needs a " + review.Code("GROUP BY ... WITH ROLLUP") + " in its query block."
	case geometryConstructors[name]:
		for _, a := range f.Args {
			if isLiteral(a) {
				return review.Code(strings.ToUpper(f.Name)) + " takes geometry values, not a literal."
			}
			if !top {
				continue
			}
			if t := b.colType(sc, a); t != "" && !isGeometryBaseType(t) {
				return review.Code(strings.ToUpper(f.Name)) + " takes geometry values, and its argument is " + review.Code(t) + "."
			}
		}
	case name == "json_value" && f.JSONValueArm && len(f.Args) > 0:
		if isNumericLiteral(f.Args[0]) {
			return review.Code("JSON_VALUE") + " reads a JSON document or a string, not a number."
		}
		if top {
			if t := b.colType(sc, f.Args[0]); notJSONType(t) {
				return review.Code("JSON_VALUE") + " reads a JSON document or a string, and its argument is " + review.Code(t) + "."
			}
		}
	}
	return ""
}

func (b *builder) castFault(sc blockScope, c *ast.CastExpr, top bool) string {
	if c.AtTimeZone == "" && !c.AtTimeZoneInterval {
		return ""
	}
	if top {
		if t := b.colType(sc, c.Expr); t != "" && t != "timestamp" {
			return review.Code("CAST(... AT TIME ZONE ...)") + " takes a " + review.Code("TIMESTAMP") + ", and its argument is " + review.Code(t) + "."
		}
	}
	zone := c.AtTimeZone
	bad := false
	if c.AtTimeZoneInterval {
		bad = zone != "+00:00" && !strings.EqualFold(zone, "UTC")
	} else {
		bad = zone != "UTC" && zone != "+00:00" && !strings.ContainsAny(zone, "0123456789")
	}
	if bad {
		return review.Code("CAST(... AT TIME ZONE ...)") + " takes only " + review.Code("'UTC'") + " or " + review.Code("INTERVAL '+00:00'") + ", not " + review.Code("'"+zone+"'") + "."
	}
	return ""
}

// selectItemFault returns the fault of one select item of a top-level query
// block: a row of more than one column (1241), or DEFAULT(c) as a column of
// the view, aliased or not, which MySQL refuses for every user (1143, the
// privilege it checks for the column, even for root; 1364 first when c has
// no default).
func (b *builder) selectItemFault(sc blockScope, expr ast.ExprNode) string {
	switch x := unparen(expr).(type) {
	case *ast.RowExpr:
		if len(x.Items) > 1 {
			return "A select item can't be a row of " + strconv.Itoa(len(x.Items)) + " columns."
		}
	case *ast.DefaultExpr:
		if x.Column == "" || sc.table == nil {
			return ""
		}
		k, ok := b.key(spaceColumn, x.Column)
		if !ok {
			return ""
		}
		if f, _ := b.l.Lookup(colKey(sc.table.db, sc.table.name, k)); f == walk.Present {
			return "A view can't take " + review.Code("DEFAULT("+ident(x.Column)+")") + " as a column."
		}
	}
	return ""
}

// valuesFault returns the fault of a VALUES query outside an INSERT: rows
// of different lengths (1136), an empty row (3942), or DEFAULT (3943).
func valuesFault(v *ast.ValuesStmt) string {
	if len(v.Rows) == 0 {
		return ""
	}
	n := len(v.Rows[0])
	for _, row := range v.Rows[1:] {
		if len(row) != n {
			return "Each row of a " + review.Code("VALUES") + " query needs the same number of columns."
		}
	}
	if n == 0 {
		return "Each row of a " + review.Code("VALUES") + " query outside an " + review.Code("INSERT") + " needs at least one column."
	}
	for _, row := range v.Rows {
		for _, e := range row {
			if d, ok := e.(*ast.DefaultExpr); ok && d.Column == "" {
				return "A " + review.Code("VALUES") + " query outside an " + review.Code("INSERT") + " can't take " + review.Code("DEFAULT") + "."
			}
		}
	}
	return ""
}

// definedWindows returns the window names a query block's WINDOW clause
// defines, folded: MySQL compares window names without case.
func definedWindows(s *ast.SelectStmt) map[string]bool {
	out := make(map[string]bool, len(s.WindowClause))
	for _, w := range s.WindowClause {
		if w != nil && w.Name != "" {
			out[strings.ToLower(w.Name)] = true
		}
	}
	return out
}

// windowClauseFault checks a query block's WINDOW clause: a name defined
// twice, without case (3591), or a definition built on a name the clause
// doesn't define (3579).
func windowClauseFault(s *ast.SelectStmt) string {
	seen := make(map[string]bool, len(s.WindowClause))
	for _, w := range s.WindowClause {
		if w == nil || w.Name == "" {
			continue
		}
		k := strings.ToLower(w.Name)
		if seen[k] {
			return "Window " + review.Code(ident(w.Name)) + " is defined twice."
		}
		seen[k] = true
	}
	for _, w := range s.WindowClause {
		if w != nil && w.RefName != "" && !seen[strings.ToLower(w.RefName)] {
			return "Window " + review.Code(ident(w.RefName)) + " isn't defined in the query's " + review.Code("WINDOW") + " clause."
		}
	}
	return ""
}

// frameFault returns the fault of a window frame MySQL 8.0 can't run: a
// GROUPS frame or an EXCLUDE clause (1235), an INTERVAL bound on ROWS
// (3596), a start at UNBOUNDED FOLLOWING (3584), an end at UNBOUNDED
// PRECEDING (3585), or an end of a kind before its start's (3586), checked
// in that order.
func frameFault(f *ast.WindowFrame) string {
	if f.Type == ast.FrameGroups {
		return engineName + " doesn't support a " + review.Code("GROUPS") + " window frame."
	}
	if f.Exclude != "" {
		return engineName + " doesn't support " + review.Code("EXCLUDE") + " in a window frame."
	}
	if f.Start == nil {
		return ""
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
	if f.Start.Type == ast.BoundUnboundedFollowing {
		return "A window frame can't start at " + review.Code("UNBOUNDED FOLLOWING") + "."
	}
	if f.End == nil {
		return ""
	}
	if f.End.Type == ast.BoundUnboundedPreceding {
		return "A window frame can't end at " + review.Code("UNBOUNDED PRECEDING") + "."
	}
	if frameRank(f.End.Type) < frameRank(f.Start.Type) {
		return "A window frame can't end before it starts."
	}
	return ""
}

// frameRank orders frame bounds by their kind, from the first row to the
// last; two bounds of one side rank alike, whatever their offsets.
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

// duplicateViewColumn returns a name two columns of the view would take,
// from its first query block, "" for none. Only a column read by name or an
// alias names one: any other expression's name is its text, which the walk
// doesn't derive, and a star's columns are the table's.
func duplicateViewColumn(sel *ast.SelectStmt) string {
	leaves := queryLeaves(sel)
	if len(leaves) == 0 {
		return ""
	}
	first := leaves[0]
	if first.TableSource != nil || first.ValuesSource != nil {
		return ""
	}
	seen := make(map[string]bool, len(first.TargetList))
	for _, item := range first.TargetList {
		name, expr := "", item
		if rt, ok := item.(*ast.ResTarget); ok {
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
