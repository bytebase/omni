package analysis

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
)

// QuerySpan captures the tables and columns a query references. It is the
// data structure consumed by Bytebase for field-level lineage and for SQL
// editor access-tracking (read/write permission checks).
type QuerySpan struct {
	// Type is the classification of the query (SELECT, DML, DDL, etc.).
	Type QueryType

	// AccessTables is the set of tables the query reads from: every physical
	// table reference in the FROM clause, JOINs, and subqueries. Entries are
	// deduplicated; CTE references (and subquery-derived aliases) are excluded.
	AccessTables []TableAccess

	// Results is the list of output columns produced by the outermost SELECT.
	// For set operations (UNION/INTERSECT/EXCEPT), the left side wins.
	Results []ColumnInfo

	// CTEs tracks the names defined by WITH clauses at any scope — useful for
	// debugging and for consumers that want to show the CTE structure.
	CTEs []string
}

// TableAccess represents a single physical table referenced in a query.
type TableAccess struct {
	Database string  // optional database qualifier (empty when none)
	Table    string  // table name
	Alias    string  // optional alias at the reference site (empty when none)
	Loc      ast.Loc // source location of the table reference
}

// ColumnInfo represents one output column from a SELECT list.
type ColumnInfo struct {
	// Name is the output column name: the alias if one is present, otherwise a
	// best-effort rendering of the expression (column name for a simple
	// ColumnRef, "*" for star items, or empty when we cannot produce one).
	Name string

	// SourceColumns is a best-effort list of underlying column references that
	// feed this output column. For the current scope we record every ColumnRef
	// we encounter while walking the expression tree; full field-level lineage
	// is out of scope for T3.1.
	SourceColumns []ColumnRef
}

// ColumnRef identifies a column by its (optional) database/table qualifier
// and column name.
type ColumnRef struct {
	Database string
	Table    string
	Column   string
}

// GetQuerySpan analyzes a SQL statement and returns its query span: the set
// of tables it reads from, the CTE names it defines, and its output columns.
//
// GetQuerySpan fails closed on parse errors. Masking and access checks
// consume the span, and a partial AST understates what the statement reads —
// before the strict Parse (BYT-10085), `SELECT secret_col[1] FROM
// sensitive_table` analyzed as a table-less SELECT. A statement (or any of
// its subqueries) that does not fully parse therefore yields an error, never
// a silently smaller span. On empty input it returns a zero-valued span with
// Type=QueryTypeUnknown.
func GetQuerySpan(statement string) (*QuerySpan, error) {
	file, errs := parser.Parse(statement)
	if len(errs) > 0 {
		return nil, &errs[0]
	}
	span := &QuerySpan{
		Type: Classify(statement),
	}
	if file == nil || len(file.Stmts) == 0 {
		return span, nil
	}

	w := newSpanWalker(span)
	for _, stmt := range file.Stmts {
		w.analyzeStmt(stmt)
	}
	w.finalize()
	if w.parseErr != nil {
		return nil, w.parseErr
	}
	return span, nil
}

// ---------------------------------------------------------------------------
// Walker
// ---------------------------------------------------------------------------

// cteScope is a linked stack of CTE name sets. When resolving a bare table
// reference we walk outwards — an inner CTE shadows an outer one of the same
// name.
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

// spanWalker performs the SELECT-tree walk that populates a QuerySpan.
//
// The walker maintains:
//   - a scope stack of CTE names so bare table references can be filtered out
//   - a deduplication map of (database, table) pairs for AccessTables
//   - a flag indicating whether the outermost SELECT has populated Results yet
type spanWalker struct {
	// parseErr records the first subquery parse failure; GetQuerySpan fails
	// closed on it rather than returning a span missing that subquery's reads.
	parseErr error

	// textBase is the absolute offset, in the original statement, of the
	// subquery text currently being walked (0 while walking the statement
	// itself). Each nesting level's TextStart is relative to its own
	// extracted text, so bases accumulate as the walker descends.
	textBase int

	span     *QuerySpan
	scope    *cteScope
	accessed map[tableKey]int // maps key -> index in span.AccessTables
	resolved bool             // Results have been captured for the outermost SELECT
}

type tableKey struct {
	Database string
	Table    string
	Alias    string
}

func newSpanWalker(span *QuerySpan) *spanWalker {
	return &spanWalker{
		span:     span,
		accessed: make(map[tableKey]int),
	}
}

// analyzeStmt dispatches on the top-level statement kind. Only SELECT /
// set-op trees produce meaningful query-span data; other statement kinds are
// recorded via their QueryType classification only.
func (w *spanWalker) analyzeStmt(node ast.Node) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		w.visitSelect(n, true /* outermost */)
	case *ast.SetOpStmt:
		w.visitSetOp(n, true /* outermost */)
	case *ast.GroupedQuery:
		w.visitGroupedQuery(n, true /* outermost */)
	default:
		w.validateEmbeddedSubqueries(node)
	}
}

// visitGroupedQuery analyzes a repeated trailing-clause group: the inner
// query first, then the group's own expressions — both clause layers carry
// table reads that must land in AccessTables.
func (w *spanWalker) visitGroupedQuery(n *ast.GroupedQuery, outermost bool) {
	if n == nil {
		return
	}
	switch q := n.Query.(type) {
	case *ast.SelectStmt:
		w.visitSelect(q, outermost)
	case *ast.SetOpStmt:
		w.visitSetOp(q, outermost)
	case *ast.GroupedQuery:
		w.visitGroupedQuery(q, outermost)
	}
	for _, o := range n.OrderBy {
		if o != nil && o.Expr != nil {
			w.walkExpr(o.Expr)
		}
	}
	if n.Limit != nil {
		w.walkExpr(n.Limit)
	}
	if n.Offset != nil {
		w.walkExpr(n.Offset)
	}
}

// noteNonQuerySubquery records the fail-closed error for a subquery
// placeholder whose body is not a query: empty (EXISTS ()), comment-only, or
// a cleanly-parsing non-query statement (EXISTS (DELETE ...)). loc is in
// outer-statement coordinates.
func (w *spanWalker) noteNonQuerySubquery(loc ast.Loc) {
	if w.parseErr == nil {
		w.parseErr = &parser.ParseError{
			Loc: loc,
			Msg: "subquery must be a SELECT statement",
		}
	}
}

// validateEmbeddedSubqueries walks a statement that produces no lineage of
// its own (DML, most DDL) and analyzes every embedded raw subquery, so the
// fail-closed contract holds there too: in
// DELETE FROM t WHERE id IN (SELECT <malformed>) the malformed subquery must
// surface as an error instead of silently returning an empty span, and a
// well-formed one contributes its table reads to AccessTables.
func (w *spanWalker) validateEmbeddedSubqueries(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch q := n.(type) {
		case *ast.SubqueryExpr:
			w.analyzeSubqueryText(q.RawText, q.TextStart)
			return false
		case *ast.RawQuery:
			// CREATE TABLE dest AS SELECT * FROM secret keeps its query as
			// raw text; the read it performs must reach AccessTables, and a
			// malformed body must fail the span like any other subquery.
			w.analyzeSubqueryText(q.RawText, q.TextStart)
			return false
		case *ast.SelectStmt:
			// A parsed query child (INSERT INTO dest SELECT * FROM secret)
			// carries real table reads; route it through the SELECT analyzer
			// so they land in AccessTables.
			w.visitSelect(q, false)
			return false
		case *ast.SetOpStmt:
			w.visitSetOp(q, false)
			return false
		case *ast.GroupedQuery:
			w.visitGroupedQuery(q, false)
			return false
		case *ast.TableRef:
			// A physical table referenced by DML — UPDATE ... FROM secret,
			// DELETE ... USING secret, MERGE ... USING secret — is a read the
			// span must report. This deliberately over-approximates by also
			// recording the write target: for access checks the safe error is
			// an extra entry, never a missing one.
			w.visitTableRef(q)
			return false
		}
		return true
	})
}

// visitSetOp walks a UNION/INTERSECT/EXCEPT tree. The left arm is the one
// whose Results we surface when this is the outermost statement — matching
// SQL's "take column names from the first SELECT" rule.
func (w *spanWalker) visitSetOp(n *ast.SetOpStmt, outermost bool) {
	if n == nil {
		return
	}

	// WITH on the leftmost SELECT scopes over the entire set operation —
	// WITH c AS (...) SELECT 1 UNION SELECT * FROM c resolves c in the right
	// arm too (engine-verified) — so install its names before walking the
	// arms, or the right arm's c is misreported as a physical table. The
	// names are installed here only; visitSelect records span.CTEs and walks
	// the bodies when it reaches the left arm.
	if with := leftmostWith(n); with != nil {
		scope := &cteScope{names: make(map[string]bool), parent: w.scope}
		for _, cte := range with.CTEs {
			if cte == nil || cte.Name == "" {
				continue
			}
			scope.names[strings.ToLower(cte.Name)] = true
		}
		saved := w.scope
		w.scope = scope
		defer func() { w.scope = saved }()
	}

	switch l := n.Left.(type) {
	case *ast.SelectStmt:
		w.visitSelect(l, outermost)
	case *ast.SetOpStmt:
		w.visitSetOp(l, outermost)
	case *ast.GroupedQuery:
		w.visitGroupedQuery(l, outermost)
	}
	switch r := n.Right.(type) {
	case *ast.SelectStmt:
		w.visitSelect(r, false)
	case *ast.SetOpStmt:
		w.visitSetOp(r, false)
	case *ast.GroupedQuery:
		w.visitGroupedQuery(r, false)
	}

	// Trailing clauses on the combined result — (SELECT 1) UNION (SELECT 2)
	// ORDER BY (SELECT x FROM secret) — carry expressions whose table reads
	// must land in AccessTables like any other clause.
	for _, o := range n.OrderBy {
		if o != nil && o.Expr != nil {
			w.walkExpr(o.Expr)
		}
	}
	if n.Limit != nil {
		w.walkExpr(n.Limit)
	}
	if n.Offset != nil {
		w.walkExpr(n.Offset)
	}
}

// leftmostWith walks the left spine of a set-operation tree to the leading
// SelectStmt and returns its WITH clause, whose names scope over the whole
// tree. A GroupedQuery stops the walk: it marks a parenthesized group (or a
// repeated clause group), and the engine scopes a WITH inside parens to that
// group only — (WITH c AS (...) SELECT 1) UNION SELECT * FROM c reads a
// physical table c (container-verified).
func leftmostWith(node ast.Node) *ast.WithClause {
	for {
		switch n := node.(type) {
		case *ast.SetOpStmt:
			node = n.Left
		case *ast.SelectStmt:
			return n.With
		default:
			return nil
		}
	}
}

// visitSelect processes one SelectStmt. It pushes a new CTE scope so the
// names defined by this WITH clause shadow any outer ones, then walks the
// FROM / WHERE / HAVING / QUALIFY / SELECT list, collecting tables and
// column references as it goes.
func (w *spanWalker) visitSelect(stmt *ast.SelectStmt, outermost bool) {
	if stmt == nil {
		return
	}

	// Push a scope. CTE bodies themselves run in the *parent* scope — a CTE
	// may not reference itself except via RECURSIVE, and sibling CTEs are
	// allowed to reference earlier siblings. To keep things simple we treat
	// all CTEs as visible to each other's bodies; that matches the legacy
	// extractor's behavior.
	scope := &cteScope{names: make(map[string]bool), parent: w.scope}
	if stmt.With != nil {
		for _, cte := range stmt.With.CTEs {
			if cte == nil || cte.Name == "" {
				continue
			}
			scope.names[strings.ToLower(cte.Name)] = true
			w.span.CTEs = append(w.span.CTEs, cte.Name)
		}
	}
	w.scope = scope

	// Walk CTE bodies first so their FROM tables land in AccessTables.
	if stmt.With != nil {
		for _, cte := range stmt.With.CTEs {
			if cte == nil || cte.Query == nil {
				continue
			}
			switch q := cte.Query.(type) {
			case *ast.SelectStmt:
				w.visitSelect(q, false)
			case *ast.SetOpStmt:
				w.visitSetOp(q, false)
			case *ast.GroupedQuery:
				w.visitGroupedQuery(q, false)
			}
		}
	}

	// FROM clause.
	for _, src := range stmt.From {
		w.visitFromItem(src)
	}

	// Predicate clauses — we walk them to pick up ColumnRef/subquery tables.
	if stmt.Where != nil {
		w.walkExpr(stmt.Where)
	}
	for _, g := range stmt.GroupBy {
		w.walkExpr(g)
	}
	if stmt.Having != nil {
		w.walkExpr(stmt.Having)
	}
	if stmt.Qualify != nil {
		w.walkExpr(stmt.Qualify)
	}
	for _, o := range stmt.OrderBy {
		if o != nil && o.Expr != nil {
			w.walkExpr(o.Expr)
		}
	}
	if stmt.Limit != nil {
		w.walkExpr(stmt.Limit)
	}
	if stmt.Offset != nil {
		w.walkExpr(stmt.Offset)
	}

	// SELECT list — walked last so Results reflect the final column order.
	for _, item := range stmt.Items {
		if item == nil {
			continue
		}
		if outermost && !w.resolved {
			w.span.Results = append(w.span.Results, w.makeColumnInfo(item))
		}
		if item.Expr != nil {
			w.walkExpr(item.Expr)
		}
	}
	if outermost {
		w.resolved = true
	}

	// Pop the scope.
	w.scope = scope.parent
}

// visitFromItem walks one entry in the FROM list. It may be a TableRef (bare
// table or subquery) or a JoinClause tree.
func (w *spanWalker) visitFromItem(node ast.Node) {
	switch n := node.(type) {
	case *ast.TableRef:
		w.visitTableRef(n)
	case *ast.JoinClause:
		w.visitFromItem(n.Left)
		w.visitFromItem(n.Right)
		if n.On != nil {
			w.walkExpr(n.On)
		}
	}
}

// visitTableRef handles a single TableRef from the FROM clause. The parser
// packs a FROM-subquery's raw text into Name.Parts[0]; we detect that case by
// trying to re-parse the text and, if it's a SELECT/set-op, recurse into it
// instead of treating it as a physical table.
func (w *spanWalker) visitTableRef(ref *ast.TableRef) {
	if ref == nil {
		return
	}
	// LATERAL VIEW generators carry real expressions (and possibly
	// subqueries); walk them regardless of the underlying source kind.
	for _, lv := range ref.LateralViews {
		if lv.Func != nil {
			w.walkExpr(lv.Func)
		}
	}
	// A table-valued function is not a physical table access. Mirror the
	// StarRocks TableFunctionRef handling: walk the call's arguments for
	// lineage and record no table — authorization would otherwise be asked
	// about a nonexistent table named BACKENDS or numbers.
	if ref.Func != nil {
		w.walkExpr(ref.Func)
		return
	}
	if ref.Name == nil || len(ref.Name.Parts) == 0 {
		return
	}

	// The parser marks FROM-subqueries explicitly; Name.Parts[0] mirrors the
	// raw text only for legacy consumers. Never classify by the text itself —
	// a table named `selected` or a quoted `SELECT` is not a subquery.
	if ref.Subquery != nil {
		w.analyzeSubqueryText(ref.Subquery.RawText, ref.Subquery.TextStart)
		return
	}

	var database, table string
	switch len(ref.Name.Parts) {
	case 1:
		table = ref.Name.Parts[0]
		// Bare name: skip if it matches an in-scope CTE.
		if w.scope.isCTE(table) {
			return
		}
	case 2:
		database = ref.Name.Parts[0]
		table = ref.Name.Parts[1]
	default:
		// catalog.db.table or longer — keep the last two parts.
		database = ref.Name.Parts[len(ref.Name.Parts)-2]
		table = ref.Name.Parts[len(ref.Name.Parts)-1]
	}

	key := tableKey{Database: database, Table: table, Alias: ref.Alias}
	if _, ok := w.accessed[key]; ok {
		return
	}
	w.accessed[key] = len(w.span.AccessTables)
	w.span.AccessTables = append(w.span.AccessTables, TableAccess{
		Database: database,
		Table:    table,
		Alias:    ref.Alias,
		Loc: ast.Loc{
			// Rebased into outer-statement coordinates: inside a reparsed
			// subquery ref.Loc is relative to the extracted text.
			Start: ref.Loc.Start + w.textBase,
			End:   ref.Loc.End + w.textBase,
		},
	})
}

// analyzeSubqueryText re-parses subquery text (SubqueryExpr.RawText) and
// recurses into any resulting SELECT. A subquery that does not fully parse
// records w.parseErr — its table reads would otherwise silently vanish from
// AccessTables, so the span fails closed instead.
//
// base is the byte offset of text within the outer statement
// (SubqueryExpr.TextStart): nested parse errors carry positions relative to
// the extracted text, and are shifted back into the outer statement's
// coordinates so editor diagnostics highlight the right spot.
func (w *spanWalker) analyzeSubqueryText(text string, base int) {
	// base is relative to the text this walker is currently inside; abs is
	// the position in the original statement, accumulated across levels.
	abs := w.textBase + base
	if strings.TrimSpace(text) == "" {
		// An empty placeholder body — EXISTS () — is not a query; the engine
		// rejects the form, and returning silently would skip the query-node
		// validation below entirely.
		w.noteNonQuerySubquery(ast.Loc{Start: abs, End: abs})
		return
	}
	file, errs := parser.Parse(text)
	if len(errs) > 0 {
		if w.parseErr == nil {
			e := errs[0]
			e.Loc.Start += abs
			e.Loc.End += abs
			w.parseErr = &e
		}
		return
	}
	if file == nil || len(file.Stmts) != 1 {
		// Comment-only bodies parse to zero statements; a body with embedded
		// delimiters parses to several. A subquery placeholder must hold
		// exactly one query — the engine rejects both shapes.
		w.noteNonQuerySubquery(ast.Loc{Start: abs, End: abs + len(text)})
		return
	}
	saved := w.textBase
	w.textBase = abs
	defer func() { w.textBase = saved }()
	for _, stmt := range file.Stmts {
		switch n := stmt.(type) {
		case *ast.SelectStmt:
			w.visitSelect(n, false)
		case *ast.SetOpStmt:
			w.visitSetOp(n, false)
		case *ast.GroupedQuery:
			// EXISTS (SELECT 1 FROM t ORDER BY 1 ORDER BY 2) is engine-valid;
			// the repeated clause group must analyze like any other query.
			w.visitGroupedQuery(n, false)
		default:
			// N2: a placeholder body that parses cleanly as something other
			// than a query (EXISTS (DELETE FROM secret) FROM public) is not a
			// valid subquery — the engine rejects it, and ignoring it here
			// would hide its reads from the span. Fail closed.
			loc := ast.NodeLoc(stmt)
			w.noteNonQuerySubquery(ast.Loc{Start: loc.Start + abs, End: loc.End + abs})
		}
	}
}

// walkExpr descends into an expression tree collecting ColumnRefs and
// recursing into any subquery nodes (SubqueryExpr, ExistsExpr, InExpr with
// a subquery value).
func (w *spanWalker) walkExpr(node ast.Node) {
	if node == nil {
		return
	}
	ast.Walk(&exprVisitor{w: w}, node)
}

type exprVisitor struct {
	w *spanWalker
}

func (v *exprVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	switch n := node.(type) {
	case *ast.SubqueryExpr:
		v.w.analyzeSubqueryText(n.RawText, n.TextStart)
		return nil // raw-text body, no parsed children
	case *ast.ExistsExpr:
		if n.Subquery != nil {
			v.w.analyzeSubqueryText(n.Subquery.RawText, n.Subquery.TextStart)
		}
		return nil
	case *ast.ColumnRef:
		// ColumnRef is recorded via the outer SELECT's SourceColumns; at this
		// point we don't know which output column it feeds, so we let
		// makeColumnInfo capture the direct references and stop here to avoid
		// double-visiting the ObjectName child.
		return nil
	}
	return v
}

// makeColumnInfo derives a ColumnInfo from a SELECT-list item. The Name is
// the alias if present; otherwise a simple rendering of the expression
// (column name for a bare ColumnRef, "*" for star). SourceColumns is the set
// of ColumnRefs the item directly references.
func (w *spanWalker) makeColumnInfo(item *ast.SelectItem) ColumnInfo {
	info := ColumnInfo{}
	if item.Star {
		if item.TableName != nil {
			info.Name = item.TableName.String() + ".*"
		} else {
			info.Name = "*"
		}
		return info
	}
	if item.Aliased {
		info.Name = item.Alias
	} else if item.Expr != nil {
		info.Name = renderColumnName(item.Expr)
	}
	if item.Expr != nil {
		info.SourceColumns = collectColumnRefs(item.Expr)
	}
	return info
}

// renderColumnName returns a short human-readable name for an expression.
// Only simple cases get a useful rendering; otherwise the empty string is
// returned and the caller can fall back to the raw source text if desired.
func renderColumnName(expr ast.Node) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if e.Name == nil || len(e.Name.Parts) == 0 {
			return ""
		}
		return e.Name.Parts[len(e.Name.Parts)-1]
	case *ast.ParenExpr:
		return renderColumnName(e.Expr)
	}
	return ""
}

// collectColumnRefs returns every ColumnRef directly mentioned in expr.
func collectColumnRefs(expr ast.Node) []ColumnRef {
	return collectColumnRefsScoped(expr, nil)
}

// collectColumnRefsScoped collects direct column references, skipping any name
// bound by an enclosing lambda. In array_map(x -> x + bonus, scores) the
// parameter x is local to the lambda and is not a column of any table, so
// reporting it would feed a nonexistent column into lineage and masking.
func collectColumnRefsScoped(expr ast.Node, bound map[string]bool) []ColumnRef {
	var refs []ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.LambdaExpr:
			// Inspect has no post-order hook to pop a scope with, so recurse
			// into the body with the widened binding set instead.
			inner := make(map[string]bool, len(bound)+len(node.Params))
			for name := range bound {
				inner[name] = true
			}
			for _, param := range node.Params {
				inner[strings.ToLower(param)] = true
			}
			refs = append(refs, collectColumnRefsScoped(node.Body, inner)...)
			return false
		case *ast.ColumnRef:
			if node.Name == nil || len(node.Name.Parts) == 0 {
				return false
			}
			// Only an unqualified name can be a lambda parameter; `t.x` is a
			// column even when a parameter happens to be called x.
			if len(node.Name.Parts) == 1 && bound[strings.ToLower(node.Name.Parts[0])] {
				return false
			}
			refs = append(refs, columnRefFromObjectName(node.Name))
			return false
		case *ast.SubqueryExpr, *ast.ExistsExpr:
			// Don't follow subqueries when collecting the *direct* column
			// references of the current SELECT item.
			return false
		}
		return true
	})
	return refs
}

func columnRefFromObjectName(name *ast.ObjectName) ColumnRef {
	parts := name.Parts
	switch len(parts) {
	case 1:
		return ColumnRef{Column: parts[0]}
	case 2:
		return ColumnRef{Table: parts[0], Column: parts[1]}
	case 3:
		return ColumnRef{Database: parts[0], Table: parts[1], Column: parts[2]}
	default:
		// Longer than 3 parts — keep the rightmost three.
		return ColumnRef{
			Database: parts[len(parts)-3],
			Table:    parts[len(parts)-2],
			Column:   parts[len(parts)-1],
		}
	}
}

// finalize performs any post-processing after the walk. AccessTables are
// already deduplicated during visitation; CTE names are left as-added so
// ordering matches declaration order.
func (w *spanWalker) finalize() {}
