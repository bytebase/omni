package parser

import (
	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// parseIf parses an IF/ELSIF/ELSE/END IF statement.
// Grammar: IF condition THEN stmts [ELSIF condition THEN stmts]... [ELSE stmts] END IF ;
func (p *Parser) parseIf() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume IF

	// Collect condition until THEN
	cond, err := p.collectUntil("THEN")
	if err != nil {
		return nil, err
	}
	if !p.isKeyword("THEN") {
		return nil, p.errorf("syntax error at or near %q, expected THEN", p.tokenText(p.cur))
	}
	p.advance() // consume THEN

	// Parse THEN body
	thenBody, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	node := &ast.PLIf{
		Condition: cond,
		ThenBody:  thenBody,
		Loc:       ast.Loc{Start: start, End: -1},
	}

	// Parse ELSIF/ELSEIF branches
	for p.isKeyword("ELSIF") || p.isKeyword("ELSEIF") {
		elsif, err := p.parseElsIf()
		if err != nil {
			return nil, err
		}
		node.ElsIfs = append(node.ElsIfs, elsif)
	}

	// Parse optional ELSE
	if p.isKeyword("ELSE") {
		p.advance() // consume ELSE
		elseBody, err := p.parseStmtList()
		if err != nil {
			return nil, err
		}
		node.ElseBody = elseBody
	}

	// Expect END IF
	if !p.isKeyword("END") {
		return nil, p.errorf("syntax error at or near %q, expected END IF", p.tokenText(p.cur))
	}
	p.advance() // consume END
	if !p.isKeyword("IF") {
		return nil, p.errorf("syntax error at or near %q, expected IF after END", p.tokenText(p.cur))
	}
	p.advance() // consume IF

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	node.Loc.End = p.prev.End
	return node, nil
}

// parseElsIf parses a single ELSIF clause.
func (p *Parser) parseElsIf() (*ast.PLElsIf, error) {
	start := p.cur.Loc
	p.advance() // consume ELSIF/ELSEIF

	cond, err := p.collectUntil("THEN")
	if err != nil {
		return nil, err
	}
	if !p.isKeyword("THEN") {
		return nil, p.errorf("syntax error at or near %q, expected THEN after ELSIF condition", p.tokenText(p.cur))
	}
	p.advance() // consume THEN

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	return &ast.PLElsIf{
		Condition: cond,
		Body:      body,
		Loc:       ast.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseCase parses a CASE statement (both simple and searched forms).
// Grammar: CASE [testexpr] WHEN expr THEN stmts ... [ELSE stmts] END CASE ;
func (p *Parser) parseCase() (ast.Node, error) {
	start := p.cur.Loc
	p.advance() // consume CASE

	node := &ast.PLCase{
		Loc: ast.Loc{Start: start, End: -1},
	}

	// Determine if this is a simple or searched CASE
	if !p.isKeyword("WHEN") {
		// Simple CASE: collect test expression until WHEN
		testExpr, err := p.collectUntil("WHEN")
		if err != nil {
			return nil, err
		}
		node.TestExpr = testExpr
		node.HasTest = true
	}

	// Must have at least one WHEN
	if !p.isKeyword("WHEN") {
		return nil, p.errorf("syntax error at or near %q, expected WHEN in CASE statement", p.tokenText(p.cur))
	}

	// Parse WHEN branches
	for p.isKeyword("WHEN") {
		when, err := p.parseCaseWhen()
		if err != nil {
			return nil, err
		}
		node.Whens = append(node.Whens, when)
	}

	// Parse optional ELSE
	if p.isKeyword("ELSE") {
		p.advance() // consume ELSE
		elseBody, err := p.parseStmtList()
		if err != nil {
			return nil, err
		}
		node.ElseBody = elseBody
	}

	// Expect END CASE
	if !p.isKeyword("END") {
		return nil, p.errorf("syntax error at or near %q, expected END CASE", p.tokenText(p.cur))
	}
	p.advance() // consume END
	if !p.isKeyword("CASE") {
		return nil, p.errorf("syntax error at or near %q, expected CASE after END", p.tokenText(p.cur))
	}
	p.advance() // consume CASE

	// Consume trailing semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	node.Loc.End = p.prev.End
	return node, nil
}

// parseCaseWhen parses a single WHEN clause in a CASE statement.
func (p *Parser) parseCaseWhen() (*ast.PLCaseWhen, error) {
	start := p.cur.Loc
	p.advance() // consume WHEN

	expr, err := p.collectUntil("THEN")
	if err != nil {
		return nil, err
	}
	if !p.isKeyword("THEN") {
		return nil, p.errorf("syntax error at or near %q, expected THEN after WHEN expression", p.tokenText(p.cur))
	}
	p.advance() // consume THEN

	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}

	return &ast.PLCaseWhen{
		Expr: expr,
		Body: body,
		Loc:  ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
