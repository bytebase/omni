package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
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
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwTINYINT:
		dt.Name = "TINYINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwSMALLINT:
		dt.Name = "SMALLINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwMEDIUMINT:
		dt.Name = "MEDIUMINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwBIGINT:
		dt.Name = "BIGINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	// Float types
	case kwFLOAT:
		dt.Name = "FLOAT"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwDOUBLE:
		dt.Name = "DOUBLE"
		p.advance()
		// DOUBLE PRECISION is a synonym
		if p.cur.Type == kwPRECISION {
			p.advance()
		}
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwREAL:
		// REAL → DOUBLE synonym
		dt.Name = "DOUBLE"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	// Decimal types
	case kwDECIMAL, kwNUMERIC:
		if p.cur.Type == kwDECIMAL {
			dt.Name = "DECIMAL"
		} else {
			dt.Name = "NUMERIC"
		}
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	case kwDEC, kwFIXED:
		// DEC and FIXED → DECIMAL synonyms
		dt.Name = "DECIMAL"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	// Boolean
	case kwBOOL, kwBOOLEAN:
		dt.Name = "BOOLEAN"
		p.advance()

	// Character string types. CHAR and CHARACTER name the same type, and
	// VARYING after either makes it VARCHAR, as the one word VARCHARACTER
	// does: https://dev.mysql.com/doc/refman/8.0/en/char.html. Every spelling
	// of VARCHAR carries a length and every spelling of CHAR may leave it out,
	// so TiDB takes `char` and `varchar(1)` but rejects `varchar`.
	case kwCHAR, kwCHARACTER:
		dt.Name = "CHAR"
		p.advance()
		if p.cur.Type == kwVARYING {
			dt.Name = "VARCHAR"
			p.advance()
			if err := p.parseRequiredLength(dt); err != nil {
				return nil, err
			}
		} else {
			if err := p.parseOptionalLength(dt); err != nil {
				return nil, err
			}
		}
		p.parseCharsetCollate(dt)

	case kwVARCHAR, kwVARCHARACTER:
		dt.Name = "VARCHAR"
		p.advance()
		if err := p.parseRequiredLength(dt); err != nil {
			return nil, err
		}
		p.parseCharsetCollate(dt)

	case kwTEXT:
		dt.Name = "TEXT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
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
		if err := p.parseOptionalLength(dt); err != nil { // fsp
			return nil, err
		}

	case kwDATETIME:
		dt.Name = "DATETIME"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil { // fsp
			return nil, err
		}

	case kwTIMESTAMP:
		dt.Name = "TIMESTAMP"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil { // fsp
			return nil, err
		}

	// YEAR, and SQL_TSI_YEAR which names the same type: the spelling
	// TIMESTAMPADD takes for its unit is a type name too
	// (https://dev.mysql.com/doc/refman/8.0/en/date-and-time-type-syntax.html).
	// Both carry an optional display width and the numeric modifiers, so
	// `YEAR ZEROFILL` and `SQL_TSI_YEAR(4) ZEROFILL` are columns; neither takes
	// a charset, so `YEAR(4) BINARY` is not.
	case kwYEAR, kwSQL_TSI_YEAR:
		dt.Name = "YEAR"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)

	// Blob types
	case kwBLOB:
		dt.Name = "BLOB"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}

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
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}

	case kwVARBINARY:
		dt.Name = "VARBINARY"
		p.advance()
		if err := p.parseRequiredLength(dt); err != nil {
			return nil, err
		}

	// Bit type
	case kwBIT:
		dt.Name = "BIT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}

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

	// NATIONAL CHAR / NATIONAL VARCHAR → utf8mb3 character set, in the same
	// spellings the bare types take.
	case kwNATIONAL:
		p.advance()
		if p.cur.Type == kwCHAR || p.cur.Type == kwCHARACTER {
			dt.Name = "CHAR"
			dt.Charset = "utf8mb3"
			p.advance()
			if p.cur.Type == kwVARYING {
				dt.Name = "VARCHAR"
				p.advance()
				if err := p.parseRequiredLength(dt); err != nil {
					return nil, err
				}
			} else {
				if err := p.parseOptionalLength(dt); err != nil {
					return nil, err
				}
			}
			p.parseCharsetCollate(dt)
		} else if p.cur.Type == kwVARCHAR || p.cur.Type == kwVARCHARACTER {
			dt.Name = "VARCHAR"
			dt.Charset = "utf8mb3"
			p.advance()
			if err := p.parseRequiredLength(dt); err != nil {
				return nil, err
			}
			p.parseCharsetCollate(dt)
		} else {
			return nil, &ParseError{
				Message:  "expected CHAR or VARCHAR after NATIONAL",
				Position: p.cur.Loc,
			}
		}

	// NCHAR → CHAR with utf8mb3 charset, and NCHAR VARCHAR / VARCHARACTER /
	// VARYING → VARCHAR with it.
	case kwNCHAR:
		dt.Charset = "utf8mb3"
		p.advance()
		if p.cur.Type == kwVARCHAR || p.cur.Type == kwVARCHARACTER || p.cur.Type == kwVARYING {
			dt.Name = "VARCHAR"
			p.advance()
			if err := p.parseRequiredLength(dt); err != nil {
				return nil, err
			}
		} else {
			dt.Name = "CHAR"
			if err := p.parseOptionalLength(dt); err != nil {
				return nil, err
			}
		}
		p.parseCharsetCollate(dt)

	// NVARCHAR → VARCHAR with utf8mb3 charset
	case kwNVARCHAR:
		dt.Name = "VARCHAR"
		dt.Charset = "utf8mb3"
		p.advance()
		if err := p.parseRequiredLength(dt); err != nil {
			return nil, err
		}
		p.parseCharsetCollate(dt)

	case kwINT1:
		dt.Name = "TINYINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwINT2:
		dt.Name = "SMALLINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwINT3:
		dt.Name = "MEDIUMINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwINT4:
		dt.Name = "INT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwINT8:
		dt.Name = "BIGINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwMIDDLEINT:
		dt.Name = "MEDIUMINT"
		p.advance()
		if err := p.parseOptionalLength(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwFLOAT4:
		dt.Name = "FLOAT"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwFLOAT8:
		dt.Name = "DOUBLE"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		p.parseUnsignedZerofill(dt)
	case kwLONG:
		// LONG and the VARCHAR spellings of it -- LONG VARCHAR, LONG
		// VARCHARACTER, LONG CHAR VARYING, LONG CHARACTER VARYING -- are
		// MEDIUMTEXT; LONG VARBINARY is MEDIUMBLOB. None of them carries a
		// length. CHAR and CHARACTER after LONG are ambiguous, because
		// `long char set utf8` is LONG with a charset clause, so the type
		// claims either only when VARYING follows it.
		p.advance()
		switch {
		case p.cur.Type == kwVARCHAR || p.cur.Type == kwVARCHARACTER:
			p.advance()
			dt.Name = "MEDIUMTEXT"
			p.parseCharsetCollate(dt)
		case (p.cur.Type == kwCHAR || p.cur.Type == kwCHARACTER) && p.peekNext().Type == kwVARYING:
			p.advance()
			p.advance()
			dt.Name = "MEDIUMTEXT"
			p.parseCharsetCollate(dt)
		case p.cur.Type == kwVARBINARY:
			p.advance()
			dt.Name = "MEDIUMBLOB"
		default:
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

	dt.Loc.End = p.prev.End
	return dt, nil
}

// parseRequiredLength parses the length a VARCHAR-shaped type must carry.
// TiDB takes CHAR with a length or without one but VARCHAR only with one, in
// every spelling: `varcharacter(1)`, `char varying(50)` and `nchar
// varying(50)` are accepted and those three without a length are not. The LONG
// forms carry no length at all, so they do not come here.
func (p *Parser) parseRequiredLength(dt *nodes.DataType) error {
	if p.cur.Type != '(' {
		return &ParseError{
			Message:  "expected a length for " + dt.Name,
			Position: p.cur.Loc,
		}
	}
	p.advance()
	if p.cur.Type != tokICONST {
		return &ParseError{
			Message:  "expected a length for " + dt.Name,
			Position: p.cur.Loc,
		}
	}
	dt.Length = int(p.cur.Ival)
	p.advance()
	if p.cur.Type != ')' {
		return &ParseError{
			Message:  "expected ) after the length for " + dt.Name,
			Position: p.cur.Loc,
		}
	}
	p.advance()
	return nil
}

// parseOptionalLength parses the `(length)` a type may leave out. The
// parentheses are optional; what goes between them is not. Every engine answers
// 1064 for `int()`, `year()` and `datetime()` as surely as for `varchar()`, and
// for a length that never closes, so once `(` is there the rest is a length
// that must carry a number and end at `)`.
func (p *Parser) parseOptionalLength(dt *nodes.DataType) error {
	if p.cur.Type != '(' {
		return nil
	}
	return p.parseRequiredLength(dt)
}

// parseOptionalPrecision parses the `(precision[, scale])` the approximate and
// exact numeric types may leave out. As with a length, the parentheses are
// optional and their contents are not: the engines answer 1064 for
// `decimal(10,)`, `decimal(,2)` and `float(10,)`, so a precision that opens
// must carry a number, a scale that follows a comma must too, and the whole
// must close.
func (p *Parser) parseOptionalPrecision(dt *nodes.DataType) error {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	if p.cur.Type != tokICONST {
		return &ParseError{
			Message:  "expected a precision for " + dt.Name,
			Position: p.cur.Loc,
		}
	}
	dt.Length = int(p.cur.Ival)
	p.advance()
	if p.cur.Type == ',' {
		p.advance()
		if p.cur.Type != tokICONST {
			return &ParseError{
				Message:  "expected a scale for " + dt.Name,
				Position: p.cur.Loc,
			}
		}
		dt.Scale = int(p.cur.Ival)
		p.advance()
	}
	if p.cur.Type != ')' {
		return &ParseError{
			Message:  "expected ) after the precision for " + dt.Name,
			Position: p.cur.Loc,
		}
	}
	p.advance()
	return nil
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

// parseCharsetCollate parses optional CHARACTER SET / CHARSET clauses, the
// standalone BINARY modifier (e.g. CHAR(10) BINARY), and an optional COLLATE.
// CHARACTER SET and BINARY may appear in any order, but COLLATE comes LAST:
// TiDB accepts `CHARACTER SET x BINARY COLLATE y` and `BINARY COLLATE y` but
// REJECTS `COLLATE y BINARY` and `CHARACTER SET x COLLATE y BINARY` (BINARY
// after COLLATE). So charset/BINARY are consumed in a loop, then a single
// trailing COLLATE — a BINARY after COLLATE is left for the caller to reject.
//
// Only the character types (CHAR/VARCHAR/TEXT/ENUM/SET/NCHAR/...) call this, so
// BINARY is naturally not accepted on BLOB/VARBINARY, matching TiDB.
//
// charset/collation names can be `binary`, which is a reserved keyword, so we
// use parseKeywordOrIdent to accept keyword tokens as names.
func (p *Parser) parseCharsetCollate(dt *nodes.DataType) {
	// Each of CHARACTER SET and the standalone BINARY modifier may appear at most
	// once; a repeat is left for the caller to reject (TiDB rejects e.g.
	// `BINARY BINARY` and `CHARACTER SET x CHARACTER SET y`). Note that
	// `CHARACTER SET binary BINARY` is still valid — the first `binary` is the
	// charset name (consumed by parseKeywordOrIdent), the second is the modifier.
	charsetSeen := false
	binarySeen := false
charsetBinary:
	for {
		switch p.cur.Type {
		case kwCHARSET:
			if charsetSeen {
				break charsetBinary
			}
			charsetSeen = true
			p.advance()
			if p.isIdentToken() || p.cur.Type == kwBINARY {
				dt.Charset, _, _ = p.parseKeywordOrIdent()
			}
		case kwCHARACTER, kwCHAR:
			// CHARACTER SET x, and MySQL's two-word CHAR SET x, which is what
			// `long char set utf8` carries. Neither token is consumed unless
			// SET follows it: a trailing CHARACTER is not a charset clause
			// with the name left off, and TiDB rejects `long character`.
			if charsetSeen || p.peekNext().Type != kwSET {
				break charsetBinary
			}
			charsetSeen = true
			p.advance()
			p.advance()
			if p.isIdentToken() || p.cur.Type == kwBINARY {
				dt.Charset, _, _ = p.parseKeywordOrIdent()
			}
		case kwBINARY:
			if binarySeen {
				break charsetBinary
			}
			binarySeen = true
			p.advance()
			dt.Binary = true
		default:
			break charsetBinary
		}
	}

	// COLLATE collation_name — must come after any CHARACTER SET / BINARY.
	if _, ok := p.match(kwCOLLATE); ok {
		if p.isIdentToken() || p.cur.Type == kwBINARY {
			dt.Collate, _, _ = p.parseKeywordOrIdent()
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
		switch p.cur.Type {
		case tokSCONST:
			vals = append(vals, p.cur.Str)
			p.advance()
		case tokXCONST, tokBCONST:
			// A member may be written as a hex or bit literal, in either
			// spelling: `ENUM(0x61, 'b')` and `ENUM(x'61')` both name the
			// member `a` (https://dev.mysql.com/doc/refman/8.0/en/hexadecimal-literals.html,
			// .../bit-value-literals.html — in string context both are their
			// bytes). The member is the bytes, not the text that wrote them, so
			// it is decoded here: leaving `0x61` to a consumer would have the
			// catalog report `enum('0x61')` and a deparse write a different
			// table than the one the engine built.
			val, err := decodeByteLiteral(p.cur)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
			p.advance()
		default:
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

// decodeByteLiteral reads the bytes a hex or bit literal holds. The lexer keeps
// the token's own text, which carries the `0x`/`0b` prefix when the literal was
// written that way and not when it was written `x'..'` or `b'..'`, so the prefix
// comes off first.
//
// Hex digits are read in pairs from the right, so an odd count takes a leading
// zero: `0x6` is one byte, 0x06. Bits are read in groups of eight the same way.
// This is MySQL's reading of either literal in string context.
func decodeByteLiteral(tok Token) (string, error) {
	digits := tok.Str
	if len(digits) > 1 && digits[0] == '0' && (digits[1] == 'x' || digits[1] == 'X' || digits[1] == 'b' || digits[1] == 'B') {
		digits = digits[2:]
	}
	if digits == "" {
		return "", &ParseError{Message: "empty byte literal", Position: tok.Loc}
	}
	base, group := 16, 2
	if tok.Type == tokBCONST {
		base, group = 2, 8
	}
	// Pad on the left so the digits divide into whole bytes.
	if r := len(digits) % group; r != 0 {
		digits = strings.Repeat("0", group-r) + digits
	}
	out := make([]byte, 0, len(digits)/group)
	for i := 0; i < len(digits); i += group {
		b, err := strconv.ParseUint(digits[i:i+group], base, 8)
		if err != nil {
			return "", &ParseError{Message: "malformed byte literal", Position: tok.Loc}
		}
		out = append(out, byte(b))
	}
	return string(out), nil
}
