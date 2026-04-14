package parser

// simpleTypenameLeadTokens is the single source of truth for tokens that
// unambiguously start a SimpleTypename production (excluding GenericType
// which is covered by the type_function_name category check).
//
// Grammar reference: SimpleTypename in postgres/src/backend/parser/gram.y:14339.
//
// DO NOT extend this list without also extending parseSimpleTypename's
// dispatch switch in type.go, and vice versa. The oracle test
// TestSimpleTypenameLeadTokensMatchPG enforces this invariant against
// PG 17.
var simpleTypenameLeadTokens = []int{
	INT_P, INTEGER, SMALLINT, BIGINT, REAL, FLOAT_P, DOUBLE_P,
	DECIMAL_P, DEC, NUMERIC,
	BIT, CHARACTER, CHAR_P, VARCHAR, NATIONAL, NCHAR,
	BOOLEAN_P, JSON,
	TIMESTAMP, TIME, INTERVAL,
}

var simpleTypenameLeadSet = func() map[int]bool {
	m := make(map[int]bool, len(simpleTypenameLeadTokens))
	for _, t := range simpleTypenameLeadTokens {
		m[t] = true
	}
	return m
}()

// isSimpleTypenameStart reports whether the current token can start a
// SimpleTypename production. This is the single source of truth for all
// SimpleTypename FIRST-set probes. Callers must use this function rather
// than hand-writing token case clusters.
func (p *Parser) isSimpleTypenameStart() bool {
	if simpleTypenameLeadSet[p.cur.Type] {
		return true
	}
	// GenericType fallthrough: type_function_name opt_type_modifiers
	return p.isTypeFunctionName()
}

// isConstTypenameStart reports whether the current token starts a
// type-cast literal at PG's AexprConst grammar position. This is the
// FIRST set used by isAExprConstTypeCast / parseTypeCastedConst, which
// handles the union of two PG productions:
//
//	ConstTypename:  Numeric | ConstBit | ConstCharacter | ConstDatetime | JsonType  (gram.y:14370)
//	ConstInterval:  INTERVAL                                                          (gram.y:14688)
//
// Strictly speaking, INTERVAL is in ConstInterval, not ConstTypename —
// PG's AexprConst rule has both `ConstTypename Sconst` and
// `ConstInterval Sconst` as separate alternatives. omni's parser routes
// both alternatives through the same parseTypeCastedConst entry point
// (see expr.go), so the predicate's FIRST set is the union. Hence the
// "Const" in the name refers to the broader "type-cast literal lead"
// rather than the strict ConstTypename production.
//
// In omni this union coincides exactly with simpleTypenameLeadSet
// (SimpleTypename = the union plus GenericType fallback plus a couple
// of historical bits), so we reuse simpleTypenameLeadSet rather than
// maintaining a second slice. The only behavioral difference between
// isSimpleTypenameStart and isConstTypenameStart is the GenericType /
// type_function_name fallthrough — SimpleTypename has it, this predicate
// does not, because typed-string-literal contexts do not accept
// arbitrary identifier-typed names.
//
// NOTE: this predicate CANNOT be validated via a PG oracle at the
// AexprConst grammar position because that position is ambiguous with
// `func_name Sconst`, which accepts almost any keyword. See
// TestIsConstTypenameStartRejectsNonTypeStarters for the negative
// coverage that replaces the oracle test.
func (p *Parser) isConstTypenameStart() bool {
	return simpleTypenameLeadSet[p.cur.Type]
}

// isTypenameStart reports whether the current token starts a Typename
// production.
//
//	Typename:
//	    SimpleTypename opt_array_bounds
//	    | SETOF SimpleTypename opt_array_bounds
//	    | SimpleTypename ARRAY '[' Iconst ']'
//	    | SETOF SimpleTypename ARRAY '[' Iconst ']'
//	    | SimpleTypename ARRAY
//	    | SETOF SimpleTypename ARRAY
//
// Composed from isSimpleTypenameStart plus SETOF. There is no
// typenameLeadTokens slice — the FIRST set is fully expressible as
// {SETOF} ∪ FIRST(SimpleTypename), and we reuse the latter.
//
// At the FIRST-set level, FIRST(Typename) equals FIRST(func_type),
// because func_type's %TYPE alternatives all start with
// type_function_name, which is already in FIRST(Typename) via
// SimpleTypename's GenericType path. The oracle test below probes the
// RETURNS position of CREATE FUNCTION, which is grammatically func_type,
// and that's a sound oracle for Typename's FIRST set.
func (p *Parser) isTypenameStart() bool {
	if p.cur.Type == SETOF {
		return true
	}
	return p.isSimpleTypenameStart()
}

// ---------------------------------------------------------------------------
// Token-variant predicates: same semantics as the receiver-on-cur versions
// above, but operate on an arbitrary Token. Used by callers that need to
// peek at lookahead tokens (e.g. parseFuncArg's peek-then-commit
// disambiguation) without consuming them.
//
// These are thin parallel implementations rather than wrappers because the
// receiver versions read p.cur.Type and p.cur.Str, and we don't want to
// briefly mutate p.cur to reuse them — that would defeat the point of
// peeking.
// ---------------------------------------------------------------------------

// isSimpleTypenameStartToken is the Token-variant of isSimpleTypenameStart.
func (p *Parser) isSimpleTypenameStartToken(tok Token) bool {
	if simpleTypenameLeadSet[tok.Type] {
		return true
	}
	return p.isTypeFunctionNameToken(tok)
}

// isTypenameStartToken is the Token-variant of isTypenameStart.
func (p *Parser) isTypenameStartToken(tok Token) bool {
	if tok.Type == SETOF {
		return true
	}
	return p.isSimpleTypenameStartToken(tok)
}

// isFuncTypeStartToken is the Token-variant of isFuncTypeStart.
// FIRST(func_type) == FIRST(Typename) — see isFuncTypeStart's doc comment
// for the grammar-level proof.
func (p *Parser) isFuncTypeStartToken(tok Token) bool {
	return p.isTypenameStartToken(tok)
}

// aExprLeadTokens is the set of tokens that start an a_expr but are
// NOT covered by isConstTypenameStart or the isColId category check.
// Includes literal tokens, expression-opener keywords, the parameter
// reference, the SQL value functions (which use func_expr_common_subexpr,
// a distinct production from columnref), and unary punctuation.
//
// Grammar reference: a_expr in postgres/src/backend/parser/gram.y around
// line 14780. The FIRST set of a_expr is the most general expression
// FIRST set in PG; this slice captures the explicit token leads that
// don't fall through to the category or const-typename predicates.
//
// COLLATE is INTENTIONALLY EXCLUDED — it is a postfix operator
// (`a_expr COLLATE any_name`), not an expression starter. Adding it
// would make isAExprStart over-permissive.
//
// DO NOT extend without also extending the renderExpression continuation
// map in first_set_oracle_test.go. The oracle test
// TestAExprLeadTokensMatchPG enforces parity against PG 17.
var aExprLeadTokens = []int{
	// Literals (lex tokens, no category)
	ICONST, FCONST, SCONST, BCONST, XCONST,
	// Reserved keywords that are literal-like
	NULL_P, TRUE_P, FALSE_P,
	// Parameter reference
	PARAM,
	// Expression-opener keywords
	CASE, EXISTS, NOT, LEAST, GREATEST, COALESCE, NULLIF,
	CAST, ARRAY, ROW, INTERVAL,
	// SQL value functions / context constants. These are matched by
	// func_expr_common_subexpr in gram.y, a distinct production from the
	// columnref path, so they must appear explicitly in the lead set
	// rather than falling through isColId.
	CURRENT_TIMESTAMP, CURRENT_DATE, CURRENT_TIME,
	CURRENT_USER, SESSION_USER, USER,
	CURRENT_CATALOG, CURRENT_ROLE, CURRENT_SCHEMA, SYSTEM_USER,
	LOCALTIME, LOCALTIMESTAMP,
	// Punctuation openers (Go rune values)
	int('('), int('+'), int('-'),
}

var aExprLeadSet = func() map[int]bool {
	m := make(map[int]bool, len(aExprLeadTokens))
	for _, t := range aExprLeadTokens {
		m[t] = true
	}
	return m
}()

// isAExprStart reports whether the current token starts an a_expr
// production. Composed from:
//   - aExprLeadSet (explicit literal/keyword/punctuation/value-func leads)
//   - isConstTypenameStart (typed literals like `int '42'`)
//   - isColId() (identifier path via columnref)
//
// NOTE on isTypeFunctionName: PG's a_expr columnref alternative uses
// ColId, NOT type_function_name. The TypeFuncNameKeyword category
// contains join/operator keywords (INNER, LEFT, LIKE, ILIKE, IS,
// ISNULL, NOTNULL, OVERLAPS, CROSS, NATURAL, FULL, RIGHT, OUTER, JOIN,
// TABLESAMPLE, BINARY, FREEZE, CONCURRENTLY, COLLATION, AUTHORIZATION,
// SIMILAR, VERBOSE) that CANNOT start an expression. The TypeFuncName
// keywords that ARE SQL value functions (CURRENT_SCHEMA, and the four
// reserved ones like CURRENT_USER) have their own func_expr_common_subexpr
// rule and are listed explicitly in aExprLeadTokens. See
// TestAExprLeadTokensMatchPG for the oracle that locks this down.
//
// COLLATE is excluded — it's a postfix operator, not an expression lead.
//
// Validated against PG 17 by TestAExprLeadTokensMatchPG.
func (p *Parser) isAExprStart() bool {
	if aExprLeadSet[p.cur.Type] {
		return true
	}
	if p.isConstTypenameStart() {
		return true
	}
	return p.isColId()
}

// tableConstraintLeadTokens is the FIRST set of tokens that start a
// TableConstraint clause inside CREATE TABLE / ALTER TABLE.
//
// Grammar reference: TableConstraint in postgres/src/backend/parser/gram.y.
//
//	TableConstraint:
//	    CONSTRAINT name ConstraintElem
//	    | ConstraintElem  // CHECK | UNIQUE | PRIMARY KEY | FOREIGN KEY | EXCLUDE
//
// DO NOT extend this list without also extending parseTableConstraint's
// dispatch in create_table.go, and vice versa.
var tableConstraintLeadTokens = []int{
	CONSTRAINT, CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE,
}

var tableConstraintLeadSet = func() map[int]bool {
	m := make(map[int]bool, len(tableConstraintLeadTokens))
	for _, t := range tableConstraintLeadTokens {
		m[t] = true
	}
	return m
}()

// isTableConstraintStart reports whether the current token starts a
// TableConstraint clause inside CREATE TABLE / ALTER TABLE.
//
// Used at the disambiguation point where parseTableElement /
// parseTypedTableElement / parseAlterTableAdd decide between a column
// definition and a table-level constraint.
//
// No PG oracle test: TableConstraint is only reachable through
// CREATE/ALTER TABLE, where the surrounding column-definition rules
// admit overlapping leads (e.g. an IDENT could start a column name).
// Negative coverage is provided by the existing parseTableElement /
// parseTypedTableElement / parseAlterTableAdd tests in the corpus,
// plus the in-Go unit test TestIsTableConstraintStartCoverage.
func (p *Parser) isTableConstraintStart() bool {
	return tableConstraintLeadSet[p.cur.Type]
}

// isFuncTypeStart reports whether the current token starts a func_type
// production.
//
//	func_type:
//	    Typename
//	    | type_function_name attrs '%' TYPE_P
//	    | SETOF type_function_name attrs '%' TYPE_P
//
// At the FIRST-set level, FIRST(func_type) = FIRST(Typename), because:
//   - The %TYPE alternatives start with type_function_name, which is
//     already in FIRST(Typename) via SimpleTypename's GenericType path
//     (parseSimpleTypename → parseGenericType → parseTypeFunctionName).
//   - The SETOF prefix is already in FIRST(Typename) via the
//     `SETOF SimpleTypename` alternative.
//
// So isFuncTypeStart is a thin alias for isTypenameStart. The separate
// name documents the call site's intent: "this point is checking whether
// the next token is a func_type lead". Renaming the concept at the call
// site is cheaper than threading the comment everywhere.
//
// Validated against PG 17 by TestFuncTypeLeadTokensMatchPG, which probes
// the parameter type position of CREATE FUNCTION (a func_type position)
// and confirms the predicate matches PG's accept set keyword-for-keyword.
func (p *Parser) isFuncTypeStart() bool {
	return p.isTypenameStart()
}
