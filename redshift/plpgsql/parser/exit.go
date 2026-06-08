package parser

import (
	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// parseExit parses: EXIT [label] [WHEN condition] ;
func (p *Parser) parseExit() (ast.Node, error) {
	return p.parseExitContinue(false)
}

// parseContinue parses: CONTINUE [label] [WHEN condition] ;
func (p *Parser) parseContinue() (ast.Node, error) {
	return p.parseExitContinue(true)
}

// parseExitContinue is the shared implementation for EXIT and CONTINUE statements.
// Syntax: EXIT|CONTINUE [label] [WHEN condition] ;
func (p *Parser) parseExitContinue(isContinue bool) (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume EXIT or CONTINUE

	node := &ast.PLExit{
		IsContinue: isContinue,
	}

	// Check for optional label: if next token is an identifier that isn't WHEN, it's a label.
	if !p.isEOF() && p.cur.Type != ';' && !p.isKeyword("WHEN") {
		if p.isIdent() || p.isAnyKeywordAsIdent() {
			node.Label = p.identText()
			p.advance()
		}
	}

	// Check for optional WHEN condition
	if p.consumeKeyword("WHEN") {
		cond, err := p.collectUntilSemicolon()
		if err != nil {
			return nil, err
		}
		node.Condition = cond
	}

	// Consume trailing semicolon
	if _, err := p.expect(';'); err != nil {
		return nil, err
	}

	node.Loc = ast.Loc{Start: startPos, End: p.prev.End}
	return node, nil
}
