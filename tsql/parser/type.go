// Package parser - type.go implements T-SQL data type parsing.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseDataType parses a T-SQL data type reference.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/data-types/data-types-transact-sql
//
//	data_type = [ schema_name '.' ] type_name [ '(' length ')' | '(' precision ',' scale ')' | '(' MAX ')' ]
//	type_name = INT | BIGINT | SMALLINT | TINYINT | BIT
//	          | FLOAT | REAL | DECIMAL | NUMERIC | MONEY | SMALLMONEY
//	          | CHAR | VARCHAR | NCHAR | NVARCHAR | TEXT | NTEXT
//	          | BINARY | VARBINARY | IMAGE
//	          | DATE | DATETIME | DATETIME2 | SMALLDATETIME | DATETIMEOFFSET | TIME
//	          | UNIQUEIDENTIFIER | XML | SQL_VARIANT | GEOGRAPHY | GEOMETRY | HIERARCHYID
//	          | user_defined_type
func (p *Parser) parseDataType() *nodes.DataType {
	loc := p.pos()

	// Get type name - could be keyword (INT, VARCHAR, etc.) or identifier
	var name string
	if p.cur.Type == tokIDENT {
		name = p.cur.Str
		p.advance()
	} else if p.cur.Type >= kwADD && p.cur.Str != "" {
		// Keywords used as type names
		name = p.cur.Str
		p.advance()
	} else {
		return nil
	}

	dt := &nodes.DataType{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	// Check for schema-qualified type: schema.typename
	if p.cur.Type == '.' {
		p.advance()
		dt.Schema = name
		if p.isIdentLike() {
			dt.Name = p.cur.Str
			p.advance()
		}
	}

	// Check for (precision[, scale]) or (MAX) or (length)
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == kwMAX {
			dt.MaxLength = true
			p.advance()
		} else {
			dt.Length = p.parseExpr()
			if _, ok := p.match(','); ok {
				dt.Scale = p.parseExpr()
				dt.Precision = dt.Length
				dt.Length = nil
			}
		}
		_, _ = p.expect(')')
	}

	dt.Loc.End = p.pos()
	return dt
}
