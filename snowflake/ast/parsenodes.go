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
