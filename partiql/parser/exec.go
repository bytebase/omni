package parser

import "github.com/bytebase/omni/partiql/ast"

// parseExecCommand parses the execCommand grammar rule:
//
//	execCommand : EXEC name=expr (args+=expr (COMMA args+=expr)*)?
//
// The procedure name is an expression in the grammar (`name=expr`), and
// the reference parser's visitor refuses any but a variable, `EXEC f` or
// `EXEC "f"` (PartiQLParserDefault.kt:663 "EXEC procedure must be a symbol
// identifier"), so omni does too. Arguments are zero or more
// comma-separated expressions.
//
// Called from parseRootStatement when the current token is tokEXEC, which
// the caller has already consumed.
//
// Termination: stops at EOF, COLON_SEMI, or any token that
// parseExprTop cannot consume as part of an expression. No explicit
// delimiter is required after the last argument.
func (p *Parser) parseExecCommand() (*ast.ExecStmt, error) {
	start := p.prev.Loc.Start // EXEC was just consumed

	name, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	if _, ok := name.(*ast.VarRef); !ok {
		return nil, &ParseError{
			Message:  "EXEC procedure must be a symbol identifier",
			Position: name.GetLoc().Start, End: name.GetLoc().End,
		}
	}

	// Arguments are optional. Grammar: (args+=expr (COMMA args+=expr)*)?
	// The first argument follows the name directly (no preceding comma);
	// subsequent arguments are comma-separated. We detect the start of the
	// argument list by checking for any token that can begin an expression
	// and is not a statement boundary (EOF, COLON_SEMI).
	var args []ast.ExprNode
	if p.cur.Type != tokEOF && p.cur.Type != tokCOLON_SEMI {
		arg, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		for p.cur.Type == tokCOMMA {
			p.advance() // consume COMMA
			arg, err = p.parseExprTop()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}

	end := p.prev.Loc.End
	return &ast.ExecStmt{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: end},
	}, nil
}
