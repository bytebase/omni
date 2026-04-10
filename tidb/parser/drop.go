package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseDropTableStmt parses a DROP TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-table.html
//
//	DROP [TEMPORARY] TABLE [IF EXISTS] tbl_name [, tbl_name] ...
//	    [RESTRICT | CASCADE]
func (p *Parser) parseDropTableStmt(temporary bool) (*nodes.DropTableStmt, error) {
	start := p.pos()
	p.advance() // consume TABLE

	stmt := &nodes.DropTableStmt{
		Loc:       nodes.Loc{Start: start},
		Temporary: temporary,
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Completion: after DROP TABLE [IF EXISTS], offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Table list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// RESTRICT | CASCADE
	if _, ok := p.match(kwRESTRICT); ok {
		stmt.Restrict = true
	} else if _, ok := p.match(kwCASCADE); ok {
		stmt.Cascade = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropIndexStmt parses a DROP INDEX statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-index.html
//
//	DROP INDEX index_name ON tbl_name
//	    [algorithm_option | lock_option] ...
func (p *Parser) parseDropIndexStmt() (*nodes.DropIndexStmt, error) {
	start := p.pos()
	p.advance() // consume INDEX

	// Completion: after DROP INDEX, offer index_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("index_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.DropIndexStmt{Loc: nodes.Loc{Start: start}}

	// Index name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ON table
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	// Completion: after DROP INDEX idx ON, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// Optional ALGORITHM and LOCK
	for p.cur.Type == kwALGORITHM || p.cur.Type == kwLOCK {
		if _, ok := p.match(kwALGORITHM); ok {
			p.match('=')
			if p.isIdentToken() {
				stmt.Algorithm, _, _ = p.parseIdentifier()
			}
		} else if _, ok := p.match(kwLOCK); ok {
			p.match('=')
			if p.isIdentToken() {
				stmt.Lock, _, _ = p.parseIdentifier()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropViewStmt parses a DROP VIEW statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-view.html
//
//	DROP VIEW [IF EXISTS] view_name [, view_name] ... [RESTRICT | CASCADE]
func (p *Parser) parseDropViewStmt() (*nodes.DropViewStmt, error) {
	start := p.pos()
	p.advance() // consume VIEW

	stmt := &nodes.DropViewStmt{Loc: nodes.Loc{Start: start}}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// Completion: after DROP VIEW [IF EXISTS], offer view_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("view_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// View list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Views = append(stmt.Views, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// RESTRICT | CASCADE
	if _, ok := p.match(kwRESTRICT); ok {
		stmt.Restrict = true
	} else if _, ok := p.match(kwCASCADE); ok {
		stmt.Cascade = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/truncate-table.html
//
//	TRUNCATE [TABLE] tbl_name
func (p *Parser) parseTruncateStmt() (*nodes.TruncateStmt, error) {
	start := p.pos()
	p.advance() // consume TRUNCATE

	// Optional TABLE keyword
	p.match(kwTABLE)

	// Completion: after TRUNCATE [TABLE], offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.TruncateStmt{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Tables: []*nodes.TableRef{ref},
	}
	return stmt, nil
}

// parseRenameTableStmt parses a RENAME TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/rename-table.html
//
//	RENAME TABLE tbl_name TO new_tbl_name [, tbl_name2 TO new_tbl_name2] ...
func (p *Parser) parseRenameTableStmt() (*nodes.RenameTableStmt, error) {
	start := p.pos()
	// Already consumed RENAME
	p.match(kwTABLE)

	// Completion: after RENAME TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.RenameTableStmt{Loc: nodes.Loc{Start: start}}

	for {
		pairStart := p.pos()
		old, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		// Completion: after RENAME TABLE t TO, identifier context (new name).
		p.checkCursor()
		if p.collectMode() {
			// No specific candidates — user defines a new name.
			return nil, &ParseError{Message: "collecting"}
		}
		newRef, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Pairs = append(stmt.Pairs, &nodes.RenameTablePair{
			Loc: nodes.Loc{Start: pairStart, End: p.pos()},
			Old: old,
			New: newRef,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
