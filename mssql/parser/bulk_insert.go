// Package parser - bulk_insert.go implements T-SQL BULK INSERT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// bulkInsertOptions defines the valid option names for BULK INSERT WITH (...).
// Sourced from SqlScriptDOM BulkInsertFlagOptionsHelper, BulkInsertIntOptionHelper,
// and BulkInsertStringOptionHelper.
var bulkInsertOptions = newOptionSet(
	// Registered keywords that are valid bulk insert options:
	kwORDER,       // ORDER ( column [ASC|DESC] ... )
	kwDATA_SOURCE, // DATA_SOURCE = '...'
	kwFORMAT,      // FORMAT = 'CSV'
).withIdents(
	// Flag options (BulkInsertFlagOptionsHelper):
	"CHECK_CONSTRAINTS",
	"FIRE_TRIGGERS",
	"KEEPIDENTITY",
	"KEEPNULLS",
	"TABLOCK",
	// Int options (BulkInsertIntOptionHelper):
	"BATCHSIZE",
	"FIRSTROW",
	"KILOBYTES_PER_BATCH",
	"LASTROW",
	"MAXERRORS",
	"ROWS_PER_BATCH",
	// String options (BulkInsertStringOptionHelper):
	"CODEPAGE",
	"DATAFILETYPE",
	"ERRORFILE",
	"ERRORFILE_DATA_SOURCE",
	"FIELDQUOTE",
	"FIELDTERMINATOR",
	"FORMATFILE",
	"FORMATFILE_DATA_SOURCE",
	"ROWTERMINATOR",
)

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
func (p *Parser) parseBulkInsertStmt() (*nodes.BulkInsertStmt, error) {
	loc := p.pos()
	p.advance() // consume BULK

	// Consume INSERT keyword
	p.matchIdentCI("INSERT")

	stmt := &nodes.BulkInsertStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name (possibly qualified: db.schema.table)
	stmt.Table , _ = p.parseTableRef()

	// FROM 'data_file'
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
			stmt.DataFile = p.cur.Str
			p.advance()
		}
	}

	// WITH ( option [,...n] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == '(' {
			p.advance() // consume '('
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				opt := p.parseBulkInsertOption()
				if opt != nil {
					opts = append(opts, opt)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBulkInsertOption parses a single BULK INSERT WITH option.
// Options can be:
//   - FLAG (e.g. TABLOCK, CHECK_CONSTRAINTS, FIRE_TRIGGERS, KEEPNULLS, KEEPIDENTITY)
//   - KEY = value (e.g. BATCHSIZE = 1000, FIELDTERMINATOR = ',', ROWTERMINATOR = '\n')
//   - KEY = 'string_value'
//   - ORDER ( column [ ASC | DESC ] [ ,...n ] )
func (p *Parser) parseBulkInsertOption() nodes.Node {
	if !p.isValidOption(bulkInsertOptions) {
		return nil
	}

	name := strings.ToUpper(p.cur.Str)
	p.advance()

	// ORDER ( column [ASC|DESC] [,...n] )
	if name == "ORDER" && p.cur.Type == '(' {
		p.advance() // consume '('
		var parts []string
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isAnyKeywordIdent() {
				col := p.cur.Str
				p.advance()
				// Optional ASC/DESC
				if p.cur.Type == kwASC {
					col += " ASC"
					p.advance()
				} else if p.cur.Type == kwDESC {
					col += " DESC"
					p.advance()
				}
				parts = append(parts, col)
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')') // consume ')'
		return &nodes.String{Str: "ORDER(" + strings.Join(parts, ", ") + ")"}
	}

	// Check for = sign indicating a key=value option
	if p.cur.Type == '=' {
		p.advance() // consume '='
		var valStr string
		switch p.cur.Type {
		case tokSCONST, tokNSCONST:
			valStr = p.cur.Str
			p.advance()
		case tokICONST:
			valStr = p.cur.Str
			p.advance()
		case tokFCONST:
			valStr = p.cur.Str
			p.advance()
		default:
			if p.isAnyKeywordIdent() {
				valStr = p.cur.Str
				p.advance()
			}
		}
		return &nodes.String{Str: name + "=" + valStr}
	}

	// Plain flag option
	return &nodes.String{Str: name}
}
