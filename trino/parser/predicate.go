package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `expressions` DAG node (with expr.go and
// function.go): it implements the predicate suffix of Trino's `booleanExpression`
// — the `predicate_` rule — that attaches to a valueExpression to form a
// `predicated` booleanExpression.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	predicate_
//	    : comparisonOperator right=valueExpression                            # comparison
//	    | comparisonOperator comparisonQuantifier LPAREN_ query RPAREN_        # quantifiedComparison
//	    | NOT? BETWEEN_ lower=valueExpression AND_ upper=valueExpression       # between
//	    | NOT? IN_ LPAREN_ expression (COMMA_ expression)* RPAREN_             # inList
//	    | NOT? IN_ LPAREN_ query RPAREN_                                       # inSubquery
//	    | NOT? LIKE_ pattern=valueExpression (ESCAPE_ escape=valueExpression)? # like
//	    | IS_ NOT_? NULL_                                                      # nullPredicate
//	    | IS_ NOT_? DISTINCT_ FROM_ right=valueExpression                     # distinctFrom
//	    ;
//	comparisonOperator  : EQ_ | NEQ_ | LT_ | LTE_ | GT_ | GTE_ ;
//	comparisonQuantifier: ALL_ | SOME_ | ANY_ ;
//
// Structural nuance (quantifiedComparison vs ANY/SOME/ALL as a function): in
// Trino 481 ANY/SOME/ALL are also non-reserved function NAMES, so after a
// comparison operator `<op> ANY ( … )` is a quantifiedComparison only when the
// parenthesized content is a query; otherwise `ANY(1)`, `ANY(1, 2)`, `ANY()` are
// ordinary function calls forming the comparison's right-hand value. This node
// builds a QuantifiedComparisonExpr with a raw-text subquery placeholder for the
// `<op> (ANY|SOME|ALL) ( … )` shape regardless of the inner content; the
// accept/reject verdict matches Trino in every probed case (the function-call and
// subquery readings are both accepted), and distinguishing the two structurally
// needs the query/expression decision that belongs to parser-select. Recorded as
// a flagged structural boundary, not an accept/reject divergence.
//
// Oracle-confirmed behavior (Trino 481): the predicate is NON-ASSOCIATIVE —
// exactly one predicate may follow a valueExpression. `a = b = c`, `a < b < c`,
// `a IN (1) IN (2)`, `a IS NULL IS NULL`, `a BETWEEN 1 AND 2 BETWEEN 3 AND 4` are
// all SYNTAX_ERRORs. parsePredicateSuffix therefore consumes at most one
// predicate and never loops; a trailing predicate operator is left for the
// caller, where it surfaces as the syntax error Trino reports. Note Trino has NO
// `IS [NOT] TRUE/FALSE/UNKNOWN` predicate (those spellings are rejected) — only
// IS [NOT] NULL and IS [NOT] DISTINCT FROM.

// ComparisonExpr is `left <op> right` where <op> is one of = <> < <= > >=
// (the comparison predicate). Op holds the source operator spelling.
type ComparisonExpr struct {
	Op    string
	Left  Expr
	Right Expr
	Loc   ast.Loc
}

func (n *ComparisonExpr) Span() ast.Loc { return n.Loc }
func (*ComparisonExpr) exprNode()       {}

// QuantifiedComparisonExpr is `left <op> (ALL|SOME|ANY) ( query )`
// (quantifiedComparison). Subquery is the raw-text placeholder (B1 in expr.go).
type QuantifiedComparisonExpr struct {
	Op         string // comparison operator spelling
	Quantifier string // "ALL", "SOME", or "ANY"
	Left       Expr
	Subquery   *SubqueryExpr
	Loc        ast.Loc
}

func (n *QuantifiedComparisonExpr) Span() ast.Loc { return n.Loc }
func (*QuantifiedComparisonExpr) exprNode()       {}

// BetweenExpr is `value [NOT] BETWEEN lower AND upper` (the between predicate).
type BetweenExpr struct {
	Not   bool
	Value Expr
	Lower Expr
	Upper Expr
	Loc   ast.Loc
}

func (n *BetweenExpr) Span() ast.Loc { return n.Loc }
func (*BetweenExpr) exprNode()       {}

// InListExpr is `value [NOT] IN ( e , … )` (the inList predicate): membership
// against an explicit value list.
type InListExpr struct {
	Not   bool
	Value Expr
	List  []Expr
	Loc   ast.Loc
}

func (n *InListExpr) Span() ast.Loc { return n.Loc }
func (*InListExpr) exprNode()       {}

// InSubqueryExpr is `value [NOT] IN ( query )` (the inSubquery predicate).
// Subquery is the raw-text placeholder (B1 in expr.go).
type InSubqueryExpr struct {
	Not      bool
	Value    Expr
	Subquery *SubqueryExpr
	Loc      ast.Loc
}

func (n *InSubqueryExpr) Span() ast.Loc { return n.Loc }
func (*InSubqueryExpr) exprNode()       {}

// LikeExpr is `value [NOT] LIKE pattern [ESCAPE escape]` (the like predicate).
type LikeExpr struct {
	Not     bool
	Value   Expr
	Pattern Expr
	Escape  Expr // nil when no ESCAPE clause
	Loc     ast.Loc
}

func (n *LikeExpr) Span() ast.Loc { return n.Loc }
func (*LikeExpr) exprNode()       {}

// IsNullExpr is `value IS [NOT] NULL` (the nullPredicate).
type IsNullExpr struct {
	Not   bool
	Value Expr
	Loc   ast.Loc
}

func (n *IsNullExpr) Span() ast.Loc { return n.Loc }
func (*IsNullExpr) exprNode()       {}

// IsDistinctFromExpr is `value IS [NOT] DISTINCT FROM right` (the distinctFrom
// predicate): null-safe (in)equality.
type IsDistinctFromExpr struct {
	Not   bool
	Value Expr
	Right Expr
	Loc   ast.Loc
}

func (n *IsDistinctFromExpr) Span() ast.Loc { return n.Loc }
func (*IsDistinctFromExpr) exprNode()       {}

// parsePredicateSuffix attaches at most one predicate to the already-parsed
// valueExpression `left`, returning `left` unchanged when no predicate follows.
// It is the `predicate_?` of the predicated rule; the non-associativity (one
// predicate only, never a loop) is the oracle-confirmed Trino behavior.
//
// Dispatch is on the leading token of each predicate alternative:
//   - a comparison operator (= <> < <= > >=) → comparison or, if a quantifier
//     follows, quantifiedComparison;
//   - NOT (followed by BETWEEN | IN | LIKE) or a bare BETWEEN/IN/LIKE;
//   - IS (→ NULL or DISTINCT FROM).
func (p *Parser) parsePredicateSuffix(left Expr) (Expr, error) {
	switch p.cur.Kind {
	case int('='), tokNotEq, int('<'), tokLessEq, int('>'), tokGreaterEq:
		return p.parseComparison(left)
	case kwNOT:
		// NOT here introduces a negated BETWEEN/IN/LIKE predicate. A NOT that
		// begins a boolean (logicalNot) was already handled in parseNot before
		// the value expression, so reaching NOT in predicate position means a
		// negated predicate keyword must follow.
		switch p.peekNext().Kind {
		case kwBETWEEN, kwIN, kwLIKE:
			notTok := p.advance() // consume NOT
			return p.parseNegatablePredicate(left, true, notTok.Loc)
		default:
			// Not a predicate NOT (e.g. `a NOT b` is malformed); leave it for the
			// caller, which reports the error.
			return left, nil
		}
	case kwBETWEEN, kwIN, kwLIKE:
		return p.parseNegatablePredicate(left, false, ast.NoLoc())
	case kwIS:
		return p.parseIsPredicate(left)
	default:
		return left, nil
	}
}

// parseComparison parses `<op> right` or `<op> (ALL|SOME|ANY) ( query )` after a
// comparison operator on `left`. The quantified form is recognised when the
// token after the operator is a comparison quantifier directly followed by '('.
func (p *Parser) parseComparison(left Expr) (Expr, error) {
	opTok := p.advance()
	op := comparisonOpText(opTok)

	if q, ok := comparisonQuantifierText(p.cur.Kind); ok && p.peekNext().Kind == int('(') {
		p.advance()            // consume the quantifier
		openTok := p.advance() // consume '('
		subq, err := p.parseSubqueryPlaceholder(openTok.Loc.Start, SubqueryQuantified)
		if err != nil {
			return nil, err
		}
		return &QuantifiedComparisonExpr{
			Op:         op,
			Quantifier: q,
			Left:       left,
			Subquery:   subq,
			Loc:        ast.Loc{Start: left.Span().Start, End: subq.Loc.End},
		}, nil
	}

	right, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	return &ComparisonExpr{
		Op:    op,
		Left:  left,
		Right: right,
		Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
	}, nil
}

// parseNegatablePredicate parses a BETWEEN / IN / LIKE predicate (already known
// to be the current token) on `left`, with `not` set by the caller when a
// preceding NOT was consumed. notLoc is the NOT token's Loc (valid only when not
// is true) so the result span starts at NOT.
func (p *Parser) parseNegatablePredicate(left Expr, not bool, notLoc ast.Loc) (Expr, error) {
	start := left.Span().Start
	if not {
		start = notLoc.Start
	}
	switch p.cur.Kind {
	case kwBETWEEN:
		return p.parseBetween(left, not, start)
	case kwIN:
		return p.parseIn(left, not, start)
	case kwLIKE:
		return p.parseLike(left, not, start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseBetween parses `BETWEEN lower AND upper` (BETWEEN already current). The
// AND here is the BETWEEN separator, not a boolean conjunction, so both bounds
// are valueExpressions (not full booleanExpressions).
func (p *Parser) parseBetween(left Expr, not bool, start int) (Expr, error) {
	p.advance() // consume BETWEEN
	lower, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAND); err != nil {
		return nil, err
	}
	upper, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	return &BetweenExpr{
		Not:   not,
		Value: left,
		Lower: lower,
		Upper: upper,
		Loc:   ast.Loc{Start: start, End: upper.Span().End},
	}, nil
}

// parseIn parses `IN ( … )` (IN already current): either an explicit value list
// `IN ( e , … )` or a subquery `IN ( query )`, disambiguated by whether the
// token after '(' begins a query. The list form requires at least one element.
func (p *Parser) parseIn(left Expr, not bool, start int) (Expr, error) {
	p.advance() // consume IN
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}

	if p.startsQuery() {
		subq, err := p.parseSubqueryPlaceholder(openTok.Loc.Start, SubqueryIn)
		if err != nil {
			return nil, err
		}
		return &InSubqueryExpr{
			Not:      not,
			Value:    left,
			Subquery: subq,
			Loc:      ast.Loc{Start: start, End: subq.Loc.End},
		}, nil
	}

	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	list := []Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &InListExpr{
		Not:   not,
		Value: left,
		List:  list,
		Loc:   ast.Loc{Start: start, End: closeTok.Loc.End},
	}, nil
}

// parseLike parses `LIKE pattern [ESCAPE escape]` (LIKE already current). The
// pattern and escape are valueExpressions.
func (p *Parser) parseLike(left Expr, not bool, start int) (Expr, error) {
	p.advance() // consume LIKE
	pattern, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	like := &LikeExpr{
		Not:     not,
		Value:   left,
		Pattern: pattern,
		Loc:     ast.Loc{Start: start, End: pattern.Span().End},
	}
	if p.cur.Kind == kwESCAPE {
		p.advance() // consume ESCAPE
		escape, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		like.Escape = escape
		like.Loc.End = escape.Span().End
	}
	return like, nil
}

// parseIsPredicate parses `IS [NOT] NULL` or `IS [NOT] DISTINCT FROM right`
// (IS already current). Trino has no other IS predicate — `IS TRUE`, `IS FALSE`,
// `IS UNKNOWN` are SYNTAX_ERRORs, so a token other than NULL/DISTINCT after the
// optional NOT is an error.
func (p *Parser) parseIsPredicate(left Expr) (Expr, error) {
	isTok := p.advance() // consume IS
	not := false
	if p.cur.Kind == kwNOT {
		p.advance() // consume NOT
		not = true
	}
	switch p.cur.Kind {
	case kwNULL:
		nullTok := p.advance()
		return &IsNullExpr{
			Not:   not,
			Value: left,
			Loc:   ast.Loc{Start: left.Span().Start, End: nullTok.Loc.End},
		}, nil
	case kwDISTINCT:
		p.advance() // consume DISTINCT
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		right, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		return &IsDistinctFromExpr{
			Not:   not,
			Value: left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}, nil
	default:
		// Point the error at IS so the diagnostic is actionable for the common
		// `IS TRUE` / `IS UNKNOWN` mistake.
		return nil, &ParseError{Loc: isTok.Loc, Msg: "expected NULL or DISTINCT FROM after IS"}
	}
}

// comparisonOpText returns the source spelling of a comparison-operator token.
func comparisonOpText(tok Token) string {
	switch tok.Kind {
	case int('='):
		return "="
	case tokNotEq:
		// The lexer collapses both <> and != to tokNotEq; prefer the source text
		// when available so deparse can echo the original spelling.
		if tok.Str != "" {
			return tok.Str
		}
		return "<>"
	case int('<'):
		return "<"
	case tokLessEq:
		return "<="
	case int('>'):
		return ">"
	case tokGreaterEq:
		return ">="
	default:
		return tok.Str
	}
}

// comparisonQuantifierText maps a quantifier keyword kind to its spelling; ok is
// false for a non-quantifier kind.
func comparisonQuantifierText(kind TokenKind) (string, bool) {
	switch kind {
	case kwALL:
		return "ALL", true
	case kwSOME:
		return "SOME", true
	case kwANY:
		return "ANY", true
	default:
		return "", false
	}
}
