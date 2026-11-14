// Package parser - bulk_insert.go implements T-SQL BULK INSERT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// bulkOptionKind says what follows a BULK INSERT or INSERT BULK option name.
type bulkOptionKind int

const (
	bulkOptionFlag     bulkOptionKind = iota // no value
	bulkOptionNumber                         // = integer or decimal
	bulkOptionString                         // = string
	bulkOptionOrder                          // ( column [ ASC | DESC ] [ ,...n ] )
	bulkOptionCSV                            // = 'CSV'
	bulkOptionFileType                       // = 'char' | 'native' | 'widechar' | 'widenative'
)

// bulkInsertOptionKinds are the options of BULK INSERT WITH (...), from
// SqlScriptDOM's BulkInsertFlagOptionsHelper, BulkInsertIntOptionHelper and
// BulkInsertStringOptionHelper, as the pinned SQL Server takes them: FORMAT
// takes only 'CSV' and DATAFILETYPE its four types (Msg 102 otherwise).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/bulk-insert-transact-sql
var bulkInsertOptionKinds = map[string]bulkOptionKind{
	"CHECK_CONSTRAINTS":      bulkOptionFlag,
	"FIRE_TRIGGERS":          bulkOptionFlag,
	"KEEPIDENTITY":           bulkOptionFlag,
	"KEEPNULLS":              bulkOptionFlag,
	"TABLOCK":                bulkOptionFlag,
	"BATCHSIZE":              bulkOptionNumber,
	"FIRSTROW":               bulkOptionNumber,
	"KILOBYTES_PER_BATCH":    bulkOptionNumber,
	"LASTROW":                bulkOptionNumber,
	"MAXERRORS":              bulkOptionNumber,
	"ROWS_PER_BATCH":         bulkOptionNumber,
	"CODEPAGE":               bulkOptionString,
	"DATA_SOURCE":            bulkOptionString,
	"ERRORFILE":              bulkOptionString,
	"ERRORFILE_DATA_SOURCE":  bulkOptionString,
	"FIELDQUOTE":             bulkOptionString,
	"FIELDTERMINATOR":        bulkOptionString,
	"FORMATFILE":             bulkOptionString,
	"FORMATFILE_DATA_SOURCE": bulkOptionString,
	"ROWTERMINATOR":          bulkOptionString,
	"DATAFILETYPE":           bulkOptionFileType,
	"FORMAT":                 bulkOptionCSV,
	"ORDER":                  bulkOptionOrder,
}

// insertBulkOptionKinds are the options of INSERT BULK WITH (...):
// simpleInsertBulkOption (TSql160.g:6171) takes a flag or a number, and the
// pinned SQL Server takes these words there and refuses BULK INSERT's other
// ones (Msg 102).
var insertBulkOptionKinds = map[string]bulkOptionKind{
	"CHECK_CONSTRAINTS":   bulkOptionFlag,
	"FIRE_TRIGGERS":       bulkOptionFlag,
	"KEEP_NULLS":          bulkOptionFlag,
	"NO_TRIGGERS":         bulkOptionFlag,
	"TABLOCK":             bulkOptionFlag,
	"KILOBYTES_PER_BATCH": bulkOptionNumber,
	"ROWS_PER_BATCH":      bulkOptionNumber,
	"ORDER":               bulkOptionOrder,
}

// bulkDataFileTypes are the values of BULK INSERT's DATAFILETYPE.
var bulkDataFileTypes = map[string]bool{"CHAR": true, "NATIVE": true, "WIDECHAR": true, "WIDENATIVE": true}

// parseBulkInsertStmt parses a BULK INSERT statement.
//
// BNF: mssql/parser/bnf/bulk-insert-transact-sql.bnf
//
//	BULK INSERT
//	   { database_name.schema_name.table_or_view_name | schema_name.table_or_view_name | table_or_view_name }
//	      FROM 'data_file'
//	     [ WITH
//	    (
//	   [ [ , ] DATA_SOURCE = 'data_source_name' ]
//	   [ [ , ] CODEPAGE = { 'RAW' | 'code_page' | 'ACP' | 'OEM' } ]
//	   [ [ , ] DATAFILETYPE = { 'char' | 'widechar' | 'native' | 'widenative' } ]
//	   [ [ , ] ROWTERMINATOR = 'row_terminator' ]
//	   [ [ , ] FIELDTERMINATOR = 'field_terminator' ]
//	   [ [ , ] FORMAT = 'CSV' ]
//	   [ [ , ] FIELDQUOTE = 'quote_characters' ]
//	   [ [ , ] FIRSTROW = first_row ]
//	   [ [ , ] LASTROW = last_row ]
//	   [ [ , ] FORMATFILE = 'format_file_path' ]
//	   [ [ , ] FORMATFILE_DATA_SOURCE = 'data_source_name' ]
//	   [ [ , ] MAXERRORS = max_errors ]
//	   [ [ , ] ERRORFILE = 'file_name' ]
//	   [ [ , ] ERRORFILE_DATA_SOURCE = 'errorfile_data_source_name' ]
//	   [ [ , ] KEEPIDENTITY ]
//	   [ [ , ] KEEPNULLS ]
//	   [ [ , ] FIRE_TRIGGERS ]
//	   [ [ , ] CHECK_CONSTRAINTS ]
//	   [ [ , ] TABLOCK ]
//	   [ [ , ] ORDER ( { column [ ASC | DESC ] } [ , ...n ] ) ]
//	   [ [ , ] ROWS_PER_BATCH = rows_per_batch ]
//	   [ [ , ] KILOBYTES_PER_BATCH = kilobytes_per_batch ]
//	   [ [ , ] BATCHSIZE = batch_size ]
//	    ) ]
//
// bulkInsertStatement (TSql160.g:6104): the table, FROM and the source are
// required, and the source is a number, a string or an identifier
// (bulkInsertFrom, TSql160.g:6126).
func (p *Parser) parseBulkInsertStmt() (*nodes.BulkInsertStmt, error) {
	loc := p.pos()
	p.advance() // consume BULK

	if _, err := p.expect(kwINSERT); err != nil {
		return nil, err
	}

	stmt := &nodes.BulkInsertStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name (possibly qualified: db.schema.table)
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Table = table

	// FROM { 'data_file' | integer | identifier }
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	switch {
	case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
		stmt.DataFile = p.cur.Str
		p.advance()
	case p.cur.Type == tokICONST && !execIsBinary(p.cur):
		stmt.DataFile = p.cur.Str
		p.advance()
	default:
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.DataFile = name
	}

	// WITH ( option [,...n] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err := p.parseBulkOptionList(bulkInsertOptionKinds)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBulkOptionList parses ( option [ ,...n ] ) of BULK INSERT
// (bulkInsertOptions, TSql160.g:6137) or INSERT BULK (insertBulkOptions,
// TSql160.g:6325), each option one of kinds and given once: SQL Server
// refuses a repeated option (Msg 102, 156, 4130).
func (p *Parser) parseBulkOptionList(kinds map[string]bulkOptionKind) (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		opt, name, err := p.parseBulkInsertOption(kinds)
		if err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, p.newParseError(optLoc, "duplicate option "+name)
		}
		seen[name] = true
		opts = append(opts, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: opts}, nil
}

// parseBulkInsertOption parses a single BULK INSERT or INSERT BULK option
// (bulkInsertOption, TSql160.g:6160; insertBulkOption, TSql160.g:6166), and
// returns it with its name. The name is an undelimited word: SQL Server
// refuses WITH ([TABLOCK]) (Msg 102).
//
//   - FLAG (e.g. TABLOCK, CHECK_CONSTRAINTS, FIRE_TRIGGERS)
//   - KEY = number (BATCHSIZE = 1000), an integer or a decimal
//     (integerOrNumeric, TSql160.g:26396)
//   - KEY = 'string_value'
//   - ORDER ( column [ ASC | DESC ] [ ,...n ] ) (bulkInsertSortOrderOption,
//     TSql160.g:6347)
func (p *Parser) parseBulkInsertOption(kinds map[string]bulkOptionKind) (nodes.Node, string, error) {
	if p.cur.Type != kwORDER && (!p.isIdentLike() || p.isDelimited(p.cur)) {
		return nil, "", p.syntaxErrorAtCur()
	}
	name := strings.ToUpper(p.cur.Str)
	kind, ok := kinds[name]
	if !ok {
		return nil, "", p.syntaxErrorAtCur()
	}
	p.advance()

	switch kind {
	case bulkOptionFlag:
		return &nodes.String{Str: name}, name, nil
	case bulkOptionOrder:
		if _, err := p.expect('('); err != nil {
			return nil, "", err
		}
		var parts []string
		for {
			col, ok := p.parseIdentifier()
			if !ok {
				return nil, "", p.syntaxErrorAtCur()
			}
			if p.cur.Type == kwASC {
				col += " ASC"
				p.advance()
			} else if p.cur.Type == kwDESC {
				col += " DESC"
				p.advance()
			}
			parts = append(parts, col)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, "", err
		}
		return &nodes.String{Str: "ORDER(" + strings.Join(parts, ", ") + ")"}, name, nil
	}

	if _, err := p.expect('='); err != nil {
		return nil, "", err
	}
	var valStr string
	switch kind {
	case bulkOptionNumber:
		if !bulkIsIntegerOrNumeric(p.cur) {
			return nil, "", p.syntaxErrorAtCur()
		}
	case bulkOptionString:
		if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
			return nil, "", p.syntaxErrorAtCur()
		}
	case bulkOptionCSV:
		if (p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST) || strings.ToUpper(p.cur.Str) != "CSV" {
			return nil, "", p.syntaxErrorAtCur()
		}
	case bulkOptionFileType:
		if (p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST) || !bulkDataFileTypes[strings.ToUpper(p.cur.Str)] {
			return nil, "", p.syntaxErrorAtCur()
		}
	}
	valStr = p.cur.Str
	p.advance()
	return &nodes.String{Str: name + "=" + valStr}, name, nil
}

// bulkIsIntegerOrNumeric reports whether tok is an integer or a decimal
// without an exponent, which SQL Server refuses here (`= 1e3`, `= 0x01`,
// `= $1` are Msg 102).
func bulkIsIntegerOrNumeric(tok Token) bool {
	switch tok.Type {
	case tokICONST:
		return !execIsBinary(tok)
	case tokFCONST:
		return !execIsBinary(tok) && !strings.ContainsAny(tok.Str, "eE")
	}
	return false
}

// parseInsertBulkStmt parses an INSERT BULK statement, after INSERT BULK.
//
// BNF: mssql/parser/bnf/insert-bulk-transact-sql.bnf
//
//	INSERT BULK schemaObjectThreePartName
//	  ( column_name data_type [ NULL | NOT NULL ] [ ,...n ] )
//	  [ WITH ( option [ = value ] [ ,...n ] ) ]
//
// insertBulkStatement (TSql160.g:6271) makes the column list optional, but
// SQL Server refuses INSERT BULK without one (Msg 102).
func (p *Parser) parseInsertBulkStmt() (*nodes.InsertBulkStmt, error) {
	stmt := &nodes.InsertBulkStmt{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	// table name (three-part name)
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Table = table

	// column definition list (coldefList, TSql160.g:6290)
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var cols []nodes.Node
	for {
		col, err := p.parseInsertBulkColumnDef()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Columns = &nodes.List{Items: cols}

	// optional WITH ( options )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err := p.parseBulkOptionList(insertBulkOptionKinds)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseInsertBulkColumnDef parses a column definition in INSERT BULK
// (coldefItem, TSql160.g:6309): a column with a type and an optional
// collation (columnDefinitionBasic, TSql160.g:29519), or TIMESTAMP alone
// (columnDefinitionEx), then NULL or NOT NULL.
//
//	{ column_name data_type [ COLLATE collation ] | TIMESTAMP } [ NULL | NOT NULL ]
func (p *Parser) parseInsertBulkColumnDef() (*nodes.InsertBulkColumnDef, error) {
	col := &nodes.InsertBulkColumnDef{
		Loc: nodes.Loc{Start: p.pos(), End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	col.Name = name

	// An undelimited TIMESTAMP stands alone; a column named [TIMESTAMP]
	// takes a type like any other, which SQL Server reads the same way.
	isTimestamp := !p.isDelimited(p.prev) && strings.ToUpper(name) == "TIMESTAMP"
	if !(isTimestamp && (p.cur.Type == ',' || p.cur.Type == ')' || p.cur.Type == kwNULL || p.cur.Type == kwNOT)) {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		if dt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		col.DataType = dt
		if p.cur.Type == kwCOLLATE {
			p.advance() // consume COLLATE
			if _, err := p.parseCollationName(); err != nil {
				return nil, err
			}
		}
	}

	// Optional NULL / NOT NULL
	if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		if _, err := p.expect(kwNULL); err != nil {
			return nil, err
		}
		f := false
		col.Nullable = &f
	} else if p.cur.Type == kwNULL {
		p.advance() // consume NULL
		t := true
		col.Nullable = &t
	}

	col.Loc.End = p.prevEnd()
	return col, nil
}
