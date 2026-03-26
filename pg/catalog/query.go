package catalog

import "fmt"

// Analyzed query types, translated from PG's parsenodes.h and primnodes.h.
// These represent the post-analysis (semantically resolved) form of a SELECT statement.
//
// pg: src/include/nodes/parsenodes.h — Query
// pg: src/include/nodes/primnodes.h — TargetEntry, Var, Const, FuncExpr, etc.

// --- Shared enumeration types (used by both query types and ruleutils) ---

// SetOpType identifies a set operation.
type SetOpType int

const (
	SetOpNone         SetOpType = iota
	SetOpUnion
	SetOpUnionAll
	SetOpIntersect
	SetOpIntersectAll
	SetOpExcept
	SetOpExceptAll
)

// JoinType identifies the kind of JOIN.
type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinFull
	JoinCross
)

// BoolOpType identifies a boolean operator.
type BoolOpType int

const (
	BoolAnd BoolOpType = iota
	BoolOr
	BoolNot
)

// ResultColumn is the output of SELECT type inference.
//
// pg: src/backend/commands/view.c — DefineVirtualRelation (column metadata)
type ResultColumn struct {
	Name      string
	TypeOID   uint32
	TypeMod   int32
	Collation uint32
}

// Query represents an analyzed SELECT statement.
//
// pg: src/include/nodes/parsenodes.h — Query
type Query struct {
	TargetList   []*TargetEntry
	RangeTable   []*RangeTableEntry
	JoinTree     *JoinTree
	GroupClause  []*SortGroupClause
	HavingQual   AnalyzedExpr
	SortClause   []*SortGroupClause
	LimitCount   AnalyzedExpr
	LimitOffset  AnalyzedExpr
	SetOp        SetOpType       // SetOpNone for simple SELECT
	AllSetOp     bool            // true for ALL variants
	LArg         *Query          // left side of set op
	RArg         *Query          // right side of set op
	Distinct     bool            // SELECT DISTINCT
	DistinctOn   []*SortGroupClause
	HasAggs      bool
	CTEList      []*CommonTableExprQ // WITH clause CTEs
	WindowClause []*WindowClauseQ    // WINDOW clause
	IsRecursive  bool                // WITH RECURSIVE
	GroupingSets []*GroupingSetQ      // GROUPING SETS/ROLLUP/CUBE (nil for plain GROUP BY)
}

// TargetEntry represents a column in the SELECT list.
//
// pg: src/include/nodes/primnodes.h — TargetEntry
type TargetEntry struct {
	Expr            AnalyzedExpr
	ResNo           int16  // position (1-based)
	ResName         string // output column name
	ResJunk         bool   // hidden (ORDER BY helper)
	ResOrigTbl      uint32 // source table OID (for provenance)
	ResOrigCol      int16  // source column number
	ResSortGroupRef uint32 // nonzero if referenced by GROUP BY/ORDER BY
}

// SortGroupClause identifies an output column for ORDER BY / GROUP BY.
//
// pg: src/include/nodes/parsenodes.h — SortGroupClause
type SortGroupClause struct {
	TLESortGroupRef uint32 // index into TargetList (0-based)
	Descending      bool
	NullsFirst      bool
}

// AnalyzedExpr is the interface for post-analysis expression nodes.
//
// pg: src/include/nodes/primnodes.h — Expr (base node)
// pg: src/backend/nodes/nodeFuncs.c — exprType, exprTypmod, exprCollation
type AnalyzedExpr interface {
	exprType() uint32
	exprTypMod() int32
	exprCollation() uint32
}

// VarExpr represents a resolved column reference.
//
// pg: src/include/nodes/primnodes.h — Var
type VarExpr struct {
	RangeIdx  int    // index into Query.RangeTable (0-based)
	AttNum    int16  // column number (1-based)
	TypeOID   uint32 // result type
	TypeMod   int32  // type modifier
	Collation uint32 // collation OID
	LevelsUp  int    // 0 = current query, 1 = parent, etc. (for correlated subqueries)
}

func (v *VarExpr) exprType() uint32      { return v.TypeOID }
func (v *VarExpr) exprTypMod() int32     { return v.TypeMod }
func (v *VarExpr) exprCollation() uint32 { return v.Collation }

// ConstExpr represents a typed constant.
//
// pg: src/include/nodes/primnodes.h — Const
type ConstExpr struct {
	TypeOID   uint32
	TypeMod   int32
	Collation uint32 // collation OID
	IsNull    bool
	Value     string // textual representation
}

func (c *ConstExpr) exprType() uint32      { return c.TypeOID }
func (c *ConstExpr) exprTypMod() int32     { return c.TypeMod }
func (c *ConstExpr) exprCollation() uint32 { return c.Collation }

// FuncCallExpr represents a resolved function call.
//
// pg: src/include/nodes/primnodes.h — FuncExpr
type FuncCallExpr struct {
	FuncOID      uint32
	FuncName     string // for deparse
	ResultType   uint32
	ResultTypMod int32
	Collation    uint32 // result collation
	Args         []AnalyzedExpr
	CoerceFormat byte // 0=normal func, 'e'=explicit cast, 'i'=implicit cast
}

func (f *FuncCallExpr) exprType() uint32      { return f.ResultType }
func (f *FuncCallExpr) exprTypMod() int32     { return f.ResultTypMod }
func (f *FuncCallExpr) exprCollation() uint32 { return f.Collation }

// AggExpr represents an aggregate function call.
//
// pg: src/include/nodes/primnodes.h — Aggref
type AggExpr struct {
	AggFuncOID  uint32
	AggName     string // for deparse
	ResultType  uint32
	Collation   uint32 // result collation
	Args        []AnalyzedExpr
	AggStar     bool // count(*)
	AggDistinct bool
}

func (a *AggExpr) exprType() uint32      { return a.ResultType }
func (a *AggExpr) exprTypMod() int32     { return -1 }
func (a *AggExpr) exprCollation() uint32 { return a.Collation }

// OpExpr represents a resolved operator expression.
//
// pg: src/include/nodes/primnodes.h — OpExpr
type OpExpr struct {
	OpOID      uint32
	OpName     string // for deparse
	ResultType uint32
	Collation  uint32 // result collation
	Left       AnalyzedExpr // nil for prefix
	Right      AnalyzedExpr
}

func (o *OpExpr) exprType() uint32      { return o.ResultType }
func (o *OpExpr) exprTypMod() int32     { return -1 }
func (o *OpExpr) exprCollation() uint32 { return o.Collation }

// RelabelExpr represents a binary-compatible type cast (no-op).
//
// pg: src/include/nodes/primnodes.h — RelabelType
type RelabelExpr struct {
	Arg        AnalyzedExpr
	ResultType uint32
	TypeMod    int32
	Collation  uint32 // result collation
	Format     byte   // 'e'=explicit, 'i'=implicit
}

func (r *RelabelExpr) exprType() uint32      { return r.ResultType }
func (r *RelabelExpr) exprTypMod() int32     { return r.TypeMod }
func (r *RelabelExpr) exprCollation() uint32 { return r.Collation }

// CoerceViaIOExpr represents a type cast via I/O conversion.
//
// pg: src/include/nodes/primnodes.h — CoerceViaIO
type CoerceViaIOExpr struct {
	Arg        AnalyzedExpr
	ResultType uint32
	Collation  uint32 // result collation
	Format     byte   // 'e'=explicit, 'i'=implicit
}

func (c *CoerceViaIOExpr) exprType() uint32      { return c.ResultType }
func (c *CoerceViaIOExpr) exprTypMod() int32     { return -1 }
func (c *CoerceViaIOExpr) exprCollation() uint32 { return c.Collation }

// CaseExprQ represents an analyzed CASE expression.
// Named CaseExprQ to avoid collision with the raw CaseExpr in pgparser nodes.
//
// pg: src/include/nodes/primnodes.h — CaseExpr
type CaseExprQ struct {
	Arg        AnalyzedExpr // nil for searched CASE
	When       []*CaseWhenQ
	Default    AnalyzedExpr // nil if no ELSE
	ResultType uint32
	Collation  uint32 // result collation
}

func (c *CaseExprQ) exprType() uint32      { return c.ResultType }
func (c *CaseExprQ) exprTypMod() int32     { return -1 }
func (c *CaseExprQ) exprCollation() uint32 { return c.Collation }

// CaseWhenQ represents a WHEN clause in a CASE expression.
//
// pg: src/include/nodes/primnodes.h — CaseWhen
type CaseWhenQ struct {
	Condition AnalyzedExpr
	Result    AnalyzedExpr
}

// CoalesceExprQ represents an analyzed COALESCE expression.
//
// pg: src/include/nodes/primnodes.h — CoalesceExpr
type CoalesceExprQ struct {
	Args       []AnalyzedExpr
	ResultType uint32
	Collation  uint32 // result collation
}

func (c *CoalesceExprQ) exprType() uint32      { return c.ResultType }
func (c *CoalesceExprQ) exprTypMod() int32     { return -1 }
func (c *CoalesceExprQ) exprCollation() uint32 { return c.Collation }

// BoolExprQ represents an analyzed boolean expression (AND/OR/NOT).
//
// pg: src/include/nodes/primnodes.h — BoolExpr
type BoolExprQ struct {
	Op   BoolOpType
	Args []AnalyzedExpr
}

func (b *BoolExprQ) exprType() uint32      { return BOOLOID }
func (b *BoolExprQ) exprTypMod() int32     { return -1 }
func (b *BoolExprQ) exprCollation() uint32 { return 0 }

// NullTestExpr represents IS [NOT] NULL.
//
// pg: src/include/nodes/primnodes.h — NullTest
type NullTestExpr struct {
	Arg    AnalyzedExpr
	IsNull bool // true=IS NULL, false=IS NOT NULL
}

func (n *NullTestExpr) exprType() uint32      { return BOOLOID }
func (n *NullTestExpr) exprTypMod() int32     { return -1 }
func (n *NullTestExpr) exprCollation() uint32 { return 0 }

// SubLinkExpr represents a subquery expression (scalar or EXISTS).
//
// pg: src/include/nodes/primnodes.h — SubLink
type SubLinkExpr struct {
	SubLinkType SubLinkType
	TestExpr    AnalyzedExpr // left side of IN/ANY/ALL; nil for EXISTS/scalar
	SubQuery    *Query
	ResultType  uint32
}

func (s *SubLinkExpr) exprType() uint32 { return s.ResultType }
func (s *SubLinkExpr) exprTypMod() int32 {
	if s.SubLinkType == SubLinkExprType && s.SubQuery != nil && len(s.SubQuery.TargetList) > 0 {
		return s.SubQuery.TargetList[0].Expr.exprTypMod()
	}
	return -1
}
func (s *SubLinkExpr) exprCollation() uint32 {
	if s.SubLinkType == SubLinkExprType && s.SubQuery != nil && len(s.SubQuery.TargetList) > 0 {
		return s.SubQuery.TargetList[0].Expr.exprCollation()
	}
	return 0
}

// SubLinkType identifies the kind of subquery.
type SubLinkType int

const (
	SubLinkExprType   SubLinkType = iota // scalar subquery
	SubLinkExistsType                    // EXISTS subquery
	SubLinkAnyType                       // ANY/IN subquery
	SubLinkAllType                       // ALL subquery
)

// NullIfExprQ represents an analyzed NULLIF expression.
// Named NullIfExprQ to avoid collision with pgparser's NullIfExpr.
//
// pg: src/include/nodes/primnodes.h — NullIfExpr (derived from OpExpr)
type NullIfExprQ struct {
	OpOID      uint32
	Args       []AnalyzedExpr // exactly 2 args
	ResultType uint32
}

func (n *NullIfExprQ) exprType() uint32 { return n.ResultType }
func (n *NullIfExprQ) exprTypMod() int32 { return -1 }
func (n *NullIfExprQ) exprCollation() uint32 {
	if len(n.Args) > 0 {
		return n.Args[0].exprCollation()
	}
	return 0
}

// MinMaxExprQ represents GREATEST/LEAST.
//
// pg: src/include/nodes/primnodes.h — MinMaxExpr
type MinMaxExprQ struct {
	Op         MinMaxOp // IS_GREATEST or IS_LEAST
	Args       []AnalyzedExpr
	ResultType uint32
	Collation  uint32 // result collation
}

func (m *MinMaxExprQ) exprType() uint32      { return m.ResultType }
func (m *MinMaxExprQ) exprTypMod() int32     { return -1 }
func (m *MinMaxExprQ) exprCollation() uint32 { return m.Collation }

// MinMaxOp identifies GREATEST vs LEAST.
type MinMaxOp int

const (
	MinMaxGreatest MinMaxOp = iota
	MinMaxLeast
)

// BooleanTestExpr represents IS [NOT] TRUE/FALSE/UNKNOWN.
//
// pg: src/include/nodes/primnodes.h — BooleanTest
type BooleanTestExpr struct {
	Arg      AnalyzedExpr
	TestType BoolTestType
}

func (b *BooleanTestExpr) exprType() uint32      { return BOOLOID }
func (b *BooleanTestExpr) exprTypMod() int32     { return -1 }
func (b *BooleanTestExpr) exprCollation() uint32 { return 0 }

// BoolTestType identifies the kind of boolean test.
type BoolTestType int

const (
	BoolIsTrue BoolTestType = iota
	BoolIsNotTrue
	BoolIsFalse
	BoolIsNotFalse
	BoolIsUnknown
	BoolIsNotUnknown
)

// SQLValueFuncExpr represents CURRENT_DATE, CURRENT_TIMESTAMP, etc.
//
// pg: src/include/nodes/primnodes.h — SQLValueFunction
type SQLValueFuncExpr struct {
	Op      SVFOp
	TypeOID uint32
	TypeMod int32
}

func (s *SQLValueFuncExpr) exprType() uint32  { return s.TypeOID }
func (s *SQLValueFuncExpr) exprTypMod() int32 { return s.TypeMod }
func (s *SQLValueFuncExpr) exprCollation() uint32 {
	// String-typed SVFs (CURRENT_USER etc.) have default collation.
	switch s.TypeOID {
	case NAMEOID, TEXTOID:
		return DEFAULT_COLLATION_OID
	}
	return 0
}

// SVFOp identifies which SQL-standard function this is.
type SVFOp int

const (
	SVFCurrentDate SVFOp = iota
	SVFCurrentTime
	SVFCurrentTimeN
	SVFCurrentTimestamp
	SVFCurrentTimestampN
	SVFLocaltime
	SVFLocaltimeN
	SVFLocaltimestamp
	SVFLocaltimestampN
	SVFCurrentRole
	SVFCurrentUser
	SVFUser
	SVFSessionUser
	SVFCurrentCatalog
	SVFCurrentSchema
)

// DistinctExprQ represents IS [NOT] DISTINCT FROM.
// In PG, DistinctExpr is derived from OpExpr; the deparser checks the node tag.
//
// pg: src/include/nodes/primnodes.h — DistinctExpr
type DistinctExprQ struct {
	OpOID      uint32
	OpName     string
	ResultType uint32
	Left       AnalyzedExpr
	Right      AnalyzedExpr
	IsNot      bool // true for IS NOT DISTINCT FROM
}

func (d *DistinctExprQ) exprType() uint32      { return d.ResultType }
func (d *DistinctExprQ) exprTypMod() int32     { return -1 }
func (d *DistinctExprQ) exprCollation() uint32 { return 0 }

// ScalarArrayOpExpr represents op ANY/ALL (array).
// Used for IN-list expansion: x IN (a,b,c) → x = ANY(ARRAY[a,b,c]).
//
// pg: src/include/nodes/primnodes.h — ScalarArrayOpExpr
type ScalarArrayOpExpr struct {
	OpOID  uint32
	OpName string
	UseOr  bool         // true for ANY (=IN), false for ALL
	Left   AnalyzedExpr // scalar argument
	Right  AnalyzedExpr // array expression (matches PG args[1])
}

func (s *ScalarArrayOpExpr) exprType() uint32      { return BOOLOID }
func (s *ScalarArrayOpExpr) exprTypMod() int32     { return -1 }
func (s *ScalarArrayOpExpr) exprCollation() uint32 { return 0 }

// ArrayExprQ represents an ARRAY[...] constructor.
//
// pg: src/include/nodes/primnodes.h — ArrayExpr
type ArrayExprQ struct {
	ElementType uint32 // element type OID
	ArrayType   uint32 // array type OID
	Elements    []AnalyzedExpr
}

func (a *ArrayExprQ) exprType() uint32  { return a.ArrayType }
func (a *ArrayExprQ) exprTypMod() int32 { return -1 }
func (a *ArrayExprQ) exprCollation() uint32 {
	// Array collation is derived from element collation.
	if len(a.Elements) > 0 {
		return a.Elements[0].exprCollation()
	}
	return 0
}

// RowExprQ represents a ROW(...) or (...) expression.
//
// pg: src/include/nodes/primnodes.h — RowExpr
type RowExprQ struct {
	Args       []AnalyzedExpr
	ResultType uint32 // RECORDOID typically
	RowFormat  byte   // 0=implicit, 'e'=explicit ROW keyword
}

func (r *RowExprQ) exprType() uint32      { return r.ResultType }
func (r *RowExprQ) exprTypMod() int32     { return -1 }
func (r *RowExprQ) exprCollation() uint32 { return 0 }

// CollateExprQ represents an expr COLLATE "name" expression.
//
// pg: src/include/nodes/primnodes.h — CollateExpr
type CollateExprQ struct {
	Arg       AnalyzedExpr
	CollOID   uint32 // collation OID (0 if not resolved)
	CollName  string // collation name for deparse
}

func (c *CollateExprQ) exprType() uint32      { return c.Arg.exprType() }
func (c *CollateExprQ) exprTypMod() int32     { return c.Arg.exprTypMod() }
func (c *CollateExprQ) exprCollation() uint32 { return c.CollOID }

// CoerceToDomainValueExpr represents a reference to the VALUE keyword in
// domain CHECK constraints. PG creates this node in domainAddCheckConstraint()
// via a parser hook that replaces "value" ColumnRef nodes.
//
// pg: src/include/nodes/primnodes.h — CoerceToDomainValue
type CoerceToDomainValueExpr struct {
	TypeOID   uint32
	TypeMod   int32
	Collation uint32
}

func (c *CoerceToDomainValueExpr) exprType() uint32      { return c.TypeOID }
func (c *CoerceToDomainValueExpr) exprTypMod() int32     { return c.TypeMod }
func (c *CoerceToDomainValueExpr) exprCollation() uint32 { return c.Collation }

// FieldSelectExprQ represents a composite.field selection.
//
// pg: src/include/nodes/primnodes.h — FieldSelect
type FieldSelectExprQ struct {
	Arg        AnalyzedExpr
	FieldNum   int16  // attribute number (1-based)
	FieldName  string // for deparse
	ResultType uint32
	TypeMod    int32
	Collation  uint32 // result collation
}

func (f *FieldSelectExprQ) exprType() uint32      { return f.ResultType }
func (f *FieldSelectExprQ) exprTypMod() int32     { return f.TypeMod }
func (f *FieldSelectExprQ) exprCollation() uint32 { return f.Collation }

// WindowFuncExpr represents a window function call.
//
// pg: src/include/nodes/primnodes.h — WindowFunc
type WindowFuncExpr struct {
	FuncOID    uint32
	FuncName   string
	ResultType uint32
	Collation  uint32 // result collation
	Args       []AnalyzedExpr
	AggStar    bool         // e.g. count(*)
	AggFilter  AnalyzedExpr // FILTER (WHERE ...), may be nil
	WinRef     uint32       // index into Query.WindowClause
}

func (w *WindowFuncExpr) exprType() uint32      { return w.ResultType }
func (w *WindowFuncExpr) exprTypMod() int32     { return -1 }
func (w *WindowFuncExpr) exprCollation() uint32 { return w.Collation }

// WindowClauseQ represents an analyzed WINDOW clause entry.
//
// pg: src/include/nodes/parsenodes.h — WindowClause
type WindowClauseQ struct {
	Name         string             // window name (empty if inline OVER)
	PartitionBy  []*SortGroupClause // PARTITION BY
	OrderBy      []*SortGroupClause // ORDER BY
	FrameOptions int                // bitmask
	StartOffset  AnalyzedExpr       // frame start expression
	EndOffset    AnalyzedExpr       // frame end expression
}

// CommonTableExprQ represents an analyzed CTE.
//
// pg: src/include/nodes/parsenodes.h — CommonTableExpr
type CommonTableExprQ struct {
	Name         string
	Aliases      []string // explicit column aliases
	Query        *Query
	Recursive    bool
	Materialized int // 0=default, 1=materialized, 2=not materialized
}

// GroupingSetQ represents GROUPING SETS / CUBE / ROLLUP.
//
// pg: src/include/nodes/parsenodes.h — GroupingSet
type GroupingSetQ struct {
	Kind    GroupingSetKind
	Content []*SortGroupClause // the columns in this set
	Sets    []*GroupingSetQ     // for GROUPING SETS containing nested sets
}

// GroupingSetKind identifies the type of grouping set.
type GroupingSetKind int

const (
	GroupingSetSimple   GroupingSetKind = iota // plain GROUP BY column
	GroupingSetRollup                          // ROLLUP(...)
	GroupingSetCube                            // CUBE(...)
	GroupingSetSets                            // GROUPING SETS((...), (...))
)

// --- Range Table and Join Tree types ---

// RTEKind identifies the kind of range table entry.
type RTEKind int

const (
	RTERelation RTEKind = iota // plain table
	RTESubquery                // subquery in FROM
	RTEJoin                    // join result
	RTECTE                     // CTE reference
	RTEFunction                // function in FROM clause
)

// RangeTableEntry represents an entry in the query's range table.
//
// pg: src/include/nodes/parsenodes.h — RangeTblEntry
type RangeTableEntry struct {
	Kind          RTEKind
	RelOID        uint32    // for RTERelation
	RelName       string    // original table name
	SchemaName    string    // schema name (for deparse)
	Alias         string    // user-provided alias (empty = none)
	ERef          string    // effective reference name for deparse
	ColNames      []string  // column names visible from this RTE
	ColTypes      []uint32  // column type OIDs
	ColTypMods    []int32   // column type modifiers
	ColCollations []uint32  // column collation OIDs
	Subquery      *Query    // for RTESubquery
	JoinType      JoinType  // for RTEJoin
	Lateral       bool      // LATERAL subquery or function
	CTEName       string    // for RTECTE: name of referenced CTE
	CTEIndex      int       // for RTECTE: index into Query.CTEList
	FuncExprs     []AnalyzedExpr // for RTEFunction: function call expressions
	Ordinality    bool           // for RTEFunction: WITH ORDINALITY
	Tablesample   *TablesampleClauseQ // for RTERelation: TABLESAMPLE clause
}

// TablesampleClauseQ represents a TABLESAMPLE clause on a relation.
//
// pg: src/include/nodes/parsenodes.h — TableSampleClause
type TablesampleClauseQ struct {
	Method     string         // sampling method name (e.g. "bernoulli", "system")
	Args       []AnalyzedExpr // arguments for sampling method
	Repeatable AnalyzedExpr   // REPEATABLE expression (nil if none)
}

// JoinTree represents the FROM clause structure.
//
// pg: src/include/nodes/primnodes.h — FromExpr
type JoinTree struct {
	FromList []JoinNode   // top-level FROM items
	Quals    AnalyzedExpr // WHERE condition (may be nil)
}

// JoinNode is the interface for FROM clause items.
type JoinNode interface {
	joinNodeTag()
}

// RangeTableRef references a range table entry.
//
// pg: src/include/nodes/primnodes.h — RangeTblRef
type RangeTableRef struct {
	RTIndex int // index into RangeTable (0-based)
}

func (r *RangeTableRef) joinNodeTag() {}

// JoinExprNode represents a JOIN expression in the FROM clause.
//
// pg: src/include/nodes/primnodes.h — JoinExpr
type JoinExprNode struct {
	JoinType    JoinType
	Left        JoinNode
	Right       JoinNode
	Quals       AnalyzedExpr // ON condition
	UsingClause []string     // USING column names (for deparse)
	RTIndex     int          // this join's RTE index
}

func (j *JoinExprNode) joinNodeTag() {}

// --- Utility functions moved from infer.go (used by production code) ---

// typeCollation returns the default collation OID for a type.
// (pgddl helper — PG uses typcollation from pg_type)
func (c *Catalog) typeCollation(oid uint32) uint32 {
	t := c.typeByOID[oid]
	if t == nil {
		return 0
	}
	if t.Collation != 0 {
		return t.Collation
	}
	// Array types inherit element type's collation.
	// pg: src/backend/catalog/pg_type.c — array types have typcollation=0
	//     but collation is derived from element type
	if t.Category == 'A' && t.Elem != 0 {
		if elem := c.typeByOID[t.Elem]; elem != nil {
			return elem.Collation
		}
	}
	return 0
}

// collationName returns a human-readable name for a collation OID,
// used in error messages for checkViewColumns.
// (pgddl helper — PG uses get_collation_name)
func (c *Catalog) collationName(oid uint32) string {
	switch oid {
	case 100:
		return "default"
	case 950:
		return "C"
	case 0:
		return ""
	default:
		return fmt.Sprintf("collation(%d)", oid)
	}
}

// selectCommonCollation picks a common collation for set operations.
// Simplified: if both are the same, use it; if one is 0 (non-collatable), use the other.
// (pgddl helper — PG uses select_common_collation with full strength/conflict resolution)
func selectCommonCollation(a, b uint32) uint32 {
	if a == b {
		return a
	}
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	// Both non-zero and different: in full PG this would be a conflict.
	// For pgddl, prefer the left side (matches PG's behavior for set-ops where
	// left branch names are used).
	return a
}

// typeName returns the human-readable name for a type OID.
func (c *Catalog) typeName(oid uint32) string {
	if t := c.typeByOID[oid]; t != nil {
		return t.TypeName
	}
	return fmt.Sprintf("oid(%d)", oid)
}

// resolveReturnType resolves polymorphic return types.
func (c *Catalog) resolveReturnType(proc *BuiltinProc, argTypes []uint32) uint32 {
	retType := proc.RetType
	if !isPolymorphic(retType) {
		return retType
	}

	// Find the actual type from polymorphic parameter matching.
	for i, paramOID := range proc.ArgTypes {
		if !isPolymorphic(paramOID) {
			continue
		}
		actualType := argTypes[i]
		if actualType == UNKNOWNOID {
			continue
		}

		switch retType {
		case ANYELEMENTOID, ANYOID, ANYNONARRAYOID, ANYCOMPATIBLEOID:
			return actualType
		case ANYARRAYOID, ANYCOMPATIBLEARRAYOID:
			// Return the array type of the actual arg type.
			if t := c.typeByOID[actualType]; t != nil && t.Array != 0 {
				return t.Array
			}
			return actualType
		}
	}

	// Fallback: if return type is ANYELEMENTOID/ANYOID and we couldn't resolve, return TEXTOID.
	return TEXTOID
}

// SubscriptingRefExpr represents an array subscript expression (arr[1]).
//
// pg: src/include/nodes/primnodes.h — SubscriptingRef
type SubscriptingRefExpr struct {
	ContainerExpr  AnalyzedExpr   // the array expression
	SubscriptExprs []AnalyzedExpr // subscript expressions
	ResultType     uint32         // element type OID
	IsSlice        bool           // true for arr[1:3]
	LowerExprs     []AnalyzedExpr // lower bounds for slice (nil entries = open)
}

func (s *SubscriptingRefExpr) exprType() uint32      { return s.ResultType }
func (s *SubscriptingRefExpr) exprTypMod() int32      { return -1 }
func (s *SubscriptingRefExpr) exprCollation() uint32 {
	if s.ContainerExpr != nil {
		return s.ContainerExpr.exprCollation()
	}
	return 0
}

// NamedArgExprQ represents a named argument in a function call (arg_name => value).
//
// pg: src/include/nodes/primnodes.h — NamedArgExpr
type NamedArgExprQ struct {
	Name   string       // argument name
	Arg    AnalyzedExpr // the argument expression
	ArgNum int          // argument number in positional notation
}

func (n *NamedArgExprQ) exprType() uint32      { return n.Arg.exprType() }
func (n *NamedArgExprQ) exprTypMod() int32     { return n.Arg.exprTypMod() }
func (n *NamedArgExprQ) exprCollation() uint32 { return n.Arg.exprCollation() }

// ArrayCoerceExprQ represents array element type coercion.
//
// pg: src/include/nodes/primnodes.h — ArrayCoerceExpr
type ArrayCoerceExprQ struct {
	Arg        AnalyzedExpr
	ResultType uint32
	ElemExpr   AnalyzedExpr // per-element coercion expression
	Format     byte         // 'e'=explicit, 'i'=implicit
}

func (a *ArrayCoerceExprQ) exprType() uint32      { return a.ResultType }
func (a *ArrayCoerceExprQ) exprTypMod() int32      { return -1 }
func (a *ArrayCoerceExprQ) exprCollation() uint32 {
	if a.Arg != nil {
		return a.Arg.exprCollation()
	}
	return 0
}

// CaseTestExprQ represents the internal CASE test placeholder.
//
// pg: src/include/nodes/primnodes.h — CaseTestExpr
type CaseTestExprQ struct {
	TypeOID uint32
	TypeMod int32
}

func (c *CaseTestExprQ) exprType() uint32      { return c.TypeOID }
func (c *CaseTestExprQ) exprTypMod() int32     { return c.TypeMod }
func (c *CaseTestExprQ) exprCollation() uint32 { return 0 }

// CoerceToDomainExpr represents a coercion to a domain type.
//
// pg: src/include/nodes/primnodes.h — CoerceToDomain
type CoerceToDomainExpr struct {
	Arg        AnalyzedExpr
	ResultType uint32
	TypeMod    int32
}

func (c *CoerceToDomainExpr) exprType() uint32      { return c.ResultType }
func (c *CoerceToDomainExpr) exprTypMod() int32     { return c.TypeMod }
func (c *CoerceToDomainExpr) exprCollation() uint32 {
	if c.Arg != nil {
		return c.Arg.exprCollation()
	}
	return 0
}

// RowCompareExprQ represents a row comparison like (a,b) < (c,d).
//
// pg: src/include/nodes/primnodes.h — RowCompareExpr
type RowCompareExprQ struct {
	RCType int              // comparison type: 1=<, 2=<=, 3=>=, 4=>, 5==, 6=<>
	LArgs  []AnalyzedExpr   // left-side expressions
	RArgs  []AnalyzedExpr   // right-side expressions
	OpNos  []uint32         // operator OIDs per column
}

func (r *RowCompareExprQ) exprType() uint32      { return BOOLOID }
func (r *RowCompareExprQ) exprTypMod() int32     { return -1 }
func (r *RowCompareExprQ) exprCollation() uint32 { return 0 }

// RowCompare type constants matching PG's RowCompareType.
const (
	RowCompareLT int = 1 // <
	RowCompareLE int = 2 // <=
	RowCompareEQ int = 3 // =
	RowCompareGE int = 4 // >=
	RowCompareGT int = 5 // >
	RowCompareNE int = 6 // <>
)

// ParamExpr represents a $N parameter placeholder.
//
// pg: src/include/nodes/primnodes.h — Param
type ParamExpr struct {
	ParamNum  int
	ParamType uint32
}

func (p *ParamExpr) exprType() uint32      { return p.ParamType }
func (p *ParamExpr) exprTypMod() int32     { return -1 }
func (p *ParamExpr) exprCollation() uint32 { return 0 }

// SetToDefaultExpr represents the DEFAULT keyword in expression context.
//
// pg: src/include/nodes/primnodes.h — SetToDefault
type SetToDefaultExpr struct {
	TypeOID uint32
}

func (s *SetToDefaultExpr) exprType() uint32      { return s.TypeOID }
func (s *SetToDefaultExpr) exprTypMod() int32     { return -1 }
func (s *SetToDefaultExpr) exprCollation() uint32 { return 0 }

// NextValueExprQ represents nextval for IDENTITY columns.
//
// pg: src/include/nodes/primnodes.h — NextValueExpr
type NextValueExprQ struct {
	SeqOID  uint32
	SeqName string // for deparse
}

func (n *NextValueExprQ) exprType() uint32      { return INT8OID }
func (n *NextValueExprQ) exprTypMod() int32     { return -1 }
func (n *NextValueExprQ) exprCollation() uint32 { return 0 }

// ConvertRowtypeExprQ represents a row type conversion between compatible composite types.
//
// pg: src/include/nodes/primnodes.h — ConvertRowtypeExpr
type ConvertRowtypeExprQ struct {
	Arg        AnalyzedExpr
	ResultType uint32
}

func (c *ConvertRowtypeExprQ) exprType() uint32      { return c.ResultType }
func (c *ConvertRowtypeExprQ) exprTypMod() int32     { return -1 }
func (c *ConvertRowtypeExprQ) exprCollation() uint32 { return 0 }

// FieldStoreExprQ represents a composite field store (updating a field of a composite value).
//
// pg: src/include/nodes/primnodes.h — FieldStore
type FieldStoreExprQ struct {
	Arg        AnalyzedExpr
	FieldNums  []int
	NewVals    []AnalyzedExpr
	ResultType uint32
}

func (f *FieldStoreExprQ) exprType() uint32      { return f.ResultType }
func (f *FieldStoreExprQ) exprTypMod() int32     { return -1 }
func (f *FieldStoreExprQ) exprCollation() uint32 { return 0 }

// GroupingFuncExpr represents a GROUPING(col1, col2, ...) expression.
//
// pg: src/include/nodes/primnodes.h — GroupingFunc
type GroupingFuncExpr struct {
	Args    []AnalyzedExpr // arguments (column references)
	Refs    []int          // ressortgrouprefs of arguments
	TypeOID uint32
}

func (g *GroupingFuncExpr) exprType() uint32      { return INT4OID }
func (g *GroupingFuncExpr) exprTypMod() int32     { return -1 }
func (g *GroupingFuncExpr) exprCollation() uint32 { return 0 }

// XmlExprOp identifies the type of XML expression.
const (
	XmlOpConcat    = 0 // XMLCONCAT(args)
	XmlOpElement   = 1 // XMLELEMENT(name, args)
	XmlOpForest    = 2 // XMLFOREST(args)
	XmlOpParse     = 3 // XMLPARSE(DOCUMENT|CONTENT expr)
	XmlOpPI        = 4 // XMLPI(name [, args])
	XmlOpRoot      = 5 // XMLROOT(xml, version, standalone)
	XmlOpSerialize = 6 // XMLSERIALIZE(DOCUMENT|CONTENT expr AS type)
	XmlOpDocument  = 7 // expr IS DOCUMENT
)

// XmlOptionType identifies DOCUMENT vs CONTENT.
const (
	XmlOptionDocument = 1
	XmlOptionContent  = 2
)

// XmlExprQ represents various SQL/XML functions.
//
// pg: src/include/nodes/primnodes.h — XmlExpr
type XmlExprQ struct {
	Op         int            // XmlOp* constant
	Name       string         // name in xml(NAME foo ...) syntaxes
	NamedArgs  []AnalyzedExpr // non-XML expressions for xml_attributes
	ArgNames   []string       // parallel list of attribute names
	Args       []AnalyzedExpr // list of expressions
	Xmloption  int            // DOCUMENT or CONTENT
	TypeOID    uint32         // result type
	TypeMod    int32          // result type modifier
}

func (x *XmlExprQ) exprType() uint32      { return x.TypeOID }
func (x *XmlExprQ) exprTypMod() int32     { return x.TypeMod }
func (x *XmlExprQ) exprCollation() uint32 { return 0 }

// isPolymorphic returns true if the given OID is a polymorphic pseudo-type.
func isPolymorphic(oid uint32) bool {
	switch oid {
	case ANYOID, ANYELEMENTOID, ANYARRAYOID, ANYNONARRAYOID, ANYENUMOID,
		ANYRANGEOID, ANYMULTIRANGEOID, ANYCOMPATIBLEOID, ANYCOMPATIBLEARRAYOID,
		ANYCOMPATIBLENONARRAYOID, ANYCOMPATIBLERANGEOID, ANYCOMPATIBLEMULTIRANGEOID:
		return true
	}
	return false
}
