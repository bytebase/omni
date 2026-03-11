package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-database.html
//
//	CREATE {DATABASE | SCHEMA} [IF NOT EXISTS] db_name
//	    [create_option] ...
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	    | [DEFAULT] COLLATE [=] collation_name
//	    | [DEFAULT] ENCRYPTION [=] {'Y' | 'N'}
func (p *Parser) parseCreateDatabaseStmt() (*nodes.CreateDatabaseStmt, error) {
	start := p.pos()
	p.advance() // consume DATABASE or SCHEMA

	stmt := &nodes.CreateDatabaseStmt{Loc: nodes.Loc{Start: start}}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// Database name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Options
	for {
		opt, ok := p.parseDatabaseOption()
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterDatabaseStmt parses an ALTER DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-database.html
//
//	ALTER {DATABASE | SCHEMA} [db_name]
//	    alter_option ...
func (p *Parser) parseAlterDatabaseStmt() (*nodes.AlterDatabaseStmt, error) {
	start := p.pos()
	p.advance() // consume DATABASE or SCHEMA

	stmt := &nodes.AlterDatabaseStmt{Loc: nodes.Loc{Start: start}}

	// Optional database name
	if p.isIdentToken() && p.cur.Type != kwDEFAULT && p.cur.Type != kwCHARACTER && p.cur.Type != kwCHARSET && p.cur.Type != kwCOLLATE {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// Options
	for {
		opt, ok := p.parseDatabaseOption()
		if !ok {
			break
		}
		stmt.Options = append(stmt.Options, opt)
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropDatabaseStmt parses a DROP DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-database.html
//
//	DROP {DATABASE | SCHEMA} [IF EXISTS] db_name
func (p *Parser) parseDropDatabaseStmt() (*nodes.DropDatabaseStmt, error) {
	start := p.pos()
	p.advance() // consume DATABASE or SCHEMA

	stmt := &nodes.DropDatabaseStmt{Loc: nodes.Loc{Start: start}}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// Database name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDatabaseOption parses a single database option.
func (p *Parser) parseDatabaseOption() (*nodes.DatabaseOption, bool) {
	start := p.pos()

	// Skip optional DEFAULT
	p.match(kwDEFAULT)

	switch {
	case p.cur.Type == kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=') // optional =
			if p.isIdentToken() {
				val, _, _ := p.parseIdentifier()
				return &nodes.DatabaseOption{
					Loc:   nodes.Loc{Start: start, End: p.pos()},
					Name:  "CHARACTER SET",
					Value: val,
				}, true
			}
		}
	case p.cur.Type == kwCHARSET:
		p.advance()
		p.match('=') // optional =
		if p.isIdentToken() {
			val, _, _ := p.parseIdentifier()
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "CHARACTER SET",
				Value: val,
			}, true
		}
	case p.cur.Type == kwCOLLATE:
		p.advance()
		p.match('=') // optional =
		if p.isIdentToken() {
			val, _, _ := p.parseIdentifier()
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "COLLATE",
				Value: val,
			}, true
		}
	}

	return nil, false
}
