package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseSetRoleStmt parses a SET ROLE statement.
//
// BNF: oracle/parser/bnf/SET-ROLE.bnf
//
//	SET ROLE
//	    { role [ IDENTIFIED BY password ]
//	      [, role [ IDENTIFIED BY password ] ]...
//	    | ALL [ EXCEPT role [, role ]... ]
//	    | NONE
//	    } ;
func (p *Parser) parseSetRoleStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume SET
	p.advance() // consume ROLE

	stmt := &nodes.SetRoleStmt{
		Loc: nodes.Loc{Start: start},
	}

	// NONE
	if p.isIdentLike() && p.cur.Str == "NONE" {
		stmt.None = true
		p.advance()
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// ALL [EXCEPT role [,...]]
	if p.cur.Type == kwALL {
		stmt.All = true
		p.advance()
		if p.cur.Type == kwEXCEPT {
			p.advance()
			for {
				name, parseErr1110 := p.parseObjectName()
				if parseErr1110 != nil {
					return nil, parseErr1110
				}
				if name != nil {
					stmt.Except = append(stmt.Except, name)
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// role [IDENTIFIED BY password] [,...]
	for {
		name, parseErr1111 := p.parseObjectName()
		if parseErr1111 != nil {
			return nil, parseErr1111
		}
		if name != nil {
			stmt.Roles = append(stmt.Roles, name)
		}
		// Skip optional IDENTIFIED BY password
		if p.isIdentLike() && p.cur.Str == "IDENTIFIED" {
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			}
			// consume password (identifier or string)
			p.advance()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetConstraintsStmt parses a SET CONSTRAINT(S) statement.
//
// Note: SET-CONSTRAINT.bnf not present; syntax derived from Oracle documentation.
//
//	SET { CONSTRAINT | CONSTRAINTS }
//	    { ALL | constraint [, constraint ]... }
//	    { IMMEDIATE | DEFERRED } ;
func (p *Parser) parseSetConstraintsStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume SET
	p.advance() // consume CONSTRAINT or CONSTRAINTS

	stmt := &nodes.SetConstraintsStmt{
		Loc: nodes.Loc{Start: start},
	}

	// ALL or specific constraints
	if p.cur.Type == kwALL {
		stmt.All = true
		p.advance()
	} else {
		for {
			name, parseErr1112 := p.parseObjectName()
			if parseErr1112 != nil {
				return nil, parseErr1112
			}
			if name != nil {
				stmt.Constraints = append(stmt.Constraints, name)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	// IMMEDIATE or DEFERRED
	if p.cur.Type == kwIMMEDIATE {
		p.advance()
	} else if p.cur.Type == kwDEFERRED {
		stmt.Deferred = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
