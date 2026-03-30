package parser

import "github.com/bytebase/omni/mongo/ast"

// parseEncryptionStatement parses db.getMongo().getKeyVault()/getClientEncryption() chains.
func (p *Parser) parseEncryptionStatement(stmtStart int) (ast.Node, error) {
	return nil, p.syntaxErrorAtCur()
}
