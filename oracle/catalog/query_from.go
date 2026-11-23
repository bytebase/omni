package catalog

import (
	"strings"
	"unicode"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
)

// fromItem analyzes one FROM item in the block whose scope is sc. sc.from
// holds the items before it, which only a lateral item may read.
func (a *queryAnalyzer) fromItem(te ast.TableExpr, sc *qscope) (*qfrom, error) {
	switch n := te.(type) {
	case *ast.TableRef:
		return a.tableRef(n, sc)
	case *ast.SubqueryRef:
		sel, ok := n.Subquery.(*ast.SelectStmt)
		if !ok {
			return &qfrom{}, nil
		}
		res, err := a.query(sel, sc.parent, sc.ctes, false)
		if err != nil {
			return nil, err
		}
		return relFrom(&qrel{name: aliasName(n.Alias), cols: res.cols}), nil
	case *ast.LateralRef:
		sel, ok := n.Subquery.(*ast.SelectStmt)
		if !ok {
			return &qfrom{}, nil
		}
		res, err := a.query(sel, sc, sc.ctes, false)
		if err != nil {
			return nil, err
		}
		return relFrom(&qrel{name: aliasName(n.Alias), cols: res.cols}), nil
	case *ast.JoinClause:
		return a.join(n, sc)
	case *ast.TableCollectionExpr:
		return a.tableCollection(n, sc)
	case *ast.PivotClause:
		return a.pivot(n, sc)
	case *ast.UnpivotClause:
		return a.unpivot(n, sc)
	case *ast.MatchRecognizeClause:
		return a.matchRecognize(n, sc)
	case *ast.InlineAnalyticView:
		if err := a.analyticView(n.Clause); err != nil {
			return nil, err
		}
		return &qfrom{}, nil
	case *ast.JsonTableRef:
		if _, _, err := a.expr(n.Expr, sc, clauseFrom); err != nil {
			return nil, err
		}
		if err := a.jsonTablePaths(n); err != nil {
			return nil, err
		}
		rel := &qrel{name: aliasName(n.Alias)}
		rel.cols = jsonTableColumns(n.Columns)
		return relFrom(rel), nil
	case *ast.XmlTableRef:
		for _, x := range []ast.ExprNode{n.XPath, n.Passing} {
			if x == nil {
				continue
			}
			if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
				return nil, err
			}
		}
		rel := &qrel{name: aliasName(n.Alias)}
		for _, item := range listItems(n.Columns) {
			c, ok := item.(*ast.XmlTableColumn)
			if !ok {
				continue
			}
			for _, x := range []ast.ExprNode{c.Path, c.Default} {
				if x == nil {
					continue
				}
				if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
					return nil, err
				}
			}
			t := varchar2Bytes(4000)
			if c.ForOrdinality {
				t = numberType()
			} else if c.TypeName != nil {
				t = typeNameType(c.TypeName)
			}
			rel.cols = append(rel.cols, qcol{name: c.Name, typ: t})
		}
		return relFrom(rel), nil
	case *ast.SubqueryExpr:
		sel, ok := n.Subquery.(*ast.SelectStmt)
		if !ok {
			return &qfrom{}, nil
		}
		res, err := a.query(sel, sc.parent, sc.ctes, false)
		if err != nil {
			return nil, err
		}
		return relFrom(&qrel{cols: res.cols}), nil
	case *ast.ContainersExpr:
		return a.tableRef(&ast.TableRef{Name: n.Name, Alias: n.Alias, Loc: n.Loc}, sc)
	}
	return &qfrom{}, nil
}

func jsonTableColumns(list *ast.List) []qcol {
	var cols []qcol
	for _, item := range listItems(list) {
		c, ok := item.(*ast.JsonTableColumn)
		if !ok {
			continue
		}
		if c.Nested != nil {
			cols = append(cols, jsonTableColumns(c.Nested.Columns)...)
			continue
		}
		t := varchar2Bytes(4000)
		switch {
		case c.ForOrdinality:
			t = numberType()
		case c.TypeName != nil:
			t = typeNameType(c.TypeName)
		}
		cols = append(cols, qcol{name: c.Name, typ: t})
	}
	return cols
}

func relFrom(r *qrel) *qfrom {
	f := &qfrom{rels: []*qrel{r}, flat: []*qrel{r}}
	for i := range r.cols {
		f.items = append(f.items, fromCol{rel: r, idx: i})
	}
	return f
}

func aliasName(al *ast.Alias) string {
	if al == nil {
		return ""
	}
	return al.Name
}

// tableRef resolves a named FROM item: a WITH element, a table, a view, a
// synonym of one, or DUAL.
func (a *queryAnalyzer) tableRef(n *ast.TableRef, sc *qscope) (*qfrom, error) {
	if n.Name == nil {
		return &qfrom{}, nil
	}
	if n.Dblink != "" || n.Name.DBLink != "" {
		return nil, oraError("02019", "connection description for remote database not found")
	}
	name, schema := n.Name.Name, n.Name.Schema
	exposed := name
	if n.Alias != nil && n.Alias.Name != "" {
		exposed = n.Alias.Name
	}
	if schema == "" {
		if d := sc.ctes.lookup(name); d != nil {
			if d.analyzing {
				// A WITH element that reads itself before its column list
				// is known: Oracle wants an initialization branch.
				return nil, oraError("32043", "recursive WITH clause needs an initialization branch")
			}
			if !d.analyzed {
				if err := a.analyzeCTE(d, sc.parent, sc.ctes); err != nil {
					return nil, err
				}
			}
			cols := make([]qcol, len(d.cols))
			copy(cols, d.cols)
			if err := a.tableClauses(n, &qrel{view: true}, sc); err != nil {
				return nil, err
			}
			return relFrom(&qrel{name: exposed, cols: cols}), nil
		}
	}
	e := a.e
	owner := e.currentSchema
	if owner == "" {
		owner = e.owner
	}
	if name == "DUAL" && (schema == "" || schema == "SYS" || schema == "PUBLIC") && e.tables["DUAL"] == nil && e.views["DUAL"] == nil {
		rel := &qrel{name: exposed, cols: dualColumns()}
		if n.Alias == nil {
			rel.schema = "SYS"
		}
		return relFrom(rel), nil
	}
	if schema != "" && schema != owner {
		return nil, oraError("00942", "table or view %q.%q does not exist", schema, name)
	}
	target := e.resolveSynonymTarget(name, schema)
	rel := &qrel{name: exposed}
	if n.Alias == nil {
		rel.schema = owner
	}
	switch {
	case e.views[target] != nil:
		if len(e.views[target].Columns) == 0 {
			// A view whose columns the catalog could not derive: Oracle
			// would read it only once it compiles.
			return nil, oraError("04063", "view %q.%q has errors", owner, target)
		}
		rel.cols = tableColumns(e.views[target].Columns)
		rel.view = true
	case e.materializedViews[target] != nil && e.tables[target] != nil:
		rel.cols = tableColumns(e.tables[target].Columns)
	case e.tables[target] != nil:
		t := e.tables[target]
		rel.cols = tableColumns(t.Columns)
		rel.partitioned = len(t.Partitions) > 0
		rel.table = t
	default:
		if schema == "" {
			schema = owner
		}
		return nil, oraError("00942", "table or view %q.%q does not exist", schema, name)
	}
	if err := a.tableClauses(n, rel, sc); err != nil {
		return nil, err
	}
	return relFrom(rel), nil
}

// tableClauses checks a table reference's partition extension, flashback
// and sample clauses.
func (a *queryAnalyzer) tableClauses(n *ast.TableRef, rel *qrel, sc *qscope) error {
	if n.PartitionExt != nil && !rel.partitioned {
		return oraError("14501", "object is not partitioned")
	}
	if pe := n.PartitionExt; pe != nil && !pe.IsFor && pe.Name != "" && rel.table != nil && !hasPartition(rel.table, pe.Name, pe.IsSubpartition) {
		return oraError("02149", "Specified partition does not exist")
	}
	if smp := n.Sample; smp != nil {
		for _, x := range []ast.ExprNode{smp.Percent, smp.Seed} {
			if x == nil {
				continue
			}
			if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
				return err
			}
		}
	}
	if fb := n.Flashback; fb != nil {
		if fb.IsPeriodFor {
			if !hasPeriod(rel, fb.PeriodColumn) {
				return oraError("55603", "invalid flashback archive or valid time period command")
			}
			// A valid-time bound is an expression like any other.
			for _, x := range []ast.ExprNode{fb.Expr, fb.VersionsLow, fb.VersionsHigh} {
				if x == nil {
					continue
				}
				if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
					return err
				}
			}
		}
		check := func(x ast.ExprNode) error {
			if x == nil {
				return nil
			}
			t, _, err := a.expr(x, sc, clauseFrom)
			if err != nil {
				return err
			}
			switch fb.Type {
			case "SCN":
				if !t.isNumeric() && !t.isCharacter() && typeFamily(t) != "" {
					return oraError("00932", "expression is of data type %s, which is incompatible with expected data type NUMBER", t.kind)
				}
				// An SCN the server's history cannot reach, below its floor.
				if lit, ok := unparen(x).(*ast.NumberLiteral); ok && !lit.IsFloat && lit.Ival < a.e.serverConfig.FlashbackSCNFloor {
					return oraError("08181", "specified number is not a valid system change number")
				}
			case "TIMESTAMP":
				if !t.isDatetime() && !t.isCharacter() && !t.isNull() {
					return oraError("00932", "expression is of data type %s, which is incompatible with expected data type TIMESTAMP", t.kind)
				}
			}
			return nil
		}
		if fb.IsVersions {
			if !fb.IsPeriodFor {
				if err := check(fb.VersionsLow); err != nil {
					return err
				}
				if err := check(fb.VersionsHigh); err != nil {
					return err
				}
			}
		} else if !fb.IsPeriodFor {
			if err := check(fb.Expr); err != nil {
				return err
			}
		}
	}
	return nil
}

// hasPeriod reports whether a table has a valid-time period of that name,
// which the catalog records as the period's check constraint.
func hasPeriod(rel *qrel, period string) bool {
	if rel.table == nil {
		return false
	}
	for _, check := range rel.table.GetCheckConstraints() {
		if name, _, ok := periodForCheckConstraint(rel.table, check); ok && strings.EqualFold(name, period) {
			return true
		}
	}
	return false
}

// join analyzes a join.
func (a *queryAnalyzer) join(n *ast.JoinClause, sc *qscope) (*qfrom, error) {
	left, err := a.fromItem(n.Left, sc)
	if err != nil {
		return nil, err
	}
	var right *qfrom
	// LATERAL, and the table functions that read the rows before them
	// (JSON_TABLE, XMLTABLE, TABLE()), see the join's left input.
	lateral := false
	switch n.Right.(type) {
	case *ast.LateralRef, *ast.JsonTableRef, *ast.XmlTableRef, *ast.TableCollectionExpr:
		lateral = true
	}
	switch {
	case lateral || n.Type == ast.JOIN_CROSS_APPLY || n.Type == ast.JOIN_OUTER_APPLY:
		// The right side reads the left side's columns.
		lat := &qscope{parent: sc.parent, ctes: sc.ctes, from: &qfrom{
			rels:  append(append([]*qrel{}, sc.from.rels...), left.rels...),
			items: append(append([]fromCol{}, sc.from.items...), left.items...),
		}}
		switch r := n.Right.(type) {
		case *ast.SubqueryRef:
			sel, ok := r.Subquery.(*ast.SelectStmt)
			if !ok {
				right = &qfrom{}
				break
			}
			res, err := a.query(sel, lat, sc.ctes, false)
			if err != nil {
				return nil, err
			}
			right = relFrom(&qrel{name: aliasName(r.Alias), cols: res.cols})
		case *ast.LateralRef:
			sel, ok := r.Subquery.(*ast.SelectStmt)
			if !ok {
				right = &qfrom{}
				break
			}
			res, err := a.query(sel, lat, sc.ctes, false)
			if err != nil {
				return nil, err
			}
			right = relFrom(&qrel{name: aliasName(r.Alias), cols: res.cols})
		default:
			right, err = a.fromItem(n.Right, lat)
			if err != nil {
				return nil, err
			}
		}
	default:
		right, err = a.fromItem(n.Right, sc)
		if err != nil {
			return nil, err
		}
	}
	out := &qfrom{rels: append(append([]*qrel{}, left.rels...), right.rels...)}

	var using []string
	natural := false
	switch n.Type {
	case ast.JOIN_NATURAL_INNER, ast.JOIN_NATURAL_LEFT, ast.JOIN_NATURAL_RIGHT, ast.JOIN_NATURAL_FULL:
		natural = true
		// The common visible columns: an invisible column is no NATURAL
		// join key.
		rightNames := map[string]bool{}
		for _, it := range right.items {
			if c := right.itemCol(it); !c.hidden {
				rightNames[c.name] = true
			}
		}
		for _, it := range left.items {
			c := left.itemCol(it)
			if !c.hidden && rightNames[c.name] && !containsString(using, c.name) {
				using = append(using, c.name)
			}
		}
	default:
		using = stringList(n.Using)
	}
	kind := n.Type
	switch kind {
	case ast.JOIN_NATURAL_INNER:
		kind = ast.JOIN_INNER
	case ast.JOIN_NATURAL_LEFT:
		kind = ast.JOIN_LEFT
	case ast.JOIN_NATURAL_RIGHT:
		kind = ast.JOIN_RIGHT
	case ast.JOIN_NATURAL_FULL:
		kind = ast.JOIN_FULL
	}
	switch kind {
	case ast.JOIN_LEFT, ast.JOIN_OUTER_APPLY:
		nullExtend(right)
	case ast.JOIN_RIGHT:
		nullExtend(left)
	case ast.JOIN_FULL:
		nullExtend(left)
		nullExtend(right)
	}
	leftFlat, rightFlat := left.flat, right.flat
	if _, nested := n.Left.(*ast.JoinClause); nested {
		leftFlat = nil
	}
	if _, nested := n.Right.(*ast.JoinClause); nested {
		rightFlat = nil
	}
	switch kind {
	case ast.JOIN_LEFT, ast.JOIN_OUTER_APPLY:
		rightFlat = nil
	case ast.JOIN_RIGHT:
		leftFlat = nil
	}
	out.flat = append(append([]*qrel{}, leftFlat...), rightFlat...)
	if len(using) > 0 && (kind == ast.JOIN_INNER || kind == ast.JOIN_FULL) {
		// A merged column of two items of one name reads neither.
		for _, l := range leftFlat {
			for _, r := range rightFlat {
				if l.name != "" && l.name == r.name {
					return nil, oraError("00918", "%s: column ambiguously specified", using[0])
				}
			}
		}
	}

	if len(using) == 0 {
		out.items = append(append([]fromCol{}, left.items...), right.items...)
	} else {
		sameName := len(left.rels) == 1 && len(right.rels) == 1 && left.rels[0].name != "" && left.rels[0].name == right.rels[0].name
		var merged []fromCol
		taken := map[string]bool{}
		for _, name := range using {
			lc, lfound, lerr := a.lookupIn(left, name)
			rc, rfound, rerr := a.lookupIn(right, name)
			if lerr != nil {
				return nil, lerr
			}
			if rerr != nil {
				return nil, rerr
			}
			if !lfound || !rfound {
				return nil, oraError("00904", "%q: invalid identifier", name)
			}
			// USING and NATURAL join on equality of the two columns.
			if err := comparable(lc.typ, rc.typ, nil); err != nil {
				return nil, err
			}
			var c qcol
			switch kind {
			case ast.JOIN_LEFT, ast.JOIN_OUTER_APPLY:
				c = lc
				if sameName {
					c = rc
				}
			case ast.JOIN_RIGHT:
				c = rc
				if sameName {
					c = lc
				}
			case ast.JOIN_FULL:
				t, err := mergeSetOpTypes(lc.typ, rc.typ)
				if err != nil {
					return nil, err
				}
				if t.kind == "NUMBER" && !sameType(lc.typ, rc.typ) {
					t = numberType()
				} else if t.kind == "NUMBER" {
					t = numberType()
				}
				c = qcol{name: name, typ: t}
			default:
				c = rc
			}
			c.name = name
			merged = append(merged, fromCol{col: c})
			taken[name] = true
			markMerged(left, name)
			markMerged(right, name)
		}
		out.items = append(out.items, merged...)
		for _, side := range []*qfrom{left, right} {
			for _, it := range side.items {
				c := side.itemCol(it)
				if taken[c.name] {
					continue
				}
				out.items = append(out.items, it)
			}
		}
	}
	_ = natural
	if n.On != nil {
		on := &qscope{parent: sc.parent, ctes: sc.ctes, from: &qfrom{
			rels:  append(append([]*qrel{}, sc.from.rels...), out.rels...),
			items: append(append([]fromCol{}, sc.from.items...), out.items...),
		}}
		if err := a.condition(n.On, on, clauseOn); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func containsString(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func nullExtend(f *qfrom) {
	for _, r := range f.rels {
		r.nullExt = true
	}
	for i := range f.items {
		if f.items[i].rel == nil {
			f.items[i].col.notNull = false
		}
	}
}

func markMerged(f *qfrom, name string) {
	for _, r := range f.rels {
		for _, c := range r.cols {
			if c.name == name {
				if r.merged == nil {
					r.merged = map[string]bool{}
				}
				r.merged[name] = true
			}
		}
	}
}

// tableCollection analyzes TABLE(collection) and a table function call.
func (a *queryAnalyzer) tableCollection(n *ast.TableCollectionExpr, sc *qscope) (*qfrom, error) {
	if n.CollectionExpr == nil && n.FunctionCall != nil && !a.startsWithKeyword(n.Loc.Start, "TABLE") {
		// name(args) in FROM: a polymorphic table function, which this
		// catalog holds none of.
		return nil, oraError("00904", "%q: invalid identifier", n.FunctionCall.FuncName.Name)
	}
	x := n.CollectionExpr
	if x == nil {
		x = n.Expr
	}
	elem, ok := collectionElementType(x)
	if !ok {
		if _, _, err := a.expr(x, sc, clauseFrom); err != nil {
			return nil, err
		}
		return nil, oraError("22905", "Cannot access rows from unsupported input to TABLE operator.")
	}
	// The constructor's arguments are expressions like any other.
	if fc, ok := unparen(x).(*ast.FuncCallExpr); ok {
		for _, item := range listItems(fc.Args) {
			if ex, ok := item.(ast.ExprNode); ok {
				if _, _, err := a.expr(ex, sc, clauseFrom); err != nil {
					return nil, err
				}
			}
		}
	}
	return relFrom(&qrel{name: aliasName(n.Alias), cols: []qcol{{name: "COLUMN_VALUE", typ: elem}}}), nil
}

// collectionElementType is the element type of a built-in collection
// constructor: the SYS.ODCI*LIST types.
func collectionElementType(x ast.ExprNode) (qtype, bool) {
	fc, ok := unparen(x).(*ast.FuncCallExpr)
	if !ok || fc.FuncName == nil {
		return qtype{}, false
	}
	if fc.FuncName.Schema != "" && fc.FuncName.Schema != "SYS" {
		return qtype{}, false
	}
	switch fc.FuncName.Name {
	case "ODCINUMBERLIST":
		return numberType(), true
	case "ODCIVARCHAR2LIST":
		return varchar2Bytes(4000), true
	case "ODCIDATELIST":
		return simpleType("DATE"), true
	case "ODCIRAWLIST":
		t := simpleType("RAW")
		t.bytes = 2000
		return t, true
	}
	return qtype{}, false
}

// pivot analyzes a PIVOT clause.
func (a *queryAnalyzer) pivot(n *ast.PivotClause, sc *qscope) (*qfrom, error) {
	src, err := a.fromItem(n.Source, sc)
	if err != nil {
		return nil, err
	}
	inner := &qscope{parent: sc.parent, ctes: sc.ctes, from: src}
	type agg struct {
		alias string
		typ   qtype
	}
	var aggs []agg
	used := map[string]bool{}
	if n.XML {
		return a.pivotXML(n, src, inner)
	}
	for _, item := range listItems(n.InItems) {
		pi, ok := item.(*ast.PivotInItem)
		if !ok {
			continue
		}
		for _, v := range listItems(pi.Values) {
			if x, ok := v.(ast.ExprNode); ok && !isConstantExpr(x) {
				return nil, oraError("56901", "non-constant expression in pivot or unpivot values clause")
			}
		}
	}
	for _, item := range listItems(n.Aggregates) {
		pa, ok := item.(*ast.PivotAggregate)
		if !ok {
			continue
		}
		fc, ok := unparen(pa.Expr).(*ast.FuncCallExpr)
		if !ok || fc.FuncName == nil || !isAggregateName(fc.FuncName.Name) || fc.Over != nil {
			return nil, oraError("56902", "non-aggregate expressions inside the PIVOT clause")
		}
		t, _, err := a.expr(pa.Expr, inner, clauseSelect)
		if err != nil {
			return nil, err
		}
		walkExpr(pa.Expr, func(x ast.Node) bool {
			if ref, ok := x.(*ast.ColumnRef); ok {
				used[ref.Column] = true
			}
			return true
		})
		aggs = append(aggs, agg{alias: pa.Alias, typ: t})
	}
	if len(aggs) == 0 {
		return nil, oraError("56902", "non-aggregate expressions inside the PIVOT clause")
	}
	for _, x := range listItems(n.ForColumns) {
		if ref, ok := x.(*ast.ColumnRef); ok {
			if _, found, err := a.lookupIn(src, ref.Column); err != nil {
				return nil, err
			} else if !found {
				return nil, oraError("00904", "%q: invalid identifier", ref.Column)
			}
			used[ref.Column] = true
		}
	}
	items := listItems(n.InItems)
	if len(items) == 0 {
		return nil, oraError("56902", "non-aggregate expressions inside the PIVOT clause")
	}
	rel := &qrel{name: aliasName(n.Alias)}
	if rel.name == "" && len(src.rels) == 1 {
		rel.name = src.rels[0].name
		rel.schema = src.rels[0].schema
	}
	for _, it := range src.items {
		c := src.itemCol(it)
		if used[c.name] || c.hidden {
			continue
		}
		rel.cols = append(rel.cols, c)
	}
	for _, item := range items {
		pi, ok := item.(*ast.PivotInItem)
		if !ok {
			continue
		}
		var parts []string
		for _, v := range listItems(pi.Values) {
			x, ok := v.(ast.ExprNode)
			if !ok {
				continue
			}
			if !isConstantExpr(x) {
				return nil, oraError("56901", "non-constant expression in pivot or unpivot values clause")
			}
			parts = append(parts, derivedName(a.e.nodeText(x)))
		}
		base := pi.Alias
		if base == "" {
			base = strings.Join(parts, "_")
		}
		for _, ag := range aggs {
			name := base
			if ag.alias != "" {
				name = base + "_" + ag.alias
			}
			rel.cols = append(rel.cols, qcol{name: name, typ: ag.typ})
		}
	}
	return relFrom(rel), nil
}

// isConstantExpr reports a literal, or an expression of literals.
func isConstantExpr(x ast.ExprNode) bool {
	constant := true
	walkExpr(x, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ColumnRef, *ast.BindVariable, *ast.SubqueryExpr, *ast.PseudoColumn:
			constant = false
			return false
		}
		return true
	})
	return constant
}

// unpivot analyzes an UNPIVOT clause.
func (a *queryAnalyzer) unpivot(n *ast.UnpivotClause, sc *qscope) (*qfrom, error) {
	src, err := a.fromItem(n.Source, sc)
	if err != nil {
		return nil, err
	}
	used := map[string]bool{}
	// One type per value column, merged across the IN mappings by position.
	var valueTypes []qtype
	var haveType []bool
	// One label type per FOR column, by position across the mappings.
	var labelBytes []int
	var labelIsNumber []bool
	growLabels := func(n int) {
		for len(labelBytes) < n {
			labelBytes = append(labelBytes, 0)
			labelIsNumber = append(labelIsNumber, false)
		}
	}
	for _, item := range listItems(n.InputMappings) {
		m, ok := item.(*ast.UnpivotInItem)
		if !ok {
			continue
		}
		var names []string
		for pos, x := range listItems(m.InputColumns) {
			ref, ok := x.(*ast.ColumnRef)
			if !ok {
				continue
			}
			c, found, err := a.lookupIn(src, ref.Column)
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, oraError("00904", "%q: invalid identifier", ref.Column)
			}
			used[ref.Column] = true
			names = append(names, ref.Column)
			for len(valueTypes) <= pos {
				valueTypes = append(valueTypes, qtype{})
				haveType = append(haveType, false)
			}
			if !haveType[pos] {
				valueTypes[pos], haveType[pos] = c.typ, true
				continue
			}
			valueTypes[pos], err = mergeSetOpTypes(valueTypes[pos], c.typ)
			if err != nil {
				return nil, err
			}
		}
		if m.Label != nil {
			labels := []ast.Node{m.Label}
			if row, ok := unparen(m.Label).(*ast.FuncCallExpr); ok && (row.FuncName == nil || row.FuncName.Name == "") {
				labels = listItems(row.Args)
			}
			growLabels(len(labels))
			for pos, l := range labels {
				ex, ok := l.(ast.ExprNode)
				if !ok {
					continue
				}
				switch v := unparen(ex).(type) {
				case *ast.StringLiteral:
					labelBytes[pos] = max(labelBytes[pos], len(v.Val))
				case *ast.NumberLiteral:
					labelIsNumber[pos] = true
				}
			}
		} else {
			growLabels(1)
			labelBytes[0] = max(labelBytes[0], len(strings.Join(names, "_")))
		}
	}
	rel := &qrel{name: aliasName(n.Alias)}
	if rel.name == "" && len(src.rels) == 1 {
		rel.name = src.rels[0].name
	}
	for _, it := range src.items {
		c := src.itemCol(it)
		if used[c.name] || c.hidden {
			continue
		}
		c.notNull = false
		rel.cols = append(rel.cols, c)
	}
	for pos, x := range unpivotLabelColumns(n) {
		ref, ok := x.(*ast.ColumnRef)
		if !ok {
			continue
		}
		growLabels(pos + 1)
		labelType := varchar2Bytes(labelBytes[pos])
		if labelIsNumber[pos] {
			labelType = numberType()
		}
		rel.cols = append(rel.cols, qcol{name: ref.Column, typ: labelType})
	}
	for pos, x := range unpivotValueColumns(n) {
		if ref, ok := x.(*ast.ColumnRef); ok {
			t := simpleType("UNKNOWN")
			if pos < len(valueTypes) && haveType[pos] {
				t = valueTypes[pos]
			}
			rel.cols = append(rel.cols, qcol{name: ref.Column, typ: t})
		}
	}
	seen := map[string]bool{}
	for _, c := range rel.cols {
		if seen[c.name] {
			return nil, oraError("00918", "%s: column ambiguously specified", c.name)
		}
		seen[c.name] = true
	}
	return relFrom(rel), nil
}

// matchRecognize analyzes a MATCH_RECOGNIZE clause.
func (a *queryAnalyzer) matchRecognize(n *ast.MatchRecognizeClause, sc *qscope) (*qfrom, error) {
	src, err := a.fromItem(n.Source, sc)
	if err != nil {
		return nil, err
	}
	vars := patternVariables(n.Pattern)
	varSet := map[string]bool{}
	for _, v := range vars {
		varSet[v] = true
	}
	subsets := map[string]bool{}
	for _, item := range listItems(n.Subsets) {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		if varSet[rt.Name] || subsets[rt.Name] {
			return nil, oraError("62503", "illegal variable definition")
		}
		// Its members are distinct primary pattern variables.
		members := map[string]bool{}
		var bad bool
		walkExpr(rt.Expr, func(x ast.Node) bool {
			if ref, ok := x.(*ast.ColumnRef); ok {
				if !varSet[ref.Column] || members[ref.Column] {
					bad = true
				}
				members[ref.Column] = true
			}
			return true
		})
		if bad {
			return nil, oraError("62503", "illegal variable definition")
		}
		subsets[rt.Name] = true
	}
	defined := map[string]bool{}
	for _, item := range listItems(n.Definitions) {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		if defined[rt.Name] || !varSet[rt.Name] || subsets[rt.Name] {
			return nil, oraError("62503", "illegal variable definition")
		}
		defined[rt.Name] = true
	}
	// A pattern variable or subset qualifies the input's columns in
	// MEASURES and DEFINE: a.id reads id of the rows a matched.
	inner := &qscope{parent: sc.parent, ctes: sc.ctes, from: &qfrom{items: src.items, flat: src.flat}}
	var inputCols []qcol
	for _, it := range src.items {
		inputCols = append(inputCols, src.itemCol(it))
	}
	// A variable the pattern names twice is one binding.
	var names []string
	for _, v := range append(append([]string{}, vars...), mapKeys(subsets)...) {
		if !containsString(names, v) {
			names = append(names, v)
		}
	}
	for _, r := range src.rels {
		// A pattern variable shadows an input item of its name.
		if !containsString(names, r.name) {
			inner.from.rels = append(inner.from.rels, r)
		}
	}
	for _, v := range names {
		inner.from.rels = append(inner.from.rels, &qrel{name: v, cols: inputCols})
	}
	for _, item := range listItems(n.Definitions) {
		if rt, ok := item.(*ast.ResTarget); ok && rt.Expr != nil {
			if err := a.condition(rt.Expr, inner, clauseSelect); err != nil {
				return nil, err
			}
		}
	}
	for _, item := range listItems(n.OrderBy) {
		if sb, ok := item.(*ast.SortBy); ok {
			if _, _, err := a.expr(sb.Expr, inner, clauseSelect); err != nil {
				return nil, err
			}
		}
	}
	rel := &qrel{name: aliasName(n.Alias)}
	var partition []qcol
	for _, x := range listItems(n.PartitionBy) {
		ex, ok := x.(ast.ExprNode)
		if !ok {
			continue
		}
		t, nn, err := a.expr(ex, inner, clauseSelect)
		if err != nil {
			return nil, err
		}
		name := derivedName(a.e.nodeText(ex))
		if ref, ok := unparen(ex).(*ast.ColumnRef); ok {
			name = ref.Column
		}
		partition = append(partition, qcol{name: name, typ: t, notNull: nn})
	}
	var measures []qcol
	for _, item := range listItems(n.Measures) {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		t, err := a.measureType(rt.Expr, inner)
		if err != nil {
			return nil, err
		}
		measures = append(measures, qcol{name: rt.Name, typ: t})
	}
	if strings.HasPrefix(strings.ToUpper(n.RowsPerMatch), "ALL ROWS") {
		// Partition columns, ordering columns, measures, then the rest of
		// the input in its order.
		var order []qcol
		for _, item := range listItems(n.OrderBy) {
			if sb, ok := item.(*ast.SortBy); ok {
				if ref, ok := unparen(sb.Expr).(*ast.ColumnRef); ok {
					if c, found, _ := a.lookupIn(src, ref.Column); found {
						c.name = ref.Column
						order = append(order, c)
					}
				}
			}
		}
		taken := map[string]bool{}
		for _, c := range append(append([]qcol{}, partition...), order...) {
			if !taken[c.name] {
				rel.cols = append(rel.cols, c)
				taken[c.name] = true
			}
		}
		rel.cols = append(rel.cols, measures...)
		for _, it := range src.items {
			c := src.itemCol(it)
			if c.hidden || taken[c.name] {
				continue
			}
			rel.cols = append(rel.cols, c)
		}
	} else {
		rel.cols = append(append(rel.cols, partition...), measures...)
	}
	if len(rel.cols) == 0 {
		return nil, oraError("30732", "table contains no user-visible columns")
	}
	seen := map[string]bool{}
	for _, c := range rel.cols {
		if seen[c.name] {
			return nil, oraError("00918", "%s: column ambiguously specified", c.name)
		}
		seen[c.name] = true
	}
	return relFrom(rel), nil
}

// measureType types a MATCH_RECOGNIZE measure: FIRST, LAST, PREV, NEXT and
// CLASSIFIER read the input row's columns.
func (a *queryAnalyzer) measureType(x ast.ExprNode, sc *qscope) (qtype, error) {
	if fc, ok := unparen(x).(*ast.FuncCallExpr); ok && fc.FuncName != nil {
		switch fc.FuncName.Name {
		case "FIRST", "LAST", "PREV", "NEXT":
			// The value is the first argument's; an offset follows it.
			result := simpleType("UNKNOWN")
			for i, item := range listItems(fc.Args) {
				arg, ok := item.(ast.ExprNode)
				if !ok {
					continue
				}
				t, _, err := a.expr(arg, sc, clauseSelect)
				if err != nil {
					return qtype{}, err
				}
				if i == 0 {
					result = t
				}
			}
			return result, nil
		case "CLASSIFIER":
			return varchar2Bytes(128), nil
		case "MATCH_NUMBER":
			return numberType(), nil
		}
	}
	t, _, err := a.expr(x, sc, clauseSelect)
	return t, err
}

// patternVariables lists the variable names a PATTERN's text holds.
func patternVariables(pattern string) []string {
	var out []string
	word := strings.Builder{}
	flush := func() {
		if word.Len() == 0 {
			return
		}
		// The parser writes the pattern normalized: an unquoted variable
		// uppercased, a quoted one as spelled.
		w := word.String()
		word.Reset()
		if w == "PERMUTE" {
			return
		}
		out = append(out, w)
	}
	runes := []rune(pattern)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '"' {
			// A quoted variable keeps its spelling.
			flush()
			j := i + 1
			for j < len(runes) && runes[j] != '"' {
				j++
			}
			out = append(out, string(runes[i+1:j]))
			i = j
			continue
		}
		if r == '_' || r == '$' || r == '#' || unicode.IsLetter(r) || (word.Len() > 0 && unicode.IsDigit(r)) {
			word.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	return out
}

// analyticView checks an inline analytic view's sub_av_clause. The catalog
// holds no analytic views, so the base view never exists.
func (a *queryAnalyzer) analyticView(c *ast.SubAVClause) error {
	if c == nil {
		return nil
	}
	for _, opt := range listItems(c.Options) {
		if o, ok := opt.(*ast.DDLOption); ok && strings.EqualFold(o.Key, "HIERARCHIES") {
			return oraError("18370", "The HIERARCHIES clause may only be specified for cubes.")
		}
	}
	schema, name := "", ""
	if c.Using != nil {
		schema, name = c.Using.Schema, c.Using.Name
	}
	if schema == "" {
		schema = a.e.currentSchema
	}
	return oraError("18307", "analytic view %q.%q does not exist", schema, name)
}

// unpivotLabelColumns lists an UNPIVOT's FOR columns, which a list in
// parentheses holds as a row constructor.
func unpivotLabelColumns(n *ast.UnpivotClause) []ast.Node {
	switch p := n.PivotColumn.(type) {
	case nil:
		return nil
	case *ast.ParenExpr:
		return []ast.Node{p.Expr}
	case *ast.FuncCallExpr:
		if p.FuncName == nil || p.FuncName.Name == "" {
			return listItems(p.Args)
		}
		return []ast.Node{p}
	default:
		return []ast.Node{p}
	}
}

// startsWithKeyword reports whether the schema text at offset begins with
// the keyword.
func (a *queryAnalyzer) startsWithKeyword(offset int, kw string) bool {
	text := a.e.schemaText
	if offset < 0 || offset+len(kw) > len(text) {
		return false
	}
	if !strings.EqualFold(text[offset:offset+len(kw)], kw) {
		return false
	}
	if offset+len(kw) == len(text) {
		return true
	}
	ch := text[offset+len(kw)]
	return !(ch == '_' || ch == '$' || ch == '#' || unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)))
}

// unpivotValueColumns lists an UNPIVOT's value columns, which a list in
// parentheses holds as a row constructor.
func unpivotValueColumns(n *ast.UnpivotClause) []ast.Node {
	var out []ast.Node
	for _, x := range listItems(n.ValueColumns) {
		if row, ok := x.(*ast.FuncCallExpr); ok && (row.FuncName == nil || row.FuncName.Name == "") {
			out = append(out, listItems(row.Args)...)
			continue
		}
		out = append(out, x)
	}
	return out
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// pivotXML analyzes PIVOT XML: the input's columns the pivot does not
// consume, then one XMLTYPE column named for the FOR columns.
func (a *queryAnalyzer) pivotXML(n *ast.PivotClause, src *qfrom, inner *qscope) (*qfrom, error) {
	// The IN list's subquery supplies the pivot values.
	for _, item := range listItems(n.InList) {
		var sub ast.StmtNode
		switch v := item.(type) {
		case *ast.SubqueryExpr:
			sub = v.Subquery
		case *ast.SelectStmt:
			sub = v
		case *ast.ResTarget:
			if sq, ok := v.Expr.(*ast.SubqueryExpr); ok {
				sub = sq.Subquery
			}
		}
		if sel, ok := sub.(*ast.SelectStmt); ok {
			if _, err := a.query(sel, inner.parent, inner.ctes, false); err != nil {
				return nil, err
			}
		}
	}
	used := map[string]bool{}
	for _, item := range listItems(n.Aggregates) {
		pa, ok := item.(*ast.PivotAggregate)
		if !ok {
			continue
		}
		fc, ok := unparen(pa.Expr).(*ast.FuncCallExpr)
		if !ok || fc.FuncName == nil || !isAggregateName(fc.FuncName.Name) || fc.Over != nil {
			return nil, oraError("56902", "non-aggregate expressions inside the PIVOT clause")
		}
		if _, _, err := a.expr(pa.Expr, inner, clauseSelect); err != nil {
			return nil, err
		}
		walkExpr(pa.Expr, func(x ast.Node) bool {
			if ref, ok := x.(*ast.ColumnRef); ok {
				used[ref.Column] = true
			}
			return true
		})
	}
	var forNames []string
	for _, x := range listItems(n.ForColumns) {
		ref, ok := x.(*ast.ColumnRef)
		if !ok {
			continue
		}
		if _, found, err := a.lookupIn(src, ref.Column); err != nil {
			return nil, err
		} else if !found {
			return nil, oraError("00904", "%q: invalid identifier", ref.Column)
		}
		used[ref.Column] = true
		forNames = append(forNames, ref.Column)
	}
	rel := &qrel{name: aliasName(n.Alias)}
	if rel.name == "" && len(src.rels) == 1 {
		rel.name = src.rels[0].name
		rel.schema = src.rels[0].schema
	}
	for _, it := range src.items {
		c := src.itemCol(it)
		if used[c.name] || c.hidden {
			continue
		}
		rel.cols = append(rel.cols, c)
	}
	rel.cols = append(rel.cols, qcol{name: strings.Join(forNames, "_") + "_XML", typ: simpleType("XMLTYPE")})
	return relFrom(rel), nil
}

// hasPartition reports whether a table has a partition, or a subpartition,
// of that name.
func hasPartition(t *metadatapb.TableMetadata, name string, sub bool) bool {
	for _, p := range t.GetPartitions() {
		if !sub && strings.EqualFold(p.GetName(), name) {
			return true
		}
		if sub {
			for _, sp := range p.GetSubpartitions() {
				if strings.EqualFold(sp.GetName(), name) {
					return true
				}
			}
		}
	}
	return false
}

// jsonTablePaths checks JSON_TABLE's row path and its columns' paths, which
// are text literals: a bind is ORA-01027, another expression ORA-40597.
func (a *queryAnalyzer) jsonTablePaths(n *ast.JsonTableRef) error {
	check := func(x ast.ExprNode) error {
		switch unparen(x).(type) {
		case nil, *ast.StringLiteral:
			return nil
		case *ast.BindVariable:
			return oraError("01027", "bind variables not allowed for data definition operations")
		}
		return oraError("40597", "JSON path expression syntax error")
	}
	if err := check(n.Path); err != nil {
		return err
	}
	for _, item := range listItems(n.Columns) {
		c, ok := item.(*ast.JsonTableColumn)
		if !ok {
			continue
		}
		if err := check(c.Path); err != nil {
			return err
		}
		if c.Nested != nil {
			if err := a.jsonTablePaths(c.Nested); err != nil {
				return err
			}
		}
	}
	return nil
}
