// Package parser - create_synonym.go implements CREATE SYNONYM parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSynonymStmt parses a CREATE SYNONYM statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-synonym-transact-sql
//
// BNF:
//
//	CREATE SYNONYM [ schema_name. ] synonym_name FOR <object>
func (p *Parser) parseCreateSynonymStmt() *nodes.CreateSynonymStmt {
	stmt := &nodes.CreateSynonymStmt{}
	stmt.Name = p.parseTableRef()
	p.match(kwFOR)
	stmt.Target = p.parseTableRef()
	stmt.Loc.End = p.pos()
	return stmt
}
