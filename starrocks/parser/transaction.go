package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// parseBeginStmt parses:
//
//	BEGIN [WITH LABEL label_name]
//
// cur must be kwBEGIN on entry; it is consumed here.
func (p *Parser) parseBeginStmt() (ast.Node, error) {
	startTok := p.advance() // consume BEGIN
	stmt := &ast.BeginStmt{}
	endLoc := startTok.Loc

	// Optional WITH LABEL label_name
	if p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect(kwLABEL); err != nil {
			return nil, err
		}
		// label_name is an identifier (possibly a non-reserved keyword)
		name, nameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Label = name
		endLoc = nameLoc
	}

	stmt.Loc = startTok.Loc.Merge(endLoc)
	return stmt, nil
}

// parseCommitStmt parses:
//
//	COMMIT [WORK] [AND [NO] CHAIN] [[NO] RELEASE]
//
// cur must be kwCOMMIT on entry; it is consumed here.
func (p *Parser) parseCommitStmt() (ast.Node, error) {
	startTok := p.advance() // consume COMMIT
	stmt := &ast.CommitStmt{}
	endLoc := startTok.Loc

	// Optional WORK
	if p.cur.Kind == kwWORK {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Work = true
	}

	// Optional AND [NO] CHAIN
	if p.cur.Kind == kwAND {
		p.advance() // consume AND
		if p.cur.Kind == kwNO {
			endLoc = p.cur.Loc
			p.advance() // consume NO
			if _, err := p.expect(kwCHAIN); err != nil {
				return nil, err
			}
			stmt.Chain = "AND NO CHAIN"
		} else {
			if _, err := p.expect(kwCHAIN); err != nil {
				return nil, err
			}
			stmt.Chain = "AND CHAIN"
		}
		endLoc = p.prev.Loc
	}

	// Optional [NO] RELEASE
	if p.cur.Kind == kwNO {
		endLoc = p.cur.Loc
		p.advance() // consume NO
		if _, err := p.expect(kwRELEASE); err != nil {
			return nil, err
		}
		stmt.Release = "NO RELEASE"
		endLoc = p.prev.Loc
	} else if p.cur.Kind == kwRELEASE {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Release = "RELEASE"
	}

	stmt.Loc = startTok.Loc.Merge(endLoc)
	return stmt, nil
}

// parseRollbackStmt parses:
//
//	ROLLBACK [WORK] [AND [NO] CHAIN] [[NO] RELEASE]
//
// cur must be kwROLLBACK on entry; it is consumed here.
func (p *Parser) parseRollbackStmt() (ast.Node, error) {
	startTok := p.advance() // consume ROLLBACK
	stmt := &ast.RollbackStmt{}
	endLoc := startTok.Loc

	// Optional WORK
	if p.cur.Kind == kwWORK {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Work = true
	}

	// Optional AND [NO] CHAIN
	if p.cur.Kind == kwAND {
		p.advance() // consume AND
		if p.cur.Kind == kwNO {
			endLoc = p.cur.Loc
			p.advance() // consume NO
			if _, err := p.expect(kwCHAIN); err != nil {
				return nil, err
			}
			stmt.Chain = "AND NO CHAIN"
		} else {
			if _, err := p.expect(kwCHAIN); err != nil {
				return nil, err
			}
			stmt.Chain = "AND CHAIN"
		}
		endLoc = p.prev.Loc
	}

	// Optional [NO] RELEASE
	if p.cur.Kind == kwNO {
		endLoc = p.cur.Loc
		p.advance() // consume NO
		if _, err := p.expect(kwRELEASE); err != nil {
			return nil, err
		}
		stmt.Release = "NO RELEASE"
		endLoc = p.prev.Loc
	} else if p.cur.Kind == kwRELEASE {
		endLoc = p.cur.Loc
		p.advance()
		stmt.Release = "RELEASE"
	}

	stmt.Loc = startTok.Loc.Merge(endLoc)
	return stmt, nil
}
