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
// GetQuerySpan is tolerant of parse errors — if the parser produces a partial
// AST, whatever was parsed is still analyzed. On empty input it returns a
// zero-valued span with Type=QueryTypeUnknown.
func GetQuerySpan(statement string) (*QuerySpan, error) {
	file, _ := parser.Parse(statement)
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
	}
}

// visitSetOp walks a UNION/INTERSECT/EXCEPT tree. The left arm is the one
// whose Results we surface when this is the outermost statement — matching
// SQL's "take column names from the first SELECT" rule.
func (w *spanWalker) visitSetOp(n *ast.SetOpStmt, outermost bool) {
	if n == nil {
		return
	}
	switch l := n.Left.(type) {
	case *ast.SelectStmt:
		w.visitSelect(l, outermost)
	case *ast.SetOpStmt:
		w.visitSetOp(l, outermost)
	}
	switch r := n.Right.(type) {
	case *ast.SelectStmt:
		w.visitSelect(r, false)
	case *ast.SetOpStmt:
		w.visitSetOp(r, false)
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
	if ref == nil || ref.Name == nil || len(ref.Name.Parts) == 0 {
		return
	}

	// Detect FROM-subquery: parser stores the subquery's raw body in Parts[0].
	// A legitimate table name cannot contain whitespace, and every subquery
	// body starts with SELECT or WITH.
	if len(ref.Name.Parts) == 1 && looksLikeSubquery(ref.Name.Parts[0]) {
		w.analyzeSubqueryText(ref.Name.Parts[0])
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
		Loc:      ref.Loc,
	})
}

// looksLikeSubquery heuristically detects text the parser stuffed into
// ObjectName.Parts when it encountered a FROM-subquery. Bare identifiers
// (even quoted ones) don't contain these leading keywords.
func looksLikeSubquery(s string) bool {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false
	}
	upper := strings.ToUpper(trimmed)
	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "WITH") ||
		strings.HasPrefix(upper, "(")
}

// analyzeSubqueryText re-parses subquery text (from SubqueryExpr.RawText or
// a FROM-subquery's packed Parts[0]) and recurses into any resulting SELECT.
// Errors are swallowed — if the subquery is unparseable, the consumer of
// QuerySpan still gets whatever tables were already discovered.
func (w *spanWalker) analyzeSubqueryText(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	file, _ := parser.Parse(text)
	if file == nil {
		return
	}
	for _, stmt := range file.Stmts {
		switch n := stmt.(type) {
		case *ast.SelectStmt:
			w.visitSelect(n, false)
		case *ast.SetOpStmt:
			w.visitSetOp(n, false)
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
		v.w.analyzeSubqueryText(n.RawText)
		return nil // raw-text body, no parsed children
	case *ast.ExistsExpr:
		if n.Subquery != nil {
			v.w.analyzeSubqueryText(n.Subquery.RawText)
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
	if item.Alias != "" {
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
	var refs []ColumnRef
	ast.Inspect(expr, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ColumnRef:
			if node.Name == nil || len(node.Name.Parts) == 0 {
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
