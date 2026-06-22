package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// parseDataType parses a CQL data type: name or name<params>.
func (p *Parser) parseDataType() (*ast.DataType, error) {
	start := p.curLoc()
	name, err := p.parseDataTypeName()
	if err != nil {
		return nil, err
	}
	dt := &ast.DataType{Name: name, Loc: p.makeLoc(start)}

	isVector := strings.EqualFold(name.Name, "vector")
	if p.cur.Type == tokLT {
		p.advance() // <
		for {
			if isVector && p.cur.Type == tokINTEGER {
				dim := &ast.IntegerLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
				p.advance()
				dt.Dimension = dim
			} else {
				param, err := p.parseDataType()
				if err != nil {
					return nil, err
				}
				dt.TypeParams = append(dt.TypeParams, param)
			}
			if !p.match(tokCOMMA) {
				break
			}
		}
		if _, err := p.expect(tokGT); err != nil {
			return nil, err
		}
		dt.Loc = p.makeLoc(start)
	}
	return dt, nil
}

// parseDataTypeName parses a type name which can be a keyword or identifier.
func (p *Parser) parseDataTypeName() (*ast.Identifier, error) {
	tok := p.cur
	switch tok.Type {
	case tokIDENT, tokQUOTED:
		p.advance()
		return &ast.Identifier{Name: tok.Str, Quoted: tok.Type == tokQUOTED, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tokASCII, tokBIGINT, tokBLOB, tokBOOLEAN, tokCOUNTER, tokDATE, tokDECIMAL,
		tokDOUBLE, tokDURATION, tokFLOATKW, tokFROZEN, tokINET, tokINT, tokLIST,
		tokMAP, tokSET, tokSMALLINT, tokTEXT, tokTIME, tokTIMESTAMP, tokTIMEUUID,
		tokTINYINT, tokTUPLE, tokVARCHAR, tokVARINT, tokUUID_KW, tokVECTOR:
		p.advance()
		return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	default:
		return nil, p.errorf("expected data type name, got %s", p.tokenDesc())
	}
}
