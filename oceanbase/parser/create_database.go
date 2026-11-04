package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-database.html, minus the
// options 5.7.25-OceanBase_CE-v4.3.5.6 answers 1064 to (read on 2026-11-04).
//
//	CREATE {DATABASE | SCHEMA} [IF NOT EXISTS] db_name
//	    [create_option] ...
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	    | [DEFAULT] COLLATE [=] collation_name
func (p *Parser) parseCreateDatabaseStmt(start int) (*nodes.CreateDatabaseStmt, error) {
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

	// Options. CREATE DATABASE takes none at all; ALTER DATABASE needs one.
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterDatabaseStmt parses an ALTER DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-database.html
//
//	ALTER {DATABASE | SCHEMA} [db_name]
//	    alter_option ...
func (p *Parser) parseAlterDatabaseStmt(start int) (*nodes.AlterDatabaseStmt, error) {
	p.advance() // consume DATABASE or SCHEMA

	stmt := &nodes.AlterDatabaseStmt{Loc: nodes.Loc{Start: start}}

	// Optional database name. CHARACTER, COLLATE, DEFAULT and READ are reserved
	// words and so start an option rather than a name. CHARSET is not reserved,
	// and this is where OceanBase departs from the engines mysql forked from:
	// CE 4.3.5.6 takes `ALTER DATABASE CHARSET = 'utf8'` and
	// `ALTER DATABASE CHARSET utf8` as the nameless arm, while MySQL 8.0,
	// MariaDB 11.8.8 and TiDB v8.5.5 read CHARSET as the database name there
	// and answer 1064 for the options they then run out of. Every other
	// identifier-capable token is a name here as it is on those three:
	// `ALTER DATABASE ENCRYPTION = 'N'` and `ALTER DATABASE COMMENT = 'x'` are
	// each a 1064 on CE 4.3.5.6 too.
	if p.isIdentToken() && p.cur.Type != kwCHARSET {
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
	// alter_option is a one-or-more list, so an ALTER DATABASE that names no
	// option is a 1064 -- `ALTER DATABASE`, `ALTER DATABASE d` and
	// `ALTER DATABASE d DEFAULT` all parsed before this check.
	if len(stmt.Options) == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropDatabaseStmt parses a DROP DATABASE/SCHEMA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-database.html
//
//	DROP {DATABASE | SCHEMA} [IF EXISTS] db_name
func (p *Parser) parseDropDatabaseStmt(start int) (*nodes.DropDatabaseStmt, error) {
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDatabaseOption parses a single database option.
//
//	create_option / alter_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	  | [DEFAULT] COLLATE [=] collation_name
//	  | READ ONLY
//	  | READ WRITE
//
// OceanBase has neither of MySQL's other two database options: 4.3.5.6 answers
// 1064 to ENCRYPTION on CREATE and on ALTER, and to the COMMENT option MariaDB
// takes. Its READ ONLY and READ WRITE carry no value and read on both
// statements. Read on 2026-11-04 against 5.7.25-OceanBase_CE-v4.3.5.6;
// databaseOptionCases holds the verdicts.
func (p *Parser) parseDatabaseOption() (*nodes.DatabaseOption, bool, error) {
	start := p.pos()

	// READ ONLY and READ WRITE, which OceanBase spells with no value at all,
	// on CREATE DATABASE as well as ALTER DATABASE and beside a charset option
	// on either side. CE 4.3.5.6 takes `ALTER DATABASE d READ ONLY`,
	// `CREATE DATABASE d READ WRITE` and `ALTER DATABASE d CHARSET = utf8 READ
	// ONLY`, and answers 1064 to `READ ONLY = 1`, `READ ONLY 1`,
	// `READ ONLY = DEFAULT`, a bare `READ`, and a comma before it. MySQL's
	// valued `READ ONLY [=] {DEFAULT | 0 | 1}` is a 1064 here, which is why
	// mysql reads one shape and this parser another.
	if p.cur.Type == kwREAD {
		switch p.peekNext().Type {
		case kwONLY, kwWRITE:
			name := "READ ONLY"
			if p.peekNext().Type == kwWRITE {
				name = "READ WRITE"
			}
			p.advance() // consume READ
			p.advance() // consume ONLY or WRITE
			return &nodes.DatabaseOption{
				Loc:  nodes.Loc{Start: start, End: p.prev.End},
				Name: name,
			}, true, nil
		}
	}

	// Optional DEFAULT prefix, taken only when an option that reads it follows.
	// Consuming it on its own dropped the word: `ALTER DATABASE d DEFAULT` and
	// `ALTER DATABASE d CHARSET = utf8 DEFAULT` both parsed, and both are a
	// 1064 on OceanBase CE 4.3.5.6 as on MySQL 8.0.
	if p.cur.Type == kwDEFAULT {
		switch p.peekNext().Type {
		case kwCHARACTER, kwCHARSET, kwCOLLATE:
			p.advance()
		default:
			return nil, false, nil
		}
	}

	switch {
	case p.cur.Type == kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.match('=') // optional =
			val, err := p.parseCharsetName()
			if err != nil {
				return nil, false, err
			}
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "CHARACTER SET",
				Value: val,
			}, true, nil
		}
		// CHARACTER with no SET after it is no option, and the word has been
		// consumed: saying "not an option" here dropped it, so
		// `ALTER DATABASE d CHARSET = utf8 CHARACTER` parsed.
		return nil, false, p.syntaxErrorAtCur()
	case p.cur.Type == kwCHARSET:
		p.advance()
		p.match('=') // optional =
		val, err := p.parseCharsetName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "CHARACTER SET",
			Value: val,
		}, true, nil
	case p.cur.Type == kwCOLLATE:
		p.advance()
		p.match('=') // optional =
		val, err := p.parseCollationName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "COLLATE",
			Value: val,
		}, true, nil
	}

	return nil, false, nil
}

// parseCharsetName parses the name a database option's CHARACTER SET or
// CHARSET takes: ident_or_text plus the reserved word BINARY. That is an
// identifier, a non-reserved keyword, a backquoted identifier, or a string
// literal in either quote. An integer and every other reserved keyword are a
// 1064 -- `CHARSET = 1` and `CHARSET = select` are a syntax error on OceanBase
// CE 4.3.5.6 as on MySQL 8.0.
//
// Ref: charset_name and ident_or_text in MySQL 8.0.46 sql/sql_yacc.yy, which
// OceanBase forked; read against 5.7.25-OceanBase_CE-v4.3.5.6.
func (p *Parser) parseCharsetName() (string, error) {
	// BINARY is reserved, so parseIdentOrText would turn `CHARSET = binary`
	// away; it is a charset name all the same, and CE 4.3.5.6 takes both
	// `CHARSET = binary` and `CHARACTER SET binary`.
	if p.cur.Type == kwBINARY {
		tok := p.advance()
		return tok.Str, nil
	}
	name, _, err := p.parseIdentOrText()
	return name, err
}

// parseCollationName parses the name a database option's COLLATE takes. This
// is where OceanBase departs from MySQL: its collation name carries no BINARY
// arm, so CE 4.3.5.6 answers 1064 to `COLLATE = binary` and to
// `COLLATE binary` while taking `CHARSET = binary`. Quoted, the same word is
// an ordinary ident_or_text name and parses, in single quotes and in
// backquotes alike.
func (p *Parser) parseCollationName() (string, error) {
	name, _, err := p.parseIdentOrText()
	return name, err
}
