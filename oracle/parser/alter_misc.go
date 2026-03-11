package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAlterStmt dispatches ALTER statements based on the next keyword.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SQL-Statements-ALTER-ANALYTIC-VIEW-to-ALTER-SYSTEM.html
//
//	ALTER SESSION SET param = value [ param = value ... ]
//	ALTER SYSTEM  SET param = value [ param = value ... ]
//	ALTER SYSTEM  KILL SESSION 'sid,serial#'
//	ALTER INDEX   name ...
//	ALTER VIEW    name ...
//	ALTER SEQUENCE name ...
func (p *Parser) parseAlterStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwSESSION:
		return p.parseAlterSessionStmt(start)
	case kwSYSTEM:
		return p.parseAlterSystemStmt(start)
	case kwINDEX:
		return p.parseAlterGeneric(start, nodes.OBJECT_INDEX)
	case kwVIEW:
		return p.parseAlterGeneric(start, nodes.OBJECT_VIEW)
	case kwSEQUENCE:
		return p.parseAlterGeneric(start, nodes.OBJECT_SEQUENCE)
	case kwTABLE:
		return p.parseAlterGeneric(start, nodes.OBJECT_TABLE)
	default:
		// Unknown ALTER target — skip to semicolon or EOF.
		p.skipToSemicolon()
		return nil
	}
}

// parseAlterSessionStmt parses ALTER SESSION SET param = value [, ...].
func (p *Parser) parseAlterSessionStmt(start int) nodes.StmtNode {
	p.advance() // consume SESSION

	stmt := &nodes.AlterSessionStmt{
		Loc: nodes.Loc{Start: start},
	}

	if _, ok := p.matchKeyword(kwSET); ok {
		stmt.SetParams = p.parseSetParams()
	} else {
		// ALTER SESSION with other clauses — skip remaining tokens.
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterSystemStmt parses ALTER SYSTEM SET/KILL ... .
func (p *Parser) parseAlterSystemStmt(start int) nodes.StmtNode {
	p.advance() // consume SYSTEM

	stmt := &nodes.AlterSystemStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch {
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		stmt.SetParams = p.parseSetParams()
	default:
		// ALTER SYSTEM KILL SESSION or other — consume remaining tokens.
		// Capture KILL SESSION 'sid,serial#' if present.
		if p.isIdentLike() && p.cur.Str == "KILL" {
			p.advance() // consume KILL
			if p.cur.Type == kwSESSION {
				p.advance() // consume SESSION
			}
			if p.cur.Type == tokSCONST {
				stmt.Kill = p.cur.Str
				p.advance()
			}
		}
		p.skipToSemicolon()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSetParams parses one or more param = value pairs.
func (p *Parser) parseSetParams() *nodes.List {
	params := &nodes.List{}
	for {
		param := p.parseSetParam()
		if param == nil {
			break
		}
		params.Items = append(params.Items, param)
		// Some Oracle ALTER SESSION SET supports multiple params without commas;
		// but also handle comma separation.
		if p.cur.Type == ',' {
			p.advance()
		}
		// Stop if we hit end of statement.
		if !p.isIdentLike() {
			break
		}
	}
	return params
}

// parseSetParam parses a single name = value parameter setting.
func (p *Parser) parseSetParam() *nodes.SetParam {
	if !p.isIdentLike() {
		return nil
	}
	start := p.pos()
	name := p.parseIdentifier()
	if name == "" {
		return nil
	}

	// Expect '='
	if p.cur.Type != '=' {
		return &nodes.SetParam{
			Name: name,
			Loc:  nodes.Loc{Start: start, End: p.pos()},
		}
	}
	p.advance() // consume '='

	value := p.parseExpr()

	return &nodes.SetParam{
		Name:  name,
		Value: value,
		Loc:   nodes.Loc{Start: start, End: p.pos()},
	}
}

// parseAlterGeneric parses ALTER INDEX/VIEW/SEQUENCE/TABLE by consuming the
// object name and skipping the rest (simplified). Returns an AlterSessionStmt
// as a placeholder — in practice these would have their own AST types, but for
// now we skip the body to avoid blocking other work.
func (p *Parser) parseAlterGeneric(start int, _ nodes.ObjectType) nodes.StmtNode {
	p.advance() // consume INDEX/VIEW/SEQUENCE/TABLE

	// Parse the object name.
	_ = p.parseObjectName()

	// Skip remainder of the statement (clauses vary greatly by object type).
	p.skipToSemicolon()

	// Return a minimal AlterSessionStmt as placeholder.
	return &nodes.AlterSessionStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}
}

// skipToSemicolon advances until a semicolon or EOF is found.
// It does NOT consume the semicolon.
func (p *Parser) skipToSemicolon() {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}
}
