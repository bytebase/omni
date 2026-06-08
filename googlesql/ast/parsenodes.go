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
	// GranteeRole is a Spanner role grantee, e.g. `TO ROLE analyst`. The legacy
	// ZetaSQL grammar's grantee_list is string_literal_or_parameter only and has
	// no role form; Spanner GRANT/REVOKE instead names roles via a single `ROLE`
	// keyword followed by a comma-separated role-name list. Value holds the role
	// name. (Authoritative oracle: the live Spanner emulator — `TO 'string'`
	// REJECTS, `TO ROLE name [, name]` ACCEPTS.)
	GranteeRole
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
	case GranteeRole:
		return "ROLE"
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
	Path          NamePath     // the ON object name (the FIRST object; == Paths[0])
	// Paths is the full ON object list. The legacy ZetaSQL grant names a single
	// object (len 1); the Spanner role grant allows a comma-separated list
	// (`ON TABLE t1, t2 TO ROLE r`, emulator-verified). Path mirrors Paths[0] for
	// the common single-object case.
	Paths    []NamePath
	Grantees []*Grantee // the TO recipients (>= 1)
	// Roles is the Spanner role-to-role grant subject list: `GRANT ROLE r [, ...]
	// TO ROLE …`. When non-nil this is the role-grant form (a Spanner extension,
	// emulator-verified), and Privileges/AllPrivileges/ObjectType/Path are unset;
	// Grantees holds the target roles. Nil for an ordinary privilege grant.
	Roles []NamePath
	Loc   Loc
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
	Path          NamePath     // the ON object name (the FIRST object; == Paths[0])
	Paths         []NamePath   // full ON object list (Spanner allows a comma list); Path == Paths[0]
	Grantees      []*Grantee   // the FROM subjects (>= 1)
	// Roles is the Spanner role-to-role revoke subject list: `REVOKE ROLE r [,
	// ...] FROM ROLE …` (mirrors GrantStmt.Roles). Non-nil ⇒ role-revoke form.
	Roles []NamePath
	Loc   Loc
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
// STRICT.
//
// Column-match suffix (one of, mutually exclusive with each other):
//   - ByName: `BY NAME [ON (cols)]` — ByName true, MatchColumns holds the
//     optional ON column list.
//   - Corresponding: `[STRICT] CORRESPONDING [BY (cols)]` — Corresponding true,
//     MatchColumns holds the optional BY column list.
//
// All these column-match forms are BigQuery-valid but Spanner feature-rejects
// them; they parse in the union grammar (oracle: "BY NAME ... is not supported"
// / "CORRESPONDING ... is not supported", both classified accept).
type SetOperation struct {
	Op            SetOp
	All           bool   // ALL (vs DISTINCT)
	OuterMode     string // "" | "FULL" | "OUTER" | "LEFT"
	Strict        bool
	Corresponding bool     // [STRICT] CORRESPONDING [BY (cols)]
	ByName        bool     // BY NAME [ON (cols)]
	MatchColumns  []string // ON/BY column list for ByName/Corresponding; nil if absent
	Left          Node     // *SelectStmt | *QueryStmt
	Right         Node     // *SelectStmt | *QueryStmt
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
//
// Pivot / Unpivot / Sample carry the optional table-source operators
// (parser-query-clauses): a PIVOT(...), an UNPIVOT(...), or a TABLESAMPLE. The
// grammar attaches at most one pivot/unpivot per source and a sample is a
// distinct suffix, so at most one of {Pivot, Unpivot, Sample} is set (a PIVOT
// and a TABLESAMPLE never combine on one source — oracle-verified). The
// pivot/unpivot operator's own trailing alias is carried on the operator node;
// Alias is the alias preceding the operator (the table's own `[AS] name`).
type TableExpr struct {
	Path            *PathExpr      // table / array-field path; nil for subquery / TVF
	Subquery        Node           // ( query ) subquery; nil otherwise (*QueryStmt)
	Func            *FuncCall      // table-valued function call; nil otherwise
	Alias           string         // [AS] alias; "" if absent
	WithOffset      bool           // trailing WITH OFFSET
	WithOffsetAlias string         // alias for WITH OFFSET; "" if absent or unaliased
	SystemTime      Node           // FOR SYSTEM_TIME AS OF expr; nil if absent
	Pivot           *PivotClause   // PIVOT(...); nil if absent
	Unpivot         *UnpivotClause // UNPIVOT(...); nil if absent
	Sample          *SampleClause  // TABLESAMPLE …; nil if absent
	Loc             Loc
}

// Tag implements Node.
func (n *TableExpr) Tag() NodeTag { return T_TableExpr }

// UnnestExpr is an `UNNEST(array_expr [, …]) [[AS] alias] [WITH OFFSET [[AS] name]]`
// FROM source (unnest_expression as a table_path_expression_base). Array is the
// UNNEST(...) call (a *FuncCall named UNNEST, as the expressions node builds it).
// Alias / WithOffset / WithOffsetAlias mirror TableExpr.
//
// Pivot / Unpivot / Sample mirror TableExpr: the legacy grammar threads an
// unnest_expression through table_path_expression, which carries
// opt_pivot_or_unpivot_clause_and_alias, and through `table_primary
// sample_clause`. These operators are always SEMANTICALLY invalid on an array
// scan (the emulator rejects "PIVOT/TABLESAMPLE is not allowed with array
// scans"), but the union grammar PARSES them, so they are retained for parse
// parity. At most one of {Pivot, Unpivot, Sample} is set.
type UnnestExpr struct {
	Array           Node           // the UNNEST(...) call (*FuncCall); the array argument(s)
	Alias           string         // [AS] alias; "" if absent
	WithOffset      bool           // trailing WITH OFFSET
	WithOffsetAlias string         // alias for WITH OFFSET; "" if absent or unaliased
	Pivot           *PivotClause   // PIVOT(...); nil if absent
	Unpivot         *UnpivotClause // UNPIVOT(...); nil if absent
	Sample          *SampleClause  // TABLESAMPLE …; nil if absent
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
// PIVOT / UNPIVOT / TABLESAMPLE table-source operators (parser-query-clauses)
// ---------------------------------------------------------------------------

// PivotClause is a PIVOT operator applied to a FROM source (pivot_clause):
//
//	PIVOT ( <agg> [[AS] alias][, …] FOR <for_expr> IN ( <value> [[AS] alias][, …] ) )
//	  [ [AS] alias ]
//
// Aggregates is the comma-separated list of aggregate expressions (each an
// expression with an optional alias). For is the pivot input column/expression
// (parsed at higher-than-AND precedence so the trailing `IN (...)` is the pivot
// value list, NOT an `x IN (...)` predicate). Values is the IN (...) value list
// (each a value expression with an optional alias). Alias is the trailing
// `[AS] name` for the pivot result ("" if absent).
type PivotClause struct {
	Aggregates []*PivotExpr // aggregate expressions (>= 1)
	For        Node         // the FOR input column/expression
	Values     []*PivotExpr // IN (...) value list (>= 1)
	Alias      string       // trailing [AS] alias; "" if absent
	Loc        Loc
}

// Tag implements Node.
func (n *PivotClause) Tag() NodeTag { return T_PivotClause }

// PivotExpr is one `expression [[AS] alias]` entry of a PIVOT aggregate list or
// IN value list (pivot_expression / pivot_value). Expr is the expression; Alias
// is the optional `[AS] name` ("" if absent).
type PivotExpr struct {
	Expr  Node
	Alias string
	Loc   Loc
}

// Tag implements Node.
func (n *PivotExpr) Tag() NodeTag { return T_PivotExpr }

// UnpivotNullsMode encodes the optional {INCLUDE|EXCLUDE} NULLS modifier on an
// UNPIVOT (unpivot_nulls_filter).
type UnpivotNullsMode int

const (
	// UnpivotNullsUnspecified is no modifier (the grammar default).
	UnpivotNullsUnspecified UnpivotNullsMode = iota
	// UnpivotIncludeNulls is INCLUDE NULLS.
	UnpivotIncludeNulls
	// UnpivotExcludeNulls is EXCLUDE NULLS.
	UnpivotExcludeNulls
)

// String returns the modifier keyword phrase ("" for unspecified).
func (m UnpivotNullsMode) String() string {
	switch m {
	case UnpivotIncludeNulls:
		return "INCLUDE NULLS"
	case UnpivotExcludeNulls:
		return "EXCLUDE NULLS"
	default:
		return ""
	}
}

// UnpivotClause is an UNPIVOT operator applied to a FROM source (unpivot_clause):
//
//	UNPIVOT [ {INCLUDE|EXCLUDE} NULLS ]
//	  ( <value_col(s)> FOR <name_col> IN ( <in_item>[, …] ) ) [ [AS] alias ]
//
// NullsMode is the optional INCLUDE/EXCLUDE NULLS. ValueColumns is the value
// column list: a single column for single-column UNPIVOT, or several for the
// parenthesized multi-column form (`(c1, c2)`). NameColumn is the FOR name
// column. Items is the IN (...) list (each an UnpivotInItem). Alias is the
// trailing `[AS] name` ("" if absent).
type UnpivotClause struct {
	NullsMode    UnpivotNullsMode // INCLUDE / EXCLUDE NULLS; default unspecified
	ValueColumns []*PathExpr      // value column(s): 1 (single) or >1 (multi-column)
	NameColumn   *PathExpr        // FOR <name_col>
	Items        []*UnpivotInItem // IN (...) source items (>= 1)
	Alias        string           // trailing [AS] alias; "" if absent
	Loc          Loc
}

// Tag implements Node.
func (n *UnpivotClause) Tag() NodeTag { return T_UnpivotClause }

// UnpivotInItem is one entry of an UNPIVOT … IN ( … ) list (unpivot_in_item):
// a column or parenthesized column group, with an optional `[AS] string|int`
// row-value alias. Columns holds the column path(s) (one for the single-column
// form, several for a parenthesized group). AliasString / AliasInt carry the
// optional `AS 'name'` / `AS 1` literal alias (HasAlias flags its presence;
// AliasIsInt selects which literal field is meaningful).
type UnpivotInItem struct {
	Columns     []*PathExpr // column path(s): 1 (single) or >1 (parenthesized group)
	HasAlias    bool        // an `[AS] string|int` alias is present
	AliasIsInt  bool        // true if the alias is an integer literal (else a string)
	AliasString string      // the string-literal alias (when HasAlias && !AliasIsInt)
	AliasInt    string      // the integer-literal alias spelling (when HasAlias && AliasIsInt)
	Loc         Loc
}

// Tag implements Node.
func (n *UnpivotInItem) Tag() NodeTag { return T_UnpivotInItem }

// SampleSizeUnit selects the TABLESAMPLE size unit (sample_size_unit).
type SampleSizeUnit int

const (
	// SampleUnitPercent is PERCENT.
	SampleUnitPercent SampleSizeUnit = iota
	// SampleUnitRows is ROWS.
	SampleUnitRows
)

// String returns the unit keyword.
func (u SampleSizeUnit) String() string {
	if u == SampleUnitRows {
		return "ROWS"
	}
	return "PERCENT"
}

// SampleClause is a TABLESAMPLE operator applied to a FROM source (sample_clause):
//
//	TABLESAMPLE <method> ( <size> { PERCENT | ROWS } [ PARTITION BY <expr>[, …] ] )
//	  [ REPEATABLE ( <seed> ) | WITH WEIGHT [[AS] alias] [ REPEATABLE ( <seed> ) ] ]
//
// Method is the sampling method identifier (SYSTEM / BERNOULLI / RESERVOIR — the
// grammar reads it as a bare identifier, so any method name is accepted). Size
// is the sample-size expression; Unit is PERCENT or ROWS. PartitionBy carries an
// optional `PARTITION BY` list inside the size clause (nil if absent).
//
// The suffix is one of: nothing; a REPEATABLE(seed) (Repeatable set, WithWeight
// false); or WITH WEIGHT [[AS] alias] [REPEATABLE(seed)] (WithWeight true,
// WeightAlias the optional alias, Repeatable the optional trailing seed).
type SampleClause struct {
	Method      string // sampling method identifier (SYSTEM / BERNOULLI / RESERVOIR / …)
	Size        Node   // sample-size expression
	Unit        SampleSizeUnit
	PartitionBy []Node // PARTITION BY exprs inside the size clause; nil if absent
	WithWeight  bool   // a WITH WEIGHT suffix is present
	WeightAlias string // WITH WEIGHT [AS] alias; "" if absent / unaliased
	Repeatable  Node   // REPEATABLE ( seed ) seed expression; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *SampleClause) Tag() NodeTag { return T_SampleClause }

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

// ===========================================================================
// Core DDL (googlesql/parser-ddl node)
//
// AST nodes for the CREATE / ALTER / DROP statements over the table-like
// objects bytebase consumes for BigQuery + Spanner: TABLE, VIEW, INDEX,
// SCHEMA, DATABASE. These mirror the legacy ANTLR DDL grammar
// (GoogleSQLParser.g4 §2.2-§2.6 — create_table_statement, create_view_statement,
// create_index_statement, create_schema_statement, create_database_statement,
// alter_statement/alter_action, drop_statement), a hand-port of ZetaSQL, serving
// the UNION of both dialects through one grammar.
//
// Dialect-specific object DDL (BigQuery FUNCTION / PROCEDURE / MATERIALIZED VIEW
// / SEARCH+VECTOR INDEX / SNAPSHOT / ROW ACCESS POLICY / CAPACITY; Spanner CHANGE
// STREAM / SEQUENCE / named-schema role DDL / LOCALITY GROUP / PROTO BUNDLE) is
// owned by the parser-ddl-bigquery and parser-ddl-spanner nodes, which depend on
// this one.
//
// Design: object names are *PathExpr (a Node — the walker descends so the
// query-span extractor finds the referenced table), and every expression /
// embedded-query child (column DEFAULT / generated / CHECK exprs, OPTIONS values,
// CTAS and VIEW bodies, partition/cluster exprs, index key exprs) is an ast.Node
// field so the code-generated walker traverses it. Structural sub-nodes
// (ColumnDef, TableConstraint, OptionsEntry, …) implement Node so the walk
// reaches their expression children. Keyword-set choices (foreign-key actions,
// SQL SECURITY kind, generated mode) are kept as small enums/strings — the
// query-span path never inspects them and the package carries no deparser.

// ---------------------------------------------------------------------------
// Shared DDL leaves
// ---------------------------------------------------------------------------

// OptionsEntry is one `name <op> value` entry of an OPTIONS(...) list
// (options_entry). Name is the option key (identifier_in_hints, normalized);
// Op is the assignment operator spelling ("=", "+=", "-="); Value is the
// option value expression (an ast.Node — usually a Literal, but the grammar
// allows any expression, plus the bare PROTO keyword which is captured as a
// Literal with the value "PROTO"). The GoogleSQL grammar accepts arbitrary
// option keys, so the parser does not validate the vocabulary (the Spanner
// emulator does, which is why option-name rejects are a non-authoritative
// over-reject — see oracle.md).
type OptionsEntry struct {
	Name  string
	Op    string
	Value Node
	Loc   Loc
}

// Tag implements Node.
func (n *OptionsEntry) Tag() NodeTag { return T_OptionsEntry }

// KeyPart is one element of a PRIMARY KEY / index key list (primary_key_element
// / column_ordering_and_options_expr): a column name with an optional
// ASC/DESC direction and NULLS FIRST/LAST order. For an index key the element
// is a general expression in the grammar; the common case is a column path,
// captured in Expr (an ast.Node so the walker descends). Name is the plain
// identifier when the element is a bare column (the PRIMARY KEY case, which is
// identifier-only); it is "" for a general index expression.
type KeyPart struct {
	Name      string          // column name (PK element / bare-column index key); "" if Expr is a non-column expression
	Expr      Node            // index key expression (general column_ordering_and_options_expr); nil for a PK element
	Collate   string          // index key COLLATE spelling; "" if absent
	Direction string          // "ASC" | "DESC" | "" (none)
	NullOrder string          // "FIRST" | "LAST" | "" (none)
	Options   []*OptionsEntry // per-index-key OPTIONS(...); nil if absent
	Loc       Loc
}

// Tag implements Node.
func (n *KeyPart) Tag() NodeTag { return T_KeyPart }

// ---------------------------------------------------------------------------
// CREATE TABLE
// ---------------------------------------------------------------------------

// ForeignKeyAction enumerates a foreign_key_action (ON UPDATE / ON DELETE):
// NO ACTION, RESTRICT, CASCADE, SET NULL, or the unset (no action specified)
// state.
type ForeignKeyAction int

const (
	FKActionNone     ForeignKeyAction = iota // no action specified
	FKActionNoAction                         // NO ACTION
	FKActionRestrict                         // RESTRICT
	FKActionCascade                          // CASCADE
	FKActionSetNull                          // SET NULL
)

// String returns the SQL spelling of the action ("" for FKActionNone).
func (a ForeignKeyAction) String() string {
	switch a {
	case FKActionNoAction:
		return "NO ACTION"
	case FKActionRestrict:
		return "RESTRICT"
	case FKActionCascade:
		return "CASCADE"
	case FKActionSetNull:
		return "SET NULL"
	default:
		return ""
	}
}

// GeneratedMode enumerates the generated_mode prefix of a generated column:
// `AS`, `GENERATED AS`, `GENERATED ALWAYS AS`, or `GENERATED BY DEFAULT AS`.
type GeneratedMode int

const (
	GenModeAs                 GeneratedMode = iota // AS (expr)  (the Spanner spelling)
	GenModeGeneratedAs                             // GENERATED AS
	GenModeGeneratedAlways                         // GENERATED ALWAYS AS
	GenModeGeneratedByDefault                      // GENERATED BY DEFAULT AS
)

// ForeignKeyRef is a REFERENCES clause (foreign_key_reference): the referenced
// table path, its referenced columns, an optional MATCH mode, and the optional
// ON UPDATE / ON DELETE actions.
type ForeignKeyRef struct {
	Table    *PathExpr // referenced table (a Node; walked for query-span)
	Columns  []string  // referenced column names
	Match    string    // "SIMPLE" | "FULL" | "NOT DISTINCT" | "" (none)
	OnUpdate ForeignKeyAction
	OnDelete ForeignKeyAction
	Loc      Loc
}

// Tag implements Node.
func (n *ForeignKeyRef) Tag() NodeTag { return T_ForeignKeyRef }

// ColumnDef is one column definition (table_column_definition): a name, a type,
// optional collation, optional NOT NULL / HIDDEN / inline PRIMARY KEY / inline
// FOREIGN KEY attributes, an optional DEFAULT or generated-column clause, and an
// optional per-column OPTIONS list.
//
// Exactly one of Default / Generated is set (the grammar rejects both with a
// dedicated error). For a generated column Generated holds the expression and
// GenMode/Stored describe the mode; for a DEFAULT column Default holds the
// expression. Type is a TypeRef (rendered type text + span) built from the
// parsed GoogleSQL type, matching the CAST-target convention.
type ColumnDef struct {
	Name       string
	Type       *TypeRef       // column type (column_schema_inner); nil only for a generated-without-type column (rare)
	Collate    string         // trailing COLLATE spelling on the type; "" if absent
	NotNull    bool           // NOT NULL attribute
	Hidden     bool           // HIDDEN attribute (Spanner)
	PrimaryKey bool           // inline PRIMARY KEY column attribute
	Enforced   string         // "ENFORCED" | "NOT ENFORCED" | "" — constraint_enforcement after the attributes
	ForeignKey *ForeignKeyRef // inline foreign_key_column_attribute REFERENCES …; nil if absent
	FKName     string         // CONSTRAINT name preceding an inline FK; "" if absent
	Default    Node           // DEFAULT (expr); nil if absent
	Generated  Node           // generated-column AS (expr) body; nil if not generated
	GenMode    GeneratedMode  // the generated_mode prefix (valid only when Generated != nil)
	Stored     string         // "STORED" | "STORED VOLATILE" | "VIRTUAL" | "" — stored_mode (valid only when Generated != nil)
	Options    []*OptionsEntry
	Loc        Loc
}

// Tag implements Node.
func (n *ColumnDef) Tag() NodeTag { return T_ColumnDef }

// TableConstraintKind discriminates a table_constraint_definition.
type TableConstraintKind int

const (
	ConstraintPrimaryKey TableConstraintKind = iota // PRIMARY KEY ( … )
	ConstraintForeignKey                            // FOREIGN KEY ( … ) REFERENCES …
	ConstraintCheck                                 // CHECK ( expr )
)

// TableConstraint is a table-level constraint (table_constraint_definition): a
// PRIMARY KEY, FOREIGN KEY, or CHECK constraint, with an optional leading
// CONSTRAINT name and trailing enforcement / OPTIONS.
type TableConstraint struct {
	Kind       TableConstraintKind
	Name       string         // CONSTRAINT name; "" if anonymous
	KeyParts   []*KeyPart     // PRIMARY KEY element list (ConstraintPrimaryKey)
	Columns    []string       // FOREIGN KEY column list (ConstraintForeignKey)
	ForeignKey *ForeignKeyRef // REFERENCES clause (ConstraintForeignKey)
	Check      Node           // CHECK expression (ConstraintCheck)
	Enforced   string         // "ENFORCED" | "NOT ENFORCED" | ""
	Options    []*OptionsEntry
	Loc        Loc
}

// Tag implements Node.
func (n *TableConstraint) Tag() NodeTag { return T_TableConstraint }

// InterleaveClause is the Spanner INTERLEAVE IN PARENT clause on CREATE TABLE
// (opt_spanner_interleave_in_parent_clause): the parent table and the
// ON DELETE action.
type InterleaveClause struct {
	Parent   *PathExpr // parent table (a Node; walked)
	OnDelete ForeignKeyAction
	Loc      Loc
}

// Tag implements Node.
func (n *InterleaveClause) Tag() NodeTag { return T_InterleaveClause }

// CreateTableStmt is a CREATE TABLE statement (create_table_statement) for the
// BigQuery + Spanner union. It carries every documented sub-clause of both
// dialects; a given statement uses the subset its dialect permits (the parser
// accepts the union — the Spanner emulator's rejects of BigQuery-only clauses
// are non-authoritative, see oracle.md).
type CreateTableStmt struct {
	OrReplace   bool      // OR REPLACE (BigQuery)
	Scope       string    // "TEMP" | "TEMPORARY" | "PUBLIC" | "PRIVATE" | "" — opt_create_scope
	IfNotExists bool      // IF NOT EXISTS
	Name        *PathExpr // table name (maybe_dashed_path_expression)
	Columns     []*ColumnDef
	Constraints []*TableConstraint
	// Spanner trailing PRIMARY KEY ( … ) after the element list
	// (opt_spanner_table_options). PrimaryKey holds its key parts; HasPrimaryKey
	// distinguishes a present-but-empty `PRIMARY KEY ()` from an absent clause.
	PrimaryKey     []*KeyPart
	HasPrimaryKey  bool
	Interleave     *InterleaveClause // Spanner INTERLEAVE IN PARENT; nil if absent
	RowDeletion    Node              // Spanner ROW DELETION POLICY (expr); nil if absent (opt_ttl_clause)
	Like           *PathExpr         // LIKE source table (BigQuery); nil if absent
	Clone          *PathExpr         // CLONE source table (BigQuery); nil if absent
	Copy           *PathExpr         // COPY source table (BigQuery); nil if absent
	DefaultCollate string            // DEFAULT COLLATE spelling; "" if absent
	PartitionBy    []Node            // PARTITION BY expressions (BigQuery); nil if absent
	ClusterBy      []Node            // CLUSTER BY expressions (BigQuery); nil if absent
	Options        []*OptionsEntry
	AsQuery        Node // AS query (CTAS); nil if absent (*QueryStmt)
	Loc            Loc
}

// Tag implements Node.
func (n *CreateTableStmt) Tag() NodeTag { return T_CreateTableStmt }

// ---------------------------------------------------------------------------
// CREATE VIEW
// ---------------------------------------------------------------------------

// ViewColumn is one entry of a view's column_with_options_list (BigQuery view
// column list): a column name with optional per-column OPTIONS.
type ViewColumn struct {
	Name    string
	Options []*OptionsEntry
	Loc     Loc
}

// Tag implements Node.
func (n *ViewColumn) Tag() NodeTag { return T_ViewColumn }

// CreateViewStmt is a CREATE VIEW statement (the plain-view alternative of
// create_view_statement). MATERIALIZED / APPROX views are owned by
// parser-ddl-bigquery; this node covers the standard view both dialects share,
// plus Spanner's SQL SECURITY clause and BigQuery's column list + OPTIONS.
type CreateViewStmt struct {
	OrReplace   bool
	Scope       string // opt_create_scope (BigQuery); "" if absent
	Recursive   bool   // RECURSIVE
	IfNotExists bool
	Name        *PathExpr
	Columns     []*ViewColumn // column_with_options_list; nil if absent
	SQLSecurity string        // "INVOKER" | "DEFINER" | "" — opt_sql_security_clause
	Options     []*OptionsEntry
	AsQuery     Node // AS query (required); *QueryStmt
	Loc         Loc
}

// Tag implements Node.
func (n *CreateViewStmt) Tag() NodeTag { return T_CreateViewStmt }

// ---------------------------------------------------------------------------
// CREATE INDEX
// ---------------------------------------------------------------------------

// CreateIndexStmt is a CREATE INDEX statement (create_index_statement). The
// plain / UNIQUE / NULL_FILTERED forms are covered here; the SEARCH / VECTOR
// index_type forms are owned by parser-ddl-{bigquery,spanner} (the IndexType
// field carries the SEARCH/VECTOR token when this node is reused, but the
// core node leaves it "").
type CreateIndexStmt struct {
	OrReplace    bool
	Unique       bool
	NullFiltered bool   // NULL_FILTERED (Spanner)
	IndexType    string // "SEARCH" | "VECTOR" | "" (plain) — index_type
	IfNotExists  bool
	Name         *PathExpr  // index name
	Table        *PathExpr  // ON <table> (a Node; walked)
	Keys         []*KeyPart // index key list (index_order_by_and_options)
	AllColumns   bool       // ( ALL COLUMNS ) form (BigQuery SEARCH index)
	Storing      []Node     // STORING ( expr, … ); nil if absent
	PartitionBy  []Node     // PARTITION BY ( … ) suffix (BigQuery); nil if absent
	Interleave   *PathExpr  // , INTERLEAVE IN <table> (Spanner); nil if absent
	Where        Node       // WHERE filter (Spanner partial index); nil if absent
	Options      []*OptionsEntry
	Loc          Loc
}

// Tag implements Node.
func (n *CreateIndexStmt) Tag() NodeTag { return T_CreateIndexStmt }

// ---------------------------------------------------------------------------
// CREATE SCHEMA / DATABASE
// ---------------------------------------------------------------------------

// CreateSchemaStmt is a CREATE SCHEMA statement (create_schema_statement),
// shared by both dialects. BigQuery adds DEFAULT COLLATE and OPTIONS.
type CreateSchemaStmt struct {
	OrReplace      bool
	IfNotExists    bool
	Name           *PathExpr
	DefaultCollate string // "" if absent
	Options        []*OptionsEntry
	Loc            Loc
}

// Tag implements Node.
func (n *CreateSchemaStmt) Tag() NodeTag { return T_CreateSchemaStmt }

// CreateDatabaseStmt is a CREATE DATABASE statement (create_database_statement;
// Spanner). It carries the database name and optional OPTIONS.
type CreateDatabaseStmt struct {
	Name    *PathExpr
	Options []*OptionsEntry
	Loc     Loc
}

// Tag implements Node.
func (n *CreateDatabaseStmt) Tag() NodeTag { return T_CreateDatabaseStmt }

// ---------------------------------------------------------------------------
// ALTER
// ---------------------------------------------------------------------------

// AlterObjectKind discriminates the object an ALTER statement targets (the
// subset this node owns: TABLE, VIEW, INDEX, SCHEMA, DATABASE).
type AlterObjectKind int

const (
	AlterTable AlterObjectKind = iota
	AlterView
	AlterIndex
	AlterSchema
	AlterDatabase
)

// String returns the object keyword.
func (k AlterObjectKind) String() string {
	switch k {
	case AlterTable:
		return "TABLE"
	case AlterView:
		return "VIEW"
	case AlterIndex:
		return "INDEX"
	case AlterSchema:
		return "SCHEMA"
	case AlterDatabase:
		return "DATABASE"
	default:
		return "UNKNOWN"
	}
}

// AlterStmt is an ALTER statement over a table-like object (alter_statement
// alternatives for TABLE / schema-object). Object selects the kind; IfExists
// flags `IF EXISTS`; Name is the object path; Actions is the comma-separated
// alter_action list (>= 1).
type AlterStmt struct {
	Object   AlterObjectKind
	IfExists bool
	Name     *PathExpr
	Actions  []*AlterAction
	Loc      Loc
}

// Tag implements Node.
func (n *AlterStmt) Tag() NodeTag { return T_AlterStmt }

// AlterActionKind enumerates the alter_action variants this node supports.
// (Generic-entity / row-access-policy / privilege-restriction actions are out
// of scope — owned by the dialect-specific DDL nodes.)
type AlterActionKind int

const (
	AlterSetOptions               AlterActionKind = iota // SET OPTIONS ( … )
	AlterAddColumn                                       // ADD COLUMN [IF NOT EXISTS] column_def [position] [FILL USING expr]
	AlterDropColumn                                      // DROP COLUMN [IF EXISTS] name
	AlterRenameColumn                                    // RENAME COLUMN [IF EXISTS] name TO name
	AlterAlterColumnType                                 // ALTER COLUMN [IF EXISTS] name SET DATA TYPE field_schema  (BigQuery)
	AlterAlterColumnOptions                              // ALTER COLUMN [IF EXISTS] name SET OPTIONS ( … )
	AlterAlterColumnSetDefault                           // ALTER COLUMN [IF EXISTS] name SET DEFAULT expr
	AlterAlterColumnDropDefault                          // ALTER COLUMN [IF EXISTS] name DROP DEFAULT
	AlterAlterColumnDropNotNull                          // ALTER COLUMN [IF EXISTS] name DROP NOT NULL  (BigQuery)
	AlterAlterColumnDropGenerated                        // ALTER COLUMN [IF EXISTS] name DROP GENERATED
	AlterSpannerAlterColumn                              // Spanner ALTER COLUMN name <schema> [NOT NULL] [generated/default] [OPTIONS]
	AlterAddConstraint                                   // ADD [CONSTRAINT [IF NOT EXISTS] name] {PK | FK | CHECK}
	AlterDropConstraint                                  // DROP CONSTRAINT [IF EXISTS] name
	AlterDropPrimaryKey                                  // DROP PRIMARY KEY [IF EXISTS]
	AlterAlterConstraintEnforce                          // ALTER CONSTRAINT [IF EXISTS] name {ENFORCED|NOT ENFORCED}
	AlterAlterConstraintOptions                          // ALTER CONSTRAINT [IF EXISTS] name SET OPTIONS ( … )
	AlterRenameTo                                        // RENAME TO path
	AlterSetDefaultCollate                               // SET DEFAULT collate_clause
	AlterAddRowDeletionPolicy                            // ADD ROW DELETION POLICY [IF NOT EXISTS] ( expr )
	AlterReplaceRowDeletionPolicy                        // REPLACE ROW DELETION POLICY [IF EXISTS] ( expr )
	AlterDropRowDeletionPolicy                           // DROP ROW DELETION POLICY [IF EXISTS]
	AlterSetOnDelete                                     // Spanner SET ON DELETE foreign_key_action
	AlterAddStoredColumn                                 // ALTER INDEX … ADD STORED COLUMN name
	AlterDropStoredColumn                                // ALTER INDEX … DROP STORED COLUMN name
)

// AlterAction is one action of an alter_action_list. Kind selects the variant;
// the relevant fields below carry its payload. Expression children (Default,
// Check, RowDeletion) and the embedded column/constraint are ast.Node-typed so
// the walker descends.
type AlterAction struct {
	Kind           AlterActionKind
	IfExists       bool             // for the column/constraint/policy variants that allow IF EXISTS
	IfNotExists    bool             // for ADD COLUMN / ADD CONSTRAINT / ADD ROW DELETION POLICY
	Column         *ColumnDef       // ADD COLUMN / Spanner ALTER COLUMN target definition
	ColumnName     string           // DROP/RENAME/ALTER COLUMN target name
	NewName        string           // RENAME COLUMN … TO <NewName>
	Position       string           // ADD COLUMN position: "PRECEDING <id>" | "FOLLOWING <id>" | ""
	FillUsing      Node             // ADD COLUMN … FILL USING expr; nil if absent
	Constraint     *TableConstraint // ADD constraint payload
	ConstraintName string           // DROP/ALTER CONSTRAINT name; ADD CONSTRAINT name
	Enforced       string           // ALTER CONSTRAINT enforcement / column SET DATA TYPE trailing enforcement: "ENFORCED"|"NOT ENFORCED"|""
	NewType        *TypeRef         // SET DATA TYPE / Spanner ALTER COLUMN new type; nil if absent
	NotNull        bool             // Spanner ALTER COLUMN inline NOT NULL
	Generated      Node             // Spanner ALTER COLUMN generated AS (expr) STORED body; nil if absent
	Default        Node             // SET DEFAULT expr; nil if absent
	RowDeletion    Node             // ADD/REPLACE ROW DELETION POLICY ( expr ); nil if absent
	Collate        string           // SET DEFAULT COLLATE spelling; "" if absent
	OnDelete       ForeignKeyAction // SET ON DELETE action
	RenameTo       *PathExpr        // RENAME TO target (a Node; walked)
	StoredColumn   string           // ALTER INDEX ADD/DROP STORED COLUMN name
	Options        []*OptionsEntry  // SET OPTIONS / column SET OPTIONS payload
	Loc            Loc
}

// Tag implements Node.
func (n *AlterAction) Tag() NodeTag { return T_AlterAction }

// ---------------------------------------------------------------------------
// DROP
// ---------------------------------------------------------------------------

// DropObjectKind discriminates the object a DROP statement targets (the subset
// this node owns).
type DropObjectKind int

const (
	DropTable DropObjectKind = iota
	DropView
	DropIndex
	DropSchema
	DropDatabase
)

// String returns the object keyword(s).
func (k DropObjectKind) String() string {
	switch k {
	case DropTable:
		return "TABLE"
	case DropView:
		return "VIEW"
	case DropIndex:
		return "INDEX"
	case DropSchema:
		return "SCHEMA"
	case DropDatabase:
		return "DATABASE"
	default:
		return "UNKNOWN"
	}
}

// DropStmt is a DROP statement over a table-like object (drop_statement
// alternatives for INDEX / TABLE / schema-object). Object selects the kind;
// IfExists flags `IF EXISTS`; Name is the object path. DropMode carries the
// trailing RESTRICT / CASCADE (opt_drop_mode — attached to every
// schema_object_kind object: VIEW / SCHEMA / DATABASE / INDEX; NOT a bare TABLE).
// OnTable holds the `ON <table>` of a DROP {SEARCH|VECTOR} INDEX (BigQuery); core
// DDL never sets it (a plain DROP INDEX has no ON), so it is currently always nil
// — it is reserved for the SEARCH/VECTOR INDEX dialect node.
type DropStmt struct {
	Object   DropObjectKind
	External bool // DROP EXTERNAL SCHEMA (BigQuery)
	IfExists bool
	Name     *PathExpr
	OnTable  *PathExpr // DROP {SEARCH|VECTOR} INDEX … ON <table>; reserved (dialect node), nil here
	DropMode string    // "RESTRICT" | "CASCADE" | "" — opt_drop_mode
	Loc      Loc
}

// Tag implements Node.
func (n *DropStmt) Tag() NodeTag { return T_DropStmt }

// Compile-time assertions that the DDL node types satisfy Node.
var (
	_ Node = (*OptionsEntry)(nil)
	_ Node = (*KeyPart)(nil)
	_ Node = (*ForeignKeyRef)(nil)
	_ Node = (*ColumnDef)(nil)
	_ Node = (*TableConstraint)(nil)
	_ Node = (*InterleaveClause)(nil)
	_ Node = (*CreateTableStmt)(nil)
	_ Node = (*ViewColumn)(nil)
	_ Node = (*CreateViewStmt)(nil)
	_ Node = (*CreateIndexStmt)(nil)
	_ Node = (*CreateSchemaStmt)(nil)
	_ Node = (*CreateDatabaseStmt)(nil)
	_ Node = (*AlterStmt)(nil)
	_ Node = (*AlterAction)(nil)
	_ Node = (*DropStmt)(nil)
)

// ===========================================================================
// BigQuery-specific DDL (parser-ddl-bigquery node)
// ===========================================================================
//
// These node types model the BigQuery-only object DDL the core parser-ddl node
// routes to a stub: CREATE [AGGREGATE] FUNCTION / CREATE TABLE FUNCTION (TVF) /
// CREATE PROCEDURE, CREATE [MATERIALIZED|APPROX] VIEW, CREATE [SEARCH|VECTOR]
// INDEX, CREATE SNAPSHOT TABLE, CREATE ROW ACCESS POLICY, and the generic-entity
// CREATE/ALTER/DROP path (CAPACITY / RESERVATION / ASSIGNMENT and any other
// extensible object), plus their ALTER (ALTER MATERIALIZED VIEW SET OPTIONS,
// ALTER VECTOR INDEX … REBUILD) and DROP (DROP {SEARCH|VECTOR} INDEX … ON,
// DROP SNAPSHOT TABLE, DROP TABLE FUNCTION, DROP ALL ROW ACCESS POLICIES) forms.
//
// They are a hand-port of the legacy ANTLR GoogleSQLParser.g4 rules
// (create_function_statement, create_procedure_statement,
// create_table_function_statement, create_view_statement materialized/approx
// alternatives, create_index_statement SEARCH/VECTOR index_type,
// create_snapshot_statement, create_row_access_policy_statement,
// create_entity_statement, the generic-entity alter/drop alternatives), the same
// ZetaSQL lineage the rest of the grammar follows.
//
// ORACLE NOTE — every form here is BigQuery-only at the GoogleSQL UNION level:
// the Spanner emulator (the differential harness) either lacks the form entirely
// (TABLE FUNCTION, PROCEDURE, MATERIALIZED VIEW, SEARCH index ALL COLUMNS,
// SNAPSHOT, ROW ACCESS POLICY, CAPACITY/RESERVATION/ASSIGNMENT, ALTER VECTOR
// INDEX REBUILD, DROP {SEARCH|VECTOR} INDEX … ON, DROP TABLE FUNCTION) or
// supports only a NARROWER shape (a bare `CREATE FUNCTION … AS (expr)` and
// `DROP FUNCTION` parse, but TEMP / AGGREGATE / LANGUAGE / DETERMINISTIC / REMOTE
// reject). So the Spanner verdict is NON-authoritative for these and they are
// triangulated against the legacy .g4 + the BigQuery truth1 corpus (oracle.md);
// the unit tests in bq_*_test.go are the structural/accept-reject gate.

// FunctionParamMode is the IN / OUT / INOUT mode prefix of a procedure parameter
// (opt_procedure_parameter_mode). Function parameters carry no mode.
type FunctionParamMode int

const (
	ParamModeNone  FunctionParamMode = iota // no mode (function parameter)
	ParamModeIn                             // IN
	ParamModeOut                            // OUT
	ParamModeInout                          // INOUT
)

// String returns the mode keyword, or "" for ParamModeNone.
func (m FunctionParamMode) String() string {
	switch m {
	case ParamModeIn:
		return "IN"
	case ParamModeOut:
		return "OUT"
	case ParamModeInout:
		return "INOUT"
	default:
		return ""
	}
}

// FunctionParam is one parameter of a CREATE FUNCTION / TABLE FUNCTION /
// PROCEDURE parameter list (function_parameter / procedure_parameter). It models
// the union of both grammars:
//
//   - function_parameter: identifier type_or_tvf_schema [AS alias]
//     [DEFAULT expr] [NOT AGGREGATE]  — or a bare type_or_tvf_schema (no name).
//   - procedure_parameter: [IN|OUT|INOUT] identifier type_or_tvf_schema.
//
// Type holds the rendered parameter type (a simple type, ANY TYPE templated
// type, or a TVF `TABLE<…>` schema — all captured via Type.Text). For a bare
// unnamed function parameter, Name is "" and Type still carries the type.
type FunctionParam struct {
	Mode         FunctionParamMode // procedure mode; ParamModeNone for functions / unset
	Name         string            // parameter name; "" for a bare unnamed function parameter
	Type         *TypeRef          // parameter type (type_or_tvf_schema); rendered text
	Alias        string            // AS alias (function_parameter opt_as_alias_with_required_as); "" if absent
	Default      Node              // DEFAULT expr (function_parameter); nil if absent
	NotAggregate bool              // NOT AGGREGATE (function_parameter)
	Loc          Loc
}

// Tag implements Node.
func (n *FunctionParam) Tag() NodeTag { return T_FunctionParam }

// CreateFunctionStmt is a CREATE [AGGREGATE] FUNCTION or CREATE TABLE FUNCTION
// statement (create_function_statement / create_table_function_statement). The
// IsTableFunc flag selects the TVF shape (CREATE TABLE FUNCTION … RETURNS TABLE<…>
// AS query) from the scalar/aggregate shape (CREATE [AGGREGATE] FUNCTION … RETURNS
// type AS (expr)|string).
type CreateFunctionStmt struct {
	OrReplace   bool
	Scope       string // opt_create_scope: "TEMP" | "TEMPORARY" | "" (PUBLIC/PRIVATE rare)
	Aggregate   bool   // AGGREGATE (scalar aggregate function)
	IsTableFunc bool   // CREATE TABLE FUNCTION (TVF)
	IfNotExists bool
	Name        *PathExpr
	Params      []*FunctionParam // function_parameters / opt_function_parameters
	Returns     *TypeRef         // RETURNS type (scalar) — rendered text; nil if absent
	// RETURNS TABLE<col …> (TVF): the column declarations. ReturnsTable is true when
	// a TVF declares an explicit RETURNS TABLE<…>; ReturnColumns holds the columns.
	ReturnsTable   bool
	ReturnColumns  []*ColumnDef
	SQLSecurity    string // SQL SECURITY {INVOKER|DEFINER}; "" if absent
	Determinism    string // DETERMINISTIC | NOT DETERMINISTIC | IMMUTABLE | STABLE | VOLATILE; "" if absent
	Language       string // LANGUAGE <id> (js/python/…); "" if absent
	Remote         bool   // REMOTE [WITH CONNECTION …]
	HasConnection  bool   // a WITH CONNECTION clause was present (the connection path is captured as text only)
	ConnectionName string // the connection path text, if any
	Options        []*OptionsEntry
	Body           Node   // scalar/aggregate AS (expr) body; nil if absent or string/TVF
	BodyString     string // AS string body (JS/Python code or quoted SQL); "" if absent
	HasBodyString  bool   // distinguishes an empty AS '' from an absent string body
	AsQuery        Node   // TVF AS query body; nil if absent (*QueryStmt)
	Loc            Loc
}

// Tag implements Node.
func (n *CreateFunctionStmt) Tag() NodeTag { return T_CreateFunctionStmt }

// CreateProcedureStmt is a CREATE PROCEDURE statement (create_procedure_statement):
// a parameter list, optional EXTERNAL SECURITY / WITH CONNECTION / OPTIONS, then
// either a BEGIN…END block body or a LANGUAGE <id> [AS code] body.
type CreateProcedureStmt struct {
	OrReplace        bool
	Scope            string // opt_create_scope
	IfNotExists      bool
	Name             *PathExpr
	Params           []*FunctionParam // procedure_parameters
	ExternalSecurity string           // EXTERNAL SECURITY {INVOKER|DEFINER}; "" if absent
	HasConnection    bool             // WITH CONNECTION present
	ConnectionName   string           // connection path text, if any
	Options          []*OptionsEntry
	// Body: a BEGIN…END block is captured as raw text (the procedural script
	// statement_list is owned by the parser-scripting node; this node records the
	// block text so the procedure round-trips and the consumer sees a body). For a
	// LANGUAGE body, Language is set and BodyString holds the AS code.
	BodyText      string // BEGIN…END block text (verbatim, incl. the BEGIN/END keywords); "" for a LANGUAGE body
	Language      string // LANGUAGE <id>; "" for a BEGIN…END body
	BodyString    string // LANGUAGE … AS <code>; "" if absent
	HasBodyString bool
	Loc           Loc
}

// Tag implements Node.
func (n *CreateProcedureStmt) Tag() NodeTag { return T_CreateProcedureStmt }

// ViewKind discriminates a CREATE VIEW variant: a plain view (owned by the core
// parser-ddl node, not produced here) vs the BigQuery MATERIALIZED / APPROX
// variants this node owns.
type ViewKind int

const (
	ViewPlain        ViewKind = iota // plain VIEW (core node)
	ViewMaterialized                 // MATERIALIZED VIEW
	ViewApprox                       // APPROX VIEW
)

// String returns the view-kind keyword phrase.
func (k ViewKind) String() string {
	switch k {
	case ViewMaterialized:
		return "MATERIALIZED VIEW"
	case ViewApprox:
		return "APPROX VIEW"
	default:
		return "VIEW"
	}
}

// CreateMaterializedViewStmt is a CREATE MATERIALIZED|APPROX VIEW statement (the
// materialized / approx alternatives of create_view_statement). It adds, over a
// plain view, PARTITION BY / CLUSTER BY (materialized only) and the
// `AS REPLICA OF <source>` body alternative (query_or_replica_source).
type CreateMaterializedViewStmt struct {
	Kind        ViewKind // ViewMaterialized | ViewApprox
	OrReplace   bool
	Recursive   bool // RECURSIVE
	IfNotExists bool
	Name        *PathExpr
	Columns     []*ViewColumn // column_with_options_list; nil if absent
	SQLSecurity string        // SQL SECURITY {INVOKER|DEFINER}; "" if absent
	PartitionBy []Node        // PARTITION BY exprs (materialized only); nil if absent
	ClusterBy   []Node        // CLUSTER BY exprs (materialized only); nil if absent
	Options     []*OptionsEntry
	AsQuery     Node      // AS query body; nil when ReplicaOf is set (*QueryStmt)
	ReplicaOf   *PathExpr // AS REPLICA OF <source> (materialized); nil if a query body
	Loc         Loc
}

// Tag implements Node.
func (n *CreateMaterializedViewStmt) Tag() NodeTag { return T_CreateMaterializedViewStmt }

// CreateSnapshotStmt is a CREATE SNAPSHOT TABLE statement (create_snapshot_statement):
//
//	CREATE [OR REPLACE] SNAPSHOT TABLE [IF NOT EXISTS] <name>
//	  CLONE <source> [FOR SYSTEM_TIME AS OF expr] [OPTIONS(…)]
type CreateSnapshotStmt struct {
	OrReplace     bool
	IfNotExists   bool
	Name          *PathExpr
	CloneSource   *PathExpr // CLONE <source>
	ForSystemTime Node      // FOR SYSTEM_TIME AS OF expr; nil if absent
	Where         Node      // optional WHERE filter on the clone source; nil if absent
	Options       []*OptionsEntry
	Loc           Loc
}

// Tag implements Node.
func (n *CreateSnapshotStmt) Tag() NodeTag { return T_CreateSnapshotStmt }

// SearchVectorIndexStmt is a CREATE [SEARCH|VECTOR] INDEX statement (the
// SEARCH/VECTOR index_type variants of create_index_statement). It carries the
// index name, target table, the column list or ALL COLUMNS, STORING list,
// PARTITION BY (vector), and OPTIONS.
type SearchVectorIndexStmt struct {
	IsVector    bool // VECTOR (vs SEARCH)
	OrReplace   bool // OR REPLACE (vector)
	IfNotExists bool
	Name        *PathExpr  // index name
	Table       *PathExpr  // ON <table>
	Keys        []*KeyPart // column list (index_order_by_and_options); nil for ALL COLUMNS
	AllColumns  bool       // ( ALL COLUMNS [WITH COLUMN OPTIONS(…)] ) — SEARCH index
	Storing     []Node     // STORING ( col, … ) — vector index; nil if absent
	PartitionBy []Node     // PARTITION BY expr — vector index; nil if absent
	Options     []*OptionsEntry
	Loc         Loc
}

// Tag implements Node.
func (n *SearchVectorIndexStmt) Tag() NodeTag { return T_SearchVectorIndexStmt }

// CreateRowAccessPolicyStmt is a CREATE ROW ACCESS POLICY statement
// (create_row_access_policy_statement):
//
//	CREATE [OR REPLACE] ROW ACCESS POLICY [IF NOT EXISTS] [name]
//	  ON <table> [GRANT TO (grantees) | TO grantees] [FILTER] USING (expr)
type CreateRowAccessPolicyStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        string     // policy name (identifier?); "" if absent
	Table       *PathExpr  // ON <table>
	Grantees    []*Grantee // GRANT TO (…) / TO …; nil if absent
	HasGrantTo  bool       // a grant-to clause was present (distinguishes empty list)
	Filter      Node       // [FILTER] USING (expr) — required; the filter expression
	Loc         Loc
}

// Tag implements Node.
func (n *CreateRowAccessPolicyStmt) Tag() NodeTag { return T_CreateRowAccessPolicyStmt }

// CreateEntityStmt is a CREATE <generic-entity> statement (create_entity_statement):
//
//	CREATE [OR REPLACE] <entity_type> [IF NOT EXISTS] <name> [OPTIONS(…)] [AS <body>]
//
// The generic-entity mechanism is how the legacy grammar parses extensible
// object kinds whose keyword is a bare identifier — including BigQuery
// CAPACITY / RESERVATION / ASSIGNMENT (DDL-024/025/026: `CREATE CAPACITY
// \`p.l.c\` OPTIONS(…)`). EntityType holds the entity keyword spelling.
type CreateEntityStmt struct {
	OrReplace   bool
	EntityType  string // the generic_entity_type identifier (e.g. "CAPACITY", "RESERVATION", "ASSIGNMENT", "PROJECT")
	IfNotExists bool
	Name        *PathExpr
	Options     []*OptionsEntry
	// BodyText: the AS <generic_entity_body> tail, captured as raw text (the body
	// is a free-form JSON/string in the grammar). "" if absent.
	BodyText string
	Loc      Loc
}

// Tag implements Node.
func (n *CreateEntityStmt) Tag() NodeTag { return T_CreateEntityStmt }

// BQAlterObjectKind discriminates the object kind of a BigQuery-only ALTER
// statement this node owns (the kinds the core parser-ddl node routes away).
type BQAlterObjectKind int

const (
	BQAlterMaterializedView BQAlterObjectKind = iota // ALTER MATERIALIZED VIEW
	BQAlterApproxView                                // ALTER APPROX VIEW
	BQAlterVectorIndex                               // ALTER VECTOR INDEX … REBUILD
	BQAlterEntity                                    // ALTER <generic-entity> (CAPACITY/RESERVATION/…)
)

// String returns the object keyword phrase.
func (k BQAlterObjectKind) String() string {
	switch k {
	case BQAlterMaterializedView:
		return "MATERIALIZED VIEW"
	case BQAlterApproxView:
		return "APPROX VIEW"
	case BQAlterVectorIndex:
		return "VECTOR INDEX"
	default:
		return "ENTITY"
	}
}

// BQAlterStmt is a BigQuery-only ALTER statement: ALTER MATERIALIZED|APPROX VIEW
// SET OPTIONS, ALTER VECTOR INDEX … ON <table> REBUILD [OPTIONS(…)], or ALTER
// <generic-entity> SET OPTIONS / SET AS <body>.
type BQAlterStmt struct {
	Object     BQAlterObjectKind
	EntityType string // generic_entity_type spelling when Object==BQAlterEntity
	IfExists   bool
	Name       *PathExpr
	OnTable    *PathExpr       // ON <table> for ALTER VECTOR INDEX; nil otherwise
	Rebuild    bool            // REBUILD (vector index)
	SetOptions []*OptionsEntry // SET OPTIONS(…); nil if absent
	SetAsBody  string          // SET AS <generic_entity_body> text (entity); "" if absent
	Loc        Loc
}

// Tag implements Node.
func (n *BQAlterStmt) Tag() NodeTag { return T_BQAlterStmt }

// BQDropObjectKind discriminates the object kind of a BigQuery-only DROP
// statement this node owns.
type BQDropObjectKind int

const (
	BQDropFunction         BQDropObjectKind = iota // DROP FUNCTION
	BQDropTableFunction                            // DROP TABLE FUNCTION
	BQDropProcedure                                // DROP PROCEDURE
	BQDropMaterializedView                         // DROP MATERIALIZED VIEW
	BQDropApproxView                               // DROP APPROX VIEW
	BQDropSnapshotTable                            // DROP SNAPSHOT TABLE
	BQDropSearchIndex                              // DROP SEARCH INDEX … ON <table>
	BQDropVectorIndex                              // DROP VECTOR INDEX … ON <table>
	BQDropRowAccessPolicy                          // DROP ROW ACCESS POLICY … ON <table>
	BQDropEntity                                   // DROP <generic-entity> (CAPACITY/RESERVATION/…)
)

// String returns the object keyword phrase.
func (k BQDropObjectKind) String() string {
	switch k {
	case BQDropFunction:
		return "FUNCTION"
	case BQDropTableFunction:
		return "TABLE FUNCTION"
	case BQDropProcedure:
		return "PROCEDURE"
	case BQDropMaterializedView:
		return "MATERIALIZED VIEW"
	case BQDropApproxView:
		return "APPROX VIEW"
	case BQDropSnapshotTable:
		return "SNAPSHOT TABLE"
	case BQDropSearchIndex:
		return "SEARCH INDEX"
	case BQDropVectorIndex:
		return "VECTOR INDEX"
	case BQDropRowAccessPolicy:
		return "ROW ACCESS POLICY"
	default:
		return "ENTITY"
	}
}

// BQDropStmt is a BigQuery-only DROP statement: DROP FUNCTION / TABLE FUNCTION /
// PROCEDURE / MATERIALIZED VIEW / APPROX VIEW / SNAPSHOT TABLE / {SEARCH|VECTOR}
// INDEX … ON <table> / ROW ACCESS POLICY … ON <table> / generic-entity.
type BQDropStmt struct {
	Object     BQDropObjectKind
	EntityType string // generic_entity_type spelling when Object==BQDropEntity
	IfExists   bool
	Name       *PathExpr // object path; for ROW ACCESS POLICY this is the policy name path
	OnTable    *PathExpr // ON <table> for {SEARCH|VECTOR} INDEX / ROW ACCESS POLICY; nil otherwise
	Loc        Loc
}

// Tag implements Node.
func (n *BQDropStmt) Tag() NodeTag { return T_BQDropStmt }

// DropAllRowAccessPoliciesStmt is a DROP ALL ROW ACCESS POLICIES statement
// (drop_all_row_access_policies_statement): `DROP ALL ROW ACCESS POLICIES ON
// <table>`.
type DropAllRowAccessPoliciesStmt struct {
	Table *PathExpr // ON <table>
	Loc   Loc
}

// Tag implements Node.
func (n *DropAllRowAccessPoliciesStmt) Tag() NodeTag { return T_DropAllRowAccessPoliciesStmt }

// Compile-time assertions that the BigQuery DDL node types satisfy Node.
var (
	_ Node = (*FunctionParam)(nil)
	_ Node = (*CreateFunctionStmt)(nil)
	_ Node = (*CreateProcedureStmt)(nil)
	_ Node = (*CreateMaterializedViewStmt)(nil)
	_ Node = (*CreateSnapshotStmt)(nil)
	_ Node = (*SearchVectorIndexStmt)(nil)
	_ Node = (*CreateRowAccessPolicyStmt)(nil)
	_ Node = (*CreateEntityStmt)(nil)
	_ Node = (*BQAlterStmt)(nil)
	_ Node = (*BQDropStmt)(nil)
	_ Node = (*DropAllRowAccessPoliciesStmt)(nil)
)

// ---------------------------------------------------------------------------
// DML — INSERT / UPDATE / DELETE / MERGE / TRUNCATE (parser-dml node)
// ---------------------------------------------------------------------------
//
// These node types model the legacy ANTLR DML family (GoogleSQLParser.g4 §2.7),
// a hand-port of Google's open-source ZetaSQL reference and the grammar bytebase
// consumes today. One omni parser serves the BigQuery + Spanner union; the
// dialect deltas are FEATURE rejections layered on top of a shared grammar (the
// grammar parses MERGE, ON CONFLICT, THEN RETURN, INSERT OR UPDATE/IGNORE for
// both, and the engine rejects the unsupported ones at execution). The
// authoritative accept/reject oracle is the canonical ZetaSQL corpus
// (parser/googlesql/examples/.../dml_*.sql) plus the live Spanner emulator for
// the shared + Spanner-only forms (BigQuery-only forms — MERGE, TRUNCATE,
// dashed paths — are triangulated against the docs + the .g4; the emulator's
// verdict on them is NON-authoritative; see the divergence ledger).
//
// The grammar is (abbreviated):
//
//	insert_statement: insert_statement_prefix column_list? insert_values_or_query           opt_assert? opt_returning?
//	                | insert_statement_prefix column_list? insert_values_list_or_table_clause on_conflict opt_assert? opt_returning?
//	                | insert_statement_prefix column_list? '(' query ')'                       on_conflict opt_assert? opt_returning?
//	insert_statement_prefix: INSERT opt_or_ignore_replace_update? INTO? <target> hint?
//	opt_or_ignore_replace_update: [OR] IGNORE | [OR] REPLACE | [OR] UPDATE
//	update_statement: UPDATE <target> hint? as_alias? with_offset? SET update_item_list from_clause? where? assert? returning?
//	delete_statement: DELETE FROM? <target> hint? as_alias? with_offset? where? assert? returning?
//	merge_statement:  MERGE INTO? <target> as_alias? USING merge_source ON expr (merge_when_clause)+
//	truncate_statement: TRUNCATE TABLE <path> where?
//
// <target> is maybe_dashed_generalized_path_expression — a dotted/dashed path
// optionally extended by generalized accessors (`.field`, `.(extension)`,
// `[index]`) for nested/array DML targets (oracle: `UPDATE t.a[0].b SET …`
// parses; the nested-target rejection is semantic, not syntactic).

// DefaultExpr is the `DEFAULT` keyword used where an expression is expected
// (expression_or_default's DEFAULT alternative): a VALUES row element, an
// update_set_value RHS, or a merge INSERT value. It carries no payload beyond
// its source position — the engine fills the column default at execution.
type DefaultExpr struct {
	Loc Loc
}

// Tag implements Node.
func (n *DefaultExpr) Tag() NodeTag { return T_DefaultExpr }

// InsertStmt is an INSERT statement (insert_statement). Exactly one of
// {Rows, Query, TableClause} carries the inserted data:
//   - Rows:        a VALUES list — one InsertRow per parenthesized row.
//   - Query:       an `INSERT … <query>` or `INSERT … ( query )` source (a
//     *QueryStmt). The parenthesized form is recorded with Query set and
//     QueryParens true (it is the on-conflict-eligible alternative).
//   - TableClause: an `INSERT … TABLE <path-or-tvf> [WHERE …]` source
//     (insert_values_list_or_table_clause's table_clause alt; on-conflict only).
//
// OrAction is the upsert modifier (opt_or_ignore_replace_update) — one of "",
// "OR IGNORE", "IGNORE", "OR REPLACE", "REPLACE", "OR UPDATE", "UPDATE" (Spanner
// upserts; the bare IGNORE/REPLACE/UPDATE forms drop the OR and the emulator
// accepts them). Into records whether the optional INTO keyword was present.
// Target is the table path (maybe_dashed_generalized_path_expression). Columns
// is the optional explicit column_list (nil if absent). OnConflict is the
// ON CONFLICT clause (only valid with the VALUES / TABLE / `( query )` forms;
// nil otherwise). AssertRows is the ASSERT_ROWS_MODIFIED count expression (nil
// if absent). Returning is the THEN RETURN clause (nil if absent).
type InsertStmt struct {
	OrAction    string       // "" | "OR IGNORE" | "IGNORE" | "OR REPLACE" | "REPLACE" | "OR UPDATE" | "UPDATE"
	Into        bool         // INTO keyword present
	Target      *PathExpr    // table path
	Columns     []string     // explicit column list; nil if absent
	Rows        []*InsertRow // VALUES rows; nil unless the VALUES form
	Query       Node         // INSERT … SELECT / ( query ) source (*QueryStmt); nil unless the query form
	QueryParens bool         // the query source was written as `( query )` (on-conflict-eligible)
	TableClause *InsertTable // INSERT … TABLE source; nil unless the table-clause form
	OnConflict  *OnConflict  // ON CONFLICT clause; nil if absent
	AssertRows  Node         // ASSERT_ROWS_MODIFIED count; nil if absent
	Returning   *Returning   // THEN RETURN clause; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *InsertStmt) Tag() NodeTag { return T_InsertStmt }

// InsertRow is one parenthesized row of an INSERT … VALUES list
// (insert_values_row): a list of expression_or_default values (a DEFAULT value
// is a *DefaultExpr element).
type InsertRow struct {
	Values []Node // expression_or_default elements (>= 1)
	Loc    Loc
}

// Tag implements Node.
func (n *InsertRow) Tag() NodeTag { return T_InsertRow }

// InsertTable is the `TABLE <path-or-tvf> [WHERE expr]` insert source
// (table_clause_unreversed → table_clause_no_keyword). Path is the source table
// path or TVF name; Func, when set, makes it a TVF call. Where is the optional
// trailing WHERE filter. This form is on-conflict-only in the grammar.
type InsertTable struct {
	Path  *PathExpr // source table path / TVF name
	Func  *FuncCall // table-valued function call; nil for a plain table
	Where Node      // trailing WHERE expr; nil if absent
	Loc   Loc
}

// Tag implements Node.
func (n *InsertTable) Tag() NodeTag { return T_InsertTable }

// OnConflict is an ON CONFLICT clause on an INSERT (on_conflict_clause):
//
//	ON CONFLICT [conflict_target] DO NOTHING
//	ON CONFLICT [conflict_target] DO UPDATE SET update_item_list [WHERE expr]
//
// Target is the optional conflict target: either an explicit Columns list
// (`( col, … )`) or a named unique constraint (ConstraintName, from
// `ON UNIQUE CONSTRAINT id`). DoNothing selects the DO NOTHING action; when
// false the action is DO UPDATE with SetItems (and optional Where). At most one
// of Columns / ConstraintName is set; both empty means no conflict target.
type OnConflict struct {
	Columns        []string      // ON CONFLICT ( col, … ); nil if none / constraint form
	ConstraintName string        // ON CONFLICT ON UNIQUE CONSTRAINT <id>; "" if none / columns form
	DoNothing      bool          // DO NOTHING (true) vs DO UPDATE (false)
	SetItems       []*UpdateItem // DO UPDATE SET items; nil for DO NOTHING
	Where          Node          // DO UPDATE … WHERE expr; nil if absent
	Loc            Loc
}

// Tag implements Node.
func (n *OnConflict) Tag() NodeTag { return T_OnConflict }

// UpdateStmt is an UPDATE statement (update_statement):
//
//	UPDATE <target> [hint] [[AS] alias] [WITH OFFSET [[AS] name]]
//	  SET <update_item_list> [FROM <from>] [WHERE expr]
//	  [ASSERT_ROWS_MODIFIED count] [THEN RETURN …]
//
// Target is the table path. Alias is the optional `[AS] alias` ("" if absent).
// WithOffset / WithOffsetAlias record the Spanner `WITH OFFSET [[AS] name]`
// companion column. Items is the non-empty SET list. From is the optional
// join-update FROM source list (nil if absent). Where / AssertRows / Returning
// mirror the other DML nodes. A leading statement `@{…}` hint and a per-target
// `@{…}` hint are skipped (not retained) like elsewhere in the parser.
type UpdateStmt struct {
	Target          *PathExpr     // table path
	Alias           string        // [AS] alias; "" if absent
	WithOffset      bool          // WITH OFFSET
	WithOffsetAlias string        // alias for WITH OFFSET; "" if absent / unaliased
	Items           []*UpdateItem // SET items (>= 1)
	From            []Node        // FROM sources (*TableExpr / *JoinExpr / *UnnestExpr); nil if absent
	Where           Node          // WHERE expr; nil if absent
	AssertRows      Node          // ASSERT_ROWS_MODIFIED count; nil if absent
	Returning       *Returning    // THEN RETURN clause; nil if absent
	Loc             Loc
}

// Tag implements Node.
func (n *UpdateStmt) Tag() NodeTag { return T_UpdateStmt }

// UpdateItem is one entry of a SET clause (update_item). Two shapes:
//   - assignment:  `generalized_path = expression_or_default`. Path holds the
//     LHS generalized path spelling (`a`, `a.b.c`, `id.(pkg.ext)`, `id[0]`);
//     Value is the RHS (an expression or a *DefaultExpr). Nested is nil.
//   - nested DML:  `( insert | update | delete )` — a nested statement run
//     against an array-valued column (Spanner). Nested holds the inner DML
//     node (*InsertStmt / *UpdateStmt / *DeleteStmt); Path is "" and Value nil.
type UpdateItem struct {
	Path   string // LHS generalized path; "" for the nested-DML form
	Value  Node   // RHS expression or *DefaultExpr; nil for the nested-DML form
	Nested Node   // nested ( dml ) statement; nil for the assignment form
	Loc    Loc
}

// Tag implements Node.
func (n *UpdateItem) Tag() NodeTag { return T_UpdateItem }

// DeleteStmt is a DELETE statement (delete_statement):
//
//	DELETE [FROM] <target> [hint] [[AS] alias] [WITH OFFSET [[AS] name]]
//	  [WHERE expr] [ASSERT_ROWS_MODIFIED count] [THEN RETURN …]
//
// The FROM keyword is optional (ZetaSQL: `DELETE T WHERE …` parses). Target /
// Alias / WithOffset / WithOffsetAlias / Where / AssertRows / Returning mirror
// UpdateStmt. From records whether the optional FROM keyword was present (it has
// no structural effect, but is retained for faithful round-trip).
type DeleteStmt struct {
	From            bool       // FROM keyword present
	Target          *PathExpr  // table path
	Alias           string     // [AS] alias; "" if absent
	WithOffset      bool       // WITH OFFSET
	WithOffsetAlias string     // alias for WITH OFFSET; "" if absent / unaliased
	Where           Node       // WHERE expr; nil if absent
	AssertRows      Node       // ASSERT_ROWS_MODIFIED count; nil if absent
	Returning       *Returning // THEN RETURN clause; nil if absent
	Loc             Loc
}

// Tag implements Node.
func (n *DeleteStmt) Tag() NodeTag { return T_DeleteStmt }

// MergeStmt is a MERGE statement (merge_statement):
//
//	MERGE [INTO] <target> [[AS] alias] USING <source> ON <expr> <when_clause>+
//
// Target is the merge target path; Alias the optional `[AS] alias`. Source is
// the USING source (a *TableExpr — a table path or a `( query )` subquery, per
// merge_source). On is the merge condition expression. Whens is the non-empty
// list of WHEN clauses. MERGE is BigQuery-only among the bytebase consumers
// (Spanner feature-rejects it after parse); the grammar — and therefore omni —
// parses it for both. The emulator gives NO authoritative syntax verdict on
// MERGE bodies (it rejects at statement dispatch), so MERGE forms are oracled
// against the ZetaSQL corpus + the docs, not the emulator.
type MergeStmt struct {
	Target *PathExpr    // merge target path
	Alias  string       // [AS] alias; "" if absent
	Source Node         // USING source (*TableExpr: path or ( query ))
	On     Node         // ON merge condition
	Whens  []*MergeWhen // WHEN clauses (>= 1)
	Loc    Loc
}

// Tag implements Node.
func (n *MergeStmt) Tag() NodeTag { return T_MergeStmt }

// MergeWhen is one WHEN clause of a MERGE (merge_when_clause):
//
//	WHEN MATCHED [AND expr] THEN <action>
//	WHEN NOT MATCHED [BY TARGET] [AND expr] THEN <action>
//	WHEN NOT MATCHED BY SOURCE [AND expr] THEN <action>
//
// Matched is true for `WHEN MATCHED`; NotMatched true for `WHEN NOT MATCHED`
// (exactly one holds). For the NOT MATCHED forms, BySource distinguishes
// `BY SOURCE` from the default/`BY TARGET`; ByTarget records whether the
// explicit `BY TARGET` words were present (default when neither). And is the
// optional `AND search_condition` (nil if absent). Action is the THEN action.
type MergeWhen struct {
	Matched    bool         // WHEN MATCHED
	NotMatched bool         // WHEN NOT MATCHED
	ByTarget   bool         // explicit BY TARGET (NOT MATCHED only)
	BySource   bool         // BY SOURCE (NOT MATCHED only)
	And        Node         // AND search_condition; nil if absent
	Action     *MergeAction // THEN action
	Loc        Loc
}

// Tag implements Node.
func (n *MergeWhen) Tag() NodeTag { return T_MergeWhen }

// MergeActionKind enumerates the THEN action of a MERGE WHEN clause.
type MergeActionKind int

const (
	// MergeInsert is `INSERT [(cols)] (VALUES (row) | ROW)`.
	MergeInsert MergeActionKind = iota
	// MergeUpdate is `UPDATE SET update_item_list`.
	MergeUpdate
	// MergeDelete is `DELETE`.
	MergeDelete
)

// MergeAction is the THEN action of a MERGE WHEN clause (merge_action):
//
//	INSERT [column_list] (VALUES insert_values_row | ROW)
//	UPDATE SET update_item_list
//	DELETE
//
// Kind selects the action. For MergeInsert: Columns is the optional column list;
// InsertRow holds the VALUES row (nil when the bare `ROW` form was used), and
// SourceRow is true for `ROW`. For MergeUpdate: SetItems is the SET list.
// DELETE carries no payload.
type MergeAction struct {
	Kind      MergeActionKind
	Columns   []string      // INSERT column list; nil if absent
	InsertRow *InsertRow    // INSERT VALUES row; nil for `ROW` / non-insert
	SourceRow bool          // INSERT ... ROW (source-row insert)
	SetItems  []*UpdateItem // UPDATE SET items; nil for non-update
	Loc       Loc
}

// Tag implements Node.
func (n *MergeAction) Tag() NodeTag { return T_MergeAction }

// TruncateStmt is a TRUNCATE TABLE statement (truncate_statement):
//
//	TRUNCATE TABLE <path> [WHERE expr]
//
// Target is the table path (maybe_dashed_path_expression). Where is the optional
// filter (ZetaSQL accepts `TRUNCATE TABLE foo WHERE bar > 3`). TRUNCATE is a
// BigQuery DML statement; Spanner has no TRUNCATE (its DDL path syntax-rejects
// it — non-authoritative), so TRUNCATE forms are oracled against the ZetaSQL
// corpus + the BigQuery docs.
type TruncateStmt struct {
	Target *PathExpr // table path
	Where  Node      // WHERE expr; nil if absent
	Loc    Loc
}

// Tag implements Node.
func (n *TruncateStmt) Tag() NodeTag { return T_TruncateStmt }

// Returning is a THEN RETURN clause on a DML statement (opt_returning_clause):
//
//	THEN RETURN [WITH ACTION [AS action_alias]] <select_list>
//
// WithAction records the `WITH ACTION` modifier (the returned action column);
// ActionAlias is its optional `AS alias` ("" if absent / no WITH ACTION). Items
// is the returned select_list (the same SelectItem shape as a SELECT list, so
// `*`, `* EXCEPT (…)`, `expr AS x` are all supported).
type Returning struct {
	WithAction  bool          // WITH ACTION
	ActionAlias string        // AS action_alias; "" if absent
	Items       []*SelectItem // returned select list
	Loc         Loc
}

// Tag implements Node.
func (n *Returning) Tag() NodeTag { return T_Returning }

// Compile-time assertions that the DML node types satisfy Node.
var (
	_ Node = (*DefaultExpr)(nil)
	_ Node = (*InsertStmt)(nil)
	_ Node = (*InsertRow)(nil)
	_ Node = (*InsertTable)(nil)
	_ Node = (*OnConflict)(nil)
	_ Node = (*UpdateStmt)(nil)
	_ Node = (*UpdateItem)(nil)
	_ Node = (*DeleteStmt)(nil)
	_ Node = (*MergeStmt)(nil)
	_ Node = (*MergeWhen)(nil)
	_ Node = (*MergeAction)(nil)
	_ Node = (*TruncateStmt)(nil)
	_ Node = (*Returning)(nil)
)

// ===========================================================================
// Transactions / batch (parser-utility node)
// ===========================================================================
//
// These node types model the legacy ANTLR transaction / batch statement family
// (GoogleSQLParser.g4 §2.9), itself a hand-port of Google's open-source ZetaSQL
// reference and the grammar bytebase consumes today:
//
//	begin_statement:        begin_transaction_keywords transaction_mode_list?
//	begin_transaction_keywords: START TRANSACTION | BEGIN TRANSACTION?
//	commit_statement:       COMMIT TRANSACTION?
//	rollback_statement:     ROLLBACK TRANSACTION?
//	transaction_mode:       READ ONLY | READ WRITE
//	                        | ISOLATION LEVEL id | ISOLATION LEVEL id id
//	start_batch_statement:  START BATCH identifier?
//	run_batch_statement:    RUN BATCH
//	abort_batch_statement:  ABORT BATCH
//
// Adjudication — the live Cloud Spanner emulator (oracle.md) parses every one of
// these (it accepts the leading form, then feature-rejects with "Statement not
// supported: BeginStatement / CommitStatement / …"), so the LEADING-FORM accept
// is oracle-authoritative. The PRECISE trailing grammar, however, is NOT: the
// emulator's recognizer for these keywords swallows arbitrary trailing tokens
// (it accepts `COMMIT WORK`, `START BATCH a b`, even `COMMIT garbage !!!`). These
// nodes therefore follow the ZetaSQL .g4 (the grammar bytebase consumes), which
// caps the tail at the documented forms; `COMMIT WORK` and a multi-word
// `START BATCH a b` are syntax errors here (divergence ledger: Spanner's shallow
// recognizer is non-authoritative for the trailing tokens of these statements).

// TransactionStmtKind discriminates the four BEGIN/COMMIT/ROLLBACK shapes that
// share TransactionStmt.
type TransactionStmtKind int

const (
	// TransactionBegin is `BEGIN [TRANSACTION]` (the begin_transaction_keywords
	// BEGIN alternative).
	TransactionBegin TransactionStmtKind = iota
	// TransactionStart is `START TRANSACTION` (the begin_transaction_keywords
	// START alternative). Modeled distinctly from BEGIN so the source verb
	// round-trips, even though both are begin_statement.
	TransactionStart
	// TransactionCommit is `COMMIT [TRANSACTION]`.
	TransactionCommit
	// TransactionRollback is `ROLLBACK [TRANSACTION]`.
	TransactionRollback
)

// String returns a human-readable name for the transaction-statement kind.
func (k TransactionStmtKind) String() string {
	switch k {
	case TransactionBegin:
		return "BEGIN"
	case TransactionStart:
		return "START TRANSACTION"
	case TransactionCommit:
		return "COMMIT"
	case TransactionRollback:
		return "ROLLBACK"
	default:
		return "UNKNOWN"
	}
}

// TransactionModeKind discriminates the transaction_mode alternatives.
type TransactionModeKind int

const (
	// TransactionModeReadOnly is `READ ONLY`.
	TransactionModeReadOnly TransactionModeKind = iota
	// TransactionModeReadWrite is `READ WRITE`.
	TransactionModeReadWrite
	// TransactionModeIsolationLevel is `ISOLATION LEVEL <id> [<id>]` — the
	// isolation-level words are kept verbatim in TransactionMode.Levels (1 or 2
	// identifiers, e.g. ["SERIALIZABLE"] or ["REPEATABLE","READ"]).
	TransactionModeIsolationLevel
)

// TransactionMode is one entry of a transaction_mode_list — the READ ONLY /
// READ WRITE / ISOLATION LEVEL modes that may follow BEGIN / START TRANSACTION.
// For the ISOLATION LEVEL kind, Levels holds the 1-or-2 isolation-level words
// (the grammar's `identifier` / `identifier identifier`); it is empty for the
// READ modes.
type TransactionMode struct {
	Kind   TransactionModeKind
	Levels []string // ISOLATION LEVEL words (1 or 2); empty otherwise
	Loc    Loc
}

// Tag implements Node.
func (n *TransactionMode) Tag() NodeTag { return T_TransactionMode }

// TransactionStmt is a BEGIN / START TRANSACTION / COMMIT / ROLLBACK statement
// (grammar: begin_statement / commit_statement / rollback_statement). Kind
// selects the verb; Transaction records whether the optional TRANSACTION keyword
// was present (it is always implied for START TRANSACTION); Modes is the
// transaction_mode_list (only ever non-empty for BEGIN / START TRANSACTION).
type TransactionStmt struct {
	Kind        TransactionStmtKind
	Transaction bool               // the TRANSACTION keyword was present
	Modes       []*TransactionMode // transaction_mode_list; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *TransactionStmt) Tag() NodeTag { return T_TransactionStmt }

// BatchStmtKind discriminates the three batch statements.
type BatchStmtKind int

const (
	// BatchStart is `START BATCH [identifier]`.
	BatchStart BatchStmtKind = iota
	// BatchRun is `RUN BATCH`.
	BatchRun
	// BatchAbort is `ABORT BATCH`.
	BatchAbort
)

// String returns a human-readable name for the batch-statement kind.
func (k BatchStmtKind) String() string {
	switch k {
	case BatchStart:
		return "START BATCH"
	case BatchRun:
		return "RUN BATCH"
	case BatchAbort:
		return "ABORT BATCH"
	default:
		return "UNKNOWN"
	}
}

// BatchStmt is a START BATCH / RUN BATCH / ABORT BATCH statement (grammar:
// start_batch_statement / run_batch_statement / abort_batch_statement). Kind
// selects the verb; Name is the optional batch-type identifier of START BATCH
// (e.g. `ddl`), "" for RUN/ABORT or a bare START BATCH.
type BatchStmt struct {
	Kind BatchStmtKind
	Name string // START BATCH <id>; "" otherwise
	Loc  Loc
}

// Tag implements Node.
func (n *BatchStmt) Tag() NodeTag { return T_BatchStmt }

// ===========================================================================
// Utility / metadata statements (parser-utility node)
// ===========================================================================
//
// ASSERT / ANALYZE / DESCRIBE / RENAME / CALL — the legacy ANTLR utility
// statement family (GoogleSQLParser.g4 §2.10 + the §2.3 rename_statement),
// itself a hand-port of Google's open-source ZetaSQL reference:
//
//	assert_statement:   ASSERT expression opt_description?       (opt_description: AS string_literal)
//	analyze_statement:  ANALYZE opt_options_list? table_and_column_info_list?
//	table_and_column_info: maybe_dashed_path_expression column_list?
//	describe_statement: describe_keyword describe_info
//	describe_info:      identifier? maybe_slashed_or_dashed_path_expression opt_from_path_expression?
//	rename_statement:   RENAME identifier path_expression TO path_expression
//	call_statement:     CALL path_expression '(' (tvf_argument (',' tvf_argument)*)? ')'
//
// Adjudication (oracle.md): the Spanner emulator parses ASSERT / DESCRIBE via its
// shallow recognizer (LEADING-FORM accept authoritative; trailing tokens
// swallowed → non-authoritative, follow the .g4); RENAME and bare ANALYZE go
// through its real DDL parser (authoritative — `RENAME TABLE a TO b` accepts,
// bare `ANALYZE` accepts). CALL is parsed in full by the emulator (authoritative
// — it validates the argument list: `CALL p(1 => 2)` and `CALL p(,)` reject).
// Spanner-specific narrowings vs the union grammar are flagged divergences:
//   - ANALYZE with targets (`ANALYZE t`) — Spanner rejects (its ANALYZE is the
//     bare whole-database form); the .g4 + BigQuery union accept the targets.
//   - RENAME with a non-TABLE object kind — Spanner only allows `RENAME TABLE`;
//     the .g4 allows any object-kind identifier.

// AssertStmt is an ASSERT statement (grammar: assert_statement):
//
//	ASSERT <expression> [AS <description-string>]
//
// Expr is the asserted boolean expression. Description is the optional `AS
// 'message'` string-literal content (opt_description: AS string_literal); it is
// nil when no AS clause is present. Per the .g4 the description MUST be a string
// literal (the live emulator's shallow recognizer accepts a non-string there, but
// that is the non-authoritative trailing-token behavior — see the node header).
type AssertStmt struct {
	Expr        Node // the asserted expression
	Description Node // AS string_literal; nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *AssertStmt) Tag() NodeTag { return T_AssertStmt }

// TableAndColumnInfo is one target of an ANALYZE statement (grammar:
// table_and_column_info: maybe_dashed_path_expression column_list?). Path is the
// (possibly dashed) table path; Columns is the optional parenthesized column
// list, nil when absent.
type TableAndColumnInfo struct {
	Path    *PathExpr // table path
	Columns []string  // optional column_list; nil if no '(' ... ')'
	Loc     Loc
}

// Tag implements Node.
func (n *TableAndColumnInfo) Tag() NodeTag { return T_TableAndColumnInfo }

// AnalyzeStmt is an ANALYZE statement (grammar: analyze_statement):
//
//	ANALYZE [OPTIONS (...)] [target [, target ...]]
//
// Options is the optional OPTIONS(...) list (nil if absent). Targets is the
// optional table_and_column_info_list (nil for a bare `ANALYZE`, which analyzes
// the whole database — the only form the Spanner emulator accepts; targets are a
// BigQuery-union form, see the node header).
type AnalyzeStmt struct {
	Options []*OptionsEntry       // OPTIONS(...); nil if absent
	Targets []*TableAndColumnInfo // table_and_column_info_list; nil if absent
	Loc     Loc
}

// Tag implements Node.
func (n *AnalyzeStmt) Tag() NodeTag { return T_AnalyzeStmt }

// DescribeStmt is a DESCRIBE / DESC statement (grammar: describe_statement +
// describe_info):
//
//	{ DESCRIBE | DESC } [<object-type>] <path> [FROM <path>]
//
// IsDesc records whether the abbreviated DESC keyword was used (vs DESCRIBE).
// ObjectType is the optional leading object-type identifier (describe_info's
// `identifier?`, e.g. `TABLE` / `INDEX`), "" if absent. Path is the described
// object's (possibly slashed/dashed) path. FromPath is the optional
// `FROM <path>` qualifier (opt_from_path_expression), nil if absent.
type DescribeStmt struct {
	IsDesc     bool      // DESC (true) vs DESCRIBE (false)
	ObjectType string    // optional leading object-type word; "" if absent
	Path       *PathExpr // the described object path
	FromPath   *PathExpr // FROM <path>; nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *DescribeStmt) Tag() NodeTag { return T_DescribeStmt }

// RenameStmt is a RENAME statement (grammar: rename_statement):
//
//	RENAME <object-kind> <path> TO <path>
//
// ObjectType is the object-kind identifier (the grammar's `identifier`, e.g.
// `TABLE`). From is the current name path; To is the new name path.
type RenameStmt struct {
	ObjectType string    // object-kind word (e.g. TABLE)
	From       *PathExpr // current name
	To         *PathExpr // new name
	Loc        Loc
}

// Tag implements Node.
func (n *RenameStmt) Tag() NodeTag { return T_RenameStmt }

// CallArgKind discriminates the non-expression shapes of a CALL argument
// (tvf_argument). A plain expression / named argument is carried directly as the
// argument Node (an expression or *NamedArg); the table / model / connection /
// descriptor clauses are wrapped in a CallArg so the kind is explicit.
type CallArgKind int

const (
	// CallArgTable is `TABLE <path>` (table_clause).
	CallArgTable CallArgKind = iota
	// CallArgModel is `MODEL <path>` (model_clause).
	CallArgModel
	// CallArgConnection is `CONNECTION <path>|DEFAULT` (connection_clause).
	CallArgConnection
	// CallArgDescriptor is `DESCRIPTOR (col, ...)` (descriptor_argument).
	CallArgDescriptor
)

// CallArg wraps a non-expression CALL argument (the table / model / connection /
// descriptor clauses of tvf_argument). Kind selects the shape; Path carries the
// referenced path for TABLE / MODEL / CONNECTION (nil for the CONNECTION DEFAULT
// form); Columns carries the descriptor column names for DESCRIPTOR.
type CallArg struct {
	Kind    CallArgKind
	Path    *PathExpr // TABLE/MODEL/CONNECTION path; nil for CONNECTION DEFAULT or DESCRIPTOR
	Default bool      // CONNECTION DEFAULT
	Columns []string  // DESCRIPTOR column names
	Loc     Loc
}

// Tag implements Node.
func (n *CallArg) Tag() NodeTag { return T_CallArg }

// CallStmt is a CALL statement (grammar: call_statement):
//
//	CALL <path> ( [arg [, arg ...]] )
//
// Proc is the procedure name path. Args is the (possibly empty) argument list;
// each entry is a plain expression, a *NamedArg (`name => value`), or a *CallArg
// (the TABLE / MODEL / CONNECTION / DESCRIPTOR clauses).
type CallStmt struct {
	Proc *PathExpr // procedure name path
	Args []Node    // tvf_argument list (expression | *NamedArg | *CallArg); may be empty
	Loc  Loc
}

// Tag implements Node.
func (n *CallStmt) Tag() NodeTag { return T_CallStmt }

// Compile-time assertions that the transaction / utility node types satisfy Node.
var (
	_ Node = (*TransactionStmt)(nil)
	_ Node = (*TransactionMode)(nil)
	_ Node = (*BatchStmt)(nil)
	_ Node = (*AssertStmt)(nil)
	_ Node = (*AnalyzeStmt)(nil)
	_ Node = (*TableAndColumnInfo)(nil)
	_ Node = (*DescribeStmt)(nil)
	_ Node = (*RenameStmt)(nil)
	_ Node = (*CallArg)(nil)
	_ Node = (*CallStmt)(nil)
)

// ===========================================================================
// BigQuery EXPORT / LOAD / CLONE DATA (googlesql/parser-dml-ext node)
// ===========================================================================
//
// The §2.8 data-movement statements of the legacy ANTLR GoogleSQLParser.g4
// (a hand-port of ZetaSQL):
//
//	export_data_statement:  EXPORT DATA with_connection_clause? opt_options_list? AS query
//	export_model_statement: EXPORT MODEL path_expression with_connection_clause? opt_options_list?
//	aux_load_data_statement: LOAD DATA (INTO|OVERWRITE) maybe_dashed_path_expression_with_scope
//	    table_element_list? load_data_partitions_clause? collate_clause?
//	    partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint? opt_options_list?
//	    aux_load_data_from_files_options_list opt_external_table_with_clauses?
//	clone_data_statement:   CLONE DATA INTO maybe_dashed_path_expression FROM clone_data_source_list
//
// All four are BigQuery-ONLY at the GoogleSQL union level (oracle.md): the
// Spanner emulator hard-syntax-rejects them, so the union parser accepts them on
// the authority of the legacy .g4 + the BigQuery truth1 corpus (OTHER-001 EXPORT
// DATA, OTHER-003 LOAD DATA). bytebase does not consume these today (P1 parity),
// but the parser must accept-and-model them so Diagnose reports them as valid
// rather than a false syntax error, and GetQuerySpan can reach the embedded
// query of EXPORT DATA / the table targets of LOAD/CLONE.

// ExportDataStmt is `EXPORT DATA [WITH CONNECTION conn] [OPTIONS(...)] AS query`
// (export_data_statement). The query is required (export_data_statement =
// export_data_no_query as_query; there is no no-query form at statement level).
type ExportDataStmt struct {
	HasConnection  bool   // a WITH CONNECTION clause was present
	ConnectionName string // connection path text, or "DEFAULT"; "" if absent
	Options        []*OptionsEntry
	Query          Node // AS query (required); *QueryStmt
	Loc            Loc
}

// Tag implements Node.
func (n *ExportDataStmt) Tag() NodeTag { return T_ExportDataStmt }

// ExportModelStmt is `EXPORT MODEL <path> [WITH CONNECTION conn] [OPTIONS(...)]`
// (export_model_statement).
type ExportModelStmt struct {
	Name           *PathExpr // model path (path_expression)
	HasConnection  bool
	ConnectionName string // connection path text, or "DEFAULT"; "" if absent
	Options        []*OptionsEntry
	Loc            Loc
}

// Tag implements Node.
func (n *ExportModelStmt) Tag() NodeTag { return T_ExportModelStmt }

// LoadDataStmt is a `LOAD DATA { INTO | OVERWRITE } …` statement
// (aux_load_data_statement). It loads external files (FROM FILES (options)) into
// a destination table, which may be a TEMP/TEMPORARY table.
type LoadDataStmt struct {
	Overwrite   bool      // true => OVERWRITE; false => INTO (append_or_overwrite)
	Temp        bool      // a TEMP / TEMPORARY TABLE scope prefix was present
	TempKeyword string    // "TEMP" | "TEMPORARY" | "" — the source spelling of the scope keyword
	Name        *PathExpr // destination table (maybe_dashed_path_expression)
	// Optional explicit table schema (table_element_list). The grammar admits both
	// column definitions and table constraints; we keep them in two slices like
	// CreateTableStmt (the original interleaving is not load-bearing here).
	Columns             []*ColumnDef
	Constraints         []*TableConstraint
	PartitionsOverwrite bool   // OVERWRITE inside the load_data_partitions_clause
	Partitions          Node   // OVERWRITE? PARTITIONS ( expr ); nil if absent
	Collate             string // collate_clause spelling; "" if absent
	PartitionBy         []Node // PARTITION BY expr (, expr)*; nil if absent
	ClusterBy           []Node // CLUSTER BY expr (, expr)*; nil if absent
	Options             []*OptionsEntry
	FromFiles           []*OptionsEntry // FROM FILES ( options_list ) — required
	// opt_external_table_with_clauses: WITH PARTITION COLUMNS [(cols)] and/or
	// WITH CONNECTION conn.
	HasPartitionColumns bool         // a WITH PARTITION COLUMNS clause was present
	PartitionColumns    []*ColumnDef // its optional explicit column list; nil = inferred
	HasConnection       bool         // a WITH CONNECTION clause was present
	ConnectionName      string       // connection path text, or "DEFAULT"; "" if absent
	Loc                 Loc
}

// Tag implements Node.
func (n *LoadDataStmt) Tag() NodeTag { return T_LoadDataStmt }

// CloneDataStmt is `CLONE DATA INTO <path> FROM <source> (UNION ALL <source>)*`
// (clone_data_statement). It copies rows from one or more source tables (each
// optionally at a system-time / filtered by WHERE) into the destination.
type CloneDataStmt struct {
	Name    *PathExpr // destination table (maybe_dashed_path_expression)
	Sources []*CloneDataSource
	Loc     Loc
}

// Tag implements Node.
func (n *CloneDataStmt) Tag() NodeTag { return T_CloneDataStmt }

// CloneDataSource is one entry of a clone_data_source_list:
// `maybe_dashed_path_expression opt_at_system_time? where_clause?`.
type CloneDataSource struct {
	Name          *PathExpr // source table
	ForSystemTime Node      // FOR SYSTEM_TIME AS OF expr; nil if absent
	Where         Node      // WHERE expr filter; nil if absent
	Loc           Loc
}

// Tag implements Node.
func (n *CloneDataSource) Tag() NodeTag { return T_CloneDataSource }

// Compile-time assertions that the data-movement node types satisfy Node.
var (
	_ Node = (*ExportDataStmt)(nil)
	_ Node = (*ExportModelStmt)(nil)
	_ Node = (*LoadDataStmt)(nil)
	_ Node = (*CloneDataStmt)(nil)
	_ Node = (*CloneDataSource)(nil)
)

// ===========================================================================
// Spanner-specific DDL (googlesql/parser-ddl-spanner node)
// ===========================================================================
//
// Cloud Spanner extends GoogleSQL with object kinds that the legacy ANTLR
// GoogleSQLParser.g4 does NOT model as first-class rules — it rides them on the
// generic-entity mechanism (or rejects them outright). The omni parser, whose
// authoritative oracle for Spanner DDL is the live Cloud Spanner emulator
// (oracle.md), parses these as dedicated statements so bytebase's Spanner
// consumers (Diagnose / GetQuerySpan / SplitSQL) see a real tree instead of a
// false "syntax error". The forms (all emulator-accept-verified) are:
//
//	CHANGE STREAM  — CREATE / ALTER / DROP CHANGE STREAM       (DDL-024/025/026)
//	SEQUENCE       — CREATE / ALTER / DROP SEQUENCE            (DDL-027/028/029)
//	ROLE           — CREATE / DROP ROLE                        (DDL-032/033)
//	role GRANT     — GRANT/REVOKE … TO/FROM ROLE; GRANT/REVOKE ROLE … (DDL-034-037)
//	LOCALITY GROUP — CREATE / ALTER / DROP LOCALITY GROUP      (DDL-041/042/043)
//	PROTO BUNDLE   — CREATE / ALTER PROTO BUNDLE               (DDL-046/047)
//
// These DIVERGE from the legacy grammar (which over-rejects them); the divergence
// is oracle-backed and recorded in the ledger. The role-based GRANT/REVOKE adds a
// GranteeRole grantee kind to the shared Grantee node and role-list fields to the
// shared GrantStmt/RevokeStmt (see those types), so the one union parser accepts
// both the legacy ZetaSQL string-grantee dialect and the Spanner role dialect.

// ChangeStreamTrackedTable is one entry of a CREATE/ALTER CHANGE STREAM's
// `FOR table_and_column` list (grammar truth1 DDL-024):
//
//	table_name                       -> AllColumns=false, Columns=nil (whole table, all columns)
//	table_name ( )                   -> ExplicitColumns=true, Columns=nil (track no non-key columns)
//	table_name ( col [, col ...] )   -> ExplicitColumns=true, Columns=[…]
//
// The empty `()` form is distinct from omitting the parens entirely (the former
// tracks only the primary key; the latter tracks the whole table), so a bare
// ExplicitColumns flag disambiguates `t` from `t()`.
type ChangeStreamTrackedTable struct {
	Name            *PathExpr // the tracked table name
	ExplicitColumns bool      // true when a `( … )` column list was given (even if empty)
	Columns         []string  // the listed column names (nil for `t` or `t()`)
	Loc             Loc
}

// Tag implements Node.
func (n *ChangeStreamTrackedTable) Tag() NodeTag { return T_ChangeStreamTrackedTable }

// CreateChangeStreamStmt is a CREATE CHANGE STREAM statement (Spanner; truth1
// DDL-024):
//
//	CREATE CHANGE STREAM name { FOR ALL | FOR table_and_column [, ...] }?
//	  [OPTIONS ( option [, ...] )]
//
// Exactly one of ForAll / ForTables describes the watched set when a FOR clause
// is present; a CHANGE STREAM created with neither (just OPTIONS, or nothing)
// watches nothing until ALTERed. HasFor records whether a FOR clause was present
// at all (to distinguish `FOR ALL` from "no FOR").
type CreateChangeStreamStmt struct {
	Name      *PathExpr
	HasFor    bool                        // a FOR clause was present
	ForAll    bool                        // FOR ALL
	ForTables []*ChangeStreamTrackedTable // FOR table_and_column [, ...]
	Options   []*OptionsEntry
	Loc       Loc
}

// Tag implements Node.
func (n *CreateChangeStreamStmt) Tag() NodeTag { return T_CreateChangeStreamStmt }

// AlterChangeStreamStmt is an ALTER CHANGE STREAM statement (Spanner; truth1
// DDL-025):
//
//	ALTER CHANGE STREAM name
//	  { SET FOR { ALL | table_and_column [, ...] }
//	  | DROP FOR ALL
//	  | SET OPTIONS ( option [, ...] )
//	  }
//
// Action selects which of the three forms this is.
type AlterChangeStreamStmt struct {
	Name      *PathExpr
	Action    ChangeStreamAlterAction
	ForAll    bool                        // for SetFor: SET FOR ALL
	ForTables []*ChangeStreamTrackedTable // for SetFor: SET FOR table_and_column [, ...]
	Options   []*OptionsEntry             // for SetOptions
	Loc       Loc
}

// Tag implements Node.
func (n *AlterChangeStreamStmt) Tag() NodeTag { return T_AlterChangeStreamStmt }

// ChangeStreamAlterAction discriminates the three ALTER CHANGE STREAM forms.
type ChangeStreamAlterAction int

const (
	// ChangeStreamSetFor is `SET FOR { ALL | table_and_column [, ...] }`.
	ChangeStreamSetFor ChangeStreamAlterAction = iota
	// ChangeStreamDropForAll is `DROP FOR ALL`.
	ChangeStreamDropForAll
	// ChangeStreamSetOptions is `SET OPTIONS ( … )`.
	ChangeStreamSetOptions
)

// DropChangeStreamStmt is a DROP CHANGE STREAM statement (Spanner; truth1
// DDL-026): `DROP CHANGE STREAM name`.
type DropChangeStreamStmt struct {
	Name *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *DropChangeStreamStmt) Tag() NodeTag { return T_DropChangeStreamStmt }

// CreateSequenceStmt is a CREATE SEQUENCE statement (Spanner; truth1 DDL-027):
//
//	CREATE SEQUENCE [IF NOT EXISTS] name [OPTIONS ( option [, ...] )]
type CreateSequenceStmt struct {
	Name        *PathExpr
	IfNotExists bool
	Options     []*OptionsEntry
	Loc         Loc
}

// Tag implements Node.
func (n *CreateSequenceStmt) Tag() NodeTag { return T_CreateSequenceStmt }

// AlterSequenceStmt is an ALTER SEQUENCE statement (Spanner; truth1 DDL-028):
//
//	ALTER SEQUENCE [IF EXISTS] name SET OPTIONS ( option [, ...] )
type AlterSequenceStmt struct {
	Name       *PathExpr
	IfExists   bool
	SetOptions []*OptionsEntry
	Loc        Loc
}

// Tag implements Node.
func (n *AlterSequenceStmt) Tag() NodeTag { return T_AlterSequenceStmt }

// DropSequenceStmt is a DROP SEQUENCE statement (Spanner; truth1 DDL-029):
//
//	DROP SEQUENCE [IF EXISTS] name
type DropSequenceStmt struct {
	Name     *PathExpr
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropSequenceStmt) Tag() NodeTag { return T_DropSequenceStmt }

// CreateRoleStmt is a CREATE ROLE statement (Spanner; truth1 DDL-032):
// `CREATE ROLE role_name`. role_name is a single identifier in the grammar.
type CreateRoleStmt struct {
	Name *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *CreateRoleStmt) Tag() NodeTag { return T_CreateRoleStmt }

// DropRoleStmt is a DROP ROLE statement (Spanner; truth1 DDL-033):
// `DROP ROLE role_name`.
type DropRoleStmt struct {
	Name *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *DropRoleStmt) Tag() NodeTag { return T_DropRoleStmt }

// CreateLocalityGroupStmt is a CREATE LOCALITY GROUP statement (Spanner; truth1
// DDL-041):
//
//	CREATE LOCALITY GROUP name [OPTIONS ( storage_option [, ...] )]
type CreateLocalityGroupStmt struct {
	Name    *PathExpr
	Options []*OptionsEntry
	Loc     Loc
}

// Tag implements Node.
func (n *CreateLocalityGroupStmt) Tag() NodeTag { return T_CreateLocalityGroupStmt }

// AlterLocalityGroupStmt is an ALTER LOCALITY GROUP statement (Spanner; truth1
// DDL-042):
//
//	ALTER LOCALITY GROUP name SET OPTIONS ( storage_option [, ...] )
type AlterLocalityGroupStmt struct {
	Name       *PathExpr
	SetOptions []*OptionsEntry
	Loc        Loc
}

// Tag implements Node.
func (n *AlterLocalityGroupStmt) Tag() NodeTag { return T_AlterLocalityGroupStmt }

// DropLocalityGroupStmt is a DROP LOCALITY GROUP statement (Spanner; truth1
// DDL-043): `DROP LOCALITY GROUP name`.
type DropLocalityGroupStmt struct {
	Name *PathExpr
	Loc  Loc
}

// Tag implements Node.
func (n *DropLocalityGroupStmt) Tag() NodeTag { return T_DropLocalityGroupStmt }

// CreateProtoBundleStmt is a CREATE PROTO BUNDLE statement (Spanner; truth1
// DDL-046):
//
//	CREATE PROTO BUNDLE ( proto_type_name [, ...] )
//
// Each proto_type_name is a (usually backtick-quoted, dotted) fully-qualified
// proto message name; captured as a NamePath so the dotted form survives.
type CreateProtoBundleStmt struct {
	Types []NamePath // the bundled proto type names (>= 1)
	Loc   Loc
}

// Tag implements Node.
func (n *CreateProtoBundleStmt) Tag() NodeTag { return T_CreateProtoBundleStmt }

// AlterProtoBundleStmt is an ALTER PROTO BUNDLE statement (Spanner; truth1
// DDL-047):
//
//	ALTER PROTO BUNDLE { INSERT ( … ) | UPDATE ( … ) | DELETE ( … ) } [, ...]
//
// The grammar allows any combination of INSERT/UPDATE/DELETE groups in one
// statement; each group's type list is collected here. An absent group is nil.
type AlterProtoBundleStmt struct {
	Insert []NamePath // INSERT ( proto_type_name [, ...] )
	Update []NamePath // UPDATE ( proto_type_name [, ...] )
	Delete []NamePath // DELETE ( proto_type_name [, ...] )
	Loc    Loc
}

// Tag implements Node.
func (n *AlterProtoBundleStmt) Tag() NodeTag { return T_AlterProtoBundleStmt }

// Compile-time assertions that the Spanner-DDL node types satisfy Node.
var (
	_ Node = (*ChangeStreamTrackedTable)(nil)
	_ Node = (*CreateChangeStreamStmt)(nil)
	_ Node = (*AlterChangeStreamStmt)(nil)
	_ Node = (*DropChangeStreamStmt)(nil)
	_ Node = (*CreateSequenceStmt)(nil)
	_ Node = (*AlterSequenceStmt)(nil)
	_ Node = (*DropSequenceStmt)(nil)
	_ Node = (*CreateRoleStmt)(nil)
	_ Node = (*DropRoleStmt)(nil)
	_ Node = (*CreateLocalityGroupStmt)(nil)
	_ Node = (*AlterLocalityGroupStmt)(nil)
	_ Node = (*DropLocalityGroupStmt)(nil)
	_ Node = (*CreateProtoBundleStmt)(nil)
	_ Node = (*AlterProtoBundleStmt)(nil)
)

// ---------------------------------------------------------------------------
// GQL — graph query language + CREATE PROPERTY GRAPH (parser-gql node)
// ---------------------------------------------------------------------------
//
// The §2.12 GQL sub-language and create_property_graph_statement from the
// legacy ANTLR grammar (GoogleSQLParser.g4) — a hand-port of ZetaSQL's
// GoogleSQL graph extension (BigQuery / Spanner Graph). One omni AST serves
// both dialects: the production accepts the BigQuery+Spanner union.
//
// ORACLE NOTE — GQL is BigQuery/Spanner-Graph syntax. The Spanner emulator's
// accept/reject is recorded by the differential harness and used where it does
// not disagree with the legacy .g4; the authoritative reference for the GQL
// grammar shape is the pinned legacy GoogleSQLParser.g4 itself (the truth1
// BigQuery corpus explicitly scopes GQL out — INDEX.md). See
// graph_query_oracle_test.go.

// GQLStmt is a top-level GQL graph query (gql_statement):
//
//	GRAPH <path> <graph_operation_block>
//
// Name is the graph name (path_expression). Blocks is the NEXT-separated list
// of composite query blocks (graph_operation_block: block (NEXT block)*); it has
// >= 1 entry. Each block is a *GraphLinearQuery or a *GraphSetOp.
type GQLStmt struct {
	Name   *PathExpr
	Blocks []Node // *GraphLinearQuery | *GraphSetOp; >= 1
	Loc    Loc
}

// Tag implements Node.
func (n *GQLStmt) Tag() NodeTag { return T_GQLStmt }

// GraphSetOp is a set-operation-joined chain of linear query operations
// (graph_composite_query_prefix): a sequence of >= 2 linear queries separated by
// set-operation metadata. Ops holds the linear queries (each a
// *GraphLinearQuery), and Metas holds the set-operation metadata BETWEEN them
// (len(Metas) == len(Ops)-1): Metas[i] joins Ops[i] and Ops[i+1].
type GraphSetOp struct {
	Ops   []Node            // *GraphLinearQuery; >= 2
	Metas []*GraphSetOpMeta // set-op between consecutive Ops; len == len(Ops)-1
	Loc   Loc
}

// Tag implements Node.
func (n *GraphSetOp) Tag() NodeTag { return T_GraphSetOp }

// GraphSetOpMeta is one set-operation metadata token pair
// (graph_set_operation_metadata: query_set_operation_type all_or_distinct).
// Op is "UNION" | "EXCEPT" | "INTERSECT"; Quantifier is "ALL" | "DISTINCT"
// (both are REQUIRED by the grammar). It is a plain value type, not a Node.
type GraphSetOpMeta struct {
	Op         string // "UNION" | "EXCEPT" | "INTERSECT"
	Quantifier string // "ALL" | "DISTINCT"
	Loc        Loc
}

// GraphLinearQuery is a graph_linear_query_operation:
//
//	[graph_linear_operator …] <graph_return_operator>
//
// Operators is the leading sequence of linear operators (MATCH / OPTIONAL MATCH
// / LET / FILTER / ORDER BY / PAGE / WITH / FOR / TABLESAMPLE), each a Node of
// the corresponding Graph*Op type; nil if there were none. Return is the
// trailing RETURN operator (required by the grammar) — a *GraphReturnOp.
type GraphLinearQuery struct {
	Operators []Node // Graph*Op operators; may be empty
	Return    *GraphReturnOp
	Loc       Loc
}

// Tag implements Node.
func (n *GraphLinearQuery) Tag() NodeTag { return T_GraphLinearQuery }

// GraphMatchOp is a (graph_match_operator / graph_optional_match_operator):
//
//	[OPTIONAL] MATCH [hint] <graph_pattern>
//
// Optional flags the OPTIONAL prefix. Pattern is the graph pattern matched.
type GraphMatchOp struct {
	Optional bool
	Pattern  *GraphPattern
	Loc      Loc
}

// Tag implements Node.
func (n *GraphMatchOp) Tag() NodeTag { return T_GraphMatchOp }

// GraphLetOp is a graph_let_operator (`LET <var defs>`). Vars has >= 1 entry.
type GraphLetOp struct {
	Vars []*GraphLetVar // >= 1
	Loc  Loc
}

// Tag implements Node.
func (n *GraphLetOp) Tag() NodeTag { return T_GraphLetOp }

// GraphLetVar is one variable definition of a LET operator
// (graph_let_variable_definition: identifier = expression).
type GraphLetVar struct {
	Name string
	Expr Node
	Loc  Loc
}

// Tag implements Node.
func (n *GraphLetVar) Tag() NodeTag { return T_GraphLetVar }

// GraphFilterOp is a graph_filter_operator:
//
//	FILTER WHERE <expr>   |   FILTER <expr>
//
// HasWhere records whether the WHERE keyword was present (`FILTER WHERE expr`)
// vs the bare `FILTER expr` form. Expr is the filter predicate.
type GraphFilterOp struct {
	HasWhere bool
	Expr     Node
	Loc      Loc
}

// Tag implements Node.
func (n *GraphFilterOp) Tag() NodeTag { return T_GraphFilterOp }

// GraphOrderByOp is a graph_order_by_operator (graph_order_by_clause):
//
//	ORDER [hint] BY <ordering_expression> (, …)
//
// Items reuses OrderItem (the shared `expr [COLLATE] [ASC|DESC|ASCENDING|
// DESCENDING] [NULLS …]` ordering element); >= 1 entry.
type GraphOrderByOp struct {
	Items []*OrderItem // >= 1
	Loc   Loc
}

// Tag implements Node.
func (n *GraphOrderByOp) Tag() NodeTag { return T_GraphOrderByOp }

// GraphPageOp is a graph_page_operator (graph_page_clause):
//
//	(OFFSET|SKIP) <n> [LIMIT <n>]   |   LIMIT <n>
//
// Skip records `OFFSET <n>` / `SKIP <n>` (nil if the bare-LIMIT form was used).
// SkipIsSkip distinguishes the SKIP spelling from OFFSET. Limit is the LIMIT
// count (nil if a bare OFFSET/SKIP with no LIMIT).
type GraphPageOp struct {
	Skip       Node // OFFSET/SKIP count; nil if bare LIMIT
	SkipIsSkip bool // true if the SKIP keyword was used (vs OFFSET)
	Limit      Node // LIMIT count; nil if bare OFFSET/SKIP
	Loc        Loc
}

// Tag implements Node.
func (n *GraphPageOp) Tag() NodeTag { return T_GraphPageOp }

// GraphWithOp is a graph_with_operator:
//
//	WITH [ALL|DISTINCT] [hint] <return_item_list> [GROUP BY …]
//
// Quantifier is "ALL" | "DISTINCT" | "" (absent). Items is the projected return
// items (>= 1). GroupBy is the optional trailing GROUP BY clause (nil if absent)
// — it reuses the shared GroupByClause.
type GraphWithOp struct {
	Quantifier string             // "ALL" | "DISTINCT" | ""
	Items      []*GraphReturnItem // >= 1
	GroupBy    *GroupByClause     // nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *GraphWithOp) Tag() NodeTag { return T_GraphWithOp }

// GraphForOp is a graph_for_operator:
//
//	FOR <id> IN <expr> [WITH OFFSET [AS alias]]
//
// Name is the loop variable. Expr is the iterated expression. WithOffset records
// the `WITH OFFSET` suffix; OffsetAlias is its optional `AS alias` ("" if absent
// or no WITH OFFSET).
type GraphForOp struct {
	Name        string
	Expr        Node
	WithOffset  bool
	OffsetAlias string
	Loc         Loc
}

// Tag implements Node.
func (n *GraphForOp) Tag() NodeTag { return T_GraphForOp }

// GraphSampleOp is a graph_sample_clause:
//
//	TABLESAMPLE <method> ( <size> (ROWS|PERCENT) ) [WITH WEIGHT [AS id]] [REPEATABLE(n)]
//
// Method is the sampling method identifier (e.g. RESERVOIR / BERNOULLI). Size is
// the sample size expression; Unit is "ROWS" | "PERCENT". WithWeight records the
// `WITH WEIGHT` suffix; WeightAlias its optional `AS id` ("" if absent).
// Repeatable holds the REPEATABLE(n) seed expression (nil if absent).
type GraphSampleOp struct {
	Method      string
	Size        Node
	Unit        string // "ROWS" | "PERCENT"
	WithWeight  bool
	WeightAlias string
	Repeatable  Node // REPEATABLE(n); nil if absent
	Loc         Loc
}

// Tag implements Node.
func (n *GraphSampleOp) Tag() NodeTag { return T_GraphSampleOp }

// GraphReturnOp is a graph_return_operator:
//
//	RETURN [hint] [ALL|DISTINCT] <return_item_list> [GROUP BY] [ORDER BY] [PAGE]
//
// Quantifier is "ALL" | "DISTINCT" | "". Items is the return list (>= 1). GroupBy
// / OrderBy / Page are the optional trailing clauses (nil/absent if not present).
type GraphReturnOp struct {
	Quantifier string             // "ALL" | "DISTINCT" | ""
	Items      []*GraphReturnItem // >= 1
	GroupBy    *GroupByClause     // nil if absent
	OrderBy    []*OrderItem       // nil if absent
	Page       *GraphPageOp       // nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *GraphReturnOp) Tag() NodeTag { return T_GraphReturnOp }

// GraphReturnItem is one entry of a graph_return_item_list (graph_return_item):
//
//	<expr> [AS alias]   |   *
//
// Star is true for the bare `*` form (Expr nil). Otherwise Expr is the projected
// expression and Alias the optional `AS name` ("" if absent).
type GraphReturnItem struct {
	Expr  Node   // nil for `*`
	Star  bool   // true for `*`
	Alias string // `AS alias`; "" if absent
	Loc   Loc
}

// Tag implements Node.
func (n *GraphReturnItem) Tag() NodeTag { return T_GraphReturnItem }

// GraphPattern is a graph_pattern (`<path_pattern_list> [WHERE]`). Paths has
// >= 1 path pattern; Where is the optional trailing WHERE predicate (nil if
// absent).
type GraphPattern struct {
	Paths []*GraphPathPattern // >= 1
	Where Node                // nil if absent
	Loc   Loc
}

// Tag implements Node.
func (n *GraphPattern) Tag() NodeTag { return T_GraphPattern }

// GraphPathPattern is a graph_path_pattern:
//
//	[<var> =] [(ANY|ALL) [SHORTEST]] [<mode> [PATH|PATHS]] <path_pattern_expr>
//
// PathVar is the optional `<id> =` path-variable assignment ("" if absent).
// Search is the optional search prefix ("ANY" | "ALL" | "ANY SHORTEST" |
// "ALL SHORTEST" | ""). Mode is the optional path-mode prefix ("WALK" | "TRAIL"
// | "SIMPLE" | "ACYCLIC", optionally suffixed " PATH"/" PATHS"; "" if absent).
// Factors is the sequence of path factors (node/edge/parenthesized patterns,
// each a Node) that make up the path expression (>= 1).
type GraphPathPattern struct {
	PathVar string // path-variable assignment id; "" if absent
	Search  string // "ANY" | "ALL" | "ANY SHORTEST" | "ALL SHORTEST" | ""
	Mode    string // "WALK"/"TRAIL"/"SIMPLE"/"ACYCLIC" [PATH|PATHS]; "" if absent
	Factors []Node // *GraphNodePattern | *GraphEdgePattern | *GraphPathPattern (parenthesized); >= 1
	// Where is the trailing `WHERE <expr>` of a parenthesized path pattern
	// (graph_parenthesized_path_pattern: '(' … graph_path_pattern where_clause? ')').
	// It is nil for a top-level (non-parenthesized) path pattern, whose WHERE
	// belongs to the enclosing graph_pattern instead.
	Where Node // parenthesized-path trailing WHERE; nil if absent
	Loc   Loc
}

// Tag implements Node.
func (n *GraphPathPattern) Tag() NodeTag { return T_GraphPathPattern }

// GraphNodePattern is a graph_node_pattern (`( <filler> )`). Filler carries the
// optional variable / label expression / property spec / inline WHERE.
type GraphNodePattern struct {
	Filler *GraphPatternFiller
	Loc    Loc
}

// Tag implements Node.
func (n *GraphNodePattern) Tag() NodeTag { return T_GraphNodePattern }

// GraphEdgeDirection enumerates the six graph_edge_pattern shapes (three
// abbreviated — EdgeAny / EdgeLeft / EdgeRight — and three full forms —
// EdgeLeftFull / EdgeRightFull / EdgeUndirectedFull).
type GraphEdgeDirection int

const (
	// EdgeAny is the undirected abbreviated edge `-` (MINUS_OPERATOR).
	EdgeAny GraphEdgeDirection = iota
	// EdgeLeft is the left-directed abbreviated edge `<-` (LT MINUS).
	EdgeLeft
	// EdgeRight is the right-directed abbreviated edge `->` (SUB_GT_BRACKET).
	EdgeRight
	// EdgeLeftFull is the left-directed full edge `<-[ filler ]-`.
	EdgeLeftFull
	// EdgeRightFull is the right-directed full edge `-[ filler ]->`.
	EdgeRightFull
	// EdgeUndirectedFull is the undirected full edge `-[ filler ]-` (the legacy
	// graph_edge_pattern's first alternative `LT_OPERATOR? MINUS … MINUS` with the
	// leading `<` ABSENT — distinct from EdgeLeftFull, which has it present).
	EdgeUndirectedFull
)

// GraphEdgePattern is a graph_edge_pattern. Direction selects the arrow shape;
// Filler is the bracketed `[ … ]` filler for the three full forms (EdgeLeftFull
// / EdgeRightFull / EdgeUndirectedFull) and nil for the three abbreviated forms
// (EdgeAny / EdgeLeft / EdgeRight).
type GraphEdgePattern struct {
	Direction GraphEdgeDirection
	Filler    *GraphPatternFiller // nil for abbreviated edges
	Loc       Loc
}

// Tag implements Node.
func (n *GraphEdgePattern) Tag() NodeTag { return T_GraphEdgePattern }

// GraphPatternFiller is a graph_element_pattern_filler — the shared interior of
// a node `( … )` or full-edge `[ … ]` pattern:
//
//	[hint] [<var>] [(IS|:) <label_expr>] [{ prop : expr, … }] [WHERE]
//
// Var is the optional element variable ("" if absent). Label is the optional
// label expression (nil if absent); LabelColon records whether the `:` spelling
// (vs `IS`) introduced it. Properties is the optional `{ … }` property
// specification (nil if absent). Where is the optional inline WHERE predicate
// (nil if absent).
type GraphPatternFiller struct {
	Var        string
	Label      *GraphLabelExpr // nil if absent
	LabelColon bool            // true if `:` spelling, false if `IS`
	Properties *GraphPropertySpec
	Where      Node // inline WHERE; nil if absent
	Loc        Loc
}

// Tag implements Node.
func (n *GraphPatternFiller) Tag() NodeTag { return T_GraphPatternFiller }

// GraphLabelKind enumerates a label_expression / label_primary shape.
type GraphLabelKind int

const (
	// LabelName is a bare label identifier (label_primary: identifier).
	LabelName GraphLabelKind = iota
	// LabelWildcard is the `%` any-label primary (MODULO_OPERATOR).
	LabelWildcard
	// LabelNot is `! <label>` (EXCLAMATION_OPERATOR label_expression).
	LabelNot
	// LabelAnd is `<label> & <label>` (BIT_AND_SYMBOL).
	LabelAnd
	// LabelOr is `<label> | <label>` (STROKE_SYMBOL).
	LabelOr
)

// GraphLabelExpr is a node of the label algebra (label_expression /
// label_primary): a bare label, the `%` wildcard, negation `!`, or the binary
// `&` / `|` combinators. Parens are dropped (parenthesized_label_expression is
// transparent — associativity is captured by the tree shape).
//
// For LabelName, Name holds the label identifier. For LabelNot, Operand is the
// negated sub-expression. For LabelAnd / LabelOr, Left and Right are the
// operands. LabelWildcard carries no children.
type GraphLabelExpr struct {
	Kind    GraphLabelKind
	Name    string          // for LabelName
	Operand *GraphLabelExpr // for LabelNot
	Left    *GraphLabelExpr // for LabelAnd / LabelOr
	Right   *GraphLabelExpr // for LabelAnd / LabelOr
	Loc     Loc
}

// Tag implements Node.
func (n *GraphLabelExpr) Tag() NodeTag { return T_GraphLabelExpr }

// GraphPropertySpec is a graph_property_specification (`{ id : expr, … }`).
// Properties has >= 1 entry.
type GraphPropertySpec struct {
	Properties []*GraphPropertyNameValue // >= 1
	Loc        Loc
}

// Tag implements Node.
func (n *GraphPropertySpec) Tag() NodeTag { return T_GraphPropertySpec }

// GraphPropertyNameValue is one entry of a property specification
// (graph_property_name_and_value: identifier : expression).
type GraphPropertyNameValue struct {
	Name string
	Expr Node
	Loc  Loc
}

// Tag implements Node.
func (n *GraphPropertyNameValue) Tag() NodeTag { return T_GraphPropertyNameValue }

// CreatePropertyGraphStmt is a CREATE PROPERTY GRAPH statement
// (create_property_graph_statement):
//
//	CREATE [OR REPLACE] PROPERTY GRAPH [IF NOT EXISTS] <path>
//	  [OPTIONS(…)] NODE TABLES ( <element_table>, … ) [EDGE TABLES ( … )]
//
// OrReplace / IfNotExists flag the prefixes. Name is the graph name. Options is
// the optional OPTIONS list. NodeTables is the required NODE TABLES element
// list (>= 1). EdgeTables is the optional EDGE TABLES element list (nil if
// absent).
type CreatePropertyGraphStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *PathExpr
	Options     []*OptionsEntry
	NodeTables  []*ElementTableDef // >= 1
	EdgeTables  []*ElementTableDef // nil if no EDGE TABLES clause
	Loc         Loc
}

// Tag implements Node.
func (n *CreatePropertyGraphStmt) Tag() NodeTag { return T_CreatePropertyGraphStmt }

// ElementTableDef is one entry of a NODE/EDGE TABLES list
// (element_table_definition):
//
//	<path> [AS alias] [KEY ( cols )]
//	  [SOURCE KEY ( cols ) REFERENCES <node> [( cols )]]
//	  [DESTINATION KEY ( cols ) REFERENCES <node> [( cols )]]
//	  [<label-and-properties> | <properties-clause>]
//
// Name is the underlying table path. Alias is the optional `AS alias` ("" if
// absent). Key is the optional `KEY ( … )` column list (nil if absent). Source /
// Dest are the optional SOURCE / DESTINATION node-table references (edge tables;
// nil if absent). Labels is the label-and-properties list (the
// opt_label_and_properties_clause); nil if absent.
type ElementTableDef struct {
	Name   *PathExpr
	Alias  string
	Key    []string         // KEY ( cols ); nil if absent
	Source *ElementTableRef // SOURCE KEY … REFERENCES …; nil if absent
	Dest   *ElementTableRef // DESTINATION KEY … REFERENCES …; nil if absent
	Labels []*LabelAndProperties
	Loc    Loc
}

// Tag implements Node.
func (n *ElementTableDef) Tag() NodeTag { return T_ElementTableDef }

// ElementTableRef is a SOURCE / DESTINATION node-table reference of an edge
// table (opt_source_node_table_clause / opt_dest_node_table_clause):
//
//	(SOURCE|DESTINATION) KEY ( cols ) REFERENCES <node> [( cols )]
//
// Columns is the local edge-table key column list. Node is the referenced node
// table name. RefColumns is the optional referenced-column list (nil if absent).
// It is a plain value type, not a Node.
type ElementTableRef struct {
	Columns    []string // KEY ( cols )
	Node       string   // REFERENCES <node>
	RefColumns []string // optional ( cols ) after the node; nil if absent
	Loc        Loc
}

// LabelKind discriminates a properties_clause shape inside a label-and-properties
// or a bare element-table properties clause.
type LabelKind int

const (
	// LabelPropsAllColumns is `PROPERTIES [ARE] ALL COLUMNS [EXCEPT ( cols )]`.
	LabelPropsAllColumns LabelKind = iota
	// LabelPropsNone is `NO PROPERTIES`.
	LabelPropsNone
	// LabelPropsList is `PROPERTIES ( <derived_property>, … )`.
	LabelPropsList
	// LabelPropsDefault is no properties clause on this label entry.
	LabelPropsDefault
)

// LabelAndProperties is one entry of a label_and_properties_list, OR the bare
// properties_clause form of opt_label_and_properties_clause (in which case
// LabelName is "" and Default is false):
//
//	[DEFAULT] LABEL <id> [<properties_clause>]
//
// Default flags the `DEFAULT LABEL` form. LabelName is the label identifier (""
// for the bare-properties form). Kind selects the properties_clause shape;
// PropsList holds the derived properties for LabelPropsList; ExceptColumns holds
// the `EXCEPT ( cols )` list for LabelPropsAllColumns (nil if none).
type LabelAndProperties struct {
	Default       bool
	LabelName     string // "" for the bare opt_label_and_properties_clause properties form
	Kind          LabelKind
	PropsList     []*DerivedProperty // for LabelPropsList
	ExceptColumns []string           // EXCEPT ( cols ) for LabelPropsAllColumns; nil if none
	Loc           Loc
}

// Tag implements Node.
func (n *LabelAndProperties) Tag() NodeTag { return T_LabelAndProperties }

// DerivedProperty is one entry of a derived_property_list
// (derived_property: expression [AS alias]).
type DerivedProperty struct {
	Expr  Node
	Alias string // `AS alias`; "" if absent
	Loc   Loc
}

// Tag implements Node.
func (n *DerivedProperty) Tag() NodeTag { return T_DerivedProperty }

// Compile-time assertions that the GQL / property-graph node types satisfy Node.
var (
	_ Node = (*GQLStmt)(nil)
	_ Node = (*GraphSetOp)(nil)
	_ Node = (*GraphLinearQuery)(nil)
	_ Node = (*GraphMatchOp)(nil)
	_ Node = (*GraphLetOp)(nil)
	_ Node = (*GraphLetVar)(nil)
	_ Node = (*GraphFilterOp)(nil)
	_ Node = (*GraphOrderByOp)(nil)
	_ Node = (*GraphPageOp)(nil)
	_ Node = (*GraphWithOp)(nil)
	_ Node = (*GraphForOp)(nil)
	_ Node = (*GraphSampleOp)(nil)
	_ Node = (*GraphReturnOp)(nil)
	_ Node = (*GraphReturnItem)(nil)
	_ Node = (*GraphPattern)(nil)
	_ Node = (*GraphPathPattern)(nil)
	_ Node = (*GraphNodePattern)(nil)
	_ Node = (*GraphEdgePattern)(nil)
	_ Node = (*GraphPatternFiller)(nil)
	_ Node = (*GraphPropertySpec)(nil)
	_ Node = (*GraphPropertyNameValue)(nil)
	_ Node = (*GraphLabelExpr)(nil)
	_ Node = (*CreatePropertyGraphStmt)(nil)
	_ Node = (*ElementTableDef)(nil)
	_ Node = (*LabelAndProperties)(nil)
	_ Node = (*DerivedProperty)(nil)
)
