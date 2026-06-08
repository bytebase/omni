package parser

import (
	"strings"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// parseDynExecute parses an EXECUTE dynamic SQL statement.
//
// Ref: https://www.postgresql.org/docs/17/plpgsql-statements.html#PLPGSQL-STATEMENTS-EXECUTING-DYN
//
//	EXECUTE command-string [ INTO [STRICT] target [, ...] ] [ USING expression [, ...] ] ;
//
// INTO and USING can appear in either order.
// Duplicate INTO or USING produces an error.
func (p *Parser) parseDynExecute() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume EXECUTE

	// Collect the query expression until INTO, USING, or semicolon (tracking parens).
	query, err := p.collectUntil("INTO", "USING")
	if err != nil {
		return nil, err
	}
	if query == "" {
		return nil, p.errorf("syntax error at or near %q, expected expression after EXECUTE", p.tokenText(p.cur))
	}

	var into []string
	var strict bool
	var params []string
	hasInto := false
	hasUsing := false

	// Parse optional INTO and USING clauses (in any order).
	for {
		if p.isKeyword("INTO") {
			if hasInto {
				return nil, p.errorf("syntax error at or near %q, duplicate INTO clause", p.tokenText(p.cur))
			}
			hasInto = true
			p.advance() // consume INTO

			// Check for optional STRICT
			if p.isKeyword("STRICT") {
				strict = true
				p.advance()
			}

			// Parse comma-separated target identifiers
			targets, err := p.parseIntoTargets()
			if err != nil {
				return nil, err
			}
			into = targets
		} else if p.isKeyword("USING") {
			if hasUsing {
				return nil, p.errorf("syntax error at or near %q, duplicate USING clause", p.tokenText(p.cur))
			}
			hasUsing = true
			p.advance() // consume USING

			// Parse comma-separated USING parameter expressions
			usingParams, err := p.parseUsingParams()
			if err != nil {
				return nil, err
			}
			params = usingParams
		} else {
			break
		}
	}

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	return &ast.PLDynExecute{
		Query:  query,
		Into:   into,
		Strict: strict,
		Params: params,
		Loc:    ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// parseIntoTargets parses one or more comma-separated identifiers for an INTO clause.
func (p *Parser) parseIntoTargets() ([]string, error) {
	var targets []string
	for {
		if !p.isIdent() && !p.isAnyKeywordAsIdent() {
			return nil, p.errorf("syntax error at or near %q, expected identifier", p.tokenText(p.cur))
		}
		targets = append(targets, p.identText())
		p.advance()

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume comma
	}
	return targets, nil
}

// parseUsingParams parses comma-separated USING parameter expressions.
// Each expression is collected until a comma, INTO, or semicolon at paren depth 0.
func (p *Parser) parseUsingParams() ([]string, error) {
	var params []string
	for {
		expr, err := p.collectUsingExpr()
		if err != nil {
			return nil, err
		}
		expr = strings.TrimSpace(expr)
		if expr == "" {
			return nil, p.errorf("syntax error at or near %q, expected expression in USING clause", p.tokenText(p.cur))
		}
		params = append(params, expr)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume comma
	}
	return params, nil
}

// collectUsingExpr collects a single USING parameter expression until
// a comma, INTO, USING, or semicolon at paren depth 0.
func (p *Parser) collectUsingExpr() (string, error) {
	start := p.pos()
	depth := 0

	for !p.isEOF() {
		if p.cur.Type == '(' {
			depth++
			p.advance()
			continue
		}
		if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
			p.advance()
			continue
		}

		if depth == 0 {
			// Comma separates USING params
			if p.cur.Type == ',' {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				return text, nil
			}
			// INTO or USING or semicolon terminates USING params
			if p.isKeyword("INTO") || p.isKeyword("USING") || p.cur.Type == ';' {
				text := strings.TrimSpace(p.source[start:p.cur.Loc])
				return text, nil
			}
		}

		p.advance()
	}

	text := strings.TrimSpace(p.source[start:p.pos()])
	return text, nil
}
