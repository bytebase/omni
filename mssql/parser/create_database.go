// Package parser - create_database.go implements T-SQL CREATE DATABASE statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateDatabaseStmt parses a CREATE DATABASE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-transact-sql
//
//	CREATE DATABASE name [options]
func (p *Parser) parseCreateDatabaseStmt() *nodes.CreateDatabaseStmt {
	loc := p.pos()

	stmt := &nodes.CreateDatabaseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Database name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt
}
