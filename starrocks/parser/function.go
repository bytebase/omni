package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// isCreateFunctionStart reports whether cur, followed by next, begins the
// part of createFunctionStatement after CREATE orReplace: GLOBAL,
// AGGREGATE, FUNCTION, or TABLE FUNCTION. FUNCTION is reserved, so no table
// takes its name and CREATE TABLE FUNCTION is never a CREATE TABLE.
func isCreateFunctionStart(cur, next TokenKind) bool {
	switch cur {
	case kwGLOBAL, kwAGGREGATE, kwFUNCTION:
		return true
	case kwTABLE:
		return next == kwFUNCTION
	}
	return false
}

// parseCreateFunction parses createFunctionStatement (StarRocks.g4:1414):
//
//	CREATE orReplace GLOBAL? functionType=(TABLE | AGGREGATE)? FUNCTION ifNotExists
//	    qualifiedName '(' typeList ')' RETURNS returnType=type
//	    (properties | inlineProperties)?? inlineFunction?
//
// On entry, CREATE has been consumed; OR REPLACE too, when orReplace is set.
//
// inlineFunction (AS ATTACHMENT, a $$-quoted body) is not read: StarRocks
// 3.4.10 refuses CREATE FUNCTION a() RETURNS BOOLEAN AS $$a$$
// (testdata/compat.tsv), so it is left for the statement's end to reject.
func (p *Parser) parseCreateFunction(createLoc ast.Loc, orReplace bool) (ast.Node, error) {
	stmt := &ast.CreateFunctionStmt{Loc: createLoc, OrReplace: orReplace}
	start := p.cur.Loc.Start

	if p.cur.Kind == kwGLOBAL {
		stmt.Global = true
		p.advance()
	}
	switch p.cur.Kind {
	case kwTABLE:
		stmt.Table = true
		p.advance()
	case kwAGGREGATE:
		stmt.Aggregate = true
		p.advance()
	}
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if stmt.Args, stmt.Variadic, err = p.parseFunctionTypeList(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	if stmt.Returns, err = p.parseDataType(); err != nil {
		return nil, err
	}
	switch {
	case p.cur.Kind == kwPROPERTIES:
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	case isIdentifierToken(p.cur.Kind):
		if stmt.Properties, err = p.parseInlineProperties(); err != nil {
			return nil, err
		}
	}

	stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : p.prev.Loc.End-p.baseOffset])
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseInlineProperties reads inlineProperties (StarRocks.g4:2719):
//
//	inlineProperties : inlineProperty (inlineProperty)*
//	inlineProperty   : key=identifier '=' value=string
func (p *Parser) parseInlineProperties() ([]*ast.Property, error) {
	var props []*ast.Property
	for {
		key, keyLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.advance()
		props = append(props, &ast.Property{
			Key:   key,
			Value: val.Str,
			Loc:   ast.Loc{Start: keyLoc.Start, End: val.Loc.End},
		})
		if !isIdentifierToken(p.cur.Kind) {
			return props, nil
		}
	}
}

// parseFunctionTypeList reads a function signature's argument list, from its
// '(' through its ')' (StarRocks.g4:1421):
//
//	typeList : type? (',' type)* (',' DOTDOTDOT)?
//
// So the first type may be left out, "(, INT)", and "..." follows a comma.
func (p *Parser) parseFunctionTypeList() ([]*ast.TypeName, bool, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, false, err
	}
	var args []*ast.TypeName
	if p.cur.Kind != int(',') && p.cur.Kind != int(')') {
		arg, err := p.parseDataType()
		if err != nil {
			return nil, false, err
		}
		args = append(args, arg)
	}
	variadic := false
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		if p.cur.Kind == tokDotDotDot {
			p.advance()
			variadic = true
			break
		}
		arg, err := p.parseDataType()
		if err != nil {
			return nil, false, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}
	return args, variadic, nil
}

// parseDropFunction parses dropFunctionStatement (StarRocks.g4:1410):
//
//	DROP GLOBAL? FUNCTION (IF EXISTS)? qualifiedName '(' typeList ')'
//
// On entry, DROP has been consumed.
func (p *Parser) parseDropFunction(dropLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.DropFunctionStmt{Loc: dropLoc}
	start := p.cur.Loc.Start

	if p.cur.Kind == kwGLOBAL {
		stmt.Global = true
		p.advance()
	}
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if stmt.Args, stmt.Variadic, err = p.parseFunctionTypeList(); err != nil {
		return nil, err
	}

	stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : p.prev.Loc.End-p.baseOffset])
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
