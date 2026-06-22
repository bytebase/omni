package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseBatch parses a BATCH statement:
//
//	BEGIN [UNLOGGED|LOGGED] BATCH [USING TIMESTAMP n]
//	  (insert | update | delete) ';'
//	  ...
//	APPLY BATCH
func (p *Parser) parseBatch() (*ast.BatchStmt, error) {
	start := p.curLoc()
	if err := p.expectKeyword(tokBEGIN); err != nil {
		return nil, err
	}

	// Parse optional batch type: UNLOGGED or LOGGED
	batchType := ast.BatchDefault
	switch p.cur.Type {
	case tokUNLOGGED:
		batchType = ast.BatchUnlogged
		p.advance()
	case tokLOGGED:
		batchType = ast.BatchLogged
		p.advance()
	}

	if err := p.expectKeyword(tokBATCH); err != nil {
		return nil, err
	}

	// Optional USING clause (typically just TIMESTAMP for batches)
	using, err := p.parseUsingClause()
	if err != nil {
		return nil, err
	}

	// Parse inner DML statements until APPLY BATCH
	var stmts []ast.StmtNode
	for {
		// Check for APPLY BATCH
		if p.cur.Type == tokAPPLY {
			break
		}
		if p.cur.Type == tokEOF {
			return nil, p.errorf("unexpected end of input, expected APPLY BATCH")
		}

		// Skip any stray semicolons between statements
		if p.cur.Type == tokSEMI {
			p.advance()
			continue
		}

		// Parse an inner DML statement (INSERT, UPDATE, or DELETE)
		var stmt ast.StmtNode
		switch p.cur.Type {
		case tokINSERT:
			s, err := p.parseInsert()
			if err != nil {
				return nil, err
			}
			stmt = s
		case tokUPDATE:
			s, err := p.parseUpdate()
			if err != nil {
				return nil, err
			}
			stmt = s
		case tokDELETE:
			s, err := p.parseDelete()
			if err != nil {
				return nil, err
			}
			stmt = s
		default:
			return nil, p.errorf("expected INSERT, UPDATE, or DELETE inside BATCH, got %s", p.tokenDesc())
		}

		stmts = append(stmts, stmt)

		// Consume optional semicolon after each statement
		p.match(tokSEMI)
	}

	// Expect APPLY BATCH
	if err := p.expectKeyword(tokAPPLY); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokBATCH); err != nil {
		return nil, err
	}

	return &ast.BatchStmt{
		Type:       batchType,
		Using:      using,
		Statements: stmts,
		Loc:        p.makeLoc(start),
	}, nil
}
