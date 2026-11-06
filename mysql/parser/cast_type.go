package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseCastDataType reads the type of CAST and CONVERT, which is its own short list and not
// a column's: MySQL 8.0 answers 1064 for CAST(a AS INT), BIGINT, TEXT, VARCHAR(3), BOOLEAN,
// TIMESTAMP and GEOMETRY, all of which a column takes.
//
// Ref: cast_type and opt_charset_with_opt_binary in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseCastDataType() (*nodes.DataType, error) {
	p.castSized = false
	// Completion: offer CAST/CONVERT type candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, &ParseError{Message: "collecting"}
	}

	dt := &nodes.DataType{Loc: nodes.Loc{Start: p.pos()}}
	switch p.cur.Type {
	case kwSIGNED, kwUNSIGNED:
		// The server's lexer reads INT4 as INT, FLOAT4 as FLOAT and FLOAT8 as DOUBLE
		// (sql/lex.h), so each alias goes where its type does. INT8 is BIGINT, which CAST
		// does not take.
		dt.Name = strings.ToUpper(p.advance().Str)
		p.match(kwINTEGER, kwINT, kwINT4)
	case kwCHAR, kwCHARACTER:
		dt.Name = "CHAR"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
		if err := p.parseCastCharOptions(dt); err != nil {
			return nil, err
		}
	case kwNCHAR:
		dt.Name, dt.Charset = "CHAR", "utf8mb3"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwNATIONAL:
		dt.Name, dt.Charset = "CHAR", "utf8mb3"
		p.advance()
		if _, ok := p.match(kwCHAR, kwCHARACTER); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
	case kwBINARY:
		dt.Name = "BINARY"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwDATE:
		dt.Name = "DATE"
		p.advance()
	case kwYEAR:
		dt.Name = "YEAR"
		p.advance()
	case kwTIME, kwDATETIME:
		dt.Name = strings.ToUpper(p.advance().Str)
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
	case kwDECIMAL, kwDEC:
		// A precision alone may be a decimal, and with a scale it is an integer: DECIMAL(1.5)
		// parses and DECIMAL(1.5,2) is 1064.
		dt.Name = "DECIMAL"
		p.advance()
		if p.cur.Type == '(' && p.peekNext().Type == tokFCONST {
			if err := p.parseCastSize(dt, true); err != nil {
				return nil, err
			}
			break
		}
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
	case kwFLOAT, kwFLOAT4:
		dt.Name = "FLOAT"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwDOUBLE, kwFLOAT8:
		dt.Name = "DOUBLE"
		p.advance()
		p.match(kwPRECISION)
	case kwREAL:
		dt.Name = "DOUBLE"
		p.advance()
	case kwJSON:
		dt.Name = "JSON"
		p.advance()
	case kwPOINT, kwLINESTRING, kwPOLYGON, kwMULTIPOINT, kwMULTILINESTRING, kwMULTIPOLYGON,
		kwGEOMETRYCOLLECTION, kwGEOMCOLLECTION:
		dt.Name = strings.ToUpper(p.advance().Str)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	dt.Loc.End = p.prev.End
	if err := castTypeRefused(dt); err != nil {
		return nil, err
	}
	return dt, nil
}

// maxCastLength is the largest length CHAR and BINARY take in a CAST.
const maxCastLength = 4294967295

// castTypeRefused holds a cast type to the limits MySQL 8.0 checks while it parses:
// 1439 for a CHAR or BINARY past maxCastLength; 1426 for a FLOAT past 53, a TIME or
// DATETIME past 6 and a DECIMAL past 65; 1425 for a DECIMAL scale past 30, and 1427 for
// one past its precision. CONVERT reads the same list.
func castTypeRefused(dt *nodes.DataType) error {
	refused := false
	switch dt.Name {
	case "CHAR", "BINARY":
		refused = dt.Length > maxCastLength
	case "FLOAT":
		refused = dt.Length > 53
	case "TIME", "DATETIME":
		refused = dt.Length > 6
	case "DECIMAL":
		refused = dt.Length > 65 || dt.Scale > 30 || dt.Scale > dt.Length
	}
	if refused {
		return &ParseError{Message: "Too big a size for CAST to " + dt.Name, Position: dt.Loc.Start}
	}
	return nil
}

// castArrayRefused holds CAST(... ARRAY) to the types MySQL 8.0 takes one of, answering
// 1235 while it parses for the rest: JSON, YEAR, DOUBLE, FLOAT, the geometry types, NCHAR,
// and a CHAR or BINARY that has no length. A CHAR takes no character set here but binary,
// which BYTE and CHARACTER SET binary both write.
func (p *Parser) castArrayRefused(dt *nodes.DataType) error {
	allowed := false
	switch dt.Name {
	case "SIGNED", "UNSIGNED", "DATE", "DATETIME", "TIME", "DECIMAL":
		allowed = true
	case "BINARY":
		allowed = p.castSized
	case "CHAR":
		allowed = p.castSized && (strings.EqualFold(dt.Charset, "binary") || (dt.Charset == "" && !dt.Binary))
	}
	if !allowed {
		return &ParseError{Message: "CAST to an array of " + dt.Name + " is not supported", Position: dt.Loc.Start}
	}
	return nil
}

// parseCastSize reads the (n) a cast type may carry, which is a plain number where a column's
// length may be followed by more. A decimal is read where the engine takes one and its whole
// part kept.
func (p *Parser) parseCastSize(dt *nodes.DataType, allowDecimal bool) error {
	if p.cur.Type != '(' {
		return nil
	}
	p.castSized = true
	p.advance()
	switch {
	case p.cur.Type == tokICONST:
		dt.Length = int(p.cur.Ival)
	case allowDecimal && p.cur.Type == tokFCONST && !strings.ContainsAny(p.cur.Str, "eE"):
		whole, _, _ := strings.Cut(p.cur.Str, ".")
		for _, digit := range whole {
			dt.Length = dt.Length*10 + int(digit-'0')
		}
	default:
		return p.syntaxErrorAtCur()
	}
	// A written length of 0 is a length, which JSON_VALUE's RETURNING keeps.
	dt.HasLength = true
	p.advance()
	_, err := p.expect(')')
	return err
}

// parseCastCharOptions reads what CHAR takes after its length: ASCII or UNICODE, with BINARY
// on either side; BYTE, alone; or a character set and BINARY in either order. ASCII is latin1,
// UNICODE is ucs2 and BYTE is binary.
func (p *Parser) parseCastCharOptions(dt *nodes.DataType) error {
	binary := false
	if _, ok := p.match(kwBINARY); ok {
		binary = true
		dt.Binary = true
	}
	switch p.cur.Type {
	case kwASCII, kwUNICODE:
		dt.Charset = "latin1"
		if p.cur.Type == kwUNICODE {
			dt.Charset = "ucs2"
		}
		p.advance()
		if !binary {
			if _, ok := p.match(kwBINARY); ok {
				dt.Binary = true
			}
		}
		return nil
	case kwBYTE:
		if binary {
			return p.syntaxErrorAtCur()
		}
		dt.Charset = "binary"
		p.advance()
		return nil
	}
	has, err := p.parseCastCharset(dt)
	if err != nil {
		return err
	}
	if has && !binary {
		if _, ok := p.match(kwBINARY); ok {
			dt.Binary = true
		}
	}
	return nil
}

// parseCastCharset reads CHARACTER SET, CHARSET or CHAR SET and the name, which is an
// identifier, a string or BINARY. The name is not optional: every engine answers 1064 for
// CAST(a AS CHAR CHARACTER SET), which the column reader lets go.
func (p *Parser) parseCastCharset(dt *nodes.DataType) (bool, error) {
	switch {
	case p.cur.Type == kwCHARSET:
		p.advance()
	case (p.cur.Type == kwCHARACTER || p.cur.Type == kwCHAR) && p.peekNext().Type == kwSET:
		p.advance()
		p.advance()
	default:
		return false, nil
	}
	switch {
	case p.cur.Type == tokSCONST:
		dt.Charset = p.advance().Str
	case p.isIdentToken() || p.cur.Type == kwBINARY:
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return false, err
		}
		dt.Charset = name
	default:
		return false, p.syntaxErrorAtCur()
	}
	return true, nil
}
