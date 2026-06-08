package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"

	pgparser "github.com/bytebase/omni/redshift/parser"
)

// --------------------------------------------------------------------------
// Section 4.1: OPEN Cursor
// --------------------------------------------------------------------------

// parseOpen parses an OPEN cursor statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-cursors.html#PLPGSQL-CURSOR-OPENING
//
//	OPEN bound_cursor [ ( [ arg [, ...] ] ) ] ;
//	OPEN unbound_cursor [ [ NO ] SCROLL ] FOR query ;
//	OPEN unbound_cursor [ [ NO ] SCROLL ] FOR EXECUTE query_string [ USING expression [, ...] ] ;
func (p *Parser) parseOpen() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume OPEN

	// Expect cursor variable name (identifier)
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return nil, p.errorf("syntax error at or near %q, expected cursor name after OPEN", p.tokenText(p.cur))
	}
	cursorVar := p.identText()
	p.advance()

	node := &ast.PLOpen{
		CursorVar: cursorVar,
	}

	// Check what follows: '(' for bound cursor with args, SCROLL/NO SCROLL/FOR for unbound, ';' for bound without args
	switch {
	case p.cur.Type == '(':
		// OPEN bound_cursor(args)
		args, err := p.parseOpenArgs()
		if err != nil {
			return nil, err
		}
		node.Args = args

	case p.isKeyword("SCROLL") || p.isKeyword("NO") || p.isKeyword("FOR"):
		// Unbound cursor: [NO] SCROLL FOR ...
		scroll := ast.ScrollNone
		if p.isKeyword("NO") {
			p.advance() // consume NO
			if err := p.expectKeyword("SCROLL"); err != nil {
				return nil, err
			}
			scroll = ast.ScrollNo
		} else if p.isKeyword("SCROLL") {
			p.advance() // consume SCROLL
			scroll = ast.ScrollYes
		}
		node.Scroll = scroll

		// Expect FOR
		if err := p.expectKeyword("FOR"); err != nil {
			return nil, err
		}

		// Check for EXECUTE (dynamic) or SELECT/WITH (static query)
		if p.isKeyword("EXECUTE") {
			p.advance() // consume EXECUTE

			// Collect dynamic query expression until USING or semicolon
			query, err := p.collectUntil("USING")
			if err != nil {
				return nil, err
			}
			node.DynQuery = query

			// Check for optional USING clause
			if p.isKeyword("USING") {
				p.advance() // consume USING
				params, err := p.parseUsingParams()
				if err != nil {
					return nil, err
				}
				node.Params = params
			}
		} else {
			// Static query: collect until semicolon
			query, err := p.collectUntilSemicolon()
			if err != nil {
				return nil, err
			}
			node.Query = query
		}

	case p.cur.Type == ';':
		// OPEN bound_cursor; (no args)

	default:
		return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
	}

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
	return node, nil
}

// parseOpenArgs parses parenthesized arguments for a bound cursor: (expr, expr, ...)
// Returns the argument text (the content between parentheses).
func (p *Parser) parseOpenArgs() (string, error) {
	p.advance() // consume '('
	start := p.pos()
	depth := 1

	for !p.isEOF() && depth > 0 {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
			if depth == 0 {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				p.advance() // consume ')'
				return text, nil
			}
		}
		p.advance()
	}

	return "", p.errorf("syntax error: unterminated argument list for cursor")
}

// --------------------------------------------------------------------------
// Section 4.2: FETCH and MOVE
// --------------------------------------------------------------------------

// parseFetch parses a FETCH cursor statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-cursors.html#PLPGSQL-CURSOR-USING
//
//	FETCH [ direction { FROM | IN } ] cursor INTO target ;
func (p *Parser) parseFetch() (ast.Node, error) {
	return p.parseFetchOrMove(false)
}

// parseMove parses a MOVE cursor statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-cursors.html#PLPGSQL-CURSOR-USING
//
//	MOVE [ direction { FROM | IN } ] cursor ;
func (p *Parser) parseMove() (ast.Node, error) {
	return p.parseFetchOrMove(true)
}

// parseFetchOrMove parses FETCH or MOVE statements.
// They share the same syntax except MOVE has no INTO clause.
func (p *Parser) parseFetchOrMove(isMove bool) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume FETCH or MOVE

	direction := ast.FetchNext
	count := ""
	cursorVar := ""

	// Try to parse direction and cursor variable.
	// The tricky part: a bare identifier could be a direction keyword OR a cursor variable name.
	// Strategy: try to detect known direction keywords, otherwise treat as cursor variable.
	dir, cnt, consumed := p.tryParseDirection()
	if consumed {
		direction = dir
		count = cnt

		// After direction, expect FROM or IN, then cursor name
		if !p.isKeyword("FROM") && !p.isKeyword("IN") {
			return nil, p.errorf("syntax error at or near %q, expected FROM or IN after direction", p.tokenText(p.cur))
		}
		p.advance() // consume FROM or IN

		if !p.isIdent() && !p.isAnyKeywordAsIdent() {
			return nil, p.errorf("syntax error at or near %q, expected cursor name", p.tokenText(p.cur))
		}
		cursorVar = p.identText()
		p.advance()
	} else {
		// No direction keyword recognized — the next token is the cursor variable name
		if !p.isIdent() && !p.isAnyKeywordAsIdent() {
			return nil, p.errorf("syntax error at or near %q, expected cursor name", p.tokenText(p.cur))
		}
		cursorVar = p.identText()
		p.advance()
	}

	// For FETCH: parse optional INTO clause
	var into []string
	if !isMove && p.isKeyword("INTO") {
		p.advance() // consume INTO
		targets, err := p.parseIntoTargets()
		if err != nil {
			return nil, err
		}
		into = targets
	}

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLFetch{
		CursorVar: cursorVar,
		Direction: direction,
		Count:     count,
		Into:      into,
		IsMove:    isMove,
		Loc:       ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// tryParseDirection attempts to parse a FETCH/MOVE direction keyword.
// Returns (direction, count_expr, true) if a direction was consumed,
// or (0, "", false) if the current token is not a recognized direction keyword.
//
// Directions: NEXT, PRIOR, FIRST, LAST, ABSOLUTE n, RELATIVE n,
// FORWARD [n|ALL], BACKWARD [n|ALL]
func (p *Parser) tryParseDirection() (ast.FetchDirection, string, bool) {
	switch {
	case p.isKeyword("NEXT"):
		next := p.peekNext()
		if next.Type == pgparser.FROM || isTokenKeyword(next, "IN") {
			p.advance() // consume NEXT
			return ast.FetchNext, "", true
		}
		return 0, "", false

	case p.isKeyword("PRIOR"):
		next := p.peekNext()
		if next.Type == pgparser.FROM || isTokenKeyword(next, "IN") {
			p.advance() // consume PRIOR
			return ast.FetchPrior, "", true
		}
		return 0, "", false

	case p.isKeyword("FIRST"):
		next := p.peekNext()
		if next.Type == pgparser.FROM || isTokenKeyword(next, "IN") {
			p.advance() // consume FIRST
			return ast.FetchFirst, "", true
		}
		return 0, "", false

	case p.isKeyword("LAST"):
		next := p.peekNext()
		if next.Type == pgparser.FROM || isTokenKeyword(next, "IN") {
			p.advance() // consume LAST
			return ast.FetchLast, "", true
		}
		return 0, "", false

	case p.isKeyword("ABSOLUTE"):
		p.advance() // consume ABSOLUTE
		count := p.collectDirectionCount()
		return ast.FetchAbsolute, count, true

	case p.isKeyword("RELATIVE"):
		p.advance() // consume RELATIVE
		count := p.collectDirectionCount()
		return ast.FetchRelative, count, true

	case p.isKeyword("FORWARD"):
		p.advance() // consume FORWARD
		if p.isKeyword("ALL") {
			p.advance() // consume ALL
			return ast.FetchForwardAll, "", true
		}
		if p.isKeyword("FROM") || p.isKeyword("IN") {
			return ast.FetchForward, "", true
		}
		count := p.collectDirectionCount()
		return ast.FetchForwardN, count, true

	case p.isKeyword("BACKWARD"):
		p.advance() // consume BACKWARD
		if p.isKeyword("ALL") {
			p.advance() // consume ALL
			return ast.FetchBackwardAll, "", true
		}
		if p.isKeyword("FROM") || p.isKeyword("IN") {
			return ast.FetchBackward, "", true
		}
		count := p.collectDirectionCount()
		return ast.FetchBackwardN, count, true
	}

	return 0, "", false
}

// collectDirectionCount collects the direction count expression for ABSOLUTE, RELATIVE, FORWARD n, BACKWARD n.
// This handles simple integers, negative integers (e.g., -1), and expressions.
func (p *Parser) collectDirectionCount() string {
	start := p.pos()

	// Handle optional minus sign (can appear as '-' char token or Op "-")
	if p.cur.Type == '-' || (p.cur.Type == pgparser.Op && p.cur.Str == "-") {
		p.advance()
	}

	// Consume the number or identifier
	if p.cur.Type == pgparser.ICONST || p.cur.Type == pgparser.FCONST || p.isIdent() || p.isAnyKeywordAsIdent() {
		end := p.cur.End
		p.advance()
		text := strings.TrimSpace(p.source[start:end])
		return text
	}

	// Fallback: use prev.End if we consumed a minus sign
	text := strings.TrimSpace(p.source[start:p.prev.End])
	return text
}

// isTokenKeyword checks if a token matches a keyword name without being the parser's current token.
func isTokenKeyword(tok pgparser.Token, name string) bool {
	lower := strings.ToLower(name)
	if tok.Type == pgparser.IDENT {
		return strings.EqualFold(tok.Str, name)
	}
	switch lower {
	case "in":
		return tok.Type == pgparser.IN_P
	case "from":
		return tok.Type == pgparser.FROM
	}
	return false
}
