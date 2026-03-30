package parser

import "github.com/bytebase/omni/mongo/ast"

// parseShowCommand parses "show dbs", "show collections", etc.
func (p *Parser) parseShowCommand() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}

// parseUseCommand parses "use <database>".
func (p *Parser) parseUseCommand() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
