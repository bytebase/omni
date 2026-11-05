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
		opt, ok, err := p.parseDatabaseOption(false)
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

	// Optional database name. Every token that can be an unquoted identifier is
	// read as one: that is how the engine resolves `ALTER DATABASE opt_db_name
	// alter_database_options` against the nameless arm with one token of
	// lookahead. CHARACTER, COLLATE, DEFAULT and READ are reserved words and so
	// start an option, while CHARSET and ENCRYPTION are not, which is why
	// `ALTER DATABASE CHARSET = 'utf8'` and `ALTER DATABASE ENCRYPTION = 'N'`
	// name a database and then run out of options -- a 1064 on MySQL 8.0 --
	// while `ALTER DATABASE COLLATE = 'utf8_bin'` and
	// `ALTER DATABASE READ ONLY = 1` parse.
	if p.isIdentToken() {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// Options
	for {
		opt, ok, err := p.parseDatabaseOption(true)
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

// parseDatabaseOption parses a single database option. alter says which
// statement it belongs to: READ ONLY is an ALTER DATABASE option and not a
// CREATE DATABASE one.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-database.html
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	  | [DEFAULT] COLLATE [=] collation_name
//	  | [DEFAULT] ENCRYPTION [=] {'Y' | 'N'}
//
//	alter_option:
//	    create_option
//	  | READ ONLY [=] {DEFAULT | 0 | 1}
func (p *Parser) parseDatabaseOption(alter bool) (*nodes.DatabaseOption, bool, error) {
	start := p.pos()

	// READ ONLY [=] {DEFAULT | 0 | 1}, which takes no DEFAULT prefix and is no
	// CREATE DATABASE option: `CREATE DATABASE d READ ONLY = 1` is a 1064 on
	// MySQL 8.0 while `ALTER DATABASE d READ ONLY = 1` parses, and so does
	// READ ONLY beside a charset option on either side. The value is required
	// and is one of those three -- a bare `READ ONLY` and `READ ONLY = 2`,
	// `= TRUE`, `= 'x'` and `= abc` are each a 1064, where the value used to
	// be whatever token came next.
	//
	// Ref: alter_database_options in MySQL 8.0.46 sql/sql_yacc.yy.
	if alter && p.cur.Type == kwREAD && p.peekNext().Type == kwONLY {
		p.advance()    // consume READ
		p.advance()    // consume ONLY
		p.matchEqual() // optional =
		var val string
		switch {
		case p.cur.Type == kwDEFAULT:
			val = "DEFAULT"
		case p.cur.Type == tokICONST && p.cur.Ival == 0:
			val = "0"
		case p.cur.Type == tokICONST && p.cur.Ival == 1:
			val = "1"
		default:
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance()
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "READ ONLY",
			Value: val,
		}, true, nil
	}

	// Optional DEFAULT prefix, taken only when an option that reads it follows.
	// Consuming it on its own dropped the word: `ALTER DATABASE d DEFAULT` and
	// `ALTER DATABASE d CHARSET = utf8 DEFAULT` both parsed, and both are a
	// 1064 on MySQL 8.0, MariaDB 11.8.8 and TiDB v8.5.5.
	if p.cur.Type == kwDEFAULT {
		switch p.peekNext().Type {
		case kwCHARACTER, kwCHARSET, kwCOLLATE, kwENCRYPTION:
			p.advance()
		default:
			return nil, false, nil
		}
	}

	switch {
	case p.cur.Type == kwCHARACTER:
		p.advance()
		if _, ok := p.match(kwSET); ok {
			p.matchEqual() // optional =
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
		p.matchEqual() // optional =
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
		p.matchEqual() // optional =
		val, err := p.parseCollationName()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "COLLATE",
			Value: val,
		}, true, nil
	case p.cur.Type == kwENCRYPTION:
		p.advance()
		p.matchEqual() // optional =
		// The value is a quoted 'Y' or 'N' and is required: `ENCRYPTION = N`
		// and a bare `ENCRYPTION` are both a 1064.
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		val := p.cur.Str
		p.advance()
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "ENCRYPTION",
			Value: val,
		}, true, nil
	}

	return nil, false, nil
}

// parseCharsetName parses the name a database option's CHARACTER SET or
// CHARSET takes: ident_or_text plus the reserved word BINARY. That is an
// identifier, a non-reserved keyword, a backquoted identifier, or a string
// literal in either quote. An integer and every other reserved keyword are a
// 1064 -- `CHARSET = 1` and `CHARSET = select` are a syntax error on
// MySQL 8.0.
//
// Ref: charset_name and ident_or_text in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseCharsetName() (string, error) {
	// BINARY is reserved, so parseIdentOrText would turn `CHARSET = binary`
	// away; it is a charset name all the same.
	if p.cur.Type == kwBINARY {
		tok := p.advance()
		return tok.Str, nil
	}
	name, _, err := p.parseIdentOrText()
	return name, err
}

// parseCollationName parses the name a database option's COLLATE takes. It is
// parseCharsetName's grammar: the collation name carries the same BINARY arm,
// so `COLLATE = binary` parses on MySQL 8.0.
//
// Ref: collation_name in MySQL 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseCollationName() (string, error) {
	return p.parseCharsetName()
}
