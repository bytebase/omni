package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Sequence DDL — CREATE / ALTER SEQUENCE (T4.4)
// ---------------------------------------------------------------------------
//
// A SEQUENCE is a small, fully-enumerated object — unlike the open-ended STAGE /
// FILE FORMAT / table-variant grammars, its clause set is closed and stable
// (START / INCREMENT / ORDER | NOORDER / COMMENT), so each clause is modeled
// explicitly rather than captured open-ended. The START / INCREMENT values are
// signed 64-bit integers (a negative INCREMENT is documented and valid), so an
// optional leading '-' is consumed before the integer. The optional WITH / WITH /
// BY / '=' connector keywords are all accepted per the docs:
//
//	START [ WITH ] [ = ] <n>     INCREMENT [ BY ] [ = ] <n>
//
// (DROP / UNDROP SEQUENCE are handled by drop.go.)

// ---------------------------------------------------------------------------
// CREATE SEQUENCE
// ---------------------------------------------------------------------------

// parseCreateSequenceStmt parses
//
//	CREATE [ OR REPLACE ] SEQUENCE [ IF NOT EXISTS ] <name>
//	  [ WITH ]
//	  [ START [ WITH ] [ = ] <initial_value> ]
//	  [ INCREMENT [ BY ] [ = ] <sequence_interval> ]
//	  [ { ORDER | NOORDER } ]
//	  [ COMMENT = '<string_literal>' ]
//
// The CREATE keyword and the optional OR REPLACE modifier have already been
// consumed by parseCreateStmt; start is the Loc of the CREATE token and cur is the
// SEQUENCE keyword. The leading [ WITH ] is an optional no-op connector; the
// START / INCREMENT / ORDER|NOORDER / COMMENT clauses follow in that documented
// order.
func (p *Parser) parseCreateSequenceStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume SEQUENCE

	stmt := &ast.CreateSequenceStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional leading WITH connector (CREATE SEQUENCE <name> WITH START ...).
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
	}

	// START [ WITH ] [ = ] <n>.
	if p.cur.Type == kwSTART {
		p.advance() // consume START
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		if p.cur.Type == '=' {
			p.advance() // consume '='
		}
		v, err := p.parseSignedInt()
		if err != nil {
			return nil, err
		}
		stmt.Start = &v
	}

	// INCREMENT [ BY ] [ = ] <n>.
	if p.cur.Type == kwINCREMENT {
		p.advance() // consume INCREMENT
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		if p.cur.Type == '=' {
			p.advance() // consume '='
		}
		v, err := p.parseSignedInt()
		if err != nil {
			return nil, err
		}
		stmt.Increment = &v
	}

	// { ORDER | NOORDER }.
	if order, ok := p.parseOptionalOrderNoorder(); ok {
		stmt.Order = order
	}

	// COMMENT = '<string_literal>'.
	if p.cur.Type == kwCOMMENT {
		c, err := p.parseCommentEqString()
		if err != nil {
			return nil, err
		}
		stmt.Comment = c
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER SEQUENCE
// ---------------------------------------------------------------------------

// parseAlterSequenceStmt parses ALTER SEQUENCE [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the SEQUENCE keyword.
//
//	RENAME TO <new_name>
//	[ SET ] INCREMENT [ BY ] [ = ] <sequence_interval>
//	SET [ { ORDER | NOORDER } ] [ COMMENT = '<string_literal>' ]
//	UNSET COMMENT
func (p *Parser) parseAlterSequenceStmt() (ast.Node, error) {
	altTok := p.advance() // consume SEQUENCE
	stmt := &ast.AlterSequenceStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO <new_name>.
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSequenceRename
		stmt.NewName = newName

	case kwINCREMENT:
		// Bare (SET-less) INCREMENT [ BY ] [ = ] <n>.
		incr, err := p.parseSequenceIncrementClause()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSequenceSetIncrement
		stmt.Increment = &incr

	case kwSET:
		p.advance() // consume SET
		// SET may carry INCREMENT, or { ORDER | NOORDER } and/or COMMENT.
		if p.cur.Type == kwINCREMENT {
			incr, err := p.parseSequenceIncrementClause()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSequenceSetIncrement
			stmt.Increment = &incr
			break
		}
		// SET [ { ORDER | NOORDER } ] [ COMMENT = '...' ]. At least one of the two
		// must be present.
		sawSomething := false
		if order, ok := p.parseOptionalOrderNoorder(); ok {
			stmt.Order = order
			sawSomething = true
		}
		if p.cur.Type == kwCOMMENT {
			c, err := p.parseCommentEqString()
			if err != nil {
				return nil, err
			}
			stmt.Comment = c
			sawSomething = true
		}
		if !sawSomething {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterSequenceSet

	case kwUNSET:
		// UNSET COMMENT.
		p.advance() // consume UNSET
		if _, err := p.expect(kwCOMMENT); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSequenceUnsetComment

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSequenceIncrementClause parses INCREMENT [ BY ] [ = ] <n> and returns the
// signed value. cur is the INCREMENT keyword.
func (p *Parser) parseSequenceIncrementClause() (int64, error) {
	p.advance() // consume INCREMENT
	if p.cur.Type == kwBY {
		p.advance() // consume BY
	}
	if p.cur.Type == '=' {
		p.advance() // consume '='
	}
	return p.parseSignedInt()
}

// ---------------------------------------------------------------------------
// Shared sequence helpers
// ---------------------------------------------------------------------------

// parseSignedInt parses an optionally negation-prefixed integer literal and
// returns its int64 value. A leading '-' lexes as a separate token, so it is
// consumed explicitly and the following integer is negated; a leading '+' is
// also tolerated. Sequence START / INCREMENT accept any non-zero 64-bit value,
// which includes negatives.
func (p *Parser) parseSignedInt() (int64, error) {
	neg := false
	switch p.cur.Type {
	case '-':
		neg = true
		p.advance() // consume '-'
	case '+':
		p.advance() // consume '+'
	}
	tok, err := p.expect(tokInt)
	if err != nil {
		return 0, err
	}
	if neg {
		return -tok.Ival, nil
	}
	return tok.Ival, nil
}

// parseOptionalOrderNoorder consumes an optional ORDER / NOORDER token, returning
// (*bool, true) when present (true for ORDER, false for NOORDER) and (nil, false)
// when the current token is neither.
func (p *Parser) parseOptionalOrderNoorder() (*bool, bool) {
	switch p.cur.Type {
	case kwORDER:
		p.advance()
		b := true
		return &b, true
	case kwNOORDER:
		p.advance()
		b := false
		return &b, true
	}
	return nil, false
}

// parseCommentEqString parses COMMENT = '<string_literal>' and returns the string.
// The '=' is required here (matching the docs' `COMMENT = '...'`). cur is the
// COMMENT keyword.
func (p *Parser) parseCommentEqString() (*string, error) {
	p.advance() // consume COMMENT
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	tok, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	s := tok.Str
	return &s, nil
}
