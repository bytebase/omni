package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// parseDropTable parses:
//
//	DROP [TEMPORARY] TABLE [IF EXISTS] [db.]name [FORCE]
//
// On entry, DROP has been consumed and cur is TEMPORARY or TABLE.
func (p *Parser) parseDropTable(dropLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.DropTableStmt{Loc: dropLoc}

	if p.cur.Kind == kwTEMPORARY {
		stmt.Temporary = true
		p.advance()
	}

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Kind == kwFORCE {
		stmt.Force = true
		p.advance()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
