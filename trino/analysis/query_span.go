package analysis

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// QuerySpan captures the tables and columns a query references. It is the data
// structure Bytebase consumes for field-level lineage (masking) and for SQL
// editor access tracking (read/write permission checks). It is the omni
// counterpart of the legacy plugin/parser/trino GetQuerySpan result.
//
// The extraction is best-effort and parser-driven: it walks the SELECT / set-op
// / CTE / FROM / expression tree the trino/parser produces. Known best-effort
// boundaries (each a deliberate scope decision, matching the doris/analysis
// precedent — full resolution belongs to a later metadata-aware layer):
//   - `SELECT *` is not expanded against catalog metadata (needs a
//     DatabaseMetadata the completion/bytebase layer supplies); a star item is
//     recorded by name only ("*").
//   - A column reference is recorded by its written (catalog/schema/)table/
//     column parts; it is not resolved to which FROM relation actually owns it.
//   - A named-window reference (`OVER w`) is not resolved to window w's
//     definition, so its partition/order columns are captured as predicate
//     columns (via the WINDOW clause walk) but not attributed to the specific
//     OVER-w result column.
type QuerySpan struct {
	// Type is the statement classification (Select, DML, DDL, ...).
	Type QueryType

	// AccessTables is the set of base tables the query reads from: every
	// physical table reference in FROM clauses, JOINs, and subqueries. Entries
	// are deduplicated on (catalog, schema, table, alias); CTE references and
	// derived-table aliases are excluded.
	AccessTables []TableAccess

	// Results is the list of output columns produced by the outermost query. For
	// a set operation (UNION/INTERSECT/EXCEPT) the left arm wins, matching SQL's
	// "column names come from the first select" rule.
	Results []ColumnInfo

	// PredicateColumns is the set of columns referenced in filter/join positions
	// (WHERE, JOIN ON, HAVING, QUALIFY). Bytebase uses these for row-level
	// masking policy evaluation. Deduplicated.
	PredicateColumns []ColumnRef

	// CTEs lists the names defined by WITH clauses at any scope, in declaration
	// order — useful for consumers and for debugging the CTE structure.
	CTEs []string
}

// TableAccess represents one physical table referenced in a query. In Trino a
// fully-qualified name is catalog.schema.table; the optional qualifiers are
// empty when the source name omitted them (resolution against the session's
// current catalog/schema is the engine's job, not the parser's).
type TableAccess struct {
	Catalog string  // optional catalog qualifier (empty when none)
	Schema  string  // optional schema qualifier (empty when none)
	Table   string  // table (or view) name
	Alias   string  // alias at the reference site (empty when none)
	Loc     ast.Loc // source location of the table reference
}

// ColumnInfo represents one output column of a query result.
type ColumnInfo struct {
	// Name is the output column name: the explicit alias if present, otherwise a
	// best-effort rendering of the expression (the column name for a bare
	// column reference, "*" for a star item, or "" when none can be derived).
	Name string

	// SourceColumns is the best-effort list of column references that directly
	// feed this output column (the column refs in the select-item expression,
	// excluding those inside nested subqueries).
	SourceColumns []ColumnRef
}

// ColumnRef identifies a column by its optional catalog/schema/table qualifier
// and its name. Trino column references are at most table.column inside a query
// (catalog/schema qualification on a column is rare but legal via a longer
// dotted chain); the rightmost component is always the column.
type ColumnRef struct {
	Catalog string
	Schema  string
	Table   string
	Column  string
}

// GetQuerySpan analyzes a single Trino SQL statement and returns its query
// span: the statement classification, the base tables it reads, the CTE names
// it defines, its output columns with their source columns, and the columns
// used in predicate positions.
//
// It is tolerant of parse errors — if the parser produces a partial AST,
// whatever parsed is still analyzed. On empty input it returns a zero-valued
// span with Type=Unknown.
func GetQuerySpan(statement string) (*QuerySpan, error) {
	file, _ := parser.Parse(statement)
	span := &QuerySpan{Type: Classify(statement)}
	if file == nil || len(file.Stmts) == 0 {
		return span, nil
	}

	w := newSpanWalker(span)
	for _, stmt := range file.Stmts {
		w.analyzeStmt(stmt)
	}

	// Deepen result-column lineage through derived relations (subqueries in
	// FROM and CTE references): the primary walk records a derived column by the
	// name written at the reference site, which has no base table to mask
	// against. This rewrites those refs to the recovered base columns, leaving
	// direct base-table references untouched.
	if len(file.Stmts) > 0 {
		resolveDerivedLineage(file.Stmts[0], span)
	}
	return span, nil
}

// ---------------------------------------------------------------------------
// CTE scope
// ---------------------------------------------------------------------------

// cteScope is a linked stack of CTE name sets. Resolving a bare table reference
// walks outwards; an inner CTE shadows an outer one of the same name.
type cteScope struct {
	names  map[string]bool
	parent *cteScope
}

func (s *cteScope) isCTE(name string) bool {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.names[strings.ToLower(name)] {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Walker
// ---------------------------------------------------------------------------

// spanWalker walks the query tree the trino/parser produces (parser-package
// node types — Query / QueryNode / Relation / Expr — not ast.Node, so it cannot
// reuse ast.Walk). It maintains:
//   - a CTE scope stack so bare table references that name a CTE are filtered
//     out of AccessTables;
//   - a dedup map of (catalog, schema, table, alias) for AccessTables;
//   - a dedup map for PredicateColumns;
//   - a flag tracking whether the outermost query has populated Results yet.
type spanWalker struct {
	span     *QuerySpan
	scope    *cteScope
	accessed map[tableKey]bool
	predSeen map[ColumnRef]bool
	resolved bool // Results captured for the outermost query
}

type tableKey struct {
	Catalog string
	Schema  string
	Table   string
	Alias   string
}

func newSpanWalker(span *QuerySpan) *spanWalker {
	return &spanWalker{
		span:     span,
		accessed: make(map[tableKey]bool),
		predSeen: make(map[ColumnRef]bool),
	}
}

// analyzeStmt dispatches on the top-level statement node. Only the query
// statement (QueryStmt) yields query-span data; other statement kinds are
// recorded via their Type classification only.
func (w *spanWalker) analyzeStmt(node ast.Node) {
	if qs, ok := node.(*parser.QueryStmt); ok && qs.Query != nil {
		w.visitQuery(qs.Query, true /* outermost */)
	}
}

// visitQuery walks a `Query` (optional WITH + a queryNoWith body with its own
// ORDER BY/OFFSET/LIMIT/FETCH). The WITH clause defines a new CTE scope visible
// to the body and to sibling CTE bodies.
func (w *spanWalker) visitQuery(q *parser.Query, outermost bool) {
	if q == nil {
		return
	}

	scope := &cteScope{names: make(map[string]bool), parent: w.scope}
	if q.With != nil {
		for _, cte := range q.With.CTEs {
			name := identName(cte.Name)
			if name == "" {
				continue
			}
			scope.names[strings.ToLower(name)] = true
			w.span.CTEs = append(w.span.CTEs, name)
		}
	}
	w.scope = scope

	// Walk CTE bodies first so their base tables land in AccessTables. They run
	// in the scope that already includes the sibling CTE names (matching the
	// legacy extractor's permissive visibility).
	if q.With != nil {
		for _, cte := range q.With.CTEs {
			w.visitQuery(cte.Query, false)
		}
	}

	// The query body (a queryNoWith) carries the result columns when outermost.
	w.visitQueryNode(q.Body, outermost)

	// queryNoWith-level ORDER BY keys may reference columns; walk them as
	// predicate-ish expressions (collecting column refs, recursing subqueries).
	for _, item := range q.OrderBy {
		w.walkExprPredicate(item.Expr)
	}
	if q.Offset != nil {
		w.walkExprPredicate(q.Offset.Count)
	}
	if q.Limit != nil {
		w.walkExprPredicate(q.Limit.Count)
	}
	if q.Fetch != nil {
		w.walkExprPredicate(q.Fetch.Count)
	}

	w.scope = scope.parent
}

// visitQueryNode dispatches on a query body / set-op operand / query primary.
func (w *spanWalker) visitQueryNode(node parser.QueryNode, outermost bool) {
	switch n := node.(type) {
	case *parser.QuerySpec:
		w.visitQuerySpec(n, outermost)
	case *parser.SetOperation:
		w.visitSetOp(n, outermost)
	case *parser.ParenQuery:
		if n.Inner != nil {
			w.visitQuery(n.Inner, outermost)
		}
	case *parser.TableQuery:
		// TABLE name == SELECT * FROM name: record the table and a single star
		// result column when this is the outermost query.
		w.recordTable(n.Name, "")
		if outermost && !w.resolved {
			w.span.Results = append(w.span.Results, ColumnInfo{Name: "*"})
			w.resolved = true
		}
	case *parser.ValuesQuery:
		for _, row := range n.Rows {
			w.walkExprPredicate(row)
		}
	}
}

// visitSetOp walks a UNION/INTERSECT/EXCEPT tree. The left arm surfaces its
// Results when this is the outermost statement (SQL's first-select rule); the
// right arm never does.
func (w *spanWalker) visitSetOp(n *parser.SetOperation, outermost bool) {
	if n == nil {
		return
	}
	w.visitQueryNode(n.Left, outermost)
	w.visitQueryNode(n.Right, false)
}

// visitQuerySpec processes one querySpecification (a SELECT block): its FROM
// relations, predicate clauses, and select list.
func (w *spanWalker) visitQuerySpec(spec *parser.QuerySpec, outermost bool) {
	if spec == nil {
		return
	}

	for _, rel := range spec.From {
		w.visitRelation(rel)
	}

	if spec.Where != nil {
		w.walkExprPredicate(spec.Where)
	}
	if spec.GroupBy != nil {
		for _, el := range spec.GroupBy.Elements {
			for _, e := range el.Exprs {
				w.walkExprPredicate(e)
			}
			for _, set := range el.Sets {
				for _, e := range set {
					w.walkExprPredicate(e)
				}
			}
		}
	}
	if spec.Having != nil {
		w.walkExprPredicate(spec.Having)
	}
	for _, win := range spec.Windows {
		ew := exprWalk{w: w, followSub: true, onColumn: w.addPredicateColumn}
		ew.walkWindowSpec(win.Spec)
	}

	// SELECT list last so Results reflects final column order.
	for _, item := range spec.Items {
		if outermost && !w.resolved {
			w.span.Results = append(w.span.Results, w.makeColumnInfo(item))
		}
		// Walk the item expression to pick up subquery tables (the direct
		// column refs are captured in makeColumnInfo; here we only need the
		// subqueries and any deeper table references).
		w.walkExprSubqueries(item.Expr)
	}
	if outermost {
		w.resolved = true
	}
}

// ---------------------------------------------------------------------------
// Relations (FROM)
// ---------------------------------------------------------------------------

// visitRelation walks one FROM relation, threading an optional alias down only
// to a direct base-table primary.
func (w *spanWalker) visitRelation(rel parser.Relation) {
	switch n := rel.(type) {
	case *parser.AliasedRelation:
		w.visitAliasedRelation(n)
	case *parser.Join:
		w.visitRelation(n.Left)
		w.visitRelation(n.Right)
		w.walkExprPredicate(n.On)
		// USING (col, …) names join-key columns — record them as predicate
		// columns (they are unqualified; resolution to a side is the engine's).
		for _, id := range n.Using {
			if name := identName(id); name != "" {
				w.addPredicateColumn(ColumnRef{Column: name})
			}
		}
	case *parser.ParenRelation:
		w.visitRelation(n.Inner)
	case *parser.TableRelation:
		w.recordTable(n.Name, "")
	case *parser.SubqueryRelation:
		w.visitQuery(n.Query, false)
	case *parser.LateralRelation:
		w.visitQuery(n.Query, false)
	case *parser.UnnestRelation:
		for _, e := range n.Exprs {
			w.walkExprPredicate(e)
		}
	case *parser.TableFunctionRelation:
		w.visitTableFunctionCall(n.Call)
	}
}

// visitAliasedRelation handles a relationPrimary wrapped in its (optional)
// alias. The alias attaches to a direct base table only; for a derived table
// (subquery / lateral / unnest / table function) the alias names the derived
// relation, not its underlying base tables, so it is not propagated.
func (w *spanWalker) visitAliasedRelation(n *parser.AliasedRelation) {
	if n == nil {
		return
	}
	alias := identName(n.Alias)
	switch inner := n.Inner.(type) {
	case *parser.TableRelation:
		w.recordTable(inner.Name, alias)
	default:
		w.visitRelation(n.Inner)
	}
}

// visitTableFunctionCall walks a table-function invocation's table arguments,
// which may themselves name base tables or wrap subqueries.
func (w *spanWalker) visitTableFunctionCall(call *parser.TableFunctionCall) {
	if call == nil {
		return
	}
	for _, arg := range call.Args {
		switch arg.Kind {
		case parser.TFArgExpr:
			w.walkExprPredicate(arg.Expr)
		case parser.TFArgTable:
			if arg.Table != nil {
				if arg.Table.Query != nil {
					w.visitQuery(arg.Table.Query, false)
				} else if arg.Table.Name != nil {
					w.recordTable(arg.Table.Name, "")
				}
				for _, e := range arg.Table.PartitionBy {
					w.walkExprPredicate(e)
				}
				for _, item := range arg.Table.OrderBy {
					w.walkExprPredicate(item.Expr)
				}
			}
		}
	}
}

// recordTable adds a base-table reference to AccessTables unless it names an
// in-scope CTE or is a duplicate. The qualified name's components map to
// (catalog, schema, table): a 1-part name is the table; 2-part is schema.table;
// 3-part (or longer) is catalog.schema.table on the rightmost three parts.
func (w *spanWalker) recordTable(name *ast.QualifiedName, alias string) {
	parts := normalizedParts(name)
	if len(parts) == 0 {
		return
	}

	var catalog, schema, table string
	switch len(parts) {
	case 1:
		table = parts[0]
		// A bare (single-part) name that matches an in-scope CTE is a CTE
		// reference, not a base table — exclude it regardless of any alias
		// (`FROM c` and `FROM c AS x` both reference the CTE c).
		if w.scope.isCTE(table) {
			return
		}
	case 2:
		schema = parts[0]
		table = parts[1]
	default:
		catalog = parts[len(parts)-3]
		schema = parts[len(parts)-2]
		table = parts[len(parts)-1]
	}

	key := tableKey{Catalog: catalog, Schema: schema, Table: table, Alias: alias}
	if w.accessed[key] {
		return
	}
	w.accessed[key] = true

	loc := ast.Loc{}
	if name != nil {
		loc = name.Loc
	}
	w.span.AccessTables = append(w.span.AccessTables, TableAccess{
		Catalog: catalog,
		Schema:  schema,
		Table:   table,
		Alias:   alias,
		Loc:     loc,
	})
}

// ---------------------------------------------------------------------------
// Expression walking
// ---------------------------------------------------------------------------

// exprWalk parameterizes a single expression traversal so the walker code is
// written once and reused for the two concerns it serves:
//   - onColumn is invoked for each column reference encountered (nil to ignore
//     columns — used when a pass only needs to discover subquery tables).
//   - followSub controls whether subquery placeholders are recursed into by
//     re-parsing their raw text (true for predicate/table discovery; false when
//     collecting a select item's DIRECT columns, which must not cross a
//     subquery boundary).
//
// Threading these through one walker (rather than two near-identical type
// switches) keeps the ~40-arm expr dispatch in a single place, so a new
// expression node type cannot be handled in one pass but silently dropped in
// the other.
type exprWalk struct {
	w         *spanWalker
	onColumn  func(ColumnRef)
	followSub bool
	// onSubquery, when non-nil, is invoked for each scalar subquery placeholder
	// (a bare `( query )`) instead of followSubquery. The lineage resolver uses
	// it to recover a scalar subquery's output-column sources; the walker's own
	// passes leave it nil and keep the followSub/followSubquery behaviour.
	onSubquery func(*parser.SubqueryExpr)
}

// addPredicateColumn appends a column reference to PredicateColumns,
// deduplicating against earlier predicate columns.
func (w *spanWalker) addPredicateColumn(ref ColumnRef) {
	if w.predSeen[ref] {
		return
	}
	w.predSeen[ref] = true
	w.span.PredicateColumns = append(w.span.PredicateColumns, ref)
}

// walkExprPredicate walks a predicate-position expression (WHERE / ON / HAVING /
// GROUP BY / ORDER BY / VALUES), collecting its column references into
// PredicateColumns and recursing into any subqueries (which may reference more
// base tables).
func (w *spanWalker) walkExprPredicate(expr parser.Expr) {
	ew := exprWalk{
		w:         w,
		followSub: true,
		onColumn:  w.addPredicateColumn,
	}
	ew.walk(expr)
}

// walkExprSubqueries walks a select-item expression purely to discover the base
// tables inside nested subqueries; the item's direct column refs are captured
// separately by makeColumnInfo, so column refs are ignored here.
func (w *spanWalker) walkExprSubqueries(expr parser.Expr) {
	ew := exprWalk{w: w, followSub: true}
	ew.walk(expr)
}

// walk is the shared expression traversal. It descends every expression child
// of the trino expr model, invoking onColumn for each column reference (a
// ColumnRef or a Dereference chain), and — when followSub is set — recursing
// into subquery placeholders by re-parsing their raw text. Expression node
// types that cannot contain a column or subquery are leaves.
func (ew exprWalk) walk(expr parser.Expr) {
	onColumn := ew.onColumn
	switch e := expr.(type) {
	case nil:
		return

	// Column references.
	case *parser.ColumnRef:
		if onColumn != nil {
			if ref, ok := columnRefFrom(e); ok {
				onColumn(ref)
			}
		}
	case *parser.Dereference:
		// A dotted name (t.a / cat.sch.t.col) is a Dereference chain over a
		// ColumnRef. Flatten it into a single qualified ColumnRef; if the base
		// is not a plain name chain (e.g. field access on a function result),
		// fall back to walking the base expression.
		if ref, ok := dereferenceRef(e); ok {
			if onColumn != nil {
				onColumn(ref)
			}
		} else {
			ew.walk(e.Base)
		}

	// Subquery placeholders — re-parse the raw inner text and recurse, but only
	// when this pass follows subqueries (a select item's direct-column pass does
	// not cross the subquery boundary).
	case *parser.SubqueryExpr:
		if ew.onSubquery != nil {
			ew.onSubquery(e)
		} else {
			ew.followSubquery(e)
		}

	// Composite expressions — recurse into children.
	case *parser.Subscript:
		ew.walk(e.Value)
		ew.walk(e.Index)
	case *parser.UnaryExpr:
		ew.walk(e.Operand)
	case *parser.BinaryExpr:
		ew.walk(e.Left)
		ew.walk(e.Right)
	case *parser.LogicalExpr:
		ew.walk(e.Left)
		ew.walk(e.Right)
	case *parser.NotExpr:
		ew.walk(e.Operand)
	case *parser.AtTimeZoneExpr:
		ew.walk(e.Value)
		ew.walk(e.Zone)
	case *parser.ParenExpr:
		ew.walk(e.Expr)
	case *parser.RowConstructor:
		ew.walkList(e.Elements)
	case *parser.ArrayConstructor:
		ew.walkList(e.Elements)
	case *parser.ComparisonExpr:
		ew.walk(e.Left)
		ew.walk(e.Right)
	case *parser.QuantifiedComparisonExpr:
		ew.walk(e.Left)
		ew.followSubquery(e.Subquery)
	case *parser.BetweenExpr:
		ew.walk(e.Value)
		ew.walk(e.Lower)
		ew.walk(e.Upper)
	case *parser.InListExpr:
		ew.walk(e.Value)
		ew.walkList(e.List)
	case *parser.InSubqueryExpr:
		ew.walk(e.Value)
		ew.followSubquery(e.Subquery)
	case *parser.LikeExpr:
		ew.walk(e.Value)
		ew.walk(e.Pattern)
		ew.walk(e.Escape)
	case *parser.IsNullExpr:
		ew.walk(e.Value)
	case *parser.IsDistinctFromExpr:
		ew.walk(e.Value)
		ew.walk(e.Right)
	case *parser.FuncCall:
		ew.walkList(e.Args)
		for _, item := range e.OrderBy {
			ew.walk(item.Expr)
		}
		ew.walk(e.Filter)
		ew.walkWindowSpec(e.Over)
	case *parser.LambdaExpr:
		// Lambda parameters (the `x` in `x -> …`) are bound variables, not table
		// columns. Walk the body with an onColumn that drops any reference whose
		// unqualified name matches a parameter, so free (outer) columns still
		// surface but the bound params do not pollute lineage.
		ew.walkLambdaBody(e)
	case *parser.CaseExpr:
		ew.walk(e.Operand)
		for _, when := range e.Whens {
			ew.walk(when.Cond)
			ew.walk(when.Result)
		}
		ew.walk(e.Else)
	case *parser.CastExpr:
		ew.walk(e.Expr)
	case *parser.ExtractExpr:
		ew.walk(e.Source)
	case *parser.SubstringExpr:
		ew.walk(e.Source)
		ew.walk(e.From)
		ew.walk(e.For)
	case *parser.TrimExpr:
		ew.walk(e.Char)
		ew.walk(e.Source)
	case *parser.NormalizeExpr:
		ew.walk(e.Source)
	case *parser.PositionExpr:
		ew.walk(e.Needle)
		ew.walk(e.Haystack)
	case *parser.ListaggExpr:
		ew.walk(e.Arg)
		ew.walk(e.Separator)
		for _, item := range e.WithinGroupBy {
			ew.walk(item.Expr)
		}
		ew.walk(e.Filter)
	case *parser.GroupingExpr:
		// GROUPING(col, …) references each named column.
		if onColumn != nil {
			for _, qn := range e.Args {
				if parts := normalizedParts(qn); len(parts) > 0 {
					onColumn(columnRefFromParts(parts))
				}
			}
		}

	// SQL/JSON functions — the path input and arguments are column positions.
	case *parser.JSONExistsExpr:
		ew.walkJSONPath(e.Path)
	case *parser.JSONValueFunc:
		ew.walkJSONPath(e.Path)
		if e.OnEmpty != nil {
			ew.walk(e.OnEmpty.Default)
		}
		if e.OnError != nil {
			ew.walk(e.OnError.Default)
		}
	case *parser.JSONQueryFunc:
		ew.walkJSONPath(e.Path)
	case *parser.JSONObjectExpr:
		for _, m := range e.Members {
			ew.walk(m.Key)
			ew.walkJSONValueExpr(m.Value)
		}
	case *parser.JSONArrayExpr:
		for _, el := range e.Elements {
			ew.walkJSONValueExpr(el)
		}

	// Leaves (literals, parameters, special date/time funcs, interval, typed
	// literals) carry no column or subquery and need no recursion.
	default:
		return
	}
}

// followSubquery recurses into a subquery placeholder by re-parsing its raw
// text, but only when this pass follows subqueries. A nil placeholder is a
// no-op.
func (ew exprWalk) followSubquery(sub *parser.SubqueryExpr) {
	if !ew.followSub || sub == nil {
		return
	}
	ew.w.analyzeSubqueryText(sub.RawText)
}

// walkLambdaBody walks a lambda's body with the bound parameter names excluded
// from column collection. A lambda parameter is an unqualified single
// identifier; a body reference to that bare name is the bound variable, not a
// table column, so it is dropped — while qualified references and free
// (outer-scope) column names still pass through to onColumn.
func (ew exprWalk) walkLambdaBody(lambda *parser.LambdaExpr) {
	if lambda == nil {
		return
	}
	if ew.onColumn == nil {
		ew.walk(lambda.Body)
		return
	}
	params := make(map[string]bool, len(lambda.Params))
	for _, p := range lambda.Params {
		if name := identName(p); name != "" {
			params[name] = true
		}
	}
	inner := ew
	outer := ew.onColumn
	inner.onColumn = func(ref ColumnRef) {
		// Drop a reference whose ROOT name is a bound parameter: either a bare
		// `x` (param used as a value) or a field access `x.field` (which flattens
		// to a ColumnRef whose leftmost qualifier is the param). Both are uses of
		// the bound variable, not table columns.
		if params[rootName(ref)] {
			return
		}
		outer(ref)
	}
	inner.walk(lambda.Body)
}

// rootName returns the leftmost (most-qualifying) name of a column reference:
// the catalog if present, else the schema, else the table, else the column. For
// a flattened dereference like x.price (ColumnRef{Table:"x", Column:"price"})
// the root is "x"; for a bare column it is the column name itself.
func rootName(ref ColumnRef) string {
	switch {
	case ref.Catalog != "":
		return ref.Catalog
	case ref.Schema != "":
		return ref.Schema
	case ref.Table != "":
		return ref.Table
	default:
		return ref.Column
	}
}

// walkJSONPath walks a SQL/JSON path invocation's input expression and PASSING
// arguments for column references (the path string itself is opaque).
func (ew exprWalk) walkJSONPath(path *parser.JSONPathInvocation) {
	if path == nil {
		return
	}
	ew.walkJSONValueExpr(path.Input)
	for _, arg := range path.Passing {
		ew.walkJSONValueExpr(arg.Value)
	}
}

// walkJSONValueExpr walks a `expression [FORMAT …]` JSON value expression.
func (ew exprWalk) walkJSONValueExpr(v *parser.JSONValueExpr) {
	if v == nil {
		return
	}
	ew.walk(v.Expr)
}

// walkList walks each expression in a slice.
func (ew exprWalk) walkList(exprs []parser.Expr) {
	for _, e := range exprs {
		ew.walk(e)
	}
}

// walkWindowSpec walks an inline/named window specification's PARTITION BY and
// ORDER BY expressions and frame bounds for column/subquery references, reusing
// this pass's onColumn/followSub settings.
func (ew exprWalk) walkWindowSpec(spec *parser.WindowSpec) {
	if spec == nil {
		return
	}
	ew.walkList(spec.PartitionBy)
	for _, item := range spec.OrderBy {
		ew.walk(item.Expr)
	}
	if spec.Frame != nil {
		ew.walk(spec.Frame.Start.Value)
		if spec.Frame.End != nil {
			ew.walk(spec.Frame.End.Value)
		}
	}
}

// analyzeSubqueryText re-parses an expression-embedded subquery's raw inner text
// and recurses into the resulting query. Errors are swallowed — an unparseable
// subquery leaves the already-discovered tables intact (best-effort).
func (w *spanWalker) analyzeSubqueryText(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	file, _ := parser.Parse(text)
	if file == nil {
		return
	}
	for _, stmt := range file.Stmts {
		if qs, ok := stmt.(*parser.QueryStmt); ok {
			w.visitQuery(qs.Query, false)
		}
	}
}

// ---------------------------------------------------------------------------
// Select-item rendering
// ---------------------------------------------------------------------------

// makeColumnInfo derives a ColumnInfo from a select-item. Name is the alias if
// present, else a simple rendering ("*"/"t.*" for a star, the column name for a
// bare column reference, "" otherwise). SourceColumns is the set of column refs
// the item's expression directly references (excluding nested subqueries).
func (w *spanWalker) makeColumnInfo(item parser.SelectItem) ColumnInfo {
	info := ColumnInfo{}
	switch item.Kind {
	case parser.SelectAll:
		info.Name = "*"
		return info
	case parser.SelectAllFrom:
		// expr.* — name it "<expr>.*" best-effort; the underlying columns are
		// the column refs in the row expression.
		if name := renderExprName(item.Expr); name != "" {
			info.Name = name + ".*"
		} else {
			info.Name = "*"
		}
		info.SourceColumns = w.collectDirectColumns(item.Expr)
		return info
	}

	if name := identName(item.Alias); name != "" {
		info.Name = name
	} else {
		info.Name = renderExprName(item.Expr)
	}
	info.SourceColumns = w.collectDirectColumns(item.Expr)
	return info
}

// collectDirectColumns returns the column references directly mentioned in expr,
// in source order, excluding any inside nested subqueries (which are followed
// for table discovery elsewhere, not for this item's source columns). It is the
// same traversal as the predicate pass with followSub disabled.
func (w *spanWalker) collectDirectColumns(expr parser.Expr) []ColumnRef {
	var refs []ColumnRef
	ew := exprWalk{
		w:         w,
		followSub: false,
		onColumn:  func(ref ColumnRef) { refs = append(refs, ref) },
	}
	ew.walk(expr)
	return refs
}

// renderExprName returns a short human-readable name for a select-item
// expression: the column name for a bare column reference (or the trailing
// field of a dotted reference), unwrapping a single layer of parentheses. Other
// expressions have no derivable name ("").
func renderExprName(expr parser.Expr) string {
	switch e := expr.(type) {
	case *parser.ColumnRef:
		return identName(e.Name)
	case *parser.Dereference:
		return identName(e.FieldName)
	case *parser.ParenExpr:
		return renderExprName(e.Expr)
	}
	return ""
}

// ---------------------------------------------------------------------------
// Name helpers
// ---------------------------------------------------------------------------

// identName returns the normalized (Trino case-folded) name of an identifier,
// or "" for a nil identifier.
func identName(id *ast.Identifier) string {
	if id == nil {
		return ""
	}
	return id.Normalize()
}

// normalizedParts returns the normalized component names of a qualified name,
// skipping nil components.
func normalizedParts(name *ast.QualifiedName) []string {
	if name == nil {
		return nil
	}
	return name.NormalizedParts()
}

// columnRefFrom builds a ColumnRef from a bare ColumnRef node (a single
// identifier). Returns ok=false for a nil node or nil/empty name.
func columnRefFrom(c *parser.ColumnRef) (ColumnRef, bool) {
	if c == nil {
		return ColumnRef{}, false
	}
	name := identName(c.Name)
	if name == "" {
		return ColumnRef{}, false
	}
	return ColumnRef{Column: name}, true
}

// dereferenceRef flattens a Dereference chain that bottoms out at a ColumnRef
// into a single qualified ColumnRef (the rightmost component is the column; the
// preceding ones map to table / schema / catalog). It returns ok=false when the
// base is not a plain name chain (e.g. a subscript or function result), so the
// caller can fall back to walking the base expression.
func dereferenceRef(d *parser.Dereference) (ColumnRef, bool) {
	parts := flattenDereference(d)
	if parts == nil {
		return ColumnRef{}, false
	}
	return columnRefFromParts(parts), true
}

// flattenDereference returns the dotted-name components of a Dereference chain
// rooted at a ColumnRef (e.g. ["t", "a"] for t.a), or nil when the chain's base
// is not a ColumnRef.
func flattenDereference(d *parser.Dereference) []string {
	if d == nil {
		return nil
	}
	field := identName(d.FieldName)
	if field == "" {
		return nil
	}
	switch base := d.Base.(type) {
	case *parser.ColumnRef:
		root := identName(base.Name)
		if root == "" {
			return nil
		}
		return []string{root, field}
	case *parser.Dereference:
		prefix := flattenDereference(base)
		if prefix == nil {
			return nil
		}
		return append(prefix, field)
	default:
		return nil
	}
}

// columnRefFromParts maps the components of a dotted column reference to a
// ColumnRef. The rightmost part is always the column; the (up to three)
// preceding parts are table, schema, catalog from right to left.
func columnRefFromParts(parts []string) ColumnRef {
	n := len(parts)
	ref := ColumnRef{Column: parts[n-1]}
	if n >= 2 {
		ref.Table = parts[n-2]
	}
	if n >= 3 {
		ref.Schema = parts[n-3]
	}
	if n >= 4 {
		ref.Catalog = parts[n-4]
	}
	return ref
}
