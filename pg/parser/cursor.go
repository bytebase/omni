package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseDeclareCursorStmt parses a DECLARE CURSOR statement.
//
// Ref: gram.y DeclareCursorStmt
//
//	DeclareCursorStmt:
//	    DECLARE cursor_name cursor_options CURSOR opt_hold FOR SelectStmt
func (p *Parser) parseDeclareCursorStmt() *nodes.DeclareCursorStmt {
	p.advance() // consume DECLARE

	name, _ := p.parseCursorName()

	options := p.parseCursorOptions()

	p.expect(CURSOR)

	hold := p.parseOptHold()

	p.expect(FOR)

	query, _ := p.parseSelectNoParens()

	return &nodes.DeclareCursorStmt{
		Portalname: name,
		Options:    int(options | hold | nodes.CURSOR_OPT_FAST_PLAN),
		Query:      query,
	}
}

// parseCursorOptions parses cursor_options.
//
//	cursor_options:
//	    /* EMPTY */
//	    | cursor_options NO SCROLL
//	    | cursor_options SCROLL
//	    | cursor_options BINARY
//	    | cursor_options ASENSITIVE
//	    | cursor_options INSENSITIVE
func (p *Parser) parseCursorOptions() int {
	opts := 0
	for {
		switch p.cur.Type {
		case NO:
			p.advance() // consume NO
			p.expect(SCROLL)
			opts |= nodes.CURSOR_OPT_NO_SCROLL
		case SCROLL:
			p.advance()
			opts |= nodes.CURSOR_OPT_SCROLL
		case BINARY:
			p.advance()
			opts |= nodes.CURSOR_OPT_BINARY
		case ASENSITIVE:
			p.advance()
			opts |= nodes.CURSOR_OPT_ASENSITIVE
		case INSENSITIVE:
			p.advance()
			opts |= nodes.CURSOR_OPT_INSENSITIVE
		default:
			return opts
		}
	}
}

// parseOptHold parses opt_hold.
//
//	opt_hold:
//	    /* EMPTY */
//	    | WITH HOLD
//	    | WITHOUT HOLD
func (p *Parser) parseOptHold() int {
	if p.cur.Type == WITH {
		p.advance() // consume WITH
		p.expect(HOLD)
		return nodes.CURSOR_OPT_HOLD
	}
	if p.cur.Type == WITHOUT {
		p.advance() // consume WITHOUT
		p.expect(HOLD)
		return 0
	}
	return 0
}

// parseFetchStmt parses a FETCH or MOVE statement.
//
// Ref: gram.y FetchStmt
//
//	FetchStmt:
//	    FETCH fetch_args
//	    | MOVE fetch_args
func (p *Parser) parseFetchStmt() *nodes.FetchStmt {
	ismove := p.cur.Type == MOVE
	p.advance() // consume FETCH or MOVE

	stmt := p.parseFetchArgs()
	if stmt != nil {
		stmt.Ismove = ismove
	}
	return stmt
}

// parseFetchArgs parses fetch_args.
//
// Ref: gram.y fetch_args
func (p *Parser) parseFetchArgs() *nodes.FetchStmt {
	switch p.cur.Type {
	case NEXT:
		p.advance() // consume NEXT
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   1,
			Portalname: name,
		}

	case PRIOR:
		p.advance() // consume PRIOR
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_BACKWARD,
			HowMany:   1,
			Portalname: name,
		}

	case FIRST_P:
		p.advance() // consume FIRST
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_ABSOLUTE,
			HowMany:   1,
			Portalname: name,
		}

	case LAST_P:
		p.advance() // consume LAST
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_ABSOLUTE,
			HowMany:   -1,
			Portalname: name,
		}

	case ABSOLUTE_P:
		p.advance() // consume ABSOLUTE
		howMany := p.parseSignedIconst()
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_ABSOLUTE,
			HowMany:   howMany,
			Portalname: name,
		}

	case RELATIVE_P:
		p.advance() // consume RELATIVE
		howMany := p.parseSignedIconst()
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_RELATIVE,
			HowMany:   howMany,
			Portalname: name,
		}

	case ALL:
		p.advance() // consume ALL
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   nodes.FETCH_ALL,
			Portalname: name,
		}

	case FORWARD:
		p.advance() // consume FORWARD
		// FORWARD ALL, FORWARD SignedIconst, or FORWARD (alone)
		if p.cur.Type == ALL {
			p.advance() // consume ALL
			p.parseOptFromIn()
			name, _ := p.parseCursorName()
			return &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   nodes.FETCH_ALL,
				Portalname: name,
			}
		}
		if p.cur.Type == ICONST || p.cur.Type == '+' || p.cur.Type == '-' {
			howMany := p.parseSignedIconst()
			p.parseOptFromIn()
			name, _ := p.parseCursorName()
			return &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   howMany,
				Portalname: name,
			}
		}
		// FORWARD opt_from_in cursor_name
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   1,
			Portalname: name,
		}

	case BACKWARD:
		p.advance() // consume BACKWARD
		// BACKWARD ALL, BACKWARD SignedIconst, or BACKWARD (alone)
		if p.cur.Type == ALL {
			p.advance() // consume ALL
			p.parseOptFromIn()
			name, _ := p.parseCursorName()
			return &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   nodes.FETCH_ALL,
				Portalname: name,
			}
		}
		if p.cur.Type == ICONST || p.cur.Type == '+' || p.cur.Type == '-' {
			howMany := p.parseSignedIconst()
			p.parseOptFromIn()
			name, _ := p.parseCursorName()
			return &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   howMany,
				Portalname: name,
			}
		}
		// BACKWARD opt_from_in cursor_name
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_BACKWARD,
			HowMany:   1,
			Portalname: name,
		}

	case FROM, IN_P:
		// from_in cursor_name
		p.advance() // consume FROM or IN
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   1,
			Portalname: name,
		}

	case ICONST, '+', '-':
		// SignedIconst opt_from_in cursor_name
		howMany := p.parseSignedIconst()
		p.parseOptFromIn()
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   howMany,
			Portalname: name,
		}

	default:
		// Just cursor_name
		name, _ := p.parseCursorName()
		return &nodes.FetchStmt{
			Direction:  nodes.FETCH_FORWARD,
			HowMany:   1,
			Portalname: name,
		}
	}
}

// parseOptFromIn parses opt_from_in (optional FROM or IN).
//
//	opt_from_in:
//	    from_in
//	    | /* EMPTY */
func (p *Parser) parseOptFromIn() {
	if p.cur.Type == FROM || p.cur.Type == IN_P {
		p.advance()
	}
}

// parseClosePortalStmt parses a CLOSE statement.
//
// Ref: gram.y ClosePortalStmt
//
//	ClosePortalStmt:
//	    CLOSE cursor_name
//	    | CLOSE ALL
func (p *Parser) parseClosePortalStmt() *nodes.ClosePortalStmt {
	p.advance() // consume CLOSE

	if p.cur.Type == ALL {
		p.advance()
		return &nodes.ClosePortalStmt{Portalname: ""}
	}

	name, _ := p.parseCursorName()
	return &nodes.ClosePortalStmt{Portalname: name}
}
