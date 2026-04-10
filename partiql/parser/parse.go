package parser

import "github.com/bytebase/omni/partiql/ast"

// Parse parses a PartiQL script and returns an *ast.List of statement nodes.
//
// The script rule is:
//
//	script : root (COLON_SEMI root)* COLON_SEMI? EOF
//	root   : (EXPLAIN (PAREN_LEFT explainOption (COMMA explainOption)* PAREN_RIGHT)?)? statement
//
// Multiple statements are separated by semicolons. A trailing semicolon is
// allowed. Each statement may optionally be prefixed with EXPLAIN.
//
// Returns an error on the first syntax error encountered. An empty input
// (or an input consisting only of whitespace/comments) is not an error —
// it returns an empty List.
func Parse(input string) (*ast.List, error) {
	p := NewParser(input)
	return p.parseScript()
}

// parseScript implements the script grammar rule. It is an internal method
// so that unit tests can drive it through NewParser without going through
// the public Parse function.
func (p *Parser) parseScript() (*ast.List, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}

	var stmts []ast.Node

	// A script may be completely empty.
	if p.cur.Type == tokEOF {
		return &ast.List{Items: stmts}, nil
	}

	// Parse the first root.
	stmt, err := p.parseRoot()
	if err != nil {
		return nil, err
	}
	stmts = append(stmts, stmt)

	// Parse subsequent roots separated by semicolons.
	for p.cur.Type == tokCOLON_SEMI {
		p.advance() // consume ;
		// Trailing semicolon: stop if we hit EOF after consuming it.
		if p.cur.Type == tokEOF {
			break
		}
		stmt, err = p.parseRoot()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: "unexpected token after statement",
			Loc:     p.cur.Loc,
		}
	}

	return &ast.List{Items: stmts}, nil
}

// parseRoot implements the root grammar rule:
//
//	root : (EXPLAIN (PAREN_LEFT explainOption (COMMA explainOption)* PAREN_RIGHT)?)? statement
//
// If the current token is EXPLAIN, the optional options list is parsed and
// discarded (ExplainStmt in the AST has no options field), and the inner
// statement is wrapped in an ExplainStmt.
func (p *Parser) parseRoot() (ast.StmtNode, error) {
	if p.cur.Type != tokEXPLAIN {
		return p.parseRootStatement()
	}

	explainStart := p.cur.Loc.Start
	p.advance() // consume EXPLAIN

	// Optional EXPLAIN options: (key value, key value, ...)
	// The AST has no options field on ExplainStmt, so we parse and discard.
	// Note: ExplainStmt only stores Inner + Loc; options are grammar-legal but
	// not representable in the current AST.
	if p.cur.Type == tokPAREN_LEFT {
		p.advance() // consume (
		for {
			// Each option is a pair of identifiers: param=IDENTIFIER value=IDENTIFIER.
			if _, err := p.expect(tokIDENT); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokIDENT); err != nil {
				return nil, err
			}
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
		if _, err := p.expect(tokPAREN_RIGHT); err != nil {
			return nil, err
		}
	}

	inner, err := p.parseRootStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ExplainStmt{
		Inner: inner,
		Loc:   ast.Loc{Start: explainStart, End: inner.GetLoc().End},
	}, nil
}

// parseRootStatement dispatches to the appropriate statement parser based
// on the current token. This is the internal statement parser used by
// parseRoot; it does not assert EOF (unlike ParseStatement).
func (p *Parser) parseRootStatement() (ast.StmtNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case tokCREATE:
		return p.parseCreateCommand()
	case tokDROP:
		return p.parseDropCommand()
	case tokINSERT:
		return p.parseInsertStmt()
	case tokUPDATE:
		return p.parseUpdateStmt()
	case tokDELETE:
		return p.parseDeleteStmt()
	case tokREPLACE:
		return p.parseReplaceStmt()
	case tokUPSERT:
		return p.parseUpsertStmt()
	case tokREMOVE:
		return p.parseRemoveStmt()
	case tokEXEC, tokEXECUTE:
		p.advance() // consume EXEC/EXECUTE
		return p.parseExecCommand()
	default:
		// DQL fallback: parse as expression. If the result is a StmtNode
		// (e.g. SelectStmt via SubLink unwrap), return it. Bare non-statement
		// expressions (e.g. `1 + 2`) are deferred — the grammar allows them
		// as dql but there is no ExprStmt wrapper in the AST.
		expr, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if sub, ok := expr.(*ast.SubLink); ok {
			return sub.Stmt, nil
		}
		if sn, ok := expr.(ast.StmtNode); ok {
			return sn, nil
		}
		return nil, p.deferredFeature("bare expression as statement", "parse-entry (DAG node 8)")
	}
}
