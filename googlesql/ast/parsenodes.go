package ast

// This file holds the concrete GoogleSQL parse-tree node types. The ast-core
// foundation ships only the File root container; later migration nodes
// (identifiers, types, expressions, SELECT core, joins, set-ops, DML, DDL,
// etc.) populate the rest, following the ZetaSQL-shaped tree.
//
// The cmd/genwalker code generator scans this file together with node.go to
// produce walk_generated.go.

// File is the root node of a parsed GoogleSQL source file. It holds the
// top-level statement list and the byte range covering the entire file.
// The parser entry point returns *File from Parse.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (f *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// DCL — GRANT / REVOKE (parser-dcl node)
// ---------------------------------------------------------------------------
//
// These node types model the legacy ANTLR grant_statement / revoke_statement
// family verbatim (GoogleSQLParser.g4 §2.8), which is itself a hand-port of
// Google's open-source ZetaSQL reference and the grammar bytebase consumes
// today. The grammar is:
//
//	grant_statement:  GRANT  privileges ON (identifier identifier?)? path_expression TO   grantee_list
//	revoke_statement: REVOKE privileges ON (identifier identifier?)? path_expression FROM grantee_list
//	privileges:       ALL PRIVILEGES? | privilege_list
//	privilege_list:   privilege (',' privilege)*
//	privilege:        privilege_name ('(' path_expression_list ')')?
//	privilege_name:   identifier | SELECT
//	grantee_list:     string_literal_or_parameter (',' string_literal_or_parameter)*
//
// The canonical ZetaSQL corpus (parser/googlesql/examples/.../grant_and_revoke.sql)
// is the accept/reject oracle for these forms. NOTE — DCL is a dialect-divergent
// zone: the Spanner emulator speaks a DIFFERENT GRANT dialect (`GRANT priv ON
// TABLE x TO ROLE r`, role-name grantees) and is therefore NON-AUTHORITATIVE
// here; see the parser-dcl divergence ledger entry. These nodes follow the
// ZetaSQL/.g4 grammar, not the Spanner emulator.
//
// Object names / privilege-column paths are kept as the lightweight NamePath
// value type (below) rather than a Node so they do NOT pre-empt the canonical
// path-expression node the expressions DAG node will own; DCL is self-contained
// (it depends only on the parser foundation).

// NamePath is a dotted identifier path (the grammar's path_expression:
// identifier ('.' identifier)*), used for a GRANT/REVOKE object name and for the
// optional per-privilege column list. It is a plain value type, NOT a Node: the
// walker does not descend into it and it intentionally does not collide with the
// canonical path-expression node owned by the expressions node.
//
// Parts holds each name component as the lexer surfaced it: a `backtick`-quoted
// identifier already has its backticks stripped, an unquoted identifier or
// keyword-as-identifier is its source word. (The lexer normalizes backtick
// identifiers the same way the legacy bigquery/spanner query-span extractors
// do.) Callers that need the original source slice can use Loc; DCL is not a
// query-span consumer, so no raw spelling is retained.
type NamePath struct {
	Parts []string // component names (>= 1)
	Loc   Loc
}

// String renders the path by joining its normalized parts with '.'. It does NOT
// re-quote; it is a debug/representation helper, not a deparser.
func (p NamePath) String() string {
	out := ""
	for i, part := range p.Parts {
		if i > 0 {
			out += "."
		}
		out += part
	}
	return out
}

// GranteeKind discriminates the four shapes of the grammar's
// string_literal_or_parameter (the entries of a grantee_list).
type GranteeKind int

const (
	// GranteeString is a string-literal grantee, e.g. 'user@google.com'.
	// Value holds the unquoted string content.
	GranteeString GranteeKind = iota
	// GranteeNamedParameter is a named query parameter, e.g. @user1 (the
	// grammar's named_parameter_expression: '@' identifier). Value holds the
	// parameter name without the leading '@'.
	GranteeNamedParameter
	// GranteePositionalParameter is a '?' positional query parameter (the
	// grammar's QUESTION). Value is empty.
	GranteePositionalParameter
	// GranteeSystemVariable is a system variable, e.g. @@user2 (the grammar's
	// system_variable_expression: '@@' path_expression). Value holds the
	// dotted variable path without the leading '@@'.
	GranteeSystemVariable
)

// String returns a human-readable name for the grantee kind.
func (k GranteeKind) String() string {
	switch k {
	case GranteeString:
		return "STRING"
	case GranteeNamedParameter:
		return "NAMED_PARAMETER"
	case GranteePositionalParameter:
		return "POSITIONAL_PARAMETER"
	case GranteeSystemVariable:
		return "SYSTEM_VARIABLE"
	default:
		return "UNKNOWN"
	}
}

// Grantee is one entry of a GRANT's TO list or a REVOKE's FROM list — the
// grammar's string_literal_or_parameter. Kind selects the shape; Value carries
// the payload (string content / parameter name / system-variable path), empty
// for a positional '?' parameter.
type Grantee struct {
	Kind  GranteeKind
	Value string
	Loc   Loc
}

// Tag implements Node.
func (n *Grantee) Tag() NodeTag { return T_Grantee }

// Privilege is one entry of a privilege_list: a privilege name with an optional
// parenthesized column list (the grammar's privilege:
// privilege_name path_expression_list_with_parens?). It is used only when a
// GRANT/REVOKE names explicit privileges; an `ALL PRIVILEGES` grant carries no
// Privilege nodes (see GrantStmt.AllPrivileges).
//
// Name is the privilege name (privilege_name: identifier | SELECT), as the lexer
// surfaced it (backtick identifiers unquoted). Columns is the optional
// per-privilege column list (e.g. the `(col1, col2)` in `insert(col1, col2)`);
// it is nil when no parentheses are present, and otherwise has >= 1 entry (the
// grammar's path_expression_list requires at least one path).
type Privilege struct {
	Name    string
	Columns []NamePath // nil when no '(' ... ')' column list was given
	Loc     Loc
}

// Tag implements Node.
func (n *Privilege) Tag() NodeTag { return T_Privilege }

// GrantStmt is a GRANT statement (grammar: grant_statement). It is the union of
// the documented forms; AllPrivileges selects ALL [PRIVILEGES] vs an explicit
// Privileges list.
//
//	GRANT { ALL [PRIVILEGES] | priv [, priv ...] }
//	  ON [ <object_type> [ <object_subtype> ] ] <path>
//	  TO <grantee> [, <grantee> ...]
//
// ObjectType holds the optional 0/1/2 object-type words (the grammar's
// (identifier identifier?)?), e.g. [] / ["table"] / ["materialized","view"].
// Path is the object name (a dotted path_expression). Grantees is the non-empty
// TO list.
type GrantStmt struct {
	AllPrivileges bool         // ALL [PRIVILEGES]
	Privileges    []*Privilege // explicit privilege list; nil when AllPrivileges
	ObjectType    []string     // 0, 1, or 2 object-type words
	Path          NamePath     // the ON object name
	Grantees      []*Grantee   // the TO recipients (>= 1)
	Loc           Loc
}

// Tag implements Node.
func (n *GrantStmt) Tag() NodeTag { return T_GrantStmt }

// RevokeStmt is a REVOKE statement (grammar: revoke_statement). It mirrors
// GrantStmt with a FROM grantee list instead of TO:
//
//	REVOKE { ALL [PRIVILEGES] | priv [, priv ...] }
//	  ON [ <object_type> [ <object_subtype> ] ] <path>
//	  FROM <grantee> [, <grantee> ...]
type RevokeStmt struct {
	AllPrivileges bool         // ALL [PRIVILEGES]
	Privileges    []*Privilege // explicit privilege list; nil when AllPrivileges
	ObjectType    []string     // 0, 1, or 2 object-type words
	Path          NamePath     // the ON object name
	Grantees      []*Grantee   // the FROM subjects (>= 1)
	Loc           Loc
}

// Tag implements Node.
func (n *RevokeStmt) Tag() NodeTag { return T_RevokeStmt }

// Compile-time assertions that the DCL node types satisfy Node.
var (
	_ Node = (*GrantStmt)(nil)
	_ Node = (*RevokeStmt)(nil)
	_ Node = (*Privilege)(nil)
	_ Node = (*Grantee)(nil)
)

// ===========================================================================
// Expressions (googlesql/expressions node)
// ===========================================================================
//
// The GoogleSQL expression tree. These node types model the legacy ANTLR
// expression grammar (GoogleSQLParser.g4 §2.17 precedence chain, §2.18
// constructors/CASE/CAST/EXTRACT/function calls), itself a hand-port of
// Google's open-source ZetaSQL reference and the grammar bytebase consumes
// today. The parser (parser/expr.go) builds these via precedence climbing.
//
// Adjudication: the LIVE Cloud Spanner emulator oracle (oracle.md) decides
// accept/reject for every shared GoogleSQL-core expression form; the ZetaSQL
// .g4 is the structural hint. Two oracle-confirmed precedence facts are baked
// into the parser and reflected in these shapes:
//
//	P1 (comparison non-associativity). The whole comparison family —
//	   comparative operators (= != <> < <= > >=), [NOT] LIKE, [NOT] IN,
//	   [NOT] BETWEEN, IS [NOT] {NULL|TRUE|FALSE|UNKNOWN}, and IS [NOT] DISTINCT
//	   FROM — is NON-ASSOCIATIVE: `a = b = c`, `1 < 2 < 3`, `x IN (1) IN (2)`,
//	   `1 IS NULL IS NULL`, `x LIKE 'a' LIKE 'b'`, `1 BETWEEN 0 AND 2 BETWEEN ..`
//	   all REJECT (oracle: "Syntax error: ..." / "Expression to the left of IS
//	   must be parenthesized"). The operands of these nodes are therefore the
//	   next-higher precedence level, never another comparison.
//
//	P2 (NOT binds looser than comparison). `NOT a = b` parses as `NOT (a = b)`
//	   and `NOT a IS NULL` as `NOT (a IS NULL)` (oracle: both accept). Prefix
//	   NOT's operand spans the whole comparison level (it sits between the
//	   comparison family and AND in precedence).
//
// Subquery seam: GoogleSQL subquery-bearing expressions — `(SELECT …)`,
// `EXISTS(…)`, `ARRAY(…)`, and IN/comparison RHS subqueries — are recognized
// here so the expression grammar accepts every oracle-accepted form, but the
// INNER query is NOT parsed by this node: the query/SELECT grammar belongs to
// the downstream parser-select node (which depends on this one). The subquery
// nodes (SubqueryExpr / ExistsExpr / ArraySubqueryExpr) capture the balanced
// parenthesized source span in RawText and leave Query nil; parser-select fills
// Query when it lands. See parser/expr.go parseParenOrSubquery.

// BinaryOp enumerates GoogleSQL binary operators that build a left-associative
// BinaryExpr (arithmetic, bitwise, shift, logical, concat). The comparison
// family is modeled separately (CompareExpr / IsExpr / InExpr / BetweenExpr /
// LikeExpr) because it is non-associative (P1).
type BinaryOp int

const (
	BinAdd        BinaryOp = iota // +
	BinSub                        // -
	BinMul                        // *
	BinDiv                        // /
	BinConcat                     // || (BOOL_OR_SYMBOL — string concat / logical-or-style)
	BinBitOr                      // |  (STROKE_SYMBOL)
	BinBitXor                     // ^  (CIRCUMFLEX_SYMBOL)
	BinBitAnd                     // &  (BIT_AND_SYMBOL)
	BinShiftLeft                  // <<
	BinShiftRight                 // >>
	BinAnd                        // AND
	BinOr                         // OR
)

// String returns the operator symbol/keyword.
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
	case BinConcat:
		return "||"
	case BinBitOr:
		return "|"
	case BinBitXor:
		return "^"
	case BinBitAnd:
		return "&"
	case BinShiftLeft:
		return "<<"
	case BinShiftRight:
		return ">>"
	case BinAnd:
		return "AND"
	case BinOr:
		return "OR"
	default:
		return "?"
	}
}

// UnaryOp enumerates GoogleSQL prefix unary operators (unary_operator + NOT).
type UnaryOp int

const (
	UnaryPlus   UnaryOp = iota // +
	UnaryMinus                 // -
	UnaryBitNot                // ~ (BITWISE_NOT_OPERATOR)
	UnaryNot                   // NOT
)

// String returns the operator symbol/keyword.
func (op UnaryOp) String() string {
	switch op {
	case UnaryPlus:
		return "+"
	case UnaryMinus:
		return "-"
	case UnaryBitNot:
		return "~"
	case UnaryNot:
		return "NOT"
	default:
		return "?"
	}
}

// CompareOp enumerates the comparative operators (comparative_operator). These
// build a non-associative CompareExpr.
type CompareOp int

const (
	CmpEq CompareOp = iota // =
	CmpNe                  // != or <>
	CmpLt                  // <
	CmpLe                  // <=
	CmpGt                  // >
	CmpGe                  // >=
)

// String returns the operator symbol.
func (op CompareOp) String() string {
	switch op {
	case CmpEq:
		return "="
	case CmpNe:
		return "!="
	case CmpLt:
		return "<"
	case CmpLe:
		return "<="
	case CmpGt:
		return ">"
	case CmpGe:
		return ">="
	default:
		return "?"
	}
}

// LiteralKind classifies a Literal's value type, mirroring the grammar's
// literal alternatives (§2.21) that carry no type keyword. Typed-prefix
// literals (NUMERIC '…', DATE '…', JSON '…', RANGE<…> '…') are modeled by
// TypedLiteral, not Literal.
type LiteralKind int

const (
	// LitNull is the NULL literal (null_literal).
	LitNull LiteralKind = iota
	// LitBool is TRUE / FALSE (boolean_literal). Value is the source spelling.
	LitBool
	// LitInt is an INTEGER_LITERAL (decimal or 0x-hex). Ival holds the value
	// (0 on int64 overflow — Value preserves the exact spelling).
	LitInt
	// LitFloat is a FLOATING_POINT_LITERAL. Value holds the source spelling.
	LitFloat
	// LitString is a STRING_LITERAL (one or more adjacent components, already
	// concatenated by the parser). Value holds the unquoted content.
	LitString
	// LitBytes is a BYTES_LITERAL. Value holds the unquoted, un-prefixed bytes.
	LitBytes
)

// String returns a human-readable name for the kind.
func (k LiteralKind) String() string {
	switch k {
	case LitNull:
		return "NULL"
	case LitBool:
		return "BOOL"
	case LitInt:
		return "INT"
	case LitFloat:
		return "FLOAT"
	case LitString:
		return "STRING"
	case LitBytes:
		return "BYTES"
	default:
		return "UNKNOWN"
	}
}

// ---------------------------------------------------------------------------
// Primary / leaf expression nodes
// ---------------------------------------------------------------------------

// Identifier is a single GoogleSQL identifier used as an expression (the
// grammar's identifier as a primary expression — a column / variable / name
// reference). Name is the normalized text: a backtick-quoted identifier has its
// backticks stripped by the lexer; an unquoted identifier or keyword-as-
// identifier is its source word.
//
// A multi-part dotted name (`a.b.c`) is a PathExpr, not an Identifier — but a
// dotted continuation in an expression is built as repeated FieldAccess over a
// leading Identifier (matching the grammar's `EHP . identifier` left-recursion),
// so a bare Identifier is always a single component here.
type Identifier struct {
	Name string
	Loc  Loc
}

// Tag implements Node.
func (n *Identifier) Tag() NodeTag { return T_Identifier }

// PathExpr is a dotted identifier path used where the grammar requires a
// path_expression as a whole (function names, sequence args, generalized-
// extension `(path)`, system-variable paths). Parts holds each normalized
// component (>= 1). General dotted field access on an arbitrary expression is
// FieldAccess, not PathExpr; PathExpr is used only in the grammar positions that
// take a path_expression token-sequence directly.
type PathExpr struct {
	Parts []string
	Loc   Loc
}

// Tag implements Node.
func (n *PathExpr) Tag() NodeTag { return T_PathExpr }

// String renders the path by joining its normalized parts with '.'. It does NOT
// re-quote; it is a debug/representation helper, not a deparser.
func (n *PathExpr) String() string {
	out := ""
	for i, part := range n.Parts {
		if i > 0 {
			out += "."
		}
		out += part
	}
	return out
}

// TypeRef is a reference to a parsed GoogleSQL data type in an expression
// position (CAST target, ARRAY<T>[…] element type, STRUCT<…>(…) / NEW type /
// braced-struct type). The `type` grammar itself is owned by the parser's
// `types` node (parser.DataType), which is a parser-package value type and not
// an ast.Node — so to keep the ast package free of a parser dependency, TypeRef
// stores only the rendered type text (parser.DataType.String(), which
// round-trips) plus its source span. Downstream consumers that need the
// structured type re-parse Text via parser.ParseDataType; bytebase's query-span
// path does not inspect types, so the rendered text is sufficient here.
type TypeRef struct {
	Text string // rendered type (e.g. "INT64", "ARRAY<STRUCT<x INT64>>")
	Loc  Loc
}

// Tag implements Node.
func (n *TypeRef) Tag() NodeTag { return T_TypeRef }

// StarExpr is a `*` used as an expression — only valid as a bare function
// argument (the grammar's MULTIPLY_OPERATOR alternative in
// function_call_expression_with_clauses_suffix, e.g. COUNT(*)). It is NOT a
// general expression primary; the parser admits it only in argument position.
// Modifiers carries an optional EXCEPT/REPLACE star_modifiers (rare in
// argument position; nil normally).
type StarExpr struct {
	Modifiers *StarModifiers // optional; nil for a bare *
	Loc       Loc
}

// Tag implements Node.
func (n *StarExpr) Tag() NodeTag { return T_StarExpr }

// StarModifiers carries the EXCEPT(...) / REPLACE(...) suffixes that can follow
// a star (star_modifiers). Except holds the EXCEPT column names; Replace holds
// the REPLACE items (each `expr AS name`). This node is defined here for the
// argument-position StarExpr; the SELECT-list star reuses it.
type StarModifiers struct {
	Except  []string           // EXCEPT (a, b, ...) column names
	Replace []*StructFieldExpr // REPLACE (expr AS name, ...) — reuses the aliased-expr shape
	Loc     Loc
}

// Tag implements Node.
func (n *StarModifiers) Tag() NodeTag { return T_StarModifiers }

// Literal is a NULL / boolean / integer / float / string / bytes literal that
// carries no type-name prefix (§2.21). Value is the source spelling (unquoted
// for string/bytes); Ival is the integer value for LitInt.
type Literal struct {
	Kind  LiteralKind
	Value string
	Ival  int64 // for LitInt
	Loc   Loc
}

// Tag implements Node.
func (n *Literal) Tag() NodeTag { return T_Literal }

// TypedLiteral is a type-prefixed literal: NUMERIC/DECIMAL/BIGNUMERIC/BIGDECIMAL
// '…', JSON '…', DATE/TIME/DATETIME/TIMESTAMP '…', or RANGE<type> '…'
// (numeric_literal / bignumeric_literal / json_literal / date_or_time_literal /
// range_literal). TypeKeyword is the leading type word (e.g. "NUMERIC", "DATE",
// "RANGE"); RangeType is the parsed element type for a RANGE literal (nil
// otherwise); Value is the unquoted string body.
type TypedLiteral struct {
	TypeKeyword string
	Value       string // the unquoted string-literal body
	Loc         Loc
}

// Tag implements Node.
func (n *TypedLiteral) Tag() NodeTag { return T_TypedLiteral }

// IntervalExpr is an interval value expression:
// `INTERVAL expr datepart [TO datepart]` (interval_expression). Value is the
// magnitude expression; From / To are the date-part identifier spellings (To is
// "" for the single-part form).
type IntervalExpr struct {
	Value Node
	From  string // leading datepart (e.g. "DAY", "YEAR")
	To    string // trailing datepart for `… TO …`; "" if absent
	Loc   Loc
}

// Tag implements Node.
func (n *IntervalExpr) Tag() NodeTag { return T_IntervalExpr }

// Parameter is a query parameter: a named parameter `@name`
// (named_parameter_expression) or a positional `?` (QUESTION_SYMBOL). Name holds
// the parameter name for the named form, "" for positional.
type Parameter struct {
	Name       string // "" for positional '?'
	Positional bool   // true for '?'
	Loc        Loc
}

// Tag implements Node.
func (n *Parameter) Tag() NodeTag { return T_Parameter }

// SystemVariable is a system-variable reference `@@path`
// (system_variable_expression: ATAT path_expression). Parts holds the dotted
// path components after the `@@` (>= 1).
type SystemVariable struct {
	Parts []string
	Loc   Loc
}

// Tag implements Node.
func (n *SystemVariable) Tag() NodeTag { return T_SystemVariable }

// ParenExpr is a parenthesized expression `( expr )`
// (parenthesized_expression_not_a_query). A parenthesized QUERY is a
// SubqueryExpr instead; the parser disambiguates by the token after '('.
type ParenExpr struct {
	Expr Node
	Loc  Loc
}

// Tag implements Node.
func (n *ParenExpr) Tag() NodeTag { return T_ParenExpr }

// ---------------------------------------------------------------------------
// Operator nodes
// ---------------------------------------------------------------------------

// UnaryExpr is a prefix unary operation `op expr` (unary_operator EHP, or
// NOT EHP). For NOT, Expr spans the whole comparison level (P2).
type UnaryExpr struct {
	Op   UnaryOp
	Expr Node
	Loc  Loc
}

// Tag implements Node.
func (n *UnaryExpr) Tag() NodeTag { return T_UnaryExpr }

// BinaryExpr is a left-associative binary operation `left op right` for the
// arithmetic / bitwise / shift / logical / concat operators. (The comparison
// family is non-associative and uses CompareExpr/IsExpr/etc., not BinaryExpr.)
type BinaryExpr struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Loc   Loc
}

// Tag implements Node.
func (n *BinaryExpr) Tag() NodeTag { return T_BinaryExpr }

// CompareExpr is a single comparative operation `left op right` for
// = != <> < <= > >= (comparative_operator). NON-ASSOCIATIVE (P1): Left and
// Right are both the next-higher precedence level, never another comparison.
type CompareExpr struct {
	Op    CompareOp
	Left  Node
	Right Node
	Loc   Loc
}

// Tag implements Node.
func (n *CompareExpr) Tag() NodeTag { return T_CompareExpr }

// IsExpr is `expr IS [NOT] {NULL | TRUE | FALSE | UNKNOWN}` or
// `expr IS [NOT] DISTINCT FROM expr` (is_operator / distinct_operator). Not
// flags the NOT. For the IS-predicate form, Pred holds the predicate keyword
// ("NULL", "TRUE", "FALSE", "UNKNOWN") and DistinctFrom is nil. For the DISTINCT
// FROM form, DistinctFrom is the right operand and Pred is "".
type IsExpr struct {
	Expr         Node
	Not          bool
	Pred         string // "NULL" | "TRUE" | "FALSE" | "UNKNOWN"; "" for DISTINCT FROM
	DistinctFrom Node   // non-nil for IS [NOT] DISTINCT FROM expr
	Loc          Loc
}

// Tag implements Node.
func (n *IsExpr) Tag() NodeTag { return T_IsExpr }

// InExpr is `expr [NOT] IN <rhs>` (in_operator). Exactly one RHS form is set:
// Values (a parenthesized expression list — one or more), Unnest (an
// UNNEST(...) call), or Subquery (a parenthesized query). Hint carries an
// optional `@{...}` between IN and a non-UNNEST RHS (the grammar rejects a hint
// before UNNEST).
type InExpr struct {
	Expr     Node
	Not      bool
	Values   []Node // parenthesized expression list (>= 1); nil for the other forms
	Unnest   Node   // UNNEST(...) RHS; nil otherwise
	Subquery Node   // parenthesized-query RHS (SubqueryExpr); nil otherwise
	Loc      Loc
}

// Tag implements Node.
func (n *InExpr) Tag() NodeTag { return T_InExpr }

// BetweenExpr is `expr [NOT] BETWEEN low AND high` (between_operator). low and
// high are at the EHP (higher-than-AND) level, so a bare OR inside is rejected
// (the grammar's "Expression in BETWEEN must be parenthesized" alt).
type BetweenExpr struct {
	Expr Node
	Low  Node
	High Node
	Not  bool
	Loc  Loc
}

// Tag implements Node.
func (n *BetweenExpr) Tag() NodeTag { return T_BetweenExpr }

// LikeExpr is `expr [NOT] LIKE <rhs>` (like_operator). The plain form sets
// Pattern. The quantified form (`LIKE ANY|SOME|ALL ...`) sets Quantifier and one
// of QuantValues (parenthesized list), QuantUnnest (UNNEST), or QuantSubquery.
type LikeExpr struct {
	Expr          Node
	Not           bool
	Pattern       Node   // plain `LIKE pattern`; nil for the quantified form
	Quantifier    string // "ANY" | "SOME" | "ALL"; "" for the plain form
	QuantValues   []Node // quantified list RHS
	QuantUnnest   Node   // quantified UNNEST RHS
	QuantSubquery Node   // quantified parenthesized-query RHS
	Loc           Loc
}

// Tag implements Node.
func (n *LikeExpr) Tag() NodeTag { return T_LikeExpr }

// ---------------------------------------------------------------------------
// CASE / CAST / EXTRACT
// ---------------------------------------------------------------------------

// CaseExpr is a CASE expression (case_expression). Operand is non-nil for the
// simple form (`CASE x WHEN …`) and nil for the searched form (`CASE WHEN …`).
// Whens has one or more branches; Else is the optional ELSE result (nil if
// absent).
type CaseExpr struct {
	Operand Node          // nil for searched CASE
	Whens   []*WhenClause // >= 1
	Else    Node          // nil if no ELSE
	Loc     Loc
}

// Tag implements Node.
func (n *CaseExpr) Tag() NodeTag { return T_CaseExpr }

// WhenClause is one `WHEN cond THEN result` branch of a CaseExpr. For a simple
// CASE, Cond is the value compared against the operand.
type WhenClause struct {
	Cond   Node
	Result Node
	Loc    Loc
}

// Tag implements Node.
func (n *WhenClause) Tag() NodeTag { return T_WhenClause }

// CastExpr is `CAST(expr AS type [FORMAT fmt [AT TIME ZONE tz]])` or the
// `SAFE_CAST(...)` variant (cast_expression). Safe flags SAFE_CAST. Type is the
// target type. Format / TimeZone are the optional FORMAT and AT TIME ZONE
// expressions (nil if absent). Type is carried as a Node wrapper produced by the
// parser (the parser embeds its *DataType in a typeNode); downstream consumers
// read the rendered type via the parser, so Type is opaque here.
type CastExpr struct {
	Expr     Node
	Type     *TypeRef // the target type (rendered); nil only on a malformed parse
	Safe     bool     // SAFE_CAST
	Format   Node     // optional FORMAT expr; nil if absent
	TimeZone Node     // optional AT TIME ZONE expr (only with FORMAT); nil if absent
	Loc      Loc
}

// Tag implements Node.
func (n *CastExpr) Tag() NodeTag { return T_CastExpr }

// ExtractExpr is `EXTRACT(part FROM expr [AT TIME ZONE tz])`
// (extract_expression). Part is the part expression (e.g. YEAR, or
// DATE(...)); From is the source; TimeZone is the optional AT TIME ZONE expr.
type ExtractExpr struct {
	Part     Node
	From     Node
	TimeZone Node // optional; nil if absent
	Loc      Loc
}

// Tag implements Node.
func (n *ExtractExpr) Tag() NodeTag { return T_ExtractExpr }

// ---------------------------------------------------------------------------
// Function calls, arguments, window
// ---------------------------------------------------------------------------

// FuncCall is a function-call expression with all its clauses
// (function_call_expression_with_clauses + suffix). Name is the function name
// (a path, or a keyword-as-function-name like IF/GROUPING/LEFT/RIGHT/COLLATE/
// RANGE rendered into Name.Parts). The modifier fields mirror the suffix grammar.
type FuncCall struct {
	Name          *PathExpr        // function name path
	Args          []Node           // positional / named / lambda / sequence args; a bare '*' is a StarExpr
	Distinct      bool             // DISTINCT before the args
	NullHandling  string           // "IGNORE NULLS" | "RESPECT NULLS"; "" if absent
	Having        *HavingModifier  // HAVING MAX/MIN expr [group-by]; nil if absent
	Clamped       *ClampedModifier // CLAMPED BETWEEN low AND high; nil if absent
	WithReport    bool             // WITH REPORT (...) present
	OrderBy       []*OrderItem     // trailing ORDER BY inside the call (aggregate ordering)
	Limit         Node             // trailing LIMIT expr; nil if absent
	LimitOffset   Node             // trailing OFFSET expr (only with Limit); nil if absent
	WithGroupRows bool             // WITH GROUP ROWS present
	Over          *WindowSpec      // OVER window; nil if not a window call
	Loc           Loc
}

// Tag implements Node.
func (n *FuncCall) Tag() NodeTag { return T_FuncCall }

// HavingModifier is the aggregate `HAVING MAX expr` / `HAVING MIN expr <group>`
// modifier (opt_having_or_group_by_modifier). IsMax selects MAX vs MIN; Expr is
// the having expression. (The MIN form's trailing group-by is rare and its items
// are not retained beyond Expr in this node — the modifier's presence is the
// load-bearing fact for parse parity.)
type HavingModifier struct {
	IsMax bool
	Expr  Node
	Loc   Loc
}

// Tag implements Node.
func (n *HavingModifier) Tag() NodeTag { return T_HavingModifier }

// ClampedModifier is the differential-privacy `CLAMPED BETWEEN low AND high`
// modifier (clamped_between_modifier).
type ClampedModifier struct {
	Low  Node
	High Node
	Loc  Loc
}

// Tag implements Node.
func (n *ClampedModifier) Tag() NodeTag { return T_ClampedModifier }

// NamedArg is a named function argument `name => value` (named_argument). Value
// is the argument expression (which may itself be a LambdaExpr).
type NamedArg struct {
	Name  string
	Value Node
	Loc   Loc
}

// Tag implements Node.
func (n *NamedArg) Tag() NodeTag { return T_NamedArg }

// LambdaExpr is a lambda function argument `params -> body` (lambda_argument).
// Params holds the parameter identifier spellings (empty for the `() -> body`
// form). Body is the lambda body expression.
type LambdaExpr struct {
	Params []string
	Body   Node
	Loc    Loc
}

// Tag implements Node.
func (n *LambdaExpr) Tag() NodeTag { return T_LambdaExpr }

// SequenceArg is the `SEQUENCE path` function argument (sequence_arg). Path is
// the sequence name path.
type SequenceArg struct {
	Path *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *SequenceArg) Tag() NodeTag { return T_SequenceArg }

// WindowSpec is an OVER window specification (over_clause + window_specification):
// either a named window reference (Name set, all else empty) or an inline spec
// `( [name] [PARTITION BY …] [ORDER BY …] [frame] )`.
type WindowSpec struct {
	Name        string       // named window reference, or the leading name of an inline spec; "" if none
	PartitionBy []Node       // PARTITION BY expressions
	OrderBy     []*OrderItem // ORDER BY items
	Frame       *WindowFrame // window frame; nil if absent
	Inline      bool         // true if the spec was a parenthesized inline spec (vs a bare name reference)
	Loc         Loc
}

// Tag implements Node.
func (n *WindowSpec) Tag() NodeTag { return T_WindowSpec }

// OrderItem is one `expr [COLLATE c] [ASC|DESC] [NULLS FIRST|LAST]` ordering
// element (ordering_expression), used in window ORDER BY and aggregate ORDER BY.
type OrderItem struct {
	Expr       Node
	Collate    Node  // optional COLLATE expr; nil if absent
	Desc       bool  // true for DESC
	HasDir     bool  // true if ASC or DESC was given explicitly
	NullsFirst *bool // nil = unspecified, true = NULLS FIRST, false = NULLS LAST
	Loc        Loc
}

// Tag implements Node.
func (n *OrderItem) Tag() NodeTag { return T_OrderItem }

// WindowFrameKind enumerates the frame unit (frame_unit).
type WindowFrameKind int

const (
	FrameRows  WindowFrameKind = iota // ROWS
	FrameRange                        // RANGE
)

// WindowBoundKind enumerates a frame bound shape (window_frame_bound).
type WindowBoundKind int

const (
	BoundUnboundedPreceding WindowBoundKind = iota // UNBOUNDED PRECEDING
	BoundUnboundedFollowing                        // UNBOUNDED FOLLOWING
	BoundCurrentRow                                // CURRENT ROW
	BoundPreceding                                 // expr PRECEDING
	BoundFollowing                                 // expr FOLLOWING
)

// WindowFrame is a `ROWS|RANGE { <bound> | BETWEEN <bound> AND <bound> }` clause
// (opt_window_frame_clause). Between is true for the BETWEEN form (Start + End);
// for the single-bound form only Start is meaningful.
type WindowFrame struct {
	Kind    WindowFrameKind
	Between bool
	Start   WindowBound
	End     WindowBound // meaningful only when Between
	Loc     Loc
}

// Tag implements Node.
func (n *WindowFrame) Tag() NodeTag { return T_WindowFrame }

// WindowBound is one end of a frame (window_frame_bound). Offset is the `expr`
// for BoundPreceding/BoundFollowing (nil otherwise). Not a Node.
type WindowBound struct {
	Kind   WindowBoundKind
	Offset Node // for BoundPreceding/BoundFollowing
}

// ---------------------------------------------------------------------------
// Constructors: array / struct / new / braced / replace-fields / with
// ---------------------------------------------------------------------------

// ArrayExpr is an array constructor: `[…]`, `ARRAY[…]`, or `ARRAY<T>[…]`
// (array_constructor). Elements holds zero or more element expressions; ElemType
// is the optional explicit element type (the T in ARRAY<T>[…]), wrapped as a
// Node by the parser (nil if absent). HasArrayKeyword records whether the ARRAY
// keyword was present (vs the bare `[…]` form).
type ArrayExpr struct {
	Elements        []Node
	ElemType        *TypeRef // explicit element type (ARRAY<T>[…]); nil otherwise
	HasArrayKeyword bool
	Loc             Loc
}

// Tag implements Node.
func (n *ArrayExpr) Tag() NodeTag { return T_ArrayExpr }

// StructExpr is a struct constructor (struct_constructor). Three source forms map
// here: the keyword form `STRUCT(args)` / `STRUCT<…>(args)` (HasStruct true,
// Type optionally set), and the bare tuple form `(e1, e2, …)` with two or more
// elements (HasStruct false, Type nil). Fields holds each `expr [AS alias]` arg.
type StructExpr struct {
	HasStruct bool               // STRUCT keyword present
	Type      *TypeRef           // explicit STRUCT<…> type; nil otherwise
	Fields    []*StructFieldExpr // the constructor args (>= 0 for STRUCT(), >= 2 for the bare tuple)
	Loc       Loc
}

// Tag implements Node.
func (n *StructExpr) Tag() NodeTag { return T_StructExpr }

// StructFieldExpr is one `expr [AS alias]` constructor argument
// (struct_constructor_arg) — reused for array/struct REPLACE items. Alias is ""
// when no `AS name` was given.
type StructFieldExpr struct {
	Value Node
	Alias string // "" if no AS alias
	Loc   Loc
}

// Tag implements Node.
func (n *StructFieldExpr) Tag() NodeTag { return T_StructFieldExpr }

// NewConstructor is a proto constructor `NEW type ( args )` (new_constructor) or
// the braced form `NEW type { … }` (braced_new_constructor). Type is the proto
// type name (wraps a parsed type_name as a Node). Args holds the parenthesized
// args; Braced holds the braced constructor (one of Args/Braced is set).
type NewConstructor struct {
	Type   *TypeRef           // the proto type_name (rendered)
	Args   []*StructFieldExpr // parenthesized form args (each `expr [AS id | AS (path)]`)
	Braced *BracedConstructor // braced form; nil for the parenthesized form
	Loc    Loc
}

// Tag implements Node.
func (n *NewConstructor) Tag() NodeTag { return T_NewConstructor }

// BracedConstructor is a proto braced constructor `{ field… }` (braced_constructor
// / struct_braced_constructor). Type is the optional leading STRUCT<…> / struct
// type for struct_braced_constructor (nil for the proto braced_constructor and
// the bare `STRUCT {…}` no-type form). The field internals are retained as raw
// field expressions in Fields (each a `path : value` or `path {…}`); the exact
// proto-field shape is not consumed by bytebase, so fields capture the value
// expressions for walk/coverage without a bespoke per-field node.
type BracedConstructor struct {
	Type   *TypeRef // optional struct type (struct_braced_constructor); nil otherwise
	Fields []Node   // field value expressions (best-effort capture for walking)
	Loc    Loc
}

// Tag implements Node.
func (n *BracedConstructor) Tag() NodeTag { return T_BracedConstructor }

// ReplaceFieldsExpr is `REPLACE_FIELDS(expr, value AS path, …)`
// (replace_fields_expression). Expr is the base; Items holds each replacement
// (`expr AS generalized-path`).
type ReplaceFieldsExpr struct {
	Expr  Node
	Items []*StructFieldExpr // each `value AS path` (Alias holds the rendered path)
	Loc   Loc
}

// Tag implements Node.
func (n *ReplaceFieldsExpr) Tag() NodeTag { return T_ReplaceFieldsExpr }

// WithExpr is the inline `WITH(name AS expr, …, body)` expression
// (with_expression). Vars holds the `name AS expr` bindings; Body is the final
// result expression.
type WithExpr struct {
	Vars []*StructFieldExpr // each `name AS expr` (Alias holds the var name, Value the bound expr)
	Body Node
	Loc  Loc
}

// Tag implements Node.
func (n *WithExpr) Tag() NodeTag { return T_WithExpr }

// ---------------------------------------------------------------------------
// Access (field / index / extension) and subqueries
// ---------------------------------------------------------------------------

// FieldAccess is dotted field access on an expression `expr . field`
// (EHP DOT identifier). Field is the accessed field name (normalized).
type FieldAccess struct {
	Expr  Node
	Field string
	Loc   Loc
}

// Tag implements Node.
func (n *FieldAccess) Tag() NodeTag { return T_FieldAccess }

// IndexAccess is array/struct subscript access `expr [ index ]`
// (EHP LS_BRACKET expression RS_BRACKET). Index is the subscript expression
// (which may itself be an OFFSET(...)/ORDINAL(...)/SAFE_OFFSET/SAFE_ORDINAL
// function call — those are ordinary FuncCall nodes, not special syntax).
type IndexAccess struct {
	Expr  Node
	Index Node
	Loc   Loc
}

// Tag implements Node.
func (n *IndexAccess) Tag() NodeTag { return T_IndexAccess }

// ExtensionAccess is proto-extension field access `expr . ( path )`
// (EHP DOT '(' path_expression ')'). Path is the parenthesized extension path.
type ExtensionAccess struct {
	Expr Node
	Path *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *ExtensionAccess) Tag() NodeTag { return T_ExtensionAccess }

// SubqueryExpr is a parenthesized query used as a scalar expression
// `( SELECT … )` (parenthesized_query in expression position). Query is the
// parsed inner query — left nil by the expressions node (the query grammar is
// owned by parser-select); RawText holds the inner query source (without the
// outer parens) so parser-select can fill Query later and so deparse/round-trip
// works in the interim.
type SubqueryExpr struct {
	Query   Node   // nil until parser-select wires the query grammar
	RawText string // inner query source (between the parens)
	Loc     Loc
}

// Tag implements Node.
func (n *SubqueryExpr) Tag() NodeTag { return T_SubqueryExpr }

// ExistsExpr is `EXISTS [hint] ( query )` (expression_subquery_with_keyword).
// Query is the parsed inner query (nil until parser-select); RawText holds the
// inner source.
type ExistsExpr struct {
	Query   Node
	RawText string
	Loc     Loc
}

// Tag implements Node.
func (n *ExistsExpr) Tag() NodeTag { return T_ExistsExpr }

// ArraySubqueryExpr is `ARRAY ( query )` (expression_subquery_with_keyword).
// Query is the parsed inner query (nil until parser-select); RawText holds the
// inner source.
type ArraySubqueryExpr struct {
	Query   Node
	RawText string
	Loc     Loc
}

// Tag implements Node.
func (n *ArraySubqueryExpr) Tag() NodeTag { return T_ArraySubqueryExpr }

// Compile-time assertions that the expression node types satisfy Node.
var (
	_ Node = (*Identifier)(nil)
	_ Node = (*PathExpr)(nil)
	_ Node = (*TypeRef)(nil)
	_ Node = (*StarExpr)(nil)
	_ Node = (*StarModifiers)(nil)
	_ Node = (*Literal)(nil)
	_ Node = (*TypedLiteral)(nil)
	_ Node = (*IntervalExpr)(nil)
	_ Node = (*Parameter)(nil)
	_ Node = (*SystemVariable)(nil)
	_ Node = (*ParenExpr)(nil)
	_ Node = (*UnaryExpr)(nil)
	_ Node = (*BinaryExpr)(nil)
	_ Node = (*CompareExpr)(nil)
	_ Node = (*IsExpr)(nil)
	_ Node = (*InExpr)(nil)
	_ Node = (*BetweenExpr)(nil)
	_ Node = (*LikeExpr)(nil)
	_ Node = (*CaseExpr)(nil)
	_ Node = (*WhenClause)(nil)
	_ Node = (*CastExpr)(nil)
	_ Node = (*ExtractExpr)(nil)
	_ Node = (*FuncCall)(nil)
	_ Node = (*NamedArg)(nil)
	_ Node = (*LambdaExpr)(nil)
	_ Node = (*SequenceArg)(nil)
	_ Node = (*WindowSpec)(nil)
	_ Node = (*WindowFrame)(nil)
	_ Node = (*OrderItem)(nil)
	_ Node = (*HavingModifier)(nil)
	_ Node = (*ClampedModifier)(nil)
	_ Node = (*ArrayExpr)(nil)
	_ Node = (*StructExpr)(nil)
	_ Node = (*StructFieldExpr)(nil)
	_ Node = (*NewConstructor)(nil)
	_ Node = (*BracedConstructor)(nil)
	_ Node = (*ReplaceFieldsExpr)(nil)
	_ Node = (*WithExpr)(nil)
	_ Node = (*FieldAccess)(nil)
	_ Node = (*IndexAccess)(nil)
	_ Node = (*ExtensionAccess)(nil)
	_ Node = (*SubqueryExpr)(nil)
	_ Node = (*ExistsExpr)(nil)
	_ Node = (*ArraySubqueryExpr)(nil)
)

// ===========================================================================
// Query / SELECT core (googlesql/parser-select node)
// ===========================================================================
//
// These node types model the legacy ANTLR query grammar (GoogleSQLParser.g4
// §2.13-§2.16 — query / query_without_pipe_operators / select / from_clause /
// joins / set-ops / with_clause), itself a hand-port of Google's open-source
// ZetaSQL reference and the grammar bytebase consumes today. The parser
// (parser/select.go, from_join.go, set_ops.go, cte.go) builds these.
//
// One omni parser serves both BigQuery and Spanner; this is the UNION grammar.
// The live Cloud Spanner emulator (oracle.md) decides accept/reject for shared
// + Spanner-only query forms; BigQuery-only forms (QUALIFY, WITH RECURSIVE,
// SELECT AS VALUE in some positions) parse in the union grammar even though the
// Spanner emulator *feature-rejects* them — the emulator's feature-reject is
// NON-AUTHORITATIVE for grammar verdicts (divergence ledger #9, #11).
//
// Grammar-faithful tree shape (a key ZetaSQL nuance): ORDER BY / LIMIT-OFFSET /
// FOR UPDATE are properties of the whole `query_without_pipe_operators`, NOT of
// the inner SELECT. So a *QueryStmt wraps the set-op/select Body and carries the
// trailing ORDER BY/LIMIT/FOR UPDATE; the inner *SelectStmt holds only
// SELECT…FROM…WHERE…GROUP BY…HAVING…QUALIFY…WINDOW. This makes
// `SELECT 1 UNION ALL SELECT 2 ORDER BY x` bind ORDER BY to the union (correct),
// not to the second SELECT.

// QueryStmt is a complete query (grammar: query / query_without_pipe_operators):
//
//	[WITH [RECURSIVE] cte, …] <body> [ORDER BY …] [LIMIT n [OFFSET m]] [FOR UPDATE]
//
// With is the optional leading WITH clause (nil if absent). Body is the query
// body — a *SelectStmt, a *SetOperation, or a nested *QueryStmt (a parenthesized
// `( query )` query_primary). OrderBy / Limit / Offset are the trailing
// query-level modifiers (nil/absent when not present). ForUpdate records a
// trailing `FOR UPDATE` (Spanner row-locking; parses in the union grammar).
type QueryStmt struct {
	With      *WithClause  // leading WITH clause; nil if absent
	Body      Node         // *SelectStmt | *SetOperation | *QueryStmt
	OrderBy   []*OrderItem // trailing query-level ORDER BY; nil if absent
	Limit     Node         // trailing LIMIT count; nil if absent
	Offset    Node         // trailing OFFSET (only with Limit); nil if absent
	ForUpdate bool         // trailing FOR UPDATE (Spanner)
	Parens    bool         // true if this query was wrapped in `( … )` as a query_primary
	Loc       Loc
}

// Tag implements Node.
func (n *QueryStmt) Tag() NodeTag { return T_QueryStmt }

// WithClause is the leading `WITH [RECURSIVE] aliased_query (, aliased_query)*`
// CTE clause (with_clause). Recursive flags the RECURSIVE keyword (parses in the
// union grammar; Spanner feature-rejects — divergence #11). CTEs has >= 1 entry.
type WithClause struct {
	Recursive bool
	CTEs      []*CTE // >= 1
	Loc       Loc
}

// Tag implements Node.
func (n *WithClause) Tag() NodeTag { return T_WithClause }

// CTE is one named subquery in a WITH clause (aliased_query):
//
//	name [( column, … )] AS ( query ) [WITH DEPTH …]
//
// Name is the CTE name (normalized). Columns is the optional explicit
// column-name list (Spanner's `cte_name ( column_name, … )` form; nil if
// absent). Query is the parsed CTE body (a *QueryStmt). Depth carries the
// optional recursion-depth modifier (`WITH DEPTH [AS alias] [BETWEEN n AND m |
// MAX n]`); nil if absent.
type CTE struct {
	Name    string
	Columns []string // explicit column list (Spanner form); nil if absent
	Query   Node     // *QueryStmt
	Depth   *RecursionDepth
	Loc     Loc
}

// Tag implements Node.
func (n *CTE) Tag() NodeTag { return T_CTE }

// RecursionDepth is the recursion-depth modifier on a recursive CTE
// (recursion_depth_modifier): `WITH DEPTH [AS alias] [BETWEEN low AND high |
// MAX n]`. Alias is the optional depth-column alias ("" if absent). Lower/Upper
// are the BETWEEN bounds (nil if the BETWEEN form was not used); Max is the
// `MAX n` bound (nil if not used). At most one of {Lower/Upper} or Max is set.
type RecursionDepth struct {
	Alias string
	Lower Node // BETWEEN low …; nil if absent
	Upper Node // … AND high; nil if absent
	Max   Node // MAX n; nil if absent
	Loc   Loc
}

// Tag implements Node.
func (n *RecursionDepth) Tag() NodeTag { return T_RecursionDepth }

// SelectAsKind enumerates the `SELECT AS …` modifier (opt_select_as_clause).
type SelectAsKind int

const (
	// SelectAsNone is no AS modifier.
	SelectAsNone SelectAsKind = iota
	// SelectAsStruct is `SELECT AS STRUCT`.
	SelectAsStruct
	// SelectAsValue is `SELECT AS VALUE`.
	SelectAsValue
	// SelectAsTypeName is `SELECT AS <path>` (a proto/type name); the name is in
	// SelectStmt.AsTypeName.
	SelectAsTypeName
)

// SelectStmt is a single SELECT block (grammar: select):
//
//	SELECT [hint] [WITH <dp/agg-threshold>] [ALL|DISTINCT] [AS STRUCT|VALUE|path]
//	  <select_list>
//	  [FROM <from>] [WHERE …] [GROUP BY …] [HAVING …] [QUALIFY …] [WINDOW …]
//
// ORDER BY / LIMIT / FOR UPDATE are NOT here — they belong to the enclosing
// QueryStmt (see that type). Distinct/All select the set quantifier. As selects
// the `AS STRUCT|VALUE|path` modifier; AsTypeName holds the path for
// SelectAsTypeName. SelectWith records a leading `WITH <name> [OPTIONS(…)]`
// differential-privacy / aggregation-threshold clause (its body is not retained
// beyond the name; the clause's presence + name is the load-bearing fact).
type SelectStmt struct {
	Distinct   bool
	All        bool
	As         SelectAsKind
	AsTypeName *PathExpr      // for SelectAsTypeName; nil otherwise
	SelectWith string         // leading WITH <name> clause name ("" if absent)
	Items      []*SelectItem  // the select list (>= 1)
	From       []Node         // FROM items: *TableExpr / *JoinExpr / *UnnestExpr; nil if no FROM
	Where      Node           // WHERE expr; nil if absent
	GroupBy    *GroupByClause // nil if absent
	Having     Node           // HAVING expr; nil if absent
	Qualify    Node           // QUALIFY expr; nil if absent
	Window     []*WindowDef   // WINDOW named windows; nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }

// SelectItem is one entry of a select_list (select_list_item). Three shapes:
//   - expression:  Expr set, Star false. Alias is the optional `[AS] name`.
//   - bare star:   `*` with optional modifiers. Star true, Expr nil.
//   - dot star:    `expr.*` with optional modifiers. Star true, Expr set.
//
// Modifiers carries the EXCEPT/REPLACE star_modifiers (nil for the expression
// form and for a star with no modifiers).
type SelectItem struct {
	Expr      Node           // the expression, or the dot-star qualifier; nil for bare `*`
	Star      bool           // true for `*` or `expr.*`
	Modifiers *StarModifiers // EXCEPT/REPLACE; nil if none
	Alias     string         // `[AS] alias` (expression form only); "" if absent
	Loc       Loc
}

// Tag implements Node.
func (n *SelectItem) Tag() NodeTag { return T_SelectItem }

// SetOp enumerates the set-operation kind (query_set_operation_type).
type SetOp int

const (
	SetOpUnion     SetOp = iota // UNION
	SetOpIntersect              // INTERSECT
	SetOpExcept                 // EXCEPT
)

// String returns the operator keyword.
func (op SetOp) String() string {
	switch op {
	case SetOpUnion:
		return "UNION"
	case SetOpIntersect:
		return "INTERSECT"
	case SetOpExcept:
		return "EXCEPT"
	default:
		return "?"
	}
}

// SetOperation is a set-operation node (query_set_operation):
//
//	<left> [corresponding-outer] <UNION|INTERSECT|EXCEPT> {ALL|DISTINCT}
//	   [STRICT] [CORRESPONDING [BY]] <right>
//
// Left-associative: `a UNION b UNION c` nests as
// SetOperation{Left: SetOperation{Left: a, Right: b}, Right: c}. Left and Right
// are query_primary nodes (*SelectStmt or a parenthesized *QueryStmt). Op is the
// operator; All flags ALL (vs DISTINCT). OuterMode records an optional
// corresponding-outer prefix ("" / "FULL" / "OUTER" / "LEFT"); Strict flags
// STRICT; Corresponding flags a trailing CORRESPONDING [BY]. The all-or-distinct
// choice is required by the grammar, so AllOrDistinctSet is always true for a
// well-formed node (kept explicit so a deparser can round-trip).
type SetOperation struct {
	Op            SetOp
	All           bool   // ALL (vs DISTINCT)
	OuterMode     string // "" | "FULL" | "OUTER" | "LEFT"
	Strict        bool
	Corresponding bool
	Left          Node // *SelectStmt | *QueryStmt
	Right         Node // *SelectStmt | *QueryStmt
	Loc           Loc
}

// Tag implements Node.
func (n *SetOperation) Tag() NodeTag { return T_SetOperation }

// ---------------------------------------------------------------------------
// FROM clause / table sources / joins
// ---------------------------------------------------------------------------

// TableExpr is a non-join FROM source (table_primary minus the join cases):
// a table path reference, a parenthesized subquery, a TVF call, or an implicit
// array-path source. Exactly one of {Path, Subquery, Func} is set.
//
//   - Path:     a (possibly dashed/slashed/dotted) table path, or a correlated
//     array field path (`t.array_col`). Held as a *PathExpr.
//   - Subquery: a parenthesized `( query )` table_subquery. Held as a *QueryStmt.
//   - Func:     a table-valued function call `name(args)`. Held as a *FuncCall.
//
// Alias is the optional `[AS] name` (alias text; "" if absent). WithOffset
// records a trailing `WITH OFFSET [[AS] name]` (the offset companion column);
// WithOffsetAlias is its optional alias. SystemTime captures a trailing
// `FOR SYSTEM[_TIME] AS OF expr` time-travel expression (nil if absent).
type TableExpr struct {
	Path            *PathExpr // table / array-field path; nil for subquery / TVF
	Subquery        Node      // ( query ) subquery; nil otherwise (*QueryStmt)
	Func            *FuncCall // table-valued function call; nil otherwise
	Alias           string    // [AS] alias; "" if absent
	WithOffset      bool      // trailing WITH OFFSET
	WithOffsetAlias string    // alias for WITH OFFSET; "" if absent or unaliased
	SystemTime      Node      // FOR SYSTEM_TIME AS OF expr; nil if absent
	Loc             Loc
}

// Tag implements Node.
func (n *TableExpr) Tag() NodeTag { return T_TableExpr }

// UnnestExpr is an `UNNEST(array_expr [, …]) [[AS] alias] [WITH OFFSET [[AS] name]]`
// FROM source (unnest_expression as a table_path_expression_base). Array is the
// UNNEST(...) call (a *FuncCall named UNNEST, as the expressions node builds it).
// Alias / WithOffset / WithOffsetAlias mirror TableExpr.
type UnnestExpr struct {
	Array           Node   // the UNNEST(...) call (*FuncCall); the array argument(s)
	Alias           string // [AS] alias; "" if absent
	WithOffset      bool   // trailing WITH OFFSET
	WithOffsetAlias string // alias for WITH OFFSET; "" if absent or unaliased
	Loc             Loc
}

// Tag implements Node.
func (n *UnnestExpr) Tag() NodeTag { return T_UnnestExpr }

// JoinType enumerates the join kind (join_type / opt_natural).
type JoinType int

const (
	JoinComma JoinType = iota // `,` (implicit cross join)
	JoinInner                 // [INNER] JOIN
	JoinCross                 // CROSS JOIN
	JoinFull                  // FULL [OUTER] JOIN
	JoinLeft                  // LEFT [OUTER] JOIN
	JoinRight                 // RIGHT [OUTER] JOIN
)

// String returns a human-readable name for the join type.
func (t JoinType) String() string {
	switch t {
	case JoinComma:
		return "COMMA"
	case JoinInner:
		return "INNER"
	case JoinCross:
		return "CROSS"
	case JoinFull:
		return "FULL"
	case JoinLeft:
		return "LEFT"
	case JoinRight:
		return "RIGHT"
	default:
		return "?"
	}
}

// JoinExpr is a join between two FROM sources (from_clause_contents_suffix /
// join_item). Left and Right are FROM sources (*TableExpr / *UnnestExpr /
// nested *JoinExpr / a parenthesized join). Type is the join kind. Natural flags
// a NATURAL prefix. JoinHint records the `HASH | LOOKUP` algorithm hint ("" if
// none). On is the ON condition (nil if USING/NATURAL/CROSS/comma). Using is the
// USING column list (nil if ON/NATURAL/CROSS/comma).
type JoinExpr struct {
	Type     JoinType
	Natural  bool
	JoinHint string // "HASH" | "LOOKUP"; "" if none
	Left     Node
	Right    Node
	On       Node     // ON expr; nil otherwise
	Using    []string // USING (col, …); nil otherwise
	Loc      Loc
}

// Tag implements Node.
func (n *JoinExpr) Tag() NodeTag { return T_JoinExpr }

// ---------------------------------------------------------------------------
// GROUP BY / WINDOW
// ---------------------------------------------------------------------------

// GroupByKind enumerates the GROUP BY shape (group_by_clause).
type GroupByKind int

const (
	// GroupByItems is `GROUP BY <grouping items>` (the general form; each item
	// is a GroupingItem).
	GroupByItems GroupByKind = iota
	// GroupByAll is `GROUP BY ALL`.
	GroupByAll
)

// GroupByClause is a GROUP BY clause (group_by_clause): either `GROUP BY ALL`
// (Kind == GroupByAll, Items nil) or `GROUP BY <grouping items>` (Kind ==
// GroupByItems, Items >= 1). AndOrder flags the `AND ORDER` infix
// (`GROUP AND ORDER BY`). Each item is a GroupingItem (a plain expression, the
// empty grouping `()`, or ROLLUP/CUBE/GROUPING SETS).
type GroupByClause struct {
	Kind     GroupByKind
	AndOrder bool            // GROUP AND ORDER BY
	Items    []*GroupingItem // nil for GROUP BY ALL
	Loc      Loc
}

// Tag implements Node.
func (n *GroupByClause) Tag() NodeTag { return T_GroupByClause }

// GroupingKind enumerates a grouping_item shape.
type GroupingKind int

const (
	// GroupingExpr is a plain `expr [AS alias] [ASC|DESC]` grouping item.
	GroupingExpr GroupingKind = iota
	// GroupingEmpty is the empty grouping `()`.
	GroupingEmpty
	// GroupingRollup is `ROLLUP ( expr, … )`.
	GroupingRollup
	// GroupingCube is `CUBE ( expr, … )`.
	GroupingCube
	// GroupingSets is `GROUPING SETS ( <set>, … )`.
	GroupingSets
)

// GroupingItem is one entry of a GROUP BY list (grouping_item). For
// GroupingExpr, Expr is the expression (Alias the optional `AS name`). For
// GroupingRollup / GroupingCube / GroupingSets, Items holds the parenthesized
// expression list (for GROUPING SETS the items may themselves be nested
// rollup/cube/empty sets, captured as their inner expressions for walking). For
// GroupingEmpty all are zero.
type GroupingItem struct {
	Kind  GroupingKind
	Expr  Node   // for GroupingExpr; nil otherwise
	Alias string // GroupingExpr `AS alias`; "" if absent
	Items []Node // for ROLLUP/CUBE/GROUPING SETS; nil otherwise
	Loc   Loc
}

// Tag implements Node.
func (n *GroupingItem) Tag() NodeTag { return T_GroupingItem }

// WindowDef is one named-window definition in a WINDOW clause (window_definition):
// `name AS <window_specification>`. Name is the window name; Spec is the
// specification (a *WindowSpec, reused from the expressions node — it may itself
// reference a base window by name).
type WindowDef struct {
	Name string
	Spec *WindowSpec
	Loc  Loc
}

// Tag implements Node.
func (n *WindowDef) Tag() NodeTag { return T_WindowDef }

// Compile-time assertions that the query node types satisfy Node.
var (
	_ Node = (*QueryStmt)(nil)
	_ Node = (*WithClause)(nil)
	_ Node = (*CTE)(nil)
	_ Node = (*RecursionDepth)(nil)
	_ Node = (*SelectStmt)(nil)
	_ Node = (*SelectItem)(nil)
	_ Node = (*SetOperation)(nil)
	_ Node = (*TableExpr)(nil)
	_ Node = (*UnnestExpr)(nil)
	_ Node = (*JoinExpr)(nil)
	_ Node = (*GroupByClause)(nil)
	_ Node = (*GroupingItem)(nil)
	_ Node = (*WindowDef)(nil)
)
