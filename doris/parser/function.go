package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// Doris's ANTLR grammar derives CREATE FUNCTION, DROP FUNCTION and SHOW
// CREATE FUNCTION from unsupportedStatement, so Doris parses them with both
// of its parsers and takes one only if both do (testdata/README.md). The
// readers here take what both grammars take: DorisParser.g4 and the CUP
// parser's DorisLegacyParser.g4, cited by line.

// parseCreateFunction parses (DorisParser.g4:799-807)
//
//	CREATE (GLOBAL | SESSION | LOCAL)? (TABLES | AGGREGATE)? FUNCTION (IF NOT EXISTS)?
//	    functionIdentifier LEFT_PAREN functionArguments? RIGHT_PAREN
//	    RETURNS dataType (INTERMEDIATE dataType)? propertyClause?
//	CREATE (GLOBAL | SESSION | LOCAL)? ALIAS FUNCTION (IF NOT EXISTS)?
//	    functionIdentifier LEFT_PAREN functionArguments? RIGHT_PAREN
//	    WITH PARAMETER LEFT_PAREN identifierSeq? RIGHT_PAREN AS expression
//
// and the CUP parser's create_stmt forms (DorisLegacyParser.g4:420-422),
// whose properties are 'k' = 'v' pairs. On entry, CREATE has been consumed
// and cur is the word after it.
func (p *Parser) parseCreateFunction(createLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CreateFunctionStmt{Loc: createLoc}
	start := p.cur.Loc.Start

	switch p.cur.Kind {
	case kwGLOBAL:
		stmt.Global = true
		p.advance()
	case kwSESSION, kwLOCAL:
		p.advance()
	}
	switch p.cur.Kind {
	case kwTABLES:
		stmt.Tables = true
		p.advance()
	case kwAGGREGATE:
		stmt.Aggregate = true
		p.advance()
	case kwALIAS:
		stmt.Alias = true
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
	name, err := p.parseFunctionIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if stmt.Args, stmt.Variadic, err = p.parseFunctionArguments(); err != nil {
		return nil, err
	}

	if stmt.Alias {
		if err := p.parseAliasFunctionBody(stmt); err != nil {
			return nil, err
		}
	} else {
		if _, err := p.expect(kwRETURNS); err != nil {
			return nil, err
		}
		if stmt.Returns, err = p.parseLegacyDataType(); err != nil {
			return nil, err
		}
		if p.cur.Kind == kwINTERMEDIATE {
			p.advance()
			if stmt.Intermediate, err = p.parseLegacyDataType(); err != nil {
				return nil, err
			}
		}
		if p.cur.Kind == kwPROPERTIES {
			if stmt.Properties, err = p.parseProperties(); err != nil {
				return nil, err
			}
		}
	}

	stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : p.prev.Loc.End-p.baseOffset])
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAliasFunctionBody reads the tail of CREATE ALIAS FUNCTION:
//
//	WITH PARAMETER LEFT_PAREN identifierSeq? RIGHT_PAREN AS expression
//
// (DorisParser.g4:806-807; the CUP parser's opt_ident_list and expr).
func (p *Parser) parseAliasFunctionBody(stmt *ast.CreateFunctionStmt) error {
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	if _, err := p.expect(kwPARAMETER); err != nil {
		return err
	}
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	if p.cur.Kind != int(')') {
		for {
			name, _, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			stmt.Parameters = append(stmt.Parameters, name)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	body, err := p.parseExpr()
	if err != nil {
		return err
	}
	stmt.Body = body
	return nil
}

// parseFunctionIdentifier reads a function's name as both grammars take it:
// DorisParser.g4's functionIdentifier, (identifier DOT)?
// functionNameIdentifier (:1609), and the CUP parser's function_name,
// type_function_name or ident '.' type_function_name, whose
// type_function_name is an ident, LEFT or RIGHT (DorisLegacyParser.g4:2206,
// :2212, :3105). So the name may be the reserved LEFT or RIGHT; the other
// reserved words functionNameIdentifier lists the CUP parser refuses.
func (p *Parser) parseFunctionIdentifier() (*ast.ObjectName, error) {
	part := func() (string, ast.Loc, error) {
		if p.cur.Kind == kwLEFT || p.cur.Kind == kwRIGHT {
			tok := p.advance()
			return tok.Str, tok.Loc, nil
		}
		return p.parseIdentifier()
	}
	if isIdentifierToken(p.cur.Kind) && p.peekNext().Kind == int('.') {
		db, dbLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		p.advance() // consume '.'
		fn, fnLoc, err := part()
		if err != nil {
			return nil, err
		}
		return &ast.ObjectName{Parts: []string{db, fn}, Loc: ast.Loc{Start: dbLoc.Start, End: fnLoc.End}}, nil
	}
	fn, loc, err := part()
	if err != nil {
		return nil, err
	}
	return &ast.ObjectName{Parts: []string{fn}, Loc: loc}, nil
}

// parseFunctionArguments reads a function signature's argument list, from
// its '(' through its ')': what DorisParser.g4's
//
//	LEFT_PAREN functionArguments? RIGHT_PAREN
//	functionArguments : functionArgument (COMMA functionArgument)*
//	functionArgument  : DOTDOTDOT | dataType
//
// (:872-879) and the CUP parser's func_args_def (DorisLegacyParser.g4:2624:
// empty, '...', or types with an optional trailing ', ...') both take. So
// "..." comes last, once.
func (p *Parser) parseFunctionArguments() ([]*ast.TypeName, bool, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, false, err
	}
	var args []*ast.TypeName
	variadic := false
	if p.cur.Kind != int(')') {
		for {
			if p.cur.Kind == tokDotDotDot {
				p.advance()
				variadic = true
				break
			}
			arg, err := p.parseLegacyDataType()
			if err != nil {
				return nil, false, err
			}
			args = append(args, arg)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}
	return args, variadic, nil
}

// parseDropFunction parses (DorisParser.g4:730-731)
//
//	DROP (GLOBAL | SESSION | LOCAL)? FUNCTION (IF EXISTS)?
//	    functionIdentifier LEFT_PAREN functionArguments? RIGHT_PAREN
//
// and the CUP parser's drop_stmt form (DorisLegacyParser.g4:962). On entry,
// DROP has been consumed.
func (p *Parser) parseDropFunction(dropLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.DropFunctionStmt{Loc: dropLoc}
	start := p.cur.Loc.Start

	switch p.cur.Kind {
	case kwGLOBAL:
		stmt.Global = true
		p.advance()
	case kwSESSION, kwLOCAL:
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
	name, err := p.parseFunctionIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if stmt.Args, stmt.Variadic, err = p.parseFunctionArguments(); err != nil {
		return nil, err
	}

	stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : p.prev.Loc.End-p.baseOffset])
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
