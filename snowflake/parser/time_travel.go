package parser

// Time-travel and change-tracking table clauses (T5.3):
//
//	{ AT | BEFORE } ( { TIMESTAMP | OFFSET | STATEMENT | STREAM } => <expr> )
//	CHANGES ( INFORMATION => { DEFAULT | APPEND_ONLY } ) { AT | BEFORE } (…) [END (…)]
//
// These attach to a table primary in the FROM clause (see select.go). The
// legacy SnowflakeParser.g4 restricts AT/BEFORE's value to `expr` for
// TIMESTAMP/OFFSET and `string` for STATEMENT/STREAM, and only allows
// STATEMENT for BEFORE; the Snowflake docs are broader (BEFORE accepts the
// same four anchors, and any expression is accepted as the value). The docs
// win — we accept all four anchors for both AT and BEFORE, with a general
// expression value.

import "github.com/bytebase/omni/snowflake/ast"

// parseTimeTravelClause parses a { AT | BEFORE } ( <anchor> => <expr> ) clause.
// The caller has verified that p.cur is kwAT or kwBEFORE.
func (p *Parser) parseTimeTravelClause() (*ast.TimeTravelClause, error) {
	kindTok := p.advance() // consume AT or BEFORE
	clause := &ast.TimeTravelClause{
		Loc: ast.Loc{Start: kindTok.Loc.Start},
	}
	switch kindTok.Type {
	case kwAT:
		clause.Kind = ast.TimeTravelAt
	case kwBEFORE:
		clause.Kind = ast.TimeTravelBefore
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	anchor, err := p.parseTimeTravelAnchor()
	if err != nil {
		return nil, err
	}
	clause.Anchor = anchor

	if _, err := p.expect(tokAssoc); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	clause.Expr = expr

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End
	return clause, nil
}

// parseTimeTravelAnchor consumes the anchor keyword (TIMESTAMP / OFFSET /
// STATEMENT / STREAM) and returns the corresponding enum.
func (p *Parser) parseTimeTravelAnchor() (ast.TimeTravelAnchor, error) {
	switch p.cur.Type {
	case kwTIMESTAMP:
		p.advance()
		return ast.TimeTravelTimestamp, nil
	case kwOFFSET:
		p.advance()
		return ast.TimeTravelOffset, nil
	case kwSTATEMENT:
		p.advance()
		return ast.TimeTravelStatement, nil
	case kwSTREAM:
		p.advance()
		return ast.TimeTravelStream, nil
	default:
		return 0, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected TIMESTAMP, OFFSET, STATEMENT, or STREAM in time-travel clause",
		}
	}
}

// parseChangesClause parses:
//
//	CHANGES ( INFORMATION => { DEFAULT | APPEND_ONLY } ) { AT | BEFORE } (…) [END (…)]
//
// The caller has verified that p.cur is kwCHANGES.
func (p *Parser) parseChangesClause() (*ast.ChangesClause, error) {
	changesTok := p.advance() // consume CHANGES
	clause := &ast.ChangesClause{
		Loc: ast.Loc{Start: changesTok.Loc.Start},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwINFORMATION); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokAssoc); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwDEFAULT:
		p.advance()
		clause.Info = ast.ChangesDefault
	case kwAPPEND_ONLY:
		p.advance()
		clause.Info = ast.ChangesAppendOnly
	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected DEFAULT or APPEND_ONLY in CHANGES(INFORMATION => …)",
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Required AT | BEFORE anchor.
	if p.cur.Type != kwAT && p.cur.Type != kwBEFORE {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected AT or BEFORE after CHANGES(…)",
		}
	}
	start, err := p.parseTimeTravelClause()
	if err != nil {
		return nil, err
	}
	clause.Start = start
	clause.Loc.End = start.Loc.End

	// Optional END ( … ) bound. END mirrors AT/BEFORE's anchor forms.
	if p.cur.Type == kwEND {
		end, err := p.parseEndClause()
		if err != nil {
			return nil, err
		}
		clause.End = end
		clause.Loc.End = end.Loc.End
	}

	return clause, nil
}

// parseEndClause parses END ( <anchor> => <expr> ), the optional upper bound
// of a CHANGES clause. It reuses the time-travel anchor grammar and records
// the result as a TimeTravelClause whose Kind is left as AT (END is not AT,
// but the anchor/value shape is identical and the Kind field is unused for
// END consumers).
func (p *Parser) parseEndClause() (*ast.TimeTravelClause, error) {
	endTok := p.advance() // consume END
	clause := &ast.TimeTravelClause{
		Kind: ast.TimeTravelAt,
		Loc:  ast.Loc{Start: endTok.Loc.Start},
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	anchor, err := p.parseTimeTravelAnchor()
	if err != nil {
		return nil, err
	}
	clause.Anchor = anchor
	if _, err := p.expect(tokAssoc); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	clause.Expr = expr
	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End
	return clause, nil
}
