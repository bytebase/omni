// Package parser - bulk_insert.go implements T-SQL BULK INSERT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBulkInsertStmt parses a BULK INSERT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/bulk-insert-transact-sql
//
//	BULK INSERT [ database_name . [ schema_name ] . | schema_name . ] table_name
//	    FROM 'data_file'
//	    [ WITH ( option [,...n] ) ]
func (p *Parser) parseBulkInsertStmt() *nodes.BulkInsertStmt {
	loc := p.pos()
	p.advance() // consume BULK

	// Consume INSERT keyword
	p.matchIdentCI("INSERT")

	stmt := &nodes.BulkInsertStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name (possibly qualified: db.schema.table)
	stmt.Table = p.parseTableRef()

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

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBulkInsertOption parses a single BULK INSERT WITH option.
// Options can be:
//   - FLAG (e.g. TABLOCK, CHECK_CONSTRAINTS, FIRE_TRIGGERS, KEEPNULLS, KEEPIDENTITY)
//   - KEY = value (e.g. BATCHSIZE = 1000, FIELDTERMINATOR = ',', ROWTERMINATOR = '\n')
//   - KEY = 'string_value'
func (p *Parser) parseBulkInsertOption() nodes.Node {
	if !p.isIdentLike() && p.cur.Type != kwFILE {
		return nil
	}

	name := strings.ToUpper(p.cur.Str)
	p.advance()

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
			if p.isIdentLike() {
				valStr = p.cur.Str
				p.advance()
			}
		}
		return &nodes.String{Str: name + "=" + valStr}
	}

	// Plain flag option
	return &nodes.String{Str: name}
}
