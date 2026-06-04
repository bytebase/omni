package ast

import "strings"

// This file holds the concrete snowflake parse-tree node types. F1 ships
// only the File root container; Tier 1+ migration nodes (identifiers,
// types, expressions, SELECT core, DDL, etc.) populate the rest.
//
// The cmd/genwalker code generator scans this file together with node.go
// to produce walk_generated.go.

// File is the root node of a parsed Snowflake source file. It holds the
// top-level statement list and the byte range covering the entire file.
// F4 (parser-entry) returns *File from Parse.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (f *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// Identifier types
// ---------------------------------------------------------------------------

// Ident represents a single identifier — a name used to reference a database
// object (table, column, schema, etc.).
//
// Name is the raw text from source: for quoted identifiers, the content
// between the double-quotes with "" un-escaped; for unquoted identifiers,
// the source bytes with case preserved.
//
// Quoted reports whether the source used "..." quoting. This matters because
// Snowflake case-folds unquoted identifiers to uppercase at resolution time,
// while quoted identifiers preserve case.
//
// Ident is a value struct, NOT a Node. It is embedded by value in parent
// nodes (e.g. ObjectName) and is not visited by the AST walker.
//
// The zero value (Name == "" && Quoted == false) represents an absent
// identifier — used by ObjectName for unused parts (e.g. a 1-part name
// has zero Database and Schema).
type Ident struct {
	Name   string
	Quoted bool
	Loc    Loc
}

// Normalize returns the canonical form of the identifier per Snowflake
// resolution rules:
//   - Quoted identifiers: returned as-is (case-sensitive)
//   - Unquoted identifiers: uppercased
func (i Ident) Normalize() string {
	if i.Quoted {
		return i.Name
	}
	return strings.ToUpper(i.Name)
}

// String returns the source form of the identifier, re-quoting if it was
// originally quoted. Inner " characters are escaped as "". Useful for
// deparse and error messages.
func (i Ident) String() string {
	if !i.Quoted {
		return i.Name
	}
	return `"` + strings.ReplaceAll(i.Name, `"`, `""`) + `"`
}

// IsEmpty reports whether the identifier is the zero value (absent).
// Used by ObjectName to check whether a part (Database, Schema) is present.
func (i Ident) IsEmpty() bool {
	return i.Name == "" && !i.Quoted
}

// ObjectName represents a qualified object name (1/2/3-part) like
// `table`, `schema.table`, or `database.schema.table`.
//
// For 1-part names, only Name is set (Database and Schema are zero Idents).
// For 2-part names, Schema and Name are set.
// For 3-part names, all three are set.
//
// ObjectName is a Node and can be used as a child in the AST tree. The
// walker visits *ObjectName but does NOT descend into the embedded Ident
// fields (they are value structs, not Nodes).
type ObjectName struct {
	Database Ident // may be zero (IsEmpty)
	Schema   Ident // may be zero (IsEmpty)
	Name     Ident // always present for a valid ObjectName
	Loc      Loc
}

// Tag implements Node.
func (n *ObjectName) Tag() NodeTag { return T_ObjectName }

// Compile-time assertion that *ObjectName satisfies Node.
var _ Node = (*ObjectName)(nil)

// Normalize returns the canonical dotted form with each non-empty part
// normalized per Snowflake resolution rules.
//
// Examples:
//   - 1-part: "TABLE"
//   - 2-part: "SCHEMA.TABLE"
//   - 3-part: "DB.SCHEMA.TABLE"
func (n ObjectName) Normalize() string {
	parts := n.Parts()
	normalized := make([]string, len(parts))
	for i, p := range parts {
		normalized[i] = p.Normalize()
	}
	return strings.Join(normalized, ".")
}

// String returns the source form with dots. Each part is re-quoted if it
// was originally quoted.
func (n ObjectName) String() string {
	parts := n.Parts()
	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = p.String()
	}
	return strings.Join(strs, ".")
}

// Parts returns the non-empty parts in order. Length is 1, 2, or 3.
func (n ObjectName) Parts() []Ident {
	switch {
	case !n.Database.IsEmpty():
		return []Ident{n.Database, n.Schema, n.Name}
	case !n.Schema.IsEmpty():
		return []Ident{n.Schema, n.Name}
	default:
		return []Ident{n.Name}
	}
}

// Matches reports whether this ObjectName suffix-matches other using
// normalized (case-folded) comparison. A 1-part name matches any other
// with the same normalized Name. A 2-part name matches any other with
// the same normalized Schema + Name. A 3-part name requires all three
// parts to match.
func (n ObjectName) Matches(other ObjectName) bool {
	if n.Name.Normalize() != other.Name.Normalize() {
		return false
	}
	if n.Schema.IsEmpty() {
		return true // 1-part match
	}
	if n.Schema.Normalize() != other.Schema.Normalize() {
		return false
	}
	if n.Database.IsEmpty() {
		return true // 2-part match
	}
	return n.Database.Normalize() == other.Database.Normalize()
}

// ---------------------------------------------------------------------------
// Data type types
// ---------------------------------------------------------------------------

// TypeKind classifies Snowflake data types into categories for fast
// switch dispatch by downstream consumers. The Name field on TypeName
// carries the exact source text for round-tripping; Kind carries the
// category for semantic analysis.
type TypeKind int

const (
	TypeUnknown      TypeKind = iota
	TypeInt                   // INT, INTEGER, SMALLINT, TINYINT, BYTEINT, BIGINT
	TypeNumber                // NUMBER, NUMERIC, DECIMAL — may have (precision, scale)
	TypeFloat                 // FLOAT, FLOAT4, FLOAT8, DOUBLE, DOUBLE PRECISION, REAL
	TypeBoolean               // BOOLEAN
	TypeDate                  // DATE
	TypeDateTime              // DATETIME — may have (precision)
	TypeTime                  // TIME — may have (precision)
	TypeTimestamp             // TIMESTAMP — may have (precision)
	TypeTimestampLTZ          // TIMESTAMP_LTZ — may have (precision)
	TypeTimestampNTZ          // TIMESTAMP_NTZ — may have (precision)
	TypeTimestampTZ           // TIMESTAMP_TZ — may have (precision)
	TypeChar                  // CHAR, NCHAR, CHARACTER — may have (length)
	TypeVarchar               // VARCHAR, CHAR VARYING, NCHAR VARYING, NVARCHAR, NVARCHAR2, STRING, TEXT
	TypeBinary                // BINARY — may have (length)
	TypeVarbinary             // VARBINARY — may have (length)
	TypeVariant               // VARIANT
	TypeObject                // OBJECT
	TypeArray                 // ARRAY — may have ElementType
	TypeGeography             // GEOGRAPHY
	TypeGeometry              // GEOMETRY
	TypeVector                // VECTOR — has ElementType + VectorDim
)

// String returns the human-readable name of the TypeKind.
func (k TypeKind) String() string {
	switch k {
	case TypeUnknown:
		return "Unknown"
	case TypeInt:
		return "Int"
	case TypeNumber:
		return "Number"
	case TypeFloat:
		return "Float"
	case TypeBoolean:
		return "Boolean"
	case TypeDate:
		return "Date"
	case TypeDateTime:
		return "DateTime"
	case TypeTime:
		return "Time"
	case TypeTimestamp:
		return "Timestamp"
	case TypeTimestampLTZ:
		return "TimestampLTZ"
	case TypeTimestampNTZ:
		return "TimestampNTZ"
	case TypeTimestampTZ:
		return "TimestampTZ"
	case TypeChar:
		return "Char"
	case TypeVarchar:
		return "Varchar"
	case TypeBinary:
		return "Binary"
	case TypeVarbinary:
		return "Varbinary"
	case TypeVariant:
		return "Variant"
	case TypeObject:
		return "Object"
	case TypeArray:
		return "Array"
	case TypeGeography:
		return "Geography"
	case TypeGeometry:
		return "Geometry"
	case TypeVector:
		return "Vector"
	default:
		return "Unknown"
	}
}

// TypeName represents a Snowflake data type as it appears in SQL source.
//
// Examples:
//
//	INT                  → Kind=TypeInt, Name="INT", Params=nil
//	NUMBER(38, 0)        → Kind=TypeNumber, Name="NUMBER", Params=[38, 0]
//	VARCHAR(100)         → Kind=TypeVarchar, Name="VARCHAR", Params=[100]
//	TIMESTAMP_LTZ(9)     → Kind=TypeTimestampLTZ, Name="TIMESTAMP_LTZ", Params=[9]
//	DOUBLE PRECISION     → Kind=TypeFloat, Name="DOUBLE PRECISION", Params=nil
//	ARRAY(VARCHAR)       → Kind=TypeArray, Name="ARRAY", ElementType=&TypeName{...}
//	VECTOR(INT, 256)     → Kind=TypeVector, Name="VECTOR", ElementType=&TypeName{...}, VectorDim=256
//
// TypeName is a Node. The walker descends into ElementType when non-nil.
type TypeName struct {
	Kind        TypeKind  // classified type category
	Name        string    // source text of the type name for round-tripping
	Params      []int     // numeric type parameters; nil if absent
	ElementType *TypeName // element type for ARRAY and VECTOR; nil otherwise
	VectorDim   int       // dimension for VECTOR(type, dim); -1 if not VECTOR
	Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion that *TypeName satisfies Node.
var _ Node = (*TypeName)(nil)

// ---------------------------------------------------------------------------
// Expression enums
// ---------------------------------------------------------------------------

// BinaryOp enumerates binary operator types.
type BinaryOp int

const (
	BinAdd    BinaryOp = iota // +
	BinSub                    // -
	BinMul                    // *
	BinDiv                    // /
	BinMod                    // %
	BinConcat                 // ||
	BinEq                     // =
	BinNe                     // <> or !=
	BinLt                     // <
	BinGt                     // >
	BinLe                     // <=
	BinGe                     // >=
	BinAnd                    // AND
	BinOr                     // OR
)

// String returns the operator symbol.
func (op BinaryOp) String() string {
	switch op {
	case BinAdd:
		return "+"
	case BinSub:
		return "-"
	case BinMul:
		return "*"
	case BinDiv:
		return "/"
	case BinMod:
		return "%"
	case BinConcat:
		return "||"
	case BinEq:
		return "="
	case BinNe:
		return "!="
	case BinLt:
		return "<"
	case BinGt:
		return ">"
	case BinLe:
		return "<="
	case BinGe:
		return ">="
	case BinAnd:
		return "AND"
	case BinOr:
		return "OR"
	default:
		return "?"
	}
}

// UnaryOp enumerates unary operator types.
type UnaryOp int

const (
	UnaryMinus UnaryOp = iota // -
	UnaryPlus                 // +
	UnaryNot                  // NOT
)

// String returns the operator symbol.
func (op UnaryOp) String() string {
	switch op {
	case UnaryMinus:
		return "-"
	case UnaryPlus:
		return "+"
	case UnaryNot:
		return "NOT"
	default:
		return "?"
	}
}

// LiteralKind classifies literal value types.
type LiteralKind int

const (
	LitNull LiteralKind = iota
	LitBool
	LitInt
	LitFloat
	LitString
)

// LikeOp enumerates LIKE operator variants.
type LikeOp int

const (
	LikeOpLike   LikeOp = iota // LIKE
	LikeOpILike                // ILIKE
	LikeOpRLike                // RLIKE
	LikeOpRegexp               // REGEXP
)

// AccessKind enumerates semi-structured access operator types.
type AccessKind int

const (
	AccessColon   AccessKind = iota // expr:field (JSON path)
	AccessBracket                   // expr[idx] (array subscript)
	AccessDot                       // expr.field (dot path)
)

// WindowFrameKind enumerates window frame types.
type WindowFrameKind int

const (
	FrameRows WindowFrameKind = iota
	FrameRange
	FrameGroups
)

// WindowBoundKind enumerates window frame bound types.
type WindowBoundKind int

const (
	BoundUnboundedPreceding WindowBoundKind = iota
	BoundPreceding
	BoundCurrentRow
	BoundFollowing
	BoundUnboundedFollowing
)

// CaseKind discriminates simple vs searched CASE expressions.
type CaseKind int

const (
	CaseSimple   CaseKind = iota // CASE expr WHEN val THEN ...
	CaseSearched                 // CASE WHEN cond THEN ...
)

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

// Literal represents an integer, float, string, boolean, or NULL literal.
type Literal struct {
	Kind  LiteralKind
	Value string // raw source text for all kinds
	Ival  int64  // for LitInt
	Bval  bool   // for LitBool
	Loc   Loc
}

func (n *Literal) Tag() NodeTag { return T_Literal }

// ColumnRef represents a 1-4 part column reference (col, t.col, s.t.col, db.s.t.col).
type ColumnRef struct {
	Parts []Ident
	Loc   Loc
}

func (n *ColumnRef) Tag() NodeTag { return T_ColumnRef }

// StarExpr represents * or qualifier.* in SELECT lists and COUNT(*).
type StarExpr struct {
	Qualifier *ObjectName // optional table qualifier; nil for bare *
	Loc       Loc
}

func (n *StarExpr) Tag() NodeTag { return T_StarExpr }

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Loc   Loc
}

func (n *BinaryExpr) Tag() NodeTag { return T_BinaryExpr }

// UnaryExpr represents a prefix unary operation: op expr.
type UnaryExpr struct {
	Op   UnaryOp
	Expr Node
	Loc  Loc
}

func (n *UnaryExpr) Tag() NodeTag { return T_UnaryExpr }

// ParenExpr represents a parenthesized expression: ( expr ).
type ParenExpr struct {
	Expr Node
	Loc  Loc
}

func (n *ParenExpr) Tag() NodeTag { return T_ParenExpr }

// CastExpr represents CAST(expr AS type), TRY_CAST(expr AS type), or expr::type.
type CastExpr struct {
	Expr       Node
	TypeName   *TypeName
	TryCast    bool // true for TRY_CAST
	ColonColon bool // true for expr::type syntax
	Loc        Loc
}

func (n *CastExpr) Tag() NodeTag { return T_CastExpr }

// CaseExpr represents a CASE expression (simple or searched).
type CaseExpr struct {
	Kind    CaseKind
	Operand Node          // non-nil for CaseSimple; nil for CaseSearched
	Whens   []*WhenClause // one or more WHEN clauses
	Else    Node          // optional ELSE; nil if absent
	Loc     Loc
}

func (n *CaseExpr) Tag() NodeTag { return T_CaseExpr }

// WhenClause is a single WHEN...THEN branch inside a CaseExpr.
// Not a top-level Node — embedded in CaseExpr.
type WhenClause struct {
	Cond   Node // the WHEN condition (or value for simple CASE)
	Result Node // the THEN result
	Loc    Loc
}

// FuncCallExpr represents a function call including aggregates and window functions.
//
// Star is true for COUNT(*). Distinct is true for COUNT(DISTINCT x).
// OrderBy is used by aggregates with WITHIN GROUP (ORDER BY ...).
// Over is non-nil for window functions (SUM(x) OVER (...)).
type FuncCallExpr struct {
	Name     ObjectName   // function name (may be qualified: schema.func)
	Args     []Node       // argument expressions
	Star     bool         // COUNT(*)
	Distinct bool         // DISTINCT keyword in aggregate
	OrderBy  []*OrderItem // for WITHIN GROUP (ORDER BY ...)
	Over     *WindowSpec  // non-nil for window functions
	Loc      Loc
}

func (n *FuncCallExpr) Tag() NodeTag { return T_FuncCallExpr }

// IffExpr represents Snowflake's IFF(condition, then_expr, else_expr).
type IffExpr struct {
	Cond Node
	Then Node
	Else Node
	Loc  Loc
}

func (n *IffExpr) Tag() NodeTag { return T_IffExpr }

// CollateExpr represents expr COLLATE collation_name.
type CollateExpr struct {
	Expr      Node
	Collation string
	Loc       Loc
}

func (n *CollateExpr) Tag() NodeTag { return T_CollateExpr }

// IsExpr represents expr IS [NOT] NULL or expr IS [NOT] DISTINCT FROM expr.
type IsExpr struct {
	Expr         Node
	Not          bool // IS NOT NULL / IS NOT DISTINCT FROM
	Null         bool // true for IS [NOT] NULL; false for IS [NOT] DISTINCT FROM
	DistinctFrom Node // non-nil for IS [NOT] DISTINCT FROM expr
	Loc          Loc
}

func (n *IsExpr) Tag() NodeTag { return T_IsExpr }

// BetweenExpr represents expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Expr Node
	Low  Node
	High Node
	Not  bool
	Loc  Loc
}

func (n *BetweenExpr) Tag() NodeTag { return T_BetweenExpr }

// InExpr represents expr [NOT] IN (value_list).
// Subquery form (expr IN (SELECT ...)) is handled by T1.4.
type InExpr struct {
	Expr   Node
	Values []Node
	Not    bool
	Loc    Loc
}

func (n *InExpr) Tag() NodeTag { return T_InExpr }

// LikeExpr represents expr [NOT] LIKE/ILIKE/RLIKE/REGEXP pattern [ESCAPE esc].
// Also handles LIKE ANY (...) via the Any + AnyValues fields.
type LikeExpr struct {
	Expr      Node
	Pattern   Node   // the pattern expression
	Escape    Node   // optional ESCAPE clause; nil if absent
	Op        LikeOp // LIKE, ILIKE, RLIKE, REGEXP
	Not       bool
	Any       bool   // true for LIKE ANY (...)
	AnyValues []Node // values for LIKE ANY; nil unless Any is true
	Loc       Loc
}

func (n *LikeExpr) Tag() NodeTag { return T_LikeExpr }

// AccessExpr represents semi-structured data access:
//   - AccessColon:   expr:field (JSON path)
//   - AccessBracket: expr[index] (array subscript)
//   - AccessDot:     expr.field (dot path chaining)
type AccessExpr struct {
	Expr  Node
	Kind  AccessKind
	Field Ident // for Colon and Dot access
	Index Node  // for Bracket access
	Loc   Loc
}

func (n *AccessExpr) Tag() NodeTag { return T_AccessExpr }

// ArrayLiteralExpr represents an array literal: [elem1, elem2, ...].
type ArrayLiteralExpr struct {
	Elements []Node
	Loc      Loc
}

func (n *ArrayLiteralExpr) Tag() NodeTag { return T_ArrayLiteralExpr }

// JsonLiteralExpr represents a JSON object literal: {key: value, ...}.
type JsonLiteralExpr struct {
	Pairs []KeyValuePair
	Loc   Loc
}

func (n *JsonLiteralExpr) Tag() NodeTag { return T_JsonLiteralExpr }

// KeyValuePair is a single key-value pair in a JsonLiteralExpr.
type KeyValuePair struct {
	Key   string
	Value Node
	Loc   Loc
}

// LambdaExpr represents a lambda expression: param -> body or (p1, p2) -> body.
type LambdaExpr struct {
	Params []Ident
	Body   Node
	Loc    Loc
}

func (n *LambdaExpr) Tag() NodeTag { return T_LambdaExpr }

// SubqueryExpr represents a parenthesized subquery expression (SELECT ...).
type SubqueryExpr struct {
	Query Node // the SELECT statement
	Loc   Loc
}

func (n *SubqueryExpr) Tag() NodeTag { return T_SubqueryExpr }

// ExistsExpr represents EXISTS (SELECT ...).
type ExistsExpr struct {
	Query Node
	Loc   Loc
}

func (n *ExistsExpr) Tag() NodeTag { return T_ExistsExpr }

// WindowSpec describes a window specification for OVER (...).
// Not a top-level Node — embedded in FuncCallExpr.
type WindowSpec struct {
	PartitionBy []Node
	OrderBy     []*OrderItem
	Frame       *WindowFrame // nil if no frame clause
	Loc         Loc
}

// OrderItem represents one element in an ORDER BY clause.
type OrderItem struct {
	Expr       Node
	Desc       bool  // true for DESC
	NullsFirst *bool // nil = unspecified, true = NULLS FIRST, false = NULLS LAST
	Loc        Loc
}

// WindowFrame represents ROWS/RANGE/GROUPS BETWEEN start AND end.
type WindowFrame struct {
	Kind  WindowFrameKind
	Start WindowBound
	End   WindowBound // End.Kind may be zero if single-bound form (not BETWEEN)
	Loc   Loc
}

// WindowBound represents one end of a window frame specification.
type WindowBound struct {
	Kind   WindowBoundKind
	Offset Node // for BoundPreceding/BoundFollowing: the N in "N PRECEDING"; nil otherwise
}

// Compile-time assertions.
var (
	_ Node = (*Literal)(nil)
	_ Node = (*ColumnRef)(nil)
	_ Node = (*StarExpr)(nil)
	_ Node = (*BinaryExpr)(nil)
	_ Node = (*UnaryExpr)(nil)
	_ Node = (*ParenExpr)(nil)
	_ Node = (*CastExpr)(nil)
	_ Node = (*CaseExpr)(nil)
	_ Node = (*FuncCallExpr)(nil)
	_ Node = (*IffExpr)(nil)
	_ Node = (*CollateExpr)(nil)
	_ Node = (*IsExpr)(nil)
	_ Node = (*BetweenExpr)(nil)
	_ Node = (*InExpr)(nil)
	_ Node = (*LikeExpr)(nil)
	_ Node = (*AccessExpr)(nil)
	_ Node = (*ArrayLiteralExpr)(nil)
	_ Node = (*JsonLiteralExpr)(nil)
	_ Node = (*LambdaExpr)(nil)
	_ Node = (*SubqueryExpr)(nil)
	_ Node = (*ExistsExpr)(nil)
)

// ---------------------------------------------------------------------------
// Statement nodes
// ---------------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	With     []*CTE          // WITH clause CTEs; nil if absent
	Distinct bool            // SELECT DISTINCT
	All      bool            // SELECT ALL
	Top      Node            // TOP n expression; nil if absent
	Targets  []*SelectTarget // SELECT list items
	From     []Node          // FROM: mixed *TableRef and *JoinExpr; nil if absent
	Where    Node            // WHERE condition; nil if absent
	GroupBy  *GroupByClause  // GROUP BY; nil if absent
	Having   Node            // HAVING condition; nil if absent
	Qualify  Node            // QUALIFY condition; nil if absent (Snowflake-specific)
	OrderBy  []*OrderItem    // ORDER BY; nil if absent
	Limit    Node            // LIMIT n; nil if absent
	Offset   Node            // OFFSET n; nil if absent
	Fetch    *FetchClause    // FETCH FIRST/NEXT; nil if absent
	Loc      Loc
}

func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }

var _ Node = (*SelectStmt)(nil)

// SelectTarget is one item in a SELECT list.
// For expressions: Expr is set, Star is false.
// For star: Star is true, Expr may be a qualifier (table.*) or nil (bare *).
type SelectTarget struct {
	Expr    Node    // expression; nil for bare *
	Alias   Ident   // AS alias; zero Ident if absent
	Star    bool    // true for * or qualifier.*
	Exclude []Ident // EXCLUDE columns; nil if absent
	Loc     Loc
}

// TableRef is a table reference in the FROM clause.
// A TableRef is polymorphic:
//   - Table: Name is set, others nil
//   - Subquery: Subquery is set, Name is nil
//   - Table function: FuncCall is set, Name is nil
//   - Any of the above can have Lateral = true
type TableRef struct {
	Name     *ObjectName   // table name; nil for subquery/func sources
	Alias    Ident         // AS alias; zero if absent
	Subquery Node          // (SELECT ...) in FROM; nil for table refs
	FuncCall *FuncCallExpr // TABLE(func(...)); nil for table refs
	Lateral  bool          // LATERAL prefix
	Loc      Loc
}

func (n *TableRef) Tag() NodeTag { return T_TableRef }

var _ Node = (*TableRef)(nil)

// JoinExpr represents a JOIN between two FROM sources.
type JoinExpr struct {
	Type           JoinType
	Left           Node    // TableRef or nested JoinExpr
	Right          Node    // TableRef or nested JoinExpr
	On             Node    // ON condition; nil for CROSS/NATURAL/USING-only
	Using          []Ident // USING columns; nil if ON or NATURAL
	Natural        bool
	Directed       bool // Snowflake DIRECTED hint
	MatchCondition Node // ASOF MATCH_CONDITION(expr); nil for non-ASOF
	Loc            Loc
}

func (n *JoinExpr) Tag() NodeTag { return T_JoinExpr }

var _ Node = (*JoinExpr)(nil)

// JoinType enumerates the kinds of JOIN.
type JoinType int

const (
	JoinInner JoinType = iota // [INNER] JOIN
	JoinLeft                  // LEFT [OUTER] JOIN
	JoinRight                 // RIGHT [OUTER] JOIN
	JoinFull                  // FULL [OUTER] JOIN
	JoinCross                 // CROSS JOIN
	JoinAsof                  // ASOF JOIN (Snowflake-specific)
)

// CTE represents a Common Table Expression in a WITH clause.
type CTE struct {
	Name      Ident   // CTE name
	Columns   []Ident // optional column aliases
	Query     Node    // the SELECT body (*SelectStmt)
	Recursive bool    // WITH RECURSIVE flag
	Loc       Loc
}

// GroupByClause represents a GROUP BY clause with optional variant.
type GroupByClause struct {
	Kind  GroupByKind
	Items []Node // group-by expressions
	Loc   Loc
}

// GroupByKind enumerates GROUP BY variants.
type GroupByKind int

const (
	GroupByNormal       GroupByKind = iota // GROUP BY a, b
	GroupByCube                            // GROUP BY CUBE (a, b)
	GroupByRollup                          // GROUP BY ROLLUP (a, b)
	GroupByGroupingSets                    // GROUP BY GROUPING SETS ((a), (b))
	GroupByAll                             // GROUP BY ALL
)

// FetchClause represents FETCH FIRST/NEXT n ROWS ONLY.
type FetchClause struct {
	Count Node // the count expression
	Loc   Loc
}

// ---------------------------------------------------------------------------
// Set operator node
// ---------------------------------------------------------------------------

// SetOp enumerates the set operator kinds.
type SetOp int

const (
	SetOpUnion     SetOp = iota // UNION
	SetOpExcept                 // EXCEPT (also MINUS)
	SetOpIntersect              // INTERSECT
)

// SetOperationStmt represents a set-operation query:
// UNION [ALL] [BY NAME] / EXCEPT / INTERSECT between two query expressions.
//
// Left and Right are either *SelectStmt (leaf) or nested *SetOperationStmt
// (chained). The chain is left-associative:
//
//	SELECT 1 UNION SELECT 2 UNION SELECT 3
//	→ SetOperationStmt{Left: SetOperationStmt{Left: S1, Right: S2}, Right: S3}
type SetOperationStmt struct {
	Op     SetOp // the operator kind
	All    bool  // true for UNION ALL
	ByName bool  // true for UNION [ALL] BY NAME (Snowflake-specific)
	Left   Node  // *SelectStmt or nested *SetOperationStmt
	Right  Node  // *SelectStmt or nested *SetOperationStmt
	Loc    Loc
}

func (n *SetOperationStmt) Tag() NodeTag { return T_SetOperationStmt }

var _ Node = (*SetOperationStmt)(nil)

// ---------------------------------------------------------------------------
// Constraint enums
// ---------------------------------------------------------------------------

// ConstraintType enumerates constraint kinds for inline and table-level constraints.
type ConstraintType int

const (
	ConstrPrimaryKey ConstraintType = iota
	ConstrForeignKey
	ConstrUnique
)

// String returns the constraint type name.
func (c ConstraintType) String() string {
	switch c {
	case ConstrPrimaryKey:
		return "PRIMARY KEY"
	case ConstrForeignKey:
		return "FOREIGN KEY"
	case ConstrUnique:
		return "UNIQUE"
	default:
		return "UNKNOWN"
	}
}

// ReferenceAction enumerates FK referential actions.
type ReferenceAction int

const (
	RefActNone       ReferenceAction = iota // not specified
	RefActCascade                           // CASCADE
	RefActSetNull                           // SET NULL
	RefActSetDefault                        // SET DEFAULT
	RefActRestrict                          // RESTRICT
	RefActNoAction                          // NO ACTION
)

// ---------------------------------------------------------------------------
// Helper structs (not Nodes)
// ---------------------------------------------------------------------------

// InlineConstraint represents a column-level constraint.
type InlineConstraint struct {
	Type       ConstraintType
	Name       Ident          // CONSTRAINT name; zero if unnamed
	References *ForeignKeyRef // for FK; nil otherwise
	Loc        Loc
}

// ForeignKeyRef holds REFERENCES clause details.
type ForeignKeyRef struct {
	Table    *ObjectName
	Columns  []Ident
	OnDelete ReferenceAction
	OnUpdate ReferenceAction
	Match    string // "FULL"/"PARTIAL"/"SIMPLE"; empty if absent
}

// IdentitySpec holds IDENTITY/AUTOINCREMENT configuration.
type IdentitySpec struct {
	Start     *int64 // START WITH value; nil if default
	Increment *int64 // INCREMENT BY value; nil if default
	Order     *bool  // true=ORDER, false=NOORDER, nil=unspecified
}

// TagAssignment is a single TAG name = 'value' pair.
type TagAssignment struct {
	Name  *ObjectName
	Value string
}

// CloneSource holds CLONE source with optional time travel.
type CloneSource struct {
	Source   *ObjectName
	AtBefore string // "AT" or "BEFORE"; empty if no time travel
	Kind     string // "TIMESTAMP"/"OFFSET"/"STATEMENT"
	Value    string // the time travel value
}

// ---------------------------------------------------------------------------
// DDL statement nodes
// ---------------------------------------------------------------------------

// CreateTableStmt represents CREATE [OR REPLACE] [TRANSIENT|TEMPORARY|VOLATILE] TABLE ...
type CreateTableStmt struct {
	OrReplace   bool
	Transient   bool
	Temporary   bool
	Volatile    bool
	IfNotExists bool
	Name        *ObjectName
	Columns     []*ColumnDef
	Constraints []*TableConstraint
	ClusterBy   []Node  // CLUSTER BY expressions; nil if absent
	Linear      bool    // CLUSTER BY LINEAR modifier
	Comment     *string // COMMENT = 'text'; nil if absent
	CopyGrants  bool
	Tags        []*TagAssignment // WITH TAG (...); nil if absent
	AsSelect    Node             // CREATE TABLE ... AS SELECT; nil if absent
	Like        *ObjectName      // CREATE TABLE ... LIKE source; nil if absent
	Clone       *CloneSource     // CREATE TABLE ... CLONE source; nil if absent
	Loc         Loc
}

func (n *CreateTableStmt) Tag() NodeTag { return T_CreateTableStmt }

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Name             Ident
	DataType         *TypeName // nil for virtual columns without explicit type
	Default          Node      // DEFAULT expr; nil if absent
	NotNull          bool
	Nullable         bool              // explicit NULL
	Identity         *IdentitySpec     // IDENTITY/AUTOINCREMENT; nil if absent
	Collate          string            // COLLATE 'name'; empty if absent
	MaskingPolicy    *ObjectName       // WITH MASKING POLICY name; nil if absent
	InlineConstraint *InlineConstraint // inline PK/FK/UNIQUE; nil if absent
	Comment          *string           // COMMENT 'text'; nil if absent
	Tags             []*TagAssignment  // WITH TAG (...); nil if absent
	VirtualExpr      Node              // AS (expr); nil if absent
	Loc              Loc
}

func (n *ColumnDef) Tag() NodeTag { return T_ColumnDef }

// TableConstraint represents a table-level constraint (out-of-line).
type TableConstraint struct {
	Type       ConstraintType // ConstrPrimaryKey/ConstrForeignKey/ConstrUnique
	Name       Ident          // CONSTRAINT name; zero if unnamed
	Columns    []Ident        // constrained column names
	References *ForeignKeyRef // FK only; nil otherwise
	Comment    *string        // inline COMMENT 'text'; nil if absent
	Loc        Loc
}

func (n *TableConstraint) Tag() NodeTag { return T_TableConstraint }

// Compile-time assertions.
var (
	_ Node = (*CreateTableStmt)(nil)
	_ Node = (*ColumnDef)(nil)
	_ Node = (*TableConstraint)(nil)
)

// ---------------------------------------------------------------------------
// DATABASE / SCHEMA DDL — enums and helpers
// ---------------------------------------------------------------------------

// AlterDatabaseAction discriminates the action variants of ALTER DATABASE.
type AlterDatabaseAction int

const (
	AlterDBRename             AlterDatabaseAction = iota // RENAME TO
	AlterDBSwap                                          // SWAP WITH
	AlterDBSet                                           // SET <properties>
	AlterDBUnset                                         // UNSET <properties>
	AlterDBSetTag                                        // SET TAG (...)
	AlterDBUnsetTag                                      // UNSET TAG (...)
	AlterDBEnableReplication                             // ENABLE REPLICATION ...
	AlterDBDisableReplication                            // DISABLE REPLICATION ...
	AlterDBEnableFailover                                // ENABLE FAILOVER ...
	AlterDBDisableFailover                               // DISABLE FAILOVER ...
	AlterDBRefresh                                       // REFRESH
	AlterDBPrimary                                       // PRIMARY
)

// AlterSchemaAction discriminates the action variants of ALTER SCHEMA.
type AlterSchemaAction int

const (
	AlterSchemaRename               AlterSchemaAction = iota // RENAME TO
	AlterSchemaSwap                                          // SWAP WITH
	AlterSchemaSet                                           // SET <properties>
	AlterSchemaUnset                                         // UNSET <properties>
	AlterSchemaSetTag                                        // SET TAG (...)
	AlterSchemaUnsetTag                                      // UNSET TAG (...)
	AlterSchemaEnableManagedAccess                           // ENABLE MANAGED ACCESS
	AlterSchemaDisableManagedAccess                          // DISABLE MANAGED ACCESS
)

// DBSchemaProps holds the optional settable properties shared by DATABASE and
// SCHEMA DDL statements.  It is NOT a Node — embedded by value or pointer in
// the statement nodes that need it.
type DBSchemaProps struct {
	DataRetention *int64  // DATA_RETENTION_TIME_IN_DAYS = n; nil if absent
	MaxDataExt    *int64  // MAX_DATA_EXTENSION_TIME_IN_DAYS = n; nil if absent
	DefaultDDLCol *string // DEFAULT_DDL_COLLATION = 'str'; nil if absent
	Comment       *string // COMMENT = 'str'; nil if absent
}

// ---------------------------------------------------------------------------
// DATABASE DDL statement nodes
// ---------------------------------------------------------------------------

// CreateDatabaseStmt represents CREATE [OR REPLACE] [TRANSIENT] DATABASE ...
type CreateDatabaseStmt struct {
	OrReplace   bool
	Transient   bool
	IfNotExists bool
	Name        *ObjectName
	Clone       *CloneSource     // CLONE source [AT|BEFORE (...)]; nil if absent
	Props       DBSchemaProps    // optional properties
	Tags        []*TagAssignment // WITH TAG (...); nil if absent
	Loc         Loc
}

func (n *CreateDatabaseStmt) Tag() NodeTag { return T_CreateDatabaseStmt }

var _ Node = (*CreateDatabaseStmt)(nil)

// AlterDatabaseStmt represents ALTER DATABASE ... (all action variants).
type AlterDatabaseStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterDatabaseAction
	NewName    *ObjectName      // RENAME TO / SWAP WITH target
	SetProps   *DBSchemaProps   // SET properties; non-nil for AlterDBSet
	UnsetProps []string         // UNSET property names
	Tags       []*TagAssignment // SET TAG (...) assignments
	UnsetTags  []*ObjectName    // UNSET TAG (...) names
	Loc        Loc
}

func (n *AlterDatabaseStmt) Tag() NodeTag { return T_AlterDatabaseStmt }

var _ Node = (*AlterDatabaseStmt)(nil)

// DropDatabaseStmt represents DROP DATABASE [IF EXISTS] name [CASCADE|RESTRICT].
type DropDatabaseStmt struct {
	IfExists bool
	Name     *ObjectName
	Cascade  bool
	Restrict bool
	Loc      Loc
}

func (n *DropDatabaseStmt) Tag() NodeTag { return T_DropDatabaseStmt }

var _ Node = (*DropDatabaseStmt)(nil)

// UndropDatabaseStmt represents UNDROP DATABASE name.
type UndropDatabaseStmt struct {
	Name *ObjectName
	Loc  Loc
}

func (n *UndropDatabaseStmt) Tag() NodeTag { return T_UndropDatabaseStmt }

var _ Node = (*UndropDatabaseStmt)(nil)

// ---------------------------------------------------------------------------
// SCHEMA DDL statement nodes
// ---------------------------------------------------------------------------

// CreateSchemaStmt represents CREATE [OR REPLACE] [TRANSIENT] SCHEMA ...
type CreateSchemaStmt struct {
	OrReplace     bool
	Transient     bool
	IfNotExists   bool
	Name          *ObjectName
	Clone         *CloneSource     // CLONE source [AT|BEFORE (...)]; nil if absent
	ManagedAccess bool             // WITH MANAGED ACCESS
	Props         DBSchemaProps    // optional properties
	Tags          []*TagAssignment // WITH TAG (...); nil if absent
	Loc           Loc
}

func (n *CreateSchemaStmt) Tag() NodeTag { return T_CreateSchemaStmt }

var _ Node = (*CreateSchemaStmt)(nil)

// AlterSchemaStmt represents ALTER SCHEMA ... (all action variants).
type AlterSchemaStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterSchemaAction
	NewName    *ObjectName      // RENAME TO / SWAP WITH target
	SetProps   *DBSchemaProps   // SET properties; non-nil for AlterSchemaSet
	UnsetProps []string         // UNSET property names
	Tags       []*TagAssignment // SET TAG (...) assignments
	UnsetTags  []*ObjectName    // UNSET TAG (...) names
	Loc        Loc
}

func (n *AlterSchemaStmt) Tag() NodeTag { return T_AlterSchemaStmt }

var _ Node = (*AlterSchemaStmt)(nil)

// DropSchemaStmt represents DROP SCHEMA [IF EXISTS] name [CASCADE|RESTRICT].
type DropSchemaStmt struct {
	IfExists bool
	Name     *ObjectName
	Cascade  bool
	Restrict bool
	Loc      Loc
}

func (n *DropSchemaStmt) Tag() NodeTag { return T_DropSchemaStmt }

var _ Node = (*DropSchemaStmt)(nil)

// UndropSchemaStmt represents UNDROP SCHEMA name.
type UndropSchemaStmt struct {
	Name *ObjectName
	Loc  Loc
}

func (n *UndropSchemaStmt) Tag() NodeTag { return T_UndropSchemaStmt }

var _ Node = (*UndropSchemaStmt)(nil)

// ---------------------------------------------------------------------------
// DROP / UNDROP statement nodes (non-DATABASE/SCHEMA object types)
// ---------------------------------------------------------------------------

// DropObjectKind enumerates the object types that can appear in a DROP
// statement handled by this parser. DATABASE and SCHEMA are handled by T2.1.
type DropObjectKind int

const (
	DropTable            DropObjectKind = iota
	DropView                            // DROP VIEW
	DropMaterializedView                // DROP MATERIALIZED VIEW
	DropDynamicTable                    // DROP DYNAMIC TABLE
	DropExternalTable                   // DROP EXTERNAL TABLE
	DropStream                          // DROP STREAM
	DropTask                            // DROP TASK
	DropSequence                        // DROP SEQUENCE
	DropStage                           // DROP STAGE
	DropFileFormat                      // DROP FILE FORMAT
	DropFunction                        // DROP FUNCTION
	DropProcedure                       // DROP PROCEDURE
	DropPipe                            // DROP PIPE
	DropTag                             // DROP TAG
	DropRole                            // DROP ROLE
	DropWarehouse                       // DROP WAREHOUSE
)

// String returns the SQL object-type keyword for the kind.
func (k DropObjectKind) String() string {
	switch k {
	case DropTable:
		return "TABLE"
	case DropView:
		return "VIEW"
	case DropMaterializedView:
		return "MATERIALIZED VIEW"
	case DropDynamicTable:
		return "DYNAMIC TABLE"
	case DropExternalTable:
		return "EXTERNAL TABLE"
	case DropStream:
		return "STREAM"
	case DropTask:
		return "TASK"
	case DropSequence:
		return "SEQUENCE"
	case DropStage:
		return "STAGE"
	case DropFileFormat:
		return "FILE FORMAT"
	case DropFunction:
		return "FUNCTION"
	case DropProcedure:
		return "PROCEDURE"
	case DropPipe:
		return "PIPE"
	case DropTag:
		return "TAG"
	case DropRole:
		return "ROLE"
	case DropWarehouse:
		return "WAREHOUSE"
	default:
		return "UNKNOWN"
	}
}

// DropStmt represents a DROP <object_type> [IF EXISTS] name [CASCADE|RESTRICT]
// statement. DATABASE and SCHEMA are handled by T2.1's DropDatabaseStmt /
// DropSchemaStmt and are NOT covered by this type.
type DropStmt struct {
	Kind     DropObjectKind
	IfExists bool
	Name     *ObjectName
	Cascade  bool // CASCADE option (mutually exclusive with Restrict)
	Restrict bool // RESTRICT option (mutually exclusive with Cascade)
	Loc      Loc
}

// Tag implements Node.
func (n *DropStmt) Tag() NodeTag { return T_DropStmt }

// UndropObjectKind enumerates the object types that can appear in an UNDROP
// statement. DATABASE and SCHEMA are handled by T2.1.
type UndropObjectKind int

const (
	UndropTable        UndropObjectKind = iota
	UndropDynamicTable                  // UNDROP DYNAMIC TABLE
	UndropTag                           // UNDROP TAG
)

// String returns the SQL object-type keyword for the kind.
func (k UndropObjectKind) String() string {
	switch k {
	case UndropTable:
		return "TABLE"
	case UndropDynamicTable:
		return "DYNAMIC TABLE"
	case UndropTag:
		return "TAG"
	default:
		return "UNKNOWN"
	}
}

// UndropStmt represents an UNDROP <object_type> name statement.
// DATABASE and SCHEMA are handled by T2.1.
type UndropStmt struct {
	Kind UndropObjectKind
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *UndropStmt) Tag() NodeTag { return T_UndropStmt }

// Compile-time assertions.
var (
	_ Node = (*DropStmt)(nil)
	_ Node = (*UndropStmt)(nil)
)

// ---------------------------------------------------------------------------
// DML statement nodes
// ---------------------------------------------------------------------------

// InsertStmt represents a single-table INSERT statement:
//
//	INSERT [OVERWRITE] INTO table [(cols)] {VALUES (exprs)[, ...] | SELECT ...}
type InsertStmt struct {
	Overwrite bool
	Target    *ObjectName
	Columns   []Ident  // optional column list; nil if not specified
	Values    [][]Node // VALUES rows; nil if SELECT form used
	Select    Node     // SELECT body; nil if VALUES form used
	Loc       Loc
}

// Tag implements Node.
func (n *InsertStmt) Tag() NodeTag { return T_InsertStmt }

// InsertMultiStmt represents INSERT ALL / INSERT FIRST (multi-table insert):
//
//	INSERT [OVERWRITE] {ALL | FIRST} [WHEN cond THEN] INTO target [(cols)] VALUES (exprs) ... SELECT ...
type InsertMultiStmt struct {
	Overwrite bool
	First     bool                 // true = INSERT FIRST, false = INSERT ALL
	Branches  []*InsertMultiBranch // INTO targets with optional WHEN guards
	Select    Node                 // driving SELECT query
	Loc       Loc
}

// Tag implements Node.
func (n *InsertMultiStmt) Tag() NodeTag { return T_InsertMultiStmt }

// InsertMultiBranch is one INTO target inside an INSERT ALL/FIRST statement.
// When is nil for unconditional branches; non-nil for WHEN cond THEN branches.
type InsertMultiBranch struct {
	When    Node // nil for unconditional ALL/FIRST; WHEN condition otherwise
	Target  *ObjectName
	Columns []Ident
	Values  []Node // single row of VALUES expressions; nil if VALUES clause omitted
	Loc     Loc
}

// UpdateSet represents one assignment in an UPDATE SET list.
// Column holds the target column name, which may be qualified (e.g. t.col).
type UpdateSet struct {
	Column *ObjectName
	Value  Node
	Loc    Loc
}

// UpdateStmt represents an UPDATE statement:
//
//	UPDATE table [alias] SET col = expr [, ...] [FROM sources] [WHERE cond]
type UpdateStmt struct {
	Target *ObjectName
	Alias  Ident // table alias; zero Ident if absent
	Sets   []*UpdateSet
	From   []Node // FROM clause items (Snowflake extension for joins); nil if absent
	Where  Node   // WHERE condition; nil if absent
	Loc    Loc
}

// Tag implements Node.
func (n *UpdateStmt) Tag() NodeTag { return T_UpdateStmt }

// DeleteStmt represents a DELETE statement:
//
//	DELETE FROM table [alias] [USING sources] [WHERE cond]
type DeleteStmt struct {
	Target *ObjectName
	Alias  Ident  // table alias; zero Ident if absent
	Using  []Node // USING clause items; nil if absent
	Where  Node   // WHERE condition; nil if absent
	Loc    Loc
}

// Tag implements Node.
func (n *DeleteStmt) Tag() NodeTag { return T_DeleteStmt }

// MergeAction identifies the action in a WHEN clause of a MERGE statement.
type MergeAction int

const (
	MergeActionUpdate MergeAction = iota // WHEN ... THEN UPDATE SET ...
	MergeActionDelete                    // WHEN ... THEN DELETE
	MergeActionInsert                    // WHEN ... THEN INSERT ...
)

// MergeWhen represents one WHEN clause in a MERGE statement.
type MergeWhen struct {
	Matched       bool // true = WHEN MATCHED, false = WHEN NOT MATCHED
	BySource      bool // WHEN NOT MATCHED BY SOURCE (Snowflake extension)
	ByTarget      bool // WHEN NOT MATCHED BY TARGET (or plain NOT MATCHED)
	AndCond       Node // optional AND condition; nil if absent
	Action        MergeAction
	Sets          []*UpdateSet // for MergeActionUpdate
	InsertCols    []Ident      // for MergeActionInsert: optional column list
	InsertVals    []Node       // for MergeActionInsert: VALUES expressions
	InsertDefault bool         // for MergeActionInsert: INSERT VALUES DEFAULT
	Loc           Loc
}

// MergeStmt represents a MERGE statement:
//
//	MERGE INTO target [alias] USING source [alias] ON cond [WHEN ...]...
type MergeStmt struct {
	Target      *ObjectName
	TargetAlias Ident // alias for target; zero Ident if absent
	Source      Node  // table ref or subquery
	SourceAlias Ident // alias for source; zero Ident if absent
	On          Node
	Whens       []*MergeWhen
	Loc         Loc
}

// Tag implements Node.
func (n *MergeStmt) Tag() NodeTag { return T_MergeStmt }

// Compile-time assertions for DML nodes.
var (
	_ Node = (*InsertStmt)(nil)
	_ Node = (*InsertMultiStmt)(nil)
	_ Node = (*UpdateStmt)(nil)
	_ Node = (*DeleteStmt)(nil)
	_ Node = (*MergeStmt)(nil)
)

// ---------------------------------------------------------------------------
// VIEW / MATERIALIZED VIEW DDL — helper structs (not Nodes)
// ---------------------------------------------------------------------------

// ViewColumn represents a single column entry in a VIEW column list or a
// view_col binding (WITH MASKING POLICY / WITH TAG).
//
// Used by both CreateViewStmt and CreateMaterializedViewStmt.
type ViewColumn struct {
	Name          Ident
	MaskingPolicy *ObjectName      // WITH MASKING POLICY p; nil if absent
	MaskingUsing  []Ident          // USING (col, ...); nil if absent
	Tags          []*TagAssignment // WITH TAG (...); nil if absent
	Comment       *string          // COMMENT 'text'; nil if absent
	Loc           Loc
}

// RowAccessPolicy holds the policy name and column list from a
// WITH ROW ACCESS POLICY name ON (col, ...) clause.
type RowAccessPolicy struct {
	PolicyName *ObjectName
	Columns    []Ident
}

// ---------------------------------------------------------------------------
// CREATE VIEW
// ---------------------------------------------------------------------------

// CreateViewStmt represents CREATE [OR REPLACE] [SECURE] [RECURSIVE] VIEW ...
type CreateViewStmt struct {
	OrReplace   bool
	Secure      bool
	Recursive   bool
	IfNotExists bool
	Name        *ObjectName
	Columns     []*ViewColumn // optional column list from ( col [COMMENT 'x'], ... )
	ViewCols    []*ViewColumn // view_col* — column masking/tag bindings (outside parens)
	CopyGrants  bool
	Comment     *string          // COMMENT = 'str'; nil if absent
	Tags        []*TagAssignment // WITH TAG (...); nil if absent
	RowPolicy   *RowAccessPolicy // WITH ROW ACCESS POLICY ...; nil if absent
	Query       Node             // the AS query body
	Loc         Loc
}

// Tag implements Node.
func (n *CreateViewStmt) Tag() NodeTag { return T_CreateViewStmt }

var _ Node = (*CreateViewStmt)(nil)

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// CreateMaterializedViewStmt represents CREATE [OR REPLACE] [SECURE] MATERIALIZED VIEW ...
type CreateMaterializedViewStmt struct {
	OrReplace   bool
	Secure      bool
	IfNotExists bool
	Name        *ObjectName
	Columns     []*ViewColumn // optional column list
	ViewCols    []*ViewColumn // view_col* — column masking/tag bindings
	CopyGrants  bool
	Comment     *string          // COMMENT = 'str'; nil if absent
	ClusterBy   []Node           // CLUSTER BY (exprs); nil if absent
	Linear      bool             // CLUSTER BY LINEAR modifier
	Tags        []*TagAssignment // WITH TAG (...); nil if absent
	RowPolicy   *RowAccessPolicy // WITH ROW ACCESS POLICY ...; nil if absent
	Query       Node             // the AS query body
	Loc         Loc
}

// Tag implements Node.
func (n *CreateMaterializedViewStmt) Tag() NodeTag { return T_CreateMaterializedViewStmt }

var _ Node = (*CreateMaterializedViewStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER VIEW — action enum and statement
// ---------------------------------------------------------------------------

// AlterViewAction discriminates the action variants of ALTER VIEW.
type AlterViewAction int

const (
	AlterViewRename                   AlterViewAction = iota // RENAME TO
	AlterViewSetComment                                      // SET COMMENT = '...'
	AlterViewUnsetComment                                    // UNSET COMMENT
	AlterViewSetSecure                                       // SET SECURE
	AlterViewUnsetSecure                                     // UNSET SECURE
	AlterViewSetTag                                          // SET TAG (...)
	AlterViewUnsetTag                                        // UNSET TAG (...)
	AlterViewAddRowAccessPolicy                              // ADD ROW ACCESS POLICY
	AlterViewDropRowAccessPolicy                             // DROP ROW ACCESS POLICY
	AlterViewDropAllRowAccessPolicies                        // DROP ALL ROW ACCESS POLICIES
	AlterViewColumnSetMaskingPolicy                          // ALTER COLUMN col SET MASKING POLICY
	AlterViewColumnUnsetMaskingPolicy                        // ALTER COLUMN col UNSET MASKING POLICY
	AlterViewColumnSetTag                                    // ALTER COLUMN col SET TAG (...)
	AlterViewColumnUnsetTag                                  // ALTER COLUMN col UNSET TAG (...)
)

// AlterViewStmt represents ALTER VIEW ... (all action variants).
type AlterViewStmt struct {
	IfExists      bool
	Name          *ObjectName
	Action        AlterViewAction
	NewName       *ObjectName      // RENAME TO
	Comment       *string          // SET COMMENT = '...'
	Secure        bool             // true = SET SECURE; false = UNSET SECURE (check Action)
	Tags          []*TagAssignment // SET TAG (...)
	UnsetTags     []*ObjectName    // UNSET TAG (...)
	PolicyName    *ObjectName      // ADD/DROP ROW ACCESS POLICY name
	PolicyCols    []Ident          // ON (col, ...) for ADD ROW ACCESS POLICY
	Column        Ident            // ALTER COLUMN col name
	MaskingPolicy *ObjectName      // SET MASKING POLICY p
	MaskingUsing  []Ident          // USING (col, ...)
	Loc           Loc
}

// Tag implements Node.
func (n *AlterViewStmt) Tag() NodeTag { return T_AlterViewStmt }

var _ Node = (*AlterViewStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW — action enum and statement
// ---------------------------------------------------------------------------

// AlterMaterializedViewAction discriminates the action variants of ALTER MATERIALIZED VIEW.
type AlterMaterializedViewAction int

const (
	AlterMVRename            AlterMaterializedViewAction = iota // RENAME TO
	AlterMVClusterBy                                            // CLUSTER BY (exprs)
	AlterMVDropClusteringKey                                    // DROP CLUSTERING KEY
	AlterMVSuspend                                              // SUSPEND
	AlterMVResume                                               // RESUME
	AlterMVSuspendRecluster                                     // SUSPEND RECLUSTER
	AlterMVResumeRecluster                                      // RESUME RECLUSTER
	AlterMVSetSecure                                            // SET SECURE
	AlterMVUnsetSecure                                          // UNSET SECURE
	AlterMVSetComment                                           // SET COMMENT = '...'
	AlterMVUnsetComment                                         // UNSET COMMENT
)

// AlterMaterializedViewStmt represents ALTER MATERIALIZED VIEW ... (all action variants).
// Note: the legacy grammar does NOT support IF EXISTS for ALTER MATERIALIZED VIEW.
type AlterMaterializedViewStmt struct {
	Name      *ObjectName
	Action    AlterMaterializedViewAction
	NewName   *ObjectName // RENAME TO
	ClusterBy []Node      // CLUSTER BY (exprs)
	Linear    bool        // CLUSTER BY LINEAR modifier
	Comment   *string     // SET COMMENT = '...'
	Secure    bool        // true = SET SECURE; false = UNSET SECURE (check Action)
	Loc       Loc
}

// Tag implements Node.
func (n *AlterMaterializedViewStmt) Tag() NodeTag { return T_AlterMaterializedViewStmt }

var _ Node = (*AlterMaterializedViewStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER TABLE DDL — enums, helpers, and statement node
// ---------------------------------------------------------------------------

// AlterTableActionKind discriminates the action variants of ALTER TABLE.
type AlterTableActionKind int

const (
	AlterTableRename                   AlterTableActionKind = iota // RENAME TO new_name
	AlterTableSwapWith                                             // SWAP WITH other_table
	AlterTableAddColumn                                            // ADD [COLUMN] [IF NOT EXISTS] col_def [, ...]
	AlterTableDropColumn                                           // DROP [COLUMN] [IF EXISTS] col [, ...]
	AlterTableRenameColumn                                         // RENAME COLUMN old TO new
	AlterTableAlterColumn                                          // ALTER/MODIFY COLUMN col ...
	AlterTableAddConstraint                                        // ADD [CONSTRAINT name] PK/UK/FK
	AlterTableDropConstraint                                       // DROP CONSTRAINT name | DROP PRIMARY KEY | DROP UNIQUE
	AlterTableRenameConstraint                                     // RENAME CONSTRAINT old TO new
	AlterTableClusterBy                                            // CLUSTER BY [LINEAR] (exprs)
	AlterTableDropClusterKey                                       // DROP CLUSTERING KEY
	AlterTableRecluster                                            // RECLUSTER [MAX_SIZE = n] [WHERE expr]
	AlterTableSuspendRecluster                                     // SUSPEND RECLUSTER
	AlterTableResumeRecluster                                      // RESUME RECLUSTER
	AlterTableSet                                                  // SET properties
	AlterTableUnset                                                // UNSET properties
	AlterTableSetTag                                               // SET TAG (...)
	AlterTableUnsetTag                                             // UNSET TAG (...)
	AlterTableAddRowAccessPolicy                                   // ADD ROW ACCESS POLICY name ON (cols)
	AlterTableDropRowAccessPolicy                                  // DROP ROW ACCESS POLICY name
	AlterTableDropAllRowAccessPolicies                             // DROP ALL ROW ACCESS POLICIES
	AlterTableAddSearchOpt                                         // ADD SEARCH OPTIMIZATION [ON ...]
	AlterTableDropSearchOpt                                        // DROP SEARCH OPTIMIZATION [ON ...]
	AlterTableSetMaskingPolicy                                     // ALTER/MODIFY COLUMN col SET MASKING POLICY
	AlterTableUnsetMaskingPolicy                                   // ALTER/MODIFY COLUMN col UNSET MASKING POLICY
	AlterTableSetColumnTag                                         // ALTER/MODIFY col SET TAG (...)
	AlterTableUnsetColumnTag                                       // ALTER/MODIFY col UNSET TAG (...)
)

// ColumnAlterKind discriminates the sub-action inside ALTER/MODIFY COLUMN.
type ColumnAlterKind int

const (
	ColumnAlterSetDataType  ColumnAlterKind = iota // SET DATA TYPE t | TYPE t | t
	ColumnAlterSetDefault                          // SET DEFAULT expr
	ColumnAlterDropDefault                         // DROP DEFAULT
	ColumnAlterSetNotNull                          // SET NOT NULL
	ColumnAlterDropNotNull                         // DROP NOT NULL
	ColumnAlterSetComment                          // COMMENT 'text'
	ColumnAlterUnsetComment                        // UNSET COMMENT
)

// ColumnAlter holds the specification for a single ALTER/MODIFY COLUMN
// sub-action (used inside AlterTableAction.ColumnAlters).
//
// ColumnAlter is NOT a Node. It is owned by AlterTableAction.
type ColumnAlter struct {
	Column      Ident // column being altered
	Kind        ColumnAlterKind
	DataType    *TypeName // for ColumnAlterSetDataType
	DefaultExpr Node      // for ColumnAlterSetDefault
	Comment     *string   // for ColumnAlterSetComment
}

// TableProp is a single SET property key=value pair used in ALTER TABLE SET.
// It is NOT a Node.
type TableProp struct {
	Name  string // uppercased property name
	Value string // raw value text
}

// AlterTableAction represents one action in an ALTER TABLE statement.
// Only the fields relevant to the Kind are populated.
//
// AlterTableAction is NOT a Node. It is owned by AlterTableStmt.
type AlterTableAction struct {
	Kind AlterTableActionKind

	// --- Rename / SwapWith ---
	NewName *ObjectName // RENAME TO target / SWAP WITH target

	// --- AddColumn ---
	Columns     []*ColumnDef // column definitions for ADD COLUMN
	IfNotExists bool         // ADD COLUMN IF NOT EXISTS guard

	// --- DropColumn ---
	DropColumnNames []Ident // column names for DROP COLUMN
	IfExists        bool    // DROP COLUMN IF EXISTS guard

	// --- RenameColumn ---
	OldName    Ident // RENAME COLUMN old
	NewColName Ident // RENAME COLUMN ... TO new

	// --- AlterColumn ---
	ColumnAlters []*ColumnAlter // one per ALTER/MODIFY COLUMN spec

	// --- AddConstraint ---
	Constraint *TableConstraint // reuse T2.2 TableConstraint

	// --- DropConstraint / RenameConstraint ---
	ConstraintName    Ident // DROP CONSTRAINT name / RENAME CONSTRAINT old
	NewConstraintName Ident // RENAME CONSTRAINT ... TO new
	IsPrimaryKey      bool  // DROP PRIMARY KEY (unnamed)
	DropUnique        bool  // DROP UNIQUE
	DropForeignKey    bool  // DROP FOREIGN KEY
	Cascade           bool
	Restrict          bool

	// --- ClusterBy ---
	ClusterBy []Node // CLUSTER BY expressions
	Linear    bool   // CLUSTER BY LINEAR

	// --- Recluster ---
	ReclusterMaxSize *int64 // MAX_SIZE = n; nil if absent
	ReclusterWhere   Node   // WHERE expr; nil if absent

	// --- Set / Unset properties ---
	Props      []*TableProp // SET property list
	UnsetProps []string     // UNSET property names

	// --- SetTag / UnsetTag ---
	Tags      []*TagAssignment // SET TAG assignments
	UnsetTags []*ObjectName    // UNSET TAG names

	// --- Row access policy ---
	PolicyName *ObjectName // ADD/DROP ROW ACCESS POLICY name
	PolicyCols []Ident     // ADD ROW ACCESS POLICY ... ON (cols)

	// --- Search optimization ---
	SearchOptOn []string // ON targets (raw text); nil = no ON clause

	// --- Masking policy (column-level) ---
	MaskColumn    Ident       // column for SET/UNSET MASKING POLICY
	MaskingPolicy *ObjectName // SET MASKING POLICY target

	// --- Column tag (ALTER/MODIFY col SET/UNSET TAG) ---
	TagColumn Ident // column for SetColumnTag / UnsetColumnTag

	Loc Loc
}

// AlterTableStmt represents ALTER TABLE [IF EXISTS] name action [, action ...].
type AlterTableStmt struct {
	IfExists bool
	Name     *ObjectName
	Actions  []*AlterTableAction
	Loc      Loc
}

// Tag implements Node.
func (n *AlterTableStmt) Tag() NodeTag { return T_AlterTableStmt }

// Compile-time assertion.
var _ Node = (*AlterTableStmt)(nil)

// ---------------------------------------------------------------------------
// DCL — GRANT / REVOKE (roles, privileges, ownership, shares) [T6.1]
// ---------------------------------------------------------------------------
//
// Snowflake's GRANT/REVOKE surface is large and grows continuously: new
// privileges (e.g. CREATE PROVISIONED THROUGHPUT, CREATE SNOWFLAKE.CORE.BUDGET)
// and new object types (e.g. NOTEBOOK, WORKSPACE, COMPUTE POOL, SEMANTIC VIEW)
// are added over time. Rather than hard-coding the legacy ANTLR grammar's
// finite privilege/object enumerations — which the official documentation
// corpus already exceeds — these nodes model privileges and object types as
// open-ended token runs (free-form uppercased names). This keeps the parser
// faithful to the docs (truth1) and resilient to Snowflake's additions; the
// catalog/semantic layer, not the parser, is responsible for validating that
// a given privilege is legal for a given object.

// GrantKind discriminates the three structural shapes of a GRANT statement.
type GrantKind int

const (
	// GrantPrivileges is GRANT <privileges> ON <target> TO <grantee>.
	GrantPrivileges GrantKind = iota
	// GrantRole is GRANT [DATABASE | APPLICATION] ROLE <name> TO <grantee>.
	GrantRole
	// GrantOwnership is GRANT OWNERSHIP ON <target> TO <grantee> [...CURRENT GRANTS].
	GrantOwnership
)

// GrantedRoleKind discriminates the kind of role being granted/revoked in a
// role-grant statement (the noun after GRANT/REVOKE).
type GrantedRoleKind int

const (
	// GrantedAccountRole is GRANT ROLE <name> ... (account-level role).
	GrantedAccountRole GrantedRoleKind = iota
	// GrantedDatabaseRole is GRANT DATABASE ROLE <name> ...
	GrantedDatabaseRole
	// GrantedApplicationRole is GRANT APPLICATION ROLE <name> ...
	GrantedApplicationRole
)

// String returns a human-readable name for the kind.
func (k GrantedRoleKind) String() string {
	switch k {
	case GrantedAccountRole:
		return "Role"
	case GrantedDatabaseRole:
		return "DatabaseRole"
	case GrantedApplicationRole:
		return "ApplicationRole"
	default:
		return "Unknown"
	}
}

// String returns a human-readable name for the kind.
func (k GrantKind) String() string {
	switch k {
	case GrantPrivileges:
		return "Privileges"
	case GrantRole:
		return "Role"
	case GrantOwnership:
		return "Ownership"
	default:
		return "Unknown"
	}
}

// RevokeKind discriminates the two structural shapes of a REVOKE statement.
// (There is no REVOKE OWNERSHIP — ownership is transferred, never revoked.)
type RevokeKind int

const (
	// RevokePrivileges is REVOKE [GRANT OPTION FOR] <privileges> ON <target> FROM <grantee>.
	RevokePrivileges RevokeKind = iota
	// RevokeRole is REVOKE [DATABASE] ROLE <name> FROM <grantee>.
	RevokeRole
)

// String returns a human-readable name for the kind.
func (k RevokeKind) String() string {
	switch k {
	case RevokePrivileges:
		return "Privileges"
	case RevokeRole:
		return "Role"
	default:
		return "Unknown"
	}
}

// Privilege is a single privilege in a GRANT/REVOKE privilege list. Name holds
// the uppercased source text of the privilege, which may be multiple words
// (e.g. "SELECT", "CREATE MATERIALIZED VIEW", "CREATE PROVISIONED THROUGHPUT",
// "READ", "CREATE SNOWFLAKE.CORE.BUDGET").
//
// Privilege is NOT a Node; it is owned by GrantStmt / RevokeStmt.
type Privilege struct {
	Name string
	Loc  Loc
}

// GranteeKind discriminates the recipient type following TO/FROM in a GRANT or
// REVOKE statement.
type GranteeKind int

const (
	GranteeRole            GranteeKind = iota // [TO|FROM] [ROLE] <name>
	GranteeDatabaseRole                       // [TO|FROM] DATABASE ROLE <name>
	GranteeUser                               // [TO|FROM] USER <name>
	GranteeShare                              // [TO|FROM] SHARE <name>
	GranteeApplication                        // TO APPLICATION <name>
	GranteeApplicationRole                    // TO APPLICATION ROLE <name>
)

// String returns the SQL keyword(s) for the grantee kind.
func (k GranteeKind) String() string {
	switch k {
	case GranteeRole:
		return "ROLE"
	case GranteeDatabaseRole:
		return "DATABASE ROLE"
	case GranteeUser:
		return "USER"
	case GranteeShare:
		return "SHARE"
	case GranteeApplication:
		return "APPLICATION"
	case GranteeApplicationRole:
		return "APPLICATION ROLE"
	default:
		return "UNKNOWN"
	}
}

// Grantee is the recipient of a GRANT (after TO) or the subject of a REVOKE
// (after FROM). Name is the (possibly qualified) role/user/share name.
//
// Grantee is a Node so the walker visits its Name. (Name is itself an
// *ObjectName Node.)
type Grantee struct {
	Kind GranteeKind
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *Grantee) Tag() NodeTag { return T_Grantee }

// GrantTargetKind discriminates the shape of the ON clause in a privilege or
// ownership GRANT/REVOKE.
type GrantTargetKind int

const (
	// GrantTargetAccount is ON ACCOUNT.
	GrantTargetAccount GrantTargetKind = iota
	// GrantTargetObject is ON <object_type> <name> [ ( signature ) ].
	GrantTargetObject
	// GrantTargetAllIn is ON ALL <object_type_plural> IN { DATABASE <db> | SCHEMA <schema> }.
	GrantTargetAllIn
	// GrantTargetFutureIn is ON FUTURE <object_type_plural> IN { DATABASE <db> | SCHEMA <schema> }.
	GrantTargetFutureIn
)

// String returns a human-readable name for the kind.
func (k GrantTargetKind) String() string {
	switch k {
	case GrantTargetAccount:
		return "Account"
	case GrantTargetObject:
		return "Object"
	case GrantTargetAllIn:
		return "AllIn"
	case GrantTargetFutureIn:
		return "FutureIn"
	default:
		return "Unknown"
	}
}

// GrantContainerKind discriminates the IN container of ALL/FUTURE ... IN forms.
type GrantContainerKind int

const (
	GrantContainerNone     GrantContainerKind = iota // not an ALL/FUTURE form
	GrantContainerDatabase                           // IN DATABASE <db>
	GrantContainerSchema                             // IN SCHEMA <schema>
)

// GrantTarget is the object the privileges/ownership apply to (the ON clause).
//
// Field population by Kind:
//   - GrantTargetAccount:  no other fields set.
//   - GrantTargetObject:   ObjectType + Name (+ optional Signature).
//   - GrantTargetAllIn / GrantTargetFutureIn: ObjectTypePlural + Container + ContainerName.
//
// ObjectType / ObjectTypePlural are uppercased free-form names (possibly
// multi-word, e.g. "MATERIALIZED VIEW", "EXTERNAL TABLE", "CORTEX SEARCH
// SERVICE", "STORAGE LIFECYCLE POLICY").
//
// GrantTarget is a Node so the walker visits Name and ContainerName. The
// Signature TypeName slice is not walker-visited (mirroring how
// CreateTableStmt.Columns is handled).
type GrantTarget struct {
	Kind GrantTargetKind

	// GrantTargetObject:
	ObjectType string      // e.g. "TABLE", "FUNCTION", "MATERIALIZED VIEW", "NOTEBOOK", "CORTEX SEARCH SERVICE"
	Name       *ObjectName // the object's qualified name
	Signature  []*TypeName // FUNCTION/PROCEDURE arg-type list; nil if no ( ... ) present, non-nil (maybe empty) if ()

	// GrantTargetAllIn / GrantTargetFutureIn:
	ObjectTypePlural string             // e.g. "TABLES", "VIEWS", "SCHEMAS"
	Container        GrantContainerKind // DATABASE or SCHEMA
	ContainerName    *ObjectName        // the database/schema name

	Loc Loc
}

// Tag implements Node.
func (n *GrantTarget) Tag() NodeTag { return T_GrantTarget }

// CurrentGrantsAction is the trailing { REVOKE | COPY } CURRENT GRANTS option
// on GRANT OWNERSHIP.
type CurrentGrantsAction int

const (
	CurrentGrantsNone   CurrentGrantsAction = iota // no trailing clause
	CurrentGrantsRevoke                            // REVOKE CURRENT GRANTS
	CurrentGrantsCopy                              // COPY CURRENT GRANTS
)

// String returns a human-readable name for the action.
func (a CurrentGrantsAction) String() string {
	switch a {
	case CurrentGrantsNone:
		return "None"
	case CurrentGrantsRevoke:
		return "RevokeCurrentGrants"
	case CurrentGrantsCopy:
		return "CopyCurrentGrants"
	default:
		return "Unknown"
	}
}

// GrantStmt represents a GRANT statement in any of its three shapes
// (privileges, role, ownership). See GrantKind.
//
//	GRANT { <privileges> | ALL [PRIVILEGES] } ON <target> TO <grantee> [WITH GRANT OPTION]
//	GRANT [DATABASE | APPLICATION] ROLE <name> TO <grantee>
//	GRANT OWNERSHIP ON <target> TO <grantee> [{REVOKE|COPY} CURRENT GRANTS]
type GrantStmt struct {
	Kind GrantKind

	// --- GrantPrivileges ---
	Privileges    []*Privilege // explicit privilege list; nil when AllPrivileges is true
	AllPrivileges bool         // ALL [PRIVILEGES]
	GrantOption   bool         // WITH GRANT OPTION

	// --- GrantRole ---
	Role     *ObjectName     // the role being granted (GrantRole)
	RoleKind GrantedRoleKind // account / database / application role

	// --- GrantOwnership ---
	CurrentGrants CurrentGrantsAction // trailing CURRENT GRANTS option

	// --- GrantPrivileges + GrantOwnership ---
	On *GrantTarget // the ON clause target; nil for GrantRole

	// --- all kinds ---
	Grantee *Grantee // the TO recipient

	Loc Loc
}

// Tag implements Node.
func (n *GrantStmt) Tag() NodeTag { return T_GrantStmt }

// RevokeStmt represents a REVOKE statement in either of its shapes
// (privileges, role). See RevokeKind.
//
//	REVOKE [GRANT OPTION FOR] { <privileges> | ALL [PRIVILEGES] } ON <target> FROM <grantee> [CASCADE|RESTRICT]
//	REVOKE [DATABASE | APPLICATION] ROLE <name> FROM <grantee>
type RevokeStmt struct {
	Kind RevokeKind

	// --- RevokePrivileges ---
	GrantOptionFor bool         // REVOKE GRANT OPTION FOR ...
	Privileges     []*Privilege // explicit privilege list; nil when AllPrivileges is true
	AllPrivileges  bool         // ALL [PRIVILEGES]
	On             *GrantTarget // the ON clause target; nil for RevokeRole
	Cascade        bool         // trailing CASCADE
	Restrict       bool         // trailing RESTRICT

	// --- RevokeRole ---
	Role     *ObjectName     // the role being revoked (RevokeRole)
	RoleKind GrantedRoleKind // account / database / application role

	// --- all kinds ---
	Grantee *Grantee // the FROM subject

	Loc Loc
}

// Tag implements Node.
func (n *RevokeStmt) Tag() NodeTag { return T_RevokeStmt }

// Compile-time assertions for DCL nodes.
var (
	_ Node = (*GrantStmt)(nil)
	_ Node = (*RevokeStmt)(nil)
	_ Node = (*Grantee)(nil)
	_ Node = (*GrantTarget)(nil)
)

// ---------------------------------------------------------------------------
// Utility / introspection — SHOW / DESCRIBE / USE / SET / UNSET / COMMENT /
// TRUNCATE (T6.3)
// ---------------------------------------------------------------------------
//
// Like GRANT/REVOKE, the SHOW and DESCRIBE/DESC surfaces are large and grow
// continuously: Snowflake documents 50+ SHOW object classes (DATABASES, TABLES,
// STREAMS, SEMANTIC VIEWS, GIT REPOSITORIES, ...) and 30+ DESCRIBE object types,
// and ships new ones (e.g. SEMANTIC VIEWS, DATASETS, CORTEX SEARCH SERVICE) over
// time. Rather than mirror the legacy ANTLR grammar's per-object rule
// enumerations — which the official documentation corpus already exceeds —
// these nodes model the object class/type as an open-ended uppercased token run
// (free-form ObjectClass / ObjectType string). Structural keywords (TERSE,
// HISTORY, LIKE, IN, STARTS WITH, LIMIT, WITH, IS, TYPE) anchor the grammar;
// validating that a given class/type is a real Snowflake object is the
// catalog/semantic layer's job, not the parser's. New object classes therefore
// parse without code changes.

// ShowFilterKind discriminates the IN scope qualifier of a SHOW statement
// (the noun after IN, e.g. ACCOUNT, DATABASE, SCHEMA).
type ShowFilterKind int

const (
	// ShowScopeNone means no IN clause was present.
	ShowScopeNone ShowFilterKind = iota
	// ShowScopeAccount is IN ACCOUNT.
	ShowScopeAccount
	// ShowScopeDatabase is IN DATABASE [<name>].
	ShowScopeDatabase
	// ShowScopeSchema is IN SCHEMA [<name>] or a bare IN <schema-name>.
	ShowScopeSchema
	// ShowScopeTable is IN TABLE [<name>] (SHOW COLUMNS / PARAMETERS).
	ShowScopeTable
	// ShowScopeView is IN VIEW [<name>] (SHOW COLUMNS).
	ShowScopeView
	// ShowScopeOther is any other IN/FOR scope keyword run captured verbatim
	// (e.g. IN APPLICATION, FOR SESSION, IN FAILOVER GROUP). The keyword run is
	// stored in ScopeText and the optional trailing name (if any) in ScopeName.
	ShowScopeOther
)

// String returns a human-readable name for the scope kind.
func (k ShowFilterKind) String() string {
	switch k {
	case ShowScopeNone:
		return "None"
	case ShowScopeAccount:
		return "Account"
	case ShowScopeDatabase:
		return "Database"
	case ShowScopeSchema:
		return "Schema"
	case ShowScopeTable:
		return "Table"
	case ShowScopeView:
		return "View"
	case ShowScopeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// ShowStmt represents a SHOW statement:
//
//	SHOW [TERSE] <object_class> [HISTORY] [LIKE '<pat>']
//	     [IN <scope>] [STARTS WITH '<s>'] [LIMIT <n> [FROM '<s>']]
//	     [WITH PRIVILEGES <priv> [, ...]]
//	SHOW GRANTS [ON <target> | TO <grantee> | OF <grantee>]
//	SHOW FUTURE GRANTS IN { DATABASE <db> | SCHEMA <schema> }
//
// ObjectClass is the open-ended, uppercased object-class token run (e.g.
// "DATABASES", "TABLES", "MATERIALIZED VIEWS", "SEMANTIC VIEWS",
// "GIT REPOSITORIES"). For SHOW GRANTS, IsGrants is true and the grant-specific
// fields (GrantsOn / GrantsTo / Future) are populated instead of the generic
// filter fields.
//
// A trailing `->> <query>` "result pipe" (SHOW ... ->> SELECT ...) is captured
// in Pipe as the piped statement Node.
//
// ShowStmt is a Node; the walker descends into LikeName? No — Like/StartsWith
// are literals stored as strings. The walker visits the structured child Nodes:
// ScopeName, GrantsOn, GrantsTo, and Pipe.
type ShowStmt struct {
	Terse       bool   // SHOW TERSE ...
	ObjectClass string // open-ended uppercased class run, e.g. "TABLES", "MATERIALIZED VIEWS"
	History     bool   // ... HISTORY

	// Generic filter options (mutually exclusive with the GRANTS fields).
	Like       string // LIKE '<pat>' raw source text incl. quotes; "" if absent
	HasLike    bool   // whether a LIKE clause was present
	Scope      ShowFilterKind
	ScopeText  string      // verbatim scope keyword run for ShowScopeOther (e.g. "FAILOVER GROUP")
	ScopeName  *ObjectName // optional name after the scope keyword (db/schema/table name)
	StartsWith string      // STARTS WITH '<s>' raw text; "" if absent
	HasStarts  bool
	Limit      string // LIMIT <n> raw integer text; "" if absent
	HasLimit   bool
	LimitFrom  string       // LIMIT n FROM '<s>' raw text; "" if absent
	Privileges []*Privilege // WITH PRIVILEGES <list>; nil if absent

	// SHOW GRANTS specifics.
	IsGrants bool         // SHOW GRANTS ...
	Future   bool         // SHOW FUTURE GRANTS ...
	GrantsOn *GrantTarget // ON <target> (SHOW GRANTS ON ... / SHOW FUTURE GRANTS IN ...)
	GrantsTo *Grantee     // TO/OF <grantee> (SHOW GRANTS TO/OF ...)

	// Result pipe: SHOW ... ->> <query>.
	Pipe Node // the piped statement (e.g. a SELECT), or nil if no ->> present

	Loc Loc
}

// Tag implements Node.
func (n *ShowStmt) Tag() NodeTag { return T_ShowStmt }

// DescribeStmt represents a DESCRIBE | DESC statement:
//
//	{ DESCRIBE | DESC } <object_type> <name> [ ( <signature> ) ] [ TYPE = <kw> ]
//	{ DESCRIBE | DESC } SEARCH OPTIMIZATION ON <name>
//	{ DESCRIBE | DESC } RESULT { '<query_id>' | LAST_QUERY_ID() }
//
// ObjectType is the open-ended, uppercased object-type token run (e.g. "TABLE",
// "MATERIALIZED VIEW", "SEARCH OPTIMIZATION", "ROW ACCESS POLICY"). The name
// may be a qualified identifier (Name) or a literal (NameLiteral, for
// DESCRIBE RESULT '<id>' / DESCRIBE TRANSACTION <num>). On is set for the
// SEARCH OPTIMIZATION ON <name> form. Signature holds an optional
// FUNCTION/PROCEDURE argument-type list. TypeOption holds the trailing
// TYPE = <kw> value (e.g. "COLUMNS", "STAGE") uppercased, "" if absent.
//
// DescribeStmt is a Node; the walker descends into Name (an *ObjectName).
type DescribeStmt struct {
	Short       bool        // true if spelled DESC, false if DESCRIBE
	ObjectType  string      // open-ended uppercased type run
	Name        *ObjectName // the object's qualified name; nil if NameLiteral set
	NameLiteral *Literal    // literal name (DESCRIBE RESULT '<id>' / TRANSACTION <num>); nil otherwise
	Signature   []*TypeName // FUNCTION/PROCEDURE arg-type list; nil if no ( ... )
	TypeOption  string      // trailing TYPE = <kw>, uppercased; "" if absent
	Loc         Loc
}

// Tag implements Node.
func (n *DescribeStmt) Tag() NodeTag { return T_DescribeStmt }

// UseTargetKind discriminates the object kind in a USE statement.
type UseTargetKind int

const (
	// UseDefault is USE <name> with no object keyword (defaults to database in
	// Snowflake, but the parser does not assume — it records the bare form).
	UseDefault UseTargetKind = iota
	// UseRole is USE ROLE <name>.
	UseRole
	// UseWarehouse is USE WAREHOUSE <name>.
	UseWarehouse
	// UseDatabase is USE DATABASE <name>.
	UseDatabase
	// UseSchema is USE [SCHEMA] <name>.
	UseSchema
	// UseSecondaryRoles is USE SECONDARY ROLES { ALL | NONE | <roles> }.
	UseSecondaryRoles
)

// String returns a human-readable name for the kind.
func (k UseTargetKind) String() string {
	switch k {
	case UseDefault:
		return "Default"
	case UseRole:
		return "Role"
	case UseWarehouse:
		return "Warehouse"
	case UseDatabase:
		return "Database"
	case UseSchema:
		return "Schema"
	case UseSecondaryRoles:
		return "SecondaryRoles"
	default:
		return "Unknown"
	}
}

// UseStmt represents a USE statement:
//
//	USE [DATABASE] <name>
//	USE ROLE <name>
//	USE [SCHEMA] <name>
//	USE WAREHOUSE <name>
//	USE SECONDARY ROLES { ALL | NONE | <role> [, ...] }
//
// For USE SECONDARY ROLES, SecondaryRoles holds the uppercased option
// ("ALL" / "NONE") or a comma-joined role-name list; Name is nil.
//
// UseStmt is a Node; the walker descends into Name (an *ObjectName).
type UseStmt struct {
	Kind           UseTargetKind
	Name           *ObjectName // the target name; nil for SecondaryRoles
	SecondaryRoles string      // USE SECONDARY ROLES <this>; "" otherwise
	Loc            Loc
}

// Tag implements Node.
func (n *UseStmt) Tag() NodeTag { return T_UseStmt }

// SetVar pairs a session-variable name with its assigned value expression in a
// SET statement. SetVar is NOT a Node; it is owned by SetStmt.
type SetVar struct {
	Name  Ident
	Value Node // the assigned expression
	Loc   Loc
}

// SetStmt represents a SET session-variable assignment:
//
//	SET <var> = <expr>
//	SET ( <var> [, ...] ) = ( <expr> [, ...] )
//
// Single-variable form has exactly one Vars entry. The parenthesized
// multi-variable form pairs each variable with the value at the same position.
//
// SetStmt is a Node; the walker descends into each Vars[i].Value.
type SetStmt struct {
	Vars  []*SetVar
	Paren bool // true for the SET (...) = (...) parenthesized form
	Loc   Loc
}

// Tag implements Node.
func (n *SetStmt) Tag() NodeTag { return T_SetStmt }

// UnsetStmt represents UNSET of one or more session variables:
//
//	UNSET <var>
//	UNSET ( <var> [, ...] )
//
// UnsetStmt is a Node but has no child Nodes (variable names are value Idents).
type UnsetStmt struct {
	Names []Ident
	Paren bool // true for the UNSET (...) parenthesized form
	Loc   Loc
}

// Tag implements Node.
func (n *UnsetStmt) Tag() NodeTag { return T_UnsetStmt }

// CommentStmt represents a COMMENT statement attaching a comment string to an
// object or a column:
//
//	COMMENT [IF EXISTS] ON <object_type> <name> [ ( <signature> ) ] IS '<string>'
//	COMMENT [IF EXISTS] ON COLUMN <column_name> IS '<string>'
//
// ObjectType is the open-ended, uppercased object-type token run (e.g.
// "TABLE", "MASKING POLICY", "ROW ACCESS POLICY"). For the object form, Name
// holds the 1/2/3-part object name. For the COLUMN form, IsColumn is true,
// ObjectType is "COLUMN", and Column holds the 1- to 4-part column reference —
// Snowflake's full_column_name allows db.schema.table.column, which exceeds
// ObjectName's 3-part shape, so a ColumnRef is used. Comment holds the raw
// quoted string source.
//
// CommentStmt is a Node; the walker descends into Name (an *ObjectName) and
// Column (a *ColumnRef).
type CommentStmt struct {
	IfExists   bool
	IsColumn   bool        // COMMENT ON COLUMN <name> ...
	ObjectType string      // open-ended uppercased type run ("COLUMN" for the column form)
	Name       *ObjectName // object form: the 1/2/3-part object name; nil for the column form
	Column     *ColumnRef  // column form: the 1- to 4-part column reference; nil for the object form
	Signature  []*TypeName // FUNCTION/PROCEDURE arg-type list; nil if no ( ... )
	Comment    string      // IS '<string>' raw source text incl. quotes
	Loc        Loc
}

// Tag implements Node.
func (n *CommentStmt) Tag() NodeTag { return T_CommentStmt }

// TruncateStmt represents a TRUNCATE statement:
//
//	TRUNCATE [TABLE] [IF EXISTS] <name>
//	TRUNCATE [TABLE] [IF EXISTS] ERROR_TABLE( <base_table_name> )
//	TRUNCATE MATERIALIZED VIEW <name>
//
// MaterializedView is true for the TRUNCATE MATERIALIZED VIEW form (which has
// no IF EXISTS in the legacy grammar). ErrorTable is true for the
// ERROR_TABLE(<base>) form (an iceberg/error-table truncate documented by
// Snowflake but absent from the legacy grammar); Name then holds the base table
// name. Otherwise it is a plain table truncate; the TABLE keyword is optional.
//
// TruncateStmt is a Node; the walker descends into Name (an *ObjectName).
type TruncateStmt struct {
	MaterializedView bool        // TRUNCATE MATERIALIZED VIEW <name>
	ErrorTable       bool        // TRUNCATE [TABLE] [IF EXISTS] ERROR_TABLE(<name>)
	IfExists         bool        // TRUNCATE [TABLE] IF EXISTS <name>
	Name             *ObjectName // the target name (or ERROR_TABLE base table name)
	Loc              Loc
}

// Tag implements Node.
func (n *TruncateStmt) Tag() NodeTag { return T_TruncateStmt }

// Compile-time assertions for utility nodes.
var (
	_ Node = (*ShowStmt)(nil)
	_ Node = (*DescribeStmt)(nil)
	_ Node = (*UseStmt)(nil)
	_ Node = (*SetStmt)(nil)
	_ Node = (*UnsetStmt)(nil)
	_ Node = (*CommentStmt)(nil)
	_ Node = (*TruncateStmt)(nil)
)
