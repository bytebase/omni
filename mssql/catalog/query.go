package catalog

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	mssqlparser "github.com/bytebase/omni/mssql/parser"
)

// This file binds the query of a view as SQL Server does when it creates one:
// it resolves the objects and columns the query names, gives the view its
// column names, and keeps the rules a view's query is held to. A query the
// binder refuses is refused with the number of the message SQL Server gives,
// as measured by testdata/queries/simulate.tsv.
//
// What the catalog does not model is not checked: an object in another
// database or on a linked server, a system view, a synonym's target, and a
// relation whose columns the binder cannot work out are opaque, and any
// column of them resolves.

// qtype is the coarse class of an expression's type, as far as the checks
// that depend on a type need it.
type qtype int

const (
	qtUnknown qtype = iota
	qtNumber
	qtString
	qtDate
	qtBinary
	qtXML
	qtOther
)

// typeName names the class in a message, as SQL Server names a type there.
func (t qtype) typeName() string {
	switch t {
	case qtNumber:
		return "int"
	case qtString:
		return "varchar"
	case qtDate:
		return "datetime"
	case qtBinary:
		return "varbinary"
	case qtXML:
		return "xml"
	default:
		return "sql_variant"
	}
}

// qcol is a column of a relation a query reads, or of a query's result.
type qcol struct {
	name     string
	typ      qtype
	identity bool
	guid     bool // uniqueidentifier, the type a ROWGUIDCOL column has
	hidden   bool // HIDDEN: named, but not in *
	// noInsert is a column an INSERT without a column list does not fill:
	// computed, rowversion, or a system-versioned table's period column.
	noInsert bool
}

// qrel is a relation in a FROM clause: a table, a view, a derived table, a
// common table expression, or the result of a table-valued function.
type qrel struct {
	name   string // the exposed name: the alias, or the object's name
	schema string // the schema of an unaliased object, "" otherwise
	alias  bool   // the exposed name is a correlation name
	cols   []qcol
	opaque bool // the columns are unknown; any column resolves
	// starOpaque is set when * may leave out columns the catalog cannot
	// tell are HIDDEN.
	starOpaque bool
	// table is the table a relation reads, nil for any other relation.
	table *metadata.TableMetadata
}

func (r *qrel) col(name string) *qcol {
	for i := range r.cols {
		if strings.EqualFold(r.cols[i].name, name) {
			return &r.cols[i]
		}
	}
	return nil
}

// rowguidCol is the column $ROWGUID may name. The catalog does not record
// ROWGUIDCOL, so it is the relation's first uniqueidentifier column, the only
// type the property takes; with none, $ROWGUID names nothing.
func (r *qrel) rowguidCol() *qcol {
	for i := range r.cols {
		if r.cols[i].guid {
			return &r.cols[i]
		}
	}
	return nil
}

func (r *qrel) identityCol() *qcol {
	for i := range r.cols {
		if r.cols[i].identity {
			return &r.cols[i]
		}
	}
	return nil
}

// qcte is a common table expression a query can name.
type qcte struct {
	name    string
	rel     *qrel // nil while its query is being bound and no anchor is known
	binding bool
}

// qscope is the names one query block sees: its own relations, then its
// enclosing blocks'.
type qscope struct {
	parent *qscope
	rels   []*qrel
	ctes   []*qcte
	block  *qblock
}

// qblock is one query specification (SELECT ... FROM ... GROUP BY ...).
type qblock struct {
	hasAgg  bool
	windows map[string]*ast.WindowDef
}

// qclause is the clause an expression is bound in.
type qclause int

const (
	clOther qclause = iota
	clSelect
	clWhere
	clOn
	clGroup
	clHaving
	clOrder
)

type ectx struct {
	clause   qclause
	aggDepth int
	block    *qblock
}

type colBinding struct {
	rel   *qrel
	block *qblock
}

// binder binds one object's query.
type binder struct {
	c           *Catalog
	source      string
	object      string          // the view or function being created, for messages
	vars        map[string]bool // declared variables, upper-cased
	inView      bool
	schemaBound bool
	bound       map[*ast.ColumnRef]colBinding
	depth       int
	// setOrder holds the last operand of each set operation, which carries
	// the ORDER BY of the whole operation.
	setOrder map[*ast.SelectStmt]bool
	// hoisted holds the first operands whose WITH a set operation bound.
	hoisted map[*ast.SelectStmt]bool
	// tableVars holds the columns of the table variables a module declares,
	// upper-cased; nil columns are unknown.
	tableVars map[string][]qcol
	// pseudo holds the tables a trigger's body reads by name: inserted and
	// deleted, upper-cased.
	pseudo map[string]*qrel
	// deferNames is set in a module's body, where SQL Server resolves a
	// user-defined function only when the module runs.
	deferNames bool
	// clauses holds the clause each enclosing subquery is in, innermost
	// last.
	clauses []qclause
	// existsBody holds the queries EXISTS tests, whose select list SQL
	// Server does not hold to their GROUP BY.
	existsBody map[*ast.SelectStmt]bool
}

func newBinder(c *Catalog, source, object string, depth int) *binder {
	return &binder{c: c, source: source, object: object, vars: map[string]bool{}, bound: map[*ast.ColumnRef]colBinding{}, depth: depth, setOrder: map[*ast.SelectStmt]bool{}, hoisted: map[*ast.SelectStmt]bool{}, existsBody: map[*ast.SelectStmt]bool{}}
}

// maxBindDepth bounds how deep views and functions that name each other are
// bound; past it a relation is opaque.
const maxBindDepth = 32

// bindView binds the query of CREATE VIEW or ALTER VIEW and returns the
// view's columns, nil when they are unknown.
func (c *Catalog) bindView(n *ast.CreateViewStmt, source, viewName string) ([]qcol, error) {
	sel, ok := n.Query.(*ast.SelectStmt)
	if !ok || sel == nil {
		return nil, nil
	}
	b := newBinder(c, source, viewName, 0)
	b.inView = true
	b.schemaBound = n.SchemaBinding || hasOptionWord(n.Options, "SCHEMABINDING")
	if err := b.viewClauses(sel); err != nil {
		return nil, err
	}
	if err := b.variables(sel); err != nil {
		return nil, err
	}
	cols, opaque, err := b.query(sel, nil)
	if err != nil || opaque {
		return nil, err
	}
	if n.Columns != nil && len(n.Columns.Items) > 0 {
		return renameColumns(cols, stringNodeList(n.Columns), viewName)
	}
	for i, col := range cols {
		if col.name == "" {
			return nil, mssqlErr(4511, "Create View or Function failed because no column name was specified for column %d.", i+1)
		}
		for _, prev := range cols[:i] {
			if strings.EqualFold(prev.name, col.name) {
				return nil, mssqlErr(4506, "Column names in each view or function must be unique. Column name '%s' in view or function '%s' is specified more than once.", col.name, viewName)
			}
		}
	}
	return cols, nil
}

func hasOptionWord(list *ast.List, word string) bool {
	if list == nil {
		return false
	}
	for _, item := range list.Items {
		if s, ok := item.(*ast.String); ok && strings.EqualFold(s.Str, word) {
			return true
		}
	}
	return false
}

// variables refuses a variable the object does not declare. SQL Server's
// parser declares variables, so this comes before any other check of the
// query.
func (b *binder) variables(sel *ast.SelectStmt) error {
	var err error
	ast.Inspect(sel, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch v := n.(type) {
		case *ast.VariableRef:
			if !strings.HasPrefix(v.Name, "@@") && !b.vars[strings.ToUpper(v.Name)] {
				err = mssqlErr(137, "Must declare the scalar variable \"%s\".", v.Name)
			}
		case *ast.SelectAssign:
			if !b.vars[strings.ToUpper(v.Variable)] {
				err = mssqlErr(137, "Must declare the scalar variable \"%s\".", v.Variable)
			}
		case *ast.TableVarRef:
			if !b.vars[strings.ToUpper(v.Name)] {
				err = mssqlErr(1087, "Must declare the table variable \"%s\".", v.Name)
			}
		case *ast.TableVarMethodCallRef:
			if !b.vars[strings.ToUpper(v.Var)] {
				err = mssqlErr(137, "Must declare the scalar variable \"%s\".", v.Var)
			}
		}
		return err == nil
	})
	return err
}

// firstOperand is the first query specification of a set operation.
func firstOperand(sel *ast.SelectStmt) *ast.SelectStmt {
	for sel != nil && sel.Op != ast.SetOpNone && sel.Larg != nil {
		sel = sel.Larg
	}
	return sel
}

// lastOperand is the last query specification of a set operation, which
// carries the ORDER BY, OFFSET and FOR of the whole operation.
func lastOperand(sel *ast.SelectStmt) *ast.SelectStmt {
	for sel != nil && sel.Op != ast.SetOpNone && sel.Rarg != nil {
		sel = sel.Rarg
	}
	return sel
}

// viewClauses refuses the clauses a view's select_statement does not have.
func (b *binder) viewClauses(sel *ast.SelectStmt) error {
	if sel.IntoTable != nil {
		return mssqlErr(156, "Incorrect syntax near the keyword 'INTO'.")
	}
	if sel.OptionClause != nil {
		return mssqlErr(156, "Incorrect syntax near the keyword 'OPTION'.")
	}
	if sel.ForClause != nil && sel.ForClause.Mode == ast.ForBrowse {
		return mssqlErr(176, "The FOR BROWSE clause is no longer supported in views.")
	}
	last := lastOperand(sel)
	if sel.Op != ast.SetOpNone && last != nil && last.OrderByClause != nil && len(last.OrderByClause.Items) > 0 && last.OffsetClause == nil && last.ForClause == nil {
		return mssqlErr(1033, "The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.")
	}
	if sel.OrderByClause != nil && len(sel.OrderByClause.Items) > 0 && sel.Top == nil && sel.OffsetClause == nil && sel.ForClause == nil {
		return mssqlErr(1033, "The ORDER BY clause is invalid in views, inline functions, derived tables, subqueries, and common table expressions, unless TOP, OFFSET or FOR XML is also specified.")
	}
	return nil
}

// query binds a query expression and returns its result columns. opaque
// reports a result whose columns are unknown.
func (b *binder) query(sel *ast.SelectStmt, parent *qscope) ([]qcol, bool, error) {
	if sel == nil {
		return nil, true, nil
	}
	scope := parent
	// The parser hangs a set operation's WITH on its first operand, but its
	// CTEs belong to the whole statement, so the outermost set operation
	// binds them for every operand.
	if sel.Op != ast.SetOpNone {
		if first := firstOperand(sel); first != nil && first.WithClause != nil && !b.hoisted[first] {
			b.hoisted[first] = true
			scope = &qscope{parent: parent}
			if parent != nil {
				scope.block = parent.block
			}
			if err := b.ctes(first.WithClause, scope); err != nil {
				return nil, false, err
			}
		}
	}
	if sel.WithClause != nil && !b.hoisted[sel] {
		scope = &qscope{parent: parent}
		if parent != nil {
			scope.block = parent.block
		}
		if err := b.ctes(sel.WithClause, scope); err != nil {
			return nil, false, err
		}
	}
	if sel.Op == ast.SetOpNone {
		return b.block(sel, scope)
	}
	// The last operand's ORDER BY belongs to the whole operation, so only
	// the outermost set operation that reaches it binds it.
	last := lastOperand(sel)
	owner := !b.setOrder[last]
	b.setOrder[last] = true
	left, lopaque, err := b.query(sel.Larg, scope)
	if err != nil {
		return nil, false, err
	}
	right, ropaque, err := b.query(sel.Rarg, scope)
	if err != nil {
		return nil, false, err
	}
	if !lopaque && !ropaque && len(left) != len(right) {
		return nil, false, mssqlErr(205, "All queries combined using a UNION, INTERSECT or EXCEPT operator must have an equal number of expressions in their target lists.")
	}
	if err := b.topAndOffset(sel, scope); err != nil {
		return nil, false, err
	}
	// A result column's type follows every operand's, by type precedence;
	// the classes keep only what the operands agree on.
	if !lopaque && !ropaque {
		merged := append([]qcol(nil), left...)
		for i := range merged {
			if merged[i].typ != right[i].typ {
				merged[i].typ = qtUnknown
			}
		}
		left = merged
	}
	if owner && last != nil && last.OrderByClause != nil && !lopaque {
		for _, item := range last.OrderByClause.Items {
			ob, ok := item.(*ast.OrderByItem)
			if !ok {
				continue
			}
			if err := b.setOrderItem(ob.Expr, left); err != nil {
				return nil, false, err
			}
		}
	}
	return left, lopaque || ropaque, nil
}

// setOrderItem resolves an ORDER BY item of a set operation, which names a
// result column or its position.
func (b *binder) setOrderItem(e ast.ExprNode, cols []qcol) error {
	if n, ok := ordinal(e); ok {
		if n < 1 || n > len(cols) {
			return mssqlErr(108, "The ORDER BY position number %d is out of range of the number of items in the select list.", n)
		}
		return nil
	}
	if ref, ok := e.(*ast.ColumnRef); ok {
		for _, c := range cols {
			if strings.EqualFold(c.name, ref.Column) {
				return nil
			}
		}
	}
	return mssqlErr(104, "ORDER BY items must appear in the select list if the statement contains a UNION, INTERSECT or EXCEPT operator.")
}

func ordinal(e ast.ExprNode) (int, bool) {
	if lit, ok := e.(*ast.Literal); ok && lit.Type == ast.LitInteger {
		return int(lit.Ival), true
	}
	return 0, false
}

// ctes binds a WITH list into scope.
func (b *binder) ctes(wc *ast.WithClause, scope *qscope) error {
	if wc.CTEs == nil {
		return nil
	}
	for _, item := range wc.CTEs.Items {
		def, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		for _, prev := range scope.ctes {
			if strings.EqualFold(prev.name, def.Name) {
				return mssqlErr(239, "Duplicate common table expression name '%s' was specified.", def.Name)
			}
		}
		q := &qcte{name: def.Name, binding: true}
		scope.ctes = append(scope.ctes, q)
		aliases := stringNodeList(def.Columns)
		sel, _ := def.Query.(*ast.SelectStmt)
		var cols []qcol
		var opaque bool
		var err error
		if sel != nil && referencesName(sel, def.Name) {
			cols, opaque, err = b.recursiveCTE(q, sel, aliases, scope)
		} else {
			cols, opaque, err = b.query(sel, scope)
		}
		if err != nil {
			return err
		}
		q.binding = false
		if opaque {
			q.rel = &qrel{name: def.Name, alias: true, opaque: true}
			continue
		}
		named, err := b.derivedColumns(cols, aliases, def.Name)
		if err != nil {
			return err
		}
		q.rel = &qrel{name: def.Name, alias: true, cols: named}
	}
	return nil
}

// recursiveCTE binds a common table expression that names itself: its anchor
// first, which gives the CTE its columns, and then the recursive members.
func (b *binder) recursiveCTE(q *qcte, sel *ast.SelectStmt, aliases []string, scope *qscope) ([]qcol, bool, error) {
	if sel.Op == ast.SetOpNone || sel.Larg == nil || referencesName(sel.Larg, q.name) {
		return nil, false, mssqlErr(252, "Recursive common table expression '%s' does not contain a top-level UNION ALL operator.", q.name)
	}
	anchor, opaque, err := b.query(sel.Larg, scope)
	if err != nil {
		return nil, false, err
	}
	if opaque {
		q.rel = &qrel{name: q.name, alias: true, opaque: true}
	} else {
		named, err := b.derivedColumns(anchor, aliases, q.name)
		if err != nil {
			return nil, false, err
		}
		q.rel = &qrel{name: q.name, alias: true, cols: named}
	}
	right, ropaque, err := b.query(sel.Rarg, scope)
	if err != nil {
		return nil, false, err
	}
	if !opaque && !ropaque && len(anchor) != len(right) {
		return nil, false, mssqlErr(205, "All queries combined using a UNION, INTERSECT or EXCEPT operator must have an equal number of expressions in their target lists.")
	}
	return anchor, opaque, nil
}

// referencesName reports whether a query names an unqualified table name.
func referencesName(n ast.Node, name string) bool {
	found := false
	ast.Inspect(n, func(node ast.Node) bool {
		if found {
			return false
		}
		if ref, ok := node.(*ast.TableRef); ok && ref.Schema == "" && ref.Database == "" && strings.EqualFold(ref.Object, name) {
			found = true
		}
		return !found
	})
	return found
}

// derivedColumns names the columns of a derived table or common table
// expression: by its column list, or by the names its query gives them.
func (b *binder) derivedColumns(cols []qcol, aliases []string, name string) ([]qcol, error) {
	if len(aliases) > 0 {
		return renameColumns(cols, aliases, name)
	}
	for i, col := range cols {
		if col.name == "" {
			return nil, mssqlErr(8155, "No column name was specified for column %d of '%s'.", i+1, name)
		}
		for _, prev := range cols[:i] {
			if strings.EqualFold(prev.name, col.name) {
				return nil, mssqlErr(8156, "The column '%s' was specified multiple times for '%s'.", col.name, name)
			}
		}
	}
	return cols, nil
}

func renameColumns(cols []qcol, aliases []string, name string) ([]qcol, error) {
	if len(aliases) < len(cols) {
		return nil, mssqlErr(8158, "'%s' has more columns than were specified in the column list.", name)
	}
	if len(aliases) > len(cols) {
		return nil, mssqlErr(8159, "'%s' has fewer columns than were specified in the column list.", name)
	}
	out := make([]qcol, len(cols))
	for i, alias := range aliases {
		for _, prev := range aliases[:i] {
			if strings.EqualFold(prev, alias) {
				return nil, mssqlErr(8156, "The column '%s' was specified multiple times for '%s'.", alias, name)
			}
		}
		out[i] = cols[i]
		out[i].name = alias
	}
	return out, nil
}

func (b *binder) topAndOffset(sel *ast.SelectStmt, sc *qscope) error {
	if sel.Top != nil && sel.OffsetClause != nil {
		return mssqlErr(10741, "A TOP can not be used in the same query or sub-query as a OFFSET.")
	}
	x := ectx{clause: clOther}
	if sel.Top != nil {
		if _, err := b.expr(sel.Top.Count, sc, x); err != nil {
			return err
		}
	}
	if sel.OffsetClause != nil {
		if _, err := b.expr(sel.OffsetClause, sc, x); err != nil {
			return err
		}
	}
	if sel.FetchClause != nil {
		if _, err := b.expr(sel.FetchClause.Count, sc, x); err != nil {
			return err
		}
	}
	return nil
}

// block binds one query specification.
func (b *binder) block(sel *ast.SelectStmt, parent *qscope) ([]qcol, bool, error) {
	blk := &qblock{windows: map[string]*ast.WindowDef{}}
	sc := &qscope{parent: parent, block: blk}
	if err := b.topAndOffset(sel, parent); err != nil {
		return nil, false, err
	}
	if sel.FromClause != nil {
		for _, item := range sel.FromClause.Items {
			rels, err := b.tableExpr(item, sc, nil)
			if err != nil {
				return nil, false, err
			}
			if err := checkExposedNames(sc.rels, rels); err != nil {
				return nil, false, err
			}
			sc.rels = append(sc.rels, rels...)
		}
	}
	if err := b.windowClause(sel, sc, blk); err != nil {
		return nil, false, err
	}
	if _, err := b.expr(sel.WhereClause, sc, ectx{clause: clWhere, block: blk}); err != nil {
		return nil, false, err
	}
	groups, err := b.groupBy(sel, sc, blk)
	if err != nil {
		return nil, false, err
	}
	if _, err := b.expr(sel.HavingClause, sc, ectx{clause: clHaving, block: blk}); err != nil {
		return nil, false, err
	}
	if sel.HavingClause != nil {
		blk.hasAgg = true
	}
	cols, opaque, err := b.selectList(sel, sc, blk)
	if err != nil {
		return nil, false, err
	}
	if !b.setOrder[sel] {
		if err := b.orderBy(sel, sc, blk, cols, opaque); err != nil {
			return nil, false, err
		}
	}
	if blk.hasAgg || groups != nil {
		if err := b.checkGrouping(sel, sc, blk, groups); err != nil {
			return nil, false, err
		}
	}
	if sel.ForClause != nil && sel.ForClause.Mode != ast.ForBrowse {
		// FOR XML and FOR JSON give one column, which has no name.
		return []qcol{{typ: qtString}}, false, nil
	}
	return cols, opaque, nil
}

func checkExposedNames(have, add []*qrel) error {
	for i, r := range add {
		for _, prev := range append(append([]*qrel(nil), have...), add[:i]...) {
			if !strings.EqualFold(prev.name, r.name) {
				continue
			}
			if !prev.alias && !r.alias {
				return mssqlErr(1013, "The objects \"%s\" and \"%s\" in the FROM clause have the same exposed names. Use correlation names to distinguish them.", r.name, prev.name)
			}
			return mssqlErr(1011, "The correlation name '%s' is specified multiple times in a FROM clause.", r.name)
		}
	}
	return nil
}

// windowClause checks the WINDOW clause's definitions name windows that
// exist.
func (b *binder) windowClause(sel *ast.SelectStmt, sc *qscope, blk *qblock) error {
	if sel.WindowClause == nil {
		return nil
	}
	var defs []*ast.WindowDef
	for _, item := range sel.WindowClause.Items {
		if def, ok := item.(*ast.WindowDef); ok {
			defs = append(defs, def)
			blk.windows[strings.ToUpper(def.Name)] = def
		}
	}
	for _, def := range defs {
		// A definition may name a window defined before or after it, but
		// not one that leads back to it.
		seen := map[string]bool{strings.ToUpper(def.Name): true}
		for ref := def.RefName; ref != ""; {
			next, ok := blk.windows[strings.ToUpper(ref)]
			if !ok || seen[strings.ToUpper(ref)] {
				return mssqlErr(5362, "Window '%s' is not defined, or its definition refers to itself.", ref)
			}
			seen[strings.ToUpper(ref)] = true
			ref = next.RefName
		}
		x := ectx{clause: clOrder, block: blk}
		if _, err := b.expr(def.PartitionBy, sc, x); err != nil {
			return err
		}
		if _, err := b.expr(def.OrderBy, sc, x); err != nil {
			return err
		}
	}
	return nil
}

// groupItem is one expression of a GROUP BY clause.
type groupItem struct {
	text string
	node ast.Node
}

// groupBy binds GROUP BY and returns its expressions, nil when there is none.
func (b *binder) groupBy(sel *ast.SelectStmt, sc *qscope, blk *qblock) ([]groupItem, error) {
	if sel.GroupByClause == nil || len(sel.GroupByClause.Items) == 0 {
		return nil, nil
	}
	groups := []groupItem{}
	var visit func(n ast.Node) error
	visit = func(n ast.Node) error {
		switch g := n.(type) {
		case *ast.GroupingSetsExpr:
			return visitList(g.Sets, visit)
		case *ast.RollupExpr:
			return visitList(g.Args, visit)
		case *ast.CubeExpr:
			return visitList(g.Args, visit)
		case *ast.List:
			return visitList(g, visit)
		case *ast.ParenExpr:
			if g.Expr == nil {
				return nil // the grand total ( )
			}
		}
		if _, err := b.expr(n, sc, ectx{clause: clGroup, block: blk}); err != nil {
			return err
		}
		if !b.hasLocalColumn(n, blk) {
			return mssqlErr(164, "Each GROUP BY expression must contain at least one column that is not an outer reference.")
		}
		groups = append(groups, groupItem{text: b.normText(n), node: n})
		return nil
	}
	if err := visitList(sel.GroupByClause, visit); err != nil {
		return nil, err
	}
	blk.hasAgg = true
	return groups, nil
}

func visitList(list *ast.List, f func(ast.Node) error) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

// hasLocalColumn reports whether n names a column of the block's own FROM
// clause.
func (b *binder) hasLocalColumn(n ast.Node, blk *qblock) bool {
	found := false
	ast.Inspect(n, func(node ast.Node) bool {
		if found {
			return false
		}
		if ref, ok := node.(*ast.ColumnRef); ok {
			if bound, ok := b.bound[ref]; ok && bound.block == blk {
				found = true
			}
		}
		return !found
	})
	return found
}

// normText is n's text with case and spacing folded and each bound column
// written as the relation and column it binds to, to compare an expression
// with a GROUP BY expression: t.a + 1 and a + 1 are the same expression.
func (b *binder) normText(n ast.Node) string {
	n = unparen(n)
	loc := ast.NodeLoc(n)
	if loc.Start < 0 || loc.End > len(b.source) || loc.Start > loc.End {
		return ""
	}
	type span struct {
		start, end int
		text       string
	}
	var spans []span
	ast.Inspect(n, func(node ast.Node) bool {
		ref, ok := node.(*ast.ColumnRef)
		if !ok {
			return true
		}
		if bound, ok := b.bound[ref]; ok && bound.rel != nil && ref.Loc.Start >= loc.Start && ref.Loc.End <= loc.End {
			spans = append(spans, span{ref.Loc.Start, ref.Loc.End, fmt.Sprintf("%p.%s", bound.rel, strings.ToLower(ref.Column))})
		}
		return false
	})
	// Compare tokens, not characters: a keyword or name folds its case and
	// loses its brackets, a literal keeps its text, and a bound column is
	// written as what it binds to.
	var out []string
	next := 0
	for _, tok := range mssqlparser.Tokenize(b.source[loc.Start:loc.End]) {
		at := loc.Start + tok.Loc
		for next < len(spans) && spans[next].end <= at {
			next++
		}
		if next < len(spans) && at >= spans[next].start {
			if at == spans[next].start || len(out) == 0 || out[len(out)-1] != spans[next].text {
				out = append(out, spans[next].text)
			}
			continue
		}
		switch tok.Type {
		case mssqlparser.SCONST, mssqlparser.NSCONST:
			out = append(out, b.source[at:loc.Start+tok.End])
		default:
			text := tok.Str
			if text == "" {
				text = b.source[at : loc.Start+tok.End]
			}
			out = append(out, strings.ToLower(text))
		}
	}
	return strings.Join(out, " ")
}

// selectList binds the select list and returns the result's columns.
func (b *binder) selectList(sel *ast.SelectStmt, sc *qscope, blk *qblock) ([]qcol, bool, error) {
	var cols []qcol
	opaque := false
	if sel.TargetList == nil {
		return nil, false, nil
	}
	x := ectx{clause: clSelect, block: blk}
	for _, item := range sel.TargetList.Items {
		switch t := item.(type) {
		case *ast.ResTarget:
			if star, ok := t.Val.(*ast.StarExpr); ok {
				expanded, o, err := b.star(star, sc)
				if err != nil {
					return nil, false, err
				}
				opaque = opaque || o
				cols = append(cols, expanded...)
				continue
			}
			typ, err := b.expr(t.Val, sc, x)
			if err != nil {
				return nil, false, err
			}
			name := t.Name
			if name == "" {
				name = b.resultName(t.Val, sc)
			}
			cols = append(cols, qcol{name: name, typ: typ})
		case *ast.StarExpr:
			expanded, o, err := b.star(t, sc)
			if err != nil {
				return nil, false, err
			}
			opaque = opaque || o
			cols = append(cols, expanded...)
		case *ast.SelectAssign:
			if !b.vars[strings.ToUpper(t.Variable)] {
				return nil, false, mssqlErr(137, "Must declare the scalar variable \"%s\".", t.Variable)
			}
			if _, err := b.expr(t.Value, sc, x); err != nil {
				return nil, false, err
			}
		default:
			if _, err := b.expr(item, sc, x); err != nil {
				return nil, false, err
			}
			cols = append(cols, qcol{})
		}
	}
	return cols, opaque, nil
}

// pseudoColumn reports whether ref is $IDENTITY, IDENTITYCOL, $ROWGUID or
// ROWGUIDCOL as a keyword, not a column delimited to bear that name:
// [IDENTITYCOL] is an ordinary column.
func (b *binder) pseudoColumn(ref *ast.ColumnRef) bool {
	switch strings.ToUpper(ref.Column) {
	case "$IDENTITY", "IDENTITYCOL", "$ROWGUID", "ROWGUIDCOL":
	default:
		return false
	}
	if ref.Loc.Start < 0 || ref.Loc.End > len(b.source) || ref.Loc.Start >= ref.Loc.End {
		return true
	}
	text := strings.TrimSpace(b.source[ref.Loc.Start:ref.Loc.End])
	return !strings.HasSuffix(text, "]") && !strings.HasSuffix(text, "\"")
}

// resultName is the name a select-list expression without an alias gives
// its column: a column's name, or none.
func (b *binder) resultName(e ast.ExprNode, sc *qscope) string {
	switch v := e.(type) {
	case *ast.ParenExpr:
		return b.resultName(v.Expr, sc)
	case *ast.ColumnRef:
		switch strings.ToUpper(v.Column) {
		case "$IDENTITY", "IDENTITYCOL":
			if !b.pseudoColumn(v) {
				return v.Column
			}
			if bound, ok := b.bound[v]; ok && bound.rel != nil {
				if c := bound.rel.identityCol(); c != nil {
					return c.name
				}
			}
			return ""
		}
		return v.Column
	}
	return ""
}

// star expands * or qualifier.*.
func (b *binder) star(star *ast.StarExpr, sc *qscope) ([]qcol, bool, error) {
	if b.schemaBound {
		return nil, false, mssqlErr(1054, "Syntax '*' is not allowed in schema-bound objects.")
	}
	if star.Qualifier == "" {
		if len(sc.rels) == 0 {
			return nil, false, mssqlErr(263, "Must specify table to select from.")
		}
		var cols []qcol
		opaque := false
		for _, r := range sc.rels {
			opaque = opaque || r.opaque || r.starOpaque
			cols = append(cols, visible(r.cols)...)
		}
		return cols, opaque, nil
	}
	// The qualifier is one identifier, the table's or its alias, and a
	// delimited one may hold dots: [a.b].*.
	// The nearest query that has the relation wins: a subquery's t.* may
	// name an enclosing query's t.
	for s := sc; s != nil; s = s.parent {
		for _, r := range s.rels {
			if relMatches(r, "", star.Qualifier) {
				return visible(r.cols), r.opaque || r.starOpaque, nil
			}
		}
	}
	return nil, false, mssqlErr(107, "The column prefix '%s' does not match with a table name or alias name used in the query.", star.Qualifier)
}

// visible are the columns * names: all but the HIDDEN ones.
func visible(cols []qcol) []qcol {
	out := make([]qcol, 0, len(cols))
	for _, c := range cols {
		if !c.hidden {
			out = append(out, c)
		}
	}
	return out
}

func relMatches(r *qrel, schema, table string) bool {
	if !strings.EqualFold(r.name, table) {
		return false
	}
	return schema == "" || (!r.alias && strings.EqualFold(r.schema, schema))
}

// orderBy binds the ORDER BY of a query specification.
func (b *binder) orderBy(sel *ast.SelectStmt, sc *qscope, blk *qblock, cols []qcol, opaque bool) error {
	if sel.OrderByClause == nil {
		return nil
	}
	x := ectx{clause: clOrder, block: blk}
	for _, item := range sel.OrderByClause.Items {
		ob, ok := item.(*ast.OrderByItem)
		if !ok {
			continue
		}
		if n, ok := ordinal(ob.Expr); ok {
			if !opaque && (n < 1 || n > len(cols)) {
				return mssqlErr(108, "The ORDER BY position number %d is out of range of the number of items in the select list.", n)
			}
			continue
		}
		if ref, ok := ob.Expr.(*ast.ColumnRef); ok && ref.Table == "" && namesResult(sel, ref.Column) {
			continue
		}
		if _, err := b.expr(ob.Expr, sc, x); err != nil {
			return err
		}
		if sel.Distinct && !opaque && !b.inSelectList(sel, sc, ob.Expr) {
			return mssqlErr(145, "ORDER BY items must appear in the select list if SELECT DISTINCT is specified.")
		}
	}
	return nil
}

// namesResult reports whether name is the alias of a select-list item.
func namesResult(sel *ast.SelectStmt, name string) bool {
	if sel.TargetList == nil {
		return false
	}
	for _, item := range sel.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok && rt.Name != "" && strings.EqualFold(rt.Name, name) {
			return true
		}
	}
	return false
}

// inSelectList reports whether an ORDER BY item of SELECT DISTINCT is in the
// select list: the same expression, or the same column, which a * names too.
func (b *binder) inSelectList(sel *ast.SelectStmt, sc *qscope, e ast.ExprNode) bool {
	want := b.normText(e)
	ref, isRef := unparen(e).(*ast.ColumnRef)
	var bound colBinding
	if isRef {
		bound, isRef = b.bound[ref]
	}
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Val == nil {
			continue
		}
		if b.normText(rt.Val) == want {
			return true
		}
		if !isRef || bound.rel == nil {
			continue
		}
		switch v := unparen(rt.Val).(type) {
		case *ast.StarExpr:
			for _, r := range sc.rels {
				if r != bound.rel {
					continue
				}
				if v.Qualifier == "" {
					return true
				}
				if relMatches(r, "", v.Qualifier) {
					return true
				}
			}
		case *ast.ColumnRef:
			if vb, ok := b.bound[v]; ok && vb.rel == bound.rel && strings.EqualFold(v.Column, ref.Column) {
				return true
			}
		}
	}
	return false
}

// checkGrouping holds the select list, HAVING and ORDER BY of an aggregated
// query to its GROUP BY.
func (b *binder) checkGrouping(sel *ast.SelectStmt, sc *qscope, blk *qblock, groups []groupItem) error {
	texts := map[string]bool{}
	for _, g := range groups {
		texts[g.text] = true
	}
	grouped := func(ref *ast.ColumnRef) bool {
		bound, ok := b.bound[ref]
		if !ok || bound.block != blk {
			return true
		}
		for _, g := range groups {
			if gref, ok := unparen(g.node).(*ast.ColumnRef); ok {
				if gb, ok := b.bound[gref]; ok && gb.rel == bound.rel && strings.EqualFold(gref.Column, ref.Column) {
					return true
				}
			}
		}
		return texts[b.normText(ref)]
	}
	var firstBad *ast.ColumnRef
	var check func(n ast.Node) bool
	check = func(n ast.Node) bool {
		if firstBad != nil {
			return false
		}
		if n == nil {
			return false
		}
		if texts[b.normText(n)] && b.normText(n) != "" {
			return false
		}
		switch v := n.(type) {
		case *ast.FuncCallExpr:
			if isAggregateCall(v) {
				return false
			}
		case *ast.ColumnRef:
			// A subquery's own columns bind to its blocks and pass; a
			// correlated reference to this block is held to its grouping.
			if !grouped(v) {
				firstBad = v
			}
			return false
		}
		return true
	}
	walk := func(n ast.Node) {
		ast.Inspect(n, func(c ast.Node) bool {
			if c == nil {
				return false
			}
			return check(c)
		})
	}
	if sel.TargetList != nil && !b.existsBody[sel] {
		for _, item := range sel.TargetList.Items {
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				continue
			}
			if star, ok := rt.Val.(*ast.StarExpr); ok {
				if err := b.starGrouped(star, sc, groups); err != nil {
					return err
				}
				continue
			}
			walk(rt.Val)
			if firstBad != nil {
				return mssqlErr(8120, "Column '%s' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.", columnText(firstBad))
			}
		}
	}
	walk(sel.HavingClause)
	if firstBad != nil {
		return mssqlErr(8121, "Column '%s' is invalid in the HAVING clause because it is not contained in either an aggregate function or the GROUP BY clause.", columnText(firstBad))
	}
	if sel.OrderByClause != nil {
		for _, item := range sel.OrderByClause.Items {
			ob, ok := item.(*ast.OrderByItem)
			if !ok {
				continue
			}
			if _, ok := ordinal(ob.Expr); ok {
				continue
			}
			if ref, ok := ob.Expr.(*ast.ColumnRef); ok && ref.Table == "" && namesResult(sel, ref.Column) {
				continue
			}
			walk(ob.Expr)
			if firstBad != nil {
				return mssqlErr(8127, "Column '%s' is invalid in the ORDER BY clause because it is not contained in either an aggregate function or the GROUP BY clause.", columnText(firstBad))
			}
		}
	}
	return nil
}

func (b *binder) starGrouped(star *ast.StarExpr, sc *qscope, groups []groupItem) error {
	for _, r := range sc.rels {
		if star.Qualifier != "" && !strings.EqualFold(r.name, star.Qualifier) {
			continue
		}
		if r.opaque {
			continue
		}
		for _, col := range visible(r.cols) {
			ok := false
			for _, g := range groups {
				if gref, isRef := unparen(g.node).(*ast.ColumnRef); isRef && strings.EqualFold(gref.Column, col.name) {
					if gb, bound := b.bound[gref]; bound && gb.rel == r {
						ok = true
					}
				}
			}
			if !ok {
				return mssqlErr(8120, "Column '%s.%s' is invalid in the select list because it is not contained in either an aggregate function or the GROUP BY clause.", r.name, col.name)
			}
		}
	}
	return nil
}

func unparen(n ast.Node) ast.Node {
	for {
		p, ok := n.(*ast.ParenExpr)
		if !ok || p.Expr == nil {
			return n
		}
		n = p.Expr
	}
}

func columnText(ref *ast.ColumnRef) string {
	parts := []string{}
	for _, p := range []string{ref.Database, ref.Schema, ref.Table, ref.Column} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

// --- FROM ---

// tableExpr binds one table source of a FROM clause and returns the
// relations it adds. lateral are the relations to its left an APPLY lets it
// see.
func (b *binder) tableExpr(te ast.Node, sc *qscope, lateral []*qrel) ([]*qrel, error) {
	outer := &qscope{parent: sc.parent, rels: lateral}
	switch t := te.(type) {
	case *ast.TableRef:
		r, err := b.objectRel(t, t.Alias, nil, outer)
		if err != nil {
			return nil, err
		}
		return []*qrel{r}, nil
	case *ast.AliasedTableRef:
		r, err := b.aliased(t, outer)
		if err != nil {
			return nil, err
		}
		return []*qrel{r}, nil
	case *ast.JoinClause:
		left, err := b.tableExpr(t.Left, sc, lateral)
		if err != nil {
			return nil, err
		}
		rightLateral := lateral
		if t.Type == ast.JoinCrossApply || t.Type == ast.JoinOuterApply {
			rightLateral = append(append([]*qrel(nil), lateral...), left...)
		}
		right, err := b.tableExpr(t.Right, sc, rightLateral)
		if err != nil {
			return nil, err
		}
		if err := checkExposedNames(left, right); err != nil {
			return nil, err
		}
		rels := append(left, right...)
		if t.Condition != nil {
			on := &qscope{parent: sc.parent, rels: append(append([]*qrel(nil), lateral...), rels...), block: sc.block}
			if _, err := b.expr(t.Condition, on, ectx{clause: clOn, block: sc.block}); err != nil {
				return nil, err
			}
		}
		return rels, nil
	case *ast.TableVarRef:
		if !b.vars[strings.ToUpper(t.Name)] {
			return nil, mssqlErr(1087, "Must declare the table variable \"%s\".", t.Name)
		}
		name := t.Alias
		if name == "" {
			name = t.Name
		}
		if cols := b.tableVars[strings.ToUpper(t.Name)]; cols != nil {
			return []*qrel{{name: name, alias: true, cols: append([]qcol(nil), cols...)}}, nil
		}
		return []*qrel{{name: name, alias: true, opaque: true}}, nil
	case *ast.TableVarMethodCallRef:
		if !b.vars[strings.ToUpper(t.Var)] {
			return nil, mssqlErr(137, "Must declare the scalar variable \"%s\".", t.Var)
		}
		return []*qrel{{name: t.Alias, alias: true, opaque: true}}, nil
	case *ast.PivotExpr:
		r, err := b.pivot(t, sc, lateral)
		if err != nil {
			return nil, err
		}
		return []*qrel{r}, nil
	case *ast.UnpivotExpr:
		r, err := b.unpivot(t, sc, lateral)
		if err != nil {
			return nil, err
		}
		return []*qrel{r}, nil
	case *ast.FullTextTableRef:
		r, err := b.objectRel(t.Table, "", nil, outer)
		if err != nil {
			return nil, err
		}
		if b.c.fromMetadata || r.table == nil || b.c.fulltextTables[r.table] {
			return []*qrel{{name: t.Alias, alias: true, cols: []qcol{{name: "KEY"}, {name: "RANK", typ: qtNumber}}}}, nil
		}
		return nil, mssqlErr(7601, "Cannot use a CONTAINS or FREETEXT predicate on table or indexed view '%s' because it is not full-text indexed.", t.Table.Object)
	case *ast.SemanticTableRef:
		r, err := b.objectRel(t.Table, "", nil, outer)
		if err != nil {
			return nil, err
		}
		if b.c.fromMetadata || r.table == nil || b.c.fulltextTables[r.table] {
			return []*qrel{{name: t.Alias, alias: true, opaque: true}}, nil
		}
		return nil, mssqlErr(41202, "The semantic language statistics database is not registered, or the table '%s' is not semantically indexed.", t.Table.Object)
	case *ast.SubqueryExpr:
		cols, opaque, err := b.query(asSelect(t.Query), outer)
		if err != nil {
			return nil, err
		}
		return []*qrel{{alias: true, cols: cols, opaque: opaque}}, nil
	case *ast.FuncCallExpr:
		r, err := b.functionRel(t, "", nil, outer)
		if err != nil {
			return nil, err
		}
		return []*qrel{r}, nil
	}
	return []*qrel{{alias: true, opaque: true}}, nil
}

func asSelect(n ast.Node) *ast.SelectStmt {
	sel, _ := n.(*ast.SelectStmt)
	return sel
}

// aliased binds a table source that carries an alias, column aliases, or a
// TABLESAMPLE clause.
func (b *binder) aliased(t *ast.AliasedTableRef, outer *qscope) (*qrel, error) {
	if t.TableSample != nil && b.inView {
		if ref, ok := t.Table.(*ast.TableRef); ok {
			if _, err := b.objectRel(ref, ref.Alias, nil, outer); err != nil {
				return nil, err
			}
		}
		return nil, mssqlErr(478, "The TABLESAMPLE clause cannot be used in a view definition or inline table function definition.")
	}
	aliases := stringNodeList(t.Columns)
	switch inner := t.Table.(type) {
	case *ast.TableRef:
		alias := t.Alias
		if alias == "" {
			alias = inner.Alias
		}
		return b.objectRel(inner, alias, t.Hints, outer)
	case *ast.SubqueryExpr:
		cols, opaque, err := b.query(asSelect(inner.Query), outer)
		if err != nil {
			return nil, err
		}
		if opaque {
			return &qrel{name: t.Alias, alias: true, opaque: true}, nil
		}
		named, err := b.derivedColumns(cols, aliases, t.Alias)
		if err != nil {
			return nil, err
		}
		return &qrel{name: t.Alias, alias: true, cols: named}, nil
	case *ast.ValuesClause:
		cols, err := b.values(inner, outer)
		if err != nil {
			return nil, err
		}
		named, err := b.derivedColumns(cols, aliases, t.Alias)
		if err != nil {
			return nil, err
		}
		return &qrel{name: t.Alias, alias: true, cols: named}, nil
	case *ast.FuncCallExpr:
		return b.functionRel(inner, t.Alias, t.Columns, outer)
	}
	rels, err := b.tableExpr(t.Table, &qscope{parent: outer.parent, rels: outer.rels}, outer.rels)
	if err != nil {
		return nil, err
	}
	if len(rels) == 1 {
		r := *rels[0]
		if t.Alias != "" {
			r.name, r.alias, r.schema = t.Alias, true, ""
		}
		return &r, nil
	}
	return &qrel{name: t.Alias, alias: true, opaque: true}, nil
}

// values binds a table value constructor.
func (b *binder) values(v *ast.ValuesClause, outer *qscope) ([]qcol, error) {
	width := -1
	var cols []qcol
	if v.Rows == nil {
		return nil, nil
	}
	for _, row := range v.Rows.Items {
		list, ok := row.(*ast.List)
		if !ok {
			continue
		}
		if width >= 0 && len(list.Items) != width {
			return nil, mssqlErr(10709, "The number of columns for each row in a table value constructor must be the same.")
		}
		width = len(list.Items)
		for _, item := range list.Items {
			if _, err := b.expr(item, outer, ectx{clause: clOther}); err != nil {
				return nil, err
			}
		}
		if cols == nil {
			cols = make([]qcol, width)
		}
	}
	return cols, nil
}

// systemSchema reports a schema whose objects the catalog does not model.
func systemSchema(name string) bool {
	return strings.EqualFold(name, "sys") || strings.EqualFold(name, "INFORMATION_SCHEMA")
}

// objectName is a TableRef as written, for a message.
func objectName(ref *ast.TableRef) string {
	parts := []string{}
	for _, p := range []string{ref.Server, ref.Database, ref.Schema, ref.Object} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

// localObject reports whether ref names an object of this database.
func (b *binder) localObject(ref *ast.TableRef) bool {
	if ref.Server != "" {
		return false
	}
	return ref.Database == "" || (b.c.database != "" && strings.EqualFold(ref.Database, b.c.database))
}

// objectRel binds a named table source: a common table expression, a table,
// or a view.
func (b *binder) objectRel(ref *ast.TableRef, alias string, hints *ast.List, outer *qscope) (*qrel, error) {
	if ref == nil {
		return &qrel{alias: true, opaque: true}, nil
	}
	exposed := func(r *qrel) *qrel {
		if alias != "" {
			r.name, r.alias, r.schema = alias, true, ""
		}
		return r
	}
	if strings.HasPrefix(ref.Object, "#") {
		if b.inView {
			return nil, mssqlErr(4508, "Views or functions are not allowed on temporary tables. Table names that begin with '#' denote temporary tables.")
		}
		return exposed(&qrel{name: ref.Object, opaque: true}), nil
	}
	if !b.localObject(ref) {
		return exposed(&qrel{name: ref.Object, schema: ref.Schema, opaque: true}), nil
	}
	if ref.Schema == "" && ref.Database == "" {
		if r, ok := b.pseudo[strings.ToUpper(ref.Object)]; ok {
			c := *r
			c.name, c.alias, c.schema = ref.Object, true, ""
			return exposed(&c), nil
		}
		for s := outer; s != nil; s = s.parent {
			for _, q := range s.ctes {
				if strings.EqualFold(q.name, ref.Object) {
					if q.rel == nil {
						return exposed(&qrel{name: q.name, alias: true, opaque: true}), nil
					}
					r := *q.rel
					r.name, r.alias = ref.Object, true
					return exposed(&r), nil
				}
			}
		}
	}
	if b.schemaBound && ref.Schema == "" {
		return nil, mssqlErr(4512, "Cannot schema bind view '%s' because name '%s' is invalid for schema binding. Names must be in two-part format and an object cannot reference itself.", b.object, ref.Object)
	}
	if systemSchema(ref.Schema) {
		if b.schemaBound {
			return nil, mssqlErr(2720, "Cannot schema bind view '%s' because it references system object '%s'.", b.object, objectName(ref))
		}
		return exposed(&qrel{name: ref.Object, schema: ref.Schema, opaque: true}), nil
	}
	schema := b.objectSchema(ref.Schema, ref.Object)
	if schema == nil {
		return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(ref))
	}
	if t := b.c.resolveTableInSchema(schema, ref.Object); t != nil {
		if err := b.tableHints(hints, ref); err != nil {
			return nil, err
		}
		if err := b.temporal(ref, schema.Name, t.Name, outer); err != nil {
			return nil, err
		}
		cols := tableColumns(t)
		period := b.c.systemTimeTables[b.c.temporalKey(schema.Name, t.Name)]
		datetime2 := 0
		for i := range cols {
			cols[i].hidden = b.c.hiddenColumns[t.Columns[i]]
			if strings.EqualFold(cols[i].name, period[0]) || strings.EqualFold(cols[i].name, period[1]) {
				cols[i].noInsert = true
			}
			if strings.HasPrefix(strings.ToLower(t.Columns[i].Type), "datetime2") {
				datetime2++
			}
		}
		// A synced table may hide two datetime2 period columns from *.
		starOpaque := b.c.fromMetadata && datetime2 >= 2
		return exposed(&qrel{name: ref.Object, schema: schema.Name, cols: cols, starOpaque: starOpaque, table: t}), nil
	}
	if idx, ok := resolveExactOrFoldIndex(schema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, ref.Object); ok {
		v := schema.Views[idx]
		if b.schemaBound && !isSchemaBound(v.Definition) {
			return nil, mssqlErr(4513, "Cannot schema bind view '%s'. '%s.%s' is not schema bound.", b.object, schema.Name, v.Name)
		}
		// FOR SYSTEM_TIME on a view applies to the temporal tables it reads.
		if ref.Temporal != nil {
			for _, e := range []ast.ExprNode{ref.Temporal.Start, ref.Temporal.End} {
				if _, err := b.expr(e, outer, ectx{clause: clOther}); err != nil {
					return nil, err
				}
			}
		}
		r := b.viewRel(v, schema.Name)
		r.name, r.schema = ref.Object, schema.Name
		return exposed(r), nil
	}
	switch b.c.lookupSchemaObject(schema, ref.Object, false) {
	case kindSynonym:
		return exposed(&qrel{name: ref.Object, schema: schema.Name, opaque: true}), nil
	case kindFunction:
		return nil, mssqlErr(216, "Parameters were not supplied for the function '%s'.", objectName(ref))
	}
	return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(ref))
}

// objectSchema is the schema a name resolves in: the one it names, or for a
// one-part name the default schema and then dbo, as SQL Server looks.
func (b *binder) objectSchema(schemaName, object string) *metadata.SchemaMetadata {
	if schemaName != "" {
		return b.c.GetSchema(schemaName)
	}
	schema := b.c.GetSchema(b.c.defaultSchema)
	if schema != nil && b.c.lookupSchemaObject(schema, object, false) != kindNone {
		return schema
	}
	if !strings.EqualFold(b.c.defaultSchema, defaultSchema) {
		if dbo := b.c.GetSchema(defaultSchema); dbo != nil && b.c.lookupSchemaObject(dbo, object, false) != kindNone {
			return dbo
		}
	}
	return schema
}

func (b *binder) tableHints(hints *ast.List, ref *ast.TableRef) error {
	for _, list := range []*ast.List{hints, ref.Hints} {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			if h, ok := item.(*ast.TableHint); ok && strings.EqualFold(h.Name, "NOEXPAND") {
				return mssqlErr(8171, "Hint 'noexpand' on object '%s' is invalid.", objectName(ref))
			}
		}
	}
	return nil
}

// temporal checks FOR SYSTEM_TIME names a system-versioned table.
func (b *binder) temporal(ref *ast.TableRef, schemaName, tableName string, outer *qscope) error {
	if ref.Temporal == nil {
		return nil
	}
	for _, e := range []ast.ExprNode{ref.Temporal.Start, ref.Temporal.End} {
		if _, err := b.expr(e, outer, ectx{clause: clOther}); err != nil {
			return err
		}
	}
	if b.c.fromMetadata || b.c.systemVersionedTables[b.c.temporalKey(schemaName, tableName)] {
		return nil
	}
	return mssqlErr(13544, "Temporal FOR SYSTEM_TIME clause can only be used with system-versioned tables. '%s.%s' is not a system-versioned table.", schemaName, tableName)
}

func tableColumns(t *metadata.TableMetadata) []qcol {
	cols := make([]qcol, 0, len(t.Columns))
	for _, c := range t.Columns {
		base := strings.Trim(strings.ToLower(c.Type), "[ ")
		cols = append(cols, qcol{
			name:     c.Name,
			typ:      typeClass(c.Type),
			identity: c.IsIdentity,
			guid:     strings.HasPrefix(base, "uniqueidentifier"),
			noInsert: isComputedColumnType(c.Type) || strings.HasPrefix(base, "timestamp") || strings.HasPrefix(base, "rowversion"),
		})
	}
	return cols
}

// viewRel binds a view's stored definition for the columns it gives. A view
// that does not bind is opaque.
func (b *binder) viewRel(v *metadata.ViewMetadata, schemaName string) *qrel {
	// A view keeps the columns it was created with until it is altered or
	// refreshed, whatever its tables became since.
	cols, known := b.c.viewColumns[v]
	if cols != nil {
		return &qrel{cols: append([]qcol(nil), cols...)}
	}
	if !known && len(v.Columns) > 0 {
		cols := make([]qcol, 0, len(v.Columns))
		for _, col := range v.Columns {
			cols = append(cols, qcol{name: col.Name, typ: typeClass(col.Type)})
		}
		return &qrel{cols: cols}
	}
	if b.depth >= maxBindDepth {
		return &qrel{opaque: true}
	}
	stmts, err := mssqlparser.Parse(v.Definition)
	if err != nil || stmts == nil {
		return &qrel{opaque: true}
	}
	for _, item := range stmts.Items {
		n, ok := item.(*ast.CreateViewStmt)
		if !ok {
			continue
		}
		sel, ok := n.Query.(*ast.SelectStmt)
		if !ok || sel == nil {
			break
		}
		inner := newBinder(b.c, v.Definition, v.Name, b.depth+1)
		inner.inView = true
		cols, opaque, err := inner.query(sel, nil)
		if err != nil || opaque {
			break
		}
		if names := stringNodeList(n.Columns); len(names) == len(cols) && len(names) > 0 {
			for i := range cols {
				cols[i].name = names[i]
			}
		}
		return &qrel{cols: cols}
	}
	return &qrel{opaque: true}
}

// builtinTableFunctions are the table-valued built-ins a FROM clause calls
// without a schema.
var builtinTableFunctions = map[string]bool{
	"STRING_SPLIT": true, "OPENJSON": true, "GENERATE_SERIES": true, "OPENXML": true,
	"OPENROWSET": true, "OPENQUERY": true, "OPENDATASOURCE": true, "PREDICT": true,
}

// functionRel binds a function call in a FROM clause.
func (b *binder) functionRel(f *ast.FuncCallExpr, alias string, columns *ast.List, outer *qscope) (*qrel, error) {
	name := ""
	schemaName := ""
	if f.Name != nil {
		name, schemaName = f.Name.Object, f.Name.Schema
	}
	rel := func(cols []qcol, opaque bool) *qrel {
		r := &qrel{name: alias, alias: alias != "", cols: cols, opaque: opaque}
		if alias == "" {
			r.name = name
		}
		return r
	}
	bindArgs := func() error {
		_, err := b.expr(f.Args, outer, ectx{clause: clOther})
		return err
	}
	if f.Name != nil && schemaName == "" && f.Name.Database == "" && builtinTableFunctions[strings.ToUpper(name)] {
		if err := bindArgs(); err != nil {
			return nil, err
		}
		switch strings.ToUpper(name) {
		case "STRING_SPLIT":
			// enable_ordinal adds an ordinal column when it is 1, and none
			// when it is 0 or NULL; any other value is unknown here.
			cols := []qcol{{name: "value", typ: qtString}}
			if f.Args != nil && len(f.Args.Items) == 3 {
				lit, ok := unparen(f.Args.Items[2]).(*ast.Literal)
				switch {
				case ok && lit.Type == ast.LitNull, ok && lit.Type == ast.LitInteger && lit.Ival == 0:
				case ok && lit.Type == ast.LitInteger && lit.Ival == 1:
					cols = append(cols, qcol{name: "ordinal", typ: qtNumber})
				default:
					return rel(nil, true), nil
				}
			}
			return rel(cols, false), nil
		case "GENERATE_SERIES":
			return rel([]qcol{{name: "value", typ: qtNumber}}, false), nil
		case "OPENJSON":
			if columns == nil || len(columns.Items) == 0 {
				return rel([]qcol{{name: "key", typ: qtString}, {name: "value", typ: qtString}, {name: "type", typ: qtNumber}}, false), nil
			}
			return b.rowsetColumns(columns, rel)
		case "OPENXML":
			if columns != nil && len(columns.Items) > 0 {
				return b.rowsetColumns(columns, rel)
			}
		}
		return rel(nil, true), nil
	}
	// col.nodes(path) reads rows of an xml column's value: the columns are
	// the ones its alias list names.
	if f.Name != nil && f.Name.Server == "" && strings.EqualFold(name, "nodes") && schemaName != "" {
		target := &ast.ColumnRef{Table: f.Name.Database, Column: schemaName, Loc: f.Name.Loc}
		if b.resolves(target, outer) {
			if err := bindArgs(); err != nil {
				return nil, err
			}
			names := stringNodeList(columns)
			if len(names) == 0 {
				return rel(nil, true), nil
			}
			cols := make([]qcol, len(names))
			for i, n := range names {
				cols[i] = qcol{name: n, typ: qtXML}
			}
			return rel(cols, false), nil
		}
	}
	if f.Name == nil || !b.localObject(f.Name) || systemSchema(schemaName) {
		if err := bindArgs(); err != nil {
			return nil, err
		}
		return rel(nil, true), nil
	}
	schema := b.objectSchema(schemaName, name)
	if schema == nil {
		return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(f.Name))
	}
	if columns != nil && len(columns.Items) > 0 {
		if kind := b.c.lookupSchemaObject(schema, name, false); kind != kindNone {
			return nil, mssqlErr(317, "Table-valued function '%s' cannot have a column alias.", name)
		}
	}
	switch b.c.lookupSchemaObject(schema, name, false) {
	case kindTable, kindView:
		if oldStyleHints(f.Args) {
			ref := *f.Name
			ref.Alias = alias
			return b.objectRel(&ref, alias, nil, outer)
		}
		if err := bindArgs(); err != nil {
			return nil, err
		}
		return nil, mssqlErr(215, "Parameters supplied for object '%s' which is not a function. If the parameters are intended as a table hint, a WITH keyword is required.", objectName(f.Name))
	case kindFunction:
	default:
		return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(f.Name))
	}
	idx, ok := resolveExactOrFoldIndex(schema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, name)
	if !ok {
		return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(f.Name))
	}
	fn := schema.Functions[idx]
	info := b.functionInfo(fn)
	if info.category == "scalar" {
		return nil, mssqlErr(ErrInvalidObjectName, "Invalid object name '%s'.", objectName(f.Name))
	}
	if b.schemaBound {
		if schemaName == "" {
			return nil, mssqlErr(4512, "Cannot schema bind view '%s' because name '%s' is invalid for schema binding. Names must be in two-part format and an object cannot reference itself.", b.object, name)
		}
		if !isSchemaBound(fn.Definition) {
			return nil, mssqlErr(4513, "Cannot schema bind view '%s'. '%s.%s' is not schema bound.", b.object, schema.Name, fn.Name)
		}
	}
	if err := b.checkArgCount(f, info, schema.Name, fn.Name); err != nil {
		return nil, err
	}
	if err := bindArgs(); err != nil {
		return nil, err
	}
	return rel(info.cols, info.cols == nil), nil
}

// tableHintWords are the table hints SQL Server reads in parentheses after a
// table's name without WITH: FROM a(FORCESEEK) (measured,
// testdata/queries/simulate.tsv).
var tableHintWords = map[string]bool{
	"FORCESCAN": true, "FORCESEEK": true, "HOLDLOCK": true, "NOEXPAND": true, "NOLOCK": true, "NOWAIT": true,
	"PAGLOCK": true, "READCOMMITTED": true, "READCOMMITTEDLOCK": true, "READPAST": true, "READUNCOMMITTED": true,
	"REPEATABLEREAD": true, "ROWLOCK": true, "SERIALIZABLE": true, "SNAPSHOT": true, "TABLOCK": true,
	"TABLOCKX": true, "UPDLOCK": true, "XLOCK": true,
}

// oldStyleHints reports whether a table's parenthesized arguments are table
// hints.
func oldStyleHints(args *ast.List) bool {
	if args == nil || len(args.Items) == 0 {
		return false
	}
	for _, item := range args.Items {
		ref, ok := item.(*ast.ColumnRef)
		if !ok || ref.Table != "" || !tableHintWords[strings.ToUpper(ref.Column)] {
			return false
		}
	}
	return true
}

func (b *binder) rowsetColumns(columns *ast.List, rel func([]qcol, bool) *qrel) (*qrel, error) {
	var cols []qcol
	for i, item := range columns.Items {
		def, ok := item.(*ast.RowsetColumnDef)
		if !ok {
			return rel(nil, true), nil
		}
		if def.DataType != nil && def.DataType.Schema == "" && !isValidDataType(def.DataType, false) {
			return nil, mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", i+1, def.DataType.Name)
		}
		cols = append(cols, qcol{name: def.Name, typ: dataTypeClass(def.DataType)})
	}
	return rel(cols, false), nil
}

// fnInfo is what a query needs to know of a user-defined function.
type fnInfo struct {
	category string // scalar, itvf, mstvf, or "" when its definition does not parse
	params   int
	cols     []qcol // the columns of a table-valued function, nil when unknown
}

func (b *binder) functionInfo(fn *metadata.FunctionMetadata) fnInfo {
	stmts, err := mssqlparser.Parse(fn.Definition)
	if err != nil || stmts == nil {
		return fnInfo{params: -1}
	}
	for _, item := range stmts.Items {
		n, ok := item.(*ast.CreateFunctionStmt)
		if !ok {
			continue
		}
		info := fnInfo{category: functionCategory(n)}
		if n.Params != nil {
			info.params = len(n.Params.Items)
		}
		switch info.category {
		case "mstvf":
			if n.ReturnsTable != nil && n.ReturnsTable.Columns != nil {
				for _, item := range n.ReturnsTable.Columns.Items {
					if cd, ok := item.(*ast.ColumnDef); ok {
						info.cols = append(info.cols, qcol{name: cd.Name, typ: dataTypeClass(cd.DataType)})
					}
				}
			}
		case "itvf":
			info.cols = b.inlineFunctionColumns(n, fn)
		}
		return info
	}
	return fnInfo{params: -1}
}

// inlineFunctionColumns binds an inline table-valued function's query for
// its columns, nil when it does not bind.
func (b *binder) inlineFunctionColumns(n *ast.CreateFunctionStmt, fn *metadata.FunctionMetadata) []qcol {
	if b.depth >= maxBindDepth {
		return nil
	}
	var sel *ast.SelectStmt
	ast.Inspect(n.Body, func(node ast.Node) bool {
		if sel != nil {
			return false
		}
		if s, ok := node.(*ast.SelectStmt); ok {
			sel = s
			return false
		}
		return true
	})
	if sel == nil {
		return nil
	}
	inner := newBinder(b.c, fn.Definition, fn.Name, b.depth+1)
	if n.Params != nil {
		for _, item := range n.Params.Items {
			if p, ok := item.(*ast.ParamDef); ok {
				inner.vars[strings.ToUpper(p.Name)] = true
			}
		}
	}
	cols, opaque, err := inner.query(sel, nil)
	if err != nil || opaque {
		return nil
	}
	return cols
}

func (b *binder) checkArgCount(f *ast.FuncCallExpr, info fnInfo, schemaName, name string) error {
	if info.params < 0 {
		return nil
	}
	args := 0
	if f.Args != nil {
		args = len(f.Args.Items)
	}
	if args > info.params {
		return mssqlErr(8144, "Procedure or function %s has too many arguments specified.", name)
	}
	if args < info.params {
		return mssqlErr(313, "An insufficient number of arguments were supplied for the procedure or function %s.%s.", schemaName, name)
	}
	return nil
}

// pivot binds PIVOT: its source, the aggregate's and FOR's columns, and
// the columns it gives.
func (b *binder) pivot(p *ast.PivotExpr, sc *qscope, lateral []*qrel) (*qrel, error) {
	src, err := b.tableExpr(p.Source, sc, lateral)
	if err != nil {
		return nil, err
	}
	if len(src) != 1 || src[0].opaque {
		return &qrel{name: p.Alias, alias: true, opaque: true}, nil
	}
	in := &qscope{parent: sc.parent, rels: src}
	used := map[string]bool{}
	var argErr error
	ast.Inspect(p.AggFunc, func(n ast.Node) bool {
		if argErr != nil {
			return false
		}
		if ref, ok := n.(*ast.ColumnRef); ok {
			if _, err := b.column(ref, in); err != nil {
				argErr = err
			}
			used[strings.ToUpper(ref.Column)] = true
			return false
		}
		return true
	})
	if argErr != nil {
		return nil, argErr
	}
	if p.ForCol != nil {
		if _, err := b.column(p.ForCol, in); err != nil {
			return nil, err
		}
		used[strings.ToUpper(p.ForCol.Column)] = true
	}
	var cols []qcol
	for _, c := range src[0].cols {
		if !used[strings.ToUpper(c.name)] {
			cols = append(cols, c)
		}
	}
	if p.InValues != nil {
		for _, item := range p.InValues.Items {
			if ref, ok := item.(*ast.ColumnRef); ok {
				cols = append(cols, qcol{name: ref.Column})
			}
		}
	}
	return &qrel{name: p.Alias, alias: true, cols: cols}, nil
}

// unpivot binds UNPIVOT.
func (b *binder) unpivot(u *ast.UnpivotExpr, sc *qscope, lateral []*qrel) (*qrel, error) {
	src, err := b.tableExpr(u.Source, sc, lateral)
	if err != nil {
		return nil, err
	}
	if len(src) != 1 || src[0].opaque {
		return &qrel{name: u.Alias, alias: true, opaque: true}, nil
	}
	in := &qscope{parent: sc.parent, rels: src}
	used := map[string]bool{}
	if u.InCols != nil {
		for _, item := range u.InCols.Items {
			if ref, ok := item.(*ast.ColumnRef); ok {
				if _, err := b.column(ref, in); err != nil {
					return nil, err
				}
				used[strings.ToUpper(ref.Column)] = true
			}
		}
	}
	var cols []qcol
	for _, c := range src[0].cols {
		if !used[strings.ToUpper(c.name)] {
			cols = append(cols, c)
		}
	}
	if u.ForCol != nil {
		cols = append(cols, qcol{name: u.ForCol.Column, typ: qtString})
	}
	if u.ValueCol != nil {
		cols = append(cols, qcol{name: u.ValueCol.Column})
	}
	return &qrel{name: u.Alias, alias: true, cols: cols}, nil
}

// --- expressions ---

// column binds a column reference.
func (b *binder) column(ref *ast.ColumnRef, sc *qscope) (qtype, error) {
	special := b.pseudoColumn(ref)
	identity := special && (strings.EqualFold(ref.Column, "$IDENTITY") || strings.EqualFold(ref.Column, "IDENTITYCOL"))
	if ref.Table != "" {
		for s := sc; s != nil; s = s.parent {
			for _, r := range s.rels {
				if !relMatches(r, ref.Schema, ref.Table) {
					continue
				}
				if r.opaque {
					b.bound[ref] = colBinding{rel: r, block: s.block}
					return qtUnknown, nil
				}
				var c *qcol
				switch {
				case identity:
					c = r.identityCol()
				case special:
					c = r.rowguidCol()
				default:
					c = r.col(ref.Column)
				}
				if c == nil {
					return qtUnknown, mssqlErr(ErrInvalidColumnName, "Invalid column name '%s'.", ref.Column)
				}
				b.bound[ref] = colBinding{rel: r, block: s.block}
				return c.typ, nil
			}
		}
		if special {
			return qtUnknown, mssqlErr(107, "The column prefix '%s' does not match with a table name or alias name used in the query.", ref.Table)
		}
		return qtUnknown, mssqlErr(4104, "The multi-part identifier \"%s\" could not be bound.", columnText(ref))
	}
	for s := sc; s != nil; s = s.parent {
		var match *qrel
		var col *qcol
		opaque := false
		for _, r := range s.rels {
			if r.opaque {
				opaque = true
				continue
			}
			var c *qcol
			switch {
			case identity:
				c = r.identityCol()
			case special:
				c = r.rowguidCol()
			default:
				c = r.col(ref.Column)
			}
			if c == nil {
				continue
			}
			if match != nil {
				return qtUnknown, mssqlErr(209, "Ambiguous column name '%s'.", ref.Column)
			}
			match, col = r, c
		}
		if match != nil {
			b.bound[ref] = colBinding{rel: match, block: s.block}
			return col.typ, nil
		}
		if opaque {
			b.bound[ref] = colBinding{block: s.block}
			return qtUnknown, nil
		}
	}
	return qtUnknown, mssqlErr(ErrInvalidColumnName, "Invalid column name '%s'.", ref.Column)
}

// children binds every node below n.
func (b *binder) children(n ast.Node, sc *qscope, x ectx) error {
	var err error
	ast.Inspect(n, func(c ast.Node) bool {
		if err != nil || c == nil {
			return false
		}
		if c == n {
			return true
		}
		_, err = b.expr(c, sc, x)
		return false
	})
	return err
}

// expr binds an expression and returns its type class.
func (b *binder) expr(n ast.Node, sc *qscope, x ectx) (qtype, error) {
	if n == nil || isNilNode(n) {
		return qtUnknown, nil
	}
	switch e := n.(type) {
	case *ast.ColumnRef:
		return b.column(e, sc)
	case *ast.VariableRef:
		if strings.HasPrefix(e.Name, "@@") || b.vars[strings.ToUpper(e.Name)] {
			return qtUnknown, nil
		}
		return qtUnknown, mssqlErr(137, "Must declare the scalar variable \"%s\".", e.Name)
	case *ast.Literal:
		switch e.Type {
		case ast.LitString:
			return qtString, nil
		case ast.LitInteger, ast.LitFloat, ast.LitMoney:
			return qtNumber, nil
		}
		return qtUnknown, nil
	case *ast.ParenExpr:
		return b.expr(e.Expr, sc, x)
	case *ast.SubqueryExpr:
		return b.subquery(e.Query, sc, x, true)
	case *ast.ExistsExpr:
		// An EXISTS projection is not checked against its GROUP BY: SELECT *
		// FROM t GROUP BY id is a valid EXISTS body.
		if sel, ok := e.Subquery.(*ast.SelectStmt); ok {
			b.existsBody[sel] = true
		}
		_, err := b.subquery(e.Subquery, sc, x, false)
		return qtNumber, err
	case *ast.InExpr:
		if _, err := b.expr(e.Expr, sc, x); err != nil {
			return qtUnknown, err
		}
		if _, err := b.expr(e.List, sc, x); err != nil {
			return qtUnknown, err
		}
		if e.Subquery != nil {
			q := e.Subquery
			if s, ok := q.(*ast.SubqueryExpr); ok {
				if _, err := b.subquery(s.Query, sc, x, true); err != nil {
					return qtUnknown, err
				}
			} else if _, err := b.expr(q, sc, x); err != nil {
				return qtUnknown, err
			}
		}
		return qtNumber, nil
	case *ast.SubqueryComparisonExpr:
		if _, err := b.expr(e.Left, sc, x); err != nil {
			return qtUnknown, err
		}
		_, err := b.subquery(e.Subquery, sc, x, true)
		return qtNumber, err
	case *ast.SelectStmt:
		return b.subquery(e, sc, x, true)
	case *ast.FuncCallExpr:
		return b.funcCall(e, sc, x)
	case *ast.CastExpr:
		return b.cast(e.Expr, nil, e.DataType, sc, x)
	case *ast.TryCastExpr:
		return b.cast(e.Expr, nil, e.DataType, sc, x)
	case *ast.ConvertExpr:
		return b.cast(e.Expr, e.Style, e.DataType, sc, x)
	case *ast.TryConvertExpr:
		return b.cast(e.Expr, e.Style, e.DataType, sc, x)
	case *ast.ParseExpr:
		if _, err := b.expr(e.Expr, sc, x); err != nil {
			return qtUnknown, err
		}
		if _, err := b.expr(e.Culture, sc, x); err != nil {
			return qtUnknown, err
		}
		if e.DataType != nil && !isValidDataType(e.DataType, false) {
			fn := "PARSE"
			if e.Try {
				fn = "TRY_PARSE"
			}
			return qtUnknown, mssqlErr(10761, "Invalid data type %s in function %s.", e.DataType.Name, fn)
		}
		return dataTypeClass(e.DataType), nil
	case *ast.MethodCallExpr:
		if e.Type != nil && e.Type.Schema == "" && !isValidDataType(e.Type, false) {
			return qtUnknown, mssqlErr(ErrUnknownType, "Type %s is not a defined system type.", e.Type.Name)
		}
		return qtUnknown, b.children(e.Args, sc, x)
	case *ast.CollateExpr:
		t, err := b.expr(e.Expr, sc, x)
		if err != nil {
			return qtUnknown, err
		}
		if t != qtString && t != qtUnknown {
			return qtUnknown, mssqlErr(447, "Expression type %s is invalid for COLLATE clause.", t.typeName())
		}
		return qtString, nil
	case *ast.NextValueForExpr:
		if b.inView {
			return qtUnknown, mssqlErr(11719, "NEXT VALUE FOR function is not allowed in check constraints, default objects, computed columns, views, user-defined functions, user-defined aggregates, user-defined table types, sub-queries, common table expressions, derived tables or return statements.")
		}
		return qtNumber, nil
	case *ast.FullTextPredicate:
		if err := b.children(e.Columns, sc, x); err != nil {
			return qtUnknown, err
		}
		if b.fulltextUnknown(e.Columns, sc) {
			return qtNumber, nil
		}
		return qtUnknown, mssqlErr(7601, "Cannot use a CONTAINS or FREETEXT predicate on table or indexed view because it is not full-text indexed.")
	case *ast.BinaryExpr:
		lt, err := b.expr(e.Left, sc, x)
		if err != nil {
			return qtUnknown, err
		}
		rt, err := b.expr(e.Right, sc, x)
		if err != nil {
			return qtUnknown, err
		}
		switch e.Op {
		case ast.BinOpAdd:
			// A string and a number add as numbers, by type precedence; the
			// classes cannot tell which number, so only two strings concatenate.
			if lt == qtString && rt == qtString {
				return qtString, nil
			}
			if lt == qtNumber && rt == qtNumber {
				return qtNumber, nil
			}
			return qtUnknown, nil
		case ast.BinOpSub, ast.BinOpMul, ast.BinOpDiv, ast.BinOpMod, ast.BinOpBitAnd, ast.BinOpBitOr, ast.BinOpBitXor, ast.BinOpShiftLeft, ast.BinOpShiftRight:
			if lt == qtNumber && rt == qtNumber {
				return qtNumber, nil
			}
			return qtUnknown, nil
		}
		return qtNumber, nil
	case *ast.UnaryExpr:
		t, err := b.expr(e.Operand, sc, x)
		if e.Op == ast.UnaryNot {
			return qtNumber, err
		}
		return t, err
	case *ast.StarExpr, *ast.DatePart:
		return qtUnknown, nil
	case *ast.OverClause:
		return qtUnknown, b.over(e, "", sc, x)
	}
	return qtUnknown, b.children(n, sc, x)
}

// fulltextUnknown reports whether a full-text predicate's columns may be
// full-text indexed: on a table Exec indexed, on a relation that is no table
// the catalog tracks, or in a catalog loaded from synced metadata.
func (b *binder) fulltextUnknown(cols *ast.List, sc *qscope) bool {
	if b.c.fromMetadata {
		return true
	}
	var rels []*qrel
	if cols != nil {
		for _, item := range cols.Items {
			switch v := item.(type) {
			case *ast.ColumnRef:
				if bound, ok := b.bound[v]; ok && bound.rel != nil {
					rels = append(rels, bound.rel)
				}
			case *ast.StarExpr:
				rels = append(rels, sc.rels...)
			}
		}
	}
	if len(rels) == 0 {
		return true
	}
	for _, r := range rels {
		if r.table == nil || b.c.fulltextTables[r.table] {
			return true
		}
	}
	return false
}

func isNilNode(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.List:
		return v == nil
	case *ast.SelectStmt:
		return v == nil
	case *ast.ColumnRef:
		return v == nil
	case *ast.FuncCallExpr:
		return v == nil
	case *ast.SubqueryExpr:
		return v == nil
	case *ast.OverClause:
		return v == nil
	}
	return false
}

// subquery binds a subquery in an expression. single holds it to one result
// column.
func (b *binder) subquery(q ast.Node, sc *qscope, x ectx, single bool) (qtype, error) {
	if x.aggDepth > 0 {
		return qtUnknown, mssqlErr(130, "Cannot perform an aggregate function on an expression containing an aggregate or a subquery.")
	}
	if x.clause == clGroup {
		return qtUnknown, mssqlErr(144, "Cannot use an aggregate or a subquery in an expression used for the group by list of a GROUP BY clause.")
	}
	sel, ok := q.(*ast.SelectStmt)
	if !ok {
		if s, isSub := q.(*ast.SubqueryExpr); isSub {
			sel, _ = s.Query.(*ast.SelectStmt)
		}
	}
	if sel == nil {
		return qtUnknown, nil
	}
	b.clauses = append(b.clauses, x.clause)
	cols, opaque, err := b.query(sel, sc)
	b.clauses = b.clauses[:len(b.clauses)-1]
	if err != nil {
		return qtUnknown, err
	}
	if single && !opaque && len(cols) != 1 {
		return qtUnknown, mssqlErr(116, "Only one expression can be specified in the select list when the subquery is not introduced with EXISTS.")
	}
	if len(cols) == 1 {
		return cols[0].typ, nil
	}
	return qtUnknown, nil
}

func (b *binder) cast(e, style ast.ExprNode, dt *ast.DataType, sc *qscope, x ectx) (qtype, error) {
	if _, err := b.expr(e, sc, x); err != nil {
		return qtUnknown, err
	}
	if _, err := b.expr(style, sc, x); err != nil {
		return qtUnknown, err
	}
	if dt != nil && dt.Schema == "" && !isValidDataType(dt, false) {
		return qtUnknown, mssqlErr(ErrUnknownType, "Type %s is not a defined system type.", dt.Name)
	}
	return dataTypeClass(dt), nil
}

// aggregateFunctions are the built-in aggregates.
var aggregateFunctions = map[string]bool{
	"AVG": true, "CHECKSUM_AGG": true, "COUNT": true, "COUNT_BIG": true, "GROUPING": true, "GROUPING_ID": true,
	"MAX": true, "MIN": true, "STDEV": true, "STDEVP": true, "STRING_AGG": true, "SUM": true, "VAR": true, "VARP": true,
	"APPROX_COUNT_DISTINCT": true, "APPROX_PERCENTILE_CONT": true, "APPROX_PERCENTILE_DISC": true,
}

// windowFunctions are the built-ins that need OVER, and whether OVER needs
// ORDER BY.
var windowFunctions = map[string]bool{
	"ROW_NUMBER": true, "RANK": true, "DENSE_RANK": true, "NTILE": true, "LAG": true, "LEAD": true,
	"FIRST_VALUE": true, "LAST_VALUE": true, "PERCENT_RANK": true, "CUME_DIST": true,
	"PERCENTILE_CONT": false, "PERCENTILE_DISC": false,
}

func isAggregateCall(f *ast.FuncCallExpr) bool {
	if f.Name == nil || f.Name.Schema != "" || f.Over != nil {
		return false
	}
	return aggregateFunctions[strings.ToUpper(f.Name.Object)]
}

// builtinArity holds the argument counts of built-ins whose count SQL Server
// checks when it binds a call: min and max.
var builtinArity = map[string][2]int{
	"ABS": {1, 1}, "ASCII": {1, 1}, "CEILING": {1, 1}, "CHAR": {1, 1}, "CHARINDEX": {2, 3}, "CONCAT": {2, 254},
	"DATEADD": {3, 3}, "DATEDIFF": {3, 3}, "DATEDIFF_BIG": {3, 3}, "DATENAME": {2, 2}, "DATEPART": {2, 2},
	"DATEFROMPARTS": {3, 3}, "DAY": {1, 1}, "DIFFERENCE": {2, 2}, "EOMONTH": {1, 2}, "FLOOR": {1, 1},
	"GETDATE": {0, 0}, "GETUTCDATE": {0, 0}, "ISNULL": {2, 2}, "LEFT": {2, 2}, "LEN": {1, 1},
	"LOWER": {1, 1}, "MONTH": {1, 1}, "NCHAR": {1, 1}, "NEWID": {0, 0}, "PATINDEX": {2, 2},
	"POWER": {2, 2}, "REPLACE": {3, 3}, "REPLICATE": {2, 2}, "REVERSE": {1, 1}, "RIGHT": {2, 2},
	"ROUND": {2, 3}, "SIGN": {1, 1}, "SOUNDEX": {1, 1}, "SPACE": {1, 1}, "SQRT": {1, 1},
	"SQUARE": {1, 1}, "STUFF": {4, 4}, "SUBSTRING": {3, 3}, "SYSDATETIME": {0, 0}, "UNICODE": {1, 1},
	"UPPER": {1, 1}, "YEAR": {1, 1},
}

// datePartFunctions are the built-ins whose first argument is a datepart.
var datePartFunctions = map[string]bool{
	"DATEADD": true, "DATEDIFF": true, "DATEDIFF_BIG": true, "DATENAME": true, "DATEPART": true, "DATETRUNC": true,
}

var dateParts = map[string]bool{
	"YEAR": true, "YY": true, "YYYY": true, "QUARTER": true, "QQ": true, "Q": true, "MONTH": true, "MM": true, "M": true,
	"DAYOFYEAR": true, "DY": true, "Y": true, "DAY": true, "DD": true, "D": true, "WEEK": true, "WK": true, "WW": true,
	"WEEKDAY": true, "DW": true, "W": true, "HOUR": true, "HH": true, "MINUTE": true, "MI": true, "N": true,
	"SECOND": true, "SS": true, "S": true, "MILLISECOND": true, "MS": true, "MICROSECOND": true, "MCS": true,
	"NANOSECOND": true, "NS": true, "TZOFFSET": true, "TZ": true, "ISO_WEEK": true, "ISOWK": true, "ISOWW": true,
}

// stringFunctions return a character string.
var stringFunctions = map[string]bool{
	"CONCAT": true, "CONCAT_WS": true, "FORMAT": true, "LEFT": true, "LOWER": true, "LTRIM": true, "REPLACE": true,
	"REPLICATE": true, "REVERSE": true, "RIGHT": true, "RTRIM": true, "STR": true, "STUFF": true, "SUBSTRING": true,
	"TRIM": true, "UPPER": true, "DATENAME": true, "STRING_AGG": true, "QUOTENAME": true, "TRANSLATE": true,
	"JSON_VALUE": true, "JSON_QUERY": true, "JSON_MODIFY": true, "SOUNDEX": true, "SPACE": true, "CHAR": true, "NCHAR": true,
}

// funcCall binds a function call in an expression.
func (b *binder) funcCall(f *ast.FuncCallExpr, sc *qscope, x ectx) (qtype, error) {
	if f.Name == nil {
		return qtUnknown, b.children(f, sc, x)
	}
	if f.Name.Schema != "" || f.Name.Database != "" || f.Name.Server != "" {
		return b.udfCall(f, sc, x)
	}
	name := strings.ToUpper(f.Name.Object)
	if name == "UPDATE" && b.pseudo != nil {
		// UPDATE ( column ) in a DML trigger names a column of the table
		// the trigger is on.
		in := &qscope{rels: []*qrel{b.pseudo["INSERTED"]}}
		_, err := b.expr(f.Args, in, x)
		return qtNumber, err
	}
	agg := aggregateFunctions[name] && f.Over == nil
	_, win := windowFunctions[name]
	win = win || f.Over != nil
	if win && name != "PERCENTILE_CONT" && name != "PERCENTILE_DISC" || f.Over != nil {
		switch x.clause {
		case clWhere, clOn, clGroup, clHaving:
			return qtUnknown, mssqlErr(4108, "Windowed functions can only appear in the SELECT or ORDER BY clauses.")
		}
	}
	if agg {
		switch {
		case x.clause == clWhere && !b.outerAggregate(f, sc, x):
			return qtUnknown, mssqlErr(147, "An aggregate may not appear in the WHERE clause unless it is in a subquery contained in a HAVING clause or a select list, and the column being aggregated is an outer reference.")
		case x.clause == clGroup:
			return qtUnknown, mssqlErr(144, "Cannot use an aggregate or a subquery in an expression used for the group by list of a GROUP BY clause.")
		case x.aggDepth > 0:
			return qtUnknown, mssqlErr(130, "Cannot perform an aggregate function on an expression containing an aggregate or a subquery.")
		}
		if x.block != nil {
			x.block.hasAgg = true
		}
	}
	nargs := 0
	if f.Args != nil {
		nargs = len(f.Args.Items)
	}
	if ar, ok := builtinArity[name]; ok && !f.Star && (nargs < ar[0] || nargs > ar[1]) {
		if ar[0] == ar[1] {
			return qtUnknown, mssqlErr(174, "The %s function requires %d argument(s).", strings.ToLower(name), ar[0])
		}
		return qtUnknown, mssqlErr(189, "The %s function requires %d to %d arguments.", strings.ToLower(name), ar[0], ar[1])
	}
	if datePartFunctions[name] && nargs > 0 {
		if dp, ok := f.Args.Items[0].(*ast.DatePart); ok && !dateParts[strings.ToUpper(dp.Name)] {
			return qtUnknown, mssqlErr(155, "'%s' is not a recognized %s option.", dp.Name, strings.ToLower(name))
		}
	}
	inner := x
	if agg {
		inner.aggDepth++
	}
	var argTypes []qtype
	if f.Args != nil {
		for _, arg := range f.Args.Items {
			t, err := b.expr(arg, sc, inner)
			if err != nil {
				return qtUnknown, err
			}
			argTypes = append(argTypes, t)
		}
	}
	switch name {
	case "TRIM":
		for i, t := range argTypes {
			if t != qtString && t != qtUnknown {
				return qtUnknown, mssqlErr(8116, "Argument data type %s is invalid for argument %d of Trim function.", t.typeName(), i+1)
			}
		}
	case "SUM", "AVG":
		if len(argTypes) == 1 && argTypes[0] != qtNumber && argTypes[0] != qtUnknown {
			return qtUnknown, mssqlErr(8117, "Operand data type %s is invalid for %s operator.", argTypes[0].typeName(), strings.ToLower(name))
		}
	}
	if _, err := b.expr(f.Within, sc, x); err != nil {
		return qtUnknown, err
	}
	if f.Over != nil {
		if err := b.over(f.Over, name, sc, x); err != nil {
			return qtUnknown, err
		}
	} else if needsOrder, ok := windowFunctions[name]; ok && needsOrder {
		return qtUnknown, mssqlErr(10753, "The function '%s' must have an OVER clause.", strings.ToLower(name))
	}
	switch {
	case stringFunctions[name]:
		return qtString, nil
	case name == "LAG" || name == "LEAD" || name == "FIRST_VALUE" || name == "LAST_VALUE":
		// Navigation functions return their value's type.
		if len(argTypes) > 0 {
			return argTypes[0], nil
		}
		return qtUnknown, nil
	case name == "COUNT" || name == "COUNT_BIG" || name == "LEN" || windowFunctions[name]:
		return qtNumber, nil
	case (name == "SUM" || name == "AVG" || name == "MIN" || name == "MAX") && len(argTypes) == 1:
		return argTypes[0], nil
	}
	return qtUnknown, nil
}

// outerAggregate reports whether an aggregate in a subquery's WHERE belongs
// to an enclosing query: its columns all bind outside the subquery, and the
// subquery is in that query's select list or HAVING (the exception Msg 147
// states).
func (b *binder) outerAggregate(f *ast.FuncCallExpr, sc *qscope, x ectx) bool {
	if len(b.clauses) == 0 {
		return false
	}
	switch b.clauses[len(b.clauses)-1] {
	case clSelect, clHaving, clOrder:
	default:
		return false
	}
	inner := x
	inner.aggDepth++
	if _, err := b.expr(f.Args, sc, inner); err != nil {
		return false
	}
	refs := 0
	outer := true
	ast.Inspect(f.Args, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok {
			refs++
			if bound, ok := b.bound[ref]; !ok || bound.block == x.block {
				outer = false
			}
		}
		return true
	})
	return refs > 0 && outer
}

// over binds an OVER clause; name is the windowed function's.
func (b *binder) over(o *ast.OverClause, name string, sc *qscope, x ectx) error {
	ox := ectx{clause: clOrder, block: x.block}
	if _, err := b.expr(o.PartitionBy, sc, ox); err != nil {
		return err
	}
	if _, err := b.expr(o.OrderBy, sc, ox); err != nil {
		return err
	}
	hasOrder := o.OrderBy != nil && len(o.OrderBy.Items) > 0
	named := o.WindowName != ""
	if named {
		var def *ast.WindowDef
		for s := sc; s != nil && def == nil; s = s.parent {
			if s.block != nil {
				def = s.block.windows[strings.ToUpper(o.WindowName)]
			}
		}
		if def == nil {
			return mssqlErr(5362, "Window '%s' is not defined, or is defined after it is referenced.", o.WindowName)
		}
		for def != nil {
			hasOrder = hasOrder || (def.OrderBy != nil && len(def.OrderBy.Items) > 0)
			if def.RefName == "" || x.block == nil {
				break
			}
			def = x.block.windows[strings.ToUpper(def.RefName)]
		}
	}
	if needsOrder := windowFunctions[name]; needsOrder && !hasOrder && named {
		return mssqlErr(5366, "The function '%s' must have a window specification with ORDER BY.", strings.ToLower(name))
	}
	if needsOrder := windowFunctions[name]; needsOrder && !hasOrder {
		return mssqlErr(4112, "The function '%s' must have an OVER clause with ORDER BY.", strings.ToLower(name))
	}
	return nil
}

// xmlMethods are the xml type's methods a column calls.
var xmlMethods = map[string]qtype{"VALUE": qtUnknown, "QUERY": qtXML, "EXIST": qtNumber, "NODES": qtUnknown, "MODIFY": qtUnknown}

// udfCall binds a qualified call: a method of a column, or a user-defined
// scalar function.
func (b *binder) udfCall(f *ast.FuncCallExpr, sc *qscope, x ectx) (qtype, error) {
	ref := f.Name
	// A method of a column's value: an xml method, or a spatial or other
	// CLR type's, which the catalog does not model.
	if ref.Server == "" {
		target := &ast.ColumnRef{Table: ref.Database, Column: ref.Schema, Loc: ref.Loc}
		if b.resolves(target, sc) {
			return xmlMethods[strings.ToUpper(ref.Object)], b.children(f.Args, sc, x)
		}
	}
	if strings.EqualFold(ref.Schema, "$PARTITION") || !b.localObject(ref) || systemSchema(ref.Schema) {
		return qtUnknown, b.children(f.Args, sc, x)
	}
	missing := mssqlErr(4121, "Cannot find either column \"%s\" or the user-defined function or aggregate \"%s.%s\", or the name is ambiguous.", ref.Schema, ref.Schema, ref.Object)
	schema := b.c.GetSchema(ref.Schema)
	idx, ok := -1, false
	if schema != nil {
		idx, ok = resolveExactOrFoldIndex(schema.Functions, func(fn *metadata.FunctionMetadata) string { return fn.Name }, ref.Object)
	}
	if !ok {
		if b.deferNames {
			return qtUnknown, b.children(f.Args, sc, x)
		}
		return qtUnknown, missing
	}
	fn := schema.Functions[idx]
	info := b.functionInfo(fn)
	if info.category != "" && info.category != "scalar" {
		return qtUnknown, missing
	}
	if b.schemaBound && !isSchemaBound(fn.Definition) {
		return qtUnknown, mssqlErr(4513, "Cannot schema bind view '%s'. '%s.%s' is not schema bound.", b.object, schema.Name, fn.Name)
	}
	if err := b.checkArgCount(f, info, schema.Name, fn.Name); err != nil {
		return qtUnknown, err
	}
	return qtUnknown, b.children(f.Args, sc, x)
}

// resolves reports whether ref binds, without reporting why not.
func (b *binder) resolves(ref *ast.ColumnRef, sc *qscope) bool {
	if ref.Column == "" {
		return false
	}
	_, err := b.column(ref, sc)
	return err == nil
}

// typeClass is the class of a column type as the catalog stores it.
func typeClass(typ string) qtype {
	base := strings.ToLower(strings.TrimSpace(typ))
	if i := strings.IndexByte(base, '('); i >= 0 {
		base = base[:i]
	}
	base = strings.Trim(strings.TrimSpace(base), "[]")
	switch base {
	case "bigint", "int", "integer", "smallint", "tinyint", "bit", "decimal", "dec", "numeric", "money", "smallmoney", "float", "real":
		return qtNumber
	case "char", "varchar", "nchar", "nvarchar", "text", "ntext", "sysname":
		return qtString
	case "date", "time", "datetime", "datetime2", "smalldatetime", "datetimeoffset":
		return qtDate
	case "binary", "varbinary", "image", "timestamp", "rowversion":
		return qtBinary
	case "xml":
		return qtXML
	case "uniqueidentifier", "sql_variant", "hierarchyid", "geometry", "geography", "cursor", "table":
		return qtOther
	}
	// An alias type, or any type the catalog does not resolve, may be any
	// of the above.
	return qtUnknown
}

func dataTypeClass(dt *ast.DataType) qtype {
	if dt == nil || dt.Schema != "" {
		return qtUnknown
	}
	return typeClass(dt.Name)
}
