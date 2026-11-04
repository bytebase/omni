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

	// Optional database name. Every token that can be an unquoted identifier is
	// read as one: that is how the engine resolves `ALTER DATABASE DBName
	// DatabaseOptionList` against the nameless arm with one token of lookahead.
	// CHARACTER, COLLATE, DEFAULT and SET are reserved words and so start an
	// option, while CHARSET and ENCRYPTION are not, which is why
	// `ALTER DATABASE CHARSET = 'utf8'` and `ALTER DATABASE ENCRYPTION = 'N'`
	// name a database and then run out of options -- a 1064 on TiDB v8.5.5 --
	// while `ALTER DATABASE COLLATE = 'utf8_bin'` and
	// `ALTER DATABASE SET TIFLASH REPLICA 1` parse.
	//
	// PLACEMENT is the one exception: it is a non-reserved keyword, yet TiDB
	// v8.5.5 reads `ALTER DATABASE PLACEMENT POLICY = 'p'` as the option rather
	// than as a database named PLACEMENT, answering 8239 for the unknown policy.
	if p.isIdentToken() && p.cur.Type != kwPLACEMENT {
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
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-database.html
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	  | [DEFAULT] COLLATE [=] collation_name
//	  | [DEFAULT] ENCRYPTION [=] {'Y' | 'N'}
//	  | [DEFAULT] PLACEMENT POLICY [=] policy_name
//	  | SET TIFLASH REPLICA count [LOCATION LABELS label, ...]
//
//	alter_option:
//	    create_option
//
// TiDB v8.5.5 has no READ ONLY option on either statement, and answers 1064 to
// `ALTER DATABASE d READ ONLY = 1`, so this parser does not read one.
func (p *Parser) parseDatabaseOption() (*nodes.DatabaseOption, bool, error) {
	// Completion: offer database option keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{
			kwDEFAULT, kwCHARACTER, kwCHARSET, kwCOLLATE, kwENCRYPTION,
			// TiDB-specific: PLACEMENT POLICY and SET TIFLASH REPLICA
			// on CREATE/ALTER DATABASE (parser.y:4482 arm).
			kwPLACEMENT, kwSET,
		} {
			p.addTokenCandidate(t)
		}
		return nil, false, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	// Optional DEFAULT prefix, taken only when an option that reads it follows.
	// Consuming it on its own dropped the word: `ALTER DATABASE d DEFAULT` and
	// `ALTER DATABASE d CHARSET = utf8 DEFAULT` both parsed, and both are a
	// 1064 on MySQL 8.0, MariaDB 11.8.8 and TiDB v8.5.5.
	// TiDB takes it before PLACEMENT POLICY as well.
	if p.cur.Type == kwDEFAULT {
		switch p.peekNext().Type {
		case kwCHARACTER, kwCHARSET, kwCOLLATE, kwENCRYPTION, kwPLACEMENT:
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
	case p.cur.Type == kwENCRYPTION:
		p.advance()
		p.match('=') // optional =
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
	case p.cur.Type == kwPLACEMENT:
		p.advance()
		if _, ok := p.match(kwPOLICY); !ok {
			return nil, false, p.syntaxErrorAtCur()
		}
		// PLACEMENT POLICY SET DEFAULT is an arm of its own, and DEFAULT is
		// the sentinel that clears the policy. Reading SET as the policy name
		// left the option holding "SET" and dropped the DEFAULT after it, so
		// `ALTER DATABASE d PLACEMENT POLICY SET DEFAULT` parsed and then came
		// back from the catalog as an unknown policy named SET.
		if _, ok := p.match(kwSET); ok {
			if _, err := p.expect(kwDEFAULT); err != nil {
				return nil, false, err
			}
			return &nodes.DatabaseOption{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Name:  "PLACEMENT POLICY",
				Value: "DEFAULT",
			}, true, nil
		}
		p.match('=')
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, false, err
		}
		return &nodes.DatabaseOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  "PLACEMENT POLICY",
			Value: val,
		}, true, nil
	case p.cur.Type == kwSET:
		// TiDB: CREATE/ALTER DATABASE ... SET TIFLASH REPLICA n [LOCATION LABELS ...]
		// Ref: parser.y:4482 DatabaseOption arm.
		// Note: no `=` sign allowed per upstream grammar. The DEFAULT prefix is
		// not allowed on this arm either (DefaultKwdOpt is missing from the
		// grammar production), which the prefix above refuses before we get
		// here by leaving DEFAULT unconsumed.
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
			Loc:                   nodes.Loc{Start: start, End: p.prev.End},
			Name:                  "TIFLASH REPLICA",
			TiFlashReplica:        count,
			TiFlashLocationLabels: labels,
		}, true, nil
	}

	return nil, false, nil
}

// parseCharsetName parses the name a database option's CHARACTER SET or
// CHARSET takes: ident_or_text plus the reserved word BINARY. That is an
// identifier, a non-reserved keyword, a backquoted identifier, or a string
// literal in either quote. An integer and every other reserved keyword are a
// 1064 -- `CHARSET = 1` and `CHARSET = select` are a syntax error on
// TiDB v8.5.5.
//
// Ref: CharsetName and StringName in TiDB v8.5.5 parser/parser.y; charset_name
// and ident_or_text in MySQL 8.0.46 sql/sql_yacc.yy.
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
// so `COLLATE = binary` parses on TiDB v8.5.5.
//
// Ref: CollationName in TiDB v8.5.5 parser/parser.y; collation_name in MySQL
// 8.0.46 sql/sql_yacc.yy.
func (p *Parser) parseCollationName() (string, error) {
	return p.parseCharsetName()
}
