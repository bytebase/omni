// Package catalog — query.go defines the analyzed-query IR for MySQL.
//
// This file is the Phase 0 deliverable of the MySQL semantic layer effort.
// See: docs/plans/2026-04-09-mysql-semantic-layer.md
//
// # Purpose
//
// These types are the post-analysis (semantically resolved) form of a SELECT
// statement, parallel to PostgreSQL's `Query` / `RangeTblEntry` / `Var` /
// `TargetEntry` family. They are produced by `AnalyzeSelectStmt` (Phase 1)
// and consumed by:
//
//   - bytebase MySQL query span (column-level lineage), via an external adapter
//   - mysql/deparse Phase 4: deparse from IR back to canonical SQL
//   - mysql/diff and mysql/sdl Phase 5: structural schema diffing
//
// The IR follows PG's shape so algorithms (lineage walking, dependency
// extraction, deparse) can be ported across engines with minimal change.
// MySQL-specific divergences are documented inline next to each affected type.
//
// # Status
//
// Phase 0 = TYPE DEFINITIONS ONLY. No analyzer, no methods beyond interface
// tags and trivial accessors. Code in this file is not consumed by any
// production path yet; the purpose is to be reviewed (user + cc + codex) and
// locked down before any consumer is built.
//
// This is **Revision 2** (2026-04-10), incorporating sub-agent review feedback.
// See `docs/plans/2026-04-09-mysql-semantic-layer.md` § 9 for the change log.
//
// # Naming convention: the `Q` suffix
//
// All Query-IR struct types in this file end in `Q` (for "Query IR"), with
// two exceptions:
//
//  1. `Query` itself does not get a Q suffix — it is the IR's namespace
//     anchor and `QueryQ` would be tautological.
//  2. The `AnalyzedExpr` interface does not get a Q suffix — interfaces are
//     category labels, not IR nodes.
//
// Enum types do NOT get a Q suffix (they are values, not structures), with
// one mechanical exception: `JoinTypeQ` collides with `mysql/ast.JoinType`
// and is forced to take the suffix.
//
// The rule mirrors PG's convention but applies it uniformly within the
// MySQL IR, rather than only to types that collide with the parser AST.
// Rationale: in `mysql/catalog/analyze.go` (Phase 1) and below, Q-suffixed
// names provide an instant visual cue that this is the analyzed namespace,
// not the catalog state machine (`Table`, `Column`, etc.) and not the parser
// AST (`mysql/ast.SelectStmt`, etc.).
//
// # ResolvedType (vs the plan doc's `ColumnType`)
//
// The plan doc § 4.1 named the type `ColumnType`. That name collides with the
// existing field `Column.ColumnType` (string) on `mysql/catalog/table.go:63`,
// which holds the full type text (e.g. `"varchar(100)"`). We use `ResolvedType`
// here to avoid shadowing and to better describe the role: the type as
// resolved by the analyzer at output time.
//
// # PG cross-reference key
//
// Each type carries a `// pg:` reference to the PG header / file it mirrors.
// When in doubt about semantics, consult `pg/catalog/query.go` for the
// translated equivalent.
package catalog

// =============================================================================
// Section 1 — Top-level Query
// =============================================================================

// Query is the analyzed form of a SELECT (and, in later phases, of an
// INSERT / UPDATE / DELETE).
//
// pg: src/include/nodes/parsenodes.h — Query
//
// Field ordering parallels pg/catalog/query.go::Query for ease of review;
// MySQL-specific fields (SQLMode, LockClause) are appended last.
type Query struct {
	// CommandType discriminates SELECT / INSERT / UPDATE / DELETE.
	// Phase 1 only emits CmdSelect; the field exists so the IR shape is stable
	// when DML analysis is added in Phase 8+.
	CommandType CmdType

	// TargetList is the analyzed SELECT list, in output order. ResNo on each
	// entry is 1-based. ResJunk entries (helper columns synthesized for
	// ORDER BY references) appear after non-junk entries.
	TargetList []*TargetEntryQ

	// RangeTable is the flattened set of input relations referenced by this
	// query level: base tables, FROM-clause subqueries, CTE references, and
	// JOIN result rows. Indices into this slice are referenced by VarExprQ
	// (RangeIdx) and JoinTreeQ (RTIndex on RangeTableRefQ / JoinExprNodeQ).
	//
	// Order is the discovery order during FROM-clause walking. Once analyzer
	// is done, the slice is immutable.
	RangeTable []*RangeTableEntryQ

	// JoinTree describes the structure of the FROM clause and the WHERE qual.
	// FromList contains top-level FROM items; nested JOINs become JoinExprNodeQ.
	// JoinTree.Quals is the WHERE expression.
	//
	// For set-operation queries (SetOp != SetOpNone), JoinTree is set to an
	// empty FromList with nil Quals. The "real" inputs are in LArg/RArg.
	JoinTree *JoinTreeQ

	// GroupClause is the GROUP BY list. Each entry references a TargetEntryQ
	// by index (1-based) into TargetList. Plain GROUP BY only — Phase 1 does
	// not model ROLLUP / CUBE / GROUPING SETS (deferred).
	GroupClause []*SortGroupClauseQ

	// HavingQual is the HAVING expression, nil if absent.
	HavingQual AnalyzedExpr

	// SortClause is the ORDER BY list. Same indexing convention as GroupClause.
	SortClause []*SortGroupClauseQ

	// LimitCount and LimitOffset are the LIMIT / OFFSET expressions.
	// MySQL syntax: `LIMIT [offset,] count` or `LIMIT count OFFSET offset`.
	// Both forms are normalized into these two fields.
	LimitCount  AnalyzedExpr
	LimitOffset AnalyzedExpr

	// Distinct is true for SELECT DISTINCT.
	// Note: PG has a DistinctOn slice for `SELECT DISTINCT ON (...)`;
	// MySQL does not have that construct, so it is not modeled here.
	Distinct bool

	// Set operations: when SetOp != SetOpNone, LArg and RArg hold the analyzed
	// arms and TargetList describes the *result* shape of the set op (column
	// names, output types). RangeTable and JoinTree are empty.
	//
	// AllSetOp distinguishes UNION from UNION ALL, etc. It applies to *this*
	// level's SetOp only — nested set ops carry their own AllSetOp on
	// LArg/RArg recursively.
	SetOp    SetOpType
	AllSetOp bool
	LArg     *Query
	RArg     *Query

	// CTEList holds WITH-clause CTEs declared at this query level. Each CTE
	// is referenced from the RangeTable by an RTECTE entry whose CTEIndex
	// indexes into this slice.
	CTEList     []*CommonTableExprQ
	IsRecursive bool // WITH RECURSIVE

	// WindowClause holds named window declarations from the WINDOW clause:
	//   SELECT ... OVER w FROM t WINDOW w AS (PARTITION BY ...)
	// References to named windows from FuncCallExprQ.Over carry only the Name
	// field; the analyzer resolves the reference against this slice.
	// Inline window definitions (`OVER (PARTITION BY ...)`) are stored
	// directly on FuncCallExprQ.Over and not added here.
	WindowClause []*WindowDefQ

	// HasAggs is true if the analyzer found at least one aggregate call in
	// TargetList / HavingQual. Used by deparse and by GROUP BY validation.
	// (Computed via FuncCallExprQ.IsAggregate, not by string-matching names.)
	HasAggs bool

	// LockClause holds FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE if present.
	// MySQL-specific: nil for SELECTs without locking clauses.
	LockClause *LockingClauseQ

	// SQLMode captures the session sql_mode at analyze time. The same SELECT
	// can have different semantics under different sql_mode values
	// (ANSI_QUOTES affects identifier quoting, ONLY_FULL_GROUP_BY affects
	// validation, PIPES_AS_CONCAT affects `||`, etc.). Required for
	// round-trip deparse fidelity.
	SQLMode SQLMode
}

// CmdType discriminates the kind of analyzed statement.
//
// pg: src/include/nodes/nodes.h — CmdType
type CmdType int

const (
	CmdSelect CmdType = iota
	CmdInsert
	CmdUpdate
	CmdDelete
)

// =============================================================================
// Section 2 — Target list
// =============================================================================

// TargetEntryQ represents one column in the SELECT list.
//
// pg: src/include/nodes/primnodes.h — TargetEntry
type TargetEntryQ struct {
	// Expr is the analyzed select-list expression. For `SELECT t.a`, this is
	// a *VarExprQ; for `SELECT a + 1`, this is an *OpExprQ; etc.
	Expr AnalyzedExpr

	// ResNo is the 1-based output position. Helper (junk) columns are
	// numbered after non-junk columns. Kept as a separate field (rather than
	// implied by slice index) because the analyzer may reorder TargetList
	// during set-operation column unification or DISTINCT processing.
	//
	// Note: PG's TargetEntry.resno is int16; we use plain int for simplicity.
	ResNo int

	// ResName is the output column name as it would appear in a result set
	// header. Either the user-provided alias (`SELECT a AS x`), the original
	// column name (`SELECT a`), or the synthesized expression text
	// (`SELECT a+1` → `a+1`).
	ResName string

	// ResJunk marks helper columns synthesized for ORDER BY references that
	// don't appear in the user-visible select list. These are stripped by
	// deparse for the SELECT list output but referenced by SortClause.
	ResJunk bool

	// Provenance — populated when Expr is a single VarExprQ that resolves to
	// a physical table column (or a chain of VarExprQs through views/CTEs
	// that bottoms out at one). Used by:
	//   - SDL view column derivation (Phase 2)
	//   - lineage shortcut path (Phase 2 in-test walker)
	//
	// Empty when the column is computed (`a+1`, `COUNT(*)`, `CASE ...`).
	// Multi-source provenance (e.g. COALESCE over two columns) is the
	// lineage walker's responsibility, not this shortcut field's.
	ResOrigDB    string
	ResOrigTable string
	ResOrigCol   string
}

// SortGroupClauseQ references a TargetEntryQ from GROUP BY / ORDER BY / DISTINCT.
//
// pg: src/include/nodes/parsenodes.h — SortGroupClause
//
// MySQL note: PG's SortGroupClause carries equality- and sort-operator OIDs;
// we don't (no operator OID space). Collation is captured on the underlying
// VarExprQ / expression rather than here.
type SortGroupClauseQ struct {
	// TargetIdx is the 1-based index into Query.TargetList.
	// Note: PG's TLESortGroupRef is 0-based; we standardize on 1-based here
	// to match TargetEntryQ.ResNo and avoid arithmetic shifts.
	TargetIdx int

	// Descending controls sort order.
	Descending bool

	// NullsFirst is preserved for symmetry with PG; MySQL has no
	// `NULLS FIRST` / `NULLS LAST` syntax. The analyzer always sets it to
	// MySQL's default behavior:
	//   - ASC  → NULLs sort first
	//   - DESC → NULLs sort last
	NullsFirst bool
}

// =============================================================================
// Section 3 — Range table
// =============================================================================

// RTEKind discriminates the variant of a RangeTableEntryQ.
//
// pg: src/include/nodes/parsenodes.h — RTEKind
type RTEKind int

const (
	// RTERelation is a base table or view referenced from FROM.
	//
	// For views: the analyzer sets IsView=true and ViewAlgorithm. The view's
	// underlying body is NOT substituted into Subquery here. Consumers that
	// want lineage transparency through MERGE views call the (Phase 2)
	// helper `(*Query).ExpandMergeViews()` to perform substitution at consume
	// time. This keeps deparse-of-original-text correct (SHOW CREATE VIEW
	// must reproduce the user's text, not the inlined body) and lets each
	// consumer decide whether to expand.
	RTERelation RTEKind = iota

	// RTESubquery is a subquery in FROM (`FROM (SELECT ...) AS x`).
	// Subquery holds the analyzed inner Query. DBName and TableName are
	// empty; the user-visible name is in Alias / ERef.
	RTESubquery

	// RTEJoin is the synthetic RTE created for the result of a JOIN
	// expression. ColNames holds the *coalesced* column list (NATURAL JOIN
	// and USING merge same-named columns); the underlying tables remain in
	// the RangeTable as separate RTEs referenced by JoinExprNodeQ.Left/Right.
	RTEJoin

	// RTECTE is a reference to a CTE declared in WITH. CTEIndex indexes into
	// the enclosing Query.CTEList; the CTE body itself is
	// `Query.CTEList[i].Query`.
	RTECTE

	// RTEFunction is a function-in-FROM clause. In MySQL the only such
	// construct is `JSON_TABLE(...)` (8.0.19+). FuncExprs holds the analyzed
	// function call expression(s). Phase 1 analyzer rejects this kind with
	// an "unsupported" error; full implementation is deferred to Phase 8+.
	// The kind exists in the IR now so callers can dispatch on it.
	RTEFunction
)

// RangeTableEntryQ represents one entry in Query.RangeTable.
//
// pg: src/include/nodes/parsenodes.h — RangeTblEntry
//
// # Identification
//
// MySQL has no schema namespace, so we identify base relations by the
// (DBName, TableName) pair instead of by an OID. See plan doc decision D1.
//
// # Per-Kind Field Applicability
//
//	Field           RTERelation  RTESubquery  RTEJoin  RTECTE  RTEFunction
//	-------------   -----------  -----------  -------  ------  -----------
//	DBName              ✓            -          -         -          -
//	TableName           ✓            -          -         -          -
//	Alias               ✓            ✓          -         ✓          ✓
//	ERef                ✓            ✓          ✓         ✓          ✓
//	ColNames            ✓            ✓          ✓         ✓          ✓
//	ColTypes            ✓            ✓          ✓         ✓          ✓
//	ColCollations       ✓            ✓          ✓         ✓          ✓
//	Subquery            -            ✓          -         ✓ (mirror) -
//	JoinType            -            -          ✓         -          -
//	JoinUsing           -            -          ✓         -          -
//	CTEIndex            -            -          -         ✓          -
//	CTEName             -            -          -         ✓          -
//	Lateral             -            ✓          -         -          ✓
//	IsView              ✓            -          -         -          -
//	ViewAlgorithm       ✓            -          -         -          -
//	FuncExprs           -            -          -         -          ✓
//
// Fields not applicable to a kind are zero-valued. The analyzer must enforce
// this; consumers may assume it.
type RangeTableEntryQ struct {
	Kind RTEKind

	// DBName / TableName — populated for RTERelation only (the underlying
	// base table or view). For other kinds these are empty; the user-visible
	// name is in Alias / ERef.
	DBName    string
	TableName string

	// Alias is the user-provided alias (`FROM t AS x` → "x"). Empty if none.
	Alias string

	// ERef is the *effective reference name* used to qualify columns from
	// this RTE. It is Alias if non-empty, else TableName, else a synthesized
	// name for unaliased subqueries. Always populated.
	ERef string

	// Column catalog — present for ALL kinds. The analyzer populates it
	// differently per kind:
	//   - RTERelation: from the catalog (table/view definition)
	//   - RTESubquery: from Subquery.TargetList (non-junk columns)
	//   - RTECTE:      from the CTE body's TargetList
	//   - RTEJoin:     from the join's coalesced column list
	//   - RTEFunction: from the function's declared output schema (JSON_TABLE)
	//
	// Contract: the three slices are parallel and equal-length to ColNames.
	// In Phase 1, ColTypes and ColCollations entries are nil/empty —
	// populated in Phase 3. Consumers must tolerate nil entries until then.
	ColNames      []string
	ColTypes      []*ResolvedType // parallel to ColNames; entries may be nil in Phase 1
	ColCollations []string        // parallel to ColNames; entries may be empty in Phase 1

	// Subquery is populated for RTESubquery, and mirrors the CTE body for
	// RTECTE (for convenience when walking).
	//
	// IMPORTANT: For RTERelation referring to a view, Subquery is NOT
	// populated by the analyzer. View body expansion happens at consume time
	// via (*Query).ExpandMergeViews() — see RTERelation doc.
	Subquery *Query

	// JoinType / JoinUsing — populated for RTEJoin only. Mirrors the
	// corresponding JoinExprNodeQ for convenience during lineage walking.
	//
	// Invariant: JoinType here MUST equal the corresponding
	// JoinExprNodeQ.JoinType, and JoinUsing MUST equal JoinExprNodeQ.UsingClause.
	// The analyzer maintains this; consumers may assume it.
	JoinType  JoinTypeQ
	JoinUsing []string // USING column names, in source order

	// CTEIndex / CTEName — populated for RTECTE.
	// CTEIndex is the index into Query.CTEList of the referenced CTE.
	CTEIndex int
	CTEName  string

	// Lateral marks an RTESubquery (or RTEFunction) as LATERAL, allowing it
	// to reference columns from earlier FROM items in the same FROM clause.
	// MySQL 8.0.14+.
	Lateral bool

	// View metadata — populated when an RTERelation references a view.
	//   IsView=false, ViewAlgorithm=ViewAlgNone   → not a view
	//   IsView=true,  ViewAlgorithm=ViewAlgMerge  → MERGE view
	//   IsView=true,  ViewAlgorithm=ViewAlgTemptable → TEMPTABLE view
	//   IsView=true,  ViewAlgorithm=ViewAlgUndefined → UNDEFINED (MySQL chooses)
	// Consumers MUST check IsView before interpreting ViewAlgorithm — the
	// zero value (ViewAlgNone) means "not a view", not "undefined view".
	IsView        bool
	ViewAlgorithm ViewAlgorithm

	// FuncExprs — populated for RTEFunction only (JSON_TABLE in MySQL 8.0.19+).
	// Holds the analyzed function call expression(s). Phase 1 analyzer
	// rejects RTEFunction; this field exists for forward compatibility.
	FuncExprs []AnalyzedExpr
}

// ViewAlgorithm mirrors the MySQL `ALGORITHM` clause on CREATE VIEW.
//
// MySQL doc: https://dev.mysql.com/doc/refman/8.0/en/view-algorithms.html
type ViewAlgorithm int

const (
	// ViewAlgNone is the zero value used when the RTE is not a view at all.
	// IsView=false implies ViewAlgorithm == ViewAlgNone.
	ViewAlgNone ViewAlgorithm = iota

	// ViewAlgUndefined — MySQL chooses MERGE if possible, else TEMPTABLE.
	// User wrote no ALGORITHM clause, or wrote ALGORITHM=UNDEFINED explicitly.
	ViewAlgUndefined

	// ViewAlgMerge — view body is rewritten into the referencing query.
	// Lineage-transparent: callers of ExpandMergeViews() will see through it.
	ViewAlgMerge

	// ViewAlgTemptable — view is materialized into a temporary table.
	// Lineage-opaque: ExpandMergeViews() does not expand TEMPTABLE views.
	ViewAlgTemptable
)

// =============================================================================
// Section 4 — Join tree
// =============================================================================

// JoinTreeQ describes the FROM clause and WHERE clause structure.
//
// pg: src/include/nodes/primnodes.h — FromExpr
type JoinTreeQ struct {
	// FromList holds the top-level FROM items, each one a JoinNode.
	// `FROM a, b, c` produces three RangeTableRefQ entries; nested JOINs
	// produce JoinExprNodeQ entries.
	//
	// Empty for set-operation queries (Query.SetOp != SetOpNone).
	FromList []JoinNode

	// Quals is the analyzed WHERE expression, nil if no WHERE clause.
	// Always nil for set-operation queries.
	Quals AnalyzedExpr
}

// JoinNode is the interface for items in a JoinTreeQ's FromList and for the
// children of a JoinExprNodeQ. Implementations: *RangeTableRefQ, *JoinExprNodeQ.
type JoinNode interface {
	joinNodeTag()
}

// RangeTableRefQ is a leaf in the join tree — a reference to a single RTE.
//
// pg: src/include/nodes/primnodes.h — RangeTblRef
type RangeTableRefQ struct {
	// RTIndex is the 0-based index into the enclosing Query.RangeTable.
	RTIndex int
}

func (*RangeTableRefQ) joinNodeTag() {}

// JoinExprNodeQ represents a JOIN expression in the FROM clause.
//
// pg: src/include/nodes/primnodes.h — JoinExpr
//
// The join itself produces an RTEJoin entry in the enclosing Query.RangeTable;
// RTIndex is the index of that synthetic RTE. The join's input rows come from
// Left and Right (each a RangeTableRefQ or another JoinExprNodeQ).
//
// Invariant: this node's JoinType MUST equal the corresponding
// `Query.RangeTable[RTIndex].JoinType`, and UsingClause MUST equal that RTE's
// JoinUsing. The analyzer maintains both copies in sync.
type JoinExprNodeQ struct {
	JoinType    JoinTypeQ
	Left        JoinNode
	Right       JoinNode
	Quals       AnalyzedExpr // ON expression, nil for CROSS / NATURAL / USING-only
	UsingClause []string     // USING (col, col, ...) — empty if not used
	Natural     bool         // NATURAL JOIN flag
	RTIndex     int          // index of this join's synthetic RTE in RangeTable
}

func (*JoinExprNodeQ) joinNodeTag() {}

// JoinTypeQ discriminates the kind of JOIN.
//
// MySQL note: includes STRAIGHT_JOIN (MySQL-only optimizer hint that forces
// left-to-right join order). The hint affects optimizer behavior but not
// lineage; deparse must preserve it.
//
// Naming: this enum is `JoinTypeQ` (with Q) because `mysql/ast.JoinType`
// already exists. The Q suffix is mechanically forced here, not stylistic.
type JoinTypeQ int

const (
	JoinInner JoinTypeQ = iota
	JoinLeft
	JoinRight
	JoinCross
	JoinStraight // STRAIGHT_JOIN — MySQL-only
	// Note: MySQL does NOT support FULL OUTER JOIN.
)

// =============================================================================
// Section 5 — CTE
// =============================================================================

// CommonTableExprQ is an analyzed CTE declared in a WITH clause.
//
// pg: src/include/nodes/parsenodes.h — CommonTableExpr
type CommonTableExprQ struct {
	// Name is the CTE name (`WITH x AS (...)` → "x").
	Name string

	// ColumnNames is the optional explicit column rename list
	// (`WITH x(a, b) AS (...)`). Empty if not specified.
	ColumnNames []string

	// Query is the analyzed body. For recursive CTEs, the analyzer handles
	// the self-reference by giving the CTE its own RTECTE entry visible to
	// its own body during analysis.
	Query *Query

	// Recursive marks WITH RECURSIVE CTEs.
	Recursive bool
}

// =============================================================================
// Section 6 — Set operations
// =============================================================================

// SetOpType discriminates the kind of set operation in a Query.
//
// MySQL gained INTERSECT and EXCEPT in 8.0.31; UNION has been supported
// since the beginning. The All distinction (UNION vs UNION ALL) is carried
// on Query.AllSetOp rather than as separate enum values.
type SetOpType int

const (
	SetOpNone SetOpType = iota
	SetOpUnion
	SetOpIntersect
	SetOpExcept
)

// =============================================================================
// Section 7 — Locking clause
// =============================================================================

// LockingClauseQ is the analyzed FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE.
//
// pg: src/include/nodes/parsenodes.h — LockingClause (PG's syntax differs)
type LockingClauseQ struct {
	Strength LockStrength

	// Tables is the OF list, if specified. Empty means "all tables in FROM".
	//
	// Constraint: Tables MUST be empty when Strength == LockInShareMode
	// (the legacy syntax does not support OF). Analyzer enforces.
	Tables []string

	// WaitPolicy controls behavior on lock contention.
	//
	// Constraint: WaitPolicy MUST be LockWaitDefault when
	// Strength == LockInShareMode (the legacy syntax does not support
	// NOWAIT or SKIP LOCKED). Analyzer enforces.
	WaitPolicy LockWaitPolicy
}

// LockStrength enumerates lock modes.
//
// IMPORTANT: LockForShare and LockInShareMode are NOT synonyms despite
// targeting the same lock semantics. They differ in their syntactic
// envelope:
//
//   - FOR SHARE (8.0+): supports OF tbl_name, NOWAIT, SKIP LOCKED
//   - LOCK IN SHARE MODE (legacy): does NOT support any of those modifiers
//
// Both enum values exist so deparse can faithfully reproduce the user's
// original syntax.
type LockStrength int

const (
	LockNone        LockStrength = iota
	LockForUpdate                // FOR UPDATE
	LockForShare                 // FOR SHARE (8.0+, supports OF/NOWAIT/SKIP LOCKED)
	LockInShareMode              // LOCK IN SHARE MODE (legacy, no modifier support)
)

// LockWaitPolicy enumerates the wait behavior on lock contention.
type LockWaitPolicy int

const (
	LockWaitDefault    LockWaitPolicy = iota // block until lock acquired (no NOWAIT/SKIP LOCKED)
	LockWaitNowait                           // NOWAIT — error on contention
	LockWaitSkipLocked                       // SKIP LOCKED — silently skip locked rows
)

// =============================================================================
// Section 8 — AnalyzedExpr interface and implementations
// =============================================================================

// AnalyzedExpr is the interface implemented by all post-analysis expression
// nodes. It exposes the resolved type and collation of the expression's
// result, plus an unexported tag method to close the interface.
//
// pg: src/include/nodes/primnodes.h — Expr (base node)
//
// Naming: this interface does NOT carry a Q suffix because interfaces are
// category labels rather than IR nodes themselves. PG also uses
// `AnalyzedExpr`.
//
// In Phase 1, exprType() returns nil for most node kinds because the
// analyzer does not yet do type inference. Phase 3 fills these in.
// Consumers must tolerate nil return values until Phase 3.
//
// Two exceptions where exprType() is non-nil even in Phase 1:
//   - BoolExprQ.exprType() returns BoolType (the package-level singleton)
//   - NullTestExprQ.exprType() returns BoolType
type AnalyzedExpr interface {
	// exprType returns the resolved result type. May be nil in Phase 1/2 for
	// most node kinds; always BoolType for boolean-result nodes.
	exprType() *ResolvedType

	// exprCollation returns the resolved result collation name (e.g.
	// "utf8mb4_0900_ai_ci"). Empty string in Phase 1/2.
	exprCollation() string

	// analyzedExprTag closes the interface to this package's types.
	analyzedExprTag()
}

// BoolType is the package-level singleton for boolean-result expressions.
// MySQL has no native BOOLEAN type; TINYINT(1) is the conventional encoding.
// Defined here so BoolExprQ and NullTestExprQ can return a non-nil exprType
// even in Phase 1, eliminating a class of nil-checks in consumers.
//
// Note that this uses BaseTypeTinyIntBool, not BaseTypeTinyInt — see
// MySQLBaseType § 10 for the rationale.
var BoolType = &ResolvedType{
	BaseType: BaseTypeTinyIntBool,
}

// ----------------------------------------------------------------------------
// VarExprQ — resolved column reference. The single most important node kind.
// ----------------------------------------------------------------------------

// VarExprQ is a resolved column reference: an `(RangeIdx, AttNum)` coordinate
// into the enclosing Query's RangeTable.
//
// pg: src/include/nodes/primnodes.h — Var
//
// Lineage extraction (Phase 2 in-test walker, bytebase adapter) traverses
// the TargetList collecting VarExprQs and resolves each one through the
// RangeTable:
//
//   - RTERelation: terminal — emit (DBName, TableName, ColNames[AttNum-1])
//   - RTESubquery: recurse into Subquery.TargetList[AttNum-1].Expr
//   - RTECTE:      recurse into the CTE body's TargetList[AttNum-1].Expr
//   - RTEJoin:     dispatch to the underlying Left/Right RTE based on which
//                  side the column was coalesced from
//
// For RTERelation that is a view (IsView=true), the walker may call
// (*Query).ExpandMergeViews() first to make MERGE views transparent.
type VarExprQ struct {
	// RangeIdx is the 0-based index into Query.RangeTable.
	RangeIdx int

	// AttNum is the 1-based column number within
	// `Query.RangeTable[RangeIdx].ColNames`. Matches PG convention so lineage
	// algorithms port without arithmetic shifts.
	AttNum int

	// LevelsUp distinguishes correlated subquery references:
	//   0 = column from this query level
	//   1 = column from immediate enclosing query
	//   ...
	LevelsUp int

	// Type and Collation — populated in Phase 3+. Phase 1 leaves them nil/empty.
	// The field is named Type rather than ResolvedType for brevity at use sites
	// (`v.Type` reads naturally; `v.ResolvedType` is verbose).
	Type      *ResolvedType
	Collation string
}

func (*VarExprQ) analyzedExprTag()          {}
func (v *VarExprQ) exprType() *ResolvedType { return v.Type }
func (v *VarExprQ) exprCollation() string   { return v.Collation }

// ----------------------------------------------------------------------------
// ConstExprQ — typed constant.
// ----------------------------------------------------------------------------

// ConstExprQ is a literal value with a resolved type.
//
// pg: src/include/nodes/primnodes.h — Const
type ConstExprQ struct {
	Type      *ResolvedType
	Collation string
	IsNull    bool
	// Value is the textual form as the lexer produced it. Quoting and
	// numeric format are preserved so deparse can reproduce the original.
	Value string
}

func (*ConstExprQ) analyzedExprTag()          {}
func (c *ConstExprQ) exprType() *ResolvedType { return c.Type }
func (c *ConstExprQ) exprCollation() string   { return c.Collation }

// ----------------------------------------------------------------------------
// FuncCallExprQ — scalar / aggregate / window function call.
// ----------------------------------------------------------------------------

// FuncCallExprQ is a resolved function invocation. Aggregates and window
// functions are NOT separated into their own node types (PG has Aggref and
// WindowFunc); instead, IsAggregate distinguishes aggregates and Over
// distinguishes window functions. This matches the MySQL parser's choice and
// keeps the IR smaller.
//
// pg: src/include/nodes/primnodes.h — FuncExpr / Aggref / WindowFunc (merged)
type FuncCallExprQ struct {
	// Name is the canonical lowercase function name (e.g. "concat", "count").
	Name string

	// Args are the analyzed arguments in source order.
	Args []AnalyzedExpr

	// IsAggregate marks this call as an aggregate function (COUNT, SUM, AVG,
	// MIN, MAX, GROUP_CONCAT, etc.). Set by the analyzer based on a function
	// catalog, NOT by string-matching `Name` at consumer sites. Lineage
	// walkers and HasAggs validation use this field directly.
	IsAggregate bool

	// Distinct marks aggregates with DISTINCT, e.g. COUNT(DISTINCT x).
	// Always false for non-aggregate functions.
	Distinct bool

	// Over is the analyzed OVER clause for window functions; nil for plain
	// scalar / aggregate calls. The analyzer sets Over for any call followed
	// by an OVER clause in the source. Window functions are detected by
	// `Over != nil`, not by IsAggregate (window funcs are NOT aggregates
	// even though some aggregates can be used as window funcs).
	Over *WindowDefQ

	// ResultType — populated in Phase 3 from the function return type table.
	// Phase 1 leaves nil.
	ResultType *ResolvedType
	Collation  string
}

func (*FuncCallExprQ) analyzedExprTag()          {}
func (f *FuncCallExprQ) exprType() *ResolvedType { return f.ResultType }
func (f *FuncCallExprQ) exprCollation() string   { return f.Collation }

// ----------------------------------------------------------------------------
// OpExprQ — binary / unary operator expression.
// ----------------------------------------------------------------------------

// OpExprQ is a resolved operator application. Operators are represented by
// their canonical text (`+`, `=`, `LIKE`, `IS DISTINCT FROM`) rather than by
// an OID — MySQL has no operator OID space, and PG-style operator
// overloading is absent.
//
// Note: `LIKE x ESCAPE y` is NOT modeled here (the ESCAPE adds a third
// operand). It is deferred — see "deferred expression nodes" below.
//
// pg: src/include/nodes/primnodes.h — OpExpr
type OpExprQ struct {
	Op string

	// Left is nil for prefix unary operators (e.g. `-x`).
	// `NOT x` uses BoolExprQ, not OpExprQ.
	Left  AnalyzedExpr
	Right AnalyzedExpr

	ResultType *ResolvedType
	Collation  string
}

func (*OpExprQ) analyzedExprTag()          {}
func (o *OpExprQ) exprType() *ResolvedType { return o.ResultType }
func (o *OpExprQ) exprCollation() string   { return o.Collation }

// ----------------------------------------------------------------------------
// BoolExprQ — AND / OR / NOT.
// ----------------------------------------------------------------------------

// BoolExprQ is a logical combinator.
//
// pg: src/include/nodes/primnodes.h — BoolExpr
type BoolExprQ struct {
	Op   BoolOpType
	Args []AnalyzedExpr // single arg for NOT
}

// BoolOpType enumerates AND / OR / NOT.
type BoolOpType int

const (
	BoolAnd BoolOpType = iota
	BoolOr
	BoolNot
)

func (*BoolExprQ) analyzedExprTag()         {}
func (*BoolExprQ) exprType() *ResolvedType  { return BoolType }
func (*BoolExprQ) exprCollation() string    { return "" }

// ----------------------------------------------------------------------------
// CaseExprQ — CASE WHEN.
// ----------------------------------------------------------------------------

// CaseExprQ is both forms of CASE:
//
//	simple:   CASE x WHEN 1 THEN ... WHEN 2 THEN ... ELSE ... END
//	searched: CASE WHEN cond1 THEN ... WHEN cond2 THEN ... ELSE ... END
//
// TestExpr is non-nil for the simple form, nil for the searched form.
//
// pg: src/include/nodes/primnodes.h — CaseExpr
type CaseExprQ struct {
	TestExpr   AnalyzedExpr // nil for searched CASE
	Args       []*CaseWhenQ
	Default    AnalyzedExpr // ELSE branch; nil if absent
	ResultType *ResolvedType
	Collation  string
}

// CaseWhenQ is one WHEN/THEN arm of a CASE expression.
//
// pg: src/include/nodes/primnodes.h — CaseWhen
type CaseWhenQ struct {
	Cond AnalyzedExpr
	Then AnalyzedExpr
}

func (*CaseExprQ) analyzedExprTag()          {}
func (c *CaseExprQ) exprType() *ResolvedType { return c.ResultType }
func (c *CaseExprQ) exprCollation() string   { return c.Collation }

// ----------------------------------------------------------------------------
// CoalesceExprQ — COALESCE / IFNULL.
// ----------------------------------------------------------------------------

// CoalesceExprQ is a COALESCE() call. IFNULL is normalized to a 2-arg
// CoalesceExprQ by the analyzer.
//
// pg: src/include/nodes/primnodes.h — CoalesceExpr
type CoalesceExprQ struct {
	Args       []AnalyzedExpr
	ResultType *ResolvedType
	Collation  string
}

func (*CoalesceExprQ) analyzedExprTag()          {}
func (c *CoalesceExprQ) exprType() *ResolvedType { return c.ResultType }
func (c *CoalesceExprQ) exprCollation() string   { return c.Collation }

// ----------------------------------------------------------------------------
// NullTestExprQ — IS NULL / IS NOT NULL.
// ----------------------------------------------------------------------------

// NullTestExprQ is the IS [NOT] NULL predicate.
//
// pg: src/include/nodes/primnodes.h — NullTest
type NullTestExprQ struct {
	Arg    AnalyzedExpr
	IsNull bool // true = IS NULL, false = IS NOT NULL
}

func (*NullTestExprQ) analyzedExprTag()         {}
func (*NullTestExprQ) exprType() *ResolvedType  { return BoolType }
func (*NullTestExprQ) exprCollation() string    { return "" }

// ----------------------------------------------------------------------------
// SubLinkExprQ — subquery expression (EXISTS, IN, scalar, ALL/ANY).
// ----------------------------------------------------------------------------

// SubLinkExprQ is a subquery used as an expression: `EXISTS (...)`,
// `x IN (SELECT ...)`, `(SELECT ...)`, `x = ANY (...)`.
//
// Note: `x IN (1, 2, 3)` (non-subquery) is NOT a SubLinkExprQ — it is an
// InListExprQ. The two are kept separate so consumers can distinguish them
// without having to peek at the right-hand side.
//
// pg: src/include/nodes/primnodes.h — SubLink
type SubLinkExprQ struct {
	Kind SubLinkKind

	// TestExpr is the left-hand side for ALL / ANY / IN forms; nil for
	// EXISTS and scalar.
	TestExpr AnalyzedExpr

	// Op is the comparison operator for ALL / ANY / IN forms ("=", "<", ...).
	// For SubLinkIn, the comparison is implicitly "=". Empty for EXISTS /
	// scalar.
	Op string

	// Subquery is the analyzed inner SELECT.
	Subquery *Query

	// ResultType is BoolType for EXISTS / IN / ALL / ANY; the inner column's
	// type for scalar. Phase 3 fills it in.
	ResultType *ResolvedType

	// Collation applies to scalar subqueries returning string types. Empty
	// for boolean-result kinds. Phase 3 fills it in.
	Collation string
}

// SubLinkKind discriminates the subquery flavor.
type SubLinkKind int

const (
	SubLinkExists SubLinkKind = iota // EXISTS (...)
	SubLinkScalar                    // (SELECT ...) used as scalar
	SubLinkIn                        // x IN (SELECT ...)
	SubLinkAny                       // x op ANY (SELECT ...)
	SubLinkAll                       // x op ALL (SELECT ...)
)

func (*SubLinkExprQ) analyzedExprTag()          {}
func (s *SubLinkExprQ) exprType() *ResolvedType { return s.ResultType }
func (s *SubLinkExprQ) exprCollation() string   { return s.Collation }

// ----------------------------------------------------------------------------
// InListExprQ — x IN (literal_list).
// ----------------------------------------------------------------------------

// InListExprQ is `x IN (1, 2, 3)` — the non-subquery form. Modeled as its
// own node (rather than lowered to a chain of OR/=) so deparse can faithfully
// reproduce the source syntax and lineage walkers can introspect the list.
//
// pg: src/include/nodes/primnodes.h — ScalarArrayOpExpr (different shape)
type InListExprQ struct {
	Arg     AnalyzedExpr
	List    []AnalyzedExpr
	Negated bool // true for NOT IN
}

func (*InListExprQ) analyzedExprTag()        {}
func (*InListExprQ) exprType() *ResolvedType { return BoolType }
func (*InListExprQ) exprCollation() string   { return "" }

// ----------------------------------------------------------------------------
// BetweenExprQ — x BETWEEN a AND b.
// ----------------------------------------------------------------------------

// BetweenExprQ is `x BETWEEN a AND b` (and the negated `NOT BETWEEN` form).
// Modeled as its own node rather than lowered to `BoolAnd(>=, <=)` so
// deparse can faithfully reproduce the source syntax. The lineage walker
// descends into Arg, Lower, and Upper.
//
// pg: src/parser/parse_expr.c — transformAExprBetween (lowered in PG; we keep)
type BetweenExprQ struct {
	Arg     AnalyzedExpr
	Lower   AnalyzedExpr
	Upper   AnalyzedExpr
	Negated bool // true for NOT BETWEEN
}

func (*BetweenExprQ) analyzedExprTag()        {}
func (*BetweenExprQ) exprType() *ResolvedType { return BoolType }
func (*BetweenExprQ) exprCollation() string   { return "" }

// ----------------------------------------------------------------------------
// RowExprQ — ROW(...) constructor.
// ----------------------------------------------------------------------------

// RowExprQ is a row constructor used in row comparisons:
// `ROW(a, b) = ROW(1, 2)` or implicit row form `(a, b) = (1, 2)`.
//
// pg: src/include/nodes/primnodes.h — RowExpr
type RowExprQ struct {
	Args       []AnalyzedExpr
	ResultType *ResolvedType
}

func (*RowExprQ) analyzedExprTag()          {}
func (r *RowExprQ) exprType() *ResolvedType { return r.ResultType }
func (*RowExprQ) exprCollation() string     { return "" }

// ----------------------------------------------------------------------------
// CastExprQ — explicit CAST / CONVERT.
// ----------------------------------------------------------------------------

// CastExprQ is an *explicit* type cast: `CAST(x AS UNSIGNED)`,
// `CONVERT(x, CHAR)`.
//
// Implicit casts inserted by MySQL during expression evaluation are NOT
// materialized as CastExprQ nodes — see decision D6 in the plan doc. The
// analyzer leaves implicit conversions implicit and lets deparse / lineage
// reason about them as needed. This decision is revisited at end of Phase 3.
//
// pg: src/include/nodes/parsenodes.h — TypeCast (raw form),
//     src/include/nodes/primnodes.h — CoerceViaIO / RelabelType (analyzed form)
type CastExprQ struct {
	Arg        AnalyzedExpr
	TargetType *ResolvedType

	// Collation applies when the cast target is a string type with an
	// explicit COLLATE clause: `CAST(x AS CHAR CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)`.
	Collation string
}

func (*CastExprQ) analyzedExprTag()          {}
func (c *CastExprQ) exprType() *ResolvedType { return c.TargetType }
func (c *CastExprQ) exprCollation() string   { return c.Collation }

// ----------------------------------------------------------------------------
// Deferred expression nodes
// ----------------------------------------------------------------------------
//
// The following expression node kinds are intentionally NOT in the IR yet.
// Each has a known consumer or scenario but is deferred to keep Phase 0/1
// scope manageable. Listed here so reviewers can challenge the omissions
// explicitly:
//
//   - LikeExprQ:        `x LIKE pattern ESCAPE c`. The 3-operand ESCAPE form
//                       does not fit OpExprQ; deferred to Phase 3.
//   - MatchAgainstExprQ: `MATCH (col) AGAINST ('text' IN NATURAL LANGUAGE MODE)`.
//                       MySQL fulltext predicate. Deferred to Phase 4 (deparse
//                       fidelity for fulltext indexes).
//   - JsonExtractExprQ: `col->'$.x'` and `col->>'$.x'` operator forms (vs the
//                       JSON_EXTRACT function form). Deferred to Phase 3 once
//                       we decide whether the analyzer normalizes one form to
//                       the other or preserves both.
//   - AssignmentExprQ:  `SET col = expr` for UPDATE statements. Deferred to
//                       Phase 8+ when DML analysis lands.
//
// If a Phase 1/2 corpus query hits one of these, the analyzer should return
// a clear "unsupported expression kind" error rather than silently dropping
// the node.

// =============================================================================
// Section 9 — Window definition
// =============================================================================

// WindowDefQ is the analyzed OVER clause of a window function, OR a named
// window declaration from the WINDOW clause. The same struct serves both
// roles:
//
//  1. Reference (`OVER w`):
//     Only `Name` is set; PartitionBy/OrderBy/FrameClause are empty.
//     Analyzer resolves the reference against `Query.WindowClause`.
//
//  2. Inline definition (`OVER (PARTITION BY ... ORDER BY ...)`):
//     `Name` is empty; PartitionBy/OrderBy/FrameClause carry the definition.
//
//  3. Named declaration (`WINDOW w AS (PARTITION BY ...)`):
//     Stored in `Query.WindowClause`. Both `Name` and the body fields are set.
//
// In Phase 1 the FrameClause is preserved as raw text — frame parsing is
// nontrivial and not required for lineage / SDL view bodies. Phase 3+ may
// upgrade this to a structured form if needed (Phase 4 deparse round-trip
// will care).
//
// pg: src/include/nodes/parsenodes.h — WindowDef
type WindowDefQ struct {
	// Name is the window's identifier:
	//   - For references: the name being referenced (`OVER w` → "w")
	//   - For inline definitions: empty
	//   - For declarations in WindowClause: the declared name
	Name string

	// PartitionBy and OrderBy are analyzed.
	PartitionBy []AnalyzedExpr
	OrderBy     []*SortGroupClauseQ

	// FrameClause is the raw text of the ROWS / RANGE / GROUPS frame
	// specification, e.g. "ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW".
	// Empty if no frame specified.
	//
	// Phase 4 round-trip caveat: raw text is from the lexer; whitespace and
	// keyword case are not normalized. Round-trip deparse must normalize on
	// output or the SDL diff will see false positives.
	FrameClause string
}

// =============================================================================
// Section 10 — ResolvedType + MySQLBaseType
// =============================================================================

// ResolvedType is the analyzed-time MySQL type representation used by
// expression nodes (VarExprQ.Type, FuncCallExprQ.ResultType, etc.) and by
// RTE column metadata (RangeTableEntryQ.ColTypes).
//
// # Naming
//
// This type was called `ColumnType` in plan doc § 4.1, but the existing
// `Column.ColumnType` (string) field on `mysql/catalog/table.go:63` shadowing
// forced a rename. The semantic role is "the resolved type of an expression
// or column at analyzer output", hence ResolvedType.
//
// # Status
//
// In Phase 1, instances of this type are not produced — the analyzer leaves
// VarExprQ.Type / RangeTableEntryQ.ColTypes nil. Phase 3 begins populating
// it from the function return type table and view column derivation.
//
// # Design principle
//
// Track only the modifiers that change one of:
//   - the value range of the type
//   - the storage layout of the type
//   - the comparison/sort behavior
//   - client-visible semantics (e.g. TINYINT(1)-as-bool)
//
// Modifiers that MySQL itself normalizes away (integer display width,
// ZEROFILL on 8.0.17+) are NOT tracked. Column-level properties
// (Nullable, AutoIncrement, Default, Comment) belong on the Column struct,
// not here.
//
// MySQL doc: https://dev.mysql.com/doc/refman/8.0/en/data-types.html
type ResolvedType struct {
	// BaseType is the type's identity. See MySQLBaseType for the value list.
	BaseType MySQLBaseType

	// Unsigned applies to integer and decimal/float families. Doubles the
	// non-negative range for integers and shifts the range for decimal/float.
	Unsigned bool

	// Length is the parameterized length for fixed/variable-length types:
	//   - VARCHAR(N), CHAR(N), VARBINARY(N), BINARY(N), BIT(N)
	//
	// NOT used for integer types (INT(N) is deprecated display width;
	// MySQL 8.0.17+ strips it. Special-cased TINYINT(1) is represented via
	// BaseType=BaseTypeTinyIntBool, not via Length.)
	//
	// NOT used for TEXT/BLOB families (the size variants are different
	// BaseTypes: TINYTEXT/TEXT/MEDIUMTEXT/LONGTEXT).
	//
	// 0 if unset.
	Length int

	// Precision and Scale are for fixed-point types: DECIMAL(P, S), NUMERIC(P, S).
	// Both 0 if unset.
	Precision int
	Scale     int

	// FSP is the fractional seconds precision for DATETIME(N), TIME(N),
	// TIMESTAMP(N). Range 0–6. Affects both storage size and precision.
	// 0 if unset (which means precision 0, not "absent" — DATETIME and
	// DATETIME(0) are the same).
	FSP int

	// Charset and Collation apply to string-like types and ENUM/SET. Required
	// for SDL diff fidelity (changing a column's collation triggers schema
	// changes in MySQL).
	Charset   string
	Collation string

	// EnumValues holds the value list for ENUM('a', 'b'). Nil for non-ENUM.
	EnumValues []string

	// SetValues holds the value list for SET('a', 'b'). Nil for non-SET.
	// Kept separate from EnumValues for clarity at use sites despite the
	// structural similarity.
	SetValues []string
}

// MySQLBaseType is the typed enumeration of MySQL base type identities.
//
// One value per distinct type as MySQL itself reports them through
// `SHOW COLUMNS` and `information_schema.COLUMNS`. Storage size, value range,
// and client-visible semantics are properties of the BaseType, not of any
// modifier slot in ResolvedType.
//
// # Special case: TINYINT(1)
//
// MySQL 8.0.17 deprecated integer display widths, so `INT(11)` and `INT(5)`
// are normalized to plain `INT`. The single exception is `TINYINT(1)`,
// which MySQL preserves because client libraries (Connector/J, mysqlclient,
// etc.) treat `TINYINT(1)` columns as boolean.
//
// We model this by giving `TINYINT(1)` its own BaseType
// (`BaseTypeTinyIntBool`) rather than carrying a `DisplayWidth` field on
// ResolvedType. This makes the special case explicit at the type level,
// keeps integer handling code free of display-width branches, and lets SDL
// diff distinguish `TINYINT(1)` from `TINYINT` automatically.
//
// MySQL doc: https://dev.mysql.com/doc/refman/8.0/en/data-types.html
type MySQLBaseType int

const (
	BaseTypeUnknown MySQLBaseType = iota

	// Integer family
	BaseTypeTinyInt
	BaseTypeTinyIntBool // TINYINT(1) — client libraries treat as boolean
	BaseTypeSmallInt
	BaseTypeMediumInt
	BaseTypeInt
	BaseTypeBigInt

	// Fixed point
	BaseTypeDecimal // DECIMAL / NUMERIC

	// Floating point
	BaseTypeFloat
	BaseTypeDouble // DOUBLE / REAL

	// Bit
	BaseTypeBit // BIT(N)

	// Date / time
	BaseTypeDate
	BaseTypeDateTime
	BaseTypeTimestamp // distinct from DATETIME (timezone semantics differ)
	BaseTypeTime
	BaseTypeYear

	// Strings
	BaseTypeChar
	BaseTypeVarchar
	BaseTypeBinary
	BaseTypeVarBinary
	BaseTypeTinyText
	BaseTypeText
	BaseTypeMediumText
	BaseTypeLongText
	BaseTypeTinyBlob
	BaseTypeBlob
	BaseTypeMediumBlob
	BaseTypeLongBlob
	BaseTypeEnum
	BaseTypeSet

	// JSON
	BaseTypeJSON

	// Spatial
	BaseTypeGeometry
	BaseTypePoint
	BaseTypeLineString
	BaseTypePolygon
	BaseTypeMultiPoint
	BaseTypeMultiLineString
	BaseTypeMultiPolygon
	BaseTypeGeometryCollection
)

// =============================================================================
// Section 11 — SQL mode bitmap
// =============================================================================

// SQLMode is a bitmap of MySQL `sql_mode` flags relevant to analyzer behavior.
//
// Captured on Query.SQLMode at analyze time so deparse can reproduce the
// exact semantics later. Includes both flags that affect parsing/analysis
// directly and flags that affect later evaluation (the latter are needed
// for SDL diff: changing sql_mode at the session level can alter how a
// view body is interpreted).
//
// MySQL doc: https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html
type SQLMode uint64

const (
	// SQLModeAnsiQuotes — `"x"` is an identifier delimiter, not a string.
	// Affects parsing.
	SQLModeAnsiQuotes SQLMode = 1 << iota

	// SQLModePipesAsConcat — `||` is string concatenation, not boolean OR.
	// Affects expression analysis.
	SQLModePipesAsConcat

	// SQLModeIgnoreSpace — allows whitespace between a function name and the
	// opening parenthesis. With this off, `count (*)` is a column ref + paren
	// expr, not a function call. Directly affects how the parse tree maps to
	// FuncCallExprQ vs other forms.
	SQLModeIgnoreSpace

	// SQLModeHighNotPrecedence — raises NOT precedence above comparison
	// operators. Changes how `NOT a BETWEEN b AND c` is parsed.
	SQLModeHighNotPrecedence

	// SQLModeNoBackslashEscapes — disables backslash escape interpretation in
	// string literals. Affects ConstExprQ.Value semantics.
	SQLModeNoBackslashEscapes

	// SQLModeOnlyFullGroupBy — every non-aggregated select-list column must
	// be functionally dependent on GROUP BY. Affects validation but not
	// parsing.
	SQLModeOnlyFullGroupBy

	// SQLModeNoUnsignedSubtraction — disables unsigned-subtract wraparound.
	SQLModeNoUnsignedSubtraction

	// SQLModeStrictAllTables, SQLModeStrictTransTables — error vs warning on
	// invalid data. Affects DDL CHECK / GENERATED expression evaluation.
	SQLModeStrictAllTables
	SQLModeStrictTransTables

	// SQLModeNoZeroDate / SQLModeNoZeroInDate — restrict zero-valued dates.
	SQLModeNoZeroDate
	SQLModeNoZeroInDate

	// SQLModeRealAsFloat — `REAL` is `FLOAT` instead of `DOUBLE`.
	SQLModeRealAsFloat

	// Add more as the analyzer encounters them. Composite modes (ANSI,
	// TRADITIONAL) are intentionally NOT modeled here — the analyzer
	// expands them at session-capture time.
)

// =============================================================================
// End of Phase 0 IR.
// =============================================================================
