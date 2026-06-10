package analysis

import (
	"fmt"
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

	// CTEReferences lists the FROM-position table references that resolved to an
	// in-scope CTE (e.g. `FROM c`), deduplicated on (catalog, schema, table, alias)
	// like AccessTables. These are kept SEPARATE from AccessTables — a CTE is not a
	// physical table, so it must not be column-expanded or counted in the user/
	// system mix — but the legacy bytebase access-table listener recorded every
	// FROM table path INCLUDING CTE references in its flat source-column set, so the
	// catalog-aware consumer unions CTEReferences into the span's table-level
	// SourceColumns to reproduce that (a CTE reference surfaces as a
	// {default-dataset, cteName} resource).
	CTEReferences []TableAccess
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
	IsSystem bool    // true when the reference is to a system/information schema
	Loc      ast.Loc // source location of the table reference
}

// ColumnInfo represents one output column of a query result — or, when
// StarSegments is set, one `*`/`rel.*` star item that the catalog-aware consumer
// expands.
type ColumnInfo struct {
	// Name is the output column name: the explicit alias if present, otherwise a
	// best-effort rendering (the column name for a bare column reference, "*"/
	// "t.*" for a star item, or "" when none can be derived).
	Name string

	// SourceColumns is the best-effort list of column references that directly
	// feed this output column (the column refs in the select-item expression,
	// excluding those inside nested subqueries). For a column omni RESOLVED through
	// a CTE / derived relation, these carry the underlying BASE-table lineage (the
	// relation alias is resolved away); for a reference omni could not resolve to a
	// relation (a bare column over a base table whose columns omni cannot
	// enumerate), they carry the written qualifier parts for the consumer to match
	// against catalog metadata.
	SourceColumns []ColumnRef

	// IsPlain mirrors the legacy QuerySpanResult.IsPlainField: true when this output
	// column is a direct base-table column passthrough (a `SELECT *` / `rel.*`
	// expansion column that was not rewrapped by a join-left side, a set-operation
	// merge, or an explicit select-list item). An explicit select item — even a
	// bare `SELECT id` or a column resolved through a CTE/derived relation — is NOT
	// plain. The consumer copies this onto base.QuerySpanResult.IsPlainField.
	IsPlain bool

	// StarSegments, when non-nil, marks this ColumnInfo as a `*` / `rel.*` star item
	// (a bare star, a qualified star, or a `SELECT *` over a CTE / derived relation)
	// that the catalog-aware consumer expands IN ORDER into one output column per
	// segment-column. A segment is either a base-table star (BaseTable set — expand
	// every column of that physical table via metadata) or an already-resolved
	// concrete projection column (a CTE/derived column with its base lineage). The
	// star's EXCEPT/REPLACE modifiers (StarExcept/StarReplace) apply across the
	// fully expanded column list. When StarSegments is set, Name/SourceColumns are
	// not consumed for output (they remain only for star-shape detection back-compat
	// — Name is "*" or "rel.*").
	StarSegments []StarSegment

	// StarExcept holds the column names of a `SELECT * EXCEPT (a, b)` modifier on a
	// star item; nil for a non-star item or a star with no EXCEPT. The catalog-aware
	// consumer (the bytebase extractor) drops these names from the star's expanded
	// column set so an EXCEPT-ed column is not surfaced (and so over-masking does not
	// occur by leaving it in).
	StarExcept []string

	// StarReplace holds the `SELECT * REPLACE (expr AS name)` substitutions on a
	// star item; nil for a non-star item or a star with no REPLACE. The consumer
	// replaces the star-expanded output column named Name with one whose lineage is
	// the replacement expression's Sources.
	StarReplace []StarReplaceItem

	// StarMerge, when non-nil, marks this output position as a set-operation merge
	// against a base-table-star arm whose arity only the metadata-aware consumer
	// knows: the consumer expands StarMerge.Table, takes its StarMerge.Index-th
	// column, unions that column into this position's SourceColumns, and (when
	// StarMerge.LeftStar) takes this position's output Name from that column. This
	// reproduces the legacy "expand the star arm against metadata, then position-
	// merge" behaviour for `concrete UNION ALL (SELECT * FROM base)`.
	StarMerge *StarMergeInfo

	// SetOpMerge, when non-nil, marks this ColumnInfo as a WHOLE deferred set-
	// operation merge whose arms could not be position-merged inline because at
	// least one arm carries an un-enumerable star (a base-table star, an
	// EXCEPT/REPLACE star, or a nested merge). The metadata-aware consumer expands
	// each arm's projection fully (recursively, reusing this same per-column
	// expansion) and then position-merges the two expanded lists: output names from
	// the LEFT arm, SourceColumns unioned per position, IsPlainField=false. This
	// reproduces the legacy "fully resolve each arm, then zip"
	// (extractTableSourceFromQuerySetOperation) for the star-involving arm
	// combinations a per-position StarMerge cannot express — a base-star UNION
	// base-star (one merged output column per expanded position, not two
	// concatenated stars) and an EXCEPT/REPLACE star arm in a set operation. When
	// set, Name is "*" and Star*/SourceColumns are not consumed for output.
	SetOpMerge *SetOpMergeInfo
}

// SetOpMergeInfo carries a deferred set-operation merge: the Left and Right arms'
// resolved projections, each a list of ColumnInfo the consumer expands fully (a
// base-table star against metadata, an EXCEPT/REPLACE star with its modifiers, a
// nested merge recursively) before merging them. See ColumnInfo.SetOpMerge.
//
// ByName marks a BY NAME / CORRESPONDING set operation (fix F3): after expanding
// both arms the consumer merges them by case-insensitive column NAME — output
// order is the left arm's columns (then any right-only names appended), each
// unioning both arms' same-named lineage — instead of by ordinal. MatchColumns,
// when non-empty, is the `ON (cols)` / `BY (cols)` restriction list: the output
// is exactly those columns, in list order. Both are zero for an ordinal merge.
type SetOpMergeInfo struct {
	Left         []ColumnInfo
	Right        []ColumnInfo
	ByName       bool
	MatchColumns []string
}

// StarMergeInfo carries a set-operation merge against a base-table-star arm: the
// base Table to expand, the output ordinal (Index) whose lineage gains that
// table's Index-th column, and LeftStar (the star arm was the LEFT one, so the
// output name is taken from the expanded base column rather than from a concrete
// arm). See ColumnInfo.StarMerge.
type StarMergeInfo struct {
	Table    ColumnRef
	Index    int
	LeftStar bool
}

// StarSegment is one element of a star item's resolved expansion (see
// ColumnInfo.StarSegments). Exactly one shape applies:
//   - BaseTable set: a base-table star — the consumer expands EVERY column of the
//     named physical table (via catalog metadata), emitting each with IsPlain
//     (Plain). This is the only piece omni leaves to the metadata-aware consumer,
//     because omni has no catalog to enumerate a physical table's columns.
//     ExceptColumns, when non-empty, lists column names the expansion must SKIP
//     (case-insensitively): a JOIN ... USING key is projected once as a coalesced
//     column, so each side's star excludes it (fix F1) — expanding it again would
//     shift every later position and misalign the positional masker.
//   - BaseTable nil: an already-resolved concrete projection column (a CTE /
//     derived / explicit column reached through a relation), with its base lineage
//     in Sources, output Name, and plainness Plain. The consumer emits it directly.
type StarSegment struct {
	BaseTable     *ColumnRef  // base table to expand; nil for a concrete segment
	ExceptColumns []string    // base-table column names to skip when expanding
	Name          string      // concrete segment output column name
	Sources       []ColumnRef // concrete segment base lineage
	Plain         bool        // IsPlainField for this segment's column(s)
}

// StarReplaceItem is one `expr AS name` entry of a star REPLACE modifier: the
// output column Name whose value is overridden, and the source columns the
// replacement expression directly references (Sources).
type StarReplaceItem struct {
	Name    string
	Sources []ColumnRef
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
	// Fail closed (structural rule): a shape the walker accepted but could not
	// resolve to CORRECT lineage (NATURAL JOIN, a USING join over a deferred
	// set-op projection) surfaces as an error rather than silently misaligned /
	// empty lineage — the masking consumer must reject the statement, not return
	// a sensitive column unmasked.
	if w.failure != nil {
		return nil, w.failure
	}
	return span, nil
}

// ---------------------------------------------------------------------------
// CTE scope
// ---------------------------------------------------------------------------

// cteScope is a linked stack of CTE definitions. Resolving a bare table
// reference walks outwards; an inner CTE shadows an outer one of the same name.
// GoogleSQL unquoted identifiers are case-insensitive for resolution, so CTE
// names are compared case-folded (matching the legacy findTableSchema's EqualFold
// CTE lookup).
//
// Each in-scope CTE retains its RESOLVED output projection (the masking-grade
// upgrade): a `SELECT *` / `rel.*` / column reference over the CTE reproduces that
// projection with base-column lineage, rather than collapsing to a catalog-blind
// star (the legacy resolver stored each CTE as a base.PseudoTable with its
// resolved columns — recordCTE).
type cteScope struct {
	names   map[string]bool
	columns map[string][]projColumn // CTE name (lower-cased) → resolved projection
	parent  *cteScope
}

func (s *cteScope) isCTE(name string) bool {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.names[strings.ToLower(name)] {
			return true
		}
	}
	return false
}

// cteColumns returns the resolved projection of the nearest in-scope CTE named
// `name`, or nil if no such CTE is in scope (or its projection was not retained,
// e.g. a recursive CTE still being defined).
func (s *cteScope) cteColumns(name string) ([]projColumn, bool) {
	lower := strings.ToLower(name)
	for cur := s; cur != nil; cur = cur.parent {
		if cur.names[lower] {
			cols, ok := cur.columns[lower]
			return cols, ok
		}
	}
	return nil, false
}

// ---------------------------------------------------------------------------
// Walker
// ---------------------------------------------------------------------------

// spanWalker walks the connected googlesql AST to populate a QuerySpan. It
// maintains a CTE scope stack (so bare references that name a CTE are filtered
// out of AccessTables), dedup maps for AccessTables and PredicateColumns, and
// the current SELECT block's FROM aliases (so a correlated array/field path off
// an earlier FROM source is not mistaken for a base table).
//
// Output columns (Results) are computed bottom-up: every query body — a SELECT,
// a set operation, or a parenthesized QueryStmt — RETURNS its output ColumnInfo
// list, and the outermost caller assigns it to span.Results. Returning (rather
// than appending the first SELECT's items as a side effect) is what lets a set
// operation MERGE both arms' per-position lineage (the masking-grade union
// semantics): an out-col[i]'s sources are left.results[i].sources ∪
// right.results[i].sources.
type spanWalker struct {
	span        *QuerySpan
	dialect     Dialect
	scope       *cteScope
	fromAliases map[string]bool // alias + bare table name of the current SELECT's FROM sources
	// leafRels is the current SELECT's leaf FROM relations (base tables, CTE
	// references, derived subqueries) by reference name, used to resolve a column
	// reference (`rel.col`, a bare `col` over a CTE/derived, `rel.*`) to its
	// relation's projection. A join contributes its two sides as separate leaves
	// here (mirroring the legacy tableSourceFrom, which held each table AND the join
	// anchor), so `x.id` over `t x JOIN t y` resolves against leaf `x`. Saved and
	// restored around each SELECT so nested blocks do not leak relations.
	leafRels []*relation
	accessed map[tableKey]bool
	predSeen map[ColumnRef]bool
	// failure, when set, marks the span as fail-closed: the walker met a shape it
	// parses but cannot resolve to correct lineage without catalog metadata
	// (NATURAL JOIN; JOIN USING over a deferred set-op projection). GetQuerySpan
	// returns it as an error so the masking consumer rejects the statement instead
	// of consuming silently wrong lineage. The first failure wins. Classification
	// paths (collectAccessTables) deliberately ignore it — classification is
	// best-effort and not masking-grade.
	failure error
}

// failClosed records a fail-closed condition (first one wins).
func (w *spanWalker) failClosed(err error) {
	if w.failure == nil {
		w.failure = err
	}
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
		w.span.Results = projColumnsToColumnInfos(w.visitQueryStmt(n))
	case *ast.SelectStmt:
		w.span.Results = projColumnsToColumnInfos(w.visitSelect(n))
	case *ast.SetOperation:
		w.span.Results = projColumnsToColumnInfos(w.visitSetOp(n))
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
func (w *spanWalker) visitQueryStmt(q *ast.QueryStmt) []projColumn {
	if q == nil {
		return nil
	}

	scope := &cteScope{names: make(map[string]bool), columns: make(map[string][]projColumn), parent: w.scope}
	w.scope = scope

	if q.With != nil {
		recursive := q.With.Recursive
		// Walk CTE bodies in declaration order, adding each name to scope at the
		// right moment so earlier-sibling (and, when RECURSIVE, self) references are
		// filtered while later-sibling names still resolve to base tables. Each
		// CTE's RESOLVED projection is retained in the scope so a `SELECT *` / column
		// reference over the CTE reproduces it (legacy recordCTE → PseudoTable).
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
			var cols []projColumn
			if cte.Query != nil {
				cols = w.resolveCTEColumns(cte, recursive, scope, lower)
			}
			cols = applyCTEColumnAliases(cols, cte.Columns)
			// Non-recursive: the name becomes visible only AFTER its body (to the
			// following siblings and the main body).
			scope.names[lower] = true
			scope.columns[lower] = cols
		}
	}

	results := w.visitBody(q.Body)

	// Query-level ORDER BY / LIMIT / OFFSET keys may reference columns.
	for _, item := range q.OrderBy {
		if item != nil {
			w.walkPredicate(item.Expr)
		}
	}
	w.walkPredicate(q.Limit)
	w.walkPredicate(q.Offset)

	w.scope = scope.parent
	return results
}

// resolveCTEColumns resolves a CTE body to its projection. For a non-recursive
// CTE (or a recursive one whose body is not an anchor/recursive set operation) it
// is a plain body resolution. For a recursive CTE shaped `(anchor) UNION [ALL]
// recursive`, it mirrors the legacy recordRecursiveCTE: resolve the anchor arm,
// publish it as the CTE's projection so the recursive arm's self-references
// resolve to the anchor's columns (and thus to the base lineage), then resolve
// the recursive arm and position-merge it into the anchor (the union of both
// arms' per-position lineage). One pass suffices — lineage is monotone, so the
// fixpoint the legacy loop computes is reached after a single recursive-arm merge.
func (w *spanWalker) resolveCTEColumns(cte *ast.CTE, recursive bool, scope *cteScope, lower string) []projColumn {
	setOp := recursiveSetOp(recursive, cte.Query)
	if setOp == nil {
		return w.visitBody(cte.Query)
	}
	// Resolve the anchor arm and publish it so the recursive arm's self-references
	// resolve against it. Then iterate the recursive arm to a fixpoint (lineage is
	// monotone, so it converges): each pass re-resolves the recursive arm against
	// the current projection and position-merges it, until no source set grows. A
	// fixpoint is needed because a self-reference to column j may pull in column k's
	// lineage, which a later column's expression then reads — the legacy
	// recordRecursiveCTE loops for the same reason. The pass count is bounded by the
	// column count (each pass can only add sources; bounded extra guard prevents a
	// pathological loop).
	anchor := w.visitBody(setOp.Left)
	scope.columns[lower] = anchor
	current := anchor
	for iter := 0; iter <= len(anchor)+1; iter++ {
		recursivePart := w.visitBody(setOp.Right)
		merged := mergeProjections(current, recursivePart)
		if projectionsEqual(current, merged) {
			break
		}
		current = merged
		scope.columns[lower] = current
	}
	return current
}

// projectionsEqual reports whether two projections have identical per-position
// output names, source-column sets (order-insensitive on the sources), AND star
// markers. Used to detect the recursive-CTE lineage fixpoint. The star-marker
// comparison matters: a base-star / starMerge / setOpMerge arm carries no concrete
// sources, so comparing only name+sources would declare a still-changing
// star-shaped projection "equal" and stop the fixpoint before the recursive arm's
// star lineage was published (the masking under-attribution this guards). Equality
// therefore also requires the per-position star shape to match.
func projectionsEqual(a, b []projColumn) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].name != b[i].name || len(a[i].sources) != len(b[i].sources) {
			return false
		}
		if !starShapeEqual(a[i], b[i]) {
			return false
		}
		seen := make(map[ColumnRef]bool, len(a[i].sources))
		for _, r := range a[i].sources {
			seen[r] = true
		}
		for _, r := range b[i].sources {
			if !seen[r] {
				return false
			}
		}
	}
	return true
}

// starShapeEqual reports whether two projection columns carry the same star
// marker shape: both concrete, both a base-table star over the same table, or
// both a starMerge over the same table/index. A starGroup or setOpMerge marker is
// treated as never-equal (its expansion is metadata-dependent and not compared
// here) so the recursive fixpoint keeps iterating until the projection stops
// changing shape. This is only used by the recursive-CTE fixpoint, whose
// iteration count is bounded, so a conservative "not equal" merely costs an extra
// bounded pass rather than risking early termination with unpublished lineage.
func starShapeEqual(a, b projColumn) bool {
	if (a.baseStar == nil) != (b.baseStar == nil) {
		return false
	}
	if a.baseStar != nil && !sameBaseTable(*a.baseStar, *b.baseStar) {
		return false
	}
	if (a.starMerge == nil) != (b.starMerge == nil) {
		return false
	}
	if a.starMerge != nil &&
		(!sameBaseTable(a.starMerge.table, b.starMerge.table) ||
			a.starMerge.index != b.starMerge.index ||
			a.starMerge.leftStar != b.starMerge.leftStar) {
		return false
	}
	if a.starGroup != nil || b.starGroup != nil {
		return false
	}
	if a.setOpMerge != nil || b.setOpMerge != nil {
		return false
	}
	return true
}

// recursiveSetOp returns the CTE body's top-level set operation when the CTE is
// recursive and its body is `(anchor) UNION/INTERSECT/EXCEPT recursive` (directly,
// or wrapped in a QueryStmt with no further WITH), else nil. A recursive CTE
// without a set-operation body has no separate anchor/recursive arms (it then
// resolves like a normal body).
func recursiveSetOp(recursive bool, body ast.Node) *ast.SetOperation {
	if !recursive || body == nil {
		return nil
	}
	switch n := body.(type) {
	case *ast.SetOperation:
		return n
	case *ast.QueryStmt:
		if n.With == nil {
			if so, ok := n.Body.(*ast.SetOperation); ok {
				return so
			}
		}
	}
	return nil
}

// applyCTEColumnAliases renames a CTE body's resolved projection columns to the
// CTE's explicit column-name list (`cte_name (a, b, …) AS (…)`, the Spanner form)
// when one is present, mirroring the legacy behaviour of projecting under the
// declared names. Lineage and plainness are preserved; only the output names
// change. A mismatch in arity (or an un-enumerable base-star body) leaves the
// projection unchanged (best-effort).
func applyCTEColumnAliases(cols []projColumn, names []string) []projColumn {
	if len(names) == 0 || len(names) != len(cols) {
		return cols
	}
	out := make([]projColumn, len(cols))
	for i, c := range cols {
		c.name = names[i]
		out[i] = c
	}
	return out
}

// visitBody dispatches a query body / set-op operand / query primary, returning
// the resolved output projection it produces.
func (w *spanWalker) visitBody(node ast.Node) []projColumn {
	switch n := node.(type) {
	case *ast.QueryStmt:
		// A parenthesized ( query ) query_primary.
		return w.visitQueryStmt(n)
	case *ast.SelectStmt:
		return w.visitSelect(n)
	case *ast.SetOperation:
		return w.visitSetOp(n)
	}
	return nil
}

// visitSetOp walks a UNION/INTERSECT/EXCEPT tree and merges the two arms'
// per-position output projection. For masking the union must be conservative:
// out-col[i]'s sources are the UNION of left[i] AND right[i] (an output column of
// a set operation reads from BOTH arms at that position). The output column NAME
// comes from the LEFT arm (SQL's "column names come from the first SELECT"), so a
// three-arm left-associative tree keeps the leftmost names while accumulating
// every arm's sources. See mergeProjections for the base-table-star arm handling.
//
// A BY NAME / CORRESPONDING operation (fix F3) merges by column NAME instead:
// merging those by ordinal attributed each output to the OTHER arm's
// same-position column (verified mis-attribution — `id` carried a secret's
// policy). See mergeProjectionsByName.
func (w *spanWalker) visitSetOp(n *ast.SetOperation) []projColumn {
	if n == nil {
		return nil
	}
	left := w.visitBody(n.Left)
	right := w.visitBody(n.Right)
	if n.ByName || n.Corresponding {
		return mergeProjectionsByName(left, right, n.MatchColumns)
	}
	return mergeProjections(left, right)
}

// unionColumnRefs returns the order-preserving, deduplicated union of two
// ColumnRef slices (left first, then any right ref not already present).
func unionColumnRefs(left, right []ColumnRef) []ColumnRef {
	if len(right) == 0 {
		return left
	}
	seen := make(map[ColumnRef]bool, len(left)+len(right))
	out := make([]ColumnRef, 0, len(left)+len(right))
	for _, r := range left {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	for _, r := range right {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

// visitSelect processes one SelectStmt: its FROM relations, predicate clauses,
// and select list. It returns the resolved output projection (one position per
// select-list output column, with base-column lineage) so a parent set operation
// can merge per-position lineage; the outermost dispatch surfaces it as
// span.Results.
func (w *spanWalker) visitSelect(stmt *ast.SelectStmt) []projColumn {
	if stmt == nil {
		return nil
	}

	// A new FROM-alias scope for this SELECT block. Populated as FROM sources are
	// visited so a later comma-source that is a correlated array/field path off an
	// earlier alias (e.g. `FROM T t, t.arr`) is recognized as a field path, not a
	// base table. Saved/restored so nested SELECTs do not leak aliases.
	prevFrom := w.fromAliases
	w.fromAliases = map[string]bool{}
	prevLeaves := w.leafRels
	w.leafRels = nil
	defer func() { w.fromAliases = prevFrom; w.leafRels = prevLeaves }()

	// Resolve the FROM clause into the comma-item relations (a join collapses into
	// one combined relation, used for a bare `*`); building them also populates
	// w.leafRels (each base table / CTE / derived leaf, used for `rel.*` and column
	// resolution) and performs the AccessTable / predicate side-effects (a base path
	// is recorded, a subquery body walked, ON/USING predicates collected).
	fromRels := w.buildFromRelations(stmt.From)

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

	// Resolve the select list against the FROM relations (relation-aware lineage),
	// then walk each item for nested-subquery tables (the direct refs are resolved
	// by the projection resolver; the subquery walk only discovers tables).
	results := w.resolveSelectProjection(stmt, fromRels)
	for _, item := range stmt.Items {
		if item == nil {
			continue
		}
		w.walkSubqueriesOnly(item.Expr)
		if item.Modifiers != nil {
			for _, r := range item.Modifiers.Replace {
				if r != nil {
					w.walkSubqueriesOnly(r.Value)
				}
			}
		}
	}
	return results
}

// ---------------------------------------------------------------------------
// FROM relations
// ---------------------------------------------------------------------------

// buildFromRelations resolves a SELECT's FROM list into the ordered leaf
// relations used for `rel.*` and column resolution. Each comma-separated FROM
// item contributes one relation (a join collapses into one combined relation —
// see buildFromItem), and a bare `*` over the SELECT reproduces the in-order
// concatenation of these relations' projections (allColumns). Building a relation
// performs the same AccessTable / predicate side-effects the prior walk did.
func (w *spanWalker) buildFromRelations(from []ast.Node) []*relation {
	var rels []*relation
	for _, src := range from {
		if rel := w.buildFromItem(src); rel != nil {
			rels = append(rels, rel)
		}
	}
	return rels
}

// buildFromItem resolves one FROM source (TableExpr / JoinExpr / UnnestExpr) into
// a relation with its resolved projection, while performing the AccessTable /
// predicate side-effects. A JoinExpr collapses its two sides into one combined
// relation (the legacy joinTable): the left/anchor side's columns are rewrapped
// as non-plain, the right side's columns keep their plainness, and a USING join
// coalesces the key columns (fix F1 — see joinRelations). An UNNEST source
// builds a one-column value relation whose element carries the array argument's
// resolved lineage (fix F2 — the legacy extractor errored on array sources;
// dropping the relation here would shift `SELECT *` positions and resolve the
// element alias to nothing, both fail-open); its argument columns remain
// predicate references too.
func (w *spanWalker) buildFromItem(node ast.Node) *relation {
	switch n := node.(type) {
	case *ast.TableExpr:
		return w.buildTableExpr(n)
	case *ast.JoinExpr:
		left := w.buildFromItem(n.Left)
		right := w.buildFromItem(n.Right)
		w.walkPredicate(n.On)
		for _, col := range n.Using {
			if col != "" {
				w.addPredicateColumn(ColumnRef{Column: col})
			}
		}
		return w.joinRelations(left, right, n)
	case *ast.UnnestExpr:
		// UNNEST(array_expr) [AS alias] [WITH OFFSET [AS o]] — the array argument's
		// columns stay predicate-ish references (and nested subqueries are walked),
		// AND the source resolves to a value relation projecting the element
		// column(s) with the argument's lineage.
		w.walkPredicate(n.Array)
		return w.finishUnnestRelation(w.unnestElementColumns(n.Array), n.Alias, n.WithOffset, n.WithOffsetAlias)
	}
	return nil
}

// joinRelations combines two join sides into one relation, mirroring the legacy
// joinTable: the anchor (left) side's columns are rebuilt as NON-plain fields
// (the legacy join rewraps them without IsPlainField), while the right side's
// columns keep their plainness. The combined relation is unnamed (a `rel.*` /
// column lookup uses the leaf relations, not the join wrapper).
//
// A `USING (keys)` join coalesces the key columns (fix F1): the real GoogleSQL
// `SELECT *` output is [each key ONCE (its value reads from BOTH sides), then
// the left side's non-key columns, then the right side's]. The legacy resolver
// only achieved this for upper-case-written keys (its key map was keyed on the
// written spelling but probed with the upper-cased field name); omni coalesces
// case-insensitively — the GoogleSQL identifier rule, and strictly safer (the
// legacy lowercase concatenation shifted positions, a masking leak). Keys-first
// ordering matches the engine; legacy emitted the coalesced key at the LEFT
// side's key position, which agrees whenever the key is the left table's first
// column (the corpus shape).
//
// NATURAL joins fail closed: both engines reject them at analysis, and the
// shared-column set is unknowable without a catalog — silently concatenating
// would misalign positions (structural rule: correct lineage or an error).
// A USING side that cannot be name-partitioned without metadata (a deferred
// set-op marker in its projection) likewise fails closed.
func (w *spanWalker) joinRelations(left, right *relation, join *ast.JoinExpr) *relation {
	if join != nil && join.Natural {
		w.failClosed(fmt.Errorf("NATURAL JOIN cannot be resolved to column lineage without catalog metadata (fail closed)"))
	} else if join != nil && len(join.Using) > 0 && left != nil && right != nil {
		if rel, ok := coalesceUsingJoin(left, right, join.Using); ok {
			return rel
		}
		w.failClosed(fmt.Errorf("JOIN USING (%s) over a deferred set-operation projection cannot be name-partitioned without catalog metadata (fail closed)", strings.Join(join.Using, ", ")))
	}
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	cols := make([]projColumn, 0, len(left.columns)+len(right.columns))
	for _, c := range left.columns {
		c.plain = false
		cols = append(cols, c)
	}
	cols = append(cols, right.columns...)
	return &relation{columns: cols}
}

// buildTableExpr resolves a non-join FROM source into a relation: a base table
// path (a single baseStar projection the consumer expands), a CTE reference (the
// CTE's retained resolved projection), a parenthesized subquery (its resolved
// projection), or a TVF / correlated array-path (no resolvable projection). It
// performs the same AccessTable / predicate side-effects as before: a base path
// is recorded, a subquery/TVF body is walked, a correlated array path contributes
// its root as a predicate, and an alias is registered in the FROM scope.
func (w *spanWalker) buildTableExpr(n *ast.TableExpr) *relation {
	if n == nil {
		return nil
	}
	var rel *relation
	switch {
	case n.Subquery != nil:
		// A FROM subquery: resolve its projection and expose it under the alias so a
		// `SELECT *` / `alias.col` over it reproduces the subquery's resolved columns.
		cols := w.visitBody(n.Subquery)
		w.registerFromAlias(n.Alias)
		rel = &relation{name: strings.ToLower(n.Alias), columns: cols}
		w.addLeaf(rel)
	case n.Func != nil:
		// Table-valued function: its arguments may reference columns and
		// subqueries. The argument columns are predicate-position references
		// (they parameterize the TVF), so collect them as predicates AND recurse
		// any subquery tables. The TVF's OUTPUT columns are unknown without a
		// function signature, so it contributes no resolved projection — matching the
		// legacy extractor, which rejected TVFs outright for lack of return-column
		// info.
		w.walkPredicate(n.Func)
		w.registerFromAlias(n.Alias)
		rel = nil
	case n.Path != nil:
		// A multi-part path whose ROOT matches an earlier FROM alias/table in this
		// SELECT is a correlated array/field unnest (`FROM T t, t.arr`), NOT a base
		// table — the legacy extractor treats the array-path source as a non-table
		// (it rejected UNNEST/array sources). Record its root as a predicate column
		// (the lineage still reflects the array column it walks) AND build the same
		// value relation an explicit UNNEST source gets (fix F2): the element column
		// carries the array path's resolved lineage, so `SELECT *` keeps its
		// positions and the element alias resolves. finishUnnestRelation registers
		// the source's alias so a chained correlated path off it (`t.arr a,
		// a.child`) is likewise recognized.
		if len(n.Path.Parts) >= 2 && w.fromAliases[strings.ToLower(n.Path.Parts[0])] {
			w.addPredicateColumn(columnRefFromParts(n.Path.Parts, w.dialect))
			_, sources := w.resolvePath(n.Path.Parts)
			elem := projColumn{name: n.Path.Parts[len(n.Path.Parts)-1], sources: sources}
			rel = w.finishUnnestRelation([]projColumn{elem}, n.Alias, n.WithOffset, n.WithOffsetAlias)
		} else {
			rel = w.buildPathRelation(n.Path, n.Alias)
		}
	}
	w.walkPredicate(n.SystemTime)
	return rel
}

// unnestElementColumns resolves an UNNEST(...) call's element column(s): one
// projColumn per array argument, named by the argument path's last component
// when the argument is a plain column path (the GoogleSQL implicit-alias rule;
// "" otherwise), with the argument's resolved column refs as lineage. A literal
// array yields a column with RESOLVED-empty lineage (correct — it reads no
// table). A named zip-mode argument (`mode => …`) is an option, not an array.
func (w *spanWalker) unnestElementColumns(array ast.Node) []projColumn {
	fc, ok := array.(*ast.FuncCall)
	if !ok {
		_, sources := w.resolveExprSources(array)
		return []projColumn{{sources: sources}}
	}
	var cols []projColumn
	for _, arg := range fc.Args {
		if _, named := arg.(*ast.NamedArg); named {
			continue
		}
		name := ""
		if parts := exprToParts(arg); len(parts) > 0 {
			name = parts[len(parts)-1]
		}
		_, sources := w.resolveExprSources(arg)
		cols = append(cols, projColumn{name: name, sources: sources})
	}
	if len(cols) == 0 {
		cols = []projColumn{{}}
	}
	return cols
}

// finishUnnestRelation completes an UNNEST / correlated-array-path value
// relation (fix F2): the explicit alias names a single element column (and the
// relation); `WITH OFFSET [AS o]` appends the positional companion column —
// named o, or GoogleSQL's default name `offset` — with NO lineage (it is a
// position, not data). The relation is registered as a leaf (so `alias.field` /
// the bare element name resolve) and its name enters the FROM scope (so a
// chained correlated path off it is recognized). The element columns are not
// plain fields (an array element is derived data, not a base-column
// passthrough).
func (w *spanWalker) finishUnnestRelation(cols []projColumn, alias string, withOffset bool, offsetAlias string) *relation {
	if alias != "" && len(cols) == 1 {
		cols[0].name = alias
	}
	if withOffset {
		name := offsetAlias
		if name == "" {
			name = "offset"
		}
		cols = append(cols, projColumn{name: name})
	}
	relName := alias
	if relName == "" && len(cols) > 0 {
		relName = cols[0].name
	}
	rel := &relation{name: strings.ToLower(relName), valueTable: true, columns: cols}
	w.addLeaf(rel)
	w.registerFromAlias(relName)
	return rel
}

// buildPathRelation resolves a base-table FROM path into a relation and records
// the access. A bare path naming an in-scope CTE reproduces that CTE's retained
// resolved projection (renamed to the alias / CTE name); a real physical table
// becomes a relation with a single baseStar projection element the metadata-aware
// consumer expands. recordPath keeps the existing AccessTable + FROM-scope
// side-effects (a CTE reference is excluded from AccessTables there).
func (w *spanWalker) buildPathRelation(path *ast.PathExpr, alias string) *relation {
	w.recordPath(path, alias)
	if path == nil || len(path.Parts) == 0 {
		return nil
	}
	cat, db, schema, table := bucketNameParts(path.Parts, w.dialect)
	refName := alias
	if refName == "" {
		refName = table
	}
	// A bare name matching an in-scope CTE reproduces the CTE's resolved projection.
	if cat == "" && db == "" && schema == "" {
		if cols, ok := w.scope.cteColumns(table); ok {
			rel := &relation{name: strings.ToLower(refName), columns: cols}
			w.addLeaf(rel)
			return rel
		}
	}
	// A physical base table: omni cannot enumerate its columns, so its projection is
	// a single baseStar the consumer expands. Plain (a base-table column passthrough
	// is a plain field) until a join-left/set-op/explicit-select rewraps it.
	baseRef := ColumnRef{Catalog: cat, Database: db, Schema: schema, Table: table}
	rel := &relation{
		name:    strings.ToLower(refName),
		isBase:  true,
		baseRef: baseRef,
		columns: []projColumn{{name: "*", plain: true, baseStar: &baseRef}},
	}
	w.addLeaf(rel)
	return rel
}

// addLeaf registers a leaf FROM relation in the current SELECT's lookup set
// (w.leafRels), used to resolve `rel.col` / `rel.*` / bare-column references. A
// nil or unnamed relation is skipped (a join wrapper is unnamed — its sides are
// registered individually).
func (w *spanWalker) addLeaf(rel *relation) {
	if rel == nil || rel.name == "" {
		return
	}
	w.leafRels = append(w.leafRels, rel)
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
		w.visitBody(n.Query)
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
		// UPDATE … FROM sources: walk for AccessTable / predicate side-effects (the
		// resolved relation is not needed — UPDATE has no result projection).
		w.buildFromItem(src)
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
	// The USING source is a TableExpr (a table path or a ( query ) subquery). Walk
	// it for AccessTable / predicate side-effects (MERGE has no result projection).
	if src, ok := n.Source.(*ast.TableExpr); ok {
		w.buildTableExpr(src)
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
		w.visitBody(n.AsQuery)
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
		w.visitBody(n.AsQuery)
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

	ta := TableAccess{
		Catalog:  cat,
		Database: db,
		Schema:   schema,
		Table:    table,
		Alias:    alias,
		Loc:      loc,
	}
	key := tableKey{
		Catalog:  ta.Catalog,
		Database: ta.Database,
		Schema:   ta.Schema,
		Table:    ta.Table,
		Alias:    ta.Alias,
	}

	// Bare name matching an in-scope CTE: a CTE reference, NOT a physical table. It
	// is recorded separately (CTEReferences) — excluded from AccessTables so it is
	// not column-expanded or counted in the user/system mix — but still surfaced so
	// the consumer can union it into the table-level SourceColumns (the legacy
	// access-table listener recorded every FROM table path, CTE references
	// included). Deduplicated against both sets so the same name is not recorded
	// twice (a CTE reference and a same-keyed physical access cannot both exist for
	// one key here — a bare in-scope CTE name always resolves to the CTE).
	if cat == "" && db == "" && schema == "" && w.scope.isCTE(table) {
		if w.accessed[key] {
			return
		}
		w.accessed[key] = true
		w.span.CTEReferences = append(w.span.CTEReferences, ta)
		return
	}

	ta.IsSystem = w.tableIsSystem(ta)
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
//   - BigQuery: X is Schema iff X==INFORMATION_SCHEMA, else X is Database
//     (dataset) — the legacy listener's INFORMATION_SCHEMA-only schema rule.
//   - Spanner: X is always Schema (named schema under one DB).
//   - 3 parts (X.Y.table):
//   - BigQuery: Y is Schema iff Y==INFORMATION_SCHEMA (then X is Database),
//     else Y is Database (dataset) and X is Catalog (project) — the
//     project.dataset.table shape.
//   - Spanner: Y is Schema, X is Database (db.schema.table — Spanner has no
//     catalog/project layer, so X is Database, not Catalog).
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
	onParts   func([]string) // raw dotted-name parts of a column reference (for relation-aware resolution); nil to ignore
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
		ew.emitParts([]string{e.Name})
	case *ast.PathExpr:
		ew.emit(columnRefFromParts(e.Parts, ew.w.dialect))
		ew.emitParts(e.Parts)
	case *ast.FieldAccess:
		// A dotted name (t.a / a.b.c) is a FieldAccess chain over a leading
		// Identifier. Flatten it into a single qualified ColumnRef; if the base is
		// not a plain name chain (field access on a function result, etc.), fall
		// back to walking the base expression.
		if parts := flattenFieldAccess(e); parts != nil {
			ew.emit(columnRefFromParts(parts, ew.w.dialect))
			ew.emitParts(parts)
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
		// Quantified form (`expr op {ANY|SOME|ALL} <rhs>`) carries the RHS in the
		// Quant* fields with Right nil — mirror the LikeExpr/InExpr handling so a
		// subquery (or list/UNNEST) RHS still contributes its tables/columns.
		ew.walkList(e.QuantValues)
		ew.walk(e.QuantUnnest)
		ew.walk(e.QuantSubquery)
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

// emitParts invokes onParts with a column reference's raw dotted-name components
// (relation-aware resolution needs the flat parts, not the dialect-bucketed
// ColumnRef, so a 2-part `rel.col` is resolved as relation.column rather than
// dataset.column). A no-op when onParts is unset or parts is empty.
func (ew *exprWalk) emitParts(parts []string) {
	if ew.onParts == nil || len(parts) == 0 {
		return
	}
	ew.onParts(append([]string{}, parts...))
}

// followSubquery recurses into a parsed subquery body when this pass follows
// subqueries. The body is a *QueryStmt (or set-op / select) the parser already
// filled, so the walk reaches it directly.
func (ew *exprWalk) followSubquery(query ast.Node) {
	if !ew.followSub || query == nil {
		return
	}
	// A subquery in an expression position is walked for its tables/predicates;
	// its output columns are discarded (a SELECT-list scalar subquery is opaque
	// to lineage — see makeColumnInfo).
	ew.w.visitBody(query)
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
// dropped; qualified and free names still pass through. Both collection hooks are
// filtered — onColumn (predicate/source ColumnRefs) and onParts (the relation-
// aware resolver's raw dotted paths) — so a bound variable is excluded from BOTH
// the column set and the resolved select-item lineage. When neither hook is set
// (or there are no bound names) the body is walked unchanged.
func (ew *exprWalk) walkBoundBody(bound []string, body ast.Node) {
	if (ew.onColumn == nil && ew.onParts == nil) || len(bound) == 0 {
		ew.walk(body)
		return
	}
	names := make(map[string]bool, len(bound))
	for _, b := range bound {
		if b != "" {
			names[strings.ToLower(b)] = true
		}
	}
	inner := &exprWalk{w: ew.w, followSub: ew.followSub}
	if ew.onColumn != nil {
		outer := ew.onColumn
		inner.onColumn = func(ref ColumnRef) {
			if names[strings.ToLower(rootName(ref))] {
				return
			}
			outer(ref)
		}
	}
	if ew.onParts != nil {
		outer := ew.onParts
		inner.onParts = func(parts []string) {
			if len(parts) > 0 && names[strings.ToLower(parts[0])] {
				return
			}
			outer(parts)
		}
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
