package analysis

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// QuerySpan captures the tables and columns a GoogleSQL statement references. It
// is the data structure Bytebase consumes for field-level lineage (masking) and
// for SQL-editor access tracking (read/write permission checks); it is the omni
// counterpart of the legacy plugin/parser/{bigquery,spanner} GetQuerySpan result
// (which the bytebase-switch node maps onto base.QuerySpan).
//
// The extraction is best-effort and parser-driven: it walks the fully-connected
// AST the googlesql/parser produces (subqueries are already parsed into real
// QueryStmt nodes by the parser, so the walk reaches them directly — no
// re-parsing of raw text). Known best-effort boundaries (each a deliberate scope
// decision matching the merged trino/doris/snowflake analysis peers — full
// catalog-aware resolution belongs to the bytebase-switch layer):
//   - `SELECT *` is not expanded against catalog metadata; a star item is
//     recorded by name only ("*", or "t.*" for a qualified star).
//   - A column reference is recorded by its written qualifier parts; it is not
//     resolved to which FROM relation actually owns it.
//   - Table-valued-function and UNNEST sources contribute their argument column
//     references (as predicates) but no resolved output columns.
type QuerySpan struct {
	// Type is the statement classification (Select / SelectInfoSchema / DML / DDL
	// / Explain / Unknown).
	Type QueryType

	// AccessTables is the set of base tables the statement reads from or writes
	// to: every physical table path in FROM clauses, JOINs, subqueries, and DML/
	// DDL targets. Entries are deduplicated on (catalog, schema, table, alias);
	// CTE references and derived-table aliases are excluded.
	AccessTables []TableAccess

	// Results is the list of output columns produced by the outermost query (for
	// a query statement). For a set operation the left arm wins, matching SQL's
	// "column names come from the first select" rule. Empty for non-query
	// statements.
	Results []ColumnInfo

	// PredicateColumns is the set of columns referenced in filter/join positions
	// (WHERE, JOIN ON / USING, HAVING, QUALIFY, GROUP BY). Deduplicated.
	PredicateColumns []ColumnRef

	// CTEs lists the names defined by WITH clauses at any scope, in declaration
	// order.
	CTEs []string
}

// TableAccess represents one physical table referenced in a statement. The
// qualifier fields are dialect-bucketed (contract.md §4): for BigQuery a
// 2-part name is dataset.table (Database.Table, no Schema) unless the qualifier
// is INFORMATION_SCHEMA; for Spanner a 2-part name is schema.table. The
// rightmost component is always the table.
type TableAccess struct {
	Catalog  string  // optional catalog/project qualifier (empty when none)
	Database string  // optional database/dataset qualifier (empty when none)
	Schema   string  // optional schema qualifier (empty when none)
	Table    string  // table (or view) name
	Alias    string  // alias at the reference site (empty when none)
	IsSystem bool     // true when the reference is to a system/information schema
	Loc      ast.Loc // source location of the table reference
}

// ColumnInfo represents one output column of a query result.
type ColumnInfo struct {
	// Name is the output column name: the explicit alias if present, otherwise a
	// best-effort rendering (the column name for a bare column reference, "*"/
	// "t.*" for a star item, or "" when none can be derived).
	Name string

	// SourceColumns is the best-effort list of column references that directly
	// feed this output column (the column refs in the select-item expression,
	// excluding those inside nested subqueries).
	SourceColumns []ColumnRef
}

// ColumnRef identifies a column by its optional qualifier parts and its name.
// The rightmost component is always the column.
type ColumnRef struct {
	Catalog  string
	Database string
	Schema   string
	Table    string
	Column   string
}

// GetQuerySpan analyzes a single GoogleSQL statement in the given dialect and
// returns its query span: the statement classification, the base tables it
// touches, the CTE names it defines, its output columns with their source
// columns, and the columns used in predicate positions.
//
// It is tolerant of parse errors — if the parser produces a partial AST, whatever
// parsed is still analyzed. On empty input it returns a zero-valued span with
// Type=Unknown.
func GetQuerySpan(statement string, dialect Dialect) (*QuerySpan, error) {
	file := parseFile(statement)
	span := &QuerySpan{Type: ClassifyFromFile(file, dialect)}
	if file == nil || len(file.Stmts) == 0 {
		return span, nil
	}

	w := newSpanWalker(span, dialect)
	w.analyzeStmt(file.Stmts[0])
	return span, nil
}

// ---------------------------------------------------------------------------
// CTE scope
// ---------------------------------------------------------------------------

// cteScope is a linked stack of CTE name sets. Resolving a bare table reference
// walks outwards; an inner CTE shadows an outer one of the same name. GoogleSQL
// unquoted identifiers are case-insensitive for resolution, so CTE names are
// compared case-folded (matching the legacy findTableSchema's EqualFold CTE
// lookup).
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

// spanWalker walks the connected googlesql AST to populate a QuerySpan. It
// maintains a CTE scope stack (so bare references that name a CTE are filtered
// out of AccessTables), dedup maps for AccessTables and PredicateColumns, the
// current SELECT block's FROM aliases (so a correlated array/field path off an
// earlier FROM source is not mistaken for a base table), and a flag tracking
// whether the outermost query has populated Results.
type spanWalker struct {
	span        *QuerySpan
	dialect     Dialect
	scope       *cteScope
	fromAliases map[string]bool // alias + bare table name of the current SELECT's FROM sources
	accessed    map[tableKey]bool
	predSeen    map[ColumnRef]bool
	resolved    bool // Results captured for the outermost query
}

type tableKey struct {
	Catalog  string
	Database string
	Schema   string
	Table    string
	Alias    string
}

func newSpanWalker(span *QuerySpan, dialect Dialect) *spanWalker {
	return &spanWalker{
		span:     span,
		dialect:  dialect,
		accessed: make(map[tableKey]bool),
		predSeen: make(map[ColumnRef]bool),
	}
}

// analyzeStmt dispatches on the top-level statement node. A query statement
// yields full query-span data (results + predicates + tables); DML/DDL
// statements have their bodies walked for the tables they touch (masking and
// access tracking care about INSERT/UPDATE/DELETE/CTAS targets and sources).
func (w *spanWalker) analyzeStmt(node ast.Node) {
	switch n := node.(type) {
	case *ast.QueryStmt:
		w.visitQueryStmt(n, true /* outermost */)
	case *ast.SelectStmt:
		w.visitSelect(n, true)
	case *ast.SetOperation:
		w.visitSetOp(n, true)
	case *ast.InsertStmt:
		w.visitInsert(n)
	case *ast.UpdateStmt:
		w.visitUpdate(n)
	case *ast.DeleteStmt:
		w.visitDelete(n)
	case *ast.MergeStmt:
		w.visitMerge(n)
	case *ast.TruncateStmt:
		w.recordTable(n.Target, "")
		w.walkPredicate(n.Where)
	default:
		// DDL and everything else: walk for table references generically. The
		// generic table-discovery pass records every PathExpr that sits in a
		// table position (CTAS source, VIEW body, the target name).
		w.discoverTables(node)
	}
}

// ---------------------------------------------------------------------------
// Query / SELECT / set-op
// ---------------------------------------------------------------------------

// visitQueryStmt walks a QueryStmt (optional WITH + body + trailing ORDER BY/
// LIMIT/OFFSET). CTE names enter scope PROGRESSIVELY: a non-recursive CTE body
// sees only the siblings declared BEFORE it (GoogleSQL non-recursive WITH
// visibility, matching the legacy extractor's sequential recordCTE — a forward
// reference to a later sibling resolves to a base table, not the CTE). A
// RECURSIVE WITH additionally makes a CTE's OWN name visible inside its body (so
// the recursive self-reference is filtered out, not recorded as a base table).
func (w *spanWalker) visitQueryStmt(q *ast.QueryStmt, outermost bool) {
	if q == nil {
		return
	}

	scope := &cteScope{names: make(map[string]bool), parent: w.scope}
	w.scope = scope

	if q.With != nil {
		recursive := q.With.Recursive
		// Walk CTE bodies in declaration order, adding each name to scope at the
		// right moment so earlier-sibling (and, when RECURSIVE, self) references are
		// filtered while later-sibling names still resolve to base tables.
		for _, cte := range q.With.CTEs {
			if cte == nil || cte.Name == "" {
				continue
			}
			w.span.CTEs = append(w.span.CTEs, cte.Name)
			lower := strings.ToLower(cte.Name)
			// RECURSIVE: the name is visible inside its own body.
			if recursive {
				scope.names[lower] = true
			}
			if cte.Query != nil {
				w.visitBody(cte.Query, false)
			}
			// Non-recursive: the name becomes visible only AFTER its body (to the
			// following siblings and the main body).
			scope.names[lower] = true
		}
	}

	w.visitBody(q.Body, outermost)

	// Query-level ORDER BY / LIMIT / OFFSET keys may reference columns.
	for _, item := range q.OrderBy {
		if item != nil {
			w.walkPredicate(item.Expr)
		}
	}
	w.walkPredicate(q.Limit)
	w.walkPredicate(q.Offset)

	w.scope = scope.parent
}

// visitBody dispatches a query body / set-op operand / query primary.
func (w *spanWalker) visitBody(node ast.Node, outermost bool) {
	switch n := node.(type) {
	case *ast.QueryStmt:
		// A parenthesized ( query ) query_primary.
		w.visitQueryStmt(n, outermost)
	case *ast.SelectStmt:
		w.visitSelect(n, outermost)
	case *ast.SetOperation:
		w.visitSetOp(n, outermost)
	}
}

// visitSetOp walks a UNION/INTERSECT/EXCEPT tree. The left arm surfaces its
// Results when this is the outermost statement (SQL's first-select rule); the
// right arm never does.
func (w *spanWalker) visitSetOp(n *ast.SetOperation, outermost bool) {
	if n == nil {
		return
	}
	w.visitBody(n.Left, outermost)
	w.visitBody(n.Right, false)
}

// visitSelect processes one SelectStmt: its FROM relations, predicate clauses,
// and select list.
func (w *spanWalker) visitSelect(stmt *ast.SelectStmt, outermost bool) {
	if stmt == nil {
		return
	}

	// A new FROM-alias scope for this SELECT block. Populated as FROM sources are
	// visited so a later comma-source that is a correlated array/field path off an
	// earlier alias (e.g. `FROM T t, t.arr`) is recognized as a field path, not a
	// base table. Saved/restored so nested SELECTs do not leak aliases.
	prevFrom := w.fromAliases
	w.fromAliases = map[string]bool{}
	defer func() { w.fromAliases = prevFrom }()

	for _, src := range stmt.From {
		w.visitFromItem(src)
	}

	w.walkPredicate(stmt.Where)
	if stmt.GroupBy != nil {
		for _, item := range stmt.GroupBy.Items {
			if item == nil {
				continue
			}
			w.walkPredicate(item.Expr)
			for _, e := range item.Items {
				w.walkPredicate(e)
			}
		}
	}
	w.walkPredicate(stmt.Having)
	w.walkPredicate(stmt.Qualify)
	for _, win := range stmt.Window {
		if win != nil {
			w.walkWindowSpec(win.Spec)
		}
	}

	// SELECT list last so Results reflects final column order.
	for _, item := range stmt.Items {
		if item == nil {
			continue
		}
		if outermost && !w.resolved {
			w.span.Results = append(w.span.Results, w.makeColumnInfo(item))
		}
		// Walk the item expression for nested-subquery tables (the direct column
		// refs are captured in makeColumnInfo).
		w.walkSubqueriesOnly(item.Expr)
		if item.Modifiers != nil {
			for _, r := range item.Modifiers.Replace {
				if r != nil {
					w.walkSubqueriesOnly(r.Value)
				}
			}
		}
	}
	if outermost {
		w.resolved = true
	}
}

// ---------------------------------------------------------------------------
// FROM relations
// ---------------------------------------------------------------------------

// visitFromItem walks one FROM source (TableExpr / JoinExpr / UnnestExpr).
func (w *spanWalker) visitFromItem(node ast.Node) {
	switch n := node.(type) {
	case *ast.TableExpr:
		w.visitTableExpr(n)
	case *ast.JoinExpr:
		w.visitFromItem(n.Left)
		w.visitFromItem(n.Right)
		w.walkPredicate(n.On)
		for _, col := range n.Using {
			if col != "" {
				w.addPredicateColumn(ColumnRef{Column: col})
			}
		}
	case *ast.UnnestExpr:
		// UNNEST(array_expr) — the array argument's columns are predicate-ish
		// references; UNNEST produces no base-table access.
		w.walkPredicate(n.Array)
	}
}

// visitTableExpr handles a non-join FROM source: a base table path, a
// parenthesized subquery, or a table-valued function call. A subquery/TVF/array-
// path alias names a derived relation, not a base table, so it is not propagated
// into AccessTables — but the alias IS registered in the FROM scope so a later
// correlated path off it (`FROM (…) s, s.arr` or `FROM T t, t.arr a, a.child`)
// is also recognized as a correlated path rather than a base table.
func (w *spanWalker) visitTableExpr(n *ast.TableExpr) {
	if n == nil {
		return
	}
	switch {
	case n.Subquery != nil:
		w.visitBody(n.Subquery, false)
		w.registerFromAlias(n.Alias)
	case n.Func != nil:
		// Table-valued function: its arguments may reference columns and
		// subqueries. The argument columns are predicate-position references
		// (they parameterize the TVF), so collect them as predicates AND recurse
		// any subquery tables. The TVF's OUTPUT columns are unknown without a
		// function signature, so it contributes no resolved Results — matching the
		// legacy extractor, which rejected TVFs outright for lack of return-column
		// info.
		w.walkPredicate(n.Func)
		w.registerFromAlias(n.Alias)
	case n.Path != nil:
		// A multi-part path whose ROOT matches an earlier FROM alias/table in this
		// SELECT is a correlated array/field unnest (`FROM T t, t.arr`), NOT a base
		// table — the legacy extractor treats the array-path source as a non-table
		// (it rejected UNNEST/array sources). Record its root as a predicate column
		// instead so the lineage still reflects the array column it walks, and
		// register this source's own alias so a chained correlated path off it
		// (`t.arr a, a.child`) is likewise recognized.
		if len(n.Path.Parts) >= 2 && w.fromAliases[strings.ToLower(n.Path.Parts[0])] {
			w.addPredicateColumn(columnRefFromParts(n.Path.Parts, w.dialect))
			w.registerFromAlias(n.Alias)
		} else {
			w.recordPath(n.Path, n.Alias)
		}
	}
	w.walkPredicate(n.SystemTime)
}

// registerFromAlias records a derived-relation alias (subquery / TVF / correlated
// array-path source) in the current SELECT's FROM scope, so a later correlated
// path off it is recognized as a field path rather than a base table. A no-op
// when the alias is empty or there is no active FROM scope.
func (w *spanWalker) registerFromAlias(alias string) {
	if w.fromAliases == nil || alias == "" {
		return
	}
	w.fromAliases[strings.ToLower(alias)] = true
}

// ---------------------------------------------------------------------------
// DML
// ---------------------------------------------------------------------------

func (w *spanWalker) visitInsert(n *ast.InsertStmt) {
	if n == nil {
		return
	}
	w.recordTable(n.Target, "")
	for _, row := range n.Rows {
		if row == nil {
			continue
		}
		for _, v := range row.Values {
			w.walkPredicate(v)
		}
	}
	if n.Query != nil {
		w.visitBody(n.Query, false)
	}
	if n.TableClause != nil {
		if n.TableClause.Path != nil {
			w.recordPath(n.TableClause.Path, "")
		}
		if n.TableClause.Func != nil {
			w.walkPredicate(n.TableClause.Func)
		}
		w.walkPredicate(n.TableClause.Where)
	}
	if n.OnConflict != nil {
		for _, it := range n.OnConflict.SetItems {
			if it != nil {
				w.walkPredicate(it.Value)
			}
		}
		w.walkPredicate(n.OnConflict.Where)
	}
	if n.Returning != nil {
		for _, item := range n.Returning.Items {
			if item != nil {
				w.walkSubqueriesOnly(item.Expr)
			}
		}
	}
}

func (w *spanWalker) visitUpdate(n *ast.UpdateStmt) {
	if n == nil {
		return
	}
	w.recordTable(n.Target, n.Alias)
	for _, it := range n.Items {
		if it == nil {
			continue
		}
		w.walkPredicate(it.Value)
		if it.Nested != nil {
			w.analyzeStmt(it.Nested)
		}
	}
	for _, src := range n.From {
		w.visitFromItem(src)
	}
	w.walkPredicate(n.Where)
	if n.Returning != nil {
		for _, item := range n.Returning.Items {
			if item != nil {
				w.walkSubqueriesOnly(item.Expr)
			}
		}
	}
}

func (w *spanWalker) visitDelete(n *ast.DeleteStmt) {
	if n == nil {
		return
	}
	w.recordTable(n.Target, n.Alias)
	w.walkPredicate(n.Where)
	if n.Returning != nil {
		for _, item := range n.Returning.Items {
			if item != nil {
				w.walkSubqueriesOnly(item.Expr)
			}
		}
	}
}

func (w *spanWalker) visitMerge(n *ast.MergeStmt) {
	if n == nil {
		return
	}
	w.recordTable(n.Target, n.Alias)
	// The USING source is a TableExpr (a table path or a ( query ) subquery).
	if src, ok := n.Source.(*ast.TableExpr); ok {
		w.visitTableExpr(src)
	} else {
		w.discoverTables(n.Source)
	}
	w.walkPredicate(n.On)
	for _, when := range n.Whens {
		if when == nil {
			continue
		}
		w.walkPredicate(when.And)
		if when.Action != nil {
			if when.Action.InsertRow != nil {
				for _, v := range when.Action.InsertRow.Values {
					w.walkPredicate(v)
				}
			}
			for _, it := range when.Action.SetItems {
				if it != nil {
					w.walkPredicate(it.Value)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Generic table discovery (DDL bodies, CTAS, VIEW bodies, merge non-TableExpr)
// ---------------------------------------------------------------------------

// discoverTables records the tables a DDL statement touches: its (table-like)
// target name, any source tables it references (CTAS query, LIKE/CLONE/COPY
// source, INTERLEAVE parent, FOREIGN KEY references), and embedded query bodies.
// A full per-clause walk is unnecessary — the masking/access concern is purely
// which tables the statement reads or writes. Embedded queries are routed through
// the normal CTE-aware path (a raw ast.Inspect would double-count CTE names as
// base tables).
//
// Only TABLE/VIEW objects contribute a target table access; SCHEMA / DATABASE /
// INDEX names are NOT table references (recording them would pollute AccessTables
// with non-table objects — Codex finding #3). An index's ON <table> target IS
// recorded, and its INTERLEAVE parent too.
func (w *spanWalker) discoverTables(node ast.Node) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.CreateTableStmt:
		w.recordTable(n.Name, "")
		w.visitBody(n.AsQuery, false)
		// Source tables a CREATE TABLE references (Codex finding #4).
		w.recordTable(n.Like, "")
		w.recordTable(n.Clone, "")
		w.recordTable(n.Copy, "")
		if n.Interleave != nil {
			w.recordTable(n.Interleave.Parent, "")
		}
		w.recordConstraintRefs(n.Constraints)
		for _, col := range n.Columns {
			if col != nil && col.ForeignKey != nil {
				w.recordTable(col.ForeignKey.Table, "")
			}
		}
	case *ast.CreateViewStmt:
		w.recordTable(n.Name, "")
		w.visitBody(n.AsQuery, false)
	case *ast.CreateIndexStmt:
		// An index references its target table (ON <table>) and, on Spanner, its
		// INTERLEAVE parent. The index NAME itself is not a table.
		w.recordTable(n.Table, "")
		w.recordTable(n.Interleave, "")
	case *ast.AlterStmt:
		// Only ALTER TABLE / ALTER VIEW name a table-like object; ALTER INDEX /
		// SCHEMA / DATABASE do not (their name is not a table). An ADD CONSTRAINT
		// FK references another table, recorded regardless of object kind guard
		// below (the guard only gates the TARGET name).
		if n.Object == ast.AlterTable || n.Object == ast.AlterView {
			w.recordTable(n.Name, "")
		}
		for _, act := range n.Actions {
			if act == nil {
				continue
			}
			if act.Constraint != nil && act.Constraint.ForeignKey != nil {
				w.recordTable(act.Constraint.ForeignKey.Table, "")
			}
			if act.Column != nil && act.Column.ForeignKey != nil {
				w.recordTable(act.Column.ForeignKey.Table, "")
			}
		}
	case *ast.DropStmt:
		// Only DROP TABLE / DROP VIEW name a table-like object (Codex finding #3:
		// DROP SCHEMA/DATABASE/INDEX names are not tables).
		if n.Object == ast.DropTable || n.Object == ast.DropView {
			w.recordTable(n.Name, "")
		}
	case *ast.CreateSchemaStmt, *ast.CreateDatabaseStmt:
		// Schema/Database creates have no table reference of lineage interest.
	default:
		// Fallback: re-dispatch known query/DML nodes; otherwise nothing.
		w.redispatch(node)
	}
}

// recordConstraintRefs records the referenced table of every FOREIGN KEY
// table-constraint in the list.
func (w *spanWalker) recordConstraintRefs(constraints []*ast.TableConstraint) {
	for _, c := range constraints {
		if c != nil && c.ForeignKey != nil {
			w.recordTable(c.ForeignKey.Table, "")
		}
	}
}

// redispatch routes a node back through analyzeStmt for the node kinds that
// carry query-span data (used when discoverTables encounters an embedded
// statement it does not special-case).
func (w *spanWalker) redispatch(node ast.Node) {
	switch node.(type) {
	case *ast.QueryStmt, *ast.SelectStmt, *ast.SetOperation,
		*ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		w.analyzeStmt(node)
	}
}

// ---------------------------------------------------------------------------
// Table recording
// ---------------------------------------------------------------------------

// recordPath records a base-table reference from a *PathExpr FROM source with an
// optional alias, AND registers the source's alias + bare table name in the
// current SELECT's FROM scope so a later correlated array/field path off this
// source (e.g. `FROM T t, t.arr`) is recognized as a field path rather than a
// base table (see visitTableExpr).
func (w *spanWalker) recordPath(path *ast.PathExpr, alias string) {
	if path == nil {
		return
	}
	w.recordParts(path.Parts, alias, path.Loc)
	if w.fromAliases != nil {
		if alias != "" {
			w.fromAliases[strings.ToLower(alias)] = true
		}
		if n := len(path.Parts); n > 0 {
			// The bare (rightmost) table name is also a valid correlation root.
			w.fromAliases[strings.ToLower(path.Parts[n-1])] = true
		}
	}
}

// recordTable records a base-table reference from a *PathExpr target (DDL/DML
// targets are *PathExpr too).
func (w *spanWalker) recordTable(path *ast.PathExpr, alias string) {
	if path == nil {
		return
	}
	w.recordParts(path.Parts, alias, path.Loc)
}

// recordParts buckets the dotted name components per the dialect's metadata
// model, applies the CTE-shadowing rule for a bare name, dedups, and appends a
// TableAccess. The qualifier bucketing is shared with column references (see
// bucketNameParts) so a column's qualifier fields line up with the table's.
func (w *spanWalker) recordParts(parts []string, alias string, loc ast.Loc) {
	if len(parts) == 0 {
		return
	}

	cat, db, schema, table := bucketNameParts(parts, w.dialect)

	// Bare name: a match to an in-scope CTE is a CTE reference, not a base table.
	if cat == "" && db == "" && schema == "" && w.scope.isCTE(table) {
		return
	}

	ta := TableAccess{
		Catalog:  cat,
		Database: db,
		Schema:   schema,
		Table:    table,
		Alias:    alias,
		Loc:      loc,
	}
	ta.IsSystem = w.tableIsSystem(ta)

	key := tableKey{
		Catalog:  ta.Catalog,
		Database: ta.Database,
		Schema:   ta.Schema,
		Table:    ta.Table,
		Alias:    ta.Alias,
	}
	if w.accessed[key] {
		return
	}
	w.accessed[key] = true
	w.span.AccessTables = append(w.span.AccessTables, ta)
}

// bucketNameParts maps a dotted name's components to (catalog, database, schema,
// table) per the dialect's metadata model. It is the single bucketing rule used
// for BOTH table paths and column references, so a column's qualifier fields are
// directly comparable to the accessed table's (the consistency Codex finding #1
// flagged).
//
// Dialect bucketing (contract.md §4, mirroring the legacy accessTableListener):
//   - 1 part: bare table (no qualifier).
//   - 2 parts (X.table):
//     - BigQuery: X is Schema iff X==INFORMATION_SCHEMA, else X is Database
//       (dataset) — the legacy listener's INFORMATION_SCHEMA-only schema rule.
//     - Spanner: X is always Schema (named schema under one DB).
//   - 3 parts (X.Y.table):
//     - BigQuery: Y is Schema iff Y==INFORMATION_SCHEMA (then X is Database),
//       else Y is Database (dataset) and X is Catalog (project) — the
//       project.dataset.table shape.
//     - Spanner: Y is Schema, X is Database (db.schema.table — Spanner has no
//       catalog/project layer, so X is Database, not Catalog).
//   - >3 parts: the rightmost three follow the 3-part rule; the one extra
//     leading part is captured as Catalog for the BigQuery info-schema shape
//     (project.dataset.INFORMATION_SCHEMA.VIEW), otherwise dropped (no
//     meaningful deeper qualifier exists in either model).
func bucketNameParts(parts []string, dialect Dialect) (catalog, database, schema, table string) {
	n := len(parts)
	if n == 0 {
		return "", "", "", ""
	}
	table = parts[n-1]
	if n == 1 {
		return "", "", "", table
	}

	mid := parts[n-2] // 2nd-to-last
	if n == 2 {
		if dialect == DialectBigQuery {
			if isSystemSchemaName(mid, dialect) {
				schema = mid
			} else {
				database = mid
			}
		} else {
			schema = mid
		}
		return catalog, database, schema, table
	}

	// n >= 3.
	outer := parts[n-3] // 3rd-to-last
	if dialect == DialectBigQuery {
		if isSystemSchemaName(mid, dialect) {
			// [project.](dataset|region.)INFORMATION_SCHEMA.VIEW — INFORMATION_SCHEMA
			// is the Schema, the part before it is the dataset/region (Database), and
			// (when present, n>=4) the part before THAT is the project (Catalog). The
			// project must be captured so two projects' INFORMATION_SCHEMA do not
			// collapse to one access identity (re-review finding #2).
			schema = mid
			database = outer
			if n >= 4 {
				catalog = parts[n-4]
			}
		} else {
			// project.dataset.table: dataset is Database, project is Catalog.
			database = mid
			catalog = outer
		}
	} else {
		// Spanner db.schema.table: schema is Schema, db is Database (no
		// catalog/project layer).
		schema = mid
		database = outer
	}
	return catalog, database, schema, table
}

// tableIsSystem reports whether a recorded TableAccess is a system/information
// schema reference for the walker's dialect, matching the legacy isSystemResource
// (BigQuery: Schema EqualFold INFORMATION_SCHEMA; Spanner: INFORMATION_SCHEMA OR
// SPANNER_SYS). For Spanner, SPANNER_SYS appears as the table's Schema (e.g.
// SPANNER_SYS.X → Schema=SPANNER_SYS); for both, INFORMATION_SCHEMA does too.
func (w *spanWalker) tableIsSystem(ta TableAccess) bool {
	return isSystemSchemaName(ta.Schema, w.dialect)
}

// ---------------------------------------------------------------------------
// Expression walking
// ---------------------------------------------------------------------------

// addPredicateColumn appends a column reference to PredicateColumns,
// deduplicating against earlier predicate columns.
func (w *spanWalker) addPredicateColumn(ref ColumnRef) {
	if w.predSeen[ref] {
		return
	}
	w.predSeen[ref] = true
	w.span.PredicateColumns = append(w.span.PredicateColumns, ref)
}

// walkPredicate walks a predicate-position expression (WHERE / ON / HAVING /
// QUALIFY / GROUP BY / ORDER BY / VALUES), collecting its column references into
// PredicateColumns and recursing into nested subqueries (which may reference more
// base tables).
func (w *spanWalker) walkPredicate(expr ast.Node) {
	if expr == nil {
		return
	}
	ew := &exprWalk{w: w, followSub: true, onColumn: w.addPredicateColumn}
	ew.walk(expr)
}

// walkSubqueriesOnly walks an expression purely to discover the base tables
// inside nested subqueries; column refs are ignored (the direct refs of a select
// item are captured separately by makeColumnInfo).
func (w *spanWalker) walkSubqueriesOnly(expr ast.Node) {
	if expr == nil {
		return
	}
	ew := &exprWalk{w: w, followSub: true}
	ew.walk(expr)
}

// walkWindowSpec walks a window specification's PARTITION BY / ORDER BY / frame
// bounds as predicate-ish expressions.
func (w *spanWalker) walkWindowSpec(spec *ast.WindowSpec) {
	if spec == nil {
		return
	}
	for _, e := range spec.PartitionBy {
		w.walkPredicate(e)
	}
	for _, item := range spec.OrderBy {
		if item != nil {
			w.walkPredicate(item.Expr)
		}
	}
	if spec.Frame != nil {
		w.walkPredicate(spec.Frame.Start.Offset)
		if spec.Frame.Between {
			w.walkPredicate(spec.Frame.End.Offset)
		}
	}
}

// exprWalk parameterizes one expression traversal. onColumn is invoked for each
// column reference (nil to ignore columns). followSub controls whether subquery
// bodies are recursed into (true for predicate/table discovery; false when
// collecting a select item's DIRECT columns, which must not cross a subquery
// boundary).
type exprWalk struct {
	w         *spanWalker
	onColumn  func(ColumnRef)
	followSub bool
}

// walk is the shared expression traversal over the googlesql expression AST. It
// recognizes the two column-reference shapes — a leading Identifier and a
// FieldAccess/PathExpr dotted chain — invoking onColumn for each, and descends
// every expression child. Subquery bodies (already parsed into QueryStmt nodes)
// are recursed when followSub is set.
func (ew *exprWalk) walk(node ast.Node) {
	if node == nil {
		return
	}
	switch e := node.(type) {

	// --- Column references ---
	case *ast.Identifier:
		ew.emit(ColumnRef{Column: e.Name})
	case *ast.PathExpr:
		ew.emit(columnRefFromParts(e.Parts, ew.w.dialect))
	case *ast.FieldAccess:
		// A dotted name (t.a / a.b.c) is a FieldAccess chain over a leading
		// Identifier. Flatten it into a single qualified ColumnRef; if the base is
		// not a plain name chain (field access on a function result, etc.), fall
		// back to walking the base expression.
		if parts := flattenFieldAccess(e); parts != nil {
			ew.emit(columnRefFromParts(parts, ew.w.dialect))
		} else {
			ew.walk(e.Expr)
		}

	// --- Subquery bodies (parsed into real QueryStmt by the parser) ---
	case *ast.SubqueryExpr:
		ew.followSubquery(e.Query)
	case *ast.ExistsExpr:
		ew.followSubquery(e.Query)
	case *ast.ArraySubqueryExpr:
		ew.followSubquery(e.Query)

	// --- Composite expressions: recurse into children ---
	case *ast.ParenExpr:
		ew.walk(e.Expr)
	case *ast.UnaryExpr:
		ew.walk(e.Expr)
	case *ast.BinaryExpr:
		ew.walk(e.Left)
		ew.walk(e.Right)
	case *ast.CompareExpr:
		ew.walk(e.Left)
		ew.walk(e.Right)
	case *ast.IsExpr:
		ew.walk(e.Expr)
		ew.walk(e.DistinctFrom)
	case *ast.InExpr:
		ew.walk(e.Expr)
		ew.walkList(e.Values)
		ew.walk(e.Unnest)
		ew.walk(e.Subquery)
	case *ast.BetweenExpr:
		ew.walk(e.Expr)
		ew.walk(e.Low)
		ew.walk(e.High)
	case *ast.LikeExpr:
		ew.walk(e.Expr)
		ew.walk(e.Pattern)
		ew.walkList(e.QuantValues)
		ew.walk(e.QuantUnnest)
		ew.walk(e.QuantSubquery)
	case *ast.CaseExpr:
		ew.walk(e.Operand)
		for _, when := range e.Whens {
			if when != nil {
				ew.walk(when.Cond)
				ew.walk(when.Result)
			}
		}
		ew.walk(e.Else)
	case *ast.CastExpr:
		ew.walk(e.Expr)
		ew.walk(e.Format)
		ew.walk(e.TimeZone)
	case *ast.ExtractExpr:
		ew.walk(e.Part)
		ew.walk(e.From)
		ew.walk(e.TimeZone)
	case *ast.IntervalExpr:
		ew.walk(e.Value)
	case *ast.FuncCall:
		ew.walkFuncCall(e)
	case *ast.NamedArg:
		ew.walk(e.Value)
	case *ast.LambdaExpr:
		ew.walkLambda(e)
	case *ast.ArrayExpr:
		ew.walkList(e.Elements)
	case *ast.StructExpr:
		for _, f := range e.Fields {
			if f != nil {
				ew.walk(f.Value)
			}
		}
	case *ast.NewConstructor:
		for _, f := range e.Args {
			if f != nil {
				ew.walk(f.Value)
			}
		}
		if e.Braced != nil {
			ew.walkList(e.Braced.Fields)
		}
	case *ast.BracedConstructor:
		ew.walkList(e.Fields)
	case *ast.ReplaceFieldsExpr:
		ew.walk(e.Expr)
		for _, f := range e.Items {
			if f != nil {
				ew.walk(f.Value)
			}
		}
	case *ast.WithExpr:
		ew.walkWithExpr(e)
	case *ast.IndexAccess:
		ew.walk(e.Expr)
		ew.walk(e.Index)
	case *ast.ExtensionAccess:
		ew.walk(e.Expr)
	case *ast.StarExpr:
		// A bare '*' (COUNT(*)) carries no column reference.

	// --- Leaves: literals, parameters, system vars, typed literals ---
	default:
		// Literal / TypedLiteral / Parameter / SystemVariable / TypeRef / etc.
		// carry no column or subquery.
	}
}

// emit invokes onColumn for a resolved column reference, when collecting columns
// and the reference has a non-empty column name.
func (ew *exprWalk) emit(ref ColumnRef) {
	if ew.onColumn == nil || ref.Column == "" {
		return
	}
	ew.onColumn(ref)
}

// followSubquery recurses into a parsed subquery body when this pass follows
// subqueries. The body is a *QueryStmt (or set-op / select) the parser already
// filled, so the walk reaches it directly.
func (ew *exprWalk) followSubquery(query ast.Node) {
	if !ew.followSub || query == nil {
		return
	}
	ew.w.visitBody(query, false)
}

// walkFuncCall walks a function call's arguments, ORDER BY, FILTER, and OVER
// window for column/subquery references.
func (ew *exprWalk) walkFuncCall(e *ast.FuncCall) {
	if e == nil {
		return
	}
	ew.walkList(e.Args)
	for _, item := range e.OrderBy {
		if item != nil {
			ew.walk(item.Expr)
		}
	}
	if e.Having != nil {
		ew.walk(e.Having.Expr)
	}
	if e.Clamped != nil {
		ew.walk(e.Clamped.Low)
		ew.walk(e.Clamped.High)
	}
	ew.walk(e.Limit)
	ew.walk(e.LimitOffset)
	if e.Over != nil {
		// Reuse the walker's window-spec pass with this pass's settings.
		ew.walkWindowSpec(e.Over)
	}
}

// walkWindowSpec (on exprWalk) walks an OVER window spec with this pass's
// onColumn/followSub settings (so a window inside a select-item subquery pass
// only discovers subqueries, while a predicate pass also collects columns).
func (ew *exprWalk) walkWindowSpec(spec *ast.WindowSpec) {
	if spec == nil {
		return
	}
	ew.walkList(spec.PartitionBy)
	for _, item := range spec.OrderBy {
		if item != nil {
			ew.walk(item.Expr)
		}
	}
	if spec.Frame != nil {
		ew.walk(spec.Frame.Start.Offset)
		if spec.Frame.Between {
			ew.walk(spec.Frame.End.Offset)
		}
	}
}

// walkLambda walks a lambda body with the bound parameter names excluded from
// column collection: a body reference to a bare parameter name is the bound
// variable, not a table column, so it is dropped — while qualified references
// and free (outer-scope) names still pass through.
func (ew *exprWalk) walkLambda(lambda *ast.LambdaExpr) {
	if lambda == nil {
		return
	}
	ew.walkBoundBody(lambda.Params, lambda.Body)
}

// walkWithExpr walks an inline `WITH(name AS expr, …, body)` expression. A
// WITH-expr binding may reference EARLIER bindings (they are in scope left to
// right), and the body sees all of them — so each binding value is walked with
// the previously-bound names excluded, and the body with every bound name
// excluded. A reference to a bound variable is the local binding, not a table
// column (Codex findings #7 and re-review #1), exactly like a lambda parameter.
func (ew *exprWalk) walkWithExpr(e *ast.WithExpr) {
	if e == nil {
		return
	}
	var bound []string
	for _, v := range e.Vars {
		if v == nil {
			continue
		}
		// The binding value sees the names bound BEFORE it (not its own name).
		ew.walkBoundBody(bound, v.Value)
		if v.Alias != "" {
			bound = append(bound, v.Alias)
		}
	}
	ew.walkBoundBody(bound, e.Body)
}

// walkBoundBody walks body with the given bound names excluded from column
// collection: a column reference whose root name is one of the bound names is a
// local binding (lambda param / WITH-expr var), not a table column, so it is
// dropped; qualified and free names still pass through. When this pass does not
// collect columns (onColumn nil), it walks the body unchanged.
func (ew *exprWalk) walkBoundBody(bound []string, body ast.Node) {
	if ew.onColumn == nil || len(bound) == 0 {
		ew.walk(body)
		return
	}
	names := make(map[string]bool, len(bound))
	for _, b := range bound {
		if b != "" {
			names[strings.ToLower(b)] = true
		}
	}
	outer := ew.onColumn
	inner := &exprWalk{
		w:         ew.w,
		followSub: ew.followSub,
		onColumn: func(ref ColumnRef) {
			if names[strings.ToLower(rootName(ref))] {
				return
			}
			outer(ref)
		},
	}
	inner.walk(body)
}

// walkList walks each node in a slice.
func (ew *exprWalk) walkList(nodes []ast.Node) {
	for _, n := range nodes {
		ew.walk(n)
	}
}

// ---------------------------------------------------------------------------
// Select-item rendering
// ---------------------------------------------------------------------------

// makeColumnInfo derives a ColumnInfo from a select-item. Name is the alias if
// present, else a best-effort rendering ("*"/"t.*" for a star, the column name
// for a bare column reference, "" otherwise). SourceColumns is the set of column
// refs the item's expression directly references (excluding nested subqueries).
func (w *spanWalker) makeColumnInfo(item *ast.SelectItem) ColumnInfo {
	info := ColumnInfo{}
	if item.Star {
		if item.Expr != nil {
			// expr.* — name it "<expr>.*"; the underlying columns are the item's
			// direct column refs.
			if name := renderExprName(item.Expr); name != "" {
				info.Name = name + ".*"
			} else {
				info.Name = "*"
			}
			info.SourceColumns = w.collectDirectColumns(item.Expr)
		} else {
			info.Name = "*"
		}
		return info
	}

	if item.Alias != "" {
		info.Name = item.Alias
	} else {
		info.Name = renderExprName(item.Expr)
	}
	info.SourceColumns = w.collectDirectColumns(item.Expr)
	return info
}

// collectDirectColumns returns the column references directly mentioned in expr,
// in source order, excluding any inside nested subqueries (followSub disabled).
func (w *spanWalker) collectDirectColumns(expr ast.Node) []ColumnRef {
	if expr == nil {
		return nil
	}
	var refs []ColumnRef
	ew := &exprWalk{
		w:         w,
		followSub: false,
		onColumn:  func(ref ColumnRef) { refs = append(refs, ref) },
	}
	ew.walk(expr)
	return refs
}

// renderExprName returns a short human-readable name for a select-item
// expression: the column name for a bare Identifier, the trailing field of a
// FieldAccess, the last part of a PathExpr, unwrapping a single layer of
// parentheses. Other expressions have no derivable name ("").
func renderExprName(expr ast.Node) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Name
	case *ast.FieldAccess:
		return e.Field
	case *ast.PathExpr:
		if len(e.Parts) > 0 {
			return e.Parts[len(e.Parts)-1]
		}
	case *ast.ParenExpr:
		return renderExprName(e.Expr)
	}
	return ""
}

// ---------------------------------------------------------------------------
// Name helpers
// ---------------------------------------------------------------------------

// flattenFieldAccess returns the dotted-name components of a FieldAccess chain
// rooted at an Identifier (e.g. ["t", "a"] for t.a), or nil when the chain's
// base is not a plain name chain (so the caller can fall back to walking it).
func flattenFieldAccess(f *ast.FieldAccess) []string {
	if f == nil || f.Field == "" {
		return nil
	}
	switch base := f.Expr.(type) {
	case *ast.Identifier:
		if base.Name == "" {
			return nil
		}
		return []string{base.Name, f.Field}
	case *ast.PathExpr:
		if len(base.Parts) == 0 {
			return nil
		}
		return append(append([]string{}, base.Parts...), f.Field)
	case *ast.FieldAccess:
		prefix := flattenFieldAccess(base)
		if prefix == nil {
			return nil
		}
		return append(prefix, f.Field)
	default:
		return nil
	}
}

// columnRefFromParts maps the components of a dotted column reference to a
// ColumnRef, in the given dialect. The rightmost part is the column; the
// preceding parts name its qualifying table and are bucketed by the SAME
// dialect rule as a table path (bucketNameParts), so a ColumnRef's
// Catalog/Database/Schema/Table fields line up exactly with the corresponding
// TableAccess (the consistency Codex finding #1 flagged). E.g. BigQuery
// `proj.ds.t.c` → {Catalog:proj, Database:ds, Table:t, Column:c}, matching the
// TableAccess for `proj.ds.t`.
func columnRefFromParts(parts []string, dialect Dialect) ColumnRef {
	n := len(parts)
	if n == 0 {
		return ColumnRef{}
	}
	ref := ColumnRef{Column: parts[n-1]}
	if n >= 2 {
		ref.Catalog, ref.Database, ref.Schema, ref.Table = bucketNameParts(parts[:n-1], dialect)
	}
	return ref
}

// rootName returns the leftmost (most-qualifying) name of a column reference.
func rootName(ref ColumnRef) string {
	switch {
	case ref.Catalog != "":
		return ref.Catalog
	case ref.Database != "":
		return ref.Database
	case ref.Schema != "":
		return ref.Schema
	case ref.Table != "":
		return ref.Table
	default:
		return ref.Column
	}
}

// collectAccessTables runs a lightweight access-table-only pass over a statement
// node and returns its TableAccess set. Used by Classify's isAllSystemQuery to
// decide SelectInfoSchema without building a full span.
func collectAccessTables(node ast.Node, dialect Dialect) []TableAccess {
	span := &QuerySpan{}
	w := newSpanWalker(span, dialect)
	w.analyzeStmt(node)
	return span.AccessTables
}
