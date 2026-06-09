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
	UnaryMinus         UnaryOp = iota // -
	UnaryPlus                         // +
	UnaryNot                          // NOT
	UnaryPrior                        // PRIOR (CONNECT BY hierarchical queries)
	UnaryConnectByRoot                // CONNECT_BY_ROOT (CONNECT BY hierarchical queries)
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
	case UnaryPrior:
		return "PRIOR"
	case UnaryConnectByRoot:
		return "CONNECT_BY_ROOT"
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

// DollarRef represents a Snowflake '$'-reference used as a primary expression:
//
//   - $N          positional column reference (Positional=true, Name="N").
//     Used in COPY-transform / stage queries (SELECT $1, $2 FROM @stage) and
//     anywhere a column may appear.
//   - $<name>     session-variable / scripting-variable reference
//     (Positional=false), e.g. $min, $Variable1.
//   - <qual>.$N   positional reference qualified by a table alias or name
//     (e.g. d.$1), with the leading qualifier captured in Qualifier.
//
// The lexer strips the leading '$' and emits a single tokVariable whose text is
// stored in Name; Positional records whether that text is all decimal digits.
type DollarRef struct {
	Qualifier  *ObjectName // optional leading qualifier (e.g. d in d.$1); nil when bare
	Name       string      // text after '$': digits for positional, identifier otherwise
	Positional bool        // true when Name is all decimal digits ($1, $2, ...)
	Loc        Loc
}

func (n *DollarRef) Tag() NodeTag { return T_DollarRef }

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
	// StartWith / ConnectBy implement hierarchical queries
	// (CONNECT BY [, START WITH]). StartWith holds the optional START WITH
	// condition; ConnectBy holds the comma-separated CONNECT BY conditions
	// (PRIOR appears as a *UnaryExpr with Op == UnaryPrior inside them).
	// ConnectBy is non-nil whenever the query has a CONNECT BY clause.
	StartWith Node         // START WITH condition; nil if absent
	ConnectBy []Node       // CONNECT BY conditions; nil if absent
	OrderBy   []*OrderItem // ORDER BY; nil if absent
	Limit     Node         // LIMIT n; nil if absent
	Offset    Node         // OFFSET n; nil if absent
	Fetch     *FetchClause // FETCH FIRST/NEXT; nil if absent
	Loc       Loc
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
	// Table-attached clauses (Snowflake). Any combination may appear; the
	// documented source order is AT/BEFORE → CHANGES → MATCH_RECOGNIZE →
	// PIVOT/UNPIVOT → alias → SAMPLE.
	TimeTravel     *TimeTravelClause     // AT(...) / BEFORE(...); nil if absent
	Changes        *ChangesClause        // CHANGES(...) AT|BEFORE(...); nil if absent
	MatchRecognize *MatchRecognizeClause // MATCH_RECOGNIZE(...); nil if absent
	Pivot          *PivotClause          // PIVOT(...); nil if absent
	Unpivot        *UnpivotClause        // UNPIVOT(...); nil if absent
	Sample         *SampleClause         // SAMPLE/TABLESAMPLE(...); nil if absent
	Loc            Loc
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
// Table-attached query clauses (T5.3)
// ---------------------------------------------------------------------------

// PivotClause represents a PIVOT applied to a table source:
//
//	PIVOT ( <agg>(<col>) [AS <alias>]
//	        FOR <col> IN ( <values> | ANY [ORDER BY …] | <subquery> )
//	        [DEFAULT ON NULL (<expr>)] ) [ [AS] alias ]
type PivotClause struct {
	Agg        *FuncCallExpr // the aggregate, e.g. SUM(amount)
	AggAlias   Ident         // alias for the aggregate (AS total); zero if absent
	ForColumn  *ColumnRef    // the FOR pivot column
	In         *PivotInClause
	DefaultVal Node  // DEFAULT ON NULL (<expr>); nil if absent
	Alias      Ident // trailing [AS] alias for the pivot result; zero if absent
	Loc        Loc
}

func (n *PivotClause) Tag() NodeTag { return T_PivotClause }

// PivotInKind distinguishes the three forms of the PIVOT … IN (…) list.
type PivotInKind int

const (
	PivotInValues   PivotInKind = iota // IN ( v1 [AS a1], v2 [AS a2], … )
	PivotInAny                         // IN ( ANY [ORDER BY …] )
	PivotInSubquery                    // IN ( <subquery> )
)

// PivotInClause is the IN (...) part of a PIVOT.
type PivotInClause struct {
	Kind     PivotInKind
	Values   []*PivotValue // for PivotInValues
	OrderBy  []*OrderItem  // for PivotInAny: ANY ORDER BY …; nil otherwise
	Subquery Node          // for PivotInSubquery
	Loc      Loc
}

func (n *PivotInClause) Tag() NodeTag { return T_PivotInClause }

// PivotValue is one value in a PIVOT … IN ( v [AS alias], … ) list.
type PivotValue struct {
	Value Node  // the value expression (typically a literal)
	Alias Ident // AS alias; zero if absent
	Loc   Loc
}

func (n *PivotValue) Tag() NodeTag { return T_PivotValue }

// UnpivotClause represents an UNPIVOT applied to a table source:
//
//	UNPIVOT [ {INCLUDE|EXCLUDE} NULLS ]
//	  ( <value_col> FOR <name_col> IN ( <col> [AS alias], … ) ) [alias]
type UnpivotClause struct {
	NullsMode   UnpivotNullsMode // INCLUDE / EXCLUDE NULLS; default = unspecified
	ValueColumn Ident            // the new value column name
	NameColumn  Ident            // the new name column name (FOR <name_col>)
	Columns     []*UnpivotColumn // the IN (...) source columns
	Alias       Ident            // trailing alias; zero if absent
	Loc         Loc
}

func (n *UnpivotClause) Tag() NodeTag { return T_UnpivotClause }

// UnpivotNullsMode encodes the optional {INCLUDE|EXCLUDE} NULLS modifier.
type UnpivotNullsMode int

const (
	UnpivotNullsUnspecified UnpivotNullsMode = iota // no modifier (defaults to EXCLUDE)
	UnpivotIncludeNulls                             // INCLUDE NULLS
	UnpivotExcludeNulls                             // EXCLUDE NULLS
)

// UnpivotColumn is one column in an UNPIVOT … IN ( col [AS alias], … ) list.
type UnpivotColumn struct {
	Column Ident // the source column
	Alias  Ident // AS alias; zero if absent
	Loc    Loc
}

func (n *UnpivotColumn) Tag() NodeTag { return T_UnpivotColumn }

// ---------------------------------------------------------------------------
// MATCH_RECOGNIZE
// ---------------------------------------------------------------------------

// MatchRecognizeClause represents a MATCH_RECOGNIZE applied to a table source.
//
//	MATCH_RECOGNIZE ( [PARTITION BY …] [ORDER BY …] [MEASURES …]
//	  [{ONE ROW | ALL ROWS} PER MATCH [match-opts]]
//	  [AFTER MATCH SKIP …] PATTERN ( <row_pattern> )
//	  [DEFINE <var> AS <cond>, …] ) [ [AS] alias ]
type MatchRecognizeClause struct {
	PartitionBy  []Node          // PARTITION BY exprs; nil if absent
	OrderBy      []*OrderItem    // ORDER BY items; nil if absent
	Measures     []*MatchMeasure // MEASURES list; nil if absent
	RowsPerMatch *RowsPerMatch   // ONE ROW / ALL ROWS PER MATCH; nil if absent
	AfterMatch   *AfterMatchSkip // AFTER MATCH SKIP …; nil if absent
	Pattern      *RowPattern     // PATTERN ( … ); nil only on malformed input
	Define       []*MatchDefine  // DEFINE list; nil if absent
	Alias        Ident           // trailing [AS] alias; zero if absent
	Loc          Loc
}

func (n *MatchRecognizeClause) Tag() NodeTag { return T_MatchRecognizeClause }

// MatchMeasure is one item in a MEASURES list: [FINAL|RUNNING] <expr> [AS] <alias>.
type MatchMeasure struct {
	Semantics MatchSemantics // FINAL / RUNNING / unspecified
	Expr      Node           // the measure expression
	Alias     Ident          // the column alias
	Loc       Loc
}

func (n *MatchMeasure) Tag() NodeTag { return T_MatchMeasure }

// MatchSemantics encodes the optional FINAL / RUNNING prefix on a measure.
type MatchSemantics int

const (
	MatchSemanticsUnspecified MatchSemantics = iota
	MatchSemanticsFinal                      // FINAL
	MatchSemanticsRunning                    // RUNNING
)

// RowsPerMatchKind selects ONE ROW vs ALL ROWS PER MATCH.
type RowsPerMatchKind int

const (
	OneRowPerMatch  RowsPerMatchKind = iota // ONE ROW PER MATCH
	AllRowsPerMatch                         // ALL ROWS PER MATCH
)

// RowsPerMatchOpt encodes the optional modifier on ALL ROWS PER MATCH.
type RowsPerMatchOpt int

const (
	RowsPerMatchOptNone       RowsPerMatchOpt = iota
	RowsPerMatchShowEmpty                     // SHOW EMPTY MATCHES
	RowsPerMatchOmitEmpty                     // OMIT EMPTY MATCHES
	RowsPerMatchWithUnmatched                 // WITH UNMATCHED ROWS
)

// RowsPerMatch represents the {ONE ROW | ALL ROWS} PER MATCH [opt] clause.
type RowsPerMatch struct {
	Kind RowsPerMatchKind
	Opt  RowsPerMatchOpt // only meaningful for ALL ROWS PER MATCH
	Loc  Loc
}

func (n *RowsPerMatch) Tag() NodeTag { return T_RowsPerMatch }

// AfterMatchKind selects the AFTER MATCH SKIP target.
type AfterMatchKind int

const (
	AfterMatchSkipPastLastRow AfterMatchKind = iota // SKIP PAST LAST ROW
	AfterMatchSkipToNextRow                         // SKIP TO NEXT ROW
	AfterMatchSkipToFirst                           // SKIP TO FIRST <var>
	AfterMatchSkipToLast                            // SKIP TO LAST <var>
	AfterMatchSkipToVar                             // SKIP TO <var>
)

// AfterMatchSkip represents the AFTER MATCH SKIP … clause.
type AfterMatchSkip struct {
	Kind   AfterMatchKind
	Symbol Ident // the target pattern variable for the TO … forms; zero otherwise
	Loc    Loc
}

func (n *AfterMatchSkip) Tag() NodeTag { return T_AfterMatchSkip }

// MatchDefine is one DEFINE entry: <var> AS <condition>.
type MatchDefine struct {
	Symbol Ident // the pattern variable being defined
	Cond   Node  // the boolean condition expression
	Loc    Loc
}

func (n *MatchDefine) Tag() NodeTag { return T_MatchDefine }

// RowPattern holds the PATTERN ( … ) body. The body is kept as raw text
// (Raw, with surrounding parentheses stripped) because the full row-pattern
// grammar — quantifiers, alternation, PERMUTE, anchors — has no executable
// oracle to validate a structured tree against; consumers that need the
// pattern re-lex Raw. RawLoc spans the inner text.
type RowPattern struct {
	Raw    string // the verbatim pattern text between the outer parentheses
	RawLoc Loc    // location of the inner pattern text
	Loc    Loc    // location spanning PATTERN ( … )
}

func (n *RowPattern) Tag() NodeTag { return T_RowPattern }

// ---------------------------------------------------------------------------
// SAMPLE / TABLESAMPLE
// ---------------------------------------------------------------------------

// SampleKeyword distinguishes the SAMPLE vs TABLESAMPLE spelling.
type SampleKeyword int

const (
	SampleKwSample      SampleKeyword = iota // SAMPLE
	SampleKwTablesample                      // TABLESAMPLE
)

// SampleMethod is the optional sampling method.
type SampleMethod int

const (
	SampleMethodUnspecified SampleMethod = iota
	SampleMethodBernoulli                // BERNOULLI
	SampleMethodRow                      // ROW (synonym of BERNOULLI)
	SampleMethodSystem                   // SYSTEM
	SampleMethodBlock                    // BLOCK (synonym of SYSTEM)
)

// SampleClause represents SAMPLE / TABLESAMPLE applied to a table source:
//
//	{SAMPLE | TABLESAMPLE} [method] ( <prob> | <n> ROWS ) [{SEED|REPEATABLE} (<n>)]
type SampleClause struct {
	Keyword  SampleKeyword
	Method   SampleMethod
	Quantity Node // probability or row count expression
	Rows     bool // true if the quantity is "<n> ROWS"
	// Seed holds the {SEED|REPEATABLE} (<n>) value; nil if absent.
	Seed     Node
	SeedKind SampleSeedKind // which keyword introduced Seed
	Loc      Loc
}

func (n *SampleClause) Tag() NodeTag { return T_SampleClause }

// SampleSeedKind records whether the seed used SEED or REPEATABLE.
type SampleSeedKind int

const (
	SampleSeedNone       SampleSeedKind = iota
	SampleSeedSeed                      // SEED (<n>)
	SampleSeedRepeatable                // REPEATABLE (<n>)
)

// ---------------------------------------------------------------------------
// Time travel: AT / BEFORE / CHANGES
// ---------------------------------------------------------------------------

// TimeTravelKind selects AT vs BEFORE.
type TimeTravelKind int

const (
	TimeTravelAt     TimeTravelKind = iota // AT (...)
	TimeTravelBefore                       // BEFORE (...)
)

// TimeTravelAnchor selects the parameter name inside AT/BEFORE.
type TimeTravelAnchor int

const (
	TimeTravelTimestamp TimeTravelAnchor = iota // TIMESTAMP => expr
	TimeTravelOffset                            // OFFSET => expr
	TimeTravelStatement                         // STATEMENT => expr
	TimeTravelStream                            // STREAM => expr
)

// TimeTravelClause represents { AT | BEFORE } ( <anchor> => <expr> ).
type TimeTravelClause struct {
	Kind   TimeTravelKind
	Anchor TimeTravelAnchor
	Expr   Node // the anchor value expression
	Loc    Loc
}

func (n *TimeTravelClause) Tag() NodeTag { return T_TimeTravelClause }

// ChangesInfo selects DEFAULT vs APPEND_ONLY for CHANGES(INFORMATION => …).
type ChangesInfo int

const (
	ChangesDefault    ChangesInfo = iota // INFORMATION => DEFAULT
	ChangesAppendOnly                    // INFORMATION => APPEND_ONLY
)

// ChangesClause represents:
//
//	CHANGES ( INFORMATION => {DEFAULT|APPEND_ONLY} ) { AT|BEFORE } (…) [END (…)]
type ChangesClause struct {
	Info  ChangesInfo
	Start *TimeTravelClause // the required AT/BEFORE anchor
	End   *TimeTravelClause // optional END (…); nil if absent
	Loc   Loc
}

func (n *ChangesClause) Tag() NodeTag { return T_ChangesClause }

// Compile-time assertions for the T5.3 clause node types.
var (
	_ Node = (*PivotClause)(nil)
	_ Node = (*PivotInClause)(nil)
	_ Node = (*PivotValue)(nil)
	_ Node = (*UnpivotClause)(nil)
	_ Node = (*UnpivotColumn)(nil)
	_ Node = (*MatchRecognizeClause)(nil)
	_ Node = (*MatchMeasure)(nil)
	_ Node = (*RowsPerMatch)(nil)
	_ Node = (*AfterMatchSkip)(nil)
	_ Node = (*MatchDefine)(nil)
	_ Node = (*RowPattern)(nil)
	_ Node = (*SampleClause)(nil)
	_ Node = (*TimeTravelClause)(nil)
	_ Node = (*ChangesClause)(nil)
)

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
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
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
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
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
	OrAlter       bool // CREATE OR ALTER (mutually exclusive with OrReplace)
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
	DropNetworkRule                     // DROP NETWORK RULE
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
	case DropNetworkRule:
		return "NETWORK RULE"
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
// TCL (transaction control) statement nodes (T6.2)
// ---------------------------------------------------------------------------

// BeginKind distinguishes the surface form that opened a transaction.
//
// Snowflake treats START TRANSACTION as a synonym for BEGIN, and the optional
// WORK / TRANSACTION modifier after BEGIN is purely for cross-database
// compatibility. The kind preserves the original keyword(s) so deparse and
// tooling can round-trip the exact form the user wrote.
type BeginKind int

const (
	// BeginBare is `BEGIN` with no WORK/TRANSACTION modifier.
	BeginBare BeginKind = iota
	// BeginWork is `BEGIN WORK`.
	BeginWork
	// BeginTransaction is `BEGIN TRANSACTION`.
	BeginTransaction
	// BeginStartTransaction is `START TRANSACTION`.
	BeginStartTransaction
)

// String returns the keyword(s) that introduced the transaction.
func (k BeginKind) String() string {
	switch k {
	case BeginBare:
		return "BEGIN"
	case BeginWork:
		return "BEGIN WORK"
	case BeginTransaction:
		return "BEGIN TRANSACTION"
	case BeginStartTransaction:
		return "START TRANSACTION"
	default:
		return "UNKNOWN"
	}
}

// BeginStmt represents the transaction-opening statement, covering both
// surface forms documented by Snowflake:
//
//	BEGIN [ { WORK | TRANSACTION } ] [ NAME <name> ]
//	START TRANSACTION [ NAME <name> ]
//
// Kind records which form/modifier was used. Name is the optional transaction
// name following the NAME keyword; it is the zero Ident (IsEmpty) when absent.
// Name is an Ident value (not a Node) and is therefore not visited by the
// walker, matching how ObjectName embeds its Ident parts.
type BeginStmt struct {
	Kind BeginKind
	Name Ident // optional; zero Ident (IsEmpty) when no NAME clause
	Loc  Loc
}

// Tag implements Node.
func (n *BeginStmt) Tag() NodeTag { return T_BeginStmt }

// CommitStmt represents `COMMIT [ WORK ]`. Work records whether the optional
// WORK keyword (a cross-database-compatibility no-op) was present.
type CommitStmt struct {
	Work bool
	Loc  Loc
}

// Tag implements Node.
func (n *CommitStmt) Tag() NodeTag { return T_CommitStmt }

// RollbackStmt represents `ROLLBACK [ WORK ]`. Work records whether the
// optional WORK keyword (a cross-database-compatibility no-op) was present.
type RollbackStmt struct {
	Work bool
	Loc  Loc
}

// Tag implements Node.
func (n *RollbackStmt) Tag() NodeTag { return T_RollbackStmt }

// Compile-time assertions.
var (
	_ Node = (*BeginStmt)(nil)
	_ Node = (*CommitStmt)(nil)
	_ Node = (*RollbackStmt)(nil)
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
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
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
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
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

// ---------------------------------------------------------------------------
// Bulk data movement — COPY INTO / PUT / GET / LIST / REMOVE (T5.2)
// ---------------------------------------------------------------------------
//
// These statements move data between Snowflake tables, stages, and external
// cloud storage. Their option vocabularies (copyOptions, file-format options,
// PUT/GET transfer options) are large and grow with new Snowflake versions, so
// — mirroring the merged GRANT/REVOKE and SHOW grammars — option names are NOT
// enumerated. copyOptions and FILE_FORMAT type-options are captured as
// open-ended `KEY = <value>` pairs (CopyOption), where the value is a verbatim
// token run, a literal, a parenthesized group, or a nested option list. The
// catalog/semantic layer, not the parser, validates that an option is legal.

// StageLocation is a file location that is the source or destination of a bulk
// data-movement statement. It captures the three syntactic flavors Snowflake
// accepts wherever a stage/location may appear:
//
//   - Stage      — an internal/external stage reference beginning with '@':
//     @[ns.]name[/path], @[ns.]%table[/path], @~[/path].
//   - External   — a cloud-storage URI single-quoted string:
//     's3://...', 'gcs://...', 'azure://...' (any 'scheme://...').
//   - LocalFile  — a file:// URI (PUT source / GET target) or a bare local path.
//
// Exactly one of the kind-specific fields is meaningful per Kind. Raw holds the
// verbatim source text of the location in every case (e.g. "@mystage/a.csv",
// "s3://b/p", "file:///tmp/x.csv"), which downstream consumers can use without
// re-deriving it.
type StageLocation struct {
	Kind StageLocationKind
	Raw  string // verbatim source text of the location
	Loc  Loc
}

// Tag implements Node.
func (n *StageLocation) Tag() NodeTag { return T_StageLocation }

// StageLocationKind classifies a StageLocation.
type StageLocationKind int

const (
	// StageRef is an '@'-prefixed internal or external stage reference.
	StageRef StageLocationKind = iota
	// StageExternal is a cloud-storage URI ('s3://...' / 'gcs://...' / 'azure://...').
	StageExternal
	// StageLocalFile is a file:// URI or a bare local filesystem path.
	StageLocalFile
)

// CopyOption is one open-ended `KEY = <value>` option used by COPY (copyOptions
// + FILE_FORMAT + FILES + PATTERN + STORAGE_INTEGRATION + CREDENTIALS + ...) and,
// for transfer statements, PUT/GET (PARALLEL, AUTO_COMPRESS, ...). The option
// name is uppercased; the value is one of the mutually-exclusive forms below.
// A bare option (e.g. HEADER, FORCE with no `=`) has Bare=true and no value.
type CopyOption struct {
	Name string // uppercased option name (e.g. "ON_ERROR", "FILE_FORMAT")
	Bare bool   // true for a value-less option (e.g. bare HEADER)

	// Exactly one of the following is set when Bare is false:
	Words string        // verbatim uppercased token run (e.g. "CASE_INSENSITIVE", "TRUE", "myint")
	Lit   *Literal      // a string/number literal value (e.g. 'RETURN_ROWS', 32000000)
	Group []*CopyOption // a nested ( k = v, ... ) group (e.g. FILE_FORMAT, CREDENTIALS, INCLUDE_METADATA)
	List  []*Literal    // a ( 'a', 'b', ... ) literal list (e.g. FILES = ('f1', 'f2'))

	Loc Loc
}

// Tag implements Node.
func (n *CopyOption) Tag() NodeTag { return T_CopyOption }

// CopyIntoTableStmt represents COPY INTO <table> (a data load):
//
//	COPY INTO [ns.]table [ (cols) ]
//	  FROM { stage | location | ( SELECT ...$col... FROM stage ) }
//	  [ FILES = (...) ] [ PATTERN = '...' ] [ FILE_FORMAT = (...) ]
//	  [ copyOptions ] [ VALIDATION_MODE = ... ]
//
// When the source is a transformation query, Transform holds the SELECT and
// From is nil; otherwise From holds the stage/location source. Options holds
// every option that followed the source (FILES / PATTERN / FILE_FORMAT /
// copyOptions / VALIDATION_MODE), each as a CopyOption, preserving source order.
type CopyIntoTableStmt struct {
	Target    *ObjectName
	Columns   []Ident        // optional ( col, ... ) before FROM; nil if absent
	From      *StageLocation // stage/location source; nil when Transform is set
	Transform *SelectStmt    // FROM ( SELECT ... ) transformation; nil for the plain form
	Options   []*CopyOption  // FILES / PATTERN / FILE_FORMAT / copyOptions / VALIDATION_MODE
	Loc       Loc
}

// Tag implements Node.
func (n *CopyIntoTableStmt) Tag() NodeTag { return T_CopyIntoTableStmt }

// CopyIntoLocationStmt represents COPY INTO <location> (a data unload):
//
//	COPY INTO { stage | location }
//	  FROM { [ns.]table | ( query ) }
//	  [ PARTITION BY <expr> ] [ FILE_FORMAT = (...) ] [ copyOptions ]
//	  [ VALIDATION_MODE = RETURN_ROWS ] [ HEADER ]
//
// Either FromTable or FromQuery is set (the unload source). Partition holds the
// PARTITION BY expression (nil if absent). Options holds FILE_FORMAT /
// copyOptions / VALIDATION_MODE / HEADER, preserving source order.
type CopyIntoLocationStmt struct {
	Into      *StageLocation
	FromTable *ObjectName   // FROM <table>; nil when FromQuery is set
	FromQuery Node          // FROM ( query ); nil when FromTable is set
	Partition Node          // PARTITION BY <expr>; nil if absent
	Options   []*CopyOption // FILE_FORMAT / copyOptions / VALIDATION_MODE / HEADER
	Loc       Loc
}

// Tag implements Node.
func (n *CopyIntoLocationStmt) Tag() NodeTag { return T_CopyIntoLocationStmt }

// PutStmt represents a PUT statement (upload a local file to an internal stage):
//
//	PUT file://<path> <stage> [ PARALLEL = n ] [ AUTO_COMPRESS = b ]
//	    [ SOURCE_COMPRESSION = kw ] [ OVERWRITE = b ]
type PutStmt struct {
	File    *StageLocation // the file:// (or bare local) source
	Stage   *StageLocation // the internal-stage destination
	Options []*CopyOption  // PARALLEL / AUTO_COMPRESS / SOURCE_COMPRESSION / OVERWRITE
	Loc     Loc
}

// Tag implements Node.
func (n *PutStmt) Tag() NodeTag { return T_PutStmt }

// GetStmt represents a GET statement (download files from an internal stage):
//
//	GET <stage> file://<dir> [ PARALLEL = n ] [ PATTERN = '...' ]
//
// Target is the local destination directory: a file:// URI or a bare path.
type GetStmt struct {
	Stage   *StageLocation // the internal-stage source
	Target  *StageLocation // the local destination directory
	Options []*CopyOption  // PARALLEL / PATTERN
	Loc     Loc
}

// Tag implements Node.
func (n *GetStmt) Tag() NodeTag { return T_GetStmt }

// ListStmt represents LIST / LS (enumerate the files in a stage):
//
//	{ LIST | LS } <stage> [ PATTERN = '<regex>' ]
//
// Short is true when the LS alias was used.
type ListStmt struct {
	Short   bool
	Stage   *StageLocation
	Pattern *Literal // PATTERN = '<regex>'; nil if absent
	Loc     Loc
}

// Tag implements Node.
func (n *ListStmt) Tag() NodeTag { return T_ListStmt }

// RemoveStmt represents REMOVE / RM (delete files from a stage):
//
//	{ REMOVE | RM } <stage> [ PATTERN = '<regex>' ]
//
// Short is true when the RM alias was used.
type RemoveStmt struct {
	Short   bool
	Stage   *StageLocation
	Pattern *Literal // PATTERN = '<regex>'; nil if absent
	Loc     Loc
}

// Tag implements Node.
func (n *RemoveStmt) Tag() NodeTag { return T_RemoveStmt }

// Compile-time assertions for bulk data-movement nodes.
var (
	_ Node = (*StageLocation)(nil)
	_ Node = (*CopyOption)(nil)
	_ Node = (*CopyIntoTableStmt)(nil)
	_ Node = (*CopyIntoLocationStmt)(nil)
	_ Node = (*PutStmt)(nil)
	_ Node = (*GetStmt)(nil)
	_ Node = (*ListStmt)(nil)
	_ Node = (*RemoveStmt)(nil)
)

// ---------------------------------------------------------------------------
// Stage DDL — CREATE / ALTER STAGE (T4.1)
// ---------------------------------------------------------------------------
//
// Like COPY, the stage grammar carries large, version-growing option
// vocabularies: the external-stage cloud params (URL, STORAGE_INTEGRATION,
// CREDENTIALS, ENCRYPTION, ENDPOINT, AWS_ACCESS_POINT_ARN,
// USE_PRIVATELINK_ENDPOINT, ...), the directory-table params (DIRECTORY = (...)),
// the file-format params (FILE_FORMAT = (FORMAT_NAME = ... | TYPE = ... ...)) and
// the copy options (COPY_OPTIONS = (...)). Rather than enumerate them — the
// legacy ANTLR grammar's finite lists are already stale relative to the docs
// (they lack AWS_ACCESS_POINT_ARN, USE_PRIVATELINK_ENDPOINT, ENDPOINT, the
// s3compat:// form, ...) — every such param is captured as an open-ended
// `KEY = <value>` pair (CopyOption), reusing the merged COPY (T5.2) machinery.
// Only the structurally-distinct WITH TAG and COMMENT clauses, and the ALTER
// action keywords (RENAME / SET / UNSET / REFRESH), anchor the grammar. The
// catalog/semantic layer, not the parser, validates that an option is legal.

// CreateStageStmt represents
//
//	CREATE [ OR REPLACE ] [ TEMP | TEMPORARY ] STAGE [ IF NOT EXISTS ] <name>
//	  <stageParams>            -- internal: ENCRYPTION=(...); external: URL=...,
//	                              STORAGE_INTEGRATION=..., CREDENTIALS=(...),
//	                              ENCRYPTION=(...), ...
//	  [ DIRECTORY = ( ... ) ]
//	  [ FILE_FORMAT = ( ... ) ]
//	  [ COPY_OPTIONS = ( ... ) ]
//	  [ COMMENT = '<string>' ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//
// All cloud/format/copy/comment params are captured open-ended in Options,
// preserving source order (an internal stage simply has no URL option). Tags
// holds the WITH TAG assignments. There is no dedicated External flag: an
// external stage is exactly one that carries a URL option, which downstream
// consumers can detect.
type CreateStageStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Temporary   bool
	IfNotExists bool
	Name        *ObjectName
	Options     []*CopyOption    // URL / STORAGE_INTEGRATION / CREDENTIALS / ENCRYPTION / DIRECTORY / FILE_FORMAT / COPY_OPTIONS / COMMENT / ...
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Loc         Loc
}

func (n *CreateStageStmt) Tag() NodeTag { return T_CreateStageStmt }

// AlterStageAction discriminates the action variants of ALTER STAGE.
type AlterStageAction int

const (
	AlterStageRename   AlterStageAction = iota // RENAME TO <new_name>
	AlterStageSet                              // SET <stageParams|FILE_FORMAT|COPY_OPTIONS|COMMENT|DIRECTORY>
	AlterStageUnset                            // UNSET <property names> (e.g. DCM PROJECT)
	AlterStageSetTag                           // SET TAG (...)
	AlterStageUnsetTag                         // UNSET TAG (...)
	AlterStageRefresh                          // REFRESH [ SUBPATH = '<relative-path>' ]
)

// AlterStageStmt represents ALTER STAGE [IF EXISTS] <name> <action>.
//
//	RENAME TO <new_name>
//	SET <options>                       -- open-ended KEY = value params
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]          -- e.g. UNSET DCM PROJECT
//	REFRESH [ SUBPATH = '<relative-path>' ]
type AlterStageStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterStageAction
	NewName    *ObjectName      // RENAME TO target; non-nil for AlterStageRename
	Options    []*CopyOption    // SET <options>; for AlterStageSet
	Tags       []*TagAssignment // SET TAG (...) assignments; for AlterStageSetTag
	UnsetTags  []*ObjectName    // UNSET TAG (...) names; for AlterStageUnsetTag
	UnsetProps []string         // UNSET <property> names; for AlterStageUnset
	Subpath    *string          // REFRESH SUBPATH = '...'; nil when absent
	Loc        Loc
}

func (n *AlterStageStmt) Tag() NodeTag { return T_AlterStageStmt }

// Compile-time assertions for stage DDL nodes.
var (
	_ Node = (*CreateStageStmt)(nil)
	_ Node = (*AlterStageStmt)(nil)
)

// ---------------------------------------------------------------------------
// File-format DDL — CREATE / ALTER FILE FORMAT (T4.2)
// ---------------------------------------------------------------------------
//
// A named file format carries a TYPE ({ CSV | JSON | AVRO | ORC | PARQUET |
// XML }) plus a large, type-dependent and version-growing vocabulary of
// formatTypeOptions (CSV: FIELD_DELIMITER / SKIP_HEADER /
// FIELD_OPTIONALLY_ENCLOSED_BY / NULL_IF / ...; JSON: STRIP_OUTER_ARRAY / ...;
// PARQUET: BINARY_AS_TEXT / USE_VECTORIZED_SCANNER / ...; etc.). Rather than
// enumerate those options — the legacy ANTLR grammar's format_type_options list
// is already stale relative to the docs (it lacks USE_VECTORIZED_SCANNER,
// USE_LOGICAL_TYPE, the per-type COMPRESSION variants, ...) — every option that
// follows the TYPE clause is captured as an open-ended `KEY = <value>` pair
// (CopyOption), reusing the merged COPY (T5.2) machinery. The TYPE keyword is
// the one structural anchor (it selects which options are legal); COMMENT is
// captured open-ended in Options. The catalog/semantic layer, not the parser,
// validates that an option is real and legal for the chosen TYPE. This mirrors
// the merged STAGE (T4.1) / COPY (T5.2) open-ended approach.

// CreateFileFormatStmt represents
//
//	CREATE [ OR REPLACE ] [ { TEMP | TEMPORARY | VOLATILE } ]
//	  FILE FORMAT [ IF NOT EXISTS ] <name>
//	  [ TYPE = { CSV | JSON | AVRO | ORC | PARQUET | XML } ]
//	  [ formatTypeOptions ]
//	  [ COMMENT = '<string_literal>' ]
//
// Type holds the uppercased format type word (e.g. "CSV"); it is "" when the
// TYPE clause is omitted (TYPE is optional per the docs, defaulting to CSV).
// Options holds the formatTypeOptions and the COMMENT clause, each a CopyOption,
// preserving source order. There is no dedicated External/Volatile distinction
// beyond the Temporary flag: VOLATILE is a synonym of TEMPORARY here and sets
// Temporary.
type CreateFileFormatStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Temporary   bool
	IfNotExists bool
	Name        *ObjectName
	Type        string        // uppercased TYPE word (CSV/JSON/AVRO/ORC/PARQUET/XML); "" if omitted
	TypeLoc     Loc           // Loc of the TYPE value word; zero when Type is ""
	Options     []*CopyOption // formatTypeOptions + COMMENT, preserving source order
	Loc         Loc
}

// Tag implements Node.
func (n *CreateFileFormatStmt) Tag() NodeTag { return T_CreateFileFormatStmt }

// AlterFileFormatAction discriminates the action variants of ALTER FILE FORMAT.
type AlterFileFormatAction int

const (
	AlterFileFormatRename AlterFileFormatAction = iota // RENAME TO <new_name>
	AlterFileFormatSet                                 // SET <formatTypeOptions> [COMMENT=...]
)

// AlterFileFormatStmt represents
//
//	ALTER FILE FORMAT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER FILE FORMAT [ IF EXISTS ] <name> SET { [ formatTypeOptions ] [ COMMENT = '<string_literal>' ] }
//
// The SET form is unparenthesized (per the docs and the legacy ANTLR grammar):
// the formatTypeOptions and COMMENT follow SET directly. ALTER FILE FORMAT
// supports no UNSET, no SET TAG, and no RENAME-plus-SET combination.
type AlterFileFormatStmt struct {
	IfExists bool
	Name     *ObjectName
	Action   AlterFileFormatAction
	NewName  *ObjectName   // RENAME TO target; non-nil for AlterFileFormatRename
	Options  []*CopyOption // SET <formatTypeOptions> [COMMENT=...]; for AlterFileFormatSet
	Loc      Loc
}

// Tag implements Node.
func (n *AlterFileFormatStmt) Tag() NodeTag { return T_AlterFileFormatStmt }

// Compile-time assertions for file-format DDL nodes.
var (
	_ Node = (*CreateFileFormatStmt)(nil)
	_ Node = (*AlterFileFormatStmt)(nil)
)

// ---------------------------------------------------------------------------
// Account-level integration-object DDL — CREATE / ALTER for STORAGE / API /
// NOTIFICATION / SECURITY INTEGRATION, RESOURCE MONITOR, SECRET, CONNECTION,
// EXTERNAL VOLUME, GIT REPOSITORY (T4.7)
// ---------------------------------------------------------------------------
//
// These account/schema-level objects all configure an external resource through
// a large, version-growing vocabulary of `KEY = <value>` parameters: a STORAGE
// INTEGRATION carries STORAGE_PROVIDER / STORAGE_AWS_ROLE_ARN /
// STORAGE_ALLOWED_LOCATIONS / ...; an API INTEGRATION API_PROVIDER /
// API_ALLOWED_PREFIXES / ...; a SECRET TYPE / OAUTH_SCOPES / SECRET_STRING / ...;
// a GIT REPOSITORY ORIGIN / API_INTEGRATION / GIT_CREDENTIALS; an EXTERNAL VOLUME
// STORAGE_LOCATIONS / ALLOW_WRITES; and so on. Rather than mirror the legacy
// ANTLR grammar's finite, already-stale per-type enumerations (its
// create_*_integration rules pin a fixed option order and lack newer params, and
// it has no rule for CREATE EXTERNAL VOLUME at all), every parameter that follows
// the object name is captured as an open-ended `KEY = <value>` pair (CopyOption),
// reusing the merged COPY (T5.2) / STAGE (T4.1) / FILE FORMAT (T4.2) machinery.
// Only the structural anchors are modeled as dedicated fields: the object Kind,
// the RESOURCE MONITOR WITH keyword + its TRIGGERS clause, the GIT REPOSITORY
// [WITH] TAG clause, and (on ALTER) the action keyword and EXTERNAL VOLUME's
// ADD/REMOVE/UPDATE STORAGE_LOCATION actions. The catalog/semantic layer, not the
// parser, validates that an option is real and legal for the chosen Kind.

// IntegrationObjectKind discriminates the account-level object types handled by
// the T4.7 CREATE / ALTER parsers.
type IntegrationObjectKind int

const (
	StorageIntegration      IntegrationObjectKind = iota // [STORAGE] INTEGRATION
	APIIntegration                                       // API INTEGRATION
	NotificationIntegration                              // NOTIFICATION INTEGRATION
	SecurityIntegration                                  // SECURITY INTEGRATION
	ResourceMonitor                                      // RESOURCE MONITOR
	Secret                                               // SECRET
	Connection                                           // CONNECTION
	ExternalVolume                                       // EXTERNAL VOLUME
	GitRepository                                        // GIT REPOSITORY
)

// String returns the SQL object keyword(s) for the kind (uppercased), used for
// diagnostics and deparse.
func (k IntegrationObjectKind) String() string {
	switch k {
	case StorageIntegration:
		return "STORAGE INTEGRATION"
	case APIIntegration:
		return "API INTEGRATION"
	case NotificationIntegration:
		return "NOTIFICATION INTEGRATION"
	case SecurityIntegration:
		return "SECURITY INTEGRATION"
	case ResourceMonitor:
		return "RESOURCE MONITOR"
	case Secret:
		return "SECRET"
	case Connection:
		return "CONNECTION"
	case ExternalVolume:
		return "EXTERNAL VOLUME"
	case GitRepository:
		return "GIT REPOSITORY"
	default:
		return "INTEGRATION OBJECT"
	}
}

// ResourceMonitorTrigger is one RESOURCE MONITOR trigger definition:
//
//	ON <threshold> PERCENT DO { SUSPEND | SUSPEND_IMMEDIATE | NOTIFY }
//
// Triggers are space-separated in the source (no comma). Action is the uppercased
// action keyword.
type ResourceMonitorTrigger struct {
	Threshold int64  // the ON <threshold> PERCENT value
	Action    string // "SUSPEND" | "SUSPEND_IMMEDIATE" | "NOTIFY"
	Loc       Loc
}

// Tag implements Node.
func (n *ResourceMonitorTrigger) Tag() NodeTag { return T_ResourceMonitorTrigger }

// ConnectionReplica holds a CREATE CONNECTION ... AS REPLICA OF
// <organization>.<account>.<connection> clause. The three parts are captured as
// an ObjectName (so quoting / dotted-name handling is shared with every other
// name in the AST).
type ConnectionReplica struct {
	Source *ObjectName // organization_name.account_name.connection_name
	Loc    Loc
}

// Tag implements Node.
func (n *ConnectionReplica) Tag() NodeTag { return T_ConnectionReplica }

// CreateIntegrationStmt represents the CREATE form of every T4.7 object:
//
//	CREATE [ OR REPLACE ] { STORAGE | API | NOTIFICATION | SECURITY } INTEGRATION [ IF NOT EXISTS ] <name> <options...>
//	CREATE [ OR REPLACE ] RESOURCE MONITOR [ IF NOT EXISTS ] <name> WITH <options...> [ TRIGGERS <trigger>... ]
//	CREATE [ OR REPLACE ] SECRET [ IF NOT EXISTS ] <name> <options...>
//	CREATE CONNECTION [ IF NOT EXISTS ] <name> [ AS REPLICA OF <a.b.c> ] [ COMMENT = '...' ]
//	CREATE [ OR REPLACE ] EXTERNAL VOLUME [ IF NOT EXISTS ] <name> STORAGE_LOCATIONS = (...) [ ALLOW_WRITES = ... ] [ COMMENT = '...' ]
//	CREATE [ OR REPLACE ] GIT REPOSITORY [ IF NOT EXISTS ] <name> <options...> [ [ WITH ] TAG (...) ]
//
// All configuration parameters are captured open-ended in Options, preserving
// source order. The structural anchors are modeled separately: With records the
// RESOURCE MONITOR WITH keyword; Triggers holds its TRIGGERS clause; Tags holds a
// GIT REPOSITORY [WITH] TAG clause; Replica holds a CONNECTION AS REPLICA OF
// clause. Per the docs OR REPLACE and IF NOT EXISTS are mutually exclusive for
// every kind, and CONNECTION supports neither OR REPLACE nor SECURE — the parser
// accepts the open-ended combination and the semantic layer enforces the
// exclusions (matching the rest of this engine's parse-permissively philosophy).
type CreateIntegrationStmt struct {
	Kind        IntegrationObjectKind
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	With        bool                      // RESOURCE MONITOR's mandatory WITH keyword was present
	Options     []*CopyOption             // open-ended KEY = value parameters, source order
	Tags        []*TagAssignment          // GIT REPOSITORY [WITH] TAG (...); nil if absent
	Triggers    []*ResourceMonitorTrigger // RESOURCE MONITOR TRIGGERS ...; nil if absent
	Replica     *ConnectionReplica        // CONNECTION AS REPLICA OF a.b.c; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateIntegrationStmt) Tag() NodeTag { return T_CreateIntegrationStmt }

// AlterIntegrationAction discriminates the action variants of the T4.7 ALTER
// parsers. Most objects share the SET / UNSET / SET TAG / UNSET TAG quartet;
// EXTERNAL VOLUME instead uses the storage-location actions, GIT REPOSITORY adds
// FETCH, RESOURCE MONITOR's SET carries NOTIFY_USERS + TRIGGERS, and CONNECTION
// adds the failover/primary actions.
type AlterIntegrationAction int

const (
	AlterIntegrationSet             AlterIntegrationAction = iota // SET <options> [ NOTIFY_USERS / TRIGGERS for RESOURCE MONITOR ]
	AlterIntegrationUnset                                         // UNSET <property> [ , ... ]
	AlterIntegrationSetTag                                        // SET TAG <tag> = '<value>' [ , ... ]
	AlterIntegrationUnsetTag                                      // UNSET TAG <tag> [ , ... ]
	AlterIntegrationFetch                                         // GIT REPOSITORY FETCH
	AlterIntegrationAddLocation                                   // EXTERNAL VOLUME ADD STORAGE_LOCATION = (...)
	AlterIntegrationRemoveLocation                                // EXTERNAL VOLUME REMOVE STORAGE_LOCATION '<name>'
	AlterIntegrationUpdateLocation                                // EXTERNAL VOLUME UPDATE STORAGE_LOCATION = '<name>' CREDENTIALS = (...)
	AlterIntegrationEnableFailover                                // CONNECTION ENABLE FAILOVER TO ACCOUNTS ...
	AlterIntegrationDisableFailover                               // CONNECTION DISABLE FAILOVER [ TO ACCOUNTS ... ]
	AlterIntegrationPrimary                                       // CONNECTION PRIMARY
)

// AlterIntegrationStmt represents the ALTER form of every T4.7 object. The set of
// legal actions depends on Kind (validated by the semantic layer, not here):
//
//	ALTER [STORAGE] INTEGRATION [ IF EXISTS ] <name> SET <options>
//	ALTER {API|NOTIFICATION|SECURITY} INTEGRATION [ IF EXISTS ] <name> SET <options>
//	ALTER ... INTEGRATION [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER ... INTEGRATION <name> { SET | UNSET } TAG ...
//	ALTER RESOURCE MONITOR [ IF EXISTS ] <name> SET <options> [ NOTIFY_USERS = (...) ] [ TRIGGERS <trigger>... ]
//	ALTER SECRET [ IF EXISTS ] <name> { SET <options> | UNSET COMMENT }
//	ALTER CONNECTION [ IF EXISTS ] <name> { SET <options> | UNSET COMMENT }
//	ALTER CONNECTION <name> { ENABLE FAILOVER TO ACCOUNTS a.b [,...] | DISABLE FAILOVER [ TO ACCOUNTS a.b [,...] ] | PRIMARY }
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> ADD STORAGE_LOCATION = (...)
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> REMOVE STORAGE_LOCATION '<name>'
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> UPDATE STORAGE_LOCATION = '<name>' CREDENTIALS = (...)
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> SET { ALLOW_WRITES | COMMENT } = ...
//	ALTER GIT REPOSITORY <name> { SET <options> | UNSET <property> [ , ... ] | FETCH }
type AlterIntegrationStmt struct {
	Kind       IntegrationObjectKind
	IfExists   bool
	Name       *ObjectName
	Action     AlterIntegrationAction
	Options    []*CopyOption             // SET <options>; for AlterIntegrationSet / Add / Update (the STORAGE_LOCATION + CREDENTIALS params)
	Tags       []*TagAssignment          // SET TAG assignments; for AlterIntegrationSetTag
	UnsetTags  []*ObjectName             // UNSET TAG names; for AlterIntegrationUnsetTag
	UnsetProps []string                  // UNSET <property> names; for AlterIntegrationUnset
	Triggers   []*ResourceMonitorTrigger // RESOURCE MONITOR SET ... TRIGGERS ...; nil if absent
	Location   string                    // EXTERNAL VOLUME REMOVE/UPDATE STORAGE_LOCATION '<name>'; "" otherwise
	Accounts   []*ObjectName             // CONNECTION {ENABLE|DISABLE} FAILOVER TO ACCOUNTS <a.b> [,...]; nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *AlterIntegrationStmt) Tag() NodeTag { return T_AlterIntegrationStmt }

// Compile-time assertions for integration-object DDL nodes.
var (
	_ Node = (*ResourceMonitorTrigger)(nil)
	_ Node = (*ConnectionReplica)(nil)
	_ Node = (*CreateIntegrationStmt)(nil)
	_ Node = (*AlterIntegrationStmt)(nil)
)

// ---------------------------------------------------------------------------
// Replication & sharing DDL — CREATE / ALTER FAILOVER GROUP / REPLICATION GROUP
// / ACCOUNT / SHARE (T4.8)
// ---------------------------------------------------------------------------
//
// These account-level objects drive Snowflake's replication and data-sharing
// features. They share a distinctive option shape with the rest of this engine's
// account-level DDL — large, version-growing `KEY = <value>` vocabularies — but
// with one structural twist the COPY/STAGE option machinery does not cover: the
// replication/sharing list parameters (OBJECT_TYPES, ALLOWED_DATABASES,
// ALLOWED_SHARES, ALLOWED_INTEGRATION_TYPES, ALLOWED_ACCOUNTS, ...) are
// UNPARENTHESIZED comma lists whose elements are multi-word object-type names
// (e.g. `ACCOUNT PARAMETERS`, `RESOURCE MONITORS`, `NETWORK POLICIES`) or dotted
// account names (`org.account`). Both the official docs (truth1) and the legacy
// ANTLR grammar (truth2) agree on the no-parentheses form — the COPY-style
// `KEY = ( a, b )` reader is therefore NOT reused for them. They are captured as
// GroupOption (a `KEY = <value-list | literal>` pair). Scalar options
// (REPLICATION_SCHEDULE = '...', OPTIMIZED_REFRESH = TRUE, ERROR_INTEGRATION =
// <name>, ...) reuse the same GroupOption type. Only the structural anchors are
// modeled as dedicated fields: the object discriminator (FAILOVER vs REPLICATION
// GROUP, MANAGED ACCOUNT), the [WITH] TAG clause, the AS REPLICA OF secondary
// form, and (on ALTER) the action keyword plus its operands (ADD/REMOVE/MOVE
// name lists, the ALLOWED_* list target, the MOVE-TO group, RENAME's new name,
// the account-policy forms, etc.). Per this engine's parse-permissively
// philosophy the parser accepts the open-ended combination and the
// catalog/semantic layer validates legality (mutually-exclusive OR REPLACE / IF
// NOT EXISTS, which options are real for the chosen object, and so on).

// GroupOption is a single `KEY = <value>` parameter of a replication/sharing
// statement. The value is one of:
//   - Values: an unparenthesized comma list of items, each a verbatim,
//     uppercased multi-word token run (OBJECT_TYPES = DATABASES, ROLES;
//     ALLOWED_DATABASES = db1, db2; ALLOWED_ACCOUNTS = org.acct1, org.acct2) or
//     a single bareword/identifier value (OPTIMIZED_REFRESH = TRUE,
//     ERROR_INTEGRATION = my_int, EDITION = STANDARD, ADMIN_NAME = admin).
//   - Lit: a string or numeric literal (REPLICATION_SCHEDULE = '10 MINUTE',
//     ADMIN_PASSWORD = 'secret', EMAIL = 'a@b.com').
//
// Exactly one of Values / Lit is set. Name is the uppercased option name.
type GroupOption struct {
	Name   string   // uppercased option name (e.g. "OBJECT_TYPES", "REPLICATION_SCHEDULE")
	Values []string // unparenthesized comma-list of verbatim uppercased items; nil when Lit is set
	Lit    *Literal // a string/number literal value; nil when Values is set
	Loc    Loc
}

// Tag implements Node.
func (n *GroupOption) Tag() NodeTag { return T_GroupOption }

// ReplicationGroupKind discriminates a FAILOVER GROUP from a REPLICATION GROUP.
// The two objects share a near-identical CREATE/ALTER grammar (the docs differ
// only in a handful of options), so they share AST node types and are told apart
// by this flag.
type ReplicationGroupKind int

const (
	ReplicationGroup ReplicationGroupKind = iota // REPLICATION GROUP
	FailoverGroup                                // FAILOVER GROUP
)

// String returns the SQL keyword(s) for the kind (uppercased).
func (k ReplicationGroupKind) String() string {
	if k == FailoverGroup {
		return "FAILOVER GROUP"
	}
	return "REPLICATION GROUP"
}

// CreateReplicationGroupStmt represents the CREATE form of a FAILOVER GROUP or
// REPLICATION GROUP (Failover discriminates):
//
//	CREATE [ OR REPLACE ] { FAILOVER | REPLICATION } GROUP [ IF NOT EXISTS ] <name>
//	  OBJECT_TYPES = <object_type> [ , ... ]
//	  [ ALLOWED_DATABASES = <db_name> [ , ... ] ]
//	  [ ALLOWED_EXTERNAL_VOLUMES = <ev_name> [ , ... ] ]
//	  [ ALLOWED_SHARES = <share_name> [ , ... ] ]
//	  [ ALLOWED_INTEGRATION_TYPES = <integration_type_name> [ , ... ] ]
//	  ALLOWED_ACCOUNTS = <org>.<account> [ , ... ]
//	  [ IGNORE EDITION CHECK ]
//	  [ REPLICATION_SCHEDULE = '...' ] [ OPTIMIZED_REFRESH = { TRUE | FALSE } ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//	  [ ERROR_INTEGRATION = <integration_name> ]
//
//	-- secondary form:
//	CREATE [ OR REPLACE ] { FAILOVER | REPLICATION } GROUP [ IF NOT EXISTS ] <name>
//	  AS REPLICA OF <org>.<source_account>.<name>
//
// All configuration parameters are captured open-ended in Options, source order.
// IgnoreEditionCheck records the bare `IGNORE EDITION CHECK` flag, Tags the
// [WITH] TAG clause, and Replica the AS REPLICA OF secondary-form source (a
// dotted ObjectName); when Replica is set the primary-form fields are empty.
type CreateReplicationGroupStmt struct {
	Failover           bool
	OrReplace          bool
	OrAlter            bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists        bool
	Name               *ObjectName
	Options            []*GroupOption   // open-ended KEY = value parameters, source order
	IgnoreEditionCheck bool             // bare IGNORE EDITION CHECK flag
	Tags               []*TagAssignment // [WITH] TAG (...); nil if absent
	Replica            *ObjectName      // AS REPLICA OF <org>.<account>.<name>; nil for primary form
	Loc                Loc
}

// Tag implements Node.
func (n *CreateReplicationGroupStmt) Tag() NodeTag { return T_CreateReplicationGroupStmt }

// AlterGroupAction discriminates the action variants of ALTER FAILOVER GROUP /
// REPLICATION GROUP.
type AlterGroupAction int

const (
	AlterGroupRename   AlterGroupAction = iota // RENAME TO <new_name>
	AlterGroupSet                              // SET <options>
	AlterGroupUnset                            // UNSET <property> [ , ... ]
	AlterGroupSetTag                           // SET TAG <tag> = '<value>' [ , ... ]
	AlterGroupUnsetTag                         // UNSET TAG <tag> [ , ... ]
	AlterGroupAdd                              // ADD <name-list> TO ALLOWED_{DATABASES|SHARES|ACCOUNTS} [ IGNORE EDITION CHECK ]
	AlterGroupRemove                           // REMOVE <name-list> FROM ALLOWED_{DATABASES|SHARES|ACCOUNTS}
	AlterGroupMove                             // MOVE { DATABASES | SHARES } <name-list> TO { FAILOVER | REPLICATION } GROUP <name>
	AlterGroupRefresh                          // REFRESH
	AlterGroupPrimary                          // PRIMARY (FAILOVER GROUP only)
	AlterGroupSuspend                          // SUSPEND [ IMMEDIATE ]
	AlterGroupResume                           // RESUME
)

// AlterReplicationGroupStmt represents the ALTER form of a FAILOVER GROUP or
// REPLICATION GROUP (Failover discriminates). The legal action set depends on
// the kind (validated by the semantic layer, not here):
//
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> SET <options>
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER { FAILOVER | REPLICATION } GROUP <name> { SET | UNSET } TAG ...
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> ADD <name-list> TO ALLOWED_{DATABASES|SHARES|ACCOUNTS} [ IGNORE EDITION CHECK ]
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> MOVE { DATABASES | SHARES } <name-list> TO { FAILOVER | REPLICATION } GROUP <name>
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> REMOVE <name-list> FROM ALLOWED_{DATABASES|SHARES|ACCOUNTS}
//	ALTER { FAILOVER | REPLICATION } GROUP [ IF EXISTS ] <name> { REFRESH | PRIMARY | SUSPEND [ IMMEDIATE ] | RESUME }
type AlterReplicationGroupStmt struct {
	Failover           bool
	IfExists           bool
	Name               *ObjectName
	Action             AlterGroupAction
	NewName            *ObjectName      // RENAME TO; for AlterGroupRename
	Options            []*GroupOption   // SET <options>; for AlterGroupSet
	Tags               []*TagAssignment // SET TAG assignments; for AlterGroupSetTag
	UnsetTags          []*ObjectName    // UNSET TAG names; for AlterGroupUnsetTag
	UnsetProps         []string         // UNSET <property> names; for AlterGroupUnset
	Names              []*ObjectName    // ADD/MOVE/REMOVE name list (databases / shares / org.account)
	ListTarget         string           // ALLOWED_DATABASES / ALLOWED_SHARES / ALLOWED_ACCOUNTS (ADD/REMOVE) — uppercased
	MoveKind           string           // MOVE DATABASES | SHARES — uppercased; for AlterGroupMove
	MoveTo             *ObjectName      // MOVE ... TO { FAILOVER | REPLICATION } GROUP <name>; for AlterGroupMove
	IgnoreEditionCheck bool             // ADD ... TO ALLOWED_ACCOUNTS IGNORE EDITION CHECK
	Immediate          bool             // SUSPEND IMMEDIATE
	Loc                Loc
}

// Tag implements Node.
func (n *AlterReplicationGroupStmt) Tag() NodeTag { return T_AlterReplicationGroupStmt }

// CreateAccountStmt represents
//
//	CREATE ACCOUNT <name> ADMIN_NAME = '...' { ADMIN_PASSWORD = '...' | ADMIN_RSA_PUBLIC_KEY = '...' }
//	  [ ADMIN_USER_TYPE = ... ] [ FIRST_NAME = '...' ] [ LAST_NAME = '...' ]
//	  EMAIL = '...' [ MUST_CHANGE_PASSWORD = { TRUE | FALSE } ]
//	  EDITION = { STANDARD | ENTERPRISE | BUSINESS_CRITICAL }
//	  [ REGION_GROUP = <id> ] [ REGION = <id> ] [ COMMENT = '...' ] [ POLARIS = { TRUE | FALSE } ]
//
//	CREATE [ OR REPLACE ] MANAGED ACCOUNT <name>
//	  ADMIN_NAME = <username>, ADMIN_PASSWORD = <password>, TYPE = READER [ , COMMENT = '...' ]
//
// Managed discriminates the MANAGED ACCOUNT form. Every parameter is captured
// open-ended in Options (source order); CREATE ACCOUNT separates them with
// whitespace, CREATE MANAGED ACCOUNT with commas (the parser tolerates both).
type CreateAccountStmt struct {
	Managed   bool
	OrReplace bool // only MANAGED ACCOUNT accepts OR REPLACE per the docs
	OrAlter   bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Name      *ObjectName
	Options   []*GroupOption // ADMIN_NAME / ADMIN_PASSWORD / EMAIL / EDITION / TYPE / COMMENT / ...
	Loc       Loc
}

// Tag implements Node.
func (n *CreateAccountStmt) Tag() NodeTag { return T_CreateAccountStmt }

// AlterAccountAction discriminates the action variants of ALTER ACCOUNT.
type AlterAccountAction int

const (
	AlterAccountSet         AlterAccountAction = iota // SET <param> = <value> [ , ... ] / SET RESOURCE_MONITOR = <name>
	AlterAccountUnset                                 // UNSET <param> [ , ... ]
	AlterAccountSetTag                                // SET TAG <tag> = '<value>' [ , ... ]
	AlterAccountUnsetTag                              // UNSET TAG <tag> [ , ... ]
	AlterAccountSetPolicy                             // SET { <policy_kind> } POLICY <name> [ <scope> ] [ FORCE ]
	AlterAccountUnsetPolicy                           // UNSET { <policy_kind> } POLICY [ <scope> ]
	AlterAccountRename                                // <name> RENAME TO <new_name> [ SAVE_OLD_URL = ... ]
	AlterAccountDropURL                               // <name> DROP OLD [ ORGANIZATION ] URL
)

// AlterAccountStmt represents ALTER ACCOUNT. The current-account forms omit the
// name (Name is nil); the cross-account forms (org admins) carry a name:
//
//	ALTER ACCOUNT SET <param> = <value> [ , ... ]
//	ALTER ACCOUNT SET RESOURCE_MONITOR = <monitor_name>
//	ALTER ACCOUNT UNSET <param> [ , ... ]
//	ALTER ACCOUNT SET TAG <tag> = '<value>' [ , ... ]   |   UNSET TAG <tag> [ , ... ]
//	ALTER ACCOUNT SET { AUTHENTICATION | SESSION } POLICY <name> [ FOR ALL { PERSON | SERVICE } USERS ] [ FORCE ]
//	ALTER ACCOUNT SET { PACKAGES | PASSWORD } POLICY <name> [ FORCE ]
//	ALTER ACCOUNT SET FEATURE POLICY <name> FOR ALL APPLICATIONS [ FORCE ]
//	ALTER ACCOUNT UNSET { <policy_kind> } POLICY [ <scope> ]
//	ALTER ACCOUNT <name> SET EDITION = ... | SET IS_ORG_ADMIN = ...
//	ALTER ACCOUNT <name> RENAME TO <new_name> [ SAVE_OLD_URL = { TRUE | FALSE } ]
//	ALTER ACCOUNT <name> DROP OLD [ ORGANIZATION ] URL
type AlterAccountStmt struct {
	Name         *ObjectName // nil for the current-account forms; set for the cross-account forms
	Action       AlterAccountAction
	Options      []*GroupOption   // SET <param> = value list; for AlterAccountSet
	Tags         []*TagAssignment // SET TAG assignments; for AlterAccountSetTag
	UnsetTags    []*ObjectName    // UNSET TAG names; for AlterAccountUnsetTag
	UnsetProps   []string         // UNSET <property> names; for AlterAccountUnset
	PolicyKind   string           // "{ AUTHENTICATION | SESSION | PACKAGES | PASSWORD | FEATURE } POLICY"; for Set/UnsetPolicy — uppercased
	PolicyName   *ObjectName      // the policy name; for AlterAccountSetPolicy
	PolicyScope  string           // the trailing FOR ALL ... clause, verbatim uppercased; "" if absent
	Force        bool             // trailing FORCE; for AlterAccountSetPolicy
	NewName      *ObjectName      // RENAME TO; for AlterAccountRename
	SaveOldURL   *bool            // RENAME ... SAVE_OLD_URL = { TRUE | FALSE }; nil if absent
	Organization bool             // DROP OLD ORGANIZATION URL (vs DROP OLD URL); for AlterAccountDropURL
	Loc          Loc
}

// Tag implements Node.
func (n *AlterAccountStmt) Tag() NodeTag { return T_AlterAccountStmt }

// CreateShareStmt represents
//
//	CREATE [ OR REPLACE ] SHARE [ IF NOT EXISTS ] <name> [ COMMENT = '<string_literal>' ]
//
// COMMENT (the only documented parameter) is captured open-ended in Options.
type CreateShareStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Options     []*GroupOption // COMMENT = '...'; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateShareStmt) Tag() NodeTag { return T_CreateShareStmt }

// AlterShareAction discriminates the action variants of ALTER SHARE.
type AlterShareAction int

const (
	AlterShareAdd      AlterShareAction = iota // ADD ACCOUNTS = <a> [ , ... ] [ SHARE_RESTRICTIONS = ... ]
	AlterShareRemove                           // REMOVE ACCOUNTS = <a> [ , ... ]
	AlterShareSet                              // SET [ ACCOUNTS = <a> [ , ... ] ] [ COMMENT = '...' ]
	AlterShareSetTag                           // SET TAG <tag> = '<value>' [ , ... ]
	AlterShareUnsetTag                         // UNSET TAG <tag> [ , ... ]
	AlterShareUnset                            // UNSET COMMENT
)

// AlterShareStmt represents
//
//	ALTER SHARE [ IF EXISTS ] <name> { ADD | REMOVE } ACCOUNTS = <a> [ , ... ] [ SHARE_RESTRICTIONS = { TRUE | FALSE } ]
//	ALTER SHARE [ IF EXISTS ] <name> SET [ ACCOUNTS = <a> [ , ... ] ] [ COMMENT = '...' ]
//	ALTER SHARE [ IF EXISTS ] <name> SET TAG <tag> = '<value>' [ , ... ]
//	ALTER SHARE <name> UNSET TAG <tag> [ , ... ]
//	ALTER SHARE [ IF EXISTS ] <name> UNSET COMMENT
//
// Accounts holds the consumer-account list (ADD/REMOVE/SET ACCOUNTS = ...);
// ShareRestrictions records the optional SHARE_RESTRICTIONS flag; Options holds
// the SET COMMENT param; UnsetProps holds the UNSET property names.
type AlterShareStmt struct {
	IfExists          bool
	Name              *ObjectName
	Action            AlterShareAction
	Accounts          []*ObjectName    // ACCOUNTS = <a> [ , ... ]; for Add/Remove/Set
	ShareRestrictions *bool            // SHARE_RESTRICTIONS = { TRUE | FALSE }; nil if absent
	Options           []*GroupOption   // SET COMMENT = '...'; nil if absent
	Tags              []*TagAssignment // SET TAG assignments; for AlterShareSetTag
	UnsetTags         []*ObjectName    // UNSET TAG names; for AlterShareUnsetTag
	UnsetProps        []string         // UNSET <property> names (e.g. COMMENT); for AlterShareUnset
	Loc               Loc
}

// Tag implements Node.
func (n *AlterShareStmt) Tag() NodeTag { return T_AlterShareStmt }

// Compile-time assertions for replication & sharing DDL nodes.
var (
	_ Node = (*GroupOption)(nil)
	_ Node = (*CreateReplicationGroupStmt)(nil)
	_ Node = (*AlterReplicationGroupStmt)(nil)
	_ Node = (*CreateAccountStmt)(nil)
	_ Node = (*AlterAccountStmt)(nil)
	_ Node = (*CreateShareStmt)(nil)
	_ Node = (*AlterShareStmt)(nil)
)

// ---------------------------------------------------------------------------
// Tag / semantic-view / dataset DDL — CREATE / ALTER (T4.9)
// ---------------------------------------------------------------------------
//
// Three governance / semantic-layer objects:
//
//   - TAG: a named label whose CREATE carries an ALLOWED_VALUES string list plus
//     a small, version-growing option vocabulary (PROPAGATE, ON_CONFLICT,
//     COMMENT). The legacy ANTLR grammar's create_tag rule modeled only
//     `tag_allowed_values? comment_clause?` and so lacks PROPAGATE / ON_CONFLICT
//     (both present in the official docs and the create-tag corpus). To avoid
//     re-staling, every clause after ALLOWED_VALUES is captured open-ended as a
//     CopyOption; only ALLOWED_VALUES — the one structural anchor the docs pin
//     to first position — is lifted into a dedicated field.
//   - SEMANTIC VIEW: a logical model over tables. Its TABLES / RELATIONSHIPS /
//     FACTS / DIMENSIONS / METRICS / AI_VERIFIED_QUERIES sections are each a
//     comma-separated definition list inside a parenthesized group whose inner
//     grammar is large and version-growing; rather than model each inner
//     definition, each section's body is captured as a balanced raw group
//     (SemanticViewSection). The trailing scalar clauses (COMMENT,
//     AI_SQL_GENERATION, AI_QUESTION_CATEGORIZATION) are open-ended options, and
//     the [WITH] TAG / COPY GRANTS clauses are the remaining structural anchors.
//   - DATASET: a newer ML object whose CREATE is, per both the docs and the
//     legacy grammar, just a name (no options).
//
// This mirrors the merged STAGE (T4.1) / FILE FORMAT (T4.2) / integration (T4.7)
// open-ended philosophy: the catalog/semantic layer, not the parser, validates
// that an option is real and legal.

// CreateTagStmt represents
//
//	CREATE [ OR REPLACE ] TAG [ IF NOT EXISTS ] <name>
//	  [ ALLOWED_VALUES '<v1>' [ , '<v2>' ... ] ]
//	  [ PROPAGATE = { ON_DEPENDENCY_AND_DATA_MOVEMENT | ON_DEPENDENCY | ON_DATA_MOVEMENT }
//	    [ ON_CONFLICT = { '<string>' | ALLOWED_VALUES_SEQUENCE } ] ]
//	  [ COMMENT = '<string_literal>' ]
//
// AllowedValues holds the ALLOWED_VALUES string list (nil when the clause is
// absent). Options captures the trailing PROPAGATE / ON_CONFLICT / COMMENT
// clauses open-ended, in source order.
type CreateTagStmt struct {
	OrReplace     bool
	OrAlter       bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists   bool
	Name          *ObjectName
	AllowedValues []string      // ALLOWED_VALUES 'v1', 'v2', ...; nil if absent
	Options       []*CopyOption // PROPAGATE / ON_CONFLICT / COMMENT; source order
	Loc           Loc
}

// Tag implements Node.
func (n *CreateTagStmt) Tag() NodeTag { return T_CreateTagStmt }

// AlterTagAction discriminates the action variants of ALTER TAG.
type AlterTagAction int

const (
	AlterTagRename             AlterTagAction = iota // RENAME TO <new_name>
	AlterTagAddAllowedValues                         // ADD ALLOWED_VALUES '<v>' [ , ... ]
	AlterTagDropAllowedValues                        // DROP ALLOWED_VALUES '<v>' [ , ... ]
	AlterTagSet                                      // SET [ ALLOWED_VALUES ... ] [ PROPAGATE ... ] [ COMMENT ... ]
	AlterTagUnset                                    // UNSET { ALLOWED_VALUES | PROPAGATE | ON_CONFLICT | COMMENT | DCM PROJECT }
	AlterTagSetMaskingPolicy                         // SET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ] [ FORCE ]
	AlterTagUnsetMaskingPolicy                       // UNSET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ]
)

// AlterTagStmt represents ALTER TAG [ IF EXISTS ] <name> <action>.
//
//	RENAME TO <new_name>
//	{ ADD | DROP } ALLOWED_VALUES '<v>' [ , ... ]
//	SET [ ALLOWED_VALUES '<v>' [ , ... ] ] [ PROPAGATE = ... [ ON_CONFLICT = ... ] ] [ COMMENT = '...' ]
//	UNSET { ALLOWED_VALUES | PROPAGATE | ON_CONFLICT | COMMENT | DCM PROJECT }
//	SET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ] [ FORCE ]
//	UNSET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ]
type AlterTagStmt struct {
	IfExists        bool
	Name            *ObjectName
	Action          AlterTagAction
	NewName         *ObjectName   // RENAME TO target; for AlterTagRename
	AllowedValues   []string      // ADD/DROP/SET ALLOWED_VALUES list
	Options         []*CopyOption // SET <options> (ALLOWED_VALUES lives in AllowedValues; PROPAGATE / ON_CONFLICT / COMMENT here)
	UnsetProps      []string      // UNSET <property> names (ALLOWED_VALUES / PROPAGATE / ON_CONFLICT / COMMENT / DCM PROJECT)
	MaskingPolicies []*ObjectName // SET/UNSET MASKING POLICY names
	Force           bool          // SET MASKING POLICY ... FORCE
	Loc             Loc
}

// Tag implements Node.
func (n *AlterTagStmt) Tag() NodeTag { return T_AlterTagStmt }

// SemanticViewSection is one section of a CREATE SEMANTIC VIEW body — a keyword
// (TABLES / RELATIONSHIPS / FACTS / DIMENSIONS / METRICS / AI_VERIFIED_QUERIES)
// followed by a parenthesized, comma-separated definition list. The inner
// grammar of each section is large and version-growing, so the body is captured
// as the verbatim source text between the section's outer parentheses (the
// parentheses themselves are not included in Body) rather than fully modeled.
type SemanticViewSection struct {
	Keyword string // uppercased section keyword: TABLES, RELATIONSHIPS, FACTS, DIMENSIONS, METRICS, AI_VERIFIED_QUERIES
	Body    string // verbatim source text inside the section's ( ... ); excludes the parens
	Loc     Loc    // spans the keyword through the closing ')'
}

// Tag implements Node.
func (n *SemanticViewSection) Tag() NodeTag { return T_SemanticViewSection }

// CreateSemanticViewStmt represents
//
//	CREATE [ OR REPLACE ] SEMANTIC VIEW [ IF NOT EXISTS ] <name>
//	  TABLES ( ... )
//	  [ RELATIONSHIPS ( ... ) ]
//	  [ FACTS ( ... ) ]
//	  [ DIMENSIONS ( ... ) ]
//	  [ METRICS ( ... ) ]
//	  [ COMMENT = '<string>' ]
//	  [ AI_SQL_GENERATION '<instructions>' ]
//	  [ AI_QUESTION_CATEGORIZATION '<instructions>' ]
//	  [ AI_VERIFIED_QUERIES ( ... ) ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//	  [ COPY GRANTS ]
//
// Sections holds the parenthesized definition-list sections (TABLES first, per
// the docs, then the optional RELATIONSHIPS / FACTS / DIMENSIONS / METRICS /
// AI_VERIFIED_QUERIES), each as a balanced raw group. Options captures the
// scalar trailing clauses (COMMENT, AI_SQL_GENERATION, AI_QUESTION_CATEGORIZATION)
// open-ended. Tags holds a [WITH] TAG clause; CopyGrants records COPY GRANTS.
type CreateSemanticViewStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Sections    []*SemanticViewSection // TABLES / RELATIONSHIPS / FACTS / DIMENSIONS / METRICS / AI_VERIFIED_QUERIES, source order
	Options     []*CopyOption          // COMMENT / AI_SQL_GENERATION / AI_QUESTION_CATEGORIZATION; source order
	Tags        []*TagAssignment       // [WITH] TAG (...); nil if absent
	CopyGrants  bool
	Loc         Loc
}

// Tag implements Node.
func (n *CreateSemanticViewStmt) Tag() NodeTag { return T_CreateSemanticViewStmt }

// AlterSemanticViewAction discriminates the action variants of ALTER SEMANTIC
// VIEW.
type AlterSemanticViewAction int

const (
	AlterSemanticViewRename   AlterSemanticViewAction = iota // RENAME TO <new_name>
	AlterSemanticViewSet                                     // SET COMMENT = '...'
	AlterSemanticViewUnset                                   // UNSET COMMENT
	AlterSemanticViewSetTag                                  // SET TAG <tag> = '<value>' [ , ... ]
	AlterSemanticViewUnsetTag                                // UNSET TAG <tag> [ , ... ]
)

// AlterSemanticViewStmt represents ALTER SEMANTIC VIEW [ IF EXISTS ] <name>
// <action>.
//
//	RENAME TO <new_name>
//	SET COMMENT = '<string_literal>'
//	UNSET COMMENT
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
type AlterSemanticViewStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterSemanticViewAction
	NewName    *ObjectName      // RENAME TO target; for AlterSemanticViewRename
	Options    []*CopyOption    // SET COMMENT; for AlterSemanticViewSet
	Tags       []*TagAssignment // SET TAG assignments; for AlterSemanticViewSetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterSemanticViewUnsetTag
	UnsetProps []string         // UNSET <property> names (COMMENT); for AlterSemanticViewUnset
	Loc        Loc
}

// Tag implements Node.
func (n *AlterSemanticViewStmt) Tag() NodeTag { return T_AlterSemanticViewStmt }

// CreateDatasetStmt represents CREATE [ OR REPLACE ] DATASET [ IF NOT EXISTS ]
// <name>. DATASET is a newer ML object whose CREATE carries no options in either
// the docs or the legacy grammar — only the name. (The official docs render
// IF NOT EXISTS before DATASET; the legacy grammar and every other CREATE object
// in this engine place it after the object keyword. The post-keyword spelling is
// accepted here for consistency; see the divergence note in semantic_view.go.)
type CreateDatasetStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Loc         Loc
}

// Tag implements Node.
func (n *CreateDatasetStmt) Tag() NodeTag { return T_CreateDatasetStmt }

// Compile-time assertions for tag / semantic-view / dataset DDL nodes.
var (
	_ Node = (*CreateTagStmt)(nil)
	_ Node = (*AlterTagStmt)(nil)
	_ Node = (*SemanticViewSection)(nil)
	_ Node = (*CreateSemanticViewStmt)(nil)
	_ Node = (*AlterSemanticViewStmt)(nil)
	_ Node = (*CreateDatasetStmt)(nil)
)

// ---------------------------------------------------------------------------
// Data-pipeline DDL — CREATE / ALTER PIPE / STREAM / TASK / ALERT (T4.3)
// ---------------------------------------------------------------------------
//
// These four objects are Snowflake's data-pipeline primitives. Like STAGE /
// FILE FORMAT / COPY before them, each carries a large, version-growing
// configuration vocabulary (PIPE: AUTO_INGEST / INTEGRATION / ERROR_INTEGRATION
// / AWS_SNS_TOPIC / ...; TASK: WAREHOUSE / SCHEDULE / CONFIG / OVERLAP_POLICY /
// ERROR_INTEGRATION / SUCCESS_INTEGRATION / LOG_LEVEL / TARGET_COMPLETION_INTERVAL
// / SERVERLESS_TASK_* / arbitrary session parameters / ...; ALERT: WAREHOUSE /
// SCHEDULE / COMMENT / CONFIG / RUNBOOK / SUSPEND_ALERT_AFTER_NUM_FAILURES /
// ...). The legacy ANTLR grammar enumerates a stale subset of each (its
// create_task lists only task_overlap?/task_timeout?/task_suspend_after_failure;
// its create_pipe lacks SUCCESS-side options; etc.), so — matching the merged
// STAGE (T4.1) / FILE FORMAT (T4.2) / COPY (T5.2) approach — every such config
// parameter is captured as an open-ended `KEY = <value>` pair (CopyOption). Only
// the structural anchors are modeled as dedicated fields: the embedded statement
// bodies (PIPE's AS COPY, TASK's AS <sql>, ALERT's IF(EXISTS(...)) THEN <action>),
// STREAM's ON <source> + AT/BEFORE time-travel, TASK's AFTER predecessor list and
// WHEN predicate, and the [WITH] TAG / COPY GRANTS clauses.
//
// Embedded bodies are handled by REUSING the existing statement parsers rather
// than reinventing them: PIPE's body is parsed by parseCopyStmt, ALERT's
// condition by parseQueryExpr, and TASK's body / ALERT's action by the top-level
// parseStmt. Because a TASK body may be a Snowflake Scripting block
// (BEGIN…END / DECLARE…BEGIN…END) or any statement parseStmt does not yet
// support (e.g. CALL), those bodies additionally retain the verbatim source text
// (BodyRaw / ActionRaw) and tolerate a nil parsed node — the statement is never
// rejected on account of an unsupported body. (F3's Split already keeps a
// `... AS BEGIN … END;` body as one segment, so the body text runs to the
// segment end.)

// PipeOption / StreamOption / TaskOption / AlertOption are all CopyOption — the
// shared open-ended `KEY = <value>` option type. The aliases below exist only to
// document intent at the field declarations; no separate type is introduced.

// CreatePipeStmt represents
//
//	CREATE [ OR REPLACE ] PIPE [ IF NOT EXISTS ] <name>
//	  [ AUTO_INGEST = { TRUE | FALSE } ]
//	  [ ERROR_INTEGRATION = <integration_name> ]
//	  [ AWS_SNS_TOPIC = '<string>' ]
//	  [ INTEGRATION = '<string>' ]
//	  [ COMMENT = '<string_literal>' ]
//	  AS <copy_statement>
//
// Options holds every config parameter that precedes AS (AUTO_INGEST /
// ERROR_INTEGRATION / AWS_SNS_TOPIC / INTEGRATION / COMMENT / ...), each a
// CopyOption, preserving source order. Copy holds the parsed `AS COPY INTO`
// body (a *CopyIntoTableStmt); the AS body is always a COPY INTO, which the
// parser fully supports, so Copy is non-nil for a well-formed statement. The
// body may be wrapped in parentheses in the source (`AS (COPY INTO …)`); the
// parens are not retained.
type CreatePipeStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Options     []*CopyOption // AUTO_INGEST / ERROR_INTEGRATION / AWS_SNS_TOPIC / INTEGRATION / COMMENT / ...
	Copy        Node          // the AS <copy_into_table> body (parseCopyStmt result)
	Loc         Loc
}

// Tag implements Node.
func (n *CreatePipeStmt) Tag() NodeTag { return T_CreatePipeStmt }

// AlterPipeAction discriminates the action variants of ALTER PIPE.
type AlterPipeAction int

const (
	AlterPipeSet      AlterPipeAction = iota // SET <options>
	AlterPipeUnset                           // UNSET <property> [, ...]
	AlterPipeSetTag                          // SET TAG <tag> = '<value>' [, ...]
	AlterPipeUnsetTag                        // UNSET TAG <tag> [, ...]
	AlterPipeRefresh                         // REFRESH [ PREFIX = '<p>' ] [ MODIFIED_AFTER = '<ts>' ]
)

// AlterPipeStmt represents ALTER PIPE [ IF EXISTS ] <name> <action>.
//
//	SET <options>                                        -- open-ended KEY = value
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]                            -- e.g. PIPE_EXECUTION_PAUSED, COMMENT
//	REFRESH [ PREFIX = '<path>' ] [ MODIFIED_AFTER = '<timestamp>' ]
type AlterPipeStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterPipeAction
	Options    []*CopyOption    // SET <options> / REFRESH PREFIX|MODIFIED_AFTER params
	Tags       []*TagAssignment // SET TAG assignments; for AlterPipeSetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterPipeUnsetTag
	UnsetProps []string         // UNSET <property> names; for AlterPipeUnset
	Loc        Loc
}

// Tag implements Node.
func (n *AlterPipeStmt) Tag() NodeTag { return T_AlterPipeStmt }

// StreamSourceKind classifies the object a stream is created ON.
type StreamSourceKind int

const (
	StreamOnTable         StreamSourceKind = iota // ON TABLE <name>
	StreamOnView                                  // ON VIEW <name>
	StreamOnStage                                 // ON STAGE <name>
	StreamOnExternalTable                         // ON EXTERNAL TABLE <name>
)

// StreamTimeTravel holds a stream's { AT | BEFORE } ( <key> => <value> ) clause.
// Per the docs the key is one of TIMESTAMP / OFFSET / STATEMENT / STREAM and the
// value is a general expression (the corpus shows function calls like
// TO_TIMESTAMP(...), arithmetic like -60*5, and string literals). Both the AT/
// BEFORE selector and the key are captured verbatim (uppercased); Value holds the
// parsed value expression.
type StreamTimeTravel struct {
	AtBefore string // "AT" or "BEFORE"
	Key      string // "TIMESTAMP" / "OFFSET" / "STATEMENT" / "STREAM" (uppercased)
	Value    Node   // the => value expression
	Loc      Loc
}

// CreateStreamStmt represents
//
//	CREATE [ OR REPLACE ] STREAM [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//	  [ COPY GRANTS ]
//	  ON { TABLE | VIEW | STAGE | EXTERNAL TABLE } <object_name>
//	  [ { AT | BEFORE } ( { TIMESTAMP => … | OFFSET => … | STATEMENT => … | STREAM => … } ) ]
//	  [ APPEND_ONLY = { TRUE | FALSE } ]      -- TABLE/VIEW
//	  [ INSERT_ONLY = TRUE ]                  -- EXTERNAL TABLE
//	  [ SHOW_INITIAL_ROWS = { TRUE | FALSE } ]
//	  [ COMMENT = '<string_literal>' ]
//
// (plus a CLONE variant: CREATE [OR REPLACE] STREAM <name> CLONE <source> [COPY GRANTS].)
//
// SourceKind + Source model the mandatory ON clause. TimeTravel holds the
// optional AT/BEFORE clause. Options holds every trailing config parameter
// (APPEND_ONLY / INSERT_ONLY / SHOW_INITIAL_ROWS / COMMENT / ...) open-ended,
// preserving source order — they are not type-validated against SourceKind here
// (the catalog/semantic layer does that). Tags holds the [WITH] TAG assignments;
// CopyGrants records the COPY GRANTS flag. For the CLONE form, Clone is set and
// SourceKind/Source/Options are zero.
type CreateStreamStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Tags        []*TagAssignment  // [WITH] TAG (...); nil if absent
	CopyGrants  bool              // COPY GRANTS
	SourceKind  StreamSourceKind  // ON { TABLE | VIEW | STAGE | EXTERNAL TABLE }
	Source      *ObjectName       // the ON <object_name>
	TimeTravel  *StreamTimeTravel // { AT | BEFORE } (...); nil if absent
	Options     []*CopyOption     // APPEND_ONLY / INSERT_ONLY / SHOW_INITIAL_ROWS / COMMENT / ...
	Clone       *ObjectName       // CLONE <source>; non-nil only for the CLONE form
	Loc         Loc
}

// Tag implements Node.
func (n *CreateStreamStmt) Tag() NodeTag { return T_CreateStreamStmt }

// AlterStreamAction discriminates the action variants of ALTER STREAM.
type AlterStreamAction int

const (
	AlterStreamSet      AlterStreamAction = iota // SET <options>
	AlterStreamUnset                             // UNSET <property> [, ...] (e.g. COMMENT)
	AlterStreamSetTag                            // SET TAG <tag> = '<value>' [, ...]
	AlterStreamUnsetTag                          // UNSET TAG <tag> [, ...]
)

// AlterStreamStmt represents ALTER STREAM [ IF EXISTS ] <name> <action>.
//
//	SET <options>                              -- open-ended KEY = value (COMMENT = '…')
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]                  -- e.g. COMMENT
type AlterStreamStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterStreamAction
	Options    []*CopyOption    // SET <options>; for AlterStreamSet
	Tags       []*TagAssignment // SET TAG assignments; for AlterStreamSetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterStreamUnsetTag
	UnsetProps []string         // UNSET <property> names; for AlterStreamUnset
	Loc        Loc
}

// Tag implements Node.
func (n *AlterStreamStmt) Tag() NodeTag { return T_AlterStreamStmt }

// CreateTaskStmt represents
//
//	CREATE [ OR REPLACE ] TASK [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( … ) ]
//	  [ config parameters: WAREHOUSE | USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE |
//	    SCHEDULE | CONFIG | OVERLAP_POLICY | <session_parameter> | USER_TASK_TIMEOUT_MS |
//	    SUSPEND_TASK_AFTER_NUM_FAILURES | ERROR_INTEGRATION | SUCCESS_INTEGRATION |
//	    LOG_LEVEL | COMMENT | FINALIZE | TASK_AUTO_RETRY_ATTEMPTS | … ]
//	  [ AFTER <predecessor> [ , <predecessor> , ... ] ]
//	  [ WHEN <boolean_expr> ]
//	  AS <sql>
//
// Options holds every config parameter that precedes AFTER/WHEN/AS, open-ended,
// preserving source order. After holds the AFTER predecessor list (object names;
// the corpus uses bare identifiers while the docs spell them as strings — object
// names cover both). When holds the optional WHEN predicate expression. Body /
// BodyRaw hold the AS body: Body is the parsed statement (a SELECT / INSERT /
// etc.) when parseStmt supports it, and nil for a Snowflake Scripting block
// (BEGIN…END / DECLARE…BEGIN…END) or any other unsupported body; BodyRaw always
// holds the verbatim body source text so consumers retain it regardless. Tags
// holds the [WITH] TAG assignments.
//
// CREATE TASK also has a CLONE form (CREATE [OR REPLACE] TASK <name> CLONE
// <source_task>), recorded via Clone with Options/After/When/Body all zero.
type CreateTaskStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Options     []*CopyOption    // WAREHOUSE / SCHEDULE / CONFIG / session params / ...
	After       []*ObjectName    // AFTER <predecessor> [, ...]; nil if absent
	When        Node             // WHEN <boolean_expr>; nil if absent
	Body        Node             // AS <sql> parsed statement; nil for an unsupported/scripting body
	BodyRaw     string           // verbatim AS body source text (always set)
	Clone       *ObjectName      // CLONE <source_task>; non-nil only for the CLONE form
	Loc         Loc
}

// Tag implements Node.
func (n *CreateTaskStmt) Tag() NodeTag { return T_CreateTaskStmt }

// AlterTaskAction discriminates the action variants of ALTER TASK.
type AlterTaskAction int

const (
	AlterTaskResume      AlterTaskAction = iota // RESUME
	AlterTaskSuspend                            // SUSPEND
	AlterTaskAddAfter                           // ADD AFTER <task> [, ...]
	AlterTaskRemoveAfter                        // REMOVE AFTER <task> [, ...]
	AlterTaskSet                                // SET <options>
	AlterTaskUnset                              // UNSET <property> [, ...]
	AlterTaskSetTag                             // SET TAG <tag> = '<value>' [, ...]
	AlterTaskUnsetTag                           // UNSET TAG <tag> [, ...]
	AlterTaskModifyAs                           // MODIFY AS <sql>
	AlterTaskModifyWhen                         // MODIFY WHEN <boolean_expr>
)

// AlterTaskStmt represents ALTER TASK [ IF EXISTS ] <name> <action>.
//
//	RESUME | SUSPEND
//	{ ADD | REMOVE } AFTER <task> [ , ... ]
//	SET <options>                              -- open-ended KEY = value
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	UNSET <property> [ , ... ]
//	MODIFY AS <sql>
//	MODIFY WHEN <boolean_expr>
type AlterTaskStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterTaskAction
	After      []*ObjectName    // {ADD|REMOVE} AFTER list; for AlterTaskAddAfter/RemoveAfter
	Options    []*CopyOption    // SET <options>; for AlterTaskSet
	Tags       []*TagAssignment // SET TAG assignments; for AlterTaskSetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterTaskUnsetTag
	UnsetProps []string         // UNSET <property> names; for AlterTaskUnset
	When       Node             // MODIFY WHEN <expr>; for AlterTaskModifyWhen
	Body       Node             // MODIFY AS <sql> parsed statement; nil for an unsupported/scripting body
	BodyRaw    string           // verbatim MODIFY AS body source text; set for AlterTaskModifyAs
	Loc        Loc
}

// Tag implements Node.
func (n *AlterTaskStmt) Tag() NodeTag { return T_AlterTaskStmt }

// CreateAlertStmt represents
//
//	CREATE [ OR REPLACE ] ALERT [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( … ) ]
//	  [ WAREHOUSE = <warehouse_name> ]        -- optional for serverless alerts
//	  SCHEDULE = '<schedule>'
//	  [ COMMENT = '…' ] [ CONFIG = '…' ] [ RUNBOOK = '…' ]
//	  [ SUSPEND_ALERT_AFTER_NUM_FAILURES = <n> ]
//	  IF ( EXISTS ( <condition> ) )
//	  THEN <action>
//
// Options holds every config parameter that precedes IF (WAREHOUSE / SCHEDULE /
// COMMENT / CONFIG / RUNBOOK / SUSPEND_ALERT_AFTER_NUM_FAILURES / ...),
// open-ended, preserving source order — WAREHOUSE and SCHEDULE are NOT lifted out
// (WAREHOUSE is optional for serverless alerts and SCHEDULE's value vocabulary
// grows), matching the open-ended philosophy. Condition holds the query inside
// IF(EXISTS(<condition>)) (a SELECT / SHOW / CALL; parsed via parseQueryExpr for
// the SELECT case). Action / ActionRaw hold the THEN body: Action is the parsed
// statement when parseStmt supports it, nil otherwise; ActionRaw always holds the
// verbatim action source text.
type CreateAlertStmt struct {
	OrReplace    bool
	OrAlter      bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists  bool
	Name         *ObjectName
	Tags         []*TagAssignment // [WITH] TAG (...); nil if absent
	Options      []*CopyOption    // WAREHOUSE / SCHEDULE / COMMENT / CONFIG / RUNBOOK / ...
	Condition    Node             // IF(EXISTS(<condition>)) query; nil only if unparsed
	ConditionRaw string           // verbatim condition source text (always set)
	Action       Node             // THEN <action> parsed statement; nil for an unsupported body
	ActionRaw    string           // verbatim THEN action source text (always set)
	Loc          Loc
}

// Tag implements Node.
func (n *CreateAlertStmt) Tag() NodeTag { return T_CreateAlertStmt }

// AlterAlertAction discriminates the action variants of ALTER ALERT.
type AlterAlertAction int

const (
	AlterAlertResume          AlterAlertAction = iota // RESUME
	AlterAlertSuspend                                 // SUSPEND
	AlterAlertSet                                     // SET <options>
	AlterAlertUnset                                   // UNSET <property> [, ...]
	AlterAlertModifyCondition                         // MODIFY CONDITION EXISTS (<query>)
	AlterAlertModifyAction                            // MODIFY ACTION <action>
)

// AlterAlertStmt represents ALTER ALERT [ IF EXISTS ] <name> <action>.
//
//	RESUME | SUSPEND
//	SET <options>                              -- WAREHOUSE = … | SCHEDULE = … | COMMENT = …
//	UNSET <property> [ , ... ]                  -- WAREHOUSE | SCHEDULE | COMMENT
//	MODIFY CONDITION EXISTS ( <query> )
//	MODIFY ACTION <action>
type AlterAlertStmt struct {
	IfExists     bool
	Name         *ObjectName
	Action       AlterAlertAction
	Options      []*CopyOption // SET <options>; for AlterAlertSet
	UnsetProps   []string      // UNSET <property> names; for AlterAlertUnset
	Condition    Node          // MODIFY CONDITION EXISTS (<query>); for AlterAlertModifyCondition
	ConditionRaw string        // verbatim condition source text; set for AlterAlertModifyCondition
	ActionBody   Node          // MODIFY ACTION <action> parsed statement; nil for an unsupported body
	ActionRaw    string        // verbatim MODIFY ACTION source text; set for AlterAlertModifyAction
	Loc          Loc
}

// Tag implements Node.
func (n *AlterAlertStmt) Tag() NodeTag { return T_AlterAlertStmt }

// Compile-time assertions for data-pipeline DDL nodes.
var (
	_ Node = (*CreatePipeStmt)(nil)
	_ Node = (*AlterPipeStmt)(nil)
	_ Node = (*CreateStreamStmt)(nil)
	_ Node = (*AlterStreamStmt)(nil)
	_ Node = (*CreateTaskStmt)(nil)
	_ Node = (*AlterTaskStmt)(nil)
	_ Node = (*CreateAlertStmt)(nil)
	_ Node = (*AlterAlertStmt)(nil)
)

// ---------------------------------------------------------------------------
// Access-control DDL — CREATE / ALTER ROLE / USER / { MASKING | ROW ACCESS |
// SESSION | PASSWORD | NETWORK | AUTHENTICATION } POLICY (T4.6)
// ---------------------------------------------------------------------------
//
// Roles, users, and the six policy objects share the open-ended `KEY = <value>`
// option-bag approach already merged for STAGE (T4.1) / FILE FORMAT (T4.2) /
// COPY (T5.2) / pipeline (T4.3): every option a CREATE/ALTER USER or
// SESSION/PASSWORD/NETWORK/AUTHENTICATION POLICY carries is a version-growing
// vocabulary the catalog/semantic layer — not the parser — validates, so each
// option is captured as a CopyOption rather than enumerated against the
// already-stale legacy ANTLR rules (which lack, e.g., the network policy's
// ALLOWED_NETWORK_RULE_LIST, the user's TYPE / RSA_PUBLIC_KEY_2, and the whole
// AUTHENTICATION POLICY object). The two structurally distinct objects —
// MASKING POLICY and ROW ACCESS POLICY — carry a typed argument list, a RETURNS
// type, and an expression body after `->`; those are modeled as dedicated fields
// (Args / Returns / Body), with the body reusing the expression parser so it
// participates in walks / query-span / analysis. COMMENT and other trailing
// `KEY = value` options after the body are captured open-ended in Options.

// PolicyKind discriminates the six access-control policy object types.
type PolicyKind int

const (
	PolicyMasking        PolicyKind = iota // MASKING POLICY
	PolicyRowAccess                        // ROW ACCESS POLICY
	PolicySession                          // SESSION POLICY
	PolicyPassword                         // PASSWORD POLICY
	PolicyNetwork                          // NETWORK POLICY
	PolicyAuthentication                   // AUTHENTICATION POLICY
)

// String returns the SQL spelling of a PolicyKind (without the trailing POLICY).
func (k PolicyKind) String() string {
	switch k {
	case PolicyMasking:
		return "MASKING"
	case PolicyRowAccess:
		return "ROW ACCESS"
	case PolicySession:
		return "SESSION"
	case PolicyPassword:
		return "PASSWORD"
	case PolicyNetwork:
		return "NETWORK"
	case PolicyAuthentication:
		return "AUTHENTICATION"
	default:
		return "UNKNOWN"
	}
}

// PolicyArg is one `<arg_name> <arg_type>` entry in a MASKING / ROW ACCESS
// policy's signature: AS ( <arg_name> <arg_type> [ , ... ] ).
type PolicyArg struct {
	Name     Ident
	DataType *TypeName
	Loc      Loc
}

// Tag implements Node.
func (n *PolicyArg) Tag() NodeTag { return T_PolicyArg }

// CreateRoleStmt represents
//
//	CREATE [ OR REPLACE ] [ DATABASE ] ROLE [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// Database is true for the CREATE DATABASE ROLE form (its <name> may be
// db-qualified, db_name.role_name). Account roles take an unqualified name.
type CreateRoleStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Database    bool // CREATE DATABASE ROLE
	IfNotExists bool
	Name        *ObjectName
	Comment     *string          // COMMENT = '...'; nil if absent
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateRoleStmt) Tag() NodeTag { return T_CreateRoleStmt }

// CreateUserStmt represents
//
//	CREATE [ OR REPLACE ] USER [ IF NOT EXISTS ] <name>
//	  [ objectProperties ] [ objectParams ] [ sessionParams ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// Every property/parameter (PASSWORD / LOGIN_NAME / DISPLAY_NAME / DEFAULT_ROLE
// / DEFAULT_WAREHOUSE / RSA_PUBLIC_KEY / DEFAULT_SECONDARY_ROLES / TYPE /
// MUST_CHANGE_PASSWORD / network-policy / session params / ...) is captured
// open-ended as a CopyOption, preserving source order.
type CreateUserStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Options     []*CopyOption    // open-ended objectProperties / objectParams / sessionParams
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateUserStmt) Tag() NodeTag { return T_CreateUserStmt }

// CreatePolicyStmt represents the CREATE form of the six access-control policy
// objects (one node, discriminated by Kind):
//
//	CREATE [ OR REPLACE ] MASKING POLICY [ IF NOT EXISTS ] <name>
//	  AS ( <arg> <type> [ , ... ] ) RETURNS <type> -> <body>
//	  [ COMMENT = '...' ] [ EXEMPT_OTHER_POLICIES = { TRUE | FALSE } ]
//
//	CREATE [ OR REPLACE ] ROW ACCESS POLICY [ IF NOT EXISTS ] <name>
//	  AS ( <arg> <type> [ , ... ] ) RETURNS BOOLEAN -> <body> [ COMMENT = '...' ]
//
//	CREATE [ OR REPLACE ] { SESSION | PASSWORD | NETWORK | AUTHENTICATION } POLICY
//	  [ IF NOT EXISTS ] <name> [ <option> ... ]
//
// For MASKING / ROW ACCESS policies Args holds the typed signature, Returns the
// RETURNS type, and Body the `-> <expr>` body (an expression node). For the
// SESSION / PASSWORD / NETWORK / AUTHENTICATION policies Args/Returns/Body are
// nil and Options holds the open-ended option bag. Trailing options after a
// MASKING / ROW ACCESS body (COMMENT / EXEMPT_OTHER_POLICIES) are also captured
// in Options.
type CreatePolicyStmt struct {
	Kind        PolicyKind
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Args        []*PolicyArg  // MASKING / ROW ACCESS signature; nil for option-bag policies
	Returns     *TypeName     // RETURNS <type>; nil for option-bag policies
	Body        Node          // `-> <expr>` body; nil for option-bag policies
	Options     []*CopyOption // option bag + trailing COMMENT / EXEMPT_OTHER_POLICIES
	Loc         Loc
}

// Tag implements Node.
func (n *CreatePolicyStmt) Tag() NodeTag { return T_CreatePolicyStmt }

// AlterRoleAction discriminates the action variants of ALTER ROLE.
type AlterRoleAction int

const (
	AlterRoleRename     AlterRoleAction = iota // RENAME TO <new_name>
	AlterRoleSetComment                        // SET COMMENT = '...'
	AlterRoleUnset                             // UNSET COMMENT
	AlterRoleSetTag                            // SET TAG <tag> = '<value>' [ , ... ]
	AlterRoleUnsetTag                          // UNSET TAG <tag> [ , ... ]
)

// AlterRoleStmt represents ALTER [ DATABASE ] ROLE [ IF EXISTS ] <name> <action>.
//
//	RENAME TO <new_name>
//	SET COMMENT = '<string>'
//	UNSET COMMENT
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//
// Database is true for ALTER DATABASE ROLE.
type AlterRoleStmt struct {
	Database  bool
	IfExists  bool
	Name      *ObjectName
	Action    AlterRoleAction
	NewName   *ObjectName      // RENAME TO target; for AlterRoleRename
	Comment   *string          // SET COMMENT = '...'; for AlterRoleSetComment
	Tags      []*TagAssignment // SET TAG (...) assignments; for AlterRoleSetTag
	UnsetTags []*ObjectName    // UNSET TAG names; for AlterRoleUnsetTag
	Loc       Loc
}

// Tag implements Node.
func (n *AlterRoleStmt) Tag() NodeTag { return T_AlterRoleStmt }

// AlterUserAction discriminates the action variants of ALTER USER.
type AlterUserAction int

const (
	AlterUserRename          AlterUserAction = iota // RENAME TO <new_name>
	AlterUserResetPassword                          // RESET PASSWORD
	AlterUserAbortQueries                           // ABORT ALL QUERIES
	AlterUserSet                                    // SET <options>
	AlterUserUnset                                  // UNSET <property> [ , ... ]
	AlterUserSetTag                                 // SET TAG <tag> = '<value>' [ , ... ]
	AlterUserUnsetTag                               // UNSET TAG <tag> [ , ... ]
	AlterUserAddDelegated                           // ADD DELEGATED AUTHORIZATION ...
	AlterUserRemoveDelegated                        // REMOVE DELEGATED ... FROM SECURITY INTEGRATION ...
)

// AlterUserStmt represents ALTER USER [ IF EXISTS ] <name> <action>.
//
//	RENAME TO <new_name>
//	RESET PASSWORD
//	ABORT ALL QUERIES
//	SET <options>                                  -- open-ended KEY = value
//	UNSET <property> [ , ... ]
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	ADD DELEGATED AUTHORIZATION OF ROLE <r> TO SECURITY INTEGRATION <i>
//	REMOVE DELEGATED { AUTHORIZATION OF ROLE <r> | AUTHORIZATIONS } FROM SECURITY INTEGRATION <i>
//
// The ADD / REMOVE DELEGATED forms are captured structurally only by Action
// (their body identifiers are consumed but not modeled), matching the
// open-ended philosophy; Raw holds their verbatim source for round-tripping.
type AlterUserStmt struct {
	IfExists   bool
	Name       *ObjectName
	Action     AlterUserAction
	NewName    *ObjectName      // RENAME TO target; for AlterUserRename
	Options    []*CopyOption    // SET <options>; for AlterUserSet
	UnsetProps []string         // UNSET <property> names; for AlterUserUnset
	Tags       []*TagAssignment // SET TAG (...) assignments; for AlterUserSetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterUserUnsetTag
	Raw        string           // verbatim tail for ADD / REMOVE DELEGATED forms
	Loc        Loc
}

// Tag implements Node.
func (n *AlterUserStmt) Tag() NodeTag { return T_AlterUserStmt }

// AlterPolicyAction discriminates the action variants of ALTER <...> POLICY.
type AlterPolicyAction int

const (
	AlterPolicyRename   AlterPolicyAction = iota // RENAME TO <new_name>
	AlterPolicySetBody                           // SET BODY -> <expr> (MASKING / ROW ACCESS)
	AlterPolicySet                               // SET <options> (COMMENT / option bag)
	AlterPolicyUnset                             // UNSET <property> [ , ... ]
	AlterPolicySetTag                            // SET TAG <tag> = '<value>' [ , ... ]
	AlterPolicyUnsetTag                          // UNSET TAG <tag> [ , ... ]
)

// AlterPolicyStmt represents the ALTER form of the six access-control policy
// objects (one node, discriminated by Kind):
//
//	RENAME TO <new_name>
//	SET BODY -> <expr>                  -- MASKING / ROW ACCESS only
//	SET { COMMENT = '...' | <options> }
//	UNSET { COMMENT | <property> } [ , ... ]
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
type AlterPolicyStmt struct {
	Kind       PolicyKind
	IfExists   bool
	Name       *ObjectName
	Action     AlterPolicyAction
	NewName    *ObjectName      // RENAME TO target; for AlterPolicyRename
	Body       Node             // SET BODY -> <expr>; for AlterPolicySetBody
	Options    []*CopyOption    // SET <options>; for AlterPolicySet
	UnsetProps []string         // UNSET <property> names; for AlterPolicyUnset
	Tags       []*TagAssignment // SET TAG (...) assignments; for AlterPolicySetTag
	UnsetTags  []*ObjectName    // UNSET TAG names; for AlterPolicyUnsetTag
	Loc        Loc
}

// Tag implements Node.
func (n *AlterPolicyStmt) Tag() NodeTag { return T_AlterPolicyStmt }

// Compile-time assertions for access-control DDL nodes.
var (
	_ Node = (*PolicyArg)(nil)
	_ Node = (*CreateRoleStmt)(nil)
	_ Node = (*CreateUserStmt)(nil)
	_ Node = (*CreatePolicyStmt)(nil)
	_ Node = (*AlterRoleStmt)(nil)
	_ Node = (*AlterUserStmt)(nil)
	_ Node = (*AlterPolicyStmt)(nil)
)

// ---------------------------------------------------------------------------
// Routine DDL — CREATE / ALTER FUNCTION & PROCEDURE (T4.5)
// ---------------------------------------------------------------------------
//
// A UDF / UDTF / stored procedure / external function carries a large,
// version-growing vocabulary of property clauses: LANGUAGE, RUNTIME_VERSION,
// HANDLER, PACKAGES, IMPORTS, TARGET_PATH, EXTERNAL_ACCESS_INTEGRATIONS,
// SECRETS, COMMENT, the volatility / null-handling modifiers (VOLATILE,
// IMMUTABLE, STRICT, MEMOIZABLE, CALLED ON NULL INPUT, RETURNS NULL ON NULL
// INPUT), and the external-function cloud params (API_INTEGRATION, HEADERS,
// CONTEXT_HEADERS, MAX_BATCH_ROWS, COMPRESSION, REQUEST_TRANSLATOR,
// RESPONSE_TRANSLATOR). Rather than mirror the legacy ANTLR grammar's finite,
// already-stale enumerations (its create_function rule lacks TARGET_PATH,
// EXTERNAL_ACCESS_INTEGRATIONS, SECRETS, the SCALA / JAVA-scala runtimes, ...
// all of which appear in the official docs corpus), every property is captured
// as an open-ended `KEY = <value>` pair or a bare modifier word (CopyOption),
// reusing the merged STAGE (T4.1) / FILE FORMAT (T4.2) / PIPE (T4.3) / COPY
// (T5.2) machinery. Only the structural anchors are modeled explicitly: the
// argument list, the RETURNS clause (a scalar data type or a TABLE (...) column
// list), the EXECUTE AS {CALLER|OWNER} clause, and the AS <body> clause. The
// body is OPAQUE — its language is SQL-scripting / JavaScript / Python / Java /
// Scala, never parsed here — and is captured verbatim including its delimiters,
// whether single-quoted ('...') or dollar-quoted ($$...$$). The catalog /
// semantic layer, not the parser, validates that a property is real and legal.

// RoutineArg is one argument of a function / procedure signature:
//
//	<arg_name> <data_type> [ DEFAULT <expr> ]
//
// It is a helper struct (not a Node); the walker does not descend into it,
// mirroring TagAssignment / CloneSource / ForeignKeyRef.
type RoutineArg struct {
	Name    Ident
	Type    *TypeName
	Default Node // DEFAULT <expr>; nil if absent
	Loc     Loc
}

// RoutineTableColumn is one column of a RETURNS TABLE ( ... ) result spec:
//
//	<column_name> <data_type>
//
// It is a helper struct (not a Node).
type RoutineTableColumn struct {
	Name Ident
	Type *TypeName
	Loc  Loc
}

// RoutineKind discriminates the object a routine statement targets.
type RoutineKind int

const (
	RoutineFunction         RoutineKind = iota // FUNCTION (UDF / UDTF)
	RoutineExternalFunction                    // EXTERNAL FUNCTION
	RoutineProcedure                           // PROCEDURE
)

// ExecuteAs discriminates the EXECUTE AS clause of a procedure.
type ExecuteAs int

const (
	ExecuteAsUnset  ExecuteAs = iota // no EXECUTE AS clause
	ExecuteAsCaller                  // EXECUTE AS CALLER
	ExecuteAsOwner                   // EXECUTE AS OWNER
)

// CreateRoutineStmt represents CREATE [OR REPLACE] [SECURE] [TEMP|TEMPORARY]
// FUNCTION / EXTERNAL FUNCTION / PROCEDURE.
//
//	CREATE [ OR REPLACE ] [ SECURE ] [ TEMP | TEMPORARY ]
//	  FUNCTION [ IF NOT EXISTS ] <name> ( [ <arg> [, ...] ] )
//	  RETURNS { <data_type> | TABLE ( <col> <type> [, ...] ) } [ [ NOT ] NULL ]
//	  [ <property clauses> ]                     -- open-ended KEY = value / bare words
//	  AS { '<body>' | $$<body>$$ }
//
//	CREATE [ OR REPLACE ] [ SECURE ] EXTERNAL FUNCTION <name> ( [ <arg> [, ...] ] )
//	  RETURNS <data_type> [ [ NOT ] NULL ]
//	  [ <property clauses> ] API_INTEGRATION = <id> [ <property clauses> ]
//	  AS '<url>'
//
//	CREATE [ OR REPLACE ] [ SECURE ] PROCEDURE <name> ( [ <arg> [, ...] ] )
//	  RETURNS <data_type> [ [ NOT ] NULL ]
//	  [ <property clauses> ] [ EXECUTE AS { CALLER | OWNER } ] [ <property clauses> ]
//	  AS { '<body>' | $$<body>$$ }
//
// ReturnTable is non-nil for the RETURNS TABLE ( ... ) form, in which case
// ReturnType is nil; otherwise ReturnType holds the scalar return type. Options
// holds every property clause (LANGUAGE / RUNTIME_VERSION / HANDLER / PACKAGES /
// IMPORTS / API_INTEGRATION / the volatility & null-handling modifiers / ...),
// open-ended, preserving source order. Body holds the verbatim body text WITH
// its delimiters ('...' or $$...$$); BodyDollar reports whether the
// dollar-quoted form was used. Body is "" only for the no-body external/SQL UDF
// forms the docs permit (e.g. a function defined entirely by its handler).
type CreateRoutineStmt struct {
	Kind          RoutineKind
	OrReplace     bool
	OrAlter       bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Secure        bool
	Temporary     bool
	IfNotExists   bool
	Name          *ObjectName
	Args          []*RoutineArg         // ( <arg> [, ...] ); nil for an empty arg list
	ReturnType    *TypeName             // scalar RETURNS type; nil for the TABLE form
	ReturnTable   []*RoutineTableColumn // RETURNS TABLE ( ... ) columns; non-nil for the TABLE form
	ReturnNotNull bool                  // RETURNS ... NOT NULL
	ReturnNull    bool                  // RETURNS ... NULL (explicit nullable)
	Options       []*CopyOption         // property clauses (LANGUAGE / HANDLER / volatility / ...)
	ExecuteAs     ExecuteAs             // procedure EXECUTE AS {CALLER|OWNER}; ExecuteAsUnset otherwise
	Body          string                // verbatim body incl. delimiters ('...' or $$...$$); "" if no body
	BodyDollar    bool                  // true when the $$...$$ form was used
	Loc           Loc
}

// Tag implements Node.
func (n *CreateRoutineStmt) Tag() NodeTag { return T_CreateRoutineStmt }

// AlterRoutineAction discriminates the action variants of ALTER FUNCTION /
// PROCEDURE.
type AlterRoutineAction int

const (
	AlterRoutineRename    AlterRoutineAction = iota // RENAME TO <new_name>
	AlterRoutineSet                                 // SET <options> (SECURE / COMMENT / LOG_LEVEL / ...)
	AlterRoutineUnset                               // UNSET <property> [, ...] (SECURE / COMMENT / ...)
	AlterRoutineExecuteAs                           // EXECUTE AS {CALLER|OWNER} (procedure only)
)

// AlterRoutineStmt represents
//
//	ALTER { FUNCTION | PROCEDURE } [ IF EXISTS ] <name> ( [ <argtype> [, ...] ] ) <action>
//
//	RENAME TO <new_name>
//	SET <options>                              -- open-ended (SECURE / COMMENT / LOG_LEVEL / ...)
//	UNSET <property> [ , ... ]                  -- SECURE / COMMENT / ...
//	EXECUTE AS { CALLER | OWNER }               -- procedure only
//
// The ALTER signature carries argument TYPES only (no names), captured in
// ArgTypes. Procedure is true for ALTER PROCEDURE, false for ALTER FUNCTION.
type AlterRoutineStmt struct {
	Procedure  bool // true for ALTER PROCEDURE, false for ALTER FUNCTION
	IfExists   bool
	Name       *ObjectName
	ArgTypes   []*TypeName // ( <argtype> [, ...] ) signature; nil for an empty list
	Action     AlterRoutineAction
	NewName    *ObjectName   // RENAME TO target; non-nil for AlterRoutineRename
	Options    []*CopyOption // SET <options>; for AlterRoutineSet
	UnsetProps []string      // UNSET <property> names; for AlterRoutineUnset
	ExecuteAs  ExecuteAs     // EXECUTE AS {CALLER|OWNER}; for AlterRoutineExecuteAs
	Loc        Loc
}

// Tag implements Node.
func (n *AlterRoutineStmt) Tag() NodeTag { return T_AlterRoutineStmt }

// Compile-time assertions for routine DDL nodes.
var (
	_ Node = (*CreateRoutineStmt)(nil)
	_ Node = (*AlterRoutineStmt)(nil)
)

// ---------------------------------------------------------------------------
// Table-variant + sequence DDL — CREATE / ALTER DYNAMIC / EXTERNAL / EVENT
// TABLE + SEQUENCE (T4.4)
// ---------------------------------------------------------------------------
//
// DYNAMIC / EXTERNAL / EVENT tables and SEQUENCEs each carry a large,
// version-growing configuration vocabulary that the legacy ANTLR grammar
// enumerates as an already-stale subset (its create_dynamic_table lacks
// REFRESH_MODE = ADAPTIVE | CUSTOM_INCREMENTAL, SCHEDULER, INITIALIZATION_WAREHOUSE,
// EXTERNAL_VOLUME / CATALOG / BASE_LOCATION / ICEBERG_VERSION for the ICEBERG
// variant, the REQUIRE USER and IMMUTABLE/FROZEN WHERE clauses, the REFRESH USING
// body, ...; its create_external_table lacks the USING TEMPLATE form). Matching
// the merged STAGE (T4.1) / FILE FORMAT (T4.2) / COPY (T5.2) / pipeline (T4.3)
// approach, every such config parameter that is a `KEY = <value>` pair is captured
// open-ended (CopyOption), and only the structurally distinct anchors are modeled
// as dedicated fields: the optional ICEBERG modifier, the column list, the
// CLUSTER BY clause, the CLONE source, the REQUIRE USER flag, the IMMUTABLE/FROZEN
// WHERE predicate, the LOCATION = @stage / PARTITION BY / USING TEMPLATE clauses,
// and the terminal AS <query> / REFRESH USING (<dml>) body. The catalog/semantic
// layer, not the parser, validates that an option is real and legal.

// CreateDynamicTableStmt represents
//
//	CREATE [ OR REPLACE ] [ TRANSIENT ] DYNAMIC [ ICEBERG ] TABLE [ IF NOT EXISTS ] <name>
//	  [ ( <col_name> [ <col_type> ] [ , ... ] ) ]
//	  [ TARGET_LAG = { '<dur>' | DOWNSTREAM } ] [ WAREHOUSE = <wh> ]
//	  [ REFRESH_MODE = ... ] [ INITIALIZE = ... ] [ CLUSTER BY ( ... ) ]
//	  [ DATA_RETENTION_TIME_IN_DAYS = <n> ] [ EXTERNAL_VOLUME = ... ]
//	  [ CATALOG = ... ] [ BASE_LOCATION = ... ] [ <other options> ]
//	  [ REQUIRE USER ] [ { IMMUTABLE | FROZEN } WHERE ( <expr> ) ]
//	  { AS <query> | REFRESH USING ( <dml_statement> ) }
//
//	CREATE [ OR REPLACE ] [ TRANSIENT ] DYNAMIC [ ICEBERG ] TABLE <name>
//	  CLONE <source> [ { AT | BEFORE } ( ... ) ]
//
// Iceberg records the optional ICEBERG modifier. Columns holds the optional
// column list (materialized columns: a name with an optional type, no full
// constraint vocabulary). Options holds every open-ended KEY = value parameter
// (TARGET_LAG / WAREHOUSE / REFRESH_MODE / INITIALIZE / DATA_RETENTION_TIME_IN_DAYS
// / EXTERNAL_VOLUME / CATALOG / BASE_LOCATION / ICEBERG_VERSION / SCHEDULER /
// INITIALIZATION_WAREHOUSE / COMMENT / ...), preserving source order. ClusterBy /
// Linear model the CLUSTER BY [LINEAR] (...) clause. RequireUser records REQUIRE
// USER. ImmutableWhere holds the IMMUTABLE/FROZEN WHERE (<expr>) predicate;
// ImmutableKind is "IMMUTABLE" or "FROZEN". AsQuery holds the AS <query> body;
// RefreshUsing holds the REFRESH USING (<dml>) body (one of the two is set for a
// non-CLONE statement). Clone holds the CLONE source for the CLONE form (with the
// other body fields zero).
type CreateDynamicTableStmt struct {
	OrReplace      bool
	OrAlter        bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	Transient      bool
	Iceberg        bool
	IfNotExists    bool
	Name           *ObjectName
	Columns        []*ColumnDef  // optional materialized column list; nil if absent
	Options        []*CopyOption // TARGET_LAG / WAREHOUSE / REFRESH_MODE / INITIALIZE / ...
	ClusterBy      []Node        // CLUSTER BY ( exprs ); nil if absent
	Linear         bool          // CLUSTER BY LINEAR modifier
	RequireUser    bool          // REQUIRE USER
	ImmutableKind  string        // "IMMUTABLE" or "FROZEN"; empty if no WHERE clause
	ImmutableWhere Node          // { IMMUTABLE | FROZEN } WHERE ( <expr> ); nil if absent
	AsQuery        Node          // AS <query>; nil for the CLONE / REFRESH USING forms
	RefreshUsing   Node          // REFRESH USING ( <dml_statement> ); nil otherwise
	Clone          *CloneSource  // CLONE <source>; non-nil only for the CLONE form
	Loc            Loc
}

// Tag implements Node.
func (n *CreateDynamicTableStmt) Tag() NodeTag { return T_CreateDynamicTableStmt }

// AlterDynamicTableAction discriminates the action variants of ALTER DYNAMIC TABLE.
type AlterDynamicTableAction int

const (
	AlterDynamicTableSuspend           AlterDynamicTableAction = iota // SUSPEND
	AlterDynamicTableResume                                           // RESUME
	AlterDynamicTableRefresh                                          // REFRESH [ COPY SESSION ]
	AlterDynamicTableRename                                           // RENAME TO <new_name>
	AlterDynamicTableSwap                                             // SWAP WITH <other_name>
	AlterDynamicTableSet                                              // SET <settable params>
	AlterDynamicTableUnset                                            // UNSET <property> [, ...]
	AlterDynamicTableSetTag                                           // SET TAG <tag> = '<value>' [, ...]
	AlterDynamicTableUnsetTag                                         // UNSET TAG <tag> [, ...]
	AlterDynamicTableClusterBy                                        // CLUSTER BY ( exprs )
	AlterDynamicTableDropClusteringKey                                // DROP CLUSTERING KEY
)

// AlterDynamicTableStmt represents ALTER DYNAMIC TABLE [ IF EXISTS ] <name> <action>.
//
//	{ SUSPEND | RESUME }
//	REFRESH [ COPY SESSION ]
//	RENAME TO <new_name>
//	SWAP WITH <target_dynamic_table_name>
//	SET <settable params>                      -- open-ended KEY = value
//	UNSET <property> [ , ... ]
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	CLUSTER BY ( <expr> [ , ... ] )
//	DROP CLUSTERING KEY
//
// RefreshCopySession records the optional COPY SESSION on the REFRESH action.
type AlterDynamicTableStmt struct {
	IfExists           bool
	Name               *ObjectName
	Action             AlterDynamicTableAction
	NewName            *ObjectName      // RENAME TO / SWAP WITH target
	Options            []*CopyOption    // SET <settable params>; for AlterDynamicTableSet
	Tags               []*TagAssignment // SET TAG assignments; for AlterDynamicTableSetTag
	UnsetTags          []*ObjectName    // UNSET TAG names; for AlterDynamicTableUnsetTag
	UnsetProps         []string         // UNSET <property> names; for AlterDynamicTableUnset
	ClusterBy          []Node           // CLUSTER BY exprs; for AlterDynamicTableClusterBy
	Linear             bool             // CLUSTER BY LINEAR modifier
	RefreshCopySession bool             // REFRESH COPY SESSION; for AlterDynamicTableRefresh
	Loc                Loc
}

// Tag implements Node.
func (n *AlterDynamicTableStmt) Tag() NodeTag { return T_AlterDynamicTableStmt }

// ExternalColumnDef represents one column of a CREATE EXTERNAL TABLE column list:
//
//	<col_name> <col_type> AS ( <expr> | <id> ) [ inlineConstraint ]
//
// Unlike an ordinary ColumnDef, an external-table column's value is always derived
// from an AS expression over the staged file (the data type is required, and the
// AS clause is mandatory). Expr holds the AS expression (parenthesized or not).
type ExternalColumnDef struct {
	Name       Ident
	DataType   *TypeName         // required (per the docs / legacy grammar)
	Expr       Node              // AS ( <expr> ) | AS <id>; the partition/derivation expression
	NotNull    bool              // inline NOT NULL
	Constraint *InlineConstraint // inline PRIMARY KEY / UNIQUE / FOREIGN KEY; nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *ExternalColumnDef) Tag() NodeTag { return T_ExternalColumnDef }

// CreateExternalTableStmt represents
//
//	CREATE [ OR REPLACE ] EXTERNAL TABLE [ IF NOT EXISTS ] <name>
//	  ( <col_name> <col_type> AS ( <expr> ) [ inlineConstraint ] [ , ... ] )
//	  [ INTEGRATION = '<integration>' ] [ PARTITION BY ( <expr> [ , ... ] ) ]
//	  [ WITH ] LOCATION = @<stage>[/<path>]
//	  [ REFRESH_ON_CREATE = { TRUE | FALSE } ] [ AUTO_REFRESH = { TRUE | FALSE } ]
//	  [ PATTERN = '<regex>' ] [ PARTITION_TYPE = USER_SPECIFIED ]
//	  FILE_FORMAT = ( ... ) [ TABLE_FORMAT = DELTA ] [ AWS_SNS_TOPIC = '...' ]
//	  [ COPY GRANTS ] [ <with_row_access_policy> ] [ [ WITH ] TAG ( ... ) ]
//	  [ COMMENT = '...' ]
//
//	CREATE [ OR REPLACE ] EXTERNAL TABLE [ IF NOT EXISTS ] <name>
//	  USING TEMPLATE ( <query> )
//	  [ WITH ] LOCATION = @<stage> [ <other options> ]
//
// Columns holds the explicit column list (nil for the USING TEMPLATE form).
// UsingTemplate holds the USING TEMPLATE ( <query> ) body (nil for the explicit
// column-list form). PartitionBy holds the PARTITION BY ( <expr> [, ...] ) clause.
// Location holds the mandatory LOCATION = @stage reference. Options holds every
// open-ended KEY = value parameter (REFRESH_ON_CREATE / AUTO_REFRESH / PATTERN /
// PARTITION_TYPE / FILE_FORMAT / TABLE_FORMAT / AWS_SNS_TOPIC / INTEGRATION /
// FILE_FORMAT / COMMENT / ...), preserving source order. CopyGrants records COPY
// GRANTS; Tags holds the [WITH] TAG assignments.
type CreateExternalTableStmt struct {
	OrReplace     bool
	OrAlter       bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists   bool
	Name          *ObjectName
	Columns       []*ExternalColumnDef // explicit column list; nil for USING TEMPLATE
	UsingTemplate Node                 // USING TEMPLATE ( <query> ); nil for the column-list form
	PartitionBy   []Node               // PARTITION BY ( exprs ); nil if absent
	Location      *StageLocation       // LOCATION = @stage; non-nil for a well-formed statement
	Options       []*CopyOption        // REFRESH_ON_CREATE / AUTO_REFRESH / PATTERN / FILE_FORMAT / ...
	CopyGrants    bool                 // COPY GRANTS
	Tags          []*TagAssignment     // [WITH] TAG (...); nil if absent
	Loc           Loc
}

// Tag implements Node.
func (n *CreateExternalTableStmt) Tag() NodeTag { return T_CreateExternalTableStmt }

// ExternalTablePartition is one `<col> = '<value>'` pair in an ALTER EXTERNAL
// TABLE ADD PARTITION (...) clause.
type ExternalTablePartition struct {
	Column Ident
	Value  string // the '<string>' partition value
}

// AlterExternalTableAction discriminates the action variants of ALTER EXTERNAL TABLE.
type AlterExternalTableAction int

const (
	AlterExternalTableRefresh       AlterExternalTableAction = iota // REFRESH [ '<path>' ]
	AlterExternalTableAddFiles                                      // ADD FILES ( '<f>' [, ...] )
	AlterExternalTableRemoveFiles                                   // REMOVE FILES ( '<f>' [, ...] )
	AlterExternalTableSet                                           // SET [ AUTO_REFRESH = ... ] [ TAG ... ]
	AlterExternalTableUnsetTag                                      // UNSET TAG <tag> [, ...]
	AlterExternalTableAddPartition                                  // ADD PARTITION ( col=val,... ) LOCATION '<path>'
	AlterExternalTableDropPartition                                 // DROP PARTITION LOCATION '<path>'
)

// AlterExternalTableStmt represents ALTER EXTERNAL TABLE [ IF EXISTS ] <name> <action>.
//
//	REFRESH [ '<relative-path>' ]
//	ADD FILES ( '<path>/<file>' [ , ... ] )
//	REMOVE FILES ( '<path>/<file>' [ , ... ] )
//	SET [ AUTO_REFRESH = { TRUE | FALSE } ] [ TAG <tag> = '<value>' [ , ... ] ]
//	UNSET TAG <tag> [ , ... ]
//	ADD PARTITION ( <col> = '<val>' [ , ... ] ) LOCATION '<path>'
//	DROP PARTITION LOCATION '<path>'
//
// RefreshPath holds the optional REFRESH path; Files holds the ADD/REMOVE FILES
// list; Options/Tags hold the SET params/tags; Partitions + Location hold the
// ADD PARTITION pairs and path; DropLocation holds the DROP PARTITION path.
type AlterExternalTableStmt struct {
	IfExists     bool
	Name         *ObjectName
	Action       AlterExternalTableAction
	RefreshPath  *string                   // REFRESH '<path>'; nil if bare REFRESH
	Files        []*Literal                // ADD/REMOVE FILES list
	Options      []*CopyOption             // SET [ AUTO_REFRESH = ... ]; for AlterExternalTableSet
	Tags         []*TagAssignment          // SET ... TAG / for AlterExternalTableSet
	UnsetTags    []*ObjectName             // UNSET TAG names; for AlterExternalTableUnsetTag
	Partitions   []*ExternalTablePartition // ADD PARTITION pairs; for AlterExternalTableAddPartition
	Location     *string                   // ADD PARTITION ... LOCATION '<path>'
	DropLocation *string                   // DROP PARTITION LOCATION '<path>'
	Loc          Loc
}

// Tag implements Node.
func (n *AlterExternalTableStmt) Tag() NodeTag { return T_AlterExternalTableStmt }

// CreateEventTableStmt represents
//
//	CREATE [ OR REPLACE ] EVENT TABLE [ IF NOT EXISTS ] <name>
//	  [ CLUSTER BY ( <expr> [ , ... ] ) ]
//	  [ DATA_RETENTION_TIME_IN_DAYS = <n> ] [ MAX_DATA_EXTENSION_TIME_IN_DAYS = <n> ]
//	  [ CHANGE_TRACKING = { TRUE | FALSE } ] [ DEFAULT_DDL_COLLATION = '<collation>' ]
//	  [ COPY GRANTS ] [ <with_row_access_policy> ] [ [ WITH ] TAG ( ... ) ]
//	  [ [ WITH ] COMMENT = '<string_literal>' ]
//
// An EVENT TABLE has a fixed, system-defined schema (no user column list).
// ClusterBy / Linear model the optional CLUSTER BY [LINEAR] (...) clause. Options
// holds every open-ended KEY = value parameter (DATA_RETENTION_TIME_IN_DAYS /
// MAX_DATA_EXTENSION_TIME_IN_DAYS / CHANGE_TRACKING / DEFAULT_DDL_COLLATION /
// COMMENT / ...), preserving source order. CopyGrants records COPY GRANTS; Tags
// holds the [WITH] TAG assignments.
type CreateEventTableStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	ClusterBy   []Node           // CLUSTER BY ( exprs ); nil if absent
	Linear      bool             // CLUSTER BY LINEAR modifier
	Options     []*CopyOption    // DATA_RETENTION_TIME_IN_DAYS / CHANGE_TRACKING / COMMENT / ...
	CopyGrants  bool             // COPY GRANTS
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateEventTableStmt) Tag() NodeTag { return T_CreateEventTableStmt }

// CreateSequenceStmt represents
//
//	CREATE [ OR REPLACE ] SEQUENCE [ IF NOT EXISTS ] <name>
//	  [ WITH ]
//	  [ START [ WITH ] [ = ] <initial_value> ]
//	  [ INCREMENT [ BY ] [ = ] <sequence_interval> ]
//	  [ { ORDER | NOORDER } ]
//	  [ COMMENT = '<string_literal>' ]
//
// Start / Increment hold the optional START / INCREMENT integer values (each may
// be negative; nil when the clause is omitted). Order is true for ORDER, false
// for NOORDER, nil when unspecified. Comment holds the COMMENT clause text.
type CreateSequenceStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Start       *int64  // START [WITH] [=] <n>; nil if absent
	Increment   *int64  // INCREMENT [BY] [=] <n>; nil if absent
	Order       *bool   // true=ORDER, false=NOORDER, nil=unspecified
	Comment     *string // COMMENT = '<text>'; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateSequenceStmt) Tag() NodeTag { return T_CreateSequenceStmt }

// AlterSequenceAction discriminates the action variants of ALTER SEQUENCE.
type AlterSequenceAction int

const (
	AlterSequenceRename       AlterSequenceAction = iota // RENAME TO <new_name>
	AlterSequenceSetIncrement                            // [ SET ] INCREMENT [ BY ] [ = ] <n>
	AlterSequenceSet                                     // SET [ { ORDER | NOORDER } ] [ COMMENT = '...' ]
	AlterSequenceUnsetComment                            // UNSET COMMENT
)

// AlterSequenceStmt represents ALTER SEQUENCE [ IF EXISTS ] <name> <action>.
//
//	RENAME TO <new_name>
//	[ SET ] INCREMENT [ BY ] [ = ] <sequence_interval>
//	SET [ { ORDER | NOORDER } ] [ COMMENT = '<string_literal>' ]
//	UNSET COMMENT
//
// NewName holds the RENAME target. Increment holds the new interval (may be
// negative) for the SET INCREMENT action. Order / Comment hold the SET
// ORDER|NOORDER / COMMENT values.
type AlterSequenceStmt struct {
	IfExists  bool
	Name      *ObjectName
	Action    AlterSequenceAction
	NewName   *ObjectName // RENAME TO target
	Increment *int64      // SET INCREMENT value; for AlterSequenceSetIncrement
	Order     *bool       // SET ORDER|NOORDER; for AlterSequenceSet
	Comment   *string     // SET COMMENT value; for AlterSequenceSet
	Loc       Loc
}

// Tag implements Node.
func (n *AlterSequenceStmt) Tag() NodeTag { return T_AlterSequenceStmt }

// Compile-time assertions for table-variant + sequence DDL nodes.
var (
	_ Node = (*CreateDynamicTableStmt)(nil)
	_ Node = (*AlterDynamicTableStmt)(nil)
	_ Node = (*ExternalColumnDef)(nil)
	_ Node = (*CreateExternalTableStmt)(nil)
	_ Node = (*AlterExternalTableStmt)(nil)
	_ Node = (*CreateEventTableStmt)(nil)
	_ Node = (*CreateSequenceStmt)(nil)
	_ Node = (*AlterSequenceStmt)(nil)
)

// ---------------------------------------------------------------------------
// Warehouse DDL — CREATE / ALTER WAREHOUSE (gap-warehouse)
// ---------------------------------------------------------------------------
//
// A warehouse carries an open-ended, version-growing vocabulary of object
// properties (WAREHOUSE_SIZE, WAREHOUSE_TYPE, RESOURCE_CONSTRAINT, AUTO_RESUME,
// AUTO_SUSPEND, INITIALLY_SUSPENDED, GENERATION, MAX_CLUSTER_COUNT,
// MIN_CLUSTER_COUNT, SCALING_POLICY, ENABLE_QUERY_ACCELERATION, COMMENT, ...)
// plus object-parameters and session-parameters that share the same
// `KEY = value` shape. Rather than mirror the legacy ANTLR grammar's finite,
// already-stale enumeration (its wh_properties / wh_common_size rules predate
// WAREHOUSE_TYPE, RESOURCE_CONSTRAINT, GENERATION, ENABLE_QUERY_ACCELERATION,
// QUERY_ACCELERATION_MAX_SCALE_FACTOR, ...), every property is captured as an
// open-ended CopyOption pair, exactly like STAGE (T4.1) / FILE FORMAT (T4.2).
// The catalog/semantic layer, not the parser, validates that a property is real
// and legal. The trailing [ WITH ] TAG (...) clause is the one structural anchor.

// CreateWarehouseStmt represents
//
//	CREATE [ OR REPLACE | OR ALTER ] WAREHOUSE [ IF NOT EXISTS ] <name>
//	  [ WITH ] <prop> [ <prop> ... ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//
// where each <prop> is an open-ended `KEY = value` pair captured in Options,
// preserving source order. The leading optional WITH (CREATE WAREHOUSE w WITH
// WAREHOUSE_SIZE=...) is purely cosmetic and not retained.
type CreateWarehouseStmt struct {
	OrReplace   bool
	OrAlter     bool // CREATE OR ALTER (mutually exclusive with OrReplace)
	IfNotExists bool
	Name        *ObjectName
	Options     []*CopyOption    // open-ended warehouse properties / object+session params
	Tags        []*TagAssignment // [WITH] TAG (...); nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *CreateWarehouseStmt) Tag() NodeTag { return T_CreateWarehouseStmt }

// AlterWarehouseAction discriminates the action variants of ALTER WAREHOUSE.
type AlterWarehouseAction int

const (
	AlterWarehouseSuspend      AlterWarehouseAction = iota // SUSPEND
	AlterWarehouseResume                                   // RESUME [ IF SUSPENDED ]
	AlterWarehouseAbort                                    // ABORT ALL QUERIES
	AlterWarehouseRename                                   // RENAME TO <new_name>
	AlterWarehouseSet                                      // SET <prop> [ <prop> ... ]
	AlterWarehouseUnset                                    // UNSET <key> [ , ... ]
	AlterWarehouseSetTag                                   // SET TAG <tag> = '<value>' [ , ... ]
	AlterWarehouseUnsetTag                                 // UNSET TAG <tag> [ , ... ]
	AlterWarehouseAddTables                                // ADD TABLES ( <id> [ , ... ] )
	AlterWarehouseRemoveTables                             // { REMOVE | DROP } TABLES ( <id> [ , ... ] )
)

// AlterWarehouseStmt represents ALTER WAREHOUSE [ IF EXISTS ] <name> <action>.
//
//	SUSPEND
//	RESUME [ IF SUSPENDED ]
//	ABORT ALL QUERIES
//	RENAME TO <new_name>
//	SET <prop> [ <prop> ... ]               -- open-ended KEY = value params
//	UNSET <key> [ , ... ]
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	ADD TABLES ( <id> [ , ... ] )           -- Unistore/interactive variant
//	{ REMOVE | DROP } TABLES ( <id> [ , ... ] )
//
// NewName holds the RENAME target. Options holds the SET properties. UnsetKeys
// holds the UNSET key names. Tags / UnsetTags hold the SET TAG / UNSET TAG
// assignments and names. Tables holds the ADD/REMOVE TABLES identifier list.
// ResumeIfSuspended records the optional IF SUSPENDED on RESUME.
type AlterWarehouseStmt struct {
	IfExists          bool
	Name              *ObjectName
	Action            AlterWarehouseAction
	NewName           *ObjectName      // RENAME TO target; for AlterWarehouseRename
	Options           []*CopyOption    // SET <props>; for AlterWarehouseSet
	UnsetKeys         []string         // UNSET <key> [, ...]; for AlterWarehouseUnset
	Tags              []*TagAssignment // SET TAG (...) assignments; for AlterWarehouseSetTag
	UnsetTags         []*ObjectName    // UNSET TAG (...) names; for AlterWarehouseUnsetTag
	Tables            []*ObjectName    // ADD/REMOVE TABLES (...) ids; for the table actions
	ResumeIfSuspended bool             // RESUME IF SUSPENDED; for AlterWarehouseResume
	Loc               Loc
}

// Tag implements Node.
func (n *AlterWarehouseStmt) Tag() NodeTag { return T_AlterWarehouseStmt }

// Compile-time assertions for warehouse DDL nodes.
var (
	_ Node = (*CreateWarehouseStmt)(nil)
	_ Node = (*AlterWarehouseStmt)(nil)
)

// ---------------------------------------------------------------------------
// NETWORK RULE DDL statement nodes (gap-network-rule)
// ---------------------------------------------------------------------------
//
// A network rule carries a small but version-growing set of `KEY = value`
// properties — TYPE (IPV4 / IPV6 / AWSVPCEID / HOST_PORT / PRIVATE_HOST_PORT /
// COMPUTE_POOL / ...), VALUE_LIST = ( '...' [ , ... ] ), MODE (INGRESS /
// EGRESS / INTERNAL_STAGE / ...), and COMMENT. Rather than enumerate the
// (already-growing) TYPE / MODE vocabularies, every property is parsed as an
// open-ended `KEY = <value>` pair (ast.CopyOption), reusing the COPY (T5.2)
// option machinery exactly as STAGE (T4.1), FILE FORMAT (T4.2) and WAREHOUSE
// (gap-warehouse) do. The parenthesized VALUE_LIST string list is already a
// native CopyOption value shape (List). Property order is free.

// CreateNetworkRuleStmt represents
//
//	CREATE [ OR REPLACE ] NETWORK RULE [ IF NOT EXISTS ] <name> <prop> [ <prop> ... ]
//
// where each <prop> is an open-ended `KEY = value` pair (TYPE, VALUE_LIST,
// MODE, COMMENT, ...). Options preserves source order.
type CreateNetworkRuleStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *ObjectName
	Options     []*CopyOption // open-ended network-rule properties (TYPE / VALUE_LIST / MODE / COMMENT / ...)
	Loc         Loc
}

// Tag implements Node.
func (n *CreateNetworkRuleStmt) Tag() NodeTag { return T_CreateNetworkRuleStmt }

// AlterNetworkRuleAction discriminates the action variants of ALTER NETWORK
// RULE.
type AlterNetworkRuleAction int

const (
	AlterNetworkRuleSet   AlterNetworkRuleAction = iota // SET <prop> [ <prop> ... ]
	AlterNetworkRuleUnset                               // UNSET <key> [ , ... ]
)

// AlterNetworkRuleStmt represents
//
//	ALTER NETWORK RULE [ IF EXISTS ] <name> { SET <prop>... | UNSET <key>,... }
//
// SET carries an open-ended `KEY = value` property run (VALUE_LIST, COMMENT,
// ...) in Options; UNSET carries the uppercased key names (VALUE_LIST,
// COMMENT, ...) in UnsetKeys.
type AlterNetworkRuleStmt struct {
	IfExists  bool
	Name      *ObjectName
	Action    AlterNetworkRuleAction
	Options   []*CopyOption // SET <props>; for AlterNetworkRuleSet
	UnsetKeys []string      // UNSET <key> [ , ... ]; for AlterNetworkRuleUnset
	Loc       Loc
}

// Tag implements Node.
func (n *AlterNetworkRuleStmt) Tag() NodeTag { return T_AlterNetworkRuleStmt }

// Compile-time assertions for network-rule DDL nodes.
var (
	_ Node = (*CreateNetworkRuleStmt)(nil)
	_ Node = (*AlterNetworkRuleStmt)(nil)
)

// ---------------------------------------------------------------------------
// CALL / EXECUTE IMMEDIATE / EXECUTE TASK / EXPLAIN statement nodes (T5.4)
// ---------------------------------------------------------------------------

// CallArg is one argument of a CALL statement. Snowflake stored procedures
// accept either positional arguments or named arguments using the `=>` syntax
// (CALL p(province => 'MB')). For a positional argument Name is the zero Ident
// (IsEmpty); for a named argument Name carries the parameter name. Value is the
// argument expression (any expression node; a subquery argument is a
// SubqueryExpr or, per Snowflake docs, a bare SELECT — see ParenlessQuery).
type CallArg struct {
	Name  Ident // optional parameter name for `name => value`; zero Ident when positional
	Value Node  // the argument expression
	Loc   Loc
}

// Tag implements Node.
func (n *CallArg) Tag() NodeTag { return T_CallArg }

// CallStmt represents `CALL <proc_name> ( [ <arg> [ , ... ] ] )`.
//
// Syntax (docs + legacy .g4):
//
//	CALL <procedure_name> ( [ <arg> , ... ] )
//
// Args carries the (possibly empty) argument list. Every element is a *CallArg
// (a positional argument has a zero Name; a named argument records its name).
// Args is typed []Node so the generated walker descends into each *CallArg —
// and thence into the argument expression — via walkNodes.
type CallStmt struct {
	Name *ObjectName
	Args []Node // each element is a *CallArg
	Loc  Loc
}

// Tag implements Node.
func (n *CallStmt) Tag() NodeTag { return T_CallStmt }

// ExecImmSource classifies the input form of an EXECUTE IMMEDIATE statement.
//
// Snowflake documents three forms (docs win over the legacy grammar, which adds
// the equivalent DBL_DOLLAR alternative):
//
//	EXECUTE IMMEDIATE '<string_literal>'   -> ExecImmString
//	EXECUTE IMMEDIATE $$<dollar_body>$$    -> ExecImmDollar
//	EXECUTE IMMEDIATE <variable>           -> ExecImmVariable  (Snowflake Scripting local var, no $)
//	EXECUTE IMMEDIATE $<session_variable>  -> ExecImmSessionVar ($-prefixed session variable)
type ExecImmSource int

const (
	// ExecImmString is a single-quoted string-literal body.
	ExecImmString ExecImmSource = iota
	// ExecImmDollar is a $$...$$ dollar-quoted body.
	ExecImmDollar
	// ExecImmVariable is a bare Snowflake Scripting local variable (no $).
	ExecImmVariable
	// ExecImmSessionVar is a $-prefixed session variable.
	ExecImmSessionVar
)

// String returns a human-readable name for the source kind.
func (k ExecImmSource) String() string {
	switch k {
	case ExecImmString:
		return "STRING"
	case ExecImmDollar:
		return "DOLLAR"
	case ExecImmVariable:
		return "VARIABLE"
	case ExecImmSessionVar:
		return "SESSION_VARIABLE"
	default:
		return "UNKNOWN"
	}
}

// ExecuteImmediateStmt represents an EXECUTE IMMEDIATE statement.
//
// Syntax (docs):
//
//	EXECUTE IMMEDIATE { '<string_literal>' | $$<body>$$ | <variable> | $<session_variable> }
//	  [ USING ( <bind_variable> [ , ... ] ) ]
//
// The body of the string / dollar forms is OPAQUE: it is captured VERBATIM
// (delimiters included for the dollar form; quotes included for the string form)
// and never parsed as SQL — its contents may be a Snowflake Scripting block, a
// single SQL statement, or any string the engine evaluates at runtime. Source
// records which form was used. For the string / dollar forms Body holds the
// verbatim text; for the variable / session-variable forms Var holds the name
// (for $<session_variable>, Var.Name excludes the leading $). Using holds the
// optional bind-variable list (bare identifiers per the legacy grammar and the
// official corpus).
type ExecuteImmediateStmt struct {
	Source ExecImmSource
	Body   string  // verbatim body for ExecImmString / ExecImmDollar (incl. delimiters); "" otherwise
	Var    Ident   // variable name for ExecImmVariable / ExecImmSessionVar; zero Ident otherwise
	Using  []Ident // optional USING ( <bind_variable> , ... ) list
	Loc    Loc
}

// Tag implements Node.
func (n *ExecuteImmediateStmt) Tag() NodeTag { return T_ExecuteImmediateStmt }

// ExecuteTaskStmt represents an EXECUTE TASK statement.
//
// Syntax (docs):
//
//	EXECUTE TASK <name> [ USING CONFIG = <configuration_string> ]
//	EXECUTE TASK <name> RETRY LAST
//
// RetryLast records the RETRY LAST form. UsingConfig holds the verbatim
// configuration string (including quotes) when `USING CONFIG = '<...>'` is
// present, or "" when absent. The two trailing forms are mutually exclusive per
// the docs.
type ExecuteTaskStmt struct {
	Name        *ObjectName
	RetryLast   bool
	UsingConfig string // verbatim USING CONFIG = value (incl. quotes); "" when absent
	Loc         Loc
}

// Tag implements Node.
func (n *ExecuteTaskStmt) Tag() NodeTag { return T_ExecuteTaskStmt }

// ExplainFormat enumerates the EXPLAIN output format selected by the optional
// USING clause. ExplainDefault is the bare `EXPLAIN <statement>` form (no USING).
type ExplainFormat int

const (
	// ExplainDefault is `EXPLAIN <statement>` with no USING clause.
	ExplainDefault ExplainFormat = iota
	// ExplainTabular is `EXPLAIN USING TABULAR <statement>`.
	ExplainTabular
	// ExplainJSON is `EXPLAIN USING JSON <statement>`.
	ExplainJSON
	// ExplainText is `EXPLAIN USING TEXT <statement>`.
	ExplainText
)

// String returns the USING-clause keyword for the format (or "DEFAULT").
func (f ExplainFormat) String() string {
	switch f {
	case ExplainDefault:
		return "DEFAULT"
	case ExplainTabular:
		return "TABULAR"
	case ExplainJSON:
		return "JSON"
	case ExplainText:
		return "TEXT"
	default:
		return "UNKNOWN"
	}
}

// ExplainStmt represents an EXPLAIN statement.
//
// Syntax (docs + legacy .g4):
//
//	EXPLAIN [ USING { TABULAR | JSON | TEXT } ] <statement>
//
// Format records the optional USING format (ExplainDefault when absent). Stmt is
// the inner statement, parsed structurally via the top-level statement parser.
type ExplainStmt struct {
	Format ExplainFormat
	Stmt   Node
	Loc    Loc
}

// Tag implements Node.
func (n *ExplainStmt) Tag() NodeTag { return T_ExplainStmt }

// Compile-time assertions for CALL / EXECUTE / EXPLAIN nodes (T5.4).
var (
	_ Node = (*CallArg)(nil)
	_ Node = (*CallStmt)(nil)
	_ Node = (*ExecuteImmediateStmt)(nil)
	_ Node = (*ExecuteTaskStmt)(nil)
	_ Node = (*ExplainStmt)(nil)
)

// ---------------------------------------------------------------------------
// Snowflake Scripting (T7.1)
//
// A Snowflake Scripting block is a structured procedural body that may appear
// as a top-level statement (an anonymous block) or as the body of a CREATE
// PROCEDURE / TASK. These nodes model the body STRUCTURALLY: the block, its
// declarations, its statements, and the control-flow constructs (IF / CASE /
// FOR / WHILE / REPEAT / LOOP / cursors / exceptions). Expressions and nested
// SQL are parsed by the engine's existing parseExpr / parseStmt machinery.
//
// Grammar (Snowflake Scripting reference; docs are authoritative — the legacy
// SnowflakeParser.g4 modeled only a thin subset: BEGIN..END, a DECLARE list of
// `id data_type`, `id := expr`, and `RETURN expr`):
//
//	[ DECLARE <declaration>; [ <declaration>; ... ] ]
//	BEGIN
//	    <statement>; [ <statement>; ... ]
//	[ EXCEPTION <handler> ... ]
//	END [ <label> ] ;
//
// Walker note: the structural sub-nodes (ScriptDeclaration, ScriptIfBranch,
// ScriptCaseWhen, ScriptExceptionHandler) ARE Nodes and are held in []Node
// slices, so the generated walker descends into the nested bodies and queries
// they carry. (This deviates from the CaseExpr/WhenClause precedent, whose WHEN
// bodies are unreachable — a gap an analysis pass over scripting must not have.)
// ---------------------------------------------------------------------------

// ScriptDeclKind classifies a single entry in a DECLARE section.
type ScriptDeclKind int

const (
	// ScriptDeclVar is `<name> [<type>] [ { DEFAULT | := } <expr> ]`.
	ScriptDeclVar ScriptDeclKind = iota
	// ScriptDeclCursor is `<name> CURSOR FOR <query>`.
	ScriptDeclCursor
	// ScriptDeclResultset is `<name> RESULTSET [ { DEFAULT | := } [ ASYNC ] ( <query> ) ]`.
	ScriptDeclResultset
	// ScriptDeclException is `<name> EXCEPTION [ ( <number>, '<message>' ) ]`.
	ScriptDeclException
)

// ScriptDeclaration is one entry in a DECLARE section. The active fields depend
// on Kind:
//
//	ScriptDeclVar       — Name, optional Type, optional Default (the := / DEFAULT value)
//	ScriptDeclCursor    — Name, Query (the SELECT after CURSOR FOR)
//	ScriptDeclResultset — Name, optional Query (the ( <query> ) after := / DEFAULT), Async
//	ScriptDeclException — Name, optional ExcArgs (the ( number, 'message' ) expr list)
type ScriptDeclaration struct {
	Kind    ScriptDeclKind
	Name    Ident
	Type    *TypeName // ScriptDeclVar only; nil when the type is inferred from the value
	Default Node      // ScriptDeclVar value (DEFAULT / := expr); nil if none
	Query   Node      // ScriptDeclCursor query; ScriptDeclResultset query; nil otherwise
	Async   bool      // ScriptDeclResultset: the ASYNC modifier was present
	ExcArgs []Node    // ScriptDeclException ( number, 'message' ) args; nil if none
	Loc     Loc
}

// Tag implements Node.
func (n *ScriptDeclaration) Tag() NodeTag { return T_ScriptDeclaration }

// ScriptBlockStmt represents a Snowflake Scripting block:
//
//	[ DECLARE <decls> ] BEGIN <stmts> [ EXCEPTION <handlers> ] END [ <label> ]
//
// Decls is nil when there is no DECLARE section (elements are
// *ScriptDeclaration). Handlers is nil when there is no EXCEPTION section
// (elements are *ScriptExceptionHandler). Label is the zero Ident when END
// carries no label.
type ScriptBlockStmt struct {
	Decls    []Node // *ScriptDeclaration entries
	Body     []Node // BEGIN ... statements
	Handlers []Node // *ScriptExceptionHandler entries
	Label    Ident  // optional label after END; zero Ident when absent
	Loc      Loc
}

// Tag implements Node.
func (n *ScriptBlockStmt) Tag() NodeTag { return T_ScriptBlockStmt }

// ScriptExceptionHandler is one WHEN clause in an EXCEPTION section:
//
//	WHEN <exc_name> [ OR <exc_name> ... ] THEN <stmts>
//	WHEN OTHER THEN <stmts>
//
// Other is true for the `WHEN OTHER` catch-all (in which case Names is nil).
type ScriptExceptionHandler struct {
	Names []Ident // exception names (OR-separated); nil when Other is true
	Other bool    // WHEN OTHER catch-all
	Body  []Node  // statements after THEN
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptExceptionHandler) Tag() NodeTag { return T_ScriptExceptionHandler }

// ScriptAssignStmt represents a scripting assignment `<var> := <expr>`. The
// target may be a qualified name (e.g. a RESULTSET) but is captured as a plain
// Ident path component; only single-identifier targets occur in practice.
type ScriptAssignStmt struct {
	Target Ident
	Value  Node
	Loc    Loc
}

// Tag implements Node.
func (n *ScriptAssignStmt) Tag() NodeTag { return T_ScriptAssignStmt }

// ScriptLetStmt represents a LET declaration-statement inside a block:
//
//	LET <var> [ <type> ] { DEFAULT | := } <expr>
//	LET <cursor> CURSOR FOR <query>
//	LET <resultset> RESULTSET { DEFAULT | := } [ ASYNC ] ( <query> )
//
// Kind reuses ScriptDeclKind (LET cannot declare an EXCEPTION; only Var /
// Cursor / Resultset are valid). The active fields mirror ScriptDeclaration.
type ScriptLetStmt struct {
	Kind    ScriptDeclKind
	Name    Ident
	Type    *TypeName // Var only; nil when inferred
	Default Node      // Var value; nil if none (a LET var always has a value, but kept Node-typed)
	Query   Node      // Cursor / Resultset query
	Async   bool      // Resultset ASYNC
	Loc     Loc
}

// Tag implements Node.
func (n *ScriptLetStmt) Tag() NodeTag { return T_ScriptLetStmt }

// ScriptIfBranch is one IF / ELSEIF arm: a condition plus its THEN body.
type ScriptIfBranch struct {
	Cond Node
	Body []Node
	Loc  Loc
}

// Tag implements Node.
func (n *ScriptIfBranch) Tag() NodeTag { return T_ScriptIfBranch }

// ScriptIfStmt represents an IF statement:
//
//	IF ( <cond> ) THEN <stmts>
//	[ ELSEIF ( <cond> ) THEN <stmts> ]*
//	[ ELSE <stmts> ]
//	END IF
//
// Branches holds the leading IF arm followed by each ELSEIF arm (in order;
// elements are *ScriptIfBranch). Else holds the optional ELSE body (nil when
// absent).
type ScriptIfStmt struct {
	Branches []Node // *ScriptIfBranch entries
	Else     []Node
	Loc      Loc
}

// Tag implements Node.
func (n *ScriptIfStmt) Tag() NodeTag { return T_ScriptIfStmt }

// ScriptCaseWhen is one WHEN arm of a CASE statement: a match/condition
// expression plus its THEN body.
type ScriptCaseWhen struct {
	Match Node // simple form: value to compare; searched form: boolean condition
	Body  []Node
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptCaseWhen) Tag() NodeTag { return T_ScriptCaseWhen }

// ScriptCaseStmt represents a CASE statement (simple or searched):
//
//	CASE [ ( <operand> ) ] WHEN <expr> THEN <stmts> [ WHEN ... ]* [ ELSE <stmts> ] END [ CASE ]
//
// Operand is non-nil for the simple form (CASE <expr> WHEN ...) and nil for the
// searched form (CASE WHEN <cond> ...). Whens elements are *ScriptCaseWhen.
// Else holds the optional ELSE body.
type ScriptCaseStmt struct {
	Operand Node   // nil for searched form
	Whens   []Node // *ScriptCaseWhen entries
	Else    []Node
	Loc     Loc
}

// Tag implements Node.
func (n *ScriptCaseStmt) Tag() NodeTag { return T_ScriptCaseStmt }

// ScriptForKind distinguishes the counter and cursor/resultset FOR forms.
type ScriptForKind int

const (
	// ScriptForCounter is `FOR <var> IN [ REVERSE ] <start> TO <end> { DO | LOOP } ... END { FOR | LOOP }`.
	ScriptForCounter ScriptForKind = iota
	// ScriptForCursor is `FOR <row> IN { <cursor> | <resultset> } DO ... END FOR`.
	ScriptForCursor
)

// ScriptForStmt represents a FOR loop in either form.
//
// Counter form fields: Var, Reverse, Start, End.
// Cursor form fields:  Var, Source (the cursor / resultset name).
type ScriptForStmt struct {
	Kind    ScriptForKind
	Var     Ident
	Reverse bool  // counter form: the REVERSE modifier
	Start   Node  // counter form: lower bound
	End     Node  // counter form: upper bound
	Source  Ident // cursor form: the cursor / resultset name
	Body    []Node
	Label   Ident // optional label after END; zero Ident when absent
	Loc     Loc
}

// Tag implements Node.
func (n *ScriptForStmt) Tag() NodeTag { return T_ScriptForStmt }

// ScriptWhileStmt represents `WHILE ( <cond> ) { DO | LOOP } ... END { WHILE | LOOP } [ <label> ]`.
type ScriptWhileStmt struct {
	Cond  Node
	Body  []Node
	Label Ident // optional label after END; zero Ident when absent
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptWhileStmt) Tag() NodeTag { return T_ScriptWhileStmt }

// ScriptRepeatStmt represents `REPEAT ... UNTIL ( <cond> ) END REPEAT [ <label> ]`.
type ScriptRepeatStmt struct {
	Body  []Node
	Cond  Node  // the UNTIL condition
	Label Ident // optional label after END; zero Ident when absent
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptRepeatStmt) Tag() NodeTag { return T_ScriptRepeatStmt }

// ScriptLoopStmt represents `LOOP ... END LOOP [ <label> ]`.
type ScriptLoopStmt struct {
	Body  []Node
	Label Ident // optional label after END; zero Ident when absent
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptLoopStmt) Tag() NodeTag { return T_ScriptLoopStmt }

// ScriptBreakStmt represents `BREAK [ <label> ]` (alias EXIT). Label is the
// zero Ident when no label follows.
type ScriptBreakStmt struct {
	Label Ident
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptBreakStmt) Tag() NodeTag { return T_ScriptBreakStmt }

// ScriptContinueStmt represents `CONTINUE [ <label> ]` (alias ITERATE). Label
// is the zero Ident when no label follows.
type ScriptContinueStmt struct {
	Label Ident
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptContinueStmt) Tag() NodeTag { return T_ScriptContinueStmt }

// ScriptReturnStmt represents `RETURN [ <expr> ]`. Value is nil for a bare
// RETURN.
type ScriptReturnStmt struct {
	Value Node
	Loc   Loc
}

// Tag implements Node.
func (n *ScriptReturnStmt) Tag() NodeTag { return T_ScriptReturnStmt }

// ScriptOpenStmt represents `OPEN <cursor> [ USING ( <bind> [ , ... ] ) ]`.
type ScriptOpenStmt struct {
	Cursor Ident
	Using  []Node // optional USING ( ... ) bind expressions; nil when absent
	Loc    Loc
}

// Tag implements Node.
func (n *ScriptOpenStmt) Tag() NodeTag { return T_ScriptOpenStmt }

// ScriptFetchStmt represents `FETCH <cursor> INTO <var> [ , ... ]`.
type ScriptFetchStmt struct {
	Cursor Ident
	Into   []Ident // target variables
	Loc    Loc
}

// Tag implements Node.
func (n *ScriptFetchStmt) Tag() NodeTag { return T_ScriptFetchStmt }

// ScriptCloseStmt represents `CLOSE <cursor>`.
type ScriptCloseStmt struct {
	Cursor Ident
	Loc    Loc
}

// Tag implements Node.
func (n *ScriptCloseStmt) Tag() NodeTag { return T_ScriptCloseStmt }

// ScriptRaiseStmt represents `RAISE [ <exc_name> ]`. Name is the zero Ident for
// a bare `RAISE` (re-raise the current exception inside a handler).
type ScriptRaiseStmt struct {
	Name Ident
	Loc  Loc
}

// Tag implements Node.
func (n *ScriptRaiseStmt) Tag() NodeTag { return T_ScriptRaiseStmt }

// Compile-time assertions for Snowflake Scripting nodes (T7.1).
var (
	_ Node = (*ScriptDeclaration)(nil)
	_ Node = (*ScriptExceptionHandler)(nil)
	_ Node = (*ScriptIfBranch)(nil)
	_ Node = (*ScriptCaseWhen)(nil)
	_ Node = (*ScriptBlockStmt)(nil)
	_ Node = (*ScriptAssignStmt)(nil)
	_ Node = (*ScriptLetStmt)(nil)
	_ Node = (*ScriptIfStmt)(nil)
	_ Node = (*ScriptCaseStmt)(nil)
	_ Node = (*ScriptForStmt)(nil)
	_ Node = (*ScriptWhileStmt)(nil)
	_ Node = (*ScriptRepeatStmt)(nil)
	_ Node = (*ScriptLoopStmt)(nil)
	_ Node = (*ScriptBreakStmt)(nil)
	_ Node = (*ScriptContinueStmt)(nil)
	_ Node = (*ScriptReturnStmt)(nil)
	_ Node = (*ScriptOpenStmt)(nil)
	_ Node = (*ScriptFetchStmt)(nil)
	_ Node = (*ScriptCloseStmt)(nil)
	_ Node = (*ScriptRaiseStmt)(nil)
)
