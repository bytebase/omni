package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseCreateFunction parses Doris CREATE [GLOBAL] [AGGREGATE | ALIAS] FUNCTION ...
// On entry, CREATE has been consumed. cur is one of:
//   - kwGLOBAL (for CREATE GLOBAL ALIAS FUNCTION)
//   - kwAGGREGATE (for CREATE AGGREGATE FUNCTION)
//   - kwALIAS (for CREATE ALIAS FUNCTION)
//   - kwFUNCTION (for CREATE FUNCTION)
//
// The body is captured as raw text (functions in Doris contain JAR refs,
// expressions, or other forms that vary widely).
func (p *Parser) parseCreateFunction(createLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CreateFunctionStmt{Loc: createLoc}
	start := p.cur.Loc.Start

	if p.cur.Kind == kwGLOBAL {
		stmt.Global = true
		p.advance()
	}
	if p.cur.Kind == kwAGGREGATE {
		stmt.Aggregate = true
		p.advance()
	} else if p.cur.Kind == kwALIAS {
		stmt.Alias = true
		p.advance()
	}

	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}

	// Best-effort: consume the remainder of the statement up to the semicolon
	// or EOF, storing it as raw text. The function signature, RETURNS, body,
	// and PROPERTIES are all preserved verbatim.
	depth := 0
	for p.cur.Kind != tokEOF {
		if depth == 0 && p.cur.Kind == int(';') {
			break
		}
		if p.cur.Kind == int('(') || p.cur.Kind == int('[') || p.cur.Kind == int('{') {
			depth++
		} else if p.cur.Kind == int(')') || p.cur.Kind == int(']') || p.cur.Kind == int('}') {
			if depth > 0 {
				depth--
			}
		}
		p.advance()
	}
	end := p.prev.Loc.End
	if end > start {
		stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : end-p.baseOffset])
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropFunction parses Doris DROP [GLOBAL] FUNCTION name(arg_types).
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

	// Function name (qualified)
	if isIdentifierToken(p.cur.Kind) {
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// Argument-type list: (type, type, ...)
	if p.cur.Kind == int('(') {
		depth := 0
		for p.cur.Kind != tokEOF {
			if p.cur.Kind == int('(') {
				depth++
			} else if p.cur.Kind == int(')') {
				depth--
				if depth == 0 {
					p.advance() // consume ')'
					break
				}
			}
			p.advance()
		}
	}
	end := p.prev.Loc.End
	if end > start {
		stmt.RawSignature = strings.TrimSpace(p.input[start-p.baseOffset : end-p.baseOffset])
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
