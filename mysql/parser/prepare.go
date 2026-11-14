package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parsePrepareStmt parses a PREPARE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/prepare.html
//
//	PREPARE stmt_name FROM preparable_stmt
func (p *Parser) parsePrepareStmt() (*nodes.PrepareStmt, error) {
	start := p.pos()
	p.advance() // consume PREPARE

	stmt := &nodes.PrepareStmt{Loc: nodes.Loc{Start: start}}

	// Statement name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// FROM
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// SQL string or variable
	if p.cur.Type == tokSCONST {
		stmt.Stmt = p.cur.Str
		p.advance()
	} else if p.isVariableRef() {
		vref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = "@" + vref.Name
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseExecuteStmt parses an EXECUTE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/execute.html
//
//	EXECUTE stmt_name [USING @var_name [, @var_name] ...]
func (p *Parser) parseExecuteStmt() (*nodes.ExecuteStmt, error) {
	start := p.pos()
	p.advance() // consume EXECUTE

	stmt := &nodes.ExecuteStmt{Loc: nodes.Loc{Start: start}}

	// Statement name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional USING: execute_var_ident is '@' ident_or_text, a user variable and nothing
	// else.
	//
	// Ref: execute_var_list and execute_var_ident in MySQL 8.0 sql/sql_yacc.yy.
	if p.cur.Type == kwUSING {
		p.advance()
		for {
			if !p.isVariableRef() {
				return nil, p.syntaxErrorAtCur()
			}
			v, err := p.parseVariableRef()
			if err != nil {
				return nil, err
			}
			if v.System {
				return nil, &ParseError{Message: "syntax error at or near \"@@\"", Position: v.Loc.Start}
			}
			stmt.Params = append(stmt.Params, v)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDeallocateStmt parses a DEALLOCATE PREPARE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/deallocate-prepare.html
//
//	{DEALLOCATE | DROP} PREPARE stmt_name
func (p *Parser) parseDeallocateStmt() (*nodes.DeallocateStmt, error) {
	start := p.pos()
	p.advance() // consume DEALLOCATE

	// PREPARE
	p.match(kwPREPARE)

	stmt := &nodes.DeallocateStmt{Loc: nodes.Loc{Start: start}}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
