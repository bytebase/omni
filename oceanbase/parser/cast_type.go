package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseCastDataType reads the type of CAST and CONVERT, which is its own list and not a
// column's: 4.3.5.6 answers 1064 for CAST(a AS TINYINT), TEXT, VARCHAR(3), BOOLEAN, TIMESTAMP
// and REAL, all of which a column takes, and takes INT, BIGINT, NUMBER, FIXED and DATE(3).
//
// Ref: cast_data_type in sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
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
	case kwBIGINT:
		dt.Name = "BIGINT"
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
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwBINARY:
		dt.Name = "BINARY"
		p.advance()
		if err := p.parseCastSize(dt, true); err != nil {
			return nil, err
		}
	case kwDATE, kwTIME, kwDATETIME, kwYEAR:
		dt.Name = strings.ToUpper(p.advance().Str)
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
	case kwDECIMAL, kwDEC, kwNUMERIC, kwFIXED, kwNUMBER:
		// All one type under other names, as the column reader has FIXED: keeping the word
		// would print a cast with no precision, since only DECIMAL's writer prints one.
		dt.Name = "DECIMAL"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
	case kwFLOAT:
		dt.Name = "FLOAT"
		p.advance()
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
	case kwDOUBLE:
		dt.Name = "DOUBLE"
		p.advance()
	case kwJSON:
		dt.Name = "JSON"
		p.advance()
	case kwARRAY:
		// ARRAY(type) and ARRAY<type>, which nest in either form, and whose element is any
		// column's type.
		dt.Name = "ARRAY"
		p.advance()
		p.castShiftHalf = false
		element, err := p.parseCastArrayElement()
		if err != nil {
			return nil, err
		}
		dt.ArrayElem = element
	case kwPOINT, kwLINESTRING, kwPOLYGON, kwMULTIPOINT, kwMULTILINESTRING, kwMULTIPOLYGON,
		kwGEOMETRYCOLLECTION, kwGEOMCOLLECTION:
		dt.Name = strings.ToUpper(p.advance().Str)
	default:
		if p.isBareWord("ROARINGBITMAP") {
			dt.Name = "ROARINGBITMAP"
			p.advance()
			dt.Loc.End = p.prev.End
			return dt, nil
		}
		return nil, p.syntaxErrorAtCur()
	}
	dt.Loc.End = p.prev.End
	return dt, nil
}

// parseCastArrayElement reads the element of a cast to ARRAY, from the opening ( or < on.
// Two closing brackets lex as one >> token, so the inner of the two closes half of it and
// the outer the rest: castShiftHalf says the first half is gone.
func (p *Parser) parseCastArrayElement() (*nodes.DataType, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	angle := p.cur.Type == '<'
	if !angle {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
	} else {
		p.advance()
	}
	var element *nodes.DataType
	if p.cur.Type == kwARRAY {
		start := p.pos()
		p.advance()
		inner, err := p.parseCastArrayElement()
		if err != nil {
			return nil, err
		}
		// A nested ARRAY's last token may be the >> it shares with the one around it.
		end := p.prev.End
		if p.castShiftHalf {
			end = p.cur.End
		}
		element = &nodes.DataType{Loc: nodes.Loc{Start: start, End: end}, Name: "ARRAY", ArrayElem: inner}
	} else {
		var err error
		if element, err = p.parseDataType(); err != nil {
			return nil, err
		}
		// An element takes no COLLATE: 4.3.5.6 answers 1064 for ARRAY(VARCHAR(3) COLLATE x).
		if element.Collate != "" {
			return nil, &ParseError{Message: "an ARRAY element takes no COLLATE", Position: element.Loc.Start}
		}
	}
	if !angle {
		_, err := p.expect(')')
		return element, err
	}
	switch {
	case p.cur.Type == tokShiftRight && !p.castShiftHalf:
		p.castShiftHalf = true
		return element, nil
	case p.cur.Type == tokShiftRight:
		p.castShiftHalf = false
		p.advance()
		return element, nil
	}
	_, err := p.expect('>')
	return element, err
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
	// A written length of 0 is a length, which JSON_VALUE's RETURNING keeps.
	dt.HasLength = true
	p.advance()
	_, err := p.expect(')')
	return err
}

// parseCastCharOptions reads what CHAR takes after its length: BINARY or a character set,
// and not both, nor ASCII, UNICODE or BYTE, which MySQL has.
func (p *Parser) parseCastCharOptions(dt *nodes.DataType) error {
	if _, ok := p.match(kwBINARY); ok {
		dt.Binary = true
		return nil
	}
	_, err := p.parseCastCharset(dt)
	return err
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
