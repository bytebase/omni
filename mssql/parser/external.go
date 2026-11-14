// Package parser - external.go implements T-SQL EXTERNAL object statement parsing.
//
// Option lists here are read by their words, as SQL Server and ScriptDOM read
// them (mssql/testdata/TSql160.g), not by their shape: an option name comes
// from a closed set, its value has the type the grammar gives it, and an
// option appears once. Where the pinned SQL Server 2022 parses a statement,
// its verdict decides; it has no PolyBase, so it refuses CREATE EXTERNAL
// TABLE, CREATE EXTERNAL FILE FORMAT, CREATE EXTERNAL STREAM and COPY INTO
// outright, and those follow the grammar and the learn.microsoft.com syntax.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// ---------- shared helpers ----------

// extIsWord reports whether the current token is the unquoted word w, given
// in upper case: its keyword token, or a plain identifier with its text. SQL
// Server matches option words by their text and refuses a bracketed or
// double-quoted one, as TYPE = [BLOB_STORAGE].
func (p *Parser) extIsWord(w string) bool {
	if kw := lookupKeyword(w); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == w
}

// extWordIn returns the upper-case word when the current token is one of
// words, or "" when it is none of them.
func (p *Parser) extWordIn(words ...string) string {
	for _, w := range words {
		if p.extIsWord(w) {
			return w
		}
	}
	return ""
}

// extExpectWordIn consumes one of words and returns it.
func (p *Parser) extExpectWordIn(words ...string) (string, error) {
	w := p.extWordIn(words...)
	if w == "" {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return w, nil
}

// extName reads a one-part object name.
func (p *Parser) extName() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// extString reads a string literal ('...' or N'...') and returns it quoted.
func (p *Parser) extString() (string, error) {
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	s := "'" + p.cur.Str + "'"
	p.advance()
	return s, nil
}

// extNonEmptyString reads a string literal that may not be empty
// (nonEmptyString, TSql160.g:33449; SQL Server's Msg 1038).
func (p *Parser) extNonEmptyString() (string, error) {
	if (p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST) && p.cur.Str == "" {
		return "", p.syntaxErrorAtCur()
	}
	return p.extString()
}

// extIsInteger reports whether the current token is a decimal integer
// (the grammar's Integer, not a 0x binary literal).
func (p *Parser) extIsInteger() bool {
	return p.cur.Type == tokICONST && !strings.HasPrefix(strings.ToLower(p.cur.Str), "0x")
}

// extInteger reads a decimal integer.
func (p *Parser) extInteger() (string, error) {
	if !p.extIsInteger() {
		return "", p.syntaxErrorAtCur()
	}
	s := p.cur.Str
	p.advance()
	return s, nil
}

// extIntegerOrNumeric reads an integer or a decimal number without an
// exponent (integerOrNumeric, TSql160.g:33376).
func (p *Parser) extIntegerOrNumeric() (string, error) {
	if p.cur.Type == tokFCONST && !strings.ContainsAny(p.cur.Str, "eExX") {
		s := p.cur.Str
		p.advance()
		return s, nil
	}
	return p.extInteger()
}

// extIdentifier reads an identifier value.
func (p *Parser) extIdentifier() (string, error) {
	return p.extName()
}

// extIdentifierOrString reads an identifier or a string literal.
func (p *Parser) extIdentifierOrString() (string, error) {
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		return p.extString()
	}
	return p.extIdentifier()
}

// extOnOff reads ON or OFF.
func (p *Parser) extOnOff() (string, error) {
	switch p.cur.Type {
	case kwON:
		p.advance()
		return "ON", nil
	case kwOFF:
		p.advance()
		return "OFF", nil
	}
	return "", p.syntaxErrorAtCur()
}

// extAuthorization reads an optional AUTHORIZATION owner_name.
func (p *Parser) extAuthorization() error {
	if p.cur.Type != kwAUTHORIZATION {
		return nil
	}
	p.advance()
	_, err := p.extName()
	return err
}

// extOptionHead reads an option name from names and the = after it, failing
// on a name not in names or one already in seen.
func (p *Parser) extOptionHead(seen map[string]bool, names ...string) (string, int, error) {
	key := p.extWordIn(names...)
	if key == "" || seen[key] {
		return "", 0, p.syntaxErrorAtCur()
	}
	seen[key] = true
	loc := p.pos()
	p.advance()
	if _, err := p.expect('='); err != nil {
		return "", 0, err
	}
	return key, loc, nil
}

// extExpectListEnd fails unless the current token ends an option or a list
// item: a comma or a closing parenthesis.
func (p *Parser) extExpectListEnd() error {
	if p.cur.Type != ',' && p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// extTableRefName renders a qualified name as its parts joined by dots.
func extTableRefName(ref *nodes.TableRef) string {
	var parts []string
	if ref.Server != "" {
		parts = append(parts, ref.Server)
	}
	if ref.Server != "" || ref.Database != "" {
		parts = append(parts, ref.Database)
	}
	if ref.Server != "" || ref.Database != "" || ref.Schema != "" {
		parts = append(parts, ref.Schema)
	}
	parts = append(parts, ref.Object)
	return strings.Join(parts, ".")
}

// extRequireTableRef reads a qualified object name.
func (p *Parser) extRequireTableRef() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return ref, nil
}

// ---------- EXTERNAL DATA SOURCE ----------

// extDataSourceCreateOptions are the options SQL Server 2022 takes in CREATE
// EXTERNAL DATA SOURCE ... WITH. ScriptDOM also takes SHARD_MAP_NAME and
// TYPE = HADOOP | RDBMS | SHARD_MAP_MANAGER (TSql160.g:23722), for Azure SQL
// and Hadoop, which SQL Server 2022 refuses (Msg 102).
var extDataSourceCreateOptions = []string{
	"TYPE", "LOCATION", "PUSHDOWN", "CREDENTIAL", "CONNECTION_OPTIONS",
	"DATABASE_NAME", "RESOURCE_MANAGER_LOCATION",
}

// extDataSourceAlterOptions are the options SQL Server 2022 takes in ALTER
// EXTERNAL DATA SOURCE ... SET; it refuses RESOURCE_MANAGER_LOCATION and
// DATABASE_NAME there (Msg 102).
var extDataSourceAlterOptions = []string{
	"LOCATION", "PUSHDOWN", "CREDENTIAL", "CONNECTION_OPTIONS",
}

// extDataSourceOptionValue reads the value of a data source option
// (TSql160.g:23722-23830): TYPE = BLOB_STORAGE, PUSHDOWN = ON | OFF,
// CREDENTIAL = identifier, and a string for the rest, which SQL Server
// refuses empty but for CONNECTION_OPTIONS (Msg 46504).
func (p *Parser) extDataSourceOptionValue(key string) (string, error) {
	switch key {
	case "TYPE":
		return p.extExpectWordIn("BLOB_STORAGE")
	case "PUSHDOWN":
		return p.extOnOff()
	case "CREDENTIAL":
		return p.extIdentifier()
	case "CONNECTION_OPTIONS":
		return p.extString()
	}
	return p.extNonEmptyString()
}

// parseCreateExternalDataSourceStmt parses CREATE EXTERNAL DATA SOURCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-data-source-transact-sql
// Grammar: TSql160.g:23660 (createExternalDataSourceStatement)
//
//	CREATE EXTERNAL DATA SOURCE <data_source_name>
//	WITH
//	  ( [ LOCATION = '<prefix>://<path>[:<port>]' ]
//	    [ [ , ] CONNECTION_OPTIONS = '<key_value_pairs>'[,...] ]
//	    [ [ , ] CREDENTIAL = <credential_name> ]
//	    [ [ , ] PUSHDOWN = { ON | OFF } ]
//	    [ [ , ] TYPE = BLOB_STORAGE ]
//	    [ [ , ] DATABASE_NAME = '<database_name>' ]
//	    [ [ , ] RESOURCE_MANAGER_LOCATION = '<resource_manager>[:<port>]' ]
//	  )
//
// SQL Server 2022 takes an empty list and stray commas in it, WITH () and
// WITH (, LOCATION = 'a',), so this does too.
func (p *Parser) parseCreateExternalDataSourceStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL DATA SOURCE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for p.cur.Type != ')' {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		key, optLoc, err := p.extOptionHead(seen, extDataSourceCreateOptions...)
		if err != nil {
			return nil, err
		}
		val, err := p.extDataSourceOptionValue(key)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ExternalOption{Key: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if err := p.extExpectListEnd(); err != nil {
			return nil, err
		}
	}
	p.advance() // consume ')'
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterExternalDataSourceStmt parses ALTER EXTERNAL DATA SOURCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-data-source-transact-sql
// Grammar: TSql160.g:23830 (alterExternalDataSourceStatement)
//
//	ALTER EXTERNAL DATA SOURCE data_source_name SET
//	    { LOCATION = '<prefix>://<path>[:<port>]'
//	    | CONNECTION_OPTIONS = '<key_value_pairs>'
//	    | CREDENTIAL = credential_name
//	    | PUSHDOWN = { ON | OFF } } [ ,...n ]
func (p *Parser) parseAlterExternalDataSourceStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL DATA SOURCE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		key, optLoc, err := p.extOptionHead(seen, extDataSourceAlterOptions...)
		if err != nil {
			return nil, err
		}
		val, err := p.extDataSourceOptionValue(key)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ExternalOption{Key: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- DROP EXTERNAL ----------

// parseDropExternalStmt parses DROP EXTERNAL { DATA SOURCE | TABLE | FILE
// FORMAT | RESOURCE POOL | LIBRARY | LANGUAGE | MODEL | STREAM | STREAMING JOB }.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-data-source-transact-sql
// Grammar: TSql160.g:15149, 23877, 24256, 22314 (dropExternalStatement)
//
//	DROP EXTERNAL DATA SOURCE external_data_source_name
//	DROP EXTERNAL TABLE [ database_name . [ schema_name ] . | schema_name . ] table_name
//	DROP EXTERNAL FILE FORMAT external_file_format_name
//	DROP EXTERNAL RESOURCE POOL pool_name
func (p *Parser) parseDropExternalStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// EXTERNAL already consumed by caller
	stmt := &nodes.SecurityStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	switch p.cur.Type {
	case kwDATA:
		p.advance()
		if _, err := p.expect(kwSOURCE); err != nil {
			return nil, err
		}
		stmt.ObjectType = "EXTERNAL DATA SOURCE"
	case kwTABLE:
		p.advance()
		stmt.ObjectType = "EXTERNAL TABLE"
		ref, err := p.extRequireTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Name = extTableRefName(ref)
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	case kwFILE:
		p.advance()
		if _, err := p.expect(kwFORMAT); err != nil {
			return nil, err
		}
		stmt.ObjectType = "EXTERNAL FILE FORMAT"
	case kwRESOURCE:
		p.advance()
		if _, err := p.expect(kwPOOL); err != nil {
			return nil, err
		}
		stmt.ObjectType = "EXTERNAL RESOURCE POOL"
	case kwLIBRARY:
		p.advance()
		return p.parseDropExternalLibraryStmt()
	case kwLANGUAGE:
		p.advance()
		return p.parseDropExternalLanguageStmt()
	case kwMODEL:
		p.advance()
		return p.parseDropExternalModelStmt()
	case kwSTREAM:
		p.advance()
		stmt.ObjectType = "EXTERNAL STREAM"
	case kwSTREAMING:
		p.advance()
		if _, err := p.expect(kwJOB); err != nil {
			return nil, err
		}
		stmt.ObjectType = "EXTERNAL STREAMING JOB"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- EXTERNAL TABLE ----------

// extTableOptions are the options of CREATE EXTERNAL TABLE ... WITH
// (TSql160.g:27679-27880, ExternalTableOptionKind). ScriptDOM also takes
// DISTRIBUTION = { SHARDED ( column ) | REPLICATED | ROUND_ROBIN }, which
// only Azure SQL Database's elastic query reads, over an RDBMS or
// SHARD_MAP_MANAGER data source that SQL Server 2022 refuses as well.
var extTableOptions = []string{
	"DATA_SOURCE", "FILE_FORMAT", "LOCATION", "SCHEMA_NAME", "OBJECT_NAME",
	"REJECTED_ROW_LOCATION", "TABLE_OPTIONS", "REJECT_TYPE", "REJECT_VALUE",
	"REJECT_SAMPLE_VALUE",
}

// extTableWithOptions reads WITH ( option [ ,...n ] ) of CREATE EXTERNAL TABLE.
//
//	DATA_SOURCE = name | FILE_FORMAT = name
//	| { LOCATION | SCHEMA_NAME | OBJECT_NAME | REJECTED_ROW_LOCATION | TABLE_OPTIONS } = 'string'
//	| REJECT_TYPE = { VALUE | PERCENTAGE }
//	| { REJECT_VALUE | REJECT_SAMPLE_VALUE } = integer_or_numeric
func (p *Parser) extTableWithOptions() ([]nodes.Node, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		key, _, err := p.extOptionHead(seen, extTableOptions...)
		if err != nil {
			return nil, err
		}
		var val string
		switch key {
		case "DATA_SOURCE", "FILE_FORMAT":
			val, err = p.extIdentifier()
		case "REJECT_TYPE":
			val, err = p.extExpectWordIn("VALUE", "PERCENTAGE")
		case "REJECT_VALUE", "REJECT_SAMPLE_VALUE":
			val, err = p.extIntegerOrNumeric()
		default:
			val, err = p.extString()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: key + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// extTableColumnDef reads a column of CREATE EXTERNAL TABLE
// (columnDefinitionCreateExternalTable, TSql160.g:28953):
//
//	column_name data_type [ COLLATE collation_name ] [ NULL | NOT NULL ]
func (p *Parser) extTableColumnDef() (*nodes.ColumnDef, error) {
	loc := p.pos()
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	col := &nodes.ColumnDef{Name: name, Loc: nodes.Loc{Start: loc, End: -1}}
	if col.DataType, err = p.parseDataType(); err != nil {
		return nil, err
	}
	if col.DataType == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwCOLLATE {
		p.advance()
		if col.Collation, err = p.parseCollationName(); err != nil {
			return nil, err
		}
	}
	switch {
	case p.cur.Type == kwNULL:
		nullLoc := p.pos()
		p.advance()
		col.Nullable = &nodes.NullableSpec{Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
	case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
		nullLoc := p.pos()
		p.advance()
		p.advance()
		col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
	}
	col.Loc.End = p.prevEnd()
	return col, nil
}

// parseCreateExternalTableOrCETAS parses CREATE EXTERNAL TABLE and CREATE
// EXTERNAL TABLE AS SELECT (CETAS). Caller has consumed CREATE EXTERNAL TABLE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-table-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-table-as-select-transact-sql
// Grammar: TSql160.g:27638 (createExternalTableStatement)
//
//	CREATE EXTERNAL TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    ( <column_definition> [ ,...n ] )
//	    WITH ( <option> [ ,...n ] )
//
//	CREATE EXTERNAL TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ ( column_name [ ,...n ] ) ]
//	    WITH ( <option> [ ,...n ] )
//	    AS <select_statement>
//
// A CETAS becomes a CreateExternalTableAsSelectStmt; the other form a
// SecurityStmt whose Options hold the column definitions, then the options.
func (p *Parser) parseCreateExternalTableOrCETAS(createLoc int) (nodes.StmtNode, error) {
	ref, err := p.extRequireTableRef()
	if err != nil {
		return nil, err
	}

	var colDefs []nodes.Node  // ( column_definition, ... )
	var colNames []nodes.Node // ( column_name, ... ) of a CETAS
	if p.cur.Type == '(' {
		p.advance()
		next := p.peekNext().Type
		namesOnly := next == ',' || next == ')'
		for {
			if namesOnly {
				nameLoc := p.pos()
				name, err := p.extName()
				if err != nil {
					return nil, err
				}
				colNames = append(colNames, &nodes.ColumnDef{Name: name, Loc: nodes.Loc{Start: nameLoc, End: p.prevEnd()}})
			} else {
				col, err := p.extTableColumnDef()
				if err != nil {
					return nil, err
				}
				colDefs = append(colDefs, col)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	opts, err := p.extTableWithOptions()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == kwAS && colDefs == nil {
		p.advance()
		if p.cur.Type != kwSELECT && p.cur.Type != kwWITH {
			return nil, p.syntaxErrorAtCur()
		}
		query, err := p.parseSelectStmtWithCTE()
		if err != nil {
			return nil, err
		}
		if query == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt := &nodes.CreateExternalTableAsSelectStmt{
			Name:    ref,
			Options: &nodes.List{Items: opts},
			Query:   query,
			Loc:     nodes.Loc{Start: createLoc, End: p.prevEnd()},
		}
		if colNames != nil {
			stmt.Columns = &nodes.List{Items: colNames}
		}
		return stmt, nil
	}
	// Without AS SELECT the table needs its column definitions.
	if colDefs == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL TABLE",
		Name:       extTableRefName(ref),
		Options:    &nodes.List{Items: append(colDefs, opts...)},
		Loc:        nodes.Loc{Start: createLoc, End: p.prevEnd()},
	}, nil
}

// ---------- EXTERNAL FILE FORMAT ----------

// extFormatOptions are the FORMAT_OPTIONS words (TSql160.g:24161-24230,
// ExternalFileFormatOptionKind).
var extFormatOptions = []string{
	"FIELD_TERMINATOR", "STRING_DELIMITER", "DATE_FORMAT", "ENCODING",
	"PARSER_VERSION", "FIRST_ROW", "USE_TYPE_DEFAULT",
}

// parseCreateExternalFileFormatStmt parses CREATE EXTERNAL FILE FORMAT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-file-format-transact-sql
// Grammar: TSql160.g:24013 (createExternalFileFormatStatement)
//
//	CREATE EXTERNAL FILE FORMAT file_format_name
//	WITH (
//	    FORMAT_TYPE = { DELIMITEDTEXT | RCFILE | ORC | PARQUET | JSON | DELTA }
//	    [ , FORMAT_OPTIONS ( <format_options> [ ,...n ] ) ]
//	    [ , SERDE_METHOD = 'SerDe_method' ]
//	    [ , DATA_COMPRESSION = 'compression_codec' ]
//	)
//
//	<format_options> ::=
//	{
//	    FIELD_TERMINATOR = 'string' | STRING_DELIMITER = 'string'
//	    | DATE_FORMAT = 'string' | ENCODING = 'string' | PARSER_VERSION = 'string'
//	    | FIRST_ROW = integer
//	    | USE_TYPE_DEFAULT = { TRUE | FALSE }
//	}
func (p *Parser) parseCreateExternalFileFormatStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL FILE FORMAT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.extExpectWordIn("FORMAT_TYPE"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	formatType, err := p.extExpectWordIn("DELIMITEDTEXT", "RCFILE", "ORC", "PARQUET", "JSON", "DELTA")
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{&nodes.String{Str: "FORMAT_TYPE=" + formatType}}
	for p.cur.Type == ',' {
		p.advance()
		key := p.extWordIn("FORMAT_OPTIONS", "SERDE_METHOD", "DATA_COMPRESSION")
		if key == "" {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if key != "FORMAT_OPTIONS" {
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			val, err := p.extString()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: key + "=" + val})
			continue
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for {
			fkey := p.extWordIn(extFormatOptions...)
			if fkey == "" {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			var val string
			switch fkey {
			case "FIRST_ROW":
				val, err = p.extInteger()
			case "USE_TYPE_DEFAULT":
				val, err = p.extExpectWordIn("TRUE", "FALSE")
			default:
				val, err = p.extString()
			}
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.String{Str: "FORMAT_OPTIONS." + fkey + "=" + val})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- EXTERNAL LIBRARY and EXTERNAL LANGUAGE ----------

// extFileSpec reads one file specification of CREATE/ALTER EXTERNAL LIBRARY
// or LANGUAGE (TSql160.g:1149, 1168, 1210):
//
//	( CONTENT = { <expression> | NONE } [ , PLATFORM = { WINDOWS | LINUX } ]
//	  [ , FILE_NAME = 'name' ] [ , PARAMETERS = 'string' ]
//	  [ , ENVIRONMENT_VARIABLES = 'string' ] )
//
// FILE_NAME, PARAMETERS and ENVIRONMENT_VARIABLES are a language's only. SQL
// Server takes the options in any order, each once (Msg 1039), and CONTENT
// as any expression, which ScriptDOM limits to a binary or a string. required
// names the options a CREATE must give: CONTENT for a library, CONTENT and
// FILE_NAME for a language (Msg 39123).
func (p *Parser) extFileSpec(language bool, required ...string) ([]nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	names := []string{"CONTENT", "PLATFORM"}
	if language {
		names = append(names, "FILE_NAME", "PARAMETERS", "ENVIRONMENT_VARIABLES")
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		key, _, err := p.extOptionHead(seen, names...)
		if err != nil {
			return nil, err
		}
		var val string
		switch key {
		case "CONTENT":
			if p.extIsWord("NONE") {
				val = "NONE"
				p.advance()
				break
			}
			start := p.pos()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.syntaxErrorAtCur()
			}
			val = p.source[start:p.prevEnd()]
		case "PLATFORM":
			val, err = p.extExpectWordIn("WINDOWS", "LINUX")
		case "FILE_NAME":
			val, err = p.extNonEmptyString()
		default:
			val, err = p.extString()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: key + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	for _, r := range required {
		if !seen[r] {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// extFileSpecs reads the one or two file specifications after FROM in CREATE
// EXTERNAL LIBRARY or LANGUAGE. SQL Server takes one per platform (Msg 39112,
// 39114), so two that name the same PLATFORM, or neither names one, are
// refused. One without PLATFORM beside one with it depends on the server's
// own platform, so it is taken.
func (p *Parser) extFileSpecs(language bool, required ...string) ([]nodes.Node, error) {
	var opts []nodes.Node
	var platforms []string
	for {
		start := p.pos()
		spec, err := p.extFileSpec(language, required...)
		if err != nil {
			return nil, err
		}
		platform := ""
		for _, o := range spec {
			if s := o.(*nodes.String).Str; strings.HasPrefix(s, "PLATFORM=") {
				platform = s
			}
		}
		if len(platforms) == 1 && platforms[0] == platform {
			return nil, &ParseError{Message: "duplicate file specification for one platform", Position: start, Source: p.source}
		}
		platforms = append(platforms, platform)
		opts = append(opts, spec...)
		if p.cur.Type != ',' || len(platforms) == 2 {
			return opts, nil
		}
		p.advance()
	}
}

// extLibraryLanguage reads WITH ( LANGUAGE = 'language' ) of an external
// library (TSql160.g:1129).
func (p *Parser) extLibraryLanguage() (nodes.Node, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLANGUAGE); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	val, err := p.extNonEmptyString()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.String{Str: "LANGUAGE=" + val}, nil
}

// parseCreateExternalLibraryStmt parses CREATE EXTERNAL LIBRARY.
// Caller has consumed CREATE EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-library-transact-sql
// Grammar: TSql160.g:1080 (createExternalLibraryStatement)
//
//	CREATE EXTERNAL LIBRARY library_name
//	[ AUTHORIZATION owner_name ]
//	FROM <file_spec> [ ,...2 ]
//	WITH ( LANGUAGE = <language> )
func (p *Parser) parseCreateExternalLibraryStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL LIBRARY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.extAuthorization(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	opts, err := p.extFileSpecs(false, "CONTENT")
	if err != nil {
		return nil, err
	}
	lang, err := p.extLibraryLanguage()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: append(opts, lang)}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterExternalLibraryStmt parses ALTER EXTERNAL LIBRARY.
// Caller has consumed ALTER EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-library-transact-sql
// Grammar: TSql160.g:1159 (alterExternalLibraryStatement)
//
//	ALTER EXTERNAL LIBRARY library_name
//	[ AUTHORIZATION owner_name ]
//	{ SET | ADD } <file_spec>
//	WITH ( LANGUAGE = <language> )
//
// ScriptDOM takes SET only; SQL Server 2022 also takes ADD.
func (p *Parser) parseAlterExternalLibraryStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL LIBRARY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.extAuthorization(); err != nil {
		return nil, err
	}
	if p.cur.Type != kwSET && p.cur.Type != kwADD {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	opts, err := p.extFileSpec(false)
	if err != nil {
		return nil, err
	}
	lang, err := p.extLibraryLanguage()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: append(opts, lang)}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateExternalLanguageStmt parses CREATE EXTERNAL LANGUAGE.
// Caller has consumed CREATE EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-language-transact-sql
// Grammar: TSql160.g:1189 (createExternalLanguageStatement)
//
//	CREATE EXTERNAL LANGUAGE language_name
//	[ AUTHORIZATION owner_name ]
//	FROM <file_spec> [ ,...2 ]
func (p *Parser) parseCreateExternalLanguageStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL LANGUAGE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.extAuthorization(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	opts, err := p.extFileSpecs(true, "CONTENT", "FILE_NAME")
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterExternalLanguageStmt parses ALTER EXTERNAL LANGUAGE.
// Caller has consumed ALTER EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-external-language-transact-sql
// Grammar: TSql160.g:1268 (alterExternalLanguageStatement)
//
//	ALTER EXTERNAL LANGUAGE language_name
//	[ AUTHORIZATION owner_name ]
//	{ SET <file_spec> | ADD <file_spec> | REMOVE PLATFORM { WINDOWS | LINUX } }
func (p *Parser) parseAlterExternalLanguageStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL LANGUAGE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.extAuthorization(); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	switch p.cur.Type {
	case kwSET, kwADD:
		p.advance()
		if opts, err = p.extFileSpec(true); err != nil {
			return nil, err
		}
	case kwREMOVE:
		p.advance()
		if _, err := p.expect(kwPLATFORM); err != nil {
			return nil, err
		}
		platform, err := p.extExpectWordIn("WINDOWS", "LINUX")
		if err != nil {
			return nil, err
		}
		opts = []nodes.Node{&nodes.String{Str: "REMOVE_PLATFORM=" + platform}}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropExternalLibraryStmt parses DROP EXTERNAL LIBRARY.
// Caller has consumed DROP EXTERNAL LIBRARY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-library-transact-sql
// Grammar: TSql160.g:1129 (dropExternalLibraryStatement)
//
//	DROP EXTERNAL LIBRARY library_name [ AUTHORIZATION owner_name ]
func (p *Parser) parseDropExternalLibraryStmt() (*nodes.SecurityStmt, error) {
	return p.extDropWithAuthorization("EXTERNAL LIBRARY")
}

// parseDropExternalLanguageStmt parses DROP EXTERNAL LANGUAGE.
// Caller has consumed DROP EXTERNAL LANGUAGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-language-transact-sql
// Grammar: TSql160.g:1311 (dropExternalLanguageStatement)
//
//	DROP EXTERNAL LANGUAGE language_name [ AUTHORIZATION owner_name ]
func (p *Parser) parseDropExternalLanguageStmt() (*nodes.SecurityStmt, error) {
	return p.extDropWithAuthorization("EXTERNAL LANGUAGE")
}

func (p *Parser) extDropWithAuthorization(objectType string) (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if err := p.extAuthorization(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- EXTERNAL RESOURCE POOL ----------

// parseExternalResourcePoolBody reads the name and the options of CREATE or
// ALTER EXTERNAL RESOURCE POOL (TSql160.g:22327-22440):
//
//	pool_name [ WITH ( <parameter> [ ,...n ] ) ]
//
//	<parameter> ::=
//	    MAX_CPU_PERCENT = 1..100 | MAX_MEMORY_PERCENT = 1..100 | MAX_PROCESSES = integer
//	  | AFFINITY { CPU = { AUTO | ( <range> [ ,...n ] ) } | NUMANODE = ( <range> [ ,...n ] ) }
//	<range> ::= integer [ TO integer ]
//
// SQL Server 2022 requires the WITH clause in ALTER, takes each parameter
// once (Msg 1039), and a percent from 1 to 100 (Msg 1080).
func (p *Parser) parseExternalResourcePoolBody(stmt *nodes.SecurityStmt, requireWith bool) error {
	name, err := p.extName()
	if err != nil {
		return err
	}
	stmt.Name = name
	if p.cur.Type != kwWITH {
		if requireWith {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		key := p.extWordIn("MAX_CPU_PERCENT", "MAX_MEMORY_PERCENT", "MAX_PROCESSES", "AFFINITY")
		if key == "" || seen[key] {
			return p.syntaxErrorAtCur()
		}
		seen[key] = true
		p.advance()
		var val string
		if key == "AFFINITY" {
			if val, err = p.extPoolAffinity(); err != nil {
				return err
			}
		} else {
			if _, err := p.expect('='); err != nil {
				return err
			}
			if key != "MAX_PROCESSES" && p.extIsInteger() && (p.cur.Ival < 1 || p.cur.Ival > 100) {
				return p.syntaxErrorAtCur()
			}
			if val, err = p.extInteger(); err != nil {
				return err
			}
		}
		opts = append(opts, &nodes.String{Str: key + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	stmt.Options = &nodes.List{Items: opts}
	return nil
}

// extPoolAffinity reads what follows AFFINITY in an external resource pool:
// CPU = { AUTO | ( ranges ) } or NUMANODE = ( ranges ).
func (p *Parser) extPoolAffinity() (string, error) {
	kind, err := p.extExpectWordIn("CPU", "NUMANODE")
	if err != nil {
		return "", err
	}
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	if kind == "CPU" && p.extIsWord("AUTO") {
		p.advance()
		return "CPU AUTO", nil
	}
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	var ranges []string
	for {
		lo, err := p.extInteger()
		if err != nil {
			return "", err
		}
		if p.cur.Type == kwTO {
			p.advance()
			hi, err := p.extInteger()
			if err != nil {
				return "", err
			}
			lo += " TO " + hi
		}
		ranges = append(ranges, lo)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return kind + " (" + strings.Join(ranges, ", ") + ")", nil
}

// ---------- EXTERNAL STREAM ----------

// parseCreateExternalStreamStmt parses CREATE EXTERNAL STREAM (Azure SQL Edge).
// Caller has consumed CREATE EXTERNAL STREAM.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-stream-transact-sql
// Grammar: TSql160.g:23900 (createExternalStreamStatement)
//
//	CREATE EXTERNAL STREAM stream_name
//	    WITH ( <option> [ ,...n ] )
//
//	<option> ::= { LOCATION | INPUT_OPTIONS | OUTPUT_OPTIONS } = 'string'
//	           | { DATA_SOURCE | FILE_FORMAT } = { name | 'string' }
func (p *Parser) parseCreateExternalStreamStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL STREAM",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.extName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		key, _, err := p.extOptionHead(seen, "LOCATION", "INPUT_OPTIONS", "OUTPUT_OPTIONS", "DATA_SOURCE", "FILE_FORMAT")
		if err != nil {
			return nil, err
		}
		var val string
		if key == "DATA_SOURCE" || key == "FILE_FORMAT" {
			val, err = p.extIdentifierOrString()
		} else {
			val, err = p.extString()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.String{Str: key + "=" + val})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// ---------- COPY INTO ----------

// extCopyOptions are the options of COPY INTO ... WITH
// (https://learn.microsoft.com/en-us/sql/t-sql/statements/copy-into-transact-sql).
var extCopyOptions = []string{
	"FILE_TYPE", "FILE_FORMAT", "CREDENTIAL", "ERRORFILE", "ERRORFILE_DATA_SOURCE",
	"ERRORFILE_CREDENTIAL", "MAXERRORS", "COMPRESSION", "FIELDQUOTE",
	"FIELDTERMINATOR", "ROWTERMINATOR", "FIRSTROW", "DATEFORMAT", "ENCODING",
	"IDENTITY_INSERT", "AUTO_CREATE_TABLE", "PARSER_VERSION", "MATCH_COLUMN_COUNT",
}

// parseCopyIntoStmt parses a COPY INTO statement (Azure Synapse / Fabric).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/copy-into-transact-sql
// Grammar: TSql160.g:280 (copyStatement)
//
//	COPY INTO [ schema. ] table_name
//	[ ( column_name [ DEFAULT default_value ] [ field_number ] [ ,...n ] ) ]
//	FROM '<external_location>' [ ,...n ]
//	[ WITH ( <option> [ ,...n ] ) ]
//
//	<option> ::= { CREDENTIAL | ERRORFILE_CREDENTIAL } = ( IDENTITY = 'string' [ , SECRET = 'string' ] )
//	           | <other option> = { identifier | integer | 'string' }
func (p *Parser) parseCopyIntoStmt() (*nodes.CopyIntoStmt, error) {
	loc := p.pos()
	p.advance() // consume COPY
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	stmt := &nodes.CopyIntoStmt{Loc: nodes.Loc{Start: loc, End: -1}}
	table, err := p.extRequireTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// ( column_name [ DEFAULT default_value ] [ field_number ] [ ,...n ] )
	// (copyColumnOption, TSql160.g:341)
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for {
			colLoc := p.pos()
			name, err := p.extName()
			if err != nil {
				return nil, err
			}
			col := &nodes.CopyIntoColumn{Name: name, Loc: nodes.Loc{Start: colLoc, End: -1}}
			if p.cur.Type == kwDEFAULT {
				p.advance()
				if col.DefaultValue, err = p.extCopyDefaultValue(); err != nil {
					return nil, err
				}
			}
			if p.extIsInteger() {
				col.FieldNumber = int(p.cur.Ival)
				p.advance()
			}
			col.Loc.End = p.prevEnd()
			cols = append(cols, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.ColumnList = &nodes.List{Items: cols}
	}

	// FROM '<external_location>' [ ,...n ]
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	var sources []nodes.Node
	for {
		if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		src, err := p.extNonEmptyString()
		if err != nil {
			return nil, err
		}
		sources = append(sources, &nodes.String{Str: src[1 : len(src)-1]})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	stmt.Sources = &nodes.List{Items: sources}

	// WITH ( option [ ,...n ] ) (copyWithClause, TSql160.g:363)
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var opts []nodes.Node
		seen := map[string]bool{}
		for {
			opt, err := p.extCopyOption(seen)
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// extCopyDefaultValue reads a COPY INTO column default: a literal or a signed
// number (defaultValueLiteral, TSql160.g:33460).
func (p *Parser) extCopyDefaultValue() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokMONEY, kwNULL:
	case '-':
		if next := p.peekNext().Type; next != tokICONST && next != tokFCONST {
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
}

// extCopyOption reads one COPY INTO option (copyOption, TSql160.g:383): a
// name from extCopyOptions, once each, and its value. CREDENTIAL and
// ERRORFILE_CREDENTIAL take ( IDENTITY = 'string' [ , SECRET = 'string' ] );
// the rest an identifier, an integer or a string (singleValueTypeCopyOption).
func (p *Parser) extCopyOption(seen map[string]bool) (nodes.Node, error) {
	key, _, err := p.extOptionHead(seen, extCopyOptions...)
	if err != nil {
		return nil, err
	}
	if key == "CREDENTIAL" || key == "ERRORFILE_CREDENTIAL" {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwIDENTITY); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		identity, err := p.extString()
		if err != nil {
			return nil, err
		}
		val := "IDENTITY = " + identity
		if p.cur.Type == ',' {
			p.advance()
			if _, err := p.extExpectWordIn("SECRET"); err != nil {
				return nil, err
			}
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			secret, err := p.extString()
			if err != nil {
				return nil, err
			}
			val += ", SECRET = " + secret
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.String{Str: key + "=(" + val + ")"}, nil
	}
	var val string
	switch {
	case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
		val, err = p.extString()
	case p.extIsInteger():
		val, err = p.extInteger()
	default:
		val, err = p.extIdentifier()
	}
	if err != nil {
		return nil, err
	}
	return &nodes.String{Str: key + "=" + val}, nil
}

// ---------- EXTERNAL MODEL and EXTERNAL STREAMING JOB ----------
//
// Neither is in TSql160.g: EXTERNAL MODEL is SQL Server 2025's and EXTERNAL
// STREAMING JOB Azure SQL Edge's.

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
//	    [ , CREDENTIAL = <credential_name> ]
//	    [ , PARAMETERS = '{"valid":"JSON"}' ]
//	    [ , LOCAL_RUNTIME_PATH = 'path to the ONNX Runtime files' ]
//	  )
//	[ ; ]
func (p *Parser) parseCreateExternalModelStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc, End: -1},
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
func (p *Parser) parseAlterExternalModelStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc, End: -1},
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropExternalModelStmt parses DROP EXTERNAL MODEL.
// Caller has consumed DROP EXTERNAL MODEL.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-external-model-transact-sql
//
//	DROP EXTERNAL MODEL external_model_object_name
//	[ ; ]
func (p *Parser) parseDropExternalModelStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "EXTERNAL MODEL",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// external_model_object_name
	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateExternalStreamingJobStmt parses CREATE EXTERNAL STREAMING JOB.
// Caller has consumed CREATE EXTERNAL STREAMING JOB.
//
//	CREATE EXTERNAL STREAMING JOB job_name
//	    WITH ( options )
//	    AS query_string
//	[ ; ]
func (p *Parser) parseCreateExternalStreamingJobStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "EXTERNAL STREAMING JOB",
		Loc:        nodes.Loc{Start: loc, End: -1},
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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
