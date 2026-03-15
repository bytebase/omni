// Package parser - create_synonym.go implements CREATE SYNONYM parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSynonymStmt parses a CREATE SYNONYM statement.
//
// BNF: mssql/parser/bnf/create-synonym-transact-sql.bnf
//
//	CREATE SYNONYM [ schema_name_1. ] synonym_name FOR <object>
//
//	<object> ::=
//	{
//	    [ server_name. [ database_name ] . [ schema_name_2 ] .
//	      | database_name. [ schema_name_2 ] .
//	      | schema_name_2. ]
//	    object_name
//	}
func (p *Parser) parseCreateSynonymStmt() *nodes.CreateSynonymStmt {
	stmt := &nodes.CreateSynonymStmt{}
	stmt.Name = p.parseTableRef()
	p.match(kwFOR)
	stmt.Target = p.parseTableRef()
	stmt.Loc.End = p.pos()
	return stmt
}
