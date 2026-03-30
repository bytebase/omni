package parser

import "github.com/bytebase/omni/mongo/ast"

// parseEncryptionStatement parses db.getMongo().getKeyVault()/getClientEncryption() chains.
// Called after db.getMongo() and '.' have been consumed in parseDbStatement.
// The '.' is still current; next token is kwGetKeyVault or kwGetClientEncryption.
func (p *Parser) parseEncryptionStatement(stmtStart int) (ast.Node, error) {
	// Consume the '.'
	p.advance()

	// Determine target
	var target string
	if p.cur.Type == kwGetKeyVault {
		target = "keyVault"
	} else {
		target = "clientEncryption"
	}
	p.advance() // consume getKeyVault / getClientEncryption

	// Expect ()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Parse optional method chain
	var chain []ast.MethodCall
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		mc, err := p.parseMethodCall()
		if err != nil {
			return nil, err
		}
		chain = append(chain, *mc)
	}

	return &ast.EncryptionStatement{
		Target:         target,
		ChainedMethods: chain,
		Loc:            ast.Loc{Start: stmtStart, End: p.prev.End},
	}, nil
}
