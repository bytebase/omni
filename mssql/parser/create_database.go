// Package parser - create_database.go implements T-SQL CREATE DATABASE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// Shared SQL Server size unit suffix: {KB | MB | GB | TB}.
var sizeUnitValues = newOptionSet().withIdents("KB", "MB", "GB", "TB")

// parseCreateDatabaseStmt parses a CREATE DATABASE statement. Every clause is
// read by its words, as SQL Server 2022 reads it.
//
// BNF: mssql/parser/bnf/create-database-transact-sql.bnf
// Grammar: TSql160.g:4427 (createDatabase), TSql160.g:4582 (recoveryUnitList),
// TSql160.g:4869 (dbAddendums), TSql160.g:4889 (withCreateDbOptions).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-transact-sql
//
//	CREATE DATABASE database_name
//	[ CONTAINMENT = { NONE | PARTIAL } ]
//	[ ON
//	      [ PRIMARY ] <filespec> [ , ...n ]
//	      [ , <filegroup> [ , ...n ] ]
//	      [ LOG ON <filespec> [ , ...n ] ]
//	]
//	[ COLLATE collation_name ]
//	[ FOR { ATTACH [ WITH <attach_database_option> [ , ...n ] ] | ATTACH_REBUILD_LOG }
//	| AS SNAPSHOT OF source_database_name ]
//	[ WITH <option> [ , ...n ] ]
//
//	<option> ::=
//	{
//	      FILESTREAM ( <filestream_option> [ , ...n ] )
//	    | DEFAULT_FULLTEXT_LANGUAGE = { lcid | language_name | language_alias }
//	    | DEFAULT_LANGUAGE = { lcid | language_name | language_alias }
//	    | NESTED_TRIGGERS = { OFF | ON }
//	    | TRANSFORM_NOISE_WORDS = { OFF | ON }
//	    | TWO_DIGIT_YEAR_CUTOFF = <two_digit_year_cutoff>
//	    | DB_CHAINING { OFF | ON }
//	    | TRUSTWORTHY { OFF | ON }
//	    | PERSISTENT_LOG_BUFFER = { OFF | ON ( DIRECTORY_NAME = 'path' ) }
//	    | LEDGER = { ON | OFF }
//	    | CATALOG_COLLATION = { DATABASE_DEFAULT | SQL_Latin1_General_CP1_CI_AS }
//	}
//
//	<attach_database_option> ::=
//	      <option>
//	    | { ENABLE_BROKER | NEW_BROKER | ERROR_BROKER_CONVERSATIONS }
//	    | RESTRICTED_USER
//	    | FILESTREAM ( DIRECTORY_NAME = { 'directory_name' | NULL } )
//
// SQL Server 2022 refuses the Azure forms ScriptDOM also reads (AS COPY OF,
// ( EDITION = ... )), DATA_RETENTION, a LOG ON with no data file (Msg 188),
// and a snapshot's WITH unless the snapshot names its files.
func (p *Parser) parseCreateDatabaseStmt() (*nodes.CreateDatabaseStmt, error) {
	loc := p.pos()

	// Completion: after CREATE DATABASE → identifier (new database name)
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.CreateDatabaseStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.dbExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CONTAINMENT = { NONE | PARTIAL } (TSql160.g:4128)
	if p.dbWordIs("CONTAINMENT") {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		v, err := p.dbExpectWordIn("NONE", "PARTIAL")
		if err != nil {
			return nil, err
		}
		stmt.Containment = v
	}

	// ON [ PRIMARY ] <filespec> [ , ...n ] [ , <filegroup> [ , ...n ] ] (TSql160.g:4593)
	hasOn := false
	if p.cur.Type == kwON {
		p.advance()
		hasOn = true
		if err := p.parseCreateDatabaseOnDisk(stmt); err != nil {
			return nil, err
		}
	}

	// LOG ON <filespec> [ , ...n ]. A log file needs a data file (Msg 188),
	// and takes no PRIMARY (Msg 1041) and no filegroup (Msg 1064).
	if hasOn && p.dbWordIs("LOG") {
		p.advance()
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		files, err := p.parseDatabaseFileSpecList(dbFileCreate)
		if err != nil {
			return nil, err
		}
		stmt.LogOn = files
	}

	// COLLATE collation_name
	if p.cur.Type == kwCOLLATE {
		p.advance()
		collation, err := p.parseCollationName()
		if err != nil {
			return nil, err
		}
		stmt.Collation = collation
	}

	// FOR { ATTACH | ATTACH_REBUILD_LOG } | AS SNAPSHOT OF name (TSql160.g:4869)
	withAllowed := true
	switch {
	case p.cur.Type == kwFOR:
		p.advance()
		w, err := p.dbExpectWordIn("ATTACH", "ATTACH_REBUILD_LOG")
		if err != nil {
			return nil, err
		}
		if w == "ATTACH" {
			stmt.ForAttach = true
		} else {
			stmt.ForAttachRebuildLog = true
		}
	case p.cur.Type == kwAS:
		p.advance()
		if err := p.dbExpectWord("SNAPSHOT"); err != nil {
			return nil, err
		}
		if err := p.dbExpectWord("OF"); err != nil {
			return nil, err
		}
		src, err := p.dbExpectName()
		if err != nil {
			return nil, err
		}
		stmt.SnapshotOf = src
		// SQL Server ends a snapshot that names no files before WITH.
		withAllowed = hasOn
	}

	// WITH <option> [ , ...n ]
	if withAllowed && p.cur.Type == kwWITH {
		p.advance()
		opts, err := p.parseDatabaseWithOptions(stmt.ForAttach)
		if err != nil {
			return nil, err
		}
		if stmt.ForAttach {
			stmt.AttachOptions = opts
		} else {
			stmt.WithOptions = opts
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateDatabaseOnDisk parses the data files after ON: files of the
// primary filegroup, each optionally marked PRIMARY, then filegroups. Once a
// filegroup starts, a following file belongs to it, and PRIMARY is refused.
//
//	[ PRIMARY ] <filespec> { , [ PRIMARY ] <filespec> }* { , <filegroup> }*
func (p *Parser) parseCreateDatabaseOnDisk(stmt *nodes.CreateDatabaseStmt) error {
	var files []nodes.Node
	for {
		if p.cur.Type == kwPRIMARY {
			p.advance()
		}
		spec, err := p.parseDatabaseFileSpec(dbFileCreate)
		if err != nil {
			return err
		}
		files = append(files, spec)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.dbWordIs("FILEGROUP") {
			break
		}
	}
	stmt.OnPrimary = &nodes.List{Items: files}

	for p.dbWordIs("FILEGROUP") {
		fg, err := p.parseDatabaseFilegroup()
		if err != nil {
			return err
		}
		if stmt.Filegroups == nil {
			stmt.Filegroups = &nodes.List{}
		}
		stmt.Filegroups.Items = append(stmt.Filegroups.Items, fg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if !p.dbWordIs("FILEGROUP") {
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// dbFileSpecMode says which statement a <filespec> is in: SQL Server requires
// FILENAME everywhere but MODIFY FILE (Msg 1036), and takes NEWNAME only in
// CREATE DATABASE and MODIFY FILE (Msg 122).
type dbFileSpecMode int

const (
	dbFileCreate dbFileSpecMode = iota // CREATE DATABASE ON / LOG ON
	dbFileAdd                          // ALTER DATABASE ADD [LOG] FILE, REBUILD LOG ON
	dbFileModify                       // ALTER DATABASE MODIFY FILE
)

// parseDatabaseFileSpecList parses <filespec> [ , ...n ] where every item
// is a parenthesized file specification (TSql160.g:4620).
func (p *Parser) parseDatabaseFileSpecList(mode dbFileSpecMode) (*nodes.List, error) {
	var files []nodes.Node
	for {
		spec, err := p.parseDatabaseFileSpec(mode)
		if err != nil {
			return nil, err
		}
		files = append(files, spec)
		if p.cur.Type != ',' || p.peekNext().Type != '(' {
			break
		}
		p.advance() // consume ','
	}
	return &nodes.List{Items: files}, nil
}

// parseDatabaseFileSpec parses a single file specification (TSql160.g:4680,
// TSql160.g:4703). Each option is given at most once (Msg 153).
//
//	<filespec> ::=
//	(
//	    [ NAME = logical_file_name , ]
//	    FILENAME = { 'os_file_name' | 'filestream_path' }
//	    [ , NEWNAME = new_logical_name ]
//	    [ , SIZE = size [ KB | MB | GB | TB ] ]
//	    [ , MAXSIZE = { max_size [ KB | MB | GB | TB ] | UNLIMITED } ]
//	    [ , FILEGROWTH = growth_increment [ KB | MB | GB | TB | % ] ]
//	    [ , OFFLINE ]
//	)
func (p *Parser) parseDatabaseFileSpec(mode dbFileSpecMode) (*nodes.DatabaseFileSpec, error) {
	loc := p.pos()
	spec := &nodes.DatabaseFileSpec{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	hasFilename := false
	for {
		key, err := p.dbExpectWordIn("NAME", "FILENAME", "SIZE", "MAXSIZE", "FILEGROWTH", "NEWNAME", "OFFLINE")
		if err != nil {
			return nil, err
		}
		if seen[key] {
			return nil, p.dbErrorAtPrev()
		}
		seen[key] = true
		if key == "NEWNAME" && mode == dbFileAdd {
			return nil, p.dbErrorAtPrev()
		}
		if key != "OFFLINE" {
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
		}
		switch key {
		case "NAME", "NEWNAME":
			// nonEmptyStringOrIdentifier (TSql160.g:4727)
			v, err := p.dbExpectNameOrString()
			if err != nil {
				return nil, err
			}
			if key == "NAME" {
				spec.Name = v
			} else {
				spec.NewName = v
			}
		case "FILENAME":
			v, err := p.dbExpectString(true)
			if err != nil {
				return nil, err
			}
			spec.Filename = v
			hasFilename = true
		case "SIZE":
			if spec.Size, err = p.parseSizeValue(false); err != nil {
				return nil, err
			}
		case "MAXSIZE":
			if p.dbWordIs("UNLIMITED") {
				p.advance()
				spec.MaxSizeUnlimited = true
			} else if spec.MaxSize, err = p.parseSizeValue(false); err != nil {
				return nil, err
			}
		case "FILEGROWTH":
			if spec.FileGrowth, err = p.parseSizeValue(true); err != nil {
				return nil, err
			}
		case "OFFLINE":
			spec.Offline = true
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if mode != dbFileModify && !hasFilename {
		return nil, p.dbErrorAtPrev()
	}
	spec.Loc.End = p.prevEnd()
	return spec, nil
}

// parseSizeValue parses an unsigned integer size with an optional unit
// (TSql160.g:4768, TSql160.g:4819). Only FILEGROWTH takes % (Msg 153).
//
//	size_value ::= integer [ KB | MB | GB | TB | % ]
func (p *Parser) parseSizeValue(allowPercent bool) (*nodes.SizeValue, error) {
	sv := &nodes.SizeValue{Loc: nodes.Loc{Start: p.pos(), End: -1}}
	v, err := p.dbExpectInt()
	if err != nil {
		return nil, err
	}
	sv.Value = v
	if allowPercent && p.cur.Type == '%' {
		sv.Unit = "%"
		p.advance()
	} else if unit, ok := p.dbWordIn("KB", "MB", "GB", "TB"); ok {
		sv.Unit = unit
		p.advance()
	}
	sv.Loc.End = p.prevEnd()
	return sv, nil
}

// parseDatabaseFilegroup parses a FILEGROUP clause (TSql160.g:4644).
//
//	FILEGROUP filegroup_name
//	    [ CONTAINS { FILESTREAM | MEMORY_OPTIMIZED_DATA } ] [ DEFAULT ]
//	    <filespec> [ , ...n ]
func (p *Parser) parseDatabaseFilegroup() (*nodes.DatabaseFilegroup, error) {
	loc := p.pos()
	fg := &nodes.DatabaseFilegroup{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	p.advance() // consume FILEGROUP

	name, err := p.dbExpectName()
	if err != nil {
		return nil, err
	}
	fg.Name = name

	if p.dbWordIs("CONTAINS") {
		p.advance()
		w, err := p.dbExpectWordIn("FILESTREAM", "MEMORY_OPTIMIZED_DATA")
		if err != nil {
			return nil, err
		}
		fg.ContainsFilestream = w == "FILESTREAM"
		fg.ContainsMemoryOptimized = w == "MEMORY_OPTIMIZED_DATA"
	}

	if p.cur.Type == kwDEFAULT {
		fg.IsDefault = true
		p.advance()
	}

	files, err := p.parseDatabaseFileSpecList(dbFileCreate)
	if err != nil {
		return nil, err
	}
	fg.Files = files

	fg.Loc.End = p.prevEnd()
	return fg, nil
}

// parseDatabaseWithOptions parses the WITH options of CREATE DATABASE
// (TSql160.g:4930, createAlterDbOption). After FOR ATTACH the list also takes
// RESTRICTED_USER and one service broker option, and FILESTREAM takes only
// DIRECTORY_NAME.
func (p *Parser) parseDatabaseWithOptions(attach bool) (*nodes.List, error) {
	var opts []nodes.Node
	seen := map[string]bool{}
	broker := false
	for {
		opt, err := p.parseOneDatabaseOption(attach)
		if err != nil {
			return nil, err
		}
		switch opt.Name {
		case "ENABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS":
			// SQL Server takes one service broker option.
			if broker {
				return nil, p.dbErrorAtPrev()
			}
			broker = true
		case "DB_CHAINING", "TRUSTWORTHY", "FILESTREAM":
			// Given twice: Msg 102, or Msg 33413 for FILESTREAM.
			if seen[opt.Name] {
				return nil, p.dbErrorAtPrev()
			}
			seen[opt.Name] = true
		}
		opts = append(opts, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return &nodes.List{Items: opts}, nil
}

// parseOneDatabaseOption parses a single CREATE DATABASE WITH option.
func (p *Parser) parseOneDatabaseOption(attach bool) (*nodes.DatabaseOption, error) {
	optLoc := p.pos()
	words := []string{
		"FILESTREAM", "DEFAULT_FULLTEXT_LANGUAGE", "DEFAULT_LANGUAGE", "NESTED_TRIGGERS",
		"TRANSFORM_NOISE_WORDS", "TWO_DIGIT_YEAR_CUTOFF", "DB_CHAINING", "TRUSTWORTHY",
		"PERSISTENT_LOG_BUFFER", "LEDGER", "CATALOG_COLLATION",
	}
	if attach {
		words = append(words, "ENABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS", "RESTRICTED_USER")
	}
	key, err := p.dbExpectWordIn(words...)
	if err != nil {
		return nil, err
	}
	opt := &nodes.DatabaseOption{
		Name: key,
		Loc:  nodes.Loc{Start: optLoc, End: -1},
	}

	switch key {
	case "FILESTREAM":
		if err := p.parseDatabaseFilestreamOption(opt, attach); err != nil {
			return nil, err
		}
	case "ENABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS", "RESTRICTED_USER":
		// flag
	case "DB_CHAINING", "TRUSTWORTHY":
		// DB_CHAINING { ON | OFF }, with no '=' (TSql160.g:4951)
		if opt.Value, err = p.dbExpectOnOff(); err != nil {
			return nil, err
		}
	default:
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		switch key {
		case "NESTED_TRIGGERS", "TRANSFORM_NOISE_WORDS", "LEDGER":
			if opt.Value, err = p.dbExpectOnOff(); err != nil {
				return nil, err
			}
		case "TWO_DIGIT_YEAR_CUTOFF":
			if opt.Value, err = p.dbExpectInt(); err != nil {
				return nil, err
			}
		case "DEFAULT_FULLTEXT_LANGUAGE", "DEFAULT_LANGUAGE":
			if opt.Value, err = p.dbExpectLanguage(); err != nil {
				return nil, err
			}
			opt.Value = strings.ToUpper(opt.Value)
		case "CATALOG_COLLATION":
			if opt.Value, err = p.dbExpectName(); err != nil {
				return nil, err
			}
			opt.Value = strings.ToUpper(opt.Value)
		case "PERSISTENT_LOG_BUFFER":
			if opt.Value, err = p.dbExpectOnOff(); err != nil {
				return nil, err
			}
			if opt.Value == "ON" {
				dir, err := p.parseDatabasePersistentLogBufferDir()
				if err != nil {
					return nil, err
				}
				opt.PersistentLogDir = dir
			}
		}
	}

	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseDatabasePersistentLogBufferDir parses the directory of
// PERSISTENT_LOG_BUFFER = ON, which SQL Server requires.
//
//	( DIRECTORY_NAME = 'path' )
func (p *Parser) parseDatabasePersistentLogBufferDir() (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	if err := p.dbExpectWord("DIRECTORY_NAME"); err != nil {
		return "", err
	}
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	dir, err := p.dbExpectString(false)
	if err != nil {
		return "", err
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return dir, nil
}

// parseDatabaseFilestreamOption parses the parenthesized options after
// FILESTREAM (TSql160.g:5106, TSql160.g:5125). Each option is given at most
// once; after FOR ATTACH only DIRECTORY_NAME is taken.
//
//	FILESTREAM ( { NON_TRANSACTED_ACCESS = { OFF | READ_ONLY | FULL }
//	             | DIRECTORY_NAME = { 'directory_name' | NULL } } [ , ...n ] )
func (p *Parser) parseDatabaseFilestreamOption(opt *nodes.DatabaseOption, attach bool) error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	words := []string{"NON_TRANSACTED_ACCESS", "DIRECTORY_NAME"}
	if attach {
		words = words[1:]
	}
	seen := map[string]bool{}
	for {
		key, err := p.dbExpectWordIn(words...)
		if err != nil {
			return err
		}
		if seen[key] {
			return p.dbErrorAtPrev()
		}
		seen[key] = true
		if _, err := p.expect('='); err != nil {
			return err
		}
		if key == "NON_TRANSACTED_ACCESS" {
			v, err := p.dbExpectWordIn("OFF", "READ_ONLY", "FULL")
			if err != nil {
				return err
			}
			opt.FilestreamAccess = v
		} else if p.cur.Type == kwNULL {
			opt.FilestreamDirName = "NULL"
			p.advance()
		} else {
			v, err := p.dbExpectString(false)
			if err != nil {
				return err
			}
			opt.FilestreamDirName = v
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	_, err := p.expect(')')
	return err
}

// dbWordIs reports whether the current token is the word w (upper case), a
// keyword or an undelimited identifier. SQL Server reads its option words
// only undelimited: [MAXDOP] is refused.
func (p *Parser) dbWordIs(w string) bool {
	if kw := lookupKeyword(w); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == w
}

// dbWordIn returns which of the words the current token is.
func (p *Parser) dbWordIn(words ...string) (string, bool) {
	for _, w := range words {
		if p.dbWordIs(w) {
			return w, true
		}
	}
	return "", false
}

// dbExpectWord consumes the word w or fails at the current token.
func (p *Parser) dbExpectWord(w string) error {
	if !p.dbWordIs(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// dbExpectWordIn consumes one of the words and returns it.
func (p *Parser) dbExpectWordIn(words ...string) (string, error) {
	w, ok := p.dbWordIn(words...)
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return w, nil
}

// dbExpectOnOff consumes ON or OFF.
func (p *Parser) dbExpectOnOff() (string, error) {
	return p.dbExpectWordIn("ON", "OFF")
}

// dbIsInt reports whether the current token is an unsigned decimal integer.
// SQL Server refuses a binary constant (0x10), a sign, and a decimal point
// where these options take an integer.
func (p *Parser) dbIsInt() bool {
	return p.cur.Type == tokICONST && !strings.HasPrefix(p.cur.Str, "0x") && !strings.HasPrefix(p.cur.Str, "0X")
}

// dbExpectInt consumes an unsigned decimal integer.
func (p *Parser) dbExpectInt() (string, error) {
	if !p.dbIsInt() {
		return "", p.syntaxErrorAtCur()
	}
	v := p.cur.Str
	p.advance()
	return v, nil
}

// dbExpectString consumes a string literal; nonEmpty refuses ” (Msg 1038).
func (p *Parser) dbExpectString(nonEmpty bool) (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	if nonEmpty && p.cur.Str == "" {
		return "", p.syntaxErrorAtCur()
	}
	v := p.cur.Str
	p.advance()
	return v, nil
}

// dbExpectName consumes an identifier.
func (p *Parser) dbExpectName() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// dbExpectNameOrString consumes an identifier or a non-empty string.
func (p *Parser) dbExpectNameOrString() (string, error) {
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		return p.dbExpectString(true)
	}
	return p.dbExpectName()
}

// dbExpectLanguage consumes a language: an lcid or a language name
// (DEFAULT_LANGUAGE, DEFAULT_FULLTEXT_LANGUAGE). SQL Server refuses a string,
// ON, OFF, and reserved words there.
func (p *Parser) dbExpectLanguage() (string, error) {
	if p.dbIsInt() {
		return p.dbExpectInt()
	}
	if p.dbWordIs("ON") || p.dbWordIs("OFF") {
		return "", p.syntaxErrorAtCur()
	}
	return p.dbExpectName()
}

// dbErrorAtPrev reports a syntax error at the token just consumed, for a
// check SQL Server makes once the token is read (a repeated option, a
// missing FILENAME).
func (p *Parser) dbErrorAtPrev() error {
	return &ParseError{
		Message:  "syntax error at or near \"" + p.tokenText(p.prev) + "\"",
		Position: p.prev.Loc,
		Source:   p.source,
	}
}
