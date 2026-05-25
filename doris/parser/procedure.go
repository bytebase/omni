package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateProcedure parses a CREATE [OR REPLACE] PROCEDURE statement.
// On entry, CREATE has already been consumed (and OR REPLACE if present).
// cur is PROCEDURE.
//
// Syntax:
//
//	CREATE [OR REPLACE] PROCEDURE [IF NOT EXISTS] proc_name([param_list])
//	    [COMMENT 'text']
//	    BEGIN
//	        ... procedure body ...
//	    END
func (p *Parser) parseCreateProcedure(startLoc ast.Loc, orReplace bool) (ast.Node, error) {
	// Consume PROCEDURE keyword
	if _, err := p.expect(kwPROCEDURE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateProcedureStmt{
		OrReplace: orReplace,
		Loc:       startLoc,
	}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Procedure name (may be qualified: db.proc_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parameter list: ( [param, ...] )
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if p.cur.Kind != int(')') {
		params, err := p.parseProcedureParams()
		if err != nil {
			return nil, err
		}
		stmt.Parameters = params
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	// Optional COMMENT 'text'
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	// Body: BEGIN ... END (collect raw text, tracking nesting)
	body, bodyLoc, err := p.collectProcedureBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body
	stmt.Loc = startLoc.Merge(bodyLoc)

	return stmt, nil
}

// parseProcedureParams parses a non-empty comma-separated list of procedure
// parameters. Called after '(' has been consumed and ')' has not yet been seen.
func (p *Parser) parseProcedureParams() ([]*ast.ProcedureParam, error) {
	var params []*ast.ProcedureParam

	for {
		param, err := p.parseProcedureParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	return params, nil
}

// parseProcedureParam parses one procedure parameter:
//
//	[IN | OUT | INOUT] param_name type
func (p *Parser) parseProcedureParam() (*ast.ProcedureParam, error) {
	paramStart := p.cur.Loc
	param := &ast.ProcedureParam{}

	// Optional direction: IN, OUT, INOUT.
	// IN is a reserved keyword (kwIN). OUT and INOUT are not in the keyword
	// table and will be lexed as tokIdent. Detect them by their string value
	// only when the token AFTER them is a valid identifier (the param name),
	// using a lookahead check.
	switch {
	case p.cur.Kind == kwIN && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "IN"
		p.advance()
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "out") && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "OUT"
		p.advance()
	case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "inout") && isIdentifierToken(p.peekNext().Kind):
		param.Direction = "INOUT"
		p.advance()
	}

	// Parameter name
	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	param.Name = name

	// Type
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	param.Type = typeName
	param.Loc = paramStart.Merge(nameLoc).Merge(ast.NodeLoc(typeName))

	return param, nil
}

// collectProcedureBody collects the raw text of a BEGIN...END block,
// tracking nesting so inner BEGIN...END blocks are captured correctly.
// Returns the raw text and its source location.
// cur must be BEGIN on entry.
func (p *Parser) collectProcedureBody() (string, ast.Loc, error) {
	beginTok, err := p.expect(kwBEGIN)
	if err != nil {
		return "", ast.Loc{}, err
	}

	bodyStart := beginTok.Loc.Start
	depth := 1

	var endTok Token
	for depth > 0 {
		if p.cur.Kind == tokEOF {
			return "", ast.Loc{}, &ParseError{
				Loc: p.cur.Loc,
				Msg: "unexpected end of input inside procedure body (missing END)",
			}
		}
		tok := p.advance()
		switch tok.Kind {
		case kwBEGIN:
			depth++
		case kwEND:
			depth--
			if depth == 0 {
				endTok = tok
			}
		}
	}

	bodyEnd := endTok.Loc.End
	bodyText := strings.TrimSpace(p.input[bodyStart:bodyEnd])
	return bodyText, ast.Loc{Start: bodyStart, End: bodyEnd}, nil
}

// parseCallProcedure parses a CALL statement.
// On entry, cur is CALL (not yet consumed).
//
// Syntax:
//
//	CALL proc_name([args])
func (p *Parser) parseCallProcedure() (ast.Node, error) {
	callTok := p.advance() // consume CALL

	stmt := &ast.CallProcedureStmt{
		Loc: callTok.Loc,
	}

	// Procedure name (may be qualified: db.proc_name)
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Argument list: ( [expr, ...] )
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var args []ast.Node
	if p.cur.Kind != int(')') {
		args, err = p.parseExprList()
		if err != nil {
			return nil, err
		}
	}
	stmt.Args = args

	closeParen, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	stmt.Loc = callTok.Loc.Merge(closeParen.Loc)
	return stmt, nil
}

// parseDropProcedure parses a DROP PROCEDURE statement.
// On entry, DROP has already been consumed; cur is PROCEDURE.
//
// Syntax:
//
//	DROP PROCEDURE [IF EXISTS] proc_name
func (p *Parser) parseDropProcedure(startLoc ast.Loc) (ast.Node, error) {
	// Consume PROCEDURE keyword
	if _, err := p.expect(kwPROCEDURE); err != nil {
		return nil, err
	}

	stmt := &ast.DropProcedureStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Procedure name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}
