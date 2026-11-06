package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseCastDataType reads the type of CAST and CONVERT, which is its own list and not a
// column's: MariaDB 11.8.8 answers 1064 for CAST(a AS BIGINT), TEXT, BOOLEAN, TIMESTAMP, YEAR,
// JSON and REAL, all of which a column takes, and takes INT, VARCHAR(3) and the types it reads
// by name. The list is what the server answered. An unknown name, which it refuses with 4161
// while it parses, is refused here as it is on a column.
func (p *Parser) parseCastDataType() (*nodes.DataType, error) {
	// Completion: offer CAST/CONVERT type candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, &ParseError{Message: "collecting"}
	}

	dt := &nodes.DataType{Loc: nodes.Loc{Start: p.pos()}}
	switch p.cur.Type {
	case kwSIGNED, kwUNSIGNED:
		dt.Name = strings.ToUpper(p.advance().Str)
		p.match(kwINTEGER, kwINT)
	case kwINT, kwINTEGER:
		dt.Name = "INT"
		p.advance()
	case kwCHAR, kwCHARACTER:
		dt.Name = "CHAR"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
		if err := p.parseCastCharOptions(dt); err != nil {
			return nil, err
		}
	case kwVARCHAR, kwVARCHARACTER:
		dt.Name = "VARCHAR"
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
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
	case kwBINARY:
		dt.Name = "BINARY"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwDATE:
		dt.Name = "DATE"
		p.advance()
	case kwTIME, kwDATETIME:
		dt.Name = strings.ToUpper(p.advance().Str)
		if err := p.parseCastSize(dt, true); err != nil {
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
	case kwFLOAT:
		dt.Name = "FLOAT"
		p.advance()
	case kwDOUBLE:
		// DOUBLE takes a precision and a scale together, or neither: 11.8.8 answers 1064
		// for DOUBLE(3).
		dt.Name = "DOUBLE"
		p.advance()
		if p.cur.Type == '(' {
			hasScale, err := p.parsePrecisionAndScale(dt)
			if err != nil {
				return nil, err
			}
			if !hasScale {
				return nil, p.syntaxErrorAtCur()
			}
		}
	case kwGEOMETRY, kwPOINT, kwLINESTRING, kwPOLYGON, kwMULTIPOINT, kwMULTILINESTRING, kwMULTIPOLYGON,
		kwGEOMETRYCOLLECTION, kwVECTOR:
		// These read as type names and then have no CAST: 4162 while it parses, and 1064 with
		// a size after them.
		return nil, &ParseError{Message: "Operator does not exist: 'CAST(expr AS " + p.cur.Str + ")'", Position: p.cur.Loc}
	default:
		// The types MariaDB reads by name are identifiers, so a quoted one is the same type.
		// They take no size here, where a column's does: 11.8.8 answers 1064 for
		// CAST(a AS inet6(3)) and CONVERT(a, uuid(3)).
		if p.isNamedScalarType() {
			dt.Name = strings.ToUpper(p.advance().Str)
			dt.Loc.End = p.prev.End
			return dt, nil
		}
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

// castTypeRefused holds a cast type to the limits MariaDB 11.8.8 checks while it parses:
// 1439 for a CHAR or BINARY past maxCastLength; 1426 for a TIME or DATETIME past 6 and a
// DECIMAL past 65; 1425 for a DECIMAL scale past 38, and 1427 for one past its precision.
func castTypeRefused(dt *nodes.DataType) error {
	refused := false
	switch dt.Name {
	case "CHAR", "BINARY":
		refused = dt.Length > maxCastLength
	case "TIME", "DATETIME":
		refused = dt.Length > 6
	case "DECIMAL":
		refused = dt.Length > 65 || dt.Scale > 38 || dt.Scale > dt.Length
	}
	if refused {
		return &ParseError{Message: "Too big a size for CAST to " + dt.Name, Position: dt.Loc.Start}
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
	p.advance()
	_, err := p.expect(')')
	return err
}

// parseCastCharOptions reads what CHAR takes after its length: ASCII or UNICODE, with BINARY
// on either side; BYTE, alone; or a character set and BINARY in either order. ASCII is latin1,
// UNICODE is ucs2 and BYTE is binary. MariaDB also closes the list with
// COLLATE, which MySQL and TiDB answer 1064 for.
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
		return p.parseCastCollate(dt)
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
	return p.parseCastCollate(dt)
}

// parseCastCollate reads the COLLATE that may close a character type, which no BINARY may
// stand before: 11.8.8 answers 1064 for CHAR BINARY COLLATE c and takes CHAR ASCII COLLATE c,
// refusing that later with 1253.
func (p *Parser) parseCastCollate(dt *nodes.DataType) error {
	if dt.Binary {
		return nil
	}
	if _, ok := p.match(kwCOLLATE); ok {
		if !p.isIdentToken() && p.cur.Type != kwBINARY {
			return p.syntaxErrorAtCur()
		}
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return err
		}
		dt.Collate = name
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
