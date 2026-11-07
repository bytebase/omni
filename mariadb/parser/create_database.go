package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
	"strings"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE/SCHEMA statement.
//
// Ref: https://mariadb.com/kb/en/create-database/
//
//	CREATE [OR REPLACE] {DATABASE | SCHEMA} [IF NOT EXISTS] db_name
//	    [create_option] ...
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	    | [DEFAULT] COLLATE [=] collation_name
func (p *Parser) parseCreateDatabaseStmt(start int, orReplace bool) (*nodes.CreateDatabaseStmt, error) {
	p.advance() // consume DATABASE or SCHEMA

	stmt := &nodes.CreateDatabaseStmt{Loc: nodes.Loc{Start: start}, OrReplace: orReplace}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

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
	if err := refuseDatabaseCharsetConflict(stmt.Options); err != nil {
		return nil, err
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
	// lookahead. CHARACTER, COLLATE and DEFAULT are reserved words and so start
	// an option, while CHARSET is not, which is why
	// `ALTER DATABASE CHARSET = 'utf8'` names a database and then runs out of
	// options -- a 1064 on MariaDB 11.8.8 -- while
	// `ALTER DATABASE COLLATE = 'utf8_bin'` parses.
	if p.isIdentToken() {
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
	if err := refuseDatabaseCharsetConflict(stmt.Options); err != nil {
		return nil, err
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
// Ref: https://mariadb.com/kb/en/alter-database/
//
//	create_option:
//	    [DEFAULT] CHARACTER SET [=] charset_name
//	  | [DEFAULT] COLLATE [=] collation_name
//
//	alter_option:
//	    create_option
//
// MariaDB 11.8.8 has neither of MySQL's other two: it answers 1064 to
// ENCRYPTION and to READ ONLY on both statements, so this parser reads
// neither. It does take a COMMENT option and DEFAULT as a charset or
// collation name, which this parser does not read yet.
func (p *Parser) parseDatabaseOption() (*nodes.DatabaseOption, bool, error) {
	start := p.pos()

	// Optional DEFAULT prefix, taken only when an option that reads it follows.
	// Consuming it on its own dropped the word: `ALTER DATABASE d DEFAULT` and
	// `ALTER DATABASE d CHARSET = utf8 DEFAULT` both parsed, and both are a
	// 1064 on MySQL 8.0, MariaDB 11.8.8 and TiDB v8.5.5.
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
			p.matchEqual() // optional =
			val, err := p.parseCharsetNameOrDefault()
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
		val, err := p.parseCharsetNameOrDefault()
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
		val, err := p.parseCharsetNameOrDefault()
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
// 1064 -- `CHARSET = 1` and `CHARSET = select` are a syntax error on
// MariaDB 11.8.8.
//
// Ref: charset_name and ident_or_text in MySQL 8.0.46 sql/sql_yacc.yy, which
// MariaDB forked; read against MariaDB 11.8.8.
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
// so `COLLATE = binary` parses on MariaDB 11.8.8.
//
// Ref: collation_name in MySQL 8.0.46 sql/sql_yacc.yy; read against MariaDB
// 11.8.8.
// parseCollationName parses a collation name, where DEFAULT names the
// collation the column would have had anyway and so settles nothing: MariaDB
// accepts `VARCHAR(50) COLLATE DEFAULT` and reads the column back with no
// collation of its own, while it refuses `CHARACTER SET DEFAULT` on a column
// with a 1064. Measured on mariadb:11.8.8.
//
// Ref: collation_name and charset_name in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseCollationName() (string, error) {
	if _, ok := p.match(kwDEFAULT); ok {
		return "", nil
	}
	return p.parseCharsetName()
}

// parseCharsetNameOrDefault is parseCharsetName where the DEFAULT keyword names the server's
// own: a database or table option takes it, and a column's CHARACTER SET does
// not (1064 on mariadb:11.8.8). A column's COLLATE does — see
// parseCollationName.
//
// Ref: charset_name_or_default and collation_name_or_default in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseCharsetNameOrDefault() (string, error) {
	if tok, ok := p.match(kwDEFAULT); ok {
		return tok.Str, nil
	}
	return p.parseCharsetName()
}

// knownCharsets holds the character sets of MariaDB 11.8, read from
// information_schema.CHARACTER_SETS, with utf8 as the alias of utf8mb3 it is.
var knownCharsets = map[string]bool{
	"armscii8": true, "ascii": true, "big5": true, "binary": true, "cp1250": true, "cp1251": true,
	"cp1256": true, "cp1257": true, "cp850": true, "cp852": true, "cp866": true, "cp932": true,
	"dec8": true, "eucjpms": true, "euckr": true, "gb2312": true, "gbk": true, "geostd8": true,
	"greek": true, "hebrew": true, "hp8": true, "keybcs2": true, "koi8r": true, "koi8u": true,
	"latin1": true, "latin2": true, "latin5": true, "latin7": true, "macce": true, "macroman": true,
	"sjis": true, "swe7": true, "tis620": true, "ucs2": true, "ujis": true, "utf16": true,
	"utf16le": true, "utf32": true, "utf8mb3": true, "utf8mb4": true,
}

func canonicalCharset(name string) string {
	name = strings.ToLower(name)
	if name == "utf8" {
		return "utf8mb3"
	}
	return name
}

// collationCharset names the character set a collation belongs to, which is what comes
// before its first underscore, and "" for a name no character set starts.
func collationCharset(name string) string {
	charset, _, _ := strings.Cut(name, "_")
	charset = canonicalCharset(charset)
	if !knownCharsets[charset] {
		return ""
	}
	return charset
}

// charsetOption is one CHARACTER SET or COLLATE among a database's or table's options.
type charsetOption struct {
	collate bool
	value   string
	pos     int
}

// refuseCharsetConflict is the action of these options, which merges each with what the
// options before it settled and answers 1302, ER_CONFLICTING_DECLARATIONS: two character
// sets that differ, or two collations, on one character set or not. DEFAULT clashes with no
// character set, and with a collation it does. Anything the server has to look up answers
// first with its own code (1115, 1253, 1273), and those stay reading: no name is checked past
// the first one that would. A collation whose name starts with a known character set is taken
// at its word.
//
// Ref: table_option and merge_charset_and_collation in MariaDB 11.8 sql/sql_yacc.yy and
// sql/lex_charset.cc.
func refuseCharsetConflict(opts []charsetOption) error {
	charset, collation := "", ""
	for _, opt := range opts {
		isDefault := strings.EqualFold(opt.value, "DEFAULT")
		if opt.collate {
			if !isDefault && collationCharset(opt.value) == "" {
				return nil
			}
			name := strings.ToLower(opt.value)
			if collation != "" && collation != name {
				return &ParseError{
					Message:  "Conflicting declarations: 'COLLATE " + collation + "' and 'COLLATE " + name + "'",
					Position: opt.pos,
				}
			}
			collation = name
			continue
		}
		if isDefault {
			continue
		}
		name := canonicalCharset(opt.value)
		if !knownCharsets[name] {
			return nil
		}
		if charset != "" && charset != name {
			return &ParseError{
				Message:  "Conflicting declarations: 'CHARACTER SET " + charset + "' and 'CHARACTER SET " + name + "'",
				Position: opt.pos,
			}
		}
		charset = name
	}
	return nil
}

func refuseDatabaseCharsetConflict(opts []*nodes.DatabaseOption) error {
	var list []charsetOption
	for _, opt := range opts {
		switch opt.Name {
		case "CHARACTER SET":
			list = append(list, charsetOption{value: opt.Value, pos: opt.Loc.Start})
		case "COLLATE":
			list = append(list, charsetOption{collate: true, value: opt.Value, pos: opt.Loc.Start})
		}
	}
	return refuseCharsetConflict(list)
}

func refuseTableCharsetConflict(opts []*nodes.TableOption) error {
	var list []charsetOption
	for _, opt := range opts {
		switch opt.Name {
		case "CHARACTER SET", "CHARSET":
			list = append(list, charsetOption{value: opt.Value, pos: opt.Loc.Start})
		case "COLLATE":
			list = append(list, charsetOption{collate: true, value: opt.Value, pos: opt.Loc.Start})
		}
	}
	return refuseCharsetConflict(list)
}
