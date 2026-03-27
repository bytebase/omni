// Package parser - go_batch.go implements T-SQL GO batch separator parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseGoStmt parses a GO batch separator.
//
//	GO [count]
func (p *Parser) parseGoStmt() (*nodes.GoStmt, error) {
	loc := p.pos()
	p.advance() // consume GO

	stmt := &nodes.GoStmt{
		Count: 1,
		Loc:   nodes.Loc{Start: loc},
	}

	// Optional count
	if p.cur.Type == tokICONST {
		stmt.Count = int(p.cur.Ival)
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
