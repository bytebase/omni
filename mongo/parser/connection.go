package parser

import "github.com/bytebase/omni/mongo/ast"

// parseConnectionStatement parses Mongo(), connect(), and db.getMongo() chains.
func (p *Parser) parseConnectionStatement() (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
