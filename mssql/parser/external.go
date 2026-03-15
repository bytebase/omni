// Package parser - external.go implements T-SQL EXTERNAL object statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateExternalDataSourceStmt parses CREATE EXTERNAL DATA SOURCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-data-source-transact-sql
//
//	CREATE EXTERNAL DATA SOURCE <data_source_name>
//	WITH
//	  ( [ LOCATION = '<prefix>://<path>[:<port>]' ]
//	    [ [ , ] CONNECTION_OPTIONS = '<key_value_pairs>'[,...] ]
//	    [ [ , ] CREDENTIAL = <credential_name> ]
//	    [ [ , ] PUSHDOWN = { ON | OFF } ]
//	    [ [ , ] TYPE = { HADOOP | BLOB_STORAGE } ]
//	    [ [ , ] RESOURCE_MANAGER_LOCATION = '<resource_manager>[:<port>]' ]
//	  )
//	[ ; ]
func (p *Parser) parseCreateExternalDataSourceStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// EXTERNAL DATA SOURCE already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL DATA SOURCE",
		Loc:        nodes.Loc{Start: loc},
	}

	// data_source_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ( options )
	stmt.Options = p.parseExternalWithOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterExternalDataSourceStmt parses ALTER EXTERNAL DATA SOURCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-data-source-transact-sql
//
//	ALTER EXTERNAL DATA SOURCE data_source_name SET
//	    {
//	        LOCATION = '<prefix>://<path>[:<port>]' [,] |
//	        RESOURCE_MANAGER_LOCATION = '<IP address;Port>' [,] |
//	        CREDENTIAL = credential_name
//	    }
//	[ ; ]
func (p *Parser) parseAlterExternalDataSourceStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// EXTERNAL DATA SOURCE already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL DATA SOURCE",
		Loc:        nodes.Loc{Start: loc},
	}

	// data_source_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// SET key = value [, key = value ...]
	if p.cur.Type == kwSET {
		p.advance()
	}

	stmt.Options = p.parseExternalSetOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropExternalStmt parses DROP EXTERNAL { DATA SOURCE | TABLE | FILE FORMAT }.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-data-source-transact-sql
//
//	DROP EXTERNAL DATA SOURCE external_data_source_name
//	DROP EXTERNAL TABLE [ database_name . [ schema_name ] . | schema_name . ] table_name
//	DROP EXTERNAL FILE FORMAT external_file_format_name
func (p *Parser) parseDropExternalStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// EXTERNAL already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc},
	}

	// DATA SOURCE | TABLE | FILE FORMAT
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DATA") {
		p.advance() // consume DATA
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SOURCE") {
			p.advance() // consume SOURCE
		}
		stmt.ObjectType = "EXTERNAL DATA SOURCE"
	} else if p.cur.Type == kwTABLE {
		p.advance() // consume TABLE
		stmt.ObjectType = "EXTERNAL TABLE"
	} else if p.cur.Type == kwFILE || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILE")) {
		p.advance() // consume FILE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FORMAT") {
			p.advance() // consume FORMAT
		}
		stmt.ObjectType = "EXTERNAL FILE FORMAT"
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE") {
		p.advance() // consume RESOURCE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
			p.advance() // consume POOL
		}
		stmt.ObjectType = "EXTERNAL RESOURCE POOL"
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIBRARY") {
		p.advance() // consume LIBRARY
		return p.parseDropExternalLibraryStmt()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LANGUAGE") {
		p.advance() // consume LANGUAGE
		return p.parseDropExternalLanguageStmt()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MODEL") {
		p.advance() // consume MODEL
		return p.parseDropExternalModelStmt()
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STREAM") {
		p.advance() // consume STREAM
		stmt.ObjectType = "EXTERNAL STREAM"
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STREAMING") {
		p.advance() // consume STREAMING
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "JOB") {
			p.advance() // consume JOB
		}
		stmt.ObjectType = "EXTERNAL STREAMING JOB"
	}

	// name (possibly qualified)
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		var parts []string
		parts = append(parts, p.cur.Str)
		p.advance()
		for p.cur.Type == '.' {
			p.advance() // consume '.'
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				parts = append(parts, p.cur.Str)
				p.advance()
			}
		}
		stmt.Name = strings.Join(parts, ".")
	}

	// optional AUTHORIZATION owner_name (for DROP EXTERNAL LIBRARY)
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalTableStmt parses CREATE EXTERNAL TABLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-table-transact-sql
//
//	CREATE EXTERNAL TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    ( <column_definition> [ , ...n ] )
//	    WITH (
//	        LOCATION = 'folder_or_filepath' ,
//	        DATA_SOURCE = external_data_source_name
//	        [ , FILE_FORMAT = external_file_format_name ]
//	        [ , <reject_options> [ , ...n ] ]
//	        [ , TABLE_OPTIONS = N'...' ]
//	        [ , SCHEMA_NAME = N'...' ]
//	        [ , OBJECT_NAME = N'...' ]
//	        [ , DISTRIBUTION = { SHARDED(col) | REPLICATED | ROUND_ROBIN } ]
//	    )
//	[ ; ]
//
//	<column_definition> ::=
//	    column_name <data_type>
//	        [ COLLATE collation_name ]
//	        [ NULL | NOT NULL ]
//
//	<reject_options> ::=
//	{
//	    REJECT_TYPE = { value | percentage }
//	    | REJECT_VALUE = reject_value
//	    | REJECT_SAMPLE_VALUE = reject_sample_value
//	    | REJECTED_ROW_LOCATION = '/REJECT_Directory'
//	}
func (p *Parser) parseCreateExternalTableStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// EXTERNAL TABLE already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL TABLE",
		Loc:        nodes.Loc{Start: loc},
	}

	// table name (possibly qualified: db.schema.table)
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		var parts []string
		parts = append(parts, p.cur.Str)
		p.advance()
		for p.cur.Type == '.' {
			p.advance() // consume '.'
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				parts = append(parts, p.cur.Str)
				p.advance()
			}
		}
		stmt.Name = strings.Join(parts, ".")
	}

	var opts []nodes.Node

	// ( column_definition [ , ... ] )
	if p.cur.Type == '(' {
		p.advance() // consume '('
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			col := p.parseColumnDef()
			if col != nil {
				opts = append(opts, col)
			}
		}
		p.match(')') // consume ')'
	}

	// WITH ( options )
	withOpts := p.parseExternalWithOptions()
	if withOpts != nil {
		opts = append(opts, withOpts.Items...)
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalTableOrCETAS dispatches between CREATE EXTERNAL TABLE and CETAS.
// If the statement has AS <select_statement> after the WITH clause, it's CETAS.
func (p *Parser) parseCreateExternalTableOrCETAS(createLoc int) nodes.StmtNode {
	cetasStmt := p.parseCreateExternalTableAsSelectStmt()
	// If it parsed a query, it's CETAS
	if cetasStmt.Query != nil {
		cetasStmt.Loc.Start = createLoc
		return cetasStmt
	}
	// Otherwise, convert to regular SecurityStmt for backward compat
	secStmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL TABLE",
		Loc:        nodes.Loc{Start: createLoc, End: cetasStmt.Loc.End},
	}
	if cetasStmt.Name != nil {
		secStmt.Name = cetasStmt.Name.Object
		if cetasStmt.Name.Schema != "" {
			secStmt.Name = cetasStmt.Name.Schema + "." + secStmt.Name
		}
		if cetasStmt.Name.Database != "" {
			secStmt.Name = cetasStmt.Name.Database + "." + secStmt.Name
		}
	}
	if cetasStmt.Columns != nil || cetasStmt.Options != nil {
		var opts []nodes.Node
		if cetasStmt.Columns != nil {
			// Include column definitions directly (already ColumnDef nodes)
			opts = append(opts, cetasStmt.Columns.Items...)
		}
		if cetasStmt.Options != nil {
			opts = append(opts, cetasStmt.Options.Items...)
		}
		if len(opts) > 0 {
			secStmt.Options = &nodes.List{Items: opts}
		}
	}
	return secStmt
}

// parseCreateExternalFileFormatStmt parses CREATE EXTERNAL FILE FORMAT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-file-format-transact-sql
//
//	CREATE EXTERNAL FILE FORMAT file_format_name
//	WITH (
//	    FORMAT_TYPE = { DELIMITEDTEXT | RCFILE | ORC | PARQUET | JSON | DELTA }
//	    [ , FORMAT_OPTIONS ( <format_options> [ ,...n ] ) ]
//	    [ , SERDE_METHOD = 'SerDe_method' ]  -- for RCFILE
//	    [ , DATA_COMPRESSION = 'compression_codec' ]
//	)
//	[ ; ]
//
//	<format_options> ::=
//	{
//	    FIELD_TERMINATOR = field_terminator
//	    | STRING_DELIMITER = string_delimiter
//	    | FIRST_ROW = integer
//	    | DATE_FORMAT = datetime_format
//	    | USE_TYPE_DEFAULT = { TRUE | FALSE }
//	    | ENCODING = { 'UTF8' | 'UTF16' }
//	    | PARSER_VERSION = 'parser_version'
//	}
func (p *Parser) parseCreateExternalFileFormatStmt() *nodes.SecurityStmt {
	loc := p.pos()
	// EXTERNAL FILE FORMAT already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL FILE FORMAT",
		Loc:        nodes.Loc{Start: loc},
	}

	// file_format_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ( options )
	stmt.Options = p.parseExternalWithOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalLibraryStmt parses CREATE EXTERNAL LIBRARY.
// Caller has consumed CREATE EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-library-transact-sql
//
//	CREATE EXTERNAL LIBRARY library_name
//	[ AUTHORIZATION owner_name ]
//	FROM <file_spec> [ ,...2 ]
//	WITH ( LANGUAGE = <language> )
//	[ ; ]
//
//	<file_spec> ::=
//	{
//	    (CONTENT = { <client_library_specifier> | <library_bits> }
//	    [, PLATFORM = <platform> ])
//	}
//
//	<client_library_specifier> :: = { '[file_path\]manifest_file_name' }
//	<library_bits> :: = { varbinary_literal | varbinary_expression }
//	<platform> :: = { WINDOWS | LINUX }
//	<language> :: = { 'R' | 'Python' | <external_language> }
func (p *Parser) parseCreateExternalLibraryStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL LIBRARY",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// FROM (CONTENT = ...) or FROM literal
	var fileSpecOpts []nodes.Node
	if p.cur.Type == kwFROM {
		p.advance()
		if p.cur.Type == '(' {
			fileSpecOpts = append(fileSpecOpts, p.parseExternalFileSpec()...)
		} else if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			fileSpecOpts = append(fileSpecOpts, &nodes.String{Str: "CONTENT='" + p.cur.Str + "'"})
			p.advance()
		}
	}

	withOpts := p.parseExternalWithOptions()
	// Merge file spec options with WITH options
	var allOpts []nodes.Node
	allOpts = append(allOpts, fileSpecOpts...)
	if withOpts != nil {
		allOpts = append(allOpts, withOpts.Items...)
	}
	if len(allOpts) > 0 {
		stmt.Options = &nodes.List{Items: allOpts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterExternalLibraryStmt parses ALTER EXTERNAL LIBRARY.
// Caller has consumed ALTER EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-library-transact-sql
//
//	ALTER EXTERNAL LIBRARY library_name
//	[ AUTHORIZATION owner_name ]
//	SET <file_spec>
//	WITH ( LANGUAGE = <language> )
//	[ ; ]
//
//	<file_spec> ::=
//	{
//	    (CONTENT = { <client_library_specifier> | <library_bits> | NONE}
//	    [, PLATFORM = <platform> ])
//	}
//
//	<client_library_specifier> :: = { '[path\]manifest_file_name' | '<relative_path_in_external_data_source>' }
//	<library_bits> :: = { varbinary_literal | varbinary_expression }
//	<platform> :: = { WINDOWS | LINUX }
//	<language> :: = { 'R' | 'Python' | <external_language> }
func (p *Parser) parseAlterExternalLibraryStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL LIBRARY",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// SET (CONTENT = ...) or ADD/REMOVE
	var fileSpecOpts []nodes.Node
	if p.cur.Type == kwSET || (p.isIdentLike() && (matchesKeywordCI(p.cur.Str, "ADD") || matchesKeywordCI(p.cur.Str, "REMOVE"))) {
		p.advance()
		if p.cur.Type == '(' {
			fileSpecOpts = append(fileSpecOpts, p.parseExternalFileSpec()...)
		}
	}

	withOpts := p.parseExternalWithOptions()
	var allOpts []nodes.Node
	allOpts = append(allOpts, fileSpecOpts...)
	if withOpts != nil {
		allOpts = append(allOpts, withOpts.Items...)
	}
	if len(allOpts) > 0 {
		stmt.Options = &nodes.List{Items: allOpts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalLanguageStmt parses CREATE EXTERNAL LANGUAGE.
// Caller has consumed CREATE EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-language-transact-sql
//
//	CREATE EXTERNAL LANGUAGE language_name
//	[ AUTHORIZATION owner_name ]
//	FROM <file_spec> [ ,...2 ]
//	[ ; ]
//
//	<file_spec> ::=
//	{
//	    ( CONTENT = { <external_lang_specifier> | <content_bits> },
//	    FILE_NAME = <external_lang_file_name>
//	    [ , PLATFORM = <platform> ]
//	    [ , PARAMETERS = <external_lang_parameters> ]
//	    [ , ENVIRONMENT_VARIABLES = <external_lang_env_variables> ] )
//	}
//
//	<external_lang_specifier> :: = { '[file_path\]os_file_name' }
//	<content_bits> :: = { varbinary_literal | varbinary_expression }
//	<external_lang_file_name> :: = 'extension_file_name'
//	<platform> :: = { WINDOWS | LINUX }
func (p *Parser) parseCreateExternalLanguageStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL LANGUAGE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// FROM <file_spec> [ ,...2 ]
	var fileSpecOpts []nodes.Node
	if p.cur.Type == kwFROM {
		p.advance()
		if p.cur.Type == '(' {
			fileSpecOpts = append(fileSpecOpts, p.parseExternalFileSpec()...)
		}
		// handle additional file_spec separated by commas
		for p.cur.Type == ',' {
			p.advance() // consume ','
			if p.cur.Type == '(' {
				fileSpecOpts = append(fileSpecOpts, p.parseExternalFileSpec()...)
			}
		}
	}

	if len(fileSpecOpts) > 0 {
		stmt.Options = &nodes.List{Items: fileSpecOpts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterExternalLanguageStmt parses ALTER EXTERNAL LANGUAGE.
// Caller has consumed ALTER EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-language-transact-sql
//
//	ALTER EXTERNAL LANGUAGE language_name
//	[ AUTHORIZATION owner_name ]
//	{
//	    SET <file_spec>
//	    | ADD <file_spec>
//	    | REMOVE PLATFORM <platform>
//	}
//	[ ; ]
//
//	<file_spec> ::=
//	{
//	    ( CONTENT = { <external_lang_specifier> | <content_bits> },
//	    FILE_NAME = <external_lang_file_name>
//	    [ , PLATFORM = <platform> ]
//	    [ , PARAMETERS = <external_lang_parameters> ]
//	    [ , ENVIRONMENT_VARIABLES = <external_lang_env_variables> ] )
//	}
//
//	<platform> :: = { WINDOWS | LINUX }
func (p *Parser) parseAlterExternalLanguageStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL LANGUAGE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AUTHORIZATION owner_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// SET | ADD | REMOVE PLATFORM
	var fileSpecOpts []nodes.Node
	if p.cur.Type == kwSET || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ADD")) {
		p.advance()
		if p.cur.Type == '(' {
			fileSpecOpts = append(fileSpecOpts, p.parseExternalFileSpec()...)
		}
	} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "REMOVE") {
		p.advance() // consume REMOVE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PLATFORM") {
			p.advance() // consume PLATFORM
		}
		if p.isIdentLike() {
			fileSpecOpts = append(fileSpecOpts, &nodes.String{Str: "REMOVE_PLATFORM=" + strings.ToUpper(p.cur.Str)})
			p.advance() // consume platform name (WINDOWS/LINUX)
		}
	}

	if len(fileSpecOpts) > 0 {
		stmt.Options = &nodes.List{Items: fileSpecOpts}
	}
	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropExternalLibraryStmt parses DROP EXTERNAL LIBRARY.
// Caller has consumed DROP EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-library-transact-sql
//
//	DROP EXTERNAL LIBRARY library_name
//	[ AUTHORIZATION owner_name ]
//	[ ; ]
func (p *Parser) parseDropExternalLibraryStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EXTERNAL LIBRARY",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// optional AUTHORIZATION owner_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUTHORIZATION") {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropExternalLanguageStmt parses DROP EXTERNAL LANGUAGE.
// Caller has consumed DROP EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-language-transact-sql
//
//	DROP EXTERNAL LANGUAGE language_name
//	[ ; ]
func (p *Parser) parseDropExternalLanguageStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EXTERNAL LANGUAGE",
		Loc:        nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExternalSetOptions parses SET key = value [, key = value ...] (no parentheses).
// Used by ALTER EXTERNAL DATA SOURCE.
func (p *Parser) parseExternalSetOptions() *nodes.List {
	var opts []nodes.Node

	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		// skip commas
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
			optLoc := p.pos()
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			val := ""
			if p.cur.Type == '=' {
				p.advance() // consume '='
				val = p.parseExternalOptionValue()
			}
			opts = append(opts, &nodes.ExternalOption{
				Key:   key,
				Value: val,
				Loc:   nodes.Loc{Start: optLoc, End: p.pos()},
			})
		} else {
			break
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseExternalWithOptions parses a WITH ( key = value [, ...] ) clause,
// handling nested parentheses (e.g., FORMAT_OPTIONS(...), DISTRIBUTION = SHARDED(col)).
func (p *Parser) parseExternalWithOptions() *nodes.List {
	var opts []nodes.Node

	if p.cur.Type != kwWITH {
		return nil
	}
	p.advance() // consume WITH

	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// skip commas
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		// key = value
		if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance() // consume '='
				val := p.parseExternalOptionValue()
				opts = append(opts, &nodes.String{Str: key + "=" + val})
			} else if p.cur.Type == '(' {
				// FORMAT_OPTIONS ( key = value [, ...] )
				innerOpts := p.parseExternalFileSpec()
				for _, o := range innerOpts {
					if s, ok := o.(*nodes.String); ok {
						opts = append(opts, &nodes.String{Str: key + "." + s.Str})
					}
				}
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
		} else {
			// skip unexpected tokens
			p.advance()
		}
	}
	p.match(')') // consume ')'

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseExternalOptionValue reads one option value, which may be:
// - a string literal
// - a numeric literal
// - an identifier (e.g., ON, OFF, HADOOP, PARQUET)
// - SHARDED(column_name) or similar function-like syntax
func (p *Parser) parseExternalOptionValue() string {
	var val string
	switch {
	case p.cur.Type == tokSCONST:
		val = "'" + p.cur.Str + "'"
		p.advance()
	case p.cur.Type == tokICONST || p.cur.Type == tokFCONST:
		val = p.cur.Str
		p.advance()
	case p.cur.Type == kwON:
		val = "ON"
		p.advance()
	case p.cur.Type == kwOFF:
		val = "OFF"
		p.advance()
	case p.cur.Type == kwNULL:
		val = "NULL"
		p.advance()
	case p.isIdentLike():
		val = strings.ToUpper(p.cur.Str)
		p.advance()
		// check for function-like syntax: SHARDED(col)
		if p.cur.Type == '(' {
			p.advance() // consume '('
			var args []string
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				if p.isIdentLike() {
					args = append(args, strings.ToUpper(p.cur.Str))
				} else {
					args = append(args, p.cur.Str)
				}
				p.advance()
			}
			p.match(')') // consume ')'
			val += "(" + strings.Join(args, ", ") + ")"
		}
	default:
		val = p.cur.Str
		p.advance()
	}
	return val
}

// parseCreateExternalModelStmt parses CREATE EXTERNAL MODEL.
// Caller has consumed CREATE EXTERNAL MODEL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-model-transact-sql
//
//	CREATE EXTERNAL MODEL external_model_object_name
//	[ AUTHORIZATION owner_name ]
//	WITH
//	  ( LOCATION = '<prefix>://<path>[:<port>]'
//	    , API_FORMAT = '<OpenAI, Azure OpenAI, etc>'
//	    , MODEL_TYPE = EMBEDDINGS
//	    , MODEL = 'text-embedding-model-name'
//	    [ , CREDENTIAL = <credential_name> ]
//	    [ , PARAMETERS = '{"valid":"JSON"}' ]
//	    [ , LOCAL_RUNTIME_PATH = 'path to the ONNX Runtime files' ]
//	  )
//	[ ; ]
func (p *Parser) parseCreateExternalModelStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc},
	}

	// external_model_object_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// [ AUTHORIZATION owner_name ]
	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		if p.isIdentLike() {
			p.advance()
		}
	}

	// WITH ( options )
	stmt.Options = p.parseExternalWithOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterExternalModelStmt parses ALTER EXTERNAL MODEL.
// Caller has consumed ALTER EXTERNAL MODEL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-model-transact-sql
//
//	ALTER EXTERNAL MODEL external_model_object_name
//	SET
//	  (   LOCATION = '<prefix>://<path>[:<port>]'
//	    , API_FORMAT = '<OpenAI, Azure OpenAI, etc>'
//	    , MODEL_TYPE = EMBEDDINGS
//	    , MODEL = 'text-embedding-ada-002'
//	    [ , CREDENTIAL = <credential_name> ]
//	    [ , PARAMETERS = '{"valid":"JSON"}' ]
//	  )
//	[ ; ]
func (p *Parser) parseAlterExternalModelStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc},
	}

	// external_model_object_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// SET
	if p.cur.Type == kwSET {
		p.advance()
	}

	// ( options )
	if p.cur.Type == '(' {
		p.advance() // consume '('
		var opts []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF {
				key := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					val := p.parseExternalOptionValue()
					opts = append(opts, &nodes.String{Str: key + "=" + val})
				} else {
					opts = append(opts, &nodes.String{Str: key})
				}
			} else {
				p.advance()
			}
		}
		p.match(')') // consume ')'
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDropExternalModelStmt parses DROP EXTERNAL MODEL.
// Caller has consumed DROP EXTERNAL MODEL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-model-transact-sql
//
//	DROP EXTERNAL MODEL external_model_object_name
//	[ ; ]
func (p *Parser) parseDropExternalModelStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc},
	}

	// external_model_object_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalStreamStmt parses CREATE EXTERNAL STREAM.
// Caller has consumed CREATE EXTERNAL STREAM.
//
// BNF: mssql/parser/bnf/create-external-stream-transact-sql.bnf
//
//	CREATE EXTERNAL STREAM stream_name
//	    WITH (
//	        [ LOCATION = 'location' , ]
//	        [ INPUT_OPTIONS = 'json_string' , ]
//	        [ OUTPUT_OPTIONS = 'json_string' , ]
//	        [ DATA_SOURCE = data_source_name , ]
//	        [ FILE_FORMAT = file_format_name ]
//	    )
//	[ ; ]
func (p *Parser) parseCreateExternalStreamStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL STREAM",
		Loc:        nodes.Loc{Start: loc},
	}

	// stream_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ( options )
	stmt.Options = p.parseExternalWithOptions()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateExternalStreamingJobStmt parses CREATE EXTERNAL STREAMING JOB.
// Caller has consumed CREATE EXTERNAL STREAMING JOB.
//
//	CREATE EXTERNAL STREAMING JOB job_name
//	    WITH ( options )
//	    AS query_string
//	[ ; ]
func (p *Parser) parseCreateExternalStreamingJobStmt() *nodes.SecurityStmt {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL STREAMING JOB",
		Loc:        nodes.Loc{Start: loc},
	}

	// job_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// WITH ( options )
	stmt.Options = p.parseExternalWithOptions()

	// AS query_string (skip the rest)
	if p.cur.Type == kwAS {
		p.advance()
		// consume query string literal
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExternalFileSpec parses a parenthesized list of key = value options.
// Used for file specs: (CONTENT = '...', PLATFORM = WINDOWS, FILE_NAME = '...')
// and FORMAT_OPTIONS: (FIELD_TERMINATOR = ',', FIRST_ROW = 2, ...)
// Returns a list of nodes.String items in "KEY=VALUE" format.
func (p *Parser) parseExternalFileSpec() []nodes.Node {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // consume '('

	var opts []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		if p.isIdentLike() || p.cur.Type == kwON || p.cur.Type == kwOFF || p.cur.Type == kwNULL {
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance() // consume '='
				val := p.parseExternalOptionValue()
				opts = append(opts, &nodes.String{Str: key + "=" + val})
			} else {
				opts = append(opts, &nodes.String{Str: key})
			}
		} else {
			// skip unexpected tokens
			p.advance()
		}
	}
	p.match(')') // consume ')'
	return opts
}

