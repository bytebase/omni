package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseCastDataType reads the type of CAST and CONVERT, which is its own short list and not
// a column's: TiDB v8.5.5 answers 1064 for CAST(a AS INT), BIGINT, TEXT, VARCHAR(3), BOOLEAN,
// TIMESTAMP, NCHAR and POINT, all of which a column takes.
//
// Ref: CastType and OptBinary in TiDB v8.5.5 pkg/parser/parser.y.
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
	case kwCHAR, kwCHARACTER:
		dt.Name = "CHAR"
		p.advance()
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
		if err := p.parseCastCharOptions(dt); err != nil {
			return nil, err
		}
	case kwBINARY:
		dt.Name = "BINARY"
		p.advance()
		if err := p.parseCastSize(dt, false); err != nil {
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
		dt.Name = "DECIMAL"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
	case kwFLOAT:
		dt.Name = "FLOAT"
		p.advance()
		if err := p.parseOptionalPrecision(dt); err != nil {
			return nil, err
		}
		// CastType's action answers 1426, ErrTooBigPrecision, for a FLOAT past 53, which a
		// column type takes no such check for.
		if dt.Length >= 54 {
			return nil, &ParseError{Message: "Too big precision specified for CAST. Maximum is 53.", Position: dt.Loc.Start}
		}
	case kwDOUBLE, kwREAL:
		dt.Name = "DOUBLE"
		p.advance()
	case kwJSON:
		dt.Name = "JSON"
		p.advance()
	case kwVECTOR:
		// VECTOR [<FLOAT>] [(dimension)], as a column writes it.
		dt.Name = "VECTOR"
		p.advance()
		if p.cur.Type == '<' {
			p.advance()
			if _, err := p.expect(kwFLOAT); err != nil {
				return nil, err
			}
			if _, err := p.expect('>'); err != nil {
				return nil, err
			}
		}
		if err := p.parseCastSize(dt, false); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	dt.Loc.End = p.prev.End
	return dt, nil
}

// parseCastSize reads the (n) a cast type may carry, which is a plain number where a column's
// length may be followed by more. A decimal is read where the engine takes one and its whole
// part kept.
func (p *Parser) parseCastSize(dt *nodes.DataType, allowDecimal bool) error {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	// The whole part of a decimal is its length, held to the largest int64 as an
	// integer's is, so a run of digits too wide for one is refused as too long.
	if p.cur.Type != tokICONST && !(allowDecimal && p.cur.Type == tokFCONST && !strings.ContainsAny(p.cur.Str, "eE")) {
		return p.syntaxErrorAtCur()
	}
	dt.Length = int(p.cur.Ival)
	p.advance()
	_, err := p.expect(')')
	return err
}

// parseCastCharOptions reads what CHAR takes after its length: a character set, BINARY, or
// the two in either order. v8.5.5 answers 1064 for COLLATE here, and for ASCII, UNICODE and
// BYTE, which MySQL has.
func (p *Parser) parseCastCharOptions(dt *nodes.DataType) error {
	if _, ok := p.match(kwBINARY); ok {
		dt.Binary = true
		_, err := p.parseCastCharset(dt)
		return err
	}
	has, err := p.parseCastCharset(dt)
	if err != nil {
		return err
	}
	if has {
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
