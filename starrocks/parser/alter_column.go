package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseAlterColumnDesc reads columnDesc (StarRocks.g4:396) for ALTER's
// ADD COLUMN and MODIFY COLUMN:
//
//	identifier type? charsetName? KEY? aggDesc? columnNullable?
//	(defaultDesc | AUTO_INCREMENT | generatedColumnDesc)? comment?
//
// The grammar makes the type optional, but the AST builder refuses a column
// of an ALTER clause without one (ALTER SYSTEM ADD COLUMN a NULL and
// ALTER SYSTEM MODIFY COLUMN a are rejects), so it is required here.
func (p *Parser) parseAlterColumnDesc() (*ast.ColumnDef, error) {
	start := p.cur.Loc
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	col := &ast.ColumnDef{Name: name, Loc: start}
	// CHAR SET is charsetName, not the CHAR type, so it leaves the type out.
	if p.cur.Kind == kwCHAR && p.peekNext().Kind == kwSET {
		return nil, p.syntaxErrorAtCur()
	}
	if col.Type, err = p.parseAlterType(); err != nil {
		return nil, err
	}
	// StarRocks.g4:402 charsetName: CHAR SET identifier | CHARSET identifier
	// | CHARACTER SET identifier
	switch {
	case p.cur.Kind == kwCHAR || p.alterAt("CHARACTER"):
		p.advance()
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		if _, _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
	case p.cur.Kind == kwCHARSET:
		p.advance()
		if _, _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwKEY {
		p.advance()
		col.IsKey = true
	}
	// StarRocks.g4:453 aggDesc
	switch {
	case p.cur.Kind == kwSUM, p.cur.Kind == kwMAX, p.cur.Kind == kwMIN, p.cur.Kind == kwREPLACE,
		p.cur.Kind == kwHLL_UNION, p.cur.Kind == kwBITMAP_UNION, p.cur.Kind == kwPERCENTILE_UNION,
		p.cur.Kind == kwREPLACE_IF_NOT_NULL:
		col.AggType = strings.ToUpper(p.advance().Str)
	case isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == int('(') &&
		p.cur.Kind != kwPROPERTIES && p.cur.Kind != kwCOMMENT:
		// StarRocks.g4:449 aggStateDesc:
		// identifier '(' typeWithNullable (',' typeWithNullable)* ')'
		aggStart := p.cur.Loc.Start
		p.advance()
		p.advance() // consume '('
		for {
			if _, err := p.parseAlterType(); err != nil {
				return nil, err
			}
			p.parseAlterColumnNullable()
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		col.AggType = p.input[aggStart-p.baseOffset : p.prev.Loc.End-p.baseOffset]
	}
	if n := p.parseAlterColumnNullable(); n != nil {
		col.Nullable = n
	}
	switch p.cur.Kind {
	case kwDEFAULT:
		// StarRocks.g4:408 defaultDesc:
		// DEFAULT (string | NULL | CURRENT_TIMESTAMP | '(' qualifiedName '(' ')' ')')
		p.advance()
		if col.Default, err = p.parseAlterDefaultDesc(); err != nil {
			return nil, err
		}
	case kwAUTO_INCREMENT:
		p.advance()
		col.AutoInc = true
	case kwAS:
		// StarRocks.g4:412 generatedColumnDesc: AS expression
		p.advance()
		if col.Generated, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance()
		tok, err := p.alterExpectString()
		if err != nil {
			return nil, err
		}
		col.Comment = tok.Str
	}
	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// parseAlterColumnNullable reads columnNullable? (StarRocks.g4:440):
// NULL | NOT NULL. A NOT not followed by NULL is left in place.
func (p *Parser) parseAlterColumnNullable() *bool {
	switch {
	case p.cur.Kind == kwNULL:
		p.advance()
		b := true
		return &b
	case p.cur.Kind == kwNOT && p.peekNext().Kind == kwNULL:
		p.advance()
		p.advance()
		b := false
		return &b
	}
	return nil
}

// parseAlterDefaultDesc reads what follows DEFAULT in defaultDesc
// (StarRocks.g4:408).
func (p *Parser) parseAlterDefaultDesc() (ast.Node, error) {
	switch p.cur.Kind {
	case tokString:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil
	case kwNULL:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitNull, Value: "NULL", Loc: tok.Loc}, nil
	case kwCURRENT_TIMESTAMP:
		tok := p.advance()
		return &ast.FuncCallExpr{
			Name: &ast.ObjectName{Parts: []string{"CURRENT_TIMESTAMP"}, Loc: tok.Loc},
			Loc:  tok.Loc,
		}, nil
	case int('('):
		open := p.advance()
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		return &ast.FuncCallExpr{Name: name, Loc: ast.Loc{Start: open.Loc.Start, End: closeTok.Loc.End}}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseAlterType reads type (StarRocks.g4:2784): baseType | decimalType |
// arrayType | structType | mapType.
func (p *Parser) parseAlterType() (*ast.TypeName, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	tok := p.cur
	typ := &ast.TypeName{Name: strings.ToUpper(tok.Str), Loc: tok.Loc}
	switch {
	case p.cur.Kind == kwARRAY:
		// StarRocks.g4:2792 arrayType: ARRAY '<' type '>'
		p.advance()
		if _, err := p.expect(int('<')); err != nil {
			return nil, err
		}
		elem, err := p.parseAlterType()
		if err != nil {
			return nil, err
		}
		typ.ElementType = elem
		if _, err := p.expectGT(); err != nil {
			return nil, err
		}
	case p.cur.Kind == kwMAP:
		// StarRocks.g4:2796 mapType: MAP '<' type ',' type '>'
		p.advance()
		if _, err := p.expect(int('<')); err != nil {
			return nil, err
		}
		key, err := p.parseAlterType()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(',')); err != nil {
			return nil, err
		}
		val, err := p.parseAlterType()
		if err != nil {
			return nil, err
		}
		typ.ElementType, typ.ValueType = key, val
		if _, err := p.expectGT(); err != nil {
			return nil, err
		}
	case p.cur.Kind == kwSTRUCT:
		// StarRocks.g4:2808 structType: STRUCT '<' subfieldDescs '>', with
		// subfieldDesc (StarRocks.g4:2800) (identifier | nestedFieldName) type.
		// The AST builder refuses two fields of one name (STRUCT < a BOOLEAN,
		// a BOOLEAN > is a reject).
		p.advance()
		seen := map[string]bool{}
		if _, err := p.expect(int('<')); err != nil {
			return nil, err
		}
		for {
			fstart := p.cur.Loc.Start
			fname, err := p.parseFieldPath()
			if err != nil {
				return nil, err
			}
			if seen[strings.ToLower(fname)] {
				return nil, &ParseError{Position: fstart, End: p.prev.Loc.End, Message: "duplicate field name " + fname + " in STRUCT"}
			}
			seen[strings.ToLower(fname)] = true
			ftyp, err := p.parseAlterType()
			if err != nil {
				return nil, err
			}
			typ.Fields = append(typ.Fields, &ast.StructField{
				Name: fname, Type: ftyp,
				Loc: ast.Loc{Start: fstart, End: ftyp.Loc.End},
			})
			if p.cur.Kind != int(',') {
				break
			}
			p.advance()
		}
		if _, err := p.expectGT(); err != nil {
			return nil, err
		}
	case p.alterAt("DECIMAL"), p.alterAt("DECIMALV2"), p.alterAt("DECIMAL32"), p.alterAt("DECIMAL64"),
		p.alterAt("DECIMAL128"), p.alterAt("NUMERIC"), p.alterAt("NUMBER"):
		// StarRocks.g4:2845 decimalType:
		// (DECIMAL | ...) ('(' precision=INTEGER_VALUE (',' scale=INTEGER_VALUE)? ')')?
		p.advance()
		if p.cur.Kind == int('(') {
			p.advance()
			n, err := p.alterExpectInt()
			if err != nil {
				return nil, err
			}
			typ.Params = append(typ.Params, int(n.Ival))
			if p.cur.Kind == int(',') {
				p.advance()
				if n, err = p.alterExpectInt(); err != nil {
					return nil, err
				}
				typ.Params = append(typ.Params, int(n.Ival))
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
		}
	case p.alterAt("SIGNED"), p.alterAt("UNSIGNED"):
		// baseType (StarRocks.g4:2816): SIGNED INT? | SIGNED INTEGER? | UNSIGNED INT? | UNSIGNED INTEGER?
		p.advance()
		if p.cur.Kind == kwINT || p.cur.Kind == kwINTEGER {
			typ.Name += " " + strings.ToUpper(p.advance().Str)
		}
	case alterIsParamBaseType(tok):
		// StarRocks.g4:2816 baseType with typeParameter? ('(' INTEGER_VALUE ')')
		p.advance()
		if p.cur.Kind == int('(') {
			p.advance()
			n, err := p.alterExpectInt()
			if err != nil {
				return nil, err
			}
			typ.Params = []int{int(n.Ival)}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
		}
	case alterIsPlainBaseType(tok):
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}
	typ.Loc.End = p.prev.Loc.End
	return typ, nil
}

// alterIsParamBaseType reports whether t is a baseType (StarRocks.g4:2816)
// that takes a typeParameter.
func alterIsParamBaseType(t Token) bool {
	for _, w := range []string{"TINYINT", "SMALLINT", "INT", "INTEGER", "BIGINT", "LARGEINT", "CHAR", "VARCHAR", "VARBINARY", "BINARY"} {
		if alterWordTok(t, w) {
			return true
		}
	}
	return false
}

// alterIsPlainBaseType reports whether t is a baseType (StarRocks.g4:2816)
// that takes no parameter.
func alterIsPlainBaseType(t Token) bool {
	for _, w := range []string{"BOOLEAN", "FLOAT", "DOUBLE", "DATE", "DATETIME", "TIME", "STRING", "TEXT", "BITMAP", "HLL", "PERCENTILE", "JSON"} {
		if alterWordTok(t, w) {
			return true
		}
	}
	return false
}
