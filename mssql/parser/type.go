// Package parser - type.go implements T-SQL data type parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
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
func (p *Parser) parseDataType() (*nodes.DataType, error) {
	loc := p.pos()
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}

	national := false
	// A type name is an identifier, which a reserved word is not unless
	// quoted (scalarDataType, TSql160.g:30159): SQL Server reads
	// `ADD a NULL` as a column with no type (Msg 173). DOUBLE PRECISION and
	// NATIONAL { CHAR | CHARACTER | TEXT } [ VARYING ] start with reserved
	// words (doubleDataType, TSql160.g:30283; sqlDataTypeWithNational,
	// TSql160.g:30192).
	switch {
	case p.cur.Type == kwDOUBLE:
		return p.parseDoublePrecisionType(loc)
	case p.cur.Type == kwNATIONAL:
		p.advance() // NATIONAL
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		national = true
	case !p.isIdentLike():
		return nil, nil
	}
	name := p.cur.Str
	p.advance()

	dt := &nodes.DataType{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
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

	// NATIONAL takes CHAR, CHARACTER and TEXT, and VARYING follows CHAR,
	// CHARACTER, NCHAR and BINARY; SQL Server refuses NATIONAL VARCHAR and
	// VARCHAR VARYING (Msg 102). The actions of sqlDataTypeWithNational and
	// sqlDataTypeWithoutNational (TSql160.g:30192, TSql160.g:30224) check
	// the name.
	base := strings.ToUpper(dt.Name)
	if national {
		if base != "CHAR" && base != "CHARACTER" && base != "TEXT" {
			return nil, p.syntaxErrorAtCur()
		}
		dt.Name = "NATIONAL " + dt.Name
	}
	if p.cur.Type == kwVARYING {
		if base != "CHAR" && base != "CHARACTER" && (national || base != "NCHAR" && base != "BINARY") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		dt.Name += " VARYING"
	}

	// Check for (precision[, scale]) or (MAX) or (length)
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == kwMAX {
			dt.MaxLength = true
			p.advance()
		} else {
			length, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			dt.Length = length
			if _, ok := p.match(','); ok {
				scale, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				dt.Scale = scale
				dt.Precision = dt.Length
				dt.Length = nil
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	dt.Loc.End = p.prevEnd()
	return dt, nil
}

// parseDoublePrecisionType parses DOUBLE PRECISION [ ( { MAX | integer } ) ]
// (doubleDataType, TSql160.g:30283), with DOUBLE the current token.
func (p *Parser) parseDoublePrecisionType(loc int) (*nodes.DataType, error) {
	p.advance() // DOUBLE
	if p.cur.Type != kwPRECISION {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // PRECISION
	dt := &nodes.DataType{Name: "DOUBLE PRECISION", Loc: nodes.Loc{Start: loc, End: -1}}
	if p.cur.Type == '(' {
		p.advance()
		switch p.cur.Type {
		case kwMAX:
			dt.MaxLength = true
			p.advance()
		case tokICONST:
			dt.Length = &nodes.Literal{Type: nodes.LitInteger, Ival: p.cur.Ival, Str: p.cur.Str, Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	dt.Loc.End = p.prevEnd()
	return dt, nil
}
