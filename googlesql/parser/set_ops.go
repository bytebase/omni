package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-select` DAG node. It implements GoogleSQL's
// set-operation chain (GoogleSQLParser.g4 §2.13 query_set_operation /
// set_operation_metadata), a hand-port of Google's open-source ZetaSQL
// reference. Set operations are left-associative and operate on query_primary
// items (a SELECT block or a parenthesized query); the chain is built by
// parseSetOpChain, called from parseQueryPrimaryOrSetOp (select.go).
//
// Grammar (set_operation_metadata):
//
//	[corresponding_outer] <UNION|EXCEPT|INTERSECT> [hint] <ALL|DISTINCT>
//	  [STRICT] [CORRESPONDING [BY]]
//
// where corresponding_outer ∈ {FULL [OUTER] | OUTER | LEFT [OUTER]}. The
// all_or_distinct choice is REQUIRED (a bare `UNION` with no ALL/DISTINCT is a
// syntax error in GoogleSQL — verified against the Spanner emulator). FROM-first
// queries directly after a set operator must be parenthesized (the grammar's
// FROM-after-set-op error alternatives); a non-parenthesized FROM there is a
// syntax error, which parseQueryPrimary already produces.

// parseSetOpChain consumes a left-associative chain of set operations following
// an already-parsed left query_primary. With no set operator it returns left
// unchanged. The chain `a UNION ALL b UNION ALL c` nests left:
// SetOperation{UNION, Left: SetOperation{UNION, a, b}, Right: c}.
//
// GoogleSQL requires every operation in a FLAT (unparenthesized) chain to be the
// SAME set operation — same operator AND same ALL/DISTINCT quantifier. Mixing
// (`UNION ALL … INTERSECT DISTINCT …`, or even `UNION ALL … UNION DISTINCT …`)
// is a syntax error: the operands must be parenthesized to group them (oracle:
// "Syntax error: Different set operations cannot be used in the same query
// without using parentheses for grouping"). A parenthesized operand is a
// complete query_primary (its own independent chain), so this scope check only
// compares the operations at THIS chain level. The corresponding-outer mode,
// STRICT, and CORRESPONDING modifiers are part of one operation's metadata and
// not compared across the chain (the operator + ALL/DISTINCT identity is the
// rule the oracle enforces).
func (p *Parser) parseSetOpChain(left ast.Node) (ast.Node, error) {
	var (
		haveFirst bool
		firstOp   ast.SetOp
		firstAll  bool
	)
	for {
		meta, ok, err := p.tryParseSetOpMetadata()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if !haveFirst {
			haveFirst = true
			firstOp = meta.Op
			firstAll = meta.All
		} else if meta.Op != firstOp || meta.All != firstAll {
			return nil, &ParseError{
				Loc: meta.Loc,
				Msg: "syntax error: different set operations cannot be used in the same query without using parentheses for grouping",
			}
		}
		right, err := p.parseQueryPrimary()
		if err != nil {
			return nil, err
		}
		meta.Left = left
		meta.Right = right
		meta.Loc = ast.Loc{Start: qLoc(left).Start, End: qLoc(right).End}
		left = meta
	}
	return left, nil
}

// tryParseSetOpMetadata parses one set_operation_metadata if the current
// position starts one. It returns (node, true, nil) with Left/Right/Loc unset
// (the caller fills them), (nil, false, nil) when no set operator is present, or
// (nil, false, err) on a malformed operator (e.g. a missing required
// ALL/DISTINCT).
//
// A set operator may be preceded by a corresponding-outer mode
// (FULL/OUTER/LEFT) — but those keywords are also join keywords, so we only
// treat a leading FULL/LEFT/OUTER as a set-op prefix when it is actually
// followed (after the optional OUTER) by UNION/INTERSECT/EXCEPT.
func (p *Parser) tryParseSetOpMetadata() (*ast.SetOperation, bool, error) {
	outerMode := ""
	startLoc := p.cur.Loc

	// Optional corresponding-outer mode: FULL [OUTER] | OUTER | LEFT [OUTER].
	// Only consume it if a set-op keyword follows (so a trailing LEFT/FULL that
	// belongs elsewhere is not eaten).
	switch p.cur.Type {
	case kwFULL, kwLEFT:
		if p.setOpFollowsOuterMode() {
			modeTok := p.advance() // FULL / LEFT
			outerMode = upperTokenName(modeTok.Type)
			if p.cur.Type == kwOUTER {
				p.advance()
			}
		}
	case kwOUTER:
		if p.peekNext().Type == kwUNION || p.peekNext().Type == kwINTERSECT || p.peekNext().Type == kwEXCEPT {
			p.advance() // OUTER
			outerMode = "OUTER"
		}
	}

	var op ast.SetOp
	switch p.cur.Type {
	case kwUNION:
		op = ast.SetOpUnion
	case kwINTERSECT:
		op = ast.SetOpIntersect
	case kwEXCEPT:
		op = ast.SetOpExcept
	default:
		if outerMode != "" {
			// We consumed a corresponding-outer mode but no set operator follows:
			// malformed.
			return nil, false, p.syntaxErrorAtCur()
		}
		return nil, false, nil
	}
	p.advance() // the set operator
	so := &ast.SetOperation{Op: op, OuterMode: outerMode, Loc: startLoc}

	// Optional hint between the operator and ALL/DISTINCT.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, false, herr
		}
	}

	// REQUIRED ALL | DISTINCT.
	switch p.cur.Type {
	case kwALL:
		p.advance()
		so.All = true
	case kwDISTINCT:
		p.advance()
		so.All = false
	default:
		return nil, false, &ParseError{Loc: p.cur.Loc, Msg: "syntax error: expected ALL or DISTINCT after set operator"}
	}

	// Optional STRICT.
	if p.cur.Type == kwSTRICT {
		p.advance()
		so.Strict = true
	}

	// Optional CORRESPONDING [BY]. (Spanner feature-rejects CORRESPONDING; it
	// parses in the union grammar — BigQuery supports it.)
	if p.cur.Type == kwCORRESPONDING {
		p.advance()
		so.Corresponding = true
		if p.cur.Type == kwBY {
			p.advance()
		}
	}

	return so, true, nil
}

// setOpFollowsOuterMode reports whether a leading FULL/LEFT (with an optional
// OUTER) is immediately followed by a UNION/INTERSECT/EXCEPT operator — i.e. it
// is a corresponding-outer set-op prefix rather than something else. It peeks
// one or two tokens ahead without consuming.
func (p *Parser) setOpFollowsOuterMode() bool {
	// FULL/LEFT immediately followed by the set operator.
	switch p.peekNext().Type {
	case kwUNION, kwINTERSECT, kwEXCEPT:
		return true
	case kwOUTER:
		// FULL OUTER / LEFT OUTER then the operator — needs a 3rd-token peek,
		// which the parser's 2-token lookahead cannot do directly. Fall back to a
		// throwaway lexer scan from the current position.
		return p.thirdTokenIsSetOp()
	}
	return false
}

// thirdTokenIsSetOp reports whether the token two positions after the current
// one is a set operator (used to recognize `FULL OUTER UNION` / `LEFT OUTER
// EXCEPT`). It scans a throwaway lexer from the current token's start so it does
// not mutate parser state.
func (p *Parser) thirdTokenIsSetOp() bool {
	startIdx := absIndex(p, p.cur.Loc.Start)
	lx := NewLexerWithOffset(p.input[startIdx:], 0)
	_ = lx.NextToken() // FULL/LEFT
	_ = lx.NextToken() // OUTER
	third := lx.NextToken()
	switch third.Type {
	case kwUNION, kwINTERSECT, kwEXCEPT:
		return true
	}
	return false
}

// upperTokenName returns the upper-case keyword spelling for a token type (used
// for the OuterMode label). It mirrors curIsWord's use of TokenName.
func upperTokenName(t int) string {
	return TokenName(t)
}
