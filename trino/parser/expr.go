package parser

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the `expressions` DAG node (with function.go and predicate.go):
// it implements Trino's expression grammar — the `expression`,
// `booleanExpression`, `valueExpression`, and `primaryExpression` rules plus
// their helpers — as a hand-written recursive-descent parser over the token
// stream, producing an Expr value tree.
//
// Like datatypes.go's DataType, the expression nodes are PARSER-PACKAGE types
// (Expr, not ast.Node): the Trino ast tag set is closed to the ast-core node
// (only File/Identifier/QualifiedName), so — matching the types-node precedent —
// expression nodes live in package parser and carry their own Loc. Later DAG
// nodes (expr-json for the SQL/JSON functions, parser-select for the statement
// layer, deparse, analysis) embed these Expr values.
//
// The legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	expression        : booleanExpression ;
//	booleanExpression : valueExpression predicate_?              # predicated
//	                  | NOT booleanExpression                    # logicalNot
//	                  | booleanExpression AND booleanExpression  # and
//	                  | booleanExpression OR  booleanExpression  # or ;
//	valueExpression   : primaryExpression                                    # valueExpressionDefault
//	                  | valueExpression AT timeZoneSpecifier                  # atTimeZone
//	                  | (MINUS|PLUS) valueExpression                          # arithmeticUnary
//	                  | valueExpression (ASTERISK|SLASH|PERCENT) valueExpression # arithmeticBinary
//	                  | valueExpression (PLUS|MINUS) valueExpression          # arithmeticBinary
//	                  | valueExpression CONCAT valueExpression                # concatenation ;
//	primaryExpression : <the big literal/function/constructor/reference set> ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Oracle-confirmed precedence facts baked in:
//
//	P1 (single, non-associative predicate). A valueExpression carries AT MOST ONE
//	   predicate suffix: `a = b = c`, `a < b < c`, `a IN (1) IN (2)`,
//	   `a IS NULL IS NULL`, `a BETWEEN 1 AND 2 BETWEEN 3 AND 4` are all
//	   SYNTAX_ERRORs in Trino 481. parsePredicated therefore parses one optional
//	   predicate, never a chain.
//	P2 (|| binds tighter than the predicate). `a || b = c` parses as `(a||b) = c`
//	   — concatenation is inside valueExpression, the predicate compares two
//	   valueExpressions. NOT wraps the whole predicated form (`NOT a = b` is
//	   `NOT (a = b)`).
//
// Two boundary placements (node-scope decisions, recorded in the migration
// divergence ledger as flagged-by-design, both following the doris/snowflake
// precedent):
//
//	B1 (subqueries are raw-text placeholders). A subquery appearing inside an
//	   expression — `(SELECT …)`, `EXISTS (query)`, `<cmp> ANY (query)`,
//	   `IN (query)` — is captured as a SubqueryExpr placeholder holding the raw
//	   text, NOT a parsed query tree. This mirrors doris/snowflake, where
//	   expression-embedded subqueries stay placeholders even after the SELECT
//	   parser exists; the analysis node recurses into RawText when it needs
//	   lineage. The `query` rule belongs to parser-select.
//	B2 (pattern-recognition window frame deferred). A function's OVER frame may,
//	   per the legacy windowFrame rule, carry the full row-pattern subsystem
//	   (MEASURES / PATTERN / DEFINE / SUBSET / AFTER MATCH / INITIAL|SEEK). That
//	   subsystem is the parser-match-recognize node (row_pattern.go); this node
//	   implements the common frame `(ROWS|RANGE|GROUPS) [BETWEEN bound AND] bound`
//	   and leaves the pattern-frame to that node. See function.go.

// ---------------------------------------------------------------------------
// Expr node hierarchy (parser-package; not ast.Node — see file header)
// ---------------------------------------------------------------------------

// Expr is the interface implemented by every Trino expression node produced by
// this parser. Span returns the source byte range; concrete fields are reached
// by a Go type switch (mirroring the way DataType is consumed). Expr values are
// embedded by later DAG nodes (parser-select, deparse, analysis).
type Expr interface {
	// Span returns the source byte range covered by the expression.
	Span() ast.Loc
	// exprNode is a marker preventing unrelated types from satisfying Expr.
	exprNode()
}

// LiteralKind classifies a Literal.
type LiteralKind int

const (
	// LiteralNull is the NULL keyword literal.
	LiteralNull LiteralKind = iota
	// LiteralBool is TRUE or FALSE.
	LiteralBool
	// LiteralString is a character-string literal ('…', U&'…').
	LiteralString
	// LiteralInteger is an INTEGER_VALUE_ number.
	LiteralInteger
	// LiteralDecimal is a DECIMAL_VALUE_ number.
	LiteralDecimal
	// LiteralDouble is a DOUBLE_VALUE_ number.
	LiteralDouble
	// LiteralBinary is an X'…' binary literal.
	LiteralBinary
)

// Literal is a primaryExpression literal: NULL, a boolean, a string, a number,
// or a binary literal. Value holds the source-faithful text (for strings, the
// DECODED content with quotes stripped; for numbers and binary, the source
// spelling). Unicode is the U&'…' flag for a string literal.
type Literal struct {
	Kind    LiteralKind
	Value   string
	Unicode bool // true for a U&'…' string literal
	Loc     ast.Loc
}

func (n *Literal) Span() ast.Loc { return n.Loc }
func (*Literal) exprNode()       {}

// Parameter is the `?` positional parameter placeholder (the `parameter`
// primaryExpression alternative).
type Parameter struct {
	Loc ast.Loc
}

func (n *Parameter) Span() ast.Loc { return n.Loc }
func (*Parameter) exprNode()       {}

// ColumnRef is a bare `identifier` used as a value (the columnReference
// alternative). A dotted reference (a.b.c) is a chain of Dereference around a
// leading ColumnRef.
type ColumnRef struct {
	Name *ast.Identifier
	Loc  ast.Loc
}

func (n *ColumnRef) Span() ast.Loc { return n.Loc }
func (*ColumnRef) exprNode()       {}

// Dereference is `base . fieldName` (the dereference alternative): field access
// on a row value, or one step of a dotted name. FieldName is the trailing
// identifier.
type Dereference struct {
	Base      Expr
	FieldName *ast.Identifier
	Loc       ast.Loc
}

func (n *Dereference) Span() ast.Loc { return n.Loc }
func (*Dereference) exprNode()       {}

// Subscript is `value [ index ]` (the subscript alternative): array/map element
// access.
type Subscript struct {
	Value Expr
	Index Expr
	Loc   ast.Loc
}

func (n *Subscript) Span() ast.Loc { return n.Loc }
func (*Subscript) exprNode()       {}

// UnaryExpr is a prefix `+`/`-` arithmetic operation (arithmeticUnary).
type UnaryExpr struct {
	Op      string // "+" or "-"
	Operand Expr
	Loc     ast.Loc
}

func (n *UnaryExpr) Span() ast.Loc { return n.Loc }
func (*UnaryExpr) exprNode()       {}

// BinaryExpr is a binary arithmetic or concatenation operation: `*` `/` `%`
// `+` `-` `||`.
type BinaryExpr struct {
	Op    string
	Left  Expr
	Right Expr
	Loc   ast.Loc
}

func (n *BinaryExpr) Span() ast.Loc { return n.Loc }
func (*BinaryExpr) exprNode()       {}

// LogicalExpr is a binary AND/OR (the and/or booleanExpression alternatives).
type LogicalExpr struct {
	Op    string // "AND" or "OR"
	Left  Expr
	Right Expr
	Loc   ast.Loc
}

func (n *LogicalExpr) Span() ast.Loc { return n.Loc }
func (*LogicalExpr) exprNode()       {}

// NotExpr is `NOT booleanExpression` (logicalNot).
type NotExpr struct {
	Operand Expr
	Loc     ast.Loc
}

func (n *NotExpr) Span() ast.Loc { return n.Loc }
func (*NotExpr) exprNode()       {}

// AtTimeZoneExpr is `value AT TIME ZONE specifier` (atTimeZone). Zone is either
// a string/interval/general expression value.
type AtTimeZoneExpr struct {
	Value Expr
	Zone  Expr
	Loc   ast.Loc
}

func (n *AtTimeZoneExpr) Span() ast.Loc { return n.Loc }
func (*AtTimeZoneExpr) exprNode()       {}

// ParenExpr is a parenthesized single expression `( expression )`
// (parenthesizedExpression).
type ParenExpr struct {
	Expr Expr
	Loc  ast.Loc
}

func (n *ParenExpr) Span() ast.Loc { return n.Loc }
func (*ParenExpr) exprNode()       {}

// RowConstructor is `( e , e (, e)* )` (the parenthesized rowConstructor) or
// `ROW ( e (, e)* )` (the keyword form). Explicit records the ROW keyword form.
type RowConstructor struct {
	Elements []Expr
	Explicit bool // true for the ROW(...) spelling
	Loc      ast.Loc
}

func (n *RowConstructor) Span() ast.Loc { return n.Loc }
func (*RowConstructor) exprNode()       {}

// ArrayConstructor is `ARRAY [ e , … ]` (arrayConstructor).
type ArrayConstructor struct {
	Elements []Expr
	Loc      ast.Loc
}

func (n *ArrayConstructor) Span() ast.Loc { return n.Loc }
func (*ArrayConstructor) exprNode()       {}

// IntervalLiteral is `INTERVAL [+|-] 'string' field [TO field]`
// (intervalLiteral / the `interval` rule). Unlike the INTERVAL *type*
// (datatypes.go), this is a value with a string body and a sign.
type IntervalLiteral struct {
	Sign  string // "", "+", or "-"
	Value string // the decoded string body
	From  IntervalField
	To    *IntervalField // nil for the single-field form
	Loc   ast.Loc
}

func (n *IntervalLiteral) Span() ast.Loc { return n.Loc }
func (*IntervalLiteral) exprNode()       {}

// TypeConstructor is `name 'string'` (typeConstructor): a typed literal such as
// `DATE '2020-01-01'`, `JSON '[1]'`, `BIGINT '7'`, or `DOUBLE PRECISION '1.5'`.
// Name is the source type name (which may be the two-word "DOUBLE PRECISION").
type TypeConstructor struct {
	Name  string
	Value string // decoded string body
	Loc   ast.Loc
}

func (n *TypeConstructor) Span() ast.Loc { return n.Loc }
func (*TypeConstructor) exprNode()       {}

// SubqueryExpr is a placeholder for a subquery appearing inside an expression
// (B1 in the file header). RawText is the subquery source between the
// parentheses, trimmed; Kind records the surrounding form so the analysis node
// can interpret it. The query is NOT parsed here.
type SubqueryExpr struct {
	Kind    SubqueryKind
	RawText string
	Loc     ast.Loc
}

func (n *SubqueryExpr) Span() ast.Loc { return n.Loc }
func (*SubqueryExpr) exprNode()       {}

// SubqueryKind records the syntactic position a subquery placeholder occupies.
type SubqueryKind int

const (
	// SubqueryScalar is a bare `( query )` used as a scalar value.
	SubqueryScalar SubqueryKind = iota
	// SubqueryExists is `EXISTS ( query )`.
	SubqueryExists
	// SubqueryIn is the `( query )` of an `IN ( query )` predicate.
	SubqueryIn
	// SubqueryQuantified is the `( query )` of a quantified comparison
	// `<op> (ALL|ANY|SOME) ( query )`.
	SubqueryQuantified
)

// ---------------------------------------------------------------------------
// Entry points
// ---------------------------------------------------------------------------

// ParseExpression parses a complete Trino expression from a standalone string,
// returning the Expr and any ParseErrors. Trailing tokens after the expression
// are reported as an error. It is the string-input counterpart of parseExpr,
// for tests and callers that hold an expression string rather than a token
// stream. Mirrors ParseDataType / ParseQualifiedName.
func ParseExpression(input string) (Expr, []ParseError) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()

	expr, err := p.parseExpr()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}
	if p.cur.Kind != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		return expr, []ParseError{{Loc: p.cur.Loc, Msg: "unexpected token after expression: " + text}}
	}
	return expr, nil
}

// parseExpr parses the `expression` rule (== booleanExpression). It is the
// single entry point used by every expression position downstream: SELECT-list
// items, WHERE/HAVING/ON conditions, function arguments, CASE branches, GROUP
// BY / ORDER BY keys, VALUES rows, DML SET targets.
func (p *Parser) parseExpr() (Expr, error) {
	return p.parseBooleanExpr()
}

// ---------------------------------------------------------------------------
// booleanExpression — OR / AND / NOT / predicated
// ---------------------------------------------------------------------------

// parseBooleanExpr parses `booleanExpression`. Precedence (lowest first): OR,
// then AND, then a prefix NOT, then the predicated value expression. OR and AND
// are left-associative; they are flattened iteratively rather than recursively
// to keep deep chains shallow on the Go stack.
func (p *Parser) parseBooleanExpr() (Expr, error) {
	return p.parseOr()
}

// parseOr parses `left OR right (OR right)*`, left-associative.
func (p *Parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwOR {
		p.advance() // consume OR
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{
			Op:    "OR",
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseAnd parses `left AND right (AND right)*`, left-associative.
func (p *Parser) parseAnd() (Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwAND {
		p.advance() // consume AND
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{
			Op:    "AND",
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseNot parses `NOT booleanExpression` (prefix, right-recursive so `NOT NOT a`
// nests) or falls through to the predicated value expression. NOT binds looser
// than every predicate/arithmetic operator: `NOT a = b` is `NOT (a = b)`.
func (p *Parser) parseNot() (Expr, error) {
	if p.cur.Kind == kwNOT {
		notTok := p.advance() // consume NOT
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &NotExpr{
			Operand: operand,
			Loc:     ast.Loc{Start: notTok.Loc.Start, End: operand.Span().End},
		}, nil
	}
	return p.parsePredicated()
}

// parsePredicated parses `valueExpression predicate_?` (the predicated
// alternative). At most ONE predicate suffix is consumed — Trino's predicate is
// non-associative (P1 in the file header), so a second predicate is left for the
// caller, which surfaces as a syntax error at the trailing operator. The
// predicate-suffix parsing itself lives in predicate.go.
func (p *Parser) parsePredicated() (Expr, error) {
	left, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	return p.parsePredicateSuffix(left)
}

// ---------------------------------------------------------------------------
// valueExpression — AT TIME ZONE / arithmetic / concat / unary
// ---------------------------------------------------------------------------

// parseValueExpr parses `valueExpression`. The layered precedence, tightest
// last, is: concatenation (||) over additive (+ -) over multiplicative (* / %)
// over a prefix unary (+ -) over the AT-TIME-ZONE postfix over primaryExpression.
//
// All binary operators here are left-associative and flattened iteratively. The
// AT-TIME-ZONE postfix and the unary prefix are folded into the operand layers
// (parseUnary / parseAtTimeZonePostfix) so they compose with arithmetic as Trino
// accepts (`-a * b`, `a * b AT TIME ZONE 'UTC'`, `a AT TIME ZONE 'UTC' * b`).
func (p *Parser) parseValueExpr() (Expr, error) {
	return p.parseConcat()
}

// parseConcat parses `left || right (|| right)*`, left-associative. `||` is the
// loosest arithmetic-level operator, so `a || b = c` groups as `(a || b) = c`
// when the predicate suffix is later applied (P2).
func (p *Parser) parseConcat() (Expr, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == tokConcat {
		p.advance() // consume ||
		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{
			Op:    "||",
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseAdditive parses `left (+|-) right (…)*`, left-associative.
func (p *Parser) parseAdditive() (Expr, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == int('+') || p.cur.Kind == int('-') {
		opTok := p.advance()
		op := string(rune(opTok.Kind))
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseMultiplicative parses `left (*|/|%) right (…)*`, left-associative.
func (p *Parser) parseMultiplicative() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == int('*') || p.cur.Kind == int('/') || p.cur.Kind == int('%') {
		opTok := p.advance()
		op := string(rune(opTok.Kind))
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parseUnary parses a prefix `+`/`-` (arithmeticUnary, right-recursive so
// `- - 1` nests) or falls through to the AT-TIME-ZONE postfix layer.
func (p *Parser) parseUnary() (Expr, error) {
	if p.cur.Kind == int('+') || p.cur.Kind == int('-') {
		opTok := p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Op:      string(rune(opTok.Kind)),
			Operand: operand,
			Loc:     ast.Loc{Start: opTok.Loc.Start, End: operand.Span().End},
		}, nil
	}
	return p.parseAtTimeZonePostfix()
}

// parseAtTimeZonePostfix parses a primaryExpression followed by zero or more
// `AT TIME ZONE specifier` postfixes (atTimeZone). The specifier is
// `TIME ZONE interval` or `TIME ZONE string` — both reduce to an expression
// value here (an INTERVAL literal or a string), parsed by parsePrimary so the
// interval form is captured as an IntervalLiteral.
func (p *Parser) parseAtTimeZonePostfix() (Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwAT {
		p.advance() // consume AT
		if _, err := p.expect(kwTIME); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwZONE); err != nil {
			return nil, err
		}
		zone, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		expr = &AtTimeZoneExpr{
			Value: expr,
			Zone:  zone,
			Loc:   ast.Loc{Start: expr.Span().Start, End: zone.Span().End},
		}
	}
	return expr, nil
}

// ---------------------------------------------------------------------------
// primaryExpression — atoms + postfix subscript/dereference
// ---------------------------------------------------------------------------

// parsePrimary parses one `primaryExpression`, including the left-recursive
// postfix subscript (`[index]`) and dereference (`.field`) which are folded into
// a loop after the atom. It dispatches keyword-led atoms (CASE, CAST, ARRAY,
// ROW, special functions, …) and otherwise reads a literal / parameter /
// identifier-or-function / parenthesized form via parsePrimaryAtom.
func (p *Parser) parsePrimary() (Expr, error) {
	atom, err := p.parsePrimaryAtom()
	if err != nil {
		return nil, err
	}
	return p.parsePostfix(atom)
}

// parsePostfix consumes zero or more `[ index ]` subscripts and `. field`
// dereferences on an already-parsed primary, left-associative. `m['k'][1].f`
// and `a.b[1]` both round-trip through this loop.
//
// A `.` is treated as a dereference only when followed by an identifier; `.*`
// (the selectAll form) is NOT a primaryExpression and belongs to the SELECT-list
// layer, so a `.` before `*` is left for the caller (it surfaces there).
func (p *Parser) parsePostfix(expr Expr) (Expr, error) {
	for {
		switch p.cur.Kind {
		case int('['):
			p.advance() // consume '['
			index, err := p.parseValueExpr()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(']'))
			if err != nil {
				return nil, err
			}
			expr = &Subscript{
				Value: expr,
				Index: index,
				Loc:   ast.Loc{Start: expr.Span().Start, End: closeTok.Loc.End},
			}
		case int('.'):
			// Only a `.identifier` dereference belongs here; `.*` is selectAll.
			if !isIdentifierStart(p.peekNext().Kind) {
				return expr, nil
			}
			p.advance() // consume '.'
			field, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			expr = &Dereference{
				Base:      expr,
				FieldName: field,
				Loc:       ast.Loc{Start: expr.Span().Start, End: field.Loc.End},
			}
		default:
			return expr, nil
		}
	}
}

// parsePrimaryAtom parses a primaryExpression atom WITHOUT the postfix
// subscript/dereference (handled by parsePostfix). It dispatches on the leading
// token. Function-call, CAST/TRY_CAST, CASE, the special built-ins
// (EXTRACT/SUBSTRING/TRIM/NORMALIZE/POSITION/GROUPING/LISTAGG/special-datetime),
// and the SQL/JSON functions are routed to function.go.
func (p *Parser) parsePrimaryAtom() (Expr, error) {
	switch p.cur.Kind {
	case kwNULL:
		tok := p.advance()
		return &Literal{Kind: LiteralNull, Value: "NULL", Loc: tok.Loc}, nil
	case kwTRUE:
		tok := p.advance()
		return &Literal{Kind: LiteralBool, Value: "TRUE", Loc: tok.Loc}, nil
	case kwFALSE:
		tok := p.advance()
		return &Literal{Kind: LiteralBool, Value: "FALSE", Loc: tok.Loc}, nil
	case tokString, tokUnicodeString:
		return p.parseStringLiteral()
	case tokInteger:
		tok := p.advance()
		return &Literal{Kind: LiteralInteger, Value: tok.Str, Loc: tok.Loc}, nil
	case tokDecimal:
		tok := p.advance()
		return &Literal{Kind: LiteralDecimal, Value: tok.Str, Loc: tok.Loc}, nil
	case tokDouble:
		tok := p.advance()
		return &Literal{Kind: LiteralDouble, Value: tok.Str, Loc: tok.Loc}, nil
	case tokBinaryLiteral:
		tok := p.advance()
		return &Literal{Kind: LiteralBinary, Value: tok.Str, Loc: tok.Loc}, nil
	case tokQuestion, int('?'):
		tok := p.advance()
		return &Parameter{Loc: tok.Loc}, nil
	case kwINTERVAL:
		return p.parseIntervalLiteral()
	case int('('):
		return p.parseParenOrSubqueryOrRow()
	case kwARRAY:
		// ARRAY[ … ] constructor — distinct from the ARRAY type. If ARRAY is not
		// followed by '[', it is an ordinary identifier (e.g. a column/function
		// named ARRAY), handled by the identifier path.
		if p.peekNext().Kind == int('[') {
			return p.parseArrayConstructor()
		}
		return p.parseIdentifierOrFunction()
	case kwROW:
		// ROW( … ) constructor — distinct from the ROW type. ROW not followed by
		// '(' is an ordinary identifier.
		if p.peekNext().Kind == int('(') {
			return p.parseRowConstructor()
		}
		return p.parseIdentifierOrFunction()
	case kwCASE:
		return p.parseCaseExpr()
	case kwCAST, kwTRY_CAST:
		return p.parseCast()
	case kwEXISTS:
		return p.parseExists()
	case kwEXTRACT:
		return p.parseExtract()
	case kwTRIM:
		return p.parseTrim()
	case kwNORMALIZE:
		return p.parseNormalize()
	case kwLISTAGG:
		return p.parseListagg()
	case kwGROUPING:
		return p.parseGrouping()
	case kwJSON_EXISTS, kwJSON_VALUE, kwJSON_QUERY, kwJSON_OBJECT, kwJSON_ARRAY:
		return p.parseJSONFunction()
	case kwRUNNING, kwFINAL:
		// processingMode prefix on a function call (RUNNING/FINAL last(x)),
		// used in MATCH_RECOGNIZE measure expressions.
		return p.parseProcessingModeFunction()
	case kwCURRENT_DATE, kwCURRENT_USER, kwCURRENT_CATALOG, kwCURRENT_SCHEMA, kwCURRENT_PATH,
		kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwLOCALTIME, kwLOCALTIMESTAMP:
		return p.parseSpecialDateTimeFunction()
	case kwDOUBLE:
		// DOUBLE PRECISION 'literal' typeConstructor. Otherwise DOUBLE is an
		// ordinary identifier handled by the identifier path.
		if p.peekNext().Kind == kwPRECISION {
			return p.parseDoublePrecisionConstructor()
		}
		return p.parseIdentifierOrFunction()
	default:
		// An identifier (or a non-reserved keyword usable as one) begins a
		// column reference, a function call, a lambda, or a typeConstructor.
		if isIdentifierStart(p.cur.Kind) {
			return p.parseIdentifierOrFunction()
		}
		return nil, p.exprError()
	}
}

// parseStringLiteral parses a basic or unicode string literal (string_ rule).
// The lexer has already decoded the body and recorded whether it was a U&'…'
// unicode string. A `UESCAPE '…'` clause, when present, has been consumed by the
// lexer as part of the unicode-string token.
func (p *Parser) parseStringLiteral() (Expr, error) {
	tok := p.advance()
	return &Literal{
		Kind:    LiteralString,
		Value:   tok.Str,
		Unicode: tok.Kind == tokUnicodeString,
		Loc:     tok.Loc,
	}, nil
}

// parseIntervalLiteral parses `INTERVAL [+|-] 'string' field [TO field]`
// (the `interval` rule / intervalLiteral primaryExpression). It is distinct from
// the INTERVAL *type* (datatypes.go): a literal requires a string body and a
// field, and allows a leading sign. The same D1 same-family/ordered qualifier
// rule as the type applies to the optional `TO to` range.
func (p *Parser) parseIntervalLiteral() (Expr, error) {
	intervalTok := p.advance() // consume INTERVAL

	sign := ""
	if p.cur.Kind == int('+') || p.cur.Kind == int('-') {
		sign = string(rune(p.advance().Kind))
	}

	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return nil, p.exprErrorAt("expected string literal in INTERVAL")
	}
	strTok := p.advance()

	from, ok := intervalFieldFromKind(p.cur.Kind)
	if !ok {
		return nil, p.exprErrorAt("expected interval field (YEAR/MONTH/DAY/HOUR/MINUTE/SECOND)")
	}
	fromTok := p.advance()

	lit := &IntervalLiteral{
		Sign:  sign,
		Value: strTok.Str,
		From:  from,
		Loc:   ast.Loc{Start: intervalTok.Loc.Start, End: fromTok.Loc.End},
	}

	if p.cur.Kind == kwTO {
		toTok := p.peekNext()
		to, ok := intervalFieldFromKind(toTok.Kind)
		if !ok {
			p.advance() // consume TO so the error points at the bad to-field
			return nil, p.exprErrorAt("expected interval field after TO")
		}
		if !ValidIntervalRange(from, to) {
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "invalid interval qualifier: " + from.String() + " TO " + to.String(),
			}
		}
		p.advance() // consume TO
		toEnd := p.advance()
		toVal := to
		lit.To = &toVal
		lit.Loc.End = toEnd.Loc.End
	}
	return lit, nil
}

// parseArrayConstructor parses `ARRAY [ (expression (, expression)*)? ]`
// (arrayConstructor). The element list may be empty (`ARRAY[]`).
func (p *Parser) parseArrayConstructor() (Expr, error) {
	arrTok := p.advance() // consume ARRAY
	if _, err := p.expect(int('[')); err != nil {
		return nil, err
	}
	elems, closeTok, err := p.parseBracketedExprList(int(']'))
	if err != nil {
		return nil, err
	}
	return &ArrayConstructor{
		Elements: elems,
		Loc:      ast.Loc{Start: arrTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseRowConstructor parses the keyword form `ROW ( expression (, expression)* )`
// (rowConstructor). The element list must be non-empty (Trino rejects ROW()).
func (p *Parser) parseRowConstructor() (Expr, error) {
	rowTok := p.advance() // consume ROW
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	elems := []Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		elems = append(elems, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &RowConstructor{
		Elements: elems,
		Explicit: true,
		Loc:      ast.Loc{Start: rowTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseParenOrSubqueryOrRow disambiguates a leading '(' among:
//   - a parenthesized-parameter lambda `( (ident (, ident)*)? ) -> body`
//     (the second lambda alternative, including the empty-parameter `() -> …`);
//   - a subquery `( query )` → SubqueryExpr placeholder (B1), when SELECT/WITH/
//     TABLE/VALUES follows the '(';
//   - a parenthesized rowConstructor `( e , e (, e)* )` when a top-level comma
//     follows the first expression;
//   - a plain parenthesized expression `( expression )` otherwise.
//
// The lambda case is detected by a speculative parse of an identifier list
// followed by `) ->`; on failure the parser rewinds and falls through to the
// subquery/row/paren readings. This is required because `(x)` alone is a
// ParenExpr and `(x, y)` alone is a rowConstructor — only the trailing `->`
// distinguishes a lambda, and the deciding token lies past the ')'.
func (p *Parser) parseParenOrSubqueryOrRow() (Expr, error) {
	openTok := p.advance() // consume '('

	if lam, ok, err := p.tryParenLambda(openTok.Loc.Start); err != nil {
		return nil, err
	} else if ok {
		return lam, nil
	}

	if p.startsQuery() {
		return p.parseSubqueryPlaceholder(openTok.Loc.Start, SubqueryScalar)
	}

	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.cur.Kind == int(',') {
		// rowConstructor: ( e , e (, e)* )
		elems := []Expr{first}
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			elems = append(elems, next)
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		return &RowConstructor{
			Elements: elems,
			Loc:      ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ParenExpr{
		Expr: first,
		Loc:  ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// tryParenLambda speculatively parses a parenthesized-parameter lambda body
// `(ident (, ident)*)? ) -> body`, the opening '(' already consumed (its offset
// is startOffset). It returns ok=true with the LambdaExpr when the input is a
// lambda, ok=false (parser rewound to just after the '(') when it is not, and a
// non-nil error only when the input WAS a lambda but its body failed to parse.
//
// A '(' that opens a lambda parameter list contains only identifiers and commas
// before the ')'; the ')' must be immediately followed by '->'. Any other shape
// (a subquery, a row/paren expression) fails the speculation and rewinds.
func (p *Parser) tryParenLambda(startOffset int) (Expr, bool, error) {
	cp := p.checkpoint()

	var params []*ast.Identifier
	if p.cur.Kind != int(')') {
		if !isIdentifierStart(p.cur.Kind) {
			p.restore(cp)
			return nil, false, nil
		}
		params = append(params, identFromToken(p.advance()))
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			if !isIdentifierStart(p.cur.Kind) {
				p.restore(cp)
				return nil, false, nil
			}
			params = append(params, identFromToken(p.advance()))
		}
	}
	if p.cur.Kind != int(')') || p.peekNext().Kind != tokArrow {
		p.restore(cp)
		return nil, false, nil
	}
	p.advance() // consume ')'
	p.advance() // consume '->'

	body, err := p.parseExpr()
	if err != nil {
		return nil, true, err
	}
	return &LambdaExpr{
		Params: params,
		Body:   body,
		Loc:    ast.Loc{Start: startOffset, End: body.Span().End},
	}, true, nil
}

// startsQuery reports whether the current token begins a `query` (used to
// recognise a subquery after a '('). Trino's queryPrimary starts with SELECT,
// WITH (cte or inline FUNCTION), TABLE, or VALUES.
func (p *Parser) startsQuery() bool {
	switch p.cur.Kind {
	case kwSELECT, kwWITH, kwTABLE, kwVALUES:
		return true
	default:
		return false
	}
}

// parseSubqueryPlaceholder consumes tokens until the matching ')' that closes a
// subquery and returns a SubqueryExpr placeholder holding the raw inner text
// (B1 in the file header). The opening '(' has already been consumed; startOffset
// is its byte offset. Nested parentheses are balanced. This mirrors
// doris/parser.parseSubqueryPlaceholder — expression-embedded subqueries are not
// parsed into a query tree here.
func (p *Parser) parseSubqueryPlaceholder(startOffset int, kind SubqueryKind) (*SubqueryExpr, error) {
	depth := 1
	subStart := p.cur.Loc.Start
	subEnd := p.cur.Loc.Start
	for depth > 0 && p.cur.Kind != tokEOF {
		switch p.cur.Kind {
		case int('('):
			depth++
		case int(')'):
			depth--
			if depth == 0 {
				subEnd = p.cur.Loc.Start
			}
		}
		if depth > 0 {
			p.advance()
		}
	}
	if depth != 0 {
		return nil, &ParseError{Loc: p.cur.Loc, Msg: "unterminated subquery"}
	}
	raw := p.sourceSlice(subStart, subEnd)
	closeTok := p.advance() // consume ')'
	return &SubqueryExpr{
		Kind:    kind,
		RawText: strings.TrimSpace(raw),
		Loc:     ast.Loc{Start: startOffset, End: closeTok.Loc.End},
	}, nil
}

// sourceSlice returns the substring of the original input spanning the absolute
// byte range [absStart, absEnd). It accounts for the parser's baseOffset so a
// statement carved out of a larger document slices correctly. Out-of-range
// inputs yield "" rather than panicking.
func (p *Parser) sourceSlice(absStart, absEnd int) string {
	s := absStart - p.baseOffset
	e := absEnd - p.baseOffset
	if s < 0 || e > len(p.input) || s > e {
		return ""
	}
	return p.input[s:e]
}

// parseBracketedExprList parses `(expression (, expression)*)? closer`, returning
// the elements and the closing token. The opening bracket has already been
// consumed. The list may be empty (e.g. ARRAY[] or GROUPING()).
func (p *Parser) parseBracketedExprList(closer TokenKind) ([]Expr, Token, error) {
	if p.cur.Kind == closer {
		closeTok := p.advance()
		return nil, closeTok, nil
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, Token{}, err
	}
	elems := []Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, Token{}, err
		}
		elems = append(elems, next)
	}
	closeTok, err := p.expect(closer)
	if err != nil {
		return nil, Token{}, err
	}
	return elems, closeTok, nil
}

// ---------------------------------------------------------------------------
// errors
// ---------------------------------------------------------------------------

// exprError returns a *ParseError describing a missing expression at the current
// token, distinct from the generic identifier/type errors so a value position
// reports "expected expression".
func (p *Parser) exprError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{Loc: p.cur.Loc, Msg: "expected expression, found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "expected expression, found " + text}
}

// exprErrorAt returns a *ParseError with a custom message at the current token.
func (p *Parser) exprErrorAt(msg string) *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: msg}
}
