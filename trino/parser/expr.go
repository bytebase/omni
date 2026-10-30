package parser

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

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
//	   lineage. The `query` rule belongs to parser-select. The strict entry
//	   point does parse each body in place to validate it, and keeps that tree
//	   in SubqueryExpr.Query.
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
// can interpret it.
type SubqueryExpr struct {
	Kind    SubqueryKind
	RawText string
	// TextStart is the absolute byte offset of RawText's first character in
	// the parser's input, so errors and locations from re-parsing RawText can
	// be shifted back into the outer statement's coordinates.
	TextStart int
	// Query is the body as the strict entry point parsed it in place, with
	// locations in the parser's input like the rest of the tree. It is nil
	// from ParseBestEffort, which leaves bodies raw, and for a body the strict
	// parse scanned raw instead. Use it rather than parsing RawText again:
	// reparsing the body at every level of n nested subqueries costs O(n²).
	Query *Query
	Loc   ast.Loc
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
	if p.nestingErr != nil {
		err = p.nestingErr // see parseSingle
	} else if err == nil && p.treeTooHigh(expr, len(input)) {
		err = &ParseError{Position: expr.Span().Start, End: expr.Span().Start, Message: nestingTooDeep}
	}
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Message: err.Error()}}
	}
	if p.cur.Kind != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		return expr, []ParseError{{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "unexpected token after expression: " + text}}
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
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
//
// A speculative reading records the result, and the reading tried after it
// replays the record (see Parser.speculating).
func (p *Parser) parsePrimary() (Expr, error) {
	start := p.cur.Loc.Start
	if r, ok := p.primaryMemo[start]; ok {
		p.replayTo(r.end)
		if r.raws != nil {
			p.rawQueries = append(p.rawQueries, rawQuery{group: r.raws})
		}
		return r.expr, r.err
	}
	raws := len(p.rawQueries)
	expr, err := p.parsePrimaryUnrecorded()
	if p.speculating > 0 && p.cur.Loc.Start != start {
		p.recordPrimary(start, raws, expr, err)
	}
	return expr, err
}

func (p *Parser) parsePrimaryUnrecorded() (Expr, error) {
	atom, err := p.parsePrimaryAtom()
	if err != nil {
		return nil, err
	}
	return p.parsePostfix(atom)
}

// recordPrimary records the primary expression that started at start, whose
// parse found raws query bodies already captured. The bodies it captured are
// folded into one group entry, so an enclosing record copies one entry for
// this whole primary instead of every body inside it.
func (p *Parser) recordPrimary(start, raws int, expr Expr, err error) {
	r := exprReplay{expr: expr, err: err}
	if err == nil && len(p.rawQueries) > raws {
		r.raws = append([]rawQuery(nil), p.rawQueries[raws:]...)
		p.rawQueries = append(p.rawQueries[:raws], rawQuery{group: r.raws})
	}
	r.end = p.checkpoint()
	if p.primaryMemo == nil {
		p.primaryMemo = make(map[int]exprReplay)
	}
	p.primaryMemo[start] = r
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
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "invalid interval qualifier: " + from.String() + " TO " + to.String(),
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
//
// Strict mode validates every body (validateSubqueries). It first parses the
// body in place: a body that parses as a query ending at the matching ')' is
// valid and is not parsed again. Validating it afterwards instead reparsed
// the text of every nested body once per level, so n nested subqueries cost
// O(n²). A body that fails here is scanned raw and validated afterwards as
// before, so its errors are unchanged.
//
// A lex error does not fail the in-place parse either. Failing on it made
// every body around the error raw, and validating each parsed the text of
// the bodies inside it again, so a lex error under n nested subqueries cost
// O(n²). Validating the body from its text reports what the bodies nested in
// it report, then its own lex errors; the bodies nested in it are captured
// in its place in rawQueries, and a lexOnly entry after them reports its lex
// errors, so the statement's errors come out as before. That holds while the
// body is one statement to the text parse: a ';' token in it (inside a nested
// raw body, in a routine) makes that parse split it, so such a body is still
// scanned raw.
//
// Nor does white space the lexer rejects (U+00A0, \v) between the query and
// the ')': the body's text is trimmed of it, so the text parse reads the
// query alone and succeeds. Failing here on it made a body raw whose
// validation then validated the bodies inside it, each in turn, so n nested
// subqueries ending in it cost O(n²). Its lex errors are not the body's but
// those of the text around it, as for the text parse.
func (p *Parser) parseSubqueryPlaceholder(startOffset int, kind SubqueryKind) (*SubqueryExpr, error) {
	subStart := p.cur.Loc.Start
	if p.strictTrailing {
		cp := p.checkpoint()
		if q, err := p.parseQuery(); err == nil {
			end := p.cur.Loc.Start // where the body's text ends, trimmed
			if p.cur.Kind == int(')') || p.skipTrimmedSpace() {
				lexErrs := p.lexer.errorCount() != cp.lexErrs
				if !lexErrs || !p.hasSemicolonIn(subStart, end) {
					sub := p.finishSubqueryPlaceholder(startOffset, subStart, kind)
					sub.Query = q
					if lexErrs {
						p.rawQueries = append(p.rawQueries, rawQuery{lexOnly: true, start: subStart, end: end})
					}
					return sub, nil
				}
			}
		}
		p.restore(cp)
		if p.cur.Kind == kwWITH && p.peekNext().Kind == kwFUNCTION {
			if sub := p.parseRoutineSubquery(startOffset, subStart, kind); sub != nil {
				return sub, nil
			}
		}
	}
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		if p.cur.Kind == int('(') {
			p.skipParenGroup()
			if p.cur.Kind == tokEOF {
				break
			}
		}
		p.advance()
	}
	p.recordParenGroupEnd(startOffset)
	if p.cur.Kind == tokEOF {
		return nil, &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "unterminated subquery"}
	}
	sub := p.finishSubqueryPlaceholder(startOffset, subStart, kind)
	p.rawQueries = append(p.rawQueries, rawQuery{text: sub.RawText, start: sub.TextStart, depth: p.depth})
	return sub, nil
}

// parseRoutineSubquery parses in place an expression subquery body that is a
// WITH FUNCTION statement, as its validation parses its text: as a
// statement, one level deeper. Trino rejects the body (WITH FUNCTION
// prefixes a statement, not a query), and its validation reports the errors
// of the bodies captured raw in it, or else that it is no query. Captured
// raw, the body was validated by parsing its text, which validated each body
// inside it the same way, so n nested ones cost O(n²). The body's entry in
// rawQueries reports the same errors from the in-place parse.
//
// It returns nil, with the parser back where it was, when the in-place parse
// does not read the body as the text parse does: when it fails, ends before
// the ')', or meets a ';', at which the text parse may split the text. It
// does too when it nests too deeply: the text parse reports that among the
// statement's errors, where it is. The parser keeps no record of a parse
// that failed so (see enterNesting), so from then on it leaves every such
// body to the text parse.
func (p *Parser) parseRoutineSubquery(startOffset, subStart int, kind SubqueryKind) *SubqueryExpr {
	if p.routineNestingFailed || p.nestingErr != nil {
		return nil
	}
	cp := p.checkpoint()
	node, err := p.parseStmt()
	if p.nestingErr != nil {
		p.nestingErr = nil
		p.routineNestingFailed = true
		p.restore(cp)
		return nil
	}
	end := p.cur.Loc.Start // where the body's text ends, trimmed
	if err != nil || p.cur.Kind != int(')') && !p.skipTrimmedSpace() || p.hasSemicolonIn(subStart, end) {
		p.restore(cp)
		return nil
	}
	nested := append([]rawQuery(nil), p.rawQueries[cp.rawqs:]...)
	p.rawQueries = p.rawQueries[:cp.rawqs]
	sub := p.finishSubqueryPlaceholder(startOffset, subStart, kind)
	p.rawQueries = append(p.rawQueries, rawQuery{
		text:    sub.RawText,
		start:   sub.TextStart,
		end:     end,
		group:   nested,
		routine: &routineBody{tooHigh: p.treeTooHigh(node, len(sub.RawText)), first: cp.cur.Loc},
	})
	return sub
}

// skipTrimmedSpace moves to the ')' after the current token when only white
// space, as strings.TrimSpace trims it, lies before it from the current
// token on, and reports whether it did. The current token is then an INVALID
// one, white space the lexer rejects.
func (p *Parser) skipTrimmedSpace() bool {
	if p.cur.Kind != tokInvalid {
		return false
	}
	from := p.cur.Loc.Start - p.baseOffset
	i := from
	for i < len(p.input) {
		r, size := utf8.DecodeRuneInString(p.input[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += size
	}
	if i == from || i >= len(p.input) || p.input[i] != ')' {
		return false
	}
	for p.cur.Kind != tokEOF && p.cur.Loc.Start-p.baseOffset < i {
		p.advance()
	}
	return p.cur.Kind == int(')')
}

// hasSemicolonIn reports whether a ';' token lies in [start, end), absolute
// offsets in the segment. The segment is lexed for them once, and only when
// it holds a ';' byte at all.
func (p *Parser) hasSemicolonIn(start, end int) bool {
	if !p.semicolonsLexed {
		p.semicolonsLexed = true
		if strings.IndexByte(p.input, ';') >= 0 {
			lx := NewLexerWithOffset(p.input, p.baseOffset)
			for tok := lx.NextToken(); tok.Kind != tokEOF; tok = lx.NextToken() {
				if tok.Kind == int(';') {
					p.semicolons = append(p.semicolons, tok.Loc.Start)
				}
			}
		}
	}
	i := sort.SearchInts(p.semicolons, start)
	return i < len(p.semicolons) && p.semicolons[i] < end
}

// skipParenGroup moves from the '(' at p.cur to its matching ')', left as
// p.cur, or to EOF when it never closes. A group scanned before is jumped
// over (see parenGroupEnd), and every group scanned is recorded there, the
// nested ones first.
//
// It keeps the groups it is inside on a stack of its own instead of calling
// itself for each nested group: the scan runs over bodies the parse could
// not read, nested arbitrarily deep, and a call per level overflowed the Go
// stack.
func (p *Parser) skipParenGroup() {
	var opens []int // offsets of the '(' groups being scanned, innermost last
	for {
		// p.cur is a '(': jump over its group when it was scanned before,
		// leaving p.cur at the group's end, else enter it.
		if r, ok := p.parenGroupEnd[p.cur.Loc.Start]; ok {
			p.replayParen(r)
			if len(opens) == 0 {
				return
			}
			if p.cur.Kind != tokEOF {
				p.advance() // consume the nested group's ')'
			}
		} else {
			opens = append(opens, p.cur.Loc.Start)
			p.advance() // consume '('
		}
		// Scan the innermost group up to a nested '(' to enter, or to its
		// end; an ended group returns the scan to the group around it.
		for {
			for p.cur.Kind != int(')') && p.cur.Kind != tokEOF && p.cur.Kind != int('(') {
				p.advance()
			}
			if p.cur.Kind == int('(') {
				break
			}
			p.recordParenGroupEnd(opens[len(opens)-1])
			opens = opens[:len(opens)-1]
			if len(opens) == 0 {
				return
			}
			if p.cur.Kind != tokEOF {
				p.advance() // consume the nested group's ')'
			}
		}
	}
}

// recordParenGroupEnd records the current token, the matching ')' or EOF, as
// the end of the '(' group at absolute offset open.
func (p *Parser) recordParenGroupEnd(open int) {
	if p.parenGroupEnd == nil {
		p.parenGroupEnd = make(map[int]parenReplay)
	}
	p.parenGroupEnd[open] = p.recordParenReplay()
}

// finishSubqueryPlaceholder consumes the ')' that closes a subquery body
// starting at subStart and returns its placeholder.
func (p *Parser) finishSubqueryPlaceholder(startOffset, subStart int, kind SubqueryKind) *SubqueryExpr {
	raw := p.sourceSlice(subStart, p.cur.Loc.Start)
	closeTok := p.advance() // consume ')'
	trimmed := strings.TrimSpace(raw)
	return &SubqueryExpr{
		Kind:      kind,
		RawText:   trimmed,
		TextStart: subStart + strings.Index(raw, trimmed),
		Loc:       ast.Loc{Start: startOffset, End: closeTok.Loc.End},
	}
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
		return &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected expression, found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected expression, found " + text}
}

// exprErrorAt returns a *ParseError with a custom message at the current token.
func (p *Parser) exprErrorAt(msg string) *ParseError {
	return &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: msg}
}
