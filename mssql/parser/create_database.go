// Package parser - create_database.go implements T-SQL CREATE DATABASE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE statement.
//
// BNF: mssql/parser/bnf/create-database-transact-sql.bnf
//
//	CREATE DATABASE database_name
//	[ CONTAINMENT = { NONE | PARTIAL } ]
//	[ ON
//	      [ PRIMARY ] <filespec> [ , ...n ]
//	      [ , <filegroup> [ , ...n ] ]
//	      [ LOG ON <filespec> [ , ...n ] ]
//	]
//	[ COLLATE collation_name ]
//	[ WITH <option> [ , ...n ] ]
//	[ ; ]
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
//	    | PERSISTENT_LOG_BUFFER = ON ( DIRECTORY_NAME = 'path-to-directory-on-a-DAX-volume' )
//	    | LEDGER = { ON | OFF }
//	}
//
//	<filestream_option> ::=
//	{
//	      NON_TRANSACTED_ACCESS = { OFF | READ_ONLY | FULL }
//	    | DIRECTORY_NAME = 'directory_name'
//	}
//
//	<filespec> ::=
//	{
//	(
//	    NAME = logical_file_name ,
//	    FILENAME = { 'os_file_name' | 'filestream_path' }
//	    [ , SIZE = size [ KB | MB | GB | TB ] ]
//	    [ , MAXSIZE = { max_size [ KB | MB | GB | TB ] | UNLIMITED } ]
//	    [ , FILEGROWTH = growth_increment [ KB | MB | GB | TB | % ] ]
//	)
//	}
//
//	<filegroup> ::=
//	{
//	FILEGROUP filegroup_name [ [ CONTAINS FILESTREAM ] [ DEFAULT ] | CONTAINS MEMORY_OPTIMIZED_DATA ]
//	    <filespec> [ , ...n ]
//	}
//
//	-- Attach form:
//	CREATE DATABASE database_name
//	    ON <filespec> [ , ...n ]
//	    FOR { { ATTACH [ WITH <attach_database_option> [ , ...n ] ] }
//	        | ATTACH_REBUILD_LOG }
//	[ ; ]
//
//	<attach_database_option> ::=
//	{
//	      <service_broker_option>
//	    | RESTRICTED_USER
//	    | FILESTREAM ( DIRECTORY_NAME = { 'directory_name' | NULL } )
//	}
//
//	<service_broker_option> ::=
//	{
//	    ENABLE_BROKER
//	  | NEW_BROKER
//	  | ERROR_BROKER_CONVERSATIONS
//	}
//
//	-- Snapshot form:
//	CREATE DATABASE database_snapshot_name
//	    ON
//	    (
//	        NAME = logical_file_name ,
//	        FILENAME = 'os_file_name'
//	    ) [ , ...n ]
//	    AS SNAPSHOT OF source_database_name
//	[ ; ]
func (p *Parser) parseCreateDatabaseStmt() (*nodes.CreateDatabaseStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateDatabaseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Database name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// CONTAINMENT = { NONE | PARTIAL }
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTAINMENT") {
		p.advance() // consume CONTAINMENT
		p.match('=')
		if id, ok := p.parseIdentifier(); ok {
			stmt.Containment = strings.ToUpper(id)
		}
	}

	// ON clause (file specifications)
	if p.cur.Type == kwON {
		p.advance() // consume ON

		// [ PRIMARY ] <filespec> [ , ...n ]
		isPrimary := false
		if p.cur.Type == kwPRIMARY {
			p.advance()
			isPrimary = true
		}

		// Parse primary filegroup file specs
		if p.cur.Type == '(' {
			var files []nodes.Node
			files = append(files, p.parseDatabaseFileSpec())
			for p.cur.Type == ',' {
				// Peek ahead: if next is '(' it's another filespec, otherwise break
				next := p.peekNext()
				if next.Type == '(' || (p.isFileSpecStart(next)) {
					p.advance() // consume comma
					if p.cur.Type == '(' {
						files = append(files, p.parseDatabaseFileSpec())
					} else {
						break
					}
				} else {
					p.advance() // consume comma
					if p.cur.Type == '(' {
						files = append(files, p.parseDatabaseFileSpec())
					} else {
						break
					}
				}
			}
			if isPrimary || len(files) > 0 {
				stmt.OnPrimary = &nodes.List{Items: files}
			}
		}

		// [ , <filegroup> [ , ...n ] ]
		for p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILEGROUP") {
			fg := p.parseDatabaseFilegroup()
			if stmt.Filegroups == nil {
				stmt.Filegroups = &nodes.List{}
			}
			stmt.Filegroups.Items = append(stmt.Filegroups.Items, fg)
			// consume optional comma before next FILEGROUP
			if p.cur.Type == ',' {
				next := p.peekNext()
				if p.isIdentLikeToken(next) && matchesKeywordCI(next.Str, "FILEGROUP") {
					p.advance()
				}
			}
		}

		// LOG ON <filespec> [ , ...n ]
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LOG") {
			next := p.peekNext()
			if next.Type == kwON {
				p.advance() // consume LOG
				p.advance() // consume ON
				var logFiles []nodes.Node
				logFiles = append(logFiles, p.parseDatabaseFileSpec())
				for p.cur.Type == ',' {
					next := p.peekNext()
					if next.Type == '(' {
						p.advance() // consume comma
						logFiles = append(logFiles, p.parseDatabaseFileSpec())
					} else {
						break
					}
				}
				stmt.LogOn = &nodes.List{Items: logFiles}
			}
		}

		// FOR ATTACH / FOR ATTACH_REBUILD_LOG / AS SNAPSHOT OF
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ATTACH_REBUILD_LOG") {
				p.advance()
				stmt.ForAttachRebuildLog = true
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ATTACH") {
				p.advance()
				stmt.ForAttach = true
				// WITH <attach_database_option>
				if p.cur.Type == kwWITH {
					p.advance() // consume WITH
					stmt.AttachOptions = p.parseDatabaseAttachOptions()
				}
			}
		}

		if p.cur.Type == kwAS {
			p.advance() // consume AS
			// SNAPSHOT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SNAPSHOT") {
				p.advance() // consume SNAPSHOT
				// OF
				if p.cur.Type == kwOF {
					p.advance() // consume OF
					if id, ok := p.parseIdentifier(); ok {
						stmt.SnapshotOf = id
					}
				}
			}
		}
	}

	// COLLATE collation_name
	if p.cur.Type == kwCOLLATE {
		p.advance() // consume COLLATE
		if id, ok := p.parseIdentifier(); ok {
			stmt.Collation = id
		}
	}

	// WITH <option> [ , ...n ]
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.WithOptions = p.parseDatabaseWithOptions()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDatabaseFileSpec parses a single file specification.
//
//	<filespec> ::=
//	(
//	    NAME = logical_file_name ,
//	    FILENAME = { 'os_file_name' | 'filestream_path' }
//	    [ , SIZE = size [ KB | MB | GB | TB ] ]
//	    [ , MAXSIZE = { max_size [ KB | MB | GB | TB ] | UNLIMITED } ]
//	    [ , FILEGROWTH = growth_increment [ KB | MB | GB | TB | % ] ]
//	)
func (p *Parser) parseDatabaseFileSpec() *nodes.DatabaseFileSpec {
	loc := p.pos()
	spec := &nodes.DatabaseFileSpec{
		Loc: nodes.Loc{Start: loc},
	}

	p.match('(') // consume '('

	// Parse comma-separated key=value pairs inside parens
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isIdentLike() {
			key := strings.ToUpper(p.cur.Str)
			p.advance() // consume key
			switch key {
			case "NAME":
				p.match('=')
				if id, ok := p.parseIdentifier(); ok {
					spec.Name = id
				}
			case "FILENAME":
				p.match('=')
				if p.cur.Type == tokSCONST {
					spec.Filename = p.cur.Str
					p.advance()
				}
			case "SIZE":
				p.match('=')
				spec.Size = p.parseSizeValue()
			case "MAXSIZE":
				p.match('=')
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "UNLIMITED") {
					spec.MaxSizeUnlimited = true
					p.advance()
				} else {
					spec.MaxSize = p.parseSizeValue()
				}
			case "FILEGROWTH":
				p.match('=')
				spec.FileGrowth = p.parseSizeValue()
			case "NEWNAME":
				p.match('=')
				if id, ok := p.parseIdentifier(); ok {
					spec.NewName = id
				}
			case "OFFLINE":
				spec.Offline = true
			}
		}
		// Consume comma separator
		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	p.match(')') // consume ')'
	spec.Loc.End = p.prevEnd()
	return spec
}

// parseSizeValue parses a structured size value like "10MB", "100GB", "5%", or just a number.
//
//	size_value ::= number [ KB | MB | GB | TB | % ]
func (p *Parser) parseSizeValue() *nodes.SizeValue {
	sv := &nodes.SizeValue{Loc: nodes.Loc{Start: p.pos()}}

	// Read the numeric part
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		sv.Value = p.cur.Str
		p.advance()
	} else if p.isIdentLike() {
		// Could be a bare identifier like a number
		sv.Value = p.cur.Str
		p.advance()
		sv.Loc.End = p.prevEnd()
		return sv
	}

	// Read the optional unit suffix: KB, MB, GB, TB, or %
	if p.cur.Type == '%' {
		sv.Unit = "%"
		p.advance()
	} else if p.isIdentLike() {
		unit := strings.ToUpper(p.cur.Str)
		switch unit {
		case "KB", "MB", "GB", "TB":
			sv.Unit = unit
			p.advance()
		}
	}
	sv.Loc.End = p.prevEnd()
	return sv
}

// parseDatabaseFilegroup parses a FILEGROUP clause.
//
//	FILEGROUP filegroup_name [ [ CONTAINS FILESTREAM ] [ DEFAULT ] | CONTAINS MEMORY_OPTIMIZED_DATA ]
//	    <filespec> [ , ...n ]
func (p *Parser) parseDatabaseFilegroup() *nodes.DatabaseFilegroup {
	loc := p.pos()
	fg := &nodes.DatabaseFilegroup{
		Loc: nodes.Loc{Start: loc},
	}

	p.advance() // consume FILEGROUP keyword

	// filegroup_name
	if id, ok := p.parseIdentifier(); ok {
		fg.Name = id
	}

	// [ CONTAINS FILESTREAM ] [ DEFAULT ] | CONTAINS MEMORY_OPTIMIZED_DATA
	if p.cur.Type == kwCONTAINS {
		p.advance() // consume CONTAINS
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILESTREAM") {
			fg.ContainsFilestream = true
			p.advance()
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MEMORY_OPTIMIZED_DATA") {
			fg.ContainsMemoryOptimized = true
			p.advance()
		}
	}

	// [ DEFAULT ]
	if p.cur.Type == kwDEFAULT {
		fg.IsDefault = true
		p.advance()
	}

	// <filespec> [ , ...n ]
	if p.cur.Type == '(' {
		var files []nodes.Node
		files = append(files, p.parseDatabaseFileSpec())
		for p.cur.Type == ',' {
			next := p.peekNext()
			if next.Type == '(' {
				p.advance() // consume comma
				files = append(files, p.parseDatabaseFileSpec())
			} else {
				break
			}
		}
		fg.Files = &nodes.List{Items: files}
	}

	fg.Loc.End = p.prevEnd()
	return fg
}

// parseDatabaseWithOptions parses the WITH options for CREATE DATABASE.
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
//	    | PERSISTENT_LOG_BUFFER = ON ( DIRECTORY_NAME = 'path-to-directory-on-a-DAX-volume' )
//	    | LEDGER = { ON | OFF }
//	    | CATALOG_COLLATION = { DATABASE_DEFAULT | SQL_Latin1_General_CP1_CI_AS }
//	}
//
//	<filestream_option> ::=
//	{
//	      NON_TRANSACTED_ACCESS = { OFF | READ_ONLY | FULL }
//	    | DIRECTORY_NAME = 'directory_name'
//	}
func (p *Parser) parseDatabaseWithOptions() *nodes.List {
	var opts []nodes.Node

	for {
		if !p.isIdentLike() {
			break
		}

		opt := p.parseOneDatabaseOption()
		if opt != nil {
			opts = append(opts, opt)
		}

		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseOneDatabaseOption parses a single CREATE DATABASE WITH option.
func (p *Parser) parseOneDatabaseOption() *nodes.DatabaseOption {
	if !p.isIdentLike() {
		return nil
	}

	optLoc := p.pos()
	key := strings.ToUpper(p.cur.Str)

	// FILESTREAM ( <filestream_option> [ , ...n ] )
	if key == "FILESTREAM" {
		return p.parseDatabaseFilestreamOption()
	}

	p.advance() // consume key

	opt := &nodes.DatabaseOption{
		Name: key,
		Loc:  nodes.Loc{Start: optLoc},
	}

	// Parse value: = value or bare ON/OFF
	if p.cur.Type == '=' {
		p.advance() // consume =
		opt.Value = p.parseDatabaseOptionValue()
	} else if p.cur.Type == kwON || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ON")) {
		opt.Value = "ON"
		p.advance()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "OFF") {
		opt.Value = "OFF"
		p.advance()
	}

	// PERSISTENT_LOG_BUFFER = ON ( DIRECTORY_NAME = 'path' )
	if key == "PERSISTENT_LOG_BUFFER" && opt.Value == "ON" && p.cur.Type == '(' {
		p.advance() // consume (
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DIRECTORY_NAME") {
				p.advance() // consume DIRECTORY_NAME
				if p.cur.Type == '=' {
					p.advance()
					if p.cur.Type == tokSCONST {
						opt.PersistentLogDir = p.cur.Str
						p.advance()
					}
				}
			} else {
				p.advance()
			}
		}
		p.match(')')
	}

	opt.Loc.End = p.prevEnd()
	return opt
}

// parseDatabaseFilestreamOption parses FILESTREAM ( <filestream_option> [ , ...n ] ).
//
//	FILESTREAM (
//	    NON_TRANSACTED_ACCESS = { OFF | READ_ONLY | FULL },
//	    DIRECTORY_NAME = 'directory_name'
//	)
func (p *Parser) parseDatabaseFilestreamOption() *nodes.DatabaseOption {
	optLoc := p.pos()
	p.advance() // consume FILESTREAM

	opt := &nodes.DatabaseOption{
		Name: "FILESTREAM",
		Loc:  nodes.Loc{Start: optLoc},
	}

	if p.cur.Type == '(' {
		p.advance() // consume (
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() {
				subKey := strings.ToUpper(p.cur.Str)
				p.advance() // consume sub-key
				if p.cur.Type == '=' {
					p.advance() // consume =
					switch subKey {
					case "NON_TRANSACTED_ACCESS":
						if p.isIdentLike() || p.cur.Type == kwOFF {
							opt.FilestreamAccess = strings.ToUpper(p.cur.Str)
							p.advance()
						}
					case "DIRECTORY_NAME":
						if p.cur.Type == tokSCONST {
							opt.FilestreamDirName = p.cur.Str
							p.advance()
						} else if p.cur.Type == kwNULL {
							opt.FilestreamDirName = "NULL"
							p.advance()
						}
					default:
						// unknown sub-option, skip value
						if p.isIdentLike() || p.cur.Type == tokSCONST {
							p.advance()
						}
					}
				}
			}
			if p.cur.Type == ',' {
				p.advance()
			} else {
				break
			}
		}
		p.match(')')
	}

	opt.Loc.End = p.prevEnd()
	return opt
}

// parseDatabaseOptionValue parses a value for a database option (identifier, number, or string).
func (p *Parser) parseDatabaseOptionValue() string {
	if p.cur.Type == tokSCONST {
		val := p.cur.Str
		p.advance()
		return val
	}
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		return val
	}
	if p.cur.Type == kwON {
		p.advance()
		return "ON"
	}
	if p.isIdentLike() {
		val := strings.ToUpper(p.cur.Str)
		p.advance()
		return val
	}
	return ""
}

// parseDatabaseAttachOptions parses attach options after FOR ATTACH WITH.
//
//	<attach_database_option> ::=
//	{
//	      <service_broker_option>
//	    | RESTRICTED_USER
//	    | FILESTREAM ( DIRECTORY_NAME = { 'directory_name' | NULL } )
//	}
//
//	<service_broker_option> ::=
//	{
//	    ENABLE_BROKER
//	  | NEW_BROKER
//	  | ERROR_BROKER_CONVERSATIONS
//	}
func (p *Parser) parseDatabaseAttachOptions() *nodes.List {
	var opts []nodes.Node
	for {
		if !p.isIdentLike() {
			break
		}
		key := strings.ToUpper(p.cur.Str)
		switch key {
		case "ENABLE_BROKER", "NEW_BROKER", "ERROR_BROKER_CONVERSATIONS", "RESTRICTED_USER":
			optLoc := p.pos()
			p.advance()
			opts = append(opts, &nodes.DatabaseOption{
				Name: key,
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
		case "FILESTREAM":
			opts = append(opts, p.parseDatabaseFilestreamOption())
		default:
			// unknown option, advance to prevent infinite loop
			optLoc := p.pos()
			p.advance()
			opts = append(opts, &nodes.DatabaseOption{
				Name: key,
				Loc:  nodes.Loc{Start: optLoc, End: p.prevEnd()},
			})
		}
		if p.cur.Type == ',' {
			p.advance()
		} else {
			break
		}
	}
	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// isFileSpecStart checks if a token looks like the start of a filespec (an identifier or keyword).
func (p *Parser) isFileSpecStart(tok Token) bool {
	return tok.Type == '('
}

// isIdentLikeToken checks if a token can be used as an identifier, without consuming it.
func (p *Parser) isIdentLikeToken(tok Token) bool {
	if tok.Type == tokIDENT {
		return true
	}
	return tok.Type >= kwADD && tok.Str != ""
}
