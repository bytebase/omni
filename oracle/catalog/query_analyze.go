package catalog

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

// This file derives what Oracle stores for a view's query: the columns of
// its result (name, data type, nullability) and its text, in which Oracle
// expands each * of the outermost query blocks into the column list it
// stands for. It refuses a query with the error Oracle refuses it with.
// testdata/queries/simulate.tsv holds both to the pinned Oracle.

// A qcol is one column of a query's result or of a FROM item.
type qcol struct {
	name    string
	typ     qtype
	notNull bool
	hidden  bool
	// expr marks a result column computed by an expression rather than read
	// from a column, which a view must alias (ORA-00998).
	expr bool
}

// A qrel is a FROM item a column can be qualified by.
type qrel struct {
	name   string // exposed name; empty for an unaliased inline view
	schema string // owner, when the item is a table exposed by its name
	cols   []qcol
	// merged holds the columns a USING or NATURAL join merged: they can no
	// longer be qualified (ORA-25154).
	merged map[string]bool
	// nullExt reports the item is null-extended by an outer join.
	nullExt bool
	// partitioned reports a partitioned table.
	partitioned bool
	// view reports the item is a stored view (no partitions, no flashback).
	view bool
}

func (r *qrel) col(i int) qcol {
	c := r.cols[i]
	if r.nullExt {
		c.notNull = false
	}
	return c
}

// A fromCol is one column of a FROM clause in the order * lists it.
type fromCol struct {
	rel *qrel // nil for a merged USING or NATURAL column
	idx int
	col qcol // the merged column, when rel is nil
}

// A qfrom is a FROM clause, or one join in it: the items it exposes and
// the columns * lists.
type qfrom struct {
	rels  []*qrel
	items []fromCol
	// flat are the items Oracle keeps in the FROM list as they are, whose
	// exposed names * must tell apart (ORA-00918). Oracle turns a nested
	// join, and the null-extended side of an outer join, into a view of
	// its own, whose columns * reads without their item's name.
	flat []*qrel
}

func (f *qfrom) itemCol(it fromCol) qcol {
	if it.rel == nil {
		return it.col
	}
	return it.rel.col(it.idx)
}

// A qscope is a query block's name space: its FROM clause and the block it
// is nested in.
type qscope struct {
	parent *qscope
	from   *qfrom
	ctes   *cteEnv
	// aliases are the block's select-list aliases, which ORDER BY (and in
	// 23ai GROUP BY) may name.
	aliases []qcol
	// hierarchical reports a CONNECT BY clause.
	hierarchical bool
	// windows are the block's WINDOW clause definitions, by name.
	windows map[string]*ast.WindowSpec
}

type cteDef struct {
	name      string
	node      *ast.CTE
	cols      []qcol
	analyzed  bool
	analyzing bool
}

type cteEnv struct {
	parent *cteEnv
	defs   map[string]*cteDef
}

func (c *cteEnv) lookup(name string) *cteDef {
	for env := c; env != nil; env = env.parent {
		if d := env.defs[name]; d != nil {
			return d
		}
	}
	return nil
}

// queryAnalyzer derives a view query's result in the extractor's schema.
type queryAnalyzer struct {
	e *oracleMetadataExtractor
	// edits are the text replacements * expansion makes, in schema text
	// offsets.
	edits []textEdit
	// plsqlFuncs are the functions a WITH clause declares, by name, with
	// their return types.
	plsqlFuncs map[string]qtype
}

type textEdit struct {
	start, end int
	text       string
}

// queryResult is a query's result columns.
type queryResult struct {
	cols []qcol
}

// analyzeViewQuery returns the columns of a view's query and the text Oracle
// stores for it.
func (e *oracleMetadataExtractor) analyzeViewQuery(query ast.StmtNode) ([]qcol, string, error) {
	sel, ok := query.(*ast.SelectStmt)
	if !ok || sel == nil {
		return nil, e.nodeText(query), nil
	}
	a := &queryAnalyzer{e: e, plsqlFuncs: map[string]qtype{}}
	res, err := a.query(sel, nil, nil, true)
	if err != nil {
		return nil, "", err
	}
	loc := ast.NodeLoc(query)
	text := e.nodeText(query)
	if loc.Start >= 0 && len(a.edits) > 0 {
		sort.Slice(a.edits, func(i, j int) bool { return a.edits[i].start > a.edits[j].start })
		for _, ed := range a.edits {
			s, t := ed.start-loc.Start, ed.end-loc.Start
			if s < 0 || t > len(text) || s > t {
				continue
			}
			text = text[:s] + ed.text + text[t:]
		}
	}
	return res.cols, text, nil
}

// setOperand is one branch of a compound query, in order.
type setOperand struct {
	sel *ast.SelectStmt
	op  ast.SetOperation // the operator before this branch
	all bool
	// tail reports the branch whose ORDER BY and FETCH belong to the
	// compound: the parser attaches them to the last branch.
	tail bool
}

// flattenSetOp lists a compound query's branches. The parser nests a
// compound to the right, leaving the first branch's query block in the
// root, and puts a trailing ORDER BY or FETCH on the last branch, or on a
// wrapper around a parenthesized last branch.
func flattenSetOp(sel *ast.SelectStmt) (ops []setOperand, orderBy *ast.List, fetch *ast.FetchFirstClause) {
	var walk func(n *ast.SelectStmt, op ast.SetOperation, all bool)
	walk = func(n *ast.SelectStmt, op ast.SetOperation, all bool) {
		if n.Larg != nil {
			walk(n.Larg, op, all)
			if n.Op == ast.SETOP_NONE {
				if n.OrderBy != nil {
					orderBy = n.OrderBy
				}
				if n.FetchFirst != nil {
					fetch = n.FetchFirst
				}
				return
			}
		} else {
			ops = append(ops, setOperand{sel: n, op: op, all: all})
		}
		if n.Op != ast.SETOP_NONE && n.Rarg != nil {
			walk(n.Rarg, n.Op, n.SetAll)
		}
	}
	walk(sel, ast.SETOP_NONE, false)
	if len(ops) > 1 {
		last := &ops[len(ops)-1]
		last.tail = true
		if last.sel.OrderBy != nil {
			orderBy = last.sel.OrderBy
		}
		if last.sel.FetchFirst != nil {
			fetch = last.sel.FetchFirst
		}
	}
	return ops, orderBy, fetch
}

// query analyzes a query, compound or not, under outer. top reports the
// view's own query, whose * Oracle expands in the stored text.
func (a *queryAnalyzer) query(sel *ast.SelectStmt, outer *qscope, ctes *cteEnv, top bool) (*queryResult, error) {
	if sel.WithClause != nil {
		env, err := a.withClause(sel.WithClause, outer, ctes)
		if err != nil {
			return nil, err
		}
		ctes = env
	}
	ops, orderBy, fetch := flattenSetOp(sel)
	if len(ops) == 1 {
		return a.block(ops[0].sel, outer, ctes, top, true)
	}
	var results []*queryResult
	for _, op := range ops {
		r, err := a.block(op.sel, outer, ctes, top, !op.tail)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	first := results[0]
	for _, r := range results[1:] {
		if len(r.cols) != len(first.cols) {
			return nil, oraError("01789", "query block has incorrect number of result columns")
		}
	}
	cols := make([]qcol, len(first.cols))
	for i := range first.cols {
		t := first.cols[i].typ
		for _, r := range results[1:] {
			merged, err := mergeSetOpTypes(t, r.cols[i].typ)
			if err != nil {
				return nil, err
			}
			t = merged
		}
		// The first branch names the column, and its expression still needs
		// an alias in a view (ORA-00998).
		cols[i] = qcol{name: first.cols[i].name, typ: t, expr: first.cols[i].expr}
	}
	// A compound query compares whole rows unless every operator is ALL.
	for i, op := range ops {
		if i == 0 {
			continue
		}
		if !op.all || op.op != ast.SETOP_UNION {
			for _, c := range cols {
				if c.typ.isLOB() {
					return nil, oraError("22848", "cannot use %s type as comparison key", c.typ.kind)
				}
			}
			break
		}
	}
	if orderBy != nil {
		for _, item := range listItems(orderBy) {
			sb, ok := item.(*ast.SortBy)
			if !ok {
				continue
			}
			if err := a.compoundOrderItem(sb.Expr, cols); err != nil {
				return nil, err
			}
		}
	}
	if fetch != nil {
		if err := a.fetchClause(fetch, &qscope{parent: outer, ctes: ctes, from: &qfrom{}}); err != nil {
			return nil, err
		}
	}
	return &queryResult{cols: cols}, nil
}

// compoundOrderItem checks an ORDER BY item of a compound query: a
// position, or a name or expression of the first branch's result.
func (a *queryAnalyzer) compoundOrderItem(expr ast.ExprNode, cols []qcol) error {
	if n, ok := unparen(expr).(*ast.NumberLiteral); ok && !n.IsFloat {
		if n.Ival < 1 || int(n.Ival) > len(cols) {
			return oraError("01785", "ORDER BY item must be the number of a SELECT-list expression")
		}
		return nil
	}
	if ref, ok := unparen(expr).(*ast.ColumnRef); ok && ref.Table == "" {
		for _, c := range cols {
			if c.name == ref.Column {
				return nil
			}
		}
		return oraError("00904", "%q: invalid identifier", ref.Column)
	}
	return oraError("01785", "ORDER BY item must be the number of a SELECT-list expression")
}

// mergeSetOpTypes is the type of a compound query's column from two
// branches' types (ORA-01790 when their type families differ).
func mergeSetOpTypes(a, b qtype) (qtype, error) {
	if a.isNull() {
		if b.isNull() {
			return a, nil
		}
		if b.kind == "CHAR" {
			b.kind = "VARCHAR2"
		}
		return b, nil
	}
	if b.isNull() {
		if a.kind == "CHAR" {
			a.kind = "VARCHAR2"
		}
		return a, nil
	}
	if sameType(a, b) {
		return a, nil
	}
	if a.kind == "UNKNOWN" || b.kind == "UNKNOWN" {
		// A type this analyzer does not model proves no mismatch.
		return simpleType("UNKNOWN"), nil
	}
	switch {
	case a.isCharacter() && b.isCharacter():
		if a.isNational() || b.isNational() {
			n := max(a.chars, b.chars)
			return characterType("NVARCHAR2", n, false), nil
		}
		bytes := max(a.bytes, b.bytes)
		t := varchar2Bytes(bytes)
		if a.charSem || b.charSem {
			t.charSem = true
			t.chars = bytes / dbCharBytes
		}
		return t, nil
	case a.kind == "NUMBER" && b.kind == "NUMBER":
		return numberType(), nil
	case a.isNumeric() && b.isNumeric():
		if a.kind == "BINARY_DOUBLE" || b.kind == "BINARY_DOUBLE" {
			return simpleType("BINARY_DOUBLE"), nil
		}
		if a.kind == "BINARY_FLOAT" || b.kind == "BINARY_FLOAT" {
			return simpleType("BINARY_FLOAT"), nil
		}
		return numberType(), nil
	case a.isDatetime() && b.isDatetime():
		if a.kind == "DATE" {
			return b, nil
		}
		if b.kind == "DATE" {
			return a, nil
		}
		if a.kind == b.kind {
			a.frac = max(a.frac, b.frac)
			return a, nil
		}
		return a, nil
	case a.kind == b.kind:
		return a, nil
	}
	return qtype{}, oraError("01790", "expression must have same datatype as corresponding expression")
}

// withClause registers a WITH clause's subqueries and functions.
func (a *queryAnalyzer) withClause(w *ast.WithClause, outer *qscope, ctes *cteEnv) (*cteEnv, error) {
	for _, item := range listItems(w.PLSQLDeclarations) {
		switch f := item.(type) {
		case *ast.CreateFunctionStmt:
			if f.Name != nil {
				a.plsqlFuncs[f.Name.Name] = typeNameType(f.ReturnType)
			}
		}
	}
	// An element sees itself (a recursive WITH) and the elements before it,
	// not the ones after.
	env := &cteEnv{parent: ctes, defs: map[string]*cteDef{}}
	for _, item := range listItems(w.CTEs) {
		c, ok := item.(*ast.CTE)
		if !ok {
			continue
		}
		d := &cteDef{name: c.Name, node: c}
		env.defs[c.Name] = d
		if err := a.analyzeCTE(d, outer, env); err != nil {
			return nil, err
		}
	}
	return env, nil
}

func (a *queryAnalyzer) analyzeCTE(d *cteDef, outer *qscope, env *cteEnv) error {
	if d.analyzed {
		return nil
	}
	c := d.node
	if c.SubAV != nil {
		using := c.SubAV.Using
		if using != nil && using.Schema == "" && using.Name == c.Name {
			return oraError("18385", "recursive WITH clause not allowed in a derived analytic view")
		}
		return a.analyticView(c.SubAV)
	}
	sel, ok := c.Query.(*ast.SelectStmt)
	if !ok || sel == nil {
		d.analyzed = true
		return nil
	}
	names := stringList(c.Columns)
	recursive := selectReferencesTable(sel, c.Name)
	if recursive {
		if len(names) == 0 {
			return oraError("32039", "missing column alias list in recursive WITH clause element %s", c.Name)
		}
		ops, _, _ := flattenSetOp(sel)
		if len(ops) != 2 || !ops[1].all || ops[1].op != ast.SETOP_UNION {
			return oraError("32040", "recursive WITH clause must use a UNION ALL operation")
		}
		anchorBlock := *ops[0].sel
		anchorBlock.Op, anchorBlock.Rarg = ast.SETOP_NONE, nil
		if selectReferencesTable(&anchorBlock, c.Name) {
			return oraError("32043", "recursive WITH clause needs an initialization branch")
		}
		d.analyzing = true
		anchor, err := a.block(ops[0].sel, outer, env, false, true)
		d.analyzing = false
		if err != nil {
			return err
		}
		if len(anchor.cols) != len(names) {
			return oraError("32038", "number of WITH clause column names does not match number of elements in select list")
		}
		d.cols = make([]qcol, len(names))
		for i, n := range names {
			d.cols[i] = qcol{name: n, typ: anchor.cols[i].typ}
		}
		d.analyzed = true
		res, err := a.query(sel, outer, env, false)
		if err != nil {
			return err
		}
		for i := range d.cols {
			if i < len(res.cols) {
				d.cols[i].typ = res.cols[i].typ
			}
		}
	} else {
		d.analyzing = true
		res, err := a.query(sel, outer, env, false)
		d.analyzing = false
		if err != nil {
			return err
		}
		if len(names) > 0 {
			if len(names) != len(res.cols) {
				return oraError("32038", "number of WITH clause column names does not match number of elements in select list")
			}
			cols := make([]qcol, len(names))
			for i, n := range names {
				cols[i] = res.cols[i]
				cols[i].name = n
			}
			d.cols = cols
		} else {
			d.cols = res.cols
		}
		d.analyzed = true
	}
	if c.Search != nil && c.Search.SetColumn != "" {
		d.cols = append(d.cols, qcol{name: c.Search.SetColumn, typ: numberType()})
	}
	if c.Cycle != nil && c.Cycle.SetColumn != "" {
		n := 1
		if s, ok := c.Cycle.CycleValue.(*ast.StringLiteral); ok {
			n = max(len(s.Val), 1)
		}
		d.cols = append(d.cols, qcol{name: c.Cycle.SetColumn, typ: varchar2Bytes(n)})
	}
	if recursive {
		for i := range d.cols {
			d.cols[i].notNull = false
		}
	}
	return nil
}

// selectReferencesTable reports whether a query names table in a FROM
// clause, as a recursive WITH element does.
func selectReferencesTable(sel *ast.SelectStmt, table string) bool {
	found := false
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		if found || !v.IsValid() {
			return
		}
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if v.IsNil() {
				return
			}
			if tr, ok := v.Interface().(*ast.TableRef); ok {
				if tr.Name != nil && tr.Name.Schema == "" && tr.Name.Name == table {
					found = true
					return
				}
			}
			walk(v.Elem())
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() {
					walk(v.Field(i))
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
			}
		}
	}
	walk(reflect.ValueOf(sel))
	return found
}

// block analyzes one query block. ownOrder reports whether its ORDER BY and
// FETCH are its own rather than a compound query's.
func (a *queryAnalyzer) block(sel *ast.SelectStmt, outer *qscope, ctes *cteEnv, top, ownOrder bool) (*queryResult, error) {
	if sel.WithClause != nil && outer != nil {
		env, err := a.withClause(sel.WithClause, outer, ctes)
		if err != nil {
			return nil, err
		}
		ctes = env
	}
	if sel.ForUpdate != nil && top {
		return nil, oraError("03048", "SQL reserved word 'FOR' is not syntactically valid")
	}
	sc := &qscope{parent: outer, ctes: ctes, from: &qfrom{}}
	if sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
		for _, item := range listItems(sel.TargetList) {
			if rt, ok := item.(*ast.ResTarget); ok {
				if _, ok := rt.Expr.(*ast.Star); ok {
					return nil, oraError("54058", "'*' not allowed in SELECT expressions without FROM clause")
				}
			}
		}
	}
	for _, item := range listItems(sel.FromClause) {
		te, ok := item.(ast.TableExpr)
		if !ok {
			continue
		}
		f, err := a.fromItem(te, sc)
		if err != nil {
			return nil, err
		}
		sc.from.rels = append(sc.from.rels, f.rels...)
		sc.from.items = append(sc.from.items, f.items...)
		sc.from.flat = append(sc.from.flat, f.flat...)
	}
	sc.hierarchical = sel.Hierarchical != nil
	for _, w := range sel.WindowDefs {
		if w == nil || w.Spec == nil {
			continue
		}
		if sc.windows == nil {
			sc.windows = map[string]*ast.WindowSpec{}
		}
		sc.windows[w.Name] = w.Spec
	}
	for _, w := range sel.WindowDefs {
		if w == nil || w.Spec == nil {
			continue
		}
		if err := a.windowSpec(w.Spec, sc, map[string]bool{w.Name: true}); err != nil {
			return nil, err
		}
	}

	// Oracle's (+) outer joins in WHERE null-extend the tables they mark.
	if sel.WhereClause != nil {
		if err := a.plusJoins(sel.WhereClause, sc); err != nil {
			return nil, err
		}
		if err := a.condition(sel.WhereClause, sc, clauseWhere); err != nil {
			return nil, err
		}
	}
	if h := sel.Hierarchical; h != nil {
		if h.StartWith != nil {
			if err := a.condition(h.StartWith, sc, clauseWhere); err != nil {
				return nil, err
			}
		}
		if h.ConnectBy != nil {
			if err := a.condition(h.ConnectBy, sc, clauseWhere); err != nil {
				return nil, err
			}
		}
	}

	// The select list.
	var cols []qcol
	type target struct {
		rt   *ast.ResTarget
		cols []qcol
	}
	// starRefs are the column references a * or t.* in the select list
	// stands for, which the grouping rules check as written columns.
	var starRefs []ast.ExprNode
	var targets []target
	for _, item := range listItems(sel.TargetList) {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		tc, err := a.selectTarget(rt, sc, top)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target{rt: rt, cols: tc})
		cols = append(cols, tc...)
		starRefs = append(starRefs, a.starReferences(rt, sc)...)
	}
	sc.aliases = cols

	// GROUP BY, HAVING and the aggregate rules.
	var groupKeys []string
	positional := false
	for _, item := range listItems(sel.GroupClause) {
		for _, g := range groupingExprs(item) {
			if n, ok := unparen(g).(*ast.NumberLiteral); ok && !n.IsFloat {
				positional = true
			}
			if err := a.groupByExpr(g, sc); err != nil {
				return nil, err
			}
			groupKeys = append(groupKeys, a.groupKey(g, sc))
			if ref, ok := unparen(g).(*ast.ColumnRef); ok && ref.Table == "" {
				// A select alias grouped by (23ai) groups by what it names.
				for _, t := range targets {
					if t.rt.Name == ref.Column && t.rt.Expr != nil && !a.resolvesInFrom(ref, sc) {
						groupKeys = append(groupKeys, a.groupKey(t.rt.Expr, sc))
					}
				}
			}
		}
	}
	if sel.HavingClause != nil {
		if err := a.condition(sel.HavingClause, sc, clauseHaving); err != nil {
			return nil, err
		}
	}
	var exprs []ast.ExprNode
	for _, t := range targets {
		if t.rt.Expr != nil {
			exprs = append(exprs, t.rt.Expr)
		}
	}
	if sel.HavingClause != nil {
		exprs = append(exprs, sel.HavingClause)
	}
	exprs = append(exprs, starRefs...)
	if ownOrder {
		// ORDER BY is held to the grouping rules too, unless it names a
		// select-list item by position or alias.
		for _, item := range listItems(sel.OrderBy) {
			sb, ok := item.(*ast.SortBy)
			if !ok || sb.Expr == nil {
				continue
			}
			if n, ok := unparen(sb.Expr).(*ast.NumberLiteral); ok && !n.IsFloat {
				continue
			}
			if ref, ok := unparen(sb.Expr).(*ast.ColumnRef); ok && ref.Table == "" && slices.ContainsFunc(cols, func(c qcol) bool { return c.name == ref.Column }) {
				continue
			}
			exprs = append(exprs, sb.Expr)
		}
	}
	grouped := len(listItems(sel.GroupClause)) > 0
	nested := false
	for _, x := range exprs {
		if containsAggregate(x, false) {
			grouped = true
		}
		if containsNestedAggregate(x) {
			nested = true
		}
	}
	if nested && len(listItems(sel.GroupClause)) == 0 {
		return nil, oraError("00978", "nested group function without GROUP BY")
	}
	if grouped {
		for _, x := range exprs {
			keyOf := func(n ast.Node) string { return a.groupKey(n, sc) }
			if err := checkGrouped(x, groupKeys, keyOf, len(listItems(sel.GroupClause)) > 0, positional, nested); err != nil {
				return nil, err
			}
		}
		// GROUP BY keeps a grouping column's NOT NULL; an aggregate result
		// is nullable, which typing already says.
	}

	// The DISTINCT and grouping keys may not be LOBs.
	if sel.Distinct || sel.UniqueKw {
		for _, c := range cols {
			if c.typ.isLOB() {
				return nil, oraError("22848", "cannot use %s type as comparison key", c.typ.kind)
			}
		}
	}
	if sel.QualifyClause != nil {
		if err := a.condition(sel.QualifyClause, sc, clauseOrder); err != nil {
			return nil, err
		}
	}
	if ownOrder {
		for _, item := range listItems(sel.OrderBy) {
			sb, ok := item.(*ast.SortBy)
			if !ok {
				continue
			}
			if err := a.orderItem(sb.Expr, sc, len(cols)); err != nil {
				return nil, err
			}
		}
		if sel.FetchFirst != nil {
			if err := a.fetchClause(sel.FetchFirst, sc); err != nil {
				return nil, err
			}
		}
	}
	if sel.ModelClause != nil {
		cols = a.modelColumns(sel, cols)
	}
	return &queryResult{cols: cols}, nil
}

// modelColumns applies a MODEL clause: a measure loses its NOT NULL, and a
// numeric measure its precision.
func (a *queryAnalyzer) modelColumns(sel *ast.SelectStmt, cols []qcol) []qcol {
	mm := sel.ModelClause.MainModel
	if mm == nil || mm.ColumnClauses == nil {
		return cols
	}
	measures := map[string]bool{}
	for _, item := range listItems(mm.ColumnClauses.Measures) {
		if rt, ok := item.(*ast.ResTarget); ok {
			name := rt.Name
			if name == "" {
				if ref, ok := unparen(rt.Expr).(*ast.ColumnRef); ok {
					name = ref.Column
				}
			}
			measures[name] = true
		}
	}
	out := make([]qcol, len(cols))
	for i, c := range cols {
		c.notNull = false
		if measures[c.name] && c.typ.kind == "NUMBER" {
			c.typ = numberType()
		}
		out[i] = c
	}
	return out
}

// fetchClause checks a row-limiting clause's offset and count, which a view
// may not bind (ORA-01027).
func (a *queryAnalyzer) fetchClause(f *ast.FetchFirstClause, sc *qscope) error {
	for _, x := range []ast.ExprNode{f.Offset, f.Count} {
		if x == nil {
			continue
		}
		if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
			return err
		}
	}
	return nil
}

func groupingExprs(item ast.Node) []ast.ExprNode {
	switch g := item.(type) {
	case *ast.RollupClause:
		var out []ast.ExprNode
		for _, x := range listItems(g.Args) {
			out = append(out, groupingExprs(x)...)
		}
		return out
	case *ast.CubeClause:
		var out []ast.ExprNode
		for _, x := range listItems(g.Args) {
			out = append(out, groupingExprs(x)...)
		}
		return out
	case *ast.GroupingSetsClause:
		var out []ast.ExprNode
		for _, x := range listItems(g.Sets) {
			out = append(out, groupingExprs(x)...)
		}
		return out
	case *ast.List:
		var out []ast.ExprNode
		for _, x := range g.Items {
			out = append(out, groupingExprs(x)...)
		}
		return out
	case *ast.ParenExpr:
		return groupingExprs(g.Expr)
	case ast.ExprNode:
		return []ast.ExprNode{g}
	}
	return nil
}

func (a *queryAnalyzer) groupByExpr(g ast.ExprNode, sc *qscope) error {
	if containsAggregate(g, false) {
		return oraError("00934", "group function is not allowed here")
	}
	if ref, ok := unparen(g).(*ast.ColumnRef); ok && ref.Table == "" && !a.resolvesInFrom(ref, sc) {
		for _, c := range sc.aliases {
			if c.name == ref.Column {
				return nil
			}
		}
	}
	t, _, err := a.expr(g, sc, clauseGroup)
	if err != nil {
		return err
	}
	if t.isLOB() {
		return oraError("22848", "cannot use %s type as comparison key", t.kind)
	}
	return nil
}

// resolvesInFrom reports whether an unqualified reference names a column of
// the block's own FROM clause.
func (a *queryAnalyzer) resolvesInFrom(ref *ast.ColumnRef, sc *qscope) bool {
	if ref.Table != "" {
		return false
	}
	_, found, _ := a.lookupIn(sc.from, ref.Column)
	return found
}

// orderItem checks an ORDER BY item of a query block.
func (a *queryAnalyzer) orderItem(expr ast.ExprNode, sc *qscope, ncols int) error {
	if n, ok := unparen(expr).(*ast.NumberLiteral); ok && !n.IsFloat {
		if n.Ival < 1 || int(n.Ival) > ncols {
			return oraError("01785", "ORDER BY item must be the number of a SELECT-list expression")
		}
		return nil
	}
	if ref, ok := unparen(expr).(*ast.ColumnRef); ok && ref.Table == "" {
		for _, c := range sc.aliases {
			if c.name == ref.Column {
				if c.typ.isLOB() {
					return oraError("22848", "cannot use %s type as comparison key", c.typ.kind)
				}
				return nil
			}
		}
	}
	t, _, err := a.expr(expr, sc, clauseOrder)
	if err != nil {
		return err
	}
	if t.isLOB() {
		return oraError("22848", "cannot use %s type as comparison key", t.kind)
	}
	return nil
}

// selectTarget returns the columns one select-list item gives.
func (a *queryAnalyzer) selectTarget(rt *ast.ResTarget, sc *qscope, top bool) ([]qcol, error) {
	if rt.AliasBind {
		return nil, oraError("01027", "bind variables not allowed for data definition operations")
	}
	switch x := rt.Expr.(type) {
	case *ast.Star:
		items, err := a.starColumns(sc.from, "", "")
		if err != nil {
			return nil, err
		}
		cols := make([]qcol, 0, len(items))
		var names []string
		for _, it := range items {
			c := sc.from.itemCol(it)
			c.expr = false
			cols = append(cols, c)
			names = append(names, quoteIdent(c.name))
		}
		if top {
			a.edits = append(a.edits, textEdit{start: x.Loc.Start, end: x.Loc.End, text: strings.Join(names, ",")})
		}
		return cols, nil
	case *ast.ColumnRef:
		if x.Column == "*" {
			items, err := a.qualifiedStar(sc, x)
			if err != nil {
				return nil, err
			}
			cols := make([]qcol, 0, len(items))
			var parts []string
			qual := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(a.e.nodeText(x)), "*"))
			for _, c := range items {
				c.expr = false
				cols = append(cols, c)
				parts = append(parts, qual+quoteIdent(c.name))
			}
			if top {
				a.edits = append(a.edits, textEdit{start: x.Loc.Start, end: x.Loc.End, text: strings.Join(parts, ",")})
			}
			return cols, nil
		}
	}
	t, notNull, err := a.expr(rt.Expr, sc, clauseSelect)
	if err != nil {
		return nil, err
	}
	c := qcol{typ: t, notNull: notNull}
	if rt.Name != "" {
		c.name = rt.Name
		return []qcol{c}, nil
	}
	if ref, ok := unparen(rt.Expr).(*ast.ColumnRef); ok && a.isColumn(ref, sc) {
		c.name = ref.Column
		return []qcol{c}, nil
	}
	c.expr = true
	c.name = derivedName(a.e.nodeText(rt.Expr))
	return []qcol{c}, nil
}

// derivedName is the name Oracle gives an unaliased expression of an inline
// view: its text, uppercased outside quotes, without blanks.
func derivedName(text string) string {
	var b strings.Builder
	inQuote := byte(0)
	for i := 0; i < len(text); i++ {
		ch := text[i]
		switch {
		case inQuote != 0:
			b.WriteByte(ch)
			if ch == inQuote {
				inQuote = 0
			}
		case ch == '\'' || ch == '"':
			inQuote = ch
			b.WriteByte(ch)
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
		default:
			b.WriteString(strings.ToUpper(string(ch)))
		}
	}
	name := b.String()
	if len(name) > 128 {
		name = name[:128]
	}
	return name
}

// starColumns lists what * stands for in a FROM clause.
func (a *queryAnalyzer) starColumns(from *qfrom, _, _ string) ([]fromCol, error) {
	seen := map[string]int{}
	flat := map[*qrel]bool{}
	for _, r := range from.flat {
		flat[r] = true
		if r.name != "" {
			seen[r.name]++
		}
	}
	var out []fromCol
	for _, it := range from.items {
		if it.rel != nil && it.rel.cols[it.idx].hidden {
			continue
		}
		if it.rel != nil && flat[it.rel] && it.rel.name != "" && seen[it.rel.name] > 1 {
			return nil, oraError("00918", "%s: column ambiguously specified", it.rel.cols[it.idx].name)
		}
		out = append(out, it)
	}
	return out, nil
}

// qualifiedStar lists what t.* or s.t.* stands for.
func (a *queryAnalyzer) qualifiedStar(sc *qscope, ref *ast.ColumnRef) ([]qcol, error) {
	var matches []*qrel
	for _, r := range sc.from.rels {
		if r.name == ref.Table && (ref.Schema == "" || ref.Schema == r.schema) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 0 {
		if ref.Schema != "" {
			return nil, oraError("00964", "table name not in FROM list")
		}
		return nil, oraError("00904", "%q: invalid identifier", ref.Table)
	}
	if len(matches) > 1 {
		return nil, oraError("00918", "column ambiguously specified")
	}
	r := matches[0]
	var cols []qcol
	for i := range r.cols {
		if r.cols[i].hidden {
			continue
		}
		cols = append(cols, r.col(i))
	}
	return cols, nil
}

// isColumn reports whether a reference names a column rather than a
// keyword such as CURRENT_DATE.
func (a *queryAnalyzer) isColumn(ref *ast.ColumnRef, sc *qscope) bool {
	if ref.Table != "" {
		return true
	}
	for s := sc; s != nil; s = s.parent {
		if _, found, _ := a.lookupIn(s.from, ref.Column); found {
			return true
		}
	}
	return false
}

// lookupIn finds an unqualified column in a FROM clause. found is false
// when no item has it; an error reports an ambiguous name.
func (a *queryAnalyzer) lookupIn(from *qfrom, name string) (qcol, bool, error) {
	if from == nil {
		return qcol{}, false, nil
	}
	var hits []qcol
	var rels []string
	for _, it := range from.items {
		if it.rel == nil {
			// A merged USING or NATURAL column is one candidate among the
			// FROM clause's columns of that name.
			if it.col.name == name {
				hits = append(hits, it.col)
				rels = append(rels, "")
			}
			continue
		}
		c := it.rel.cols[it.idx]
		if c.name != name || it.rel.merged[name] {
			continue
		}
		hits = append(hits, it.rel.col(it.idx))
		rels = append(rels, it.rel.name)
	}
	switch len(hits) {
	case 0:
		return qcol{}, false, nil
	case 1:
		return hits[0], true, nil
	}
	return qcol{}, true, oraError("00918", "%s: column ambiguously specified - appears in %s", name, strings.Join(rels, " and "))
}

// lookupQualified finds t.c or s.t.c in a FROM clause.
func (a *queryAnalyzer) lookupQualified(from *qfrom, ref *ast.ColumnRef) (qcol, bool, error) {
	if from == nil {
		return qcol{}, false, nil
	}
	var matches []*qrel
	for _, r := range from.rels {
		if r.name == ref.Table && (ref.Schema == "" || ref.Schema == r.schema) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 0 {
		return qcol{}, false, nil
	}
	r := matches[0]
	for i, c := range r.cols {
		if c.name != ref.Column {
			continue
		}
		if r.merged[c.name] {
			return qcol{}, true, oraError("25154", "column part of USING clause cannot have qualifier")
		}
		if len(matches) > 1 {
			return qcol{}, true, oraError("00918", "%s: column ambiguously specified", ref.Column)
		}
		return r.col(i), true, nil
	}
	// ROWID and ORA_ROWSCN qualify by the item they read.
	switch ref.Column {
	case "ROWID":
		return qcol{name: ref.Column, typ: simpleType("ROWID")}, true, nil
	case "ORA_ROWSCN":
		return qcol{name: ref.Column, typ: numberType()}, true, nil
	}
	return qcol{}, true, oraError("00904", "%q.%q: invalid identifier", ref.Table, ref.Column)
}

// resolveColumn resolves a column reference in the scope chain.
func (a *queryAnalyzer) resolveColumn(ref *ast.ColumnRef, sc *qscope) (qcol, bool, error) {
	for s := sc; s != nil; s = s.parent {
		var (
			c     qcol
			found bool
			err   error
		)
		if ref.Table == "" {
			c, found, err = a.lookupIn(s.from, ref.Column)
		} else {
			c, found, err = a.lookupQualified(s.from, ref)
		}
		if err != nil {
			return qcol{}, true, err
		}
		if found {
			return c, true, nil
		}
	}
	return qcol{}, false, nil
}

// plusJoins applies Oracle's (+) outer joins in a WHERE clause: each
// conjunct may mark one table, which the join null-extends.
func (a *queryAnalyzer) plusJoins(where ast.ExprNode, sc *qscope) error {
	for _, conj := range conjuncts(where) {
		marked := map[*qrel]bool{}
		walkExpr(conj, func(n ast.Node) bool {
			if _, ok := n.(*ast.SubqueryExpr); ok {
				return false
			}
			ref, ok := n.(*ast.ColumnRef)
			if !ok || !ref.OuterJoin {
				return true
			}
			if r := a.relOf(ref, sc.from); r != nil {
				marked[r] = true
			}
			return true
		})
		if len(marked) > 1 {
			return oraError("01468", "a predicate may reference only one outer-joined table")
		}
		for r := range marked {
			r.nullExt = true
		}
	}
	return nil
}

// relOf is the FROM item a column reference reads.
func (a *queryAnalyzer) relOf(ref *ast.ColumnRef, from *qfrom) *qrel {
	for _, r := range from.rels {
		if ref.Table != "" {
			if r.name == ref.Table {
				return r
			}
			continue
		}
		for _, c := range r.cols {
			if c.name == ref.Column {
				return r
			}
		}
	}
	return nil
}

func conjuncts(x ast.ExprNode) []ast.ExprNode {
	switch b := unparen(x).(type) {
	case *ast.BoolExpr:
		if b.Boolop == ast.BOOL_AND {
			var out []ast.ExprNode
			for _, arg := range listItems(b.Args) {
				if e, ok := arg.(ast.ExprNode); ok {
					out = append(out, conjuncts(e)...)
				}
			}
			return out
		}
	case *ast.BinaryExpr:
		if strings.EqualFold(b.Op, "AND") {
			return append(conjuncts(b.Left), conjuncts(b.Right)...)
		}
	}
	return []ast.ExprNode{x}
}

func unparen(x ast.ExprNode) ast.ExprNode {
	for {
		p, ok := x.(*ast.ParenExpr)
		if !ok || p == nil {
			return x
		}
		x = p.Expr
	}
}

// groupKey renders an expression's tree without positions, each column
// reference as the column it resolves to, so two expressions written alike,
// or t.id and id, compare equal in a GROUP BY check.
func (a *queryAnalyzer) groupKey(n ast.Node, sc *qscope) string {
	var b strings.Builder
	writeKey(reflect.ValueOf(n), &b, func(ref *ast.ColumnRef) string { return a.columnIdentity(ref, sc) })
	return b.String()
}

// columnIdentity names the column a reference resolves to, or "" when it
// resolves to none.
func (a *queryAnalyzer) columnIdentity(ref *ast.ColumnRef, sc *qscope) string {
	for s := sc; s != nil; s = s.parent {
		if s.from == nil {
			continue
		}
		if ref.Table == "" {
			var hits []string
			for i, it := range s.from.items {
				if it.rel == nil {
					if it.col.name == ref.Column {
						hits = append(hits, fmt.Sprintf("merged:%p#%d", s.from, i))
					}
					continue
				}
				if it.rel.cols[it.idx].name == ref.Column && !it.rel.merged[ref.Column] {
					hits = append(hits, fmt.Sprintf("col:%p.%s", it.rel, ref.Column))
				}
			}
			if len(hits) == 1 {
				return hits[0]
			}
			if len(hits) > 1 {
				return ""
			}
			continue
		}
		for _, r := range s.from.rels {
			if r.name == ref.Table && (ref.Schema == "" || ref.Schema == r.schema) {
				return fmt.Sprintf("col:%p.%s", r, ref.Column)
			}
		}
	}
	return ""
}

func writeKey(v reflect.Value, b *strings.Builder, colKey func(*ast.ColumnRef) string) {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			b.WriteString("nil")
			return
		}
		if ref, ok := v.Interface().(*ast.ColumnRef); ok && colKey != nil {
			if k := colKey(ref); k != "" {
				b.WriteString(k)
				return
			}
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(ast.ParenExpr{}) {
			writeKey(v.FieldByName("Expr"), b, colKey)
			return
		}
		b.WriteString(v.Type().Name())
		b.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if f.Name == "Loc" || !f.IsExported() {
				continue
			}
			b.WriteString(f.Name)
			b.WriteByte(':')
			writeKey(v.Field(i), b, colKey)
			b.WriteByte(' ')
		}
		b.WriteByte('}')
	case reflect.Slice:
		b.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			writeKey(v.Index(i), b, colKey)
			b.WriteByte(',')
		}
		b.WriteByte(']')
	default:
		fmt.Fprintf(b, "%v", v.Interface())
	}
}

// walkExpr visits an expression tree depth first; fn returns false to skip
// a node's children.
func walkExpr(n ast.Node, fn func(ast.Node) bool) {
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		if !v.IsValid() {
			return
		}
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if v.IsNil() {
				return
			}
			if node, ok := v.Interface().(ast.Node); ok {
				if !fn(node) {
					return
				}
			}
			walk(v.Elem())
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() && v.Type().Field(i).Name != "Loc" {
					walk(v.Field(i))
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
			}
		}
	}
	walk(reflect.ValueOf(n))
}

// tableColumns reads a catalog table's or view's columns.
func tableColumns(cols []*metadatapb.ColumnMetadata) []qcol {
	out := make([]qcol, 0, len(cols))
	for _, c := range cols {
		out = append(out, qcol{
			name:    c.GetName(),
			typ:     parseColumnType(c.GetType()),
			notNull: !c.GetNullable(),
			hidden:  c.GetIsInvisible(),
		})
	}
	return out
}

// dualColumns is SYS.DUAL's one column.
func dualColumns() []qcol {
	return []qcol{{name: "DUMMY", typ: varchar2Bytes(1)}}
}

// viewColumnsFor names a view's columns as the catalog holds them.
func viewColumnsFor(cols []qcol, explicit []string, hidden map[int]bool) []*metadatapb.ColumnMetadata {
	out := make([]*metadatapb.ColumnMetadata, 0, len(cols))
	for i, c := range cols {
		name := c.name
		if i < len(explicit) {
			name = explicit[i]
		}
		typ := c.typ.String()
		if c.typ.kind == "UNKNOWN" {
			typ = ""
		}
		out = append(out, &metadatapb.ColumnMetadata{
			Name:        name,
			Position:    int32(i + 1),
			Type:        typ,
			Nullable:    !c.notNull,
			IsInvisible: hidden[i],
		})
	}
	return out
}

// checkViewColumnNames applies a view's naming rules to its derived
// columns: an expression needs an alias (ORA-00998), and names are unique
// (ORA-00957).
func checkViewColumnNames(cols []qcol) error {
	for _, c := range cols {
		if c.expr {
			return oraError("00998", "must name this expression with a column alias")
		}
	}
	seen := map[string]bool{}
	for _, c := range cols {
		if seen[c.name] {
			return oraError("00957", "duplicate column name")
		}
		seen[c.name] = true
	}
	return nil
}

// quoteIdent quotes an identifier as Oracle writes it in a view's text.
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// starReferences lists the column references a select-list * or t.* stands
// for, qualified by the item that holds each column.
func (a *queryAnalyzer) starReferences(rt *ast.ResTarget, sc *qscope) []ast.ExprNode {
	var qualifier string
	switch x := rt.Expr.(type) {
	case *ast.Star:
	case *ast.ColumnRef:
		if x.Column != "*" {
			return nil
		}
		qualifier = x.Table
	default:
		return nil
	}
	var refs []ast.ExprNode
	for _, it := range sc.from.items {
		if it.rel == nil {
			if qualifier == "" {
				refs = append(refs, &ast.ColumnRef{Column: it.col.name, Loc: ast.NoLoc()})
			}
			continue
		}
		if qualifier != "" && it.rel.name != qualifier {
			continue
		}
		c := it.rel.cols[it.idx]
		if c.hidden || it.rel.merged[c.name] {
			continue
		}
		refs = append(refs, &ast.ColumnRef{Table: it.rel.name, Column: c.name, Loc: ast.NoLoc()})
	}
	return refs
}
