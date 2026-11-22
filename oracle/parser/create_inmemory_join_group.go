package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateInmemoryJoinGroupStmt parses a CREATE INMEMORY JOIN GROUP
// statement. It is called after CREATE INMEMORY JOIN GROUP.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INMEMORY-JOIN-GROUP.html
//
//	CREATE INMEMORY JOIN GROUP [ IF NOT EXISTS ] [ schema. ] join_group
//	    ( join_group_member [, join_group_member ]... )
//
// An IF not followed by NOT EXISTS is refused at the IF (ORA-00900 on Oracle
// 23ai, which reads it as the group's name).
func (p *Parser) parseCreateInmemoryJoinGroupStmt(start int) (*nodes.CreateInmemoryJoinGroupStmt, error) {
	stmt := &nodes.CreateInmemoryJoinGroupStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		ifTok := p.cur
		p.advance() // consume IF
		if p.cur.Type != kwNOT || p.peekNext().Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance() // consume NOT
		p.advance() // consume EXISTS
		stmt.IfNotExists = true
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if stmt.Members, err = p.parseJoinGroupMembers(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterInmemoryJoinGroupStmt parses an ALTER INMEMORY JOIN GROUP
// statement. It is called after ALTER INMEMORY JOIN GROUP.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INMEMORY-JOIN-GROUP.html
//
//	ALTER INMEMORY JOIN GROUP [ IF EXISTS ] [ schema. ] join_group
//	    { ADD | REMOVE } ( join_group_member [, join_group_member ]... )
//
// The page takes one member; Oracle 23ai takes a list, as CREATE does.
func (p *Parser) parseAlterInmemoryJoinGroupStmt(start int) (*nodes.AlterInmemoryJoinGroupStmt, error) {
	stmt := &nodes.AlterInmemoryJoinGroupStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	switch {
	case p.cur.Type == kwADD:
		stmt.Action = "ADD"
	case p.isKeywordStr("REMOVE"):
		stmt.Action = "REMOVE"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ADD or REMOVE
	if stmt.Members, err = p.parseJoinGroupMembers(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseJoinGroupMembers parses the parenthesized member list of CREATE and
// ALTER INMEMORY JOIN GROUP:
//
//	( join_group_member [, join_group_member ]... )
//
//	join_group_member:
//	    [ schema. ] table ( column [, column ]... )
//
// The page names one column per table; Oracle 23ai takes a list. It refuses
// a reserved word as table (ORA-00903) or column (ORA-00900), a database link
// (ORA-00906), and an empty list.
func (p *Parser) parseJoinGroupMembers() ([]*nodes.JoinGroupMember, error) {
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	var members []*nodes.JoinGroupMember
	var firstTableTok Token
	for {
		if len(members) == 0 {
			firstTableTok = p.cur
		}
		m := &nodes.JoinGroupMember{Loc: nodes.Loc{Start: p.pos()}}
		table, err := p.parsePlainObjectName()
		if err != nil {
			return nil, err
		}
		m.Table = table
		if err := p.expectCreateTableToken('('); err != nil {
			return nil, err
		}
		for {
			if !p.isCreateTableNameStart() {
				return nil, p.syntaxErrorAtCur()
			}
			m.Columns = append(m.Columns, p.cur.Str)
			p.advance()
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ,
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return nil, err
		}
		m.Loc.End = p.prev.End
		members = append(members, m)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	if p.cur.Type != ')' {
		if p.cur.Type == tokEOF || (p.isCreateTableNameStart() && p.peekNext().Type == '(') {
			return nil, p.syntaxErrorAtTok(firstTableTok)
		}
		if p.peekNext().Type == ')' {
			p.advance()
		}
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume )
	return members, nil
}
