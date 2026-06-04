package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `parser-match-recognize` DAG node (with
// match_recognize.go): it implements Trino's row-pattern grammar — the
// `rowPattern`, `patternPrimary`, and `patternQuantifier` rules that appear
// inside a `MATCH_RECOGNIZE ( … PATTERN ( rowPattern ) … )` clause (and inside
// the windowFrame's PATTERN). match_recognize.go owns the surrounding
// patternRecognition clause and the two integration hooks; this file owns the
// pattern expression itself.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	rowPattern
//	    : patternPrimary patternQuantifier?   # quantifiedPrimary
//	    | rowPattern rowPattern               # patternConcatenation
//	    | rowPattern VBAR_ rowPattern         # patternAlternation ;
//	patternPrimary
//	    : identifier                                          # patternVariable
//	    | LPAREN_ RPAREN_                                     # emptyPattern
//	    | PERMUTE_ LPAREN_ rowPattern (COMMA_ rowPattern)* RPAREN_ # patternPermutation
//	    | LPAREN_ rowPattern RPAREN_                          # groupedPattern
//	    | CARET_                                              # partitionStartAnchor
//	    | DOLLAR_                                             # partitionEndAnchor
//	    | LCURLYHYPHEN_ rowPattern RCURLYHYPHEN_              # excludedPattern ;
//	patternQuantifier
//	    : ASTERISK_ QUESTION_MARK_?                                       # zeroOrMore
//	    | PLUS_ QUESTION_MARK_?                                           # oneOrMore
//	    | QUESTION_MARK_ QUESTION_MARK_?                                  # zeroOrOne
//	    | LCURLY_ exactly=INTEGER_VALUE_ RCURLY_ QUESTION_MARK_?          # range
//	    | LCURLY_ atLeast=INTEGER_VALUE_? COMMA_ atMost=INTEGER_VALUE_? RCURLY_ QUESTION_MARK_? # range ;
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed structure
// (probed; see oracle_match_recognize_test.go for the differential corpus):
//
//	R1 (precedence & associativity). The ANTLR rule is ambiguous; Trino resolves
//	   it as a classic regex: a quantifier binds to its immediately-preceding
//	   patternPrimary, concatenation binds tighter than `|` alternation, and both
//	   concatenation and alternation are left-associative. So `A B | C D` is
//	   `(A B) | (C D)`, `A | B | C` is `((A | B) | C)`, and `A B C` is
//	   `((A B) C)`. parseRowPattern is a two-level precedence climber:
//	   parseRowPattern (alternation) → parsePatternConcat (concatenation) →
//	   parsePatternFactor (primary + optional single quantifier).
//	R2 (at most one quantifier per primary). `A**` / `A{2}{3}` are REJECTED — a
//	   patternQuantifier may not directly follow another. A quantified group
//	   `(A*)*` is the legal way to stack. parsePatternFactor consumes at most one
//	   quantifier and does not loop.
//	R3 (concatenation boundary). Concatenation continues while the next token can
//	   begin a patternPrimary (identifier-start, '(', PERMUTE, '^', '$', '{-').
//	   It stops at ')', '|', '-}', ',' (PERMUTE arg separator) and the trailing
//	   MATCH_RECOGNIZE keywords — none of which begin a primary.
//	R4 (empty pattern vs grouped pattern). A leading '(' is the emptyPattern `()`
//	   when immediately followed by ')', else a groupedPattern `( rowPattern )`.
//	R5 (anchors and the empty pattern are full primaries). `^`, `$`, and `()` may
//	   themselves carry a quantifier (`^+`, `()*`) and participate in
//	   concatenation/alternation — grammatically; later semantic analysis may
//	   reject (INVALID_LABEL), which is NOT a syntax error.

// ---------------------------------------------------------------------------
// Row-pattern node hierarchy (parser-package types; not ast.Node — matching the
// Relation / Expr / DataType convention, see relation.go / expr.go headers)
// ---------------------------------------------------------------------------

// RowPattern is the interface implemented by every node of a MATCH_RECOGNIZE
// row pattern (the rowPattern / patternPrimary rules). Span returns the source
// byte range; concrete fields are reached by a Go type switch.
type RowPattern interface {
	// Span returns the source byte range covered by the pattern.
	Span() ast.Loc
	// rowPatternNode is a marker preventing unrelated types from satisfying it.
	rowPatternNode()
}

// PatternConcat is a concatenation `left right` of two adjacent patterns (the
// patternConcatenation alternative). Built left-associatively, so a longer run
// `A B C` nests as PatternConcat{PatternConcat{A,B}, C}.
type PatternConcat struct {
	Left  RowPattern
	Right RowPattern
	Loc   ast.Loc
}

func (n *PatternConcat) Span() ast.Loc { return n.Loc }
func (*PatternConcat) rowPatternNode() {}

// PatternAlternation is an alternation `left | right` (the patternAlternation
// alternative). Built left-associatively (`A | B | C` nests as
// PatternAlternation{PatternAlternation{A,B}, C}); alternation binds looser than
// concatenation.
type PatternAlternation struct {
	Left  RowPattern
	Right RowPattern
	Loc   ast.Loc
}

func (n *PatternAlternation) Span() ast.Loc { return n.Loc }
func (*PatternAlternation) rowPatternNode() {}

// QuantifiedPattern is a patternPrimary with an optional trailing quantifier
// (the quantifiedPrimary alternative). Quantifier is nil for a bare primary.
type QuantifiedPattern struct {
	Primary    RowPattern
	Quantifier *PatternQuantifier // nil when the primary is unquantified
	Loc        ast.Loc
}

func (n *QuantifiedPattern) Span() ast.Loc { return n.Loc }
func (*QuantifiedPattern) rowPatternNode() {}

// PatternVariable is a primary pattern variable `identifier` (the patternVariable
// alternative): a reference to a row-pattern variable that is given meaning by a
// DEFINE (or the implicit always-true variable).
type PatternVariable struct {
	Name *ast.Identifier
	Loc  ast.Loc
}

func (n *PatternVariable) Span() ast.Loc { return n.Loc }
func (*PatternVariable) rowPatternNode() {}

// EmptyPattern is the empty pattern `()` (the emptyPattern alternative): matches
// an empty sequence of rows.
type EmptyPattern struct {
	Loc ast.Loc
}

func (n *EmptyPattern) Span() ast.Loc { return n.Loc }
func (*EmptyPattern) rowPatternNode() {}

// PatternPermutation is `PERMUTE ( rowPattern (, rowPattern)* )` (the
// patternPermutation alternative): matches any permutation of its operand
// patterns. Patterns holds the operands (at least one).
type PatternPermutation struct {
	Patterns []RowPattern
	Loc      ast.Loc
}

func (n *PatternPermutation) Span() ast.Loc { return n.Loc }
func (*PatternPermutation) rowPatternNode() {}

// GroupedPattern is a parenthesized pattern `( rowPattern )` (the groupedPattern
// alternative): groups a sub-pattern, e.g. so a quantifier applies to the whole.
type GroupedPattern struct {
	Inner RowPattern
	Loc   ast.Loc
}

func (n *GroupedPattern) Span() ast.Loc { return n.Loc }
func (*GroupedPattern) rowPatternNode() {}

// AnchorPattern is a partition-boundary anchor — `^` (start) or `$` (end) (the
// partitionStartAnchor / partitionEndAnchor alternatives). Start marks `^`;
// otherwise it is `$`.
type AnchorPattern struct {
	Start bool // true for `^` (partition start), false for `$` (partition end)
	Loc   ast.Loc
}

func (n *AnchorPattern) Span() ast.Loc { return n.Loc }
func (*AnchorPattern) rowPatternNode() {}

// ExcludedPattern is an exclusion `{- rowPattern -}` (the excludedPattern
// alternative): the matched rows of Inner are excluded from the output.
type ExcludedPattern struct {
	Inner RowPattern
	Loc   ast.Loc
}

func (n *ExcludedPattern) Span() ast.Loc { return n.Loc }
func (*ExcludedPattern) rowPatternNode() {}

// PatternQuantifierKind enumerates the patternQuantifier shapes.
type PatternQuantifierKind int

const (
	// QuantZeroOrMore is `*` (zeroOrMoreQuantifier).
	QuantZeroOrMore PatternQuantifierKind = iota
	// QuantOneOrMore is `+` (oneOrMoreQuantifier).
	QuantOneOrMore
	// QuantZeroOrOne is `?` (zeroOrOneQuantifier).
	QuantZeroOrOne
	// QuantRange is a `{ … }` range (rangeQuantifier), with the bounds in the
	// PatternQuantifier's AtLeast/AtMost/Exactly fields.
	QuantRange
)

// PatternQuantifier is a row-pattern quantifier (the patternQuantifier rule).
// Kind is the quantifier shape. Reluctant marks the trailing `?` (reluctant /
// non-greedy) marker. For QuantRange, Exactly is set for the `{n}` form (and
// AtLeast/AtMost are nil), otherwise AtLeast/AtMost hold the `{m,n}` bounds (each
// nil when its side is omitted: `{,n}` / `{m,}` / `{,}`).
type PatternQuantifier struct {
	Kind      PatternQuantifierKind
	Reluctant bool
	Exactly   *int64 // {n} exact count, nil unless the exact form
	AtLeast   *int64 // {m,…} lower bound, nil when omitted
	AtMost    *int64 // {…,n} upper bound, nil when omitted
	Loc       ast.Loc
}

// ---------------------------------------------------------------------------
// rowPattern parsing (precedence climbing: alternation → concat → factor)
// ---------------------------------------------------------------------------

// parseRowPattern parses a full `rowPattern`: the loosest-binding alternation
// level. It parses a concatenation, then folds a left-associative chain of
// `| concatenation` alternatives (R1). It is the entry point used by
// match_recognize.go's PATTERN ( … ) and PERMUTE / grouped / excluded
// sub-patterns.
func (p *Parser) parseRowPattern() (RowPattern, error) {
	left, err := p.parsePatternConcat()
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == int('|') {
		p.advance() // consume '|'
		right, err := p.parsePatternConcat()
		if err != nil {
			return nil, err
		}
		left = &PatternAlternation{
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// parsePatternConcat parses a `patternConcatenation`: one or more pattern
// factors juxtaposed. It parses a first factor, then folds a left-associative
// chain of additional factors while the next token begins a patternPrimary (R3).
func (p *Parser) parsePatternConcat() (RowPattern, error) {
	left, err := p.parsePatternFactor()
	if err != nil {
		return nil, err
	}
	for p.startsPatternPrimary() {
		right, err := p.parsePatternFactor()
		if err != nil {
			return nil, err
		}
		left = &PatternConcat{
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Span().Start, End: right.Span().End},
		}
	}
	return left, nil
}

// startsPatternPrimary reports whether the current token can begin a
// patternPrimary — i.e. whether concatenation should continue (R3). A primary
// starts with an identifier (a pattern variable, including non-reserved keywords
// usable as identifiers), '(' (empty or grouped), PERMUTE, '^', '$', or '{-'
// (an excluded pattern). Everything else — ')', '|', '-}', ',', and the
// MATCH_RECOGNIZE trailing keywords (SUBSET / DEFINE) — does not, ending the
// concatenation.
func (p *Parser) startsPatternPrimary() bool {
	switch p.cur.Kind {
	case int('('), kwPERMUTE, int('^'), int('$'), tokLCurlyHyphen:
		return true
	default:
		return isIdentifierStart(p.cur.Kind)
	}
}

// parsePatternFactor parses a `patternPrimary patternQuantifier?` (the
// quantifiedPrimary alternative): a single primary with at most one trailing
// quantifier (R2 — a second quantifier is not consumed, so `A**` fails when the
// stray `*` is later rejected as not beginning a primary/clause).
func (p *Parser) parsePatternFactor() (RowPattern, error) {
	primary, err := p.parsePatternPrimary()
	if err != nil {
		return nil, err
	}
	q, ok, err := p.tryParsePatternQuantifier()
	if err != nil {
		return nil, err
	}
	if !ok {
		return primary, nil
	}
	return &QuantifiedPattern{
		Primary:    primary,
		Quantifier: q,
		Loc:        ast.Loc{Start: primary.Span().Start, End: q.Loc.End},
	}, nil
}

// parsePatternPrimary parses one `patternPrimary`:
//
//   - identifier                  → PatternVariable
//   - ( )                         → EmptyPattern        (R4: '(' then ')')
//   - ( rowPattern )              → GroupedPattern      (R4: '(' then a pattern)
//   - PERMUTE ( rowPattern, … )   → PatternPermutation  (only when followed by '(')
//   - ^ | $                       → AnchorPattern
//   - {- rowPattern -}            → ExcludedPattern
//
// PERMUTE is a NON-RESERVED keyword, so a bare `permute` (not followed by '(') is
// an ordinary pattern variable, not the permutation form — oracle-confirmed:
// Trino 481 accepts `PATTERN (permute) DEFINE permute AS true`. Only `PERMUTE (`
// enters parsePatternPermutation; otherwise PERMUTE falls through to the
// identifier (PatternVariable) branch.
func (p *Parser) parsePatternPrimary() (RowPattern, error) {
	switch {
	case p.cur.Kind == int('('):
		return p.parseParenPattern()
	case p.cur.Kind == kwPERMUTE && p.peekNext().Kind == int('('):
		return p.parsePatternPermutation()
	case p.cur.Kind == int('^'):
		tok := p.advance()
		return &AnchorPattern{Start: true, Loc: tok.Loc}, nil
	case p.cur.Kind == int('$'):
		tok := p.advance()
		return &AnchorPattern{Start: false, Loc: tok.Loc}, nil
	case p.cur.Kind == tokLCurlyHyphen:
		return p.parseExcludedPattern()
	default:
		if !isIdentifierStart(p.cur.Kind) {
			return nil, p.exprErrorAt("expected a row-pattern variable, '(', PERMUTE, '^', '$', or '{-'")
		}
		id, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &PatternVariable{Name: id, Loc: id.Loc}, nil
	}
}

// parseParenPattern parses `( )` (emptyPattern) or `( rowPattern )`
// (groupedPattern); the '(' is the current token. The token after '(' decides
// (R4): an immediate ')' is the empty pattern; otherwise a grouped sub-pattern.
func (p *Parser) parseParenPattern() (RowPattern, error) {
	openTok := p.advance() // consume '('
	if p.cur.Kind == int(')') {
		closeTok := p.advance() // consume ')'
		return &EmptyPattern{Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}, nil
	}
	inner, err := p.parseRowPattern()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &GroupedPattern{
		Inner: inner,
		Loc:   ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parsePatternPermutation parses `PERMUTE ( ( rowPattern (, rowPattern)* )? )`
// (the patternPermutation alternative; PERMUTE is current). The operand list may
// be EMPTY — Trino 481 accepts `PERMUTE()` (divergence "match-recognize empty
// PERMUTE": the live oracle is more permissive than the legacy ANTLR grammar,
// which required `rowPattern (, rowPattern)*`; oracle decides). A leading or
// trailing comma is still rejected (`PERMUTE(,)` / `PERMUTE(A,)`), so the list,
// when present, is the standard comma-separated form.
func (p *Parser) parsePatternPermutation() (RowPattern, error) {
	permuteTok := p.advance() // consume PERMUTE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var patterns []RowPattern
	if p.cur.Kind != int(')') {
		first, err := p.parseRowPattern()
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseRowPattern()
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, next)
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &PatternPermutation{
		Patterns: patterns,
		Loc:      ast.Loc{Start: permuteTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseExcludedPattern parses `{- rowPattern -}` (the excludedPattern
// alternative; the '{-' token is current). The body is a full rowPattern (an
// empty `{- -}` is a syntax error — there is no inner emptyPattern shortcut).
func (p *Parser) parseExcludedPattern() (RowPattern, error) {
	openTok := p.advance() // consume '{-'
	inner, err := p.parseRowPattern()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(tokRCurlyHyphen)
	if err != nil {
		return nil, err
	}
	return &ExcludedPattern{
		Inner: inner,
		Loc:   ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// tryParsePatternQuantifier parses an optional `patternQuantifier` after a
// primary. It returns ok=false when the current token does not begin a
// quantifier (`*`, `+`, `?`, or `{`). It consumes exactly one quantifier; a
// second adjacent quantifier is left for the caller, which does not accept it
// (R2).
func (p *Parser) tryParsePatternQuantifier() (*PatternQuantifier, bool, error) {
	switch p.cur.Kind {
	case int('*'):
		tok := p.advance()
		q := &PatternQuantifier{Kind: QuantZeroOrMore, Loc: tok.Loc}
		p.consumeReluctant(q)
		return q, true, nil
	case int('+'):
		tok := p.advance()
		q := &PatternQuantifier{Kind: QuantOneOrMore, Loc: tok.Loc}
		p.consumeReluctant(q)
		return q, true, nil
	case tokQuestion, int('?'):
		tok := p.advance()
		q := &PatternQuantifier{Kind: QuantZeroOrOne, Loc: tok.Loc}
		p.consumeReluctant(q)
		return q, true, nil
	case int('{'):
		q, err := p.parseRangeQuantifier()
		if err != nil {
			return nil, false, err
		}
		return q, true, nil
	default:
		return nil, false, nil
	}
}

// consumeReluctant consumes an optional trailing reluctant `?` marker on a
// quantifier and records it (extending the quantifier's span).
func (p *Parser) consumeReluctant(q *PatternQuantifier) {
	if p.cur.Kind == tokQuestion || p.cur.Kind == int('?') {
		tok := p.advance()
		q.Reluctant = true
		q.Loc.End = tok.Loc.End
	}
}

// parseRangeQuantifier parses a `{ … }` range quantifier (the rangeQuantifier
// alternatives; '{' is current):
//
//	{ n }          → Exactly = n
//	{ m , n }      → AtLeast = m, AtMost = n
//	{ m , }        → AtLeast = m
//	{ , n }        → AtMost = n
//	{ , }          → both nil (unbounded both sides)
//
// followed by an optional reluctant `?`. The two grammar alternatives are
// distinguished by whether a ',' appears: `{ n }` has none; every comma form is
// the bounded-range alternative. (Trino's lexer tokenizes `{-` specially, so a
// lone '{' here is unambiguously the range-quantifier brace.)
func (p *Parser) parseRangeQuantifier() (*PatternQuantifier, error) {
	openTok := p.advance() // consume '{'
	q := &PatternQuantifier{Kind: QuantRange, Loc: ast.Loc{Start: openTok.Loc.Start}}

	// Optional leading integer (the `{n}` count or the `{m,…}` lower bound).
	var leading *int64
	if p.cur.Kind == tokInteger {
		v := p.advance().Ival
		leading = &v
	}

	if p.cur.Kind == int(',') {
		// Bounded-range form `{ m? , n? }`.
		p.advance() // consume ','
		q.AtLeast = leading
		if p.cur.Kind == tokInteger {
			v := p.advance().Ival
			q.AtMost = &v
		}
	} else {
		// Exact form `{ n }` — the leading integer is mandatory here.
		if leading == nil {
			return nil, p.exprErrorAt("expected an integer or ',' in row-pattern range quantifier")
		}
		q.Exactly = leading
	}

	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	q.Loc.End = closeTok.Loc.End
	p.consumeReluctant(q)
	return q, nil
}
