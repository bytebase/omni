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
func (p *Parser) parsePrepareStmt() (nodes.Node, error) {
	p.advance() // consume PREPARE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// prep_type_clause: '(' type_list ')' | EMPTY
	var argtypes *nodes.List
	if p.cur.Type == '(' {
		p.advance()
		argtypes, err = p.parseTypeList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(AS); err != nil {
		return nil, err
	}

	// PreparableStmt: SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | MergeStmt
	query, err := p.parsePreparableStmt()
	if err != nil {
		return nil, err
	}

	return &nodes.PrepareStmt{
		Name:     name,
		Argtypes: argtypes,
		Query:    query,
	}, nil
}

// parseExecuteStmt parses an EXECUTE statement.
// The current token is EXECUTE (not yet consumed).
//
// Ref: gram.y ExecuteStmt
//
//	EXECUTE name execute_param_clause
func (p *Parser) parseExecuteStmt() (nodes.Node, error) {
	p.advance() // consume EXECUTE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	params, _ := p.parseExecuteParamClause()

	return &nodes.ExecuteStmt{
		Name:   name,
		Params: params,
	}, nil
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
func (p *Parser) parseDeallocateStmt() (nodes.Node, error) {
	p.advance() // consume DEALLOCATE

	// Optional PREPARE keyword
	if p.cur.Type == PREPARE {
		p.advance()
	}

	if p.cur.Type == ALL {
		p.advance()
		return &nodes.DeallocateStmt{
			IsAll: true,
		}, nil
	}

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.DeallocateStmt{
		Name: name,
	}, nil
}
