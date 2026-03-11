// Package parser - drop.go implements T-SQL DROP statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseDropStmt parses a DROP statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-table-transact-sql
//
//	DROP TABLE|VIEW|INDEX|PROCEDURE|FUNCTION|DATABASE [IF EXISTS] name [, ...]
func (p *Parser) parseDropStmt() *nodes.DropStmt {
	loc := p.pos()
	p.advance() // consume DROP

	stmt := &nodes.DropStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Object type
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.DropTable
		p.advance()
	case kwVIEW:
		stmt.ObjectType = nodes.DropView
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.DropIndex
		p.advance()
	case kwPROCEDURE, kwPROC:
		stmt.ObjectType = nodes.DropProcedure
		p.advance()
	case kwFUNCTION:
		stmt.ObjectType = nodes.DropFunction
		p.advance()
	case kwDATABASE:
		stmt.ObjectType = nodes.DropDatabase
		p.advance()
	case kwSCHEMA:
		stmt.ObjectType = nodes.DropSchema
		p.advance()
	case kwTRIGGER:
		stmt.ObjectType = nodes.DropTrigger
		p.advance()
	case kwTYPE:
		stmt.ObjectType = nodes.DropType
		p.advance()
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		}
	}

	// Names (comma-separated)
	var names []nodes.Node
	for {
		name := p.parseTableRef()
		if name == nil {
			break
		}
		// For DROP INDEX: index_name ON table_name
		if stmt.ObjectType == nodes.DropIndex {
			if _, ok := p.match(kwON); ok {
				p.parseTableRef() // consume table name
			}
		}
		names = append(names, name)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Names = &nodes.List{Items: names}

	// Some DROP types also support CASCADE / RESTRICT
	if p.cur.Type == kwCASCADE {
		p.advance()
	} else if p.cur.Type == kwRESTRICT {
		p.advance()
	} else if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "cascade") {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
