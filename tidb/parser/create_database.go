package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
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
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Completion: after CREATE DATABASE [IF NOT EXISTS], identifier context.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new database name.
		return nil, &ParseError{Message: "collecting"}
	}

	// Database name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Options
	for {
		opt, ok, err := p.parseDatabaseOption()
		if err != nil {
			return nil, err
		}
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
	if p.isIdentToken() && p.cur.Type != kwDEFAULT && p.cur.Type != kwCHARACTER && p.cur.Type != kwCHARSET && p.cur.Type != kwCOLLATE && p.cur.Type != kwPLACEMENT {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// Options
	for {
		opt, ok, err := p.parseDatabaseOption()
		if err != nil {
			return nil, err
		}
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
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Completion: after DROP DATABASE [IF EXISTS], offer database_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Database name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDatabaseOption parses a single database option.
//
//	create_option / alter_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	  | [DEFAULT] COLLATE [=] collation_name
//	  | [DEFAULT] ENCRYPTION [=] {'Y' | 'N'}
//	  | READ ONLY [=] {DEFAULT | 0 | 1}        (ALTER DATABASE only)
func (p *Parser) parseDatabaseOption() (*nodes.DatabaseOption, bool, error) {
	// Completion: offer database option keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwDEFAULT, kwCHARACTER, kwCHARSET, kwCOLLATE, kwENCRYPTION,
			kwREAD,
			// TiDB-specific: PLACEMENT POLICY and SET TIFLASH REPLICA
			// on CREATE/ALTER DATABASE (parser.y:4482 arm).
			kwPLACEMENT, kwSET,
		} {
			p.addTokenCandidate(t)
		}
		return nil, false, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	// READ ONLY [=] {DEFAULT | 0 | 1}  (no DEFAULT prefix)
	if p.cur.Type == kwREAD {
		p.advance()
		if _, ok := p.match(kwONLY); ok {
			p.match('=') // optional =
			val := p.cur.Str
			p.advance()
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "READ ONLY",
				Value: val,
			}, true, nil
		}
	}

	// Skip optional DEFAULT
	p.match(kwDEFAULT)

	switch {
	case p.cur.Type == kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=') // optional =
			val, _, err := p.parseIdent()
			if err != nil {
				return nil, false, err
			}
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "CHARACTER SET",
				Value: val,
			}, true, nil
		}
	case p.cur.Type == kwCHARSET:
		p.advance()
		p.match('=') // optional =
		val, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.pos()},
			Name:  "CHARACTER SET",
			Value: val,
		}, true, nil
	case p.cur.Type == kwCOLLATE:
		p.advance()
		p.match('=') // optional =
		val, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.pos()},
			Name:  "COLLATE",
			Value: val,
		}, true, nil
	case p.cur.Type == kwENCRYPTION:
		p.advance()
		p.match('=') // optional =
		if p.cur.Type == tokSCONST {
			val := p.cur.Str
			p.advance()
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
				Name:  "ENCRYPTION",
				Value: val,
			}, true, nil
		}
	case p.cur.Type == kwPLACEMENT:
		p.advance()
		if _, ok := p.match(kwPOLICY); !ok {
			return nil, false, p.syntaxErrorAtCur()
		}
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.pos()},
			Name:  "PLACEMENT POLICY",
			Value: val,
		}, true, nil
	case p.cur.Type == kwSET:
		// TiDB: CREATE/ALTER DATABASE ... SET TIFLASH REPLICA n [LOCATION LABELS ...]
		// Ref: parser.y:4482 DatabaseOption arm.
		// Note: no `=` sign allowed per upstream grammar.
		p.advance()
		if _, ok := p.match(kwTIFLASH); !ok {
			return nil, false, p.syntaxErrorAtCur()
		}
		if _, ok := p.match(kwREPLICA); !ok {
			return nil, false, p.syntaxErrorAtCur()
		}
		if p.cur.Type != tokICONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		count := int(p.cur.Ival)
		p.advance()
		labels, err := p.parseLocationLabelList()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:                   nodes.Loc{Start: start, End: p.pos()},
			Name:                  "TIFLASH REPLICA",
			TiFlashReplica:        count,
			TiFlashLocationLabels: labels,
		}, true, nil
	}

	return nil, false, nil
}
