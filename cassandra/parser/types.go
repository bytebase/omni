package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// typeKind is what parseDataType read, for the FROZEN<> check.
type typeKind int

const (
	typeNative typeKind = iota
	typeCollection
	typeTuple
	typeVector
	typeUser
)

// reservedTypeNames are the names Parser.g refuses for a user-defined type
// (Parser.g:30, checked by non_type_ident at Parser.g:1940 against the
// name as written).
var reservedTypeNames = map[string]bool{
	"byte": true, "complex": true, "enum": true, "date": true,
	"interval": true, "macaddr": true, "bitstring": true,
}

// isNativeType reports whether a token is a native_type (Parser.g:1873).
func isNativeType(typ int) bool {
	switch typ {
	case tokASCII, tokBIGINT, tokBLOB, tokBOOLEAN, tokCOUNTER, tokDECIMAL,
		tokDOUBLE, tokDURATION, tokFLOATKW, tokINET, tokINT, tokSMALLINT,
		tokTEXT, tokTIMESTAMP, tokTINYINT, tokUUID_KW, tokVARCHAR, tokVARINT,
		tokTIMEUUID, tokDATE, tokTIME:
		return true
	}
	return false
}

// parseDataType parses comparatorType (Parser.g:1847): a native type, a
// collection, a tuple, a vector, FROZEN<type> or a user-defined type,
// [ks.]name. A custom type given as a string is left out: Cassandra loads
// the class it names while parsing, which omni cannot check.
func (p *Parser) parseDataType() (*ast.DataType, error) {
	dt, _, err := p.parseDataTypeKind()
	return dt, err
}

func (p *Parser) parseDataTypeKind() (*ast.DataType, typeKind, error) {
	if err := p.enterNesting(); err != nil {
		return nil, 0, err
	}
	defer p.leaveNesting()
	start := p.curLoc()
	tok := p.cur
	next := p.peekNext()

	// A native type, unless it names the keyspace of a user type (ascii.t).
	if isNativeType(tok.Type) && next.Type != tokDOT {
		p.advance()
		return &ast.DataType{Name: tokIdent(tok), Loc: p.makeLoc(start)}, typeNative, nil
	}

	if next.Type == tokLT {
		var params, kind = 0, typeCollection
		switch tok.Type {
		case tokMAP:
			params = 2
		case tokLIST, tokSET, tokFROZEN:
			params = 1
		case tokTUPLE:
			params, kind = -1, typeTuple
		case tokVECTOR:
			kind = typeVector
		default:
			return nil, 0, p.errorf("expected data type name, got %s", p.tokenDesc())
		}
		p.advance() // the type keyword
		p.advance() // <
		dt := &ast.DataType{Name: tokIdent(tok)}
		if kind == typeVector {
			// VECTOR '<' comparatorType ',' INTEGER '>' (Parser.g:1916)
			elem, err := p.parseDataType()
			if err != nil {
				return nil, 0, err
			}
			dt.TypeParams = append(dt.TypeParams, elem)
			if _, err := p.expect(tokCOMMA); err != nil {
				return nil, 0, err
			}
			if p.cur.Type != tokINTEGER {
				return nil, 0, p.errorf("expected integer dimension for VECTOR, got %s", p.tokenDesc())
			}
			dt.Dimension = &ast.IntegerLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
			p.advance()
		} else {
			var paramKind typeKind
			for {
				paramTok := p.cur
				var param *ast.DataType
				var err error
				param, paramKind, err = p.parseDataTypeKind()
				if err != nil {
					return nil, 0, err
				}
				// freeze() refuses a native type (Parser.g:1854,
				// CQL3Type.Raw.freeze).
				if tok.Type == tokFROZEN && paramKind == typeNative {
					return nil, 0, p.errorAt(paramTok, "frozen<> is only allowed on collections, tuples, and user-defined types (got %s)", paramTok.Str)
				}
				dt.TypeParams = append(dt.TypeParams, param)
				if len(dt.TypeParams) == params || !p.match(tokCOMMA) {
					break
				}
			}
			if params > 0 && len(dt.TypeParams) != params {
				return nil, 0, p.errorf("expected ',', got %s", p.tokenDesc())
			}
			if tok.Type == tokFROZEN {
				kind = paramKind // FROZEN<x> is of the kind of x
			}
		}
		if _, err := p.expect(tokGT); err != nil {
			return nil, 0, err
		}
		dt.Loc = p.makeLoc(start)
		return dt, kind, nil
	}

	if !isCQLName(tok) {
		return nil, 0, p.errorf("expected data type name, got %s", p.tokenDesc())
	}
	ks, name, err := p.parseUserTypeName()
	if err != nil {
		return nil, 0, err
	}
	return &ast.DataType{Name: name, Keyspace: ks, Loc: p.makeLoc(start)}, typeUser, nil
}

// parseUserTypeName parses userTypeName (Parser.g:1489):
// [noncol_ident '.'] non_type_ident.
func (p *Parser) parseUserTypeName() (*ast.Identifier, *ast.Identifier, error) {
	var ks, name *ast.Identifier
	var err error
	if p.peekNext().Type == tokDOT {
		// noncol_ident '.' (Parser.g:1467)
		if ks, err = p.parseNameIf(isCQLName); err != nil {
			return nil, nil, err
		}
		p.advance() // .
	}
	if name, err = p.parseNameIf(isNonTypeIdent); err != nil {
		return nil, nil, err
	}
	if !name.Quoted && reservedTypeNames[name.Name] {
		return nil, nil, p.errorAt(Token{Type: tokIDENT, Str: name.Name, Loc: name.Loc.Start, End: name.Loc.End}, "Invalid (reserved) user type name %s", name.Name)
	}
	return ks, name, nil
}

// tokIdent is a token as an unquoted Identifier.
func tokIdent(tok Token) *ast.Identifier {
	return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}
}
