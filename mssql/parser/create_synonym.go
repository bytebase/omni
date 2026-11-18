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
func (p *Parser) parseCreateSynonymStmt() (*nodes.CreateSynonymStmt, error) {
	stmt := &nodes.CreateSynonymStmt{}
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil {
		return nil, p.syntaxErrorAtCur()
	}
	// FOR and the target are required: CREATE SYNONYM FOR a is Msg 156,
	// CREATE SYNONYM a a is Msg 102, and CREATE SYNONYM a FOR is Msg 102
	// (measured).
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	stmt.Target, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Target == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
