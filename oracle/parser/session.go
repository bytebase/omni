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
			roles, err := p.parseSetRoleList()
			if err != nil {
				return nil, err
			}
			stmt.Except = roles
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// role [IDENTIFIED BY password] [,...]
	roles, err := p.parseSetRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetRoleList parses role [ IDENTIFIED BY password ] [, ...]. A role
// is one name (ORA-00933 at a dot); a password is a word, a quoted name or
// a string. Oracle 23ai takes IDENTIFIED BY after a role of ALL EXCEPT too.
func (p *Parser) parseSetRoleList() ([]*nodes.ObjectName, error) {
	var roles []*nodes.ObjectName
	for {
		start := p.pos()
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		roles = append(roles, &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}})
		if p.isKeywordStr("IDENTIFIED") {
			p.advance()
			if p.cur.Type != kwBY {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != tokSCONST && !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // password
		}
		if p.cur.Type != ',' {
			return roles, nil
		}
		p.advance()
	}
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
			// IMMEDIATE is reserved (ORA-00931 at it); DEFERRED is not.
			if p.cur.Type == kwIMMEDIATE {
				return nil, p.syntaxErrorAtCur()
			}
			name, parseErr1112 := p.parseObjectName()
			if parseErr1112 != nil {
				return nil, parseErr1112
			}
			if name == nil || name.Name == "" {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Constraints = append(stmt.Constraints, name)
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
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
