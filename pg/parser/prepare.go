package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parsePrepareStmt parses a PREPARE statement.
// The current token is PREPARE (not yet consumed).
//
// Ref: gram.y PrepareStmt
//
//	PREPARE name prep_type_clause AS PreparableStmt
func (p *Parser) parsePrepareStmt() nodes.Node {
	p.advance() // consume PREPARE

	name, _ := p.parseName()

	// prep_type_clause: '(' type_list ')' | EMPTY
	var argtypes *nodes.List
	if p.cur.Type == '(' {
		p.advance()
		argtypes, _ = p.parseTypeList()
		p.expect(')')
	}

	p.expect(AS)

	// PreparableStmt: SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | MergeStmt
	query := p.parsePreparableStmt()

	return &nodes.PrepareStmt{
		Name:     name,
		Argtypes: argtypes,
		Query:    query,
	}
}

// parseExecuteStmt parses an EXECUTE statement.
// The current token is EXECUTE (not yet consumed).
//
// Ref: gram.y ExecuteStmt
//
//	EXECUTE name execute_param_clause
func (p *Parser) parseExecuteStmt() nodes.Node {
	p.advance() // consume EXECUTE

	name, _ := p.parseName()
	params := p.parseExecuteParamClause()

	return &nodes.ExecuteStmt{
		Name:   name,
		Params: params,
	}
}

// parseDeallocateStmt parses a DEALLOCATE statement.
// The current token is DEALLOCATE (not yet consumed).
//
// Ref: gram.y DeallocateStmt
//
//	DEALLOCATE name
//	| DEALLOCATE PREPARE name
//	| DEALLOCATE ALL
//	| DEALLOCATE PREPARE ALL
func (p *Parser) parseDeallocateStmt() nodes.Node {
	p.advance() // consume DEALLOCATE

	// Optional PREPARE keyword
	if p.cur.Type == PREPARE {
		p.advance()
	}

	if p.cur.Type == ALL {
		p.advance()
		return &nodes.DeallocateStmt{
			IsAll: true,
		}
	}

	name, _ := p.parseName()
	return &nodes.DeallocateStmt{
		Name: name,
	}
}
