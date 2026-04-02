package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseDataType parses a MySQL data type specification.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/data-types.html
//
//	data_type:
//	    numeric_type
//	    | string_type
//	    | date_time_type
//	    | blob_type
//	    | enum_type
//	    | set_type
//	    | json_type
//	    | bit_type
//	    | binary_type
//
//	numeric_type:
//	    (INT | INTEGER | TINYINT | SMALLINT | MEDIUMINT | BIGINT) [(length)] [UNSIGNED] [ZEROFILL]
//	    | (FLOAT | DOUBLE) [(length [, decimals])] [UNSIGNED] [ZEROFILL]
//	    | DOUBLE PRECISION [(length [, decimals])] [UNSIGNED] [ZEROFILL]
//	    | (DECIMAL | NUMERIC) [(length [, decimals])] [UNSIGNED] [ZEROFILL]
//	    | (BOOL | BOOLEAN)
//
//	string_type:
//	    (CHAR | VARCHAR) [(length)] [CHARACTER SET charset] [COLLATE collation]
//	    | (TEXT | TINYTEXT | MEDIUMTEXT | LONGTEXT) [CHARACTER SET charset] [COLLATE collation]
//
//	date_time_type:
//	    DATE | TIME [(fsp)] | DATETIME [(fsp)] | TIMESTAMP [(fsp)] | YEAR [(length)]
//
//	blob_type:
//	    BLOB | TINYBLOB | MEDIUMBLOB | LONGBLOB
//
//	binary_type:
//	    BINARY [(length)] | VARBINARY (length)
//
//	bit_type:
//	    BIT [(length)]
//
//	enum_type:
//	    ENUM (value_list)
//
//	set_type:
//	    SET (value_list)
//
//	json_type:
//	    JSON
func (p *Parser) parseDataType() (*nodes.DataType, error) {
	start := p.pos()
	dt := &nodes.DataType{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	// Integer types
	case kwINT, kwINTEGER:
		dt.Name = "INT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)

	case kwTINYINT:
		dt.Name = "TINYINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)

	case kwSMALLINT:
		dt.Name = "SMALLINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)

	case kwMEDIUMINT:
		dt.Name = "MEDIUMINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)

	case kwBIGINT:
		dt.Name = "BIGINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)

	// Float types
	case kwFLOAT:
		dt.Name = "FLOAT"
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)

	case kwDOUBLE:
		dt.Name = "DOUBLE"
		p.advance()
		// DOUBLE PRECISION is a synonym
		if p.cur.Type == kwPRECISION {
			p.advance()
		}
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)

	case kwREAL:
		// REAL → DOUBLE synonym
		dt.Name = "DOUBLE"
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)

	// Decimal types
	case kwDECIMAL, kwNUMERIC:
		if p.cur.Type == kwDECIMAL {
			dt.Name = "DECIMAL"
		} else {
			dt.Name = "NUMERIC"
		}
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)

	case kwDEC, kwFIXED:
		// DEC and FIXED → DECIMAL synonyms
		dt.Name = "DECIMAL"
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)

	// Boolean
	case kwBOOL, kwBOOLEAN:
		dt.Name = "BOOLEAN"
		p.advance()

	// Character string types
	case kwCHAR:
		dt.Name = "CHAR"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseCharsetCollate(dt)

	case kwVARCHAR:
		dt.Name = "VARCHAR"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseCharsetCollate(dt)

	case kwTEXT:
		dt.Name = "TEXT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseCharsetCollate(dt)

	case kwTINYTEXT:
		dt.Name = "TINYTEXT"
		p.advance()
		p.parseCharsetCollate(dt)

	case kwMEDIUMTEXT:
		dt.Name = "MEDIUMTEXT"
		p.advance()
		p.parseCharsetCollate(dt)

	case kwLONGTEXT:
		dt.Name = "LONGTEXT"
		p.advance()
		p.parseCharsetCollate(dt)

	// Date/time types
	case kwDATE:
		dt.Name = "DATE"
		p.advance()

	case kwTIME:
		dt.Name = "TIME"
		p.advance()
		p.parseOptionalLength(dt) // fsp

	case kwDATETIME:
		dt.Name = "DATETIME"
		p.advance()
		p.parseOptionalLength(dt) // fsp

	case kwTIMESTAMP:
		dt.Name = "TIMESTAMP"
		p.advance()
		p.parseOptionalLength(dt) // fsp

	case kwYEAR:
		dt.Name = "YEAR"
		p.advance()
		p.parseOptionalLength(dt)

	// Blob types
	case kwBLOB:
		dt.Name = "BLOB"
		p.advance()
		p.parseOptionalLength(dt)

	case kwTINYBLOB:
		dt.Name = "TINYBLOB"
		p.advance()

	case kwMEDIUMBLOB:
		dt.Name = "MEDIUMBLOB"
		p.advance()

	case kwLONGBLOB:
		dt.Name = "LONGBLOB"
		p.advance()

	// Binary types
	case kwBINARY:
		dt.Name = "BINARY"
		p.advance()
		p.parseOptionalLength(dt)

	case kwVARBINARY:
		dt.Name = "VARBINARY"
		p.advance()
		p.parseOptionalLength(dt)

	// Bit type
	case kwBIT:
		dt.Name = "BIT"
		p.advance()
		p.parseOptionalLength(dt)

	// ENUM type
	case kwENUM:
		dt.Name = "ENUM"
		p.advance()
		vals, err := p.parseStringValueList()
		if err != nil {
			return nil, err
		}
		dt.EnumValues = vals
		p.parseCharsetCollate(dt)

	// SET type (the SQL SET data type, not the SET statement keyword)
	case kwSET:
		// We must be careful: SET is also a statement keyword.
		// In data type context, SET is followed by '('.
		if p.peekNext().Type != '(' {
			return nil, &ParseError{
				Message:  "expected data type",
				Position: p.cur.Loc,
			}
		}
		dt.Name = "SET"
		p.advance()
		vals, err := p.parseStringValueList()
		if err != nil {
			return nil, err
		}
		dt.EnumValues = vals
		p.parseCharsetCollate(dt)

	// JSON type
	case kwJSON:
		dt.Name = "JSON"
		p.advance()

	// Spatial types
	case kwGEOMETRY, kwPOINT, kwLINESTRING, kwPOLYGON,
		kwMULTIPOINT, kwMULTILINESTRING, kwMULTIPOLYGON,
		kwGEOMETRYCOLLECTION:
		dt.Name = p.cur.Str
		p.advance()

	// SERIAL → BIGINT UNSIGNED AUTO_INCREMENT UNIQUE
	case kwSERIAL:
		dt.Name = p.cur.Str
		p.advance()

	// NATIONAL CHAR / NATIONAL VARCHAR → utf8mb3 character set
	case kwNATIONAL:
		p.advance()
		if p.cur.Type == kwCHAR {
			dt.Name = "CHAR"
			dt.Charset = "utf8mb3"
			p.advance()
			p.parseOptionalLength(dt)
			p.parseCharsetCollate(dt)
		} else if p.cur.Type == kwVARCHAR {
			dt.Name = "VARCHAR"
			dt.Charset = "utf8mb3"
			p.advance()
			p.parseOptionalLength(dt)
			p.parseCharsetCollate(dt)
		} else {
			return nil, &ParseError{
				Message:  "expected CHAR or VARCHAR after NATIONAL",
				Position: p.cur.Loc,
			}
		}

	// NCHAR → CHAR with utf8mb3 charset
	case kwNCHAR:
		dt.Name = "CHAR"
		dt.Charset = "utf8mb3"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseCharsetCollate(dt)

	// NVARCHAR → VARCHAR with utf8mb3 charset
	case kwNVARCHAR:
		dt.Name = "VARCHAR"
		dt.Charset = "utf8mb3"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseCharsetCollate(dt)

	case kwINT1:
		dt.Name = "TINYINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwINT2:
		dt.Name = "SMALLINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwINT3:
		dt.Name = "MEDIUMINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwINT4:
		dt.Name = "INT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwINT8:
		dt.Name = "BIGINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwMIDDLEINT:
		dt.Name = "MEDIUMINT"
		p.advance()
		p.parseOptionalLength(dt)
		p.parseUnsignedZerofill(dt)
	case kwFLOAT4:
		dt.Name = "FLOAT"
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)
	case kwFLOAT8:
		dt.Name = "DOUBLE"
		p.advance()
		p.parseOptionalPrecision(dt)
		p.parseUnsignedZerofill(dt)
	case kwLONG:
		// LONG / LONG VARCHAR → MEDIUMTEXT, LONG VARBINARY → MEDIUMBLOB
		p.advance()
		if p.cur.Type == kwVARCHAR {
			p.advance()
			dt.Name = "MEDIUMTEXT"
			p.parseCharsetCollate(dt)
		} else if p.cur.Type == kwVARBINARY {
			p.advance()
			dt.Name = "MEDIUMBLOB"
		} else {
			dt.Name = "MEDIUMTEXT"
			p.parseCharsetCollate(dt)
		}

	default:
		if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, &ParseError{
			Message:  "expected data type",
			Position: p.cur.Loc,
		}
	}

	dt.Loc.End = p.pos()
	return dt, nil
}

// parseOptionalLength parses an optional (N) length specification.
func (p *Parser) parseOptionalLength(dt *nodes.DataType) {
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == tokICONST {
			dt.Length = int(p.cur.Ival)
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}
}

// parseOptionalPrecision parses an optional (M) or (M,D) precision specification.
func (p *Parser) parseOptionalPrecision(dt *nodes.DataType) {
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == tokICONST {
			dt.Length = int(p.cur.Ival)
			p.advance()
		}
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == tokICONST {
				dt.Scale = int(p.cur.Ival)
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}
}

// parseUnsignedZerofill parses optional SIGNED, UNSIGNED, and ZEROFILL modifiers.
// SIGNED is the default for all numeric types, so we accept and ignore it.
func (p *Parser) parseUnsignedZerofill(dt *nodes.DataType) {
	// SIGNED is a keyword token (kwSIGNED) — accept and ignore it (it's the default).
	if p.cur.Type == kwSIGNED {
		p.advance()
	} else if _, ok := p.match(kwUNSIGNED); ok {
		dt.Unsigned = true
	}
	if _, ok := p.match(kwZEROFILL); ok {
		dt.Zerofill = true
	}
}

// parseCharsetCollate parses optional CHARACTER SET and COLLATE clauses.
func (p *Parser) parseCharsetCollate(dt *nodes.DataType) {
	// CHARACTER SET charset_name | CHARSET charset_name
	if _, ok := p.match(kwCHARSET); ok {
		if p.isIdentToken() {
			dt.Charset, _, _ = p.parseIdent()
		}
	} else if p.cur.Type == kwCHARACTER {
		p.advance()
		if _, ok := p.match(kwSET); ok {
			if p.isIdentToken() {
				dt.Charset, _, _ = p.parseIdent()
			}
		}
	}

	// COLLATE collation_name
	if _, ok := p.match(kwCOLLATE); ok {
		if p.isIdentToken() {
			dt.Collate, _, _ = p.parseIdent()
		}
	}
}

// parseStringValueList parses a parenthesized list of string literals: ('a', 'b', 'c')
func (p *Parser) parseStringValueList() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var vals []string
	for {
		if p.cur.Type == tokSCONST {
			vals = append(vals, p.cur.Str)
			p.advance()
		} else {
			return nil, &ParseError{
				Message:  "expected string literal in value list",
				Position: p.cur.Loc,
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return vals, nil
}
